// Package redirects validates the `redirects` array in docs.json against
// the on-disk docs/ tree.
//
// Mintlify serves docs/foo/bar.mdx at URL /docs/foo/bar. The `redirects`
// array in docs.json rewrites URLs at the edge. This rule catches:
//   - Dead redirects (source already serves a real page).
//   - Missing destinations (destination has no real page).
//   - Bad prefixes / unsupported schemes / malformed paths.
//
// Diagnostic messages are stable strings asserted by tests/test-lint.sh.
package redirects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// wildcardSegment matches a single Mintlify-style wildcard URL segment,
// e.g. ":slug" or ":slug*".
var wildcardSegment = regexp.MustCompile(`^:[A-Za-z][A-Za-z0-9]*\*?$`)

// redirectsArrayKey locates the opening of the "redirects" array in
// docs.json source text.
var redirectsArrayKey = regexp.MustCompile(`"redirects"\s*:\s*\[`)

// sourceKeyPattern locates `"source":` literals to anchor each redirect
// entry to a 1-based line number.
var sourceKeyPattern = regexp.MustCompile(`"source"\s*:`)

// Check reads ${contentRoot}/docs.json and returns violations for
// dead/missing/malformed redirects. If docs.json is absent, returns nil.
// If docs.json is present but unparseable, returns a single violation
// with a parse-error message.
func Check(contentRoot string) []analysis.Violation {
	docsJSONPath := filepath.Join(contentRoot, "docs.json")
	data, err := os.ReadFile(docsJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    1,
			Col:     1,
			Message: fmt.Sprintf("read error: %s", err),
		}}
	}
	return validate(data, contentRoot)
}

// validate parses docs.json bytes and returns redirect violations. Pure
// (no I/O of its own beyond filesystem checks under contentRoot) so that
// tests can drive it directly with byte literals.
func validate(docsJSONBytes []byte, contentRoot string) []analysis.Violation {
	var parsed map[string]any
	if err := json.Unmarshal(docsJSONBytes, &parsed); err != nil {
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    1,
			Col:     1,
			Message: fmt.Sprintf("invalid JSON: %s", err),
		}}
	}

	rawRedirects, ok := parsed["redirects"]
	if !ok {
		return nil
	}
	redirects, ok := rawRedirects.([]any)
	if !ok {
		return nil
	}
	if len(redirects) == 0 {
		return nil
	}

	docsJSONText := string(docsJSONBytes)
	lines := lineLookup(docsJSONText, len(redirects))
	openAPISources := collectOpenAPISources(parsed)

	var violations []analysis.Violation
	for i, raw := range redirects {
		line := lines[i]
		entry, ok := raw.(map[string]any)
		if !ok {
			violations = append(violations, analysis.Violation{
				File:    "docs.json",
				Line:    line,
				Col:     1,
				Message: formatMessage(i, "entry", "", "not an object"),
			})
			continue
		}

		source, sourceOk := stringField(entry, "source")
		destination, destOk := stringField(entry, "destination")

		if !sourceOk {
			violations = append(violations, analysis.Violation{
				File:    "docs.json",
				Line:    line,
				Col:     1,
				Message: formatMessage(i, "source", source, "must be a non-empty string"),
			})
		}
		if !destOk {
			violations = append(violations, analysis.Violation{
				File:    "docs.json",
				Line:    line,
				Col:     1,
				Message: formatMessage(i, "destination", destination, "must be a non-empty string"),
			})
		}
		if !sourceOk || !destOk {
			continue
		}

		// Rule (b): source must start with '/'; destination must start
		// with '/' or http(s)://.
		if !strings.HasPrefix(source, "/") {
			violations = append(violations, analysis.Violation{
				File:    "docs.json",
				Line:    line,
				Col:     1,
				Message: formatMessage(i, "source", source, "bad path: must start with '/'"),
			})
		}

		destIsHTTP := strings.HasPrefix(destination, "http://") ||
			strings.HasPrefix(destination, "https://")
		if !strings.HasPrefix(destination, "/") && !destIsHTTP {
			violations = append(violations, analysis.Violation{
				File:    "docs.json",
				Line:    line,
				Col:     1,
				Message: formatMessage(i, "destination", destination, "bad path: must start with '/' (or http(s)://)"),
			})
		}

		if strings.HasPrefix(source, "/") {
			violations = append(violations, validateSource(i, source, contentRoot, line)...)
		}
		if !destIsHTTP && strings.HasPrefix(destination, "/") {
			violations = append(violations, validateDestination(i, destination, contentRoot, openAPISources, line)...)
		}
	}

	return violations
}

// stringField returns the string value at key (and true) when it is a
// non-empty string. When the key is missing, the value is not a string,
// or the string is empty, it returns ("", false). When the value is a
// non-string scalar like a number, it falls back to ("", false) — the
// Node.js reference reports an empty value in that case too.
func stringField(entry map[string]any, key string) (string, bool) {
	raw, present := entry[key]
	if !present || raw == nil {
		return "", false
	}
	s, ok := raw.(string)
	if !ok {
		return "", false
	}
	if s == "" {
		return "", false
	}
	return s, true
}

// validateSource emits filesystem-related violations for a source URL
// that already starts with '/'. The source must additionally start with
// /docs/ (or be exactly /docs); after stripping the prefix it is
// resolved against contentRoot to detect dead redirects.
func validateSource(i int, source, contentRoot string, line int) []analysis.Violation {
	cleaned := cleanPath(source)
	if !strings.HasPrefix(cleaned, "docs/") && cleaned != "docs" {
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    line,
			Col:     1,
			Message: formatMessage(i, "source", source, "bad prefix (must start with /docs/)"),
		}}
	}

	prefix, hasWildcard := splitWildcard(cleaned)
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))

	if !hasWildcard {
		if pageExists(prefixFs) {
			return []analysis.Violation{{
				File:    "docs.json",
				Line:    line,
				Col:     1,
				Message: formatMessage(i, "source", source, "dead redirect (source resolves to a real page)"),
			}}
		}
		return nil
	}

	// Wildcard source: prefix must not be a real page or a directory
	// containing pages.
	if pageExists(prefixFs) {
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    line,
			Col:     1,
			Message: formatMessage(i, "source", source, "dead redirect (wildcard source prefix resolves to a real page)"),
		}}
	}
	if dirExists(prefixFs) && dirHasPages(prefixFs) {
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    line,
			Col:     1,
			Message: formatMessage(i, "source", source, "dead redirect (wildcard source prefix has real pages)"),
		}}
	}
	return nil
}

// validateDestination emits filesystem-related violations for a
// destination URL that already starts with '/'. Honors the OpenAPI
// escape hatch: a wildcard destination whose on-disk prefix is not a
// directory is still accepted when ${prefix}.yaml is registered as a
// nav.*.openapi.source value somewhere in docs.json.
func validateDestination(i int, destination, contentRoot string, openAPISources map[string]bool, line int) []analysis.Violation {
	cleaned := cleanPath(destination)
	if !strings.HasPrefix(cleaned, "docs/") && cleaned != "docs" {
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    line,
			Col:     1,
			Message: formatMessage(i, "destination", destination, "bad prefix (must start with /docs/)"),
		}}
	}

	prefix, hasWildcard := splitWildcard(cleaned)
	prefixRel := strings.Join(prefix, "/")
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))

	if !hasWildcard {
		if pageExists(prefixFs) {
			return nil
		}
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    line,
			Col:     1,
			Message: formatMessage(i, "destination", destination, "missing destination (no .mdx or .md page exists)"),
		}}
	}

	if dirExists(prefixFs) {
		return nil
	}
	yamlRel := prefixRel + ".yaml"
	yamlFs := prefixFs + ".yaml"
	if openAPISources[yamlRel] && fileExists(yamlFs) {
		return nil
	}
	return []analysis.Violation{{
		File:    "docs.json",
		Line:    line,
		Col:     1,
		Message: formatMessage(i, "destination", destination, "wildcard prefix not a directory"),
	}}
}

// cleanPath strips, in order, any "?..." query string, any "#..."
// fragment, the leading '/', and a trailing '/'. The result is suitable
// for splitting on '/' and joining with contentRoot.
func cleanPath(p string) string {
	s := p
	if idx := strings.Index(s, "?"); idx != -1 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "#"); idx != -1 {
		s = s[:idx]
	}
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")
	return s
}

// splitWildcard returns the prefix segments preceding the first
// wildcard segment (per wildcardSegment) and a flag indicating whether
// any wildcard segment was present.
func splitWildcard(cleaned string) ([]string, bool) {
	rawSegments := strings.Split(cleaned, "/")
	var segments []string
	for _, seg := range rawSegments {
		if seg != "" {
			segments = append(segments, seg)
		}
	}
	var prefix []string
	hasWildcard := false
	for _, seg := range segments {
		if wildcardSegment.MatchString(seg) {
			hasWildcard = true
			break
		}
		prefix = append(prefix, seg)
	}
	return prefix, hasWildcard
}

// pageExists returns true when `${prefixFs}.mdx` or `${prefixFs}.md`
// exists as a regular file.
func pageExists(prefixFs string) bool {
	return fileExists(prefixFs+".mdx") || fileExists(prefixFs+".md")
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	if err != nil {
		return false
	}
	return st.Mode().IsRegular()
}

func dirExists(p string) bool {
	st, err := os.Stat(p)
	if err != nil {
		return false
	}
	return st.IsDir()
}

// dirHasPages returns true when dir contains any .mdx or .md file at
// any depth.
func dirHasPages(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, ent := range entries {
		full := filepath.Join(dir, ent.Name())
		if ent.IsDir() {
			if dirHasPages(full) {
				return true
			}
			continue
		}
		name := ent.Name()
		if strings.HasSuffix(name, ".mdx") || strings.HasSuffix(name, ".md") {
			return true
		}
	}
	return false
}

// collectOpenAPISources walks any object node looking for
// `openapi.source` string values and returns them as a set. Mintlify
// generates pages from these yaml files at build time, so a wildcard
// destination may target a virtual directory that only exists at build
// time; the values here form the OpenAPI escape hatch.
func collectOpenAPISources(parsed map[string]any) map[string]bool {
	result := map[string]bool{}
	walkOpenAPISources(parsed, result)
	return result
}

func walkOpenAPISources(node any, result map[string]bool) {
	switch n := node.(type) {
	case map[string]any:
		if openapi, ok := n["openapi"].(map[string]any); ok {
			if src, ok := openapi["source"].(string); ok {
				result[src] = true
			}
		}
		for _, v := range n {
			walkOpenAPISources(v, result)
		}
	case []any:
		for _, item := range n {
			walkOpenAPISources(item, result)
		}
	}
}

// lineLookup returns a slice of length count whose i-th element is the
// 1-based line number of the i-th redirect entry's `"source":` literal
// in docsJSONText. Entries that cannot be anchored (e.g. non-object
// entries lacking "source":) fall back to 1.
func lineLookup(docsJSONText string, count int) []int {
	result := make([]int, count)
	for i := range result {
		result[i] = 1
	}
	arrayMatch := redirectsArrayKey.FindStringIndex(docsJSONText)
	if arrayMatch == nil {
		return result
	}
	// Start scanning from the offset of the opening '['.
	start := arrayMatch[1] - 1
	matches := sourceKeyPattern.FindAllStringIndex(docsJSONText[start:], count)
	for i, m := range matches {
		if i >= count {
			break
		}
		offset := start + m[0]
		result[i] = lineForOffset(docsJSONText, offset)
	}
	return result
}

// lineForOffset returns the 1-based line number of byteOffset in text.
func lineForOffset(text string, byteOffset int) int {
	line := 1
	for j := 0; j < byteOffset && j < len(text); j++ {
		if text[j] == '\n' {
			line++
		}
	}
	return line
}

// formatMessage returns the canonical diagnostic message body for a
// single redirect violation. The format is asserted (substring match)
// by tests/test-lint.sh.
func formatMessage(i int, field, value, message string) string {
	return fmt.Sprintf("redirects[%d] %s %q: %s", i, field, value, message)
}

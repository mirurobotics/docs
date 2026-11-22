package redirects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func wildcardSegment() *regexp.Regexp {
	return regexp.MustCompile(`^:[A-Za-z][A-Za-z0-9]*\*?$`)
}

// Check validates the redirects in contentRoot/docs.json; a missing file is skipped.
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
		violations = append(
			violations,
			validateEntry(i, raw, lines[i], contentRoot, openAPISources)...,
		)
	}

	return violations
}

type entryFields struct {
	i           int
	line        int
	source      string
	sourceOk    bool
	destination string
	destOk      bool
}

func validateEntry(
	i int,
	raw any,
	line int,
	contentRoot string,
	openAPISources map[string]bool,
) []analysis.Violation {
	entry, ok := raw.(map[string]any)
	if !ok {
		return []analysis.Violation{{
			File:    "docs.json",
			Line:    line,
			Col:     1,
			Message: formatMessage(i, "entry", "", "not an object"),
		}}
	}

	src, srcOk := stringField(entry, "source")
	dst, dstOk := stringField(entry, "destination")
	fields := entryFields{
		i:           i,
		line:        line,
		source:      src,
		sourceOk:    srcOk,
		destination: dst,
		destOk:      dstOk,
	}

	violations := validateRequiredFields(fields)
	if !srcOk || !dstOk {
		return violations
	}

	violations = append(violations, validatePathFormat(fields)...)
	fsViolations := validateFileSystem(fields, contentRoot, openAPISources)
	return append(violations, fsViolations...)
}

func validateRequiredFields(fields entryFields) []analysis.Violation {
	const msg = "must be a non-empty string"
	var violations []analysis.Violation
	if !fields.sourceOk {
		violations = append(violations, analysis.Violation{
			File:    "docs.json",
			Line:    fields.line,
			Col:     1,
			Message: formatMessage(fields.i, "source", fields.source, msg),
		})
	}
	if !fields.destOk {
		violations = append(violations, analysis.Violation{
			File:    "docs.json",
			Line:    fields.line,
			Col:     1,
			Message: formatMessage(fields.i, "destination", fields.destination, msg),
		})
	}
	return violations
}

func validatePathFormat(fields entryFields) []analysis.Violation {
	var violations []analysis.Violation
	if !strings.HasPrefix(fields.source, "/") {
		const msg = "bad path: must start with '/'"
		violations = append(violations, analysis.Violation{
			File:    "docs.json",
			Line:    fields.line,
			Col:     1,
			Message: formatMessage(fields.i, "source", fields.source, msg),
		})
	}
	if !strings.HasPrefix(fields.destination, "/") && !destIsHTTP(fields.destination) {
		const msg = "bad path: must start with '/' (or http(s)://)"
		violations = append(violations, analysis.Violation{
			File:    "docs.json",
			Line:    fields.line,
			Col:     1,
			Message: formatMessage(fields.i, "destination", fields.destination, msg),
		})
	}
	return violations
}

func validateFileSystem(
	fields entryFields,
	contentRoot string,
	openAPISources map[string]bool,
) []analysis.Violation {
	var violations []analysis.Violation
	if strings.HasPrefix(fields.source, "/") {
		srcPrefix, _ := splitWildcard(cleanPath(fields.source))
		if containsTraversalSegment(srcPrefix) {
			violations = append(violations, analysis.Violation{
				File: "docs.json",
				Line: fields.line,
				Col:  1,
				Message: formatMessage(
					fields.i, "source", fields.source,
					"bad path: contains '..' or '.' segment",
				),
			})
		} else {
			srcViolations := validateSource(
				fields.i, fields.source, contentRoot, openAPISources, fields.line,
			)
			violations = append(violations, srcViolations...)
		}
	}
	if !destIsHTTP(fields.destination) && strings.HasPrefix(fields.destination, "/") {
		dstPrefix, _ := splitWildcard(cleanPath(fields.destination))
		if containsTraversalSegment(dstPrefix) {
			violations = append(violations, analysis.Violation{
				File: "docs.json",
				Line: fields.line,
				Col:  1,
				Message: formatMessage(
					fields.i, "destination", fields.destination,
					"bad path: contains '..' or '.' segment",
				),
			})
		} else {
			dstViolations := validateDestination(
				fields.i, fields.destination, contentRoot, openAPISources, fields.line,
			)
			violations = append(violations, dstViolations...)
		}
	}
	return violations
}

func containsTraversalSegment(segments []string) bool {
	for _, seg := range segments {
		if seg == ".." || seg == "." {
			return true
		}
	}
	return false
}

func destIsHTTP(destination string) bool {
	return strings.HasPrefix(destination, "http://") ||
		strings.HasPrefix(destination, "https://")
}

func stringField(entry map[string]any, key string) (string, bool) {
	raw, present := entry[key]
	if !present || raw == nil {
		return "", false
	}
	value, ok := raw.(string)
	if !ok {
		return "", false
	}
	if value == "" {
		return "", false
	}
	return value, true
}

func validateSource(
	i int,
	source, contentRoot string,
	openAPISources map[string]bool,
	line int,
) []analysis.Violation {
	cleaned := cleanPath(source)
	prefix, hasWildcard := splitWildcard(cleaned)
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))

	if !hasWildcard {
		if pageExists(prefixFs) {
			return []analysis.Violation{{
				File: "docs.json",
				Line: line,
				Col:  1,
				Message: formatMessage(
					i, "source", source,
					"dead redirect (source resolves to a real page)",
				),
			}}
		}
		return nil
	}

	return validateWildcardSource(wildcardSourceCtx{
		i:              i,
		source:         source,
		prefix:         prefix,
		prefixFs:       prefixFs,
		openAPISources: openAPISources,
		line:           line,
	})
}

type wildcardSourceCtx struct {
	i              int
	source         string
	prefix         []string
	prefixFs       string
	openAPISources map[string]bool
	line           int
}

func validateWildcardSource(in wildcardSourceCtx) []analysis.Violation {
	i, source, prefix, prefixFs := in.i, in.source, in.prefix, in.prefixFs
	openAPISources, line := in.openAPISources, in.line
	if pageExists(prefixFs) {
		return []analysis.Violation{{
			File: "docs.json",
			Line: line,
			Col:  1,
			Message: formatMessage(
				i, "source", source,
				"dead redirect (wildcard source prefix resolves to a real page)",
			),
		}}
	}
	if dirExists(prefixFs) && dirHasPages(prefixFs) {
		return []analysis.Violation{{
			File: "docs.json",
			Line: line,
			Col:  1,
			Message: formatMessage(
				i, "source", source,
				"dead redirect (wildcard source prefix has real pages)",
			),
		}}
	}
	yamlRel := strings.Join(prefix, "/") + ".yaml"
	yamlFs := prefixFs + ".yaml"
	if openAPISources[yamlRel] && fileExists(yamlFs) {
		return []analysis.Violation{{
			File: "docs.json",
			Line: line,
			Col:  1,
			Message: formatMessage(
				i, "source", source,
				"dead redirect (wildcard source prefix has Mintlify-generated pages)",
			),
		}}
	}
	return nil
}

func validateDestination(
	i int,
	destination, contentRoot string,
	openAPISources map[string]bool,
	line int,
) []analysis.Violation {
	cleaned := cleanPath(destination)
	prefix, hasWildcard := splitWildcard(cleaned)
	prefixRel := strings.Join(prefix, "/")
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))

	if !hasWildcard {
		if pageExists(prefixFs) {
			return nil
		}
		return []analysis.Violation{{
			File: "docs.json",
			Line: line,
			Col:  1,
			Message: formatMessage(
				i, "destination", destination,
				"missing destination (no .mdx or .md page exists)",
			),
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
		File: "docs.json",
		Line: line,
		Col:  1,
		Message: formatMessage(
			i, "destination", destination,
			"wildcard prefix not a directory",
		),
	}}
}

func cleanPath(raw string) string {
	cleaned := raw
	if idx := strings.Index(cleaned, "?"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	if idx := strings.Index(cleaned, "#"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	cleaned = strings.TrimPrefix(cleaned, "/")
	cleaned = strings.TrimSuffix(cleaned, "/")
	return cleaned
}

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
		if wildcardSegment().MatchString(seg) {
			hasWildcard = true
			break
		}
		prefix = append(prefix, seg)
	}
	return prefix, hasWildcard
}

func pageExists(prefixFs string) bool {
	return fileExists(prefixFs+".mdx") || fileExists(prefixFs+".md")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

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

func collectOpenAPISources(parsed map[string]any) map[string]bool {
	result := map[string]bool{}
	walkOpenAPISources(parsed, result)
	return result
}

func walkOpenAPISources(node any, result map[string]bool) {
	switch node := node.(type) {
	case map[string]any:
		if openapi, ok := node["openapi"].(map[string]any); ok {
			if src, ok := openapi["source"].(string); ok {
				result[src] = true
			}
		}
		for _, child := range node {
			walkOpenAPISources(child, result)
		}
	case []any:
		for _, item := range node {
			walkOpenAPISources(item, result)
		}
	}
}

func lineLookup(docsJSONText string, count int) []int {
	result := make([]int, count)
	for i := range result {
		result[i] = 1
	}
	if count == 0 {
		return result
	}

	dec := json.NewDecoder(bytes.NewReader([]byte(docsJSONText)))
	if !advanceToRedirectsArray(dec) {
		return result
	}
	scanRedirectEntries(dec, docsJSONText, result, count)
	return result
}

func advanceToRedirectsArray(dec *json.Decoder) bool {
	for {
		tok, err := dec.Token()
		if err != nil {
			return false
		}
		key, ok := tok.(string)
		if !ok || key != "redirects" {
			continue
		}
		next, err := dec.Token()
		if err != nil {
			return false
		}
		if delim, ok := next.(json.Delim); ok && delim == '[' {
			return true
		}
	}
}

func scanRedirectEntries(
	dec *json.Decoder,
	docsJSONText string,
	result []int,
	count int,
) {
	arrayDepth := 1
	entryIdx := 0
	for {
		offsetBefore := dec.InputOffset()
		tok, err := dec.Token()
		if err != nil {
			return
		}

		if delim, ok := tok.(json.Delim); ok {
			switch delim {
			case '{', '[':
				if arrayDepth == 1 && entryIdx < count {
					result[entryIdx] = lineForOffset(
						docsJSONText,
						tokenStart(docsJSONText, int(offsetBefore)),
					)
					entryIdx++
				}
				arrayDepth++
			case '}', ']':
				arrayDepth--
				if arrayDepth == 0 {
					return
				}
			}
			continue
		}

		if arrayDepth == 1 && entryIdx < count {
			result[entryIdx] = lineForOffset(
				docsJSONText,
				tokenStart(docsJSONText, int(offsetBefore)),
			)
			entryIdx++
		}
	}
}

func tokenStart(text string, offset int) int {
	for offset < len(text) {
		char := text[offset]
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' || char == ',' {
			offset++
			continue
		}
		break
	}
	return offset
}

func lineForOffset(text string, byteOffset int) int {
	line := 1
	for i := 0; i < byteOffset; i++ {
		if text[i] == '\n' {
			line++
		}
	}
	return line
}

func formatMessage(i int, field, value, message string) string {
	return fmt.Sprintf("redirects[%d] %s %q: %s", i, field, value, message)
}

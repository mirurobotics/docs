package redirects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

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

func formatMessage(i int, field, value, message string) string {
	return fmt.Sprintf("redirects[%d] %s %q: %s", i, field, value, message)
}

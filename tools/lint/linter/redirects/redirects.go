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

func collectOpenAPISources(parsed map[string]any) map[string]bool {
	result := map[string]bool{}
	walkOpenAPISources(parsed, result)
	return result
}

func walkOpenAPISources(node any, result map[string]bool) {
	switch node := node.(type) {
	case map[string]any:
		openapi, _ := node["openapi"].(map[string]any)
		if src, ok := openapi["source"].(string); ok {
			result[src] = true
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
		return entryViolation(i, line, "entry", "", "not an object")
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

func validateRequiredFields(fields entryFields) []analysis.Violation {
	const msg = "must be a non-empty string"
	i, line := fields.i, fields.line
	var violations []analysis.Violation
	if !fields.sourceOk {
		violations = entryViolation(i, line, "source", fields.source, msg)
	}
	if !fields.destOk {
		dst := entryViolation(i, line, "destination", fields.destination, msg)
		violations = append(violations, dst...)
	}
	return violations
}

func validatePathFormat(fields entryFields) []analysis.Violation {
	i, line := fields.i, fields.line
	var violations []analysis.Violation
	if !strings.HasPrefix(fields.source, "/") {
		const msg = "bad path: must start with '/'"
		violations = entryViolation(i, line, "source", fields.source, msg)
	}
	if !strings.HasPrefix(fields.destination, "/") && !destIsHTTP(fields.destination) {
		const msg = "bad path: must start with '/' (or http(s)://)"
		dst := entryViolation(i, line, "destination", fields.destination, msg)
		violations = append(violations, dst...)
	}
	return violations
}

func entryViolation(i, line int, field, value, msg string) []analysis.Violation {
	return []analysis.Violation{{
		File:    "docs.json",
		Line:    line,
		Col:     1,
		Message: formatMessage(i, field, value, msg),
	}}
}

func formatMessage(i int, field, value, message string) string {
	return fmt.Sprintf("redirects[%d] %s %q: %s", i, field, value, message)
}

func destIsHTTP(dest string) bool {
	return strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://")
}

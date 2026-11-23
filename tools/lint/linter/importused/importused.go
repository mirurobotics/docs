package importused

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Check verifies that every imported name is used in the document body.
func Check(file string, lines []string) []analysis.Violation {
	imports := analysis.ParseImports(lines)
	if len(imports) == 0 {
		return nil
	}
	body := strings.Join(analysis.BodyLines(lines), "\n")
	var violations []analysis.Violation
	for _, imp := range imports {
		violations = append(violations, checkUsed(file, imp, body)...)
	}
	return violations
}

func checkUsed(
	file string,
	imp analysis.ParsedImport,
	body string,
) []analysis.Violation {
	var violations []analysis.Violation
	for _, binding := range imp.Names {
		pattern := regexp.MustCompile(
			fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(binding.Name)),
		)
		if pattern.MatchString(body) {
			continue
		}
		msg := fmt.Sprintf("import-used: %q is imported but never used", binding.Name)
		violations = append(violations, analysis.Violation{
			File: file, Line: imp.Line, Col: 1, Message: msg,
		})
	}
	return violations
}

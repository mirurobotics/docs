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
		for _, n := range imp.Names {
			re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(n.Name)))
			if re.MatchString(body) {
				continue
			}
			msg := fmt.Sprintf("import-used: %q is imported but never used", n.Name)
			violations = append(violations, analysis.Violation{
				File: file, Line: imp.Line, Col: 1, Message: msg,
			})
		}
	}
	return violations
}

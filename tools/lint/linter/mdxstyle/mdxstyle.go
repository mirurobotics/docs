package mdxstyle

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Check validates style of imports ending in .mdx.
func Check(file string, lines []string) []analysis.Violation {
	imports := analysis.ParseImports(lines)
	var violations []analysis.Violation
	for _, imp := range imports {
		if !strings.HasSuffix(imp.Path, ".mdx") {
			continue
		}
		if strings.Contains(imp.Raw, "{") {
			msg := "import-mdx-style: .mdx imports must use default import syntax"
			violations = append(violations, analysis.Violation{
				File: file, Line: imp.Line, Col: 1, Message: msg,
			})
		}
		if !strings.HasSuffix(strings.TrimRight(imp.Raw, " \t"), ";") {
			violations = append(violations, analysis.Violation{
				File:    file,
				Line:    imp.Line,
				Col:     1,
				Message: "import-mdx-style: import statement must end with ';'",
			})
		}
	}
	return violations
}

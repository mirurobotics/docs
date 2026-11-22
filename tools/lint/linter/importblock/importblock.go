package importblock

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Check verifies that no blank lines appear inside a contiguous import block.
func Check(file string, lines []string) []analysis.Violation {
	imports := analysis.ParseImports(lines)
	if len(imports) < 2 {
		return nil
	}
	var violations []analysis.Violation
	for i := 1; i < len(imports); i++ {
		prev := imports[i-1].Line
		curr := imports[i].Line
		onlyBlanks := true
		for lineIdx := prev; lineIdx < curr-1; lineIdx++ {
			line := lines[lineIdx]
			if strings.TrimSpace(line) != "" && !analysis.IsImportLine(line) {
				onlyBlanks = false
				break
			}
		}
		if !onlyBlanks {
			continue
		}
		for lineIdx := prev; lineIdx < curr-1; lineIdx++ {
			if strings.TrimSpace(lines[lineIdx]) == "" {
				violations = append(violations, analysis.Violation{
					File:    file,
					Line:    lineIdx + 1,
					Col:     1,
					Message: "import-block-contiguous: blank line inside import block",
				})
			}
		}
	}
	return violations
}

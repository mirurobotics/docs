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
		gap := lines[prev : curr-1]
		if !hasOnlyBlankOrImportLines(gap) {
			continue
		}
		violations = append(violations, checkBlankLines(file, gap, prev)...)
	}
	return violations
}

func hasOnlyBlankOrImportLines(gap []string) bool {
	for _, line := range gap {
		if strings.TrimSpace(line) != "" && !analysis.IsImportLine(line) {
			return false
		}
	}
	return true
}

func checkBlankLines(file string, gap []string, start int) []analysis.Violation {
	var violations []analysis.Violation
	for offset, line := range gap {
		if strings.TrimSpace(line) == "" {
			violations = append(violations, analysis.Violation{
				File:    file,
				Line:    start + offset + 1,
				Col:     1,
				Message: "import-block-contiguous: blank line inside import block",
			})
		}
	}
	return violations
}

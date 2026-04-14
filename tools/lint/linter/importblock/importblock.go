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
		prev := imports[i-1].Line // 1-based
		curr := imports[i].Line   // 1-based
		// Check whether the lines strictly between prev and curr are all blank.
		// If any non-blank, non-import line exists, these imports are in
		// separate blocks and no violation is reported.
		onlyBlanks := true
		for j := prev; j < curr-1; j++ { // j is 0-based index = line (j+1)
			line := lines[j]
			if strings.TrimSpace(line) != "" && !analysis.IsImportLine(line) {
				onlyBlanks = false
				break
			}
		}
		if !onlyBlanks {
			continue
		}
		// All inter-import lines are blank — flag them.
		for j := prev; j < curr-1; j++ {
			if strings.TrimSpace(lines[j]) == "" {
				violations = append(violations, analysis.Violation{
					File:    file,
					Line:    j + 1,
					Col:     1,
					Message: "import-block-contiguous: blank line inside import block",
				})
			}
		}
	}
	return violations
}

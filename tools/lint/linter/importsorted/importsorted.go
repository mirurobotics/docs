package importsorted

import (
	"fmt"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Check verifies imports are sorted by path (case-insensitive, ascending).
func Check(file string, lines []string) []analysis.Violation {
	imports := analysis.ParseImports(lines)
	for i := 1; i < len(imports); i++ {
		if strings.ToLower(imports[i].Path) < strings.ToLower(imports[i-1].Path) {
			msg := fmt.Sprintf(
				"import-sorted: import path %q is out of order (expected after %q)",
				imports[i].Path,
				imports[i-1].Path,
			)
			return []analysis.Violation{{
				File:    file,
				Line:    imports[i].Line,
				Col:     1,
				Message: msg,
			}}
		}
	}
	return nil
}

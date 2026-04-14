package importresolves

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Check verifies that each imported absolute path exists on disk.
func Check(file string, lines []string, contentRoot string) []analysis.Violation {
	imports := analysis.ParseImports(lines)
	var violations []analysis.Violation
	for _, imp := range imports {
		if !strings.HasPrefix(imp.Path, "/") {
			continue // relative path — skip
		}
		absPath := filepath.Join(contentRoot, imp.Path)
		if _, err := os.Stat(absPath); err != nil {
			violations = append(violations, analysis.Violation{
				File: file,
				Line: imp.Line,
				Col:  1,
				Message: fmt.Sprintf(
					"import-resolves: path %q does not exist on disk",
					imp.Path,
				),
			})
		}
	}
	return violations
}

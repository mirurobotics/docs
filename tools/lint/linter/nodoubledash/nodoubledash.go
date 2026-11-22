package nodoubledash

import "github.com/mirurobotics/docs/tools/lint/linter/analysis"

// Check flags a lone "--" in prose; spans[i] holds the prose spans of line i+1.
func Check(file string, spans [][]analysis.ProseSpan) []analysis.Violation {
	var violations []analysis.Violation
	for lineIdx, lineSpans := range spans {
		line := lineIdx + 1
		for _, span := range lineSpans {
			text := span.Text
			for i := 0; i < len(text)-1; i++ {
				if text[i] != '-' || text[i+1] != '-' {
					continue
				}
				if i > 0 && text[i-1] == '-' {
					continue
				}
				if i+2 < len(text) && text[i+2] == '-' {
					continue
				}
				violations = append(violations, analysis.Violation{
					File:    file,
					Line:    line,
					Col:     span.StartCol + i,
					Message: "no-double-dash: use em dash '\u2014' instead of '--'",
				})
			}
		}
	}
	return violations
}

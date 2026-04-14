package nodoubledash

import "github.com/mirurobotics/docs/tools/lint/linter/analysis"

// Check flags occurrences of exactly two consecutive hyphens ("--") in prose
// that are not part of a longer dash sequence ("---", "----", etc.).
// spans is indexed by line: spans[i] contains prose spans for line i+1.
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
				// Check it's not part of a longer sequence
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

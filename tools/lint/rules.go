package main

// Violation represents a single lint finding.
type Violation struct {
	File    string
	Line    int
	Col     int // 1-based byte column
	Message string
}

// Rule checks prose spans on a single line and returns any violations found.
type Rule interface {
	Check(file string, line int, spans []ProseSpan) []Violation
}

// NoDoubleDash flags occurrences of exactly two consecutive hyphens ("--")
// in prose that are not part of a longer dash sequence ("---", "----", etc.).
type NoDoubleDash struct{}

func (r NoDoubleDash) Check(file string, line int, spans []ProseSpan) []Violation {
	var violations []Violation
	for _, span := range spans {
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
			violations = append(violations, Violation{
				File:    file,
				Line:    line,
				Col:     span.StartCol + i,
				Message: "no-double-dash: use em dash '\u2014' instead of '--'",
			})
		}
	}
	return violations
}

// FileRule checks an entire file's lines at once and returns any violations.
// It is used for rules that require multi-line context, such as import analysis.
type FileRule interface {
	CheckFile(path string, lines []string) []Violation
}

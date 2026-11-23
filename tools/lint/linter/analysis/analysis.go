package analysis

// Violation represents a single lint finding.
type Violation struct {
	File    string
	Line    int
	Col     int // 1-based byte column
	Message string
}

// ProseSpan is a run of prose in a line, starting at 1-based byte column StartCol.
type ProseSpan struct {
	StartCol int
	Text     string
}

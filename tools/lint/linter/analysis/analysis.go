package analysis

// Violation represents a single lint finding.
type Violation struct {
	File    string
	Line    int
	Col     int // 1-based byte column
	Message string
}

// ProseSpan represents a contiguous segment of prose text within a line.
// StartCol is the 1-based byte offset of the span's first character in
// the original line.
type ProseSpan struct {
	StartCol int
	Text     string
}

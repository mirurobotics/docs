package main

import "strings"

// ProseSpan represents a contiguous segment of prose text within a line.
// StartCol is the 1-based byte offset of the span's first character in the original line.
type ProseSpan struct {
	StartCol int
	Text     string
}

type zone int

const (
	zoneProse zone = iota
	zoneFrontmatter
	zoneCodeBlock
	zoneHTMLComment
)

// Scanner is a line-by-line state machine that classifies regions of an MDX
// file as prose or excluded (frontmatter, code blocks, comments, etc.).
type Scanner struct {
	zone      zone
	codeFence string // the fence marker (e.g. "```" or "~~~") to match for closing
	lineNum   int
}

// NewScanner returns a scanner ready to process lines from the start of a file.
func NewScanner() *Scanner {
	return &Scanner{}
}

// LineNum returns the 1-based line number of the most recently scanned line.
func (s *Scanner) LineNum() int {
	return s.lineNum
}

// ScanLine processes the next line of input and returns the prose spans found
// on that line. Returns nil if the entire line is excluded.
func (s *Scanner) ScanLine(line string) []ProseSpan {
	s.lineNum++

	switch s.zone {
	case zoneFrontmatter:
		if strings.TrimRight(line, " \t") == "---" {
			s.zone = zoneProse
		}
		return nil

	case zoneCodeBlock:
		trimmed := strings.TrimSpace(line)
		if trimmed == s.codeFence {
			s.zone = zoneProse
		}
		return nil

	case zoneHTMLComment:
		idx := strings.Index(line, "-->")
		if idx < 0 {
			return nil
		}
		s.zone = zoneProse
		rest := line[idx+3:]
		if len(strings.TrimSpace(rest)) == 0 {
			return nil
		}
		return s.maskInlineRegions(rest, idx+3)
	}

	// zoneProse: check for zone entry
	if s.lineNum == 1 && strings.TrimRight(line, " \t") == "---" {
		s.zone = zoneFrontmatter
		return nil
	}

	if fence, ok := codeFenceOpen(line); ok {
		s.zone = zoneCodeBlock
		s.codeFence = fence
		return nil
	}

	if isThematicBreak(line) {
		return nil
	}

	if isTableSeparator(line) {
		return nil
	}

	if isImportExport(line) {
		return nil
	}

	return s.maskInlineRegions(line, 0)
}

// codeFenceOpen checks if a line opens a fenced code block.
// Returns the bare fence string (e.g. "```") and true, or ("", false).
// MDX allows code fences at any indentation level (nested inside JSX),
// so we don't enforce the CommonMark 0-3 space indentation limit.
func codeFenceOpen(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, " \t")

	if len(trimmed) >= 3 && trimmed[:3] == "```" {
		return "```", true
	}
	if len(trimmed) >= 3 && trimmed[:3] == "~~~" {
		return "~~~", true
	}
	return "", false
}

// isThematicBreak returns true for lines that are markdown thematic breaks
// made of dashes (e.g. "---", "----", "  ---").
func isThematicBreak(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return false
	}
	for _, c := range trimmed {
		if c != '-' {
			return false
		}
	}
	return true
}

// isTableSeparator returns true for markdown table separator rows
// (lines containing only |, -, :, and spaces with at least one |).
func isTableSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}
	hasPipe := false
	for _, c := range trimmed {
		switch c {
		case '|':
			hasPipe = true
		case '-', ':', ' ':
			// allowed
		default:
			return false
		}
	}
	return hasPipe
}

// isImportExport returns true for MDX import/export statements.
func isImportExport(line string) bool {
	return strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "export ")
}

// maskInlineRegions takes a prose line and returns spans with inline code,
// HTML tags, and HTML comments masked out. baseCol is the 0-based byte
// offset of line[0] within the original file line (used when processing
// the tail of an HTML comment line). If an HTML comment opens but does not
// close on this line, the scanner transitions to zoneHTMLComment.
func (s *Scanner) maskInlineRegions(line string, baseCol int) []ProseSpan {
	var spans []ProseSpan
	i := 0
	spanStart := 0

	for i < len(line) {
		switch {
		// HTML comment start
		case i+4 <= len(line) && line[i:i+4] == "<!--":
			if i > spanStart {
				spans = appendSpan(spans, line[spanStart:i], baseCol+spanStart)
			}
			end := strings.Index(line[i+4:], "-->")
			if end >= 0 {
				i = i + 4 + end + 3
			} else {
				// Comment doesn't close on this line — enter comment zone.
				s.zone = zoneHTMLComment
				return spans
			}
			spanStart = i

		// HTML/JSX tag
		case line[i] == '<':
			if i > spanStart {
				spans = appendSpan(spans, line[spanStart:i], baseCol+spanStart)
			}
			end := findTagEnd(line, i)
			i = end
			spanStart = i

		// Inline code
		case line[i] == '`':
			if i > spanStart {
				spans = appendSpan(spans, line[spanStart:i], baseCol+spanStart)
			}
			end := findInlineCodeEnd(line, i)
			i = end
			spanStart = i

		default:
			i++
		}
	}

	if spanStart < len(line) {
		spans = appendSpan(spans, line[spanStart:], baseCol+spanStart)
	}
	return spans
}

// findTagEnd returns the byte index just past the closing '>' of an HTML/JSX
// tag starting at pos. If no closing '>' is found, returns len(line).
func findTagEnd(line string, pos int) int {
	i := pos + 1
	for i < len(line) {
		switch line[i] {
		case '>':
			return i + 1
		case '"':
			i++
			for i < len(line) && line[i] != '"' {
				i++
			}
		case '\'':
			i++
			for i < len(line) && line[i] != '\'' {
				i++
			}
		}
		if i < len(line) {
			i++
		}
	}
	return len(line)
}

// findInlineCodeEnd returns the byte index just past the closing backtick(s)
// of inline code starting at pos. Handles both single and double backtick
// delimiters (` and ``).
func findInlineCodeEnd(line string, pos int) int {
	// Count opening backticks
	ticks := 0
	i := pos
	for i < len(line) && line[i] == '`' {
		ticks++
		i++
	}

	// Find matching closing backticks
	for i < len(line) {
		if line[i] == '`' {
			closeTicks := 0
			j := i
			for j < len(line) && line[j] == '`' {
				closeTicks++
				j++
			}
			if closeTicks == ticks {
				return j
			}
			i = j
		} else {
			i++
		}
	}

	// No closing backticks found; treat the opening backticks as literal text.
	return pos + ticks
}

func appendSpan(spans []ProseSpan, text string, col int) []ProseSpan {
	if len(strings.TrimSpace(text)) == 0 {
		return spans
	}
	return append(spans, ProseSpan{
		StartCol: col + 1, // 1-based
		Text:     text,
	})
}

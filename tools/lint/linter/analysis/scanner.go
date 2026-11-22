package analysis

import "strings"

type zone int

const (
	zoneProse zone = iota
	zoneFrontmatter
	zoneCodeBlock
	zoneHTMLComment
)

// Scanner classifies MDX lines as prose or excluded (frontmatter, code, comments).
type Scanner struct {
	zone      zone
	codeFence string
	lineNum   int
}

// NewScanner returns a scanner ready to process lines from the start of a file.
func NewScanner() *Scanner { return &Scanner{} }

// LineNum returns the 1-based line number of the most recently scanned line.
func (s *Scanner) LineNum() int { return s.lineNum }

// ScanLine returns the prose spans of the next line, or nil if all of it is excluded.
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

	if s.lineNum == 1 && strings.TrimRight(line, " \t") == "---" {
		s.zone = zoneFrontmatter
		return nil
	}

	if fence, ok := CodeFenceOpen(line); ok {
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

// CodeFenceOpen returns the fence a line opens at any indent, as MDX nests them in JSX.
func CodeFenceOpen(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, " \t")

	count := 0
	for count < len(trimmed) && trimmed[count] == '`' {
		count++
	}
	if count >= 3 {
		return trimmed[:count], true
	}

	count = 0
	for count < len(trimmed) && trimmed[count] == '~' {
		count++
	}
	if count >= 3 {
		return trimmed[:count], true
	}

	return "", false
}

func isThematicBreak(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return false
	}
	for _, char := range trimmed {
		if char != '-' {
			return false
		}
	}
	return true
}

func isTableSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}
	hasPipe := false
	for _, char := range trimmed {
		switch char {
		case '|':
			hasPipe = true
		case '-', ':', ' ':
		default:
			return false
		}
	}
	return hasPipe
}

func isImportExport(line string) bool {
	return strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "export ")
}

func (s *Scanner) maskInlineRegions(line string, baseCol int) []ProseSpan {
	var spans []ProseSpan
	i := 0
	spanStart := 0

	for i < len(line) {
		switch {
		case i+4 <= len(line) && line[i:i+4] == "<!--":
			if i > spanStart {
				spans = appendSpan(spans, line[spanStart:i], baseCol+spanStart)
			}
			end := strings.Index(line[i+4:], "-->")
			if end >= 0 {
				i = i + 4 + end + 3
			} else {
				s.zone = zoneHTMLComment
				return spans
			}
			spanStart = i

		case line[i] == '<':
			if i > spanStart {
				spans = appendSpan(spans, line[spanStart:i], baseCol+spanStart)
			}
			end := findTagEnd(line, i)
			i = end
			spanStart = i

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

func findInlineCodeEnd(line string, pos int) int {
	ticks := 0
	i := pos
	for i < len(line) && line[i] == '`' {
		ticks++
		i++
	}

	for i < len(line) {
		if line[i] == '`' {
			closeTicks := 0
			end := i
			for end < len(line) && line[end] == '`' {
				closeTicks++
				end++
			}
			if closeTicks == ticks {
				return end
			}
			i = end
		} else {
			i++
		}
	}

	return pos + ticks
}

func appendSpan(spans []ProseSpan, text string, col int) []ProseSpan {
	if len(strings.TrimSpace(text)) == 0 {
		return spans
	}
	return append(spans, ProseSpan{
		StartCol: col + 1,
		Text:     text,
	})
}

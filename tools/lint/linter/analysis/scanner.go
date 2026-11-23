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
	if s.zone != zoneProse {
		return s.scanZone(line)
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
	if isThematicBreak(line) || isTableSeparator(line) || isImportExport(line) {
		return nil
	}
	return s.maskInlineRegions(line, 0)
}

func (s *Scanner) scanZone(line string) []ProseSpan {
	switch s.zone {
	case zoneFrontmatter:
		if strings.TrimRight(line, " \t") == "---" {
			s.zone = zoneProse
		}
	case zoneCodeBlock:
		if strings.TrimSpace(line) == s.codeFence {
			s.zone = zoneProse
		}
	case zoneHTMLComment:
		end, ok := commentEnd(line, 0)
		if ok {
			s.zone = zoneProse
			return s.maskInlineRegions(line[end:], end)
		}
	}
	return nil
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
		end, ok := s.regionEnd(line, i)
		if !ok {
			i++
			continue
		}
		spans = appendSpan(spans, line[spanStart:i], baseCol+spanStart)
		i = end
		spanStart = i
	}
	if spanStart < len(line) {
		spans = appendSpan(spans, line[spanStart:], baseCol+spanStart)
	}
	return spans
}

// regionEnd runs an unclosed comment to the line's end and into the next line.
func (s *Scanner) regionEnd(line string, pos int) (int, bool) {
	if !strings.HasPrefix(line[pos:], "<!--") {
		return inlineRegionEnd(line, pos)
	}
	end, closed := commentEnd(line, pos+4)
	if !closed {
		s.zone = zoneHTMLComment
		return len(line), true
	}
	return end, true
}

func inlineRegionEnd(line string, pos int) (int, bool) {
	switch line[pos] {
	case '<':
		return findTagEnd(line, pos), true
	case '`':
		return findInlineCodeEnd(line, pos), true
	}
	return pos, false
}

func commentEnd(line string, pos int) (int, bool) {
	idx := strings.Index(line[pos:], "-->")
	if idx < 0 {
		return pos, false
	}
	return pos + idx + 3, true
}

func findTagEnd(line string, pos int) int {
	i := pos + 1
	for i < len(line) {
		switch line[i] {
		case '>':
			return i + 1
		case '"', '\'':
			i = findClosingQuote(line, i)
		}
		if i < len(line) {
			i++
		}
	}
	return len(line)
}

func findClosingQuote(line string, pos int) int {
	i := pos + 1
	for i < len(line) && line[i] != line[pos] {
		i++
	}
	return i
}

func findInlineCodeEnd(line string, pos int) int {
	ticks := countRun(line, pos)
	i := pos + ticks
	for i < len(line) {
		if line[i] != '`' {
			i++
			continue
		}
		closeTicks := countRun(line, i)
		if closeTicks == ticks {
			return i + closeTicks
		}
		i += closeTicks
	}
	return pos + ticks
}

func countRun(line string, pos int) int {
	end := pos
	for end < len(line) && line[end] == line[pos] {
		end++
	}
	return end - pos
}

func appendSpan(spans []ProseSpan, text string, col int) []ProseSpan {
	if len(strings.TrimSpace(text)) == 0 {
		return spans
	}
	return append(spans, ProseSpan{StartCol: col + 1, Text: text})
}

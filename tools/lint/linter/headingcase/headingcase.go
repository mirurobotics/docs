package headingcase

import (
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Message is the diagnostic emitted for any heading-case violation.
const Message = "heading-case: heading must be sentence-case (first letter uppercase, all other letters lowercase); proper nouns/acronyms are not yet supported"

var (
	titleRe       = regexp.MustCompile(`^title:\s*(.*)$`)
	headingRe     = regexp.MustCompile(`^(#{1,6})[ \t]+(.+?)[ \t]*$`)
	inlineCodeRe  = regexp.MustCompile("`[^`]*`")
	htmlTagRe     = regexp.MustCompile(`<[^>]*>`)
	mdLinkRe      = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	trailingPunct = ".?!:"
)

// Check enforces strict sentence-case on the front-matter title and on every
// Markdown body heading. lines is the raw file content split by line; spans is
// indexed by line and gates body-heading detection (lines without prose spans —
// e.g. fenced code blocks, frontmatter — are skipped).
func Check(file string, lines []string, spans [][]analysis.ProseSpan) []analysis.Violation {
	var violations []analysis.Violation

	// Front-matter title path.
	if end := analysis.FrontmatterEnd(lines); end >= 1 {
		for i := 1; i < end; i++ {
			line := lines[i]
			m := titleRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			value, valueByteOffset := titleValueAndOffset(line)
			if casingViolation(value) {
				violations = append(violations, analysis.Violation{
					File:    file,
					Line:    i + 1,
					Col:     valueByteOffset + 1,
					Message: Message,
				})
			}
			break
		}
	}

	// Body heading path.
	for i, line := range lines {
		if i >= len(spans) || len(spans[i]) == 0 {
			continue
		}
		m := headingRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		text := m[2]
		masked := maskHeading(text)
		if !casingViolation(masked) {
			continue
		}
		violations = append(violations, analysis.Violation{
			File:    file,
			Line:    i + 1,
			Col:     headingTextCol(line),
			Message: Message,
		})
	}

	return violations
}

// titleValueAndOffset returns the captured title value (with one matching pair
// of surrounding quotes stripped if present) and the 0-based byte offset of
// the first character of the actual title text on the raw line.
func titleValueAndOffset(line string) (string, int) {
	// Find offset of the first non-space/tab char after "title:".
	// titleRe has already matched, so "title:" is at the very start.
	idx := len("title:")
	for idx < len(line) && (line[idx] == ' ' || line[idx] == '\t') {
		idx++
	}
	value := line[idx:]
	if len(value) >= 2 {
		first := value[0]
		last := value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			value = value[1 : len(value)-1]
			idx++
		}
	}
	return value, idx
}

// headingTextCol returns the 1-based byte column of the first non-'#',
// non-space/tab character on the raw heading line.
func headingTextCol(line string) int {
	i := 0
	for i < len(line) && line[i] == '#' {
		i++
	}
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	return i + 1
}

// maskHeading removes content from heading text that should not contribute to
// the casing check: inline code, HTML/JSX tags, and the URL portion of
// Markdown links (the link text is preserved).
func maskHeading(s string) string {
	s = inlineCodeRe.ReplaceAllString(s, "")
	s = htmlTagRe.ReplaceAllString(s, "")
	s = mdLinkRe.ReplaceAllString(s, "$1")
	return s
}

// casingViolation reports whether s violates strict sentence-case after
// trimming whitespace and trailing terminal punctuation. The first ASCII
// letter must be uppercase; every subsequent ASCII letter must be lowercase.
// Non-letters are ignored. An empty result after trimming is not a violation.
func casingViolation(s string) bool {
	s = strings.TrimSpace(s)
	for len(s) > 0 && strings.ContainsRune(trailingPunct, rune(s[len(s)-1])) {
		s = s[:len(s)-1]
	}
	if s == "" {
		return false
	}
	seenFirst := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		isUpper := c >= 'A' && c <= 'Z'
		isLower := c >= 'a' && c <= 'z'
		if !isUpper && !isLower {
			continue
		}
		if !seenFirst {
			if !isUpper {
				return true
			}
			seenFirst = true
			continue
		}
		if isUpper {
			return true
		}
	}
	return false
}

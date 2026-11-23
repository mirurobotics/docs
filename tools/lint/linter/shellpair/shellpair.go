package shellpair

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

const (
	bashTitle       = "Bash"
	powershellTitle = "PowerShell"
)

// Check flags wrong-shell continuations and missing or differing Bash/PowerShell twins.
func Check(file string, lines []string) []analysis.Violation {
	blocks, groups := scan(lines)
	var violations []analysis.Violation
	for _, b := range blocks {
		if msg := wrongContinuation(b); msg != "" {
			violations = append(violations, violation(file, b.line, msg)...)
		}
	}
	for _, group := range groups {
		violations = append(violations, checkGroup(file, group)...)
	}
	return violations
}

type block struct {
	lang  string
	title string
	line  int
	body  []string
}

func scan(lines []string) ([]*block, [][]*block) {
	var scn scanner
	for i, line := range lines {
		scn.scanLine(i+1, line)
	}
	return scn.blocks, scn.groups
}

func checkGroup(file string, group []*block) []analysis.Violation {
	bash := findBlock(group, bashTitle)
	pwsh := findBlock(group, powershellTitle)
	bashMulti := bash != nil && hasContinuation(bash.body, '\\')
	pwshMulti := pwsh != nil && hasContinuation(pwsh.body, '`')
	switch {
	case bashMulti && pwsh == nil:
		return violation(file, bash.line,
			"multi-line Bash command needs a PowerShell twin")
	case pwshMulti && bash == nil:
		return violation(file, pwsh.line,
			"multi-line PowerShell command needs a Bash twin")
	case !bashMulti && !pwshMulti:
		return nil
	case wrongContinuation(bash) != "" || wrongContinuation(pwsh) != "":
		return nil
	case areTwinsDifferent(bash, pwsh):
		return violation(file, pwsh.line,
			"PowerShell command differs from its Bash twin")
	}
	return nil
}

func findBlock(group []*block, title string) *block {
	for _, b := range group {
		if b.title == title {
			return b
		}
	}
	return nil
}

func violation(file string, line int, msg string) []analysis.Violation {
	return []analysis.Violation{{
		File: file, Line: line, Col: 1, Message: "shell-pair: " + msg,
	}}
}

func wrongContinuation(b *block) string {
	switch {
	case b.lang == "bash" && hasContinuation(b.body, '`'):
		return "bash block continues a line with PowerShell's backtick"
	case b.lang == "powershell" && hasContinuation(b.body, '\\'):
		return "powershell block continues a line with a Bash backslash"
	}
	return ""
}

func hasContinuation(body []string, cont byte) bool {
	for _, line := range body {
		if strings.HasSuffix(strings.TrimRight(line, " \t"), " "+string(cont)) {
			return true
		}
	}
	return false
}

func areTwinsDifferent(bash, pwsh *block) bool {
	pwshCommand := strings.ReplaceAll(normalize(pwsh.body, '`'), `\`, "/")
	return normalize(bash.body, '\\') != pwshCommand
}

func normalize(body []string, cont byte) string {
	var parts []string
	for _, line := range body {
		line = strings.TrimRight(line, " \t")
		line = strings.TrimSuffix(line, string(cont))
		if line = strings.TrimSpace(line); line != "" {
			parts = append(parts, line)
		}
	}
	return strings.Join(parts, " ")
}

type scanner struct {
	blocks  []*block
	groups  [][]*block
	current []*block
	open    *block
	fence   string
	inGroup bool
}

func (s *scanner) scanLine(lineNum int, line string) {
	trimmed := strings.TrimSpace(line)
	if s.open != nil && trimmed == s.fence {
		s.open = nil
		return
	}
	if s.open != nil {
		s.open.body = append(s.open.body, line)
		return
	}
	if marker, ok := analysis.CodeFenceOpen(line); ok {
		s.openBlock(lineNum, trimmed, marker)
		return
	}
	switch {
	case strings.HasPrefix(trimmed, "<CodeGroup"):
		s.inGroup, s.current = true, nil
	case trimmed == "</CodeGroup>" && s.inGroup:
		s.groups = append(s.groups, s.current)
		s.inGroup = false
	}
}

func (s *scanner) openBlock(lineNum int, trimmed, marker string) {
	s.fence = marker
	lang, title := fenceInfo(trimmed, marker)
	s.open = &block{lang: lang, title: title, line: lineNum}
	s.blocks = append(s.blocks, s.open)
	if s.inGroup {
		s.current = append(s.current, s.open)
	}
}

func fenceInfo(trimmed, fence string) (lang, title string) {
	fields := strings.Fields(strings.TrimPrefix(trimmed, fence))
	if len(fields) > 0 {
		lang = fields[0]
	}
	if len(fields) > 1 {
		title = fields[1]
	}
	return lang, title
}

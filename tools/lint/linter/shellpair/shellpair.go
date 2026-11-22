package shellpair

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

const (
	bashTitle       = "Bash"
	powershellTitle = "PowerShell"
)

type block struct {
	lang  string
	title string
	line  int
	body  []string
}

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

func scan(lines []string) ([]*block, [][]*block) {
	var blocks []*block
	var groups [][]*block
	var current []*block
	var open *block
	inGroup, fence := false, ""
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if open != nil {
			if trimmed == fence {
				open = nil
			} else {
				open.body = append(open.body, line)
			}
			continue
		}
		if marker, ok := analysis.CodeFenceOpen(line); ok {
			fence = marker
			lang, title := fenceInfo(trimmed, marker)
			open = &block{lang: lang, title: title, line: i + 1}
			blocks = append(blocks, open)
			if inGroup {
				current = append(current, open)
			}
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "<CodeGroup"):
			inGroup, current = true, nil
		case trimmed == "</CodeGroup>" && inGroup:
			groups = append(groups, current)
			inGroup = false
		}
	}
	return blocks, groups
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

func wrongContinuation(b *block) string {
	switch {
	case b.lang == "bash" && hasContinuation(b.body, '`'):
		return "bash block continues a line with PowerShell's backtick"
	case b.lang == "powershell" && hasContinuation(b.body, '\\'):
		return "powershell block continues a line with a Bash backslash"
	}
	return ""
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
	case (bashMulti || pwshMulti) && bash != nil && pwsh != nil &&
		wrongContinuation(bash) == "" && wrongContinuation(pwsh) == "" &&
		normalize(bash.body, '\\') !=
			strings.ReplaceAll(normalize(pwsh.body, '`'), `\`, "/"):
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

func hasContinuation(body []string, cont byte) bool {
	for _, line := range body {
		if strings.HasSuffix(strings.TrimRight(line, " \t"), " "+string(cont)) {
			return true
		}
	}
	return false
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

func violation(file string, line int, msg string) []analysis.Violation {
	return []analysis.Violation{{
		File: file, Line: line, Col: 1, Message: "shell-pair: " + msg,
	}}
}

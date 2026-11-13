// Package shellpair keeps the Unix and PowerShell versions of a
// multi-line command correct and in step. A bash block must not continue
// lines with PowerShell's trailing "`", and a powershell or pwsh block
// must not continue them with a trailing "\". Inside a <CodeGroup>, a fenced block
// titled "Unix" that continues lines with "\" needs a block titled
// "PowerShell" with the same command, and a "PowerShell" block that
// continues lines with "`" needs a "Unix" twin. The two match when they
// are equal once each line's continuation character and surrounding
// whitespace are removed.
package shellpair

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

const (
	unixTitle       = "Unix"
	powershellTitle = "PowerShell"
)

// block is a fenced code block.
type block struct {
	lang  string
	title string
	line  int // 1-based line of the opening fence
	body  []string
}

// Check flags bash and PowerShell blocks that use the other shell's line
// continuation, and Unix and PowerShell blocks in the same CodeGroup whose
// multi-line command lacks a twin or differs from it.
func Check(file string, lines []string) []analysis.Violation {
	blocks, groups := scan(lines)
	var violations []analysis.Violation
	for _, b := range blocks {
		violations = append(violations, checkContinuation(file, b)...)
	}
	for _, group := range groups {
		violations = append(violations, checkGroup(file, group)...)
	}
	return violations
}

// scan returns every fenced block in the file, and the blocks of each
// CodeGroup.
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
		if f, ok := analysis.CodeFenceOpen(line); ok {
			fence = f
			lang, title := fenceInfo(trimmed, f)
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

// fenceInfo returns the language and title of an opening fence such as
// "```bash Unix highlight={3}", which are its first two words.
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

func checkContinuation(file string, b *block) []analysis.Violation {
	switch {
	case b.lang == "bash" && hasContinuation(b.body, '`'):
		return violation(file, b.line,
			"bash block continues a line with PowerShell's backtick")
	case isPowerShell(b.lang) && hasContinuation(b.body, '\\'):
		return violation(file, b.line,
			b.lang+" block continues a line with a Unix backslash")
	}
	return nil
}

// isPowerShell reports whether the fence language is PowerShell: the
// built-in grammar, or the docs' pwsh grammar that also colors miru commands.
func isPowerShell(lang string) bool { return lang == "powershell" || lang == "pwsh" }

func checkGroup(file string, group []*block) []analysis.Violation {
	unix := findBlock(group, unixTitle)
	ps := findBlock(group, powershellTitle)
	unixMulti := unix != nil && hasContinuation(unix.body, '\\')
	psMulti := ps != nil && hasContinuation(ps.body, '`')
	switch {
	case unixMulti && ps == nil:
		return violation(file, unix.line,
			"multi-line Unix command needs a PowerShell twin")
	case psMulti && unix == nil:
		return violation(file, ps.line,
			"multi-line PowerShell command needs a Unix twin")
	case (unixMulti || psMulti) && unix != nil && ps != nil &&
		normalize(unix.body, '\\') != normalize(ps.body, '`'):
		return violation(file, ps.line, "PowerShell command differs from its Unix twin")
	}
	return nil
}

// findBlock returns the group's first block with the title. Titles name
// a CodeGroup's tabs, so each appears once.
func findBlock(group []*block, title string) *block {
	for _, b := range group {
		if b.title == title {
			return b
		}
	}
	return nil
}

// hasContinuation reports whether any line ends with a space and the
// line-continuation character.
func hasContinuation(body []string, cont byte) bool {
	for _, line := range body {
		if strings.HasSuffix(strings.TrimRight(line, " \t"), " "+string(cont)) {
			return true
		}
	}
	return false
}

// normalize joins the command's lines with single spaces, dropping
// each line's continuation character and surrounding whitespace.
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

// Package shellpair keeps the Unix and PowerShell versions of a
// multi-line command in step. Inside a <CodeGroup>, a fenced block
// titled "Unix" that continues lines with a trailing "\" needs a block
// titled "PowerShell" with the same command, and a "PowerShell" block
// that continues lines with a trailing "`" needs a "Unix" twin. The two
// match when they are equal once each line's continuation character
// and surrounding whitespace are removed.
package shellpair

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

const (
	unixTitle       = "Unix"
	powershellTitle = "PowerShell"
)

// block is a fenced code block inside a CodeGroup.
type block struct {
	title string
	line  int // 1-based line of the opening fence
	body  []string
}

// Check flags Unix and PowerShell blocks in the same CodeGroup whose
// multi-line command lacks a twin or differs from it.
func Check(file string, lines []string) []analysis.Violation {
	var violations []analysis.Violation
	for _, group := range codeGroups(lines) {
		violations = append(violations, checkGroup(file, group)...)
	}
	return violations
}

// codeGroups returns the fenced blocks of each CodeGroup in the file.
func codeGroups(lines []string) [][]block {
	var groups [][]block
	var current []block
	inGroup, inFence := false, false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case inFence && trimmed == "```":
			inFence = false
		case inFence:
			if inGroup {
				last := &current[len(current)-1]
				last.body = append(last.body, line)
			}
		case strings.HasPrefix(trimmed, "```"):
			inFence = true
			if inGroup {
				b := block{title: fenceTitle(trimmed), line: i + 1}
				current = append(current, b)
			}
		case strings.HasPrefix(trimmed, "<CodeGroup"):
			inGroup, current = true, nil
		case trimmed == "</CodeGroup>":
			if inGroup {
				groups = append(groups, current)
			}
			inGroup = false
		}
	}
	return groups
}

// fenceTitle returns the title of an opening fence such as
// "```bash Unix highlight={3}", which is the word after the language.
func fenceTitle(fence string) string {
	fields := strings.Fields(strings.TrimPrefix(fence, "```"))
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}

func checkGroup(file string, group []block) []analysis.Violation {
	unix, hasUnix := findBlock(group, unixTitle)
	ps, hasPS := findBlock(group, powershellTitle)
	unixMulti := hasUnix && isMultiLine(unix.body, '\\')
	psMulti := hasPS && isMultiLine(ps.body, '`')
	var line int
	var msg string
	switch {
	case unixMulti && !hasPS:
		line, msg = unix.line, "multi-line Unix command needs a PowerShell twin"
	case psMulti && !hasUnix:
		line, msg = ps.line, "multi-line PowerShell command needs a Unix twin"
	case (unixMulti || psMulti) && hasUnix && hasPS &&
		normalize(unix.body, '\\') != normalize(ps.body, '`'):
		line, msg = ps.line, "PowerShell command differs from its Unix twin"
	default:
		return nil
	}
	return []analysis.Violation{{
		File: file, Line: line, Col: 1, Message: "shell-pair: " + msg,
	}}
}

func findBlock(group []block, title string) (block, bool) {
	for _, b := range group {
		if b.title == title {
			return b, true
		}
	}
	return block{}, false
}

// isMultiLine reports whether any line ends with a space and the
// shell's line-continuation character.
func isMultiLine(body []string, cont byte) bool {
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

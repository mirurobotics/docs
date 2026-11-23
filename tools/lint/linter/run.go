package linter

import (
	"bufio"
	"os"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
	"github.com/mirurobotics/docs/tools/lint/linter/componentstyle"
	"github.com/mirurobotics/docs/tools/lint/linter/headingcase"
	"github.com/mirurobotics/docs/tools/lint/linter/imagedomain"
	"github.com/mirurobotics/docs/tools/lint/linter/importblock"
	"github.com/mirurobotics/docs/tools/lint/linter/importresolves"
	"github.com/mirurobotics/docs/tools/lint/linter/importsorted"
	"github.com/mirurobotics/docs/tools/lint/linter/importused"
	"github.com/mirurobotics/docs/tools/lint/linter/mdxstyle"
	"github.com/mirurobotics/docs/tools/lint/linter/nodoubledash"
	"github.com/mirurobotics/docs/tools/lint/linter/shellpair"
)

const maxScanTokenSize = 16 * 1024 * 1024

// Rule identifies a linter rule.
type Rule string

const (
	RuleNoDoubleDash   Rule = "no-double-dash"
	RuleHeadingCase    Rule = "heading-case"
	RuleImportResolves Rule = "import-resolves"
	RuleImportUsed     Rule = "import-used"
	RuleImportSorted   Rule = "import-sorted"
	RuleComponentStyle Rule = "component-style"
	RuleMDXStyle       Rule = "mdx-style"
	RuleImportBlock    Rule = "import-block"
	RuleImageDomain    Rule = "image-domain"
	RuleShellPair      Rule = "shell-pair"
)

type checkInput struct {
	path        string
	lines       []string
	spans       [][]analysis.ProseSpan
	contentRoot string
}

// ProcessFile lints a single file and returns violations.
func ProcessFile(path, contentRoot string) ([]analysis.Violation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), maxScanTokenSize)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	mdxScanner := analysis.NewScanner()
	spans := make([][]analysis.ProseSpan, len(lines))
	for i, line := range lines {
		spans[i] = mdxScanner.ScanLine(line)
	}

	in := checkInput{
		path:        path,
		lines:       lines,
		spans:       spans,
		contentRoot: contentRoot,
	}
	return runChecks(in), nil
}

func runChecks(in checkInput) []analysis.Violation {
	var violations []analysis.Violation
	for _, entry := range ruleCheckers() {
		violations = append(violations, entry.check(in)...)
	}
	return violations
}

type ruleEntry struct {
	rule  Rule
	check func(checkInput) []analysis.Violation
}

func ruleCheckers() []ruleEntry {
	return []ruleEntry{
		{RuleNoDoubleDash, func(in checkInput) []analysis.Violation {
			return nodoubledash.Check(in.path, in.spans)
		}},
		{RuleHeadingCase, func(in checkInput) []analysis.Violation {
			return headingcase.Check(in.path, in.lines, in.spans)
		}},
		{RuleImportResolves, func(in checkInput) []analysis.Violation {
			return importresolves.Check(in.path, in.lines, in.contentRoot)
		}},
		{RuleImportUsed, func(in checkInput) []analysis.Violation {
			return importused.Check(in.path, in.lines)
		}},
		{RuleImportSorted, func(in checkInput) []analysis.Violation {
			return importsorted.Check(in.path, in.lines)
		}},
		{RuleComponentStyle, func(in checkInput) []analysis.Violation {
			return componentstyle.Check(in.path, in.lines)
		}},
		{RuleMDXStyle, func(in checkInput) []analysis.Violation {
			return mdxstyle.Check(in.path, in.lines)
		}},
		{RuleImportBlock, func(in checkInput) []analysis.Violation {
			return importblock.Check(in.path, in.lines)
		}},
		{RuleImageDomain, func(in checkInput) []analysis.Violation {
			return imagedomain.Check(in.path, in.lines)
		}},
		{RuleShellPair, func(in checkInput) []analysis.Violation {
			return shellpair.Check(in.path, in.lines)
		}},
	}
}

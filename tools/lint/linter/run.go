package linter

import (
	"bufio"
	"os"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
	"github.com/mirurobotics/docs/tools/lint/linter/componentstyle"
	"github.com/mirurobotics/docs/tools/lint/linter/importblock"
	"github.com/mirurobotics/docs/tools/lint/linter/importresolves"
	"github.com/mirurobotics/docs/tools/lint/linter/importsorted"
	"github.com/mirurobotics/docs/tools/lint/linter/importused"
	"github.com/mirurobotics/docs/tools/lint/linter/mdxstyle"
	"github.com/mirurobotics/docs/tools/lint/linter/nodoubledash"
)

// Rule identifies a linter rule.
type Rule string

const (
	RuleNoDoubleDash   Rule = "no-double-dash"
	RuleImportResolves Rule = "import-resolves"
	RuleImportUsed     Rule = "import-used"
	RuleImportSorted   Rule = "import-sorted"
	RuleComponentStyle Rule = "component-style"
	RuleMDXStyle       Rule = "mdx-style"
	RuleImportBlock    Rule = "import-block"
	RuleRedirects      Rule = "redirects"
)

// AllRules returns every valid rule name.
func AllRules() []Rule {
	return []Rule{
		RuleNoDoubleDash,
		RuleImportResolves, RuleImportUsed, RuleImportSorted,
		RuleComponentStyle, RuleMDXStyle, RuleImportBlock,
		RuleRedirects,
	}
}

type checkInput struct {
	path        string
	lines       []string
	spans       [][]analysis.ProseSpan
	contentRoot string
}

type ruleEntry struct {
	rule  Rule
	check func(checkInput) []analysis.Violation
}

// Redirects is invoked once per run from main.go (see linter.ProcessDocsJSON), not per-file via ruleCheckers, because it operates on docs.json once.
func ruleCheckers() []ruleEntry {
	return []ruleEntry{
		{RuleNoDoubleDash, func(in checkInput) []analysis.Violation {
			return nodoubledash.Check(in.path, in.spans)
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
	}
}

func runChecks(in checkInput) []analysis.Violation {
	var violations []analysis.Violation
	for _, rc := range ruleCheckers() {
		violations = append(violations, rc.check(in)...)
	}
	return violations
}

// ProcessFile lints a single file and returns violations.
func ProcessFile(path, contentRoot string) ([]analysis.Violation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var lines []string
	scanner := bufio.NewScanner(f)
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

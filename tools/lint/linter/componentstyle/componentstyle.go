package componentstyle

import (
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Check validates style of imports from /snippets/components/.
func Check(file string, lines []string) []analysis.Violation {
	imports := analysis.ParseImports(lines)
	var violations []analysis.Violation
	for _, imp := range imports {
		if !strings.HasPrefix(imp.Path, "/snippets/components/") {
			continue
		}
		violations = append(violations, checkStyle(file, imp)...)
	}
	return violations
}

func checkStyle(file string, imp analysis.ParsedImport) []analysis.Violation {
	raw := imp.Raw
	lineNum := imp.Line
	var violations []analysis.Violation

	openBrace := strings.Index(raw, "{")
	if openBrace < 0 {
		msg := "import-component-style: component import must use named import syntax { }"
		violations = append(violations, analysis.Violation{
			File: file, Line: lineNum, Col: 1, Message: msg,
		})
	} else {
		violations = append(
			violations, checkBracketSpacing(file, lineNum, raw, openBrace)...,
		)
	}

	if !strings.HasSuffix(imp.Path, ".jsx") {
		msg := "import-component-style: component import path must end in '.jsx'"
		violations = append(violations, analysis.Violation{
			File: file, Line: lineNum, Col: 1, Message: msg,
		})
	}
	if !strings.HasSuffix(strings.TrimRight(raw, " \t"), ";") {
		msg := "import-component-style: import statement must end with ';'"
		violations = append(violations, analysis.Violation{
			File: file, Line: lineNum, Col: 1, Message: msg,
		})
	}
	return violations
}

func checkBracketSpacing(
	file string,
	lineNum int,
	raw string,
	openBrace int,
) []analysis.Violation {
	closeBrace := strings.LastIndex(raw, "}")
	if closeBrace <= openBrace {
		msg := "import-component-style: component import must use named import syntax { }"
		return []analysis.Violation{{File: file, Line: lineNum, Col: 1, Message: msg}}
	}
	body := raw[openBrace+1 : closeBrace]
	var violations []analysis.Violation
	if !strings.HasPrefix(body, " ") {
		violations = append(violations, analysis.Violation{
			File:    file,
			Line:    lineNum,
			Col:     openBrace + 2,
			Message: "import-component-style: missing space after '{'",
		})
	}
	if !strings.HasSuffix(body, " ") {
		violations = append(violations, analysis.Violation{
			File:    file,
			Line:    lineNum,
			Col:     closeBrace + 1,
			Message: "import-component-style: missing space before '}'",
		})
	}
	violations = append(
		violations, checkCommaSpacing(file, lineNum, body, openBrace)...,
	)
	return violations
}

func checkCommaSpacing(
	file string,
	lineNum int,
	body string,
	openBrace int,
) []analysis.Violation {
	var violations []analysis.Violation
	for i := range body {
		if body[i] != ',' {
			continue
		}
		if i > 0 && body[i-1] == ' ' {
			violations = append(violations, analysis.Violation{
				File:    file,
				Line:    lineNum,
				Col:     openBrace + 1 + i,
				Message: "import-component-style: unexpected space before ','",
			})
		}
		after := i + 1
		tooFew := after >= len(body) || body[after] != ' '
		tooMany := !tooFew && after+1 < len(body) && body[after+1] == ' '
		if tooFew || tooMany {
			violations = append(violations, analysis.Violation{
				File:    file,
				Line:    lineNum,
				Col:     openBrace + 2 + i,
				Message: "import-component-style: expected single space after ','",
			})
		}
	}
	return violations
}

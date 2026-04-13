package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ImportedName represents a single binding from an import statement.
type ImportedName struct {
	Name    string
	IsNamed bool // true if from { } syntax (named import)
}

// ParsedImport holds the parsed representation of one MDX import line.
type ParsedImport struct {
	Line  int            // 1-based line number
	Raw   string         // original line text
	Names []ImportedName // extracted bindings
	Path  string         // import path without quotes or trailing semicolon
}

// parseSingleImport parses one import line. Returns nil if unparseable.
func parseSingleImport(lineNum int, line string) *ParsedImport {
	if !strings.HasPrefix(line, "import ") {
		return nil
	}
	rest := strings.TrimLeft(strings.TrimPrefix(line, "import "), " \t")

	var names []ImportedName
	if strings.HasPrefix(rest, "{") {
		end := strings.Index(rest, "}")
		if end < 0 {
			return nil
		}
		for _, part := range strings.Split(rest[1:end], ",") {
			if name := strings.TrimSpace(part); name != "" {
				names = append(names, ImportedName{Name: name, IsNamed: true})
			}
		}
		rest = rest[end+1:]
	} else {
		fromIdx := strings.Index(rest, " from ")
		if fromIdx < 0 {
			return nil
		}
		token := strings.TrimSpace(rest[:fromIdx])
		if token == "" {
			return nil
		}
		names = append(names, ImportedName{Name: token})
		rest = rest[fromIdx:]
	}

	path, ok := extractImportPath(rest)
	if !ok {
		return nil
	}
	return &ParsedImport{Line: lineNum, Raw: line, Names: names, Path: path}
}

// extractImportPath finds the last quoted string in rest (the "from '...'" part).
// Returns the unquoted path and true, or ("", false) if not found.
func extractImportPath(rest string) (string, bool) {
	lastQuote := -1
	var quoteChar byte
	for i := len(rest) - 1; i >= 0; i-- {
		if rest[i] == '\'' || rest[i] == '"' {
			lastQuote = i
			quoteChar = rest[i]
			break
		}
	}
	if lastQuote < 0 {
		return "", false
	}
	for i := lastQuote - 1; i >= 0; i-- {
		if rest[i] == quoteChar {
			return rest[i+1 : lastQuote], true
		}
	}
	return "", false
}

// parseImports extracts all successfully parsed imports from a slice of lines.
func parseImports(lines []string) []ParsedImport {
	var imports []ParsedImport
	for i, line := range lines {
		if !isImportLine(line) {
			continue
		}
		if pi := parseSingleImport(i+1, line); pi != nil {
			imports = append(imports, *pi)
		}
	}
	return imports
}

func isImportLine(line string) bool { return strings.HasPrefix(line, "import ") }

// frontmatterEnd returns the 0-based index of the closing "---" line,
// or -1 if no frontmatter block is present.
func frontmatterEnd(lines []string) int {
	if len(lines) == 0 || strings.TrimRight(lines[0], " \t") != "---" {
		return -1
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], " \t") == "---" {
			return i
		}
	}
	return -1
}

// bodyLines returns non-frontmatter, non-import lines for word-search purposes.
func bodyLines(lines []string) []string {
	fmEnd := frontmatterEnd(lines)
	var body []string
	for i, line := range lines {
		if fmEnd >= 0 && i <= fmEnd {
			continue
		}
		if isImportLine(line) {
			continue
		}
		body = append(body, line)
	}
	return body
}

// ImportResolvesRule checks that each imported path exists on disk.
type ImportResolvesRule struct {
	ContentRoot string
}

func (r ImportResolvesRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	var violations []Violation
	for _, imp := range imports {
		if !strings.HasPrefix(imp.Path, "/") {
			continue // relative path — skip
		}
		absPath := filepath.Join(r.ContentRoot, imp.Path)
		if _, err := os.Stat(absPath); err != nil {
			violations = append(violations, Violation{
				File: path,
				Line: imp.Line,
				Col:  1,
				Message: fmt.Sprintf(
					"import-resolves: path %q does not exist on disk",
					imp.Path,
				),
			})
		}
	}
	return violations
}

// ImportUsedRule checks that every imported name is used in the body.
type ImportUsedRule struct{}

func (r ImportUsedRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	if len(imports) == 0 {
		return nil
	}
	body := strings.Join(bodyLines(lines), "\n")
	var violations []Violation
	for _, imp := range imports {
		for _, n := range imp.Names {
			re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(n.Name)))
			if re.MatchString(body) {
				continue
			}
			msg := fmt.Sprintf("import-used: %q is imported but never used", n.Name)
			violations = append(violations, Violation{
				File: path, Line: imp.Line, Col: 1, Message: msg,
			})
		}
	}
	return violations
}

// ImportSortedRule checks imports are sorted by path (case-insensitive, ascending).
type ImportSortedRule struct{}

func (r ImportSortedRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	for i := 1; i < len(imports); i++ {
		if strings.ToLower(imports[i].Path) < strings.ToLower(imports[i-1].Path) {
			msg := fmt.Sprintf(
				"import-sorted: import path %q is out of order (expected after %q)",
				imports[i].Path,
				imports[i-1].Path,
			)
			return []Violation{{
				File:    path,
				Line:    imports[i].Line,
				Col:     1,
				Message: msg,
			}}
		}
	}
	return nil
}

// ComponentImportStyleRule checks style of imports from /snippets/components/.
type ComponentImportStyleRule struct{}

func (r ComponentImportStyleRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	var violations []Violation
	for _, imp := range imports {
		if !strings.HasPrefix(imp.Path, "/snippets/components/") {
			continue
		}
		violations = append(violations, checkComponentStyle(path, imp)...)
	}
	return violations
}

// checkComponentStyle validates a single component import line.
func checkComponentStyle(path string, imp ParsedImport) []Violation {
	raw := imp.Raw
	ln := imp.Line
	var vs []Violation

	openBrace := strings.Index(raw, "{")
	if openBrace < 0 {
		msg := "import-component-style: component import must use named import syntax { }"
		vs = append(vs, Violation{File: path, Line: ln, Col: 1, Message: msg})
	} else {
		vs = append(vs, checkBracketSpacing(path, ln, raw, openBrace)...)
	}

	if !strings.HasSuffix(imp.Path, ".jsx") {
		msg := "import-component-style: component import path must end in '.jsx'"
		vs = append(vs, Violation{File: path, Line: ln, Col: 1, Message: msg})
	}
	if !strings.HasSuffix(strings.TrimRight(raw, " \t"), ";") {
		msg := "import-component-style: import statement must end with ';'"
		vs = append(vs, Violation{File: path, Line: ln, Col: 1, Message: msg})
	}
	return vs
}

// checkBracketSpacing validates brace spacing and comma-space style.
func checkBracketSpacing(
	path string,
	lineNum int,
	raw string,
	openBrace int,
) []Violation {
	closeBrace := strings.LastIndex(raw, "}")
	if closeBrace < 0 {
		msg := "import-component-style: component import must use named import syntax { }"
		return []Violation{{File: path, Line: lineNum, Col: 1, Message: msg}}
	}
	body := raw[openBrace+1 : closeBrace]
	var vs []Violation
	if !strings.HasPrefix(body, " ") {
		vs = append(vs, Violation{
			File:    path,
			Line:    lineNum,
			Col:     openBrace + 2,
			Message: "import-component-style: missing space after '{'",
		})
	}
	if !strings.HasSuffix(body, " ") {
		vs = append(vs, Violation{
			File:    path,
			Line:    lineNum,
			Col:     closeBrace + 1,
			Message: "import-component-style: missing space before '}'",
		})
	}
	vs = append(vs, checkCommaSpacing(path, lineNum, body, openBrace)...)
	return vs
}

// checkCommaSpacing validates commas have no preceding space and one following space.
func checkCommaSpacing(
	path string,
	lineNum int,
	body string,
	openBrace int,
) []Violation {
	var vs []Violation
	for i := range body {
		if body[i] != ',' {
			continue
		}
		if i > 0 && body[i-1] == ' ' {
			vs = append(vs, Violation{
				File:    path,
				Line:    lineNum,
				Col:     openBrace + 1 + i,
				Message: "import-component-style: unexpected space before ','",
			})
		}
		after := i + 1
		tooFew := after >= len(body) || body[after] != ' '
		tooMany := !tooFew && after+1 < len(body) && body[after+1] == ' '
		if tooFew || tooMany {
			vs = append(vs, Violation{
				File:    path,
				Line:    lineNum,
				Col:     openBrace + 2 + i,
				Message: "import-component-style: expected single space after ','",
			})
		}
	}
	return vs
}

// MDXImportStyleRule checks style of imports ending in .mdx.
type MDXImportStyleRule struct{}

func (r MDXImportStyleRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	var violations []Violation
	for _, imp := range imports {
		if !strings.HasSuffix(imp.Path, ".mdx") {
			continue
		}
		if strings.Contains(imp.Raw, "{") {
			msg := "import-mdx-style: .mdx imports must use default import syntax"
			violations = append(violations, Violation{
				File: path, Line: imp.Line, Col: 1, Message: msg,
			})
		}
		if !strings.HasSuffix(strings.TrimRight(imp.Raw, " \t"), ";") {
			violations = append(violations, Violation{
				File:    path,
				Line:    imp.Line,
				Col:     1,
				Message: "import-mdx-style: import statement must end with ';'",
			})
		}
	}
	return violations
}

// ImportBlockContiguousRule checks that no blank lines appear inside a contiguous
// import block. Two imports are "in the same block" when only blank lines (and
// no other content) appear between them.
type ImportBlockContiguousRule struct{}

func (r ImportBlockContiguousRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	if len(imports) < 2 {
		return nil
	}
	var violations []Violation
	for i := 1; i < len(imports); i++ {
		prev := imports[i-1].Line // 1-based
		curr := imports[i].Line   // 1-based
		// Check whether the lines strictly between prev and curr are all blank.
		// If any non-blank, non-import line exists, these imports are in
		// separate blocks and no violation is reported.
		onlyBlanks := true
		for j := prev; j < curr-1; j++ { // j is 0-based index = line (j+1)
			line := lines[j]
			if strings.TrimSpace(line) != "" && !isImportLine(line) {
				onlyBlanks = false
				break
			}
		}
		if !onlyBlanks {
			continue
		}
		// All inter-import lines are blank — flag them.
		for j := prev; j < curr-1; j++ {
			if strings.TrimSpace(lines[j]) == "" {
				violations = append(violations, Violation{
					File:    path,
					Line:    j + 1,
					Col:     1,
					Message: "import-block-contiguous: blank line inside import block",
				})
			}
		}
	}
	return violations
}

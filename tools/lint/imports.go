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
	rest := strings.TrimPrefix(line, "import ")
	rest = strings.TrimLeft(rest, " \t")

	var names []ImportedName
	var isNamed bool

	if strings.HasPrefix(rest, "{") {
		// Named import: import { A, B } from '/path'
		isNamed = true
		end := strings.Index(rest, "}")
		if end < 0 {
			return nil
		}
		bracketContent := rest[1:end]
		for _, part := range strings.Split(bracketContent, ",") {
			name := strings.TrimSpace(part)
			if name == "" {
				continue
			}
			names = append(names, ImportedName{Name: name, IsNamed: true})
		}
		rest = rest[end+1:]
	} else {
		// Default import: import Name from '/path'
		fromIdx := strings.Index(rest, " from ")
		if fromIdx < 0 {
			return nil
		}
		token := strings.TrimSpace(rest[:fromIdx])
		if token == "" {
			return nil
		}
		names = append(names, ImportedName{Name: token, IsNamed: false})
		rest = rest[fromIdx:]
	}

	// Find the last quoted string for the path
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
		return nil
	}
	// Find matching opening quote scanning backwards from lastQuote-1
	openQuote := -1
	for i := lastQuote - 1; i >= 0; i-- {
		if rest[i] == quoteChar {
			openQuote = i
			break
		}
	}
	if openQuote < 0 {
		return nil
	}
	path := rest[openQuote+1 : lastQuote]
	// Strip trailing semicolon from path (shouldn't be inside quotes, but be safe)
	path = strings.TrimSuffix(path, ";")

	_ = isNamed
	return &ParsedImport{
		Line:  lineNum,
		Raw:   line,
		Names: names,
		Path:  path,
	}
}

// parseImports extracts all successfully parsed imports from a slice of lines.
func parseImports(lines []string) []ParsedImport {
	var imports []ParsedImport
	for i, line := range lines {
		if !isImportLine(line) {
			continue
		}
		pi := parseSingleImport(i+1, line)
		if pi != nil {
			imports = append(imports, *pi)
		}
	}
	return imports
}

// isImportLine returns true if the line is an MDX import statement.
func isImportLine(line string) bool {
	return strings.HasPrefix(line, "import ")
}

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

// bodyLines returns the non-frontmatter, non-import lines as a joined string
// for word-search purposes, along with a map from 0-based line index to content.
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
			// Relative path — skip
			continue
		}
		absPath := filepath.Join(r.ContentRoot, imp.Path)
		if _, err := os.Stat(absPath); err != nil {
			violations = append(violations, Violation{
				File:    path,
				Line:    imp.Line,
				Col:     1,
				Message: fmt.Sprintf("import-resolves: path %q does not exist on disk", imp.Path),
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
			if !re.MatchString(body) {
				violations = append(violations, Violation{
					File:    path,
					Line:    imp.Line,
					Col:     1,
					Message: fmt.Sprintf("import-used: %q is imported but never used", n.Name),
				})
			}
		}
	}
	return violations
}

// ImportSortedRule checks that imports are sorted by path (case-insensitive, ascending).
type ImportSortedRule struct{}

func (r ImportSortedRule) CheckFile(path string, lines []string) []Violation {
	imports := parseImports(lines)
	if len(imports) < 2 {
		return nil
	}
	for i := 1; i < len(imports); i++ {
		prev := strings.ToLower(imports[i-1].Path)
		curr := strings.ToLower(imports[i].Path)
		if curr < prev {
			return []Violation{{
				File:    path,
				Line:    imports[i].Line,
				Col:     1,
				Message: fmt.Sprintf("import-sorted: import path %q is out of order (expected after %q)", imports[i].Path, imports[i-1].Path),
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
		raw := imp.Raw

		// Check 1: must use named import (braces)
		openBrace := strings.Index(raw, "{")
		if openBrace < 0 {
			violations = append(violations, Violation{
				File:    path,
				Line:    imp.Line,
				Col:     1,
				Message: "import-component-style: component import must use named import syntax { }",
			})
			// Without braces, remaining checks can't apply
			goto checkPath
		}

		{
			closeBrace := strings.LastIndex(raw, "}")
			if closeBrace < 0 {
				violations = append(violations, Violation{
					File:    path,
					Line:    imp.Line,
					Col:     1,
					Message: "import-component-style: component import must use named import syntax { }",
				})
				goto checkPath
			}
			bracketBody := raw[openBrace+1 : closeBrace]

			// Check 2: space after {
			if !strings.HasPrefix(bracketBody, " ") {
				violations = append(violations, Violation{
					File:    path,
					Line:    imp.Line,
					Col:     openBrace + 2, // 1-based col of char after {
					Message: "import-component-style: missing space after '{'",
				})
			}

			// Check 3: space before }
			if !strings.HasSuffix(bracketBody, " ") {
				violations = append(violations, Violation{
					File:    path,
					Line:    imp.Line,
					Col:     closeBrace + 1, // 1-based col of }
					Message: "import-component-style: missing space before '}'",
				})
			}

			// Check 4: each comma must have exactly one space after and no space before
			for i := 0; i < len(bracketBody); i++ {
				if bracketBody[i] != ',' {
					continue
				}
				// No space before comma
				if i > 0 && bracketBody[i-1] == ' ' {
					violations = append(violations, Violation{
						File:    path,
						Line:    imp.Line,
						Col:     openBrace + 1 + i, // 1-based col of char before comma
						Message: "import-component-style: unexpected space before ','",
					})
				}
				// Exactly one space after comma
				if i+1 >= len(bracketBody) || bracketBody[i+1] != ' ' {
					violations = append(violations, Violation{
						File:    path,
						Line:    imp.Line,
						Col:     openBrace + 1 + i + 2, // 1-based col of char after comma
						Message: "import-component-style: expected single space after ','",
					})
				} else if i+2 < len(bracketBody) && bracketBody[i+2] == ' ' {
					violations = append(violations, Violation{
						File:    path,
						Line:    imp.Line,
						Col:     openBrace + 1 + i + 2, // 1-based col of extra space
						Message: "import-component-style: expected single space after ','",
					})
				}
			}
		}

	checkPath:
		// Check 5: path ends in .jsx
		if !strings.HasSuffix(imp.Path, ".jsx") {
			violations = append(violations, Violation{
				File:    path,
				Line:    imp.Line,
				Col:     1,
				Message: "import-component-style: component import path must end in '.jsx'",
			})
		}

		// Check 6: raw line (trimmed) ends with ;
		if !strings.HasSuffix(strings.TrimRight(raw, " \t"), ";") {
			violations = append(violations, Violation{
				File:    path,
				Line:    imp.Line,
				Col:     1,
				Message: "import-component-style: import statement must end with ';'",
			})
		}
	}
	return violations
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
		// Check 1: must NOT use named import (no braces)
		if strings.Contains(imp.Raw, "{") {
			violations = append(violations, Violation{
				File:    path,
				Line:    imp.Line,
				Col:     1,
				Message: "import-mdx-style: MDX import must use default import syntax (no braces)",
			})
		}
		// Check 2: raw line (trimmed) ends with ;
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

// ImportBlockContiguousRule checks that no blank lines appear inside the import block.
type ImportBlockContiguousRule struct{}

func (r ImportBlockContiguousRule) CheckFile(path string, lines []string) []Violation {
	// Find first and last import line indices (0-based)
	first := -1
	last := -1
	for i, line := range lines {
		if isImportLine(line) {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	if first < 0 || first == last {
		return nil
	}
	var violations []Violation
	for i := first + 1; i < last; i++ {
		if strings.TrimSpace(lines[i]) == "" {
			violations = append(violations, Violation{
				File:    path,
				Line:    i + 1, // 1-based
				Col:     1,
				Message: "import-block-contiguous: blank line inside import block",
			})
		}
	}
	return violations
}


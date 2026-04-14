package analysis

import "strings"

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

// ParseSingleImport parses one import line. Returns nil if unparseable.
func ParseSingleImport(lineNum int, line string) *ParsedImport {
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

	path, ok := ExtractImportPath(rest)
	if !ok {
		return nil
	}
	return &ParsedImport{Line: lineNum, Raw: line, Names: names, Path: path}
}

// ExtractImportPath finds the last quoted string in rest (the "from '...'" part).
// Returns the unquoted path and true, or ("", false) if not found.
func ExtractImportPath(rest string) (string, bool) {
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

// ParseImports extracts all successfully parsed imports from a slice of lines.
func ParseImports(lines []string) []ParsedImport {
	var imports []ParsedImport
	for i, line := range lines {
		if !IsImportLine(line) {
			continue
		}
		if pi := ParseSingleImport(i+1, line); pi != nil {
			imports = append(imports, *pi)
		}
	}
	return imports
}

// IsImportLine returns true if the line starts with "import ".
func IsImportLine(line string) bool { return strings.HasPrefix(line, "import ") }

// FrontmatterEnd returns the 0-based index of the closing "---" line,
// or -1 if no frontmatter block is present.
func FrontmatterEnd(lines []string) int {
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

// BodyLines returns non-frontmatter, non-import lines for word-search purposes.
func BodyLines(lines []string) []string {
	fmEnd := FrontmatterEnd(lines)
	var body []string
	for i, line := range lines {
		if fmEnd >= 0 && i <= fmEnd {
			continue
		}
		if IsImportLine(line) {
			continue
		}
		body = append(body, line)
	}
	return body
}

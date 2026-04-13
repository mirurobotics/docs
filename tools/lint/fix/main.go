// fix is a one-shot tool that rewrites MDX files to satisfy the lint rules.
// It performs the following transformations in order:
//  1. Fix brace spacing in named component imports: {Foo} → { Foo }
//  2. Add missing semicolons to .mdx and .jsx imports
//  3. Remove blank lines inside import blocks
//  4. Sort imports by path (case-insensitive)
//  5. Remove unused imports
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fix <file>...")
		os.Exit(2)
	}
	for _, path := range os.Args[1:] {
		if err := fixFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "fix: %s: %v\n", path, err)
		}
	}
}

func fixFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(raw), "\n")
	// Remove trailing empty element from final newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	lines = fixBraceSpacing(lines)
	lines = fixSemicolons(lines)
	lines = removeBlankLinesInImportBlock(lines)
	lines = sortImports(lines)
	lines = removeUnusedImports(lines)

	out := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(out), 0o644)
}

// reImportBraces matches: import {content} from 'path';
var reImportBraces = regexp.MustCompile(`^(import\s*)\{([^}]*)\}(\s+from\s+.*)$`)

func fixBraceSpacing(lines []string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		if !strings.HasPrefix(line, "import ") {
			result[i] = line
			continue
		}
		m := reImportBraces.FindStringSubmatchIndex(line)
		if m == nil {
			result[i] = line
			continue
		}
		prefix := line[m[2]:m[3]]
		content := line[m[4]:m[5]]
		suffix := line[m[6]:m[7]]

		// Fix comma spacing: ensure "A, B" not "A,B" or "A , B"
		parts := strings.Split(content, ",")
		for j, p := range parts {
			parts[j] = strings.TrimSpace(p)
		}
		fixedContent := strings.Join(parts, ", ")

		result[i] = prefix + "{ " + fixedContent + " }" + suffix
	}
	return result
}

func fixSemicolons(lines []string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		if !strings.HasPrefix(line, "import ") {
			result[i] = line
			continue
		}
		trimmed := strings.TrimRight(line, " \t")
		if strings.HasSuffix(trimmed, ";") {
			result[i] = line
			continue
		}
		result[i] = trimmed + ";"
	}
	return result
}

func removeBlankLinesInImportBlock(lines []string) []string {
	// Find first and last import line indices
	first, last := -1, -1
	for i, line := range lines {
		if strings.HasPrefix(line, "import ") {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	if first < 0 || first == last {
		return lines
	}
	result := make([]string, 0, len(lines))
	for i, line := range lines {
		if i > first && i < last && strings.TrimSpace(line) == "" {
			continue
		}
		result = append(result, line)
	}
	return result
}

type importLine struct {
	raw  string
	path string
}

func extractPath(line string) string {
	// Find last quoted string
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] == '\'' || line[i] == '"' {
			q := line[i]
			for j := i - 1; j >= 0; j-- {
				if line[j] == q {
					p := line[j+1 : i]
					return strings.TrimSuffix(p, ";")
				}
			}
		}
	}
	return line
}

func sortImports(lines []string) []string {
	// Find the import block indices
	first, last := -1, -1
	for i, line := range lines {
		if strings.HasPrefix(line, "import ") {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	if first < 0 || first == last {
		return lines
	}

	// Collect import lines in their positions
	var imports []importLine
	importIndices := map[int]bool{}
	for i := first; i <= last; i++ {
		if strings.HasPrefix(lines[i], "import ") {
			imports = append(imports, importLine{raw: lines[i], path: extractPath(lines[i])})
			importIndices[i] = true
		}
	}

	sort.SliceStable(imports, func(a, b int) bool {
		return strings.ToLower(imports[a].path) < strings.ToLower(imports[b].path)
	})

	result := make([]string, len(lines))
	copy(result, lines)
	importIdx := 0
	for i := first; i <= last; i++ {
		if importIndices[i] {
			result[i] = imports[importIdx].raw
			importIdx++
		}
	}
	return result
}

func removeUnusedImports(lines []string) []string {
	imports := parseImports(lines)
	if len(imports) == 0 {
		return lines
	}

	// Build body (non-frontmatter, non-import lines)
	fmEnd := frontmatterEnd(lines)
	var bodyParts []string
	for i, line := range lines {
		if fmEnd >= 0 && i <= fmEnd {
			continue
		}
		if strings.HasPrefix(line, "import ") {
			continue
		}
		bodyParts = append(bodyParts, line)
	}
	body := strings.Join(bodyParts, "\n")

	// lineReplace maps 0-based line index to replacement (empty string = delete)
	lineReplace := map[int]string{}

	for _, imp := range imports {
		var usedNames []importedName
		for _, n := range imp.Names {
			re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(n.Name)))
			if re.MatchString(body) {
				usedNames = append(usedNames, n)
			}
		}

		lineIdx := imp.Line - 1 // 0-based

		if len(usedNames) == 0 {
			// Remove the entire line
			lineReplace[lineIdx] = ""
		} else if len(usedNames) < len(imp.Names) {
			// Named import with some unused names — rebuild the import line
			// keeping only used names
			// Only applies to named imports
			if !usedNames[0].IsNamed {
				// Default import is either used or not, should have been caught above
				continue
			}
			nameStrs := make([]string, len(usedNames))
			for i, n := range usedNames {
				nameStrs[i] = n.Name
			}
			// Reconstruct: import { A, B } from 'path';
			// Find the 'from' keyword position
			raw := imp.Raw
			fromIdx := strings.Index(raw, "} from ")
			if fromIdx < 0 {
				continue
			}
			suffix := raw[fromIdx+1:] // "from 'path';"
			lineReplace[lineIdx] = "import { " + strings.Join(nameStrs, ", ") + " " + suffix
		}
	}

	if len(lineReplace) == 0 {
		return lines
	}

	result := make([]string, 0, len(lines))
	for i, line := range lines {
		if replacement, ok := lineReplace[i]; ok {
			if replacement != "" {
				result = append(result, replacement)
			}
			// else skip (delete)
			continue
		}
		result = append(result, line)
	}
	return result
}

// --- minimal copies of the parser helpers from the lint package ---

type importedName struct {
	Name    string
	IsNamed bool
}

type parsedImport struct {
	Line  int
	Raw   string
	Names []importedName
	Path  string
}

func parseSingleImport(lineNum int, line string) *parsedImport {
	if !strings.HasPrefix(line, "import ") {
		return nil
	}
	rest := strings.TrimPrefix(line, "import ")
	rest = strings.TrimLeft(rest, " \t")

	var names []importedName

	if strings.HasPrefix(rest, "{") {
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
			names = append(names, importedName{Name: name, IsNamed: true})
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
		names = append(names, importedName{Name: token, IsNamed: false})
		rest = rest[fromIdx:]
	}

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
	path = strings.TrimSuffix(path, ";")

	return &parsedImport{
		Line:  lineNum,
		Raw:   line,
		Names: names,
		Path:  path,
	}
}

func parseImports(lines []string) []parsedImport {
	var imports []parsedImport
	for i, line := range lines {
		if !strings.HasPrefix(line, "import ") {
			continue
		}
		pi := parseSingleImport(i+1, line)
		if pi != nil {
			imports = append(imports, *pi)
		}
	}
	return imports
}

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

// Ensure filepath is used
var _ = filepath.Join

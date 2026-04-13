package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lint <file>...")
		os.Exit(2)
	}

	rules := []Rule{NoDoubleDash{}}

	contentRoot, err := findContentRoot(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint: cannot determine content root: %v\n", err)
		os.Exit(2)
	}

	fileRules := []FileRule{
		ImportResolvesRule{ContentRoot: contentRoot},
		ImportUsedRule{},
		ImportSortedRule{},
		ComponentImportStyleRule{},
		MDXImportStyleRule{},
		ImportBlockContiguousRule{},
	}

	var allViolations []Violation
	exitCode := 0

	for _, path := range os.Args[1:] {
		violations, err := lintFile(path, rules, fileRules)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lint: %s: %v\n", path, err)
			exitCode = 2
			continue
		}
		allViolations = append(allViolations, violations...)
	}

	for _, v := range allViolations {
		fmt.Printf("%s:%d:%d: %s\n", v.File, v.Line, v.Col, v.Message)
	}

	if len(allViolations) > 0 && exitCode == 0 {
		exitCode = 1
	}
	os.Exit(exitCode)
}

func lintFile(path string, rules []Rule, fileRules []FileRule) ([]Violation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var lines []string
	lineScanner := bufio.NewScanner(f)
	for lineScanner.Scan() {
		lines = append(lines, lineScanner.Text())
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}

	scanner := NewScanner()
	var violations []Violation
	for _, line := range lines {
		spans := scanner.ScanLine(line)
		lineNum := scanner.LineNum()
		for _, rule := range rules {
			violations = append(violations, rule.Check(path, lineNum, spans)...)
		}
	}

	for _, fr := range fileRules {
		violations = append(violations, fr.CheckFile(path, lines)...)
	}

	return violations, nil
}

// findContentRoot walks up the directory tree from startPath until it finds
// a directory that contains a "snippets/" subdirectory.
func findContentRoot(startPath string) (string, error) {
	dir := filepath.Dir(filepath.Clean(startPath))
	for {
		info, err := os.Stat(filepath.Join(dir, "snippets"))
		if err == nil && info.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no snippets/ directory found above %s", startPath)
		}
		dir = parent
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mirurobotics/docs/tools/lint/linter"
	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lint <file>...")
		os.Exit(2)
	}

	contentRoot, err := findContentRoot(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint: cannot determine content root: %v\n", err)
		os.Exit(2)
	}

	var allViolations []analysis.Violation
	exitCode := 0

	for _, path := range os.Args[1:] {
		violations, err := linter.ProcessFile(path, contentRoot)
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

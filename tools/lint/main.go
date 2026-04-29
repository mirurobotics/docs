package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mirurobotics/docs/tools/lint/linter"
	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
	"github.com/mirurobotics/docs/tools/lint/linter/redirects"
)

func main() { os.Exit(run(os.Args, os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		_, _ = fmt.Fprintln(stderr, "usage: lint <file>...")
		return 2
	}

	contentRoot, err := findContentRoot(args[1])
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "lint: cannot determine content root: %v\n", err)
		return 2
	}

	var allViolations []analysis.Violation
	exitCode := 0

	for _, path := range args[1:] {
		violations, err := linter.ProcessFile(path, contentRoot)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "lint: %s: %v\n", path, err)
			exitCode = 2
			continue
		}
		allViolations = append(allViolations, violations...)
	}

	// Redirect violations are appended after the per-file pass so the
	// per-file output order is preserved.
	allViolations = append(allViolations, redirects.Check(contentRoot)...)

	for _, v := range allViolations {
		_, _ = fmt.Fprintf(stdout, "%s:%d:%d: %s\n", v.File, v.Line, v.Col, v.Message)
	}

	if len(allViolations) > 0 && exitCode == 0 {
		exitCode = 1
	}
	return exitCode
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

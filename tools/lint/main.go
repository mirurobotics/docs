package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lint <file>...")
		os.Exit(2)
	}

	rules := []Rule{NoDoubleDash{}}

	var allViolations []Violation
	exitCode := 0

	for _, path := range os.Args[1:] {
		violations, err := lintFile(path, rules)
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

func lintFile(path string, rules []Rule) ([]Violation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	scanner := NewScanner()
	lineScanner := bufio.NewScanner(f)
	var violations []Violation

	for lineScanner.Scan() {
		spans := scanner.ScanLine(lineScanner.Text())
		lineNum := scanner.LineNum()
		for _, rule := range rules {
			violations = append(violations, rule.Check(path, lineNum, spans)...)
		}
	}

	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	return violations, nil
}

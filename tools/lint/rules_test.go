package main

import (
	"strings"
	"testing"
)

func TestNoDoubleDash(t *testing.T) {
	rule := NoDoubleDash{}

	tests := []struct {
		name       string
		input      string
		wantCount  int
		wantCols   []int // 1-based columns of expected violations
	}{
		{
			name:      "double dash in prose",
			input:     "config type--it determines",
			wantCount: 1,
			wantCols:  []int{12},
		},
		{
			name:      "triple dash not flagged",
			input:     "use ---",
			wantCount: 0,
		},
		{
			name:      "quadruple dash not flagged",
			input:     "use ----",
			wantCount: 0,
		},
		{
			name:      "single dash not flagged",
			input:     "a-b",
			wantCount: 0,
		},
		{
			name:      "two separate double dashes",
			input:     "a--b and c--d",
			wantCount: 2,
			wantCols:  []int{2, 11},
		},
		{
			name:      "double dash at start of line",
			input:     "--start",
			wantCount: 1,
			wantCols:  []int{1},
		},
		{
			name:      "double dash at end of line",
			input:     "end--",
			wantCount: 1,
			wantCols:  []int{4},
		},
		{
			name:      "just double dash",
			input:     "--",
			wantCount: 1,
			wantCols:  []int{1},
		},
		{
			name:      "no dashes",
			input:     "hello world",
			wantCount: 0,
		},
		{
			name:      "empty string",
			input:     "",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spans := []ProseSpan{{StartCol: 1, Text: tt.input}}
			if tt.input == "" {
				spans = nil
			}
			violations := rule.Check("test.mdx", 1, spans)
			if len(violations) != tt.wantCount {
				t.Errorf("expected %d violations, got %d: %v", tt.wantCount, len(violations), violations)
				return
			}
			for i, wantCol := range tt.wantCols {
				if violations[i].Col != wantCol {
					t.Errorf("violation %d: expected col %d, got %d", i, wantCol, violations[i].Col)
				}
			}
		})
	}
}

func TestNoDoubleDashWithOffset(t *testing.T) {
	rule := NoDoubleDash{}

	// Simulate a span that starts at column 10 (e.g. after inline code)
	spans := []ProseSpan{{StartCol: 10, Text: "a--b"}}
	violations := rule.Check("test.mdx", 1, spans)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Col != 11 {
		t.Errorf("expected col 11, got %d", violations[0].Col)
	}
}

func TestNoDoubleDashIntegration(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
		wantLines []int
	}{
		{
			name: "double dash in prose after frontmatter",
			content: "---\ntitle: Test\n---\n\nconfig type--it determines",
			wantCount: 1,
			wantLines: []int{5},
		},
		{
			name: "double dash in code block not flagged",
			content: "---\ntitle: Test\n---\n\n```bash\nmiru --version\n```",
			wantCount: 0,
		},
		{
			name: "double dash in inline code not flagged",
			content: "---\ntitle: Test\n---\n\nUse `--version` flag",
			wantCount: 0,
		},
		{
			name: "double dash in frontmatter not flagged",
			content: "---\ntitle: Test--Title\n---\n\nClean prose",
			wantCount: 0,
		},
		{
			name: "double dash in JSX attribute not flagged",
			content: "---\ntitle: Test\n---\n\n<ParamField path=\"--version\" type=\"string\">",
			wantCount: 0,
		},
		{
			name: "frontmatter dashes not flagged",
			content: "---\ntitle: Test\n---\n\nClean",
			wantCount: 0,
		},
		{
			name: "thematic break not flagged",
			content: "---\ntitle: Test\n---\n\n---\n\nClean",
			wantCount: 0,
		},
		{
			name: "table separator not flagged",
			content: "---\ntitle: Test\n---\n\n| A | B |\n|---|---|\n| 1 | 2 |",
			wantCount: 0,
		},
		{
			name: "HTML comment not flagged",
			content: "---\ntitle: Test\n---\n\n<!-- --test -->",
			wantCount: 0,
		},
		{
			name: "multiline HTML comment not flagged",
			content: "---\ntitle: Test\n---\n\n<!-- \n--test\n-->",
			wantCount: 0,
		},
		{
			name: "import not flagged",
			content: "---\ntitle: Test\n---\n\nimport Foo from '--bar'",
			wantCount: 0,
		},
		{
			name: "mixed line with inline code and prose dash",
			content: "---\ntitle: Test\n---\n\nUse `--flag` for type--detection",
			wantCount: 1,
			wantLines: []int{5},
		},
	}

	rule := NoDoubleDash{}
	rules := []Rule{rule}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner()
			var violations []Violation
			for i, line := range strings.Split(tt.content, "\n") {
				spans := scanner.ScanLine(line)
				for _, r := range rules {
					violations = append(violations, r.Check("test.mdx", i+1, spans)...)
				}
			}
			if len(violations) != tt.wantCount {
				t.Errorf("expected %d violations, got %d: %v", tt.wantCount, len(violations), violations)
				return
			}
			for i, wantLine := range tt.wantLines {
				if violations[i].Line != wantLine {
					t.Errorf("violation %d: expected line %d, got %d", i, wantLine, violations[i].Line)
				}
			}
		})
	}
}

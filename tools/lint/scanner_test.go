package main

import (
	"testing"
)

func TestFrontmatter(t *testing.T) {
	s := NewScanner()

	if spans := s.ScanLine("---"); spans != nil {
		t.Errorf("frontmatter open: expected nil, got %v", spans)
	}
	if spans := s.ScanLine(`title: "Test"`); spans != nil {
		t.Errorf("frontmatter body: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("---"); spans != nil {
		t.Errorf("frontmatter close: expected nil, got %v", spans)
	}

	spans := s.ScanLine("Hello world")
	if len(spans) != 1 || spans[0].Text != "Hello world" {
		t.Errorf("after frontmatter: expected prose, got %v", spans)
	}
}

func TestFencedCodeBlock(t *testing.T) {
	s := NewScanner()

	if spans := s.ScanLine("```bash"); spans != nil {
		t.Errorf("code fence open: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("miru release create --version v1.0.0"); spans != nil {
		t.Errorf("code block body: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("```"); spans != nil {
		t.Errorf("code fence close: expected nil, got %v", spans)
	}

	spans := s.ScanLine("After code")
	if len(spans) != 1 || spans[0].Text != "After code" {
		t.Errorf("after code block: expected prose, got %v", spans)
	}
}

func TestTildeFence(t *testing.T) {
	s := NewScanner()

	if spans := s.ScanLine("~~~"); spans != nil {
		t.Errorf("tilde fence open: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("--flag"); spans != nil {
		t.Errorf("tilde block body: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("~~~"); spans != nil {
		t.Errorf("tilde fence close: expected nil, got %v", spans)
	}
}

func TestInlineCode(t *testing.T) {
	s := NewScanner()
	// Skip frontmatter
	s.ScanLine("---")
	s.ScanLine("---")

	spans := s.ScanLine("Use the `--version` flag here")
	// Should get: "Use the " and " flag here"
	if len(spans) != 2 {
		t.Fatalf("inline code: expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0].Text != "Use the " {
		t.Errorf("span 0: expected 'Use the ', got %q", spans[0].Text)
	}
	if spans[0].StartCol != 1 {
		t.Errorf("span 0 col: expected 1, got %d", spans[0].StartCol)
	}
	if spans[1].Text != " flag here" {
		t.Errorf("span 1: expected ' flag here', got %q", spans[1].Text)
	}
}

func TestDoubleBacktickInlineCode(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	spans := s.ScanLine("Use ``--flag`` here")
	if len(spans) != 2 {
		t.Fatalf("double backtick: expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0].Text != "Use " {
		t.Errorf("span 0: expected 'Use ', got %q", spans[0].Text)
	}
	if spans[1].Text != " here" {
		t.Errorf("span 1: expected ' here', got %q", spans[1].Text)
	}
}

func TestHTMLComment(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	spans := s.ScanLine("before <!-- --test --> after")
	if len(spans) != 2 {
		t.Fatalf("html comment: expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0].Text != "before " {
		t.Errorf("span 0: expected 'before ', got %q", spans[0].Text)
	}
	if spans[1].Text != " after" {
		t.Errorf("span 1: expected ' after', got %q", spans[1].Text)
	}
}

func TestMultilineHTMLComment(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	spans := s.ScanLine("before <!-- start")
	if len(spans) != 1 || spans[0].Text != "before " {
		t.Errorf("comment open line: expected 'before ', got %v", spans)
	}

	if spans := s.ScanLine("-- middle"); spans != nil {
		t.Errorf("comment body: expected nil, got %v", spans)
	}

	spans = s.ScanLine("--> after")
	if len(spans) != 1 || spans[0].Text != " after" {
		t.Errorf("comment close line: expected ' after', got %v", spans)
	}
}

func TestJSXTag(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	spans := s.ScanLine(`before <ParamField path="--version" type="string"> after`)
	if len(spans) != 2 {
		t.Fatalf("jsx tag: expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0].Text != "before " {
		t.Errorf("span 0: expected 'before ', got %q", spans[0].Text)
	}
	if spans[1].Text != " after" {
		t.Errorf("span 1: expected ' after', got %q", spans[1].Text)
	}
}

func TestThematicBreak(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	// Thematic break (not frontmatter since we're past it)
	if spans := s.ScanLine("---"); spans != nil {
		t.Errorf("thematic break: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("----"); spans != nil {
		t.Errorf("long thematic break: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("  ---"); spans != nil {
		t.Errorf("indented thematic break: expected nil, got %v", spans)
	}
}

func TestTableSeparator(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	if spans := s.ScanLine("|---|---|"); spans != nil {
		t.Errorf("table separator: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("| --- | --- |"); spans != nil {
		t.Errorf("spaced table separator: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("|:---:|:---:|"); spans != nil {
		t.Errorf("aligned table separator: expected nil, got %v", spans)
	}
}

func TestImportExport(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	if spans := s.ScanLine("import Foo from './foo'"); spans != nil {
		t.Errorf("import: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("export default Bar"); spans != nil {
		t.Errorf("export: expected nil, got %v", spans)
	}
}

func TestIndentedCodeFence(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	// MDX nests code blocks inside JSX at arbitrary indentation
	if spans := s.ScanLine("    ```bash Empty Schemas"); spans != nil {
		t.Errorf("indented fence open: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("      --version v1.0.0"); spans != nil {
		t.Errorf("indented fence body: expected nil, got %v", spans)
	}
	if spans := s.ScanLine("    ```"); spans != nil {
		t.Errorf("indented fence close: expected nil, got %v", spans)
	}

	spans := s.ScanLine("After indented code")
	if len(spans) != 1 || spans[0].Text != "After indented code" {
		t.Errorf("after indented code: expected prose, got %v", spans)
	}
}

func TestUnclosedQuoteInTag(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	// Must not panic on malformed tags with unclosed quotes
	spans := s.ScanLine(`<tag attr="unclosed`)
	if spans != nil {
		t.Errorf("unclosed double quote tag: expected nil, got %v", spans)
	}

	spans = s.ScanLine(`<tag attr='unclosed`)
	if spans != nil {
		t.Errorf("unclosed single quote tag: expected nil, got %v", spans)
	}
}

func TestLineNum(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")
	s.ScanLine("line three")

	if s.LineNum() != 3 {
		t.Errorf("expected LineNum 3, got %d", s.LineNum())
	}
}

func TestProseLineColumns(t *testing.T) {
	s := NewScanner()
	s.ScanLine("---")
	s.ScanLine("---")

	spans := s.ScanLine("plain prose line")
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].StartCol != 1 {
		t.Errorf("expected col 1, got %d", spans[0].StartCol)
	}
	if spans[0].Text != "plain prose line" {
		t.Errorf("expected 'plain prose line', got %q", spans[0].Text)
	}
}

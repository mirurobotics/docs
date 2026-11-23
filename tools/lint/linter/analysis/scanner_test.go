package analysis

import (
	"testing"
)

func TestFrontmatter(t *testing.T) {
	scanner := NewScanner()

	if spans := scanner.ScanLine("---"); spans != nil {
		t.Errorf("frontmatter open: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine(`title: "Test"`); spans != nil {
		t.Errorf("frontmatter body: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("---"); spans != nil {
		t.Errorf("frontmatter close: expected nil, got %v", spans)
	}

	spans := scanner.ScanLine("Hello world")
	if len(spans) != 1 || spans[0].Text != "Hello world" {
		t.Errorf("after frontmatter: expected prose, got %v", spans)
	}
}

func TestFencedCodeBlock(t *testing.T) {
	scanner := NewScanner()

	if spans := scanner.ScanLine("```bash"); spans != nil {
		t.Errorf("code fence open: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("miru release create --version v1.0.0"); spans != nil {
		t.Errorf("code block body: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("```"); spans != nil {
		t.Errorf("code fence close: expected nil, got %v", spans)
	}

	spans := scanner.ScanLine("After code")
	if len(spans) != 1 || spans[0].Text != "After code" {
		t.Errorf("after code block: expected prose, got %v", spans)
	}
}

func TestTildeFence(t *testing.T) {
	scanner := NewScanner()

	if spans := scanner.ScanLine("~~~"); spans != nil {
		t.Errorf("tilde fence open: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("--flag"); spans != nil {
		t.Errorf("tilde block body: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("~~~"); spans != nil {
		t.Errorf("tilde fence close: expected nil, got %v", spans)
	}
}

func TestFourBacktickFence(t *testing.T) {
	scanner := NewScanner()

	if spans := scanner.ScanLine("````mdx"); spans != nil {
		t.Errorf("4-backtick fence open: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("--flag"); spans != nil {
		t.Errorf("4-backtick block body: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("````"); spans != nil {
		t.Errorf("4-backtick fence close: expected nil, got %v", spans)
	}

	spans := scanner.ScanLine("After 4-backtick block")
	if len(spans) != 1 || spans[0].Text != "After 4-backtick block" {
		t.Errorf("after 4-backtick block: expected prose, got %v", spans)
	}
}

func TestFourBacktickFenceInnerTriple(t *testing.T) {
	scanner := NewScanner()

	if spans := scanner.ScanLine("````"); spans != nil {
		t.Errorf("4-backtick open: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("```"); spans != nil {
		t.Errorf("inner triple backtick: expected nil (body), got %v", spans)
	}
	if spans := scanner.ScanLine("````"); spans != nil {
		t.Errorf("4-backtick close: expected nil, got %v", spans)
	}

	spans := scanner.ScanLine("prose after")
	if len(spans) != 1 || spans[0].Text != "prose after" {
		t.Errorf("after block: expected prose, got %v", spans)
	}
}

func TestThematicBreak(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	if spans := scanner.ScanLine("---"); spans != nil {
		t.Errorf("thematic break: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("----"); spans != nil {
		t.Errorf("long thematic break: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("  ---"); spans != nil {
		t.Errorf("indented thematic break: expected nil, got %v", spans)
	}
}

func TestTableSeparator(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	if spans := scanner.ScanLine("|---|---|"); spans != nil {
		t.Errorf("table separator: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("| --- | --- |"); spans != nil {
		t.Errorf("spaced table separator: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("|:---:|:---:|"); spans != nil {
		t.Errorf("aligned table separator: expected nil, got %v", spans)
	}
}

func TestImportExport(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	if spans := scanner.ScanLine("import Foo from './foo'"); spans != nil {
		t.Errorf("import: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("export default Bar"); spans != nil {
		t.Errorf("export: expected nil, got %v", spans)
	}
}

func TestIndentedCodeFence(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	if spans := scanner.ScanLine("    ```bash Empty Schemas"); spans != nil {
		t.Errorf("indented fence open: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("      --version v1.0.0"); spans != nil {
		t.Errorf("indented fence body: expected nil, got %v", spans)
	}
	if spans := scanner.ScanLine("    ```"); spans != nil {
		t.Errorf("indented fence close: expected nil, got %v", spans)
	}

	spans := scanner.ScanLine("After indented code")
	if len(spans) != 1 || spans[0].Text != "After indented code" {
		t.Errorf("after indented code: expected prose, got %v", spans)
	}
}

func TestLineNum(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")
	scanner.ScanLine("line three")

	if scanner.LineNum() != 3 {
		t.Errorf("expected LineNum 3, got %d", scanner.LineNum())
	}
}

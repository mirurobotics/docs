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

func TestInlineCode(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("Use the `--version` flag here")
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
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("Use ``--flag`` here")
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
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("before <!-- --test --> after")
	if len(spans) != 2 {
		t.Fatalf("html comment: expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0] != (ProseSpan{StartCol: 1, Text: "before "}) {
		t.Errorf("span 0: expected 'before ' at col 1, got %+v", spans[0])
	}
	if spans[1] != (ProseSpan{StartCol: 23, Text: " after"}) {
		t.Errorf("span 1: expected ' after' at col 23, got %+v", spans[1])
	}
}

func TestMultilineHTMLComment(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("before <!-- start")
	if len(spans) != 1 || spans[0].Text != "before " {
		t.Errorf("comment open line: expected 'before ', got %v", spans)
	}

	if spans := scanner.ScanLine("-- middle"); spans != nil {
		t.Errorf("comment body: expected nil, got %v", spans)
	}

	spans = scanner.ScanLine("--> after")
	if len(spans) != 1 || spans[0] != (ProseSpan{StartCol: 4, Text: " after"}) {
		t.Errorf("comment close line: expected ' after' at col 4, got %+v", spans)
	}
}

func TestJSXTag(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine(
		`before <ParamField path="--version" type="string"> after`,
	)
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

func TestUnclosedQuoteInTag(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine(`<tag attr="unclosed`)
	if spans != nil {
		t.Errorf("unclosed double quote tag: expected nil, got %v", spans)
	}

	spans = scanner.ScanLine(`<tag attr='unclosed`)
	if spans != nil {
		t.Errorf("unclosed single quote tag: expected nil, got %v", spans)
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

func TestProseLineColumns(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("plain prose line")
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

// TestWhitespaceOnlyBeforeTag verifies a whitespace-only span before a tag is dropped.
func TestWhitespaceOnlyBeforeTag(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("   <Tag> prose text")
	if len(spans) != 1 {
		t.Fatalf("expected 1 span (spaces dropped), got %d: %v", len(spans), spans)
	}
	if spans[0].Text != " prose text" {
		t.Errorf("expected ' prose text', got %q", spans[0].Text)
	}
}

// TestUnclosedBacktick verifies an unmatched backtick is literal and the rest is prose.
func TestUnclosedBacktick(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("before `unclosed")
	if len(spans) != 2 {
		t.Fatalf("expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0].Text != "before " {
		t.Errorf("span 0: expected 'before ', got %q", spans[0].Text)
	}
	if spans[1].Text != "unclosed" {
		t.Errorf("span 1: expected 'unclosed', got %q", spans[1].Text)
	}
}

// TestMismatchedBacktickCount verifies one backtick cannot close double-backtick code.
func TestMismatchedBacktickCount(t *testing.T) {
	scanner := NewScanner()
	scanner.ScanLine("---")
	scanner.ScanLine("---")

	spans := scanner.ScanLine("before ``code`more`` after")
	if len(spans) != 2 {
		t.Fatalf("expected 2 spans, got %d: %v", len(spans), spans)
	}
	if spans[0].Text != "before " {
		t.Errorf("span 0: expected 'before ', got %q", spans[0].Text)
	}
	if spans[1].Text != " after" {
		t.Errorf("span 1: expected ' after', got %q", spans[1].Text)
	}
}

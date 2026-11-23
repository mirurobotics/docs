package analysis

import (
	"testing"
)

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

	scanner.ScanLine("<!-- again")
	spans = scanner.ScanLine("--> a  b <!-- y")
	if len(spans) != 1 || spans[0] != (ProseSpan{StartCol: 4, Text: " a  b "}) {
		t.Errorf("close and reopen line: expected ' a  b ' at col 4, got %+v", spans)
	}
	if spans := scanner.ScanLine("-- z"); spans != nil {
		t.Errorf("reopened comment body: expected nil, got %v", spans)
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

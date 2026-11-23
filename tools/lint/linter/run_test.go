package linter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func TestProcessFile_Reading(t *testing.T) {
	t.Run("missing file returns error", func(t *testing.T) {
		_, err := ProcessFile("/no/such/file.mdx", "/tmp")
		if err == nil {
			t.Error("expected error for missing file")
		}
	})

	t.Run("no violations", func(t *testing.T) {
		content := "# Title\n\nHello world.\n"
		violations := processContent(t, content)
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(violations), violations)
		}
	})

	t.Run("very long line does not error", func(t *testing.T) {
		longLine := strings.Repeat("x", 200*1024)
		content := "# Title\n\n" + longLine + "\n"
		processContent(t, content)
	})
}

func TestProcessFile_ProseAndFileRules(t *testing.T) {
	t.Run("prose rule violation", func(t *testing.T) {
		content := "# Title\n\nBad -- prose.\n"
		violations := processContent(t, content)
		found := false
		for _, got := range violations {
			if got.Message == "no-double-dash: use em dash '\u2014' instead of '--'" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected no-double-dash violation, got %v", violations)
		}
	})

	t.Run("file rule violation", func(t *testing.T) {
		content := "import Missing from '/snippets/missing.mdx';\n\nHello.\n"
		violations := processContent(t, content)
		found := false
		for _, got := range violations {
			if got.Line == 1 {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation on line 1, got %v", violations)
		}
	})
}

func TestProcessFile_LineRules(t *testing.T) {
	t.Run("image-domain rule violation", func(t *testing.T) {
		content := "# Title\n\n![diagram](/images/x.png)\n"
		violations := processContent(t, content)
		found := false
		for _, got := range violations {
			if got.Line == 3 && strings.HasPrefix(got.Message, "image-domain:") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected image-domain violation on line 3, got %v", violations)
		}
	})

	t.Run("shell-pair rule violation", func(t *testing.T) {
		content := "# Title\n\n<CodeGroup>\n" +
			"```bash Bash\nmiru a \\\n  --b\n```\n</CodeGroup>\n"
		violations := processContent(t, content)
		found := false
		for _, got := range violations {
			if got.Line == 4 && strings.HasPrefix(got.Message, "shell-pair:") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected shell-pair violation on line 4, got %v", violations)
		}
	})
}

func processContent(t *testing.T, content string) []analysis.Violation {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.mdx")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	violations, err := ProcessFile(path, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return violations
}

package linter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessFile(t *testing.T) {
	t.Run("missing file returns error", func(t *testing.T) {
		_, err := ProcessFile("/no/such/file.mdx", "/tmp")
		if err == nil {
			t.Error("expected error for missing file")
		}
	})

	t.Run("no violations", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.mdx")
		content := []byte("# Title\n\nHello world.\n")
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}
		vs, err := ProcessFile(path, dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("prose rule violation", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.mdx")
		content := []byte("# Title\n\nBad -- prose.\n")
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}
		vs, err := ProcessFile(path, dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		found := false
		for _, v := range vs {
			if v.Message == "no-double-dash: use em dash '\u2014' instead of '--'" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected no-double-dash violation, got %v", vs)
		}
	})

	t.Run("very long line does not error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.mdx")
		longLine := strings.Repeat("x", 200*1024)
		content := "# Title\n\n" + longLine + "\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ProcessFile(path, dir); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("image-domain rule violation", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.mdx")
		content := "# Title\n\n![diagram](/images/x.png)\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		vs, err := ProcessFile(path, dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		found := false
		for _, v := range vs {
			if v.Line == 3 && strings.HasPrefix(v.Message, "image-domain:") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected image-domain violation on line 3, got %v", vs)
		}
	})

	t.Run("file rule violation", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.mdx")
		content := "import Missing from '/snippets/missing.mdx';\n\nHello.\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		vs, err := ProcessFile(path, dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		found := false
		for _, v := range vs {
			if v.Line == 1 {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation on line 1, got %v", vs)
		}
	})
}

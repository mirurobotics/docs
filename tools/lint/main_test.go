package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindContentRoot(t *testing.T) {
	t.Run("snippets in same dir", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "snippets"), 0o755); err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(root, "file.mdx")
		if err := os.WriteFile(file, []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := findContentRoot(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != root {
			t.Errorf("expected %q, got %q", root, got)
		}
	})

	t.Run("snippets two levels up", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "snippets"), 0o755); err != nil {
			t.Fatal(err)
		}
		subdir := filepath.Join(root, "docs", "sub")
		if err := os.MkdirAll(subdir, 0o755); err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(subdir, "intro.mdx")
		if err := os.WriteFile(file, []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := findContentRoot(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != root {
			t.Errorf("expected %q, got %q", root, got)
		}
	})

	t.Run("no snippets dir returns error", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "file.mdx")
		if err := os.WriteFile(file, []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := findContentRoot(file)
		if err == nil {
			t.Error("expected error when snippets/ not found")
		}
	})
}

func TestLintFile(t *testing.T) {
	t.Run("missing file returns error", func(t *testing.T) {
		_, err := lintFile("/no/such/file.mdx", nil, nil)
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
		vs, err := lintFile(path, []Rule{NoDoubleDash{}}, nil)
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
		vs, err := lintFile(path, []Rule{NoDoubleDash{}}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
	})

	t.Run("file rule violation", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.mdx")
		content := "import Missing from '/snippets/missing.mdx';\n\nHello.\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		fr := ImportResolvesRule{ContentRoot: dir}
		vs, err := lintFile(path, nil, []FileRule{fr})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
	})
}

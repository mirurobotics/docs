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

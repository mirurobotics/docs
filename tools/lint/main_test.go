package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestRun(t *testing.T) {
	t.Run("no args returns 2", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint"}, &stdout, &stderr)
		if got != 2 {
			t.Errorf("exit code = %d, want 2", got)
		}
		if !strings.Contains(stderr.String(), "usage:") {
			t.Errorf("stderr = %q, want contains 'usage:'", stderr.String())
		}
	})

	t.Run("missing snippets returns 2", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "x.mdx")
		if err := os.WriteFile(file, []byte("# x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", file}, &stdout, &stderr)
		if got != 2 {
			t.Errorf("exit code = %d, want 2", got)
		}
		if !strings.Contains(stderr.String(), "cannot determine content root") {
			t.Errorf(
				"stderr = %q, want contains 'cannot determine content root'",
				stderr.String(),
			)
		}
	})

	t.Run("clean run returns 0", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "snippets"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(root, "docs", "x.mdx")
		if err := os.WriteFile(file, []byte("# x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", file}, &stdout, &stderr)
		if got != 0 {
			t.Errorf("exit code = %d, want 0; stderr=%q", got, stderr.String())
		}
	})

	t.Run("nonexistent file returns 2", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "snippets"), 0o755); err != nil {
			t.Fatal(err)
		}
		good := filepath.Join(root, "good.mdx")
		if err := os.WriteFile(good, []byte("# x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		missing := filepath.Join(root, "missing.mdx")
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", good, missing}, &stdout, &stderr)
		if got != 2 {
			t.Errorf("exit code = %d, want 2", got)
		}
	})

	t.Run("redirect violation returns 1", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "snippets"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(root, "docs", "x.mdx")
		if err := os.WriteFile(file, []byte("# x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		// docs.json with a redirect destined for a missing page
		docsJSON := `{"redirects":[{"source":"/docs/old",` +
			`"destination":"/docs/missing"}]}`
		docsJSONPath := filepath.Join(root, "docs.json")
		if err := os.WriteFile(docsJSONPath, []byte(docsJSON), 0o644); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", file}, &stdout, &stderr)
		if got != 1 {
			t.Errorf(
				"exit code = %d, want 1; stdout=%q stderr=%q",
				got, stdout.String(), stderr.String(),
			)
		}
		if !strings.Contains(stdout.String(), "missing destination") {
			t.Errorf(
				"stdout = %q, want contains 'missing destination'",
				stdout.String(),
			)
		}
	})
}

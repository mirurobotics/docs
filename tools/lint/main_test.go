package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindContentRoot_Found(t *testing.T) {
	t.Run("snippets in same dir", func(t *testing.T) {
		root := setupContentRoot(t)
		file := filepath.Join(root, "file.mdx")
		writeContentFile(t, file, "")
		got, err := findContentRoot(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != root {
			t.Errorf("expected %q, got %q", root, got)
		}
	})

	t.Run("snippets two levels up", func(t *testing.T) {
		root := setupContentRoot(t)
		subdir := filepath.Join(root, "section", "sub")
		if err := os.MkdirAll(subdir, 0o755); err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(subdir, "intro.mdx")
		writeContentFile(t, file, "")
		got, err := findContentRoot(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != root {
			t.Errorf("expected %q, got %q", root, got)
		}
	})
}

func TestFindContentRoot_NotFound(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file.mdx")
	writeContentFile(t, file, "")
	_, err := findContentRoot(file)
	if err == nil {
		t.Error("expected error when snippets/ not found")
	}
}

func TestRun_Usage(t *testing.T) {
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
		writeContentFile(t, file, "# x\n")
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
}

func TestRun_Files(t *testing.T) {
	t.Run("clean run returns 0", func(t *testing.T) {
		file := newContentFile(t, "# Hello\n")
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", file}, &stdout, &stderr)
		if got != 0 {
			t.Errorf("exit code = %d, want 0; stderr=%q", got, stderr.String())
		}
	})

	t.Run("nonexistent file returns 2", func(t *testing.T) {
		root := setupContentRoot(t)
		good := filepath.Join(root, "good.mdx")
		writeContentFile(t, good, "# x\n")
		missing := filepath.Join(root, "missing.mdx")
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", good, missing}, &stdout, &stderr)
		if got != 2 {
			t.Errorf("exit code = %d, want 2", got)
		}
	})
}

func TestRun_RedirectViolation(t *testing.T) {
	root := setupContentRoot(t)
	file := filepath.Join(root, "x.mdx")
	writeContentFile(t, file, "# x\n")
	docsJSON := `{"redirects":[{"source":"/old",` +
		`"destination":"/missing"}]}`
	docsJSONPath := filepath.Join(root, "docs.json")
	writeContentFile(t, docsJSONPath, docsJSON)
	var stdout, stderr bytes.Buffer
	got := run([]string{"lint", file}, &stdout, &stderr)
	if got != 1 {
		t.Errorf(
			"exit code = %d, want 1; stdout=%q stderr=%q",
			got, stdout.String(), stderr.String(),
		)
	}
	if !strings.Contains(stdout.String(), "missing destination") {
		t.Errorf("stdout = %q, want contains 'missing destination'", stdout.String())
	}
}

func TestRun_HeadingCaseViolation(t *testing.T) {
	content := "---\ntitle: \"User Management\"\n---\n\n## Configure deployments\n"
	file := newContentFile(t, content)
	var stdout, stderr bytes.Buffer
	got := run([]string{"lint", file}, &stdout, &stderr)
	if got != 1 {
		t.Errorf(
			"exit code = %d, want 1; stdout=%q stderr=%q",
			got, stdout.String(), stderr.String(),
		)
	}
	if !strings.Contains(stdout.String(), "heading-case:") {
		t.Errorf("stdout = %q, want contains 'heading-case:'", stdout.String())
	}
	if !strings.Contains(stdout.String(), ":2:") {
		t.Errorf("stdout = %q, want contains ':2:'", stdout.String())
	}
	if !strings.Contains(stdout.String(), file) {
		t.Errorf("stdout = %q, want contains %q", stdout.String(), file)
	}
	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRun_CleanHeadings(t *testing.T) {
	t.Run("clean headings return 0", func(t *testing.T) {
		content := "---\ntitle: \"Workspace\"\n---\n\n## Configure deployments\n\n### What is a config?\n"
		file := newContentFile(t, content)
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", file}, &stdout, &stderr)
		if got != 0 {
			t.Errorf(
				"exit code = %d, want 0; stdout=%q stderr=%q",
				got, stdout.String(), stderr.String(),
			)
		}
		if stdout.String() != "" {
			t.Errorf("stdout = %q, want empty", stdout.String())
		}
	})

	t.Run("clean allowlisted acronym title returns 0", func(t *testing.T) {
		content := "---\ntitle: \"API keys\"\n---\n\n## OpenAPI specifications\n"
		file := newContentFile(t, content)
		var stdout, stderr bytes.Buffer
		got := run([]string{"lint", file}, &stdout, &stderr)
		if got != 0 {
			t.Errorf(
				"exit code = %d, want 0; stdout=%q stderr=%q",
				got, stdout.String(), stderr.String(),
			)
		}
		if stdout.String() != "" {
			t.Errorf("stdout = %q, want empty", stdout.String())
		}
	})
}

func newContentFile(t *testing.T, content string) string {
	t.Helper()
	file := filepath.Join(setupContentRoot(t), "x.mdx")
	writeContentFile(t, file, content)
	return file
}

func setupContentRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "snippets"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeContentFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

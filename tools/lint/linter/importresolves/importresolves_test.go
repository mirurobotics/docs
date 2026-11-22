package importresolves

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheck(t *testing.T) {
	root := t.TempDir()
	snippetsDir := filepath.Join(root, "snippets", "components")
	if err := os.MkdirAll(snippetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	realFile := filepath.Join(snippetsDir, "framed.jsx")
	if err := os.WriteFile(realFile, []byte("// framed"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("file exists", func(t *testing.T) {
		lines := []string{"import { Framed } from '/snippets/components/framed.jsx';"}
		violations := Check("test.mdx", lines, root)
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(violations), violations)
		}
	})

	t.Run("file missing", func(t *testing.T) {
		lines := []string{"import { Missing } from '/snippets/components/missing.jsx';"}
		violations := Check("test.mdx", lines, root)
		if len(violations) != 1 {
			t.Errorf("expected 1 violation, got %d", len(violations))
		}
		if len(violations) > 0 && violations[0].Line != 1 {
			t.Errorf("expected violation on line 1, got line %d", violations[0].Line)
		}
	})

	t.Run("relative path skipped", func(t *testing.T) {
		lines := []string{"import Foo from './relative.mdx';"}
		violations := Check("test.mdx", lines, root)
		if len(violations) != 0 {
			t.Errorf("expected 0 violations for relative path, got %d", len(violations))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		violations := Check("test.mdx", []string{"# Hello world"}, root)
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})
}

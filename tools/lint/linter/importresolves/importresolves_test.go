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
		vs := Check("test.mdx", lines, root)
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("file missing", func(t *testing.T) {
		lines := []string{"import { Missing } from '/snippets/components/missing.jsx';"}
		vs := Check("test.mdx", lines, root)
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d", len(vs))
		}
		if len(vs) > 0 && vs[0].Line != 1 {
			t.Errorf("expected violation on line 1, got line %d", vs[0].Line)
		}
	})

	t.Run("relative path skipped", func(t *testing.T) {
		lines := []string{"import Foo from './relative.mdx';"}
		vs := Check("test.mdx", lines, root)
		if len(vs) != 0 {
			t.Errorf("expected 0 violations for relative path, got %d", len(vs))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := Check("test.mdx", []string{"# Hello world"}, root)
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})
}

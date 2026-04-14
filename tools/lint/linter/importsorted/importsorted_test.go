package importsorted

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Run("single import", func(t *testing.T) {
		vs := Check("test.mdx", []string{"import A from '/snippets/a.mdx';"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("two imports in order", func(t *testing.T) {
		content := "import A from '/snippets/a.mdx';\nimport B from '/snippets/b.mdx';"
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("two imports out of order", func(t *testing.T) {
		content := "import B from '/snippets/b.mdx';\nimport A from '/snippets/a.mdx';"
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d", len(vs))
		}
		if len(vs) > 0 && vs[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", vs[0].Line)
		}
	})

	t.Run("three imports first pair out of order", func(t *testing.T) {
		content := strings.Join([]string{
			"import B from '/snippets/b.mdx';",
			"import A from '/snippets/a.mdx';",
			"import C from '/snippets/c.mdx';",
		}, "\n")
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation (only first out-of-order), got %d", len(vs))
		}
		if len(vs) > 0 && vs[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", vs[0].Line)
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := Check("test.mdx", []string{"# Hello"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})
}

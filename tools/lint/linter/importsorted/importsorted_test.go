package importsorted

import (
	"strings"
	"testing"
)

func TestCheck_NoViolations(t *testing.T) {
	t.Run("single import", func(t *testing.T) {
		violations := Check("test.mdx", []string{"import A from '/snippets/a.mdx';"})
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})

	t.Run("two imports in order", func(t *testing.T) {
		content := "import A from '/snippets/a.mdx';\nimport B from '/snippets/b.mdx';"
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		violations := Check("test.mdx", []string{"# Hello"})
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})
}

func TestCheck_OutOfOrder(t *testing.T) {
	t.Run("two imports out of order", func(t *testing.T) {
		content := "import B from '/snippets/b.mdx';\nimport A from '/snippets/a.mdx';"
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 1 {
			t.Errorf("expected 1 violation, got %d", len(violations))
		}
		if len(violations) > 0 && violations[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", violations[0].Line)
		}
	})

	t.Run("three imports first pair out of order", func(t *testing.T) {
		content := strings.Join([]string{
			"import B from '/snippets/b.mdx';",
			"import A from '/snippets/a.mdx';",
			"import C from '/snippets/c.mdx';",
		}, "\n")
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 1 {
			t.Errorf(
				"expected 1 violation (only first out-of-order), got %d",
				len(violations),
			)
		}
		if len(violations) > 0 && violations[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", violations[0].Line)
		}
	})
}

package importblock

import (
	"strings"
	"testing"
)

func TestCheck_NoBlankLines(t *testing.T) {
	t.Run("no blank lines between imports", func(t *testing.T) {
		content := "import A from '/snippets/a.mdx';\nimport B from '/snippets/b.mdx';"
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})

	t.Run("non-blank content between imports no violation", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"## Heading",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 0 {
			t.Errorf(
				"expected 0 violations with non-blank content between imports, got %d",
				len(violations),
			)
		}
	})
}

func TestCheck_BlankInsideBlock(t *testing.T) {
	t.Run("one blank line between first and second import", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(violations), violations)
		}
		if len(violations) > 0 && violations[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", violations[0].Line)
		}
	})

	t.Run("two blank lines in import block", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"",
			"",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 2 {
			t.Errorf("expected 2 violations, got %d: %v", len(violations), violations)
		}
	})
}

func TestCheck_BlankOutsideBlock(t *testing.T) {
	t.Run("blank line before first import", func(t *testing.T) {
		content := strings.Join([]string{
			"",
			"import A from '/snippets/a.mdx';",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 0 {
			t.Errorf(
				"expected 0 violations (blank before first), got %d",
				len(violations),
			)
		}
	})

	t.Run("blank line after last import", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"import B from '/snippets/b.mdx';",
			"",
			"Some content",
		}, "\n")
		violations := Check("test.mdx", strings.Split(content, "\n"))
		if len(violations) != 0 {
			t.Errorf(
				"expected 0 violations (blank after last), got %d",
				len(violations),
			)
		}
	})
}

func TestCheck_FewerThanTwoImports(t *testing.T) {
	t.Run("single import", func(t *testing.T) {
		violations := Check("test.mdx", []string{"import A from '/snippets/a.mdx';"})
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		violations := Check("test.mdx", []string{"# Hello world"})
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d", len(violations))
		}
	})
}

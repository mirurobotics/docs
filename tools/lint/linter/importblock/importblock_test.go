package importblock

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Run("no blank lines between imports", func(t *testing.T) {
		content := "import A from '/snippets/a.mdx';\nimport B from '/snippets/b.mdx';"
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("one blank line between first and second import", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
		if len(vs) > 0 && vs[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", vs[0].Line)
		}
	})

	t.Run("blank line before first import", func(t *testing.T) {
		content := strings.Join([]string{
			"",
			"import A from '/snippets/a.mdx';",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations (blank before first), got %d", len(vs))
		}
	})

	t.Run("blank line after last import", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"import B from '/snippets/b.mdx';",
			"",
			"Some content",
		}, "\n")
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations (blank after last), got %d", len(vs))
		}
	})

	t.Run("two blank lines in import block", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"",
			"",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 2 {
			t.Errorf("expected 2 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("single import", func(t *testing.T) {
		vs := Check("test.mdx", []string{"import A from '/snippets/a.mdx';"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := Check("test.mdx", []string{"# Hello world"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("non-blank content between imports no violation", func(t *testing.T) {
		content := strings.Join([]string{
			"import A from '/snippets/a.mdx';",
			"## Heading",
			"import B from '/snippets/b.mdx';",
		}, "\n")
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf(
				"expected 0 violations with non-blank content between imports, got %d",
				len(vs),
			)
		}
	})
}

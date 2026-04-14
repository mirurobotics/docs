package importused

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Run("name used as self-closing JSX", func(t *testing.T) {
		content := `import Framed from '/snippets/components/framed.jsx';

<Framed />`
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("name used as open tag", func(t *testing.T) {
		content := `import Framed from '/snippets/components/framed.jsx';

<Framed>some content</Framed>`
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("name not in body", func(t *testing.T) {
		content := `import Unused from '/snippets/foo.mdx';

Some text without the component.`
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d", len(vs))
		}
	})

	t.Run("named import two names one unused", func(t *testing.T) {
		content := `import { A, B } from '/snippets/components/badges.jsx';

<A />`
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
		if len(vs) > 0 && vs[0].Line != 1 {
			t.Errorf("expected violation on line 1, got line %d", vs[0].Line)
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := Check("test.mdx", []string{"# Hello"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("name in frontmatter only is unused", func(t *testing.T) {
		content := `---
title: "Framed content"
---
import Framed from '/snippets/components/framed.jsx';

No component used here.`
		vs := Check("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf(
				"expected 1 violation (frontmatter excluded), got %d: %v",
				len(vs), vs,
			)
		}
	})
}

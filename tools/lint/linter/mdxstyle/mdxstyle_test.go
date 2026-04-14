package mdxstyle

import "testing"

func TestCheck(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx';"
		vs := Check("test.mdx", []string{line})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("missing semicolon", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx'"
		vs := Check("test.mdx", []string{line})
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
	})

	t.Run("named import violation", func(t *testing.T) {
		line := "import { DeviceDef } from '/snippets/definitions/device.mdx';"
		vs := Check("test.mdx", []string{line})
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
	})

	t.Run("non-mdx import no violations", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx';"
		vs := Check("test.mdx", []string{line})
		if len(vs) != 0 {
			t.Errorf(
				"expected 0 violations for non-mdx import, got %d: %v",
				len(vs), vs,
			)
		}
	})
}

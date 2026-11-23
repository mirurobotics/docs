package mdxstyle

import "testing"

func TestCheck_NoViolations(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx';"
		violations := Check("test.mdx", []string{line})
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(violations), violations)
		}
	})

	t.Run("non-mdx import no violations", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx';"
		violations := Check("test.mdx", []string{line})
		if len(violations) != 0 {
			t.Errorf(
				"expected 0 violations for non-mdx import, got %d: %v",
				len(violations), violations,
			)
		}
	})
}

func TestCheck_Violations(t *testing.T) {
	t.Run("missing semicolon", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx'"
		violations := Check("test.mdx", []string{line})
		if len(violations) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(violations), violations)
		}
	})

	t.Run("named import violation", func(t *testing.T) {
		line := "import { DeviceDef } from '/snippets/definitions/device.mdx';"
		violations := Check("test.mdx", []string{line})
		if len(violations) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(violations), violations)
		}
	})
}

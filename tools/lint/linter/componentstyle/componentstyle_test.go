package componentstyle

import (
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func TestCheck_Accepted(t *testing.T) {
	good := "import { Framed } from '/snippets/components/framed.jsx';"

	t.Run("correct", func(t *testing.T) {
		violations := Check("test.mdx", []string{good})
		if len(violations) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(violations), violations)
		}
	})

	t.Run("non-component import no violations", func(t *testing.T) {
		line := "import Foo from '/snippets/definitions/foo.mdx';"
		violations := Check("test.mdx", []string{line})
		if len(violations) != 0 {
			t.Errorf(
				"non-component import: expected 0 violations, got %d",
				len(violations),
			)
		}
	})
}

func TestCheck_BraceSpacing(t *testing.T) {
	t.Run("missing space after brace", func(t *testing.T) {
		line := "import {Framed} from '/snippets/components/framed.jsx';"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "after '{'") {
			t.Error("expected violation for missing space after '{'")
		}
		if !hasMessage(violations, "before '}'") {
			t.Error("expected violation for missing space before '}'")
		}
	})

	t.Run("missing space before close brace", func(t *testing.T) {
		line := "import { Framed} from '/snippets/components/framed.jsx';"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "before '}'") {
			t.Error("expected violation for missing space before '}'")
		}
	})
}

func TestCheck_CommaSpacing(t *testing.T) {
	t.Run("no comma-space", func(t *testing.T) {
		line := "import { A,B } from '/snippets/components/field-badges.jsx';"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "single space after ','") {
			t.Errorf(
				"expected violation for missing space after ',', got %v",
				violations,
			)
		}
	})

	t.Run("space before comma", func(t *testing.T) {
		line := "import { A , B } from '/snippets/components/field-badges.jsx';"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "before ','") {
			t.Errorf("expected violation for space before ',', got %v", violations)
		}
	})
}

func TestCheck_ImportSyntax(t *testing.T) {
	t.Run("default import used", func(t *testing.T) {
		line := "import Framed from '/snippets/components/framed.jsx';"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "named import syntax") {
			t.Errorf(
				"expected violation for default import of component, got %v",
				violations,
			)
		}
	})

	t.Run("malformed brace order does not panic", func(t *testing.T) {
		line := "import Foo} from '/snippets/components/x.jsx' // {Bar"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "named import syntax") {
			t.Errorf("expected named-import-syntax violation, got %v", violations)
		}
	})
}

func TestCheck_PathAndSemicolon(t *testing.T) {
	t.Run("path ends in .mdx not .jsx", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.mdx';"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, ".jsx") {
			t.Errorf("expected violation for .mdx extension, got %v", violations)
		}
	})

	t.Run("missing semicolon", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx'"
		violations := Check("test.mdx", []string{line})
		if !hasMessage(violations, "';'") {
			t.Errorf("expected violation for missing semicolon, got %v", violations)
		}
	})
}

func hasMessage(violations []analysis.Violation, substr string) bool {
	for _, got := range violations {
		if strings.Contains(got.Message, substr) {
			return true
		}
	}
	return false
}

func TestCheckBracketSpacingNoClosure(t *testing.T) {
	raw := "import {Foo from '/snippets/components/foo.jsx'"
	violations := checkBracketSpacing("f.mdx", 1, raw, strings.Index(raw, "{"))
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if !strings.Contains(violations[0].Message, "named import syntax") {
		t.Errorf("unexpected message: %q", violations[0].Message)
	}
}

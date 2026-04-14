package componentstyle

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	good := "import { Framed } from '/snippets/components/framed.jsx';"

	t.Run("correct", func(t *testing.T) {
		vs := Check("test.mdx", []string{good})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("missing space after brace", func(t *testing.T) {
		line := "import {Framed} from '/snippets/components/framed.jsx';"
		vs := Check("test.mdx", []string{line})
		hasAfter, hasBefore := false, false
		for _, v := range vs {
			if strings.Contains(v.Message, "after '{'") {
				hasAfter = true
			}
			if strings.Contains(v.Message, "before '}'") {
				hasBefore = true
			}
		}
		if !hasAfter {
			t.Error("expected violation for missing space after '{'")
		}
		if !hasBefore {
			t.Error("expected violation for missing space before '}'")
		}
	})

	t.Run("missing space before close brace", func(t *testing.T) {
		line := "import { Framed} from '/snippets/components/framed.jsx';"
		vs := Check("test.mdx", []string{line})
		found := false
		for _, v := range vs {
			if strings.Contains(v.Message, "before '}'") {
				found = true
			}
		}
		if !found {
			t.Error("expected violation for missing space before '}'")
		}
	})

	t.Run("no comma-space", func(t *testing.T) {
		line := "import { A,B } from '/snippets/components/badges.jsx';"
		vs := Check("test.mdx", []string{line})
		found := false
		for _, v := range vs {
			if strings.Contains(v.Message, "single space after ','") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation for missing space after ',', got %v", vs)
		}
	})

	t.Run("space before comma", func(t *testing.T) {
		line := "import { A , B } from '/snippets/components/badges.jsx';"
		vs := Check("test.mdx", []string{line})
		found := false
		for _, v := range vs {
			if strings.Contains(v.Message, "before ','") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation for space before ',', got %v", vs)
		}
	})

	t.Run("default import used", func(t *testing.T) {
		line := "import Framed from '/snippets/components/framed.jsx';"
		vs := Check("test.mdx", []string{line})
		found := false
		for _, v := range vs {
			if strings.Contains(v.Message, "named import syntax") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation for default import of component, got %v", vs)
		}
	})

	t.Run("path ends in .mdx not .jsx", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.mdx';"
		vs := Check("test.mdx", []string{line})
		found := false
		for _, v := range vs {
			if strings.Contains(v.Message, ".jsx") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation for .mdx extension, got %v", vs)
		}
	})

	t.Run("missing semicolon", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx'"
		vs := Check("test.mdx", []string{line})
		found := false
		for _, v := range vs {
			if strings.Contains(v.Message, "';'") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected violation for missing semicolon, got %v", vs)
		}
	})

	t.Run("non-component import no violations", func(t *testing.T) {
		line := "import Foo from '/snippets/definitions/foo.mdx';"
		vs := Check("test.mdx", []string{line})
		if len(vs) != 0 {
			t.Errorf("non-component import: expected 0 violations, got %d", len(vs))
		}
	})
}

func TestCheckBracketSpacingNoClosure(t *testing.T) {
	// Raw has opening brace but no closing brace; should report
	// "must use named import syntax { }" rather than spacing errors.
	raw := "import {Foo from '/snippets/components/foo.jsx'"
	vs := checkBracketSpacing("f.mdx", 1, raw, strings.Index(raw, "{"))
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(vs), vs)
	}
	if !strings.Contains(vs[0].Message, "named import syntax") {
		t.Errorf("unexpected message: %q", vs[0].Message)
	}
}

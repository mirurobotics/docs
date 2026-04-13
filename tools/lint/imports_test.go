package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSingleImport(t *testing.T) {
	t.Run("named single", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx';"
		pi := parseSingleImport(1, line)
		if pi == nil {
			t.Fatal("expected non-nil ParsedImport")
		}
		if len(pi.Names) != 1 {
			t.Fatalf("expected 1 name, got %d", len(pi.Names))
		}
		if pi.Names[0].Name != "Framed" {
			t.Errorf("expected name Framed, got %q", pi.Names[0].Name)
		}
		if !pi.Names[0].IsNamed {
			t.Error("expected IsNamed=true")
		}
		if pi.Path != "/snippets/components/framed.jsx" {
			t.Errorf("expected path /snippets/components/framed.jsx, got %q", pi.Path)
		}
	})

	t.Run("named multi", func(t *testing.T) {
		line := "import { A, B } from '/snippets/components/badges.jsx'"
		pi := parseSingleImport(2, line)
		if pi == nil {
			t.Fatal("expected non-nil ParsedImport")
		}
		if len(pi.Names) != 2 {
			t.Fatalf("expected 2 names, got %d", len(pi.Names))
		}
		if pi.Names[0].Name != "A" || pi.Names[1].Name != "B" {
			t.Errorf("expected names [A, B], got %v", pi.Names)
		}
	})

	t.Run("default import", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx'"
		pi := parseSingleImport(3, line)
		if pi == nil {
			t.Fatal("expected non-nil ParsedImport")
		}
		if len(pi.Names) != 1 {
			t.Fatalf("expected 1 name, got %d", len(pi.Names))
		}
		if pi.Names[0].Name != "DeviceDef" {
			t.Errorf("expected name DeviceDef, got %q", pi.Names[0].Name)
		}
		if pi.Names[0].IsNamed {
			t.Error("expected IsNamed=false for default import")
		}
	})

	t.Run("double quotes", func(t *testing.T) {
		line := `import Usage from "/snippets/foo.mdx";`
		pi := parseSingleImport(1, line)
		if pi == nil {
			t.Fatal("expected non-nil ParsedImport")
		}
		if pi.Path != "/snippets/foo.mdx" {
			t.Errorf("expected path /snippets/foo.mdx, got %q", pi.Path)
		}
	})

	t.Run("non-import line returns nil", func(t *testing.T) {
		pi := parseSingleImport(1, "This is not an import")
		if pi != nil {
			t.Error("expected nil for non-import line")
		}
	})

	t.Run("named import unclosed brace", func(t *testing.T) {
		pi := parseSingleImport(1, "import { Foo from '/snippets/a.mdx';")
		if pi != nil {
			t.Error("expected nil for named import with no closing brace")
		}
	})

	t.Run("default import missing from keyword", func(t *testing.T) {
		pi := parseSingleImport(1, "import Foo '/snippets/a.mdx';")
		if pi != nil {
			t.Error("expected nil for import without 'from' keyword")
		}
	})

	t.Run("no quotes in import path", func(t *testing.T) {
		pi := parseSingleImport(1, "import { Foo } from /snippets/a.mdx;")
		if pi != nil {
			t.Error("expected nil when import path has no quotes")
		}
	})

	t.Run("unmatched closing quote in path", func(t *testing.T) {
		pi := parseSingleImport(1, "import { Foo } from path.mdx';")
		if pi != nil {
			t.Error("expected nil for unmatched closing quote in path")
		}
	})
}

func TestImportResolvesRule(t *testing.T) {
	root := t.TempDir()
	snippetsDir := filepath.Join(root, "snippets", "components")
	if err := os.MkdirAll(snippetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	realFile := filepath.Join(snippetsDir, "framed.jsx")
	if err := os.WriteFile(realFile, []byte("// framed"), 0o644); err != nil {
		t.Fatal(err)
	}

	rule := ImportResolvesRule{ContentRoot: root}

	t.Run("file exists", func(t *testing.T) {
		lines := []string{"import { Framed } from '/snippets/components/framed.jsx';"}
		vs := rule.CheckFile("test.mdx", lines)
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("file missing", func(t *testing.T) {
		lines := []string{"import { Missing } from '/snippets/components/missing.jsx';"}
		vs := rule.CheckFile("test.mdx", lines)
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d", len(vs))
		}
		if len(vs) > 0 && vs[0].Line != 1 {
			t.Errorf("expected violation on line 1, got line %d", vs[0].Line)
		}
	})

	t.Run("relative path skipped", func(t *testing.T) {
		lines := []string{"import Foo from './relative.mdx';"}
		vs := rule.CheckFile("test.mdx", lines)
		if len(vs) != 0 {
			t.Errorf("expected 0 violations for relative path, got %d", len(vs))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{"# Hello world"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})
}

func TestImportUsedRule(t *testing.T) {
	rule := ImportUsedRule{}

	t.Run("name used as self-closing JSX", func(t *testing.T) {
		content := `import Framed from '/snippets/components/framed.jsx';

<Framed />`
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("name used as open tag", func(t *testing.T) {
		content := `import Framed from '/snippets/components/framed.jsx';

<Framed>some content</Framed>`
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("name not in body", func(t *testing.T) {
		content := `import Unused from '/snippets/foo.mdx';

Some text without the component.`
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d", len(vs))
		}
	})

	t.Run("named import two names one unused", func(t *testing.T) {
		content := `import { A, B } from '/snippets/components/badges.jsx';

<A />`
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
		if len(vs) > 0 && vs[0].Line != 1 {
			t.Errorf("expected violation on line 1, got line %d", vs[0].Line)
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{"# Hello"})
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf(
				"expected 1 violation (frontmatter excluded), got %d: %v",
				len(vs), vs,
			)
		}
	})
}

func TestImportSortedRule(t *testing.T) {
	rule := ImportSortedRule{}

	t.Run("single import", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{"import A from '/snippets/a.mdx';"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("two imports in order", func(t *testing.T) {
		content := "import A from '/snippets/a.mdx';\nimport B from '/snippets/b.mdx';"
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("two imports out of order", func(t *testing.T) {
		content := "import B from '/snippets/b.mdx';\nimport A from '/snippets/a.mdx';"
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 1 {
			t.Errorf("expected 1 violation (only first out-of-order), got %d", len(vs))
		}
		if len(vs) > 0 && vs[0].Line != 2 {
			t.Errorf("expected violation on line 2, got line %d", vs[0].Line)
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{"# Hello"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})
}

func TestComponentImportStyleRule(t *testing.T) {
	rule := ComponentImportStyleRule{}

	good := "import { Framed } from '/snippets/components/framed.jsx';"

	t.Run("correct", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{good})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("missing space after brace", func(t *testing.T) {
		line := "import {Framed} from '/snippets/components/framed.jsx';"
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
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
		vs := rule.CheckFile("test.mdx", []string{line})
		if len(vs) != 0 {
			t.Errorf("non-component import: expected 0 violations, got %d", len(vs))
		}
	})
}

func TestMDXImportStyleRule(t *testing.T) {
	rule := MDXImportStyleRule{}

	t.Run("correct", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx';"
		vs := rule.CheckFile("test.mdx", []string{line})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("missing semicolon", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx'"
		vs := rule.CheckFile("test.mdx", []string{line})
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
	})

	t.Run("named import violation", func(t *testing.T) {
		line := "import { DeviceDef } from '/snippets/definitions/device.mdx';"
		vs := rule.CheckFile("test.mdx", []string{line})
		if len(vs) != 1 {
			t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
		}
	})

	t.Run("non-mdx import no violations", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx';"
		vs := rule.CheckFile("test.mdx", []string{line})
		if len(vs) != 0 {
			t.Errorf(
				"expected 0 violations for non-mdx import, got %d: %v",
				len(vs), vs,
			)
		}
	})
}

func TestImportBlockContiguousRule(t *testing.T) {
	rule := ImportBlockContiguousRule{}

	t.Run("no blank lines between imports", func(t *testing.T) {
		content := "import A from '/snippets/a.mdx';\nimport B from '/snippets/b.mdx';"
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 2 {
			t.Errorf("expected 2 violations, got %d: %v", len(vs), vs)
		}
	})

	t.Run("single import", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{"import A from '/snippets/a.mdx';"})
		if len(vs) != 0 {
			t.Errorf("expected 0 violations, got %d", len(vs))
		}
	})

	t.Run("no imports", func(t *testing.T) {
		vs := rule.CheckFile("test.mdx", []string{"# Hello world"})
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
		vs := rule.CheckFile("test.mdx", strings.Split(content, "\n"))
		if len(vs) != 0 {
			t.Errorf(
				"expected 0 violations with non-blank content between imports, got %d",
				len(vs),
			)
		}
	})
}

func TestFrontmatterEndNeverCloses(t *testing.T) {
	lines := []string{"---", "title: foo", "description: bar"}
	got := frontmatterEnd(lines)
	if got != -1 {
		t.Errorf("expected -1 for unclosed frontmatter, got %d", got)
	}
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

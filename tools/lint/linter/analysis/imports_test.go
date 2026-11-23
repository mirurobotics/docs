package analysis

import (
	"testing"
)

func TestParseSingleImport_Named(t *testing.T) {
	t.Run("named single", func(t *testing.T) {
		line := "import { Framed } from '/snippets/components/framed.jsx';"
		imp := mustParseImport(t, 1, line)
		if len(imp.Names) != 1 {
			t.Fatalf("expected 1 name, got %d", len(imp.Names))
		}
		if imp.Names[0].Name != "Framed" {
			t.Errorf("expected name Framed, got %q", imp.Names[0].Name)
		}
		if !imp.Names[0].IsNamed {
			t.Error("expected IsNamed=true")
		}
		if imp.Path != "/snippets/components/framed.jsx" {
			t.Errorf("expected path /snippets/components/framed.jsx, got %q", imp.Path)
		}
	})

	t.Run("named multi", func(t *testing.T) {
		line := "import { A, B } from '/snippets/components/field-badges.jsx'"
		imp := mustParseImport(t, 2, line)
		if len(imp.Names) != 2 {
			t.Fatalf("expected 2 names, got %d", len(imp.Names))
		}
		if imp.Names[0].Name != "A" || imp.Names[1].Name != "B" {
			t.Errorf("expected names [A, B], got %v", imp.Names)
		}
	})
}

func TestParseSingleImport_Default(t *testing.T) {
	t.Run("default import", func(t *testing.T) {
		line := "import DeviceDef from '/snippets/definitions/device.mdx'"
		imp := mustParseImport(t, 3, line)
		if len(imp.Names) != 1 {
			t.Fatalf("expected 1 name, got %d", len(imp.Names))
		}
		if imp.Names[0].Name != "DeviceDef" {
			t.Errorf("expected name DeviceDef, got %q", imp.Names[0].Name)
		}
		if imp.Names[0].IsNamed {
			t.Error("expected IsNamed=false for default import")
		}
	})

	t.Run("double quotes", func(t *testing.T) {
		line := `import Usage from "/snippets/foo.mdx";`
		imp := mustParseImport(t, 1, line)
		if imp.Path != "/snippets/foo.mdx" {
			t.Errorf("expected path /snippets/foo.mdx, got %q", imp.Path)
		}
	})
}

func mustParseImport(t *testing.T, lineNum int, line string) *ParsedImport {
	t.Helper()
	imp := ParseSingleImport(lineNum, line)
	if imp == nil {
		t.Fatal("expected non-nil ParsedImport")
	}
	return imp
}

func TestParseSingleImport_Rejected(t *testing.T) {
	cases := []struct {
		name, line string
	}{
		{"non-import line returns nil", "This is not an import"},
		{"named import unclosed brace", "import { Foo from '/snippets/a.mdx';"},
		{"default import missing from keyword", "import Foo '/snippets/a.mdx';"},
		{"default import blank name", "import \u00a0 from '/snippets/a.mdx';"},
		{"no quotes in import path", "import { Foo } from /snippets/a.mdx;"},
		{"unmatched closing quote in path", "import { Foo } from path.mdx';"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if imp := ParseSingleImport(1, tc.line); imp != nil {
				t.Errorf("expected nil for %q, got %+v", tc.line, imp)
			}
		})
	}
}

func TestParseImports(t *testing.T) {
	t.Run("multiple imports", func(t *testing.T) {
		lines := []string{
			"import A from '/snippets/a.mdx';",
			"# Heading",
			"import { B } from '/snippets/components/b.jsx';",
		}
		imports := ParseImports(lines)
		if len(imports) != 2 {
			t.Fatalf("expected 2 imports, got %d", len(imports))
		}
		if imports[0].Path != "/snippets/a.mdx" {
			t.Errorf("import 0: expected path /snippets/a.mdx, got %q", imports[0].Path)
		}
		if imports[1].Path != "/snippets/components/b.jsx" {
			t.Errorf("import 1: expected path /snippets/components/b.jsx, got %q",
				imports[1].Path)
		}
	})

	t.Run("no imports", func(t *testing.T) {
		lines := []string{"# Hello", "Some text"}
		imports := ParseImports(lines)
		if len(imports) != 0 {
			t.Errorf("expected 0 imports, got %d", len(imports))
		}
	})
}

func TestIsImportLine(t *testing.T) {
	if !IsImportLine("import Foo from './foo'") {
		t.Error("expected true for import line")
	}
	if IsImportLine("# Hello") {
		t.Error("expected false for non-import line")
	}
	if IsImportLine("exporting something") {
		t.Error("expected false for export line")
	}
}

func TestFrontmatterEnd(t *testing.T) {
	t.Run("valid frontmatter", func(t *testing.T) {
		lines := []string{"---", "title: foo", "---", "content"}
		got := FrontmatterEnd(lines)
		if got != 2 {
			t.Errorf("expected 2, got %d", got)
		}
	})

	t.Run("no frontmatter", func(t *testing.T) {
		lines := []string{"# Hello"}
		got := FrontmatterEnd(lines)
		if got != -1 {
			t.Errorf("expected -1, got %d", got)
		}
	})

	t.Run("empty lines", func(t *testing.T) {
		got := FrontmatterEnd(nil)
		if got != -1 {
			t.Errorf("expected -1, got %d", got)
		}
	})

	t.Run("unclosed frontmatter", func(t *testing.T) {
		lines := []string{"---", "title: foo", "description: bar"}
		got := FrontmatterEnd(lines)
		if got != -1 {
			t.Errorf("expected -1 for unclosed frontmatter, got %d", got)
		}
	})
}

func TestBodyLines(t *testing.T) {
	t.Run("with frontmatter and imports", func(t *testing.T) {
		lines := []string{
			"---",
			"title: Test",
			"---",
			"import Foo from '/snippets/foo.mdx';",
			"",
			"Body content here.",
		}
		body := BodyLines(lines)
		if len(body) != 2 {
			t.Fatalf("expected 2 body lines, got %d: %v", len(body), body)
		}
		if body[0] != "" {
			t.Errorf("expected empty line, got %q", body[0])
		}
		if body[1] != "Body content here." {
			t.Errorf("expected body content, got %q", body[1])
		}
	})

	t.Run("no frontmatter", func(t *testing.T) {
		lines := []string{"import A from './a';", "content"}
		body := BodyLines(lines)
		if len(body) != 1 {
			t.Fatalf("expected 1 body line, got %d", len(body))
		}
		if body[0] != "content" {
			t.Errorf("expected 'content', got %q", body[0])
		}
	})
}

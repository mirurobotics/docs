package redirects

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func TestCheck_NoViolations(t *testing.T) {
	t.Run("missing_docs_json", func(t *testing.T) {
		root := t.TempDir()
		violations := Check(root)
		if violations != nil {
			t.Errorf("expected nil, got %v", violations)
		}
	})

	t.Run("success_no_violations", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"docs.json": `{"redirects": []}`})
		violations := Check(root)
		if violations != nil {
			t.Errorf("expected nil, got %v", violations)
		}
	})
}

func TestCheck_ReadError(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	violations := Check(root)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	got := violations[0]
	if !strings.HasPrefix(got.Message, "read error:") {
		t.Errorf("message: got %q, want prefix %q", got.Message, "read error:")
	}
	if got.Line != 1 {
		t.Errorf("Line: got %d, want 1", got.Line)
	}
	if got.Col != 1 {
		t.Errorf("Col: got %d, want 1", got.Col)
	}
	if got.File != "docs.json" {
		t.Errorf("File: got %q, want docs.json", got.File)
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "docs.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	violations := Check(root)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if !strings.HasPrefix(violations[0].Message, "invalid JSON:") {
		t.Errorf(
			"message: got %q, want prefix %q",
			violations[0].Message, "invalid JSON:",
		)
	}
}

func TestCheck_DelegatesToValidate(t *testing.T) {
	docsJSON := `{
  "redirects": [
    {"source": "/api/foo", "destination": "/missing"}
  ]
}`
	root := setupContentRoot(t, map[string]string{"docs.json": docsJSON})
	violations := Check(root)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	assertViolation(t, violations[0], 0, "missing destination")
}

func TestValidate_LineNumbers(t *testing.T) {
	docsJSON := "{\n" +
		"  \"redirects\": [\n" +
		"    {\"source\": \"/a\", \"destination\": \"/missing1\"},\n" +
		"    \"bad-string-entry\",\n" +
		"    {\"source\": \"/b\", \"destination\": \"/missing2\"}\n" +
		"  ]\n" +
		"}\n"
	root := t.TempDir()
	violations := validate([]byte(docsJSON), root)
	if len(violations) != 3 {
		t.Fatalf("expected 3 violations, got %d: %v", len(violations), violations)
	}
	assertViolation(t, violations[0], 3, "missing destination")
	assertViolation(t, violations[1], 4, "not an object")
	assertViolation(t, violations[2], 5, "missing destination")
}

func TestCollectOpenAPISources_Found(t *testing.T) {
	t.Run("nav_nested_maps", func(t *testing.T) {
		input := map[string]any{
			"nav": map[string]any{"openapi": map[string]any{"source": "api/a.yaml"}},
		}
		got := collectOpenAPISources(input)
		want := map[string]bool{"api/a.yaml": true}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("nav_array_of_maps", func(t *testing.T) {
		input := map[string]any{
			"nav": []any{
				map[string]any{"openapi": map[string]any{"source": "api/b.yaml"}},
				map[string]any{"openapi": map[string]any{"source": "api/c.yaml"}},
			},
		}
		got := collectOpenAPISources(input)
		want := map[string]bool{"api/b.yaml": true, "api/c.yaml": true}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("multiple_levels_combined", func(t *testing.T) {
		input := map[string]any{
			"nav": []any{
				map[string]any{
					"groups": []any{
						map[string]any{
							"openapi": map[string]any{"source": "api/d.yaml"},
						},
					},
				},
			},
			"openapi": map[string]any{"source": "api/e.yaml"},
		}
		got := collectOpenAPISources(input)
		want := map[string]bool{"api/d.yaml": true, "api/e.yaml": true}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestCollectOpenAPISources_Empty(t *testing.T) {
	t.Run("openapi_non_map_defensive", func(t *testing.T) {
		input := map[string]any{"openapi": "scalar-not-a-map"}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("openapi_source_non_string_defensive", func(t *testing.T) {
		input := map[string]any{"openapi": map[string]any{"source": 42}}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("no_nav_at_all", func(t *testing.T) {
		input := map[string]any{"name": "x"}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("nav_is_scalar", func(t *testing.T) {
		input := map[string]any{"nav": "scalar"}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})
}

func TestFormatMessage(t *testing.T) {
	got := formatMessage(2, "source", "/x", "bad prefix")
	want := `redirects[2] source "/x": bad prefix`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

type wantViolation struct {
	line   int
	substr string
}

type validateCase struct {
	name     string
	docsJSON string
	files    map[string]string
	wants    []wantViolation
}

func runValidateCases(t *testing.T, cases []validateCase) {
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{}
			for rel, contents := range tc.files {
				files[rel] = contents
			}
			root := setupContentRoot(t, files)
			violations := validate([]byte(tc.docsJSON), root)

			if len(tc.wants) == 0 {
				if len(violations) != 0 {
					t.Errorf(
						"expected no violations, got %d: %v",
						len(violations), violations,
					)
				}
				return
			}
			if len(violations) != len(tc.wants) {
				t.Fatalf(
					"expected %d violations, got %d: %v",
					len(tc.wants), len(violations), violations,
				)
			}
			for i, wanted := range tc.wants {
				assertViolation(t, violations[i], wanted.line, wanted.substr)
			}
		})
	}
}

func setupContentRoot(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, contents := range files {
		full := filepath.Join(root, rel)
		if !strings.HasSuffix(rel, "/") {
			writeContentFile(t, full, contents)
			continue
		}
		if err := os.MkdirAll(full, 0o755); err != nil {
			t.Fatalf("mkdir %q: %v", full, err)
		}
	}
	return root
}

func writeContentFile(t *testing.T, full, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir parent of %q: %v", full, err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %q: %v", full, err)
	}
}

func assertViolation(
	t *testing.T,
	got analysis.Violation,
	wantLine int,
	substr string,
) {
	t.Helper()
	if got.File != "docs.json" {
		t.Errorf("File: got %q, want %q", got.File, "docs.json")
	}
	if got.Col != 1 {
		t.Errorf("Col: got %d, want 1", got.Col)
	}
	if wantLine > 0 && got.Line != wantLine {
		t.Errorf("Line: got %d, want %d", got.Line, wantLine)
	}
	if !strings.Contains(got.Message, substr) {
		t.Errorf("Message: %q does not contain %q", got.Message, substr)
	}
}

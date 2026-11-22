package redirects

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func setupContentRoot(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, contents := range files {
		full := filepath.Join(root, rel)
		if strings.HasSuffix(rel, "/") {
			if err := os.MkdirAll(full, 0o755); err != nil {
				t.Fatalf("mkdir %q: %v", full, err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir parent of %q: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
			t.Fatalf("write %q: %v", full, err)
		}
	}
	return root
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

func TestCheck(t *testing.T) {
	t.Run("missing_docs_json", func(t *testing.T) {
		root := t.TempDir()
		violations := Check(root)
		if violations != nil {
			t.Errorf("expected nil, got %v", violations)
		}
	})

	t.Run("read_error_eisdir", func(t *testing.T) {
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
	})

	t.Run("invalid_json", func(t *testing.T) {
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
	})

	t.Run("delegates_with_violations", func(t *testing.T) {
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
	})

	t.Run("success_no_violations", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"docs.json": `{"redirects": []}`})
		violations := Check(root)
		if violations != nil {
			t.Errorf("expected nil, got %v", violations)
		}
	})
}

func TestValidate(t *testing.T) {
	type want struct {
		line   int
		substr string
	}

	cases := []struct {
		name     string
		docsJSON string
		files    map[string]string
		wants    []want
	}{
		{
			name:     "top_level_not_object",
			docsJSON: `[1,2,3]`,
			wants:    []want{{line: 1, substr: "invalid JSON:"}},
		},
		{name: "no_redirects_key", docsJSON: `{"name": "x"}`, wants: nil},
		{
			name:     "redirects_not_array",
			docsJSON: `{"redirects": {"a": 1}}`,
			wants:    nil,
		},
		{name: "redirects_empty", docsJSON: `{"redirects": []}`, wants: nil},
		{
			name:     "entry_not_object",
			docsJSON: `{"redirects": ["bad"]}`,
			wants:    []want{{substr: `redirects[0] entry "": not an object`}},
		},
		{
			name:     "missing_source",
			docsJSON: `{"redirects": [{"destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants: []want{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			name:     "missing_destination",
			docsJSON: `{"redirects": [{"source": "/x"}]}`,
			wants: []want{
				{substr: `redirects[0] destination "": must be a non-empty string`},
			},
		},
		{
			name:     "source_non_string_scalar",
			docsJSON: `{"redirects": [{"source": 5, "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants: []want{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			name:     "source_empty_string",
			docsJSON: `{"redirects": [{"source": "", "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants: []want{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			name:     "source_no_leading_slash",
			docsJSON: `{"redirects": [{"source": "x", "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants:    []want{{substr: `must start with '/'`}},
		},
		{
			name:     "destination_no_leading_slash",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "y"}]}`,
			wants:    []want{{substr: `must start with '/' (or http(s)://)`}},
		},
		{
			name:     "source_root_relative_ok",
			docsJSON: `{"redirects": [{"source": "/api/x", "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants:    nil,
		},
		{
			name:     "source_exactly_docs",
			docsJSON: `{"redirects": [{"source": "/docs", "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants:    nil,
		},
		{
			name: "source_dead_redirect_mdx",
			docsJSON: `{"redirects": [{"source": "/dead",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"dead.mdx": "dead", "y.mdx": "y"},
			wants: []want{
				{substr: `dead redirect (source resolves to a real page)`},
			},
		},
		{
			name: "source_dead_redirect_md",
			docsJSON: `{"redirects": [{"source": "/dead",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"dead.md": "dead", "y.mdx": "y"},
			wants: []want{
				{substr: `dead redirect (source resolves to a real page)`},
			},
		},
		{
			name: "source_non_wildcard_ok",
			docsJSON: `{"redirects": [{"source": "/old",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"new.mdx": "new"},
			wants: nil,
		},
		{
			name: "wildcard_source_prefix_real_page",
			docsJSON: `{"redirects": [{"source": "/old/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"old.mdx": "old", "new.mdx": "new"},
			wants: []want{
				{substr: `dead redirect ` +
					`(wildcard source prefix resolves to a real page)`},
			},
		},
		{
			name: "wildcard_source_no_dir_ok",
			docsJSON: `{"redirects": [{"source": "/old/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"new.mdx": "new"},
			wants: nil,
		},
		{
			name: "destination_missing_page",
			docsJSON: `{"redirects": [{"source": "/old",` +
				` "destination": "/missing"}]}`,
			wants: []want{
				{substr: `missing destination (no .mdx or .md page exists)`},
			},
		},
		{
			name: "destination_md_ok",
			docsJSON: `{"redirects": [{"source": "/old",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"new.md": "new"},
			wants: nil,
		},
		{
			name: "wildcard_source_dir_has_pages",
			docsJSON: `{"redirects": [{"source": "/old/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{
				"old/sub/page.mdx": "page",
				"new.mdx":          "new",
			},
			wants: []want{
				{substr: `dead redirect (wildcard source prefix has real pages)`},
			},
		},
		{
			name: "wildcard_source_dir_empty",
			docsJSON: `{"redirects": [{"source": "/old/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"old/": "", "new.mdx": "new"},
			wants: nil,
		},
		{
			name:     "destination_root_relative_missing",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "/api/y"}]}`,
			wants:    []want{{substr: `missing destination`}},
		},
		{
			name:     "destination_http_ok",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "https://example.com/x"}]}`,
			wants:    nil,
		},
		{
			name:     "destination_exactly_docs",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "/docs"}]}`,
			wants:    []want{{substr: `missing destination`}},
		},
		{
			name: "destination_wildcard_dir_ok",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "/new/:slug*"}]}`,
			files: map[string]string{"new/sub.mdx": "sub"},
			wants: nil,
		},
		{
			name: "destination_wildcard_openapi_yaml_ok",
			docsJSON: `{
  "nav": [{"openapi": {"source": "api/spec.yaml"}}],
  "redirects": [{"source": "/x", "destination": "/api/spec/:slug*"}]
}`,
			files: map[string]string{"api/spec.yaml": "openapi: 3.0"},
			wants: nil,
		},
		{
			name: "destination_wildcard_openapi_yaml_missing",
			docsJSON: `{
  "nav": [{"openapi": {"source": "api/spec.yaml"}}],
  "redirects": [{"source": "/x", "destination": "/api/spec/:slug*"}]
}`,
			wants: []want{{substr: `wildcard prefix not a directory`}},
		},
		{
			name: "destination_wildcard_not_dir_no_yaml",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "/missing/:slug*"}]}`,
			wants: []want{{substr: `wildcard prefix not a directory`}},
		},
		{
			name: "path_with_query_and_fragment",
			docsJSON: `{"redirects": [{"source": "/old?x=1#frag",` +
				` "destination": "/new?y=2#bar"}]}`,
			files: map[string]string{"new.mdx": "new"},
			wants: nil,
		},
		{
			name: "trailing_slash_stripped",
			docsJSON: `{"redirects": [{"source": "/old/",` +
				` "destination": "/new/"}]}`,
			files: map[string]string{"new.mdx": "new"},
			wants: nil,
		},
		{
			name: "source_no_slash_and_destination_no_slash",
			docsJSON: `{"redirects": [{"source": "x",` +
				` "destination": "y"}]}`,
			wants: []want{
				{substr: `source "x": bad path: must start with '/'`},
				{substr: `destination "y":` +
					` bad path: must start with '/' (or http(s)://)`},
			},
		},
		{
			name: "destination_http_uppercase_rejected",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "HTTPS://example.com"}]}`,
			wants: []want{{substr: `must start with '/' (or http(s)://)`}},
		},
		{
			name: "source_resolves_to_dir",
			docsJSON: `{"redirects": [{"source": "/old",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"old.mdx/": "", "new.mdx": "new"},
			wants: nil,
		},
		{
			name: "wildcard_source_empty_prefix",
			docsJSON: `{"redirects": [{"source": "/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"new.mdx": "new"},
			wants: []want{
				{substr: `dead redirect (wildcard source prefix has real pages)`},
			},
		},
		{
			name: "source_with_dot_dot_segment",
			docsJSON: `{"redirects": [{"source": "/../README",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"y.mdx": "y"},
			wants: []want{
				{substr: `source "/../README":` +
					` bad path: contains '..' or '.' segment`},
			},
		},
		{
			name: "destination_with_dot_dot_segment",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "/x/../../escape"}]}`,
			wants: []want{
				{substr: `destination "/x/../../escape":` +
					` bad path: contains '..' or '.' segment`},
			},
		},
		{
			name: "source_with_dot_segment",
			docsJSON: `{"redirects": [{"source": "/./foo",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"y.mdx": "y"},
			wants: []want{
				{substr: `source "/./foo":` +
					` bad path: contains '..' or '.' segment`},
			},
		},
		{
			name: "wildcard_source_prefix_is_openapi_yaml",
			docsJSON: `{
  "nav": [{"openapi": {"source": "api/spec.yaml"}}],
  "redirects": [{"source": "/api/spec/:slug*",` +
				` "destination": "/y"}]
}`,
			files: map[string]string{
				"api/spec.yaml": "openapi: 3.0",
				"y.mdx":         "y",
			},
			wants: []want{
				{substr: `dead redirect ` +
					`(wildcard source prefix has Mintlify-generated pages)`},
			},
		},
		{
			name: "wildcard_source_prefix_yaml_registered_but_yaml_missing",
			docsJSON: `{
  "nav": [{"openapi": {"source": "api/spec.yaml"}}],
  "redirects": [{"source": "/api/spec/:slug*",` +
				` "destination": "/y"}]
}`,
			files: map[string]string{"y.mdx": "y"},
			wants: nil,
		},
		{
			name: "wildcard_source_prefix_yaml_not_registered",
			docsJSON: `{"redirects": [{"source": "/api/spec/:slug*",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"y.mdx": "y"},
			wants: nil,
		},
		{
			name: "old_docs_wildcard_to_root",
			docsJSON: `{"redirects": [{"source": "/docs/:slug*",` +
				` "destination": "/:slug*"}]}`,
			wants: nil,
		},
		{
			name: "old_docs_specific_page_to_root",
			docsJSON: `{"redirects": [{"source": "/docs/references/cli/install",` +
				` "destination": "/developers/cli/install"}]}`,
			files: map[string]string{"developers/cli/install.mdx": "install"},
			wants: nil,
		},
	}

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

	t.Run("line_numbers_interleaved", func(t *testing.T) {
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
	})
}

func TestCollectOpenAPISources(t *testing.T) {
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

func TestCleanPath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"/x", "x"},
		{"/x/", "x"},
		{"/x?q=1", "x"},
		{"/x#frag", "x"},
		{"/x?q=1#frag", "x"},
		{"/", ""},
		{"x", "x"},
		{"/x/?q=1", "x"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := cleanPath(tc.in)
			if got != tc.want {
				t.Errorf("cleanPath(%q): got %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSplitWildcard(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		wantPrefix  []string
		wantHasWild bool
	}{
		{"no_wildcard", "x/y", []string{"x", "y"}, false},
		{"trailing_wildcard", "x/:slug*", []string{"x"}, true},
		{"trailing_wildcard_no_star", "x/:slug", []string{"x"}, true},
		{"wildcard_first_segment", ":slug*", nil, true},
		{"empty", "", nil, false},
		{"only_wildcard", ":slug*", nil, true},
		{"empty_segments_skipped", "//x", []string{"x"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotPrefix, gotHas := splitWildcard(tc.in)
			if !reflect.DeepEqual(gotPrefix, tc.wantPrefix) {
				t.Errorf("prefix: got %v, want %v", gotPrefix, tc.wantPrefix)
			}
			if gotHas != tc.wantHasWild {
				t.Errorf("hasWildcard: got %v, want %v", gotHas, tc.wantHasWild)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	got := formatMessage(2, "source", "/x", "bad prefix")
	want := `redirects[2] source "/x": bad prefix`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLineForOffset(t *testing.T) {
	text := "abc\ndef\nghi"
	cases := []struct {
		off  int
		want int
	}{
		{0, 1},
		{2, 1},
		{3, 1},
		{4, 2},
		{7, 2},
		{8, 3},
		{11, 3},
	}
	for _, tc := range cases {
		got := lineForOffset(text, tc.off)
		if got != tc.want {
			t.Errorf("lineForOffset(%d): got %d, want %d", tc.off, got, tc.want)
		}
	}
}

func TestLineLookup(t *testing.T) {
	t.Run("count_zero", func(t *testing.T) {
		got := lineLookup("{}", 0)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("decoder_error", func(t *testing.T) {
		got := lineLookup("not json", 1)
		want := []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("redirects_token_then_eof", func(t *testing.T) {
		got := lineLookup(`{"redirects":`, 1)
		want := []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("array_open_then_eof_mid_scan", func(t *testing.T) {
		got := lineLookup(`{"redirects": [`, 1)
		want := []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("more_entries_than_count", func(t *testing.T) {
		text := "{\n" +
			"  \"redirects\": [\n" +
			"    {\"a\":1},\n" +
			"    {\"b\":2},\n" +
			"    {\"c\":3}\n" +
			"  ]\n" +
			"}\n"
		got := lineLookup(text, 2)
		if len(got) != 2 {
			t.Fatalf("got len %d, want 2", len(got))
		}
		if got[0] != 3 || got[1] != 4 {
			t.Errorf("got %v, want [3 4]", got)
		}
	})

	t.Run("nested_objects", func(t *testing.T) {
		text := "{\n" +
			"  \"redirects\": [\n" +
			"    {\"x\": {\"nested\": 1}}\n" +
			"  ]\n" +
			"}\n"
		got := lineLookup(text, 1)
		if len(got) != 1 || got[0] != 3 {
			t.Errorf("got %v, want [3]", got)
		}
	})
}

func TestDirHasPages(t *testing.T) {
	t.Run("missing_dir", func(t *testing.T) {
		root := t.TempDir()
		if dirHasPages(filepath.Join(root, "nope")) {
			t.Error("expected false for missing dir")
		}
	})

	t.Run("md_extension", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"a.md": "x"})
		if !dirHasPages(root) {
			t.Error("expected true for dir with .md")
		}
	})

	t.Run("non_page_files_only", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"a.txt": "x"})
		if dirHasPages(root) {
			t.Error("expected false for dir with .txt only")
		}
	})

	t.Run("subdir_no_pages", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"0empty/": "", "z.mdx": "x"})
		if !dirHasPages(root) {
			t.Error("expected true for dir with .mdx sibling and empty subdir")
		}
	})
}

func TestTokenStart(t *testing.T) {
	t.Run("whitespace_only", func(t *testing.T) {
		got := tokenStart("   abc", 0)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})

	t.Run("comma_then_brace", func(t *testing.T) {
		got := tokenStart(",   {", 0)
		if got != 4 {
			t.Errorf("got %d, want 4", got)
		}
	})

	t.Run("offset_at_eof", func(t *testing.T) {
		got := tokenStart("abc", 3)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})

	t.Run("tab_and_cr", func(t *testing.T) {
		got := tokenStart("\t\r\n{", 0)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})
}

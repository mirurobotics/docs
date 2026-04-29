package redirects

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// setupContentRoot creates a tempdir, writes each file at its relative path
// (creating parent dirs as needed), and returns the root.
//   - keys ending in "/" create an empty directory.
//   - keys not ending in "/" create a regular file with the given contents.
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

// assertViolation checks a single violation's File="docs.json", Col=1,
// Line=wantLine (when wantLine > 0), and Message contains substr.
func assertViolation(t *testing.T, v analysis.Violation, wantLine int, substr string) {
	t.Helper()
	if v.File != "docs.json" {
		t.Errorf("File: got %q, want %q", v.File, "docs.json")
	}
	if v.Col != 1 {
		t.Errorf("Col: got %d, want 1", v.Col)
	}
	if wantLine > 0 && v.Line != wantLine {
		t.Errorf("Line: got %d, want %d", v.Line, wantLine)
	}
	if !strings.Contains(v.Message, substr) {
		t.Errorf("Message: %q does not contain %q", v.Message, substr)
	}
}

func TestCheck(t *testing.T) {
	t.Run("missing_docs_json", func(t *testing.T) {
		root := t.TempDir()
		vs := Check(root)
		if vs != nil {
			t.Errorf("expected nil, got %v", vs)
		}
	})

	t.Run("read_error_eisdir", func(t *testing.T) {
		root := t.TempDir()
		// docs.json is a directory, so ReadFile returns EISDIR (not IsNotExist).
		if err := os.MkdirAll(filepath.Join(root, "docs.json"), 0o755); err != nil {
			t.Fatal(err)
		}
		vs := Check(root)
		if len(vs) != 1 {
			t.Fatalf("expected 1 violation, got %d: %v", len(vs), vs)
		}
		v := vs[0]
		if !strings.HasPrefix(v.Message, "read error:") {
			t.Errorf("message: got %q, want prefix %q", v.Message, "read error:")
		}
		if v.Line != 1 {
			t.Errorf("Line: got %d, want 1", v.Line)
		}
		if v.Col != 1 {
			t.Errorf("Col: got %d, want 1", v.Col)
		}
		if v.File != "docs.json" {
			t.Errorf("File: got %q, want docs.json", v.File)
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "docs.json"), []byte("{not json"), 0o644); err != nil {
			t.Fatal(err)
		}
		vs := Check(root)
		if len(vs) != 1 {
			t.Fatalf("expected 1 violation, got %d: %v", len(vs), vs)
		}
		if !strings.HasPrefix(vs[0].Message, "invalid JSON:") {
			t.Errorf("message: got %q, want prefix %q", vs[0].Message, "invalid JSON:")
		}
	})

	t.Run("delegates_with_violations", func(t *testing.T) {
		docsJSON := `{
  "redirects": [
    {"source": "/api/foo", "destination": "/docs/x"}
  ]
}`
		root := setupContentRoot(t, map[string]string{
			"docs.json":  docsJSON,
			"docs/x.mdx": "# x",
		})
		vs := Check(root)
		if len(vs) != 1 {
			t.Fatalf("expected 1 violation, got %d: %v", len(vs), vs)
		}
		assertViolation(t, vs[0], 0, "bad prefix")
	})

	t.Run("success_no_violations", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{
			"docs.json": `{"redirects": []}`,
		})
		vs := Check(root)
		if vs != nil {
			t.Errorf("expected nil, got %v", vs)
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
			// Case 1: top-level not an object
			name:     "top_level_not_object",
			docsJSON: `[1,2,3]`,
			wants: []want{
				{line: 1, substr: "invalid JSON:"},
			},
		},
		{
			// Case 2: no redirects key
			name:     "no_redirects_key",
			docsJSON: `{"name": "x"}`,
			wants:    nil,
		},
		{
			// Case 3: redirects not an array
			name:     "redirects_not_array",
			docsJSON: `{"redirects": {"a": 1}}`,
			wants:    nil,
		},
		{
			// Case 4a: redirects is empty
			name:     "redirects_empty",
			docsJSON: `{"redirects": []}`,
			wants:    nil,
		},
		{
			// Case 4b: redirects entry not an object (string)
			name:     "entry_not_object",
			docsJSON: `{"redirects": ["bad"]}`,
			wants: []want{
				{substr: `redirects[0] entry "": not an object`},
			},
		},
		{
			// Case 5: missing source
			name:     "missing_source",
			docsJSON: `{"redirects": [{"destination": "/docs/y"}]}`,
			files:    map[string]string{"docs/y.mdx": "y"},
			wants: []want{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			// Case 5b: missing destination
			name:     "missing_destination",
			docsJSON: `{"redirects": [{"source": "/docs/x"}]}`,
			wants: []want{
				{substr: `redirects[0] destination "": must be a non-empty string`},
			},
		},
		{
			// Case 5c: source non-string scalar (number) - reports empty value
			name:     "source_non_string_scalar",
			docsJSON: `{"redirects": [{"source": 5, "destination": "/docs/y"}]}`,
			files:    map[string]string{"docs/y.mdx": "y"},
			wants: []want{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			// Case 5d: source empty string
			name:     "source_empty_string",
			docsJSON: `{"redirects": [{"source": "", "destination": "/docs/y"}]}`,
			files:    map[string]string{"docs/y.mdx": "y"},
			wants: []want{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			// Case 6: source missing leading slash
			name:     "source_no_leading_slash",
			docsJSON: `{"redirects": [{"source": "docs/x", "destination": "/docs/y"}]}`,
			files:    map[string]string{"docs/y.mdx": "y"},
			wants: []want{
				{substr: `must start with '/'`},
			},
		},
		{
			// Case 7: destination missing leading slash and not http
			name:     "destination_no_leading_slash",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "docs/y"}]}`,
			wants: []want{
				{substr: `must start with '/' (or http(s)://)`},
			},
		},
		{
			// Case 8: source bad prefix (not /docs/)
			name:     "source_bad_prefix",
			docsJSON: `{"redirects": [{"source": "/api/x", "destination": "/docs/y"}]}`,
			files:    map[string]string{"docs/y.mdx": "y"},
			wants: []want{
				{substr: `bad prefix (must start with /docs/)`},
			},
		},
		{
			// Case 8b: source exactly /docs (no slash)
			name:     "source_exactly_docs",
			docsJSON: `{"redirects": [{"source": "/docs", "destination": "/docs/y"}]}`,
			files:    map[string]string{"docs/y.mdx": "y"},
			wants:    nil,
		},
		{
			// Case 9: source dead redirect (resolves to a real .mdx page)
			name:     "source_dead_redirect_mdx",
			docsJSON: `{"redirects": [{"source": "/docs/dead", "destination": "/docs/y"}]}`,
			files: map[string]string{
				"docs/dead.mdx": "dead",
				"docs/y.mdx":    "y",
			},
			wants: []want{
				{substr: `dead redirect (source resolves to a real page)`},
			},
		},
		{
			// Case 10: source dead redirect via .md
			name:     "source_dead_redirect_md",
			docsJSON: `{"redirects": [{"source": "/docs/dead", "destination": "/docs/y"}]}`,
			files: map[string]string{
				"docs/dead.md": "dead",
				"docs/y.mdx":   "y",
			},
			wants: []want{
				{substr: `dead redirect (source resolves to a real page)`},
			},
		},
		{
			// Case 11: source non-wildcard, no real page (OK)
			name:     "source_non_wildcard_ok",
			docsJSON: `{"redirects": [{"source": "/docs/old", "destination": "/docs/new"}]}`,
			files:    map[string]string{"docs/new.mdx": "new"},
			wants:    nil,
		},
		{
			// Case 12: wildcard source prefix resolves to a real page
			name:     "wildcard_source_prefix_real_page",
			docsJSON: `{"redirects": [{"source": "/docs/old/:slug*", "destination": "/docs/new"}]}`,
			files: map[string]string{
				"docs/old.mdx": "old",
				"docs/new.mdx": "new",
			},
			wants: []want{
				{substr: `dead redirect (wildcard source prefix resolves to a real page)`},
			},
		},
		{
			// Case 13: canonical positive case — wildcard source dir doesn't exist
			name:     "wildcard_source_no_dir_ok",
			docsJSON: `{"redirects": [{"source": "/docs/old/:slug*", "destination": "/docs/new"}]}`,
			files:    map[string]string{"docs/new.mdx": "new"},
			wants:    nil,
		},
		{
			// Case 14: destination missing
			name:     "destination_missing_page",
			docsJSON: `{"redirects": [{"source": "/docs/old", "destination": "/docs/missing"}]}`,
			wants: []want{
				{substr: `missing destination (no .mdx or .md page exists)`},
			},
		},
		{
			// Case 15: canonical positive — destination exists as .md
			name:     "destination_md_ok",
			docsJSON: `{"redirects": [{"source": "/docs/old", "destination": "/docs/new"}]}`,
			files:    map[string]string{"docs/new.md": "new"},
			wants:    nil,
		},
		{
			// Case 16: wildcard source dir with pages → recursion in dirHasPages
			name:     "wildcard_source_dir_has_pages",
			docsJSON: `{"redirects": [{"source": "/docs/old/:slug*", "destination": "/docs/new"}]}`,
			files: map[string]string{
				"docs/old/sub/page.mdx": "page",
				"docs/new.mdx":          "new",
			},
			wants: []want{
				{substr: `dead redirect (wildcard source prefix has real pages)`},
			},
		},
		{
			// Case 17: wildcard source dir exists but is empty (no pages)
			name:     "wildcard_source_dir_empty",
			docsJSON: `{"redirects": [{"source": "/docs/old/:slug*", "destination": "/docs/new"}]}`,
			files: map[string]string{
				"docs/old/":    "",
				"docs/new.mdx": "new",
			},
			wants: nil,
		},
		{
			// Case 18: destination bad prefix
			name:     "destination_bad_prefix",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "/api/y"}]}`,
			wants: []want{
				{substr: `bad prefix (must start with /docs/)`},
			},
		},
		{
			// Case 18b: destination http(s):// — no fs check
			name:     "destination_http_ok",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "https://example.com/x"}]}`,
			wants:    nil,
		},
		{
			// Case 19: destination exactly /docs (canonical form), no fs file
			name:     "destination_exactly_docs",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "/docs"}]}`,
			wants: []want{
				{substr: `missing destination`},
			},
		},
		{
			// Case 20: canonical positive — destination wildcard prefix is a directory
			name:     "destination_wildcard_dir_ok",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "/docs/new/:slug*"}]}`,
			files: map[string]string{
				"docs/new/sub.mdx": "sub",
			},
			wants: nil,
		},
		{
			// Case 21: OpenAPI escape hatch — wildcard prefix is yaml file registered in nav
			name: "destination_wildcard_openapi_yaml_ok",
			docsJSON: `{
  "nav": [{"openapi": {"source": "docs/api/spec.yaml"}}],
  "redirects": [{"source": "/docs/x", "destination": "/docs/api/spec/:slug*"}]
}`,
			files: map[string]string{
				"docs/api/spec.yaml": "openapi: 3.0",
			},
			wants: nil,
		},
		{
			// Case 22: OpenAPI escape hatch — yaml file registered but missing on disk
			name: "destination_wildcard_openapi_yaml_missing",
			docsJSON: `{
  "nav": [{"openapi": {"source": "docs/api/spec.yaml"}}],
  "redirects": [{"source": "/docs/x", "destination": "/docs/api/spec/:slug*"}]
}`,
			wants: []want{
				{substr: `wildcard prefix not a directory`},
			},
		},
		{
			// Case 23: destination wildcard, prefix not dir, not registered yaml
			name:     "destination_wildcard_not_dir_no_yaml",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "/docs/missing/:slug*"}]}`,
			wants: []want{
				{substr: `wildcard prefix not a directory`},
			},
		},
		{
			// Case 24: query string and fragment stripped from path
			name:     "path_with_query_and_fragment",
			docsJSON: `{"redirects": [{"source": "/docs/old?x=1#frag", "destination": "/docs/new?y=2#bar"}]}`,
			files:    map[string]string{"docs/new.mdx": "new"},
			wants:    nil,
		},
		{
			// Case 25: trailing slash stripped
			name:     "trailing_slash_stripped",
			docsJSON: `{"redirects": [{"source": "/docs/old/", "destination": "/docs/new/"}]}`,
			files:    map[string]string{"docs/new.mdx": "new"},
			wants:    nil,
		},
		{
			// Case 26: source with both bad-prefix path issues — both reported
			name:     "source_no_slash_and_destination_no_slash",
			docsJSON: `{"redirects": [{"source": "docs/x", "destination": "docs/y"}]}`,
			wants: []want{
				{substr: `source "docs/x": bad path: must start with '/'`},
				{substr: `destination "docs/y": bad path: must start with '/' (or http(s)://)`},
			},
		},
		{
			// Case 27: destination http (uppercase) — only / and lowercase http(s):// honored.
			name:     "destination_http_uppercase_rejected",
			docsJSON: `{"redirects": [{"source": "/docs/x", "destination": "HTTPS://example.com"}]}`,
			wants: []want{
				{substr: `must start with '/' (or http(s)://)`},
			},
		},
		{
			// Case 28: source resolves to dir (not regular file) — pageExists false
			name:     "source_resolves_to_dir",
			docsJSON: `{"redirects": [{"source": "/docs/old", "destination": "/docs/new"}]}`,
			files: map[string]string{
				"docs/old.mdx/": "",
				"docs/new.mdx":  "new",
			},
			wants: nil,
		},
		{
			// Case 29: wildcard source — first segment is wildcard (empty prefix)
			name:     "wildcard_source_empty_prefix",
			docsJSON: `{"redirects": [{"source": "/docs/:slug*", "destination": "/docs/new"}]}`,
			files: map[string]string{
				"docs/new.mdx": "new",
			},
			wants: []want{
				// docs/ exists as a dir and contains pages
				{substr: `dead redirect (wildcard source prefix has real pages)`},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{}
			for k, v := range tc.files {
				files[k] = v
			}
			root := setupContentRoot(t, files)
			vs := validate([]byte(tc.docsJSON), root)

			if len(tc.wants) == 0 {
				if len(vs) != 0 {
					t.Errorf("expected no violations, got %d: %v", len(vs), vs)
				}
				return
			}
			if len(vs) != len(tc.wants) {
				t.Fatalf("expected %d violations, got %d: %v", len(tc.wants), len(vs), vs)
			}
			for i, w := range tc.wants {
				assertViolation(t, vs[i], w.line, w.substr)
			}
		})
	}

	// Case 30: line numbers across non-object interleaving.
	t.Run("line_numbers_interleaved", func(t *testing.T) {
		docsJSON := "{\n" +
			"  \"redirects\": [\n" +
			"    {\"source\": \"/docs/a\", \"destination\": \"/docs/missing1\"},\n" +
			"    \"bad-string-entry\",\n" +
			"    {\"source\": \"/docs/b\", \"destination\": \"/docs/missing2\"}\n" +
			"  ]\n" +
			"}\n"
		// Lines (1-based):
		// 1: {
		// 2:   "redirects": [
		// 3:     {"source":"/docs/a", "destination":"/docs/missing1"},
		// 4:     "bad-string-entry",
		// 5:     {"source":"/docs/b", "destination":"/docs/missing2"}
		root := t.TempDir()
		vs := validate([]byte(docsJSON), root)
		if len(vs) != 3 {
			t.Fatalf("expected 3 violations, got %d: %v", len(vs), vs)
		}
		// Expect lines 3, 4, 5.
		assertViolation(t, vs[0], 3, "missing destination")
		assertViolation(t, vs[1], 4, "not an object")
		assertViolation(t, vs[2], 5, "missing destination")
	})
}

func TestCollectOpenAPISources(t *testing.T) {
	t.Run("nav_nested_maps", func(t *testing.T) {
		input := map[string]any{
			"nav": map[string]any{
				"openapi": map[string]any{"source": "docs/api/a.yaml"},
			},
		}
		got := collectOpenAPISources(input)
		want := map[string]bool{"docs/api/a.yaml": true}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("nav_array_of_maps", func(t *testing.T) {
		input := map[string]any{
			"nav": []any{
				map[string]any{"openapi": map[string]any{"source": "docs/api/b.yaml"}},
				map[string]any{"openapi": map[string]any{"source": "docs/api/c.yaml"}},
			},
		}
		got := collectOpenAPISources(input)
		want := map[string]bool{
			"docs/api/b.yaml": true,
			"docs/api/c.yaml": true,
		}
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
							"openapi": map[string]any{"source": "docs/api/d.yaml"},
						},
					},
				},
			},
			"openapi": map[string]any{"source": "docs/api/e.yaml"},
		}
		got := collectOpenAPISources(input)
		want := map[string]bool{
			"docs/api/d.yaml": true,
			"docs/api/e.yaml": true,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("openapi_non_map_defensive", func(t *testing.T) {
		input := map[string]any{
			"openapi": "scalar-not-a-map",
		}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("openapi_source_non_string_defensive", func(t *testing.T) {
		input := map[string]any{
			"openapi": map[string]any{"source": 42},
		}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("no_nav_at_all", func(t *testing.T) {
		input := map[string]any{
			"name": "x",
		}
		got := collectOpenAPISources(input)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("nav_is_scalar", func(t *testing.T) {
		input := map[string]any{
			"nav": "scalar",
		}
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
		{"/docs/x", "docs/x"},
		{"/docs/x/", "docs/x"},
		{"/docs/x?q=1", "docs/x"},
		{"/docs/x#frag", "docs/x"},
		{"/docs/x?q=1#frag", "docs/x"},
		{"/", ""},
		{"docs/x", "docs/x"},
		{"/docs/x/?q=1", "docs/x"},
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
		{"no_wildcard", "docs/x/y", []string{"docs", "x", "y"}, false},
		{"trailing_wildcard", "docs/x/:slug*", []string{"docs", "x"}, true},
		{"trailing_wildcard_no_star", "docs/x/:slug", []string{"docs", "x"}, true},
		{"wildcard_first_segment", "docs/:slug*", []string{"docs"}, true},
		{"empty", "", nil, false},
		{"only_wildcard", ":slug*", nil, true},
		{"empty_segments_skipped", "docs//x", []string{"docs", "x"}, false},
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
	got := formatMessage(2, "source", "/docs/x", "bad prefix")
	want := `redirects[2] source "/docs/x": bad prefix`
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
		{3, 1},  // the '\n' itself; line increments after
		{4, 2},  // first byte of 'def'
		{7, 2},  // second '\n'
		{8, 3},  // first byte of 'ghi'
		{11, 3}, // past end
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
		root := setupContentRoot(t, map[string]string{
			"a.md": "x",
		})
		if !dirHasPages(root) {
			t.Error("expected true for dir with .md")
		}
	})

	t.Run("non_page_files_only", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{
			"a.txt": "x",
		})
		if dirHasPages(root) {
			t.Error("expected false for dir with .txt only")
		}
	})

	t.Run("subdir_no_pages", func(t *testing.T) {
		// "0empty/" sorts before "z.mdx" so the empty-subdir branch is
		// iterated first and the `continue` after a false dirHasPages is
		// exercised before the .mdx sibling causes return true.
		root := setupContentRoot(t, map[string]string{
			"0empty/": "",
			"z.mdx":   "x",
		})
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

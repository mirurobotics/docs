package redirects

import "testing"

func TestValidate_TraversalSegments(t *testing.T) {
	runValidateCases(t, []validateCase{
		{
			name: "source_with_dot_dot_segment",
			docsJSON: `{"redirects": [{"source": "/../README",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"y.mdx": "y"},
			wants: []wantViolation{
				{substr: `source "/../README":` +
					` bad path: contains '..' or '.' segment`},
			},
		},
		{
			name: "destination_with_dot_dot_segment",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "/x/../../escape"}]}`,
			wants: []wantViolation{
				{substr: `destination "/x/../../escape":` +
					` bad path: contains '..' or '.' segment`},
			},
		},
		{
			name: "source_with_dot_segment",
			docsJSON: `{"redirects": [{"source": "/./foo",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"y.mdx": "y"},
			wants: []wantViolation{
				{substr: `source "/./foo":` +
					` bad path: contains '..' or '.' segment`},
			},
		},
	})
}

func TestValidate_Sources(t *testing.T) {
	runValidateCases(t, []validateCase{
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
			wants: []wantViolation{
				{substr: `dead redirect (source resolves to a real page)`},
			},
		},
		{
			name: "source_dead_redirect_md",
			docsJSON: `{"redirects": [{"source": "/dead",` +
				` "destination": "/y"}]}`,
			files: map[string]string{"dead.md": "dead", "y.mdx": "y"},
			wants: []wantViolation{
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
			name: "source_resolves_to_dir",
			docsJSON: `{"redirects": [{"source": "/old",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"old.mdx/": "", "new.mdx": "new"},
			wants: nil,
		},
	})
}

func TestValidate_WildcardSources(t *testing.T) {
	runValidateCases(t, []validateCase{
		{
			name: "wildcard_source_prefix_real_page",
			docsJSON: `{"redirects": [{"source": "/old/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"old.mdx": "old", "new.mdx": "new"},
			wants: []wantViolation{
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
			name: "wildcard_source_dir_has_pages",
			docsJSON: `{"redirects": [{"source": "/old/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{
				"old/sub/page.mdx": "page",
				"new.mdx":          "new",
			},
			wants: []wantViolation{
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
			name: "wildcard_source_empty_prefix",
			docsJSON: `{"redirects": [{"source": "/:slug*",` +
				` "destination": "/new"}]}`,
			files: map[string]string{"new.mdx": "new"},
			wants: []wantViolation{
				{substr: `dead redirect (wildcard source prefix has real pages)`},
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
			wants: []wantViolation{
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
	})
}

func TestValidate_Destinations(t *testing.T) {
	runValidateCases(t, []validateCase{
		{
			name: "destination_missing_page",
			docsJSON: `{"redirects": [{"source": "/old",` +
				` "destination": "/missing"}]}`,
			wants: []wantViolation{
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
			name:     "destination_root_relative_missing",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "/api/y"}]}`,
			wants:    []wantViolation{{substr: `missing destination`}},
		},
		{
			name:     "destination_http_ok",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "https://example.com/x"}]}`,
			wants:    nil,
		},
		{
			name:     "destination_exactly_docs",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "/docs"}]}`,
			wants:    []wantViolation{{substr: `missing destination`}},
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
			name: "old_docs_specific_page_to_root",
			docsJSON: `{"redirects": [{"source": "/docs/references/cli/install",` +
				` "destination": "/developers/cli/install"}]}`,
			files: map[string]string{"developers/cli/install.mdx": "install"},
			wants: nil,
		},
	})
}

func TestValidate_WildcardDestinations(t *testing.T) {
	runValidateCases(t, []validateCase{
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
			wants: []wantViolation{{substr: `wildcard prefix not a directory`}},
		},
		{
			name: "destination_wildcard_not_dir_no_yaml",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "/missing/:slug*"}]}`,
			wants: []wantViolation{{substr: `wildcard prefix not a directory`}},
		},
		{
			name: "old_docs_wildcard_to_root",
			docsJSON: `{"redirects": [{"source": "/docs/:slug*",` +
				` "destination": "/:slug*"}]}`,
			wants: nil,
		},
	})
}

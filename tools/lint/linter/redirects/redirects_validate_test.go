package redirects

import "testing"

func TestValidate_Document(t *testing.T) {
	runValidateCases(t, []validateCase{
		{
			name:     "top_level_not_object",
			docsJSON: `[1,2,3]`,
			wants:    []wantViolation{{line: 1, substr: "invalid JSON:"}},
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
			wants:    []wantViolation{{substr: `redirects[0] entry "": not an object`}},
		},
	})
}

func TestValidate_RequiredFields(t *testing.T) {
	runValidateCases(t, []validateCase{
		{
			name:     "missing_source",
			docsJSON: `{"redirects": [{"destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants: []wantViolation{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			name:     "missing_destination",
			docsJSON: `{"redirects": [{"source": "/x"}]}`,
			wants: []wantViolation{
				{substr: `redirects[0] destination "": must be a non-empty string`},
			},
		},
		{
			name:     "source_non_string_scalar",
			docsJSON: `{"redirects": [{"source": 5, "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants: []wantViolation{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
		{
			name:     "source_empty_string",
			docsJSON: `{"redirects": [{"source": "", "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants: []wantViolation{
				{substr: `redirects[0] source "": must be a non-empty string`},
			},
		},
	})
}

func TestValidate_PathFormat(t *testing.T) {
	runValidateCases(t, []validateCase{
		{
			name:     "source_no_leading_slash",
			docsJSON: `{"redirects": [{"source": "x", "destination": "/y"}]}`,
			files:    map[string]string{"y.mdx": "y"},
			wants:    []wantViolation{{substr: `must start with '/'`}},
		},
		{
			name:     "destination_no_leading_slash",
			docsJSON: `{"redirects": [{"source": "/x", "destination": "y"}]}`,
			wants:    []wantViolation{{substr: `must start with '/' (or http(s)://)`}},
		},
		{
			name: "source_no_slash_and_destination_no_slash",
			docsJSON: `{"redirects": [{"source": "x",` +
				` "destination": "y"}]}`,
			wants: []wantViolation{
				{substr: `source "x": bad path: must start with '/'`},
				{substr: `destination "y":` +
					` bad path: must start with '/' (or http(s)://)`},
			},
		},
		{
			name: "destination_http_uppercase_rejected",
			docsJSON: `{"redirects": [{"source": "/x",` +
				` "destination": "HTTPS://example.com"}]}`,
			wants: []wantViolation{{substr: `must start with '/' (or http(s)://)`}},
		},
	})
}

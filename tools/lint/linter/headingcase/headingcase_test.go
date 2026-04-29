package headingcase

import (
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// build returns (lines, spans) for a multi-line file content string,
// running the same Scanner the linter uses in production.
func build(content string) ([]string, [][]analysis.ProseSpan) {
	lines := strings.Split(content, "\n")
	scanner := analysis.NewScanner()
	spans := make([][]analysis.ProseSpan, len(lines))
	for i, line := range lines {
		spans[i] = scanner.ScanLine(line)
	}
	return lines, spans
}

func TestCheck_Headings(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
		wantLine  int
		wantCol   int
	}{
		{
			name:      "clean sentence-case",
			content:   "## Configure deployments\n",
			wantCount: 0,
		},
		{
			name:      "bad title-case",
			content:   "## Configure Deployments\n",
			wantCount: 1,
			wantLine:  1,
			wantCol:   4,
		},
		{
			name:      "bad acronym (v1 limitation)",
			content:   "## API reference\n",
			wantCount: 1,
			wantLine:  1,
			wantCol:   4,
		},
		{
			name:      "clean apostrophe",
			content:   "## Don't be afraid\n",
			wantCount: 0,
		},
		{
			name:      "bad apostrophe + title-case",
			content:   "## Don't Be Afraid\n",
			wantCount: 1,
			wantLine:  1,
			wantCol:   4,
		},
		{
			name:      "clean trailing question mark",
			content:   "### What is a config?\n",
			wantCount: 0,
		},
		{
			name:      "clean inline-code mask",
			content:   "## The `--version` flag\n",
			wantCount: 0,
		},
		{
			name:      "clean markdown link",
			content:   "## [Learn more](/x)\n",
			wantCount: 0,
		},
		{
			name:      "bad markdown link",
			content:   "## [Learn More](/x)\n",
			wantCount: 1,
			wantLine:  1,
			wantCol:   4,
		},
		{
			name:      "heading inside fenced code block",
			content:   "```\n## not a heading\n```\n",
			wantCount: 0,
		},
		{
			name:      "empty heading text after masking",
			content:   "## <Tooltip />\n",
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines, spans := build(tc.content)
			violations := Check("test.mdx", lines, spans)
			if len(violations) != tc.wantCount {
				t.Fatalf(
					"violations count = %d, want %d; violations=%+v",
					len(violations), tc.wantCount, violations,
				)
			}
			if tc.wantCount == 1 {
				v := violations[0]
				if v.Line != tc.wantLine {
					t.Errorf("Line = %d, want %d", v.Line, tc.wantLine)
				}
				if v.Col != tc.wantCol {
					t.Errorf("Col = %d, want %d", v.Col, tc.wantCol)
				}
				if v.Message != Message {
					t.Errorf("Message = %q, want %q", v.Message, Message)
				}
				if v.File != "test.mdx" {
					t.Errorf("File = %q, want %q", v.File, "test.mdx")
				}
			}
		})
	}
}

func TestCheck_FrontmatterTitle(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
		wantLine  int
		wantCol   int
	}{
		{
			name:      "quoted clean",
			content:   "---\ntitle: \"Workspace\"\n---\n",
			wantCount: 0,
		},
		{
			name:      "quoted bad",
			content:   "---\ntitle: \"User Management\"\n---\n",
			wantCount: 1,
			wantLine:  2,
			wantCol:   9,
		},
		{
			name:      "single-quoted clean",
			content:   "---\ntitle: 'Deployments'\n---\n",
			wantCount: 0,
		},
		{
			name:      "unquoted clean",
			content:   "---\ntitle: Deployments\n---\n",
			wantCount: 0,
		},
		{
			name:      "unquoted bad",
			content:   "---\ntitle: API Reference\n---\n",
			wantCount: 1,
			wantLine:  2,
			wantCol:   8,
		},
		{
			name:      "frontmatter without title line",
			content:   "---\nslug: foo\n---\n",
			wantCount: 0,
		},
		{
			name:      "no frontmatter",
			content:   "## Foo\n",
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines, spans := build(tc.content)
			violations := Check("test.mdx", lines, spans)
			if len(violations) != tc.wantCount {
				t.Fatalf(
					"violations count = %d, want %d; violations=%+v",
					len(violations), tc.wantCount, violations,
				)
			}
			if tc.wantCount == 1 {
				v := violations[0]
				if v.Line != tc.wantLine {
					t.Errorf("Line = %d, want %d", v.Line, tc.wantLine)
				}
				if v.Col != tc.wantCol {
					t.Errorf("Col = %d, want %d", v.Col, tc.wantCol)
				}
				if v.Message != Message {
					t.Errorf("Message = %q, want %q", v.Message, Message)
				}
				if v.File != "test.mdx" {
					t.Errorf("File = %q, want %q", v.File, "test.mdx")
				}
			}
		})
	}
}

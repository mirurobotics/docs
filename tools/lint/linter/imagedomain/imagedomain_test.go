package imagedomain

import (
	"strings"
	"testing"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
		wantLines []int
		wantCols  []int
		wantMsg   string
	}{
		{
			name:      "markdown assets ok",
			content:   "![d](https://assets.mirurobotics.com/docs/a.png)",
			wantCount: 0,
		},
		{
			name:      "markdown local path",
			content:   "![d](/images/a.png)",
			wantCount: 1,
			wantCols:  []int{6},
			wantMsg: "image-domain: image must be hosted on " +
				"https://assets.mirurobotics.com (got \"/images/a.png\")",
		},
		{
			name:      "markdown other domain",
			content:   "![d](https://example.com/a.png)",
			wantCount: 1,
		},
		{
			name:      "markdown relative with title strips title",
			content:   `![d](./a.webp "title")`,
			wantCount: 1,
			wantMsg: "image-domain: image must be hosted on " +
				"https://assets.mirurobotics.com (got \"./a.webp\")",
		},

		{
			name:      "img src assets ok",
			content:   `<img src="https://assets.mirurobotics.com/docs/a.png" />`,
			wantCount: 0,
		},
		{
			name:      "src local path",
			content:   `src="/images/a.png"`,
			wantCount: 1,
			wantCols:  []int{6},
		},
		{
			name:      "src http not https",
			content:   `src="http://assets.mirurobotics.com/a.png"`,
			wantCount: 1,
		},
		{
			name:      "src protocol relative",
			content:   `src="//assets.mirurobotics.com/a.png"`,
			wantCount: 1,
		},
		{
			name:      "src assets mp4 ok",
			content:   `src="https://assets.mirurobotics.com/docs/v.mp4"`,
			wantCount: 0,
		},
		{
			name:      "src local mp4 ok non image ext",
			content:   `src="/videos/v.mp4"`,
			wantCount: 0,
		},
		{
			name:      "src local no extension ok",
			content:   `src="/nofile"`,
			wantCount: 0,
		},
		{
			name:      "src image with query still flagged",
			content:   `src="/images/a.png?v=2"`,
			wantCount: 1,
		},

		{
			name:      "image local png",
			content:   `image="/images/changelog/x.png"`,
			wantCount: 1,
		},
		{
			name:      "image no extension strict",
			content:   `image="/images/x"`,
			wantCount: 1,
		},

		{
			name:      "background assets dark svg ok",
			content:   `background="https://assets.mirurobotics.com/docs/bg.dark.svg"`,
			wantCount: 0,
		},

		{
			name:      "poster local jpg",
			content:   `poster="/images/p.jpg"`,
			wantCount: 1,
		},

		{
			name:      "src jsx brace local",
			content:   `src={"/images/a.png"}`,
			wantCount: 1,
		},
		{
			name:      "src jsx brace assets ok",
			content:   `src={"https://assets.mirurobotics.com/docs/a.png"}`,
			wantCount: 0,
		},

		{
			name:      "image colon in filename ok",
			content:   `image="https://assets.mirurobotics.com/docs/releases/header:page.png"`,
			wantCount: 0,
		},

		{
			name:      "inline backticks still flagged",
			content:   "`![d](/images/a.png)`",
			wantCount: 1,
		},

		{
			name:      "multiple candidates flagged",
			content:   `<Framed image="/a.png" background="/b.png" />`,
			wantCount: 2,
		},
		{
			name:      "mixed verdicts same line",
			content:   `<img src="/x.mp4" image="/y.png" />`,
			wantCount: 1,
		},

		{
			name:      "frontmatter attribute skipped",
			content:   "---\nsrc=\"/images/a.png\"\n---\n\nClean",
			wantCount: 0,
		},
		{
			name:      "body violation after frontmatter",
			content:   "---\ntitle: Test\n---\n\n![d](/images/a.png)",
			wantCount: 1,
			wantLines: []int{5},
		},

		{
			name:      "inside fence skipped",
			content:   "```\n![d](/images/a.png)\n```",
			wantCount: 0,
		},
		{
			name:      "flagged after fence closes",
			content:   "```\n![d](/images/a.png)\n```\n![d](/images/b.png)",
			wantCount: 1,
			wantLines: []int{4},
		},
		{
			name:      "four backtick fence skipped",
			content:   "````\n![d](/images/a.png)\n````",
			wantCount: 0,
		},

		{
			name:      "suppression skips next line",
			content:   "{/* lint-ignore image-domain */}\n![d](/images/a.png)",
			wantCount: 0,
		},
		{
			name: "suppression covers only the next line",
			content: "{/* lint-ignore image-domain */}\n" +
				"![d](/images/a.png)\n![d](/images/b.png)",
			wantCount: 1,
			wantLines: []int{3},
		},
		{
			name: "suppression covers all violations on suppressed line",
			content: "{/* lint-ignore image-domain */}\n" +
				`<Framed image="/a.png" background="/b.png" />`,
			wantCount: 0,
		},
		{
			name:      "indented suppression directive recognized",
			content:   "  {/* lint-ignore image-domain */}  \n![d](/images/a.png)",
			wantCount: 0,
		},

		{name: "empty string", content: "", wantCount: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.content, "\n")
			violations := Check("test.mdx", lines)
			if len(violations) != tc.wantCount {
				t.Fatalf(
					"expected %d violations, got %d: %v",
					tc.wantCount, len(violations), violations,
				)
			}
			assertPositions(t, violations, tc.wantLines, tc.wantCols)
			if tc.wantMsg != "" && violations[0].Message != tc.wantMsg {
				t.Errorf(
					"expected message %q, got %q",
					tc.wantMsg, violations[0].Message,
				)
			}
			if tc.wantCount > 0 && violations[0].File != "test.mdx" {
				t.Errorf("expected file test.mdx, got %q", violations[0].File)
			}
		})
	}
}

func assertPositions(
	t *testing.T,
	violations []analysis.Violation,
	wantLines, wantCols []int,
) {
	t.Helper()
	for i, wantLine := range wantLines {
		if violations[i].Line != wantLine {
			t.Errorf(
				"violation %d: expected line %d, got %d",
				i, wantLine, violations[i].Line,
			)
		}
	}
	for i, wantCol := range wantCols {
		if violations[i].Col != wantCol {
			t.Errorf(
				"violation %d: expected col %d, got %d",
				i, wantCol, violations[i].Col,
			)
		}
	}
}

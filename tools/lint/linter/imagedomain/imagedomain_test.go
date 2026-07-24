package imagedomain

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
		wantLines []int  // 1-based lines, checked positionally when set
		wantCols  []int  // 1-based columns, checked positionally when set
		wantMsg   string // if set, compared against violations[0].Message
	}{
		// Markdown images.
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

		// img/src attributes.
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

		// image attribute (strict, any URL off-domain flagged).
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

		// background attribute.
		{
			name:      "background assets dark svg ok",
			content:   `background="https://assets.mirurobotics.com/docs/bg.dark.svg"`,
			wantCount: 0,
		},

		// poster attribute.
		{
			name:      "poster local jpg",
			content:   `poster="/images/p.jpg"`,
			wantCount: 1,
		},

		// JSX brace-quoted values.
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

		// Colon in filename tolerated.
		{
			name:      "image colon in filename ok",
			content:   `image="https://assets.mirurobotics.com/docs/releases/header:page.png"`,
			wantCount: 0,
		},

		// Accepted false positive: violating URL inside inline backticks.
		{
			name:      "inline backticks still flagged",
			content:   "`![d](/images/a.png)`",
			wantCount: 1,
		},

		// Multiple candidates on one line, each flagged.
		{
			name:      "multiple candidates flagged",
			content:   `<Framed image="/a.png" background="/b.png" />`,
			wantCount: 2,
		},
		{
			name:      "mixed verdicts same line",
			content:   `<img src="/x.mp4" image="/y.png" />`,
			wantCount: 1, // src mp4 ok, image png violation
		},

		// Frontmatter skipping.
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

		// Code-fence toggling.
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

		// Suppression directive (next line only).
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

		// Empty input.
		{name: "empty string", content: "", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := strings.Split(tt.content, "\n")
			violations := Check("test.mdx", lines)
			if len(violations) != tt.wantCount {
				t.Fatalf(
					"expected %d violations, got %d: %v",
					tt.wantCount, len(violations), violations,
				)
			}
			for i, wantLine := range tt.wantLines {
				if violations[i].Line != wantLine {
					t.Errorf(
						"violation %d: expected line %d, got %d",
						i, wantLine, violations[i].Line,
					)
				}
			}
			for i, wantCol := range tt.wantCols {
				if violations[i].Col != wantCol {
					t.Errorf(
						"violation %d: expected col %d, got %d",
						i, wantCol, violations[i].Col,
					)
				}
			}
			if tt.wantMsg != "" && violations[0].Message != tt.wantMsg {
				t.Errorf(
					"expected message %q, got %q",
					tt.wantMsg, violations[0].Message,
				)
			}
			if tt.wantCount > 0 && violations[0].File != "test.mdx" {
				t.Errorf("expected file test.mdx, got %q", violations[0].File)
			}
		})
	}
}

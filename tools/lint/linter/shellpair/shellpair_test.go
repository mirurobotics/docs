package shellpair

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantLines []int
		wantMsg   string // if set, compared against the first violation
	}{
		{
			name: "matching pair",
			content: "<CodeGroup>\n" +
				"```bash Unix highlight={3}\n" +
				"miru create \\\n  --version v1 \\\n  --schemas ./s/\n```\n" +
				"```powershell PowerShell highlight={3}\n" +
				"miru create `\n  --version v1 `\n  --schemas ./s/\n```\n" +
				"</CodeGroup>",
		},
		{
			name: "pair differs",
			content: "<CodeGroup>\n" +
				"```bash Unix\nmiru create \\\n  --version v1\n```\n" +
				"```powershell PowerShell\nmiru create `\n  --version v2\n```\n" +
				"</CodeGroup>",
			wantLines: []int{6},
			wantMsg:   "shell-pair: PowerShell command differs from its Unix twin",
		},
		{
			name: "PowerShell twin still uses backslashes",
			content: "<CodeGroup>\n" +
				"```bash Unix\nmiru create \\\n  --version v1\n```\n" +
				"```powershell PowerShell\nmiru create \\\n  --version v1\n```\n" +
				"</CodeGroup>",
			wantLines: []int{6, 6},
			wantMsg: "shell-pair: powershell block continues a line " +
				"with a Unix backslash",
		},
		{
			name: "multi-line Unix without twin",
			content: "  <CodeGroup>\n" +
				"  ```bash Unix\n  miru create \\\n    --version v1\n  ```\n" +
				"  </CodeGroup>",
			wantLines: []int{2},
			wantMsg:   "shell-pair: multi-line Unix command needs a PowerShell twin",
		},
		{
			name: "multi-line PowerShell without twin",
			content: "<CodeGroup>\n" +
				"```powershell PowerShell\nmiru create `\n  --version v1\n```\n" +
				"</CodeGroup>",
			wantLines: []int{2},
			wantMsg:   "shell-pair: multi-line PowerShell command needs a Unix twin",
		},
		{
			name: "single-line commands may differ",
			content: "<CodeGroup>\n" +
				"```bash Unix\nexport MIRU_API_KEY=<key>\n```\n" +
				"```powershell PowerShell\n$env:MIRU_API_KEY = \"<key>\"\n```\n" +
				"</CodeGroup>",
		},
		{
			name: "other titles are ignored",
			content: "<CodeGroup>\n" +
				"```bash Empty Schemas (Unix)\nmiru create \\\n  --version v1\n```\n" +
				"```bash\nmiru create \\\n  --version v1\n```\n" +
				"</CodeGroup>",
		},
		{
			name:    "blocks outside a CodeGroup are ignored",
			content: "```bash Unix\nmiru create \\\n  --version v1\n```\n</CodeGroup>",
		},
		{
			name:      "powershell block with a backslash continuation",
			content:   "```powershell\nmiru create \\\n  --version v1\n```",
			wantLines: []int{1},
		},
		{
			name:      "bash block with a backtick continuation",
			content:   "```bash\nmiru create `\n  --version v1\n```",
			wantLines: []int{1},
			wantMsg: "shell-pair: bash block continues a line " +
				"with PowerShell's backtick",
		},
		{
			name: "each CodeGroup is checked on its own",
			content: "<CodeGroup>\n" +
				"```bash Unix\nmiru a \\\n  --b\n```\n" +
				"```powershell PowerShell\nmiru a `\n  --b\n```\n" +
				"</CodeGroup>\n" +
				"<CodeGroup>\n" +
				"```bash Unix\nmiru c \\\n  --d\n```\n" +
				"</CodeGroup>",
			wantLines: []int{12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := Check("test.mdx", strings.Split(tt.content, "\n"))
			if len(violations) != len(tt.wantLines) {
				t.Fatalf("expected %d violations, got %d: %v",
					len(tt.wantLines), len(violations), violations)
			}
			for i, wantLine := range tt.wantLines {
				if violations[i].Line != wantLine {
					t.Errorf("violation %d: expected line %d, got %d",
						i, wantLine, violations[i].Line)
				}
				if violations[i].File != "test.mdx" || violations[i].Col != 1 {
					t.Errorf("violation %d: unexpected position %+v", i, violations[i])
				}
			}
			if tt.wantMsg != "" && violations[0].Message != tt.wantMsg {
				t.Errorf("expected message %q, got %q",
					tt.wantMsg, violations[0].Message)
			}
		})
	}
}

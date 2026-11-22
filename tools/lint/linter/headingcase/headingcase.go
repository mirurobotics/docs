package headingcase

import (
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// Message is the diagnostic emitted for any heading-case violation.
const Message = "heading-case: heading must be sentence-case " +
	"(first letter uppercase, all other letters lowercase)"

const trailingPunct = ".?!:"

const outerPunct = ",;()[]{}\"'"

func titleRe() *regexp.Regexp { return regexp.MustCompile(`^title:\s*(.*)$`) }

func headingRe() *regexp.Regexp {
	return regexp.MustCompile(`^(#{1,6})[ \t]+(.+?)[ \t]*$`)
}

func versionTagRe() *regexp.Regexp { return regexp.MustCompile(`^v?\d+([.-]\w+)*$`) }

func inlineCodeRe() *regexp.Regexp { return regexp.MustCompile("`[^`]*`") }

func htmlTagRe() *regexp.Regexp { return regexp.MustCompile(`<[^>]*>`) }

func mdLinkRe() *regexp.Regexp { return regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`) }

func mdxCommentRe() *regexp.Regexp { return regexp.MustCompile(`(?s)\{/\*.*?\*/\}`) }

func allowlist() map[string]struct{} {
	return map[string]struct{}{
		"API":     {},
		"APIs":    {},
		"CLI":     {},
		"CI":      {},
		"SDK":     {},
		"SDKs":    {},
		"CUE":     {},
		"JSON":    {},
		"MQTT":    {},
		"TLS":     {},
		"HTTPS":   {},
		"REST":    {},
		"GUI":     {},
		"URL":     {},
		"ACLs":    {},
		"SSE":     {},
		"OpenAPI": {},
		"AWS":     {},
		"GCP":     {},
		"GCS":     {},
		"WIF":     {},
		"STS":     {},
		"IAM":     {},
		"S3":      {},
		"ARN":     {},
		"SigV4":   {},

		"Miru":    {},
		"GitHub":  {},
		"Agent":   {},
		"Unix":    {},
		"Git":     {},
		"Python":  {},
		"Schema":  {},
		"Base":    {},
		"Head":    {},
		"Cloud":   {},
		"Storage": {},

		"tetons": {},
		"zion":   {},

		"deployment.deployed": {},
		"deployment.removed":  {},
	}
}

type checker struct {
	title, heading, version *regexp.Regexp
	inlineCode, htmlTag     *regexp.Regexp
	mdLink, mdxComment      *regexp.Regexp
	allow                   map[string]struct{}
}

func newChecker() *checker {
	return &checker{
		title:      titleRe(),
		heading:    headingRe(),
		version:    versionTagRe(),
		inlineCode: inlineCodeRe(),
		htmlTag:    htmlTagRe(),
		mdLink:     mdLinkRe(),
		mdxComment: mdxCommentRe(),
		allow:      allowlist(),
	}
}

// Check flags a frontmatter title or prose heading that is not in sentence case.
func Check(
	file string,
	lines []string,
	spans [][]analysis.ProseSpan,
) []analysis.Violation {
	chk := newChecker()
	var violations []analysis.Violation

	if end := analysis.FrontmatterEnd(lines); end >= 1 {
		for i := 1; i < end; i++ {
			line := lines[i]
			if chk.title.FindStringSubmatch(line) == nil {
				continue
			}
			value, valueByteOffset := titleValueAndOffset(line)
			if chk.casingViolation(value) {
				violations = append(violations, analysis.Violation{
					File:    file,
					Line:    i + 1,
					Col:     valueByteOffset + 1,
					Message: Message,
				})
			}
			break
		}
	}

	for i, line := range lines {
		if i >= len(spans) || len(spans[i]) == 0 {
			continue
		}
		match := chk.heading.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		masked := chk.maskHeading(match[2])
		if !chk.casingViolation(masked) {
			continue
		}
		violations = append(violations, analysis.Violation{
			File:    file,
			Line:    i + 1,
			Col:     headingTextCol(line),
			Message: Message,
		})
	}

	return violations
}

func titleValueAndOffset(line string) (string, int) {
	idx := len("title:")
	for idx < len(line) && (line[idx] == ' ' || line[idx] == '\t') {
		idx++
	}
	value := line[idx:]
	if len(value) >= 2 {
		first := value[0]
		last := value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			value = value[1 : len(value)-1]
			idx++
		}
	}
	return value, idx
}

func headingTextCol(line string) int {
	i := 0
	for i < len(line) && line[i] == '#' {
		i++
	}
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	return i + 1
}

func (c *checker) maskHeading(text string) string {
	text = c.mdxComment.ReplaceAllString(text, "")
	text = c.inlineCode.ReplaceAllString(text, "")
	text = c.htmlTag.ReplaceAllString(text, "")
	text = c.mdLink.ReplaceAllString(text, "$1")
	return text
}

func stripOuterPunct(token string) string { return strings.Trim(token, outerPunct) }

func (c *checker) casingViolation(text string) bool {
	text = strings.TrimSpace(text)
	for len(text) > 0 && strings.ContainsRune(trailingPunct, rune(text[len(text)-1])) {
		text = text[:len(text)-1]
	}
	if text == "" {
		return false
	}
	if c.version.MatchString(text) {
		return false
	}
	tokens := strings.Fields(text)
	for i, tok := range tokens {
		core := stripOuterPunct(tok)
		if core == "" {
			continue
		}
		if c.version.MatchString(core) {
			continue
		}
		if _, ok := c.allow[core]; ok {
			continue
		}
		if i == 0 {
			if tokenViolatesLeading(core) {
				return true
			}
		} else if tokenViolatesLater(core) {
			return true
		}
	}
	return false
}

func tokenViolatesLeading(token string) bool {
	seenFirst := false
	for i := 0; i < len(token); i++ {
		char := token[i]
		isUpper := char >= 'A' && char <= 'Z'
		isLower := char >= 'a' && char <= 'z'
		if !isUpper && !isLower {
			continue
		}
		if !seenFirst {
			if !isUpper {
				return true
			}
			seenFirst = true
			continue
		}
		if isUpper {
			return true
		}
	}
	return false
}

func tokenViolatesLater(token string) bool {
	for i := 0; i < len(token); i++ {
		char := token[i]
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

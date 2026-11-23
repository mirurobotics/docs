package imagedomain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

const allowedPrefix = "https://assets.mirurobotics.com/"

const ignoreDirective = "{/* lint-ignore image-domain */}"

func mdImageRe() *regexp.Regexp {
	return regexp.MustCompile(`!\[[^\]]*\]\(\s*([^)\s]+)`)
}

func attrRe() *regexp.Regexp {
	return regexp.MustCompile(
		`\b(src|image|background|poster)\s*=\s*\{?\s*["']([^"']*)["']`)
}

func imageExts() map[string]struct{} {
	return map[string]struct{}{
		".png":  {},
		".jpg":  {},
		".jpeg": {},
		".gif":  {},
		".svg":  {},
		".webp": {},
	}
}

type candidate struct {
	url     string
	col     int
	srcOnly bool
}

type checker struct {
	mdImage, attr *regexp.Regexp
	exts          map[string]struct{}
}

func newChecker() *checker {
	return &checker{mdImage: mdImageRe(), attr: attrRe(), exts: imageExts()}
}

// Check flags image references whose URL is not on the Miru assets domain.
func Check(file string, lines []string) []analysis.Violation {
	chk := newChecker()
	fmEnd := analysis.FrontmatterEnd(lines)
	inFence := false
	var violations []analysis.Violation
	for i, line := range lines {
		if i <= fmEnd {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
			continue
		}
		if inFence || suppressed(lines, i) {
			continue
		}
		for _, cand := range chk.candidates(line) {
			if !chk.violates(cand) {
				continue
			}
			msg := fmt.Sprintf("image-domain: image must be hosted on "+
				"https://assets.mirurobotics.com (got %q)", cand.url)
			violations = append(violations, analysis.Violation{
				File: file, Line: i + 1, Col: cand.col, Message: msg,
			})
		}
	}
	return violations
}

func suppressed(lines []string, i int) bool {
	return i > 0 && strings.TrimSpace(lines[i-1]) == ignoreDirective
}

func (c *checker) candidates(line string) []candidate {
	var cands []candidate
	for _, match := range c.mdImage.FindAllStringSubmatchIndex(line, -1) {
		cands = append(cands, candidate{
			url: line[match[2]:match[3]], col: match[2] + 1,
		})
	}
	for _, match := range c.attr.FindAllStringSubmatchIndex(line, -1) {
		cands = append(cands, candidate{
			url:     line[match[4]:match[5]],
			col:     match[4] + 1,
			srcOnly: line[match[2]:match[3]] == "src",
		})
	}
	return cands
}

func (c *checker) violates(cand candidate) bool {
	if strings.HasPrefix(cand.url, allowedPrefix) {
		return false
	}
	if cand.srcOnly {
		return c.hasImageExt(cand.url)
	}
	return true
}

func (c *checker) hasImageExt(url string) bool {
	if i := strings.IndexAny(url, "?#"); i >= 0 {
		url = url[:i]
	}
	url = strings.ToLower(url)
	dot := strings.LastIndex(url, ".")
	if dot < 0 {
		return false
	}
	_, ok := c.exts[url[dot:]]
	return ok
}

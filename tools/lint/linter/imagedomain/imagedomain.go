// Package imagedomain enforces that every image referenced from MDX
// content is hosted on the Miru assets domain: the URL must start with
// "https://assets.mirurobotics.com/". It inspects Markdown images
// (![alt](url)) and image-bearing attributes (src, image, background,
// poster; plain or JSX brace-quoted). Relative/local paths, other
// domains, protocol-relative "//..." URLs, and "http://..." URLs are
// all violations. For src attributes only URLs with an image file
// extension (.png, .jpg, .jpeg, .gif, .svg, .webp) are checked, so
// videos and other non-image sources are never flagged. Lines inside
// fenced code blocks and YAML frontmatter are skipped.
//
// A finding is suppressed by placing the directive
//
//	{/* lint-ignore image-domain */}
//
// on its own on the line immediately before the offending line. The
// suppression applies to the next line only.
package imagedomain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

// allowedPrefix is the required prefix for every image URL.
const allowedPrefix = "https://assets.mirurobotics.com/"

// ignoreDirective suppresses findings on the immediately following
// line when it is the previous line's entire trimmed content.
const ignoreDirective = "{/* lint-ignore image-domain */}"

// The Re-suffixed helpers below return freshly-compiled regexes; this
// avoids package-level mutable state per the project's lint
// conventions. Callers should compile once per Check via newChecker
// rather than per line.

// mdImageRe matches a Markdown image and captures its URL (group 1).
// The URL stops at whitespace so optional titles are excluded.
func mdImageRe() *regexp.Regexp {
	return regexp.MustCompile(`!\[[^\]]*\]\(\s*([^)\s]+)`)
}

// attrRe matches an image-bearing attribute, capturing the attribute
// name (group 1) and the quoted URL (group 2). The optional "{"
// also matches JSX brace-quoted values like src={"https://..."}.
func attrRe() *regexp.Regexp {
	return regexp.MustCompile(
		`\b(src|image|background|poster)\s*=\s*\{?\s*["']([^"']*)["']`)
}

// imageExts returns the set of lowercase file extensions treated as
// images when deciding whether a src attribute is in scope. It is a
// function rather than a package-level variable to avoid mutable
// global state.
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

// candidate is a single image URL occurrence found on a line.
type candidate struct {
	url string
	col int // 1-based byte column of the URL's first character
	// srcOnly marks src-attribute candidates, which are flagged only
	// when the URL has an image file extension.
	srcOnly bool
}

// checker bundles compiled regexes and the extension set for a single
// Check run so they're built once instead of once per line.
type checker struct {
	mdImage, attr *regexp.Regexp
	exts          map[string]struct{}
}

func newChecker() *checker {
	return &checker{mdImage: mdImageRe(), attr: attrRe(), exts: imageExts()}
}

// Check flags image references whose URL is not hosted on
// https://assets.mirurobotics.com/. lines is the raw file content
// split by line.
func Check(file string, lines []string) []analysis.Violation {
	c := newChecker()
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
		for _, cand := range c.candidates(line) {
			if !c.violates(cand) {
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

// suppressed reports whether the previous line's trimmed content is
// exactly the lint-ignore directive for this rule.
func suppressed(lines []string, i int) bool {
	return i > 0 && strings.TrimSpace(lines[i-1]) == ignoreDirective
}

// candidates collects every image URL occurrence on the line, from
// both Markdown image syntax and image-bearing attributes.
func (c *checker) candidates(line string) []candidate {
	var cands []candidate
	for _, m := range c.mdImage.FindAllStringSubmatchIndex(line, -1) {
		cands = append(cands, candidate{url: line[m[2]:m[3]], col: m[2] + 1})
	}
	for _, m := range c.attr.FindAllStringSubmatchIndex(line, -1) {
		cands = append(cands, candidate{
			url:     line[m[4]:m[5]],
			col:     m[4] + 1,
			srcOnly: line[m[2]:m[3]] == "src",
		})
	}
	return cands
}

// violates reports whether the candidate's URL breaks the rule. src
// attributes are only flagged when the URL has an image extension;
// Markdown images and the image/background/poster attributes are
// flagged for any URL off the assets domain.
func (c *checker) violates(cand candidate) bool {
	if strings.HasPrefix(cand.url, allowedPrefix) {
		return false
	}
	if cand.srcOnly {
		return c.hasImageExt(cand.url)
	}
	return true
}

// hasImageExt reports whether the URL, with any query string or
// fragment stripped and lowercased, ends in an image file extension.
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

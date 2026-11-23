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

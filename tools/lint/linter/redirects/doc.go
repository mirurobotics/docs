// Package redirects validates the `redirects` array in docs.json against
// the on-disk Mintlify content tree.
//
// Mintlify serves foo/bar.mdx at URL /foo/bar relative to the configured
// content root. The `redirects` array in docs.json rewrites URLs at the
// edge. This rule catches:
//   - Dead redirects (source already serves a real page).
//   - Missing destinations (destination has no real page). A wildcard
//     destination may also target pages Mintlify generates from a yaml
//     registered as an openapi.source in docs.json.
//   - Unsupported schemes / malformed paths, including '.' and '..'
//     segments, which would resolve outside the content root.
package redirects

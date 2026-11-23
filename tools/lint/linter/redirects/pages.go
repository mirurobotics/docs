package redirects

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mirurobotics/docs/tools/lint/linter/analysis"
)

func validateFileSystem(
	fields entryFields,
	contentRoot string,
	openAPISources map[string]bool,
) []analysis.Violation {
	var violations []analysis.Violation
	if strings.HasPrefix(fields.source, "/") {
		violations = validateSource(fields, contentRoot, openAPISources)
	}
	if !destIsHTTP(fields.destination) && strings.HasPrefix(fields.destination, "/") {
		dst := validateDestination(fields, contentRoot, openAPISources)
		violations = append(violations, dst...)
	}
	return violations
}

func validateSource(
	fields entryFields,
	contentRoot string,
	openAPISources map[string]bool,
) []analysis.Violation {
	src := resolvePath(fields.source, contentRoot)
	var msg string
	switch {
	case src.hasTraversal():
		msg = "bad path: contains '..' or '.' segment"
	case !src.hasWildcard && src.hasPage():
		msg = "dead redirect (source resolves to a real page)"
	case !src.hasWildcard:
		return nil
	case src.hasPage():
		msg = "dead redirect (wildcard source prefix resolves to a real page)"
	case dirExists(src.fsPath) && dirHasPages(src.fsPath):
		msg = "dead redirect (wildcard source prefix has real pages)"
	case src.hasGeneratedPages(openAPISources):
		msg = "dead redirect (wildcard source prefix has Mintlify-generated pages)"
	default:
		return nil
	}
	return entryViolation(fields.i, fields.line, "source", fields.source, msg)
}

func dirHasPages(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() && dirHasPages(filepath.Join(dir, name)) {
			return true
		}
		if !ent.IsDir() && isPageFile(name) {
			return true
		}
	}
	return false
}

func isPageFile(name string) bool {
	return strings.HasSuffix(name, ".mdx") || strings.HasSuffix(name, ".md")
}

func validateDestination(
	fields entryFields,
	contentRoot string,
	openAPISources map[string]bool,
) []analysis.Violation {
	dst := resolvePath(fields.destination, contentRoot)
	var msg string
	switch {
	case dst.hasTraversal():
		msg = "bad path: contains '..' or '.' segment"
	case !dst.hasWildcard && dst.hasPage():
		return nil
	case !dst.hasWildcard:
		msg = "missing destination (no .mdx or .md page exists)"
	case dirExists(dst.fsPath):
		return nil
	case dst.hasGeneratedPages(openAPISources):
		return nil
	default:
		msg = "wildcard prefix not a directory"
	}
	return entryViolation(fields.i, fields.line, "destination", fields.destination, msg)
}

func resolvePath(raw, contentRoot string) resolvedPath {
	prefix, hasWildcard := splitWildcard(cleanPath(raw))
	return resolvedPath{
		prefix:      prefix,
		hasWildcard: hasWildcard,
		fsPath:      filepath.Join(contentRoot, filepath.Join(prefix...)),
	}
}

func splitWildcard(cleaned string) ([]string, bool) {
	rawSegments := strings.Split(cleaned, "/")
	var segments []string
	for _, seg := range rawSegments {
		if seg != "" {
			segments = append(segments, seg)
		}
	}
	var prefix []string
	hasWildcard := false
	for _, seg := range segments {
		if wildcardSegment().MatchString(seg) {
			hasWildcard = true
			break
		}
		prefix = append(prefix, seg)
	}
	return prefix, hasWildcard
}

func wildcardSegment() *regexp.Regexp {
	return regexp.MustCompile(`^:[A-Za-z][A-Za-z0-9]*\*?$`)
}

func cleanPath(raw string) string {
	cleaned := raw
	if idx := strings.Index(cleaned, "?"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	if idx := strings.Index(cleaned, "#"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	cleaned = strings.TrimPrefix(cleaned, "/")
	cleaned = strings.TrimSuffix(cleaned, "/")
	return cleaned
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

type resolvedPath struct {
	prefix      []string
	hasWildcard bool
	fsPath      string
}

func (r resolvedPath) hasTraversal() bool { return containsTraversalSegment(r.prefix) }

func (r resolvedPath) hasPage() bool { return pageExists(r.fsPath) }

func (r resolvedPath) hasGeneratedPages(openAPISources map[string]bool) bool {
	yamlRel := strings.Join(r.prefix, "/") + ".yaml"
	return openAPISources[yamlRel] && fileExists(r.fsPath+".yaml")
}

func containsTraversalSegment(segments []string) bool {
	for _, seg := range segments {
		if seg == ".." || seg == "." {
			return true
		}
	}
	return false
}

func pageExists(prefixFs string) bool {
	return fileExists(prefixFs+".mdx") || fileExists(prefixFs+".md")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

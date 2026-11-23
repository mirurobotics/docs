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
	prefix, hasWildcard := splitWildcard(cleanPath(fields.source))
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))
	yamlRel := strings.Join(prefix, "/") + ".yaml"
	var msg string
	switch {
	case containsTraversalSegment(prefix):
		msg = "bad path: contains '..' or '.' segment"
	case !hasWildcard && pageExists(prefixFs):
		msg = "dead redirect (source resolves to a real page)"
	case !hasWildcard:
		return nil
	case pageExists(prefixFs):
		msg = "dead redirect (wildcard source prefix resolves to a real page)"
	case dirExists(prefixFs) && dirHasPages(prefixFs):
		msg = "dead redirect (wildcard source prefix has real pages)"
	case openAPISources[yamlRel] && fileExists(prefixFs+".yaml"):
		msg = "dead redirect (wildcard source prefix has Mintlify-generated pages)"
	default:
		return nil
	}
	return entryViolation(fields.i, fields.line, "source", fields.source, msg)
}

func validateDestination(
	fields entryFields,
	contentRoot string,
	openAPISources map[string]bool,
) []analysis.Violation {
	prefix, hasWildcard := splitWildcard(cleanPath(fields.destination))
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))
	yamlRel := strings.Join(prefix, "/") + ".yaml"
	var msg string
	switch {
	case containsTraversalSegment(prefix):
		msg = "bad path: contains '..' or '.' segment"
	case !hasWildcard && pageExists(prefixFs):
		return nil
	case !hasWildcard:
		msg = "missing destination (no .mdx or .md page exists)"
	case dirExists(prefixFs):
		return nil
	case openAPISources[yamlRel] && fileExists(prefixFs+".yaml"):
		return nil
	default:
		msg = "wildcard prefix not a directory"
	}
	return entryViolation(fields.i, fields.line, "destination", fields.destination, msg)
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

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
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

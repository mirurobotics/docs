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
		srcPrefix, _ := splitWildcard(cleanPath(fields.source))
		if containsTraversalSegment(srcPrefix) {
			violations = append(violations, analysis.Violation{
				File: "docs.json",
				Line: fields.line,
				Col:  1,
				Message: formatMessage(
					fields.i, "source", fields.source,
					"bad path: contains '..' or '.' segment",
				),
			})
		} else {
			srcViolations := validateSource(
				fields.i, fields.source, contentRoot, openAPISources, fields.line,
			)
			violations = append(violations, srcViolations...)
		}
	}
	if !destIsHTTP(fields.destination) && strings.HasPrefix(fields.destination, "/") {
		dstPrefix, _ := splitWildcard(cleanPath(fields.destination))
		if containsTraversalSegment(dstPrefix) {
			violations = append(violations, analysis.Violation{
				File: "docs.json",
				Line: fields.line,
				Col:  1,
				Message: formatMessage(
					fields.i, "destination", fields.destination,
					"bad path: contains '..' or '.' segment",
				),
			})
		} else {
			dstViolations := validateDestination(
				fields.i, fields.destination, contentRoot, openAPISources, fields.line,
			)
			violations = append(violations, dstViolations...)
		}
	}
	return violations
}

func validateSource(
	i int,
	source, contentRoot string,
	openAPISources map[string]bool,
	line int,
) []analysis.Violation {
	cleaned := cleanPath(source)
	prefix, hasWildcard := splitWildcard(cleaned)
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))

	if !hasWildcard {
		if pageExists(prefixFs) {
			return []analysis.Violation{{
				File: "docs.json",
				Line: line,
				Col:  1,
				Message: formatMessage(
					i, "source", source,
					"dead redirect (source resolves to a real page)",
				),
			}}
		}
		return nil
	}

	return validateWildcardSource(wildcardSourceCtx{
		i:              i,
		source:         source,
		prefix:         prefix,
		prefixFs:       prefixFs,
		openAPISources: openAPISources,
		line:           line,
	})
}

func validateDestination(
	i int,
	destination, contentRoot string,
	openAPISources map[string]bool,
	line int,
) []analysis.Violation {
	cleaned := cleanPath(destination)
	prefix, hasWildcard := splitWildcard(cleaned)
	prefixRel := strings.Join(prefix, "/")
	prefixFs := filepath.Join(contentRoot, filepath.Join(prefix...))

	if !hasWildcard {
		if pageExists(prefixFs) {
			return nil
		}
		return []analysis.Violation{{
			File: "docs.json",
			Line: line,
			Col:  1,
			Message: formatMessage(
				i, "destination", destination,
				"missing destination (no .mdx or .md page exists)",
			),
		}}
	}

	if dirExists(prefixFs) {
		return nil
	}
	yamlRel := prefixRel + ".yaml"
	yamlFs := prefixFs + ".yaml"
	if openAPISources[yamlRel] && fileExists(yamlFs) {
		return nil
	}
	return []analysis.Violation{{
		File: "docs.json",
		Line: line,
		Col:  1,
		Message: formatMessage(
			i, "destination", destination,
			"wildcard prefix not a directory",
		),
	}}
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
		full := filepath.Join(dir, ent.Name())
		if ent.IsDir() {
			if dirHasPages(full) {
				return true
			}
			continue
		}
		name := ent.Name()
		if strings.HasSuffix(name, ".mdx") || strings.HasSuffix(name, ".md") {
			return true
		}
	}
	return false
}

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

# Fix codex findings on redirect lint rule (PR #75)

## Scope

This plan lives in the docs repo at `/home/ben/miru/workbench1/repos/docs/plans/backlog/20260429-redirect-rule-codex-fixes.md`. The docs repo is the sole read-write repository. Branch: `feat/redirect-lint-rule`. Date: 2026-04-29. One milestone, one commit.

| Repository | Access     | Description                          |
|------------|------------|--------------------------------------|
| docs       | read-write | Mintlify docs site; redirect linter  |

All paths are absolute. Repo root: `/home/ben/miru/workbench1/repos/docs/`.

## Purpose / Big Picture

PR #75 ports the redirects lint rule to Go. Codex review flagged two correctness bugs in `tools/lint/linter/redirects/redirects.go`:

- **P1**: Dot-segment paths like `/docs/../README` pass the `HasPrefix(cleaned, "docs/")` gate, then `filepath.Join` resolves them outside the docs root, silently inspecting arbitrary repo files.
- **P2**: Wildcard sources lack the OpenAPI escape hatch that destinations have. A wildcard source pointing at a Mintlify-generated route (`${prefix}.yaml` registered in `nav.*.openapi.source`) is wrongly classified as alive, missing real dead-redirect detections.

After this work, both source and destination path validation reject `..`/`.` segments with a stable diagnostic, and `validateSource` mirrors `validateDestination`'s OpenAPI awareness.

## Progress

Add entries as work proceeds.

## Surprises & Discoveries

Add entries as work proceeds.

## Decision Log

Add entries as work proceeds.

## Outcomes & Retrospective

Add entries as work proceeds.

## Context and Orientation

`tools/lint/linter/redirects/redirects.go` (608 lines) implements the redirects rule. Key entry points:

- `validate` (line 57) — parses `docs.json`, collects `openAPISources`, dispatches per entry.
- `validateEntry` → `validateFileSystem` (line 197) — branches to `validateSource` (currently no `openAPISources` arg, line 247) and `validateDestination` (line 310, takes `openAPISources`).
- `cleanPath` (line 370) — strips query/fragment/leading-and-trailing slashes; does NOT reject traversal.
- `splitWildcard` (line 386) — returns prefix segments (skipping empty segments) and a wildcard flag.
- `pageExists` / `dirExists` / `dirHasPages` / `fileExists` (lines 408-449) — fs probes.
- `collectOpenAPISources` / `walkOpenAPISources` (lines 456-478) — yields the `map[string]bool` of registered yaml relpaths.

Destination escape hatch (lines 351-355): if `prefixRel + ".yaml"` is in `openAPISources` AND `fileExists(prefixFs+".yaml")`, accept the wildcard destination. The fix mirrors this for sources.

`tools/lint/linter/redirects/redirects_test.go` uses table-driven `TestValidate` cases over a `(name, docsJSON, files, wants)` shape. Helpers: `setupContentRoot(t, files)` writes a tempdir (keys ending `/` create dirs, others write file contents), and `assertViolation(t, v, wantLine, substr)` checks `File="docs.json"`, `Col=1`, optional line, and a substring of `Message`. Existing case 21 (`destination_wildcard_openapi_yaml_ok`) is the template for wiring an `openapi.source` value into a fixture's `docs.json`.

Coverage gate: `tools/lint/scripts/covgate.sh` enforces 100.0% on the redirects package — every new code branch needs a test.

## Plan of Work

One milestone, one commit. Source edits and tests land together so coverage stays at 100% on every commit.

### M1: Reject traversal segments and add OpenAPI source escape hatch

Add `containsTraversalSegment([]string) bool` helper. Have `validateFileSystem` (or each side's call site) call `splitWildcard(cleanPath(path))` first; if the helper returns true, emit a diagnostic with the original path and skip downstream fs checks for that side. Apply on both source and destination.

Extend `validateSource` to accept `openAPISources map[string]bool` (matching `validateDestination`'s signature). After the existing wildcard checks (page-exists then dir-with-pages), add the OpenAPI check: if `openAPISources[prefixRel+".yaml"]` AND `fileExists(prefixFs+".yaml")`, emit the new dead-redirect diagnostic. Plumb the map through `validateFileSystem`.

Extend `TestValidate` with the six new cases in Concrete Steps. The existing 9 positive-case contract messages stay byte-for-byte unchanged.

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs/` unless noted.

Step 1. Edit `tools/lint/linter/redirects/redirects.go`:

  - Add `containsTraversalSegment(segments []string) bool` returning true if any segment equals `..` or `.`.
  - In `validateFileSystem`, for the source branch: compute `prefix, _ := splitWildcard(cleanPath(f.source))`; if `containsTraversalSegment(prefix)`, append a violation with `field="source"`, `value=f.source`, message `"bad path: contains '..' or '.' segment"` and skip the call to `validateSource`. Mirror the same check on the destination branch (skip `validateDestination`). Order: do the traversal check before the existing fs dispatch.
  - Change `validateSource` signature to `validateSource(i int, source, contentRoot string, openAPISources map[string]bool, line int) []analysis.Violation`. Inside the wildcard branch, after the existing `pageExists` and `dirExists && dirHasPages` checks, compute `prefixRel := strings.Join(prefix, "/")`; if `openAPISources[prefixRel+".yaml"] && fileExists(prefixFs+".yaml")`, return a violation with message `"dead redirect (wildcard source prefix has Mintlify-generated pages)"`.
  - Update `validateFileSystem`'s call to `validateSource` to pass `openAPISources`.

Step 2. Extend `tools/lint/linter/redirects/redirects_test.go` `TestValidate.cases` with six new entries:

  - `source_with_dot_dot_segment`: `{"source":"/docs/../README","destination":"/docs/y"}`, files `{"docs/y.mdx":"y"}`, want substr `bad path: contains '..' or '.' segment` on source.
  - `destination_with_dot_dot_segment`: `{"source":"/docs/x","destination":"/docs/x/../../escape"}`, want substr `bad path: contains '..' or '.' segment` on destination.
  - `source_with_dot_segment`: `{"source":"/docs/./foo","destination":"/docs/y"}`, files `{"docs/y.mdx":"y"}`, want substr `bad path: contains '..' or '.' segment` on source.
  - `wildcard_source_prefix_is_openapi_yaml`: docsJSON `{"nav":[{"openapi":{"source":"docs/api/spec.yaml"}}],"redirects":[{"source":"/docs/api/spec/:slug*","destination":"/docs/y"}]}`, files `{"docs/api/spec.yaml":"openapi: 3.0","docs/y.mdx":"y"}`, want substr `dead redirect (wildcard source prefix has Mintlify-generated pages)`.
  - `wildcard_source_prefix_yaml_registered_but_yaml_missing`: same docsJSON as above, files `{"docs/y.mdx":"y"}` (no yaml file), want `nil` (symmetric with case 22).
  - `wildcard_source_prefix_yaml_not_registered`: `{"redirects":[{"source":"/docs/api/spec/:slug*","destination":"/docs/y"}]}`, files `{"docs/y.mdx":"y"}`, want `nil` (false-positive guard; no nav, no .mdx, no dir-with-pages).

Step 3. From `tools/lint/`, run unit tests and coverage gate:

    cd tools/lint
    go test ./linter/redirects/...
    ./scripts/covgate.sh

Both must pass. Covgate must report `redirects` at `100.0%`.

Step 4. From repo root, run preflight:

    ./scripts/preflight.sh

Must report clean before push (modulo the pre-existing local `bats: command not found` env issue; CI installs bats).

Step 5. Commit:

    git add tools/lint/linter/redirects/redirects.go tools/lint/linter/redirects/redirects_test.go
    git commit -m "fix(lint): reject dot-segment paths and detect OpenAPI-generated source dead-redirects"

## Validation and Acceptance

- `go test ./linter/redirects/...` exits 0.
- `./tools/lint/scripts/covgate.sh` reports `redirects` at `100.0%`.
- `./scripts/preflight.sh` reports clean before push (modulo pre-existing local bats env issue).
- The 9 existing positive-case contract message strings in `TestValidate` continue to pass unchanged.
- The 2 new diagnostic substrings — `bad path: contains '..' or '.' segment` and `dead redirect (wildcard source prefix has Mintlify-generated pages)` — appear in violations as asserted by the new test cases.
- A wildcard source pointing at a registered `${prefix}.yaml` whose file is missing produces zero violations (symmetric with destination case 22).

## Idempotence and Recovery

- Re-running steps after a partial failure is safe: edits are surgical and tests are additive — re-applying the same edit is a no-op once present.
- Reverting via `git revert <commit-sha>` restores original behavior cleanly; no schema or external state changes.
- No data migration; no destructive operations; no infrastructure touched.

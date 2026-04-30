# Heading-case allowlist for acronyms and version tags

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `repos/docs/` | read-write | Extend the `heading-case` lint rule with an acronym/proper-noun allowlist plus a version-tag exemption, and rewrite ~22 genuine title-case offending headings so the lint pipeline exits 0. |

This plan lives in `repos/docs/plans/backlog/` because all code and content edits are inside `repos/docs/`. Work proceeds on the existing branch `feat/title-casing-lint-rule` (PR #76).

## Purpose / Big Picture

PR #76 introduces a lint rule that requires markdown headings to be sentence-case. The rule is currently too strict: it flags any heading whose later tokens are uppercase, including headings like `## API reference` whose uppercase token is a legitimate acronym. The rule also has no exemption for changelog-style headings such as `# v0.6.0` or `# 2026-03-09.tetons`.

After this change, the `heading-case` rule allows a hardcoded set of acronyms, proper nouns, and codenames, and exempts version-tag headings via regex. Approximately 22 genuine title-case offenders (headings whose uppercase tokens are NOT acronyms — e.g. `## Communication Protocols`, `## Default Permissions`) are rewritten to sentence-case in the same change. With both edits shipping together, `./scripts/lint.sh` exits 0 with zero diagnostics, and CI for PR #76 turns green.

User-visible behavior: a writer who lands a new heading like `## API reference` in any `.mdx` file under `docs/` no longer gets a lint failure; a writer who lands `## API Key` still gets a lint failure (because `Key` is uppercase and not allowlisted); a writer who lands `# v0.6.0` in a changelog has no failure.

## Progress

- [ ] M1 — Extend `headingcase.go` with allowlist, version-tag regex, tokenization-based check, and updated Message constant. Commit.
- [ ] M2 — Update unit tests in `headingcase_test.go` and add an E2E sub-test in `tools/lint/main_test.go`. Commit.
- [ ] M3 — Rewrite the 22 listed heading offenders across ~13 `.mdx` files. Commit.
- [ ] M4 — Preflight: `go test ./...`, `go vet ./...`, and `./scripts/lint.sh` from repo root all exit 0. No commit unless preflight surfaces an issue.

Use timestamps when steps complete. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

- Decision: Hardcode the allowlist as a package-scope `var allowlist = map[string]struct{}{...}` in `headingcase.go` rather than reading from an external config file (e.g. a YAML or `.lint.yml` next to the binary).
  Rationale: minimum-change, fastest path to green CI for PR #76. The set of acronyms is small and stable; externalizing it adds config-loading code, file-discovery rules, and another surface to test. A future PR can promote the allowlist to a config file once we have a real second consumer.
  Date/Author: 2026-04-29 / agents@miruml.com

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

The lint binary lives under `repos/docs/tools/lint/`. The package layout:

- `tools/lint/main.go` — CLI entrypoint. Walks `docs/` and runs each registered linter against every `.mdx` file.
- `tools/lint/main_test.go` — End-to-end tests that exec the built binary against fixture trees and assert exit code + stdout.
- `tools/lint/linter/headingcase/headingcase.go` — The rule under change. Exports `Linter` and a `Message` constant.
- `tools/lint/linter/headingcase/headingcase_test.go` — Table-driven unit tests using a `build` helper and `analysis.NewScanner().ScanLine` (mirrors how `main.go` feeds lines).
- `tools/lint/analysis/` — Provides the line-by-line scanner with frontmatter-aware state.

The current rule (pre-change) is roughly:

- For markdown headings (`#`, `##`, `###`, …) and frontmatter `title:` values: take the heading text, strip leading/trailing whitespace, strip trailing `.?!:`, and require the first ASCII letter to be uppercase and every subsequent ASCII letter to be lowercase. No allowlist; no version exemption.
- The current `Message` ends with `; proper nouns/acronyms are not yet supported`.

The lint pipeline is invoked from `repos/docs/scripts/lint.sh`, which builds the Go binary in `tools/lint/` and runs it across `docs/`.

Definitions used below:

- "Token": a whitespace-delimited substring of the heading text after trimming leading/trailing whitespace.
- "Core": the result of stripping surrounding non-letter punctuation `,;()[]{}"'` and trailing `.` from a token (with the exception that a `.` between alphanumerics, e.g. `deployment.deployed`, is preserved as part of the core).
- "Rule-bound token": a token whose core is non-empty, is not version-tag-like, and is not in the allowlist. These are the tokens whose casing is checked.

## Plan of Work

### 1. Extend the rule (`tools/lint/linter/headingcase/headingcase.go`)

Add at package scope:

- `var versionTagRe = regexp.MustCompile(\`^v?\d+([.-]\w+)*\`)` — note the regex must be anchored at both ends; build it as `^v?\d+([.-]\w+)*$`.
- `var allowlist = map[string]struct{}{...}` populated with the tokens listed below.

Allowlist contents (case-sensitive exact-token match):

- Acronyms: `API`, `APIs`, `CLI`, `CI`, `SDK`, `SDKs`, `CUE`, `JSON`, `MQTT`, `TLS`, `HTTPS`, `REST`, `GUI`, `URL`, `ACLs`, `SSE`, `OpenAPI`
- Proper nouns: `Miru`, `GitHub`, `Agent`, `Unix`, `Git`, `Python`, `Schema`, `Base`, `Head`
- Codenames: `tetons`, `zion`
- Specific event identifiers used as page titles: `deployment.deployed`, `deployment.removed`

Replace the existing `Message` constant with:

    heading-case: heading must be sentence-case (first letter uppercase, all other letters lowercase)

Rewrite `casingViolation` (or whatever the current internal predicate is named) to be word-aware. Algorithm specification, verbatim:

1. Trim leading/trailing whitespace from the heading text.
2. Trim trailing punctuation in the set `.?!:` (one or more, repeatedly).
3. If the trimmed result is empty, return false (no violation).
4. If the trimmed result matches the version-tag regex, return false (whole heading exempt).
5. Tokenize on whitespace into `tokens[]`.
6. Walk tokens left to right. For each token at index `i`:
   - Strip surrounding non-letter punctuation from the token: `,;()[]{}"'` and trailing `.` (but only when not part of an event-name like `deployment.deployed`). Treat the token's "core" as the result.
   - If the core is empty after stripping, skip and do NOT increment a "non-skipped" counter.
   - If the core matches the version-tag regex, skip.
   - If the core is in the allowlist, skip.
   - Otherwise (core is a "rule-bound" token):
     - If `i == 0` (the LEFTMOST token of the heading): require first ASCII letter uppercase, every subsequent ASCII letter lowercase.
     - Else (any later token, regardless of how many earlier ones were skipped): require all ASCII letters lowercase.
   - On any violation, return true (one violation per heading).

Note: the "first leftmost token" rule deliberately uses index 0 — if the leftmost token is allowlisted/version, the whole heading begins with that allowlisted text and the visual sentence-case-first-letter requirement is satisfied by the allowlisted token. Subsequent rule-bound tokens are all "non-first" and must be all-lowercase.

Algorithm worked examples (verification table):

- `## API reference` → t0=API allowed; t1=reference all-lower → PASS
- `## Configure deployments` → t0 first-cap-rest-lower; t1 all-lower → PASS
- `## Configure Deployments` → t1 has uppercase, not allowed → FAIL
- `## API Key authentication` → t1=Key uppercase, not allowed → FAIL
- `## Create the API key` → t0 cap+lower; t1=the lower; t2=API allowed; t3=key lower → PASS
- `## Create the API Key` → t3=Key uppercase, not allowed → FAIL
- `## reference API` → t0=reference, not allowed, not first-cap → FAIL
- `# v0.6.0` → version-regex on whole heading → PASS
- `## CUE support` → t0=CUE allowed; t1 lower → PASS
- `### Approachable GUI` → t0 cap+lower; t1=GUI allowed → PASS
- `### Approachable Gui` → t1=Gui (NOT in allowlist; only `GUI` is) → FAIL

### 2. Update tests

In `tools/lint/linter/headingcase/headingcase_test.go`:

(a) Update existing case in `TestCheck_Headings`:
- The case currently named `bad acronym (v1 limitation)` with content `"## API reference\n"` and `wantCount: 1` MUST be changed to `wantCount: 0` and renamed to `clean acronym (allowlisted)`. Drop the `wantLine` and `wantCol` fields for this case (no violation expected).

(b) New table-driven test `TestCheck_Allowlist` with these cases (each multi-line content fed through `analysis.NewScanner().ScanLine`, same `build` helper as existing tests):
- `clean acronym leading`: `"## API reference\n"` → 0
- `clean acronym mid-sentence`: `"## Set up the CLI\n"` → 0
- `bad uppercase non-allowlisted after allowed`: `"## API Key authentication\n"` → 1 (Col=4)
- `bad lowercase leading non-allowlisted`: `"## reference API\n"` → 1 (Col=4)
- `clean version tag`: `"# v0.6.0\n"` → 0
- `clean date-prefixed codename`: `"# 2026-03-09.tetons\n"` → 0
- `clean OpenAPI specifications`: `"## OpenAPI specifications\n"` → 0
- `clean GitHub actions`: `"## GitHub actions\n"` → 0
- `bad GitHub Actions`: `"## GitHub Actions\n"` → 1 (Col=4)
- `clean Miru help`: `"## How does Miru help?\n"` → 0
- `clean Base/Head Git refs`: `"### Base vs. Head\n"` → 0
- `bad lowercase Gui (only uppercase GUI is allowed)`: `"### Approachable Gui\n"` → 1 (Col=5)

(c) New cases in `TestCheck_FrontmatterTitle`:
- `clean SDKs title (allowlisted)`: `"---\ntitle: \"SDKs\"\n---\n"` → 0
- `clean event identifier title (allowlisted)`: `"---\ntitle: \"deployment.deployed\"\n---\n"` → 0
- `clean GitHub allowlist title`: `"---\ntitle: \"GitHub actions\"\n---\n"` → 0
- `bad GitHub Actions title`: `"---\ntitle: \"GitHub Actions\"\n---\n"` → 1, line=2, col=9

(d) Update assertions to use the NEW Message constant value (no proper-noun disclaimer).

In `tools/lint/main_test.go`:

- The existing `clean headings return 0` sub-test continues to pass (no changes needed). Verify by re-running.
- The existing `heading-case violation returns 1` sub-test (`title: "User Management"`) continues to fail correctly — the new rule still flags `Management` as uppercase non-first non-allowlisted. No changes needed.
- Add ONE new sub-test: `t.Run("clean allowlisted acronym title returns 0", ...)`. Content: `"---\ntitle: \"API keys\"\n---\n\n## OpenAPI specifications\n"`. Expect exit 0, empty stdout. This proves the allowlist fires through the full pipeline.

### 3. Rewrite genuine title-case offenders

For each row below, locate the heading by TEXT (not line number) — line numbers may have drifted slightly. Use grep/Read to verify before editing. Heading rewrites must be MINIMAL — only the heading text changes; no surrounding prose edits, no formatting churn, no link or anchor changes.

| # | File | Current | New |
|---|------|---------|-----|
| 1 | `docs/changelogs/agent.mdx` | `# Agent Install Script Update` | `# Agent install script update` |
| 2 | `docs/admin/apikeys.mdx` | `title: "API Keys"` | `title: "API keys"` |
| 3 | `docs/developers/agent/architecture.mdx` | `## Communication Protocols` | `## Communication protocols` |
| 4 | `docs/developers/agent/file-permissions.mdx` | `title: 'File Permissions'` | `title: 'File permissions'` |
| 5 | `docs/developers/agent/file-permissions.mdx` | `## Default Permissions` | `## Default permissions` |
| 6 | `docs/developers/agent/file-permissions.mdx` | `## Custom File Paths` | `## Custom file paths` |
| 7 | `docs/developers/agent/file-permissions.mdx` | `### Required Permissions` | `### Required permissions` |
| 8 | `docs/developers/ci/gh-actions.mdx` | `title: "GitHub Actions"` | `title: "GitHub actions"` |
| 9 | `docs/developers/ci/gh-actions.mdx` | `## Supported Platforms` | `## Supported platforms` |
| 10 | `docs/learn/deployments.mdx` | `## Target Status` | `## Target status` |
| 11 | `docs/learn/deployments.mdx` | `## Activity Status` | `## Activity status` |
| 12 | `docs/learn/deployments.mdx` | `## Error Status` | `## Error status` |
| 13 | `docs/learn/devices/provision.mdx` | `### Create a Device` | `### Create a device` |
| 14 | `docs/learn/devices/provision.mdx` | `### Verify the Installation` | `### Verify the installation` |
| 15 | `docs/learn/devices/provision.mdx` | `## API Keys` | `## API keys` |
| 16 | `docs/learn/devices/provision.mdx` | `### Create the API Key` | `### Create the API key` |
| 17 | `docs/learn/devices/provision.mdx` | `### Provision the Device` | `### Provision the device` |
| 18 | `docs/learn/devices/provision.mdx` | `### Verify Installation` | `### Verify installation` |
| 19 | `docs/learn/devices/provision.mdx` | `## Poor Connectivity` | `## Poor connectivity` |
| 20 | `docs/references/cli/release-create.mdx` | `### Schema Annotations` | `### Schema annotations` |
| 21 | `docs/references/device-api/v0.2.1/events/deployment-deployed.mdx` | `## Event Data` | `## Event data` |
| 22 | `docs/references/device-api/v0.2.1/events/deployment-removed.mdx` | `## Event Data` | `## Event data` |

## Concrete Steps

One commit per milestone. All commands run from `repos/docs/` unless otherwise stated.

### M1 — Extend the rule

1. Edit `tools/lint/linter/headingcase/headingcase.go`:
   - Add the `versionTagRe` package var (regex: `^v?\d+([.-]\w+)*$`).
   - Add the `allowlist` package var with the tokens listed in Plan of Work.
   - Replace the `Message` constant with the new value.
   - Rewrite the casing check per the algorithm spec.
2. Build:

       cd tools/lint && go build -o lint .

   Expect: no output, exit 0.
3. Commit:

       git add tools/lint/linter/headingcase/headingcase.go
       git commit -m "feat(lint): heading-case allowlist for acronyms and version tags"

### M2 — Update tests

1. Edit `tools/lint/linter/headingcase/headingcase_test.go` per (a)+(b)+(c)+(d).
2. Edit `tools/lint/main_test.go` to add the new sub-test `clean allowlisted acronym title returns 0`.
3. Run unit tests:

       cd tools/lint && go test ./...

   Expect: `ok` for every package; the new `TestCheck_Allowlist` cases all pass.
4. Commit:

       git add tools/lint/linter/headingcase/headingcase_test.go tools/lint/main_test.go
       git commit -m "test(lint): cover heading-case allowlist and version-tag exemption"

### M3 — Rewrite genuine title-case offenders

1. For each row in the rewrite table:
   - Use `grep -n "<current heading text>" <file>` to locate the heading.
   - Use `Edit` to apply the rewrite — only the heading text changes.
   - Verify with `Read` that the surrounding prose is unchanged.
2. After all 22 rewrites, sample-verify with grep:

       grep -n "Communication protocols" docs/developers/agent/architecture.mdx
       grep -n "Default permissions" docs/developers/agent/file-permissions.mdx
       grep -n "Supported platforms" docs/developers/ci/gh-actions.mdx

   Expect: each grep prints exactly one line.
3. Commit:

       git add docs/
       git commit -m "docs: rewrite headings to satisfy heading-case rule"

### M4 — Preflight gate

1. From `tools/lint/`:

       go test ./...
       go vet ./...

   Both must exit 0.
2. From repo root (`repos/docs/`):

       ./scripts/lint.sh

   Expect: exit 0 with zero diagnostic lines. Verify with:

       ./scripts/lint.sh 2>&1 | grep -E ':[0-9]+:[0-9]+:' | wc -l

   Expect: `0`.
3. If any heading-case violations remain, audit and either rewrite the heading or extend the allowlist (re-run M1/M3 commits as needed). No new commit required unless preflight surfaces an issue.

## Validation and Acceptance

The change is acceptable when ALL hold:

1. From `tools/lint/`, `go test ./...` exits 0; the new `TestCheck_Allowlist` cases all pass.
2. From repo root, `./scripts/lint.sh` exits 0 and produces zero diagnostic lines:

       ./scripts/lint.sh 2>&1 | grep -E ':[0-9]+:[0-9]+:' | wc -l

   prints `0`.
3. The diagnostic emitted by the rule (e.g., on a hand-crafted bad heading like `## Configure Deployments`) uses the NEW `Message` text without the proper-noun disclaimer. Quick check:

       printf '## Configure Deployments\n' > /tmp/bad.mdx && ./tools/lint/lint /tmp/bad.mdx || true

   Expect a diagnostic ending in `heading must be sentence-case (first letter uppercase, all other letters lowercase)` and NOT containing `proper nouns/acronyms are not yet supported`.
4. The 22 listed heading rewrites are present in the working tree. Verify a sample of 3-4 with `grep -n` after M3 (e.g. `grep -n "Communication protocols" docs/developers/agent/architecture.mdx`).

### Non-negotiable constraints

- Allowlist is hardcoded in `headingcase.go` as `var allowlist = map[string]struct{}{...}` at package scope. No external config file in this iteration.
- Version-tag regex compiled once with `regexp.MustCompile` at package scope.
- Heading rewrites must be MINIMAL — only the heading text changes; no surrounding prose edits, no formatting churn, no link or anchor changes.
- Do NOT add anything under `tests/lint-fixtures/`.

## Idempotence and Recovery

- M1 (rule extension) is idempotent: re-applying the edit yields the same file. If the build fails, fix and rebuild — no rollback needed.
- M2 (tests) is idempotent: re-running `go test ./...` is safe and repeatable.
- M3 (heading rewrites) is idempotent per row: each row's grep-then-edit either matches the old text and rewrites, or matches the new text and is a no-op. If a heading appears more than once in a file, audit by hand before applying — the rewrite table assumes one occurrence per file.
- M4 (preflight) is read-only (runs tests + lint); no recovery needed.

If a milestone commit needs to be redone, prefer `git commit --amend` only if the previous commit on the branch is the one being fixed AND the branch hasn't been pushed since. Otherwise, create a follow-up commit ("fixup:" prefix) and squash interactively before merge. Never force-push to PR #76 after others have reviewed.

If `./scripts/lint.sh` surfaces a heading-case violation that wasn't in the rewrite table, options:

1. The heading is a genuine offender → rewrite to sentence-case in a new commit `docs: rewrite additional heading missed in initial pass`.
2. The heading uses an acronym not in the allowlist → extend `allowlist` in `headingcase.go` and add a unit-test case in `TestCheck_Allowlist`, then commit `feat(lint): extend heading-case allowlist with <token>`.

Pick option (1) by default; reach for (2) only when the token is unambiguously a proper noun or industry-standard acronym.

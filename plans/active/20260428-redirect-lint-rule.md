# Redirect Lint Rule for docs.json

Validate every redirect in `docs.json` against the on-disk `docs/` tree so
unreachable or broken redirects fail lint locally and in CI.

## Scope

| Repo | Path | Access | Description |
| --- | --- | --- | --- |
| `mirurobotics/docs` | `/home/ben/miru/workbench1/repos/docs/` | read-write | Branch `feat/redirect-lint-rule` already checked out. |

Out of scope: other repos, OpenAPI page regeneration, modifying existing
redirects in `docs.json`, migrating other lint checks.

## Purpose / Big Picture

Mintlify serves `docs/foo/bar.mdx` at URL `/docs/foo/bar`. The `redirects`
array in `docs.json` rewrites URLs at the edge. Two undetected failure modes:

1. A `source` matching a real page is **dead config** — Mintlify serves the
   real page and the redirect never fires.
2. A `destination` not matching a real page **404s** users.

This plan adds a check to `scripts/lint.sh` (and therefore CI) that catches
both, plus invalid prefixes and unsupported schemes. Running `pnpm run lint`
fails on violations with a `file:line: message` diagnostic.

## Progress

- [ ] M1: Implement `scripts/check-redirects.mjs`
- [ ] M2: Wire into `scripts/lint.sh`
- [ ] M3: Add `bad-redirects` test fixture and test runner case
- [ ] M4: Validate against real `docs.json` and finalize

## Surprises & Discoveries

Add entries as work proceeds.

## Decision Log

Add entries as work proceeds.

## Outcomes & Retrospective

Fill in when the plan is moved to `plans/completed/`.

## Context and Orientation

Repo layout (only the parts this plan touches):

- `docs.json` — Mintlify site config. The `"redirects"` array (around lines
  361-374 today) holds `{source, destination}` objects. No `permanent` field.
  Three entries today; two use `:slug*` wildcards.
- `docs/` — 66 `.mdx` files in subdirs (`admin/`, `changelogs/`,
  `developers/`, `getting-started/`, `learn/`, `references/`). File
  `docs/foo/bar.mdx` is served at `/docs/foo/bar`. No `index.mdx` files;
  directories are not pages.
- `api/`, `pub/`, `snippets/` — NOT served as pages. Redirects whose
  `destination` resolves outside `docs/` are invalid.
- `scripts/lint.sh` — Bash entry point, `set -euo pipefail`. Resolves
  `script_dir` and `repo_root` dynamically; honors `DOCS_LINT_ROOT` to point
  checks at fixture roots. Each sub-check prints `echo "== <Name> =="`. Exits
  non-zero if any sub-check fails.
- `tests/test-lint.sh` — Bash test runner with `run_expect_pass(fixture)` and
  `run_expect_fail(fixture, pattern)` helpers. Each fixture is self-contained
  at `tests/lint-fixtures/<name>/`.
- `tools/lint/` — Go MDX linter. Not modified; only its output format is mimicked.
- `.github/workflows/ci.yml` — `lint` job runs `pnpm run test:lint` then
  `./scripts/lint.sh`. Adding the check to `scripts/lint.sh` puts it in CI.
- `package.json` — `"lint": "./scripts/lint.sh"`,
  `"test:lint": "./tests/test-lint.sh"`. `pnpm@10.17.0` (corepack), Node 22.

Design choices baked into this plan:

- **Language: Node.js** (`scripts/check-redirects.mjs`, ESM, no deps). Node
  is already required for Mintlify and pnpm. Bash needs `jq`; Go would force
  rebuilding for a 60-line JSON-only check.
- **Output format**: `docs.json:<line>: <message>`, matching the
  `file:line:col:` style emitted by `tools/lint/main.go`. Line numbers come
  from scanning `docs.json` text for each redirect's object literal. On
  lookup failure, fall back to `docs.json: redirects[<index>]: <message>`
  and print a one-line `warning:`-prefixed message.
- **Wildcard semantics**: a path segment matching
  `/^:[A-Za-z][A-Za-z0-9]*\*?$/` is a wildcard. The "prefix" is everything
  before the first wildcard segment. For `source`, the prefix MUST NOT be a
  real page and MUST NOT be a real directory containing any `.mdx`/`.md`
  page. For `destination`, the prefix MUST be a real directory under `docs/`.
- **Fixture strategy**: a NEW `tests/lint-fixtures/bad-redirects/` copied
  from `good/` (inherits passing baseline for unrelated checks); only
  `docs.json` and supporting `docs/` files are mutated. The real `docs.json`
  exercised by the top-level lint already covers the passing case.
- **Override hook**: respect `DOCS_LINT_ROOT`. The script reads
  `${DOCS_LINT_ROOT}/docs.json` and resolves files under
  `${DOCS_LINT_ROOT}/docs/`. Defaults to repo root when unset.

URL ↔ filesystem mapping rules:

1. Strip a leading `/`; remainder is repo-relative.
2. Strip any trailing `/`, `?...` query, `#...` fragment.
3. The path MUST start with `docs/`. Anything else is invalid.
4. Append `.mdx` to get the candidate filename; if absent, fall back to
   `.md`. Filesystem checks are case-sensitive. (Today the tree is 100%
   `.mdx`; the resolver tries both for forward compatibility.)
5. Destinations starting with `http://` or `https://` are skipped.
6. Wildcard segments are not part of the filesystem path — only the prefix
   (segments before the first wildcard segment) is checked.

Edge cases:

- `docs/references/device-api/v0.2.1/` and
  `docs/references/platform-api/2026-03-09/` are populated by
  `api/generate_event_pages.py` from OpenAPI specs. Fixtures and unit
  checks must NOT require running the generator. The real-tree run only
  needs the directories to exist for wildcard prefix validation; if
  generator output is uncommitted, document in Surprises & Discoveries
  rather than relaxing the rule silently.

## Plan of Work

Four milestones, each one commit. See Concrete Steps for command-level detail.

- M1: Author `scripts/check-redirects.mjs`.
- M2: Wire into `scripts/lint.sh` under `== Redirects ==`.
- M3: Add `bad-redirects` fixture and `run_expect_fail` case.
- M4: Run full lint/test stack plus preflight; fix any fallout.

Wiring into `scripts/lint.sh` (rather than a new workflow step) means CI
picks the check up via the existing `lint` job — no `ci.yml` change. The
script remains independently runnable as `node scripts/check-redirects.mjs`.

## Concrete Steps

Run all commands from `/home/ben/miru/workbench1/repos/docs/`. Branch
`feat/redirect-lint-rule` is already checked out.

### M1 — `scripts/check-redirects.mjs`

1. Create `scripts/check-redirects.mjs`. Required behavior:

   - Derive `__dirname` using the standard ESM idiom:

         import path from 'node:path';
         import { fileURLToPath } from 'node:url';
         const __dirname = path.dirname(fileURLToPath(import.meta.url));

   - Resolve `root = process.env.DOCS_LINT_ROOT || path.resolve(__dirname, '..')`.
   - Read `${root}/docs.json` as text. Parse JSON for the `redirects`
     array. Keep the original text for line-number lookup.
   - Input handling:
     - If `${root}/docs.json` does not exist, exit 0 and print
       `No docs.json at ${root}; nothing to check`.
     - If it exists but cannot be parsed, exit 1 and print
       `docs.json: invalid JSON: <parser-message>` on stderr.
     - If parsed JSON has no `redirects` key, or `redirects` is empty,
       exit 0 and print `Checked 0 redirects: OK`.
   - For each entry, validate:
     a. Both `source` and `destination` are non-empty strings.
     b. `source` MUST start with `/`. `destination` MUST start with `/`
        OR `http://` OR `https://`. Otherwise emit `bad path: must
        start with '/'` (the `destination` message also notes
        `http(s)://` is accepted).
     c. `destination` starting with `http://` or `https://` is skipped
        for remaining (filesystem) checks.
     d. Strip leading `/`, trailing `/`, `?...`, `#...` from each path.
     e. The remainder MUST start with `docs/`. Otherwise emit
        `bad prefix (must start with /docs/)`.
     f. Split into segments. The first segment matching
        `/^:[A-Za-z][A-Za-z0-9]*\*?$/` is the wildcard boundary;
        segments before it form the `prefix`.
     g. Resolve `prefixFs = path.join(root, prefix.join('/'))`.
     h. **Source rules**:
        - No wildcard: neither `${prefixFs}.mdx` nor `${prefixFs}.md`
          may exist (otherwise dead).
        - Wildcard: `${prefixFs}` MUST NOT be a directory containing
          any `.mdx` or `.md` files (recursive). It MUST also NOT be
          a file at `${prefixFs}.mdx` or `${prefixFs}.md`.
     i. **Destination rules**:
        - No wildcard: a file at `${prefixFs}.mdx` OR `${prefixFs}.md`
          MUST exist.
        - Wildcard: `${prefixFs}` MUST be an existing directory.
   - For each violation, emit
     `docs.json:<line>: redirects[<i>] <field> "<value>": <message>`
     on stderr. Compute `<line>` by:
     1. Locating the `"redirects"` key in the original text (e.g. via
        `/"redirects"\s*:\s*\[/`) and taking the offset of the opening `[`.
     2. From that offset, walking the text with a cursor: for each
        `"source":` literal encountered, advance past the match so the
        n-th `"source":` under `"redirects"` locates the n-th entry.
        Bound the scan by the parsed redirect count.
     3. Reporting the 1-based line of the located occurrence.
     Anchoring to `"redirects"` is required because `docs.json`
     contains other `"source":` keys (OpenAPI nav block, ~10
     occurrences before redirects). On lookup failure, emit
     `docs.json: redirects[<i>] ...` plus a single `warning:`-prefixed
     stderr line per script run.
   - Exit 0 if no violations, 1 otherwise. On success print
     `Checked N redirects: OK`.

2. Make it executable and verify:

       chmod +x scripts/check-redirects.mjs
       node scripts/check-redirects.mjs

   The script must pass against the real repo (the three existing
   redirects are valid). If not, investigate and record in
   Surprises & Discoveries.

3. Commit:

       git add scripts/check-redirects.mjs
       git commit -m "feat(lint): add redirect validator script for docs.json"

### M2 — Wire into `scripts/lint.sh`

1. Open `scripts/lint.sh` and add a new section after the MDX prose
   section, using the same pattern:

       echo "== Redirects =="
       node "${repo_root}/scripts/check-redirects.mjs"

   The script honors `DOCS_LINT_ROOT` itself; no extra plumbing.

2. Run `./scripts/lint.sh` end-to-end and confirm the new section
   appears and passes.

3. Commit:

       git add scripts/lint.sh
       git commit -m "feat(lint): run redirect validator from lint.sh"

### M3 — Test fixture and runner case

1. Create `tests/lint-fixtures/bad-redirects/` by copying
   `tests/lint-fixtures/good/` wholesale, then mutate it:

   - Add a `docs.json` whose `redirects` array contains, as
     **separate distinct entries** (one violation per entry — do not
     bundle), at least:
     1. Source equals an existing page (dead redirect).
     2. Destination points to a missing page.
     3. Source uses a bad prefix (e.g. `/api/foo`).
     4. Destination uses a bad prefix (e.g. `/api/foo`).
     5. Source missing the leading `/`.
     6. Destination missing the leading `/` and not `http(s)://`.
     7. Wildcard destination prefix that does not exist as a directory.
     8. Wildcard source whose prefix directory contains real pages.
     Plus at least one valid redirect and one external `https://`
     destination so skip and pass-mixed-with-fail paths are exercised.
   - Add only the `docs/` files the violations rely on (e.g.
     `docs/admin/exists.mdx`, `docs/wild/page.mdx`). Keep each file
     to a single H1 line. Use plain words like `Page` so `cspell.json`
     (consulted from repo root, not fixture root), ESLint-MDX, and the
     Go linter all pass; only the redirect validator should fail.
   - The new `docs.json` is freshly authored — keep minimal: only the
     `redirects` array is required (per M1's input handling, missing
     top-level keys do not fail the redirect check).

2. Edit `tests/test-lint.sh` to add a `run_expect_fail` invocation for
   the new fixture. For each violating entry, assert the output
   contains both the per-entry prefix `redirects[<i>]` (with `<i>`
   matching the entry's array index in the fixture's `docs.json`) AND
   the rule-specific substring (`dead redirect`, `missing destination`,
   `bad prefix`, `bad path: must start with '/'`,
   `wildcard prefix not a directory`, etc.). If the helper only
   supports one pattern, run it once per (index, substring) pair.

3. Run `./tests/test-lint.sh` and confirm the new case fails as
   expected and all previously passing fixtures still pass.

4. Commit:

       git add tests/lint-fixtures/bad-redirects tests/test-lint.sh
       git commit -m "test(lint): add bad-redirects fixture for redirect validator"

### M4 — Full validation

1. Run the full stack:

       pnpm run test:lint
       pnpm run lint

   All must pass.

2. Run preflight:

       ./scripts/preflight.sh

   Must report clean before publishing the branch.

3. If validation required tweaks, commit them separately:

       git commit -am "fix(lint): <short description>"

4. Move the plan from `plans/backlog/` to `plans/active/` once
   implementation begins, and to `plans/completed/` once merged. Fill
   in Outcomes & Retrospective at completion.

## Validation and Acceptance

Accepted when **all** hold:

1. `node scripts/check-redirects.mjs` exits 0 against the real
   `docs.json`/`docs/` tree and prints `Checked 3 redirects: OK` (or
   the current count).
2. `pnpm run lint` succeeds locally and shows the `== Redirects ==`
   section.
3. `pnpm run test:lint` succeeds, including the new `bad-redirects`
   case, surfacing every expected diagnostic substring.
4. CI's `lint` job succeeds on the branch.
5. Manually injecting `{"source": "/docs/getting-started/intro",
   "destination": "/docs/nope"}` into a scratch copy of `docs.json`
   makes `pnpm run lint` fail with two diagnostics (dead source and
   missing destination). Revert before committing.
6. **`./scripts/preflight.sh` reports clean before push.** Mandatory.

## Idempotence and Recovery

- **Read-only.** The script never mutates `docs.json` or files under
  `docs/`. Re-runs are safe and identical.
- **Revert-by-commit rollback.** Each milestone is one commit; revert
  to roll back. No data migration.
- **Fixture-narrowing tip.** Point `DOCS_LINT_ROOT` at one fixture
  (e.g. `DOCS_LINT_ROOT=tests/lint-fixtures/bad-redirects node
  scripts/check-redirects.mjs`) to inspect only that tree.

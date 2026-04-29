# Redirect Lint Rule for docs.json

Validate every redirect in `docs.json` against the on-disk `docs/` tree so that
unreachable or broken redirects fail the lint job locally and in CI.

## Scope

| Repo | Path | Access | Description |
| --- | --- | --- | --- |
| `mirurobotics/docs` | `/home/ben/miru/workbench1/repos/docs/` | read-write | Single repo for all changes. Branch `feat/redirect-lint-rule` already checked out. |

Out of scope: changes to other repositories, regenerating OpenAPI pages, modifying
existing redirects in `docs.json`, or migrating other lint checks.

## Purpose / Big Picture

Mintlify serves `docs/foo/bar.mdx` at URL `/docs/foo/bar`. The `redirects` array
in `docs.json` rewrites incoming URLs at the edge. Two failure modes are
currently undetected:

1. A redirect whose `source` matches an existing real page is **dead config** —
   Mintlify serves the real page and the redirect never fires.
2. A redirect whose `destination` does not match a real page **404s** users.

This plan adds an automated check to `scripts/lint.sh` (and therefore CI) that
catches both failure modes, plus invalid prefixes (e.g. `/api/...`) and
unsupported schemes. The check runs in well under a second on a 66-page tree
and uses only Node built-ins.

User-visible outcome: running `pnpm run lint` (or pushing to GitHub) fails when
any redirect in `docs.json` violates the rules above, with a `file:line: message`
diagnostic pointing to the offending entry.

## Progress

Add entries as work proceeds. One row per milestone.

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

- `docs.json` — Mintlify site config. The `"redirects"` array (around
  lines 361-374 today) holds objects of shape `{source, destination}`. There is
  no `permanent` field. Three entries exist; two use `:slug*` wildcards.
- `docs/` — 66 `.mdx` files in subdirs (`admin/`, `changelogs/`, `developers/`,
  `getting-started/`, `learn/`, `references/`). File `docs/foo/bar.mdx` is
  served at URL `/docs/foo/bar`. There are no `index.mdx` files; directories
  are not pages.
- `api/`, `pub/`, `snippets/` — NOT served as pages. Redirects whose
  `destination` resolves outside `docs/` are invalid.
- `scripts/lint.sh` — Bash entry point, `set -euo pipefail`. Resolves
  `script_dir` and `repo_root` dynamically; honors `DOCS_LINT_ROOT` env var to
  point checks at fixture roots. Each sub-check prints a header
  (`echo "== <Name> =="`). Exits non-zero if any sub-check fails.
- `tests/test-lint.sh` — Bash test runner with `run_expect_pass(fixture)` and
  `run_expect_fail(fixture, pattern)` helpers. Each fixture is a self-contained
  tree at `tests/lint-fixtures/<name>/` with its own `docs/` and any other
  inputs the check reads.
- `tools/lint/` — Go-based custom MDX linter. Not modified by this plan; only
  its output format is mimicked.
- `.github/workflows/ci.yml` — `lint` job runs `pnpm run test:lint` then
  `./scripts/lint.sh`. Adding the check to `scripts/lint.sh` puts it in CI for
  free.
- `package.json` — `"lint": "./scripts/lint.sh"`,
  `"test:lint": "./tests/test-lint.sh"`. Pinned `pnpm@10.17.0` (corepack),
  Node 22.

Design choices baked into this plan:

- **Implementation language: Node.js** (`scripts/check-redirects.mjs`, ESM, no
  external deps). Node is already required for Mintlify and `pnpm`; JSON
  parsing is trivial; avoids adding a Go rule package or spreading JSON logic
  into Bash. Bash was rejected because JSON parsing in Bash needs `jq` (extra
  dep) and the wildcard / line-number logic is ugly. Go was rejected because
  the existing `tools/lint/` rules operate on MDX content; adding a JSON-only
  rule there would force everyone to rebuild Go for a 60-line check.
- **Output format**: `docs.json:<line>: <message>` to match the
  `file:line:col:` style emitted by `tools/lint/main.go`. Line numbers are
  computed by scanning `docs.json` text for the JSON object literal of each
  redirect entry. If the scan cannot locate a unique line, the script falls
  back to `docs.json: redirects[<index>]: <message>` and prints a one-line
  `warning:`-prefixed message so the limitation is visible.
- **Wildcard semantics**: a path segment matching `/^:[A-Za-z][A-Za-z0-9]*\*?$/`
  is treated as a wildcard. The "prefix path" is everything before the first
  wildcard segment. For `source`, the prefix MUST NOT be a real page and MUST
  NOT be a real directory containing any `.mdx` (or `.md`) page (otherwise
  live pages exist that the redirect cannot mask). For `destination`, the
  prefix MUST be a real directory under `docs/`.
- **Fixture strategy**: add a NEW `tests/lint-fixtures/bad-redirects/` rather
  than mutating `tests/lint-fixtures/good/` in place. The fixture starts as a
  copy of `good/` so it inherits a passing baseline for unrelated checks;
  only `docs.json` (and supporting `docs/` files) are then mutated. The real
  `docs.json` exercised by the top-level lint run already covers the passing
  case; adding a duplicate `good-redirects` fixture would be redundant.
- **Override hook**: respect `DOCS_LINT_ROOT` env var to match
  `scripts/lint.sh`'s existing convention. The script reads
  `${DOCS_LINT_ROOT}/docs.json` and resolves files under
  `${DOCS_LINT_ROOT}/docs/`. Defaults to the repo root when unset.

URL ↔ filesystem mapping rules:

1. Strip a leading `/` from the redirect path; the remainder is repo-relative.
2. Strip any trailing `/`, `?...` query string, or `#...` fragment.
3. The path MUST start with `docs/`. Anything else (`api/...`, `snippets/...`,
   bare names) is invalid.
4. Append `.mdx` to get the candidate filename; if that does not exist, fall
   back to `.md`. Filesystem checks are case-sensitive. (Today the tree is
   100% `.mdx` — `find docs -type f -name '*.md' \! -name '*.mdx'` returns
   zero — but the resolver tries both extensions to stay forward-compatible
   with future Markdown-only pages.)
5. External destinations starting with `http://` or `https://` are skipped
   (none exist today; defensive only).
6. Wildcard segments (containing `:`, optionally suffixed `*`) are not part of
   the filesystem path — only the prefix (segments before the first wildcard
   segment) is checked.

Edge cases captured from research:

- The `docs/references/device-api/v0.2.1/` and
  `docs/references/platform-api/2026-03-09/` trees are populated by
  `api/generate_event_pages.py` from OpenAPI specs. The fixtures and unit-style
  checks must NOT require running the OpenAPI generator. The real-tree run
  (against the repo's actual `docs/`) only needs the directories themselves to
  exist for wildcard prefix validation; if at any point the generator output is
  not committed, document the failure in Surprises & Discoveries and reconsider
  the wildcard rule rather than relaxing it silently.

## Plan of Work

Four milestones, each one commit. See Concrete Steps for command-level detail.

- M1: Author `scripts/check-redirects.mjs` (Node ESM, no deps).
- M2: Wire the script into `scripts/lint.sh` under a `== Redirects ==` section.
- M3: Add the `bad-redirects` fixture and a `run_expect_fail` case in
  `tests/test-lint.sh`.
- M4: Run the full lint and test stack plus preflight; fix any fallout.
- **Wire-in point detail**: appending `== Redirects ==` to `scripts/lint.sh`
  (rather than a separate workflow step) means CI picks up the check via the
  existing `lint` job — no `.github/workflows/ci.yml` change required. The
  script remains independently runnable as `node scripts/check-redirects.mjs`.

## Concrete Steps

Run all commands from `/home/ben/miru/workbench1/repos/docs/` unless noted.
Branch `feat/redirect-lint-rule` is already checked out.

### M1 — `scripts/check-redirects.mjs`

1. Create `scripts/check-redirects.mjs`. Required behavior:

   - Derive `__dirname` using the standard ESM idiom:

         import path from 'node:path';
         import { fileURLToPath } from 'node:url';
         const __dirname = path.dirname(fileURLToPath(import.meta.url));

   - Resolve `root = process.env.DOCS_LINT_ROOT || path.resolve(__dirname, '..')`.
   - Read `${root}/docs.json` as text. Parse JSON for the `redirects` array.
     Keep the original text for line-number lookup.
   - Input handling:
     - If `${root}/docs.json` does not exist, exit 0 and print
       `No docs.json at ${root}; nothing to check`. (This keeps the script
       usable against fixtures that don't ship a `docs.json`, including the
       existing `tests/lint-fixtures/good/` tree.)
     - If `${root}/docs.json` exists but cannot be parsed as JSON, exit 1 and
       print `docs.json: invalid JSON: <parser-message>` on stderr.
     - If parsed JSON has no `redirects` key, or `redirects` is an empty
       array, exit 0 and print `Checked 0 redirects: OK`.
   - For each entry, validate:
     a. Both `source` and `destination` are non-empty strings.
     b. `source` MUST start with `/`. `destination` MUST start with `/` OR
        `http://` OR `https://`. Otherwise emit
        `bad path: must start with '/'` (the message for `destination`
        additionally notes that `http(s)://` is also accepted).
     c. `destination` starting with `http://` or `https://` is skipped for
        the remaining (filesystem) checks.
     d. Strip leading `/`, trailing `/`, `?...`, `#...` from each path.
     e. The remainder MUST start with `docs/`. Otherwise emit
        `bad prefix (must start with /docs/)`.
     f. Split into segments. The first segment matching `/^:[A-Za-z][A-Za-z0-9]*\*?$/`
        is the wildcard boundary; segments before it form the `prefix`.
     g. Resolve `prefixFs = path.join(root, prefix.join('/'))`.
     h. **Source rules**:
        - If no wildcard: neither `${prefixFs}.mdx` nor `${prefixFs}.md` may
          exist (otherwise the redirect is dead).
        - If wildcard: `${prefixFs}` MUST NOT be a directory containing any
          `.mdx` or `.md` files (recursive). It MUST also NOT be a file at
          `${prefixFs}.mdx` or `${prefixFs}.md`.
     i. **Destination rules**:
        - If no wildcard: a file at `${prefixFs}.mdx` OR `${prefixFs}.md`
          MUST exist.
        - If wildcard: `${prefixFs}` MUST be an existing directory.
   - For each violation, emit `docs.json:<line>: redirects[<i>] <field>
     "<value>": <message>` on stderr. Compute `<line>` by:
     1. Locating the `"redirects"` key in the original text (e.g. via a regex
        like `/"redirects"\s*:\s*\[/`) and taking the offset of the opening `[`.
     2. From that offset onward, walking the text with an offset cursor: for
        each `"source":` literal encountered, advance the cursor past the match
        so that the n-th `"source":` literal under the `"redirects"` array
        locates the n-th redirect entry. Stop scanning at the matching closing
        `]` for the array (track bracket depth so nested `[ ]` inside string
        values do not confuse the count).
     3. Reporting the 1-based line of the located occurrence.
     Anchoring the scan to the `"redirects"` array is required because
     `docs.json` contains `"source":` keys elsewhere (e.g. in the OpenAPI
     navigation block at the time of writing, ~10 occurrences before the
     redirects). If lookup fails for any entry, emit
     `docs.json: redirects[<i>] ...` and a one-line stderr warning prefixed
     `warning:`, printed at most once per script run.
   - Exit 0 if no violations, 1 otherwise. Print a one-line summary on success
     (`Checked N redirects: OK`).

2. Make it executable and verify the shebang works:

       chmod +x scripts/check-redirects.mjs
       node scripts/check-redirects.mjs

   The script must currently pass against the real repo (the three existing
   redirects are valid given today's tree). If it does not, investigate before
   moving on — the script may be wrong, or a real redirect may already be
   broken. Record the result in Surprises & Discoveries.

3. Commit:

       git add scripts/check-redirects.mjs
       git commit -m "feat(lint): add redirect validator script for docs.json"

### M2 — Wire into `scripts/lint.sh`

1. Open `scripts/lint.sh` and add a new section (preferred near the other
   structural checks, after the MDX prose section). Use the same pattern as
   existing checks:

       echo "== Redirects =="
       node "${repo_root}/scripts/check-redirects.mjs"

   The script honors `DOCS_LINT_ROOT` itself, so no extra plumbing is needed.

2. Run `./scripts/lint.sh` end-to-end and confirm the new section appears and
   passes.

3. Commit:

       git add scripts/lint.sh
       git commit -m "feat(lint): run redirect validator from lint.sh"

### M3 — Test fixture and runner case

1. Create `tests/lint-fixtures/bad-redirects/` by copying
   `tests/lint-fixtures/good/` wholesale (it already passes every other lint
   check), then mutate it:

   - Add a `docs.json` whose `redirects` array contains, as
     **separate distinct entries** (one violation per entry — do not bundle
     multiple violations into a single redirect), at least:
     1. Source equals an existing page (dead redirect).
     2. Destination points to a missing page.
     3. Source uses a bad prefix (e.g. `/api/foo`) — must start with `/docs/`.
     4. Destination uses a bad prefix (e.g. `/api/foo`) — must start with
        `/docs/` or `http(s)://`.
     5. Source missing the leading `/` (bad path).
     6. Destination missing the leading `/` and not `http(s)://` (bad path).
     7. Wildcard destination prefix that does not exist as a directory.
     8. Wildcard source whose prefix directory contains real pages.
     Plus at least one valid redirect (and one external `https://` destination)
     so the script's "skip" and "pass mixed in with fail" paths are exercised.
   - Add only the `docs/` files the violations rely on (e.g.
     `docs/admin/exists.mdx`, `docs/wild/page.mdx`). Keep each file tiny —
     a single H1 line is enough.
   - Because the fixture starts as a copy of `good/`, it inherits the MDX/CSpell/
     OpenAPI baseline that already passes those lint checks. The `good/` tree
     has no `docs.json` of its own, so the new `docs.json` is freshly authored
     here — keep it minimal: only the `redirects` array is required for the
     redirect check (per Fix-2 input handling, missing top-level keys do not
     fail the redirect check). If a future lint check requires other `docs.json`
     keys against fixtures, copy the minimum subset from the real `docs.json`.

2. Edit `tests/test-lint.sh` to add a `run_expect_fail` invocation for the new
   fixture. For each violating entry, assert that the output contains both the
   per-entry diagnostic prefix `redirects[<i>]` (with `<i>` matching the
   entry's array index in the fixture's `docs.json`) AND the rule-specific
   substring (`dead redirect`, `missing destination`, `bad prefix`,
   `bad path: must start with '/'`, `wildcard prefix not a directory`, etc.).
   Pairing the index with the substring proves each rule fires independently
   on its own entry rather than one entry tripping multiple diagnostics. Use
   the helper's existing pattern for multi-pattern assertions; if it only
   supports one pattern, run it once per (index, substring) pair.

3. Run `./tests/test-lint.sh` and confirm the new case fails as expected and
   all previously passing fixtures still pass.

4. Commit:

       git add tests/lint-fixtures/bad-redirects tests/test-lint.sh
       git commit -m "test(lint): add bad-redirects fixture for redirect validator"

### M4 — Full validation

1. Run the entire lint and test stack:

       pnpm run test:lint
       pnpm run lint

   All must pass.

2. Run preflight:

       ./scripts/preflight.sh

   Must report clean before publishing the branch.

3. If anything was changed during validation (e.g. tweaks to the script),
   commit those fixes as a separate small commit:

       git commit -am "fix(lint): <short description>"

4. Move the plan from `plans/backlog/` to `plans/active/` once implementation
   begins, and to `plans/completed/` once merged. Fill in
   Outcomes & Retrospective at completion time.

## Validation and Acceptance

The change is accepted when **all** of the following hold:

1. `node scripts/check-redirects.mjs` exits 0 against the real
   `docs.json`/`docs/` tree and prints `Checked 3 redirects: OK` (or whatever
   the current count is).
2. `pnpm run lint` succeeds locally and shows the `== Redirects ==` section.
3. `pnpm run test:lint` succeeds, including the new `bad-redirects` case which
   must fail the lint script and surface every expected diagnostic substring.
4. CI's `lint` job succeeds on the branch (it auto-picks up the new check via
   `scripts/lint.sh`).
5. Manually injecting a known-bad redirect into a scratch copy of `docs.json`
   (e.g. `{"source": "/docs/getting-started/intro", "destination":
   "/docs/nope"}`) makes `pnpm run lint` fail with two diagnostics: one for
   the dead source and one for the missing destination. Revert the scratch
   change before committing.
6. **`./scripts/preflight.sh` reports clean before the branch is pushed for
   review.** This is mandatory and gates publication.

## Idempotence and Recovery

- **Read-only.** The script never mutates `docs.json` or any file under
  `docs/`. Re-running it is safe and produces identical output.
- **Revert-by-commit rollback.** Each milestone is one commit; revert the
  failing commit (and any later ones) to roll back cleanly. No data migration
  is involved.
- **Fixture-narrowing tip.** To isolate a failure, point `DOCS_LINT_ROOT` at
  a single fixture (e.g. `DOCS_LINT_ROOT=tests/lint-fixtures/bad-redirects
  node scripts/check-redirects.mjs`) so the script only inspects that tree.

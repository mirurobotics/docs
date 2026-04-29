# Port redirect lint check to the Go custom linter

This ExecPlan is a living document. The sections Progress, Surprises &
Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to
date as work proceeds.

This plan supersedes `plans/active/20260428-redirect-lint-rule.md`. That
plan implemented the redirect validator in Node.js
(`scripts/check-redirects.mjs`). This plan refactors it into a proper rule
inside the Go custom linter at `tools/lint/`. The previous plan stays as
historical record; this one captures the corrected architectural choice.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `mirurobotics/docs` (`/home/ben/miru/workbench1/repos/docs/`) | read-write | Branch `feat/redirect-lint-rule` is checked out; PR #75 is open. New Go package at `tools/lint/linter/redirects/`; small wiring change in `tools/lint/main.go`; deletes `scripts/check-redirects.mjs` and the `== Redirects ==` block in `scripts/lint.sh`; updates the prior plan's living-doc sections. |

This plan lives in `mirurobotics/docs` because all changes are confined to
that repo. Out of scope: any other repo, OpenAPI page regeneration,
modifying the redirects in `docs.json` themselves, touching other lint rules.

## Purpose / Big Picture

After this change, the redirect validator runs as a normal rule inside the
Go custom linter (`tools/lint/`), alongside `importresolves`,
`componentstyle`, `mdxstyle`, etc. The user-visible behavior is unchanged
versus the Node.js implementation:

- `pnpm run lint` still fails on dead/missing/malformed redirects with
  the same diagnostic substrings.
- The 11 existing `run_expect_fail` assertions in `tests/test-lint.sh`
  against the `bad-redirects` fixture continue to pass without
  modification — that is the binding contract for output format.

What is gained by the refactor:

- One linter, one language. The project already has a Go linter with rule
  registration, filesystem-aware sibling rules (`importresolves`), 100%
  coverage gates per rule package, and CI integration. A second linter
  written in Node.js for one small check is not justified.
- New Go package gets the standard 100% covgate automatically.
- No new CI wiring — the existing `lint-custom-linter` and
  `test-custom-linter` jobs are gated on the `tools/lint/**` path filter
  and pick the new package up automatically.

## Progress

- [ ] M1: Create `tools/lint/linter/redirects/` package (rule + helpers + tests + `.covgate=100`).
- [ ] M2: Wire `redirects.Check(contentRoot)` into `tools/lint/main.go`.
- [ ] M3: Delete `scripts/check-redirects.mjs` and remove its block from `scripts/lint.sh`.
- [ ] M4: Confirm `tests/lint-fixtures/bad-redirects/` still produces all 11 expected diagnostics from the Go linter.
- [ ] M5: Update prior plan's Decision Log + Surprises & Discoveries entries; run preflight; fix any fallout.

Use timestamps when steps complete. Split partial work into "done" and
"remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

- Observation: …
  Evidence: …

## Decision Log

(Add entries as you go.)

- Decision: …
  Rationale: …
  Date/Author: …

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

Working directory for every command in this plan is
`/home/ben/miru/workbench1/repos/docs/` unless stated otherwise. Branch
`feat/redirect-lint-rule` is already checked out; PR #75 is open against
`main`. Today's date (UTC) is 2026-04-28.

### The Go custom linter (`tools/lint/`)

`tools/lint/go.mod` declares Go 1.25.3, module
`github.com/mirurobotics/docs/tools/lint`, no runtime deps beyond the Go
standard library (dev tools include golangci-lint v2, gofumpt, and the
miru CLI for lint/covgate scripts).

Entry point: `tools/lint/main.go`.

- `findContentRoot` walks upward from the directory of the first input
  file looking for a `snippets/` sibling. The discovered directory is
  `contentRoot`.
- The `main` function iterates `os.Args[1:]` and calls
  `linter.ProcessFile(path, contentRoot)` per file, accumulating
  `analysis.Violation` records.
- It prints `file:line:col: message` per violation to stdout and exits
  `0` (clean) / `1` (violations) / `2` (usage).
- `tools/lint/main_test.go` already covers `findContentRoot`.

Rule registration: `tools/lint/linter/run.go`.

- `ruleCheckers()` returns a slice of `ruleEntry{rule Rule, check
  func(checkInput) []analysis.Violation}`.
- `checkInput` carries `path`, `lines`, `spans`, `contentRoot`.
- `analysis.Violation` fields: `File`, `Line`, `Col`, `Message`.
- All current rules are per-file; no global/once-per-run rule concept
  exists.

Sibling rule precedent: `tools/lint/linter/importresolves/`.

- Files: `importresolves.go`, `importresolves_test.go`, `.covgate`
  (`100.0`).
- Signature: `func Check(file string, lines []string, contentRoot
  string) []analysis.Violation`.
- Uses `os.Stat(filepath.Join(contentRoot, importPath))` for
  filesystem-aware checks — the same pattern this plan needs.
- Tests build fixtures in-memory using `t.TempDir()`, `os.WriteFile`,
  `os.MkdirAll`. There is no `testdata/` directory.

Other rule packages (8 total): `analysis` (utilities, 97.8% covgate),
`componentstyle`, `importblock`, `importresolves`, `importsorted`,
`importused`, `mdxstyle`, `nodoubledash` — each at 100% covgate. The
top-level `linter` package is 93.3%.

### Build and CI

- `tools/lint/scripts/lint.sh` runs `go tool miru lint --paths=...`.
- `tools/lint/scripts/covgate.sh` runs
  `go tool miru covgate --packages="./..." --default-threshold="${1:-90.0}"`.
  Per-package overrides come from each package's `.covgate` file.
- `.github/workflows/ci.yml` jobs `lint-custom-linter` and
  `test-custom-linter` are gated on the `tools/lint/**` path filter, so
  a new package under `tools/lint/linter/redirects/` is picked up
  automatically — no workflow change required.

### The current Node.js implementation

`scripts/check-redirects.mjs` (335 lines) plus its invocation in
`scripts/lint.sh`:

- `scripts/lint.sh` line 65-66 invokes the Go linter against `.mdx`
  files (and snippets).
- `scripts/lint.sh` line 81 runs `node scripts/check-redirects.mjs`
  under an `echo "== Redirects =="` heading. This wire-in is what is
  removed by M3.

`tests/test-lint.sh` invokes `scripts/lint.sh` with
`DOCS_LINT_ROOT=tests/lint-fixtures/<name>` per fixture. The
`bad-redirects` fixture has 11 `run_expect_fail` assertions checking
substrings such as:

    redirects[0] source "/docs/admin/exists": dead redirect (source resolves to a real page)
    redirects[1] destination "/docs/admin/gone": missing destination (no .mdx or .md page exists)

The Go rule MUST emit identical strings so these assertions pass
unchanged. That is the contract.

### Rule branches to port (verbatim from current behavior)

Seven distinct branches plus the OpenAPI escape hatch, plus the
non-object case:

1. (a) Empty / non-string `source` or `destination` →
   `must be a non-empty string`.
2. (b) Leading-slash requirement: `source` MUST start with `/`;
   `destination` MUST start with `/` or `http(s)://`. Otherwise emit
   `bad path: must start with '/'` (the destination message also notes
   `(or http(s)://)`).
3. (e) Bad prefix: after stripping leading `/`, path MUST start with
   `docs/`. Otherwise emit `bad prefix (must start with /docs/)`.
4. Source non-wildcard: file at `${prefixFs}.mdx` or `.md` MUST NOT
   exist → `dead redirect (source resolves to a real page)`.
5. Source wildcard: prefix MUST NOT be a directory containing any
   `.mdx`/`.md` (recursive) AND MUST NOT be a file at
   `${prefixFs}.mdx`/`.md` →
   `dead redirect (wildcard source prefix has real pages)` or
   `dead redirect (wildcard source prefix resolves to a real page)`.
6. Destination non-wildcard: file at `${prefixFs}.mdx` or `.md` MUST
   exist → `missing destination (no .mdx or .md page exists)`.
7. Destination wildcard: `${prefixFs}` MUST be an existing directory OR
   `${prefixFs}.yaml` MUST be referenced as `nav.*.openapi.source`
   somewhere in `docs.json` → `wildcard prefix not a directory`.
8. Non-object entry → `not an object`. Emit with `field="entry"` and
   `value=""`.

Wildcard segment regex: `^:[A-Za-z][A-Za-z0-9]*\*?$`. Segments BEFORE
the first wildcard segment form the "prefix".

### Diagnostic format

The full `Violation.Message` (everything after the `file:line:col:`
prefix) is:

    redirects[<i>] <field> "<value>": <message>

Examples of the entire `Message` field (`redirects[0]` is the FIRST
token of the message, NOT a separate prefix):

    redirects[0] source "/docs/admin/exists": dead redirect (source resolves to a real page)
    redirects[1] destination "/docs/admin/gone": missing destination (no .mdx or .md page exists)
    redirects[10] source "": must be a non-empty string
    redirects[11] entry "": not an object

Format string: `fmt.Sprintf("redirects[%d] %s %q: %s", index, field,
value, message)`. Verify `%q` matches the Node.js double-quoted form for
the ASCII inputs used in the fixture (it does for ASCII).

### Line anchoring

To compute the `line` field, locate the `"redirects"\s*:\s*\[` literal
in `docs.json` text. From the `[` offset, the n-th `"source":` literal
encountered corresponds to the n-th redirect entry. Bound the scan by
the parsed redirect count. On lookup failure (e.g. malformed entries
that lack a `"source":` key), fall back to `Line: 1` so output stays
deterministic and `file:line:col:` formatting is preserved. No stderr
warning is emitted; the fallback case is rare and the violation message
itself is sufficient.

### Where `docs.json` is found

The current `findContentRoot` walks up from the first input file's
directory. The redirect rule needs `contentRoot` even when the linter is
invoked without `.mdx` files (e.g. against the `bad-redirects` fixture
that contains no `.mdx` of its own — though in practice it inherits
`docs/example.mdx` from `good/`, so contentRoot is bootstrappable from
the per-file loop's first arg).

For robustness when no file args are given: respect the
`DOCS_LINT_ROOT` env var (matches the existing test runner convention)
and treat that as `contentRoot`. Otherwise require at least one file
arg so the existing walk works.

For `scripts/lint.sh`: it already passes `.mdx` files to the linter, so
`contentRoot` is found via the existing path; the redirect rule then
runs alongside.

For `tests/test-lint.sh`: it sets `DOCS_LINT_ROOT` and invokes
`scripts/lint.sh`; the per-file loop discovers `contentRoot` from the
fixture's inherited `.mdx` files. The redirect rule then runs against
`${contentRoot}/docs.json`.

### Why the Node.js implementation was wrong

It introduced a second linter in a second language for one small check
that fits naturally inside the existing Go linter. The Go linter
already has filesystem-aware rules, per-package coverage gates, and a
dedicated CI job — none of which the Node.js script benefited from.

## Plan of Work

### M1 — New Go package `tools/lint/linter/redirects/`

Files to create:

- `tools/lint/linter/redirects/redirects.go`
- `tools/lint/linter/redirects/redirects_test.go`
- `tools/lint/linter/redirects/.covgate` containing `100.0`

Public API of the package:

    package redirects

    import "github.com/mirurobotics/docs/tools/lint/linter/analysis"

    // Check reads ${contentRoot}/docs.json and returns violations for
    // dead/missing/malformed redirects. If docs.json is absent it
    // returns nil. If docs.json is present but unparseable, it returns
    // a single violation with a parse-error message.
    func Check(contentRoot string) []analysis.Violation

Internal structure (all unexported):

- `validate(docsJSONBytes []byte, contentRoot string) []analysis.Violation`
  — pure function; the tests target this so they don't have to write
  `docs.json` to a temp dir for every case.
- `cleanPath(p string) string` — strips a leading `/`, trailing `/`,
  `?...`, `#...`.
- `splitWildcard(segments []string) (prefix []string, hasWildcard bool)`
  — uses regex `^:[A-Za-z][A-Za-z0-9]*\*?$`.
- `validateSource(i int, source, contentRoot string) []analysis.Violation`
- `validateDestination(i int, destination, contentRoot string) []analysis.Violation`
- `collectOpenAPISources(docsJSON map[string]any) map[string]bool` —
  walks `nav.*.openapi.source` recursively (any depth) and returns the
  set of source paths (as found in JSON; relative to repo root).
- `lineLookup(docsJSONText string, count int) []int` — returns 1-based
  line numbers for each redirect entry, with a fallback of `1` per
  unmatched entry.
- `formatMessage(i int, field, value, message string) string` —
  returns `fmt.Sprintf("redirects[%d] %s %q: %s", i, field, value,
  message)`.

The exported `Check(contentRoot string)`:

1. Reads `${contentRoot}/docs.json`. If missing, returns `nil`.
2. On unmarshal error, returns one violation with `Line: 1`, `Col: 1`,
   `File: "docs.json"`, `Message: "invalid JSON: <err>"`.
3. Calls `validate(bytes, contentRoot)`.

Tests (`redirects_test.go`) are table-driven against `validate`. Each
case constructs:

- `docsJSON []byte` literal (small inline snippet).
- A `t.TempDir()` prepopulated by `os.MkdirAll` / `os.WriteFile` to
  reflect the on-disk pages required by the case.
- Expected `[]analysis.Violation` slice.

Cover at minimum:

- All seven rule branches above (each both positive and negative).
- Non-object entry (`"some-string"` inside the array).
- OpenAPI escape hatch (wildcard destination resolved by a
  `nav.*.openapi.source: docs/.../foo.yaml` reference).
- `redirects` key absent, `redirects` empty, `redirects` not an array
  (defensive).
- Line-number anchoring: redirect entries return correct 1-based lines
  and fallback to `1` when no `"source":` key (non-object entry case).
- Trailing `/`, `?...`, `#...` stripped.
- `http(s)://` destinations skipped from filesystem check.

Coverage gate: the `.covgate` file MUST contain `100.0`. Run
`./tools/lint/scripts/covgate.sh` from `tools/lint/` to confirm.

### M2 — Wire into `tools/lint/main.go`

In `tools/lint/main.go`:

1. Import `github.com/mirurobotics/docs/tools/lint/linter/redirects`.
2. Determine `contentRoot`:
   - If at least one file arg is present, keep current `findContentRoot`
     bootstrap from it.
   - If no file args but `DOCS_LINT_ROOT` is set in the environment,
     use that as `contentRoot`.
   - Otherwise behave exactly as today (usage exit 2).
3. Once `contentRoot` is known, call `redirects.Check(contentRoot)` and
   append its violations to the `allViolations` slice. Order: redirect
   violations are appended after the per-file pass so per-file output
   order is preserved.
4. Print path is unchanged: each violation goes to stdout as
   `File:Line:Col: Message`. The `redirects.Check` violations have
   `File: "docs.json"` and an integer line.

Update `tools/lint/main_test.go` if the `findContentRoot` /
`DOCS_LINT_ROOT` handling reaches into a tested function. If the new
logic is in `main()` itself (where it is hard to unit-test), exercise
it via the integration test in M4 instead.

### M3 — Remove the Node.js implementation

1. Delete `scripts/check-redirects.mjs`.
2. In `scripts/lint.sh`, remove the `== Redirects ==` block (lines 80-81
   in current HEAD; verify with `git diff` before commit).
3. Run `pnpm run lint` against the real repo and confirm:
   - The `== Lint ==` (Go linter) section now shows no redirect
     violations against the real `docs.json`.
   - No `== Redirects ==` section remains.
   - Overall exit code is 0.

### M4 — Verify the existing fixture still produces all 11 diagnostics

The fixture at `tests/lint-fixtures/bad-redirects/` and the 11
`run_expect_fail` assertions in `tests/test-lint.sh` (lines 58-68) MUST
keep passing without modification. They serve as the integration-level
contract test for the Go rule.

If any assertion fails, the Go rule's diagnostic format diverges from
the Node.js script's. Fix the Go rule, not the assertions.

Do not modify `tests/lint-fixtures/bad-redirects/` or any of the 11
assertions. Their purpose for the duration of this plan is to detect
behavioral drift.

### M5 — Update the superseded plan and run preflight

1. Edit `plans/active/20260428-redirect-lint-rule.md`:
   - Append a Decision Log entry dated 2026-04-28 stating the
     implementation was refactored from Node.js to Go (see this plan).
     Rationale: the project already has a Go linter with rule
     registration, filesystem-aware sibling rules (`importresolves`),
     100% coverage gates per rule package, and CI integration. Two
     linters in two languages was unjustified.
   - Append a Surprises & Discoveries entry noting the architectural
     mistake and the recovery (a follow-up plan rather than a rewrite
     of this one, so the audit trail is preserved).
   - Do NOT move the plan to `plans/completed/` here. That happens at
     the very end of the combined work after preflight is clean.
2. Run `./scripts/preflight.sh` and resolve any findings. Preflight
   must report clean before the branch is published.

### Commit strategy

Make a series of new commits on top of the existing branch. Do NOT
rewrite history. The Node.js commits remain in the audit trail; the
refactor commits supersede them. Reviewers see the migration as a clear
sequence of steps. One commit per milestone — five commits total.

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs/` unless
stated otherwise.

### M1 — Create the Go rule package

1. Create the directory and stub files:

       mkdir -p tools/lint/linter/redirects
       printf '100.0\n' > tools/lint/linter/redirects/.covgate

2. Author `tools/lint/linter/redirects/redirects.go` with the public
   `Check(contentRoot string) []analysis.Violation` and helpers
   described in Plan of Work / M1. Match the file/import style of the
   sibling `importresolves.go` package.

3. Author `tools/lint/linter/redirects/redirects_test.go` as a
   table-driven test against the unexported `validate`. Use
   `t.TempDir()` + `os.MkdirAll` + `os.WriteFile` for filesystem
   fixtures. Cover every rule branch listed in Plan of Work.

4. Run package tests and coverage:

       cd tools/lint
       go test ./linter/redirects/...
       ./scripts/covgate.sh

   Expected: tests pass; covgate reports `redirects` at >= 100.0
   (i.e., no failure for that package).

5. Run formatters and the custom-linter lint:

       cd tools/lint
       ./scripts/lint.sh

   Resolve any findings until clean.

6. Commit (one commit per milestone — see policy):

       git add tools/lint/linter/redirects
       git commit -m "feat(lint): add redirects rule package to Go linter"

### M2 — Wire the rule into main.go

1. Edit `tools/lint/main.go` per Plan of Work / M2.

2. Run main_test and the redirects test together:

       cd tools/lint
       go test ./...
       ./scripts/covgate.sh

   Expected: all green. The top-level `linter`/`main` coverage may
   shift slightly but should remain above its `.covgate` (or the
   default 90.0 if no override).

3. Smoke-test against the real repo:

       cd /home/ben/miru/workbench1/repos/docs
       ./scripts/lint.sh

   Expected: exits 0; no redirect violations against the real
   `docs.json`.

4. Commit:

       git add tools/lint/main.go tools/lint/main_test.go
       git commit -m "feat(lint): invoke redirects rule from Go linter main"

### M3 — Remove the Node.js implementation

1. Inspect the wire-in to confirm line numbers before editing:

       grep -n 'check-redirects\|== Redirects ==' scripts/lint.sh

2. Edit `scripts/lint.sh`: remove the `echo "== Redirects =="` line and
   the following `node ... check-redirects.mjs` invocation. Do not
   touch any other section.

3. Delete the Node.js script:

       git rm scripts/check-redirects.mjs

4. Verify `pnpm run lint` and `pnpm run test:lint` against the real
   repo:

       pnpm run lint
       pnpm run test:lint

   Expected: both exit 0. Lint output shows the Go linter section but
   no `== Redirects ==` block.

5. Commit:

       git add scripts/lint.sh
       git commit -m "chore(lint): remove Node.js redirect script (now in Go linter)"

### M4 — Verify the fixture assertions still pass

1. Re-run the test runner to confirm all 11 `run_expect_fail`
   assertions for `bad-redirects` still pass:

       ./tests/test-lint.sh

   Expected: full pass. Each of the 11 expected substrings is matched
   by the Go linter's stdout for the `bad-redirects` fixture.

2. If any assertion fails, do NOT modify `tests/test-lint.sh` or the
   fixture — fix the Go rule's output to match the documented format.
   Re-run M1's `go test`, then re-run `./tests/test-lint.sh`. When all
   green, commit:

       # Only if the rule needed an output-format fix:
       git commit -am "fix(lint): align redirects rule output with fixture assertions"

   If no fix was required, no commit is made for M4. Mark M4 done in
   Progress with the `./tests/test-lint.sh` transcript noted in
   Surprises & Discoveries (or leave Progress checked with a brief
   note that no fixup was required).

### M5 — Update the superseded plan and run preflight

1. Edit `plans/active/20260428-redirect-lint-rule.md`:
   - Append the Decision Log entry described in Plan of Work / M5.
   - Append the Surprises & Discoveries entry described in Plan of
     Work / M5.

2. Run preflight:

       ./scripts/preflight.sh

   Resolve findings and re-run until it reports clean. Preflight MUST
   report clean before publishing. This is the binding gate.

3. Commit:

       git add plans/active/20260428-redirect-lint-rule.md
       git commit -m "docs(plans): note redirect rule refactored from Node.js to Go"

   If preflight required other fixups, commit them separately with a
   `fix(...)` or `chore(...)` message before the plan-update commit.

4. Push and verify CI on PR #75:

       git push

   Expected: GitHub Actions `lint`, `lint-custom-linter`, and
   `test-custom-linter` jobs all pass.

5. After merge, move BOTH plans to `plans/completed/` (this plan and
   the superseded one). Fill Outcomes & Retrospective in this plan
   summarizing what was achieved and any lessons.

## Validation and Acceptance

Accepted when ALL of the following hold. Each is something a human can
run and observe:

1. From `tools/lint/`:

       go test ./...
       ./scripts/covgate.sh
       ./scripts/lint.sh

   All three exit 0. The new `redirects` package shows >= 100.0
   coverage in covgate output.

2. From the repo root:

       pnpm run lint

   Exits 0. Output contains the Go linter section and does NOT contain
   `== Redirects ==`.

3. From the repo root:

       pnpm run test:lint

   Exits 0. The `bad-redirects` fixture still triggers all 11
   `run_expect_fail` assertions in `tests/test-lint.sh` without those
   assertions being modified. (This is the binding contract: the Go
   rule's output substrings match the Node.js implementation's
   exactly.)

4. Manually corrupting a copy of `docs.json` by injecting

       {"source": "/docs/getting-started/intro", "destination": "/docs/nope"}

   and running `pnpm run lint` produces two diagnostics:
   - `redirects[<n>] source "/docs/getting-started/intro": dead redirect (source resolves to a real page)`
   - `redirects[<n>] destination "/docs/nope": missing destination (no .mdx or .md page exists)`

   Revert the change before committing.

5. CI on PR #75 passes — specifically the `lint`, `lint-custom-linter`,
   and `test-custom-linter` jobs.

6. **`./scripts/preflight.sh` reports clean before push.** Mandatory.
   This gate cannot be skipped.

## Idempotence and Recovery

- **Read-only against content.** The rule never mutates `docs.json` or
  files under `docs/`. Re-runs are safe and identical.
- **Per-milestone commits enable rollback.** Each milestone is one
  commit; revert to roll back. No data migration.
- **Re-run safety.** All listed commands (`go test`, `./scripts/lint.sh`,
  `./scripts/covgate.sh`, `pnpm run lint`, `pnpm run test:lint`,
  `./scripts/preflight.sh`) are idempotent and can be invoked any number
  of times.
- **Recovery from a broken `tools/lint/main.go` edit.** If M2 leaves the
  binary unable to bootstrap `contentRoot`, run
  `git checkout tools/lint/main.go` and redo M2 from a clean state. The
  M1 package commit is independent and need not be touched.
- **Recovery from a broken rule output.** If M4 surfaces an output
  mismatch, fix the rule package (M1 code) — do not modify the fixture
  or the assertions. Commit the fix on top per M4's instructions.
- **Fixture-narrowing tip.** Point the Go linter directly at the fixture
  by setting `DOCS_LINT_ROOT=tests/lint-fixtures/bad-redirects` and
  invoking the linter binary against the inherited `.mdx` files; the
  redirect rule then runs against
  `tests/lint-fixtures/bad-redirects/docs.json`.

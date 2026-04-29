# Port redirect lint check to the Go custom linter

This ExecPlan is a living document; Progress, Surprises & Discoveries,
Decision Log, and Outcomes & Retrospective must be kept current.

Supersedes `plans/active/20260428-redirect-lint-rule.md` (which
implemented the redirect validator in Node.js as
`scripts/check-redirects.mjs`); refactors it into a rule inside the Go
custom linter at `tools/lint/`.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `mirurobotics/docs` (`/home/ben/miru/workbench1/repos/docs/`) | read-write | Branch `feat/redirect-lint-rule`, PR #75 open. New Go package `tools/lint/linter/redirects/`; wiring change in `tools/lint/main.go`; delete `scripts/check-redirects.mjs` and `== Redirects ==` block in `scripts/lint.sh`; update prior plan's living-doc sections. |

Out of scope: other repos, OpenAPI page regeneration, redirect content
in `docs.json`, other lint rules.

## Purpose / Big Picture

The redirect validator runs as a normal rule inside the Go custom
linter (`tools/lint/`), alongside `importresolves`, `componentstyle`,
`mdxstyle`, etc. User-visible behavior is unchanged: `pnpm run lint`
still fails on dead/missing/malformed redirects with the same
diagnostic substrings, and the 11 `run_expect_fail` assertions in
`tests/test-lint.sh` against `bad-redirects` pass without modification
— the binding contract for output format. Gains: one linter / one
language; standard 100% covgate; no new CI wiring (existing
`lint-custom-linter` / `test-custom-linter` jobs already path-filter on
`tools/lint/**`).

## Progress

- [ ] M1: Create `tools/lint/linter/redirects/` package.
- [ ] M2: Wire `redirects.Check(contentRoot)` into `tools/lint/main.go`.
- [ ] M3: Remove `scripts/check-redirects.mjs` and its block from `scripts/lint.sh`.
- [ ] M4: Confirm `bad-redirects` fixture produces all 11 expected diagnostics.
- [ ] M5: Update prior plan; run preflight; fix fallout.

Use timestamps when steps complete; split partial work into "done"
and "remaining" as needed.

## Surprises & Discoveries

Add entries as work proceeds.

## Decision Log

Add entries as work proceeds.

## Outcomes & Retrospective

Add entries as work proceeds.

## Context and Orientation

Working directory for every command is
`/home/ben/miru/workbench1/repos/docs/` unless stated otherwise. Branch
`feat/redirect-lint-rule` is checked out; PR #75 is open against
`main`. Today is 2026-04-28.

### The Go custom linter (`tools/lint/`)

`go.mod`: Go 1.25.3, module
`github.com/mirurobotics/docs/tools/lint`, stdlib-only at runtime.

`tools/lint/main.go`: `findContentRoot` walks upward from the first
input file's directory looking for a `snippets/` sibling (that
directory is `contentRoot`). `main` iterates `os.Args[1:]` calling
`linter.ProcessFile(path, contentRoot)`, prints `file:line:col:
message` per violation, exits `0`/`1`/`2` (clean/violations/usage).
`tools/lint/main_test.go` covers `findContentRoot`.

Rule registration in `tools/lint/linter/run.go`: `ruleCheckers()`
returns `[]ruleEntry{rule Rule, check func(checkInput)
[]analysis.Violation}`. `checkInput` carries `path`, `lines`, `spans`,
`contentRoot`. `analysis.Violation` fields: `File`, `Line`, `Col`,
`Message`. All current rules are per-file; no global/once-per-run
concept exists.

Sibling rule `tools/lint/linter/importresolves/` (`.covgate` 100.0):
signature `func Check(file string, lines []string, contentRoot string)
[]analysis.Violation`; uses `os.Stat(filepath.Join(contentRoot,
importPath))` — the filesystem-aware pattern this plan needs. Tests
build fixtures in-memory via `t.TempDir()`, `os.WriteFile`,
`os.MkdirAll`; no `testdata/`. Other 100%-covgate packages:
`componentstyle`, `importblock`, `importsorted`, `importused`,
`mdxstyle`, `nodoubledash`. `analysis` 97.8%, top-level `linter` 93.3%.

### Build and CI

`tools/lint/scripts/lint.sh` runs `go tool miru lint --paths=...`.
`tools/lint/scripts/covgate.sh` runs `go tool miru covgate
--packages="./..." --default-threshold="${1:-90.0}"`; per-package
overrides come from `.covgate` files. CI jobs `lint-custom-linter`
and `test-custom-linter` in `.github/workflows/ci.yml` are gated on
the `tools/lint/**` path filter.

### The current Node.js implementation

`scripts/check-redirects.mjs` (335 lines) is invoked from
`scripts/lint.sh` line 81 under `echo "== Redirects =="` (M3 removes
this); `scripts/lint.sh` line 65-66 invokes the Go linter against
`.mdx` files. `tests/test-lint.sh` invokes `scripts/lint.sh` with
`DOCS_LINT_ROOT=tests/lint-fixtures/<name>`; `bad-redirects` has 11
`run_expect_fail` substring assertions (examples below).

### Rule branches to port (verbatim)

The message strings below are the contract with `tests/test-lint.sh`
lines 58-68. Do not paraphrase. Single quotes in messages 3-4 are
LITERAL (NOT from `%q`).

1. Empty / non-string `source` or `destination` → `must be a non-empty string`.
2. Non-object entry → emit with `field="entry"`, `value=""`, message: `not an object`.
3. Source missing leading `/` → `bad path: must start with '/'`.
4. Destination missing leading `/` and not `http(s)://` → `bad path: must start with '/' (or http(s)://)`.
5. After stripping leading `/`, path doesn't start with `docs/` → `bad prefix (must start with /docs/)`.
6. Source non-wildcard, `${prefixFs}.mdx` or `.md` exists → `dead redirect (source resolves to a real page)`.
7. Source WILDCARD:
   - Prefix is a directory containing `.mdx`/`.md` pages (recursive) → `dead redirect (wildcard source prefix has real pages)`.
   - Prefix is a file (`.mdx`/`.md` exists at the prefix path) → `dead redirect (wildcard source prefix resolves to a real page)`.
8. Destination non-wildcard, no `${prefixFs}.mdx` or `.md` → `missing destination (no .mdx or .md page exists)`.
9. Destination wildcard, prefix not a directory and `${prefixFs}.yaml` not referenced as `nav.*.openapi.source` anywhere in `docs.json` → `wildcard prefix not a directory`.

Wildcard segment regex `^:[A-Za-z][A-Za-z0-9]*\*?$`; segments BEFORE
the first wildcard form the "prefix".

### Diagnostic format

`Violation.Message` (full message after `file:line:col:`) is
`redirects[<i>] <field> "<value>": <message>` — `redirects[0]` is the
FIRST token, not a separate prefix. Format string:
`fmt.Sprintf("redirects[%d] %s %q: %s", index, field, value, message)`.
`%q` matches the Node.js double-quoted form for ASCII fixture inputs.
Additional fixture shapes:

    redirects[10] source "": must be a non-empty string
    redirects[11] entry "": not an object

### Line anchoring

Locate the `"redirects"\s*:\s*\[` literal in `docs.json` text. From the
`[` offset, the n-th `"source":` corresponds to the n-th redirect
entry. Bound the scan by the parsed redirect count. On lookup failure
(e.g. non-object entries with no `"source":` key), fall back to
`Line: 1` silently.

### Where `docs.json` is found

`contentRoot` comes from `findContentRoot`'s walk-up from the first
MDX argument; all invocation paths pass MDX files. `bad-redirects`
inherits `docs/example.mdx` from `good/`. `DOCS_LINT_ROOT` is consumed
by `scripts/lint.sh` only; the Go binary does NOT read it.

## Plan of Work

Five milestones, each one commit:

- M1: Create `tools/lint/linter/redirects/` package.
- M2: Wire `redirects.Check(contentRoot)` into `tools/lint/main.go`.
- M3: Delete `scripts/check-redirects.mjs`; remove its block from `scripts/lint.sh`.
- M4: Verify `bad-redirects` fixture produces all 11 expected diagnostics.
- M5: Update prior plan; run preflight; fix fallout.

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs/`.

### M1 — Create the Go rule package

Files: `tools/lint/linter/redirects/redirects.go`,
`redirects_test.go`, `.covgate` (containing `100.0`).

Public API in `package redirects`, import path
`github.com/mirurobotics/docs/tools/lint/linter/analysis`:

    func Check(contentRoot string) []analysis.Violation

`Check` reads `${contentRoot}/docs.json` (missing → nil); on unmarshal
error returns one violation `{File: "docs.json", Line: 1, Col: 1,
Message: "invalid JSON: <err>"}`; else calls
`validate(bytes, contentRoot)`. `Col` is always `1`.

Unexported helpers: `validate(bytes, contentRoot)` (pure, tested
directly); `cleanPath` (strips leading `/`, trailing `/`, `?...`,
`#...`); `splitWildcard` (regex `^:[A-Za-z][A-Za-z0-9]*\*?$`);
`validateSource`, `validateDestination`; `collectOpenAPISources(map)
map[string]bool` (walks `nav.*.openapi.source` at any depth);
`lineLookup(text, count) []int` (1-based, fallback `1`);
`formatMessage(i, field, value, msg)` returning
`fmt.Sprintf("redirects[%d] %s %q: %s", i, field, value, message)`.

Tests are table-driven against `validate`. Each case constructs a
`docsJSON []byte` literal, a `t.TempDir()` populated via `os.MkdirAll`
/ `os.WriteFile`, and an expected `[]analysis.Violation`. Cover: all
seven rule branches (positive and negative); non-object entry
(`"some-string"`); OpenAPI escape hatch (wildcard destination resolved
by `nav.*.openapi.source: docs/.../foo.yaml`); wildcard destination
prefix neither directory nor registered yaml — pin
`wildcard prefix not a directory` (`redirects[12]` fixture); wildcard
source prefix resolving to an `.mdx` file (not a directory) — pin
`dead redirect (wildcard source prefix resolves to a real page)`
(`redirects[6]` fixture); `redirects` key absent / empty / non-array;
line anchoring (correct lines, fallback `1` for non-object); trailing
`/`, `?...`, `#...` stripped; `http(s)://` destinations skipped from
filesystem check.

Commands:

1. Stub the package:

       mkdir -p tools/lint/linter/redirects
       printf '100.0\n' > tools/lint/linter/redirects/.covgate

2. Author `redirects.go` and `redirects_test.go`, matching
   `importresolves.go` style.

3. Run tests/coverage/lint (covgate must report `redirects` >= 100.0):

       cd tools/lint
       go test ./linter/redirects/...
       ./scripts/covgate.sh
       ./scripts/lint.sh

4. Commit:

       git add tools/lint/linter/redirects
       git commit -m "feat(lint): add redirects rule package to Go linter"

### M2 — Wire the rule into main.go

1. Edit `tools/lint/main.go`: import
   `github.com/mirurobotics/docs/tools/lint/linter/redirects`. After
   `contentRoot` is discovered, call `redirects.Check(contentRoot)` and
   append violations to `allViolations` after per-file violations
   (preserving output order). Redirect violations carry
   `File: "docs.json"` and use the existing print path. Do NOT add a
   `DOCS_LINT_ROOT` fallback in `main.go`.

2. Add `RuleRedirects = "redirects"` to `tools/lint/linter/run.go`'s
   `Rule` block and extend `AllRules()` to include it. Add a one-line
   comment near `ruleCheckers()`:
   `// Redirects is invoked once per run from main.go (see linter.ProcessDocsJSON), not per-file via ruleCheckers, because it operates on docs.json once.`

3. Run tests and smoke-test (all exit 0):

       cd tools/lint && go test ./... && ./scripts/covgate.sh
       cd /home/ben/miru/workbench1/repos/docs && ./scripts/lint.sh

4. Commit:

       git add tools/lint/main.go tools/lint/main_test.go tools/lint/linter/run.go
       git commit -m "feat(lint): invoke redirects rule from Go linter main"

### M3 — Remove the Node.js implementation

1. Confirm location, then edit `scripts/lint.sh` to remove the
   `echo "== Redirects =="` line and the following
   `node ... check-redirects.mjs` invocation:

       grep -n 'check-redirects\|== Redirects ==' scripts/lint.sh

2. Delete the script and verify (both exit 0; no `== Redirects ==`):

       git rm scripts/check-redirects.mjs
       pnpm run lint
       pnpm run test:lint

3. Commit:

       git add scripts/lint.sh
       git commit -m "chore(lint): remove Node.js redirect script (now in Go linter)"

### M4 — Verify the fixture assertions still pass

1. Run; all 11 `run_expect_fail` assertions for `bad-redirects` must
   pass:

       ./tests/test-lint.sh

2. If any fail, do NOT modify `tests/test-lint.sh` or the fixture —
   fix the Go rule's output, re-run M1's `go test`, then re-run
   `./tests/test-lint.sh`. When green (only if a fix was required):

       git commit -am "fix(lint): align redirects rule output with fixture assertions"

### M5 — Update the superseded plan and run preflight

1. Edit `plans/active/20260428-redirect-lint-rule.md`: append a
   Decision Log entry dated 2026-04-28 stating the Node.js→Go refactor
   (rationale: existing Go linter with rule registration,
   filesystem-aware sibling rules, 100% covgate, CI integration); and a
   Surprises & Discoveries entry noting the architectural mistake and
   recovery via follow-up plan. Do NOT move the plan to
   `plans/completed/` here.

2. Run preflight (binding gate; re-run until clean), then commit
   (preflight fixups go into separate `fix(...)` / `chore(...)` commits
   first):

       ./scripts/preflight.sh
       git add plans/active/20260428-redirect-lint-rule.md
       git commit -m "docs(plans): note redirect rule refactored from Node.js to Go"

3. Push and verify `lint`, `lint-custom-linter`, `test-custom-linter`
   on PR #75:

       git push

4. After merge, move BOTH plans to `plans/completed/`; fill Outcomes &
   Retrospective.

## Validation and Acceptance

Accepted when ALL hold:

1. From `tools/lint/`, all three exit 0 and `redirects` shows >= 100.0
   coverage:

       go test ./...
       ./scripts/covgate.sh
       ./scripts/lint.sh

2. From the repo root, both exit 0 — `pnpm run lint` shows no
   `== Redirects ==` block; `pnpm run test:lint` passes all 11
   `run_expect_fail` assertions for `bad-redirects` unmodified
   (binding contract):

       pnpm run lint
       pnpm run test:lint

3. Injecting `{"source": "/docs/getting-started/intro", "destination":
   "/docs/nope"}` into a copy of `docs.json` and running `pnpm run
   lint` produces:
   - `redirects[<n>] source "/docs/getting-started/intro": dead redirect (source resolves to a real page)`
   - `redirects[<n>] destination "/docs/nope": missing destination (no .mdx or .md page exists)`

   Revert before committing.

4. CI on PR #75 passes (`lint`, `lint-custom-linter`,
   `test-custom-linter`).

5. **`./scripts/preflight.sh` reports clean before push.** Mandatory.

## Idempotence and Recovery

- **Read-only.** The rule never mutates `docs.json` or files under
  `docs/`; re-runs are safe and identical. All listed commands are
  idempotent.
- **Per-milestone commits enable rollback.** No data migration.
- **Broken `tools/lint/main.go`.** `git checkout tools/lint/main.go`
  and redo M2; the M1 package commit is independent.
- **Broken rule output at M4.** Fix the rule package (M1 code); do not
  modify the fixture or assertions.
- **Fixture-narrowing tip.** Set
  `DOCS_LINT_ROOT=tests/lint-fixtures/bad-redirects` and invoke the
  linter binary against the inherited `.mdx` files; the redirect rule
  then runs against `tests/lint-fixtures/bad-redirects/docs.json`.

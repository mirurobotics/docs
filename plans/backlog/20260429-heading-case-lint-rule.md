# Heading-case lint rule

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add a new rule (`heading-case`) to the custom Go linter under `tools/lint/`, plus unit and end-to-end tests. |

## Purpose / Big Picture

The docs linter refuses any heading that is not strict sentence-case. A "heading" means both the YAML front-matter `title:` field and every Markdown body heading (`#`..`######`) outside fenced code blocks.

Strict sentence-case: the heading text starts with one uppercase ASCII letter, and every other ASCII letter is lowercase. Proper-noun and acronym handling is out of scope for v1 — the rule will flag `## API reference` and the diagnostic says so explicitly. A future v2 may add an allowlist.

User-visible behavior: `./scripts/lint.sh` (or `pnpm lint`) on a file containing `## Configure Deployments` prints

    docs/x.mdx:5:4: heading-case: heading must be sentence-case (first letter uppercase, all other letters lowercase); proper nouns/acronyms are not yet supported

and exits 1. On `## Configure deployments` it prints nothing and exits 0.

## Progress

- [ ] M1: Skaffold the `headingcase` package and wire it into the rule registry; stubs compile.
- [ ] M2: Implement the front-matter `title:` check + unit tests.
- [ ] M3: Implement the body-heading check + unit tests.
- [ ] M4: Add end-to-end cases to `tools/lint/main_test.go`.
- [ ] M5: Full preflight: `go test ./...` and `./scripts/lint.sh` from repo root both clean.

Use timestamps when you complete steps. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

- Observation:
  Evidence:

## Decision Log

(Add entries as you go.)

- Decision: v1 of `heading-case` is strict — every letter except the first must be lowercase. Proper nouns and acronyms are not detected.
  Rationale: Detection requires a curated allowlist (e.g. "API", "Miru", "macOS") plus context-aware tokenisation. That is a larger design task; punting it to v2 keeps this rule small and reviewable. The diagnostic message says so explicitly so authors are not confused.
  Date/Author: 2026-04-29 / agents@miruml.com

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

The docs repo (`/home/ben/miru/workbench1/repos/docs`) ships a Go linter under `tools/lint/`, built with `go build -o lint .` and invoked per file by `scripts/lint.sh`.

### Layout

- `tools/lint/main.go` — entry point. Calls `linter.ProcessFile`, prints `file:line:col: message` (see `main.go:46`).
- `tools/lint/linter/run.go` — rule registry and dispatch. Defines:
  - The `Rule` string type and one constant per rule (`run.go:17-29`).
  - `AllRules()` listing every rule (`run.go:32-39`).
  - `checkInput{path, lines, spans, contentRoot}` (`run.go:41-46`).
  - `ruleCheckers()` returning a slice of `ruleEntry{Rule, func(checkInput) []analysis.Violation}` (`run.go:55-79`).
  - `ProcessFile(path, contentRoot)` which reads lines, builds `[][]analysis.ProseSpan` via `analysis.Scanner`, and runs every checker.
- `tools/lint/linter/analysis/analysis.go` — `Violation{File, Line, Col(1-based byte), Message}` and `ProseSpan{StartCol(1-based byte), Text(masked)}`.
- `tools/lint/linter/analysis/scanner.go` — `Scanner.ScanLine(line)` returns prose spans, masking inline code, JSX tags, HTML comments, and skipping fenced-code-block and frontmatter lines.
- `tools/lint/linter/analysis/imports.go` — `FrontmatterEnd(lines []string) int` returns the 0-based index of the closing `---` line, or `-1` if no frontmatter (`imports.go:100-110`).
- Per-rule packages: `tools/lint/linter/<rulename>/<rulename>.go` plus `<rulename>_test.go`.

### Existing rule examples to copy from

- `tools/lint/linter/nodoubledash/nodoubledash.go` — prose-span rule. Signature: `func Check(file string, spans [][]analysis.ProseSpan) []analysis.Violation`. 1-based column = `span.StartCol + i`.
- `tools/lint/linter/importsorted/importsorted.go` — line-based rule. Signature: `func Check(file string, lines []string) []analysis.Violation`.
- `tools/lint/linter/importblock/importblock.go` — same shape.

Existing rule constants (the new one is added alongside): `RuleNoDoubleDash`, `RuleImportResolves`, `RuleImportUsed`, `RuleImportSorted`, `RuleComponentStyle`, `RuleMDXStyle`, `RuleImportBlock`, `RuleRedirects`.

### Front-matter shape

A typical `.mdx` starts at line 1:

    ---
    title: "Some title"
    ---

    ## Body heading

The `title:` value may be `"…"`, `'…'`, or unquoted. There is no body H1 — the front-matter title is the page title; body headings start at `##`.

### Test layout

- E2E: `tools/lint/main_test.go` uses `t.TempDir()`, writes `docs/<name>.mdx` (and a sibling `snippets/` directory because `findContentRoot` walks up looking for it), and calls `run([]string{"lint", file}, &stdout, &stderr)`. Asserts on exit code (0 clean, 1 violations, 2 error) and stdout substrings. See `main_test.go:96-113` for the existing clean-case pattern.
- Per-rule unit tests: `tools/lint/linter/<rule>/<rule>_test.go`. The `nodoubledash` tests show both (a) hand-crafted `ProseSpan` slice cases and (b) full-content cases that drive `analysis.Scanner.ScanLine`.
- Bash fixtures under `tests/lint-fixtures/` are being phased out. **DO NOT add anything there.**

### Build / test / run commands

- Build: from `tools/lint/` run `go build -o lint .`
- Tests: from `tools/lint/` run `go test ./...`
- Full pipeline: from repo root `./scripts/lint.sh` (or `pnpm lint`)

## Plan of Work

1. **Create `tools/lint/linter/headingcase/headingcase.go`.**
   Define one exported function:

       func Check(file string, lines []string, spans [][]analysis.ProseSpan) []analysis.Violation

   Behavior:

   - **Front-matter title.** Call `analysis.FrontmatterEnd(lines)`. If the result is `>= 1`, scan lines `1..end-1` (0-based, exclusive of the closing `---`) for the first line whose trimmed form matches `^title:\s*(.*)$`. Strip surrounding `"…"` or `'…'` from the captured value. Apply the casing check (see below). Emit a violation at `Line = lineIdx+1`, `Col = `the 1-based byte column of the first character of the title value (after `title:` and any whitespace, after stripping the opening quote if present).
   - **Body headings.** For every line `i`:
     - **Gate on the scanner.** If `spans[i]` is empty, skip. The scanner emits no spans for lines inside fenced code blocks, so this rules out `# comment` lines inside Bash/etc. code blocks (and also frontmatter lines).
     - **Extract from the raw line.** Match the raw `lines[i]` against `^(#{1,6})[ \t]+(.+?)[ \t]*$`. If it doesn't match, skip. Otherwise capture the heading text from group 2.
     - **Mask the captured text manually** before the casing check:
       - Strip inline code spans: replace `` `[^`]*` `` with the empty string.
       - Strip JSX/HTML tag pairs and self-closing tags: replace `<[^>]*>` with the empty string.
       - Replace Markdown link syntax `\[([^\]]*)\]\([^)]*\)` with `$1` so URLs do NOT enter the casing check.
     - Apply the casing check to the masked string.
     - Emit at `Line = i+1`, `Col = ` the 1-based byte column of the first non-`#`, non-space character on the raw line (i.e. the start of the heading text).

     Read raw lines rather than concatenating prose spans — the scanner preserves the `## ` prefix and may split a heading across spans.
   - **Casing check.** Given a string `s`:
     - Trim leading/trailing whitespace.
     - Trim trailing punctuation in the set `.?!:` (one or more).
     - If `s` is empty after trimming, return no violation (silent skip).
     - Walk through the bytes. The first ASCII letter encountered must be uppercase `A-Z`. Every subsequent ASCII letter must be lowercase `a-z`. Non-letters (spaces, digits, hyphens, apostrophes, punctuation) are ignored.
     - If the rule is violated, return a single violation per heading.
   - **Diagnostic.** Message text:

         heading-case: heading must be sentence-case (first letter uppercase, all other letters lowercase); proper nouns/acronyms are not yet supported

2. **Create `tools/lint/linter/headingcase/headingcase_test.go`.** Table-driven Go tests. See "Test plan" below for the case list.

3. **Edit `tools/lint/linter/run.go`.** Three small additions:
   - Add `"github.com/mirurobotics/docs/tools/lint/linter/headingcase"` to the import block (keep it alphabetical).
   - Add `RuleHeadingCase Rule = "heading-case"` to the `const (...)` block.
   - Add `RuleHeadingCase` to the slice returned by `AllRules()`.
   - Append a new `ruleEntry` to `ruleCheckers()`:

         {RuleHeadingCase, func(in checkInput) []analysis.Violation {
             return headingcase.Check(in.path, in.lines, in.spans)
         }},

   `checkInput` already carries `path`, `lines`, and `spans`, so no struct change is needed.

4. **Edit `tools/lint/main_test.go`.** Add two new sub-tests inside `TestRun` (see "Test plan" below).

### Test plan

#### Unit tests in `tools/lint/linter/headingcase/headingcase_test.go`

Use a table-driven shape. For each row, build the input either as a hand-crafted `[][]analysis.ProseSpan` plus matching `[]string`, OR by feeding a multi-line string through `analysis.NewScanner().ScanLine(...)` line-by-line (the `nodoubledash` tests show both shapes — pick the more readable one per case).

`TestCheck_Headings` cases:

- clean: `## Configure deployments` → 0 violations.
- bad: `## Configure Deployments` → 1 violation; assert `Line=1`, `Col=4`, `Message` starts with `heading-case:`.
- bad: `## API reference` → 1 violation (documents the acronym limitation).
- clean: `## Don't be afraid` → 0 violations (apostrophes are non-letters and are ignored).
- bad: `## Don't Be Afraid` → 1 violation.
- clean: `### What is a config?` → 0 violations (trailing `?` is allowed).
- clean: `## The \`--version\` flag` → 0 violations (the backtick code is masked; "The flag" remains).
- clean: `## [Learn more](/x)` → 0 violations (link text "Learn more" is what gets evaluated).
- bad: `## [Learn More](/x)` → 1 violation.
- heading inside fenced code block (e.g. a four-line input where line 2 is ` ``` `, line 3 is `## not a heading`, line 4 is ` ``` `) → 0 violations.
- empty heading text after masking (e.g. `## <Tooltip />`) → 0 violations and no panic.

`TestCheck_FrontmatterTitle` cases (input is a 3+ line string starting with `---` … `---`):

- quoted clean: `title: "Workspace"` → 0.
- quoted bad: `title: "User Management"` → 1; assert `Line=2`, `Col` points at the `U` of `User`, `Message` starts with `heading-case:`.
- single-quoted clean: `title: 'Deployments'` → 0.
- unquoted clean: `title: Deployments` → 0.
- unquoted bad: `title: API Reference` → 1; assert `Line=2`, `Col=8` (the `A` of `API`: `title: ` is 7 bytes, so `A` is at byte 8; columns are 1-based), `Message` starts with `heading-case:`.
- frontmatter present but no `title:` line → 0 (no panic).
- no frontmatter at all (file starts with `## Foo`) → 0 (no panic from this code path; body heading logic still runs).

For each case verify `Line`, `Col`, and `Message` fields explicitly, not just the count.

#### E2E tests in `tools/lint/main_test.go`

Add two sub-tests to `TestRun`. Each test creates `t.TempDir()`, makes `snippets/` and `docs/` subdirectories, writes a file under `docs/`, calls `run([]string{"lint", file}, &stdout, &stderr)`, and asserts on exit code and stdout.

- `t.Run("heading-case violation returns 1", ...)`. File contents:

      ---
      title: "User Management"
      ---

      ## Configure deployments

  Expect exit `1`. Expect `stdout` to contain `heading-case:` and the file path, and stdout contains `:2:` (the title line). Expect `stderr` empty.

- `t.Run("clean headings return 0", ...)`. File contents:

      ---
      title: "Workspace"
      ---

      ## Configure deployments

      ### What is a config?

  Expect exit `0`. Expect `stdout` empty.

## Concrete Steps

All commands assume the working tree is clean before starting and that you are working on a feature branch (e.g. `feat/heading-case-lint`). The repo root is `/home/ben/miru/workbench1/repos/docs`.

### M1 — Skaffold and register

1. Create the package directory and a stub rule file:

       cd /home/ben/miru/workbench1/repos/docs
       mkdir -p tools/lint/linter/headingcase

   Create `tools/lint/linter/headingcase/headingcase.go` with the package declaration, the `analysis` import, and a stub `func Check(file string, lines []string, spans [][]analysis.ProseSpan) []analysis.Violation { return nil }`.

2. Edit `tools/lint/linter/run.go` per "Plan of Work" item 3 (import, constant, `AllRules()` entry, `ruleCheckers()` entry).

3. Create placeholder `tools/lint/linter/headingcase/headingcase_test.go` with `package headingcase`, `import "testing"`, and a no-op `func TestPlaceholder(t *testing.T) {}` so the package builds.

4. Verify build and tests:

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go build -o lint .
       go test ./...

   Expected: build succeeds, all existing tests pass, the placeholder test passes.

5. Commit:

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/linter/headingcase tools/lint/linter/run.go
       git commit -m "feat(lint): scaffold heading-case rule package and registry wiring"

### M2 — Front-matter title check

1. In `tools/lint/linter/headingcase/headingcase.go`, implement the front-matter title path from "Plan of Work" item 1, sub-bullet "Front-matter title". Add the shared `casingViolation(s string) bool` helper here — body headings reuse it in M3.

2. In `headingcase_test.go`, replace the placeholder with the `TestCheck_FrontmatterTitle` table from "Test plan".

3. Run tests:

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test ./linter/headingcase/...
       go test ./...

   Expected: every `TestCheck_FrontmatterTitle` case passes; nothing else regresses.

4. Commit:

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/linter/headingcase
       git commit -m "feat(lint): heading-case checks frontmatter title"

### M3 — Body heading check

1. Extend `Check` with the body-heading path from "Plan of Work" item 1, sub-bullet "Body headings". Reuse `casingViolation`.

2. Add the `TestCheck_Headings` table from "Test plan" to `headingcase_test.go`.

3. Run tests:

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test ./linter/headingcase/...
       go test ./...

   Expected: all `TestCheck_Headings` cases pass; nothing else regresses.

4. Commit:

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/linter/headingcase
       git commit -m "feat(lint): heading-case checks body headings"

### M4 — End-to-end tests

1. Audit existing sub-tests in `tools/lint/main_test.go` for clean-exit fixtures with headings. The current `t.Run("clean run returns 0", ...)` block writes `# x\n` — after `heading-case` is registered this becomes a violation (lowercase first letter). Update that fixture to either (a) contain no heading, or (b) use a sentence-case heading (e.g., `"# Hello\n"`). Sub-tests that only assert non-zero exit (`nonexistent file returns 2`, `missing snippets returns 2`) are unaffected — leave them alone.

2. Add the two sub-tests from "Test plan" → "E2E tests in `tools/lint/main_test.go`".

3. Run E2E tests:

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test -run TestRun ./...

   Expected: both new sub-tests pass; existing `TestRun` sub-tests still pass.

4. Commit:

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/main_test.go
       git commit -m "test(lint): end-to-end coverage for heading-case rule"

### M5 — Full preflight

1. Run the full Go test suite:

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test ./...

   Expected: `ok` for every package.

2. Run the repo-wide lint pipeline:

       cd /home/ben/miru/workbench1/repos/docs
       ./scripts/lint.sh

   The script rebuilds the linter and runs it against the docs corpus. It may now report `heading-case:` violations for pre-existing offending headings in real docs — that is expected (Decision Log) and out of scope. If exit is non-zero ONLY because of pre-existing `heading-case:` violations, capture them in Surprises & Discoveries and proceed. Any other failure must be debugged before merging.

   To confirm no other rule is firing:

       ./scripts/lint.sh 2>&1 | grep -E ':[0-9]+:[0-9]+:' | grep -v 'heading-case:' | head -20

   Expected: empty (or only pre-existing non-`heading-case` violations).

3. No commit needed unless preflight surfaces an issue you fix; if so commit as `fix(lint): <what>`.

## Validation and Acceptance

The change is acceptable when all of these hold:

1. From `tools/lint/`, `go test ./...` exits 0. The new `TestCheck_Headings`, `TestCheck_FrontmatterTitle`, and the two new `TestRun` sub-tests all pass — and fail before the implementation lands.

2. From the repo root, `./scripts/lint.sh` builds cleanly. Exit-1 caused by `heading-case` violations in pre-existing real docs is expected (Decision Log) and logged in Surprises & Discoveries; any other exit code, or violations from rules other than `heading-case`, must be investigated.

3. The linter against a hand-written file with `## Configure Deployments` matches the example output in Purpose (`<path>:5:4: heading-case: …`) and exits 1; replacing `Deployments` with `deployments` makes it exit 0 with no output.

4. A file with front-matter `title: "User Management"` produces a `heading-case:` diagnostic on line 2 and exits 1.

## Idempotence and Recovery

- All steps are safe to repeat; `go test ./...` and `./scripts/lint.sh` are idempotent.
- M1-M4 edits are additive. Re-running a milestone's instructions overwrites the rule file and re-applies registry edits cleanly. Roll back via `git restore` on listed paths or `git reset --hard` to the previous milestone's commit.
- If `go test ./...` regresses outside the `headingcase` package after wiring the rule, audit other tests for non-sentence-case heading fixtures (the existing `# x` fixture is the known instance, addressed in M4).
- If a pre-commit hook fails, do NOT amend; fix the underlying issue and create a new commit on top.

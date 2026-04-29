# Heading-case lint rule

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add a new rule (`heading-case`) to the custom Go linter under `tools/lint/`, plus unit and end-to-end tests. |

This plan lives in `docs/plans/backlog/` because all changes are confined to the docs repository's linter.

## Purpose / Big Picture

After this change, the docs linter will refuse any heading that is not strict sentence-case. A "heading" means both:

- the YAML front-matter `title:` field at the top of an `.mdx`/`.md` file, and
- every Markdown body heading (`#`, `##`, `###`, `####`, `#####`, `######`) outside fenced code blocks.

Strict sentence-case means: the heading text starts with one uppercase ASCII letter, and every other ASCII letter in the heading text is lowercase. Proper-noun and acronym handling is explicitly out of scope for v1 — the rule will flag `## API reference` and the diagnostic will tell the author so.

The user-visible behavior is: running `./scripts/lint.sh` (or `pnpm lint`) on a file containing `## Configure Deployments` prints

    docs/x.mdx:5:4: heading-case: heading must be sentence-case (first letter uppercase, all other letters lowercase); proper nouns/acronyms are not yet supported

and exits with code 1. Running it on `## Configure deployments` prints nothing and exits 0.

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

The docs repo (`/home/ben/miru/workbench1/repos/docs`) ships its own Go linter under `tools/lint/`. The linter is a single binary built with `go build -o lint .` from `tools/lint/`. It is called per file by `scripts/lint.sh`.

### Layout

- `tools/lint/main.go` — entry point. Iterates over file arguments, calls `linter.ProcessFile`, prints each violation in the format `file:line:col: message` (see `main.go:46`).
- `tools/lint/linter/run.go` — rule registry and dispatch. Defines:
  - The `Rule` string type and one constant per rule (block at `run.go:17-29`).
  - `AllRules()` listing every rule (at `run.go:32-39`).
  - `checkInput` struct with fields `path`, `lines`, `spans`, `contentRoot` (at `run.go:41-46`).
  - `ruleCheckers()` returning a slice of `ruleEntry` — each entry is `{Rule, func(checkInput) []analysis.Violation}` (at `run.go:55-79`).
  - `ProcessFile(path, contentRoot)` which reads lines, builds `[][]analysis.ProseSpan` via `analysis.Scanner`, and runs every checker.
- `tools/lint/linter/analysis/analysis.go` — shared types:
  - `Violation{File, Line, Col(1-based byte), Message}`.
  - `ProseSpan{StartCol(1-based byte), Text(masked)}`.
- `tools/lint/linter/analysis/scanner.go` — `Scanner.ScanLine(line)` returns prose spans for one line, masking inline code, JSX tags, HTML comments, and skipping fenced-code-block and frontmatter lines.
- `tools/lint/linter/analysis/imports.go` — `FrontmatterEnd(lines []string) int` returns the 0-based index of the closing `---` line, or `-1` if no frontmatter is present (at `imports.go:100-110`).
- Per-rule packages live under `tools/lint/linter/<rulename>/<rulename>.go` plus `<rulename>_test.go`.

### Existing rule examples to copy from

- `tools/lint/linter/nodoubledash/nodoubledash.go` — prose-span rule. Signature: `func Check(file string, spans [][]analysis.ProseSpan) []analysis.Violation`. Computes 1-based column as `span.StartCol + i`.
- `tools/lint/linter/importsorted/importsorted.go` — line-based rule. Signature: `func Check(file string, lines []string) []analysis.Violation`.
- `tools/lint/linter/importblock/importblock.go` — same shape.

### Existing rule constants (the new one is added alongside)

- `RuleNoDoubleDash`, `RuleImportResolves`, `RuleImportUsed`, `RuleImportSorted`, `RuleComponentStyle`, `RuleMDXStyle`, `RuleImportBlock`, `RuleRedirects`.

### Front-matter shape (from real docs files)

A typical `.mdx` starts at line 1:

    ---
    title: "Some title"
    ---

    ## Body heading

The `title:` value can be wrapped in `"…"`, `'…'`, or unquoted. There is no body H1 — the front-matter title is the page title; body headings start at `##`.

### Heading edge cases the rule must handle

- MDX components in headings, e.g. `## <Tooltip>foo</Tooltip>` — the scanner masks JSX so the prose span yields clean text.
- Inline code, e.g. `## The \`--version\` flag` — backticks masked by the scanner; only the surrounding prose is checked.
- Markdown links, e.g. `## [Learn more](/x)` — the scanner does NOT mask link syntax. The rule itself must extract just the link text; the simplest tactic is to regex-strip `\[([^\]]*)\]\([^)]*\)` to `$1` before the casing check.
- `#`-prefixed lines inside fenced code blocks (e.g. Bash comments) must be ignored. Using prose spans as the source of truth handles this for free: the scanner produces no spans for code-fence-internal lines.

### Test layout

- E2E: `tools/lint/main_test.go` uses `t.TempDir()`, writes `docs/<name>.mdx` (and a sibling `snippets/` directory because `findContentRoot` walks up looking for it), and calls `run([]string{"lint", file}, &stdout, &stderr)`. It asserts on exit code (0 clean, 1 violations, 2 error) and stdout substrings. See `main_test.go:96-113` for the existing clean-case sub-test pattern.
- Per-rule unit tests: `tools/lint/linter/<rule>/<rule>_test.go`. The `nodoubledash` tests have BOTH (a) hand-crafted `ProseSpan` slice cases and (b) full-content cases that drive `analysis.Scanner.ScanLine` to produce real spans.
- Bash fixtures under `tests/lint-fixtures/` driven by `tests/test-lint.sh` are being phased out. **DO NOT add anything there.** All new lint logic gets Go unit tests + Go E2E.

### Build / test / run commands

- Build linter: from `tools/lint/` run `go build -o lint .`
- Go tests: from `tools/lint/` run `go test ./...`
- Full lint pipeline: from repo root `./scripts/lint.sh` (also via `pnpm lint`)

### Limitations (v1 — important)

The rule does not understand proper nouns, acronyms, brand names, code identifiers, or language tokens. Every letter except the first letter of the heading text must be lowercase ASCII. That means the rule will (intentionally) flag headings that are otherwise correct English, including:

- `## API reference` (acronym)
- `## Using macOS` (brand)
- `## Configure Miru workspaces` (proper noun)

The diagnostic message states this so authors know the rule is the limitation, not their wording. A future v2 may add an allowlist; the current rule is deliberately strict and small.

## Plan of Work

The work is six small edits across two new files and one existing file, plus tests.

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

     Why read raw lines instead of concatenating prose spans: the scanner preserves the `## ` prefix in span text and may split a heading into multiple non-contiguous spans when JSX or inline code is present. Deriving the heading text from the raw line and masking manually is simpler — it avoids a prefix-strip step and avoids reasoning about inter-span gaps.
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

5. Run `go test ./...` from `tools/lint/` and `./scripts/lint.sh` from the repo root.

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

1. Create the directory and the rule file with a stub `Check` that returns `nil`.

       cd /home/ben/miru/workbench1/repos/docs
       mkdir -p tools/lint/linter/headingcase

   Then create `tools/lint/linter/headingcase/headingcase.go` containing only the package declaration, the `analysis` import, and a stub `func Check(file string, lines []string, spans [][]analysis.ProseSpan) []analysis.Violation { return nil }`.

2. Edit `tools/lint/linter/run.go` per "Plan of Work" item 3: add the import, the constant, the `AllRules()` entry, and the `ruleCheckers()` entry.

3. Create an empty placeholder test file `tools/lint/linter/headingcase/headingcase_test.go` containing only `package headingcase` and a single `import "testing"` plus a no-op `func TestPlaceholder(t *testing.T) {}` so the package builds.

4. Verify everything compiles and existing tests still pass.

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go build -o lint .
       go test ./...

   Expected: build succeeds, all existing tests pass, the new placeholder test passes.

5. Commit the milestone.

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/linter/headingcase tools/lint/linter/run.go
       git commit -m "feat(lint): scaffold heading-case rule package and registry wiring"

### M2 — Front-matter title check

1. In `tools/lint/linter/headingcase/headingcase.go`, implement the front-matter title path described in "Plan of Work" item 1, sub-bullet "Front-matter title". Implement the shared `casingViolation(s string) bool` helper here too — body headings will reuse it in M3.

2. In `tools/lint/linter/headingcase/headingcase_test.go`, replace the placeholder with the `TestCheck_FrontmatterTitle` table from "Test plan".

3. Run unit tests.

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test ./linter/headingcase/...

   Expected: every case in `TestCheck_FrontmatterTitle` passes. Then run the full suite to confirm nothing else broke:

       go test ./...

4. Commit.

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/linter/headingcase
       git commit -m "feat(lint): heading-case checks frontmatter title"

### M3 — Body heading check

1. Extend `Check` in `tools/lint/linter/headingcase/headingcase.go` with the body-heading path described in "Plan of Work" item 1, sub-bullet "Body headings". Reuse `casingViolation`.

2. Add the `TestCheck_Headings` table from "Test plan" to `headingcase_test.go`.

3. Run unit tests.

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test ./linter/headingcase/...
       go test ./...

   Expected: all `TestCheck_Headings` cases pass; nothing else regresses.

4. Commit.

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/linter/headingcase
       git commit -m "feat(lint): heading-case checks body headings"

### M4 — End-to-end tests

1. Audit the existing sub-tests in `tools/lint/main_test.go` for heading-shaped fixtures whose stdout is asserted clean (exit 0). The current `t.Run("clean run returns 0", ...)` block writes `# x\n` — after `heading-case` is registered, `# x` becomes a violation (lowercase first letter). Update that fixture to a string that either (a) contains no heading, or (b) uses a sentence-case heading (e.g., `"# Hello\n"`). Sub-tests that only assert non-zero exit (`nonexistent file returns 2`, `missing snippets returns 2`) are unaffected because they don't assert clean stdout — leave them alone.

2. Edit `tools/lint/main_test.go` to add the two sub-tests described in "Test plan" → "E2E tests in `tools/lint/main_test.go`".

3. Run E2E tests.

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test -run TestRun ./...

   Expected: both new sub-tests pass; existing `TestRun` sub-tests still pass.

4. Commit.

       cd /home/ben/miru/workbench1/repos/docs
       git add tools/lint/main_test.go
       git commit -m "test(lint): end-to-end coverage for heading-case rule"

### M5 — Full preflight

1. Run the full Go test suite from the linter root.

       cd /home/ben/miru/workbench1/repos/docs/tools/lint
       go test ./...

   Expected: `ok` for every package.

2. Run the repo-wide lint pipeline.

       cd /home/ben/miru/workbench1/repos/docs
       ./scripts/lint.sh

   This script rebuilds the linter and runs it against the docs corpus. **Expected behavior:** the script may now report `heading-case:` violations for pre-existing offending headings in the real docs. That is expected and explicitly out of scope for this plan (see Decision Log). If `./scripts/lint.sh` exits non-zero ONLY because of pre-existing `heading-case:` violations in real docs, that does not block the milestone — capture them in a brief note in Surprises & Discoveries and proceed. If it exits non-zero for any other reason, debug that separately before merging.

   To confirm at a glance that no rule other than `heading-case` is firing, also run:

       ./scripts/lint.sh 2>&1 | grep -E ':[0-9]+:[0-9]+:' | grep -v 'heading-case:' | head -20

   Expected output: empty (or only non-`heading-case` violations that pre-existed). If non-empty and not `heading-case`, investigate.

3. No commit needed for M5 unless preflight surfaces an issue you fix. If you do fix something, commit it as `fix(lint): <what you fixed>`.

## Validation and Acceptance

The change is acceptable when all of these hold:

1. From `tools/lint/`, `go test ./...` exits 0. The new `TestCheck_Headings`, `TestCheck_FrontmatterTitle`, and the two new `TestRun` sub-tests all pass. Concretely the new tests fail before the implementation lands and pass after.

2. From the repo root, `./scripts/lint.sh` builds the linter cleanly and runs the new rule across the corpus. Any exit-1 caused by `heading-case` violations in pre-existing real docs is expected (Decision Log) and is logged in Surprises & Discoveries; an exit code other than 0 or 1, or violations from rules other than `heading-case`, must be investigated.

3. Running the linter against a hand-written file with `## Configure Deployments` produces:

       <path>:5:4: heading-case: heading must be sentence-case (first letter uppercase, all other letters lowercase); proper nouns/acronyms are not yet supported

   and exits 1. Replacing `Deployments` with `deployments` makes the linter exit 0 with no output.

4. Running the linter against a file whose front-matter is `title: "User Management"` produces a `heading-case:` diagnostic on line 2 and exits 1.

## Idempotence and Recovery

- All steps are safe to repeat. `go test ./...` and `./scripts/lint.sh` are idempotent.
- The Go file edits in M1-M4 are additive; if a milestone is half-applied, re-running the M's instructions overwrites the rule file and re-applies the registry edits cleanly. If you need to roll back a milestone, use `git restore` on the listed paths or `git reset --hard` to the previous milestone's commit.
- Adding the rule may surface pre-existing `heading-case` violations in real docs. **Fixing those is out of scope for this plan** (see Decision Log). A follow-up cleanup PR will sweep them. If `./scripts/lint.sh` exits 1 only because of pre-existing offenders, that does not block the milestone; record the count in Surprises & Discoveries and move on.
- If `go test ./...` regresses outside the `headingcase` package after wiring the rule, audit other tests for fixtures whose headings are not sentence-case (the existing `# x` fixture is one known instance, addressed in M4).
- If a pre-commit hook fails, do NOT amend; fix the underlying issue and create a new commit on top, per repo policy.

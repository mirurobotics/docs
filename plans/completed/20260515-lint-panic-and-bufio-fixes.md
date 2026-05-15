# Fix lint panic in component-style rule and bufio 64KB limit

## Context

Two defects in the Go documentation linter (`tools/lint`, module
`github.com/mirurobotics/docs/tools/lint`, Go 1.25.3) can crash or silently
disable linting for an entire invocation. Both abort `main.run` with a
non-zero exit (panic / exit 2) instead of reporting a per-file violation,
so a single malformed or large doc file blocks CI for the whole repo.

This plan fixes both bugs with minimal, surgical edits and adds focused
unit tests. Base branch: `main`.

## Bug 1 — panic in component-style rule

File: `tools/lint/linter/componentstyle/componentstyle.go`
Function: `checkBracketSpacing` (lines 46-77).

`checkStyle` (line 27) computes `openBrace := strings.Index(raw, "{")`
(FIRST `{`) and passes it to `checkBracketSpacing`. There,
`closeBrace := strings.LastIndex(raw, "}")` (LAST `}`, line 52). The guard
on line 53 is `if closeBrace < 0`. For a malformed import such as:

    import Foo} from '/snippets/components/x.jsx' // {Bar

`openBrace` points at the `{` inside the trailing comment and `closeBrace`
points at the earlier `}`, so `openBrace > closeBrace`. Line 57,
`body := raw[openBrace+1 : closeBrace]`, then panics with
`slice bounds out of range`, crashing the entire lint binary.

### Fix

In `tools/lint/linter/componentstyle/componentstyle.go`, change the guard
on line 53 from:

    if closeBrace < 0 {

to:

    if closeBrace <= openBrace {

This treats a missing OR mis-ordered closing brace identically to the
existing "must use named import syntax { }" case. The violation message
string on line 54 must remain EXACTLY:

    "import-component-style: component import must use named import syntax { }"

(unchanged; it is asserted by `tests/test-lint.sh` fixtures and existing
unit tests). No other lines in this function change. The existing
`TestCheckBracketSpacingNoClosure` test (raw with `{` but no `}`,
`closeBrace == -1`, `openBrace >= 0`) still passes because `-1 <= openBrace`
holds, producing the same single named-import-syntax violation.

## Bug 2 — bufio.Scanner 64KB token limit

File: `tools/lint/linter/run.go`
Function: `ProcessFile` (lines 96-125).

Lines 104-110 create a default `bufio.NewScanner(f)` with no buffer
override. `bufio.Scanner`'s default `MaxScanTokenSize` is 64KB. A docs MDX
file with one very long line (embedded base64 image / data URI, a long
HTML table row, or inline SVG) makes `scanner.Scan()` stop early and
`scanner.Err()` return `bufio.ErrTooLong`. `ProcessFile` then returns that
error, the file is never linted, and `main.run` exits 2 for the whole
invocation.

### Fix

In `tools/lint/linter/run.go`:

1. Add a package-level const near the top of the file (after the import
   block, before `type Rule string` on line 19), with a short rationale
   comment:

       // maxScanTokenSize caps bufio.Scanner's per-line buffer. Docs MDX
       // files can contain a single very long line (embedded base64 data
       // URIs, inline SVG, wide HTML table rows) that exceeds bufio's
       // default 64KB token limit; 16 MB comfortably covers realistic
       // documentation lines without unbounded memory growth.
       const maxScanTokenSize = 16 * 1024 * 1024

2. In `ProcessFile`, immediately after line 104
   (`scanner := bufio.NewScanner(f)`) and before the scan loop, add:

       scanner.Buffer(make([]byte, 0, 64*1024), maxScanTokenSize)

   The initial buffer keeps the small-line common case cheap (starts at
   64KB, grows as needed up to `maxScanTokenSize`). The existing
   `scanner.Err()` check on lines 108-110 is unchanged; lines longer than
   16 MB still surface an error rather than truncating silently.

No other logic in `ProcessFile` changes.

## Tests

### Test 1 — componentstyle malformed import does not panic

File: `tools/lint/linter/componentstyle/componentstyle_test.go`.

The file uses `t.Run` subtests inside `TestCheck` (each subtest builds a
`line`, calls `Check("test.mdx", []string{line})`, and scans returned
violations for a substring). Add a new `t.Run` subtest inside `TestCheck`
following the same shape:

    t.Run("malformed brace order does not panic", func(t *testing.T) {
        line := "import Foo} from '/snippets/components/x.jsx' // {Bar"
        vs := Check("test.mdx", []string{line})
        found := false
        for _, v := range vs {
            if strings.Contains(v.Message, "named import syntax") {
                found = true
            }
        }
        if !found {
            t.Errorf("expected named-import-syntax violation, got %v", vs)
        }
    })

This input has the first `{` (in the trailing comment) AFTER the `}`, so
`openBrace > closeBrace` — the exact panic case. Match the file's actual
test structure; adjust shape if it differs.

### Test 2 — ProcessFile handles lines longer than 64KB

File: `tools/lint/linter/run_test.go`.

Existing pattern: create `dir := t.TempDir()`, write a file with
`os.WriteFile`, then call `ProcessFile(path, dir)` passing the temp dir as
`contentRoot`. Use plain prose content so no rule fires on the long line.

Add a new `t.Run` subtest inside `TestProcessFile`:

    t.Run("very long line does not error", func(t *testing.T) {
        dir := t.TempDir()
        path := filepath.Join(dir, "test.mdx")
        longLine := strings.Repeat("x", 200*1024)
        content := "# Title\n\n" + longLine + "\n"
        if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
            t.Fatal(err)
        }
        if _, err := ProcessFile(path, dir); err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
    })

Add `"strings"` to `run_test.go`'s import block if absent. Before the Bug
2 fix this subtest fails with `bufio.ErrTooLong`; after it passes. Match
the file's actual harness conventions.

## Validation

Run from the repo root unless noted. All must pass before publishing.

1. Go unit tests: `cd tools/lint && go test ./...`
2. Go test wrapper: `./tools/lint/scripts/test.sh`
3. Lint fixture suite: `tests/test-lint.sh`
4. Coverage ratchet — `.covgate` files exist in touched packages
   (`tools/lint/linter/componentstyle/.covgate`,
   `tools/lint/linter/.covgate`, `tools/lint/.covgate`). Coverage MUST
   NOT regress; verify with `./tools/lint/scripts/covgate.sh`. Add test
   cases rather than lowering any `.covgate` value.
5. Go lint: `LINT_FIX=0 ./tools/lint/scripts/lint.sh`
6. Full preflight: `./scripts/preflight.sh`

Preflight MUST report `clean` before these changes are published. Do not
open a PR or push until `./scripts/preflight.sh` exits successfully.

## Implementation order

1. Apply Bug 1 fix (`componentstyle.go` guard).
2. Apply Bug 2 fix (`run.go`: const + `scanner.Buffer(...)`).
3. Add Test 1 (`componentstyle_test.go`).
4. Add Test 2 (`run_test.go` + `strings` import).
5. Run Validation steps 1-6; resolve any covgate or lint findings.
6. Confirm preflight `clean`.

## Risk / rollback

Both fixes are localized (one guard condition; one const + one method
call). No public API or message-string changes. Rollback is reverting the
two source edits and the two test additions.

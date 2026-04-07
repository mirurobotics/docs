# Lint the docs MDX linter with gotools

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add a Go lint script for `tools/lint/`, wire it into `scripts/preflight.sh` and `.github/workflows/lint.yml`, fix all violations in `tools/lint/*.go`, and bump `tools/lint/go.mod`. |
| `gotools/` | read-only | Provides `github.com/mirurobotics/gotools/cmd/miru` (invoked as `go tool miru lint`). No edits. |
| `cli-private/` | read-only | Reference for the `tool` directive + `go tool miru lint` invocation pattern (see `cli-private/go.mod` and `cli-private/scripts/lint.sh`). |
| `core/` | read-only | Same reference pattern (see `core/go.mod` and `core/scripts/lint.sh`). |
| `agent/` | read-only | Reference for the "tools/lint has its own scripts/lint.sh" structure (see `agent/tools/lint/scripts/lint.sh`, which is a Rust binary but shows the directory layout). |

This plan lives in `docs/plans/backlog/` because all code changes are made inside the docs repo — specifically inside `docs/tools/lint/`, `docs/scripts/`, and `docs/.github/workflows/`. There are no changes to any other repo.

## Purpose / Big Picture

The docs repo ships a custom Go-based MDX prose linter at `docs/tools/lint/` that scans `docs/**/*.mdx` and `docs/snippets/**/*.mdx` for issues like double-dashes that should be em-dashes. Today the MDX content is linted, but the Go source of the linter itself is not. This plan lints the linter: it wires the public `gotools` linter (`github.com/mirurobotics/gotools/cmd/miru lint`) into `docs` so that `tools/lint/*.go` is continuously checked for style, formatting, line-length, and common Go mistakes — exactly the same gate that `cli-private`, `core`, and other Go repos already enforce on their code.

After this change, a contributor editing any file under `docs/tools/lint/` will:

1. Have their changes automatically reformatted by `LINT_FIX=1 ./tools/lint/scripts/lint.sh` (the default when run locally).
2. Be blocked from landing a PR if any `custom-linter`, `gofumpt`, or `golangci-lint` violation remains, via both `docs/scripts/preflight.sh` (local preflight) and `docs/.github/workflows/lint.yml` (CI).
3. Observe the new output section in the preflight transcript that reads roughly:

        == Go Lint (tools/lint) ==
        Running custom linter on /home/.../docs/tools/lint...
        Running gofumpt...
        Running golangci-lint...
        0 issues.

4. See zero violations the very first time they run it, because this plan also fixes all 15 currently-reported violations as part of the same change set.

## Progress

- [ ] Milestone 1: Add `tool` directive + helper script. Bump `docs/tools/lint/go.mod` to `go 1.25.3`, add the `tool ( github.com/mirurobotics/gotools/cmd/miru )` directive, and create `docs/tools/lint/scripts/lint.sh`.
- [ ] Milestone 2: Run `LINT_FIX=1 ./tools/lint/scripts/lint.sh` to auto-fix the collapsible-expression, inlinable-function, and incidental formatting violations.
- [ ] Milestone 3: Manually fix the 7 remaining `line is N columns wide` violations in `main.go`, `rules_test.go`, and `scanner.go`. Confirm `go test ./...` inside `docs/tools/lint/` still passes.
- [ ] Milestone 4: Wire `LINT_FIX=0 ./tools/lint/scripts/lint.sh` into `docs/scripts/preflight.sh` (as a new `== Go Lint (tools/lint) ==` step before the existing `./scripts/lint.sh` call).
- [ ] Milestone 5: Wire `LINT_FIX=0 ./tools/lint/scripts/lint.sh` into `docs/.github/workflows/lint.yml` (as a new step in the `lint` job, before the existing "Run documentation lint" step).
- [ ] Milestone 6: Run full `./scripts/preflight.sh` from `docs/` and confirm all sections report clean, then commit the final state.

Use timestamps when you complete steps. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

- Observation: …
  Evidence: …

## Decision Log

- Decision: Use a `tool` directive in `docs/tools/lint/go.mod` rather than `go run github.com/mirurobotics/gotools/cmd/miru@<version>`.
  Rationale: This matches the pattern already in use in `cli-private/go.mod` lines 262-267 and `core/go.mod`. It gives a clean `go tool miru lint` invocation with pinned versions resolved through `go.sum`, avoids network access on every CI run, and keeps the gotools version upgrade story uniform across the monorepo. The downside (adding ~250 indirect dependencies to `docs/tools/lint/go.sum`) is the same trade-off already accepted by sibling Go repos, and the docs linter module is isolated from the Mintlify content, so bloat there does not affect docs build times.
  Date/Author: 2026-04-06 / author

- Decision: Bump `docs/tools/lint/go.mod` from `go 1.24` to `go 1.25.3`.
  Rationale: The `gotools` module declares `go 1.25.3` in its own `go.mod`. A consumer cannot take a `tool` dependency on it with a lower go directive — Go will refuse to build. CI resolves its Go toolchain from `go-version-file: tools/lint/go.mod` (see `docs/.github/workflows/lint.yml` line 24), so bumping the module's go directive also bumps the CI toolchain version, which is required for `go tool miru lint` to work. This matches `cli-private/go.mod` and `core/go.mod`, both of which are on `go 1.25.3`.
  Date/Author: 2026-04-06 / author

- Decision: Create `docs/tools/lint/scripts/lint.sh` as a per-tool entry script (mirroring `agent/tools/lint/scripts/lint.sh`) rather than adding the `go tool miru lint` call directly inside `docs/scripts/lint.sh`.
  Rationale: `go tool miru lint` only resolves inside the Go module directory — the script must `cd` into `docs/tools/lint/` before running. Wrapping that cwd switch in a dedicated per-tool script keeps `docs/scripts/lint.sh` (which handles MDX prose + CSpell + ESLint + OpenAPI) focused on documentation content, and makes the Go-lint step independently runnable: `LINT_FIX=0 ./tools/lint/scripts/lint.sh` works from the docs repo root without any wrapper. The agent repo already uses exactly this layout for its Rust linter.
  Date/Author: 2026-04-06 / author

- Decision: Honor the `LINT_FIX` environment variable with the convention `LINT_FIX=1` (default) → `--fix`, `LINT_FIX=0` → `--fix=false`.
  Rationale: This matches `gotools/scripts/lint.sh` lines 4-7, `cli-private/scripts/lint.sh`, `core/scripts/lint.sh`, and `agent/tools/lint/scripts/lint.sh`. Local runs (including the default `./tools/lint/scripts/lint.sh`) auto-fix trivially fixable issues so contributors do not hand-format; preflight and CI pass `LINT_FIX=0` to enforce check-only mode and fail on any unfixed violation.
  Date/Author: 2026-04-06 / author

- Decision: Do not pass `--exclude` to `miru lint` for the docs linter.
  Rationale: `gotools/scripts/lint.sh` excludes `nofmt,paramcount` for its own internal build reasons, but `docs/tools/lint/` is a tiny 5-file module with simple functions — a full-strength lint pass is both feasible and desirable. If a future rule produces unhelpful noise, a follow-up plan can narrow the scope with `--exclude`.
  Date/Author: 2026-04-06 / author

- Decision: Do not pass `--deadcode` for the docs linter.
  Rationale: `cli-private/scripts/lint.sh` runs `--deadcode` because it has a large `internal,cmd,mock,tests` surface area where dead code accumulates. The docs linter has ~850 lines across 5 files and every exported symbol is reachable from `main.go`. Running deadcode would add build time with no signal. A future plan can add `--deadcode` if the module grows.
  Date/Author: 2026-04-06 / author

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

Read this section first if you have never touched the docs repo before.

### Repo layout relevant to this plan

- `/home/ben/miru/workbench3/docs/` — the docs repo root. Mintlify documentation site with MDX content in `docs/` and `snippets/`, plus a custom Go-based prose linter under `tools/lint/`.
- `/home/ben/miru/workbench3/docs/tools/lint/` — a standalone Go module. `go.mod` declares `module github.com/mirurobotics/docs/tools/lint` and `go 1.24`. Contains:
  - `main.go` (64 lines) — entry point. Reads file paths from `os.Args`, builds a list of `Rule` values, calls `lintFile` on each path.
  - `rules.go` (44 lines) — defines the `Violation` struct, the `Rule` interface, and the `NoDoubleDash` rule that flags `--` in prose.
  - `rules_test.go` (203 lines) — table-driven tests for `NoDoubleDash`, including a full integration test that feeds MDX through the `Scanner`.
  - `scanner.go` (279 lines) — stateful line-by-line scanner that classifies regions of an MDX file into prose / frontmatter / code block / HTML comment and masks out inline code and JSX.
  - `scanner_test.go` (263 lines) — tests for `Scanner.ScanLine`.
  - `lint` — the compiled binary (gitignored; produced by `go build -o lint .` inside the module directory).
- `/home/ben/miru/workbench3/docs/scripts/` — shell scripts invoked during local development.
  - `lint.sh` — the documentation content linter. Builds `tools/lint/lint` and runs it against `*.mdx` files, then runs ESLint, CSpell, and Mintlify's OpenAPI check.
  - `preflight.sh` — the local PR gate. Currently runs `pnpm run test:lint`, `./scripts/lint.sh`, `./scripts/audit.sh`. This plan adds a Go-lint step.
  - `audit.sh` — runs `pnpm audit`. Untouched.
- `/home/ben/miru/workbench3/docs/.github/workflows/lint.yml` — CI workflow with two jobs: `lint` (installs Go + Node, runs `pnpm run test:lint` and `./scripts/lint.sh`) and `audit` (runs `./scripts/audit.sh`). This plan adds a Go-lint step to the `lint` job.

### The gotools linter

- Module path: `github.com/mirurobotics/gotools`. Source on disk: `/home/ben/miru/workbench3/gotools/`. Public.
- Entrypoint: `./cmd/miru`. The `lint` subcommand runs three linters in sequence: a custom linter defined in `gotools/internal/customlinter/`, then `gofumpt`, then `golangci-lint`. All three must pass for the subcommand to exit 0.
- Key flags:
  - `--paths=<dirs>` — comma-separated directories the custom linter scans. Required.
  - `--fix` / `--fix=false` — auto-fix vs check-only. Default `--fix`. CI and preflight use `--fix=false`.
  - `--exclude=<rules>` — comma-separated rules to skip. Not used here.
  - `--max-line-width=88` — default soft limit for the line-width rule.
  - `--max-func-len=50` — default soft limit for the function-length rule.
- Invocation pattern (from `gotools/scripts/lint.sh`):

        FIX="--fix"
        if [ "${LINT_FIX:-1}" = "0" ]; then
            FIX="--fix=false"
        fi
        exec go run ./cmd/miru lint --paths=internal --exclude=nofmt,paramcount $FIX

- When imported via a `tool` directive in another module's `go.mod`, the invocation becomes `go tool miru lint --paths=... $FIX` (see `cli-private/scripts/lint.sh`).

### What a "tool" directive looks like

From `/home/ben/miru/workbench3/cli-private/go.mod` lines 262-267:

        tool (
            github.com/golangci/golangci-lint/v2/cmd/golangci-lint
            github.com/mirurobotics/gotools/cmd/miru
            golang.org/x/tools/cmd/deadcode
            mvdan.cc/gofumpt
        )

Adding this block to a module's `go.mod` (via `go get -tool github.com/mirurobotics/gotools/cmd/miru@<version>`) allows that module to invoke the tool with `go tool miru <subcommand>`. The tool binary is built on demand and cached under `GOCACHE`.

### The 15 violations that must be fixed

Running `go run ./cmd/miru lint --paths=/home/ben/miru/workbench3/docs/tools/lint --fix=false` (from `/home/ben/miru/workbench3/gotools`) today reports:

    /home/ben/miru/workbench3/docs/tools/lint/main.go:56: line is 90 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/main.go:15: multi-line expression can be collapsed to single line
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:82: line is 105 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:87: line is 100 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:138: line is 96 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:193: line is 105 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:198: line is 103 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:23: multi-line expression can be collapsed to single line
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:28: multi-line expression can be collapsed to single line
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:33: multi-line expression can be collapsed to single line
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:62: multi-line expression can be collapsed to single line
    /home/ben/miru/workbench3/docs/tools/lint/rules_test.go:67: multi-line expression can be collapsed to single line
    /home/ben/miru/workbench3/docs/tools/lint/scanner.go:6: line is 90 columns wide, exceeds 88-column limit
    /home/ben/miru/workbench3/docs/tools/lint/scanner.go:30: function body can be inlined to single line
    /home/ben/miru/workbench3/docs/tools/lint/scanner.go:35: function body can be inlined to single line

Total: 15 violations from the custom linter. `gofumpt` and `golangci-lint` are both clean. `scanner_test.go` and `rules.go` have no violations.

A dry-run of `--fix` mode on a copy of the module showed:

- 8 violations auto-fix: all 5 "multi-line expression can be collapsed" in `rules_test.go`, both "function body can be inlined" in `scanner.go`, and `main.go:15` (which collapses `rules := []Rule{ NoDoubleDash{}, }` to `rules := []Rule{NoDoubleDash{}}`).
- 7 violations remain and need manual wrapping, all "line is N columns wide":
  - `main.go` (after collapse): 1 long line where `violations = append(violations, rule.Check(path, scanner.LineNum(), spans)...)` is 90 columns wide. Wrap by extracting the `rule.Check(...)` call into a local variable.
  - `rules_test.go`: 5 long lines — two `t.Errorf("expected %d violations, got %d: %v", ...)` calls (93 columns each), two `t.Errorf("violation %d: expected col %d, got %d", ...)` / `t.Errorf("violation %d: expected line %d, got %d", ...)` calls (85-88 columns each), and one test table entry with a long `content:` field containing `"---\ntitle: Test\n---\n\n<ParamField path=\"--version\" type=\"string\">"` (96 columns). Wrap the `t.Errorf` calls across multiple lines; for the test table entry, split the content string into concatenated pieces or move it to a package-level constant.
  - `scanner.go:6`: the comment `// StartCol is the 1-based byte offset of the span's first character in the original line.` is 90 columns wide. Rewrap to two lines.

Note that after auto-fix, some line numbers shift. The above mapping is the "what needs to be done conceptually" — when actually running the lint, re-read the violation output, because line numbers will differ from the pre-fix list.

### Existing CI shape

From `/home/ben/miru/workbench3/docs/.github/workflows/lint.yml`:

        jobs:
          lint:
            runs-on: ubuntu-latest
            steps:
              - name: Check out repository
                uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2
              - name: Set up Go
                uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6.4.0
                with:
                  go-version-file: tools/lint/go.mod
              - name: Enable corepack
                run: corepack enable
              - name: Set up Node.js
                uses: actions/setup-node@53b83947a5a98c8d113130e565377fae1a50d02f # v6.3.0
                with:
                  node-version: 22
                  cache: pnpm
                  cache-dependency-path: pnpm-lock.yaml
              - name: Install dependencies
                run: pnpm install --frozen-lockfile
              - name: Ensure scripts are executable
                run: chmod +x scripts/lint.sh scripts/audit.sh tests/test-lint.sh
              - name: Run lint smoke tests
                run: pnpm run test:lint
              - name: Run documentation lint
                run: ./scripts/lint.sh

The `Set up Go` step already picks up the Go toolchain from `tools/lint/go.mod`, so bumping that file's `go` directive to `1.25.3` automatically upgrades CI's Go toolchain to match — no workflow change needed for the Go version bump. Only a new run step is needed.

## Plan of Work

The change divides into six milestones. Each milestone produces exactly one commit inside the `docs/` repo.

### Milestone 1: Tool directive + helper script

Edit `/home/ben/miru/workbench3/docs/tools/lint/go.mod`:

- Change `go 1.24` to `go 1.25.3` on line 3.
- After the new `go 1.25.3` line, add (separated by a blank line):

        tool (
            github.com/mirurobotics/gotools/cmd/miru
        )

- Do not hand-edit the indirect requirement block — let `go get` populate it.

From inside `/home/ben/miru/workbench3/docs/tools/lint/`, run:

        go get -tool github.com/mirurobotics/gotools/cmd/miru@latest

This resolves the gotools version, downloads it and its transitive dependencies, and writes them to `go.sum`. If a specific version is preferred (e.g. `v0.1.3` which is what `cli-private` pins to), substitute `@v0.1.3`.

Create `/home/ben/miru/workbench3/docs/tools/lint/scripts/lint.sh` with tab indentation and the following content:

        #!/usr/bin/env bash
        set -euo pipefail

        usage() {
            cat <<'EOF'
        Usage: LINT_FIX=0 ./tools/lint/scripts/lint.sh

        Runs gotools (custom linter + gofumpt + golangci-lint) against
        docs/tools/lint/. Set LINT_FIX=0 to run in check-only mode (for CI
        and preflight); omit or set LINT_FIX=1 for auto-fix mode (default
        for local runs).
        EOF
        }

        case "${1:-}" in
            -h|--help)
                usage
                exit 0
                ;;
        esac

        script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
        lint_dir="$(cd -- "${script_dir}/.." && pwd)"

        FIX="--fix"
        if [ "${LINT_FIX:-1}" = "0" ]; then
            FIX="--fix=false"
        fi

        cd "${lint_dir}"
        exec go tool miru lint --paths=. ${FIX}

Mark it executable with `chmod +x docs/tools/lint/scripts/lint.sh`.

Note the use of `--paths=.` — from inside `docs/tools/lint/`, the current directory IS the module, so passing `.` tells the custom linter to scan `main.go`, `rules.go`, `rules_test.go`, `scanner.go`, and `scanner_test.go`.

**Commit Milestone 1:** from `/home/ben/miru/workbench3/docs/`, stage `tools/lint/go.mod`, `tools/lint/go.sum`, and `tools/lint/scripts/lint.sh`, then commit with message "lint: add gotools tool directive and per-tool lint script".

### Milestone 2: Run auto-fix

From `/home/ben/miru/workbench3/docs/`:

        LINT_FIX=1 ./tools/lint/scripts/lint.sh

Expected: the custom linter prints "fixed: <path>" for each file it auto-fixes (at least `main.go`, `rules_test.go`, and `scanner.go`), then re-runs and reports the remaining (line-width) violations. The script exits non-zero because violations remain — this is expected for this milestone.

Auto-fix is idempotent: running twice produces the same result as running once. If the script hangs or errors out on network access, re-run with `GOFLAGS=-mod=mod` to force module resolution.

Inspect the diff to confirm:

- `main.go` line 15 area: `rules := []Rule{NoDoubleDash{}}` is now a single line.
- `rules_test.go`: the five flagged test table entries (`triple dash not flagged`, `quadruple dash not flagged`, `single dash not flagged`, `no dashes`, `empty string`) are now single-line literals.
- `scanner.go`: `NewScanner` and `LineNum` are one-line function bodies.

Then run the tests to confirm nothing regressed:

        cd tools/lint
        go test ./...

Expected: `ok  	github.com/mirurobotics/docs/tools/lint	0.0??s`.

Return to `docs/` root.

**Commit Milestone 2:** stage the modified `.go` files under `tools/lint/`, commit with message "lint(tools): auto-fix collapsible expressions and inlinable bodies".

### Milestone 3: Manually fix remaining line-width violations

From `/home/ben/miru/workbench3/docs/`, rerun:

        LINT_FIX=0 ./tools/lint/scripts/lint.sh

Read the remaining violations. Expect roughly 7 "line is N columns wide" findings (exact line numbers depend on what the auto-fix did in Milestone 2). For each one, edit the source file to bring it under 88 columns visually (tabs count as 4). Strategies:

- **`main.go` long append line**: rewrite

        violations = append(violations, rule.Check(path, scanner.LineNum(), spans)...)

  as

        out := rule.Check(path, scanner.LineNum(), spans)
        violations = append(violations, out...)

  Alternatively, cache `scanner.LineNum()` into a local variable before the inner `for _, rule := range rules` loop and inline that into the call. Pick whichever reads more naturally.

- **`rules_test.go` `t.Errorf` calls**: wrap the format-string calls across multiple lines, e.g.

        t.Errorf(
            "expected %d violations, got %d: %v",
            tt.wantCount, len(violations), violations,
        )

  `gofumpt` will accept this layout.

- **`rules_test.go` long `content:` field**: the 96-column line is inside a test table entry. Break the string into concatenated pieces, e.g.

        content: "---\ntitle: Test\n---\n\n" +
            "<ParamField path=\"--version\" type=\"string\">",

  or hoist it to a package-level `const` and reference by name.

- **`scanner.go:6` long doc comment**: rewrap the comment to two lines, e.g.

        // StartCol is the 1-based byte offset of the span's first
        // character in the original line.

After each edit, rerun `LINT_FIX=0 ./tools/lint/scripts/lint.sh` from `/home/ben/miru/workbench3/docs/` until it prints `0 issues.` and exits with status 0. The final expected transcript tail is:

        Running custom linter on /home/ben/miru/workbench3/docs/tools/lint...
        0 violation(s) found in ...
        Running gofumpt...
        Running golangci-lint...
        0 issues.

Then re-run the unit tests from `docs/tools/lint/`:

        cd tools/lint && go test ./... && cd ../..

Expected: `ok` line with no failures.

**Commit Milestone 3:** stage the modified `.go` files, commit with message "lint(tools): manually wrap remaining long lines".

### Milestone 4: Wire into `scripts/preflight.sh`

Edit `/home/ben/miru/workbench3/docs/scripts/preflight.sh`. Before the existing `echo "=== Lint ==="` block, insert a new block:

        echo "=== Go Lint (tools/lint) ==="
        LINT_FIX=0 ./tools/lint/scripts/lint.sh
        echo ""

The full new file should look like:

        #!/usr/bin/env bash
        set -euo pipefail

        REPO_ROOT=$(git rev-parse --show-toplevel)
        cd "$REPO_ROOT"

        echo "=== Lint Smoke Tests ==="
        pnpm run test:lint
        echo ""

        echo "=== Go Lint (tools/lint) ==="
        LINT_FIX=0 ./tools/lint/scripts/lint.sh
        echo ""

        echo "=== Lint ==="
        ./scripts/lint.sh
        echo ""

        echo "=== Audit ==="
        ./scripts/audit.sh

Run the full preflight from the docs repo root:

        ./scripts/preflight.sh

Expected: preflight prints the four section headers in order, each section completes successfully, the Go Lint section prints `0 issues.`, and the script exits 0.

**Commit Milestone 4:** stage `scripts/preflight.sh`, commit with message "lint: gate preflight on tools/lint Go lint".

### Milestone 5: Wire into `.github/workflows/lint.yml`

Edit `/home/ben/miru/workbench3/docs/.github/workflows/lint.yml`. In the `lint` job, after the existing "Ensure scripts are executable" step and before the "Run lint smoke tests" step, add:

        - name: Ensure tools/lint script is executable
          run: chmod +x tools/lint/scripts/lint.sh

        - name: Run Go lint on tools/lint
          env:
            LINT_FIX: "0"
          run: ./tools/lint/scripts/lint.sh

Rationale for placement: running the Go lint early (before pnpm-dependent steps) means CI fails fast on Go issues; it also matches the preflight ordering. The `LINT_FIX: "0"` env var forces check-only mode.

Do not touch the `audit` job.

Verify the YAML parses by running:

        python3 -c "import yaml, sys; yaml.safe_load(open('.github/workflows/lint.yml'))" && echo OK

Expected: `OK`.

Verify the new step names appear with a literal search:

        grep -n 'Run Go lint on tools/lint' .github/workflows/lint.yml

Expected: one line with the match.

**Commit Milestone 5:** stage `.github/workflows/lint.yml`, commit with message "ci(lint): run tools/lint Go lint in CI".

### Milestone 6: Full preflight + final commit

From `/home/ben/miru/workbench3/docs/`:

        ./scripts/preflight.sh

Expected: the script prints the four section headers, each section exits 0, and the transcript ends roughly with:

        All documentation lint checks passed.

If preflight reports anything other than clean — any section fails or prints non-zero issue counts — stop and fix the root cause before proceeding. This is the hard gate (see Validation and Acceptance below).

If there are no additional code changes to commit at this milestone (which is the expected steady state since the previous milestones already committed their artifacts), Milestone 6 has no commit and progress is marked done. If preflight did surface a late fix, stage only that fix and commit with message "lint: fix preflight-surfaced issue in <file>".

## Concrete Steps

All commands use absolute paths. The working directory is shown before each command.

### Preconditions

Working directory: `/home/ben/miru/workbench3/docs`.

    git rev-parse --show-toplevel

Expected output: `/home/ben/miru/workbench3/docs`.

    go version

Expected: `go version go1.25.3 linux/amd64` or newer. If older, install Go 1.25.3 before continuing — without it, `go get -tool` will refuse the version bump.

    git status --short

Expected: either clean or only the plan file itself. If there are unrelated changes, stash them.

### Milestone 1 commands

Working directory: `/home/ben/miru/workbench3/docs/tools/lint`.

Edit `go.mod` by hand to change `go 1.24` to `go 1.25.3`:

    sed -i 's/^go 1\.24$/go 1.25.3/' go.mod

Verify:

    head -5 go.mod

Expected:

    module github.com/mirurobotics/docs/tools/lint

    go 1.25.3

Add the tool directive:

    go get -tool github.com/mirurobotics/gotools/cmd/miru@latest

Expected output: `go: added github.com/mirurobotics/gotools v<X.Y.Z>` plus a long list of `go: added <indirect>` lines. `go.mod` now contains a `tool ( ... )` block and a large `require` block. `go.sum` now exists (or has many new entries).

Verify the directive is present:

    grep -A 2 '^tool (' go.mod

Expected:

    tool (
        github.com/mirurobotics/gotools/cmd/miru
    )

Create the per-tool lint script:

Working directory: `/home/ben/miru/workbench3/docs`.

    mkdir -p tools/lint/scripts

Then create `tools/lint/scripts/lint.sh` with the content shown in the Plan of Work, Milestone 1 section. Make it executable:

    chmod +x tools/lint/scripts/lint.sh

Verify the script at least parses:

    bash -n tools/lint/scripts/lint.sh && echo OK

Expected: `OK`.

Smoke-test it (this may still report violations — that is expected until Milestone 3 is complete):

    LINT_FIX=0 ./tools/lint/scripts/lint.sh || true

Expected: the script runs, prints "Running custom linter on /home/ben/miru/workbench3/docs/tools/lint...", lists 15 violations, runs gofumpt and golangci-lint successfully, and exits non-zero. That non-zero is fine at this milestone; it proves the wiring works.

Commit:

    git add tools/lint/go.mod tools/lint/go.sum tools/lint/scripts/lint.sh
    git commit -m "lint: add gotools tool directive and per-tool lint script"

Verify:

    git log -1 --stat

### Milestone 2 commands

Working directory: `/home/ben/miru/workbench3/docs`.

    LINT_FIX=1 ./tools/lint/scripts/lint.sh || true

Expected: the script prints lines like `fixed: /home/ben/miru/workbench3/docs/tools/lint/main.go`, `fixed: .../rules_test.go`, `fixed: .../scanner.go`, then reports the remaining line-width violations (roughly 7), and exits non-zero. Still fine.

Verify the fixes took effect by inspecting the diff:

    git diff --stat tools/lint

Expected: `main.go`, `rules_test.go`, `scanner.go` all show a non-zero line delta.

Run the unit tests:

    cd tools/lint && go test ./... && cd ..

Expected: `ok  	github.com/mirurobotics/docs/tools/lint	0.???s`.

Commit:

    git add tools/lint/main.go tools/lint/rules_test.go tools/lint/scanner.go
    git commit -m "lint(tools): auto-fix collapsible expressions and inlinable bodies"

### Milestone 3 commands

Working directory: `/home/ben/miru/workbench3/docs`.

    LINT_FIX=0 ./tools/lint/scripts/lint.sh 2>&1 | tee /tmp/docs-go-lint.txt

Expected: ~7 "line is N columns wide" violations listed, exit non-zero.

For each violation, open the file, rewrite the flagged line per the strategies in Plan of Work Milestone 3, save, then rerun:

    LINT_FIX=0 ./tools/lint/scripts/lint.sh

Repeat until the output is:

    Running custom linter on /home/ben/miru/workbench3/docs/tools/lint...
    0 violation(s) found in /home/ben/miru/workbench3/docs/tools/lint.
    Running gofumpt...
    Running golangci-lint...
    0 issues.

And the exit status is 0.

Sanity-check tests again:

    cd tools/lint && go test ./... && cd ..

Expected: `ok`.

Commit:

    git add tools/lint/main.go tools/lint/rules_test.go tools/lint/scanner.go
    git commit -m "lint(tools): manually wrap remaining long lines"

### Milestone 4 commands

Working directory: `/home/ben/miru/workbench3/docs`.

Edit `scripts/preflight.sh` to insert the new block per Plan of Work Milestone 4.

Verify the script parses:

    bash -n scripts/preflight.sh && echo OK

Expected: `OK`.

Verify the new step is present:

    grep -n 'Go Lint (tools/lint)' scripts/preflight.sh

Expected:

    <N>:echo "=== Go Lint (tools/lint) ==="

Run it end-to-end:

    ./scripts/preflight.sh

Expected: all four sections pass, final line is `All documentation lint checks passed.`, exit status 0.

Commit:

    git add scripts/preflight.sh
    git commit -m "lint: gate preflight on tools/lint Go lint"

### Milestone 5 commands

Working directory: `/home/ben/miru/workbench3/docs`.

Edit `.github/workflows/lint.yml` per Plan of Work Milestone 5.

Verify YAML:

    python3 -c "import yaml; yaml.safe_load(open('.github/workflows/lint.yml'))" && echo OK

Expected: `OK`.

Confirm the step name landed:

    grep -n 'Run Go lint on tools/lint' .github/workflows/lint.yml

Expected: one match.

Confirm the `LINT_FIX: "0"` env var landed:

    grep -n 'LINT_FIX: "0"' .github/workflows/lint.yml

Expected: one match.

Confirm the existing step name is still present (sanity check that the existing workflow was not clobbered):

    grep -n 'Run documentation lint' .github/workflows/lint.yml

Expected: still one match.

Commit:

    git add .github/workflows/lint.yml
    git commit -m "ci(lint): run tools/lint Go lint in CI"

### Milestone 6 commands

Working directory: `/home/ben/miru/workbench3/docs`.

Run the full preflight:

    ./scripts/preflight.sh

Expected transcript tail:

    === Go Lint (tools/lint) ===
    ... 0 issues.

    === Lint ===
    ... All documentation lint checks passed.

    === Audit ===
    ...

Exit status 0. If any section prints non-zero issue counts or the script exits non-zero, treat it as a hard gate failure (see Validation and Acceptance) and fix the root cause.

If a late fix was needed, stage and commit only that fix; otherwise this milestone is a no-op commit-wise, and Progress is simply checked off.

    git log --oneline -7

Expected: at least the five milestone commits in reverse chronological order, plus any earlier commits.

## Validation and Acceptance

**Hard gate: Preflight must report clean before any PR is opened.** This means `./scripts/preflight.sh` from `/home/ben/miru/workbench3/docs/` must exit with status 0 and print successful completion for all four sections (`Lint Smoke Tests`, `Go Lint (tools/lint)`, `Lint`, `Audit`). This is non-negotiable. If any later subagent runs preflight and sees a non-zero exit or non-clean output, that subagent must stop implementing and fix the root cause before opening a PR. Downstream skills (`$implement`, `$preflight`, `$pr`) must respect this gate.

Acceptance criteria (observable behavior a reviewer or subagent can verify):

1. **Go lint runs standalone.** From `/home/ben/miru/workbench3/docs/`:

        LINT_FIX=0 ./tools/lint/scripts/lint.sh; echo "exit=$?"

   Expected: final lines include `0 issues.` and `exit=0`.

2. **Go tests pass.** From `/home/ben/miru/workbench3/docs/tools/lint/`:

        go test ./...

   Expected: `ok  	github.com/mirurobotics/docs/tools/lint	0.???s`. Zero test failures. Before this plan, the tests pass as well, so this is a regression gate rather than a new pass — specifically, the test counts of `TestNoDoubleDash` (10 subtests), `TestNoDoubleDashWithOffset` (1), and `TestNoDoubleDashIntegration` (12+ subtests) plus `scanner_test.go` tests must all still pass post-fix.

3. **Preflight runs clean.** From `/home/ben/miru/workbench3/docs/`:

        ./scripts/preflight.sh; echo "exit=$?"

   Expected: all four section headers visible in order (`=== Lint Smoke Tests ===`, `=== Go Lint (tools/lint) ===`, `=== Lint ===`, `=== Audit ===`), followed by `exit=0`.

4. **Preflight wiring is present and ordered correctly.** From `/home/ben/miru/workbench3/docs/`:

        grep -n 'Go Lint (tools/lint)' scripts/preflight.sh
        awk '/=== Lint Smoke Tests ===/{a=NR} /=== Go Lint/{b=NR} /=== Lint ===/{c=NR} /=== Audit ===/{d=NR} END{print a,b,c,d}' scripts/preflight.sh

   Expected: the first command prints one matching line; the second prints four ascending line numbers.

5. **CI wiring is present.** From `/home/ben/miru/workbench3/docs/`:

        grep -n 'Run Go lint on tools/lint' .github/workflows/lint.yml
        grep -n 'LINT_FIX: "0"' .github/workflows/lint.yml

   Expected: one match per grep.

6. **Negative test — lint catches a fresh violation.** To prove the gate is live, introduce a temporary 100-column line at the top of `docs/tools/lint/rules.go` (for example a very long comment), then run:

        LINT_FIX=0 ./tools/lint/scripts/lint.sh; echo "exit=$?"

   Expected: `line is ??? columns wide` appears in the output and `exit=1`. Revert the temporary edit; confirm `exit=0` once more. This step is a hand-run smoke check and should not be committed.

7. **The CI lint workflow still finishes end-to-end** after the change. If running in a feature branch with CI enabled, the `lint` job in GitHub Actions should reach a green state and show the `Run Go lint on tools/lint` step in the job log.

## Idempotence and Recovery

- **Re-running the Milestone 1 `go get -tool` command** is safe: Go sees the directive is already present and either no-ops or upgrades to a newer resolution. If you need to pin to an older version, run `go get -tool github.com/mirurobotics/gotools/cmd/miru@<version>` explicitly.

- **Re-running auto-fix** (`LINT_FIX=1 ./tools/lint/scripts/lint.sh`) is safe: the custom linter rewrites files deterministically, so a second run produces no further diff after the first converges.

- **Manual long-line fixes** can be iterated: make an edit, rerun `LINT_FIX=0 ./tools/lint/scripts/lint.sh`, repeat. No state is persisted between runs beyond the Go build cache.

- **Rollback for the `go 1.25.3` bump.** If the bump breaks something unexpected, revert Milestone 1's commit with `git revert <sha>`. This removes the `tool` directive, removes the `go.sum` entries, and restores `go 1.24`. The per-tool lint script can stay in place since it only fails when invoked, not at file-load time.

- **Rollback for the preflight/CI wiring.** Each milestone is committed separately. To roll back just the preflight or CI change, `git revert <sha>` the relevant commit. The Go lint fixes can stay even if the wiring is reverted — they are pure improvements to the source.

- **Recovery from a mid-milestone failure.** If the process dies between edits and a commit, `git status` will show the in-progress changes. Re-run the milestone's commands; they are designed to be idempotent. If the working tree is corrupted, `git stash` the WIP and re-plan from the last clean commit.

- **Failed `go tool miru lint` invocation with "tool not declared"** usually means the helper script was run from a directory outside the `tools/lint` module. The script already `cd`s into `${lint_dir}` — if editing it, keep that `cd`. If running `go tool miru lint` directly, first `cd` into `/home/ben/miru/workbench3/docs/tools/lint`.

- **Failed YAML parse in CI** after editing `.github/workflows/lint.yml` is caught locally by the `python3 -c "import yaml; yaml.safe_load(...)"` check in Milestone 5. Always run that check before committing.

# Remove redundant bad-redirects bash fixture

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|------------|--------|-------------|
| `docs/` | read-write | Delete fixture directory and corresponding bash test assertions. |

This plan lives in `docs/plans/backlog/` because all changes are confined to the docs repo. Working dir: `/home/ben/miru/workbench1/repos/docs/`. Branch: `feat/redirect-lint-rule` (PR #75 open).

## Purpose / Big Picture

The redirect lint rule's behavior is fully covered by Go tests (`tools/lint/linter/redirects/redirects_test.go` with 76 subtests at 100% covgate, plus `tools/lint/main_test.go::TestRun/redirect_violation_returns_1` end-to-end). The bash-driven `bad-redirects` integration fixture adds no coverage and just slows the suite. After this change, `pnpm run test:lint` runs only the still-useful fixtures (`good`, `bad-mdx`, `bad-spelling`, `bad-openapi`) and the redirect rule remains fully tested via Go.

## Progress

- [ ] Sanity-grep `bad-redirects` lines in `tests/test-lint.sh`
- [ ] `git rm -r tests/lint-fixtures/bad-redirects`
- [ ] Remove the 11-line `run_expect_fail "bad-redirects" ...` block from `tests/test-lint.sh`
- [ ] `pnpm run test:lint` exits 0
- [ ] `go test ./tools/lint/...` passes
- [ ] `pnpm run lint` exits 0
- [ ] `./scripts/preflight.sh` clean (modulo pre-existing local bats issue)
- [ ] Commit

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

(Add entries as you go.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

- Redirect rule: `tools/lint/linter/redirects/redirects.go`, invoked from `run()` in `tools/lint/main.go` via `redirects.Check(contentRoot)`.
- Go test coverage: `tools/lint/linter/redirects/redirects_test.go` (table-driven `TestValidate`, 76 subtests, 100% covgate); `tools/lint/main_test.go::TestRun/redirect_violation_returns_1` builds a tempdir with a `docs.json` containing a missing-destination redirect, calls `run()`, asserts exit 1 and stdout contains `missing destination`.
- Bash integration runner: `tests/test-lint.sh` shells out the `lint` binary against fixtures under `tests/lint-fixtures/<name>/`. The `bad-redirects` fixture and its 11 `run_expect_fail "bad-redirects" ...` lines (currently lines 58-68) duplicate Go-side coverage.
- Files to delete under `tests/lint-fixtures/bad-redirects/`: `docs.json`, `docs/admin/exists.mdx`, `docs/example.mdx`, `docs/references/example.yaml`, `docs/wild/page.mdx`, `snippets/.gitkeep`.
- Untouched: the `redirects` Go package; `tools/lint/main.go` and `tools/lint/main_test.go`; the four other fixtures (`good`, `bad-mdx`, `bad-spelling`, `bad-openapi`); all other lines in `tests/test-lint.sh`.

## Plan of Work

Single milestone, executed in order.

1. Confirm the assumed line range of the `bad-redirects` block in `tests/test-lint.sh` with grep.
2. Stage deletion of the entire `tests/lint-fixtures/bad-redirects/` directory via `git rm -r`.
3. Edit `tests/test-lint.sh` to remove only the 11 `run_expect_fail "bad-redirects" ...` lines (lines 58-68 per the assumption from step 1; adjust if grep shows different). Leave all other lines unchanged.
4. Re-grep to confirm zero `bad-redirects` references remain in `tests/test-lint.sh`.
5. Run `pnpm run test:lint`, then `go test ./tools/lint/...`, then `pnpm run lint`, then `./scripts/preflight.sh`. All must succeed (preflight modulo pre-existing local `bats: command not found`; CI installs bats).
6. Commit the change.

## Concrete Steps

All commands from `/home/ben/miru/workbench1/repos/docs/`.

Step 1. Sanity-grep:

    grep -n 'bad-redirects' tests/test-lint.sh

Expected: 11 lines, contiguous, around lines 58-68, all of the form `run_expect_fail "bad-redirects" ...`.

Step 2. Delete the fixture directory:

    git rm -r tests/lint-fixtures/bad-redirects

Expected: six `rm` lines, one per file listed in Context and Orientation.

Step 3. Edit `tests/test-lint.sh` to remove the 11 lines from step 1. Use the line range from step 1's actual output (do not hardcode if grep showed something different).

Step 4. Verify removal:

    grep -n 'bad-redirects' tests/test-lint.sh

Expected: no output (exit code 1).

Step 5. Run validation:

    pnpm run test:lint
    go test ./tools/lint/...
    pnpm run lint
    ./scripts/preflight.sh

Expected: all exit 0; `pnpm run test:lint` output references only `good`, `bad-mdx`, `bad-spelling`, `bad-openapi`. Preflight may emit `bats: command not found` locally — this matches main and is acceptable; CI installs bats.

Step 6. Commit (the `git rm` is already staged):

    git add tests/test-lint.sh
    git commit -m "test(lint): remove bad-redirects bash fixture (covered by Go tests)"

## Validation and Acceptance

- `pnpm run test:lint` exits 0 and the runner output lists only `good`, `bad-mdx`, `bad-spelling`, `bad-openapi`. No `bad-redirects` references appear.
- `go test ./tools/lint/...` passes including covgate; this confirms the Go-side coverage we are relying on is still green (no Go file changed).
- `pnpm run lint` exits 0 against the live repo.
- `./scripts/preflight.sh` reports clean (modulo the pre-existing local `bats: command not found` env issue that also exists on `main`). This is a mandatory gate before pushing.
- `grep -n 'bad-redirects' tests/test-lint.sh` returns no matches; `ls tests/lint-fixtures/bad-redirects 2>&1` reports "No such file or directory".
- PR #75 description should be updated post-merge to drop fixture references — that belongs to the deliver step in the task workflow, not this plan.

## Idempotence and Recovery

- `git rm -r` on an already-deleted path is a no-op against working tree but errors on missing index entry; safe to re-run only if the prior commit hasn't landed. After a partial failure, inspect `git status` and resume from the next unfinished step.
- The `tests/test-lint.sh` edit is idempotent: re-running the grep + removal sequence is safe because matching lines are absent after the first successful pass.
- No data migration, no destructive infra change. Recovery path: `git revert <commit-sha>` or `git checkout HEAD~1 -- tests/lint-fixtures/bad-redirects tests/test-lint.sh` restores the fixture and the bash assertions verbatim.

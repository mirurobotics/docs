# Add scripts/codeql-prune.sh to remove stale CodeQL code-scanning configurations

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.


## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/user/docs`) | read-write | Add `scripts/codeql-prune.sh` and `scripts/codeql-prune_test.bats`; add a comment to `.github/workflows/codeql-analysis.yml`; wire the new bats file into `.github/workflows/ci.yml` (`shell-tests` job) and `scripts/preflight.sh`. |

This plan lives in `docs/plans/` because every write is to this repo.

Working branch: `claude/upbeat-davinci-tyzi5x` (already checked out; equals `origin/main` at `b697e10`). Do not create another branch. Base branch for the PR: `main`.


## Purpose / Big Picture

On `https://github.com/mirurobotics/docs/security/code-scanning/tools` (Security > Code scanning > Tools > CodeQL) three configurations are listed. One is healthy (`.github/workflows/codeql-analysis.yml:codeql`, scanned within the last few days). Two are stale and flagged "Code Scanning results may be out of date":

- `.github/workflows/codeql-analysis.yml:analyze/language:javascript`
- `.github/workflows/codeql-analysis.yml:codeql/language:javascript`

Both were last scanned 2026-03-23. Nothing in the workflow file can remove them; GitHub never garbage-collects old configurations. They can only go away by deleting their analyses through the UI or the REST API, which needs a token with code-scanning write access that this automated session does not have (the endpoint returns HTTP 403 `Resource not accessible by integration` for it).

After this change the repo ships `scripts/codeql-prune.sh`. A maintainer with the right token runs it once to list the stale categories (dry run), then once more with `--delete` to remove them, after which the Tools page shows only the healthy configuration. The workflow file carries a comment so the next person who renames the job knows to run the script, and the script is covered by a bats test that runs in CI.


## Progress

- [x] Milestone 1: Write `scripts/codeql-prune.sh`; commit. (2026-09-20)
- [x] Milestone 2: Write `scripts/codeql-prune_test.bats`; wire into `ci.yml` and `preflight.sh`; commit. (2026-09-20)
- [x] Milestone 3: Add the job-id comment to `codeql-analysis.yml`; commit. (2026-09-20)
- [ ] Milestone 4: Push, open the draft PR, drive CI to CLEAN, mark ready for review.
- [ ] Post-merge (maintainer, not a PR gate): run the script with `--delete`; confirm the Tools page shows one configuration.


## Surprises & Discoveries

- 2026-09-20: `bats` and `shellcheck` installed fine via `apt-get` in the implementation session, so the test ran locally: `bats scripts/codeql-prune_test.bats` reports `5 tests, 0 failures` and `shellcheck scripts/codeql-prune.sh` is clean.
- 2026-09-20: The combined CI command `bats pub/scripts/agent/check-miru-access_test.bats scripts/codeql-prune_test.bats` fails locally only because the session runs as uid 0 and `check-miru-access.sh` refuses to run as root (18 pre-existing tests fail for that reason). CI runners are non-root, so this is environmental; all 4 new tests pass in the combined run.


## Decision Log

- Decision: Hardcode the default keep category (`.github/workflows/codeql-analysis.yml:codeql`) in the script instead of parsing the workflow YAML. / Rationale: parsing job ids out of YAML with grep is fragile and the workflow comment already tells whoever renames the job to update the script; `--keep` overrides the default when needed. / Date/Author: 2026-09-20, plan author.
- Decision: The script only considers analyses whose `tool.name` is `CodeQL`. / Rationale: the default keep list is CodeQL-specific; if another SARIF uploader is ever added its categories must not be pruned by accident. / Date/Author: 2026-09-20, plan author.
- Decision: No README change. / Rationale: the root `README.md` is a content index for the docs site with no contributor-tooling section; the script's `--help` is its documentation. / Date/Author: 2026-09-20, plan author.
- Decision: `die` prefixes messages with `codeql-prune.sh:` and `process_ref` prints a `pruned N categories` summary line in `--delete` mode. / Rationale: the prefix makes stderr attributable when the script is chained in CI; the summary keeps the dry-run and delete outputs symmetric. Neither changes the strings the tests assert on. / Date/Author: 2026-09-20, implementer.
- Decision: `process_ref` captures the listing in a variable instead of reading from a process substitution, and a fifth bats test (`gh api failure aborts instead of reporting nothing to prune`) covers it. / Rationale: `set -e` does not see failures inside `< <(...)`, so an HTTP 403 from `gh api` would have printed `nothing to prune` and exited 0; assigning the command substitution to a variable makes the failure fatal as the plan intended. / Date/Author: 2026-09-20, implementer.
- Decision: The deletion itself is out of scope for the PR. / Rationale: the session token gets HTTP 403 on `/code-scanning/analyses`; the PR delivers the tool and the PR body carries the run instructions for a maintainer. / Date/Author: 2026-09-20, plan author.


## Outcomes & Retrospective

(Summarize at completion.)


## Context and Orientation

**Code-scanning categories.** When the CodeQL action uploads a SARIF file, GitHub tags it with a *category* string derived from the workflow: `<workflow path>:<job id>` plus `/<matrix key>:<value>` for each matrix variable. GitHub shows one "configuration" per category on the Tools page and keeps every category it has ever seen. In this repo commit `36590ca` (2026-03-23, #13) renamed the job `analyze` -> `codeql` and commit `93f4af4` (2026-03-24, #14) removed the `language: [javascript]` matrix; each change created a new category and orphaned the previous one. The current workflow `.github/workflows/codeql-analysis.yml` has a single job with id `codeql`, `languages: javascript`, no matrix, so the live category is `.github/workflows/codeql-analysis.yml:codeql`. It runs on pushes to `main`, `staging`, `uat`, `production`, weekly, and via `workflow_call`, so the stale categories may also exist on the three environment branches.

**Deleting analyses (REST).** `GET /repos/{owner}/{repo}/code-scanning/analyses?ref=refs/heads/main&per_page=100` lists analyses for a ref; each object has `id`, `ref`, `commit_sha`, `category`, `created_at`, `deletable`, `tool.name`, `url`. Analyses form a chain per category; only the most recent one in a chain can be deleted. `DELETE /repos/{owner}/{repo}/code-scanning/analyses/{id}?confirm_delete` deletes that analysis and returns `{ "next_analysis_url": ..., "confirm_delete_url": ... }` pointing at the next-older analysis of the same category; both are `null` when the chain is exhausted. So pruning a category = take its newest analysis, delete, follow `confirm_delete_url` until null. `deletable: false` analyses are skipped. Both endpoints need a token with `security-events: write` (fine-grained: "Code scanning alerts: Read and write"; classic: `repo` + `security_events`). This session's token gets 403, so the script is tested against a stub, and the real run is done by a maintainer.

**Repo conventions.** Shell scripts live in `scripts/` (`audit.sh`, `lint.sh`, `preflight.sh`, `promote.sh`). `promote.sh` is the option-taking sibling: `#!/usr/bin/env bash`, `set -euo pipefail`, `usage()` heredoc printed on `-h`/`--help`, exit 1 on error. The one existing bats file is `pub/scripts/agent/check-miru-access_test.bats` (tab-indented, `run bash "$SCRIPT" ...`, asserts on `$status` and `$output`, `-h`/`--help` return 0 and print `Usage:`); usage errors there exit 2. The new script and test follow those files plus the workbench script rules: tab indentation, `while/case` option parsing.

**CI.** `.github/workflows/ci.yml` runs on every PR: `lint` (`pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`), `audit`, `shell-tests` (`sudo apt-get install -y bats` then `bats pub/scripts/agent/check-miru-access_test.bats`), and two `tools/lint` Go jobs that are skipped unless `tools/lint/**` changes. `scripts/preflight.sh` chains the same checks locally and ends with the same `bats` line. Runners are `blacksmith-2vcpu-ubuntu-2404` (GitHub runner-image compatible: `jq` and `gh` preinstalled). `./scripts/lint.sh` (MDX lint, ESLint, cspell, OpenAPI check) only looks at `docs/**/*.mdx`, `docs/rclone.md`, `docs/README.md`, and `docs/references/*.yaml`; it never touches `scripts/` or `.github/`, so the new files need no cspell words. **CLEAN** means every GitHub Actions check is green on the pushed branch head SHA.

**This environment.** `jq` and `curl` are installed; `gh`, `bats`, and `shellcheck` are not. `GH_TOKEN` is set but lacks code-scanning access. Watch CI with the GitHub MCP tools (`mcp__github__pull_request_read` with `get_status` / `get_checks`, `mcp__github__get_job_logs` with `failed_only: true`) or `curl -sS -H "Authorization: Bearer $GH_TOKEN" https://api.github.com/repos/mirurobotics/docs/commits/<sha>/check-runs`.


## Plan of Work

**New file `scripts/codeql-prune.sh`** (mode 0755, tab-indented, roughly 120 lines). Behaviour:

    Usage: codeql-prune.sh [--repo <owner/name>] [--ref <ref>]... [--keep <category>]... [--delete]

    Lists CodeQL code-scanning categories for a ref and prunes the ones not in
    the keep set. Dry run by default; --delete removes them.

    Options:
      --repo <owner/name>  Repository (default: gh repo view of the current checkout)
      --ref <ref>          Ref to inspect, repeatable (default: refs/heads/main)
      --keep <category>    Category to keep, repeatable (default:
                           .github/workflows/codeql-analysis.yml:codeql)
      --delete             Delete every analysis of each non-kept category
      -h, --help           Show this help

    Requires gh (authenticated with a token that has code scanning write
    access, e.g. GH_TOKEN) and jq.

    Exit codes: 0 success, 1 error, 2 usage error.

Structure, in order: `set -euo pipefail`; `usage()`; `die <code> <message>` (message to stderr, usage to stderr as well when the code is 2, then `exit <code>`); `while/case` option parsing into `repo`, `refs` array, `keeps` array, `delete=false`, starting from empty arrays and applying the defaults only after parsing when an array is still empty, so a `--ref`/`--keep` on the command line replaces the default rather than adding to it (`--repo`, `--ref`, `--keep` without a value -> `die 2 "<option> requires a value"`; anything else starting with `-` -> `die 2 "unknown option: <arg>"`); *then* the tool checks (`command -v gh` and `command -v jq`, each failing with `<tool> is required by codeql-prune.sh` and exit 1) — parsing before tool checks keeps `--help` working without `gh`, and the missing-tool test relies on the script running no external command before the check; default `repo` from `gh repo view --json nameWithOwner --jq .nameWithOwner`; then `main`.

Helper functions, each under 20 lines:

- `latest_per_category <repo> <ref>`: `gh api --paginate "repos/${repo}/code-scanning/analyses?ref=${ref}&per_page=100" | jq -r '...'` where the filter is `[.[] | select(.tool.name == "CodeQL")] | group_by(.category) | map(max_by(.created_at)) | .[] | [.category, .created_at, .commit_sha[0:7], .id, .deletable] | @tsv` wrapped in `jq -s 'add // []' | jq -r '<filter>'` (`--paginate` emits one JSON array per page; `-s 'add'` concatenates them). Output: one TSV line per category.
- `is_kept <category>`: loops over `keeps`, returns 0 on exact match.
- `prune_category <repo> <id>`: `url="repos/${repo}/code-scanning/analyses/${id}?confirm_delete"`; `while [[ -n "${url}" ]]; do url="$(gh api -X DELETE "${url}" --jq '.confirm_delete_url // empty')"; done`. `gh api` accepts the full `https://api.github.com/...` URL that `confirm_delete_url` contains. With `set -e`, a failed DELETE aborts the script with gh's error.
- `process_ref <repo> <ref>`: prints a header `== ${repo} ${ref} ==`, reads the TSV lines, and for each prints one of `KEEP`, `PRUNE`, `SKIP` (deletable is `false`) followed by category, `created_at`, short sha, and `id=<id>`; when `delete=true` it calls `prune_category` for every `PRUNE` line and prints `deleted <category>`; otherwise counts them and ends with `dry run: N categories would be pruned; re-run with --delete` (or `nothing to prune` when N is 0). A `gh api` failure (for example HTTP 403 from a token without code-scanning access) aborts the script with gh's own error because of `set -eo pipefail`.
- `main`: iterates `refs`, calling `process_ref`.

Do not add colours, tables, or a confirmation prompt: the dry-run default is the safety net.

**New file `scripts/codeql-prune_test.bats`** (tab-indented, `#!/usr/bin/env bats`). `SCRIPT="$BATS_TEST_DIRNAME/codeql-prune.sh"`. `setup()` creates `STUB_BIN="$BATS_TEST_TMPDIR/bin"` and an empty `$BATS_TEST_TMPDIR/empty`, writes a fake `gh` into `STUB_BIN` (shebang `#!/usr/bin/env bash`, chmod +x) and exports `GH_STUB_LOG="$BATS_TEST_TMPDIR/gh.log"`. Every test invokes the script through `run env PATH=... bash "$SCRIPT" ...` so the `PATH` override is explicit and reaches both the script and the stub. The fake `gh` handles three shapes of arguments:

- `repo view ...` -> prints `mirurobotics/docs`.
- `api --paginate repos/.../code-scanning/analyses?...` -> prints a canned JSON array of four analyses: id 300 category `.github/workflows/codeql-analysis.yml:codeql` created `2026-09-17T07:31:00Z` sha `b697e10...` deletable true; id 200 category `.github/workflows/codeql-analysis.yml:codeql/language:javascript` created `2026-03-23T20:56:00Z` deletable true; id 100 category `.github/workflows/codeql-analysis.yml:analyze/language:javascript` created `2026-03-23T03:00:00Z` deletable true; id 50 category `.github/workflows/codeql-analysis.yml:analyze/language:javascript` created `2026-03-22T03:00:00Z` deletable true (older analysis of the same category, so `max_by` must pick id 100). Every object carries `"tool": {"name": "CodeQL"}`.
- `api -X DELETE <url> ...` -> appends `<url>` to `$GH_STUB_LOG`; if `<url>` contains `/analyses/100?` it prints `{"next_analysis_url":"https://api.github.com/repos/mirurobotics/docs/code-scanning/analyses/50","confirm_delete_url":"https://api.github.com/repos/mirurobotics/docs/code-scanning/analyses/50?confirm_delete"}`, otherwise `{"next_analysis_url":null,"confirm_delete_url":null}`. The stub honours `--jq <expr>` by piping its output through the real `jq` with that expression.

Tests (four):

1. `--help prints usage and exits 0`: `run bash "$SCRIPT" --help`; `$status` 0; output contains `Usage:`.
2. `missing gh fails with a clear message`: `run env PATH="$BATS_TEST_TMPDIR/empty" "$BASH" "$SCRIPT" --repo o/r` (an empty dir created in `setup`); `$status` 1; output contains `gh is required`.
3. `dry run lists categories and deletes nothing`: `run env PATH="$STUB_BIN:$PATH" bash "$SCRIPT" --repo mirurobotics/docs`; `$status` 0; output contains `KEEP` + the `:codeql` category, `PRUNE` + `:codeql/language:javascript`, `PRUNE` + `:analyze/language:javascript`, `id=100`, and `dry run: 2 categories would be pruned`; `[ ! -s "$GH_STUB_LOG" ]`.
4. `--delete follows confirm_delete_url`: same invocation with `--delete`; `$status` 0; output contains `deleted`; `$GH_STUB_LOG` has exactly 3 lines (`[ "$(wc -l < "$GH_STUB_LOG")" -eq 3 ]`) containing `analyses/200?confirm_delete`, `analyses/100?confirm_delete`, `analyses/50?confirm_delete`, and `analyses/300` appears nowhere in it.

**Edit `.github/workflows/ci.yml`**, job `shell-tests`, last step: change the run line to `bats pub/scripts/agent/check-miru-access_test.bats scripts/codeql-prune_test.bats`.

**Edit `scripts/preflight.sh`**, last line: the same two-file `bats` command.

**Edit `.github/workflows/codeql-analysis.yml`**: insert this comment between `jobs:` and `  codeql:` (two-space indent to match the file):

      # The job id below is part of the code-scanning category
      # (".github/workflows/codeql-analysis.yml:codeql"). Renaming it or adding a
      # matrix orphans the old category as a stale "results may be out of date"
      # configuration under Security > Code scanning > Tools. After such a change,
      # update the default keep category in scripts/codeql-prune.sh and run it to
      # prune the old one.

No other files change. Do not touch `cspell.json`, `README.md`, or the CodeQL config.


## Concrete Steps

### Milestone 1 — The script

From `/home/user/docs`:

    git branch --show-current          # claude/upbeat-davinci-tyzi5x
    git status --short                 # clean apart from this plan file

Create `scripts/codeql-prune.sh` as described in Plan of Work, then:

    chmod +x scripts/codeql-prune.sh
    bash -n scripts/codeql-prune.sh
    ./scripts/codeql-prune.sh --help
    ./scripts/codeql-prune.sh --bogus; echo "exit=$?"
    ./scripts/codeql-prune.sh --repo mirurobotics/docs; echo "exit=$?"

Expect: `bash -n` silent; `--help` prints the usage text and exits 0; `--bogus` prints `unknown option: --bogus` plus the usage to stderr and `exit=2`; the last command prints `gh is required by codeql-prune.sh` and `exit=1` (no `gh` here). If `shellcheck` happens to be available, `shellcheck scripts/codeql-prune.sh` should be clean; it is optional.

Commit:

    git add scripts/codeql-prune.sh
    git commit -m "feat(scripts): add codeql-prune.sh to remove stale code-scanning configurations"

### Milestone 2 — The bats test and wiring

Create `scripts/codeql-prune_test.bats`, edit `ci.yml` and `preflight.sh` as described. `bats` is not installed here; try `sudo apt-get install -y bats` (or `apt-get install -y bats` as root). If it installs, from `/home/user/docs`:

    bats scripts/codeql-prune_test.bats

Expect `4 tests, 0 failures` (bats prints `ok 1 ...` through `ok 4 ...`). If bats cannot be installed, exercise the stub by hand so the logic is verified before CI runs it:

    stub=$(mktemp -d); mkdir "$stub/bin" "$stub/empty"
    # write the same fake gh from the test file into $stub/bin/gh and chmod +x it
    GH_STUB_LOG="$stub/gh.log" PATH="$stub/bin:$PATH" ./scripts/codeql-prune.sh --repo mirurobotics/docs
    GH_STUB_LOG="$stub/gh.log" PATH="$stub/bin:$PATH" ./scripts/codeql-prune.sh --repo mirurobotics/docs --delete
    cat "$stub/gh.log"

Expect the dry run to print three category lines (one `KEEP`, two `PRUNE`) and `dry run: 2 categories would be pruned; re-run with --delete`, and the log to contain three DELETE URLs (200, 100, 50) after the second run. Record in Surprises & Discoveries whether bats ran locally.

Also confirm the wiring edits:

    git diff -- .github/workflows/ci.yml scripts/preflight.sh

Expect exactly one changed line in each file (the `bats` command now names both test files).

Commit:

    git add scripts/codeql-prune_test.bats .github/workflows/ci.yml scripts/preflight.sh
    git commit -m "test(scripts): cover codeql-prune.sh with bats and run it in CI"

### Milestone 3 — Workflow comment

Add the comment to `.github/workflows/codeql-analysis.yml`, then `git diff -- .github/workflows/codeql-analysis.yml` should show only added comment lines (six lines, no other change). Commit:

    git add .github/workflows/codeql-analysis.yml
    git commit -m "docs(ci): note that the codeql job id is the code-scanning category"

### Milestone 4 — Publish and preflight

From `/home/user/docs`:

    git fetch origin main
    git rebase origin/main
    git push -u origin HEAD
    git rev-parse HEAD

Use `--force-with-lease` only if the rebase rewrote already-pushed commits; never plain `--force`. Open a draft PR with `mcp__github__create_pull_request` (owner `mirurobotics`, repo `docs`, base `main`, head `claude/upbeat-davinci-tyzi5x`, `draft: true`), title `feat(scripts): add codeql-prune.sh to remove stale code-scanning configurations`, body along these lines (unindented when pasted):

    ## Summary
    - Security > Code scanning > Tools > CodeQL lists two stale `language:javascript` configurations left behind when #13 renamed the job `analyze` -> `codeql` and #14 removed the matrix. GitHub never garbage-collects old categories; the only fix is deleting their analyses.
    - Add `scripts/codeql-prune.sh` (dry run by default, `--delete` to prune, `--ref` for the environment branches) plus a bats test wired into the `shell-tests` job and `scripts/preflight.sh`.
    - Comment the job id in `codeql-analysis.yml` so the next rename runs the script.

    ## After merge (needs a token with code scanning write access)
    - `./scripts/codeql-prune.sh` (dry run) then `./scripts/codeql-prune.sh --delete`
    - Repeat with `--ref refs/heads/staging --ref refs/heads/uat --ref refs/heads/production` if those branches list stale categories.
    - Confirm the Tools page shows only `.github/workflows/codeql-analysis.yml:codeql`.

    ## Test plan
    - [ ] `bats scripts/codeql-prune_test.bats` passes (CI `shell-tests`)
    - [ ] CI green on the pushed head (preflight CLEAN)

End the body with the attribution lines required for this session (`🤖 Generated with [Claude Code](https://claude.com/claude-code)` and the session URL).

Watch CI for the pushed head SHA (see Context: MCP tools or `curl` against `/commits/<sha>/check-runs`). Expect `lint`, `audit`, and `shell-tests` to succeed; the `tools/lint` jobs are skipped. If `shell-tests` is red, read its log with `mcp__github__get_job_logs` (`failed_only: true`): a `jq: command not found` inside the stub means the runner image lacks `jq`, in which case change the install step to `sudo apt-get install -y bats jq`; any other failure is a script or test bug to fix. Fix in a new commit (never amend a pushed commit), push, re-watch.

Preflight must report `CLEAN` (every check green on the pushed branch head) before the PR leaves draft and before this task is reported complete. Only then mark the PR ready for review (`mcp__github__update_pull_request` with `draft: false`).


## Validation and Acceptance

1. `./scripts/codeql-prune.sh --help` exits 0 and prints the usage text; `--bogus` and `--keep` without a value exit 2 with a message naming the problem.
2. With no `gh` on `PATH` the script exits 1 printing `gh is required by codeql-prune.sh` and makes no other external call first (this is what test 2 checks).
3. Against the stub `gh`, the dry run prints one `KEEP` line for `.github/workflows/codeql-analysis.yml:codeql` and `PRUNE` lines for the two `language:javascript` categories, each with the newest analysis's date, short sha, and id (id 100, not 50, for the `analyze` category), ends with `dry run: 2 categories would be pruned; re-run with --delete`, and issues no DELETE.
4. Against the stub `gh`, `--delete` issues exactly three DELETE calls (ids 200, 100, 50, in that order per category) and never touches id 300.
5. `bats scripts/codeql-prune_test.bats` reports `4 tests, 0 failures`; `bats pub/scripts/agent/check-miru-access_test.bats scripts/codeql-prune_test.bats` (the CI command) passes as a whole.
6. `git diff main --stat` shows only `scripts/codeql-prune.sh`, `scripts/codeql-prune_test.bats`, `scripts/preflight.sh`, `.github/workflows/ci.yml`, `.github/workflows/codeql-analysis.yml`, and this plan file; the workflow diff is comment-only, so the live category stays `.github/workflows/codeql-analysis.yml:codeql`.
7. **Preflight reports `CLEAN`:** `lint`, `audit`, and `shell-tests` are green on the current pushed head SHA of `claude/upbeat-davinci-tyzi5x`. Criteria 1-6 and this must hold before the PR leaves draft and before the task is reported complete.
8. Post-merge, by a maintainer with a code-scanning-write token (not a PR gate): `./scripts/codeql-prune.sh --delete` prints `deleted` for the two stale categories, and `https://github.com/mirurobotics/docs/security/code-scanning/tools` lists a single CodeQL configuration with no "results may be out of date" warning.


## Idempotence and Recovery

All file edits are plain additions and can be reapplied or reverted with `git checkout -- <file>` / `git revert <commit>`; each milestone is one commit so a revert is clean. The script is safe to rerun: the dry run is read-only, and `--delete` on an already-pruned category is a no-op because the category no longer appears in the listing. If a `--delete` run is interrupted mid-chain, rerunning it resumes from the newest remaining analysis of that category (the API only ever lets you delete the newest one, so partial chains cannot get stuck); a DELETE that fails with 403 means the token lacks code-scanning write access and nothing has been changed. The bats stub writes only under `$BATS_TEST_TMPDIR`. The only permitted force-push is `--force-with-lease` after rebasing onto `origin/main`; never amend a pushed commit.

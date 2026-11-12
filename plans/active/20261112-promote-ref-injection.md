# Stop shell injection through the Promote workflow's `ref` input

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Pass the Promote workflow's inputs to `scripts/promote.sh` through environment variables, and make the script reject a ref that is not a plain tag, branch or commit SHA name. |

This plan lives in `docs/plans/` because both files it changes are in the docs repo. Work on the existing branch `fix/eng-1866-promote-ref-injection` (do not create another); the PR targets `main`. Tracking: Linear ENG-1866, security finding SEC-1111-3.

## Purpose / Big Picture

Today anyone who can start the Promote workflow can run arbitrary shell commands with the Waterfall Promoter GitHub App token, and so push any commit to the `production` docs branch, skipping the script's "reachable from main" check. After this change the `ref` input is only ever data: GitHub never pastes it into a shell script, and `scripts/promote.sh` refuses any ref that is not made of letters, digits, `.`, `_`, `/` and `-` starting with a letter or digit. Promoting `main`, a tag such as `v1.2.3`, or a 40-character commit SHA works exactly as before.

## Decisions

### Names

| New name | Alternatives | Why this one |
|----------|--------------|--------------|
| `REF` (step env var in `promote.yml`) | `PROMOTE_REF`, `INPUT_REF` | Short UPPER_SNAKE name like `BRANCH` in `preview.yml`; `INPUT_*` is the prefix the Actions runner uses for action inputs, so it would read as something it is not. |
| `ENVIRONMENT` (step env var in `promote.yml`) | `TARGET_ENVIRONMENT`, `PROMOTE_ENV` | Matches the input name `environment` and the script's `environment` variable. |

### Left out

- Validating `environment` in `scripts/promote.sh`: it is a `choice` input limited to `staging`, `uat` and `production`, and after this change it also reaches the shell only through an env var.
- Changing `run-name`, the job `name:`, the step `name:` or the `notify` job's `with:` block: GitHub renders these as text or passes them as action inputs; none of them is shell.
- Refs with `~`, `^` or `@` (for example `HEAD~1`): the input's description asks for a tag or commit SHA, and none of the documented uses needs revision syntax.
- A test harness for `promote.sh` or the workflow, and adding shellcheck or actionlint to CI: the repo has neither today, and the fix is a few lines.

### Open questions

- Should the new error message echo the rejected ref? Default: no. A crafted ref could contain a newline that the runner would read as a new workflow command, so the message names the allowed characters instead.

## Progress

- [x] Milestone 1: edit `.github/workflows/promote.yml` and `scripts/promote.sh`, then commit.
- [x] Source written
- [ ] Push, open the draft PR, and run `$preflight` until it reports `CLEAN`.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

(Add entries as you go.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

`.github/workflows/promote.yml` is a manually started (`workflow_dispatch`) workflow with two inputs: `ref`, a free-text string, and `environment`, a choice of `staging`, `uat` or `production`. Its `promote` job creates a token for the Waterfall Promoter GitHub App, checks out `main` with that token (the checkout step leaves the token in the local git config, so later `git push` calls use it), and then runs:

    - name: Promote ${{ inputs.ref }} → ${{ inputs.environment }}
      id: promote
      run: ./scripts/promote.sh "${{ inputs.ref }}" "${{ inputs.environment }}"

GitHub replaces each `${{ ... }}` expression with its text before bash sees the script. A ref such as `main"; git push ...; echo "` closes the quote and runs its own commands. The Waterfall Promoter app is the bypass actor on the ruleset that protects the environment branches, so those commands can push anything to `production`.

The safe pattern, already used in `.github/workflows/preview.yml`, is to map each expression to an environment variable in the step's `env:` block and refer to the variable, in double quotes, in `run:`. Bash then expands the variable as data and never parses its contents as code.

`scripts/promote.sh` takes `<ref> <environment>`, runs `git fetch origin --tags --force`, resolves the ref with `git rev-parse --verify "${ref}^{commit}"`, checks with `git merge-base --is-ancestor` that the commit is on `origin/main`, and force-pushes it to the environment branch. It reports failures with `echo "::error::..."` followed by `exit 1`. A ref beginning with `-` would be read by git as an option, which is why the new check requires a letter or digit first.

CI (`.github/workflows/ci.yml`) runs on every `pull_request`; its `lint` and `audit` jobs run on this PR, and the `lint-custom-linter` and `test-custom-linter` jobs are skipped because nothing under `tools/lint/` changes. No CI job runs shellcheck or actionlint, and nothing tests `promote.sh` or `promote.yml`.

## Plan of Work

In `.github/workflows/promote.yml`, in the `promote` job's step with `id: promote`, add an `env:` block between `id:` and `run:` that sets `REF: ${{ inputs.ref }}` and `ENVIRONMENT: ${{ inputs.environment }}`, and change `run:` to `./scripts/promote.sh "${REF}" "${ENVIRONMENT}"`. Leave every other line of the file as it is.

In `scripts/promote.sh`, at the start of section 1 ("Resolve ref and verify it is reachable from main"), before `git fetch origin --tags --force`, add:

    if [[ ! "${ref}" =~ ^[A-Za-z0-9][A-Za-z0-9._/-]*$ ]]; then
      echo "::error::Ref must be a tag, branch or commit SHA: letters, digits, '.', '_', '/' or '-', starting with a letter or digit."
      exit 1
    fi

The check runs before the fetch so a bad ref fails without touching the network, and before `git rev-parse` so git never sees it. No comment is needed: the error message states the rule. Add no other helpers, files or comments.

## Concrete Steps

From the docs repo root (`/home/user/docs`), confirm the branch:

    git branch --show-current
    fix/eng-1866-promote-ref-injection

Make the two edits described in Plan of Work. Then, as a one-off check of the pattern (not part of any test suite), from any directory:

    for r in main v1.2.3 0123456789abcdef0123456789abcdef01234567 'main"; echo pwned; echo "' --output=x ''; do
      if [[ "$r" =~ ^[A-Za-z0-9][A-Za-z0-9._/-]*$ ]]; then echo "accept: $r"; else echo "reject: $r"; fi
    done

Expected output:

    accept: main
    accept: v1.2.3
    accept: 0123456789abcdef0123456789abcdef01234567
    reject: main"; echo pwned; echo "
    reject: --output=x
    reject: 

Review the diff by eye: `promote.yml` must contain no `${{ inputs.` inside any `run:`, and the `env:` keys must match the names used in `run:`.

    git diff -- .github/workflows/promote.yml scripts/promote.sh

Commit milestone 1 (from the docs repo root):

    git add .github/workflows/promote.yml scripts/promote.sh
    git commit -m "fix(ci): pass promote inputs through env and validate ref"

Then push the branch and open a draft PR with `$pr`, and run `$preflight`.

## Validation and Acceptance

Tests: the repo has no tests for `scripts/promote.sh` or the workflows, and CI runs no shellcheck or actionlint, so this plan adds no test and no harness. Validation is CI green on the PR plus a manual review of the diff and the one-off pattern check above.

CI: the `ci.yml` jobs `lint` and `audit` must pass on the pushed branch head. They do not exercise the Promote workflow; a green run shows the change breaks nothing else. Validation is complete only when `$preflight` reports `CLEAN` (CI green on the pushed branch head); the PR must not leave draft, and the task must not be reported complete, before that.

Acceptance, after merge, for a maintainer (the agent does not start Promote runs, because they deploy): running Promote with `ref` = `main` and `environment` = `staging` succeeds as before and the step summary shows `main` resolved to a SHA. Running it with `ref` = `main"; echo pwned; echo "` fails in the "Promote" step with the error `Ref must be a tag, branch or commit SHA: ...`, and the log never prints `pwned`.

## Idempotence and Recovery

The edits are plain text changes and can be reapplied safely; re-running the one-off check has no side effects. If CI fails, fix forward on the same branch and push again. To undo after merge, revert the commit; nothing else depends on it.

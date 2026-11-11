# Run docs CI on GitHub-hosted `ubuntu-latest` runners

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (mirurobotics/docs) | read-write | Change the runner label `blacksmith-2vcpu-ubuntu-2404` to `ubuntu-latest` in the five workflow files under `.github/workflows/`. No other files change. |
| `cli/`, `ansible-collection-agent/` | read-only | Precedent: both already run their jobs, including the same pinned reusable AI review workflow, on `ubuntu-latest`. |

This plan lives in `docs/plans/` because the only edits are in the docs repo. Work on the checked-out branch `claude/youthful-lamport-yn0sf2` (base: `main`).

Out of scope: the same move in the other public repos (agent, gotools, setup-cli), which follow this exemplar in their own PRs, and the `blacksmith.sh` URLs inside the vendored API specs under `docs/references/*.yaml`, which are unrelated to CI.

## Purpose / Big Picture

mirurobotics/docs is a public repository. GitHub-hosted standard runners are free for public repositories (Linux: 4 vCPU, 16 GB) and do not draw from the GitHub Enterprise Cloud included compute allowance ($300/month), while every Blacksmith minute spent here is billed. After this change every job in the docs workflows runs on GitHub's `ubuntu-latest` runner, so docs CI costs nothing. The jobs keep their names and steps, so branch protection's required checks are unchanged. A reader sees it working when the job logs on the PR show "Runner name: GitHub Actions …" and an image of `ubuntu-24.04` instead of a Blacksmith runner, and all checks are green.

This PR is the exemplar for the same change in agent, gotools and setup-cli.

## Decisions

### Precedent

| New file | Mirrors |
|----------|---------|
| `plans/backlog/20261111-github-hosted-runners.md` | `plans/completed/20261109-readable-upload-paths.md` and the other docs plans |

No source file is added. The edited `runs-on` lines mirror `cli/.github/workflows/ai-review.yml` (`with: runs-on: ubuntu-latest` passed to `mirurobotics/.github/.github/workflows/ai-review.yml@bf0916104b6b5b5373c09a5f5b413ec6bded2cb5`, the same SHA docs pins), `cli/.github/workflows/codeql-analysis.yml` and `ansible-collection-agent/.github/workflows/ci.yml`.

### Left out

- The agent, gotools and setup-cli workflows: proposed follow-up PRs that cite this one, at most three open at a time.
- An `actions/setup-go` step in the `ci.yml` `lint` job (which builds the Go prose linter with the image's preinstalled `go`): GitHub's `ubuntu-24.04` image ships Go on `PATH`, like the Blacksmith image it replaces. Add it only if CI shows `go is required for MDX prose linting`.
- `.github/actionlint.yml` or any runner-label config: docs has none, and no docs workflow runs actionlint, so a standard label needs no registration.
- Larger or ARM runners, timeout or concurrency changes, and bumping the pinned reusable-workflow SHA: not needed to move off Blacksmith.

### Open questions

- `ubuntu-latest` or the pinned `ubuntu-24.04`? Default: `ubuntu-latest`, as the goal states and as cli and ansible-collection-agent use. `ubuntu-latest` resolves to `ubuntu-24.04` today, the same release as `blacksmith-2vcpu-ubuntu-2404`.

## Progress

- [x] Replace the 11 `blacksmith-2vcpu-ubuntu-2404` labels with `ubuntu-latest`.
- [x] Commit `ci: run docs workflows on GitHub-hosted ubuntu-latest runners`.
- [x] Source written
- [ ] Push and open the draft PR; CI, Preview and a dispatched CodeQL run pass on the branch head.
- [ ] Preflight reports `CLEAN`.
- [ ] After merge: AI Review on the next PR, the `ci.yml` push run on `main`, and the next Promote run pass on `ubuntu-latest`.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

(Add entries as work proceeds.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

A runner is the virtual machine a GitHub Actions job runs on, chosen by the job's `runs-on` key. `blacksmith-2vcpu-ubuntu-2404` is a label served by Blacksmith, a paid third-party runner provider; `ubuntu-latest` is GitHub's own hosted Ubuntu runner. Blacksmith's Ubuntu 24.04 image mirrors GitHub's `ubuntu-24.04` image, so the same preinstalled tools (`git`, `curl`, `jq`, `gh`, Go, Node, corepack) are present.

The label appears 11 times, all in `.github/workflows/` (line numbers as of `main` at f9152a7):

- `ai-review.yml:30`, job `ai-review`. This is not a `runs-on` key but a `with: runs-on:` input to the reusable workflow `mirurobotics/.github/.github/workflows/ai-review.yml`, which uses it as its job's runner.
- `ci.yml:19` (`changes`), `:35` (`lint`), `:64` (`audit`), `:90` (`lint-custom-linter`), `:110` (`test-custom-linter`).
- `codeql-analysis.yml:23` (`changes`), `:39` (`codeql`).
- `preview.yml:18` (`preview`).
- `promote.yml:27` (`promote`), `:57` (`notify`).

Reading every workflow found no Blacksmith-only behavior: no `useblacksmith/*` actions, no Blacksmith cache or sticky disks, and no Docker layer caching. Caching uses `actions/setup-node` (`cache: pnpm`) and `actions/setup-go`, which use the standard GitHub Actions cache on any runner. The other actions are `actions/checkout`, `dorny/paths-filter`, `github/codeql-action`, `actions/create-github-app-token` and the composite actions `mirurobotics/.github/actions/codeql-gate` and `notify-deploys`, none tied to a runner provider. No `actionlint.yml` exists in the repo.

When each workflow runs, which decides what the PR itself can prove:

- `ci.yml`: on `pull_request` and on push to `main`, `staging`, `uat`, `production`. On a PR, `lint-custom-linter` and `test-custom-linter` are skipped unless `tools/lint/**` changed, so they first run on the new runner on the push to `main` after merge.
- `preview.yml`: on `pull_request` types `opened` and `reopened` only, for same-repo, non-Dependabot PRs. It uses the PR's version of the file, so it runs on the new runner when the PR opens after the change is pushed.
- `ai-review.yml`: on `pull_request_target`, which runs the file from the base branch. The PR's own AI review still runs on Blacksmith; the new runner takes effect on the first PR after merge.
- `codeql-analysis.yml`: on push to `staging`/`uat`/`production`, nightly, and `workflow_dispatch`. It does not run on PRs, but can be dispatched against the branch.
- `promote.yml`: `workflow_dispatch` only, and it promotes a ref to a deployment environment, so it must not be dispatched from the branch. It is verified on the next real promotion after merge.

## Plan of Work

In each of the five files under `.github/workflows/`, replace `blacksmith-2vcpu-ubuntu-2404` with `ubuntu-latest` on the lines listed above, leaving indentation and every other line untouched. One commit, one milestone.

## Concrete Steps

All commands run from the docs repo root, `/home/user/docs`.

1. Confirm the branch and the starting count:

       git branch --show-current          # expect: claude/youthful-lamport-yn0sf2
       grep -rc 'blacksmith-2vcpu-ubuntu-2404' .github/workflows/
       # expect: ai-review.yml:1 ci.yml:5 codeql-analysis.yml:2 preview.yml:1 promote.yml:2

   If every count is 0, the change is already applied; skip to step 4.

2. Replace the label:

       sed -i 's/blacksmith-2vcpu-ubuntu-2404/ubuntu-latest/' .github/workflows/*.yml

3. Check the result (read-only inspection, not validation tooling):

       grep -rn 'blacksmith' .github/          # expect: no output
       grep -rc 'runs-on: ubuntu-latest' .github/workflows/
       # expect: ai-review.yml:1 ci.yml:5 codeql-analysis.yml:2 preview.yml:1 promote.yml:2
       git diff --stat                         # expect: 5 files, 11 insertions, 11 deletions

4. Commit the milestone:

       git add .github/workflows plans/
       git commit -m "ci: run docs workflows on GitHub-hosted ubuntu-latest runners"

5. Push and open the draft PR (via `$pr`). In the PR body, name the follow-ups for agent, gotools and setup-cli, and state which workflows can only be verified after merge.

6. Dispatch CodeQL against the branch, which runs CI on GitHub, not locally:

       gh workflow run codeql-analysis.yml --ref claude/youthful-lamport-yn0sf2
       gh run list --workflow codeql-analysis.yml --branch claude/youthful-lamport-yn0sf2 --limit 1

7. Run `$preflight` until it reports `CLEAN`. Fix any failure from `gh run view --log-failed`.

No local actionlint, `pnpm` scripts, builds or tests are run.

## Validation and Acceptance

The repo has no test suite for workflow configuration, and this change adds no test. The test is that every touched workflow still runs and passes on GitHub's runner:

- `ci.yml` on the PR: `changes`, `lint` and `audit` pass, and each job log's "Set up job" section shows a GitHub-hosted runner with image `ubuntu-24.04`, not Blacksmith. `lint-custom-linter` and `test-custom-linter` are skipped on the PR because `tools/lint/**` is unchanged.
- `preview.yml`: the `preview` job runs on PR open and posts a "Mintlify preview:" comment. If the PR was opened before the change was pushed, close and reopen it to rerun Preview on the new runner.
- `codeql-analysis.yml`: the dispatched run from step 6 passes `changes` and runs `codeql` on `ubuntu-latest`. If the gate reports nothing to scan and skips `codeql`, note it in Surprises and verify on the next nightly run after merge.
- Gate: `$preflight` must report `CLEAN`, meaning CI is green on the pushed branch head, before the PR leaves draft.

Verified only after merge, and recorded in Outcomes:

- `ai-review.yml`: the AI review on the next PR after merge runs on `ubuntu-latest` and posts its `ai-review` status. The cli repo already passes `runs-on: ubuntu-latest` to the same reusable workflow at the same SHA.
- `ci.yml` push to `main`: all five jobs, including `lint-custom-linter` and `test-custom-linter`, pass.
- `promote.yml`: the next promotion runs `promote` and `notify` on `ubuntu-latest` and posts to #deploys.

## Idempotence and Recovery

Step 2 is safe to rerun; once the label is gone the `sed` matches nothing. Before committing, `git checkout -- .github/workflows` restores the files. After merge, a failure that only `ubuntu-latest` causes is fixed forward in a new PR (for example, adding `actions/setup-go` to the `lint` job), or reverted with `git revert <sha>` in a PR. Never force-push `main`.

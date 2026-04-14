# Remove pull_request trigger from CodeQL workflow

This ExecPlan is a living document.

## Scope
| Repository | Access | Description |
|-----------|--------|-------------|
| `docs` | read-write | Edit CodeQL workflow triggers |

This plan lives in `docs/plans/` because the workflow file is in this repo.

## Purpose / Big Picture
CodeQL should only run on push to main — not on PRs — to reduce CI noise and align with the team-wide migration.

## Progress
- [ ] Edit `.github/workflows/codeql-analysis.yml` to remove `pull_request:` trigger
- [ ] Run preflight

## Surprises & Discoveries
(Add entries as work proceeds.)

## Decision Log
(Add entries as work proceeds.)

## Outcomes & Retrospective
(Fill at completion.)

## Context and Orientation
The file `.github/workflows/codeql-analysis.yml` currently triggers on `pull_request` (all branches), `push` to main/staging/uat/production, `schedule`, and `workflow_call`. The `pull_request` trigger needs removal.

## Plan of Work
Edit `.github/workflows/codeql-analysis.yml`: remove the `pull_request:` line from the `on:` block. Keep push, schedule, and workflow_call.

## Concrete Steps
1. Edit `.github/workflows/codeql-analysis.yml` — remove the `pull_request:` line.
2. Commit: `ci: remove pull_request trigger from codeql`
3. Run `./scripts/preflight.sh`

## Validation and Acceptance
After the edit, the `on:` block should have only `push`, `schedule`, and `workflow_call` — no `pull_request`. Preflight passes clean.

## Idempotence and Recovery
Safe to re-run; the edit is idempotent. If something goes wrong, `git checkout main` restores the original.

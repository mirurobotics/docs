# Prepare Public Agent Changelog for the Next Docs Release

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `./` (docs repo) | read-write | Owns this ExecPlan and the changelog file edit. |
| `../agent/` | read-only | Source of truth for release tags, commits, and user-visible behavior changes (sibling repo in meta workspace). |

This plan lives in `docs/.ai/exec-plans/` because `docs/` is the only repository being edited. `agent/` is read-only research context.

## Purpose / Big Picture

After this work, the public docs changelog for the Miru Agent clearly describes what external users get in the next agent release, including feature additions, behavior changes, and fixes that matter for operators and integrators. A reader should be able to open the Agent changelog page and understand what changed without reading raw commit history.

The changelog entry must be accurate to the actual release cut (tag or explicit commit), avoid internal-only refactor noise, and include migration context if any workflow changed.

## Progress

- [x] (2026-03-05 18:43Z) Loaded ExecPlan policy and meta-repo rules.
- [x] (2026-03-05 18:43Z) Located changelog target file and current placeholder in `docs/changelogs/agent.mdx`.
- [x] (2026-03-05 18:43Z) Verified agent release tags and current branch state to shape research workflow.
- [x] (2026-03-05 18:49Z) Moved ExecPlan from meta backlog to `docs/.ai/exec-plans/backlog/` to match write ownership.
- [x] (2026-03-05 18:52Z) Set `TO_REF=origin/v0.7` (`ba79556`) because the docs site already targets unreleased `v0.7.0` and the agent repo does not yet have a final `v0.7.0` tag.
- [x] (2026-03-05 18:53Z) Built a candidate change list from `../agent` commit history in the `v0.6.1..origin/v0.7` window.
- [x] (2026-03-05 18:54Z) Triaged candidate commits into `include` vs `exclude` based on external user impact.
- [x] (2026-03-05 18:54Z) Drafted public-facing changelog copy in `docs/changelogs/agent.mdx` (summary + categorized bullets).
- [x] (2026-03-05 18:55Z) Validated placeholder removal and link/path references manually; `mintlify dev` could not bind a local port and `mintlify broken-links` was blocked by an unrelated existing MDX parse error in another file.
- [x] (2026-03-05 18:55Z) Prepared source-commit mapping for each published changelog bullet.

## Surprises & Discoveries

- Observation: The current docs entry for `# v0.7.0` is still placeholder text (`TBD` and TODO note), so this release requires a full write-up, not a minor edit.
  Evidence: `docs/changelogs/agent.mdx`.

- Observation: `agent/` currently has many local uncommitted changes on `v0.7`, so changelog research must rely on immutable refs (tags/SHAs), not working-tree diff.
  Evidence: `git -C ../agent status --short` output includes many modified files and deletions.

- Observation: The release window from `v0.6.1` to `v0.7.0-beta.1` includes many `chore`, `test`, and `refactor` commits, so a strict inclusion rubric is necessary to avoid publishing internal-only noise.
  Evidence: commit-type frequency from `git -C ../agent log ... | cut -d: -f1 | sort | uniq -c`.

- Observation: `mintlify broken-links` in `docs/` is currently blocked by an unrelated MDX parse error in `docs/.ai/exec-plans/active/20260304-changelog-visual-redesign.md`, even though `docs/.mintignore` lists `.ai/`.
  Evidence: `mintlify broken-links` reports `Unexpected character '7' ...` at line 262 of that file.

## Decision Log

- Decision: Use an explicit release boundary (`FROM_REF=v0.6.1`, `TO_REF=<confirmed release tag or SHA>`) before writing copy.
  Rationale: Prevents mixing unreleased work into a public changelog and keeps claims tied to what ships.
  Date/Author: 2026-03-05 / Codex

- Decision: Include only user-observable behavior changes in the published entry; keep internal refactors/tests/coverage updates out unless they directly change external reliability or upgrade behavior.
  Rationale: Public changelog should optimize for customer comprehension, not internal implementation details.
  Date/Author: 2026-03-05 / Codex

- Decision: Keep proof mapping from each published bullet to at least one commit SHA in working notes.
  Rationale: Provides auditability for review and reduces risk of inaccurate changelog statements.
  Date/Author: 2026-03-05 / Codex

- Decision: Treat the release as `v0.7.0` with date label `Unreleased`, and use `origin/v0.7` as the source boundary.
  Rationale: The docs repo already references `v0.7.0` as the compatibility target for Device API `v0.2.0`, but the agent repo has not yet published a final `v0.7.0` tag.
  Date/Author: 2026-03-05 / Codex

## Outcomes & Retrospective

Completed the `v0.7.0` agent changelog draft in `docs/changelogs/agent.mdx`. The placeholder entry now describes the Device API `v0.2.0` rollout, the main migration-facing breaking changes, and the reliability/security fixes that external users would notice.

Published bullet mapping:

- Device API `v0.2` routes and richer `GET /version` metadata: `85d7f3d`
- Deployment metadata adds `deployed_at` and `archived_at`: `2ef6513`
- Automatic retries for transient network connection errors: `a541a4b`
- Cooldown reporting and clearing fixes: `9daf41d`, `7c9f609`, `fcaf2b8`
- Deployment staging moved to `/srv/miru/.temp`: `6f8da9d`, `6b62e12`
- Token values redacted from debug output: `33377b6`

Residual issue: full Mintlify repo validation is still blocked by a pre-existing MDX parse error in `docs/.ai/exec-plans/active/20260304-changelog-visual-redesign.md`. The agent changelog file itself passed placeholder and link-target sanity checks.

## Context and Orientation

Miru publishes component-level changelogs in this repository under `docs/changelogs/`. The target file for this task is `docs/changelogs/agent.mdx`.

Current state:

- The newest complete agent changelog entry is `v0.6.1` (January 22, 2026).
- The `v0.7.0` section exists but is a placeholder and currently references Device API v0.2.0 with TODO text.
- Agent release tags currently include `v0.7.0-alpha.1`, `v0.7.0-beta.1`, and older stable tags.

Definitions for this plan:

- **Release window**: commits included in the upcoming changelog, represented by `FROM_REF..TO_REF`.
- **External user impact**: any change that modifies installation, runtime behavior, deployment/sync semantics, API compatibility, logging/error behavior visible to operators, or upgrade/migration requirements.
- **Internal-only change**: test-only, refactor-only, CI-only, coverage-only, or docs-only changes with no externally observable behavior shift.

Primary files to touch:

- Research source: `../agent/` git history and affected Rust modules.
- Publication target: `docs/changelogs/agent.mdx`.

## Plan of Work

First, lock the release boundary so research is deterministic. Use the last documented version (`v0.6.1`) as the lower bound and confirm the upper bound (`TO_REF`) with release ownership before drafting any prose.

Second, extract the commit set for that window and classify each candidate by user impact. Prioritize commits scoped as `feat` and `fix`, then inspect `refactor` or `build/openapi` commits only when they clearly change external behavior (for example API shape, deployment handling, or install/runtime reliability). Record accepted bullets with supporting SHAs.

Third, draft the `v0.7.0` section in `docs/changelogs/agent.mdx` using the existing docs style: release header, date line, short summary paragraph, and categorized sections such as `Features`, `Improvements`, `Fixes`, and `Breaking changes` when needed.

Fourth, run lightweight docs validation and manual quality checks: remove placeholder text, verify internal links, and ensure wording is public-facing and non-speculative. Finalize with a brief reviewer note listing source SHAs for each changelog bullet.

## Concrete Steps

From the docs repo root (`/home/ben/miru/miru/docs`), identify the current changelog state:

    rg '^# v|TBD|TODO' docs/changelogs/agent.mdx

Expected: shows existing version headers and confirms placeholder text that must be replaced.

From the docs repo root, inspect available agent release refs:

    git -C ../agent tag --sort=-creatordate | head -n 20

Expected: includes `v0.6.1`, `v0.7.0-alpha.1`, and `v0.7.0-beta.1` at minimum.

Set and validate commit range refs (replace `<TO_REF>` with the confirmed release tag/SHA):

    FROM_REF=v0.6.1
    TO_REF=<TO_REF>
    git -C ../agent rev-parse "$FROM_REF" "$TO_REF"

Expected: command prints two SHAs and exits successfully.

Generate a research log containing commit subjects and touched files:

    git -C ../agent log --no-merges --name-only --pretty='format:=== %h %s' "$FROM_REF..$TO_REF" > /tmp/agent-changelog-log.txt

Expected: `/tmp/agent-changelog-log.txt` contains commit blocks starting with `=== <sha> <subject>`.

Triage candidate user-facing commits:

    rg '^=== ' /tmp/agent-changelog-log.txt

Then inspect individual candidates in detail:

    git -C ../agent show --stat --patch <sha>

Expected: enough evidence to mark each candidate as include/exclude for public changelog.

Edit the docs changelog entry in `docs/changelogs/agent.mdx`:

    # update the # v0.7.0 block: replace TBD/TODO with finalized summary and categorized bullets

Check for unresolved placeholders after editing:

    rg 'TBD|TODO' docs/changelogs/agent.mdx

Expected: no matches for the new release block.

From the docs repo root, run local docs preview validation:

    mintlify dev

Expected: local preview starts successfully and `/docs/changelogs/agent` renders with the updated entry.

If `mintlify` is unavailable in PATH, use:

    npx mintlify dev

## Validation and Acceptance

Acceptance is complete when all of the following are true:

1. `docs/changelogs/agent.mdx` has a fully written release entry (no placeholders) for the target release.
2. Every published changelog bullet maps to at least one commit SHA in the confirmed `FROM_REF..TO_REF` range.
3. The entry excludes clearly internal-only changes unless they have direct operator-visible impact.
4. Local docs preview renders the updated changelog page without MDX or navigation errors.
5. The final text clearly communicates user-visible behavior, fixes, and migration/breaking implications (if any) in plain language.

## Idempotence and Recovery

Most research commands are read-only and safe to rerun with the same refs.

If the wrong `TO_REF` is used, update `TO_REF` and regenerate `/tmp/agent-changelog-log.txt`; this cleanly replaces the prior research artifact.

Editing `docs/changelogs/agent.mdx` is repeatable. If a draft becomes inconsistent, restore the file and reapply changes from the commit-backed notes:

    git restore -- docs/changelogs/agent.mdx

Before using restore, confirm no unrelated docs changes need to be kept:

    git status --short

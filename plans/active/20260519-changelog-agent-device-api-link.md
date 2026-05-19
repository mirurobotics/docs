# Fix dead `/developers/device-api` link in `docs/changelog/agent.mdx`

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit `docs/changelog/agent.mdx` to replace one dead internal link. |

This plan lives in `docs/plans/backlog/` because the only code change is in the docs repo (`docs/changelog/agent.mdx`).

Base branch for the eventual PR: `main`.

## Purpose / Big Picture

`docs/changelog/agent.mdx` line 26 currently links to `/developers/device-api`, but no such page exists in the docs site. The `developers/device-api/` directory contains only section children (`overview.mdx`, `sdks.mdx`, `versions.mdx`, `authn.mdx`, `events.mdx`); there is no `developers/device-api.mdx` or `developers/device-api.md`. The site's `docs.json` navigation lists those children explicitly with no bare `/developers/device-api` entry, and its `redirects` block has no rule that would catch `/developers/device-api`. The link therefore 404s.

Other changelog entries (notably `docs/changelog/product.mdx` "Device API Documentation »") already link directly to `/developers/device-api/overview`. This plan brings `agent.mdx` into line with that existing pattern by replacing the single dead link with `/developers/device-api/overview`.

After this change, the affected sentence still reads identically to the user — `[Device API](...)` — but the link target points to the existing section landing page instead of a 404.

## Progress

- [ ] Confirm the current state of line 26 in `docs/changelog/agent.mdx`.
- [ ] Replace `/developers/device-api` with `/developers/device-api/overview` on that line.
- [ ] Run `grep` assertions (see Concrete Steps) to verify the edit landed precisely.
- [ ] Run `pnpm lint`.
- [ ] Run `pnpm test:lint`.
- [ ] Run `./scripts/preflight.sh`.
- [ ] Commit the single-file edit.

Use timestamps when you complete steps. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add observations here as work proceeds. Include evidence — exact commands/output — and a date/author tag.)

## Decision Log

- Decision: Replace the dead link with `/developers/device-api/overview` rather than removing it or adding a redirect in `docs.json`.
  Rationale: The `overview.mdx` page is the natural landing page for the Device API section, and other changelog entries (e.g. `docs/changelog/product.mdx`) already use this exact target. Adding a redirect would expand scope and would mask similar broken links elsewhere; out of scope for this plan.
  Date/Author: (fill in on execution)

## Outcomes & Retrospective

(Fill in on completion: link to commit SHA, lint/preflight results, and any deferred follow-ups.)

## Context and Orientation

This repo is the Miru documentation site, hosted on Mintlify. Pages live under `docs/` as `.mdx` files. The site's navigation and redirect rules are defined in `docs.json`.

Files and concepts relevant to this change:

- `docs/changelog/agent.mdx` — the file being edited. Line 26 currently reads:

      - Resources fetched from on-device applications via the [Device API](/developers/device-api) are now fetched from the Miru control plane if not found locally

  After the edit, the same line will read:

      - Resources fetched from on-device applications via the [Device API](/developers/device-api/overview) are now fetched from the Miru control plane if not found locally

- `docs/developers/device-api/` — contains `overview.mdx`, `sdks.mdx`, `versions.mdx`, `authn.mdx`, `events.mdx`. There is no `docs/developers/device-api.mdx` or `.md`, which is why the bare path 404s.

- `docs.json` — site nav and redirects. Verified by the caller to have no entry for `/developers/device-api` and no redirect rule that would catch it.

- `docs/changelog/product.mdx` — uses the canonical target `/developers/device-api/overview` for the same section. This plan aligns `agent.mdx` with that existing usage.

- `package.json` — defines lint scripts. Use whichever scripts are wired up for documentation lint and preflight (discover at execution time; do not invent).

Out of scope (do NOT do these as part of this plan):

- Auditing other internal links in the docs corpus.
- Adding a redirect rule in `docs.json`.
- Editing any other changelog file or any other MDX content.
- Adding tests or new lint rules.

## Plan of Work

Make a single edit to `docs/changelog/agent.mdx`:

- File: `docs/changelog/agent.mdx`
- Location: line 26 — the line containing `[Device API](/developers/device-api)`. Re-locate by content, not by line number, in case the file has shifted.
- Change: replace the substring `(/developers/device-api)` with `(/developers/device-api/overview)` on that line. The link text `Device API` is unchanged.

Do NOT touch any other file.

## Concrete Steps

All commands run from the docs repo root.

1. Confirm the current state of the target line.

       grep -n "(/developers/device-api)" docs/changelog/agent.mdx

   Expected output (line number may differ; the content is what matters):

       26:- Resources fetched from on-device applications via the [Device API](/developers/device-api) are now fetched from the Miru control plane if not found locally

2. Edit `docs/changelog/agent.mdx`. On the matched line, change `(/developers/device-api)` to `(/developers/device-api/overview)`. The link text `[Device API]` stays the same.

3. Verify the dead link is gone:

       grep -n "(/developers/device-api)" docs/changelog/agent.mdx

   Expected: no matches (exit code 1, no output). The trailing `)` immediately after `device-api` is what makes this assertion specific; the corrected link ends in `overview)` and will not match.

4. Verify the corrected link is present on the patched line:

       grep -n "/developers/device-api/overview" docs/changelog/agent.mdx

   Expected: at least one match, on the same line as the original edit.

5. Confirm the diff is limited to the one file and one line:

       git diff --stat
       git diff docs/changelog/agent.mdx

   Expected: `docs/changelog/agent.mdx | 2 +- 1 file changed, 1 insertion(+), 1 deletion(-)` (or similar). No other files changed.

6. Run the project's documentation lint script (discover the exact invocation from `package.json` — typically `pnpm lint`).

7. Run the project's lint-test script if present (typically `pnpm test:lint`).

8. Run preflight (typically `./scripts/preflight.sh` if present). **Preflight MUST report clean (exit 0) before changes are published.** Do not skip or work around any preflight finding — resolve the underlying issue and re-run until clean.

9. Commit the single-file edit:

       git add docs/changelog/agent.mdx
       git commit -m "fix(changelog/agent): point Device API link to /overview"

   This is a single-milestone plan; this commit is the milestone-end commit. Do NOT amend later changes onto it — if further fixes are needed, create a new commit.

10. Confirm the commit:

        git log -1 --stat

    Expected: one file changed (`docs/changelog/agent.mdx`), one insertion, one deletion.

## Validation and Acceptance

Acceptance is observable, not structural:

- **Dead link is gone.** `grep -n "(/developers/device-api)" docs/changelog/agent.mdx` returns no matches (exit 1, no output).
- **Corrected link is present.** `grep -n "/developers/device-api/overview" docs/changelog/agent.mdx` returns at least one match, on the patched line.
- **Lint passes.** The repo's documentation lint script (e.g. `pnpm lint`) exits 0.
- **Preflight is clean.** The repo's preflight script (e.g. `./scripts/preflight.sh`) exits 0 with no failures. **This is the gate for publishing — preflight must report clean before changes are published.**
- **Diff scope is correct.** `git diff origin/main -- docs/changelog/agent.mdx` shows exactly one line changed, exchanging `(/developers/device-api)` for `(/developers/device-api/overview)`. No other files changed by this plan.

There is no existing test in this repo that asserts internal-link validity for this specific link, so there is no "test X fails before and passes after" assertion to record. A one-character link fix under the existing lint regime does not justify adding a new test harness.

## Idempotence and Recovery

- The edit is a single substring replacement. Re-running the same edit on an already-edited file is a no-op (the source substring `(/developers/device-api)` is no longer present).
- If lint or preflight fails: read the failure, fix the underlying cause, and re-run. Do not bypass with `--no-verify` or env overrides.
- To revert: `git revert <commit-sha>` of the single commit produced by this plan. There are no migrations, no data, and no other files involved, so revert is safe.

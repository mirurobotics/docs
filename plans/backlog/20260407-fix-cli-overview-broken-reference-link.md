# Fix broken CLI reference link in CLI overview page

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Fix broken href in CLI overview page |

This plan lives in `docs/plans/` because the only change is in the docs repo.

## Purpose / Big Picture

The CLI overview page (`docs/developers/cli/overview.mdx`) contains a `CardNewTab` component linking to `/docs/references/cli`. That path has no corresponding file — the only CLI reference page is `docs/references/cli/release-create.mdx`. Visitors who click the "CLI Reference" card land on a 404. After this fix, the card links directly to the `release-create` reference page.

## Progress

- [ ] Update `href` in `docs/developers/cli/overview.mdx` line 25
- [ ] Commit the change
- [ ] Verify preflight passes

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

- Decision: Point the broken link directly at `docs/references/cli/release-create` rather than creating a new index page.
  Rationale: Only one CLI reference page exists. Creating an index would be scope creep; fixing the link is the minimal correct change.
  Date/Author: 2026-04-07 / hunt scan

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

**Docs site:** Mintlify site in the `mirurobotics/docs` repo, checked out at `/home/ben/miru/workbench4/docs/`.

**Broken file:** `docs/developers/cli/overview.mdx` — the CLI section overview page. Line 25 contains:

    <CardNewTab
      title="CLI Reference"
      icon="terminal"
      href="/docs/references/cli"
      arrow
    >

The path `/docs/references/cli` does not map to any file. Mintlify resolves pages from MDX files; without an `index.mdx` or a file at that exact path, the link is a 404.

**Correct target:** `docs/references/cli/release-create.mdx` — the only CLI reference page, registered in `docs.json` as `docs/references/cli/release-create`.

## Plan of Work

Edit `docs/developers/cli/overview.mdx`, line 25: change `href="/docs/references/cli"` to `href="/docs/references/cli/release-create"`. No other files need changing.

## Concrete Steps

From `docs/` (i.e. `/home/ben/miru/workbench4/docs/`):

1. Open `docs/developers/cli/overview.mdx` and on line 25 change:

       href="/docs/references/cli"

   to:

       href="/docs/references/cli/release-create"

2. Verify the change looks correct:

       grep -n "references/cli" docs/developers/cli/overview.mdx
       # Expected: 25:  href="/docs/references/cli/release-create"

3. Commit (Milestone 1):

       git add docs/developers/cli/overview.mdx
       git commit -m "fix(docs): correct broken CLI reference link in CLI overview"

## Validation and Acceptance

**Lint/preflight:** Run `pnpm run lint` from `docs/`. Expect zero errors. Preflight must report `clean` before a PR is opened.

**Link check:** Confirm `docs/references/cli/release-create.mdx` exists (from `docs/` working directory):

    ls docs/references/cli/release-create.mdx
    # Expected: file listed, no error

**Manual verify:** After the fix, `grep -n "references/cli" docs/developers/cli/overview.mdx` should show `href="/docs/references/cli/release-create"` with no mention of the bare `/docs/references/cli` path.

## Idempotence and Recovery

The single-line edit is safe to apply repeatedly. To revert: change `href="/docs/references/cli/release-create"` back to `href="/docs/references/cli"` (though that restores the broken state). Git history is the rollback path.

# Update Platform API "latest" pointers to 2026-05-06.rainier

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Title

Repoint the two "latest" Platform API references from `2026-03-09.tetons` to the newly published `2026-05-06.rainier`.

## Goal

After the rainier changelog entry was published (see `plans/completed/20260506-platform-api-changelog-rainier.md`), two locations in the repo still hardcode the previous version (`2026-03-09`) as the "latest" target. Both must point at `2026-05-06` so that:

1. The redirect `/docs/references/platform-api/latest/:slug*` lands readers on the rainier reference pages.
2. The `<PlatformApiLink>` snippet component (used inline throughout the docs to deep-link into the latest Platform API reference) generates `href`s into the rainier reference pages.

This is a minimal two-line edit. Historical references to `2026-03-09.tetons` in the changelog body and in the supported-versions table are intentional and MUST NOT be touched.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit two files: `docs.json` and `snippets/components/platform-api-link.jsx`. |

This plan lives in `docs/plans/` because all edits are confined to the `docs` repo.

## Files

Two files edited:

- `docs.json` — the Mintlify site config; the `redirects` array entry that maps `latest/:slug*` to a concrete version.
- `snippets/components/platform-api-link.jsx` — JSX snippet component that builds `href`s for inline Platform API endpoint links.

Files explicitly OUT OF SCOPE (intentional historical references; do NOT touch):

- `docs/changelogs/platform-api.mdx` — the historical changelog body that references `2026-03-09.tetons` in its tetons section, code samples, SDK notes, and migration steps.
- `docs/changelogs/product.mdx` — the product-changelog blurb that announced `2026-03-09.tetons`.
- `docs/developers/platform-api/versioning.mdx` — the supported-versions table; the `2026-03-09.tetons` row stays as a still-supported version.
- `docs/developers/platform-api/sdks.mdx` — SDK compatibility notes that pair specific SDK ranges to `2026-03-09.tetons`.
- `docs.json` lines 292–295 — the API-reference dropdown source/directory entries for the `2026-03-09.tetons` reference set itself (these point at the per-version OpenAPI spec, not at "latest").

## Context and Orientation

The `docs` repo is a Mintlify documentation site. Two mechanisms surface the "latest Platform API" reference to readers:

1. A redirect in `docs.json`: `/docs/references/platform-api/latest/:slug*` → a concrete version path. The destination needs to track whichever version is currently the latest (rainier as of 2026-05-06).
2. A reusable JSX snippet `<PlatformApiLink endpoint="..." />` in `snippets/components/platform-api-link.jsx` that builds an anchor into the latest Platform API reference page set. The component already documents itself as "Link to a Platform API endpoint using the latest API version" in its JSDoc, so the comment is correct and only the hardcoded path needs updating.

Both files currently still hardcode `2026-03-09`. The newly published version is `2026-05-06.rainier`; its reference pages are served under `/docs/references/platform-api/2026-05-06/...` (consistent with the existing `2026-03-09` URL shape — only the date prefix changes; the codename `.rainier` is not part of the URL path).

## Plan of Work

The work is a single content-edit milestone followed by a preflight milestone.

### Milestone 1 — Repoint the two latest pointers

1. `docs.json` line 384: change

       "destination": "/docs/references/platform-api/2026-03-09/:slug*"

   to

       "destination": "/docs/references/platform-api/2026-05-06/:slug*"

2. `snippets/components/platform-api-link.jsx` line 9: change

       const href = `/docs/references/platform-api/2026-03-09/endpoints/${endpoint}`;

   to

       const href = `/docs/references/platform-api/2026-05-06/endpoints/${endpoint}`;

The JSDoc comment on line 2 ("Link to a Platform API endpoint using the latest API version") remains correct and is NOT edited.

### Milestone 2 — Preflight

Run `./scripts/preflight.sh` from the repo root. The change touches only string literals inside well-formed JSON and JSX, so no lint/cspell failures are expected. Re-run until exit code 0. **Preflight must report `clean` before changes are published.**

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs/` unless otherwise stated.

### Setup

1. Confirm the working branch:

       git branch --show-current
       # expect: docs/platform-api-changelog-rainier

### Milestone 1: Repoint the two latest pointers

1. Edit `docs.json` line 384 — update the `destination` from `/docs/references/platform-api/2026-03-09/:slug*` to `/docs/references/platform-api/2026-05-06/:slug*`.

2. Edit `snippets/components/platform-api-link.jsx` line 9 — update the `href` template from `/docs/references/platform-api/2026-03-09/endpoints/${endpoint}` to `/docs/references/platform-api/2026-05-06/endpoints/${endpoint}`.

3. Verify the diff is exactly two lines and only changes `2026-03-09` → `2026-05-06` in the two intended locations:

       git diff --stat
       git diff docs.json snippets/components/platform-api-link.jsx

4. Confirm both files still parse cleanly:

       node -e "JSON.parse(require('fs').readFileSync('docs.json','utf8')); console.log('docs.json: OK')"
       node --check snippets/components/platform-api-link.jsx
       # JSX is not native node; if --check rejects on JSX, rely on pnpm run test:lint (Mintlify validation) instead.

### Milestone 2: Preflight

1. Run preflight:

       ./scripts/preflight.sh

   Expected: exits 0 and prints the lint/audit/shell-test success lines. Specifically `pnpm run test:lint` (which includes Mintlify config validation and ESLint over `snippets/`) must pass — this is the load-bearing signal for both `docs.json` schema correctness and `platform-api-link.jsx` JSX validity.

2. Confirm git working tree is clean (no preflight artifacts):

       git status

## Test steps

1. Lint check (Mintlify validation + ESLint, included in preflight but called out separately):

       pnpm install   # only if node_modules is missing or stale
       pnpm run test:lint

   Expected: exits 0 with no errors. This validates the `docs.json` redirects schema and lints the JSX snippet.

2. Confirm `docs.json` parses as JSON:

       node -e "JSON.parse(require('fs').readFileSync('docs.json','utf8'))"

3. Verify no other "latest" pointers still reference `2026-03-09`. After the edits, the only remaining occurrences of `2026-03-09` in source should be intentional historical references:

       grep -rn "2026-03-09" docs.json snippets/

   Expected matches: only `docs.json` lines 292/294/295 (the per-version reference entries — out of scope).

       grep -rn "2026-03-09" docs/

   Expected matches: only the historical changelog/product/versioning/sdks entries listed under "Files OUT OF SCOPE" above. None of these should be edited.

4. Optionally render locally (not required):

       pnpm exec mint dev

   Then visit `/docs/references/platform-api/latest/` and confirm it redirects to a `/docs/references/platform-api/2026-05-06/...` URL, and that any page using `<PlatformApiLink>` (e.g. pages in `docs/developers/platform-api/`) generates anchors into the `2026-05-06` reference pages.

## Validation

Acceptance criteria — each item must be observably true:

1. `docs.json` line 384 reads:

       "destination": "/docs/references/platform-api/2026-05-06/:slug*"

2. `snippets/components/platform-api-link.jsx` line 9 reads:

       const href = `/docs/references/platform-api/2026-05-06/endpoints/${endpoint}`;

3. The JSDoc on line 2 of `platform-api-link.jsx` is unchanged and still reads "Link to a Platform API endpoint using the latest API version".

4. `docs.json` parses as valid JSON (`node -e "JSON.parse(...)"` exits 0).

5. `pnpm run test:lint` exits 0, including Mintlify config validation and ESLint over `snippets/`.

6. `grep -rn "2026-03-09" docs.json snippets/` returns only the three per-version reference entries in `docs.json` (lines ~292/294/295). No `snippets/` matches remain.

7. The historical references in `docs/changelogs/platform-api.mdx`, `docs/changelogs/product.mdx`, `docs/developers/platform-api/versioning.mdx`, and `docs/developers/platform-api/sdks.mdx` are unmodified (verified via `git diff --stat` showing only `docs.json` and `snippets/components/platform-api-link.jsx` changed).

8. **Preflight reports `clean`**: `./scripts/preflight.sh` exits 0 with no warnings. **Preflight must report `clean` before changes are published.**

## Idempotence and Recovery

Both edits are pure text edits to two well-isolated string literals. No external state is mutated.

- Re-running is safe: revert the commit (`git revert <sha>`) or edit the strings back to `2026-03-09` to restore the file to its pre-change state.
- If preflight fails after a commit, fix the underlying issue and add a NEW commit (do not amend).

## Progress

- [ ] Milestone 1: Repoint `docs.json` redirect destination and `platform-api-link.jsx` href template to `2026-05-06`.
- [ ] Milestone 2: Run preflight; address any findings; confirm `clean`.

## Surprises & Discoveries

(none yet)

## Decision Log

- Decision: Use `2026-05-06` (the date prefix only) in both URL paths, not `2026-05-06.rainier`.
  Rationale: The existing `2026-03-09` pointers use the date prefix only — the codename suffix is not part of the URL shape. The corresponding reference pages live under `/docs/references/platform-api/2026-05-06/`, mirroring the `2026-03-09` directory layout.
  Date/Author: 2026-05-06 / planner.

- Decision: Do NOT modify the historical `2026-03-09.tetons` references in the changelogs, supported-versions table, or SDK compatibility notes.
  Rationale: Those references are intentional and time-stamped — the changelog body documents what was true at the tetons release, and the supported-versions table lists tetons as a still-supported (non-latest) version. Repointing them would falsify history.
  Date/Author: 2026-05-06 / planner.

- Decision: Leave the JSDoc comment on `platform-api-link.jsx` line 2 ("Link to a Platform API endpoint using the latest API version") unchanged.
  Rationale: The comment is version-agnostic and remains accurate after the edit. No comment churn needed.
  Date/Author: 2026-05-06 / planner.

## Outcomes & Retrospective

(populate after implementation)

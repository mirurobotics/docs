# Change changelog URLs to singular `/changelog`

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `mirurobotics/docs` (`/home/ben/miru/workbench4/repos/docs/`) | read-write | Rename the Mintlify changelog route from `/changelogs/...` to `/changelog/...`, update links, and keep compatibility redirects. |

This plan lives in `plans/backlog/` in the docs repo because all implementation work is owned by `mirurobotics/docs`. The requested base branch is `main`.

## Purpose / Big Picture

After this change, public documentation changelog pages use the singular URL path. For example, `https://docs.mirurobotics.com/changelog/product` serves the product changelog, and internal links point to `/changelog/...`. Existing plural URLs such as `/changelogs/product` continue to redirect to the new singular path.

## Progress

- [x] Prepare a branch from `main`.
- [x] Rename the changelog content route and update docs navigation, navbar, links, and redirects.
- [x] Run tests, Mintlify checks, and preflight; publish only after preflight reports `clean`.

## Surprises & Discoveries

- The working tree was already on `fix/changelog-singular-url`, based on `main`, with the active plan committed.
- The existing `.gitignore` pattern `changelog/` ignored the new `docs/changelog/` content route, so the pattern was anchored to `/changelog/`.
- `mint dev` printed `TypeError: controller[kState].transformAlgorithm is not a function` while starting, but the preview server stayed up and served the tested routes.

## Decision Log

- Decision: Rename the on-disk Mintlify content folder from `docs/changelogs/` to `docs/changelog/` instead of only adding redirects.
  Rationale: Mintlify derives page routes from file paths relative to the `docs/` content root, so the singular public route should be represented by the content path.
  Date/Author: 2026-04-30 / Codex
- Decision: Anchor the generated-asset ignore pattern as `/changelog/` instead of `changelog/`.
  Rationale: The unanchored pattern ignores any directory named `changelog`, which would block normal future additions under the new tracked `docs/changelog/` content route.
  Date/Author: 2026-04-30 / Codex

## Outcomes & Retrospective

Changelog content now lives under `docs/changelog/`, docs navigation and the navbar point at singular `/changelog/...` URLs, and internal documentation links have been rewritten to singular routes. Explicit redirects preserve compatibility for `/changelogs/:slug*` and `/docs/changelogs/:slug*`.

Focused validation passed: `pnpm run test:lint`, `pnpm run lint`, `../node_modules/.bin/mint broken-links`, local Mintlify route checks, and `./scripts/preflight.sh`.

## Context and Orientation

The repo root is `/home/ben/miru/workbench4/repos/docs/`. Make commits from this repo root, not from the Miru workbench root. The current `main` branch already serves docs without a `/docs` prefix, so this plan only changes the changelog segment from plural to singular.

Mintlify content lives under `docs/`, and `docs/docs.json` is the Mintlify configuration. The current changelog pages are:

- `docs/changelogs/product.mdx`
- `docs/changelogs/cli.mdx`
- `docs/changelogs/agent.mdx`
- `docs/changelogs/device-api.mdx`
- `docs/changelogs/platform-api.mdx`

The current plural route appears in `docs/docs.json` navigation page strings and the navbar href. Internal MDX links also point to `/changelogs/...`, especially from developer pages and changelog pages. Asset URLs such as `https://assets.mirurobotics.com/docs/changelog/...` are not documentation routes and must not be changed.

The existing `docs/docs.json` redirects include `/docs/:slug*` to `/:slug*` for old root-prefix compatibility. This change should add explicit plural-to-singular changelog redirects so old links avoid a redirect chain where practical.

## Plan of Work

Use `git mv` to rename `docs/changelogs/` to `docs/changelog/`.

In `docs/docs.json`, change Changelog navigation page entries from `changelogs/product`, `changelogs/cli`, `changelogs/agent`, `changelogs/device-api`, and `changelogs/platform-api` to the matching `changelog/...` paths. Change the navbar Changelog href from `https://docs.mirurobotics.com/changelogs/product` to `https://docs.mirurobotics.com/changelog/product`.

In `docs/docs.json`, add redirects before the generic `/docs/:slug*` redirect:

    {
      "source": "/changelogs/:slug*",
      "destination": "/changelog/:slug*"
    },
    {
      "source": "/docs/changelogs/:slug*",
      "destination": "/changelog/:slug*"
    }

Update internal docs links in `docs/**/*.mdx` from `/changelogs/...` to `/changelog/...`. Do not rewrite `https://assets.mirurobotics.com/docs/changelog/...`, CSS class names such as `changelog-page`, or historical references inside `plans/`.

## Concrete Steps

All commands run from `/home/ben/miru/workbench4/repos/docs/` unless a step says otherwise.

1. Start from the requested base branch:

       git fetch origin
       git switch main
       git pull --ff-only
       git switch -c docs/changelog-singular-url

   Expected output includes `Switched to a new branch 'docs/changelog-singular-url'`. If that branch already exists, use `git switch docs/changelog-singular-url` after confirming it is based on `main`.

2. Rename the content directory:

       git mv docs/changelogs docs/changelog
       test -f docs/changelog/product.mdx
       test ! -e docs/changelogs

3. Edit `docs/docs.json` and MDX links as described in Plan of Work. Use these searches to drive and verify the edits:

       rg -n 'changelogs|docs\.mirurobotics\.com/changelogs' docs/docs.json docs

   Expected after editing: no `changelogs` matches remain in MDX content or navigation page entries. The only intentional matches should be redirect `source` values in `docs/docs.json`, such as `/changelogs/:slug*` and `/docs/changelogs/:slug*`.

4. Confirm the singular files and links exist:

       test -f docs/changelog/product.mdx
       test -f docs/changelog/cli.mdx
       test -f docs/changelog/agent.mdx
       test -f docs/changelog/device-api.mdx
       test -f docs/changelog/platform-api.mdx
       rg -n '"changelog/product"|/changelog/product|docs\.mirurobotics\.com/changelog/product' docs/docs.json docs

   Expected output includes `docs/docs.json` navigation or navbar references and any legitimate MDX links to `/changelog/product`.

5. Run focused validation:

       pnpm run test:lint
       pnpm run lint

   Expected output: `pnpm run test:lint` exits `0`; `pnpm run lint` exits `0` and ends with `All documentation lint checks passed.`

6. Run Mintlify link validation from the content root. If dependencies are missing, first run `pnpm install --frozen-lockfile` from the repo root.

       cd docs
       ../node_modules/.bin/mint broken-links

   Expected output: the command exits `0` with no broken links.

7. Preview the routes locally:

       cd /home/ben/miru/workbench4/repos/docs/docs
       ../node_modules/.bin/mint dev

   Expected behavior in the printed local preview URL:

   - `/changelog/product` renders the product changelog.
   - `/changelog/cli`, `/changelog/agent`, `/changelog/device-api`, and `/changelog/platform-api` render their pages.
   - `/changelogs/product` redirects to `/changelog/product`.
   - `/docs/changelogs/product` redirects to `/changelog/product` if local Mintlify dev applies redirects.

   Stop the dev server with `Ctrl-C`.

8. Run full preflight before publishing:

       cd /home/ben/miru/workbench4/repos/docs
       ./scripts/preflight.sh

   Expected shell result: the script exits `0`. Also invoke the repo preflight workflow if available to the implementing agent; it must report `clean`. Do not publish, merge, or deploy until preflight reports `clean`.

9. Commit the milestone:

       git status --short
       git add docs/docs.json docs/changelog
       git commit -m "docs(changelog): use singular changelog routes"

   Expected result: one commit containing the folder rename, link updates, and redirects.

## Validation and Acceptance

Acceptance is user-visible URL behavior:

- `https://docs.mirurobotics.com/changelog/product` serves the product changelog.
- `https://docs.mirurobotics.com/changelog/cli`, `/changelog/agent`, `/changelog/device-api`, and `/changelog/platform-api` serve the corresponding changelog pages.
- The docs navbar Changelog link points to `https://docs.mirurobotics.com/changelog/product`.
- Internal links that previously pointed to `/changelogs/...` now point to `/changelog/...`.
- Old plural URLs such as `/changelogs/product` redirect to `/changelog/product`.
- Old prefixed plural URLs such as `/docs/changelogs/product` redirect to `/changelog/product`.

Required test steps:

- From `/home/ben/miru/workbench4/repos/docs/`, run `pnpm run test:lint` and expect exit code `0`.
- From `/home/ben/miru/workbench4/repos/docs/`, run `pnpm run lint` and expect `All documentation lint checks passed.`
- From `/home/ben/miru/workbench4/repos/docs/docs/`, run `../node_modules/.bin/mint broken-links` and expect exit code `0`.
- From `/home/ben/miru/workbench4/repos/docs/docs/`, run `../node_modules/.bin/mint dev` and manually verify the acceptance URLs in the local preview.
- Before changes are published, preflight must report `clean`. If preflight reports anything other than `clean`, fix the findings and rerun preflight until it reports `clean`.

## Idempotence and Recovery

The branch setup and search steps are safe to repeat. The `git mv docs/changelogs docs/changelog` step is safe from a clean tree; if it was already completed, do not run it again. If the rename was done incorrectly before commit, reverse it with `git mv docs/changelog docs/changelogs` and repeat the rename.

If `pnpm run lint` reports a dead redirect for `/changelogs/:slug*`, verify that `docs/changelogs/` no longer exists and that redirect destinations use `/changelog/:slug*`. If Mintlify rejects wildcard redirects, replace them with explicit redirects for the five changelog pages: product, cli, agent, device-api, and platform-api.

Do not use broad search-and-replace over every `changelog` substring. Asset URLs under `https://assets.mirurobotics.com/docs/changelog/...` and CSS class names like `changelog-page` are intentionally singular already and are not route links.

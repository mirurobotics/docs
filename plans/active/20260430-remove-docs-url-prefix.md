# Remove `/docs` from public documentation URLs

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `mirurobotics/docs` (`/home/ben/miru/workbench4/repos/docs/`) | read-write | Reorganize the Mintlify documentation root, update Mintlify paths and redirects, and update repo-local lint/test tooling. |
| Mintlify project settings | external deployment setting | Confirm the deployed project uses `docs/` as its documentation root before publishing the repo change. |

This plan lives in `plans/backlog/` in the docs repo because all repo changes are owned by `mirurobotics/docs`.

## Purpose / Big Picture

After this change, the published docs no longer expose the repository folder name in URLs. A page currently served at `https://docs.mirurobotics.com/docs/getting-started/intro` is served at `https://docs.mirurobotics.com/getting-started/intro`, and old `/docs/...` links redirect to the new root-relative URLs. The actual Mintlify content still lives under the repo's `docs/` directory, separate from repo plumbing such as `scripts/`, `tests/`, `tools/`, `plans/`, and package files.

## Progress

- [x] Create or update an implementation branch from `main`.
- [x] Make `docs/` the Mintlify content root and move Mintlify-owned config, snippets, styling, and public assets under it.
- [x] Remove the public `/docs` prefix from `docs/docs.json` navigation, OpenAPI paths, navbar links, redirects, and internal content links.
- [x] Update lint tooling and fixtures to use `docs/` as the content root.
- [x] Run tests, local Mintlify checks, and preflight; publish only after preflight reports `clean`.

## Surprises & Discoveries

- `pnpm exec mint broken-links` from `docs/` still resolved the CLI from the package root and scanned repo-local planning and fixture files. Running `../node_modules/.bin/mint broken-links` from `docs/` honored the intended Mintlify content root and passed.
- `mint dev` served and redirected the checked routes from `docs/`, but printed a transient `TypeError: controller[kState].transformAlgorithm is not a function` while preparing the preview. The preview recovered and responded successfully.
- The external Mintlify project setting still must be confirmed outside this repo before publishing: the documentation root must be `docs/`.

## Decision Log

- Decision: Keep `docs/` as the separate Mintlify content root instead of moving documentation files into the repository root.
  Rationale: Mintlify derives page URLs from paths relative to its documentation root. Making `docs/` the project root removes the public `/docs` prefix while keeping content isolated from repo tooling.
  Date/Author: 2026-04-30 / Codex

## Outcomes & Retrospective

Implemented the repo-side migration. Local validation confirmed new root URLs render, old `/docs/...` URLs redirect, and lint/test checks pass against the new content-root layout.

## Context and Orientation

The repo root is `/home/ben/miru/workbench4/repos/docs/`. Start from branch `main` and make commits from this repo root, not from the Miru workbench root.

Mintlify page paths map to file paths relative to the Mintlify documentation root. Today `docs.json` sits at repo root and lists pages such as `docs/getting-started/intro`, so Mintlify serves that file at `/docs/getting-started/intro`. The desired state is for Mintlify's documentation root to be the repo's `docs/` directory. Then `docs/getting-started/intro.mdx` is addressed inside Mintlify as `getting-started/intro` and is served publicly at `/getting-started/intro`.

Current key files and directories:

- `docs.json`: Mintlify configuration at repo root. It contains navigation page paths with `docs/`, OpenAPI `source` and `directory` values under `docs/references/...`, a Changelog navbar link to `https://docs.mirurobotics.com/docs/changelogs/product`, and redirects whose `source` and `destination` begin with `/docs/`.
- `docs/`: current documentation page tree. This directory should become the Mintlify root. Its children such as `getting-started/`, `learn/`, `developers/`, `admin/`, `references/`, and `changelogs/` should remain under `docs/`.
- `snippets/`: reusable MDX and JSX snippets imported as `/snippets/...`. Move this directory to `docs/snippets/` so those imports keep working when `docs/` is the Mintlify root.
- `logo/`, `favicon-black.svg`, `favicon-white.svg`, and `style.css`: Mintlify public assets/config-adjacent files currently at repo root. Move them under `docs/` so `/logo/...`, `/favicon-*.svg`, and custom CSS remain visible to Mintlify after the root changes.
- `.mintignore`: currently ignores `.ai/`, `api/`, and `tests/`. Once Mintlify uses `docs/` as its root, root-level repo plumbing is no longer in the deployed content tree. Keep this file only if local `mint` commands still consult the repo root; otherwise add a minimal `docs/.mintignore` only if local validation shows it is needed.
- `scripts/lint.sh`: docs lint entrypoint. It currently sets `content_root="${DOCS_LINT_ROOT:-${repo_root}}"`, searches `${content_root}/docs` and `${content_root}/snippets`, and checks OpenAPI specs under `${content_root}/docs/references`.
- `tests/test-lint.sh` and `tests/lint-fixtures/*`: fixture-driven lint tests. Fixtures currently put sample pages under `tests/lint-fixtures/<name>/docs/example.mdx`.
- `tools/lint/main.go`: custom Go linter. It finds the content root by walking upward from an MDX file until it finds a sibling `snippets/` directory.
- `tools/lint/linter/redirects/redirects.go`: redirect validator. It currently assumes public redirect paths must start with `/docs/` and maps `/docs/foo` to an on-disk `docs/foo.mdx` page under the content root. This must change because public paths should be root-relative while on-disk files live under the new content root.

Do not rewrite external asset URLs such as `https://assets.mirurobotics.com/docs/v03/...`; those URLs are not Mintlify page routes and intentionally contain `/docs/` in the asset service path.

## Plan of Work

Milestone 1 changes the repository layout for Mintlify. Move `docs.json` to `docs/docs.json`. Move `snippets/` to `docs/snippets/`, `logo/` to `docs/logo/`, `favicon-black.svg` and `favicon-white.svg` to `docs/`, and `style.css` to `docs/`. Leave scripts, tests, tools, plans, package files, and CI configuration at repo root.

Milestone 2 updates Mintlify routing data and content links. In `docs/docs.json`, remove the leading `docs/` from navigation page strings and OpenAPI `source`/`directory` values. Change the Changelog navbar URL from `https://docs.mirurobotics.com/docs/changelogs/product` to `https://docs.mirurobotics.com/changelogs/product`. Update redirects so current version aliases point to root-relative destinations, for example `/references/device-api/latest/:slug*` to `/references/device-api/v0.2.1/:slug*`, and add old URL compatibility redirects from `/docs/...` to `/...`. Add a wildcard redirect from `/docs/:slug*` to `/:slug*` unless Mintlify validation reports that a more explicit list is required. Update MDX and YAML content links from `/docs/...` to `/...`, but leave `/snippets/...`, `/logo/...`, `/favicon-*.svg`, and external asset URLs unchanged.

Milestone 3 updates tooling for the new content root. In `scripts/lint.sh`, default `content_root` to `${repo_root}/docs`, collect MDX files under the content root while excluding generated or hidden folders if any, collect snippets from `${content_root}/snippets`, and collect OpenAPI specs under `${content_root}/references`. Keep `DOCS_LINT_ROOT` as a fixture override pointing directly at a Mintlify content root. Update `tests/lint-fixtures/*` so each fixture root mirrors the new content-root shape, with `example.mdx` and any needed `snippets/` at the fixture root rather than under a nested `docs/` directory.

Milestone 4 updates redirect lint semantics and tests. In `tools/lint/linter/redirects/redirects.go`, remove the requirement that redirect source and destination paths begin with `/docs/`. For any non-HTTP path that starts with `/`, strip the leading slash and resolve that path under `contentRoot`. Continue to reject missing leading slashes, `.` or `..` path segments, missing destinations, dead redirects, and invalid wildcard prefixes. Add or update tests in `tools/lint/linter/redirects/redirects_test.go` for root-relative redirects, `/docs/:slug*` old-to-new redirects, root-relative OpenAPI escape hatches, and rejection of a destination that does not resolve under the new `docs/` content root. Update `tools/lint/main_test.go` if its fixture setup still assumes nested `docs/` paths.

Milestone 5 validates locally and records the deployment prerequisite. Run lint and tests from the repo root. Run Mintlify local checks from the `docs/` content root so pages resolve as the deployed project will see them. Before publishing, run the repo preflight workflow and require it to report `clean`. Also confirm in Mintlify project settings, before merge or deploy, that the documentation root is configured as `docs/`; otherwise the repo change will publish pages under the wrong root.

## Concrete Steps

All commands in this section run from `/home/ben/miru/workbench4/repos/docs/` unless a step says otherwise.

### Milestone 1 - Move Mintlify-owned files under `docs/`

1. Start from the requested base branch:

       git fetch origin
       git switch main
       git pull --ff-only
       git switch -c docs/remove-docs-url-prefix

   Expected output includes `Switched to a new branch 'docs/remove-docs-url-prefix'`.

2. Move Mintlify-owned files:

       git mv docs.json docs/docs.json
       git mv snippets docs/snippets
       git mv logo docs/logo
       git mv favicon-black.svg docs/favicon-black.svg
       git mv favicon-white.svg docs/favicon-white.svg
       git mv style.css docs/style.css

3. Confirm repo plumbing remains outside `docs/`:

       test -d scripts
       test -d tests
       test -d tools
       test -f package.json
       test -f docs/docs.json
       test -d docs/snippets

4. Commit:

       git add docs docs.json snippets logo favicon-black.svg favicon-white.svg style.css
       git commit -m "chore(docs): make docs directory the Mintlify root"

### Milestone 2 - Remove public `/docs` paths

1. Edit `docs/docs.json`:

   - Replace navigation page strings like `docs/getting-started/intro` with `getting-started/intro`.
   - Replace OpenAPI `source` and `directory` values like `docs/references/device-api/v0.2.1/api.yaml` with `references/device-api/v0.2.1/api.yaml`.
   - Replace the Changelog navbar URL with `https://docs.mirurobotics.com/changelogs/product`.
   - Update existing redirects to root-relative paths and add old `/docs/...` compatibility redirects.

2. Update internal links in `docs/**/*.mdx`, `docs/**/*.yaml`, and `docs/snippets/**/*`:

       rg -n '\]\(/docs/|href="/docs/|href="/docs"|https://docs\.mirurobotics\.com/docs/' docs

   For each match that is a Mintlify page link, remove only the public `/docs` prefix. For example, change `[Miru Agent](/docs/developers/agent/overview)` to `[Miru Agent](/developers/agent/overview)`. Do not change external asset URLs under `https://assets.mirurobotics.com/docs/...`.

3. Check for remaining public `/docs` page links:

       rg -n '\]\(/docs/|href="/docs/|href="/docs"|https://docs\.mirurobotics\.com/docs/' docs

   Expected output: no matches except any intentionally documented legacy URL examples. If intentional examples remain, add an inline comment or nearby text making that intent clear.

4. Commit:

       git add docs
       git commit -m "feat(docs): serve documentation routes from the root path"

### Milestone 3 - Update lint scripts and fixtures

1. Edit `scripts/lint.sh`:

   - Set `content_root="${DOCS_LINT_ROOT:-${repo_root}/docs}"`.
   - Collect page MDX files from the content root, excluding `${content_root}/snippets`.
   - Collect snippet MDX files from `${content_root}/snippets`.
   - Collect OpenAPI YAML files from `${content_root}/references`.
   - Keep diagnostic text accurate, for example `No OpenAPI specs found under ${content_root}/references.`.

2. Update `tests/lint-fixtures/*` to mirror the new content root. For example, move `tests/lint-fixtures/good/docs/example.mdx` to `tests/lint-fixtures/good/example.mdx` and make equivalent moves for `bad-mdx`, `bad-openapi`, and `bad-spelling`.

3. Run fixture tests:

       pnpm run test:lint

   Expected output ends with no error and each fixture assertion completes. If the script is quiet on success, the command exits `0`.

4. Commit:

       git add scripts/lint.sh tests/lint-fixtures tests/test-lint.sh
       git commit -m "test(docs): point lint fixtures at the Mintlify root"

### Milestone 4 - Update redirect lint behavior

1. Edit `tools/lint/linter/redirects/redirects.go` and `tools/lint/linter/redirects/redirects_test.go`:

   - Remove hard-coded `/docs/` prefix validation messages.
   - Resolve `/foo/bar` to `${contentRoot}/foo/bar.mdx` or `${contentRoot}/foo/bar.md`.
   - Keep HTTP and HTTPS destinations exempt from filesystem destination checks.
   - Keep traversal checks for `.` and `..`.
   - Add tests for redirects from `/docs/:slug*` to `/:slug*`, a specific old page such as `/docs/references/cli/install` to `/developers/cli/install`, and root-relative OpenAPI wildcard paths.

2. Update `tools/lint/main_test.go` if it still constructs sample redirect destinations like `/docs/missing`; use root-relative paths matching the new fixture layout.

3. Run Go tests and custom linter checks:

       cd tools/lint
       go test ./...
       ./scripts/covgate.sh
       ./scripts/lint.sh

   Expected output: `go test ./...` exits `0`; `covgate.sh` reports all packages at or above their thresholds; `lint.sh` exits `0`.

4. Commit from the repo root:

       cd /home/ben/miru/workbench4/repos/docs
       git add tools/lint
       git commit -m "fix(lint): validate root-relative docs redirects"

### Milestone 5 - Validate and prepare to publish

1. Run full repo lint:

       pnpm run lint

   Expected output ends with `All documentation lint checks passed.`

2. Run local Mintlify checks from the content root:

       cd docs
       pnpm exec mint broken-links
       pnpm exec mint dev

   Expected behavior: `mint broken-links` exits `0`. In `mint dev`, visit the local preview URL it prints and verify `/getting-started/intro`, `/developers/cli/install`, `/references/cli/release-create`, `/references/device-api/latest/health`, and `/changelogs/product` render without the `/docs` prefix. Stop the dev server with `Ctrl-C`.

3. Check redirects locally if Mintlify dev supports redirect handling in this repo:

       cd docs
       pnpm exec mint dev

   Expected behavior: visiting `/docs/getting-started/intro` redirects to `/getting-started/intro`; `/docs/references/cli/install` redirects to `/developers/cli/install`; `/docs/references/device-api/latest/health` redirects to `/references/device-api/v0.2.1/health`. If local dev does not apply redirects, record that limitation in Surprises & Discoveries and rely on `mint broken-links`, `docs/docs.json` validation, and staging/preview verification.

4. Run the repo preflight workflow before publishing:

       $preflight

   Expected result: preflight reports `clean`. Do not publish, merge, or deploy these changes until preflight reports `clean`.

5. Confirm the external Mintlify deployment setting:

       Verify in Mintlify project settings that the documentation root is `docs/`.

   Expected result: Mintlify reads `docs/docs.json` as the site configuration. If the setting cannot be confirmed, do not publish; the repo layout depends on this setting.

6. Commit any validation-only fixes discovered in this milestone:

       git status --short
       git add <files changed by validation fixes>
       git commit -m "chore(docs): finalize root path docs validation"

   If there are no validation fixes, skip this commit and record the clean validation in Progress.

## Validation and Acceptance

Acceptance is user-visible URL behavior:

- `https://docs.mirurobotics.com/getting-started/intro` serves the page currently known as `docs/getting-started/intro.mdx`.
- `https://docs.mirurobotics.com/developers/cli/install` serves the CLI install page.
- `https://docs.mirurobotics.com/references/cli/release-create` serves the CLI reference page.
- `https://docs.mirurobotics.com/changelogs/product` serves the product changelog, and the navbar Changelog link points to that URL.
- Old URLs such as `https://docs.mirurobotics.com/docs/getting-started/intro` and `https://docs.mirurobotics.com/docs/references/cli/install` redirect to the corresponding non-`/docs` URL.
- Version alias redirects still work at the new root, for example `/references/device-api/latest/:slug*` to `/references/device-api/v0.2.1/:slug*`.
- Repo plumbing remains outside the Mintlify content root: `scripts/`, `tests/`, `tools/`, `plans/`, `package.json`, and `pnpm-lock.yaml` stay at repo root; Mintlify content and assets live under `docs/`.

Required test steps:

- From `/home/ben/miru/workbench4/repos/docs/`, run `pnpm run test:lint` and expect exit code `0`.
- From `/home/ben/miru/workbench4/repos/docs/tools/lint/`, run `go test ./...`, `./scripts/covgate.sh`, and `./scripts/lint.sh`; expect all to exit `0`.
- From `/home/ben/miru/workbench4/repos/docs/`, run `pnpm run lint` and expect `All documentation lint checks passed.`
- From `/home/ben/miru/workbench4/repos/docs/docs/`, run `pnpm exec mint broken-links` and expect exit code `0`.
- From `/home/ben/miru/workbench4/repos/docs/docs/`, run `pnpm exec mint dev` and manually verify the acceptance URLs above in the local preview.
- Before publishing, run the repo preflight workflow and require it to report `clean`. Publishing is blocked until preflight reports `clean`.

## Idempotence and Recovery

The `git mv` steps are safe before commit and can be inspected with `git status --short`. If a move is made incorrectly before commit, move the file back with `git mv <new> <old>` and rerun the milestone. Do not use `git reset --hard` because this repo may contain user work.

Path rewrite steps are repeatable if they are driven by searches for remaining `/docs` page links. Guard against over-replacement by searching only for page-link patterns such as `](/docs/`, `href="/docs/`, and `https://docs.mirurobotics.com/docs/`; do not bulk replace every `docs/` substring because external asset URLs intentionally contain that path segment.

If `mint broken-links` or `mint dev` shows that a wildcard `/docs/:slug*` redirect is unsupported or too broad, replace it with explicit redirects for each top-level section: `/docs/getting-started/:slug*`, `/docs/learn/:slug*`, `/docs/developers/:slug*`, `/docs/admin/:slug*`, `/docs/references/:slug*`, and `/docs/changelogs/:slug*`. Rerun lint and Mintlify validation after changing redirects.

If the external Mintlify project root cannot be set to `docs/`, stop and revise this plan. Serving root-relative URLs while keeping content in a separate folder depends on Mintlify reading `docs/` as the documentation root.

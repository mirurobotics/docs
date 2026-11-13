# Code quality assessment of the docs repo

An audit of `mirurobotics/docs` at `f9152a7` (2026-11-10): where the repo's tooling, CI, generated content and docs fall short, and what to do about it, in priority order. It changes no code. Each finding is meant to become its own ExecPlan in `plans/active/` when picked up.

## How to read this

Each finding gives what is wrong, why it matters, and a suggested fix, with a rough effort (S: under a day, M: a few days, L: a week or more).

- **P1: fix first.** Broken or wrong docs can ship without CI noticing, or the deploy path has a security gap.
- **P2: fix soon.** Manual, copy-paste work and missing checks that make every change slower and riskier.
- **P3: tidy up.** Hygiene gaps that are cheap to close.

| # | Finding | Priority | Effort |
|---|---------|----------|--------|
| 1 | Promote puts a free-text input straight into a shell command | P1 | S |
| 2 | Internal links are not checked, and three are broken today | P1 | S |
| 3 | Generated API reference files have drifted from their generators | P1 | M |
| 4 | Shipping an API version is a copy-paste job with no consistency check | P2 | M |
| 5 | The OpenAPI lint step is slow, deprecated and mostly redundant | P2 | S |
| 6 | Tool versions differ between CI, the devcontainer and preflight | P2 | S |
| 7 | Most pages have no `description` | P2 | M |
| 8 | Image URLs and external links are never checked | P2 | S |
| 9 | The README does not explain how to work in the repo | P2 | S |
| 10 | Shell, Python and JSX code is not linted | P2 | S |
| 11 | CI and CodeQL workflow hygiene | P3 | S |
| 12 | Dead config left over from the old layout | P3 | S |
| 13 | Coverage gates sit well below actual coverage | P3 | S |
| 14 | Nine images have no alt text | P3 | S |
| 15 | Finished plans are still in `plans/active/` | P3 | S |

## What is already good

These are worth keeping as the repo changes:

- The custom Go linter in `tools/lint` enforces ten house rules (heading case, import hygiene, component style, image domain, redirects and more). Every package has a coverage gate, golangci-lint reports 0 issues, and per-package coverage is 93.9% to 100%.
- `tests/test-lint.sh` runs the whole lint pipeline against good and bad fixtures, so a broken check is caught, not only a broken page.
- The redirects rule checks every `docs.json` redirect against the content tree, and the nav lists every operation in every vendored spec today (checked for all 11 specs that use page lists).
- Every third-party action is pinned to a commit SHA, and Dependabot covers actions, npm and Go modules.
- `main` requires `ai-review`, `audit`, `lint` and both custom-linter checks, plus a code owner review from one of two owners.
- `scripts/promote.sh` only promotes commits reachable from `main` and uses `--force-with-lease`.
- The README documents the dependency audit policy, including why each out-of-range override exists and when to remove it.
- Scripts resolve the repo root from their own location, use `set -euo pipefail`, and print usage text.

---

## P1: fix first

### 1. Promote puts a free-text input straight into a shell command

**What.** `.github/workflows/promote.yml` runs `./scripts/promote.sh "${{ inputs.ref }}" "${{ inputs.environment }}"`. `inputs.ref` is a free-text `string`, and GitHub substitutes it into the script before bash runs. The job holds a GitHub App token that can force-push the `staging`, `uat` and `production` branches. The workflow also has no `concurrency` group.

**Why it matters.** A ref such as `main"; <command>; echo "` runs arbitrary commands with that token. Only people with write access can dispatch the workflow, and the environment gate applies, so the risk is limited, but this is the one path to the production docs site and the fix is a few lines. Two promotes to one environment at once also race; `--force-with-lease` stops a clobber, but the second run fails with a confusing error.

**Fix.** Pass inputs through `env:` and quote them in bash:

```yaml
env:
  REF: ${{ inputs.ref }}
  ENVIRONMENT: ${{ inputs.environment }}
run: ./scripts/promote.sh "$REF" "$ENVIRONMENT"
```

Add `concurrency: { group: promote-${{ inputs.environment }}, cancel-in-progress: false }`. Run actionlint in CI (see #10), which flags this pattern.

### 2. Internal links are not checked, and three are broken today

**What.** Nothing in CI checks links between pages. The redirects rule checks only `docs.json` redirects, and `mint validate` does not check links. Running `mint broken-links` from `docs/` reports three broken links, all in `docs/admin/backups/overview.mdx`: two to `/references/platform-api/latest/endpoints` and one to `/references/platform-api/2026-11-04`. Neither is a page: endpoint pages live one level deeper, at `/references/platform-api/<version>/endpoints/<group>/<operation>`.

Links built inside JSX components are not seen by `mint broken-links` either. That covers `DeviceApiLink`, `PlatformApiLink` and `CliLink` (17 uses), and `DeviceApiReleaseLinks` and `PlatformApiReleaseLinks` (15 uses), which link to `/references/<api>/<version>`, the same bare-version pattern `mint broken-links` reports as broken above. `docs/developers/device-api/versions.mdx` links to `/references/device-api/v0.2.0` and `v0.1.0` the same way, through `LinkNewTab`.

**Why it matters.** Readers following "latest Platform API" from the backups page land on a 404. Every page rename or API version bump can break more links, and the component links are the ones most likely to break, since they are built from strings.

**Fix.** Fix the three links. Add `mint broken-links` to `scripts/lint.sh` (it runs offline). Confirm on a preview where the bare version links (`/references/<api>/<version>`) land, and point the release-link components at a real page if they 404. For component links, extend the Go linter with a rule that collects `href`, `endpoint` and `path` props from the known link components and resolves them with the same page and redirect logic the redirects rule already has.

### 3. Generated API reference files have drifted from their generators

**What.**

- **Event pages.** `api/generate_event_pages.py` writes ` ```json Example`. Every committed event page under `docs/references/device-api/v0.2.*/events/` says ` ```json Example nocopy`, so they were edited by hand after generation. Re-running `api/gen-events.sh` for any of the five versions reverts all ten pages (checked by regenerating each one and diffing).
- **The fetch pipeline.** `api/pull-stainless.sh` calls `python3 add_unix_socket_curl.py` and the other helpers by bare filename and creates `.venv` in the current directory, so it only works when run from inside `api/`. It writes to `api/*.stainless.yaml` and `api/events/` (gitignored), and nothing documents or scripts how those files reach `docs/references/`. `plans/completed/20260816-platform-api-changelog-everglades.md` says the pipeline had to be "reverse-engineered" from an old commit.
- **Dependencies.** The helpers need PyYAML, which the scripts `pip install` unpinned. There is no `requirements.txt`, and `gen-events.sh` and `pull-stainless.sh` use different virtualenv locations.
- **Markers.** Generated YAML and MDX carry no "generated, do not edit" header. Several plans warn against hand edits instead.

**Why it matters.** `docs/references/` is about half of what the site serves. A fix made by hand is silently lost the next time someone regenerates, and the next person to ship an API version has to rediscover the process. CI checks the committed files but not that they still match the generators.

**Fix.** Replace the two entry points with one script, for example `api/sync-reference.sh <device|platform> <version>`, that fetches the spec, runs the post-processing, writes straight into `docs/references/...`, and generates the event pages. Move `nocopy` into `generate_event_pages.py`. Pin PyYAML in `api/requirements.txt` (or inline script metadata run with `uv run`). Add a CI check that re-runs `generate_event_pages.py` and `inject_scopes.py` against every vendored spec and fails on any diff; it needs no network and takes seconds (`inject_scopes.py` is already a no-op on the vendored `2026-11-04.yaml`). Write the steps into the README.

---

## P2: fix soon

### 4. Shipping an API version is a copy-paste job with no consistency check

**What.** `docs/docs.json` is 1,782 lines. Lines 334 to 1,467 (about 1,130 lines, 63%) are the Device API and Platform API version dropdowns. Every group in every version repeats the same `openapi.source` and `directory` pair, and a new version means copying a 100 to 150 line block. The two oldest Device API versions (`v0.2.0`, `v0.1.0`) use a different style (whole spec, no page list) and a flat `.yaml` file instead of a `<version>/api.yaml` folder. Which version is "latest" is set by two redirects that must match the first dropdown. The README's release checklist lists about fifteen manual steps across version tables, compatibility matrices and a "latest" component, with no detail on any of them.

**Why it matters.** Each of these is a place to forget something. The nav matches every spec today, but nothing checks it: a new operation left out of the nav simply does not appear, and a `latest` redirect left on the old version passes lint.

**Fix.** Add a Go lint rule (call it `api-versions`) that checks, for each OpenAPI-backed dropdown, that every operation in the spec appears in the nav exactly once and nothing else does; that the `latest` redirects point at the first dropdown; and that every spec under `docs/references/` is used by the nav. Then reduce the repetition: Mintlify lets `openapi` be set on a dropdown, so check that groups can inherit it, or generate the version block with a small script. Expand the README checklist into short steps with file paths.

### 5. The OpenAPI lint step is slow, deprecated and mostly redundant

**What.** `scripts/lint.sh` runs `mint openapi-check` on all 13 vendored specs, one after another. Each call takes about 10 seconds, prints "openapi-check is deprecated, use `mintlify validate` instead" twice, and the whole step is most of the 2.5 minute lint run. The specs for released versions never change. `tests/test-lint.sh` runs the full `lint.sh`, including this step, for each of its five fixtures, and CI already runs `pnpm run validate` (`mint validate`) right after.

**Why it matters.** Slow lint runs on every PR, and the command will eventually be removed from the `mint` CLI, breaking CI with no change in this repo.

**Fix.** Check that `mint validate` fails on the `bad-openapi` fixture. If it does, drop the OpenAPI step and point that fixture at `mint validate`. If not, check only specs changed against `origin/main`, and run them in parallel.

### 6. Tool versions differ between CI, the devcontainer and preflight

**What.**

- **Go.** The CI `lint` job runs `scripts/lint.sh`, which builds the Go linter, but has no `setup-go` step. It uses whatever Go the runner image has, and `GOTOOLCHAIN` then downloads Go 1.25.3 on every run (seen locally as "go: downloading go1.25.3"), with no module cache. The devcontainer installs Go `1.25`, which floats.
- **Node.** CI uses Node 22. The devcontainer image is `typescript-node:1-20-bookworm` (Node 20). `package.json` has no `engines` field and there is no `.nvmrc`.
- **Preflight.** `scripts/preflight.sh` is meant to mirror CI but skips `mint validate`, and it is not in the `package.json` scripts.

**Why it matters.** A change can pass locally and fail in CI, or the reverse. The Go download adds time and a network dependency to every lint run.

**Fix.** Add `actions/setup-go` with `go-version-file: tools/lint/go.mod` and caching to the `lint` job. Add `.nvmrc` (or `engines.node`) at 22 and read it in CI (`node-version-file`) and the devcontainer. Make preflight run the same steps as CI, and add `"preflight": "./scripts/preflight.sh"` to `package.json`.

### 7. Most pages have no `description`

**What.** 99 of the 104 hand-written pages (all pages except snippets and generated event pages) have no `description` in their frontmatter. That includes every CLI reference page and every page under `developers/agent/`.

**Why it matters.** Mintlify uses `description` for the page's meta description and the subtitle under its title. Without it, search engines and AI search tools pick their own snippet, often the first line of a code block or a callout.

**Fix.** Write descriptions, starting with the most-visited pages, then add a `frontmatter` rule to the Go linter that requires `title` and `description` on every non-snippet page.

### 8. Image URLs and external links are never checked

**What.** All 176 images are hosted on `assets.mirurobotics.com` (Cloudflare R2), outside the repo. The image-domain rule checks the host but not that the file exists. Image sets are versioned by folder (`docs/v04/images/...`) and copied by hand, following `rclone.md`. External links (GitHub, vendor docs, the Miru app) are not checked either.

**Why it matters.** A renamed or deleted object in the bucket, or a typo in a new path, shows as a broken image, and nothing reports it.

**Fix.** Add a scheduled workflow (nightly or weekly, not per-PR, so third-party outages don't block merges) that runs `lychee` over `docs/**/*.mdx` and the JSX snippets, and posts to Slack or opens an issue on failure.

### 9. The README does not explain how to work in the repo

**What.**

- **README.** `README.md` (38 lines) opens with "### Agent Updates". It is a release checklist plus the audit policy. It does not say what the repo is, how to install (corepack, pnpm, Go), how to preview (`pnpm dev`), how to lint or run preflight, what lives in `api/`, `tools/lint/` and `plans/`, or how the site is deployed (Promote to the `staging`, `uat` and `production` branches).
- **Agent instructions.** There is no `AGENTS.md` or `CLAUDE.md`, though agents write most plans and many PRs. The plans follow an ExecPlan format that is not defined anywhere in this repo.
- **Skills.** The `campaign`, `loops-cli` and `loops-lmx` skills live in `.cursor/skills/`, which only Cursor loads.
- **PR template.** There is no PR template.

**Why it matters.** Each new contributor and each agent session works the setup out again, and the rules that matter most here (never hand-edit `docs/references/`, use the link components, follow the heading-case rule) are only written down inside old plans.

**Fix.** Rewrite the README around setup, layout, lint, preview and deploy, and move the release checklist under its own heading (see #4). Add an `AGENTS.md`, with `CLAUDE.md` pointing to it, covering the house rules, where plans go and the ExecPlan format (or a link to where it is defined). Decide whether the Loops skills belong in this repo or in the shared agents repo, and make them loadable by the agents that use them.

### 10. Shell, Python and JSX code is not linted

**What.**

- **Shell.** CI runs no shellcheck. Run locally, it reports `SC2295` on `scripts/lint.sh` lines 85 and 86 (`${spec#${repo_root}/}` should quote the inner expansion).
- **Python.** The four Python helpers in `api/` have no lint or format check. `ruff check` reports 7 findings (unsorted imports, two bare `except Exception`, a non-executable file with a shebang), and `ruff format --check` would reformat all four files.
- **JSX.** The 22 JSX components in `docs/snippets/components/` (38 exports) are not linted; `eslint.config.mjs` only targets `.mdx`.
- **Workflows.** CI runs no actionlint or yamllint.

**Fix.** Add shellcheck, `ruff check` and `ruff format --check`, ESLint for `docs/snippets/**/*.jsx` (with `eslint-plugin-react`), and actionlint to `scripts/lint.sh` or a new CI job, and fix what they report.

---

## P3: tidy up

### 11. CI and CodeQL workflow hygiene

- **Permissions.** `ci.yml` has no top-level `permissions`, so the `lint` and `audit` jobs get the repo's default token permissions. Add `permissions: { contents: read }`.
- **Timeouts.** The `lint` and `audit` jobs have no `timeout-minutes`.
- **Redundant runs.** CI runs on pushes to `staging`, `uat` and `production`, which only ever receive commits already tested on `main`.
- **CodeQL.** It runs on pushes to the environment branches and nightly, not on PRs, and it does not scan the Python in `api/`. Its config is named "JavaScript CodeQL Configuration" and ignores a `dist` folder that does not exist. `codeql-action` is pinned with a `# v4` comment, and the org actions with `# main`, so reviewers can't see which version a SHA is.
- **Path filter.** The custom-linter jobs only run when `tools/lint/**` changes, so a PR that edits `ci.yml` or those jobs' scripts skips them.

### 12. Dead config left over from the old layout

- `.mintignore` sits at the repo root, but Mintlify runs from `docs/`, so it appears to have no effect.
- `eslint.config.mjs` and `cspell.json` ignore `.ai/**`, `api/**`, `images/**`, `logo/**`, `output/**`, `videos/**` and `changelog/**` at the repo root. None of these match anything ESLint or CSpell is given.
- `scripts/lint.sh` spell-checks `docs/rclone.md` and `docs/README.md`, which don't exist, so the root `README.md` and `rclone.md` are never spell-checked.
- `.gitignore` lists `node_modules` twice and has entries for `openapitools.json` and `output/` that nothing produces.
- `rclone.md`, an internal runbook, sits at the repo root. Move it under a `runbooks/` folder or into the README.
- `docs/developers/agent/security.mdx` ends with a commented-out TODO about an APT repository.

### 13. Coverage gates sit well below actual coverage

**What.** The root package's `.covgate` is 30% against 96.9% actual coverage, `linter/imagedomain` is 90% against 100%, and `linter/headingcase` has no `.covgate` (it falls back to the 90% default against 95.8%). `tools/lint/scripts/ratchet-covgates.sh` exists but nothing runs it.

**Fix.** Run the ratchet once and commit the results. Add a `.covgate` for `headingcase`. Optionally have CI fail when a gate is more than a few points below actual, so gates keep up.

### 14. Nine images have no alt text

**What.** Six images in `docs/cfg-mgmt/deploy/staging-area.mdx` use `![](...)`, and three `<img>` tags in `docs/changelog/product.mdx` (around line 1063) have no `alt`.

**Fix.** Add alt text, then add an `image-alt` rule to the Go linter so new images need it.

### 15. Finished plans are still in `plans/active/`

**What.** All three plans in `plans/active/` describe work that has shipped: the agent `v0.10.1` changelog entry is at `docs/changelog/agent.mdx:168`, `v0.6.x` and `v0.5.x` show `<EndOfLifeBadge />` in `docs/developers/agent/versions.mdx`, and `instance_format` is documented in `docs/cfg-mgmt/concepts/schemas/languages/opaque.mdx`. Separately, 33 of the 58 plans name absolute paths on one person's machine (for example `/home/ben/miru/workbench2/repos/docs`).

**Fix.** Check off and move the three plans to `plans/completed/`. Ask new plans to use repo-relative paths (part of the `AGENTS.md` in #9).

---

## Suggested order

1. **One small PR for quick wins:** #1, #2 (fix the three links and add `mint broken-links`), #5, #6, #11, #12, #13, #15.
2. **The reference pipeline:** #3, then #4, since the consistency rule is easier to write once the sync script exists.
3. **Linting the code around the docs:** #10, and the component-link rule from #2.
4. **Content, any time:** #7, #8, #9, #14.

## How this was checked

The whole repo was read at `f9152a7`. `./scripts/lint.sh`, `pnpm run test:lint`, `tools/lint/scripts/lint.sh` (golangci-lint, 0 issues) and `tools/lint/scripts/covgate.sh` all pass. `mint broken-links` reports the three links in #2. The event pages were regenerated for all five Device API versions and diffed (#3), and `inject_scopes.py` was run on the latest Platform API spec (no change). Nav entries were compared with the operations in every spec (#4). shellcheck and ruff were run locally (#10). `mint validate` was not run, because it needs network access the audit environment did not have; CI runs it on every PR. Live-site behavior (where the bare version links land, canonical and meta tags) was not checked for the same reason.

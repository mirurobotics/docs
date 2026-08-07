# Link every Platform API operation from the main documentation pages

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Goal

Every callable operation in the Platform API (version `2026-05-06.rainier`) must be explicitly referenced and linked from the prose documentation page that describes the corresponding product operation, using the existing `<PlatformApiLink>` pattern. The pattern already exists in two places — "Edit a device" on `docs/primitives/devices.mdx` links `devices/update`, and "View a release" on `docs/primitives/releases.mdx` links `releases/get` — and must be extended to the remaining 28 operations. One page (`docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx`) links an endpoint with a hardcoded versioned URL; convert it to `<PlatformApiLink>` so it tracks the latest version automatically.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add `**Platform API**` link blocks to nine `.mdx` pages. No spec, config, or component changes. |

This plan lives in `docs/plans/` because all edits are confined to the `docs` repo. Work happens on branch `claude/docs-api-operations-links-9kdem5` off `main`.

## Context and Orientation

The `docs` repo is a Mintlify documentation site. The Platform API reference is generated from the OpenAPI spec at `docs/references/platform-api/2026-05-06.yaml` and served under `/references/platform-api/2026-05-06/endpoints/<tag-slug>/<operation-slug>` (e.g. `.../endpoints/devices/update`). Tag and operation slugs are the spec's tag name and operation summary, lowercased and hyphenated — so `Download Content` under `Config Instances` becomes `config-instances/download-content`.

The reusable snippet `docs/snippets/components/platform-api-link.jsx` exports `<PlatformApiLink endpoint="...">`, which builds an `href` of `/references/platform-api/2026-05-06/endpoints/${endpoint}`. Pages import it with:

    import { PlatformApiLink } from '/snippets/components/platform-api-link.jsx';

The established block style (see "Edit a device" on `docs/primitives/devices.mdx`) is a bold `**Platform API**` label followed by one sentence per endpoint, ending with the link component whose text ends in ` »` and no trailing period:

    **Platform API**

    To edit a device programmatically, use the <PlatformApiLink endpoint="devices/update">update device endpoint »</PlatformApiLink>

### Operation inventory and target locations

The `paths:` section of `docs/references/platform-api/2026-05-06.yaml` defines exactly 30 operations across 9 tags. Each maps to a doc page and section as follows (pre-existing links marked):

| Tag | Operation (slug) | Doc page | Section |
|-----|------------------|----------|---------|
| Config Instances | `config-instances/get` | `docs/cfg-mgmt/primitives/config-instances.mdx` | View a config instance (new section) |
| Config Instances | `config-instances/list` | same | View a config instance (new section) |
| Config Instances | `config-instances/download-content` | same | View a config instance (new section) |
| Config Instances | `config-instances/create` | same | Create a config instance (new section) |
| Config Schemas | `config-schemas/create` | `docs/cfg-mgmt/primitives/schemas/manage.mdx` | Create a schema |
| Config Schemas | `config-schemas/get` | same | View a schema |
| Config Schemas | `config-schemas/list` | same | View a schema |
| Config Types | `config-types/get` | `docs/cfg-mgmt/primitives/config-types.mdx` | View a config type (new section) |
| Config Types | `config-types/list` | same | View a config type (new section) |
| Config Types | `config-types/create` | same | Create a config type |
| Config Types | `config-types/update` | same | Edit a config type |
| Deployments | `deployments/create` | `docs/cfg-mgmt/deploy/staging-area.mdx` | Stage a deployment (`target_status: staged`) AND Patch a deployment (`parent_id`) |
| Deployments | `deployments/get` | same | View a deployment |
| Deployments | `deployments/list` | same | View a deployment |
| Deployments | `deployments/archive` | same | Archive a deployment |
| Deployments | `deployments/deploy` | same | Deploy a deployment |
| Deployments | `deployments/drifts` | same | Review a deployment (end of section) |
| Devices | `devices/get` | `docs/primitives/devices.mdx` | View a device |
| Devices | `devices/list` | same | View a device |
| Devices | `devices/create` | same | Provision a device |
| Devices | `devices/update` | same | Edit a device (**pre-existing** — the exemplar) |
| Devices | `devices/ping` | same | Ping a device |
| Provisioning Tokens | `provisioning-tokens/create` | `docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx` | Create a provisioning token (convert hardcoded URL); also cross-linked from Provision a device on `docs/primitives/devices.mdx` |
| Git Commits | `git-commits/create` | `docs/cfg-mgmt/create-a-release.mdx` | Git commit |
| Git Commits | `git-commits/get` | same | Git commit |
| Git Commits | `git-commits/list` | same | Git commit |
| Principals | `principals/self` | `docs/developers/platform-api/authn.mdx` | end of page, before the authorization pointer |
| Releases | `releases/get` | `docs/primitives/releases.mdx` | View a release (**pre-existing**) |
| Releases | `releases/list` | same | View a release |
| Releases | `releases/create` | same | Create a release |

Irregular slugs to watch (slug is not guessable from the operationId or URL path): `config-instances/download-content` (operationId `getConfigInstanceContent`, path `/config_instances/{id}/content`), `principals/self` (operationId `getSelf`, path `/principal`), `deployments/drifts` (operationId `getDeploymentDrifts`), and the `git-commits/*` family (paths use `git_commits` with an underscore; slugs use a hyphen).

### Doc operations that intentionally get NO link

These dashboard operations have no Platform API endpoint in the 2026-05-06 spec and must be left without a `**Platform API**` block:

- Groups, entirely (`docs/primitives/groups.mdx`: create/edit/delete group, move device, add/edit/remove member) — the API has no Groups tag.
- Delete a device and Move a device (`docs/primitives/devices.mdx`).
- Duplicate a release and Delete a release (`docs/primitives/releases.mdx`).
- Delete a config type (`docs/cfg-mgmt/primitives/config-types.mdx`).

Additionally, `deployment.validate` (`validateDeployment`) appears in the spec under `x-webhooks:`, not `paths:` — it is an event Miru delivers to a customer-hosted endpoint, not a callable API operation, and is out of scope.

## Plan of Work

### Milestone 1 — Add the link blocks

For each of the nine pages, add the `PlatformApiLink` import (in alphabetical position among the existing `/snippets/components/` imports) where not already present, then add `**Platform API**` blocks per the mapping table above:

1. `docs/primitives/devices.mdx` — add blocks under "View a device" (get, list), "Provision a device" (create, plus a sentence pointing at provisioning tokens and linking `provisioning-tokens/create`), and "Ping a device" (ping). "Edit a device" already has its block.
2. `docs/primitives/releases.mdx` — add `releases/list` under "View a release" (below the existing `releases/get` sentence) and a new block under "Create a release" (create). Import already present.
3. `docs/cfg-mgmt/deploy/staging-area.mdx` — six blocks: "Stage a deployment" (create, noting `target_status: staged`), "View a deployment" (get, list), "Patch a deployment" (create, noting `parent_id` set to the deployment being replaced), "Archive a deployment" (archive), "Deploy a deployment" (deploy), and end of "Review a deployment" (drifts).
4. `docs/cfg-mgmt/primitives/config-types.mdx` — add a new "View a config type" section (dashboard pointer plus get/list block) before "Create a config type"; add blocks under "Create a config type" (create) and "Edit a config type" (update). Because these sections now cover two surfaces, label the existing dashboard instructions with a `**Dashboard**` lead-in.
5. `docs/cfg-mgmt/primitives/config-instances.mdx` — add new "View a config instance" (dashboard pointers to the config editor, staging area, and device history, plus get/list/download-content block) and "Create a config instance" (dashboard pointer plus create block) sections at the end of the page.
6. `docs/cfg-mgmt/primitives/schemas/manage.mdx` — block under "Create a schema" (create) and at the end of "View a schema" (get, list).
7. `docs/cfg-mgmt/create-a-release.mdx` — one block at the end of the "Git commit" section (before "Supported Git providers") linking `git-commits/create`, `git-commits/get`, `git-commits/list`.
8. `docs/developers/platform-api/authn.mdx` — one sentence linking `principals/self` for verifying credentials, before the closing authorization pointer.
9. `docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx` — replace the hardcoded `[Platform API](/references/platform-api/2026-05-06/endpoints/provisioning-tokens/create)` markdown link with `<PlatformApiLink endpoint="provisioning-tokens/create">Platform API</PlatformApiLink>`.

Commit the milestone as a single commit: `docs: link platform api endpoints from operation docs`.

### Milestone 2 — Validate

Run lint locally, spot-check endpoint slugs against the live reference, push, and drive preflight to CLEAN (see Validation and Acceptance). Fixes, if any, land as new commits.

## Concrete Steps

All commands run from the `docs` repo root.

1. Create the working branch:

       git checkout -b claude/docs-api-operations-links-9kdem5 main

2. Enumerate the API's operations to confirm the inventory table above is complete:

       grep -nE '^(  /|      summary:|      operationId:)' docs/references/platform-api/2026-05-06.yaml

   Expect 30 operations under `paths:` plus one `x-webhooks` entry (`validateDeployment`).

3. Apply the edits from Plan of Work. Follow the exemplar block style exactly (bold label, one sentence per endpoint, link text ending in ` »`, no trailing period after the component).

4. Verify coverage — every slug in the inventory table appears at least once:

       grep -rhoE 'PlatformApiLink endpoint="[a-z-]+/[a-z-]+' docs/ | sort -u

   Expect all 30 slugs from the table (28 newly added, `devices/update` and `releases/get` pre-existing). No page listed under "intentionally NO link" may appear in the diff.

5. Commit:

       git add docs/
       git commit -m "docs: link platform api endpoints from operation docs"

   Expected stat: 9 files changed (`create-a-release.mdx`, `staging-area.mdx`, `config-instances.mdx`, `config-types.mdx`, `schemas/manage.mdx`, `provisioning-tokens.mdx`, `authn.mdx`, `devices.mdx`, `releases.mdx`), insertions only except the one-line provisioning-tokens conversion.

## Test Steps

1. Lint (the repo's single local validation entrypoint):

       pnpm install --frozen-lockfile
       ./scripts/lint.sh

   Expected: exit 0 (cspell, MDX/ESLint, and prose lint all pass). MDX syntax errors in the new blocks — an unclosed component, a bad import path — fail here.

2. Endpoint slug resolution. `<PlatformApiLink>` builds URLs at render time, so a typo'd slug produces a 404 rather than a lint failure. Spot-check 8 slugs against the live reference, prioritizing the irregular ones:

       for e in config-instances/download-content principals/self deployments/drifts \
                git-commits/create devices/ping deployments/archive \
                config-types/update releases/list; do
         curl -s -o /dev/null -w "%{http_code} $e\n" \
           "https://docs.mirurobotics.com/references/platform-api/2026-05-06/endpoints/$e"
       done

   Expected: `200` for all 8.

3. Confirm the no-endpoint pages are untouched:

       git diff main --stat -- docs/primitives/groups.mdx docs/primitives/deployments.mdx

   Expected: no output.

## Validation and Acceptance

Acceptance criteria — each item must be observably true:

1. Every one of the 30 operations in the inventory table is linked from its listed page and section via `<PlatformApiLink endpoint="<slug>">`; the grep in Concrete Steps step 4 returns exactly the 30 slugs.
2. `docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx` contains no hardcoded `/references/platform-api/2026-05-06/...` URL; the provisioning-token link goes through `<PlatformApiLink>`.
3. The pages and sections listed under "intentionally NO link" carry no `**Platform API**` block and are absent from the diff.
4. `pnpm install --frozen-lockfile && ./scripts/lint.sh` exits 0.
5. The 8 spot-checked endpoint URLs return HTTP 200 on https://docs.mirurobotics.com.
6. **Preflight reports CLEAN — CI green on the pushed branch head — before the PR leaves draft or the task is reported complete.**

## Idempotence and Recovery

All edits are additive prose blocks plus one link conversion; no state outside the working tree is touched. Re-running any step is safe. To roll back, `git revert <sha>` restores every page (the conversion in provisioning-tokens.mdx reverts to the hardcoded URL, which still resolves). If lint or CI fails after the milestone commit, fix forward with a new commit — do not amend.

## Progress

- [x] Milestone 1: Add link blocks to the nine pages; verify slug coverage; commit. Landed as `cd164fd` (9 files, insertions only except the one-line provisioning-tokens conversion). Coverage grep returns exactly the 30 slugs from the inventory table.
- [ ] Milestone 2: Lint passes locally (done, exit 0) and 8 live-slug spot checks return 200 (done); preflight CLEAN on the pushed branch head still pending — branch not yet pushed.

## Surprises & Discoveries

- `PlatformApiLink` was already in use beyond the exemplar pages: `docs/developers/platform-api/query-params/*.mdx` and `docs/developers/agent/versions.mdx` link several of the same slugs (some with `#anchor` suffixes, e.g. `devices/list#query-parameters`). Repo-wide coverage greps therefore count some slugs more than once; the acceptance grep in Concrete Steps step 4 (`sort -u` on the bare slug) still yields exactly 30.
- Local `main` is behind `origin/main`; `git diff main` includes upstream commits (#125, #140) already merged into the branch's history. Verification of "no-link pages absent from the diff" must look at the work commits (`cd164fd`, `f50c42f`), which touch only the nine pages plus this plan.

## Decision Log

- Decision: Place deployment operation links on `docs/cfg-mgmt/deploy/staging-area.mdx`, not `docs/primitives/deployments.mdx`.
  Rationale: The primitives page documents properties and statuses only; the staging-area page is where stage/view/patch/archive/deploy/drift are actually described as operations, so the links sit next to the matching dashboard instructions.
  Date/Author: 2026-07-28 / planner.

- Decision: Add "View"/"Create" sections to `config-types.mdx` and `config-instances.mdx` (with `**Dashboard**` lead-ins where a section now covers both surfaces) rather than linking endpoints from unrelated prose.
  Rationale: Keeps the established section-per-operation shape used by `devices.mdx` and `releases.mdx`, so every endpoint link lives under the operation it implements.
  Date/Author: 2026-07-28 / planner.

- Decision: Exclude `validateDeployment` from the linking work.
  Rationale: It is declared under `x-webhooks:` in the spec — an event Miru sends to a customer-hosted endpoint, not an operation a client can call. The task covers callable operations only.
  Date/Author: 2026-07-28 / planner.

- Decision: Leave groups, device delete/move, release duplicate/delete, and config type delete without links.
  Rationale: The 2026-05-06 spec exposes no endpoints for these operations; adding placeholder links would 404. Revisit when the API grows the corresponding endpoints.
  Date/Author: 2026-07-28 / planner.

- Decision: Convert the provisioning-tokens hardcoded URL to `<PlatformApiLink>` as part of this task.
  Rationale: It is the one endpoint reference on a main doc page bypassing the component, so it would silently go stale on the next API version bump; the component tracks the latest version in one place.
  Date/Author: 2026-07-28 / planner.

## Outcomes & Retrospective

Implementation landed in `cd164fd` and passed the verification pass with no fixes needed:

- All 30 inventory slugs appear via `<PlatformApiLink endpoint="...">` on their mapped page and section; slugs cross-checked against the tag names and summaries in `docs/references/platform-api/2026-05-06.yaml` (including the irregular ones: `config-instances/download-content`, `principals/self`, `deployments/drifts`, `git-commits/*`).
- Block style matches the exemplar: `**Platform API**` label, one sentence per endpoint, link text ending in ` »`, no trailing period. Two plan-specified exceptions: `authn.mdx` uses a bare sentence (the page is itself Platform API docs) and `provisioning-tokens.mdx` keeps its inline "Platform API" link text through the component.
- Each touched page imports `platform-api-link.jsx` exactly once; no hardcoded `/references/platform-api/2026-05-06/endpoints/...` URLs remain in prose pages.
- New internal anchors resolve: `#view-a-deployment` and `#stage-a-deployment` (staging-area.mdx headings) and `#view-a-config-instance` (device-history.mdx heading); `/cfg-mgmt/deploy/config-editor` exists.
- No-link pages (`groups.mdx`, `deployments.mdx`, device delete/move, release duplicate/delete, config type delete) carry no blocks and are absent from the work commits.
- `pnpm install --frozen-lockfile && ./scripts/lint.sh` exits 0; all 8 live spot-check URLs return HTTP 200.

Remaining: push the branch and drive preflight to CLEAN (acceptance criterion 6).

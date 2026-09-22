# Document archiving for releases, devices, and config types

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`mirurobotics/docs`, this repo) | read-write | Add archive sections to `docs/primitives/releases.mdx`, `docs/primitives/devices.mdx`, and `docs/cfg-mgmt/primitives/config-types.mdx`; add the `archived` device status; correct the "Delete a device" section; list the new operations in `docs/admin/users/access-control.mdx`. Touch `cspell.json` only if CSpell flags a new word. |
| `backend/` (`mirurobotics/backend`) | read-only | Source of truth for archive behavior (paths cited below, commit `d4799e6`). Do not modify. |
| `openapi/` (`mirurobotics/openapi`) | read-only | Confirms which API surfaces expose archiving (`apis/configs/components/schemas/{release,device,config-type}.yaml`, `apis/common/params.yaml#archived`). Do not modify. |
| `cli-private/` | read-only | Confirmed there is no `miru` CLI command to archive or unarchive these resources. Do not modify. |

This plan lives in `docs/plans/` because every write is a documentation edit in this repo. Work happens on the already-checked-out branch `claude/hopeful-goldberg-pp1rkb` (base `main`). Do not create another branch.

## Purpose / Big Picture

Miru's dashboard lets users archive releases, devices, and config types, but the docs only cover archiving buckets (`docs/data-uploads/primitives/buckets.mdx`, `## Archive a bucket`) and deployments (`docs/cfg-mgmt/deploy/staging-area.mdx`, `## Archive a deployment`). After this change, a reader can learn from each primitive's page what archiving does, what it blocks, whether and how it can be undone, and which role can do it. The device page also says that devices with deployment history must be archived instead of deleted.

## Progress

- [ ] Milestone 1: apply Edits 1-5, pass local checks, commit.
- [ ] Milestone 2: push, open a draft PR, get preflight to `CLEAN`.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

- Decision: Add sections to the existing primitive pages instead of a new "Archiving" page. No `docs/docs.json` change.
  Rationale: Buckets and deployments document archiving as a `## Archive a …` section on their own page, next to `## Delete a …`. Following that keeps one place per resource.
  Date/Author: 2026-09-22, plan author.
- Decision: Mark every new section with `<PlatformUnsupportedBadge />` and do not mention the CLI.
  Rationale: Release, device, and config type archive endpoints exist only in the frontend (dashboard) API (`backend/internal/servers/frontend/endpoints/{releases,devices,config_types}.go`). The CLI server has device archive routes (`backend/internal/servers/cli/endpoints/entry.go:122,133`), but `cli-private` has no command that calls them. The Platform API lists return archived and unarchived rows alike (`backend/internal/servers/platform/templates/{releases,devices,config_types}.go`, `qryreqs.ArchivedAny`) and only expose the device `archived` status.
  Date/Author: 2026-09-22, plan author.
- Decision: No new screenshots. Dashboard steps are short prose that follow the existing "ellipses (...) → action → confirm" pattern.
  Rationale: The frontend repo is not available, so new UI images cannot be captured. Existing images (for example `releases/ellipses-dropdown.png`) predate archiving and would not show the **Archive** item.
  Date/Author: 2026-09-22, plan author.
- Decision: No changelog entry.
  Rationale: Archiving shipped earlier (backend `CHANGELOG.md`: "archive/unarchive for releases (#162)", "archive/unarchive for config_types (#161)"). Changelog entries are made at release time, not for documentation catch-up (precedent: `plans/completed/20260922-file-rule-deployment-behavior-docs.md`).
  Date/Author: 2026-09-22, plan author.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX (Markdown plus JSX components) under `docs/`. Navigation is in `docs/docs.json`. Role badges (`PublisherBadge`, `ProvisionerBadge`, `AdminBadge`) come from `/snippets/components/role-badges.jsx`. `PlatformUnsupportedBadge` comes from `/snippets/components/platform-api-link.jsx` and marks dashboard-only operations. All three target pages already import both.

Terms:

- **Archiving** sets a resource's `archived_at` timestamp. It never deletes data.
- A **deployment** puts one release on one device. Its **activity status** is one of `drifted`, `staged`, `queued`, `deployed`, `removing`, `archived` (see `docs/primitives/deployments.mdx`). Only `staged` and `drifted` deployments can be archived by hand (`backend/internal/configs/domain/deployments/status.go:47-49`). Deployed ones become `archived` once a newer deployment replaces them on the device.

### Verified behavior (backend `mirurobotics/backend` at `d4799e6`)

Shared:

- Every list endpoint for these resources in the frontend API takes `archived=false|true|any`, with `false` as the default. So archived items are hidden unless requested (`internal/pkg/query/request/archived_filter.go`; `internal/servers/frontend/templates/{releases,devices,config_types}.go` use `ListWithArchivedFilter`; `openapi/apis/common/params.yaml` `archived`).
- Names stay reserved. The unique indexes on device name, config type name and slug, and release version per workspace ignore `archived_at` (`tools/supabase/migrations/20250915183477_root.sql:1101,1379,1385,1391`). An archived resource's name or version cannot be reused.

Releases (`internal/configs/services/releases/`):

1. Archive and unarchive are both supported. Bulk archive or unarchive takes at most 100 releases per request (`archive.go`, `unarchive.go`, `bulk_archive.go`, `bulk_unarchive.go`; `errors.go:21` `MaxBulkArchiveReleases = 100`).
2. Precondition: a release can be archived only if every one of its deployments is `archived`. Otherwise the request fails with `release_has_unarchived_deployments`, "this release cannot be archived since it has unarchived deployments" (`bulk_archive.go` `verifyNoUnarchivedDeployments`; `errors.go:89-111`; `internal/configs/db/release.go:107-116`). The dashboard tooltip reads "Cannot archive releases that have staged or active deployments" (`internal/authz/actions/release.go:59-75`).
3. Effect: new deployments of an archived release are rejected with `deployment_release_archived`, "cannot create deployment: release '<id>' is archived" (`internal/configs/services/deployments/create/entry.go:162,249-254`; `internal/configs/domain/deployments/errors.go:22,304-319`). Existing, already-archived deployments and history are kept.
4. Unarchive has no precondition beyond the release being archived (`bulk_unarchive.go`; `internal/authz/actions/release.go:77-91`).
5. Permission: publisher (and admin/owner) roles (`internal/authz/permissions/configs/role_grants.go:48,107`).

Config types (`internal/configs/services/config_types/`):

1. Archive and unarchive are both supported, one config type at a time (`archive.go`, `unarchive.go`; bulk was removed, backend `CHANGELOG.md` #377). There is no precondition other than not already being archived (`internal/authz/actions/config_type.go:73-103`).
2. Effect: new releases that include an archived config type are rejected with `release_config_type_archived`, "cannot create release: config type '<id>' is archived" (`internal/configs/services/releases/create.go:223-232`, `errors.go:158-183`). New schemas for it are rejected with `config_schema_config_type_archived`, "cannot create config schema: config type '<id>' is archived" (`internal/configs/services/config_schemas/create.go:603-604`, `errors.go:117-142`). Existing schemas, releases, and deployments are not changed.
3. Permission: publisher (and admin/owner) (`role_grants.go:46,103`).

Devices (`internal/configs/services/devices/`):

1. Archive only. There is no unarchive endpoint. Bulk archive takes at most 100 devices (`archive.go`, `bulk_archive.go`, `errors.go:16-20`).
2. Effect on the device record: status becomes `archived` and `archived_at` is set (`bulk_archive.go` `bulkArchiveEntry`; `internal/configs/enums/devices/devices.go:11`).
3. Effect on the Miru Agent: the device's MQTT session is banned, and device-authenticated requests are rejected with HTTP 403 `archived_device_for_auth`, "device is archived" (`bulk_archive.go` `banArchivedDeviceSessions`; `internal/configs/authn/devices/authn.go:89-90`, `errors.go`; `internal/configs/authn/devicekeys/authn.go:51-52`). The device stops receiving config updates. Archiving does not change its deployments or remove files already on the machine; nothing in the archive path touches deployments or config instances.
4. Effect on deployments: deploying to an archived device is rejected with `device_is_inactive` (`internal/configs/services/deployments/create/entry.go:361-374`, `internal/configs/services/deployments/bulk_deploy.go:189`; `internal/configs/domain/devices/status.go:21-25`).
5. Device limit: archived devices do not count toward the workspace plan's device limit (`limit.go:31-35`).
6. Restoring: provisioning with the archived device's name, reprovisioning it, or activating it again restores it. The status goes to `activating`, then `online`, and `archived_at` is cleared. Restoring checks the device limit first (`provision.go:189-217,264-273`; `reprovision.go:177-182`; `activate.go:174-185,248-266`).
7. Delete vs archive: a device with any deployment record cannot be deleted. Delete fails with "device has deployments and cannot be deleted; archive the device instead" (`delete.go:68-71`, `errors.go:177`; `internal/configs/db/device.go:139-146`). The current "Delete a device" text in `docs/primitives/devices.mdx` wrongly says deleting "removes all deployment and configuration history".
8. Permission: provisioner (also group manager, admin, owner), same as delete (`role_grants.go:94`).

Not documented (internal or unverifiable): the `archived` list filter query parameter (frontend API only), the CLI server's device archive route, and exact dashboard labels for viewing archived items (frontend not available).

### Files to change

- `docs/primitives/releases.mdx`: `## Delete a release <AdminBadge />` starts at line 79.
- `docs/cfg-mgmt/primitives/config-types.mdx`: `## Delete a config type <AdminBadge />` starts at line 83.
- `docs/primitives/devices.mdx`: `## Status` table at lines 95-100 and the paragraphs after it (102-106); `## Delete a device <ProvisionerBadge />` at line 183, with its first paragraph at 187.
- `docs/admin/users/access-control.mdx`: publisher operation list ("- Edit a config type" at line 86, "- Duplicate a release" at line 110) and provisioner list ("- Delete a device" at line 140).

### Tooling

CI is `.github/workflows/ci.yml` (workflow `CI`). Jobs:

- `changes`
- `lint`: `pnpm install --frozen-lockfile`, `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`
- `audit`: `./scripts/audit.sh`
- `shell-tests`: bats
- `lint-custom-linter` and `test-custom-linter`: skipped unless `tools/lint/**` changes

`./scripts/lint.sh` runs the Go MDX linter in `tools/lint/` (sentence-case headings, no `--`, import and redirect rules), ESLint MDX, CSpell with `cspell.json` (`unarchived` is already allowed; `unarchive` may need adding), and the OpenAPI checks. `scripts/preflight.sh` runs all of these locally. GitHub access is via the GitHub MCP tools (no `gh` CLI), repo `mirurobotics/docs`. Commits must be signed; never pass `--no-gpg-sign`.

**Preflight `CLEAN`** means every CI check run on the pushed branch head SHA has completed with `success`, or `skipped` for the two custom-linter jobs.

## Plan of Work

**Edit 1: `docs/primitives/releases.mdx`.** Insert before `## Delete a release` (line 79), followed by one blank line:

    ## Archive a release  <PublisherBadge />

    <PlatformUnsupportedBadge />

    Archiving a release hides it from the releases list and prevents it from being
    deployed. The release, its config schemas, and its deployment history are kept, and
    an archived release can be unarchived at any time.

    A release can only be archived once all of its deployments are archived. Archive any
    [staged or drifted deployments](/cfg-mgmt/deploy/staging-area#archive-a-deployment)
    first. A deployment that is queued, deployed, or being removed is archived once the
    device moves to a deployment of another release.

    To archive a release, click the ellipses (...) on the release and select **Archive**,
    then confirm. To archive several releases at once (up to 100), select them and choose
    **Archive** from the bulk actions.

    ### Unarchive a release

    Unarchiving a release makes it available for deployments again. Find the release
    among the archived releases, click the ellipses (...), and select **Unarchive**.

**Edit 2: `docs/cfg-mgmt/primitives/config-types.mdx`.** Insert before `## Delete a config type` (line 83):

    ## Archive a config type  <PublisherBadge />

    <PlatformUnsupportedBadge />

    Archiving a config type retires it from future releases. While a config type is
    archived:

    - new [releases](/primitives/releases) that include it can't be created
    - new [schemas](/cfg-mgmt/primitives/schemas/manage) can't be added to it

    Existing schemas, releases, and deployments that use the config type are not
    affected, and devices keep their current configs. An archived config type can be
    unarchived at any time.

    Since most config types can't be deleted, archiving is the usual way to retire a
    config type you no longer use.

    To archive a config type, click the ellipses (...) on the config type and select
    **Archive**, then confirm.

    ### Unarchive a config type

    Unarchiving a config type lets you create schemas and releases with it again. Click
    the ellipses (...) on the archived config type and select **Unarchive**.

**Edit 3: `docs/primitives/devices.mdx`.**

(a) In the `## Status` table, add a row after `offline`:

    | `archived`   | Device has been archived and can no longer connect to Miru  |

Keep the column padding aligned. After the paragraph ending "…before transitioning back to `online`." (line 104), add:

    A device enters the `archived` state when it is [archived](#archive-a-device), and
    leaves it when it is provisioned again.

(b) Replace the first paragraph of `## Delete a device` (line 187) with:

    Deleting a device permanently removes it from Miru and immediately invalidates its
    authentication session. This action is irreversible. Please proceed with caution.

    A device can only be deleted if it has never had a deployment. Devices with
    deployment history must be [archived](#archive-a-device) instead.

(c) Insert before `## Delete a device` (line 183):

    ## Archive a device  <ProvisionerBadge />

    <PlatformUnsupportedBadge />

    Archive a device when it is retired, replaced, or no longer managed by Miru. Archiving
    a device:

    - sets its [status](#status) to `archived` and hides it from the devices list
    - immediately invalidates its authentication session, so the Miru Agent can no longer
      connect to Miru or receive config updates
    - prevents new deployments to the device
    - frees up a device slot in your plan's device limit

    Archiving keeps the device's deployment and configuration history. It does not remove
    configs that were already written to the machine. The device's name stays reserved in
    the workspace.

    To archive a device, click the ellipses (...) on the device and select **Archive**,
    then confirm. To archive several devices at once (up to 100), select them and choose
    **Archive** from the bulk actions.

    ### Restore an archived device

    Archived devices can't be unarchived directly. To bring an archived device back,
    [reprovision](/cfg-mgmt/provision-devices/reprovision) it, or
    [provision](/cfg-mgmt/provision-devices/overview) a machine with the archived
    device's name. The device returns to `activating` and then `online`, and keeps its
    deployment history. Restoring a device counts toward your plan's device limit.

**Edit 4: `docs/admin/users/access-control.mdx`.** Under publisher, add "- Archive a config type" and "- Unarchive a config type" after "- Edit a config type". Add "- Archive a release" and "- Unarchive a release" after "- Duplicate a release". Under provisioner "[Manage devices](/primitives/devices)", add "- Archive a device" after "- Delete a device".

**Edit 5 (only if CSpell fails):** add `unarchive` (and `Unarchive`/`unarchiving` if flagged) to `words` in `cspell.json`, next to `unarchived` (line 76).

Deliberately not changed: `docs/docs.json` (no new page), `docs/references/**` (generated from OpenAPI), `docs/changelog/*`, and the Platform API property lists (the Platform API does not expose `archived_at` for releases or config types).

## Concrete Steps

All commands run from the docs repo root, `/home/user/docs`.

### Milestone 1: edits, local checks, commit

1. Confirm the start state:

        git branch --show-current   # expect: claude/hopeful-goldberg-pp1rkb
        git status --short          # expect: only plans/ changes, if any

2. Re-verify the backend facts this plan cites (test step, read-only). From `/home/user/backend`:

        grep -n 'MaxBulkArchiveReleases = 100\|release_has_unarchived_deployments\|release_config_type_archived' internal/configs/services/releases/errors.go
        grep -n 'MaxBulkArchiveDevices = 100\|archive the device instead' internal/configs/services/devices/errors.go
        grep -n 'device is archived' internal/configs/authn/devices/errors.go
        grep -n 'deployment_release_archived' internal/configs/domain/deployments/errors.go
        ls internal/configs/services/devices/ | grep -c unarchive     # expect 0 (no device unarchive)

    Each grep must print at least one line. If one prints nothing, re-read the cited code, correct the affected text, and record it in Surprises & Discoveries. Optionally run the backend archive integration tests with `go test ./tests/servers/...  -run 'Archive'` (needs the backend's local database from its README; skip if unavailable). `git -C /home/user/backend status --short` must stay empty.

3. Apply Edits 1-4 from Plan of Work.

4. Content checks:

        grep -n '^## Archive a release\|^### Unarchive a release\|^## Delete a release' docs/primitives/releases.mdx
        # expect three lines, in that order
        grep -n '^## Archive a config type\|^### Unarchive a config type\|^## Delete a config type' docs/cfg-mgmt/primitives/config-types.mdx
        # expect three lines, in that order
        grep -n '^## Archive a device\|^### Restore an archived device\|^## Delete a device\|`archived`   |' docs/primitives/devices.mdx
        # expect the status row first, then the three headings in that order
        grep -c 'rchive a' docs/admin/users/access-control.mdx
        # expect 7 (bucket and deployment, plus 5 new lines)

5. Local checks (a fast mirror of the CI `lint` job):

        pnpm install --frozen-lockfile
        pnpm run test:lint     # expect exit 0
        ./scripts/lint.sh      # expect last line: All documentation lint checks passed.
        pnpm run validate      # expect exit 0 (mint validate; catches broken internal links and anchors)

    If CSpell flags a word, apply Edit 5 and rerun. If the registry is unreachable, build and run the custom linter on the changed files and rely on CI for the rest:

        cd tools/lint && go build -o lint . && cd ../..
        tools/lint/lint docs/primitives/releases.mdx docs/primitives/devices.mdx docs/cfg-mgmt/primitives/config-types.mdx docs/admin/users/access-control.mdx

    Optionally run `pnpm run dev` and open `http://localhost:3000/primitives/devices#archive-a-device` to check rendering and anchors.

6. Commit one signed commit, with the session's required trailers:

        git add docs/primitives/releases.mdx docs/primitives/devices.mdx docs/cfg-mgmt/primitives/config-types.mdx docs/admin/users/access-control.mdx plans/
        git add cspell.json   # only if changed
        git commit -m "docs: document archiving for releases, devices, and config types"

### Milestone 2: push, draft PR, preflight to CLEAN

1. `git push -u origin claude/hopeful-goldberg-pp1rkb`.
2. Load the GitHub MCP tools with `ToolSearch` (`select:mcp__github__create_pull_request,mcp__github__pull_request_read,mcp__github__actions_list,mcp__github__get_job_logs,mcp__github__update_pull_request`). Open a **draft** PR with owner `mirurobotics`, repo `docs`, head `claude/hopeful-goldberg-pp1rkb`, base `main`, and the commit subject as title. The body summarizes the edits and the verified behavior. It asks a reviewer to confirm the dashboard labels (**Archive**, **Unarchive**, bulk actions, and how archived items are shown), because the frontend was not available. It ends with the session's PR attribution lines.
3. Poll CI for `git rev-parse HEAD`. For any failed job, read the job logs, fix the cause locally, rerun the matching step 5 check, then commit (signed, new commit) and push. Never amend or force-push pushed commits. If `audit` fails on an advisory unrelated to this diff, do not widen the diff: record it and leave the PR in draft.
4. Only when preflight reports `CLEAN`, set `draft: false`.

## Validation and Acceptance

1. `/primitives/releases` has `## Archive a release` (with `### Unarchive a release`) before `## Delete a release`. From it, a reader learns that archiving hides the release and blocks new deployments, that all deployments must be archived first, that it is reversible, that up to 100 can be archived in bulk, and that the publisher role can do it.
2. `/cfg-mgmt/primitives/config-types` has `## Archive a config type` (with `### Unarchive a config type`). From it, a reader learns that archiving blocks new releases and schemas but leaves existing ones alone, and that it is reversible.
3. `/primitives/devices` lists `archived` in the status table and has `## Archive a device` (with `### Restore an archived device`). From it, a reader learns that the agent is disconnected and deployments are blocked, that there is no direct unarchive (reprovision or provision by name instead), about the device limit, and that devices with deployments must be archived rather than deleted.
4. `/admin/users/access-control` lists the five new operations under the right roles.
5. The Milestone 1 step 2 greps all match, and the step 4 greps return the stated orders and counts.
6. `git diff --stat main...HEAD` lists only the four docs files, the plan file, and optionally `cspell.json`.
7. Locally, `pnpm run test:lint`, `./scripts/lint.sh`, and `pnpm run validate` exit 0 (or the custom-linter fallback passes).
8. **Preflight reports `CLEAN`**: every CI check run on the pushed branch head SHA passed (`changes`, `lint`, `audit`, `shell-tests`; the custom-linter jobs are skipped). This must hold before the PR leaves draft and before the task is reported complete.

## Idempotence and Recovery

- Step 2 is read-only and repeatable.
- Edits are insertions plus one paragraph replacement. Check the step 4 greps before reapplying and skip any edit already present.
- Before commit, undo an edit with `git checkout -- <file>`. After commit, use `git revert <sha>`. Fix CI failures forward with new commits. The only allowed force-push is `git push --force-with-lease` after rebasing onto `origin/main` to resolve a conflict.
- The `backend`, `openapi`, and `cli-private` repos are never written to.

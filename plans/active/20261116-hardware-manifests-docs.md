# Document hardware manifests

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`mirurobotics/docs`, this repo) | read-write | New concept page `docs/cfg-mgmt/concepts/hardware-manifests.mdx` and its definition snippet, a `docs/docs.json` nav entry, and edits to `docs/admin/users/access-control.mdx`, `docs/concepts/devices/manage.mdx` and `docs/changelog/product.mdx`. `cspell.json` only if CSpell flags a word. |
| `openapi/` (`mirurobotics/openapi`) | read-only | Source of truth for manifest behavior: `apis/configs/components/schemas/hardware-manifest.yaml`, `apis/configs/components/requests/hardware-manifest.yaml`, `apis/configs/paths/hardware-manifests.yaml`, `apis/configs/paths/devices.yaml` (`/devices/set_hardware_manifest/bulk`), `apis/configs/components/requests/device.yaml`, `apis/configs/components/schemas/component.yaml`, `apis/configs/paths/components.yaml`. |

This plan lives in `docs/plans/` because every write is a docs edit. Work happens on the checked-out branch `docs/hardware-manifests` (base `main`). The PR is "docs(hardware-manifests): document hardware manifests".

## Purpose / Big Picture

Hardware manifests are new in the dashboard. A hardware manifest is the set of components a kind of device should have calibrated, for example "Picker v2" with an arm and a camera. Each device has zero or one manifest, and the dashboard shows each manifest component on the device as Not calibrated, Calibrated or Invalid. After this change a reader can learn from one concept page what a manifest is, what the statuses mean, who can create, edit, delete and set manifests, and why a manifest in use can't be deleted. The access control page lists the new Publisher operations, the device page lists the new bulk action, and the product changelog announces the feature. In the Mintlify preview, **Config Management > Concepts** shows **Hardware manifests** after **Calibrations**.

## Decisions

### Names

| New name | Alternatives | Why this one |
|----------|--------------|--------------|
| Page path `cfg-mgmt/concepts/hardware-manifests` (URL `/cfg-mgmt/concepts/hardware-manifests`) | `cfg-mgmt/concepts/manifests`; `concepts/devices/hardware-manifests` | Plural like its siblings `components` and `calibrations`; the full name matches the API resource `hardware_manifests`; it sits with the calibration concepts it builds on |
| Snippet `docs/snippets/definitions/hardware-manifest.mdx`, imported as `HardwareManifestDef` | `manifest.mdx` / `ManifestDef`; inline definition, no snippet | Singular file plus `<Name>Def`, like `component.mdx` / `ComponentDef` and `calibration.mdx` / `CalibrationDef` |
| Status labels **Not calibrated**, **Calibrated**, **Invalid**, and the **Not in manifest** list | `not_calibrated` / `calibrated` / `invalid` (code style); "Uncalibrated" | The dashboard labels named for this feature; the RFC's statuses in plain words |
| Bulk action **Set hardware manifest** | **Assign manifest**; **Set manifest** | Matches the API's `bulkSetHardwareManifest` ("set" also covers clearing) and the requested label |

### Precedent

| New file | Mirrors |
|----------|---------|
| `docs/cfg-mgmt/concepts/hardware-manifests.mdx` | `docs/cfg-mgmt/concepts/components.mdx` and `calibrations.mdx` (frontmatter, definition snippet, `ParamField` properties with field badges); operation sections with a role badge in the heading and `<PlatformUnsupportedBadge />` beneath, from `docs/cfg-mgmt/concepts/config-types.mdx` |
| `docs/snippets/definitions/hardware-manifest.mdx` | `docs/snippets/definitions/component.mdx` (one sentence, bold term, links to related concepts) |

### Left out

- Screenshots and exact click paths: the dashboard has no manifest UI yet. The page describes the UI generically. Follow-up: add images and click paths once the UI ships.
- The components page rewrite for flat components (single slug, at most 48 characters, created by Publishers): another PR documents it; this PR only links to `/cfg-mgmt/concepts/components`.
- The Device API route that marks a calibration invalid: it ships in a later agent release and another PR documents it. The page mentions the Invalid status only.
- Manifest IDs (`hw_mnf_…`), timestamps and creator fields: dashboard-only resource, and the sibling concept pages list no IDs or audit fields.
- Changing the "Edit a device" section of `docs/concepts/devices/manage.mdx`: its click path for the manifest field is unknown, and its `devices/update` Platform API badge describes the Platform route, which has no manifest field. The concept page and access control page say that editing a device sets its manifest.
- The device list's calibration summary column and a fleet filter: not built.
- A cross-link from `calibrations.mdx`: not required; the manifest page links to calibrations.
- Rewording the Publisher intro sentence in `access-control.mdx`: it already omits buckets and upload collections, which its list includes.

### Open questions

- Can Publishers delete a component? The request asks to list "Delete a component" under Publisher, but the openapi spec has only `POST /components` and `PATCH /components/{component_id}`, and `Component`'s description says "components are never archived or deleted through the API." Default: list **Create a component** and **Edit a component** only, and say so in the PR body. Add the delete line if a reviewer confirms a delete exists.
- Which example components? The request's example is "Picker v2" with `aligner1` and `gripper`; the sibling pages use `kuka-arm` and `realsense-d435`. Default: "Picker v2" with `kuka-arm` and `realsense-d435`, so the docs use one set of examples.
- How to word when invalidation is available? No agent version is known. Default: "Marking a calibration invalid requires an upcoming Miru Agent release," with no version or route.
- Is the bulk action button labeled **Set hardware manifest**? Unverified against the dashboard. Default: use it.

## Progress

- [ ] Milestone 1: concept page, definition snippet, nav entry; commit.
- [ ] Milestone 2: access control, device bulk actions, changelog; commit.
- [ ] Milestone 3 (tests): content check against the openapi spec and cross-reference check; commit any fixes.
- [ ] Preflight `CLEAN` on the pushed head.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

(Add entries as work proceeds.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The repo is the Mintlify docs site for Miru. Pages are MDX (Markdown plus JSX components) under `docs/`, served at their path without `.mdx`. Navigation is `docs/docs.json`; the Config Management > Concepts group currently ends with `"cfg-mgmt/concepts/components",` and `"cfg-mgmt/concepts/calibrations"` (around line 120). Shared JSX lives in `docs/snippets/components/`: role badges (`PublisherBadge`, `ProvisionerBadge`, …) in `role-badges.jsx`, field badges (`EditableBadge`, `ImmutableBadge`, `NullableBadge`) in `field-badges.jsx`, and `PlatformUnsupportedBadge` in `platform-api-link.jsx`, which marks operations the Platform API (Miru's public REST API) does not offer. One-sentence definitions live in `docs/snippets/definitions/` and are imported as components.

Terms. A **component** is a kind of hardware devices report calibrations for (`docs/cfg-mgmt/concepts/components.mdx`). A **calibration** is a result a device reports for a component; the most recently reported one is the device's current, or latest, calibration of it (`docs/cfg-mgmt/concepts/calibrations.mdx`, "History"). **Publisher** and **Provisioner** are member roles (`docs/admin/users/access-control.mdx`); owners and admins can do everything either role can.

Facts from openapi (read-only, `/home/user/openapi`):

- A hardware manifest has a `name`, unique in the workspace, and a list of components: at most 100, each at most once (`requests/hardware-manifest.yaml`: `uniqueItems: true`, `maxItems: 100`). Create and update take existing component IDs, so a component must exist before a manifest can list it; `POST /components` creates one.
- Update renames the manifest or replaces its components. Delete is refused while any device uses the manifest (`paths/hardware-manifests.yaml`: "A hardware manifest that any device uses can't be deleted"). The RFC's "deleting unassigns devices" is superseded.
- A device has zero or one manifest: nullable `hardware_manifest_id` on device update; null clears it. `POST /devices/set_hardware_manifest/bulk` sets or clears it on 1 to 100 devices (`requests/device.yaml`, `BulkSetHardwareManifestRequest`).
- Roles: Publishers (and admins and owners) create, edit and delete manifests and create and edit components. Setting a device's manifest is a device update (`devices:update`), so Provisioners (and admins and owners) can do it, singly or in bulk.
- Status per manifest component on a device, from its latest calibration of that component by reported time: none is Not calibrated; latest invalidated is Invalid; otherwise Calibrated. A new calibration becomes the latest and clears Invalid. Reported components outside the manifest are listed under Not in manifest.
- Only the Frontend (dashboard) API has manifests; the Platform API does not.

Lint rules that apply (`scripts/lint.sh`, Go linter in `tools/lint/linter/`): frontmatter titles and page headings in sentence case (changelog headings indented inside `<Update>` use title case, like `## Device Staging`); imports sorted by path, contiguous, all used, with `{ A, B }` spacing; `.mdx` imports end with `;`; no `--` (use an em dash). ESLint MDX and CSpell (`cspell.json`) also run.

## Plan of Work

Milestone 1, concept page.

Create `docs/snippets/definitions/hardware-manifest.mdx` with one line:

    A **hardware manifest** is the set of [components](/cfg-mgmt/concepts/components) a kind of device should have calibrated, such as a "Picker v2" manifest listing `kuka-arm` and `realsense-d435`.

Create `docs/cfg-mgmt/concepts/hardware-manifests.mdx`:

    ---
    title: "Hardware manifests"
    ---

    import { EditableBadge } from '/snippets/components/field-badges.jsx';
    import { PlatformUnsupportedBadge } from '/snippets/components/platform-api-link.jsx';
    import { ProvisionerBadge, PublisherBadge } from '/snippets/components/role-badges.jsx';
    import HardwareManifestDef from '/snippets/definitions/hardware-manifest.mdx';

    <HardwareManifestDef />

    A device has at most one hardware manifest. Miru compares the device's
    [calibrations](/cfg-mgmt/concepts/calibrations) with its manifest to show which
    components are calibrated.

    ## Properties

    <ParamField path="name" type="string">
      <EditableBadge />

      The name of the kind of device the manifest describes. Unique within the workspace.

      Example: `Picker v2`
    </ParamField>

    <ParamField path="components" type="[] Component">
      <EditableBadge />

      The [components](/cfg-mgmt/concepts/components) a device with this manifest should
      have calibrated. A manifest lists at most 100 components, each at most once.

      Examples: `kuka-arm`, `realsense-d435`
    </ParamField>

    ## Calibration status

    For each component in a device's manifest, Miru shows a status based on the device's
    latest calibration of that component:

    | Status | Meaning |
    |--------|---------|
    | Not calibrated | The device hasn't reported a calibration for the component |
    | Calibrated | The latest calibration is valid |
    | Invalid | The latest calibration was marked invalid |

    A new calibration becomes the latest, so it clears **Invalid**. Calibrations for
    components that aren't in the manifest are listed under **Not in manifest**.

    Applications on the device mark a calibration invalid, which requires an upcoming
    Miru Agent release.

    ## Create a hardware manifest  <PublisherBadge />

    <PlatformUnsupportedBadge />

    To create a hardware manifest in the dashboard, name it and choose its components. A
    manifest can only list components that already exist, so
    [create a component](/cfg-mgmt/concepts/components) first if no device has reported a
    calibration for it.

    ## Edit a hardware manifest  <PublisherBadge />

    <PlatformUnsupportedBadge />

    You can rename a hardware manifest or change its components. The change applies to
    every device that uses the manifest.

    ## Set a device's hardware manifest  <ProvisionerBadge />

    <PlatformUnsupportedBadge />

    Set a device's hardware manifest when you
    [edit the device](/concepts/devices/manage#edit-a-device), or set it on up to 100
    devices at once with [bulk actions](/concepts/devices/manage#bulk-actions). Clearing
    the manifest leaves the device without one.

    ## Delete a hardware manifest  <PublisherBadge />

    <PlatformUnsupportedBadge />

    A hardware manifest can't be deleted while any device uses it. Change or clear the
    manifest on those devices first.

Sections follow the order of use (create, edit, set, delete), so the delete section can point back to the section above it without an anchor.

In `docs/docs.json`, change `"cfg-mgmt/concepts/calibrations"` to `"cfg-mgmt/concepts/calibrations",` and add `"cfg-mgmt/concepts/hardware-manifests"` on the next line, at the same indentation.

Milestone 2, linked pages.

In `docs/admin/users/access-control.mdx`, under `#### Publisher`, insert after the `[Config schemas]` block (after `- Create a config schema`), keeping one blank line between blocks:

    [Components](/cfg-mgmt/concepts/components)

    - Create a component
    - Edit a component

    [Hardware manifests](/cfg-mgmt/concepts/hardware-manifests)

    - Create a hardware manifest
    - Edit a hardware manifest
    - Delete a hardware manifest

Under `#### Provisioner`, change `- Edit a device` (line 144) to:

    - Edit a device, including its [hardware manifest](/cfg-mgmt/concepts/hardware-manifests)

In `docs/concepts/devices/manage.mdx`, under `## Bulk actions`, append a bullet after the **Archive** bullet (which ends `See [Archive a device](#archive-a-device).`), wrapped like its siblings:

    - **Set hardware manifest** <ProvisionerBadge />: set or clear the devices'
      [hardware manifest](/cfg-mgmt/concepts/hardware-manifests), up to 100 at a time.

`ProvisionerBadge` is already imported there.

In `docs/changelog/product.mdx`, insert a new entry as the first child of `<div className="changelog-page">` (line 12), before `<Update label="October 27, 2026">`, indented like the existing entries, with a blank line after it:

    <Update label="November 16, 2026">
      ## Hardware Manifests

      Hardware manifests list the components each kind of device should have calibrated,
      so you can see which devices are missing a calibration.

      Create a manifest, such as "Picker v2" with `kuka-arm` and `realsense-d435`, and set
      it on your devices one at a time or up to 100 at once. Each component in a device's
      manifest then shows as **Calibrated**, **Not calibrated**, or **Invalid**, based on
      the device's latest calibration of it.

      Publishers can now create components in the dashboard, so a manifest can list a
      component before any device has reported a calibration for it.

      Hardware manifests work with every device that reports calibrations
      ([Miru Agent v0.12.0](/changelog/agent#v0-12-0) or later). Marking a calibration
      invalid requires an upcoming Miru Agent release.

      [Hardware manifests »](/cfg-mgmt/concepts/hardware-manifests)
    </Update>

No `<Framed>` image: there is no UI to capture yet (see Left out).

Milestone 3, tests. Docs have no unit tests. The tests are CI's lint, spellcheck and build checks (see Validation) plus a content check of the new text against the openapi spec, done by reading files. Check each claim in the new text against the facts in Context and Orientation and the openapi files named in Scope: name unique in the workspace; at most 100 components, each once; components must exist first; delete refused while in use; bulk set and clear on up to 100 devices; Publisher and Provisioner roles; status rules; no Platform API. Check that every link target in the new text exists as a page and every anchor (`#edit-a-device`, `#bulk-actions`, `#v0-12-0`) exists as a heading on its target. Fix any mismatch and record it in Surprises & Discoveries.

## Concrete Steps

All commands run from the docs repo root, `/home/user/docs`, on branch `docs/hardware-manifests`.

Milestone 1: make the edits above, then:

    git add docs/snippets/definitions/hardware-manifest.mdx docs/cfg-mgmt/concepts/hardware-manifests.mdx docs/docs.json
    git commit -m "docs(hardware-manifests): add the hardware manifests concept page"

Milestone 2: make the edits above, then:

    git add docs/admin/users/access-control.mdx docs/concepts/devices/manage.mdx docs/changelog/product.mdx
    git commit -m "docs(hardware-manifests): link manifests from roles, devices and changelog"

Milestone 3: run the read-only checks, for example:

    grep -n "maxItems\|uniqueItems\|can't be deleted" /home/user/openapi/apis/configs/components/requests/hardware-manifest.yaml /home/user/openapi/apis/configs/components/requests/device.yaml /home/user/openapi/apis/configs/paths/hardware-manifests.yaml
    ls docs/cfg-mgmt/concepts/components.mdx docs/cfg-mgmt/concepts/calibrations.mdx docs/concepts/devices/manage.mdx
    grep -n "^## Edit a device\|^## Bulk actions" docs/concepts/devices/manage.mdx
    grep -n "v0.12.0" docs/changelog/agent.mdx

Expected: `maxItems: 100` and `uniqueItems: true` on `component_ids`, `maxItems: 100` on `device_ids`, the "can't be deleted" delete description, the listed files exist, and both headings and the agent version are present. If a check led to an edit:

    git add <changed files>
    git commit -m "docs(hardware-manifests): fix claims found by the spec check"

Never run `./scripts/lint.sh`, `pnpm run validate`, `pnpm run test:lint`, `./scripts/preflight.sh` or `mint` locally; CI runs them. Let any installed git hook run on commit; never pass `--no-verify`.

## Validation and Acceptance

CI is the only validation. Workflow `ci.yml` (`CI`) runs on `pull_request`. Job `lint` runs `pnpm run test:lint`, `./scripts/lint.sh` (Go MDX prose linter, ESLint MDX, CSpell, `mint openapi-check`) and `pnpm run validate` (`mint validate`, the site build). Job `audit` runs `./scripts/audit.sh`. Jobs `lint-custom-linter` and `test-custom-linter` report success with their steps skipped, since `tools/lint/**` doesn't change. `$preflight` must report `CLEAN`, every check on the pushed head green, before the PR leaves draft. If CSpell flags a word in the new text, add it to `words` in `cspell.json` (keeping the list sorted) in a fix commit.

Acceptance in the PR's Mintlify preview: Config Management > Concepts lists **Hardware manifests** after **Calibrations**, and the page shows the definition, the two properties, the status table, and the four operation sections with their role badges and the crossed-out Platform API badge. The access control page lists Components and Hardware manifests under Publisher and the manifest note under Provisioner's "Edit a device". The Manage devices page lists **Set hardware manifest** under Bulk actions. The Product changelog's first entry is "November 16, 2026" with **Hardware Manifests**. Milestone 3's content check finds no claim that differs from the openapi spec.

## Idempotence and Recovery

All steps are file edits and safe to repeat; re-applying an edit that is already present changes nothing. To undo a milestone, `git revert` its commit. Nothing here is destructive or pushed by this plan.

# Add `2026-08-17.everglades` entry to the Platform API changelog

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Title

Add the `2026-08-17.everglades` entry to `docs/changelog/platform-api.mdx`.

## Goal

Add a new top-level changelog section for the `2026-08-17.everglades` Platform API version at the top of `docs/changelog/platform-api.mdx`, above the existing `# 2026-05-06.rainier` section. The new section follows the exact structure of the rainier entry (`# <version>` heading, italic date line, one-paragraph summary, `<PlatformApiReleaseLinks version="..." />`, `## New endpoints`, `## Breaking changes`, `## Additive changes`, `## Migration steps`, trailing `---`) and accurately captures the API delta between `platform/2026-05-06.rainier` and `platform/2026-08-17.everglades-beta.2`.

The Platform API *reference* pages (`docs/references/platform-api/2026-08-17*`), the `docs.json` reference dropdown, the `versioning.mdx` supported-versions table, the `sdks.mdx` compatibility matrix, and the product changelog entry are **out of scope** for this plan; see "Follow-ups and known dependencies".

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit a single MDX file: `docs/changelog/platform-api.mdx`. |

Working branch: `docs/platform-api-changelog-everglades` (already checked out, based on `origin/main` at `64b1715`). Base branch for the PR: `main`.

This plan lives in `plans/backlog/` because all edits are confined to the `docs` repo. Move it to `plans/active/` when work starts and to `plans/completed/` when it lands.

## Files

Single file edited:

- `docs/changelog/platform-api.mdx`

Reference files (read-only during implementation):

- `docs/snippets/components/api-dropdown.jsx`, `api-links.jsx`, `endpoint.jsx`, `separator.jsx` — already imported at the top of `platform-api.mdx`; **no new imports are needed**. (`<Steps>` / `<Step>` are Mintlify built-ins and are already used unimported in the rainier entry.)
- `docs/changelog/cli.mdx` (v0.11.0) and `docs/changelog/product.mdx` (August 14, 2026) — existing instance-slots coverage that this entry must stay consistent with.

## Context and Orientation

`docs` is a Mintlify site. The Platform API changelog is one MDX file listing releases newest-first, each section separated by `---`.

**Sources of truth used to build this entry** (all verified during planning):

1. `oasdiff` release asset `platform.changelog.stable.md` from tag `platform/2026-08-17.everglades-beta.2` — the `rainier → everglades-beta.2` machine diff. This is the authoritative delta list.
2. The release spec asset itself (`info.version: 2026-08-17.everglades`, `info.x-release-version: 2026-08-17.everglades-beta.2`).
3. `repos/backend/plans/completed/20260814-add-platform-everglades-api.md` — the backend implementation plan, whose Decision Log explains each delta.

**The `Miru-Version` header value is `2026-08-17.everglades`, not the beta tag.** The backend derives `API_VERSION` from `info.version`, which carries no beta suffix and will not change when the stable tag is cut (backend plan, "Purpose / Big Picture", fact 1).

**`<PlatformApiReleaseLinks version="2026-08-17.everglades" />`** splits on `.` and uses only the first segment, so it emits `/references/platform-api/2026-08-17` and `https://assets.mirurobotics.com/docs/openapi/platform/2026-08-17.yaml`. Neither target exists yet. See "Follow-ups and known dependencies" — this is a deliberate, recorded risk, not an oversight.

## Follow-ups and known dependencies

None of these are done by this plan, but all must be tracked:

1. **Reference page + nav.** For rainier, commit `97bd94c` (2026-05-12) landed the changelog entry, `docs/references/platform-api/2026-05-06.yaml`, the `docs.json` `"dropdown": "2026-05-06.rainier"` block, the `versioning.mdx` table row, and the `sdks.mdx` matrix row **all in one commit**. This plan intentionally splits the changelog out. Until the reference PR lands, the "API reference" and "OpenAPI spec" links inside `<PlatformApiReleaseLinks>` 404. **Do not promote this PR to `production` before the reference PR lands.**
2. **Backend release gate.** Backend Decision Log Q4: `groups:read` is not in the public API-key scope catalog until backend PR #620 merges. `GET /groups` 403s for every real customer key until then. **The everglades version must not be announced to customers before #620 is on `main`.** Re-check with `grep -n "GroupNode" internal/authz/scopes/scopes.go` on backend `main`.
3. **SDK version.** `sdks.mdx` maps API versions to Python SDK versions. The everglades-compatible SDK version is unknown at planning time — see "Content gaps".
4. **Product changelog.** `docs/changelog/product.mdx:421` shows the precedent of cross-linking a Platform API release from the product changelog. Optional follow-up.

## Plan of Work

### Milestone 1 — Insert the everglades section

Edit `docs/changelog/platform-api.mdx`. Insert the block below immediately after the import block (currently ends at line 10) and before `# 2026-05-06.rainier` (currently line 12). The block ends with `---` on its own line, mirroring how rainier ends before tetons.

After this milestone the file structure is: imports → everglades → `---` → rainier → `---` → tetons → `---` → zion.

### Milestone 2 — Preflight and CI

Run `./scripts/preflight.sh`, fix findings, re-run until exit 0. Push and confirm CI is green on the branch head.

## The changelog entry (verbatim, ready to paste)

> **Before pasting:** confirm the date line. See Decision Log Q1. `*August 17, 2026*` is written below as the planned publish date; if the `platform/2026-08-17.everglades` stable tag is cut on a different day, change this line to the actual tag date before the PR leaves draft.

```mdx
# 2026-08-17.everglades

*August 17, 2026*

The `everglades` release introduces read-only device groups and adds group membership to devices, so you can read and filter your fleet by the same hierarchy you manage in the dashboard. It also brings the Platform API up to date with three months of configuration changes: config schemas now declare instance slots instead of a single instance file path, config instances carry a slot key, and devices report an `archived` status.

<PlatformApiReleaseLinks version="2026-08-17.everglades" />

## New endpoints

<Dropdown title="Groups">
  - <GET /> `/groups` — list groups
  - <GET /> `/groups/{group_id}` — retrieve a group
</Dropdown>

## Breaking changes

<Dropdown title="Config schema instance slots">

  A config schema no longer declares a single `instance_filepath`. It declares `instance_slots` — one or more file system destinations that share the schema's validation — and an `instance_format` for the files written to them. The `instance_filepath` field is removed from every config schema response and from the `POST /config_schemas` request body.

  ```diff
    # Config Schema
      "id": "cfg_sch_123",
      "config_type_name": "Motion Control",
  -   "instance_filepath": "/srv/miru/configs/motion-control.json",
  +   "instance_format": "json",
  +   "instance_slots": [
  +     {
  +       "key": "default",
  +       "name": "Default",
  +       "filepath": "/srv/miru/configs/motion-control.json",
  +       "required": true
  +     }
  +   ],
  ```

  `instance_slots` is optional when creating a config schema. If you omit it, the server creates a single required slot with key `default` and a filepath derived from the config type slug and the schema format — **not** the path you used to pass in `instance_filepath`. To keep an existing path, send it explicitly as a slot.

  ```diff
    POST /config_schemas
    {
      "config_type_ref": { "slug": "motion-control" },
      "language": "jsonschema",
      "format": "json",
  -   "instance_filepath": "/srv/miru/configs/motion-control.json",
  +   "instance_format": "json",
  +   "instance_slots": [
  +     {
  +       "key": "default",
  +       "name": "Default",
  +       "filepath": "/srv/miru/configs/motion-control.json",
  +       "required": true
  +     }
  +   ],
      "documents": [ ... ]
    }
  ```

  Slot keys and filepaths must each be unique within a schema, and slot filepaths must be unique across every config schema in a release. Visit the [instance slots](/cfg-mgmt/primitives/schemas/overview#instance-slots) documentation for the full field reference.

</Dropdown>
<Separator />
<Dropdown title="Config type slug is immutable">

  `PATCH /config_types/{config_type_id}` no longer accepts a `slug`. A config type's slug is fixed when the config type is created. Only `name` can be updated.

  ```diff
    PATCH /config_types/{config_type_id}
    {
      "name": "Motion Control",
  -   "slug": "motion-control"
    }
  ```

</Dropdown>
<Separator />
<Dropdown title="Archived device status">

  `archived` is added to the device `status` enum. Previously, archived devices were reported as `inactive`, so a device that used to read `inactive` may now read `archived`. Update any code that branches on device status.

  ```diff
    # Device status enum
      "inactive"
      "activating"
      "online"
      "offline"
  +   "archived"
  ```

</Dropdown>
<Separator />
<Dropdown title="Deployment parent_id is nullable">

  `parent_id` on `POST /deployments` is now nullable and has three distinct meanings. It is an optimistic-concurrency token: it lets you assert what the device's current target deployment is, so a concurrent deployment cannot silently overwrite yours.

  ```diff
    POST /deployments
  - "parent_id" accepts a deployment ID or is omitted.
  + "parent_id": "dpl_123"  — assert the device's current target is dpl_123; 409 if it is not.
  + "parent_id": null       — assert the device has no current target; 409 if it has one.
  + omitted                 — no concurrency check; the server fills in the device's current target.
  ```

  Omitting the field keeps the previous behavior. Send `null` only when you intend to assert that the device has no current target deployment.

</Dropdown>

## Additive changes

<Dropdown title="Config instance slot key">

  A config instance is bound to one slot of its config schema. `slot_key` is now returned on every config instance.

  ```diff
    # Config Instance
      "id": "cfg_inst_123",
      "config_schema_id": "cfg_sch_123",
      "filepath": "/srv/miru/configs/motion-control.json",
  +   "slot_key": "default",
  ```

  When creating a config instance, `slot_key` is optional if the schema has exactly one slot, and required otherwise.

  ```diff
    POST /config_instances
    {
      "config_schema_id": "cfg_sch_123",
  +   "slot_key": "controller_2",
      "content": { ... }
    }
  ```

  Config instances can also be filtered by slot.

  ```diff
    GET /config_instances
  + GET /config_instances?slot_key=controller_2
  ```

</Dropdown>
<Separator />
<Dropdown title="Device group membership">

  Devices now report the [group](/primitives/groups) they belong to. `group_id` is always present and is `null` for unassigned devices.

  ```diff
    # Device
      "id": "dvc_123",
      "name": "My Device",
  +   "group_id": "grp_123",
  ```

  The full group object can be expanded, and devices can be filtered by group.

  ```diff
    GET /devices
  + GET /devices?expand=group
  + GET /devices?group_id=grp_123
  ```

  On `GET /groups` and `GET /groups/{group_id}`, a group's `parent` and its root-first chain of `ancestors` can be expanded.

  ```diff
    GET /groups
  + GET /groups?expand=parent
  + GET /groups?expand=ancestors
  ```

</Dropdown>
<Separator />
<Dropdown title="Device description, status filter, and status sorting">

  Devices now return a free-form `description`, which is `null` when unset.

  ```diff
    # Device
      "id": "dvc_123",
      "name": "My Device",
  +   "description": "Welder on line 3",
  ```

  `GET /devices` gained a `status` filter and `status` sort options.

  ```diff
    GET /devices
  + GET /devices?status=online
  + GET /devices?order_by=status:asc
  ```

</Dropdown>
<Separator />
<Dropdown title="Opaque schema language">

  `opaque` is added to the config schema `language` enum. An opaque schema carries metadata only and treats every config instance as valid. Visit the [opaque schema language](/cfg-mgmt/primitives/schemas/languages/opaque) documentation for details.

  ```diff
    # Schema Language enum
      "jsonschema"
      "cue"
  +   "opaque"
  ```

</Dropdown>
<Separator />
<Dropdown title="XML and plain text instance formats">

  `xml` and `text` are added to the instance content `format` enum, in both config instance responses and the `POST /config_instances` request body. Visit the [config instance file formats](/cfg-mgmt/primitives/config-instances#file-formats) documentation for which formats each schema language supports.

  ```diff
    # Instance Content format enum
      "json"
      "yaml"
      "jsonc"
  +   "xml"
  +   "text"
  ```

</Dropdown>
<Separator />
<Dropdown title="File rules on releases">

  `POST /releases` accepts an optional `file_rule_ids` array, so a release can pin the [file rules](/data-uploads/primitives/file-rules) that apply to devices running it.

  ```diff
    POST /releases
    {
      "version": "v1.0.0",
      "config_schema_ids": ["cfg_sch_123"],
  +   "file_rule_ids": ["file_rule_123"]
    }
  ```

</Dropdown>

## Migration steps

<Steps>
  <Step title="Update the Miru-Version header">
    Set `Miru-Version: 2026-08-17.everglades` on all API requests. For SDK version compatibility, visit the [Platform SDKs](/developers/platform-api/sdks) page.
  </Step>

  <Step title="Migrate from instance_filepath to instance slots">
    Replace reads of `instance_filepath` with `instance_slots[].filepath`, and replace `instance_filepath` in `POST /config_schemas` request bodies with an explicit `instance_slots` array. Omitting `instance_slots` makes the server derive a default filepath rather than reuse your previous one. See [instance slots](/cfg-mgmt/primitives/schemas/overview#instance-slots).
  </Step>

  <Step title="Stop sending slug when updating a config type">
    Remove `slug` from `PATCH /config_types/{config_type_id}` request bodies. Config type slugs are immutable.
  </Step>

  <Step title="Handle the archived device status">
    Add an `archived` branch anywhere you switch on device `status`. Devices that previously reported `inactive` because they were archived now report `archived`.
  </Step>

  <Step title="Review parent_id when creating deployments">
    No change is needed if you already omit `parent_id` or always send a deployment ID. Send `null` only to assert that the device has no current target deployment.
  </Step>

  <Step title="Optionally adopt the new additive surface">
    - Read `slot_key` on config instances, and set it when a schema declares more than one slot.
    - Read `group_id` on devices, expand `group`, and filter with `group_id=grp_123`.
    - Read `description` on devices, and use the `status` filter and `status` sort on `GET /devices`.
    - Handle the `opaque` schema language and the `xml` and `text` instance formats where needed.
    - Include `file_rule_ids` when creating releases.
  </Step>
</Steps>

---
```

## Concrete Steps

All commands run from `/home/ben/miru/workbench2/repos/docs/`.

### Setup

1. Confirm the branch and a clean tree:

       git branch --show-current   # expect: docs/platform-api-changelog-everglades
       git status --short

2. Confirm the insertion point:

       head -14 docs/changelog/platform-api.mdx

   Expect the import block through line 10 and `# 2026-05-06.rainier` on line 12.

3. **Resolve the date line before editing.** Check whether the stable tag exists:

       gh release view platform/2026-08-17.everglades -R mirurobotics/openapi --json tagName,publishedAt

   If it exists, use its publish date. If it does not, keep `*August 17, 2026*` and hold the PR in draft (Decision Log Q1).

### Milestone 1: Insert the everglades section

1. Insert the block from "The changelog entry (verbatim, ready to paste)" between the import block and `# 2026-05-06.rainier`. Do not modify any existing section.

2. Verify heading order:

       grep -n "^# " docs/changelog/platform-api.mdx

   Expected, in order: `2026-08-17.everglades`, `2026-05-06.rainier`, `2026-03-09.tetons`, `2025-10-21.zion`.

3. Confirm the release-links component:

       grep -c 'version="2026-08-17.everglades"' docs/changelog/platform-api.mdx   # expect 1

4. Confirm no unintended edits:

       git diff --stat   # expect only docs/changelog/platform-api.mdx, insertions only

### Milestone 2: Preflight and CI

1. Install dependencies if needed, then run the full preflight:

       pnpm install --frozen-lockfile
       ./scripts/preflight.sh

2. Fix any finding at its source and re-run. Do not skip checks. Only add a word to `cspell.json` if it is a genuine project-wide proper noun.

3. Push and confirm CI is green on the branch head:

       git push -u origin docs/platform-api-changelog-everglades
       gh pr checks --watch

## Test / Verification Steps

These mirror exactly what `.github/workflows/ci.yml` enforces. `./scripts/preflight.sh` runs all of them in one pass; the individual commands are listed so a failure can be isolated.

| CI job / step | Command |
|---|---|
| `lint` → Run lint smoke tests | `pnpm run test:lint` |
| `lint` → Run documentation lint (MDX prose linter, ESLint-MDX, cspell, OpenAPI) | `./scripts/lint.sh` |
| `lint` → Validate documentation build (Mintlify) | `pnpm run validate` |
| `audit` | `./scripts/audit.sh` |
| `shell-tests` | `bats pub/scripts/agent/check-miru-access_test.bats` |
| `lint-custom-linter` (only if `tools/lint/**` changed — it is not) | `LINT_FIX=0 ./tools/lint/scripts/lint.sh` |
| `test-custom-linter` (only if `tools/lint/**` changed — it is not) | `./tools/lint/scripts/covgate.sh` |

Note that `pnpm run validate` runs `mint validate` from `docs/`. This is the closest thing the repo has to a build and is the real MDX/JSX compile check. There is no separate link-checker job; `mint validate` is the link signal, and it cannot see through the template literals inside `PlatformApiReleaseLinks`, so it will **not** catch the missing `/references/platform-api/2026-08-17` target. That gap is covered by the manual step below.

Manual verification (not enforced by CI):

1. `pnpm dev`, open `/changelog/platform-api`, and confirm:
   - `2026-08-17.everglades` is the topmost entry, followed by rainier, tetons, zion.
   - The italic date renders.
   - Every `<Dropdown>` expands and collapses, and paragraphs inside them are visibly spaced.
   - Every ```diff block renders with `+` / `-` highlighting.
   - `## Migration steps` renders six numbered steps.
2. Click each in-page link and confirm it resolves:
   `/cfg-mgmt/primitives/schemas/overview#instance-slots`, `/primitives/groups`,
   `/cfg-mgmt/primitives/schemas/languages/opaque`, `/cfg-mgmt/primitives/config-instances#file-formats`,
   `/data-uploads/primitives/file-rules`, `/developers/platform-api/sdks`.
   (All six targets were verified to exist on `main` at planning time.)
3. Confirm the two `<PlatformApiReleaseLinks>` links 404 as expected, and that the reference-page follow-up is filed.

## Validation

Acceptance criteria — each must be observably true:

1. `# 2026-08-17.everglades` appears in `docs/changelog/platform-api.mdx` directly after the import block, before `# 2026-05-06.rainier`.
2. `grep -n "^# " docs/changelog/platform-api.mdx` lists everglades, rainier, tetons, zion in that order.
3. `<PlatformApiReleaseLinks version="2026-08-17.everglades" />` is present exactly once.
4. The entry contains `## New endpoints` (groups), four `## Breaking changes` dropdowns (instance slots, config type slug, archived status, deployment `parent_id`), six `## Additive changes` dropdowns (slot key, device groups, device description/status, opaque language, xml/text formats, release file rules), and a six-step `## Migration steps`.
5. The entry does **not** mention upload collections (Decision Log Q3).
6. The date line matches the actual `platform/2026-08-17.everglades` stable tag date, or the PR is still in draft pending that tag (Decision Log Q1).
7. The section ends with `---` before the rainier section; the rainier, tetons, and zion sections are byte-identical to `main`.
8. `git diff --stat` against `main` shows exactly one changed file, insertions only.
9. **Preflight reports CLEAN.** `./scripts/preflight.sh` exits 0 with no warnings, **and** CI is green on the pushed branch head (`gh pr checks` all passing). Both must hold before the PR leaves draft and before this task is reported complete.

## Idempotence and Recovery

A single text insertion. Re-running is safe: delete the inserted section or `git revert` the commit to restore the pre-change file. No external state is mutated. If preflight or CI fails after a commit, fix the cause and add a **new** commit — do not amend a pushed commit.

## Progress

- [x] Milestone 1: Insert the everglades section into `docs/changelog/platform-api.mdx`.
- [x] Milestone 3 (added — see Decision Log Q6): vendor `docs/references/platform-api/2026-08-17.yaml`.
- [x] Milestone 3: add the `2026-08-17.everglades` dropdown and repoint the `latest` redirect in `docs/docs.json`.
- [x] Milestone 3: add the `versioning.mdx` supported-versions row and bump the `Miru-Version` examples site-wide.
- [x] Milestone 3: add the `sdks.mdx` compatibility row.
- [x] Milestone 2: Preflight CLEAN (`./scripts/preflight.sh` exit 0, `pnpm run validate` success) + CI green on the pushed branch head.
- [ ] **Before leaving draft:** confirm the date line against the real `platform/2026-08-17.everglades` stable tag (Decision Log Q1 / Q7).
- [ ] **Before leaving draft:** replace the `sdks.mdx` "Not yet released" cell once an everglades Python SDK ships, and cap the rainier row's `v0.10.x+`.
- [ ] **Before leaving draft:** confirm `https://assets.mirurobotics.com/docs/openapi/platform/2026-08-17.yaml` is published (the bucket is populated outside this repo).

## Surprises & Discoveries

Pre-recorded during planning, because they contradict or extend the task framing:

- **The rainier→everglades machine diff contains zero upload-collection entries.** `platform.changelog.stable.md` (rainier → everglades-beta.2) has no `/upload_collections` sections and no `UploadCollection*` component removals. Only `platform.changelog.latest.md` (beta.1 → beta.2) shows them, because they existed only between the two betas. Corroborated by the backend plan's Decision Log Q2: "`grep -i upload api/specs/platform/v20260506.yaml` returns nothing... The backend's platform audience has never served them." Nothing was removed from any released contract.

- **The task brief's delta list was incomplete.** The machine diff also shows, all of which are covered in the entry above: `GET /devices` gains a `status` query filter and `status:asc` / `status:desc` order-by values; `POST /config_instances` accepts `slot_key` in the request body and `xml` / `text` in `content.format`; `POST /config_schemas` accepts `opaque` for `language` and both `instance_format` and `instance_slots` (optional) in the request body.

- **`instance_slots` is optional on create, but omitting it does not preserve the old path.** `CreateConfigSchemaRequest.instance_slots` description: when omitted the server creates a single required `default` slot at `/srv/miru/configs/<config_type_slug>.<json|yaml>`. A caller who previously sent a custom `instance_filepath` and now simply drops the field silently gets a *different* on-disk path. The entry calls this out explicitly in both the Breaking dropdown and the migration step, because it is the one way a caller can migrate "successfully" and still break their devices.

- **`jsonc` is not removed.** The backend plan's Decision Log Q5 records a deliberate decision to keep mapping `jsonc` rather than erroring on it, precisely because the changelog lists only `xml` and `text` as additions. The entry's format diff keeps `jsonc` as unchanged context.

## Decision Log

### Q1 — What date goes on the italic date line?

**Answer: the date the stable release tag is cut, which is also the day the docs entry lands. Plan for `*August 17, 2026*`, and verify against the real tag before the PR leaves draft.**

Evidence, all from git:

- `git tag --list "platform/*"` in `repos/openapi` with creation dates:
  `platform/2026-05-06.rainier` → **2026-05-12**; `rainier-beta.3` → 2026-05-06; `everglades-beta.1` → 2026-08-14; `everglades-beta.2` → 2026-08-16. No `platform/2026-08-17.everglades` stable tag exists.
- The rainier docs entry landed in `97bd94c`, authored **2026-05-12** — the same day the stable tag was cut — and reads `*May 12, 2026*`.
- The *plan* for rainier (`plans/completed/20260506-platform-api-changelog-rainier.md`) specified `*May 6, 2026*` (the version date). The shipped entry says May 12. The date was corrected during implementation to the actual publish date.

So the precedent is unambiguous on two points: the date is the publish/stable-tag date, not the version-string date; and the entry was published **on**, not before, the stable release. `2026-05-06.rainier` slipped six days from its version date to its stable tag, so `2026-08-17.everglades` may well slip too.

**Resolution: write `*August 17, 2026*` and hold the PR in draft until `platform/2026-08-17.everglades` exists.** If the tag lands on a different day, change the date line to that day before merging. Do not merge with an unverified date. This is the safest option that is still consistent with precedent — the alternative, dating it 2026-08-16 (today) or leaving it blank, would either misstate the release date or ship an incomplete entry.

### Q2 — How does this entry stay consistent with the existing instance-slots coverage?

**Answer: the platform entry describes only the API-contract change and links to the single canonical concept page. It duplicates nothing.**

What already exists on `main` (verified):

- `docs/cfg-mgmt/primitives/schemas/overview.mdx#instance-slots` — the canonical reference (heading at line 107, anchor referenced from line 68).
- `docs/changelog/cli.mdx` v0.11.0 (*August 12, 2026*) — one Features bullet: CLI support for the instance-slots schema annotation, linking to that same anchor.
- `docs/changelog/product.mdx` August 14, 2026 — the customer-facing feature announcement, also linking to that anchor, and stating that instance slots need CLI v0.11.0+.
- Supporting coverage in `docs/cfg-mgmt/primitives/config-instances.mdx`, `docs/cfg-mgmt/create-a-release.mdx`, and `docs/snippets/references/cli/releases/create/schema-annotations.mdx`.

On the branch question: **the `docs/instance-slots` branch is NOT merged** — `git branch --contains 9c09d2b` lists only `docs/instance-slots` and `origin/docs/instance-slots`. But its *content* is on `main`, landed through squash-merged PRs (`73d8e05` #150 "document config schema instance slots", `1a386cd` #153, `45865f0` #156). `git diff --stat main origin/docs/instance-slots -- docs/` shows a residual delta across six files, which is the pre-squash history, not unpublished content. **Do not merge or rebase `docs/instance-slots`, and do not branch from it.** Branch from `main`, which this branch already does.

Consistency rules applied to the entry:

- The three existing surfaces cover *what instance slots are* and *how the CLI uses them*. The platform entry covers only *what changed in the HTTP contract* — `instance_filepath` removed, `instance_slots` + `instance_format` added, `slot_key` on instances — and defers the concept to `#instance-slots`.
- No contradiction risk on version claims: the product entry pins instance slots to CLI v0.11.0+; the platform entry pins them to `Miru-Version: 2026-08-17.everglades`. These are different clients of the same feature, so both statements are true and neither is restated in the other's entry.
- The entry uses the same anchor (`/cfg-mgmt/primitives/schemas/overview#instance-slots`) the other two use, so all three converge on one page.

### Q3 — How is the upload-collections withdrawal framed?

**Answer: it is not mentioned at all.**

The Platform API changelog documents the contract a customer can select with `Miru-Version`. Upload collections were never in a selectable version: they appeared in `platform/2026-08-17.everglades-beta.1` and were withdrawn in `beta.2`, both pre-stable. The `rainier → everglades-beta.2` oasdiff — the diff that describes what a customer upgrading from the previous stable actually experiences — lists no upload-collection change whatsoever.

Mentioning the removal would be actively misleading in two ways: it would imply customers are losing an endpoint they could have called (they could not), and it would collide with `docs/data-uploads/primitives/upload-collections.mdx`, which correctly documents upload collections as a live *product* concept reached through the dashboard and file rules. A "removed" line in the Platform API changelog would read as a product deprecation.

No precedent exists for documenting beta-only churn: the tetons and rainier entries both describe stable-to-stable deltas only, and rainier's own reference spec was vendored from `rainier-beta.3` without any beta being mentioned in the entry.

### Q4 — `archived` device status: breaking or additive?

**Answer: breaking.** The rainier entry filed pure enum additions (`yaml`/`jsonc` formats, the `removing` deployment status) under Additive, and this entry follows that precedent for `opaque`, `xml`, and `text`. `archived` is different: backend Decision Log Q3 records that `v20260506` deliberately collapses `archived → inactive`, and `v20260817` stops collapsing. A device whose response previously read `inactive` now reads `archived`. That is a changed value on an existing field, not just a new one, and it silently breaks exhaustive branches. It belongs under Breaking with the collapse explained.

### Q5 — Scope: changelog only, no reference pages

**Answer: changelog only**, matching the rainier plan's explicit scoping ("The 'Platform API documentation' (the OpenAPI reference pages) is OUT OF SCOPE; that lands in a follow-up PR"). The consequence — the two `<PlatformApiReleaseLinks>` links 404 until the reference PR lands — is recorded in "Follow-ups and known dependencies" as a merge gate rather than silently accepted.

### Q6 — Scope expanded at implementation time: the full release-docs bundle

**Answer: Q5's "changelog only" scoping was overridden by the user before implementation. This PR lands the whole release surface in one change, exactly as `97bd94c` did for rainier.**

`97bd94c` is the template. Mirroring it produced these additions beyond the changelog entry:

- `docs/references/platform-api/2026-08-17.yaml` — the vendored reference spec.
- `docs/docs.json` — a `"dropdown": "2026-08-17.everglades"` block at the top of the Platform API Reference product, **and** the `/references/platform-api/latest/:slug*` redirect repointed to `2026-08-17`.
- `docs/snippets/components/platform-api-link.jsx` — `PlatformApiLink` repointed to the `2026-08-17` endpoint pages (`97bd94c` did the same bump for rainier).
- `docs/developers/platform-api/versioning.mdx` — a new `<SupportedBadge />` row, and the `Miru-Version` example bumped.
- `docs/developers/platform-api/sdks.mdx` — a new compatibility-matrix row (see Q8).
- `docs/developers/platform-api/authn.mdx`, `query-params/{filtering,sorting,pagination,expansions}.mdx`, and `cfg-mgmt/provision-devices/provisioning-tokens.mdx` — `Miru-Version:` curl examples bumped from `2026-05-06.rainier` to `2026-08-17.everglades`, matching `97bd94c`'s tetons→rainier sweep.

`docs/changelog/product.mdx:421` is deliberately untouched — it is a historical entry referring to the rainier release.

Consequence: the `<PlatformApiReleaseLinks>` "API reference" link now resolves (the reference pages ship in this PR). The "OpenAPI spec" link still points at `https://assets.mirurobotics.com/docs/openapi/platform/2026-08-17.yaml`, which is published from outside this repo and must be confirmed before merge.

**How the reference YAML was produced.** The vendoring pipeline was reverse-engineered from `97bd94c` and confirmed byte-for-byte: `docs/references/platform-api/2026-05-06.yaml` equals the `platform/2026-05-06.rainier-beta.3` release asset run through `api/inject_scopes.py`, plus Stainless-generated `x-codeSamples`. The everglades file was produced by copying the `platform/2026-08-17.everglades-beta.2` `platform.yaml` asset and running `api/inject_scopes.py` on it (31 operations injected, including the two new `groups` operations).

**It carries no `x-codeSamples`.** `api/pull-stainless.sh` pulls `https://app.stainless.com/api/spec/documented/miru-platform/openapi.documented.yml`, which still serves `version: 2026-05-06.rainier` — Stainless has not been regenerated for everglades. Code samples cannot be produced without a Stainless everglades run, and would be wrong if hand-written. The reference pages therefore render without Python samples until the everglades SDK ships; adding them is a follow-up, not a blocker. `mint openapi-check` and `mint validate` both pass on the file as vendored.

### Q7 — The date line is a placeholder

**Answer: `*August 17, 2026*` is written as a placeholder and MUST be re-verified before merge.**

`platform/2026-08-17.everglades` does not exist as a stable tag as of implementation (2026-08-16); only `everglades-beta.1` (2026-08-14) and `everglades-beta.2` (2026-08-16) exist. Precedent says the entry is dated the day the stable tag is cut, and rainier's stable tag slipped **six days** past its version date (`2026-05-06` version → `2026-05-12` tag → `*May 12, 2026*` in the entry). Everglades may slip the same way.

**Action required at merge time:** run `gh release view platform/2026-08-17.everglades -R mirurobotics/openapi --json tagName,publishedAt`. If the tag date is not 2026-08-17, change the date line in `docs/changelog/platform-api.mdx` to the tag date before the PR leaves draft. This is why the PR stays in draft even with CI green.

### Q8 — What goes in the `sdks.mdx` row when no SDK exists?

**Answer: the literal text `Not yet released`. No version number is invented.**

`gh release list -R mirurobotics/python-platform-sdk` shows `v0.10.0` (2026-05-13) as the newest release — that is the rainier SDK. There is no everglades SDK. The rainier row is left at `v0.10.x+`, which remains true. When the everglades SDK ships, replace `Not yet released` with its version series and cap the rainier row (e.g. `v0.10.x - v0.11.x`), mirroring how `97bd94c` capped tetons at `v0.7.x - v0.9.x`.

Consistently, the changelog's "Update the Miru-Version header" migration step omits the rainier-style SDK version sentence and links to `/developers/platform-api/sdks` instead.

## Content gaps

Facts that could not be determined from the available sources and must be resolved before merge:

1. **The stable `platform/2026-08-17.everglades` tag date.** Does not exist yet. Blocks the date line (Q1).
2. **The Python Platform SDK version for everglades.** The rainier migration step names `v0.10.0` explicitly. No equivalent is knowable today — no everglades SDK release exists in `mirurobotics/python-platform-sdk` at planning time. The entry therefore omits the version and links to `/developers/platform-api/sdks` instead. If an SDK version is published before merge, add it back in the rainier phrasing: "[Python SDK](/developers/platform-api/sdks#python) `vX.Y.Z` marks the major version rollout of the `2026-08-17.everglades` API."
3. **Whether groups are writable in a later release.** The everglades spec is read-only (`GET` only), while `docs/primitives/groups.mdx` documents create/edit/delete through the dashboard. The entry says "read-only device groups" in the summary and lists only the two `GET` endpoints; it makes no forward-looking promise about write endpoints.

## Outcomes & Retrospective

Landed as draft PR [mirurobotics/docs#158](https://github.com/mirurobotics/docs/pull/158) across four commits: the vendored spec, the changelog entry, the current-version switch, and this plan update. Preflight exits 0 and all CI checks pass (`lint`, `audit`, `shell-tests`, `changes`; the custom-linter jobs correctly skip because `tools/lint/**` is untouched).

The PR is deliberately held in draft. Three things must be true before it merges, none of which CI can check: the `platform/2026-08-17.everglades` stable tag must exist and its date must match the entry's date line; backend PR #620 (`groups:read` scope registration, now merged) must be deployed, or `GET /groups` 403s for every customer key; and `https://assets.mirurobotics.com/docs/openapi/platform/2026-08-17.yaml` must be published from outside this repo.

What the plan got right: the verbatim entry pasted in unmodified and passed the MDX prose linter, ESLint-MDX, cspell, and `mint validate` on the first try. Predicting the CI command set from `.github/workflows/ci.yml` meant no surprise failures.

What the plan under-scoped: Q5's "changelog only" was reversed before implementation began (Q6). Splitting the changelog from the reference pages would have shipped an entry whose "API reference" link 404s, which the plan itself flagged as a merge gate — landing them together removes that gate entirely and matches what `97bd94c` actually did. The plan should have read the precedent commit's full file list rather than only its changelog hunk.

The one genuinely unknowable fact remained unknowable: no everglades Python SDK exists, so `sdks.mdx` says `Not yet released` rather than inventing a version (Q8). The absent `x-codeSamples` on the vendored spec has the same root cause — Stainless has not been regenerated for everglades — and resolves the same way, when the SDK ships.

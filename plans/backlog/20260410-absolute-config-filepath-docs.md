# Absolute config filepath documentation update

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | All documentation changes live here |
| `agent` (mirurobotics/agent PR #22) | read-only | Context for agent v0.8.0 behavioral changes |
| `cli-private` (mirurobotics/cli-private PR #25) | read-only | Context for CLI v0.10.0 changes |
| `backend` (mirurobotics/backend PRs #38, #39) | read-only | Context for backend migration and Device API v0.3 |
| `openapi` (mirurobotics/openapi PR #9) | read-only | Context for OpenAPI spec description changes |

This plan lives in `docs/plans/` because all code changes are in the docs repo.

## Purpose / Big Picture

After this change, all Miru documentation accurately describes the new absolute config filepath system. A reader of the docs will understand that config instance filepaths are absolute paths on the device filesystem (e.g., `/srv/miru/configs/mobility.json`) rather than relative paths joined to a base directory. The CLI Reference supports versioned documentation (v0.9 for relative paths, v0.10 for absolute paths). All changelogs document the breaking changes in CLI v0.10.0, Agent v0.8.0, and Device API v0.3.0.

## Progress

- [ ] M1: Structural changes (docs.json versioning, new CLI v0.10 reference page)
- [ ] M2: Getting-started and learn docs updates
- [ ] M3: Changelog entries (CLI, Agent, Device API, Platform API)
- [ ] M4: Developer docs updates (GH Actions, security, versions tables)
- [ ] M5: Platform API spec and query-param docs updates
- [ ] M6: Validation — grep for stale references

## Surprises & Discoveries

Add entries as work proceeds.

## Decision Log

Add entries as work proceeds.

## Outcomes & Retrospective

Summarize at completion.

## Context and Orientation

**Framework:** Mintlify static docs site. Content is MDX files with JSX components. Navigation is defined in `docs.json` at the repo root.

**Key directories:**
- `docs/` — all page content (MDX files)
- `snippets/` — shared reusable content imported by pages
- `docs.json` — navigation structure, products, version dropdowns

**Current filepath semantics:** Config instance filepaths are described as relative to `/srv/miru/config_instances/`. Examples use leading-slash relative paths like `/v1/motion-control.json`. The default filepath for a schema is `{config-type-slug}.json`.

**New filepath semantics:** Filepaths are absolute. The new default directory is `/srv/miru/configs/` (not `config_instances`). Default filepath is `/srv/miru/configs/{config-type-slug}.json`. Examples: `/srv/miru/configs/v1/motion-control.json`, `/srv/miru/configs/safety.yaml`.

**Version changes:**
- CLI: v0.9.x → v0.10.0
- Agent: v0.7.1 → v0.8.0
- Device API: v0.2.1 → v0.3.0
- Platform API: no new version — update 2026-03-09.tetons in-place
- Legacy Platform API (2025-10-21.yaml): leave as-is

**CLI Reference versioning:** Currently a flat product with no version dropdown. Device API and Platform API already use the `dropdowns` pattern in docs.json. The CLI Reference needs the same pattern to support v0.9 and v0.10 side-by-side.

**Snippet system:** Pages import shared snippets from `snippets/`. The CLI reference page (`docs/references/cli/release-create.mdx`) imports 6 snippets from `snippets/references/cli/releases/create/`. The getting-started pages import schema example snippets from `snippets/getting-started/`.

**Changelog format:** Version heading (`# v0.X.0`), italic date (`*Month DD, YYYY*`), description paragraph, then subsections (`## Breaking changes`, `## Features`, `## Improvements`, `## Fixes`). Device API changelog uses `<Dropdown>`, `<GET />`, `<Separator />`, and `<DeviceApiReleaseLinks />` components.

## Plan of Work

### M1: Structural changes

**CLI Reference versioning in docs.json:** Restructure the CLI Reference product from flat `groups` to `dropdowns` with `v0.9` and `v0.10` entries. The v0.9 dropdown points to the existing page. The v0.10 dropdown points to a new page.

**New v0.10 CLI reference page:** Create `docs/references/cli/v0.10/release-create.mdx`. This page imports the same shared snippets as v0.9 except for `schema-annotations.mdx`, which needs a v0.10-specific version. Create `snippets/references/cli/v0.10/releases/create/schema-annotations.mdx` with absolute-path language.

**Device API v0.3 in docs.json:** Add a `v0.3.0` dropdown entry. The spec file doesn't exist yet — just wire up the docs.json reference as a placeholder.

### M2: Getting-started and learn docs

Update all `instance_filepath` annotation values and descriptive text from relative to absolute paths.

**Files to edit:**
- `snippets/getting-started/empty-cue-schemas.mdx` — 3 annotations
- `snippets/getting-started/strict-cue-schemas.mdx` — 3 annotations
- `snippets/getting-started/empty-json-schemas.mdx` — 3 annotations
- `snippets/getting-started/strict-json-schemas.mdx` — 3 annotations
- `docs/learn/config-instances.mdx` — file path property description
- `docs/learn/schemas/overview.mdx` — instance file path property description
- `docs/getting-started/quick-start/deploy-configs.mdx` — verification flow
- `docs/getting-started/quick-start/create-release.mdx` — instance file path note

### M3: Changelog entries

Add new version entries to the top of each changelog file.

**Files to edit:**
- `docs/changelogs/cli.mdx` — v0.10.0 entry
- `docs/changelogs/agent.mdx` — v0.8.0 entry
- `docs/changelogs/device-api.mdx` — v0.3.0 entry
- `docs/changelogs/platform-api.mdx` — semantic change note

### M4: Developer docs

**Files to edit:**
- `docs/developers/ci/gh-actions.mdx` — bump CLI version references from v0.9 to v0.10
- `docs/developers/agent/security.mdx` — rewrite filesystem restrictions section
- `docs/developers/agent/versions.mdx` — add v0.8.x row
- `docs/developers/device-api/versions.mdx` — add v0.3.0 row + compatibility matrix entry
- `docs/developers/cli/overview.mdx` — add version note

### M5: Platform API spec and query-param docs

**Files to edit:**
- `docs/references/platform-api/2026-03-09.yaml` — update instance_filepath descriptions and examples
- `docs/developers/platform-api/query-params/expansions.mdx` — update example value

### M6: Validation

Grep the entire docs repo for any remaining references to `/srv/miru/config_instances` (should only appear in `2025-10-21.yaml` and possibly changelogs documenting the old behavior). Grep for bare relative `instance_filepath` examples that weren't updated.

## Concrete Steps

### M1: Structural changes

**Step 1.1: Create v0.10 schema annotations snippet.**
From `docs/`: create directory and file.

    mkdir -p snippets/references/cli/v0.10/releases/create

Copy the existing snippet as a starting point:

    cp snippets/references/cli/releases/create/schema-annotations.mdx \
       snippets/references/cli/v0.10/releases/create/schema-annotations.mdx

Then edit `snippets/references/cli/v0.10/releases/create/schema-annotations.mdx`:
- Change "relative to the `/srv/miru/config_instances` directory" to "The instance file path is the absolute filesystem path where config instances for this schema are written on the device."
- Change "which deploys config instances to `/srv/miru/config_instances/{config-type-slug}.json`" to "which writes config instances to `/srv/miru/configs/{config-type-slug}.json`"
- Change examples from `/v1/mobility.json`, `/safety.yaml` to `/srv/miru/configs/v1/mobility.json`, `/srv/miru/configs/safety.yaml`

**Step 1.2: Create v0.10 CLI reference page.**
Create `docs/references/cli/v0.10/release-create.mdx` by copying `docs/references/cli/release-create.mdx`:

    mkdir -p docs/references/cli/v0.10
    cp docs/references/cli/release-create.mdx docs/references/cli/v0.10/release-create.mdx

Edit the new file to import the v0.10 annotations snippet instead of the v0.9 one:

    - import Annotations from "/snippets/references/cli/releases/create/schema-annotations.mdx";
    + import Annotations from "/snippets/references/cli/v0.10/releases/create/schema-annotations.mdx";

**Step 1.3: Update docs.json — CLI Reference.**
Replace the CLI Reference product (lines 165-174) with a dropdown structure:

```json
{
  "product": "CLI Reference",
  "dropdowns": [
    {
      "dropdown": "v0.10",
      "groups": [
        {
          "group": "Releases",
          "pages": [
            "docs/references/cli/v0.10/release-create"
          ]
        }
      ]
    },
    {
      "dropdown": "v0.9",
      "groups": [
        {
          "group": "Releases",
          "pages": [
            "docs/references/cli/release-create"
          ]
        }
      ]
    }
  ]
}
```

**Step 1.4: Update docs.json — Device API v0.3.**
Add a v0.3.0 dropdown as the first entry in the Device API Reference dropdowns array (before v0.2.1). Use a minimal placeholder structure:

```json
{
  "dropdown": "v0.3.0",
  "openapi": {
    "source": "docs/references/device-api/v0.3.0/api.yaml",
    "directory": "docs/references/device-api/v0.3.0/endpoints"
  }
}
```

**Step 1.5: Add redirect for CLI latest.**
Add a redirect entry in docs.json `redirects` array:

```json
{
  "source": "/docs/references/cli/latest/:slug*",
  "destination": "/docs/references/cli/v0.10/:slug*"
}
```

**Step 1.6: Commit M1.**
From `docs/`:

    git add docs.json docs/references/cli/v0.10/ snippets/references/cli/v0.10/
    git commit -m "docs(cli): add CLI v0.10 reference with versioned dropdowns and Device API v0.3 placeholder"

### M2: Getting-started and learn docs

**Step 2.1: Update getting-started schema snippets.**
Edit the 4 snippet files to change `instance_filepath` annotation values from relative to absolute.

In `snippets/getting-started/empty-cue-schemas.mdx`, change:
- `instance_filepath="communication.yaml"` → `instance_filepath="/srv/miru/configs/communication.yaml"`
- `instance_filepath="mobility.json"` → `instance_filepath="/srv/miru/configs/mobility.json"`
- `instance_filepath="planning.json"` → `instance_filepath="/srv/miru/configs/planning.json"`

Apply the same pattern to `strict-cue-schemas.mdx`.

In `snippets/getting-started/empty-json-schemas.mdx`, change:
- `x-miru-instance-filepath: "communication.yaml"` → `x-miru-instance-filepath: "/srv/miru/configs/communication.yaml"`
- `x-miru-instance-filepath: "mobility.json"` → `x-miru-instance-filepath: "/srv/miru/configs/mobility.json"`
- `x-miru-instance-filepath: "planning.json"` → `x-miru-instance-filepath: "/srv/miru/configs/planning.json"`

Apply the same pattern to `strict-json-schemas.mdx`.

**Step 2.2: Update learn/config-instances.mdx.**
Edit `docs/learn/config-instances.mdx` line 24: replace "The file path the config instance is deployed to relative to the `/srv/miru/config_instances` directory." with "The absolute filesystem path where the config instance is written on the device."

Replace lines 26-28 (the example sentence and examples) with:

    Examples: `/srv/miru/configs/v1/motion-control.json`, `/srv/miru/configs/safety.yaml`

**Step 2.3: Update learn/schemas/overview.mdx.**
Edit `docs/learn/schemas/overview.mdx`:
- Line 63: Replace "relative to the `/srv/miru/config_instances` directory" with "The absolute filesystem path where config instances (for this schema) are written on the device."
- Line 67: Replace "The default instance file path for a schema is `{config-type-slug}.json`, which deploys config instances to `/srv/miru/config_instances/{config-type-slug}.json`." with "The default instance file path is `/srv/miru/configs/{config-type-slug}.json`."
- Line 69: Replace examples with `/srv/miru/configs/v1/mobility.json`, `/srv/miru/configs/safety.yaml`

**Step 2.4: Update getting-started/quick-start/deploy-configs.mdx.**
Edit `docs/getting-started/quick-start/deploy-configs.mdx`:
- Line 43: Replace with "The `File Path` field shows the absolute filesystem path where the config instance is written on the device."
- Line 45: Replace with "Since the `File Path` is `/srv/miru/configs/mobility.json`, the config instance is written directly to that path on the device."
- Lines 47-56: Replace the terminal verification block. Remove the `cd` command and relative `cat`. Replace with:

```
To verify the config instance is deployed to the device's file system, open a terminal on the device and cat the file path displayed in the editor.

```bash
cat /srv/miru/configs/mobility.json
```
```

**Step 2.5: Update getting-started/quick-start/create-release.mdx.**
Edit `docs/getting-started/quick-start/create-release.mdx` lines 50-52: Replace the Note with:

```
<Note>
  The instance file path annotation is the absolute filesystem path where config instances are written on the device. Currently, JSON (`.json`) and YAML (`.yaml`, `.yml`) are supported.
</Note>
```

**Step 2.6: Commit M2.**
From `docs/`:

    git add snippets/getting-started/ docs/learn/ docs/getting-started/
    git commit -m "docs(learn): update filepath references from relative to absolute"

### M3: Changelog entries

**Step 3.1: Add CLI v0.10.0 changelog entry.**
Edit `docs/changelogs/cli.mdx`. Insert before the `# v0.9.2` heading:

```markdown
# v0.10.0

*Unreleased*

`v0.10.0` requires config schema `instance_filepath` annotations to be absolute filesystem paths, aligning the CLI with the new absolute-path model used by Agent v0.8.0 and the platform.

## Breaking changes

**Absolute instance file paths**

The `instance_filepath` annotation must now be an absolute filesystem path. Relative paths are rejected by the backend.

The default instance file path changed from `{config-type-slug}.json` to `/srv/miru/configs/{config-type-slug}.json`.

Update all schema annotations:

```diff
# JSON Schema
- x-miru-instance-filepath: "mobility.json"
+ x-miru-instance-filepath: "/srv/miru/configs/mobility.json"

# CUE
- @miru(instance_filepath="mobility.json")
+ @miru(instance_filepath="/srv/miru/configs/mobility.json")
```

---
```

**Step 3.2: Add Agent v0.8.0 changelog entry.**
Edit `docs/changelogs/agent.mdx`. Insert before the `# v0.7.1` heading:

```markdown
# v0.8.0

*Unreleased*

`v0.8.0` switches the Miru Agent to absolute config instance filepaths, replaces the staging-directory deployment model with per-file transactional writes, and upgrades to Device API v0.3.0.

[Device API v0.3.0 changelog »](/docs/changelogs/device-api#v0-3-0)

## Breaking changes

- Config instance filepaths must be absolute paths (e.g., `/srv/miru/configs/mobility.json`). Relative paths are rejected with a terminal deployment failure.
- Removed `ProtectSystem=strict` and `ProtectHome=true` from the default systemd unit to support writing config files to arbitrary absolute paths. An opt-in [`lockdown.conf.example`](https://github.com/mirurobotics/agent/blob/main/build/debian/lockdown.conf.example) drop-in is provided for operators who want to re-enable filesystem sandboxing.
- Upgraded to Device API `v0.3.0`, visit the [v0.3.0 changelog](/docs/changelogs/device-api#v0-3-0) for more details

## Features

- Config instance deployments are now written using per-file atomic writes with snapshot-based rollback. If any file fails to write, all previously written files in the deployment are restored to their prior state.
- Deployments being removed now transition through a `Removing` intermediate status before archival
- Config instance files are now deleted from disk when a deployment is removed (previously only the deployment state was updated)
- Operator-friendly error messages for filesystem permission errors (EACCES) and read-only filesystem errors (EROFS)

---
```

**Step 3.3: Add Device API v0.3.0 changelog entry.**
Edit `docs/changelogs/device-api.mdx`. Insert before the `# v0.2.1` heading. This changelog uses the `<Dropdown>`, `<GET />`, `<Separator />`, and `<DeviceApiReleaseLinks />` components already imported at the top of the file.

```markdown
# v0.3.0

*Unreleased*

The `v0.3.0` release switches config instance filepaths from relative paths to absolute filesystem paths and adds new single-resource GET endpoints.

<DeviceApiReleaseLinks version="v0.3.0" />

## New endpoints

<Dropdown title="Deployments">
  - <GET /> `/deployments/{deployment_id}` — get a deployment by ID
</Dropdown>
<Separator />
<Dropdown title="Releases">
  - <GET /> `/releases/{release_id}` — get a release by ID
</Dropdown>
<Separator />
<Dropdown title="Git commits">
  - <GET /> `/git_commits/{git_commit_id}` — get a git commit by ID
</Dropdown>

## Breaking changes

<Dropdown title="Absolute config instance filepaths">
The `filepath` field on config instances is now an absolute filesystem path (e.g., `/srv/miru/configs/mobility.json`). In previous API versions, this field contained a path relative to `/srv/miru/config_instances`.

Update any on-device code that joins filepaths to a base directory — the filepath is now usable as-is.
</Dropdown>

## Additive changes

<Dropdown title="YAML config instance format">
The `format` field on config instances now includes `yaml` as a valid value alongside `json`.
</Dropdown>

---
```

**Step 3.4: Add Platform API semantic change note.**
Edit `docs/changelogs/platform-api.mdx`. Insert a note at the top, before the first version heading. Match the existing format.

```markdown
# Config filepath change

*Unreleased*

The `instance_filepath` field on config schemas and the `filepath` field on config instances now contain absolute filesystem paths (e.g., `/srv/miru/configs/v1/motion-control.json`). Previously, these fields contained paths relative to `/srv/miru/config_instances`. No new API version is required — the field type and name are unchanged.

---
```

**Step 3.5: Commit M3.**
From `docs/`:

    git add docs/changelogs/
    git commit -m "docs(changelog): add v0.10.0 CLI, v0.8.0 Agent, v0.3.0 Device API, and Platform API entries"

### M4: Developer docs

**Step 4.1: Update GitHub Actions docs.**
Edit `docs/developers/ci/gh-actions.mdx`:
- Line 34: `version: 'v0.9'` → `version: 'v0.10'`
- Line 46: "pinning to CLI version `v0.9`" → `v0.10`; "`v0.9.x`" → "`v0.10.x`"
- Line 54: `(latest v0.9.x)` → `(latest v0.10.x)`
- Line 57: `version: 'v0.9'` → `version: 'v0.10'`
- Line 71: `version: 'v0.9.1'` → `version: 'v0.10.0'`
- Line 77: Update examples from `v0.9.1`/`v0.9` to `v0.10.0`/`v0.10`
- Line 123: `version: 'v0.9'` → `version: 'v0.10'`

**Step 4.2: Update agent security docs.**
Edit `docs/developers/agent/security.mdx`. Replace the "Filesystem restrictions" paragraph and table (lines 107-117) with:

```markdown
**Filesystem restrictions**

By default, the agent's systemd unit does not restrict filesystem access beyond standard Unix permissions. The agent writes config instance files to absolute paths specified by the platform, so it needs write access to those target directories.

For operators who want stricter confinement, the agent ships a [`lockdown.conf.example`](https://github.com/mirurobotics/agent/blob/main/build/debian/lockdown.conf.example) systemd drop-in that re-enables `ProtectSystem=strict`, `ProtectHome=true`, and configurable `ReadWritePaths`. Copy it to `/etc/systemd/system/miru.service.d/lockdown.conf` and adjust `ReadWritePaths` to include the directories your config instance filepaths target.

The agent always requires write access to these internal directories:

| Path | Purpose |
|------|---------|
| `/var/lib/miru` | Agent state / credentials (internal only) |
| `/var/log/miru` | Log files (internal only) |

The agent receives a private `/tmp` (`PrivateTmp=true`).
```

**Step 4.3: Add agent v0.8.x to versions table.**
Edit `docs/developers/agent/versions.mdx`. Add a new row at the top of the table (after the header row, before `v0.7.x`):

    | `v0.8.x`    | Unreleased   | —            | <SupportedBadge />  |

**Step 4.4: Update Device API versions page.**
Edit `docs/developers/device-api/versions.mdx`.

Add to the compatibility matrix (after header, before `v0.7.1`):

    | `v0.8.0`       | `v0.3.0`       |

Add to the supported versions table (after header, before `v0.2.1`):

    | <LinkNewTab href="/docs/references/device-api/v0.3.0">v0.3.0</LinkNewTab>    | Unreleased   | —            | <SupportedBadge /> |

**Step 4.5: Add version note to CLI overview.**
Edit `docs/developers/cli/overview.mdx`. Add after the introductory paragraph (after line 19, before the CardNewTab):

```markdown
<Info>
  This documentation covers CLI **v0.10**. For older CLI versions, visit the [CLI Reference](/docs/references/cli) version selector.
</Info>
```

**Step 4.6: Commit M4.**
From `docs/`:

    git add docs/developers/
    git commit -m "docs(developers): update GH Actions versions, agent security, and version tables"

### M5: Platform API spec and query-param docs

**Step 5.1: Update Platform API spec.**
Edit `docs/references/platform-api/2026-03-09.yaml`:

Lines 1373-1374: Replace the description "The file path to deploy the config instance relative to `/srv/miru/config_instances`. `v1/motion-control.json` would deploy to `/srv/miru/config_instances/v1/motion-control.json`." with "The absolute filesystem path where this config instance is written."

Lines 1486-1487: Replace similar description with "The absolute filesystem path where config instances for this schema are written."

Line 1588: Change example `/v1/motion-control.json` to `/srv/miru/configs/v1/motion-control.json`.

Line 2232: Change example `/v1/motion-control.json` to `/srv/miru/configs/v1/motion-control.json`.

Line 2242: Change example `/v1/localization.json` to `/srv/miru/configs/v1/localization.json`.

**Step 5.2: Update expansions example.**
Edit `docs/developers/platform-api/query-params/expansions.mdx` line 184: Change `"instance_filepath": "/v1/motion-control.json"` to `"instance_filepath": "/srv/miru/configs/v1/motion-control.json"`.

**Step 5.3: Commit M5.**
From `docs/`:

    git add docs/references/platform-api/2026-03-09.yaml docs/developers/platform-api/
    git commit -m "docs(platform-api): update instance_filepath to absolute paths in spec and examples"

### M6: Validation

**Step 6.1: Grep for stale references.**
From `docs/`:

    grep -rn "config_instances" --include="*.mdx" --include="*.yaml" .

Expected: matches only in `docs/references/platform-api/2025-10-21.yaml` (legacy, intentionally unchanged) and possibly in changelog entries that document the old behavior.

    grep -rn 'instance_filepath.*"[^/]' --include="*.mdx" --include="*.yaml" .

Expected: no matches for relative instance_filepath values outside of legacy/changelog files.

**Step 6.2: Fix any remaining stale references found in step 6.1.**

**Step 6.3: Commit M6 (only if fixes were needed).**

## Validation and Acceptance

**Test 1 — No stale relative filepath references.** From `docs/`:

    grep -rn "/srv/miru/config_instances" --include="*.mdx" .

Expected: zero matches. All MDX files should use `/srv/miru/configs/` or not reference a base directory at all.

    grep -rn "/srv/miru/config_instances" --include="*.yaml" . | grep -v "2025-10-21"

Expected: zero matches. Only the legacy Platform API spec should reference the old directory.

**Test 2 — No bare relative instance_filepath values in snippets.** From `docs/`:

    grep -rn 'instance_filepath.*"[a-z]' snippets/getting-started/

Expected: zero matches. All annotation values should start with `/srv/miru/configs/`.

**Test 3 — CLI Reference has two versions.** Verify `docs.json` contains both `v0.9` and `v0.10` dropdown entries under CLI Reference. Verify both referenced pages exist:

    test -f docs/references/cli/release-create.mdx && echo "v0.9 page exists"
    test -f docs/references/cli/v0.10/release-create.mdx && echo "v0.10 page exists"

Expected: both print their respective messages.

**Test 4 — v0.10 annotations snippet uses absolute paths.** From `docs/`:

    grep "config_instances" snippets/references/cli/v0.10/releases/create/schema-annotations.mdx

Expected: zero matches.

    grep "/srv/miru/configs/" snippets/references/cli/v0.10/releases/create/schema-annotations.mdx

Expected: at least one match.

**Test 5 — Changelog entries exist.**

    grep -c "# v0.10.0" docs/changelogs/cli.mdx
    grep -c "# v0.8.0" docs/changelogs/agent.mdx
    grep -c "# v0.3.0" docs/changelogs/device-api.mdx
    grep -c "Config filepath change" docs/changelogs/platform-api.mdx

Expected: each returns `1`.

**Test 6 — Version tables updated.**

    grep "v0.8.x" docs/developers/agent/versions.mdx
    grep "v0.3.0" docs/developers/device-api/versions.mdx
    grep "v0.8.0" docs/developers/device-api/versions.mdx

Expected: each returns at least one match.

**Test 7 — GH Actions uses v0.10.**

    grep "version: 'v0.9'" docs/developers/ci/gh-actions.mdx

Expected: zero matches. All should be `v0.10` or `v0.10.0`.

**Preflight gate:** Preflight must report `clean` before a PR is opened. This is a hard gate — do not open a PR if preflight reports any failures.

## Idempotence and Recovery

All changes are text edits to documentation files. Every step can be re-run safely — edits are idempotent (replacing specific strings with specific replacements). If a step fails partway through, re-read the file and apply the remaining edits.

The docs.json restructuring is the only structural change. If it goes wrong, `git checkout docs.json` restores the original and the step can be retried.

No database migrations, no deployments, no destructive operations.

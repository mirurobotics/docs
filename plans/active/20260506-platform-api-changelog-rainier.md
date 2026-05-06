# Add 2026-05-06.rainier entry to Platform API changelog

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Title

Add `2026-05-06.rainier` entry to the Platform API changelog.

## Goal

Add a new top-level changelog section for the `2026-05-06.rainier` Platform API version to `docs/changelogs/platform-api.mdx`. The new section sits ABOVE the existing `# 2026-03-09.tetons` section (newest first), follows the same MDX conventions (`<PlatformApiReleaseLinks>`, `<Dropdown>`, `<Separator />`, `<Steps>`, italic date subtitle), and accurately captures the API diff between `2026-03-09.tetons` and `2026-05-06.rainier`:

1. New `POST /provisioning_tokens` endpoint that replaces the removed `POST /devices/{device_id}/activation_token`.
2. Two additive enum values (`yaml`, `jsonc`) on Instance Content `format` across config-instance and deployment responses + the `POST /config_instances` request body.
3. Restoration of the `removing` deployment status / activity_status (a partial reversal of the tetons removal).
4. New `expand` query parameter on device endpoints (`current_release`, `current_deployment`), new `current_release_id` filter on `GET /devices`, and new device subschema for the expandable fields (also embedded in deployment responses).

The "Platform API documentation" (the OpenAPI reference pages) is OUT OF SCOPE; that lands in a follow-up PR. Only the changelog page changes here.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit a single MDX file: `docs/changelogs/platform-api.mdx`. |

This plan lives in `docs/plans/backlog/` because all edits are confined to the `docs` repo.

## Files

Single file edited:

- `docs/changelogs/platform-api.mdx`

Reference files (read-only during implementation):

- `docs/changelogs/device-api.mdx` — for changelog tone (the v0.2.1 entry is a small, tight example).
- `snippets/components/api-dropdown.jsx`, `snippets/components/api-links.jsx`, `snippets/components/endpoint.jsx`, `snippets/components/separator.jsx` — already imported at the top of `platform-api.mdx`; no new imports needed.

## Context and Orientation

The `docs` repo is a Mintlify documentation site. The Platform API changelog is one MDX file (`docs/changelogs/platform-api.mdx`) that lists Platform API releases newest-first. Each release section follows this structure (verified by reading the existing `2026-03-09.tetons` entry):

1. `# <version>` — H1 with the dotted version (e.g. `# 2026-03-09.tetons`).
2. `*<Month D, YYYY>*` — italic release date on its own line.
3. One short summary paragraph describing the theme of the release.
4. `<PlatformApiReleaseLinks version="<version>" />` — JSX component.
5. `## New endpoints` — list of newly added endpoints, grouped by resource in `<Dropdown title="...">` blocks separated by `<Separator />`. Endpoints are bullets formatted as `<METHOD /> ` `/path` ` — description`.
6. `## Breaking changes` — same `<Dropdown>` + `<Separator />` pattern; one Dropdown per breaking change.
7. `## Additive changes` — same pattern.
8. `## Migration steps` — `<Steps>` with `<Step title="...">` children, each describing one ordered migration step.

The existing `# 2025-10-21.zion` (initial release) entry at the bottom uses a slightly different structure (no breaking/additive sections, since it's the initial release) and ends with `---`. The `2026-03-09.tetons` entry also ends with `---` followed by the zion section.

Convention: each version section is separated by a `---` horizontal rule. The new rainier section must end with `---` before the existing `# 2026-03-09.tetons` heading.

The component `<PlatformApiReleaseLinks version="..." />` is already imported (line 8 of the file). No new imports are required.

## Plan of Work

The work is a single content-edit milestone followed by a preflight milestone.

### Milestone 1 — Insert the new rainier section

Edit `docs/changelogs/platform-api.mdx`. Insert the new section IMMEDIATELY AFTER the import block (after line 10, before `# 2026-03-09.tetons` on line 12). The new section ends with a `---` separator on its own line, mirroring how the existing tetons section ends before the zion section.

The new section structure:

- `# 2026-05-06.rainier`
- `*May 6, 2026*`
- One-paragraph summary: "The `rainier` release introduces a provisioning-token flow that replaces the activation-token endpoint, expands supported config-instance content formats, restores the `removing` deployment status, and adds expandable `current_release` / `current_deployment` fields on devices."
- `<PlatformApiReleaseLinks version="2026-05-06.rainier" />`
- `## New endpoints` — `<Dropdown title="Provisioning">` with `POST /provisioning_tokens`.
- `## Breaking changes` — `<Dropdown title="Activation tokens replaced by provisioning tokens">` describing the removal of `POST /devices/{device_id}/activation_token`, the removal of the `IssueActivationTokenRequest` and `TokenResponse` schemas, and pointing to `POST /provisioning_tokens` as the replacement.
- `## Additive changes` — three Dropdowns separated by `<Separator />`:
  1. `Content formats: yaml and jsonc` — covering the new `yaml` and `jsonc` enum values on Instance Content `format` across all listed responses and the `POST /config_instances` request body, with a note about exhaustive switches.
  2. `Deployment removing status restored` — explicitly framed as a partial reversal of the tetons removal, covering `status`, `activity_status`, and the `activity_status` query parameter on `GET /deployments`.
  3. `Device expansions: current_release and current_deployment` — covering the new `expand` query param across device endpoints, the new `current_release` and `current_deployment` enum values, the new `current_release_id` filter, and the new device subschema (also embedded in deployment responses).
- `## Migration steps` — `<Steps>` with three `<Step>`s: 1) update the `Miru-Version` header; 2) replace `activation_token` calls with `POST /provisioning_tokens`; 3) optionally adopt the additive surface (yaml/jsonc handling, `removing` status handling, device expansions / `current_release_id` filter).
- Trailing `---` separator.

After this milestone, the file's structure is: imports → rainier section → `---` → tetons section → `---` → zion section.

### Milestone 2 — Preflight

Run `./scripts/preflight.sh` from the repo root. Address any lint or cspell findings (likely few, since the new entry only uses MDX patterns and components already used elsewhere in the file). Re-run until exit code 0. The terms `provisioning`, `rainier`, `jsonc`, `yaml`, `tetons`, `removing`, `current_release_id` should already be either real words or already-spelled tokens elsewhere in the repo; if cspell flags any, prefer adding to `cspell.json` only when the term is intentional and a project-wide noun.

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs/` unless otherwise stated.

### Setup

1. Confirm working branch:

       git branch --show-current
       # expect: docs/platform-api-changelog-rainier

2. Read the current changelog to confirm insertion point:

       head -20 docs/changelogs/platform-api.mdx

   Expect imports through line 10 and `# 2026-03-09.tetons` on line 12.

### Milestone 1: Insert rainier section

1. Edit `docs/changelogs/platform-api.mdx`. Insert the new section between the import block and the `# 2026-03-09.tetons` heading.

2. Verify file integrity:

       grep -n "^# " docs/changelogs/platform-api.mdx

   Expected output, in this order:

       <line>:# 2026-05-06.rainier
       <line>:# 2026-03-09.tetons
       <line>:# 2025-10-21.zion

3. Confirm the new `<PlatformApiReleaseLinks version="2026-05-06.rainier" />` line is present:

       grep -F 'version="2026-05-06.rainier"' docs/changelogs/platform-api.mdx

### Milestone 2: Preflight

1. Run preflight:

       ./scripts/preflight.sh

   Expected: exit 0 with `clean` reported.

2. If any check fails, fix the underlying issue and rerun. Do not skip checks.

3. Confirm git working tree is clean (no preflight artifacts left behind):

       git status

## Test steps

1. Lint check:

       pnpm install   # only if node_modules is missing or stale
       pnpm lint

   The repo's `package.json` exposes `lint` (`./scripts/lint.sh`); there is no separate `pnpm build` script. Mintlify pages are not "built" in the traditional sense — `pnpm lint` plus the preflight pipeline is the correct compile/lint check for MDX changes.

2. Optionally render locally with the Mintlify dev server (not required):

       pnpm exec mint dev

   Then open the local URL, navigate to the Platform API changelog page, and confirm:

   - The new `2026-05-06.rainier` section is the topmost entry.
   - The italic date `*May 6, 2026*` renders.
   - The `<PlatformApiReleaseLinks>` component renders the API-reference / OpenAPI / changelog links for the rainier version.
   - All `<Dropdown>` blocks expand and collapse cleanly.
   - The `<Steps>` block renders three numbered steps in `## Migration steps`.
   - Section ordering: rainier → tetons → zion.

3. Visual review of the rendered file diff (`git diff docs/changelogs/platform-api.mdx`) confirms only the new section is added; no existing tetons or zion content is modified.

## Validation

Acceptance criteria — each item must be observably true:

1. The new section header `# 2026-05-06.rainier` appears in `docs/changelogs/platform-api.mdx` directly after the import block, before `# 2026-03-09.tetons`.

2. `grep -n "^# " docs/changelogs/platform-api.mdx` lists rainier, tetons, zion in that order.

3. `<PlatformApiReleaseLinks version="2026-05-06.rainier" />` is present exactly once.

4. The new section contains, at minimum:
   - `## New endpoints` Dropdown for `Provisioning` mentioning `POST /provisioning_tokens`.
   - `## Breaking changes` Dropdown explaining the `activation_token` removal and pointing to `POST /provisioning_tokens`.
   - `## Additive changes` Dropdowns for: yaml/jsonc content formats, `removing` status restoration, and device expansions (`current_release`, `current_deployment`, `current_release_id` filter).
   - `## Migration steps` `<Steps>` with three steps: header update, activation_token migration, and optional adoption of additive surface.

5. The italic date is `*May 6, 2026*`.

6. The new section ends with a `---` separator before the existing tetons section.

7. The historical tetons section is unmodified.

8. **Preflight reports `clean`**: `./scripts/preflight.sh` exits 0 with no warnings. **Preflight must report `clean` before changes are published.**

9. `pnpm lint` exits 0 (covered by preflight, but called out separately as the lint signal for Mintlify MDX).

## Idempotence and Recovery

This is a single text-edit milestone. Re-running is safe: revert the commit (`git revert <sha>`) or remove the inserted section to restore the file to its pre-change state. No external state is mutated.

If preflight fails after a commit, fix the underlying issue and add a NEW commit (do not amend).

## Progress

- [x] Milestone 1: Insert rainier section into `docs/changelogs/platform-api.mdx`.
- [ ] Milestone 2: Run preflight; address any findings; confirm `clean`.

## Surprises & Discoveries

(none)

## Decision Log

- Decision: Place `current_release_id` filter under "Device expansions" rather than its own Dropdown.
  Rationale: The filter is conceptually paired with the new expandable device subresources; combining them keeps the changelog terse.
  Date/Author: 2026-05-06 / planner.

- Decision: Frame the `removing` status restoration as a "partial reversal" with an explicit reference to the tetons change rather than a standalone additive note.
  Rationale: Readers who already migrated past tetons will notice the reversal; calling it out prevents confusion.
  Date/Author: 2026-05-06 / planner.

- Decision: Do NOT modify the historical `2026-03-09.tetons` entry, even though it still lists the now-removed `POST /devices/{id}/activation_token` endpoint and the now-restored removal of `removing`.
  Rationale: Changelog entries are historical artifacts. Modifying past entries would mislead readers about what was true at that release. The new rainier entry calls out the deltas explicitly.
  Date/Author: 2026-05-06 / planner.

## Outcomes & Retrospective

(populate after implementation)

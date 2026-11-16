# Describe flat components and 48-character slugs in the calibration docs

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit `docs/cfg-mgmt/concepts/components.mdx`, `docs/cfg-mgmt/concepts/calibrations.mdx` and `docs/cfg-mgmt/report-calibrations.mdx`. |
| `backend/` | read-only | Source of truth for flat components and slug joining (`tools/supabase/migrations/20261113120000_flatten_components.sql`, `internal/configs/domain/components/slug.go`, `internal/configs/services/components/find_or_create.go`). |

This plan lives in `docs/plans/` because every edit is a docs page. Work on branch `docs/flat-components` (from `origin/main`); the PR targets `main` with title `docs(calibrations): describe flat components and 48-character slugs`.

Out of scope (other PRs): Device API v0.3, calibration invalidation, hardware manifests, the Platform API changelog, every page under `docs/changelog/`, and the vendored API references under `docs/references/`.

## Purpose / Big Picture

Components are flat now: roles and aspects no longer exist as resources, and a component is one slug. Every slug may now be 48 characters, up from 32. The docs still describe roles and aspects as parts of a component and cap slugs at 32. After this change a reader of the three calibration pages learns that each unit or separately calibrated part gets its own component slug (such as `arm-left` or `realsense-d435-intrinsics`), that the optional `role` and `aspect` parameters still work but are joined into the component slug, that a device keeps one history per component, and that slugs may be 48 characters. No page links to the removed `#param-roles` or `#param-aspects` anchors.

## Decisions

### Names

| New name | Alternatives | Why this one |
|----------|--------------|--------------|
| Heading `### Multiple units and parts` in `report-calibrations.mdx` (anchor `#multiple-units-and-parts`) | `### Units and parts`, `### Roles and aspects` | Names both cases the section covers; "Roles and aspects" would name the legacy parameters instead of the recommended practice. |

### Left out

- Changelog entries for the flattening and the 48-character limit: owned by other PRs.
- Updating `docs/references/` or the `/calibrations/latest` wording in `docs/changelog/device-api.mdx`: vendored or out of scope.
- Docs for creating a component through `POST /components`: there is no dashboard UI or customer-facing Platform API release yet, so "Create a component" keeps its `PlatformUnsupportedBadge`.
- A claim that a joined slug over 48 characters is rejected: the backend does not validate the joined length, so the docs only advise staying within 48.

### Open questions

- Should the required `component` parameter's example change from `kuka-arm` to a joined slug? Default: no; the new subsection already shows joined slugs.

## Progress

- [ ] Milestone 1: edit the three pages, grep shows no role or aspect anchors, commit.
- [ ] Milestone 2: push, open draft PR, `$preflight` reports `CLEAN`.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

(Add entries as work proceeds.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The docs repo is a Mintlify site. Pages are `.mdx` files under `docs/`; a page's URL is its path without `docs/` and the extension, so `docs/cfg-mgmt/report-calibrations.mdx` is `/cfg-mgmt/report-calibrations`. A heading gets an anchor from its text (`### Multiple units and parts` is `#multiple-units-and-parts`), and `<ParamField path="slug">` gets `#param-slug`. `<ParamField>` is built in; badges such as `ImmutableBadge`, `NullableBadge` and `PlatformUnsupportedBadge` are imported at the top of each page.

CI lint (`scripts/lint.sh`, rules in `tools/lint/linter/`) enforces sentence-case headings (`headingcase`), no `--` in prose (`nodoubledash`), and that every import is used (`importused`). CSpell checks words against `cspell.json`.

Terms: a **component** is a kind of hardware a device reports calibrations for, named by a **slug** (lowercase letters, digits and hyphens, starting and ending with a letter or digit). A **calibration** is a result the device reports through the Miru Agent's Device API `POST /calibrations` ("create calibration") endpoint.

Facts (verified against `backend/`):

1. Roles and aspects no longer exist as resources; a component is one slug.
2. Device API v0.2 create calibration still accepts optional `role` and `aspect`. Miru joins component, role and aspect with `-`, skipping absent ones, and finds or creates the component with that slug: `arm` + role `left` is `arm-left`; `optic-flipper` + aspect `pose` is `optic-flipper-pose`; `arm` + `left` + `pose` is `arm-left-pose`.
3. Component, role and aspect are each validated as a slug of 1 to 48 characters. The joined slug's length is not validated, so it can exceed 48.
4. A device keeps a separate history per component.

Current state of the files:

- `docs/cfg-mgmt/concepts/components.mdx`: `slug` ParamField says "between 1 and 32 characters"; `roles` and `aspects` ParamFields follow it; "Create a component" says "Components, roles and aspects cannot be created directly."
- `docs/cfg-mgmt/concepts/calibrations.mdx`: `component` ParamField (examples `kuka-arm`, `realsense-d435`), then `role` and `aspect` ParamFields linking `#param-roles` and `#param-aspects`; History says "a separate history for each component, role and aspect".
- `docs/cfg-mgmt/report-calibrations.mdx`: `## Components` has `### Roles` and `### Aspects` subsections and a closing paragraph on optional roles and aspects; Optional parameters has `role` and `aspect` ParamFields linking the dead anchors.

## Plan of Work

Milestone 1. In `docs/cfg-mgmt/concepts/components.mdx`, change "between 1 and 32 characters" to "between 1 and 48 characters". After the slug examples, add one paragraph: "A calibration reported with a `role` or `aspect` lands on the component whose slug joins them with `-`, such as `arm-left`. See [multiple units and parts](/cfg-mgmt/report-calibrations#multiple-units-and-parts)." Delete the `roles` and `aspects` ParamFields and the extra blank line after them. Change the Create a component text to "Components cannot be created directly. Instead, they are created when a device first [reports a calibration](/cfg-mgmt/report-calibrations) that uses them."

In `docs/cfg-mgmt/concepts/calibrations.mdx`, delete the `role` and `aspect` ParamFields, and change the `component` examples to `kuka-arm`, `realsense-d435`, `arm-left`. Change the History paragraph to "A device keeps a separate history for each component it reports. The most recently reported calibration in a history is that component's current calibration." Keep the `NullableBadge` import; other fields still use it.

In `docs/cfg-mgmt/report-calibrations.mdx`, keep the two `## Components` paragraphs (drop their trailing spaces) and replace everything from `### Roles` through "Miru creates a role or aspect the first time a calibration reports it." with:

    ### Multiple units and parts

    Devices often have several units of the same hardware, such as two arms, or calibrate parts of a unit separately, such as a camera's intrinsics and extrinsics. Give each its own component slug, such as `arm-left`, `arm-right`, or `realsense-d435-intrinsics`. A device keeps a separate history for each component it reports.

    For older clients, create calibration in Device API v0.2 still accepts the optional `role` and `aspect` parameters. Miru joins the component, role and aspect with `-`, skipping any that are unset, and reports the calibration to that component. For example, component `arm` with role `left` becomes `arm-left`, and adding aspect `pose` makes `arm-left-pose`.

    We recommend reporting the full slug in `component` and leaving `role` and `aspect` unset. Keep a joined slug within 48 characters so it fits the component [slug](/cfg-mgmt/concepts/components#param-slug) limit.

In Optional parameters, replace the `role` body with "Joined into the component slug after `component`. We recommend putting the full slug in `component` instead. See [multiple units and parts](#multiple-units-and-parts)." and example `left`. Replace the `aspect` body with "Joined into the component slug after `component` and `role`. We recommend putting the full slug in `component` instead. See [multiple units and parts](#multiple-units-and-parts)." keeping example `pose`.

Milestone 2. Push the branch, open a draft PR with `$pr`, and run `$preflight`.

## Concrete Steps

All commands run from the docs repo root (`docs/`).

Milestone 1:

    git mv plans/backlog/20261116-flat-components-docs.md plans/active/20261116-flat-components-docs.md
    # make the edits in Plan of Work
    grep -rnE "param-roles|param-aspects|report-calibrations#(roles|aspects)" docs
    # expected: no output, exit status 1
    grep -rn "32 characters" docs/cfg-mgmt
    # expected: no output
    git add -A
    git commit -m "docs(calibrations): describe flat components and 48-character slugs"

Milestone 2:

    git push -u origin docs/flat-components
    # open the draft PR via $pr, then run $preflight until it reports CLEAN

## Validation and Acceptance

Tests for this docs change are the CI checks plus the repo-wide grep. The grep in Milestone 1 must print nothing: no page links to `#param-roles`, `#param-aspects`, or the removed `#roles` and `#aspects` headings. CI workflow `.github/workflows/ci.yml` job `lint` must pass on the pushed head; it runs the lint smoke tests, `scripts/lint.sh` (heading case, imports, double dashes, CSpell) and `pnpm run validate` (`mint validate`, the Mintlify build check). The `audit` job must also pass. Because no CI job checks anchors, the grep is the test that no dead role or aspect link remains. `$preflight` must report `CLEAN` (CI green on the pushed branch head) before the PR leaves draft. No lint, build or test runs locally.

Acceptance in the preview: `/cfg-mgmt/concepts/components` lists only `name` and `slug`, with a 48-character limit and a link that jumps to "Multiple units and parts" on `/cfg-mgmt/report-calibrations`; `/cfg-mgmt/concepts/calibrations` has no `role` or `aspect` field and says history is per component; `/cfg-mgmt/report-calibrations` has no Roles or Aspects headings and its `role` and `aspect` parameters say they are joined into the component slug.

## Idempotence and Recovery

Every step is a text edit or a read-only grep and can be repeated. If CI fails, read `gh run view --log-failed`, fix the page, commit and push again. To undo, `git revert` the milestone commit.

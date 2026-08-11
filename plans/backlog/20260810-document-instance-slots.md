# Document config schema instance slots

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `repos/docs/` | read-write | All customer-facing documentation changes, plus this plan file. |
| `repos/core/` | read-only | Source of truth for the schema-annotation syntax and parse-time errors. |
| `repos/cli-private/` | read-only | Confirms the CLI gains no new commands, flags, or output. |
| `repos/openapi/` | read-only | Source of truth for the slot object fields and API behavior. |
| `repos/backend/` | read-only | Source of truth for server-side slot rules (uniqueness, cardinality, digest, and the slot `key` character pattern). Also holds `origin/refactor/config-schema-instance-slots-frontend`. |
| `repos/frontend/` | read-only | Dashboard UI. As of 2026-08-10 it contains **no** instance-slot code and no slots branch. |

This plan lives in `repos/docs/plans/` because every file that changes is in `repos/docs`. No code changes are made in any other repository.

**Where the frontend slots branch lives.** Three edits in this plan are gated on dashboard behaviour. The branch named for that work is `origin/refactor/config-schema-instance-slots-frontend` and it lives in **`repos/backend`**, not `repos/frontend`. Confirm and inspect it read-only with:

    cd /home/ben/miru/workbench2/repos/backend && git fetch origin
    git ls-tree -r --name-only origin/refactor/config-schema-instance-slots-frontend | grep -i slot

Be aware of what that branch is and is not: it contains the frontend-facing **v08 API models** (for example `internal/servers/frontend/v08/models/model_instance_slot.go`), not dashboard markup. It can tell you that the API exposes slots; it **cannot** tell you what label the dashboard renders. `repos/frontend` `origin/main` has no slot code at all — `grep -rli instanceSlot` returns nothing. Therefore any edit in this plan that depends on a visible dashboard string is verifiable **only** against a running dashboard. If you do not have one, take the documented skip path; that is the expected outcome, not a failure.

## Purpose / Big Picture

Today the Miru docs say a config schema writes exactly **one** file, at one path, on the device. That is no longer true: a schema can declare several **instance slots**, and each slot is a distinct filesystem destination for the same schema. A robot with four identical motor controllers can now use one `motion-control` schema and four slots, instead of four near-identical config types.

After this change a reader can:

- Open <https://docs.mirurobotics.com/cfg-mgmt/primitives/schemas/overview> and see an `instance slots` property that explains what a slot is, what its five fields mean, and how it relates to the older `instance file path`.
- Open the same page's **Schema annotations** section and copy a working slots annotation in whichever of the three schema languages they use — JSON Schema, Opaque, or CUE.
- Open <https://docs.mirurobotics.com/cfg-mgmt/primitives/config-instances> and understand that an instance is bound to exactly one slot, and that its file path comes from that slot.
- Open <https://docs.mirurobotics.com/changelog/cli> and see the CLI version in which slot annotations became available.

Nothing in the CLI's command surface changes, so no new commands or flags are documented. The visible artifact is documentation only.

## Progress

- [ ] Milestone 0 — Orientation and baseline (branch, install, baseline preflight, verify merge status of the CUE and CLI work).
- [ ] Milestone 1 — Extend the shared schema-annotations snippet with the `instance slots` annotation.
- [ ] Milestone 2 — Add the `instance slots` and `slot key` properties, and correct the singular-file prose on the schema and config-instance primitives pages.
- [ ] Milestone 3 — Add slot examples and caveats to the three schema-language pages and the CUE packages snippet.
- [ ] Milestone 4 — Fix downstream drift: release creation, quick start, deploy pages, deployment constraints.
- [ ] Milestone 5 — CLI reference cross-link and a new CLI changelog entry.
- [ ] Milestone 6 — Validation, render check, push, PR, CI green, leave draft.

Add timestamps (`YYYY-MM-DD HH:MMZ`) as steps complete. Split any partially finished milestone into "done" and "remaining".

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

Entries below were made during authoring on 2026-08-10 by Ben Smidt. Add new dated entries as work proceeds.

- **2026-08-10 — Write all three language tabs as plain documentation; gate the *ship*, not the *prose*.** JSON Schema and Opaque slots are merged to `core` `main` (`f1be215` #151, `9d9845d` #150); CUE is on `core`'s unmerged `origin/feat/instance-slots-cue` (#152) and the CLI wiring on **`cli-private`**'s unmerged `origin/feat/instance-slots-annotation` (head `cc4a01b`). Beware: `core` has a *different* branch of the same name (head `0600dd0`) — finding it there does not mean the CLI wiring has landed. A two-tab `<CodeGroup>` would read as "CUE cannot do this" (false), and "coming soon" phrasing is against house style, so exposure is controlled by keeping the PR in **draft** until #152 and the CLI wiring are merged and tagged (see "Publishing is not automatic" in Context and Orientation). The Milestone 5 changelog entry is pinned to that exact version, and the PR body must state the dependency.
- **2026-08-10 — State rules in prose; document no error codes or verbatim error messages.** `repos/docs` has no precedent for error identifiers such as `invalid_instance_slots`; see the house-style note in Context and Orientation.
- **2026-08-10 — No new page; extend existing pages.** A slot is a property of a config schema, not a new primitive. Adding a page would require editing `docs/docs.json` navigation and re-running the redirects lint rule for no reader benefit. The two prior comparable plans (`plans/active/20260805-document-all-schema-annotations.md`, `plans/completed/20260805-config-instance-file-formats.md`) both extended existing pages. `docs/docs.json` is therefore **not** touched by this plan.
- **2026-08-10 — Leave the opaque `language:` vs `schema_language:` spelling alone.** `core` `main` renamed the key (`1756cc2` #134) and kept `language` as `DeprecatedSchemaLangKey` in `repos/core/pkg/schemas/opaque/compile.go:18`, so both spellings still parse. The docs currently show `language:` in three places (`docs/snippets/references/cli/releases/create/schema-annotations.mdx:33`, `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx:21` and `:28`). Changing that is independent drift with its own release-timing question and would enlarge this PR's review surface. It is recorded under Known drift in Context and Orientation, not fixed here. New opaque slot examples in this plan therefore do **not** restate the language key at all.
- **2026-08-10 — `instance file path` is not deprecated.** It is removed from the *new API contract*, but it remains a current, supported *annotation*: the single-slot shorthand. Docs must present it that way. Never write "deprecated" next to it.

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

### What this repository is

`repos/docs` is the Mintlify-rendered public documentation site for Miru. Content lives under `repos/docs/docs/`. Mintlify serves `docs/foo/bar.mdx` at the URL `/foo/bar` — there is no `/docs` prefix in URLs. Files are MDX: Markdown plus JSX components and `import` statements.

There is no `CLAUDE.md`, `AGENTS.md`, or written style guide in this repo. Conventions are enforced entirely by tooling (see Validation and Acceptance) and by imitation of neighboring files.

Publishing is **not** automatic. Merging to `main` does not put anything on the public site; a human runs the manual `promote.yml` GitHub Actions workflow to promote `main` → staging → uat → production. Because of this, `main` may legitimately be ahead of the released CLI: merging a docs PR never exposes unreleased behaviour to readers.

**House style: no error codes.** State the *rule* in prose, never a machine-readable error identifier or a verbatim error message. The only precedents in this repo are HTTP status tables in `docs/developers/device-api/events.mdx` and a status-enum table in `docs/primitives/deployments.mdx`. Follow `cue-packages.mdx`'s model: "Annotating multiple files or no files will result in an error."

### What an instance slot is

Every config schema has instance slots. Today every schema has exactly one. This feature relaxes that from N=1 to N≥1. It is **not** a mode, a flag, or a parallel code path — there is no "multi-instance mode", no `is_dynamic` boolean. Most schemas will continue to have one slot and their authors will never write the annotation.

**A slot is a distinct filesystem destination for the same schema.** All slots on a schema share that schema's validation and differ only in *where the file is written*. If two files need different validation, they are two different config types, not two slots of one schema.

A **config instance** is bound to exactly one slot. The instance's file path *is* the bound slot's file path.

A slot has five fields, identical across the API and all three schema languages:

| Field | Type | Required | Meaning |
|---|---|---|---|
| `key` | string | yes | Immutable, code-friendly identifier, unique within the schema. Lowercase letters, digits, hyphens, and underscores; must start with a lowercase letter or digit; at most 128 characters. |
| `name` | string | yes | Human-readable display name. |
| `filepath` | string | yes | Absolute path on the device. Unique within the schema **and** across every schema in the release. |
| `required` | boolean | yes | Whether every deployment of a release containing this schema must include an instance for this slot. |
| `description` | string | no | Free text. |

`required` is **mandatory on every surface**. Some plan documents in `repos/core` and `repos/cli-private` still describe it as optional-defaulting-to-true; those are stale and the code requires it. Do not document a default.

**The `key` character rules are enforced by the platform, not by the CLI.** `core` deliberately does not enforce the pattern client side. On `core` `origin/main` the merged JSON Schema and Opaque surfaces require only that `key` be present and non-empty, and each pins that with a test named "instance slot key vocabulary is not enforced client side" (`pkg/schemas/jsonschema/annotations_test.go:453`, `pkg/schemas/opaque/compile_test.go:534`). The commit that recorded the decision — `0600dd0`, "refactor(schemas): stop enforcing the slot key pattern client side" — sits on `core`'s **unmerged** `origin/feat/instance-slots-annotation`, not on `origin/main`; do not go looking for it on `main`. `^[a-z0-9][a-z0-9_-]*$` and the 128-character limit live in `repos/openapi` (`apis/configs/components/schemas/config-schema.yaml:91`) and `repos/backend` (`internal/configs/ugc/config_schemas.go:168`), and are applied when the release is created server side. Document the rules as constraints on a valid slot key — do **not** imply that `miru release create` rejects a bad key locally, and do not describe them as parse-time or annotation-time validation.

If a schema declares no slots, the server synthesizes one: `key: default`, `name: Default`, `filepath: /srv/miru/configs/{config-type-slug}.{json|yaml}`, `required: true` (`.yaml` when the schema's file format is YAML, otherwise `.json`).

### Rules an author must satisfy (state these as prose; do not name error codes)

1. A schema may declare `instance file path` **or** instance slots, never both. Declaring both is a hard error, and neither wins.
2. Declaring an `instance file path` alone is exactly equivalent to one slot: `key: default`, `name: Default`, that file path, `required: true`.
3. If the slots key is present it must list at least one slot. To use the default slot, **omit the key** rather than writing an empty list.
4. Slot keys must be unique within a schema. Slot file paths must be unique within a schema and across every schema in a release.
5. A schema has exactly one instance format. If the schema declares `instance format`, that wins and the slot file paths' extensions are never compared. If it does not, the format is inferred from the **first** slot's file path, and every slot file path must then imply that same format — mixing `.json` and `.yaml` slot paths without an explicit `instance format` is an error. (`repos/backend` `origin/main` `internal/configs/services/config_schemas/create.go`: `resolveInstFormat` returns on the explicit format before it ever calls `verifyMatchingSlotFileTypes`, which is its only caller. `core` does not check this on any surface, so the error is raised server side when the release is created.) Never state the same-format rule unconditionally.
6. `required: true` is **fleet-wide** — it means every deployment of a release containing this schema, on any device. It is not a per-device rule. Per-device cardinality is deliberately not supported.
7. A deployment may contain **at most one** config instance per (schema, slot) pair, and **must** contain an instance for every slot marked `required: true`.
8. A schema's digest — which is how Miru deduplicates schemas — covers only each slot's `key`, `filepath`, and `required`. A push that changes only a slot's `name` or `description` deduplicates onto the existing schema and the change is silently ignored. Authors must know this.

### Annotation syntax, per language

JSON Schema — top-level `x-miru-instance-slots`, an array of objects:

    {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "x-miru-config-type": "motion-control",
        "x-miru-instance-slots": [
            {
                "key": "left-arm",
                "name": "Left Arm",
                "filepath": "/srv/miru/configs/v1/left-arm.json",
                "required": true,
                "description": "Config for the left arm controller"
            },
            {
                "key": "right-arm",
                "name": "Right Arm",
                "filepath": "/srv/miru/configs/v1/right-arm.json",
                "required": false
            }
        ],
        "type": "object"
    }

Opaque — top-level `instance_slots`, an array of objects:

    config_type: motion-control
    instance_slots:
      - key: left-arm
        name: Left Arm
        filepath: /srv/miru/configs/v1/left-arm.json
        required: true
        description: Config for the left arm controller
      - key: right-arm
        name: Right Arm
        filepath: /srv/miru/configs/v1/right-arm.json
        required: false

CUE — a repeated **declaration** attribute `@miru_slot(...)`, one per slot, applied in source order, sitting alongside `@miru(...)` rather than inside it:

    @miru(config_type="motion-control")
    @miru_slot(key="left-arm", name="Left Arm", filepath="/srv/miru/configs/v1/left-arm.json", required="true", description="Config for the left arm controller")
    @miru_slot(key="right-arm", name="Right Arm", filepath="/srv/miru/configs/v1/right-arm.json", required="false")
    {
        name: string
    }

CUE-specific facts worth documenting: a `@miru_slot` attached to a *field* rather than to the file's top-level declarations is silently ignored; both `required=true` and `required="true"` are accepted because CUE attribute values are always strings; two byte-identical `@miru_slot` lines are deduplicated by CUE itself before Miru sees them, so slots must differ; and an empty `description=""` is rejected in CUE although it is accepted in JSON Schema and Opaque. In a multi-file CUE package, `@miru_slot` attributes are collected from **every** file in the package — they do not have to sit on the file that carries `@miru`. (`CompilePackage` in `core` `pkg/schemas/cue/compile.go` unifies all documents, and `GetInstanceSlots` in `pkg/schemas/cue/attributes.go` reads declaration attributes off that unified value; verified empirically against `cuelang.org/go v0.17.1`.) Slot order is source order **within** a file; across the files of a package the order follows CUE's own file ordering, not the order the files are passed to the CLI. Because the instance format is inferred from the *first* slot's file path, a package that spreads slots across files should either keep every `@miru_slot` in one file or declare `instance format` explicitly. Do **not** write that all of a package's slots must live on the `@miru`-annotated file — no such rule exists.

### What does **not** change

- **The CLI.** `miru release create` gains no commands, no flags, and no output changes. Its output is byte-identical whether a schema declares one file path or five slots. Slots are authored purely in the schema annotation and observed in the dashboard.
- **The Miru Agent.** The agent does not model slots at all. It writes each config instance to its own absolute file path — which is the bound slot's file path — so multi-slot deployments land correctly with no agent change.
- **The published Platform API.** Version `2026-05-06` does **not** expose slots; the platform hard cut has not landed. Do not describe slots as a Platform API feature.

### Known drift, deliberately not fixed here

- The opaque `language:` → `schema_language:` rename (`core` `1756cc2` #134) is not reflected in the docs. Both spellings still parse (`DeprecatedSchemaLangKey`), and the docs show `language:` at `docs/snippets/references/cli/releases/create/schema-annotations.mdx:33` and `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx:21` and `:28`. Fixing it has its own release-timing question and would enlarge this PR. New opaque slot examples in this plan therefore do not restate the language key at all.
- `docs/snippets/agent/yaml-support.mdx` says Agent **v0.7.0** while `docs/changelog/agent.mdx` says v0.7.1. Pre-existing, carried over from `plans/completed/20260805-config-instance-file-formats.md`.

### Files that must not be touched

- `repos/docs/docs/references/platform-api/*.yaml` and everything under `repos/docs/docs/references/device-api/` — these are vendored from Stainless and hand edits are overwritten.
- `repos/docs/tools/lint/**` — editing it triggers two extra CI jobs plus a 90% Go coverage gate, and its diagnostic strings are asserted by `tests/test-lint.sh`. This plan requires no lint-rule change. Verify at the end with `git diff --stat main...HEAD -- tools/lint`, which must print nothing.
- `repos/docs/docs/docs.json` — no new page is added, so navigation does not change.

### Existing conventions you must match

**Headings** are strict sentence case and machine-checked (`heading-case`). The checker allowlists acronyms including `API CLI CUE JSON Miru GitHub Agent Git Schema`. `YAML` and `XML` are **not** allowlisted. `## Instance slots` is fine; `## Slots in JSON Schema` is fine; `## Slots in YAML` would fail. Version tags like `# v0.10.4` are exempt.

**Dashes**: em dash `—` only. Two consecutive hyphens in a prose span is a lint error (`no-double-dash`). Hyphens inside code spans, fenced blocks, and JSX attributes are fine.

**Voice**: second person for instructions, present tense, active. UI elements in **bold**, field and value names in `backticks`.

**Paths in examples** are realistic Miru paths: `/srv/miru/configs/mobility.json`, `/srv/miru/configs/v1/left-arm.json`. Any non-`/srv/miru` path must be paired with the existing `<Warning>` linking to `/developers/agent/filesys-access#configs`. The established config-type slug for the multi-controller example is `motion-control` — reuse it.

**The trailing-two-space `<CodeGroup>` separator.** In `docs/snippets/references/cli/releases/create/schema-annotations.mdx`, consecutive fenced blocks inside a `<CodeGroup>` are separated by a blank line followed by a line containing **exactly two spaces**. This is load-bearing for rendering and **no lint rule catches it**. The only guard is:

    cd /home/ben/miru/workbench2/repos/docs
    cat -A docs/snippets/references/cli/releases/create/schema-annotations.mdx | grep -c '^  \$'

It prints `6` today (2 separators × 3 `<CodeGroup>` blocks). After Milestone 1 it must print `8`. **Copy an existing `<CodeGroup>` block and edit the fence bodies** — do not type a new one, because editors that trim trailing whitespace on save will silently break it.

**`<ParamField>` anchors.** `<ParamField path="instance slots">` renders the anchor `#param-instance-slots`, which other pages can deep-link to.

## Plan of Work

### Milestone 1 — the shared annotations snippet

File: `repos/docs/docs/snippets/references/cli/releases/create/schema-annotations.mdx` (101 lines, four `<ParamField>` blocks, no headings, no imports). It is rendered by exactly three pages: `docs/cfg-mgmt/primitives/schemas/overview.mdx:9`, `docs/cfg-mgmt/create-a-release.mdx:9`, and `docs/getting-started/quick-start/create-release.mdx:12`. Anything added here appears on all three, including the beginner quick start — so keep it tight.

Add a fifth `<ParamField path="instance slots">` block **after** the `instance file path` block (which ends at line 73) and **before** the `instance format` block (line 75). Body indented two spaces, matching the neighbors. Content, in order:

1. One paragraph: a slot is a distinct file system destination for this schema; declare several to write several files from one schema; each slot has a `key`, `name`, `filepath`, `required`, and optional `description`.
2. A `<Warning>` stating that a schema may declare an instance file path **or** instance slots, never both, and that declaring both is an error. (Precedent: the existing `<Warning>` in the `instance file path` block.)
3. One sentence: this annotation is optional; omitting it is equivalent to a single slot at the instance file path, and the instance file path annotation is the shorthand for the single-slot case.
4. A `<CodeGroup>` with three tabs in the established order and tab titles — `` ```yaml JSON Schema ``, `` ```cue CUE ``, `` ```yaml Opaque `` — copied structurally from the `instance file path` block so the two-space separator lines survive. Use the `motion-control` two-slot example from Context and Orientation, trimmed to what fits a reference block.
5. A closing line noting that slot file paths must be unique across every schema in the release, and that unless the schema declares an `instance format`, every slot file path must imply the same format.

Then update the neighboring blocks so they stay true:

- `instance file path` (line 46): "This annotation is optional and defaults to `/srv/miru/configs/{config-type-slug}.json`." — add that it declares the schema's single slot, and cross-link to `#param-instance-slots` for the multi-destination case.
- `instance format` (line 80): "inferred from the instance file path's extension" — change to say that when this annotation is omitted the format is inferred from the schema's first slot file path, and that inference requires every slot file path to imply the same format. Do not say slot paths must always match: declaring `instance format` removes that requirement.

### Milestone 2 — primitives properties

File: `repos/docs/docs/cfg-mgmt/primitives/schemas/overview.mdx`.

The `## Properties` section (line 16) mirrors the annotation snippet; the prior plan established that when the snippet gains an annotation, this section must gain the matching property. Insert a new `<ParamField path="instance slots" type="InstanceSlot[]">` after the `instance file path` block (ends line 72) and before `instance format` (line 74). Open it with `<ImmutableBadge />` like every sibling. Document the five slot fields as a small Markdown table, using the `slug` block in `docs/cfg-mgmt/primitives/config-types.mdx:29-33` as the model **for formatting only** (a lead-in sentence plus a bulleted constraint list). Do not copy its framing: `slug` is validated client side, `key` is not — see "The `key` character rules are enforced by the platform" in Context and Orientation. State the uniqueness rules and that `required` is fleet-wide.

Also on this page:

- `instance file path` (lines 65-69): rewrite from "the absolute file system path where config instances are written" to make clear it is the file path of the schema's single slot, and point to the slots property when there is more than one.
- `instance format` (line 81): "inferred from the instance file path's extension" → when omitted, inferred from the first slot's file path; the same-format constraint across slot file paths applies only in that omitted case.
- `## Immutability` (line 116) needs no change, but add the digest caveat where it belongs — see Milestone 4, `create-a-release.mdx`.

File: `repos/docs/docs/cfg-mgmt/primitives/config-instances.mdx`.

Add a `<ParamField path="slot key" type="string">` with `<ImmutableBadge />` to `## Properties`. Place it after `file path` (lines 22-28) and before `config schema` (line 30). Say: every config instance is bound to exactly one instance slot on its schema, identified by the slot's key; the instance's file path is that slot's file path. Then amend the existing `file path` block to say the path comes from the bound slot rather than being set on the instance.

File: `repos/docs/docs/cfg-mgmt/primitives/config-types.mdx`. Review only. The imported definition `docs/snippets/definitions/config-type.mdx` says each config type has one versioned schema its instances adhere to — still true under slots. Change nothing unless a concrete falsehood is found; record the review in Progress.

File: `repos/docs/docs/cfg-mgmt/primitives/schemas/manage.mdx`. Under `## View a schema` (line 14), the `Metadata` tab lists what the dashboard shows for a schema. If the dashboard renders a slots list there (see "Where the frontend slots branch lives" in Scope — the branch is in `repos/backend` and shows only API models, so this needs a running dashboard), add one sentence naming it. **Verify against a running dashboard before writing** — do not describe UI you have not seen. If unverifiable, skip this file and note it in Surprises & Discoveries.

### Milestone 3 — schema-language pages

File: `repos/docs/docs/cfg-mgmt/primitives/schemas/languages/jsonschema.mdx`. The `## Example` block (line 16) carries no `x-miru-*` annotation today, so do not retrofit one. Instead add a short `## Instance slots` section after `## Example` (before `## Supported drafts` at line 150) with the two-slot `x-miru-instance-slots` example and one paragraph of explanation, cross-linking to `/cfg-mgmt/primitives/schemas/overview#param-instance-slots`.

File: `repos/docs/docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`. Same shape: an `## Instance slots` section after `## Example` (which ends around line 36, just before `## File formats` at line 38) showing the `instance_slots` YAML block. Do **not** touch the `language: opaque` lines (Decision Log).

File: `repos/docs/docs/cfg-mgmt/primitives/schemas/languages/cue.mdx`. This is the language with the most surface area, because `@miru_slot` is a repeated declaration attribute rather than an argument to `@miru`. Add an `## Instance slots` section at the end of `## Example` — after the paragraph beginning "As you can see, CUE inlines constraints directly with type definitions" — and before `## CUE version`. Do not insert it at the "Although many of CUE's capabilities are omitted" lead-in; that would split `## Example` in half. It must cover: the attribute is separate from `@miru`, not a parameter of it; one attribute per slot; source order is preserved; an attribute placed on a field is ignored; `required` may be written bare or quoted; two identical `@miru_slot` lines collapse into one, so slots must differ; and an empty `description=""` is **rejected** in CUE although it is accepted in JSON Schema and Opaque (CUE rejects any empty attribute value generically). That last item is the only cross-language asymmetry in this feature and must not be dropped.

File: `repos/docs/docs/snippets/references/cli/releases/create/cue-packages.mdx`. **No required change.** Its closing sentence — "To annotate a CUE package, annotate **exactly one file** in the package. Annotating multiple files or no files will result in an error." — is about the `@miru` attribute and stays true. It must **not** gain a clause requiring `@miru_slot` attributes to sit on that same file: Miru collects `@miru_slot` from every file in the package (see "Annotation syntax, per language"). If anything is added here, it may only be a recommendation, never a rule — for example: "`@miru_slot` attributes may be declared on any file in the package. Keep them in one file so their order is predictable, since the instance format is inferred from the first slot's file path." Adding that sentence is optional; if it is skipped, record the decision in Progress. This snippet is imported only by `cue.mdx`.

File: `repos/docs/docs/cfg-mgmt/primitives/schemas/languages/overview.mdx` (16 lines, no headings). No change. Adding a heading here to house a slots comparison would be the only heading on the page and adds nothing the per-language sections do not. Record the decision not to touch it.

### Milestone 4 — downstream drift

File: `repos/docs/docs/cfg-mgmt/create-a-release.mdx`.

- Line 41: "Schemas may belong to multiple releases and are identified by hashing their content. Schemas with equivalent content are considered identical (even if comments, spacing, or other formatting is different)." Append the slot-digest caveat: the hash covers each slot's key, file path, and required flag, so changing only a slot's name or description produces the same schema and the change is not applied.
- Line 72: the lead-in to `### Schema annotations` names "the `instance file path` annotation, which defines where config instances for this schema are written to disk". Extend to mention instance slots as the way to write more than one file from one schema.
- Consider adding, under `## Config schemas` (line 37), the release-wide file path uniqueness rule: no two slots anywhere in a release may share a file path.

File: `repos/docs/docs/getting-started/quick-start/create-release.mdx`, line 53: "Each of the above schema examples are annotated with the `x-miru-config-type` and `x-miru-instance-filepath` fields." Already stale (it omits `instance format`); with a fifth annotation block rendering below it, it drifts further. Reword to something version-proof such as naming the config type annotation as the required one and referring the reader to the blocks below for the rest.

File: `repos/docs/docs/snippets/definitions/deployment-constraints.mdx`. Imported only by `docs/primitives/deployments.mdx:8`. Add a fifth numbered constraint for slot cardinality: a deployment contains at most one config instance per schema slot, and must contain an instance for every slot marked required. Match the existing format exactly — each item is written as `1.` (the list auto-numbers), a bolded one-line title, blank line, then an indented explanatory paragraph.

File: `repos/docs/docs/cfg-mgmt/deploy/initial-deployment.mdx`, lines 46-48 and 67. "it automatically fills with default values from your release's config schemas" and "For a first deployment, every file is listed as added." Add one clause noting that a schema with several slots contributes one file per slot.

File: `repos/docs/docs/cfg-mgmt/deploy/staging-area.mdx`, around line 110 ("To view the details of any config instances in the deployment, click into the config instance in the **Configurations** section"). Add that the list has one entry per schema slot. Verify the wording against a running dashboard before writing (not against the backend branch — see Scope); if unverifiable, keep the edit to the schema-agnostic statement and note it.

File: `repos/docs/docs/cfg-mgmt/deploy/config-editor.mdx`, lines 50-54 and 157-164. This page already speaks in terms of *files*, not schemas, so it mostly fits without contradiction. **Do not guess what label the editor tab shows for a slot** (slot `name` versus file name). Either verify against a running dashboard and write the verified label — `repos/backend` `origin/refactor/config-schema-instance-slots-frontend` carries API models only and cannot answer this — or make no change and record the open question in Surprises & Discoveries. Guessing here is worse than silence.

### Milestone 5 — CLI reference and changelog

File: `repos/docs/docs/references/cli/release-create.mdx`. The CLI gains no flags and no output change, so `<Flags />`, `<Usage />`, and the example snippets need no edit. Under `### Requirements` (line 15), the existing bullet "Schemas must be annotated with their config types (see [annotations](...))" is the right place for one added sentence: schemas that write more than one file must declare instance slots, linking to `/cfg-mgmt/primitives/schemas/overview#param-instance-slots`. That is the whole CLI reference change.

File: `repos/docs/docs/changelog/cli.mdx`. This file is hand-written and is **not** derived from `repos/cli-private/CHANGELOG.md` (which has no slots entry at all). Changelogs are historical records — never rewrite an existing entry; add a new one at the top.

Format, copied from the existing top entry:

    # v0.10.4

    *Month D, YYYY*

    ## Features

    - Add support for the instance slots schema annotation, letting one config schema write config instances to several file paths

    ---

The version must be the CLI version that actually ships the wiring. `repos/cli-private` currently has tags up to `v0.10.4-beta.2`, so `v0.10.4` is the expected number, but **confirm the tag exists before the PR leaves draft** and correct both the heading and the date if it differs. The date is the release date, not the day you write it.

If the release turns out to include a breaking change, the `v0.10.0` entry in this same file is the precedent for wrapping migration steps in `<Dropdown title="...">` around a ` ```diff ` block with `# JSON Schema` / `# CUE` comment separators. Slots are additive, so this is not expected to be needed.

`repos/docs/docs/changelog/product.mdx` is a separate, richer changelog using `<Update label="...">` blocks. Adding a product entry there is **optional** and out of this plan's required scope; if added, note that headings inside `<Update>` blocks use Title Case, which is only legal inside JSX blocks.

## Concrete Steps

All commands run from `/home/ben/miru/workbench2/repos/docs` unless stated otherwise.

Every `./scripts/lint.sh` invocation below must print the four banners `== MDX Prose ==`, `== ESLint (MDX) ==`, `== CSpell ==`, `== OpenAPI ==` and end with the exact line `All documentation lint checks passed.` Anything else is a failure; do not commit.

### Milestone 0 — baseline

    cd /home/ben/miru/workbench2/repos/docs
    git status --porcelain          # expect empty
    git branch --show-current       # expect docs/instance-slots
    git fetch origin && git log --oneline -1 origin/main

If the branch does not exist yet:

    git checkout main && git pull && git checkout -b docs/instance-slots

Install dependencies. Skipping this produces a confusing `Command "eslint" not found` failure later:

    pnpm install --frozen-lockfile

Establish the audit baseline **before making any edit**, because `scripts/preflight.sh` uses `set -euo pipefail` and a failing audit aborts the run before the later stages. There is no gitignored scratch directory in this repo, so use a temporary one and keep both outputs on disk. These two snapshots are a starting reference, not the draft-exit gate itself: the gate in Validation and Acceptance re-captures both sides at the head you are shipping and diffs those.

    AUDIT_DIR="$(mktemp -d)"; echo "AUDIT_DIR=$AUDIT_DIR"   # record this path in Surprises & Discoveries
    ./scripts/audit.sh > "$AUDIT_DIR/audit-branch.txt" 2>&1; echo "exit=$?"

As of 2026-08-10 this exits **1** with `3 vulnerabilities found / Severity: 3 high (1 ignored)` — transitive `js-yaml` advisories (GHSA-5p4m-2wfm-xmqj) reached only through the `mint` dev dependency. This is **pre-existing and unrelated to documentation content**. Confirm pre-existence rather than assuming it:

    git worktree add "$AUDIT_DIR/docs-main" main
    (cd "$AUDIT_DIR/docs-main" && pnpm install --frozen-lockfile && ./scripts/audit.sh > "$AUDIT_DIR/audit-main.txt" 2>&1; echo "exit=$?")
    git worktree remove "$AUDIT_DIR/docs-main"
    diff "$AUDIT_DIR/audit-main.txt" "$AUDIT_DIR/audit-branch.txt" && echo "IDENTICAL"

`$AUDIT_DIR` does not survive a new shell — write the literal path into Surprises & Discoveries along with both exit codes and the advisory list. If `diff` prints anything, a *new* advisory has appeared; that is a regression and must be resolved, not documented.

Confirm the release dependencies:

    cd /home/ben/miru/workbench2/repos/core && git fetch origin && git log --oneline -1 origin/main
    git branch -r | grep instance-slots       # core: expect feat/instance-slots-cue and feat/instance-slots-annotation
    cd /home/ben/miru/workbench2/repos/cli-private && git fetch --tags origin
    git branch -r | grep instance-slots       # the CLI wiring branch lives HERE, not in core
    git log --oneline -1 origin/feat/instance-slots-annotation
    git tag --sort=-v:refname | head -5

As of 2026-08-10: `core` `origin/main` head is `f1be215` (#151, opaque slots) with #150 (JSON Schema slots) beneath it; `origin/feat/instance-slots-cue` (#152) is **unmerged**; `cli-private` pins `github.com/mirurobotics/core v0.9.1` in `go.mod` and its newest tag is `v0.10.4-beta.2`. Re-check these; if #152 and the CLI wiring have landed since, note it and the draft gate in Milestone 6 becomes trivially satisfiable.

Commit nothing in this milestone.

### Milestone 1 — annotations snippet

Edit `docs/snippets/references/cli/releases/create/schema-annotations.mdx` per Plan of Work. Then, before anything else:

    cat -A docs/snippets/references/cli/releases/create/schema-annotations.mdx | grep -c '^  \$'

Expected output:

    8

If it prints `6` or `7`, the trailing-space separator lines in the new `<CodeGroup>` were lost. Fix by copying an existing block again rather than retyping.

    ./scripts/lint.sh

Commit:

    git add docs/snippets/references/cli/releases/create/schema-annotations.mdx
    git commit -m "docs(snippets): document the instance slots schema annotation"

### Milestone 2 — primitives properties

Edit `docs/cfg-mgmt/primitives/schemas/overview.mdx` and `docs/cfg-mgmt/primitives/config-instances.mdx`; review `docs/cfg-mgmt/primitives/config-types.mdx` and `docs/cfg-mgmt/primitives/schemas/manage.mdx`.

    ./scripts/lint.sh
    git add docs/cfg-mgmt/primitives
    git commit -m "docs(schemas): add instance slots and slot key properties"

### Milestone 3 — schema-language pages

Edit the three language pages, and optionally `docs/snippets/references/cli/releases/create/cue-packages.mdx` (see Plan of Work — that snippet may legitimately end up unchanged; drop it from the `git add` below if untouched).

New headings introduced here are `## Instance slots` on three pages — all sentence case, all pass `heading-case`. Confirm no heading uses `YAML` or `XML`, which are not allowlisted.

    ./scripts/lint.sh
    git add docs/cfg-mgmt/primitives/schemas/languages docs/snippets/references/cli/releases/create/cue-packages.mdx
    git commit -m "docs(schemas): show instance slot annotations in each schema language"

### Milestone 4 — downstream drift

Edit `docs/cfg-mgmt/create-a-release.mdx`, `docs/getting-started/quick-start/create-release.mdx`, `docs/snippets/definitions/deployment-constraints.mdx`, and the three `docs/cfg-mgmt/deploy/*.mdx` pages.

    ./scripts/lint.sh
    git add docs/cfg-mgmt docs/getting-started docs/snippets/definitions/deployment-constraints.mdx
    git commit -m "docs(cfg-mgmt): reflect instance slots in release and deployment guides"

### Milestone 5 — CLI reference and changelog

Edit `docs/references/cli/release-create.mdx` and `docs/changelog/cli.mdx`.

    ./scripts/lint.sh
    git add docs/references/cli/release-create.mdx docs/changelog/cli.mdx
    git commit -m "docs(cli): document instance slots in the CLI reference and changelog"

### Milestone 6 — validation and ship

Confirm the linter itself was never modified:

    git diff --stat main...HEAD -- tools/lint

Expected output: nothing at all.

Confirm the navigation was not modified:

    git diff --stat main...HEAD -- docs/docs.json

Expected output: nothing at all.

Confirm no vendored spec was touched:

    git diff --name-only main...HEAD | grep -E 'references/(platform|device)-api' || echo "clean"

Expected output: `clean`

Full preflight:

    ./scripts/preflight.sh

Render check — start the local site and visually confirm the new content:

    pnpm dev

Then open, in a browser, and confirm each renders with the new content and working `<CodeGroup>` tabs:

- `http://localhost:3000/cfg-mgmt/primitives/schemas/overview` — the `instance slots` property in **Properties**, and the fifth annotation block under **Schema annotations** with three working tabs.
- `http://localhost:3000/cfg-mgmt/create-a-release` — same annotation block under **Schema annotations**.
- `http://localhost:3000/getting-started/quick-start/create-release` — same block; confirm the reworded lead-in reads sensibly for a beginner.
- `http://localhost:3000/cfg-mgmt/primitives/config-instances` — the `slot key` property.
- `http://localhost:3000/cfg-mgmt/primitives/schemas/languages/cue` — the `@miru_slot` section.
- `http://localhost:3000/changelog/cli` — the new version entry at the top, above `v0.10.3`.

Commit the plan file itself with its Progress and Decision Log filled in, then ship:

    git add plans/
    git commit -m "docs(plans): record the instance slots documentation plan"
    git push -u origin docs/instance-slots
    gh pr create --draft --base main --title "docs: document config schema instance slots"

The PR body must state the release dependency explicitly, in words like: *this documents the `instance slots` schema annotation, which requires `core` #152 (CUE) and the `cli-private` wiring to be merged and released as CLI `vX.Y.Z`. Merging this PR does not publish — publishing is the manual `promote.yml` workflow — but do not promote until that CLI version is released.*

Then:

    gh pr checks --watch

When all checks pass and the CLI dependency has shipped:

    gh pr ready

Note: `gh pr edit` fails in this org because of the Projects-classic deprecation. To edit the PR body, use `gh api --method PATCH repos/mirurobotics/docs/pulls/<N> -f body=...` instead.

## Validation and Acceptance

**Preflight must report CLEAN before this task is reported complete and before the PR leaves draft.** CLEAN means, on the head SHA you actually pushed: the `changes`, `lint`, and `shell-tests` CI jobs are green, **and** the `audit` job is either green or red with an advisory set identical to pristine `main` **measured at the same moment as the head you are shipping**. The Milestone 0 pair is a starting reference only: both snapshots are taken before any edit and before the push, so diffing them proves nothing about the pushed head, and the advisory set drifts while the PR waits in draft for the CLI release. Immediately before `gh pr ready`, re-capture both sides and compare:

    cd /home/ben/miru/workbench2/repos/docs
    git checkout docs/instance-slots && git pull --ff-only   # be on the pushed head
    ./scripts/audit.sh > "$AUDIT_DIR/audit-head.txt" 2>&1; echo "exit=$?"
    git worktree add "$AUDIT_DIR/docs-main-recheck" main
    (cd "$AUDIT_DIR/docs-main-recheck" && pnpm install --frozen-lockfile && ./scripts/audit.sh > "$AUDIT_DIR/audit-main-recheck.txt" 2>&1; echo "exit=$?")
    git worktree remove "$AUDIT_DIR/docs-main-recheck"
    diff "$AUDIT_DIR/audit-main-recheck.txt" "$AUDIT_DIR/audit-head.txt" && echo "IDENTICAL"

`diff` must print nothing. If it does, an advisory exists on this branch that pristine `main` does not have; resolve it rather than documenting it. A local-only pass is not sufficient, and a green run on an older commit is not sufficient. Confirm the job states with `gh pr checks --watch` against the head SHA you actually pushed. This is the only definition of CLEAN in this plan.

A docs-only PR that does not touch `tools/lint/**` runs exactly four CI jobs: `changes`, `lint`, `audit`, and `shell-tests`. The two custom-linter jobs are gated off.

### Known pre-existing failure — do not mistake this for a regression

`scripts/audit.sh` runs `pnpm audit --ignore-registry-errors`. As of 2026-08-10 it exits **1** on transitive `js-yaml` advisories (`GHSA-5p4m-2wfm-xmqj`) reachable only through the `mint` dev dependency, reporting `3 vulnerabilities found / Severity: 3 high (1 ignored)`. These advisories were published after the last green CI run on `main` (`6d99108`, green 2026-08-06), so a red `audit` job on this branch may have nothing to do with documentation content — whether that is acceptable is decided solely by the CLEAN definition above (an advisory set identical to pristine `main` re-measured at the head you are shipping, not the Milestone 0 snapshot). Because `preflight.sh` uses `set -euo pipefail`, this also means `./scripts/preflight.sh` aborts at `=== Audit ===` and never reaches `=== Shell Script Tests ===`.

Handling:

1. Prove pre-existence with the pristine-`main` worktree recipe in Milestone 0 and paste both outputs into Surprises & Discoveries.
2. Run the stages preflight skipped, by hand, so nothing goes unverified:

        ./scripts/lint.sh
        pnpm run test:lint
        bats pub/scripts/agent/check-miru-access_test.bats

3. Note the pre-existing red `audit` job in the PR body, with the advisory ID and the proof, so a reviewer does not read it as caused by this change. Fixing the advisory (a dependency bump, likely already open as a Dependabot PR) is a **separate** PR — do not bundle it here.

If the advisory set has changed by the time you run this, re-derive the baseline; do not copy the list above verbatim into the PR.

### Acceptance criteria

Each of these is checkable by a human:

1. `cat -A docs/snippets/references/cli/releases/create/schema-annotations.mdx | grep -c '^  \$'` prints `8` (it printed `6` before the change).
2. `./scripts/lint.sh` ends with `All documentation lint checks passed.`
3. `pnpm run test:lint` passes.
4. `bats pub/scripts/agent/check-miru-access_test.bats` passes.
5. `git diff --stat main...HEAD -- tools/lint docs/docs.json` prints nothing.
6. `git diff --name-only main...HEAD | grep -E 'references/(platform|device)-api'` finds nothing.
7. Under `pnpm dev`, the six URLs listed in Milestone 6 render, and the new `<CodeGroup>` on the schemas overview page shows three selectable tabs labelled **JSON Schema**, **CUE**, and **Opaque**, each containing a syntactically valid annotation.
8. `grep -rn "coming soon" docs --include='*.mdx'` finds no new occurrence introduced by this change.
9. `grep -rni "multi-instance mode\|is_dynamic\|deprecated" docs/cfg-mgmt --include='*.mdx'` finds no new occurrence — slots are not a mode, and `instance file path` is not deprecated.
10. `grep -rn "invalid_instance_slots\|conflicting_instance_target\|instance_slot_" docs --include='*.mdx'` finds nothing — error codes are deliberately not documented.
11. CLEAN as defined at the top of this section holds on the pushed head SHA, including the draft-exit re-capture (`diff "$AUDIT_DIR/audit-main-recheck.txt" "$AUDIT_DIR/audit-head.txt"` prints nothing).
12. The CLI changelog entry's version heading matches a real tag in `repos/cli-private`.

Items 13-19 are confirmed by reading the diff, without opening any other repository:

13. `required` is documented as a mandatory slot field with no default.
14. `required: true` is described as fleet-wide, not per-device.
15. The instance-file-path / instance-slots mutual exclusion is stated, and neither is described as winning.
16. Omitting the slots key is described as the way to get the default slot; an empty list is described as invalid.
17. The digest caveat (name/description-only changes are silently ignored) appears on `create-a-release.mdx`.
18. Slots appear nowhere as a Platform API feature.
19. The CUE section says `@miru_slot` sits alongside `@miru` and is not one of its arguments, and states that `description=""` is rejected in CUE but accepted in JSON Schema and Opaque.

## Idempotence and Recovery

Every step is a file edit or a read-only command, so the whole plan is safe to re-run. Specifically:

- `pnpm install --frozen-lockfile`, `./scripts/lint.sh`, `./scripts/audit.sh`, `pnpm run test:lint`, and `bats ...` are read-only and repeatable.
- `pnpm dev` starts a local server on port 3000; stop it with Ctrl-C. It writes nothing to the repository.
- The Milestone 0 baseline check creates a `mktemp -d` directory and a worktree inside it; undo with `git worktree remove "$AUDIT_DIR/docs-main"` (and `git worktree remove "$AUDIT_DIR/docs-main-recheck"` for the draft-exit re-capture). If that fails because the directory is dirty, use `git worktree remove --force "$AUDIT_DIR/docs-main"`, then `git worktree prune`. Delete `$AUDIT_DIR` only after the PR is out of draft — the gate needs the saved outputs.
- To undo a single milestone: `git revert <sha>` for that milestone's commit, or `git reset --hard HEAD~1` if it has not been pushed.
- To restart from scratch: `git checkout main && git branch -D docs/instance-slots` and begin at Milestone 0. Nothing outside this repository is modified, so there is no external state to clean up.

The one step with a non-obvious failure mode is the `<CodeGroup>` trailing-space separator in Milestone 1. It fails silently — lint passes, rendering breaks. If `cat -A ... | grep -c '^  \$'` does not print `8`, revert the file with `git checkout -- docs/snippets/references/cli/releases/create/schema-annotations.mdx` and redo the edit by duplicating an existing `<CodeGroup>` block.

Nothing in this plan is destructive, and no migration or data change is involved.

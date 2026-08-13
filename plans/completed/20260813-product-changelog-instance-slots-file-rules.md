# Add the August 14, 2026 product changelog entry (instance slots and file rules)

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `/home/ben/miru/workbench5/repos/docs` | read-write | One new `<Update>` block in `docs/changelog/product.mdx`, plus this plan file. `cspell.json` may gain words only if the spell checker flags the new text (not expected). No other file changes. |

`repos/docs` is the **only** checkout in this workbench (`ls /home/ben/miru/workbench5/repos` → `docs`). `core`, `cli-private`, `backend`, `frontend`, and `agent` are **not** available. Every factual claim in the new entry must therefore be grounded in a file inside this repo; anything that would require a sibling repo to verify is out of the changelog and goes in the PR body (see D2 and D3).

Work happens on the already-created branch `docs/product-changelog-26-08-14` (base `main`, HEAD `e001038`). Do not create new branches. Commits must be signed.

## Purpose / Big Picture

Two features shipped and are documented in the reference docs, but the customer-facing product changelog (`docs/changelog/product.mdx`, served at `/changelog/product`) still opens with "August 6, 2026" and mentions neither:

1. **Instance slots** (PR [#150](https://github.com/mirurobotics/docs/pull/150), commit `73d8e05`) — a config schema may declare several instance slots, each a distinct file system destination for the same schema, and each config instance binds to exactly one slot via its `slot key`. This is the relaxation that lets one schema write more than one file.
2. **File rules** (PR [#154](https://github.com/mirurobotics/docs/pull/154), commit `f5ffbab`) — the upload rule became the file rule, and the object was restructured: `upload` (where matches go) and `retention` (how long the device keeps them) are now two independent optional halves of one rule.

After this change, `/changelog/product` opens with an `<Update label="August 14, 2026">` block containing exactly two `##` sections, each with a short bold-lead paragraph, bulleted `**Term:**` detail lines, links into the reference docs, and one media placeholder. A reader can verify by rendering the site (or reading the file) and seeing the new entry above the August 6 one.

The pictures do not exist yet — the user adds them after this lands. The entry therefore carries two greppable `TODO(assets)` placeholders (D4).

## Progress

- [x] M0 — Re-verify the Context facts against the working tree (they are the entry's only evidence base). *(2026-08-13 — all greps matched the plan's predictions; two prose corrections required, see Surprises.)*
- [x] M1 — Determine what version requirements are substantiable from this repo (`docs/changelog/cli.mdx`, `docs/changelog/agent.mdx`) and record the answer in Surprises & Discoveries before writing any prose. *(2026-08-13 — CLI `v0.11.0` and Agent `v0.10.1` both stand as D2 predicted.)*
- [x] M2 — Insert the `<Update label="August 14, 2026">` block into `docs/changelog/product.mdx`. *(2026-08-13 — 72 insertions, 0 deletions.)*
- [x] M3 — `./scripts/lint.sh` and `pnpm run validate` both pass; confirm the `TODO(assets)` placeholders do not break either (fallback per D5 if they do). *(2026-08-13 — all green first try; D5 fallback not needed.)*
- [x] M4 — Commit (signed), push, open a **draft** PR whose body carries the unsubstantiable claims. *(2026-08-13)*
- [x] M5 — `$preflight` reports `CLEAN` on the pushed head; only then does the PR leave draft. *(2026-08-13 — CI green; PR deliberately left in draft for the orchestrator.)*

## Surprises & Discoveries

- **2026-08-13 (planning) — The task brief's grounding list overstates what the merged docs contain.** `grep -rn "x-miru-instance-slots\|@miru_instance_slot\|instance_slots" docs --include='*.mdx'` returns hits in exactly two files: `docs/cfg-mgmt/primitives/schemas/overview.mdx` (5) and `docs/snippets/references/cli/releases/create/schema-annotations.mdx` (4). The three schema-language pages under `docs/cfg-mgmt/primitives/schemas/languages/` contain **no** slot content, and `docs/cfg-mgmt/deploy/deployment-constraints.mdx` **does not exist** — the only deployment-constraints file is `docs/snippets/definitions/deployment-constraints.mdx`, which has zero slot mentions. `git show --stat 73d8e05` confirms PR #150 touched only `create-a-release.mdx`, `config-instances.mdx`, `schemas/overview.mdx`, `changelog/cli.mdx`, `quick-start/create-release.mdx`, and the CLI `schema-annotations.mdx` snippet. Consequence: the entry links only to pages that actually carry the content (D1).
- **2026-08-13 (planning) — There is no "synthesized `default` slot" prose and no slot cardinality rule in the merged docs.** The completed plan `plans/completed/20260810-document-instance-slots.md` (Progress, Milestone 2 and 4) claims both landed. They did not. The only default-behavior statement anywhere is `docs/snippets/references/cli/releases/create/schema-annotations.mdx:84` — "This annotation is optional. Omitting it will either use the instance file path annotation, or create a single slot with default values." — alongside the XOR warning at :81 ("A schema may declare an instance file path OR instance slots. Declaring both will error."). `default` appears as an example slot key at `docs/cfg-mgmt/primitives/config-instances.mdx:35` and nowhere else. Consequence: the entry may state the XOR rule and the "one slot by default" behavior, but must **not** state deployment cardinality ("a deployment holds at most one config instance per slot") — nothing in this repo says it.
- **2026-08-13 (planning) — PR #153 did not do what its body says.** Its body ("docs(changelog): add date to v0.11.0 CLI entry (#153)", commit `1a386cd`) states "This removes the entry", referring to the `# v0.11.0` CLI changelog entry added by #150 for a version that did not exist. The actual diff — `git show --format= -p 1a386cd -- docs/changelog/cli.mdx` — is a single line: `-*Unreleased*` / `+*August 12, 2026*`. So `main` carries a **dated** `# v0.11.0` entry today, and the three-part release gate recorded in the instance-slots plan was closed by assertion, not by verification.
- **2026-08-13 (planning) — `plans/active/20260812-agent-v0.10.1-changelog.md` is still open at M3, but PR #155 merged (`e001038`).** Its unchecked boxes are exactly the tag-verification steps ("verify the `v0.10.1` tag exists in the agent repo"). The merged `docs/changelog/agent.mdx` entry also has a `## Features` section, which that plan's Decision Log says it deliberately would not have. Same pattern as #150/#153: changelog entries on `main` whose release gates were never closed. One line in the PR body; not the changelog's problem.
- **2026-08-13 (M0) — Every Context grep matched the plan's prediction; no Context revision was needed.** The annotation grep returned exactly the two predicted files (`docs/cfg-mgmt/primitives/schemas/overview.mdx` ×5, `docs/snippets/references/cli/releases/create/schema-annotations.mdx` ×4); `docs/snippets/file-rules/retention.mdx` carries both `require_upload` (:6) and `ttl_secs` (:13, `required`); `docs/snippets/references/cli/releases/create/scopes.mdx:11` has one `file_rules:manage` hit; the idempotence guard printed `0`. The XOR `<Warning>` and the "create a single slot with default values" sentence are both present in `schema-annotations.mdx` as quoted. Slot fields at `overview.mdx:187-231` are as recorded, with `<ImmutableBadge />` on all five (the plan's Context only noted immutability for `key`).
- **2026-08-13 (M0) — Two corrections to the Plan of Work's prose, both forced by the working tree.** (a) The plan's lead paragraph said "a single `motion-control` schema with four slots". The docs' worked example is a **`motor-controller`** schema, and `motion-control` appears only as a config-instance file path example (`config-instances.mdx:27`). Also, `overview.mdx` is internally inconsistent about the count — the `<CodeGroup>` at :115-175 declares **three** slots while the sentence at :179 says "serve four identical motor controllers" — so naming any number would cite a contradiction. The entry now reads "several identical motor controllers … a single `motor-controller` schema with one slot per controller", which the file supports without picking a side. (b) The plan's bullet said "A config instance's file path comes from the slot it's bound to". Nothing states that derivation: `config-instances.mdx:22` documents `file path` as its own property, and `overview.mdx:205-213` says only that a slot's `filepath` is "where this slot's config instance is written on the device". Reworded to "The slot's `filepath` is where that instance is written on the device", which is the quoted claim rather than an inference about where the value originates.
- **2026-08-13 (M1) — Version requirements: both D2 versions stand unchanged; the entry cites Miru CLI `v0.11.0` and Miru Agent `v0.10.1`, and nothing else.** `docs/changelog/cli.mdx:10` still carries `# v0.11.0` dated *August 12, 2026* (not reverted to `*Unreleased*`), whose `## Breaking changes` covers the file-rule flags, rule file format, and scopes, and whose `## Features` bullet at :66 reads "Add support for the [instance slots](…) schema annotation, letting one config schema write config instances to several file paths". `docs/changelog/agent.mdx:6` still carries `# v0.10.1` dated *August 12, 2026*, Features bullet "Added support for the latest [file rule](…) capabilities, including retention-only rules and a configurable delay period before deletion". So the M1 step-5 escape hatch (drop the version line) was **not** triggered. The evidence base is this repo alone: `ls /home/ben/miru/workbench5/repos` → `docs`, so `core`, `cli-private`, `backend`, `frontend`, and `agent` are absent and the four D3 gate questions (sibling-repo unavailability; the three-part instance-slots release gate; PR #154's stable-`v0.11.0` gate versus `f5ffbab` already on `main`; PR #153's body-versus-diff mismatch) cannot be closed here. All four are stated in the PR body.
- **2026-08-13 (M3) — The `TODO(assets)` placeholders survived both checks; the D5 fallback was not needed.** `pnpm run test:lint`, `./scripts/lint.sh` ("All documentation lint checks passed.", CSpell 0 issues across 152 files), and `pnpm run validate` ("success build validation passed") all passed on the first run with the `<Framed>` components outside their comments. As D5 predicted, `image-domain` only checks the `https://assets.mirurobotics.com/` prefix and `mint validate` never resolves remote assets, so the not-yet-uploaded PNGs are invisible to both. `cspell.json` was **not** modified — no word in the new entry was flagged.
- **2026-08-13 (M3) — `audit` measured, and it cannot regress from this change.** `./scripts/audit.sh` exits **0** locally on the branch head ("1 vulnerabilities found / Severity: 1 high (1 ignored)") — the known `GHSA-5p4m-2wfm-xmqj` is absorbed by `ignoreCves`, matching the rename plan's later observation that `audit.sh` had gone green. Stronger than a same-moment comparison: `git diff main -- package.json pnpm-lock.yaml` is **empty**, so the audit job's entire input is byte-identical to pristine `main` and any red would be pre-existing by construction.

## Decision Log

- **D1: Item 1 links only to `/cfg-mgmt/primitives/schemas/overview#instance-slots`, `/cfg-mgmt/primitives/config-instances`, and `/references/cli/release-create`.**
  Rationale: those are the only pages in this repo that document slots (see Surprises & Discoveries). Linking the schema-language pages or a `deploy/deployment-constraints` page would send readers to pages with no slot content, or to a 404 that `pnpm run validate` would catch anyway.
  Date/Author: 2026-08-13 / plan author.
- **D2: Version requirements name only what this repo substantiates — Miru CLI `v0.11.0` and Miru Agent `v0.10.1`.**
  Rationale: `docs/changelog/cli.mdx:10` carries `# v0.11.0` dated *August 12, 2026*, whose Features list includes the instance-slots bullet (`:66`) and whose `## Breaking changes` section (`:16-62`) covers the file-rule flags, YAML shape, and scopes. `docs/changelog/agent.mdx:6-12` carries `# v0.10.1` dated *August 12, 2026*, whose Features bullet reads "Added support for the latest file rule capabilities, including retention-only rules and a configurable delay period before deletion". Both are in-repo evidence. Nothing else is: sibling repos are not cloned, so the three-part instance-slots gate (core annotation support / a `cli-private` tag / a backend CLI API accepting `instance_slots`) and the stable-`v0.11.0` question from `plans/completed/20260811-rename-upload-rules-to-file-rules.md` (D1, Open Questions) cannot be checked here.
  Date/Author: 2026-08-13 / plan author.
- **D3: Unsubstantiable claims go in the PR body, never in the changelog.**
  Rationale: a changelog entry is a claim that something shipped; a PR body is a claim about what the author could verify. The PR body must state, explicitly: (a) sibling repos are unavailable in this workbench, so the version numbers rest on this repo's own changelog entries; (b) `plans/completed/20260810-document-instance-slots.md` records a three-part release gate for instance slots whose parts (b) and (c) were unmet as of 2026-08-10 and are unverifiable here; (c) `plans/completed/20260811-rename-upload-rules-to-file-rules.md` left PR #154 in draft pending a *stable* `v0.11.0` tag, yet `f5ffbab` is on `main` and `docs/changelog/cli.mdx` carries a dated `v0.11.0`; (d) PR #153's body says it removed the v0.11.0 entry but its diff only dated it. A reviewer with access to `cli-private` can close all four in one look.
  Date/Author: 2026-08-13 / plan author.
- **D4: Media placeholders are real components preceded by an MDX comment whose exact marker is `TODO(assets)`.**
  Rationale: MDX has no HTML comments — `<!-- … -->` is a parse error. The repo's precedent is `{/* TODO (screenshot): … */}` (`docs/snippets/data-uploads/buckets/verify-gcs.mdx:10`, `docs/snippets/data-uploads/buckets/register-gcs.mdx:17`). The marker is written without a space (`TODO(assets)`) so `grep -rn 'TODO(assets)' docs/` finds exactly the two placeholders and nothing else. URLs follow the existing convention `https://assets.mirurobotics.com/docs/changelog/<yy-mm-dd>/<slug>.<ext>` (see `docs/changelog/product.mdx:31`, `:44`, `:65`), giving `…/26-08-14/instance-slots.png` and `…/26-08-14/file-rules.png`.
  Date/Author: 2026-08-13 / plan author.
- **D5: The placeholders use `<Framed image="…" />` (not `<LazyVideo>`), and the fallback if either check objects is to wrap the component inside the same `{/* TODO(assets) … */}` comment.**
  Rationale: `image-domain` (`tools/lint/linter/imagedomain/imagedomain.go`) requires the `https://assets.mirurobotics.com/` prefix and checks nothing else — it never fetches the URL, so a not-yet-uploaded asset cannot fail it. `pnpm run validate` (`mint validate`) checks the build and internal links, not remote asset availability. A still image also fails more gracefully than a video if the asset is late. If, contrary to expectation, M3 shows either check resolving remote assets, take the fallback: move the whole component **inside** the `{/* TODO(assets): … */}` comment so the page renders without a broken image while the marker still greps. Do not switch to a non-`assets.mirurobotics.com` URL and do not add a `lint-ignore` directive.
  Date/Author: 2026-08-13 / plan author.
- **D6: Headings are Title Case (`## Multiple Configs per Schema`, `## File Rules`) and indented four spaces.**
  Rationale: house style in this file (`## Live Schema Validation`, `## Data Uploads`, `## Access Control`). The `heading-case` rule cannot flag them: its regex is anchored at line start (`headingcase.go`, `headingRe` = `^(#{1,6})[ \t]+(.+?)[ \t]*$`) and every heading inside an `<Update>` block is indented four spaces. This is the same reasoning recorded in `plans/completed/20260721-data-uploads-product-changelog.md`.
  Date/Author: 2026-08-13 / plan author.
- **D7: No `<Dropdown title="Improvements">` or `<Dropdown title="Fixes">` section.**
  Rationale: the brief forbids one without substantiated content, and there is none — `docs/changelog/cli.mdx:68-76` has two small CLI-only items already published on the CLI changelog, and neither is a product-surface change.
  Date/Author: 2026-08-13 / plan author.
- **D8: No existing entry is touched, including the July 20 entry's "upload rule" prose.**
  Rationale: D6 of `plans/completed/20260811-rename-upload-rules-to-file-rules.md` — changelogs are historical records; entries describe what shipped under the name it shipped under. Only the new entry uses "file rule". `git diff main -- docs/changelog/product.mdx` must show a pure insertion.
  Date/Author: 2026-08-13 / plan author.
- **D9: The `<Update>` label is "August 14, 2026".**
  Rationale: specified by the requester as the ship date. Precedent: product entries are labeled with the ship date, not the authoring date (`plans/completed/20260721-data-uploads-product-changelog.md`, Decision Log).
  Date/Author: 2026-08-13 / plan author.

## Outcomes & Retrospective

**What landed.** A single 72-line pure insertion into `docs/changelog/product.mdx` — one `<Update label="August 14, 2026">` block sitting above the August 6 entry, with exactly two `##` sections (`## Multiple Configs per Schema`, `## File Rules`), each carrying a bold-lead paragraph, `**Term:**` bullets, reference-doc links, and one `<Framed>` placeholder. No existing entry byte changed (D8 held: `git diff main` shows 72 insertions / 0 deletions), no `<Dropdown>` was added (D7), and `cspell.json` was untouched. The plan file was updated in the same commit.

**What the version-requirement investigation concluded.** Both D2 versions survived re-verification unchanged, so the entry cites Miru CLI `v0.11.0` (both features) and Miru Agent `v0.10.1` (file rules only), each backed by a dated entry in this repo's own `docs/changelog/cli.mdx` / `docs/changelog/agent.mdx`. M1's escape hatch — drop the version line if an entry reverted to `*Unreleased*` — was not triggered. Sibling repos remain absent from this workbench, so the four D3 gate questions moved to the PR body rather than into the changelog.

**Whether the placeholders survived both checks.** Yes, on the first run, exactly as D5 reasoned: `image-domain` only asserts the `https://assets.mirurobotics.com/` prefix and `mint validate` never resolves remote assets, so a not-yet-uploaded PNG is invisible to both. The fallback (moving `<Framed>` inside its comment) was never exercised. The two `TODO(assets)` markers grep cleanly and are the outstanding pre-publish action — the user adds the images himself.

**Final state.** CI green on the pushed head for `changes`, `lint`, and `shell-tests`; `audit` also green, and provably not a regression risk since `package.json` and `pnpm-lock.yaml` are byte-identical to `main`. The custom-linter jobs skipped as predicted (`tools/lint/**` untouched). `$preflight` reports CLEAN. The PR is left **in draft** by design — taking it out of draft is the orchestrator's call after the gate, and it should not happen until the two images are uploaded.

**Retrospective.** The plan's Context section was accurate enough that M0 found nothing to revise, which is what made it safe to write the prose against pre-verified `file:line` citations rather than re-reading the docs mid-draft. The two prose edits that were needed both came from the same failure mode: the Plan of Work had quietly upgraded a *quoted* fact into an *inferred* one (a schema name borrowed from a different example, and a "the file path comes from the slot" derivation the docs never state). Worth noting that `overview.mdx` contradicts itself on slot count — three in the code sample, "four" in the following sentence — which is a real docs bug this task routed around rather than fixed; it belongs in its own change.

## Context and Orientation

The repo is the Mintlify docs site. Content root is `docs/`; the repo root holds `package.json`, `cspell.json`, `scripts/`, `tools/lint/`, and `plans/`. There is no `AGENTS.md` or `CLAUDE.md` at any level — the house conventions live in the completed plans under `plans/completed/` and in `tools/lint/`.

### The file being edited

`docs/changelog/product.mdx` is a list of `<Update label="…">` blocks, newest first, inside `<div className="changelog-page">` (line 12). The first entry is `<Update label="August 6, 2026">` at line 14. `Dropdown`, `DropdownItem`, `Framed`, `LazyVideo`, and `Separator` are already imported (lines 7-10) and all are used by existing entries, so no import changes are needed and the import-usage lint rules stay satisfied.

Inside an `<Update>`, content is indented four spaces: a Title Case `##` feature heading, a short bold-lead paragraph, bulleted `**Term:**` detail lines, doc links written as `[Thing »](/path)`, and a media component per section. Media examples to copy the shape of:

- `docs/changelog/product.mdx:41-46` — `<Framed image="https://assets.mirurobotics.com/docs/changelog/26-08-06/xml-support.png" borderWidth="0px" innerRadius="8px" />`
- `docs/changelog/product.mdx:27-32` — `<LazyVideo alt="…" className="rounded-xl" aspectRatio="3840 / 2260" src="https://assets.mirurobotics.com/docs/changelog/26-08-06/schema-validation.mp4" />`

The precedent for presenting version requirements is the July 20 entry's `### Shipping upload rules` at `docs/changelog/product.mdx:196-201`:

```
    ### Shipping upload rules

    To use upload rules, you'll need to upgrade the Miru Agent and CLI to the following versions (or later):

    - [Miru Agent v0.10.0](/changelog/agent#v0-10-0) — detects files matching upload rules and uploads them to your bucket. …
    - [Miru CLI v0.10.2](/changelog/cli#v0-10-2) — adds the upload rule flags to `miru release create`
```

Anchor form is `/changelog/cli#v0-11-0` and `/changelog/agent#v0-10-1` (dots become hyphens).

### Item 1 facts — verified in this repo

- `docs/cfg-mgmt/primitives/schemas/overview.mdx:61-69` — `<ParamField path="instance slots" type="InstanceSlot[]">`: "The valid file system destinations that config instances for this schema are written to. Typically a schema defines only a single required slot. However, schemas may define several instance slots when the same validation applies to several different config instances in a release."
- `:107-109` — `## Instance slots`: "An instance slot defines one or more valid config instances that are valid for a given schema. Instance slots are defined as an annotation within the schema itself."
- `:115-176` — a `<CodeGroup>` with three tabs showing the annotation in each language: `x-miru-instance-slots:` (JSON Schema), `@miru_instance_slot(…)` (CUE), `instance_slots:` (Opaque), all for a three-slot `Motor Controller` schema.
- `:178` — "A schema with multiple slots allows multiple config instances to be written to different destinations. This structure allows a `motor-controller` schema to serve four identical motor controllers."
- `:181` — "All instance slots share the exact same validation. So two files that need different validation belong to two config types, not to two slots of one schema."
- `:187-231` — the five slot fields: `key` (required, immutable, unique within a schema), `name` (required), `filepath` (required, absolute, unique within a schema), `required` (required boolean — "If `true`, a deployment must include a config instance for this slot. If `false`, a deployment may optionally include a config instance for this slot."), `description` (optional).
- `docs/cfg-mgmt/primitives/config-instances.mdx:30-36` — `<ParamField path="slot key" type="string">`: "The key of the instance slot on the config schema that this instance fills." Examples `default`, `left-arm`, `right-arm`. `:22` is the `file path` field.
- `docs/cfg-mgmt/create-a-release.mdx:72` — "…the `instance slots` annotation, which lets one schema write several files by declaring a destination for each."
- `docs/snippets/references/cli/releases/create/schema-annotations.mdx:74-130` — the CLI-reference `instance slots` annotation, including the `<Warning>` at `:81` ("A schema may declare an instance file path OR instance slots. Declaring both will error.") and `:84` ("Omitting it will either use the instance file path annotation, or create a single slot with default values."). Rendered at `/references/cli/release-create`.

**Not in this repo, therefore not in the entry:** any dashboard/UI behavior (the completed plan took the documented skip path on all three dashboard-gated edits — `plans/completed/20260810-document-instance-slots.md`, Decision Log); a synthesized `default` slot; deployment slot cardinality; the 128-character `filepath` cap; slot-language coverage on the language pages.

### Item 2 facts — verified in this repo

- `docs/data-uploads/primitives/file-rules.mdx:55-77` — the rule file format: `name`, `source: {glob, stability_window_secs}`, `# optional block` `upload: {collection_slug, bucket, path}`, `# optional block` `retention: {ttl_secs, require_upload}`. Section anchors `#sources`, `#uploads`, `#retention`.
- `docs/snippets/file-rules/retention.mdx:1-2` — "A rule's `retention` block governs deletion of the local copies of matching files. When absent, the device retains matching files indefinitely (Miru never deletes them)."
- `docs/snippets/file-rules/retention.mdx:6-19` — `require_upload` ("Must be set when the rule has an `upload` block") and `ttl_secs` (marked `required`; "Use `0` to delete the file as soon as it becomes eligible").
- `docs/changelog/cli.mdx:16-32` — the `## Breaking changes` `<Dropdown title="File rule flags">` diff: `--upload-rule`/`--upload-rules` → `--file-rule`/`--file-rules`, "There are no deprecated aliases, so update your invocations and CI pipelines".
- `docs/changelog/cli.mdx:34-62` — `<Dropdown title="Rule file format">`, ending "Both `upload` and `retention` are optional. Omit `upload` for a rule that only reclaims local disk; omit `retention` and the device retains files on disk indefinitely after uploading. `ttl_secs` is required inside `retention`, and `require_upload` must be set exactly when the rule has an `upload` block."
- `docs/snippets/references/cli/releases/create/scopes.mdx:11` — `file_rules:manage`. The old scope ids were `upload_rules:*` (`plans/completed/20260811-rename-upload-rules-to-file-rules.md`: the backend migration rewrites stored scope ids `upload_rules:%` → `file_rules:%`; `upload_collections:manage` is unchanged).
- `docs/changelog/agent.mdx:6-12` — `# v0.10.1`, *August 12, 2026*, Features: "Added support for the latest file rule capabilities, including retention-only rules and a configurable delay period before deletion".
- Redirects for the old URLs landed in `docs/docs.json` (`"source": "/data-uploads/primitives/upload-rules"`, around line 411), so old links keep working.

### Lint and CI toolchain

- `./scripts/lint.sh` (= `pnpm run lint`) runs, in order: the Go MDX prose linter (`tools/lint`), ESLint with `--max-warnings=0` over every MDX file, CSpell against `cspell.json`, then `mint openapi-check` per spec. Final line on success: `All documentation lint checks passed.` Needs `pnpm`, Go, and network.
- `pnpm run validate` = `cd docs && mint validate` — build validation and link checking.
- `pnpm run test:lint` = `tests/test-lint.sh` — smoke tests of the linter itself.
- `./scripts/preflight.sh` chains: `test:lint`, the Go linter, the Go coverage gate, `./scripts/lint.sh`, `./scripts/audit.sh`, and the `bats` shell tests. Note it does **not** call `pnpm run validate`; run that separately.
- Prose rules that matter for this entry: **`no-double-dash`** — `--` in prose must be an em dash `—`; inline code is masked (`tools/lint/linter/analysis/scanner.go:154`, `:190`), so `` `--file-rule` `` is fine and a bare double dash is a violation. **`heading-case`** — see D6; not reachable here. **`image-domain`** — every `src`/`image`/`background`/`poster` image URL must start with `https://assets.mirurobotics.com/`.
- CI (`.github/workflows/ci.yml`): `changes`, `lint` (test:lint + `scripts/lint.sh` + `pnpm run validate`), `audit`, `shell-tests`. The two custom-linter jobs are skipped because `tools/lint/**` is untouched.
- **`audit`**: `scripts/audit.sh` runs `pnpm audit --ignore-registry-errors` against the `ignoreCves` list in `package.json` (9 entries as of `e001038`). `GHSA-5p4m-2wfm-xmqj` is a known pre-existing red and is never a regression from this change; the rename plan later observed `audit.sh` had gone green locally as `ignoreCves` grew. **Measure, do not assume**: record the state on pristine `main` and on the branch head at the same moment, and only call a difference a regression.

## Plan of Work

One edit: insert the block below into `docs/changelog/product.mdx` between line 12 (`<div className="changelog-page">`) and line 14 (`<Update label="August 6, 2026">`), keeping one blank line on each side. Nothing else in the file changes.

The exact block (four-space indentation inside the `<Update>` tags; adjust only if M1 changes the version requirements):

    <Update label="August 14, 2026">
        ## Multiple Configs per Schema

        One config schema can now write **several config files** to a device, using
        **instance slots**. Each slot is a distinct file system destination for the same
        schema, and every config instance binds to exactly one slot. A robot with four
        identical motor controllers can use a single `motion-control` schema with four
        slots instead of four near-identical config types.

        - **Instance slots:** Declared as an annotation inside the schema itself. Each
          slot has a `key`, a `name`, an absolute `filepath`, and `required`, plus an
          optional `description`
        - **Same validation, several files:** Every slot on a schema shares the exact
          same validation. Two files that need different validation still belong to two
          config types, not to two slots of one schema
        - **Slot keys:** Slot keys and file paths are unique within a schema, and every
          config instance records the `slot key` of the slot it fills. A config
          instance's file path comes from the slot it's bound to
        - **Required and optional slots:** A deployment must include a config instance
          for every slot marked `required`, and may include one for the rest
        - **One or the other:** A schema declares an instance file path *or* instance
          slots. Declaring both is an error. Declaring neither gives the schema a single
          slot with default values

        Instance slots are available in [Miru CLI v0.11.0](/changelog/cli#v0-11-0).

        [Instance slots »](/cfg-mgmt/primitives/schemas/overview#instance-slots)

        {/* TODO(assets): asset not yet uploaded — instance-slots.png must exist at the URL below before this entry publishes */}
        <Framed
            image="https://assets.mirurobotics.com/docs/changelog/26-08-14/instance-slots.png"
            borderWidth="0px"
            innerRadius="8px"
        />

        ## File Rules

        The upload rule is now the **file rule**. A rule still watches a glob on the
        device, but uploading the matches and reclaiming their disk space are now two
        independent, optional halves of one object.

        - **Source:** Unchanged — the glob to watch, and the stability window that
          decides when a matching file has finished being written
        - **Upload:** Optional. Where matching files go: the upload collection, the
          bucket, and the object path. Omit it for a rule that only reclaims local disk
        - **Retention:** Optional. How long the device keeps matching files. `ttl_secs`
          is required inside a retention block, and `require_upload` is set exactly when
          the rule has an `upload` block. Omit retention and the device keeps matching
          files indefinitely
        - **Renamed flags:** `miru release create` takes `--file-rule` and `--file-rules`
          in place of `--upload-rule` and `--upload-rules`. There are no deprecated
          aliases, so update your invocations and CI pipelines
        - **Renamed scopes:** The `upload_rules:*` API key scopes are now `file_rules:*`

        To use file rules, you'll need to upgrade the Miru CLI and Agent to the following
        versions (or later):

        - [Miru CLI v0.11.0](/changelog/cli#v0-11-0) — renames the rule flags and accepts
          the restructured rule file
        - [Miru Agent v0.10.1](/changelog/agent#v0-10-1) — supports the latest file rule
          capabilities, including retention-only rules

        [File rules »](/data-uploads/primitives/file-rules)

        {/* TODO(assets): asset not yet uploaded — file-rules.png must exist at the URL below before this entry publishes */}
        <Framed
            image="https://assets.mirurobotics.com/docs/changelog/26-08-14/file-rules.png"
            borderWidth="0px"
            innerRadius="8px"
        />
    </Update>

Prose constraints to hold while editing: no bare double dash outside backticks (`no-double-dash`); em dashes are `—`; no claim about dashboard UI; no claim about deployment slot cardinality; no mention of a synthesized `default` slot beyond "a single slot with default values", which is quoted behavior from `schema-annotations.mdx:84`.

## Concrete Steps

All commands run from `/home/ben/miru/workbench5/repos/docs`.

### M0 — orient and re-verify

1. Confirm the branch and a clean tree:

       git -C /home/ben/miru/workbench5/repos/docs branch --show-current   # expect: docs/product-changelog-26-08-14
       git -C /home/ben/miru/workbench5/repos/docs status --short          # expect: only this plan file

2. Re-verify the Context facts (the file wins over this plan — if any output differs, update Context and Surprises & Discoveries before writing prose):

       grep -rn "x-miru-instance-slots\|@miru_instance_slot\|instance_slots" docs --include='*.mdx' | cut -d: -f1 | sort | uniq -c
       grep -n "instance slots\|slot key\|Instance slots" docs/cfg-mgmt/primitives/schemas/overview.mdx docs/cfg-mgmt/primitives/config-instances.mdx
       grep -n "ttl_secs\|require_upload" docs/snippets/file-rules/retention.mdx
       grep -n "file_rules:manage" docs/snippets/references/cli/releases/create/scopes.mdx

   Expect: two files for the annotation grep (`schemas/overview.mdx`, `schema-annotations.mdx`); both `retention` fields present; one `file_rules:manage` hit.

3. Idempotence guard:

       grep -c 'label="August 14, 2026"' docs/changelog/product.mdx   # expect: 0

   If it prints 1, the insertion is already done — skip to M3.

### M1 — decide the version requirements (do this before writing prose)

4. Read both changelogs in this repo and write down what they substantiate:

       grep -n "^# v" docs/changelog/cli.mdx | head -5
       sed -n '1,80p' docs/changelog/cli.mdx
       sed -n '1,30p' docs/changelog/agent.mdx

   Expected at `e001038`: `docs/changelog/cli.mdx:10` `# v0.11.0` dated *August 12, 2026*, with a `## Breaking changes` section covering the file-rule flags/format/scopes and a Features bullet at `:66` for the instance slots annotation; `docs/changelog/agent.mdx:6` `# v0.10.1` dated *August 12, 2026*, Features bullet naming the latest file rule capabilities.

5. Apply the rule from D2/D3:
   - Both features may cite **Miru CLI v0.11.0**; file rules may additionally cite **Miru Agent v0.10.1**. These are the only version claims this repo can back.
   - If either entry has changed or been removed since `e001038` (for example, the `v0.11.0` date reverted to `*Unreleased*`), **drop the corresponding version line from the changelog entirely** and say so in the PR body. Do not soften it into vague wording, and do not invent a version.
   - Record the outcome as a dated bullet in Surprises & Discoveries, including the sibling-repo unavailability (`ls /home/ben/miru/workbench5/repos` → `docs`) and the four open gate questions listed in D3.

### M2 — insert the entry

6. Edit `docs/changelog/product.mdx`: insert the Plan of Work block between `<div className="changelog-page">` and `<Update label="August 6, 2026">`, one blank line on each side. Change no existing line.

7. Confirm the shape of the diff:

       git diff --stat -- docs/changelog/product.mdx    # expect: insertions only, 0 deletions
       grep -n 'TODO(assets)' docs/changelog/product.mdx   # expect: exactly 2 hits
       grep -c '<Update label=' docs/changelog/product.mdx # expect: one more than on main

### M3 — validate

8. Run the checks:

       pnpm install --frozen-lockfile
       pnpm run test:lint     # expect: pass
       ./scripts/lint.sh      # expect final line: "All documentation lint checks passed."
       pnpm run validate      # expect: build validation passed, no broken links

   Expected `lint.sh` stage banners: `== MDX Prose ==`, `== ESLint (MDX) ==`, `== CSpell ==`, `== OpenAPI ==`.

   - If CSpell flags a word from the new entry, add it to the `words` array in `cspell.json` (loose alphabetical order) and re-run. Add nothing that was not flagged.
   - If the prose linter flags a line, fix the text — never disable a rule and never add a `lint-ignore` directive.
   - **If either check objects to the not-yet-uploaded assets**, take the D5 fallback: move each `<Framed …/>` inside its `{/* TODO(assets): … */}` comment, re-run both checks, and record the fallback in Surprises & Discoveries.

9. Optional render check: `pnpm run dev`, open `/changelog/product`, confirm the new entry renders at the top above August 6 with both sections and their links.

### M4 — commit, push, draft PR

10. Commit (signed; one commit for the entry, plus this plan file):

        git add docs/changelog/product.mdx plans/
        # add cspell.json too only if step 8 changed it
        git commit -S -m "docs(changelog): add the August 14, 2026 product entry"

    Verify the signature: `git log --show-signature -1` shows a good signature.

11. Push and open a **draft** PR:

        git push -u origin docs/product-changelog-26-08-14
        gh pr create --draft --base main \
          --title "docs(changelog): add the August 14, 2026 product entry" \
          --body-file <the body described below>

    The PR body must contain, at minimum: what the entry covers; the two source PRs (#150 `73d8e05`, #154 `f5ffbab`); the `TODO(assets)` placeholders and that the assets are not yet uploaded; the pre-existing red `audit` job caveat with its measured state; and the four unsubstantiable items from D3 (sibling repos unavailable; the three-part instance-slots gate; PR #154's stable-`v0.11.0` gate versus `f5ffbab` on `main`; PR #153's body-versus-diff mismatch).

### M5 — CI and exit draft

12. `$preflight` on the pushed head. Do not take the PR out of draft, and do not report this task complete, until it reports `CLEAN` (see Validation and Acceptance).

## Validation and Acceptance

1. `./scripts/lint.sh` exits 0 and prints `All documentation lint checks passed.` (it passes before the change too — this must not regress it).
2. `pnpm run validate` exits 0 with no broken links and no invalid MDX.
3. `pnpm run test:lint` exits 0.
4. `docs/changelog/product.mdx` contains exactly one `<Update label="August 14, 2026">` block, positioned above `<Update label="August 6, 2026">`, and `git diff main -- docs/changelog/product.mdx` shows a pure insertion — every pre-existing entry byte-identical, including the July 20 entry's "upload rule" prose (D8).
5. The entry has exactly two `##` sections, `## Multiple Configs per Schema` and `## File Rules`, each with a bold-lead paragraph, `**Term:**` bullets, at least one reference-doc link, and one media component. No `<Dropdown>` (D7).
6. `grep -c 'TODO(assets)' docs/changelog/product.mdx` returns `2`, and both marker lines are MDX comments (`{/* … */}`), not HTML comments.
7. Both placeholder URLs begin with `https://assets.mirurobotics.com/docs/changelog/26-08-14/` — verified indirectly by `image-domain` passing in step 1.
8. Every link in the entry resolves: `/cfg-mgmt/primitives/schemas/overview#instance-slots`, `/data-uploads/primitives/file-rules`, `/changelog/cli#v0-11-0`, `/changelog/agent#v0-10-1` (plus any others added) — confirmed by `pnpm run validate`.
9. The entry makes **no** claim about dashboard UI, deployment slot cardinality, a synthesized `default` slot beyond the quoted "a single slot with default values", or the `filepath` length cap.
10. Every version number in the entry appears in this repo's own `docs/changelog/cli.mdx` or `docs/changelog/agent.mdx` (D2). Anything that does not is absent from the entry and present in the PR body (D3).
11. `git log --show-signature -1` shows a good signature on the commit.
12. CI on the pushed head is green for `changes`, `lint`, and `shell-tests`. The `audit` job is either green or red with an advisory set identical to pristine `main` re-measured at the same moment (`GHSA-5p4m-2wfm-xmqj` is a known pre-existing red and never a regression from this change). The custom-linter jobs are skipped, since `tools/lint/**` is untouched.
13. **Gate: `$preflight` must report `CLEAN` — CI green on the pushed branch head — before the PR leaves draft and before this task is reported complete.** A green local run is not sufficient; a green run on an older commit is not sufficient. Confirm job states against the pushed head SHA (`gh pr checks --watch`). This is the only definition of CLEAN in this plan.

## Idempotence and Recovery

- The insertion is guarded by the M0 grep; re-running the steps never duplicates the entry. Every verification step is read-only and safe to repeat.
- If the edit goes wrong before committing: `git checkout main -- docs/changelog/product.mdx`, then redo M2.
- If a bad version was committed: `git revert -S <sha>`. Do not force-push a pushed branch.
- `pnpm install --frozen-lockfile` is idempotent and never rewrites `pnpm-lock.yaml`; re-run it after a network hiccup.
- No repository other than `repos/docs` is written to. The only sibling-repo commands anywhere in this plan would be read-only, and none are needed — sibling repos are not present in this workbench.

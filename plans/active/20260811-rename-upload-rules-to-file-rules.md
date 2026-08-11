# Document file rules instead of upload rules

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `/home/ben/miru/workbench1/repos/docs` (branch `docs/rename-upload-rules-to-file-rules`, base `main` at `1a386cd`) | read-write | Every edit lands here: MDX pages, snippets, `docs/docs.json` nav + redirects, changelog link targets, this plan file. |
| `/home/ben/miru/workbench1/repos/openapi` | read-only | Authoritative source for the file-rule object shape. Current spec files win over older plan prose. |
| `/home/ben/miru/workbench1/repos/cli-private` | read-only | Authoritative source for the rule-file YAML format and the `miru release create` flags. |
| `/home/ben/miru/workbench1/repos/backend` | read-only | Authoritative source for API-key scope ids and bucket-delete/archive constraints. |
| `/home/ben/miru/workbench1/repos/agent` | read-only | Authoritative source for what the device actually does with `retention`. Constrains what this plan may claim. |

This plan lives in `repos/docs/plans/` because every changed file is in `repos/docs`. No code changes in any other repository.

## Purpose / Big Picture

The platform renamed the "upload rule" to a **file rule** and restructured its body. An upload rule was "watch this glob, upload matches to this bucket, optionally delete after upload." A file rule is "watch this glob; here is optionally where the matches go (`upload`), and optionally how long the device keeps them (`retention`)." Uploading and reclaiming disk are now two independent, optional halves of one object.

The docs still describe the old object end to end: a rule that belongs to an upload collection, with a `source` and a `destination`, whose `destination.delete_policy` enum decides deletion. Every reader-facing surface — the primitive page, the authoring guide, the CLI reference, the agent file-system-permissions page, the releases primitive, the access-control matrix — is wrong in vocabulary, in shape, and in the CLI flags it tells people to type.

After this change a reader can:

- Open `/data-uploads/primitives/file-rules` and learn the whole object: `name`, `source`, optional `upload`, optional `retention`, `digest`, immutability, Git provenance.
- Open `/data-uploads/define-file-rules` and copy a rule file whose YAML actually parses with the shipped CLI, and a `miru release create` invocation whose flags actually exist.
- Follow any old bookmarked URL (`/data-uploads/primitives/upload-rules`, `/data-uploads/define-upload-rules`) and land on the renamed page instead of a 404.
- Read `/admin/users/access-control` and `/references/cli/release-create` and see the real scope ids (`file_rules:manage`) and real flags (`--file-rule`, `--file-rules`).

Observable outcome: `grep -rin 'upload rule\|upload-rule\|upload_rule' docs/` returns hits **only** inside historical changelog prose (see Decision D6); `./scripts/lint.sh` and `pnpm run validate` pass; the two redirects resolve; and `$preflight` reports `CLEAN` on the pushed head.

## Progress

- [x] Milestone 0 — Orientation and baseline (branch, install, baseline lint/audit, re-verify upstream facts) — 2026-08-11, no commit (read-only)
- [x] Milestone 1 — The file-rule primitive page, snippets, and definition snippet — 2026-08-11, `a9c3be0`
- [ ] Milestone 2 — The authoring guide (`define-file-rules`) and nav + redirects
- [ ] Milestone 3 — Inbound repoints across data-uploads, developers, primitives, cfg-mgmt, admin
- [ ] Milestone 4 — CLI reference (flags, usage, scopes) and the unreleased CLI changelog entry
- [ ] Milestone 5 — Changelog link targets (hrefs only, prose untouched)
- [ ] Milestone 6 — Validation, push, `$preflight` CLEAN, PR out of draft

## Surprises & Discoveries

- **M0 (2026-08-11): `audit.sh` is GREEN locally, contrary to the known pre-existing-red condition.**
  `./scripts/audit.sh` exits 0 with `No known vulnerabilities found`. The `pnpm.auditConfig.ignoreCves` list in `package.json` now carries 8 CVEs, which evidently covers the advisories that previously made the job red. So `./scripts/preflight.sh` can run end to end; the by-hand fallback stage list in Concrete Steps M0.3 is not needed. Baseline `./scripts/lint.sh` and `pnpm run validate` also exit 0.

- **M0 (2026-08-11): every upstream fact in the Context section re-verified unchanged.** No Context edits were needed.
  - `repos/openapi` HEAD for `apis/configs/components/schemas/file-rule.yaml` is still `651db31` (`refactor(file-rule)!: make retention optional and drop the policy concept (#234)`), the same commit the plan was written against. Re-read both `schemas/file-rule.yaml` and `requests/file-rule.yaml`: `name`/`digest`/`source`/optional `upload`/optional `retention`, `retention.ttl_secs` required inside the block, `require_upload` present exactly when `upload` is, 9-attempt abandon, workspace-scoped digest dedup, and the guarantee-vs-enforcement two-phase wording all match the plan verbatim.
  - `repos/cli-private`: `internal/domain/filerules/spec.go` HEAD is still `8ba7471`; the YAML shape (`name`, `source.{glob,stability_window_secs}`, `upload.{collection_slug,bucket,path}`, `retention.{ttl_secs,require_upload}`) matches the plan's Context YAML field for field. `internal/commands/release/flags.go` registers `file-rule` and `file-rules` with no aliases.
  - `repos/agent` HEAD is still `2b24b4f`; `agent/src/models/file_rule.rs:76` still reads `impl From<backend_client::BaseUploadRule> for FileRule`, and the workspace version is still `0.10.0`. D3's premise holds, so the D3 **fallback** applies (see Decision Log addendum).
  - Inventory baseline: `grep -rn "upload rule\|Upload rule\|Upload Rule\|upload-rule\|upload_rule\|uploadRule" docs/ | wc -l` → **87** hits across the exact 21 files the plan's inventory lists — no files added or dropped.

- **M1 (2026-08-11): the `importresolves` breakage came from Milestone 2's page, not Milestone 3's.**
  Concrete Steps M1.5 anticipated folding "steps 1–3 of Milestone 3" into the M1 commit to avoid committing red. In fact the only importer of the renamed snippets was `docs/data-uploads/define-upload-rules.mdx` (a Milestone 2 file); no Milestone 3 file imports them. Rather than merge M1 and M2 into one commit and lose the two commit messages the plan specifies, the M1 commit carries a **minimal three-line repoint** of that page's import block (`/snippets/upload-rules/{sources,destinations}.mdx` → `/snippets/file-rules/{sources,upload}.mdx`, identifier `Destinations` → `Upload`). The page's rename, rewrite, and the `Retention` import stay in Milestone 2. `./scripts/lint.sh` is green at `a9c3be0`.
  Note `pnpm run validate` is *not* run at M1: `docs/docs.json` still names the now-deleted `data-uploads/primitives/upload-rules` nav path, which Milestone 2 fixes. This is the M1/M2 coupling the Idempotence section already calls out.

- **M1 (2026-08-11): the retention snippet's cross-references are absolute, not `#sources`/`#uploads`.**
  `snippets/file-rules/retention.mdx` is rendered on two pages whose heading levels and slugs differ (`## Sources` on the primitive page, `### Source` on the authoring guide), so bare fragment links would resolve on one page and 404 on the other. Both links are written as `/data-uploads/primitives/file-rules#sources` and `#uploads` in full.

- **M0 (2026-08-11): the D1 gate is NOT satisfied and will not be during this task.**
  `git tag --contains 8ba7471` → `v0.11.0-beta.1` only; `git tag --list 'v0.11.0'` → empty. There is no **stable** `v0.11.0` tag. Per D1 the PR therefore stays in **draft**. The CLI reference is still updated (Milestone 4) as D1 directs — the alternative D1 explicitly forbids is softening the docs to straddle both flag sets.

## Decision Log

- **D1: The CLI reference *is* updated, because the CLI has already adopted the rename.**
  The task brief assumed the CLI had not. It has. `repos/cli-private` `main` carries `8ba7471` — `feat!: bump cli-sdk to v0.12.1, migrate upload rules to file rules (#123)`, dated 2026-08-09 — which renames the flags `--upload-rule` → `--file-rule` and `--upload-rules` → `--file-rules` **with no deprecated aliases** (`internal/commands/release/flags.go`), and replaces the rule-file YAML shape (`internal/domain/filerules/spec.go`). `git tag --contains 8ba7471` → `v0.11.0-beta.1`. Separately, `docs/changelog/cli.mdx` already carries a `# v0.11.0` entry dated *August 12, 2026*, so this repo already stages v0.11.0 content ahead of the stable tag. Leaving the CLI reference on `--upload-rule` would document a flag that no longer exists in the version the docs already announce.
  **Gate:** before `gh pr ready` (and before any promotion past `staging`), re-run `cd /home/ben/miru/workbench1/repos/cli-private && git tag --list 'v0.11.0'` and confirm a **stable** `v0.11.0` tag exists and contains `8ba7471` (`git tag --contains 8ba7471`). If v0.11.0 slips or ships without the rename, hold the PR in draft — do not soften the docs to straddle both flag sets.

- **D2: The CLI reference snippets are hand-written, not generated.**
  Verified: nothing under `scripts/` or `.github/workflows/` reads or writes `docs/snippets/references/`, `repos/cli-private` has no docs-export target, and the files' entire history is two hand edits (`git log -- docs/snippets/references/cli/releases/create/flags.mdx` → `f1e2019`, `23cb824`). So there is no generator whose output would be clobbered, and no "document the CLI as it ships" conflict beyond the version gate in D1. They are edited by hand here.

- **D3: `retention` is documented as a field reference, but no device-side enforcement beyond delete-after-upload is promised.**
  `repos/agent` `main` (`2b24b4f`) restructured its *internal* model into `FileRule`/`FileRuleRetention` (`agent/src/models/file_rule.rs`) but still deserializes the **old** wire object (`impl From<backend_client::BaseUploadRule> for FileRule`, mapping `delete_policy: after_upload` → `retention { require_upload: true, ttl_secs: 0 }`). The only deletion path is `agent/src/upload/executor.rs::delete_source_file`, which fires **only** for the exact value `require_upload: true, ttl_secs: 0`. Workspace version is still `0.10.0` (`Cargo.toml:13`) and no agent release ships file rules yet.
  Therefore: document `retention.require_upload` and `retention.ttl_secs` as the platform contract (they are real, accepted, and stored), state the two-phase model in the API's own terms, and add a `<Note>` that local deletion is performed by the Miru Agent. Do **not** write "the device deletes the file `ttl_secs` after it becomes eligible" as shipped behavior for `ttl_secs > 0`, and do **not** describe retention-only (no `upload` block) rules as enforced on-device.
  **Open Question (must be resolved before `gh pr ready`):** the Miru Agent version that enforces general `ttl_secs` and retention-only rules. If it is known, name it in the `<Note>` with a `/changelog/agent#v-x-y-z` link, matching how `product.mdx` gates uploads on Agent v0.10.0. **If it is not known, take the fallback**: keep the `<Note>` version-free ("Local deletion is performed by the Miru Agent; deleting a file requires write access to its parent directory — see file system access") and say nothing about timing. The fallback is a complete, shippable state; do not block on the answer.

- **D3-resolution (2026-08-11): the fallback is taken.** Re-verification at M0 confirmed `repos/agent` still deserializes `BaseUploadRule` and is still at workspace version `0.10.0`, so no shipped Miru Agent version enforces general `ttl_secs` or retention-only rules. The Open Question ("which agent version?") therefore has **no answer to give**, and the plan's own sanctioned fallback applies: the retention `<Note>` is version-free, names no timing, and says only that local deletion is performed by the Miru Agent and requires write access to the parent directory, linking `/developers/agent/filesys-access#data-uploads`. Per D3 this is "a complete, shippable state".

- **D4: The section stays `/data-uploads/**` and the nav group stays "Data Uploads".**
  Only two page slugs change (below). Renaming the whole section would touch every page under it, need a wildcard redirect, and re-open a product-naming question the platform has not asked for — `upload collections`, `uploads`, and `buckets` all keep their names upstream (the openapi rename plan explicitly left them alone; the CLI PR body says "Upload *collections* are untouched").

- **D5: Page renames and redirects.**
  `docs/data-uploads/primitives/upload-rules.mdx` → `file-rules.mdx`; `docs/data-uploads/define-upload-rules.mdx` → `define-file-rules.mdx`. Both old URLs are live in production, so both get a `redirects` entry in `docs/docs.json`. The `redirects` lint rule (`tools/lint/linter/redirects`) requires the `source` to *not* resolve to a real page and the `destination` to resolve — both hold only after the old files are deleted and the new ones exist, so redirects and renames must land in the same commit.

- **D6: Historical changelog entries keep their prose; only their link `href`s are repointed.**
  `docs/changelog/cli.mdx` `# v0.10.2`, `docs/changelog/agent.mdx` `# v0.10.0`, and the `docs/changelog/product.mdx` 2026-07-23 `<Update>` describe what shipped at the time under the name it shipped under. Rewriting them would falsify the record (precedent: `plans/completed/20260810-document-instance-slots.md` — "Changelogs are historical records — never rewrite an existing entry"). But their links point at `/data-uploads/primitives/upload-rules`, which will no longer exist; a redirect covers a browser, not a docs-internal link check. So: change **only** the URL inside `](…)`, never the visible words. The `# v0.11.0` CLI entry is *not* historical — no stable `v0.11.0` tag exists yet and its date is in the future — so it is amended in place with a Breaking changes section (Milestone 4).

- **D7: `docs/references/platform-api/*.yaml` and `docs/references/device-api/**` are not touched.**
  Verified they contain zero occurrences of `upload_rule`/`UploadRule`/`file_rule` — the published API versions predate the rename. They are vendored from Stainless; hand edits are overwritten.

## Outcomes & Retrospective

- (fill in on completion)

## Context and Orientation

### Repo layout and tooling

The Mintlify site root is `docs/` inside this repo (pages at `docs/docs/...` from the repo root; site URLs drop the `docs/` prefix). Nav and redirects live in `docs/docs.json`. `./scripts/lint.sh` runs the Go prose linter (`tools/lint`), ESLint-MDX, CSpell, and `mint openapi-check`; `pnpm run validate` runs `mint validate` (the docs build, added to CI in `1566a65`). `./scripts/preflight.sh` runs the lint smoke tests, the Go linter and coverage gate, lint, audit, and the bats shell tests.

CI (`.github/workflows/ci.yml`) on a docs-only PR runs exactly four jobs: `changes`, `lint` (test:lint → lint.sh → **`pnpm run validate`**), `audit`, `shell-tests`. The two custom-linter jobs are gated off unless `tools/lint/**` changes — this plan must not change `tools/lint/**`.

Lint rules that bite this change specifically:

- **`redirects`** — every `redirects[].source` must NOT resolve to a real page and every `destination` MUST resolve. See D5.
- **`importresolves`** — a snippet import must resolve on disk. Renaming `docs/snippets/upload-rules/` breaks every importer in the same commit unless all imports are updated together.
- **`importsorted`** / **`importblock`** — import blocks are alphabetically sorted and contiguous. Renaming `UploadRuleDef` → `FileRuleDef` and `/snippets/upload-rules/…` → `/snippets/file-rules/…` **changes the sort order** in `docs/data-uploads/primitives/file-rules.mdx` and `docs/data-uploads/define-file-rules.mdx`. Re-sort deliberately; do not assume the old order still holds.
- **`heading-case`** — strict sentence case with a case-sensitive allowlist (`tools/lint/linter/headingcase/headingcase.go:47`: `API APIs CLI CI SDK SDKs CUE JSON MQTT TLS HTTPS REST GUI URL ACLs SSE OpenAPI AWS GCP GCS WIF STS IAM S3 ARN SigV4 Miru GitHub Agent Unix Git Python Schema Base Head Cloud Storage`). `YAML` and **`TTL`** are NOT allowlisted — `## File rules`, `## Sources`, `## Uploads`, `## Retention`, `## Git provenance` are fine; `## TTL and deletion` would fail. Version headings like `# v0.11.0` are exempt.
- **`no-double-dash`** — em dash `—` only in prose.
- **`image-domain`** — images must come from the approved asset domain. The three dashboard screenshots on the primitive page (`https://assets.mirurobotics.com/docs/v04/images/releases/upload-rule-*.png`) keep their URLs; asset filenames are not renameable from this repo. Only the alt text changes.

### The new object (authoritative: `repos/openapi` current spec files)

From `apis/configs/components/schemas/file-rule.yaml` and `apis/configs/components/requests/file-rule.yaml` — these **supersede** the older rename plan's prose and the intermediate `policy` / `delete_delay_secs` / `forever` / `until_stable` shapes from `20260725-restructure-file-rule-retention.md`, `20260807-retention-policy-value-rename.md`, `20260807-retention-two-phase-prose.md`, and `20260808-retention-optional.md`. The final shape is:

- `object: "file_rule"`, `id: file_rule_123`, `name`, `digest`, `source`, optional `upload`, optional `retention`, `created_at`, `updated_at`, expandable `file_rule_git_commits`.
- **`name`** (new, required): human-readable, **not unique** in a workspace, and it **participates in the digest** — renaming a rule mints a new rule.
- **`digest`**: file rules are immutable and deduplicated by digest **within a workspace** (the old docs said "within their collection" — that is now wrong).
- **`source`**: `glob` (required; absolute, ≤1024 bytes, no `..` segments, no empty segments) + `stability_window_secs` (default `60`).
- **`upload`** (optional; absent ⇒ retention-only rule): `upload_collection_id`, `upload_collection_name`, `bucket_id`, `bucket_name`, `path`. `path` **must contain `{upload_id}`**, ≤1024 bytes, balanced braces, no `..`/empty segments. Variables: `{device_id} {device_name} {file_name} {upload_id} {year} {month} {day} {hour} {minute}` — identical to the set already documented.
- **`retention`** (optional; absent ⇒ **the device keeps matching files indefinitely, Miru never deletes them**): `require_upload` (present exactly when the rule has an `upload` block; `true` ⇒ a file is never deleted before its upload is durably confirmed, and the device abandons an upload after **9 attempts**, so a permanently failing upload retains the file indefinitely) and `ttl_secs` (**required inside the block**; `0` ⇒ delete as soon as eligible).
- **Two-phase model** (use this vocabulary verbatim): the retention **guarantee** ends when the file goes quiescent (size and mtime unchanged for `source.stability_window_secs`) and, when `require_upload` is `true`, when its upload is durably confirmed; the file is then **eligible for deletion**. `ttl_secs` schedules the **enforcement** that follows — it defers the deletion, not the end of the guarantee.
- **Gone entirely**: `destination`, `destination.delete_policy`, `never`/`after_upload`, `delete_delay_secs`, `max_age_secs`, `policy`, `forever`, `until_stable`. None of these may appear in the new docs.
- Git provenance: `file_rule_git_commits[]` with `filepath` (the rule file's path relative to the repo root) + `git_commit_id`. Same accrual behavior as before.

### The rule file YAML (authoritative: `repos/cli-private/internal/domain/filerules/spec.go`)

```yaml
name: Robot Logs                       # required; unique across the loaded file set
source:
  glob: /var/log/robot/*.log           # required
  stability_window_secs: 60            # optional, default 60
upload:                                # optional
  collection_slug: robot-logs          # required inside upload
  bucket: my-uploads-bucket            # required inside upload
  path: logs/{device_id}/{upload_id}/{file_name}   # optional; server generates a default
retention:                             # optional
  ttl_secs: 604800                     # required inside retention
  require_upload: true                 # set exactly when `upload` is present
```

Note the asymmetry to document carefully: the **rule file** says `upload.collection_slug` and `upload.bucket` (names/slugs), while the **API object** returns `upload.upload_collection_id`/`_name` and `bucket_id`/`bucket_name`. That mirrors the old `collection_slug` / `destination.bucket` split and is not a mistake.

CLI-side constraints worth one sentence each: rule `name` must be unique across the file set loaded in one `miru release create` (this **replaces** the old "one rule per collection slug" invariant — two rules may now share a collection slug); `require_upload` without an `upload` block is rejected.

### API-key scopes (authoritative: `repos/backend`)

`tools/supabase/migrations/20260805120000_rename_upload_rules_to_file_rules.sql:111` rewrites stored scope ids `upload_rules:%` → `file_rules:%`, so `upload_rules:manage` is now **`file_rules:manage`**. `upload_collections:manage` is unchanged.

### ID prefix

The CLI's abbreviated id prefix changed `UPR` → `FLR` (`8ba7471` PR body: `✓ Lidar Logs · FLR-a1b2c (new)`). `docs/primitives/releases.mdx:38` still shows `UPR-FmoDN`-style examples.

### Complete inventory of files that mention the old concept

From `grep -ril "upload.rule\|uploadRule" . --exclude-dir=node_modules --exclude-dir=.git`, minus `plans/completed/**` (see Files that must NOT be touched):

**Renamed:** `docs/data-uploads/primitives/upload-rules.mdx`, `docs/data-uploads/define-upload-rules.mdx`, `docs/snippets/definitions/upload-rule.mdx`, `docs/snippets/upload-rules/{sources,destinations}.mdx`.
**Edited in place:** `docs/data-uploads/overview.mdx`, `docs/data-uploads/primitives/{uploads,buckets,upload-collections}.mdx`, `docs/snippets/definitions/upload-collection.mdx`, `docs/developers/agent/{overview,filesys-access}.mdx`, `docs/primitives/releases.mdx`, `docs/cfg-mgmt/create-a-release.mdx`, `docs/admin/users/access-control.mdx`, `docs/references/cli/release-create.mdx`, `docs/snippets/references/cli/releases/create/{flags,usage,scopes}.mdx`, `docs/docs.json`.
**Link targets only:** `docs/changelog/{cli,agent,product}.mdx` (plus the `# v0.11.0` amendment in `cli.mdx`).

### Files that must NOT be touched

- **`plans/completed/**`** (7 files match the grep) — historical documents. Precedent is explicit in the openapi plans ("Any edit to `plans/completed/` — historical documents stay as they are").
- **`docs/references/platform-api/*.yaml`, `docs/references/device-api/**`** — vendored, and they contain zero hits anyway (D7).
- **`tools/lint/**`** — would trigger two extra CI jobs plus a Go coverage gate; no rule change is needed. Verify with `git diff --stat main...HEAD -- tools/lint` printing nothing.
- **The historical prose in `docs/changelog/*.mdx`** — hrefs only (D6).
- **Asset image URLs** under `assets.mirurobotics.com` — not renameable from this repo.
- **The `/data-uploads/**` URL space and the "Data Uploads" nav group** (D4).

## Plan of Work

### Milestone 1 — the primitive page, its snippets, and the definition

**`docs/snippets/definitions/upload-rule.mdx` → `docs/snippets/definitions/file-rule.mdx`.** Rewrite:
a **file rule** is a standing instruction about files on a device — which files to match, optionally where to upload them, and optionally how long the device keeps them. Link `[bucket](/data-uploads/primitives/buckets)` and `[release](/primitives/releases)`. Delete the claim "Every rule belongs to an upload collection" — only upload-bearing rules name a collection.

**`docs/snippets/upload-rules/sources.mdx` → `docs/snippets/file-rules/sources.mdx`.** Mostly carries over. Change "the files to upload" → "the files this rule manages"; keep the glob constraints and the `<Warning>` linking `/developers/agent/filesys-access#data-uploads`; reword `stability_window_secs` so it describes quiescence generally (it now gates both upload eligibility and the end of the retention guarantee), keeping the `/data-uploads/primitives/uploads` link.

**`docs/snippets/upload-rules/destinations.mdx` → split.**
- `docs/snippets/file-rules/upload.mdx`: intro ("a rule's optional `upload` block declares which registered bucket matching files are written to, and the object path within it; omit it for a rule that only manages local retention"), then `collection_slug`, `bucket`, and `path` `<ParamField>`s. Carry the variable table and the `{upload_id}` `<Warning>` **verbatim**. **Delete the `delete_policy` `<ParamField>` entirely.**
- `docs/snippets/file-rules/retention.mdx` (new): the two-phase guarantee/enforcement paragraph from Context, then `<ParamField path="require_upload" type="boolean">` and `<ParamField path="ttl_secs" type="integer" required>`. Include the 9-attempt abandon fact and the "`0` deletes as soon as eligible" fact. Include the D3 `<Note>` (agent-version fallback applies).

**`docs/data-uploads/primitives/upload-rules.mdx` → `docs/data-uploads/primitives/file-rules.mdx`**, `title: "File rules"`. Section outline (headings are anchor contracts — other pages link to them):

| Section | Content |
| --- | --- |
| intro | `<FileRuleDef />` + one paragraph: a rule always has a `source`; `upload` and `retention` are independent optional halves; a rule with neither is not useful, a rule with both uploads and then reclaims. Created by releasing rule files via the CLI → `/data-uploads/define-file-rules`. |
| `## Properties` | `name` (immutable, not unique, participates in the digest → renaming mints a new rule), `digest` (dedup **within a workspace**), `source` → `#sources`, `upload` → `#uploads`, `retention` → `#retention`. Drop the standalone `upload collection` property — the collection now lives inside `upload`; say so and link `/data-uploads/primitives/upload-collections`. |
| `## File format` | The rule-file YAML from Context verbatim, plus the `collection_slug`/`bucket`-by-name note and a pointer to `/data-uploads/define-file-rules`. |
| `## Sources` | `<Sources />` |
| `## Uploads` | `<Upload />` (renamed from Destinations — **anchor changes `#destinations` → `#uploads`**) |
| `## Retention` | `<Retention />` (new) |
| `## Immutability` | Immutable; created once, shipped in a release, never edited in place; change = new rule + release + deploy. Add: because `name` is part of the digest, a rename is a new rule. |
| `## Git provenance` | Unchanged in substance; keep the `/data-uploads/define-file-rules` link. |
| `## Create a file rule` `<PublisherBadge />` | Unchanged in substance. |
| `## View a file rule` | Unchanged; keep the three `assets.mirurobotics.com` screenshot URLs as-is, update only alt text (`Release File Rules List`, `Release File Rule Details`). |

Re-sort the import block after renaming the imported identifiers (`importsorted`).

### Milestone 2 — the authoring guide, nav, and redirects

**`docs/data-uploads/define-upload-rules.mdx` → `docs/data-uploads/define-file-rules.mdx`**, `title: "Define file rules"`. Changes beyond vocabulary:

- The **Prerequisites** section currently says a connected bucket is required, full stop. That is now only true for rules with an `upload` block. Reword: a rule that uploads needs a connected bucket; a retention-only rule does not.
- Replace the example YAML with the Context YAML (`name`, `source`, `upload`, `retention`). Show a second, short example of a retention-only rule so the optionality is concrete.
- Keep "identified by hashing their content, so rules with equivalent content are deduplicated"; add that `name` is part of that content.
- Add one sentence: rule `name`s must be unique across the rule files included in a single `miru release create`.
- Flags section: `--upload-rule`/`--upload-rules` → `--file-rule`/`--file-rules`; update the example invocation and the `./upload-rules/` example directory path → `./file-rules/`.
- Swap `<Sources />`/`<Destinations />` for `<Sources />`/`<Upload />`/`<Retention />` under `### Source`, `### Upload`, `### Retention`.

**`docs/docs.json`**, two nav strings inside the `"Data Uploads"` group:
`"data-uploads/primitives/upload-rules"` → `"data-uploads/primitives/file-rules"`, and `"data-uploads/define-upload-rules"` → `"data-uploads/define-file-rules"`. Position and group names unchanged.

**`docs/docs.json` → `redirects`**, appended (order in the array does not matter; match the existing object style exactly):

```json
{ "source": "/data-uploads/primitives/upload-rules", "destination": "/data-uploads/primitives/file-rules" },
{ "source": "/data-uploads/define-upload-rules", "destination": "/data-uploads/define-file-rules" }
```

These must land in the **same commit** as the file renames (D5), or the `redirects` lint rule fails in both directions.

### Milestone 3 — inbound repoints

Every link target `/data-uploads/primitives/upload-rules` → `/data-uploads/primitives/file-rules`, `/data-uploads/define-upload-rules` → `/data-uploads/define-file-rules`, and **`#destinations` → `#uploads`**. Per file, beyond the mechanical swap:

- `docs/data-uploads/overview.mdx` — `:7` "upload rules" → "file rules" and reword the opening so it does not promise that every rule uploads; `:31` link + anchor; `:41` link. The page's overall framing (data uploads move data off devices) stays — this page is about the uploads half.
- `docs/data-uploads/primitives/uploads.mdx` — `:10` "Once a file rule with an `upload` block is deployed…"; `:65` `<ParamField path="upload rule" type="Upload Rule">` → `<ParamField path="file rule" type="File Rule">` (**this changes the rendered anchor `#param-upload-rule` → `#param-file-rule`; grep for inbound links to it before and after**); `:90` `#destinations` → `#uploads`.
- `docs/data-uploads/primitives/buckets.mdx:68` — "no upload rules associated" → "no file rules associated".
- `docs/data-uploads/primitives/upload-collections.mdx` — `:34` "upload rules are annotated with their collection's `slug`" → file rules, via `upload.collection_slug`; `:42` "created automatically when you create an upload rule" → "…when a file rule that uploads names them"; link to `/data-uploads/define-file-rules`.
- `docs/snippets/definitions/upload-collection.mdx` — "a named grouping of [upload rules]" → "[file rules]"; soften "Every upload and upload rule belongs to exactly one collection" to cover only upload-bearing rules.
- `docs/developers/agent/overview.mdx:11` — link + wording.
- `docs/developers/agent/filesys-access.mdx` — `:95` and `:107` vocabulary + links; **`### Deleting after upload`** (`:144`) is rewritten: the trigger is no longer `destination.delete_policy: after_upload` but "a file rule with a `retention` block". Link `/data-uploads/primitives/file-rules#retention`. The permission table and the write-on-parent-directory requirement are unchanged and are the point of the section — keep them. Respect D3: describe *that* the agent deletes and *what permissions that needs*, not *when*.
- `docs/primitives/releases.mdx` — `:33` `<ParamField path="upload rules" type="[] Upload Rule">` → `path="file rules" type="[] File Rule"` (anchor `#param-upload-rules` → `#param-file-rules`; grep for inbound links); `:36` link + `#properties`; **`:38` examples `UPR-FmoDN, UPR-3JcLq, UPR-HHDRn` → `FLR-…`**; `:59` link.
- `docs/cfg-mgmt/create-a-release.mdx` — `:100`, `:104`, `:106` vocabulary only.
- `docs/admin/users/access-control.mdx` — `:103` heading link `[Upload rules](…)` → `[File rules](/data-uploads/primitives/file-rules)`; `:105` "Create an upload rule" → "Create a file rule".

### Milestone 4 — CLI reference and the unreleased changelog entry

Gate: D1 applies to this whole milestone.

- `docs/snippets/references/cli/releases/create/flags.mdx` — `--upload-rules` → `--file-rules`, `--upload-rule` → `--file-rule`; body text "upload rule YAML files" → "file rule YAML files"; **fix the two links, which are currently relative and broken** (`](data-uploads/primitives/upload-rules#file-format)` has no leading `/`) → `](/data-uploads/primitives/file-rules#file-format)`; example paths `./upload-rules/` → `./file-rules/`, `./upload-rules/robot-logs.yaml` → `./file-rules/robot-logs.yaml`.
- `docs/snippets/references/cli/releases/create/usage.mdx` — `<Tab title="Upload Rules">` → `<Tab title="File Rules">` (Title Case is legal inside JSX attributes; `heading-case` only inspects Markdown headings), prose + both flags in the fenced example.
- `docs/snippets/references/cli/releases/create/scopes.mdx` — "If creating a release with upload rules" → "with file rules"; **`upload_rules:manage` → `file_rules:manage`**; `upload_collections:manage` unchanged.
- `docs/references/cli/release-create.mdx` — `:11`, `:13`, `:17`, `:20` vocabulary + the `/data-uploads/define-file-rules` link.
- `docs/changelog/cli.mdx`, the existing `# v0.11.0` entry (**amend, do not add a new version** — this release has not shipped; see D6). Add a `## Breaking changes` section **above** `## Features`, using the `# v0.10.0` entry in the same file as the format precedent (`<Dropdown title="…">` wrapping a ` ```diff ` block):
  - Dropdown "File rule flags": `--upload-rule` → `--file-rule`, `--upload-rules` → `--file-rules`, no aliases.
  - Dropdown "Rule file format": diff from the old shape to the new one (`name` added and required; `collection_slug` moves under `upload`; `destination` → `upload`; `delete_policy` removed; `retention` added; `path` must contain `{upload_id}`).
  - One line: the `upload_rules:*` API key scopes are now `file_rules:*`.
  - One line: rules are keyed by `name`, which must be unique across the files in one release; two rules may now share a collection slug.
  Confirm the date against the real release before `gh pr ready` (the entry currently reads *August 12, 2026*).
- A new `docs/changelog/product.mdx` `<Update>` entry for the rename is **optional and out of required scope** (precedent: `20260810-document-instance-slots.md`). If added, note headings inside `<Update>` blocks use Title Case.

### Milestone 5 — changelog link targets

Hrefs only, prose untouched (D6):
- `docs/changelog/cli.mdx:43` — `](/data-uploads/primitives/upload-rules)` → `](/data-uploads/primitives/file-rules)`; the words "upload rules" and the `--upload-rule` flag names in the v0.10.2 entry **stay** (that is what v0.10.2 shipped).
- `docs/changelog/agent.mdx:11` — href only; "upload rule" wording stays.
- `docs/changelog/product.mdx:166` — href only. `:170`, `:172`, `:175`, `:180`, `:196`, `:198`, `:200`, `:201` keep every word, including the `### Upload rules` heading and the image URL. `:182` `[Define upload rules »](/data-uploads/define-upload-rules)` — href → `/data-uploads/define-file-rules`, **link text unchanged**.

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs` unless stated otherwise. Every `./scripts/lint.sh` run must print `== MDX Prose ==`, `== ESLint (MDX) ==`, `== CSpell ==`, `== OpenAPI ==` and end with exactly `All documentation lint checks passed.` Anything else is a failure — do not commit.

### Milestone 0 — baseline

1. `git status` (clean), `git rev-parse --abbrev-ref HEAD` → `docs/rename-upload-rules-to-file-rules`.
2. `pnpm install --frozen-lockfile`.
3. Baseline: `./scripts/lint.sh; echo "exit=$?"`, `pnpm run validate; echo "exit=$?"`, `./scripts/audit.sh > /tmp/audit-base.txt 2>&1; echo "exit=$?"`. If `audit.sh` is non-zero, prove pre-existence against a pristine-`main` worktree and record both outputs in Surprises & Discoveries — `preflight.sh` uses `set -euo pipefail` and will abort at `=== Audit ===` before the bats stage, so run `pnpm run test:lint` and `bats pub/scripts/agent/check-miru-access_test.bats` by hand in that case.
4. Re-verify the upstream facts (they may have moved since 2026-08-11):
   - `cd ../openapi && git log --oneline -3 -- apis/configs/components/schemas/file-rule.yaml` and re-read that file plus `apis/configs/components/requests/file-rule.yaml`. **The file wins over this plan.** If `ttl_secs`/`require_upload` have changed again, update this plan's Context before writing any prose.
   - `cd ../cli-private && git log --oneline -1 -- internal/domain/filerules/spec.go && git tag --contains 8ba7471 && git tag --list 'v0.11.0'`.
   - `cd ../agent && grep -n "BaseUploadRule\|BaseFileRule" agent/src/models/file_rule.rs` — if the agent now consumes `BaseFileRule` and enforces general TTL, revisit D3 before writing the retention `<Note>`.
5. Capture the inventory baseline:
   `grep -rn "upload rule\|Upload rule\|Upload Rule\|upload-rule\|upload_rule\|uploadRule" docs/ | wc -l`

### Milestone 1

1. `git mv docs/snippets/definitions/upload-rule.mdx docs/snippets/definitions/file-rule.mdx`
2. `mkdir docs/snippets/file-rules && git mv docs/snippets/upload-rules/sources.mdx docs/snippets/file-rules/sources.mdx && git mv docs/snippets/upload-rules/destinations.mdx docs/snippets/file-rules/upload.mdx && rmdir docs/snippets/upload-rules`
3. Create `docs/snippets/file-rules/retention.mdx`; edit the three snippets per Plan of Work.
4. `git mv docs/data-uploads/primitives/upload-rules.mdx docs/data-uploads/primitives/file-rules.mdx`; rewrite it; re-sort the import block.
5. `./scripts/lint.sh` (expect `importresolves` failures until Milestone 3's importers are fixed if any remain — if so, fold steps 1–3 of Milestone 3 into this commit rather than committing red).
6. Commit: `docs(data-uploads): rename the upload rule primitive to file rules`

### Milestone 2

1. `git mv docs/data-uploads/define-upload-rules.mdx docs/data-uploads/define-file-rules.mdx`; rewrite.
2. Edit `docs/docs.json`: two nav strings + two redirect objects.
3. `./scripts/lint.sh` (the `redirects` rule now proves both directions) and `pnpm run validate`.
4. Commit: `docs(data-uploads): rename the file rule authoring guide and add redirects`

### Milestone 3

1. Apply the per-file edits listed in Plan of Work.
2. `grep -rn "primitives/upload-rules\|define-upload-rules\|#destinations" docs/ --include='*.mdx'` → only the changelog hits from Milestone 5 remain.
3. `grep -rn "param-upload-rule\|param-upload-rules" docs/` → nothing.
4. `./scripts/lint.sh` and `pnpm run validate`.
5. Commit: `docs: repoint upload rule references to file rules`

### Milestone 4

1. Apply the CLI reference and `# v0.11.0` changelog edits.
2. `grep -rn '\-\-upload-rule\|upload_rules:' docs/` → nothing.
3. `./scripts/lint.sh`.
4. Commit: `docs(references): document the file rule CLI flags and scopes`

### Milestone 5

1. Repoint the changelog hrefs.
2. Verify prose is untouched: `git diff main...HEAD -- docs/changelog/agent.mdx docs/changelog/product.mdx` shows **only** changed URLs inside `](…)` (and, in `cli.mdx`, only the URL plus the new v0.11.0 Breaking changes block).
3. `./scripts/lint.sh` and `pnpm run validate`.
4. Commit: `docs(changelog): repoint file rule links without rewriting history`

### Milestone 6 — validate and ship

1. `./scripts/preflight.sh` (or the by-hand stage list if `audit.sh` is pre-existing red).
2. `git push -u origin docs/rename-upload-rules-to-file-rules`; open a **draft** PR.
3. `gh pr checks --watch` against the pushed head SHA.
4. Resolve the D1 gate and the D3 Open Question, then `gh pr ready`.

## Validation and Acceptance

**`$preflight` must report `CLEAN` before the PR leaves draft and before this task is reported complete.** CLEAN means, on the head SHA you actually pushed: the `changes`, `lint`, and `shell-tests` CI jobs are green, **and** the `audit` job is either green or red with an advisory set identical to pristine `main` re-measured at the same moment as the head you are shipping. A green local run is not sufficient; a green run on an older commit is not sufficient. Confirm job states with `gh pr checks --watch` against the pushed head SHA. This is the only definition of CLEAN in this plan.

Note that the `lint` CI job now includes `pnpm run validate` (added in `1566a65`), so a broken internal link or a bad nav path fails CI, not just review.

Checks (all run from the repo root; every one must pass):

1. `./scripts/lint.sh` ends with `All documentation lint checks passed.`
2. `pnpm run validate` exits 0.
3. `pnpm run test:lint` passes; `bats pub/scripts/agent/check-miru-access_test.bats` passes.
4. **No stale concept outside the sanctioned places:**
   `grep -rn "upload rule\|Upload rule\|Upload Rule\|upload-rule\|upload_rule\|uploadRule" docs/ --include='*.mdx' --include='*.json'`
   returns hits **only** in `docs/changelog/{cli,agent,product}.mdx` historical prose (D6) and only as visible words — never inside a `](…)` target. Verify the latter with
   `grep -rn "](.*upload-rule" docs/` → nothing.
5. **No stale files/dirs:** `test ! -e docs/data-uploads/primitives/upload-rules.mdx && test ! -e docs/data-uploads/define-upload-rules.mdx && test ! -d docs/snippets/upload-rules && test ! -e docs/snippets/definitions/upload-rule.mdx && echo ok` → `ok`.
6. **Removed vocabulary:** `grep -rn "delete_policy\|after_upload\b\|delete_delay_secs\|max_age_secs\|until_stable\|policy: forever" docs/ --include='*.mdx'` returns nothing outside `docs/changelog/` historical prose.
7. **Nav:**
   `python3 -c "import json;s=json.dumps(json.load(open('docs/docs.json'))['navigation']);assert 'data-uploads/primitives/file-rules' in s and 'data-uploads/define-file-rules' in s and 'upload-rules' not in s and 'define-upload-rules' not in s;print('nav ok')"`
8. **Redirects:**
   `python3 -c "import json;r=json.load(open('docs/docs.json'))['redirects'];m={x['source']:x['destination'] for x in r};assert m['/data-uploads/primitives/upload-rules']=='/data-uploads/primitives/file-rules';assert m['/data-uploads/define-upload-rules']=='/data-uploads/define-file-rules';print('redirects ok')"`
   plus the `redirects` lint rule passing inside check 1 (it independently proves the sources 404 and the destinations resolve).
9. **Anchor contract:** `grep -n '^## ' docs/data-uploads/primitives/file-rules.mdx` shows, in order: `Properties`, `File format`, `Sources`, `Uploads`, `Retention`, `Immutability`, `Git provenance`, `Create a file rule`, `View a file rule`. Every cross-page link into this page targets one of these slugs.
10. **Untouched trees:** `git diff --stat main...HEAD -- tools/lint plans/completed docs/references/platform-api docs/references/device-api` prints nothing.
11. **Scopes and flags:** `grep -rn "file_rules:manage" docs/snippets/references/cli/releases/create/scopes.mdx` → 1 hit; `grep -rn -- "--file-rule\b" docs/` → hits in `flags.mdx`, `usage.mdx`, `define-file-rules.mdx`.
12. **ID prefix:** `grep -rn "UPR-" docs/` → nothing.
13. Behavior check under `pnpm dev` (skip with a note if the environment is non-interactive): `/data-uploads/primitives/file-rules` and `/data-uploads/define-file-rules` render with the sidebar entries in their old positions; `/data-uploads/primitives/upload-rules` and `/data-uploads/define-upload-rules` redirect rather than 404; the deep links from `/data-uploads/overview`, `/data-uploads/primitives/uploads`, and `/developers/agent/filesys-access` land on the right sections.
14. CLEAN as defined above holds on the pushed head SHA.

Content acceptance, confirmed by reading the diff alone:

15. `retention` absent is documented as "the device keeps matching files indefinitely — Miru never deletes them", not as an unspecified default.
16. `ttl_secs` is documented as required **inside** the `retention` block, with `0` meaning delete as soon as eligible.
17. `require_upload` is documented as present exactly when the rule has an `upload` block, including the 9-attempt abandon consequence.
18. Digest dedup is scoped to the **workspace**, not the collection.
19. `name` is documented as part of the digest, so renaming mints a new rule.
20. No claim is made about device-side enforcement of `ttl_secs > 0` or of retention-only rules beyond what D3 permits.
21. No page states or implies that every file rule uploads, or that every file rule belongs to an upload collection.
22. No historical changelog sentence has changed — only URLs (and the unreleased `# v0.11.0` entry).

## Idempotence and Recovery

Every step is a text edit, a `git mv`, or a read-only command, so the plan is safe to re-run. `pnpm install --frozen-lockfile`, `./scripts/lint.sh`, `pnpm run validate`, `./scripts/audit.sh`, `pnpm run test:lint`, and `bats …` are read-only and repeatable. `pnpm dev` serves on port 3000 and writes nothing to the repo.

Each milestone is one commit, so `git revert <sha>` rolls a milestone back independently — **except** that Milestones 1 and 2 are coupled by the `redirects` and `importresolves` rules: reverting Milestone 2 alone leaves redirect destinations pointing at a page that exists but a nav entry that does not, so revert them together or not at all. To restart: `git checkout main && git branch -D docs/rename-upload-rules-to-file-rules`. Nothing outside this repository is modified, so there is no external state to clean up.

If `heading-case` rejects a heading after a rename, change the heading text — but then re-run the anchor greps in Validation checks 4, 9, and 13, because every cross-page deep link is derived from the final heading text.

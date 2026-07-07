# Refactor data-uploads docs to the current feature state, mirroring cfg-mgmt

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (this repo, branch `feat/data-recording`) | read-write | Rewrite and restructure `docs/data-uploads/`, update `docs/docs.json` navigation and `docs/snippets/definitions/` |
| `openapi/` (sibling checkout `../openapi`, branch `main`) | read-only | Source of truth for upload / upload rule / upload collection / bucket schemas and endpoints |
| `backend/` (sibling checkout `../backend`, branch `main`) | read-only | Destination-path validation rules, bucket verification flow, dedup semantics |
| `cli-private/` (sibling checkout `../cli-private`, branch `main`) | read-only | `miru release create` upload-rule flags and the upload-rule YAML file spec |
| `agent/` (sibling checkout `../agent`) | read-only | Device-side behavior context (scanner, digest, delivery of rules via releases) |

This plan lives in `docs/plans/` because every change is a docs-repo change. Sibling repos are referenced from the workbench layout (`/home/ben/miru/workbench5/repos/<repo>`); adjust the prefix if your checkout differs.

## Purpose / Big Picture

The `docs/data-uploads/` section was drafted in June 2026 and has drifted from the shipped API: object properties are wrong (the user called out **properties** and **integration** pages specifically), the upload flow changed from presigned URLs to token-only downscoped credentials, a whole parent object (upload collections) is missing, and the AWS bucket setup documents the wrong order of operations. After this plan, a reader of docs.mirurobotics.com can set up a bucket, write upload rule files, release them, and audit uploads using instructions that match the live API — and the section's navigation has the same object/operations/audit shape as Config Management (`docs/cfg-mgmt/`).

## Progress

- [ ] Milestone 1: Restructure section and navigation (delete orphans, move audit page, scaffold new pages, update docs.json)
- [ ] Milestone 2: Upload collections page and definition snippets
- [ ] Milestone 3: Upload rules group rewrite (overview, sources, destinations, manage)
- [ ] Milestone 4: Uploads page and section overview rewrite (token-only flow, current properties)
- [ ] Milestone 5: Buckets group rewrite (overview, AWS, GCS)
- [ ] Milestone 6: Releasing upload rules rewrite (CLI YAML flow)
- [ ] Milestone 7: Audit page content and final validation (preflight clean)

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

- Decision: Keep `data-uploads/overview.mdx` as the first nav entry even though cfg-mgmt has no overview page.
  Rationale: data uploads is a multi-object pipeline (rule → credentials → object in customer bucket); one orientation page earns its keep. Everything else mirrors cfg-mgmt.
  Date/Author: 2026-07-07, planning.
- Decision: Delete `docs/data-uploads/quick-start/` instead of rewriting it.
  Rationale: its four pages are orphaned (absent from docs.json) duplicates of `docs/getting-started/quick-start/`, left behind by the branch restructure; cfg-mgmt has no quick-start either.
  Date/Author: 2026-07-07, planning.
- Decision: No redirects needed for renames/moves inside `data-uploads/`.
  Rationale: `git ls-tree -r main --name-only | grep data-uploads` returns nothing — the section exists only on this unmerged branch, so no published URLs can break.
  Date/Author: 2026-07-07, planning.
- Decision: Document dashboard operations with the existing screenshot-placeholder convention.
  Rationale: the frontend has no uploads UI on its main branch yet; the branch already uses screenshot TODOs (commit f541569). API-derived facts, not UI labels, are the contract.
  Date/Author: 2026-07-07, planning.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

### The feature, in current terms

Miru's data recording feature uploads files from devices to customer-owned cloud buckets. Four objects, all defined in the sibling `openapi` repo under `apis/configs/components/schemas/`:

- **Upload collection** (`upload_collection`, `upl_col_…`): named container that upload rules belong to — the analog of a config type (`name` + immutable `slug`). New since the docs were drafted.
- **Upload rule** (`upload_rule`, `upl_rule_…`): immutable rule matching device files (source) to a bucket path (destination). Deduplicated by digest **within a collection**. Shipped to devices as part of releases; linked to git commits for provenance.
- **Bucket** (`bucket`, `bkt_…`): a registered customer cloud bucket (GCS or AWS S3) with keyless credentials and a verification status.
- **Upload** (`upload`, `upl_…`): ledger entry for one file upload.

The device flow is **token-only** (openapi commit 2e23d0a): the agent calls `POST /uploads` on the device API with file details; the backend creates a `pending` ledger entry, dedups by digest, and vends short-lived **downscoped cloud credentials** scoped to one object key (S3: STS session credentials used with the native SDK's multipart upload; GCS: a downscoped OAuth2 access token via Credential Access Boundary used with resumable upload). The device uploads via the native cloud SDK, then calls `POST /uploads/{upload_id}/confirm` (status → `uploaded`, `uploaded_at` set). `POST /uploads/{upload_id}/credentials` re-vends credentials mid-upload. There are **no presigned URLs** in the current flow — the old docs' "signed URL + HTTPS PUT" story is stale. (Note: the `UploadStatus` enum description in the spec still says "presigned URL"; do not copy that stale wording.)

Platform API endpoints (see `../openapi/apis/configs/paths/`): uploads — list/get with filters (device, rule, collection, status); upload_collections — create/list/get/update(name)/delete/archive/unarchive; upload_rules — create/list/get (no update/delete; immutable); buckets — create/list/get/delete/verify/archive/unarchive (no update; immutable).

### Source-of-truth files (read these while writing each page)

- `../openapi/apis/configs/components/schemas/upload.yaml`, `upload-rule.yaml`, `upload-collection.yaml`, `bucket.yaml`
- `../openapi/apis/configs/components/enums/upload-status.yaml`, `upload-delete-policy.yaml`, `bucket-provider.yaml`, `bucket-status.yaml`
- `../openapi/apis/configs/components/requests/upload-collection.yaml`, `upload-rule.yaml`, `bucket.yaml`
- `../openapi/apis/apps/backend-server/agent/paths/uploads.yaml` and `.../agent/components/schemas/upload.yaml` (device flow, credentials, metadata map)
- `../backend/internal/configs/domain/uploadrules/destpath.go` (path template variables and validation)
- `../backend/internal/configs/services/buckets/verify.go` and `create.go` (verification flow; server-issued AWS external ID)
- `../cli-private/internal/domain/uploadrules/spec.go` (upload-rule YAML file spec) and `../cli-private/internal/commands/release/flags.go` (`--upload-rule`, `--upload-rules` flags)

### Current object shapes (condensed from the spec; the pages must match these)

Upload (`BaseUpload`): `id`, `device_id`, `upload_rule_id`, `upload_collection_name`, `source{file_path, file_modified_at}`, `destination{bucket_id, object_key}`, `digest`, `size`, `status` (`pending` | `uploaded`), `incomplete` (bool, independent of status), `release_id`, `deployment_id`, `uploaded_at` (nullable), `created_at`, `updated_at`, `workspace_id`. Device-API responses add `metadata`: a generic string map stamped as cloud object metadata (illustrative keys: `device_id`, `release_id`, `release_version`, `deployment_id`, `digest`, `file_modified_at`; size limits S3 ~2 KB / GCS ~8 KB).

Upload rule (`BaseUploadRule`): `id`, `upload_collection_id`, `upload_collection_name`, `digest`, `source{glob, stability_window_secs}`, `destination{bucket_id, bucket_name, path, delete_policy}`, `created_at`, `updated_at`; expandable `upload_collection` and `upload_rule_git_commits` (each commit link: `filepath`, `git_commit_id`, `created_by_id`).

Upload collection: `id`, `name`, `slug` (immutable), `created_by_id`, `updated_by_id`, `archived_at`, `actions{update, delete, archive, unarchive}`, timestamps.

Bucket: `id`, `name` (immutable, unique per workspace — this IS the cloud bucket's name), `provider` (`gcs` | `aws`), `config` (discriminated by provider), `status` (`pending` | `verified` | `invalid`), `last_verified_at`, `verification_error`, `created_by_id`, `updated_by_id`, `archived_at`, `actions{delete, verify, archive, unarchive}`, timestamps. `GcsBucketConfig`: `project_id`, `wip_provider`, `service_account_email`. `AwsBucketConfig`: `region`, `role_arn`, `external_id` — and `external_id` is **issued by Miru at create** (`CreateAwsBucketConfig` omits it).

Destination path (`destpath.go`): optional at create; server default `{device_id}/{year}-{month}-{day}/{upload_id}/{file_name}`. Supported variables: `{device_id}`, `{device_name}`, `{file_name}`, `{upload_id}`, `{year}`, `{month}`, `{day}`, `{hour}`, `{minute}`. `{upload_id}` is required by backend validation and `{device_id}` by the spec pattern — document both as required. Max 1024 bytes; no empty or `..` segments; no control characters.

### Stale claims found (docs page → what is wrong)

`docs/data-uploads/uploads.mdx`:
1. Documents status enum `pending`/`uploading`/`uploaded` and an `uploading_at` property — the `uploading` state and `uploading_at` do not exist; the enum is `pending` | `uploaded`.
2. Documents `device`, `upload rule`, `bucket`, `release` as expanded object properties — the API returns plain `device_id`, `upload_rule_id`, `release_id`; the bucket appears only as `destination.bucket_id`. Upload has no expandable fields.
3. Missing properties: `id`, `upload_collection_name`, `deployment_id`, `source{…}`, `destination{…}`, `created_at`, `updated_at`, `workspace_id`, and the device-API `metadata` map.
4. Top-level `object_key` — actually `destination.object_key`.
5. Flow described as "short-lived signed URL … uploads via HTTPS" — actual flow is token-only downscoped credentials + native SDK (multipart/resumable) + explicit confirm.
6. Provenance described as fixed stamped fields — now a generic string metadata map.

`docs/data-uploads/upload-rules/{overview,sources,destinations}.mdx`:
7. `source.poll_interval` documented — removed entirely (backend migration `20260701120000_drop_upload_rule_poll_interval.sql`).
8. `stability_window: 60s` (duration string) — now `stability_window_secs` (integer seconds).
9. No mention of upload collections; dedup described as global — rules belong to a collection and dedup by digest is scoped within it.
10. Destination missing `bucket_name`; `path` documented as required with only `{device_id}` mandatory — actually optional at create with a server default, nine template variables, and `{upload_id}` also required.
11. No mention of git provenance (`upload_rule_git_commits`).

`docs/data-uploads/buckets/{overview,aws,gcs}.mdx` (the "integration" docs):
12. AWS and GCS configs both document a `bucket_name` config field — neither config has one; the top-level bucket `name` is the cloud bucket name.
13. AWS `external_id` documented as user-supplied before registration — it is server-issued at create (backend commit d1eea9fc8). The setup order must become: create the bucket in Miru (name, region, role_arn) → Miru returns `external_id` → add it to the IAM role trust policy as the `sts:ExternalId` condition → run Verify.
14. Bucket properties list only `provider` + `config` — missing `name`, `status`/`last_verified_at`/`verification_error`, audit fields, and `actions` (verify/archive/unarchive/delete; no edit — buckets are immutable).
15. Credential-minting narrative says Miru mints presigned SigV4 / V4-signed URLs — devices now receive downscoped credentials/tokens and upload via native SDKs; verify wording against `../backend/internal/configs/services/buckets/verify.go` and the credential-vending code (also confirm the exact S3 IAM actions needed for multipart, e.g. whether `s3:AbortMultipartUpload` is required, before editing the policy JSON in `aws.mdx`).
16. Verify described loosely — verify now returns the full Bucket with updated `status`/`verification_error` (spec commit c412a7f).

`docs/data-uploads/releasing-upload-rules.mdx`:
17. Says rules are "included in a release by ID" — actually `miru release create` reads upload-rule **YAML files** from the repo (`--upload-rule <file>` repeatable, `--upload-rules <dir>` repeatable), finds-or-creates rules by digest, resolves collections by `collection_slug`, and records git commit provenance. File spec (from `spec.go`): `collection_slug`; `source.glob`, `source.stability_window_secs`; `destination.bucket` (bucket name), optional `destination.path`, optional `destination.delete_policy`.

Structure gaps:
18. `data-uploads/overview.mdx` exists but is orphaned — not in `docs/docs.json`.
19. `data-uploads/quick-start/` (4 pages) orphaned duplicates of `getting-started/quick-start/`.
20. `audit-uploads.mdx` is an empty stub, flat rather than in an audit group.
21. No upload collections page (cfg-mgmt's analog, `config-types.mdx`, leads its group).
22. No "manage"/viewing page for rules (cfg-mgmt analog: `schemas/manage.mdx`).

### Docs-repo conventions that apply

Mintlify site under `docs/`; nav in `docs/docs.json`. Pages are MDX with `title` frontmatter; properties documented with `<ParamField>` plus `<ImmutableBadge />`/`<MutableBadge />` (import from `/snippets/field-badges.jsx`); definition snippets in `docs/snippets/definitions/`; steps with `<Steps>`; screenshots in `<Frame>` with placeholder TODOs where UI is unbuilt. Headings are sentence case, enforced by the Go linter (`tools/lint/`; allowlist already includes AWS/GCS/WIF/STS/MCAP). Internal links are absolute paths without `.mdx`. Copy house style from the cfg-mgmt pages being mirrored.

## Plan of Work

Restructure the "Data Uploads" docs.json group to mirror the "Config Management" group's shape — container object page first, then the rule/definition group, then the produced-artifact page, the releasing page, the operations group, and the audit group:

| cfg-mgmt shape | data-uploads target |
|---|---|
| — (no overview) | `data-uploads/overview` (kept; see Decision Log) |
| `config-types` | `data-uploads/upload-collections` (new) |
| Config schemas group (`overview`, languages, `manage`) | Upload rules group (`overview`, `sources`, `destinations`, `manage` (new)) |
| `config-instances` | `data-uploads/uploads` |
| `releasing-config-schemas` | `data-uploads/releasing-upload-rules` |
| Deploy configs group | Buckets group (`overview`, `gcs`, `aws`) |
| Audit configs group | Audit uploads group (`audit/upload-history`) |

File-by-file:

| File (under `docs/`) | Action |
|---|---|
| `data-uploads/overview.mdx` | Rewrite flow narrative to token-only credentials; add to nav |
| `data-uploads/upload-collections.mdx` | **Add**: properties (`name`, immutable `slug`), operations (create, rename, archive/unarchive, delete), relationship to rules; modeled on `cfg-mgmt/config-types.mdx` |
| `data-uploads/upload-rules/overview.mdx` | Rewrite: collection parentage, digest dedup within collection, immutability, YAML example with `stability_window_secs`, git provenance |
| `data-uploads/upload-rules/sources.mdx` | Rewrite: `glob`, `stability_window_secs`; delete `poll_interval` |
| `data-uploads/upload-rules/destinations.mdx` | Rewrite: `bucket_id`/`bucket_name`, optional `path` + default + full variable table + `{device_id}`/`{upload_id}` required, `delete_policy` |
| `data-uploads/upload-rules/manage.mdx` | **Add**: viewing rules (created only via releases; no edit/delete; git commit metadata); modeled on `cfg-mgmt/schemas/manage.mdx` |
| `data-uploads/uploads.mdx` | Rewrite properties to `BaseUpload` + `metadata` map; two-state lifecycle; token flow + confirm; keep guarantees section but re-verify each claim |
| `data-uploads/releasing-upload-rules.mdx` | Rewrite around CLI YAML file flow; modeled on `cfg-mgmt/releasing-config-schemas.mdx` |
| `data-uploads/buckets/overview.mdx` | Rewrite: add `name`, verification status fields, actions; register → verify flow |
| `data-uploads/buckets/aws.mdx` | Rewrite: config = `region` + `role_arn` (+ Miru-issued `external_id`); reorder setup steps (register first); drop `bucket_name` field |
| `data-uploads/buckets/gcs.mdx` | Rewrite: config = `project_id` + `wip_provider` + `service_account_email`; drop `bucket_name` field; update credential story |
| `data-uploads/audit-uploads.mdx` | **Move** to `data-uploads/audit/upload-history.mdx` and write content: filter uploads (device/rule/collection/status), `incomplete` flag, object metadata provenance; modeled on `cfg-mgmt/audit/device-history.mdx` |
| `data-uploads/quick-start/` (4 files) | **Delete** (orphaned duplicates) |
| `snippets/definitions/upload.mdx`, `upload-rule.mdx`, `bucket.mdx` | Update wording to current model |
| `snippets/definitions/upload-collection.mdx` | **Add** (pattern: `snippets/definitions/config-type.mdx`) |
| `docs.json` | Replace the "Data Uploads" group with the block below |

Target `docs.json` "Data Uploads" group:

    {
      "group": "Data Uploads",
      "pages": [
        "data-uploads/overview",
        "data-uploads/upload-collections",
        {
          "group": "Upload rules",
          "pages": [
            "data-uploads/upload-rules/overview",
            "data-uploads/upload-rules/sources",
            "data-uploads/upload-rules/destinations",
            "data-uploads/upload-rules/manage"
          ]
        },
        "data-uploads/uploads",
        "data-uploads/releasing-upload-rules",
        {
          "group": "Buckets",
          "pages": [
            "data-uploads/buckets/overview",
            "data-uploads/buckets/gcs",
            "data-uploads/buckets/aws"
          ]
        },
        {
          "group": "Audit uploads",
          "pages": [
            "data-uploads/audit/upload-history"
          ]
        }
      ]
    }

Inbound links that must keep resolving (paths unchanged, verify anyway): `docs/developers/agent/overview.mdx` → `/data-uploads/upload-rules/overview`; `snippets/definitions/upload.mdx` and `upload-rule.mdx` → `/data-uploads/buckets/overview`. Add any new vocabulary (`downscoped`, `unarchive`, …) to `cspell.json` only when the spell check flags it.

## Concrete Steps

All commands run from the docs repo root, `/home/ben/miru/workbench5/repos/docs`, on branch `feat/data-recording`. After each milestone run `./scripts/lint.sh` (expect exit 0) and commit — one commit per milestone.

Milestone 1 — restructure and navigation:

    git rm -r docs/data-uploads/quick-start
    mkdir -p docs/data-uploads/audit
    git mv docs/data-uploads/audit-uploads.mdx docs/data-uploads/audit/upload-history.mdx
    # create docs/data-uploads/upload-collections.mdx and
    # docs/data-uploads/upload-rules/manage.mdx as frontmatter-only stubs
    # edit docs/docs.json: replace the "Data Uploads" group with the target block
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): restructure section to mirror cfg-mgmt"

Milestone 2 — upload collections:

    # write docs/data-uploads/upload-collections.mdx (model: docs/cfg-mgmt/config-types.mdx)
    # add docs/snippets/definitions/upload-collection.mdx; update upload.mdx, upload-rule.mdx,
    # bucket.mdx definition snippets to current wording
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): document upload collections"

Milestone 3 — upload rules group:

    # rewrite upload-rules/overview.mdx, sources.mdx, destinations.mdx; write manage.mdx
    # verify every property against ../openapi/apis/configs/components/schemas/upload-rule.yaml
    # and ../backend/internal/configs/domain/uploadrules/destpath.go
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): re-sync upload rule pages to current schema"

Milestone 4 — uploads and overview:

    # rewrite docs/data-uploads/uploads.mdx and docs/data-uploads/overview.mdx
    # verify against ../openapi/apis/configs/components/schemas/upload.yaml and
    # ../openapi/apis/apps/backend-server/agent/{paths,components/schemas}/upload*.yaml
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): document token-only upload flow and current upload object"

Milestone 5 — buckets (integration setup):

    # rewrite buckets/overview.mdx, aws.mdx (reordered steps, Miru-issued external ID), gcs.mdx
    # verify against ../openapi/.../schemas/bucket.yaml and
    # ../backend/internal/configs/services/buckets/{create,verify}.go
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): fix bucket properties and integration setup flows"

Milestone 6 — releasing upload rules:

    # rewrite releasing-upload-rules.mdx around the CLI YAML file flow
    # verify against ../cli-private/internal/domain/uploadrules/spec.go and
    # ../cli-private/internal/commands/release/flags.go
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): document releasing upload rules via the CLI"

Milestone 7 — audit page and final validation:

    # write docs/data-uploads/audit/upload-history.mdx content
    ./scripts/lint.sh
    npm run test:lint
    ./scripts/preflight.sh
    git add -A && git commit -m "docs(data-uploads): add upload history audit page"

## Validation and Acceptance

Run from `/home/ben/miru/workbench5/repos/docs`:

1. `./scripts/lint.sh` exits 0 (prose lint, ESLint MDX with `--max-warnings=0`, cspell, OpenAPI check all pass).
2. `npm run test:lint` exits 0.
3. `./scripts/preflight.sh` exits 0 — **preflight must report clean before these changes are published or a PR is opened.**
4. Navigation completeness — expect no output:

       cd docs && for f in $(find data-uploads -name '*.mdx'); do p="${f%.mdx}"; grep -q "\"$p\"" docs.json || echo "ORPHAN: $p"; done

5. Internal links from the section resolve — expect no output (also covers `href=` card links):

       cd docs && grep -rhoE '(\]\(|href=")/[a-z0-9/-]+' --include='*.mdx' data-uploads | sed -E 's/^(\]\(|href=")//' | sort -u | while read -r l; do [ -f ".${l}.mdx" ] || echo "CHECK: $l"; done

   Anything printed must be a known non-page path (anchors, API-reference URLs) or fixed.
6. Content acceptance (each returns no matches):

       grep -rn "poll_interval" docs/data-uploads
       grep -rn "uploading_at" docs/data-uploads
       grep -rn "stability_window[^_]" docs/data-uploads

7. Behavioral spot checks: `docs/data-uploads/uploads.mdx` documents exactly two statuses (`pending`, `uploaded`) and the `source`/`destination`/`metadata` shapes; `buckets/aws.mdx` registers the bucket in Miru *before* the trust-policy step and calls `external_id` Miru-issued; `upload-collections.mdx` documents `name` + immutable `slug`; the rendered nav (`cd docs && mint dev`, open http://localhost:3000) shows the seven-entry Data Uploads group matching the target block.

## Idempotence and Recovery

Every step is a git-tracked edit on branch `feat/data-recording`; nothing touches generated files or other repos. Re-running lint/preflight is always safe. If a milestone goes wrong before its commit, `git checkout -- <path>` (or `git reset --hard HEAD` for the whole milestone) returns to the last good state; after a commit, `git revert <sha>` undoes exactly one milestone. The `git mv`/`git rm` steps in Milestone 1 fail loudly if re-run (source already moved/deleted) — skip them on retry. The docs.json group replacement is idempotent: the target block fully replaces the old group, so re-applying it yields the same file.

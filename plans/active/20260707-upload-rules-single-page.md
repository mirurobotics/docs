# Collapse upload rules to one page, add schema provenance sections, add cloud-CLI bucket tabs

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (this repo, branch `feat/data-recording`) | read-write | Merge `docs/data-uploads/upload-rules/*` into one page, update `docs/docs.json`, repoint inbound links, extend `docs/cfg-mgmt/schemas/overview.mdx`, add CLI tabs to `docs/data-uploads/buckets/integrate/{gcs,aws}.mdx` |
| `openapi/` (sibling checkout `../openapi`, branch `main`) | read-only | Source of truth for config-schema / upload-rule git-commit shapes |
| `backend/` (sibling checkout `../backend`, branch `main`) | read-only | Glob validation rules, schema dedup and git-commit-link accrual, schema delete semantics |
| `infra/` (sibling checkout `../infra`, branch `main`) | read-only | Live GCS/AWS bucket fixtures the CLI-tab commands are cross-checked against |

This plan lives in `docs/plans/` because every change is a docs-repo change. Sibling repos are referenced from the workbench layout (`/home/ben/miru/workbench5/repos/<repo>`); adjust the prefix if your checkout differs.

## Purpose / Big Picture

Three reader-facing improvements to the data-recording docs (all on the unmerged `feat/data-recording` branch, base `main` at `35e7cb6`):

1. The four upload-rules pages (`overview`, `sources`, `destinations`, `manage`) become **one page** at `/data-uploads/upload-rules` with a fixed section outline, so a reader learns the whole object without hopping a four-page nav group.
2. The config-schemas overview gains the same **Immutability and releases** and **Git provenance** sections the upload-rules page has, so the two release-shipped, content-identified objects read alike.
3. The GCS and AWS bucket integrate pages gain a third **CLI** tab (between Console and Terraform) with pure `gcloud` / `aws` command sequences, so an infra engineer can connect a bucket without touching a console or writing Terraform.

Observable outcome: on `mint dev` (or the built site), `/data-uploads/upload-rules` renders all seven sections; the old four URLs are gone from the nav and no page links to them; `/cfg-mgmt/schemas/overview` ends with the two new sections; both integrate pages show three tabs whose CLI steps mirror the Terraform tab resource-for-resource.

## Progress

- [x] Milestone 1: Merged upload-rules page, nav update, repo-wide link repoints, delete `upload-rules/` dir
- [x] Milestone 2: Config schemas "Immutability and releases" and "Git provenance" sections
- [x] Milestone 3: CLI tabs on `gcs.mdx` and `aws.mdx`
- [ ] Final validation: `./scripts/preflight.sh` clean

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

- Decision: Merged page lives at `docs/data-uploads/upload-rules.mdx` (flat file); the `docs/data-uploads/upload-rules/` directory is deleted.
  Rationale: URL `/data-uploads/upload-rules` is the natural group-name URL; Mintlify resolves a flat `.mdx` fine once the directory is gone; keeps the nav entry a plain string.
  Date/Author: 2026-07-07, planning.
- Decision: No redirects for the four removed URLs.
  Rationale: `data-uploads/` does not exist on `main` (verified in the predecessor plan `plans/completed/20260707-data-recording-docs-refactor.md`); no published URL can break.
  Date/Author: 2026-07-07, planning.
- Decision: "File formats" on the merged page = the YAML **rule-file** format (one rule per YAML file), showing the example rule file and linking to `/data-uploads/defining-releases#upload-rule-files` for the field-by-field reference.
  Rationale: mirrors the config-schemas overview, whose "File formats" section describes the on-disk file format of the object's definition; the full rule-file field reference already lives on defining-releases and must not be duplicated.
  Date/Author: 2026-07-07, planning.
- Decision: Schemas' new sections go at the end of `docs/cfg-mgmt/schemas/overview.mdx` (after "Empty schemas"); `manage.mdx` is untouched.
  Rationale: overview owns object semantics, manage owns dashboard viewing (its Metadata tab already shows commit history); end-of-page matches the merged upload-rules page where these sections follow the content sections.
  Date/Author: 2026-07-07, planning.
- Decision: The schemas "Immutability and releases" section says schemas are never edited in place but does **not** claim they cannot be deleted.
  Rationale: unlike upload rules (no delete endpoint), the dashboard API exposes `DELETE /config_schemas/{id}` (`backend/internal/servers/frontend/endpoints.go:117`), allowed only when no config instances reference the schema (`backend/internal/configs/services/config_schemas/delete.go`). Deletion is undocumented today; stay silent rather than assert falsely.
  Date/Author: 2026-07-07, planning.
- Decision: Console tabs on the integrate pages keep their inline command blocks; the new CLI tab is additive.
  Rationale: each tab is a self-contained path (the register/verify snippets already repeat across tabs); the Console tab's click-path-plus-equivalent-command pairing is deliberate and just shipped. Stripping it is separate cleanup if ever wanted.
  Date/Author: 2026-07-07, planning.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

### Repo layout

The Mintlify site root is `docs/` inside this repo: pages live at `docs/docs/...` relative to the repo root, navigation in `docs/docs.json`, shared snippets in `docs/snippets/`. Site URLs drop the `docs/` prefix (page `docs/data-uploads/uploads.mdx` → `/data-uploads/uploads`). Lint (`./scripts/lint.sh`) runs a Go prose linter (including a `heading-case` rule enforcing sentence-case headings), ESLint for MDX, CSpell, and `mint openapi-check`; `./scripts/preflight.sh` runs all of that plus lint-tool tests.

### Current state (workstream A)

`docs/docs.json` has, inside the "Data Uploads" group, a nested "Upload rules" group with four pages: `data-uploads/upload-rules/{overview,sources,destinations,manage}`. Their content:

- `overview.mdx`: definition snippet (`/snippets/definitions/upload-rule.mdx`), intro with the YAML rule-file example, **Properties** (ParamFields: `upload_collection_id`, `digest`, `source`, `destination`), **Immutability and releases**, **Git provenance**.
- `sources.mdx`: intro + ParamFields `glob`, `stability_window_secs`.
- `destinations.mdx`: intro (destination config vs. Buckets pages) + ParamFields `bucket_id`, `bucket_name`, `path`, `delete_policy`, then **Path template variables** (variable table + validation bullets).
- `manage.mdx`: **Create an upload rule** (pointer to defining-releases + dedup note), **View an upload rule** (release → rule details; has a role/scope TODO comment and a dashboard-screenshot TODO), **Git commit metadata** (commit-link accrual), **Editing and deleting rules** (immutable; retire by releasing without the rule).

Inbound links to the four URLs (complete list, from `grep -rn 'data-uploads/upload-rules' docs/`):

- `docs/docs.json` (nav, lines ~105–108)
- `docs/snippets/definitions/upload-collection.mdx:1` → `/overview`
- `docs/developers/agent/overview.mdx:11` → `/overview`
- `docs/data-uploads/overview.mdx:7` (prose) and `:40` (Card href) → `/overview`
- `docs/data-uploads/uploads.mdx:10,36` → `/overview`; `:55` → `/destinations#path-template-variables`; `:181` → `/sources`
- `docs/data-uploads/defining-releases.mdx:15` → `/overview`; `:65,70` → `/sources`; `:81,86` → `/destinations`; `:101` → `/manage#git-commit-metadata`
- `docs/data-uploads/buckets/overview.mdx:65` → `/destinations`
- plus self-references inside the four merged pages (become in-page anchors).

Beware false positives: `./upload-rules/` appears in CLI example paths (`docs/data-uploads/defining-releases.mdx`, `docs/snippets/references/cli/releases/create/flags.mdx`) — those are shell paths, not links; leave them alone.

Glob validation rules (for the new **Globs** subsection) — verified in `backend/internal/configs/domain/uploadrules/spec.go` (`validateSrcGlob`): the glob must be non-empty, start with `/`, be at most 1024 bytes, contain no control characters, no `..` segments, and no empty segments (`//` or trailing `/`). Glob *matching* semantics (e.g. `**`) are not defined server-side; do not invent them.

### Current state (workstream B)

`docs/cfg-mgmt/schemas/overview.mdx` sections: Properties (all ParamFields carry `<ImmutableBadge />`; `digest` describes canonicalized-content hashing), Schema languages, File formats, Validating instances, Empty schemas. `docs/cfg-mgmt/schemas/manage.mdx` covers creating (pointer to defining-releases) and viewing (release → schema details; the **Metadata** tab shows "the schema's metadata and commit history").

Verified behavior to document:

- Schemas are content-identified and deduplicated **within their config type**: creation is find-or-create on (config type, digest) — `backend/internal/configs/services/config_schemas/create.go` (`FindByCfgTypeAndDigest`, `FindOrCreateWithJoin`). Canonicalization ignores comments/whitespace (digest ParamField; `references/cli/release-create.mdx:16-18` documents the idempotent creation).
- Git provenance exists: `ConfigSchema` has expandable `config_schema_git_commits` (`../openapi/apis/configs/components/schemas/config-schema.yaml:102`), each link carrying `file_paths` (plural — CUE package schemas can span multiple files; JSON Schema is a single file) and a `git_commit_id` (`config-schema-git-commit.yaml`). A git commit records SHA, message, repository owner/name/type/URL, commit URL (`git-commit.yaml`).
- Dedup hits **accrue** commit links: `FindOrCreateWithJoin` joins a new `ConfigSchemaGitCommit` when the same digest is released from a not-yet-linked commit — same behavior the upload-rules docs describe for rules.
- The CLI's Git-capture requirements are already documented at `/cfg-mgmt/defining-releases#git-commit`; link, don't restate.

### Current state (workstream C)

`docs/data-uploads/buckets/integrate/gcs.mdx` and `aws.mdx` each have a two-tab `<Tabs>` block (Console, Terraform) inside "Connecting your bucket". Both tabs reuse shared step snippets `docs/snippets/data-uploads/buckets/{register,verify}-{gcs,aws}.mdx`. The Console tabs already embed the equivalent `gcloud`/`aws` commands per step; the Terraform tabs provision the same resources. The AWS flow is two-phase: role exists first → register in Miru (issues the external ID) → tighten the trust policy with the external ID → verify. Cross-check fixtures: `../infra/deploy/terraform/gcp-bktverify-test.tf` (same WIF pool/provider/attribute-mapping/condition/impersonation-binding shape) and `../infra/cicd/tooling/bktverify-test.tf` (same trust-policy + `sts:ExternalId` condition shape).

## Plan of Work

### Milestone 1 — merged upload-rules page

**Create `docs/data-uploads/upload-rules.mdx`** (frontmatter `title: "Upload rules"`; imports `ImmutableBadge`, `UploadRuleDef`) with exactly this outline and content mapping:

| Section | Content source |
| --- | --- |
| (intro, no heading) | `<UploadRuleDef />` + overview.mdx's intro sentence (a rule = `source` + `destination`, belongs to a collection) + one sentence: rules are created by releasing rule files via the CLI (link `/data-uploads/defining-releases`) — absorbs manage.mdx "Create an upload rule". Drop the YAML block from the intro (moves to File formats). |
| `## Properties` | overview.mdx Properties ParamFields verbatim; repoint `source`/`destination` links to in-page `#sources` / `#destinations`. |
| `## File formats` | New short section: each rule is defined as a single YAML **rule file** in your Git repository; the YAML example from overview.mdx's intro; note rule files reference the collection by `collection_slug` and the bucket by name (`destination.bucket`); link to `/data-uploads/defining-releases#upload-rule-files` for the field-by-field reference. Keep to a paragraph + code block — do not duplicate defining-releases. |
| `## Sources` | sources.mdx intro + both ParamFields (`glob`, `stability_window_secs`). |
| `### Globs` | New subsection under Sources: glob constraints verified from `spec.go` — absolute (must start with `/`), ≤ 1024 bytes, no control characters, no `..` segments, no empty segments (`//` or trailing `/`); example `/var/log/robot/*.log`. |
| `## Destinations` | destinations.mdx intro (destination config vs. Buckets distinction) + four ParamFields; the `path` ParamField's "supported variables" link → `#path-templates`. |
| `### Path templates` | destinations.mdx "Path template variables" section (variable table + validation bullets), demoted to H3 and renamed. |
| `## Immutability and releases` | overview.mdx section of the same name + manage.mdx "Editing and deleting rules" merged: immutable, shipped via releases, never edited in place; digest dedup within the collection (re-release never duplicates); rules cannot be edited or deleted — retire a rule by shipping a release that no longer includes it. |
| `## Git provenance` | overview.mdx "Git provenance" + manage.mdx "Git commit metadata" merged: each rule records the commits it was released from (commit + rule-file path relative to repo root); re-releasing the identical rule from a later commit accrues a link to that commit; link `/data-uploads/defining-releases#git-commit` for what the CLI captures. Drop the old self-link to `manage#git-commit-metadata`. |
| `## Viewing upload rules` | manage.mdx "View an upload rule": releases list their rules alongside schemas; open a release → select a rule to inspect its source, destination, digest, collection. Preserve both TODO comments (role/scope; dashboard screenshots). Repoint intra-links to `#sources` / `#destinations`. |

**Delete** `docs/data-uploads/upload-rules/overview.mdx`, `sources.mdx`, `destinations.mdx`, `manage.mdx` (and the now-empty directory).

**Update `docs/docs.json`** — target Data Uploads block:

    {
     "group": "Data Uploads",
     "pages": [
      "data-uploads/overview",
      "data-uploads/upload-collections",
      "data-uploads/upload-rules",
      {
       "group": "Buckets",
       "pages": [
        "data-uploads/buckets/overview",
        "data-uploads/buckets/manage",
        {
         "group": "Integrate",
         "pages": [
          "data-uploads/buckets/integrate/gcs",
          "data-uploads/buckets/integrate/aws"
         ]
        }
       ]
      },
      "data-uploads/uploads",
      "data-uploads/defining-releases"
     ]
    },

**Repoint every inbound link** (old → new):

- `.../upload-rules/overview` → `/data-uploads/upload-rules` — in `docs/snippets/definitions/upload-collection.mdx`, `docs/developers/agent/overview.mdx`, `docs/data-uploads/overview.mdx` (prose + Card), `docs/data-uploads/uploads.mdx:10,36`, `docs/data-uploads/defining-releases.mdx:15`
- `.../upload-rules/overview#properties` → `/data-uploads/upload-rules#properties` (none outside the merged pages today; check anyway)
- `.../upload-rules/sources` → `/data-uploads/upload-rules#sources` — `docs/data-uploads/uploads.mdx:181`, `docs/data-uploads/defining-releases.mdx:65,70`
- `.../upload-rules/destinations` → `/data-uploads/upload-rules#destinations` — `docs/data-uploads/defining-releases.mdx:81,86`
- `.../upload-rules/destinations#path-template-variables` → `/data-uploads/upload-rules#path-templates` — `docs/data-uploads/uploads.mdx:55`, and `docs/data-uploads/buckets/overview.mdx:65` (its bare `/destinations` link is about the path template's `{device_id}` variable, so target `#path-templates`)
- `.../upload-rules/manage#git-commit-metadata` → `/data-uploads/upload-rules#git-provenance` — `docs/data-uploads/defining-releases.mdx:101` (also update the link text "Git commit metadata" → "Git provenance")

### Milestone 2 — config schemas sections

Append two sections to `docs/cfg-mgmt/schemas/overview.mdx` after "Empty schemas", mirroring the tone and paragraph order of the merged upload-rules page:

- `## Immutability and releases`: schemas are immutable — created once as part of a [release](/fundamentals/releases/overview) via the CLI and never edited in place; to change what a config type accepts, release a revised schema file and the old schema simply stops being referenced. Schemas are deduplicated by digest within their config type: releasing a schema whose *content* is equivalent to an existing one (comments, whitespace, and formatting ignored — see the `digest` property above) attaches the existing schema to the release instead of creating a duplicate, so re-releasing unchanged schema files is always safe. Do **not** claim schemas cannot be deleted (see Decision Log).
- `## Git provenance`: schemas are defined as files in your Git repository and created by [releasing them via the CLI](/cfg-mgmt/defining-releases). Each schema records the Git commits it was released from — the commit and the schema file's path(s) relative to the repository root (a CUE package schema may span multiple files; a JSON Schema is a single file) — so any schema in Miru traces back to the exact lines that defined it. If the same schema content is released again from a later commit, the existing schema accrues a link to that commit too. Link `/cfg-mgmt/defining-releases#git-commit` for the captured metadata and `/cfg-mgmt/schemas/manage#view-a-schema` (Metadata tab) for viewing the commit history.

### Milestone 3 — CLI tabs on the integrate pages

On both pages, update the "two ways" sentence in the intro before `<Tabs>` to describe three tabs (Console click-by-click; **CLI** runs the same setup as plain `gcloud`/`aws` commands; Terraform as code), then insert `<Tab title="CLI">` between Console and Terraform. Each CLI tab is a `<Steps>` sequence ending with the shared register/verify snippets, exactly like the other tabs. Reuse the verification-probe `<Note>` where the other tabs have it.

**`gcs.mdx` CLI tab** (commands identical to the ones already embedded in the Console tab, in the same order, plus two lookups):

1. Create the bucket — `gcloud storage buckets create gs://<bucket-name> --project=<project-id> --location=<region>`
2. Create the Workload Identity Pool — `gcloud iam workload-identity-pools create miru-pool --project=<project-id> --location=global --display-name="Miru pool"`
3. Add the AWS provider — the `gcloud iam workload-identity-pools providers create-aws miru-provider ...` command already embedded in the Console tab's "Add an AWS provider" step, with its exact `--attribute-mapping` and `--attribute-condition` strings; follow with `gcloud iam workload-identity-pools providers describe miru-provider --project=<project-id> --location=global --workload-identity-pool=miru-pool --format="value(name)"` to print the full resource name — this is the `wip_provider` value to register.
4. Create the uploader service account — `gcloud iam service-accounts create miru-uploader --project=<project-id> --display-name="Miru uploader"`
5. Grant object creation on the bucket — `gcloud storage buckets add-iam-policy-binding gs://<bucket-name> --member="serviceAccount:miru-uploader@<project-id>.iam.gserviceaccount.com" --role="roles/storage.objectCreator"` (+ probe-cleanup Note)
6. Bind Workload Identity User — first `gcloud projects describe <project-id> --format="value(projectNumber)"` (the principal set needs the project *number*), then `gcloud iam service-accounts add-iam-policy-binding miru-uploader@<project-id>.iam.gserviceaccount.com --project=<project-id> --role="roles/iam.workloadIdentityUser" --member="principalSet://iam.googleapis.com/projects/<project-number>/locations/global/workloadIdentityPools/miru-pool/attribute.aws_role/<miru-aws-role-arn>"`
7. `<RegisterGcs />` then `<VerifyGcs />`

**`aws.mdx` CLI tab** (preserve the two-phase external-ID ordering; make each step copy-pasteable by writing the policy JSON via heredoc in the same code block as the command):

1. Create the bucket — `aws s3api create-bucket --bucket <bucket-name> --region <region> --create-bucket-configuration LocationConstraint=<region>`, with a one-line note: in `us-east-1`, omit `--create-bucket-configuration` (the API rejects a location constraint there).
2. Create the cross-account role — heredoc `trust-policy.json` (principal = `<miru-integration-role-arn>`, `sts:AssumeRole`, **no** external-ID condition yet) then `aws iam create-role --role-name miru-uploader --assume-role-policy-document file://trust-policy.json`; the returned ARN is the `role_arn` to register.
3. Grant write access — heredoc `put-policy.json` (`s3:PutObject` + `s3:AbortMultipartUpload` on `arn:aws:s3:::<bucket-name>/*`) then `aws iam put-role-policy --role-name miru-uploader --policy-name miru-put-object --policy-document file://put-policy.json` (+ probe-cleanup Note; carry over the SSE-KMS TODO comment).
4. `<RegisterAws />` (Miru issues the external ID)
5. Tighten the trust policy — heredoc the trust policy again with the `"Condition": {"StringEquals": {"sts:ExternalId": "<external-id>"}}` block, then `aws iam update-assume-role-policy --role-name miru-uploader --policy-document file://trust-policy.json`.
6. `<VerifyAws />`

All command names and flags above already appear verbatim in the shipped Console tabs except: `providers describe --format="value(name)"`, `projects describe --format="value(projectNumber)"` (both standard gcloud), and the heredoc wrappers. Cross-check resource shapes against the infra fixtures named in Context. Do not add Miru-CLI commands — this tab is cloud-provider CLIs only.

## Concrete Steps

All commands run from the repo root (`/home/ben/miru/workbench5/repos/docs`) unless noted.

Milestone 1:

1. Write `docs/data-uploads/upload-rules.mdx` per the mapping table; `git rm` the four old pages.
2. Edit the Data Uploads block in `docs/docs.json` to the target block above.
3. Repoint the inbound links listed in Plan of Work (edit each file; exact old→new strings given there).
4. Check: `grep -rn 'upload-rules/overview\|upload-rules/sources\|upload-rules/destinations\|upload-rules/manage' docs/` → no output; `test ! -d docs/data-uploads/upload-rules && echo ok` → `ok`.
5. `./scripts/lint.sh` → ends with `All documentation lint checks passed.` (add any newly flagged terms to `cspell.json` only if legitimately spelled).
6. Commit: `git add -A && git commit -m "docs(data-uploads): collapse upload rules into a single page"`.

Milestone 2:

1. Append the two sections to `docs/cfg-mgmt/schemas/overview.mdx`.
2. `./scripts/lint.sh` → passes.
3. Commit: `git commit -am "docs(cfg-mgmt): add immutability and git provenance sections to schemas overview"`.

Milestone 3:

1. Edit `docs/data-uploads/buckets/integrate/gcs.mdx`: update the tabs intro paragraph, insert the CLI tab after `</Tab>` of Console.
2. Same for `aws.mdx`.
3. `./scripts/lint.sh` → passes.
4. Commit: `git commit -am "docs(data-uploads): add cloud-CLI tabs to bucket integrate pages"`.

Final: `./scripts/preflight.sh` → every stage passes (Lint Smoke Tests, Go Lint, Go Coverage, Lint, Audit, Shell Script Tests).

## Validation and Acceptance

Test steps (run from repo root; all must pass before the work is done):

1. `./scripts/lint.sh` exits 0 after each milestone; `./scripts/preflight.sh` exits 0 at the end.
2. No stale subpage references: `grep -rn 'upload-rules/overview\|upload-rules/sources\|upload-rules/destinations\|upload-rules/manage' docs/` returns nothing (the `./upload-rules/` shell paths in CLI examples don't match this pattern and stay).
3. Nav shape: `python3 -c "import json; s=json.dumps(json.load(open('docs/docs.json'))['navigation']); assert '\"data-uploads/upload-rules\"' in s and 'upload-rules/overview' not in s; print('nav ok')"` prints `nav ok`.
4. Anchor targets exist — every heading the repointed links rely on is present: `grep -n '^## Properties\|^## File formats\|^## Sources\|^### Globs\|^## Destinations\|^### Path templates\|^## Immutability and releases\|^## Git provenance\|^## Viewing upload rules' docs/data-uploads/upload-rules.mdx` shows all nine, in that order.
5. Schemas sections: `grep -n '^## Immutability and releases\|^## Git provenance' docs/cfg-mgmt/schemas/overview.mdx` shows both.
6. Tab structure: `grep -c '<Tab title=' docs/data-uploads/buckets/integrate/gcs.mdx` → `3` (same for `aws.mdx`), with `CLI` second; each CLI tab contains the register and verify snippet components (`grep -n 'RegisterGcs\|VerifyGcs' ...` inside the tab for GCS; `RegisterAws` appears **between** the write-policy step and the trust-tighten step for AWS).
7. Behavior check (if `mint dev` is available): `/data-uploads/upload-rules` renders with the sidebar entry replacing the old group in the same position; clicking the repointed links from `/data-uploads/uploads` and `/data-uploads/defining-releases` lands on the correct section anchors; both integrate pages render three tabs.

Acceptance: all seven checks pass; content claims match the verified sources named in Context (no claim about schema deletion; no invented glob matching semantics; CLI commands mirror the Terraform resources one-for-one).

## Idempotence and Recovery

All steps are plain-text edits on a feature branch; each is safely re-runnable. Every milestone is one commit, so `git revert <sha>` (or `git checkout 35e7cb6 -- <path>` before committing) rolls back a milestone independently. The only deletion (the four old pages) is recoverable from git history at `35e7cb6`. If lint flags a heading after renaming, adjust the heading text — the anchors in step 4 of Validation are derived from the final headings, so re-run the repoint greps after any heading change.

# Anchor bucket provider pages to the config field and document console + Terraform connection workflows

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (this repo, branch `feat/data-recording`) | read-write | Rework `docs/data-uploads/buckets/integrate/{gcs,aws}.mdx`: rename each page's Properties section to Config (anchored to the bucket object's `config` field), and restructure "Connecting your bucket" into two tabbed workflows (Console, Terraform) with shared register/verify step snippets. Possibly touch `cspell.json`. |
| `/home/ben/miru/workbench5/repos/openapi` | read-only | Source of truth for the `GcsBucketConfig` / `AwsBucketConfig` shapes (`apis/configs/components/schemas/bucket.yaml`) and the bucket endpoints (`apis/configs/paths/buckets.yaml`). |
| `/home/ben/miru/workbench5/repos/cli-private` | read-only | Verified the CLI has **no bucket commands** (see Context). |
| `/home/ben/miru/workbench5/repos/infra` | read-only | Source of truth for the working GCS WIF and AWS cross-account role setups the Terraform snippets are derived from. |
| `/home/ben/miru/workbench5/repos/backend` | read-only | Confirmed the verification probe prefix (`internal/configs/bktverify/probe.go`). |

This plan lives in `docs/plans/` because all edits are docs-repo edits. Paths below are relative to `/home/ben/miru/workbench5/repos/docs` unless prefixed otherwise.

## Purpose / Big Picture

After this change, a reader of `docs.mirurobotics.com/data-uploads/buckets/integrate/gcs` (or `/aws`):

1. Sees a **Config** section that explicitly documents the provider-specific shape of the bucket object's `config` property (documented at `/data-uploads/buckets/overview#properties`), field-for-field faithful to `GcsBucketConfig` / `AwsBucketConfig` in the OpenAPI spec — including the `provider` discriminator that lives *inside* `config`.
2. Can follow "Connecting your bucket" through either of two complete workflows, presented as tabs: **Console** (the existing click-by-click console + cloud-CLI steps) or **Terraform** (one HCL config provisioning the customer-side cloud resources). Both tabs end with the same "register in Miru dashboard → verify" steps, and both preserve the required ordering — for AWS: role exists → register in Miru returns `external_id` → tighten trust policy → verify; for GCS: all GCP resources exist → register → verify.
3. Is never promised a CLI command that doesn't exist — the Miru CLI has no bucket support today, so there is no CLI tab (see Decision Log).

## Progress

- [x] Milestone 1: Rename Properties → Config on both integrate pages, anchored to the `config` field
- [ ] Milestone 2: GCS — tabbed Console/Terraform workflows + WIF corrections from infra
- [ ] Milestone 3: AWS — tabbed Console/Terraform workflows with two-phase external-ID flow

## Surprises & Discoveries

(Add entries as work proceeds. Findings from planning research are in Context.)

## Decision Log

- Decision: Rename the section to `## Config` (anchor `#config`) with a lead-in linking the bucket object's `config` property at `/data-uploads/buckets/overview#properties`.
  Rationale: matches the field name exactly; no inbound links target `integrate/{gcs,aws}#properties` (verified: `grep -rn 'integrate/gcs#\|integrate/aws#' docs` returns nothing), so the heading rename breaks nothing.
  Date/Author: 2026-07-07, planning.
- Decision: Add a `provider` const ParamField to each Config section (`gcs` / `aws`).
  Rationale: `GcsBucketConfig`/`AwsBucketConfig` both require `provider`, and the create request (`CreateBucketRequest` in `openapi/apis/configs/components/requests/bucket.yaml`) has **no top-level provider** — the provider is chosen via `config.provider`. Documenting it here makes the Config section the complete, faithful shape.
  Date/Author: 2026-07-07, planning.
- Decision: **No CLI workflow tab.** `cli-private/internal/commands/root/root.go` registers only `login`, `version`, `env`, and `release` commands; buckets appear in the CLI solely as upload-rule destination references (`internal/services/uplrules/push.go`). Do not document the raw platform API either: the bucket endpoints exist in the OpenAPI *source* (`openapi/apis/configs/paths/buckets.yaml` — create/list/get/delete/verify/archive/unarchive) but the docs' published platform-API reference (`docs/references/platform-api/2026-05-06.yaml`) contains no bucket paths, so an API tab would point at endpoints invisible in our own reference. Revisit when CLI/API bucket support ships.
  Date/Author: 2026-07-07, planning.
- Decision: Present the workflows as top-level `<Tabs>` inside "Connecting your bucket" (`Console` / `Terraform`), each containing a full `<Steps>` sequence, rather than per-step tabs or separate pages (the `fundamentals/devices/provision/*` pattern).
  Rationale: Terraform replaces *all* cloud-side steps with one config, so per-step tabs don't fit; separate pages would triple the page count for ~60 shared lines. Full sequences per tab keep the AWS two-phase ordering intact inside each tab.
  Date/Author: 2026-07-07, planning.
- Decision: Share the "Register the bucket in Miru" and "Verify the bucket" steps between tabs via snippets under `docs/snippets/data-uploads/buckets/`, and reword them to reference values **by name** (e.g. "the `wip_provider` provider resource name"), not by step number.
  Rationale: prevents drift between tabs; step numbers differ between tabs (register is step 7 in the GCS Console tab but step 2 in the Terraform tab).
  Date/Author: 2026-07-07, planning.
- Decision: Fix the GCS manual steps to match the **working** WIF setup in infra: attribute mapping `google.subject = assertion.arn` plus an `attribute.aws_role` extraction, attribute condition `assertion.arn.startsWith('<miru-aws-role-arn>/')`, and an impersonation binding of **`roles/iam.workloadIdentityUser`** (not `roles/iam.serviceAccountTokenCreator`) to the `principalSet` keyed on `attribute.aws_role`. This resolves the page's two `TODO (product)` comments about the CEL condition and the principalSet selector.
  Rationale: `infra/deploy/terraform/gcp-bktverify-test.tf` is the fixture the backend's live GCP verification test passes against — it is the only known-working reference. The page's current `google.subject == '<miru-aws-role-arn>'` equality would fail in practice: the federated subject is the *assumed-role STS ARN with a per-session suffix* (hence `startsWith` and the suffix-stripping `attribute.aws_role` mapping).
  Date/Author: 2026-07-07, planning.
- Decision: Keep docs' minimal AWS permission policy (`s3:PutObject` + `s3:AbortMultipartUpload`, optional `s3:DeleteObject` on `.miru/probe/*`) and GCS `roles/storage.objectCreator`, even though the infra fixtures grant more (`objectAdmin`; `DeleteObject`/`ListBucket`/`GetBucketLocation`).
  Rationale: `backend/internal/configs/bktverify/probe.go` shows verification only uploads a zero-byte object under `.miru/probe/` and best-effort deletes it; the fixtures' extra grants are test-suite convenience, not requirements. The AWS region comes from `config.region`, so no `GetBucketLocation` is needed.
  Date/Author: 2026-07-07, planning.
- Decision: Resolve the "confirm where registration surfaces — dashboard or CLI" TODOs to **dashboard**; keep a trimmed `TODO (product)` about the exact form of the Miru-provided identity values shown at registration (account ID vs role ARN vs assumed-role ARN; where the AWS `external_id` is displayed).
  Rationale: CLI bucket support verifiably does not exist; the dashboard-display details are still unconfirmed.
  Date/Author: 2026-07-07, planning.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

This repo is the Mintlify docs site: content under `docs/`, nav + redirects in `docs/docs.json`, prose linter in `tools/lint/`. Internal links are absolute site paths without `.mdx`; anchors are kebab-cased heading text. Snippet components are imported by absolute path (e.g. `import Foo from '/snippets/definitions/bucket.mdx';`) and can contain any MDX, including `<Step>` elements. Tooling: `./scripts/lint.sh` (Go prose linter + ESLint + cspell + `mint openapi-check`; expect final line `All documentation lint checks passed.`) and `./scripts/preflight.sh` (lint smoke tests, Go coverage gate, lint.sh, `pnpm audit`, bats). The headingcase rule checks only `title:` frontmatter and `#` headings — `<Tab title>` / `<Step title>` attributes are exempt (existing pages prove this). `cspell.json` has a `words` allowlist; Terraform-specific tokens (e.g. `tfvars`) may need adding.

The two target pages (current state, after the completed restructure plan `plans/completed/20260707-buckets-section-restructure.md`):

- `docs/data-uploads/buckets/integrate/gcs.mdx` — intro (WIF summary), `## Properties` (ParamFields: `project_id`, `wip_provider`, `service_account_email`), `## Connecting your bucket` (exchange table + 8 `<Steps>`: bucket, pool, AWS provider, uploader SA, objectCreator grant, impersonation grant, register in Miru, verify), `### How device credentials are minted at runtime`. Carries `TODO (product)` comments on the attribute condition (line ~118), the principalSet selector (~206), and registration surface (~225).
- `docs/data-uploads/buckets/integrate/aws.mdx` — intro (AssumeRole summary), `## Properties` (`region`, `role_arn`, `external_id`), `## Connecting your bucket` (exchange table + 6 `<Steps>`: bucket, IAM role with open trust policy, write-access policy, register in Miru (Miru issues `external_id`), add external ID to trust policy, verify), `### How device credentials are minted at runtime`. `TODO (product)` on registration surface (~186) and SSE-KMS (~158, keep).

Schema source of truth — `/home/ben/miru/workbench5/repos/openapi/apis/configs/components/schemas/bucket.yaml`:

- `GcsBucketConfig` (required: `provider` const `"gcs"`, `project_id`, `wip_provider`, `service_account_email`).
- `AwsBucketConfig` (required: `provider` const `"aws"`, `region`, `role_arn`, `external_id`); `CreateAwsBucketConfig` omits `external_id` (server-issued, returned in the response) — the page's existing prose on `external_id` already says this and stays.
- The pages' existing ParamFields already match field-for-field except the missing `provider`.
- (Noted upstream drift, out of scope here: the schema descriptions still say "V4-presigned upload URLs" while the actual flow mints downscoped tokens — an openapi-repo fix, not a docs edit.)

Infra source of truth for the Terraform snippets:

- GCS: `infra/deploy/terraform/gcs-integration-role.tf` (Miru-side `aws_iam_role.gcs_integration`, name `MiruGCSIntegrationRole`, deliberately permission-less identity) and `infra/deploy/terraform/gcp-bktverify-test.tf` (simulated **customer** side: `google_storage_bucket`, `google_service_account` `bktverify-test-customer`, `google_storage_bucket_iam_member`, `google_iam_workload_identity_pool` + `..._provider` with the `attribute_mapping`/`attribute_condition` quoted in the Decision Log, `google_service_account_iam_member` binding `roles/iam.workloadIdentityUser` to `principalSet://iam.googleapis.com/<pool>/attribute.aws_role/arn:aws:sts::<acct>:assumed-role/MiruGCSIntegrationRole`).
- AWS: `infra/deploy/terraform/s3-integration-role.tf` (Miru-side `MiruS3IntegrationRole`, permission `sts:AssumeRole` on `*`) and `infra/cicd/tooling/bktverify-test.tf` (simulated **customer** side: `aws_s3_bucket`, `aws_iam_role` `MiruBktVerifyTestCustomerRole` whose trust policy allows `MiruS3IntegrationRole` with an `sts:ExternalId` `StringEquals` condition, plus a bucket-scoped write policy).
- **No Miru Terraform provider exists** (searched infra and openapi; the Stainless sdkgen targets are cli/device/platform only). Terraform therefore covers only the customer's cloud-side resources; the Miru registration step remains a dashboard action, and each Terraform tab must say so explicitly.

Backend confirmation: `backend/internal/configs/bktverify/probe.go` writes `".miru/probe/" + rand.String()` and best-effort deletes it — the pages' probe-cleanup `<Note>`s are already correct (the `.miru-connectivity-check/` prefix in the infra fixtures is stale infra, not a docs bug).

## Plan of Work

### Milestone 1 — Config section rename (both pages)

In `docs/data-uploads/buckets/integrate/gcs.mdx` and `aws.mdx`:

1. Rename `## Properties` → `## Config`.
2. Replace the section lead-in with one that anchors the section to the `config` field, e.g. (GCS; AWS analogous):

       These are the fields of the bucket's
       [`config`](/data-uploads/buckets/overview#properties) property when the
       bucket's `provider` is `gcs`. The top-level bucket
       [`name`](/data-uploads/buckets/overview#properties) is the name of the GCS
       bucket itself; `config` carries everything provider-specific:

3. Add a leading `provider` ParamField to each (`type="string"`, `required`, `<ImmutableBadge />`): "The provider discriminator, always `gcs` for this shape." / "…always `aws`…". Keep the existing ParamFields untouched (already schema-faithful).

### Milestone 2 — GCS workflows

1. Create shared snippets (content lifted from the current steps 7–8 of `gcs.mdx`, rewording value references by name instead of step number; keep `<Frame>` images and `{/* TODO (screenshot) */}` comments verbatim):
   - `docs/snippets/data-uploads/buckets/register-gcs.mdx` — the `<Step title="Register the bucket in Miru">`; resolve its registration-surface TODO per the Decision Log (dashboard; keep a trimmed TODO about the exact displayed identity values).
   - `docs/snippets/data-uploads/buckets/verify-gcs.mdx` — the `<Step title="Verify the bucket">`.
2. Restructure `## Connecting your bucket` in `gcs.mdx`: keep the intro + exchange table; replace "Each gives the GCP Console path and the equivalent gcloud command" with a sentence introducing the two tabs (Console for click-by-click + `gcloud`, Terraform for infrastructure-as-code); then:

       <Tabs>
         <Tab title="Console">
           <Steps> existing steps 1–6 (corrected, below) + <RegisterGcs /> + <VerifyGcs /> </Steps>
         </Tab>
         <Tab title="Terraform">
           lead-in: Terraform provisions the GCP-side resources only — registering
           the bucket in Miru still happens in the dashboard.
           <Steps>
             <Step title="Provision the GCP resources"> HCL below + terraform init/apply;
               the outputs are the `wip_provider` and `service_account_email` you register. </Step>
             <RegisterGcs /> + <VerifyGcs />
           </Steps>
         </Tab>
       </Tabs>

3. Correct the Console steps per the Decision Log (source: `infra/deploy/terraform/gcp-bktverify-test.tf`):
   - Step "Add an AWS provider…": document the attribute mapping and replace the condition with `assertion.arn.startsWith('<miru-aws-role-arn>/')`, where `<miru-aws-role-arn>` is the assumed-role ARN Miru provides (`arn:aws:sts::<miru-aws-account-id>:assumed-role/...`). Update the `gcloud iam workload-identity-pools providers create-aws` command with `--attribute-mapping` and the new `--attribute-condition`; delete the resolved TODO.
   - Step "Let Miru's federated principal impersonate the service account": role becomes `roles/iam.workloadIdentityUser`; member becomes `principalSet://iam.googleapis.com/projects/<project-number>/locations/global/workloadIdentityPools/miru-pool/attribute.aws_role/<miru-aws-role-arn>`; delete the resolved TODO.
4. Terraform tab HCL (derived from the fixture; placeholders follow the page's `<…>` convention):

       variable "miru_aws_account_id" {
         description = "Miru's AWS account ID (provided by Miru)"
         type        = string
       }

       variable "miru_aws_role_arn" {
         description = "Assumed-role ARN of Miru's integration identity (provided by Miru)"
         type        = string
       }

       resource "google_storage_bucket" "uploads" {
         name                        = "<bucket-name>"
         project                     = "<project-id>"
         location                    = "<region>"
         uniform_bucket_level_access = true
       }

       resource "google_service_account" "miru_uploader" {
         project      = "<project-id>"
         account_id   = "miru-uploader"
         display_name = "Miru uploader"
       }

       resource "google_storage_bucket_iam_member" "miru_uploader_creator" {
         bucket = google_storage_bucket.uploads.name
         role   = "roles/storage.objectCreator"
         member = "serviceAccount:${google_service_account.miru_uploader.email}"
       }

       resource "google_iam_workload_identity_pool" "miru" {
         project                   = "<project-id>"
         workload_identity_pool_id = "miru-pool"
       }

       resource "google_iam_workload_identity_pool_provider" "miru_aws" {
         project                            = "<project-id>"
         workload_identity_pool_id          = google_iam_workload_identity_pool.miru.workload_identity_pool_id
         workload_identity_pool_provider_id = "miru-provider"

         attribute_mapping = {
           "google.subject"     = "assertion.arn"
           "attribute.aws_role" = "assertion.arn.contains('assumed-role') ? assertion.arn.extract('{account_arn}assumed-role/') + 'assumed-role/' + assertion.arn.extract('assumed-role/{role_name}/') : assertion.arn"
         }

         attribute_condition = "assertion.arn.startsWith('${var.miru_aws_role_arn}/')"

         aws {
           account_id = var.miru_aws_account_id
         }
       }

       resource "google_service_account_iam_member" "miru_wif_impersonation" {
         service_account_id = google_service_account.miru_uploader.name
         role               = "roles/iam.workloadIdentityUser"
         member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.miru.name}/attribute.aws_role/${var.miru_aws_role_arn}"
       }

       output "wip_provider" {
         value = google_iam_workload_identity_pool_provider.miru_aws.name
       }

       output "service_account_email" {
         value = google_service_account.miru_uploader.email
       }

   Mention the optional probe-cleanup grant (same `<Note>` as the Console step 5) in the step body or keep the Note inside the tab.

### Milestone 3 — AWS workflows

1. Snippets `docs/snippets/data-uploads/buckets/register-aws.mdx` and `verify-aws.mdx` (from current steps 4 and 6 of `aws.mdx`; register step keeps "Miru issues the bucket's external ID — copy it for the next step", rewords "role ARN from step 2" → "the role's ARN"; resolve the registration-surface TODO as in GCS).
2. Restructure `## Connecting your bucket` in `aws.mdx` the same way. Console tab = existing steps 1–3 + `<RegisterAws />` + existing step 5 (external-ID trust-policy tightening) + `<VerifyAws />` — ordering unchanged. Terraform tab preserves the same two-phase ordering:

       <Steps>
         1. Provision the AWS resources — first `terraform apply` with
            `miru_external_id` unset (trust policy allows Miru's integration role,
            no external-ID condition yet); output is the `role_arn` to register.
         2. <RegisterAws />  (Miru issues the external ID)
         3. Tighten the trust policy — set `miru_external_id` in terraform.tfvars
            and `terraform apply` again; the dynamic condition adds the
            `sts:ExternalId` requirement.
         4. <VerifyAws />
       </Steps>

3. Terraform tab HCL (derived from `infra/cicd/tooling/bktverify-test.tf`; the dynamic block is what makes the two-phase flow one config):

       variable "miru_integration_role_arn" {
         description = "Miru's integration role ARN (provided by Miru)"
         type        = string
       }

       variable "miru_external_id" {
         description = "External ID issued by Miru at registration; leave unset for the first apply"
         type        = string
         default     = null
       }

       resource "aws_s3_bucket" "uploads" {
         bucket = "<bucket-name>"
       }

       data "aws_iam_policy_document" "miru_trust" {
         statement {
           effect  = "Allow"
           actions = ["sts:AssumeRole"]

           principals {
             type        = "AWS"
             identifiers = [var.miru_integration_role_arn]
           }

           dynamic "condition" {
             for_each = var.miru_external_id == null ? [] : [1]
             content {
               test     = "StringEquals"
               variable = "sts:ExternalId"
               values   = [var.miru_external_id]
             }
           }
         }
       }

       resource "aws_iam_role" "miru_uploader" {
         name               = "miru-uploader"
         assume_role_policy = data.aws_iam_policy_document.miru_trust.json
       }

       data "aws_iam_policy_document" "miru_uploader_write" {
         statement {
           effect    = "Allow"
           actions   = ["s3:PutObject", "s3:AbortMultipartUpload"]
           resources = ["${aws_s3_bucket.uploads.arn}/*"]
         }
       }

       resource "aws_iam_role_policy" "miru_uploader_write" {
         name   = "miru-put-object"
         role   = aws_iam_role.miru_uploader.name
         policy = data.aws_iam_policy_document.miru_uploader_write.json
       }

       output "role_arn" {
         value = aws_iam_role.miru_uploader.arn
       }

   Keep the probe-cleanup and SSE-KMS notes available in this tab too (the SSE-KMS `TODO (product)` stays).

## Concrete Steps

All commands run from `/home/ben/miru/workbench5/repos/docs` on branch `feat/data-recording`. One commit per milestone.

Milestone 1:

    # edit docs/data-uploads/buckets/integrate/{gcs,aws}.mdx per Plan of Work M1
    ./scripts/lint.sh          # expect: All documentation lint checks passed.
    grep -n '## Properties' docs/data-uploads/buckets/integrate/*.mdx   # expect no output
    git add -A && git commit -m "docs(data-uploads): anchor bucket provider pages to the config field"

Milestone 2:

    mkdir -p docs/snippets/data-uploads/buckets
    # create register-gcs.mdx and verify-gcs.mdx; restructure gcs.mdx per M2
    ./scripts/lint.sh          # add any flagged Terraform tokens (e.g. tfvars) to cspell.json "words"
    git add -A && git commit -m "docs(data-uploads): document console and terraform workflows for GCS buckets"

Milestone 3:

    # create register-aws.mdx and verify-aws.mdx; restructure aws.mdx per M3
    ./scripts/lint.sh
    ./scripts/preflight.sh     # expect exit 0
    git add -A && git commit -m "docs(data-uploads): document console and terraform workflows for AWS buckets"

Do not push or open a PR as part of this plan.

## Validation and Acceptance

Run from `/home/ben/miru/workbench5/repos/docs`; all must pass before publishing:

1. `./scripts/lint.sh` exits 0 after each milestone; `./scripts/preflight.sh` exits 0 at the end (`pnpm run test:lint` too if `cspell.json` or `tools/lint/` changed).
2. Stale anchors/headings — each expects **no output**:

       grep -rn 'integrate/gcs#properties\|integrate/aws#properties' docs
       grep -n '## Properties' docs/data-uploads/buckets/integrate/*.mdx
       grep -rn 'serviceAccountTokenCreator' docs/data-uploads/buckets/integrate/gcs.mdx

3. Snippet wiring: all four files exist under `docs/snippets/data-uploads/buckets/` and each integrate page imports and renders its two; `grep -c 'Register the bucket in Miru' docs/data-uploads/buckets/integrate/gcs.mdx` returns 0 (the step text lives only in the snippet).
4. Content acceptance (read the rendered pages via `cd docs && mint dev`, or the raw MDX):
   - Each page has `## Config` opening with a lead-in that links `/data-uploads/buckets/overview#properties` and names the `provider` value, followed by ParamFields exactly matching the OpenAPI config shape (GCS: `provider`, `project_id`, `wip_provider`, `service_account_email`; AWS: `provider`, `region`, `role_arn`, `external_id`).
   - Each page's "Connecting your bucket" renders exactly two tabs, Console and Terraform; no CLI tab and no invented `miru` commands anywhere (`grep -n 'miru bucket' docs -r` → no output).
   - The AWS ordering inside **both** tabs is: role exists → register (external ID issued) → tighten trust policy → verify. The GCS ordering in both tabs is: all GCP resources exist → register → verify.
   - Each Terraform tab states that Terraform covers only the cloud-side resources and registration happens in the Miru dashboard.
   - GCS Console step for the provider uses the `startsWith` attribute condition + attribute mapping, and the impersonation step grants `roles/iam.workloadIdentityUser`; the two resolved `TODO (product)` comments are gone.

## Idempotence and Recovery

Every step is a git-tracked text edit in this repo; lint/preflight/grep checks are safe to re-run indefinitely. Before a milestone's commit, `git checkout -- <path>` (or `git reset --hard HEAD`) restores the last good state; after, `git revert <sha>` undoes exactly one milestone. `mkdir -p` and re-writing the snippet files are idempotent. If a later milestone reveals a problem in a snippet, fix the snippet — both tabs and both pages pick it up, which is the point of the extraction.

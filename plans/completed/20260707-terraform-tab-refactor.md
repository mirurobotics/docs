# Refactor the bucket-integration Terraform tabs (GCS and AWS)

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Rewrite the `<Tab title="Terraform">` sections of `docs/data-uploads/buckets/integrate/gcs.mdx` and `docs/data-uploads/buckets/integrate/aws.mdx`. No other files change. |
| `backend/` | read-only | `internal/configs/bktverify/{probe.go,aws.go,gcs.go}` confirm the exact permissions verification and uploads need. |
| `infra/` | read-only | `deploy/terraform/gcp-bktverify-test.tf` and `cicd/tooling/bktverify-test.tf` are the reference resource shapes for the drafted HCL. |

This plan lives in `docs/plans/` because all edits are confined to the `docs` repo. Work happens on branch `feat/data-recording` (clean at `98c357f`).

## Purpose / Big Picture

The two bucket-integration pages each have Console / CLI / Terraform tabs. Today the Terraform tabs contain working but half-parameterized HCL: literal `<bucket-name>`-style placeholders sit inside resource blocks, the GCS config emits only two of the four values the Miru registration form asks for, and the AWS config emits one of three. After this change, each Terraform tab is one complete, copy-paste config in the shape established by vendors like Datadog and CrowdStrike: a `variable` block for every knob (with descriptions and sensible defaults), least-privilege resources, and an `output` block whose names match the Miru registration form **field for field**. A user pastes the file, fills `terraform.tfvars`, runs `terraform apply`, and `terraform output` prints exactly what the registration form asks for — closing the loop between apply and register. The AWS tab keeps its two-phase external-ID flow but now says *why* the second apply exists.

## Progress

- [x] Milestone 1: Rewrite the GCS Terraform tab; lint; commit. (2026-07-08: spliced replacement block verbatim from this plan; `terraform fmt -check` + `validate` passed on the extracted HCL before commit; lint clean.)
- [x] Milestone 2: Rewrite the AWS Terraform tab; lint; commit. (2026-07-08: same splice-from-plan procedure; `terraform fmt -check` + `validate` passed on the extracted HCL before commit; lint clean.)
- [x] Validate drafted HCL with `terraform fmt`/`validate`; run `./scripts/preflight.sh` clean. (2026-07-08: both extracted configs — `fmt -check` clean, `validate` "Success! The configuration is valid."; preflight exit 0, all sections pass.)

## Surprises & Discoveries

- None of substance. The sandbox had network access, so `terraform init` downloaded the hashicorp providers and full `validate` ran for both configs — the no-network fallback documented in Concrete Steps step 6 was not needed. Validation order deviated slightly from Concrete Steps: each config was validated *before* its milestone commit rather than after both, so no amends were ever needed.
- To rule out transcription errors, the replacement `<Tab>` blocks were spliced into the MDX pages programmatically from this plan (strip four leading spaces) rather than retyped.

## Decision Log

- Decision: One complete self-contained config per tab, everything variable-driven, with outputs named identically to the registration form fields.
  Rationale: This is the dominant pattern across the vendor survey (see Research findings). Outputs that mirror the paste-back values are the Snowflake-community answer to the same "apply, then transcribe values into the vendor form" loop we have.
  Date/Author: 2026-07-07 / planner.

- Decision: Keep the AWS two-phase external-ID flow as `variable "miru_external_id"` with `default = null` plus a `dynamic "condition"` block and a second apply.
  Rationale: Datadog and CrowdStrike escape the two-phase dance only because they ship a Terraform provider whose resource returns the external ID in-graph; Miru has no Terraform provider. Snowflake, whose docs face the identical chicken-and-egg, sequences it as create → retrieve ID → tighten trust policy. The null-default variable is the Terraform-native version of that sequence and is already what the tab does — it just never explained why.
  Date/Author: 2026-07-07 / planner.

- Decision: Deliberately omit `terraform {}` / `provider {}` blocks from both configs; state in prose that an authenticated `aws` / `google` provider is assumed.
  Rationale: Users overwhelmingly paste these snippets into an existing workspace that already configures the provider; only the stock hashicorp providers are needed (unlike Datadog, which must declare its own provider). For AWS this also decides where the bucket's region comes from: the ambient provider config, reported back via the `region` output.
  Date/Author: 2026-07-07 / planner.

- Decision: Leave the shared register/verify snippets (`docs/snippets/data-uploads/buckets/register-{aws,gcs}.mdx`, `verify-{aws,gcs}.mdx`) untouched.
  Rationale: The snippet field lists (`name`/`region`/`role_arn`; `name`/`project_id`/`wip_provider`/`service_account_email`) already match the new output names exactly, so the Terraform tabs need only one connecting sentence ("`terraform output` prints exactly these values"), and the Console/CLI tabs that share the snippets are unaffected.
  Date/Author: 2026-07-07 / planner.

- Decision: Keep the permission grants exactly as they are today (AWS: `s3:PutObject` + `s3:AbortMultipartUpload` on `<bucket-arn>/*`; GCS: `roles/storage.objectCreator` on the bucket).
  Rationale: Verified against `backend/internal/configs/bktverify/probe.go` — the verification probe uploads a zero-byte object under `.miru/probe/` and delete-cleanup is explicitly best-effort (a warning log, not a failure). Runtime uploads are multipart (covered by `s3:PutObject`; abort needs its own action) / resumable (covered by object creation). The infra fixtures grant more (`s3:DeleteObject`, `s3:ListBucket`, `roles/storage.objectAdmin`), but they are CI test scaffolding, not the customer contract; customer docs stay least-privilege, with the existing `<Note>` covering optional delete grants.
  Date/Author: 2026-07-07 / planner.

## Outcomes & Retrospective

Completed 2026-07-08 on `feat/data-recording` in two milestone commits (`53e121d` GCS, `dfec86b` AWS), exactly as planned. Both Terraform tabs are now complete variable-driven configs whose outputs mirror the registration form field for field (GCS: `name`, `project_id`, `wip_provider`, `service_account_email`; AWS: `name`, `region`, `role_arn` — verified mechanically against the register snippets, which are byte-identical to `main`). Both extracted configs passed `terraform fmt -check` and `terraform validate` (Terraform v1.14.x, real provider downloads); `./scripts/lint.sh` passed after each milestone and `./scripts/preflight.sh` exited 0 at the end. No deviations from the drafted HCL; no scope creep beyond the two `<Tab title="Terraform">` elements and their Steps prose.

## Context and Orientation

The `docs` repo is a Mintlify documentation site; pages live under `docs/` and shared step snippets under `docs/snippets/`. The two pages in scope:

- `docs/data-uploads/buckets/integrate/gcs.mdx` — GCS integration via Workload Identity Federation (WIF): Miru presents an AWS identity to Google, a Workload Identity Pool (WIP) provider pinned to Miru's role accepts it, and Miru impersonates an uploader service account. The Terraform tab is the `<Tab title="Terraform">` element (lines 356–450 at `98c357f`). Its HCL declares two variables (`miru_aws_account_id`, `miru_aws_role_arn`) but hardcodes `<bucket-name>` / `<project-id>` / `<region>` placeholders inside resources, and outputs only `wip_provider` and `service_account_email` — the registration form also asks for `name` and `project_id`.
- `docs/data-uploads/buckets/integrate/aws.mdx` — S3 integration via cross-account STS AssumeRole gated by a Miru-issued external ID. The Terraform tab (lines 362–469 at `98c357f`) already has the `miru_external_id` null-default + `dynamic "condition"` two-phase pattern, but hardcodes `<bucket-name>`, outputs only `role_arn` (the form also asks for `name` and `region`), and never explains why the second apply exists.

Shared snippets embedded in every tab of both pages (do not edit):

- `docs/snippets/data-uploads/buckets/register-aws.mdx` — dashboard registration step; lists form fields `name`, `region`, `role_arn`; says Miru issues the external ID on create.
- `docs/snippets/data-uploads/buckets/register-gcs.mdx` — lists form fields `name`, `project_id`, `wip_provider`, `service_account_email`.
- `docs/snippets/data-uploads/buckets/verify-aws.mdx`, `verify-gcs.mdx` — the Verify step.

Ground truth for what the HCL must permit, in other repos (read-only):

- `backend/internal/configs/bktverify/probe.go` — verification probe: `Upload` of a zero-byte object at `.miru/probe/<random>`, then best-effort `Delete` (failure only logs a warning).
- `backend/internal/configs/bktverify/aws.go` — the AssumeRole chain: Miru's integration role assumes the customer `role_arn` presenting `cfg.ExternalID`.
- `infra/deploy/terraform/gcp-bktverify-test.tf` — canonical GCP resource shapes: `google_iam_workload_identity_pool{,_provider}` with the exact `attribute_mapping` / `attribute_condition` strings, `google_service_account_iam_member` with the `principalSet://...attribute.aws_role/...` member.
- `infra/cicd/tooling/bktverify-test.tf` — canonical AWS shapes: `data "aws_iam_policy_document"` for trust with `sts:ExternalId` condition, `aws_iam_role`, bucket-scoped policy resources.

### Research findings (vendor survey, 2026-07-07)

How established platforms ship Terraform for "connect your cloud account/bucket" onboarding:

- **Datadog** (<https://docs.datadoghq.com/integrations/guide/aws-terraform-setup/>, <https://registry.terraform.io/providers/DataDog/datadog/latest/docs/resources/integration_aws_account>): docs present a **complete, self-contained file** ready for `terraform apply` — `aws_iam_policy_document` trust policy, `aws_iam_policy`, `aws_iam_role`, attachments, plus the `datadog_integration_aws_account` provider resource. The external ID never touches the user: the trust policy references `datadog_integration_aws_account...auth_config.aws_auth_config_role.external_id` and Terraform resolves the cycle in one apply. Only possible because Datadog ships a Terraform provider. Notably their example uses literal `<ACCOUNT_ID>` placeholders rather than variables — the one part we do better.
- **CrowdStrike** (<https://registry.terraform.io/modules/CrowdStrike/cloud-registration/aws/latest>, <https://github.com/CrowdStrike/terraform-aws-cloud-registration>): a **published registry module**; the README is a full copy-paste example with a `terraform {}` block, sensitive credential variables with descriptions, provider config, and a `crowdstrike_cloud_aws_account` provider resource that supplies `external_id` into the module. Module exposes typed inputs/outputs. Wiz is similar via community providers/modules (e.g. <https://github.com/cisagov/cool-master-wiz>, `wiz_connector_aws` resources).
- **Snowflake** (<https://docs.snowflake.com/en/user-guide/data-load-s3-config-storage-integration>): console/SQL, **explicitly two-phase** like ours: create IAM policy + role with a *placeholder* external ID (`0000`), `CREATE STORAGE INTEGRATION`, then `DESC INTEGRATION` to read `STORAGE_AWS_IAM_USER_ARN` / `STORAGE_AWS_EXTERNAL_ID`, then edit the role's trust policy with the real values. Community Terraform ports of this flow run two applies and surface the paste-back values as **outputs**.
- **Fivetran** (<https://fivetran.com/docs/connectors/files/amazon-s3/setup-guide>): avoids two-phase by showing the auto-generated external ID **in the setup form before role creation**; minimal read-only policy (`s3:GetObject`, `s3:ListBucket`, `s3:GetBucketLocation`) with optional KMS grants for encrypted buckets; console-only, no Terraform in the connector docs. **Hightouch** (<https://hightouch.com/docs/security/aws>) likewise publishes its account ID + external ID up front.
- **Google WIF** (<https://cloud.google.com/blog/products/devops-sre/infrastructure-as-code-with-terraform-and-identity-federation/>, <https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/iam_workload_identity_pool>): official AWS→GCP federation examples use exactly our resource trio — pool, provider with `attribute_mapping`/`attribute_condition`, and a `roles/iam.workloadIdentityUser` binding on the impersonated service account.

Recurring patterns adopted (see Decision Log for rationale): one complete file; variables with descriptions/defaults up top; outputs that emit exactly the vendor-form paste-back values; least-privilege resource-scoped policies; vendor-issued external ID via null-default variable + second apply; comments sparingly, on the *why* (the two-phase dance, the write-only scoping). Patterns noted but **not** adopted: a published registry module and a vendor Terraform provider (Miru has neither today — registration stays in the dashboard; both tabs say so in their intro line).

## Plan of Work

Two milestones, one commit each, on `feat/data-recording`. Each milestone replaces the entire `<Tab title="Terraform">…</Tab>` element of one page; nothing outside those elements changes. In both replacement blocks below, every line is indented four extra spaces for plan formatting — strip exactly four leading spaces from each line when copying (the `<Tab>` tag itself sits at a two-space indent in the file, matching its siblings).

### Milestone 1 — GCS Terraform tab

In `docs/data-uploads/buckets/integrate/gcs.mdx`, replace the `<Tab title="Terraform">` element (lines 356–450 at `98c357f`, from the opening `<Tab title="Terraform">` through its closing `</Tab>` just before `</Tabs>`) with:

      <Tab title="Terraform">
        Terraform provisions the GCP-side resources only — registering the bucket in
        Miru still happens in the dashboard: Miru does not publish a Terraform module
        or provider today, so the config below is yours to paste and own.

        <Steps>
          <Step title="Write and apply the Terraform config">
            Save the config below in a new or existing workspace; it assumes an
            authenticated `google` provider is already configured. It creates
            everything the Console tab does: the bucket, the uploader service
            account and its object-creation grant, the Workload Identity Pool with
            an AWS provider pinned to Miru's role, and the impersonation binding.
            The variables up top are the only knobs; the outputs at the bottom
            match the Miru registration form field for field.

            ```hcl
            variable "project_id" {
              description = "GCP project ID that owns the bucket and the Workload Identity Federation resources"
              type        = string
            }

            variable "bucket_name" {
              description = "Name of the GCS bucket device uploads land in"
              type        = string
            }

            variable "bucket_location" {
              description = "Location of the bucket (region or multi-region)"
              type        = string
              default     = "US"
            }

            variable "pool_id" {
              description = "ID for the Workload Identity Pool"
              type        = string
              default     = "miru-pool"
            }

            variable "provider_id" {
              description = "ID for the pool's AWS provider"
              type        = string
              default     = "miru-provider"
            }

            variable "service_account_id" {
              description = "Account ID for the uploader service account"
              type        = string
              default     = "miru-uploader"
            }

            variable "miru_aws_account_id" {
              description = "Miru's AWS account ID (provided by Miru)"
              type        = string
            }

            variable "miru_aws_role_arn" {
              description = "Assumed-role ARN of Miru's integration identity (provided by Miru)"
              type        = string
            }

            resource "google_storage_bucket" "uploads" {
              name                        = var.bucket_name
              project                     = var.project_id
              location                    = var.bucket_location
              uniform_bucket_level_access = true
            }

            resource "google_service_account" "miru_uploader" {
              project      = var.project_id
              account_id   = var.service_account_id
              display_name = "Miru uploader"
            }

            # Create-only: the uploader can write objects but never read, list,
            # or delete your data.
            resource "google_storage_bucket_iam_member" "miru_uploader_creator" {
              bucket = google_storage_bucket.uploads.name
              role   = "roles/storage.objectCreator"
              member = "serviceAccount:${google_service_account.miru_uploader.email}"
            }

            resource "google_iam_workload_identity_pool" "miru" {
              project                   = var.project_id
              workload_identity_pool_id = var.pool_id
            }

            # attribute.aws_role strips the per-session suffix from the
            # assumed-role ARN; the attribute condition accepts only Miru's
            # exact role.
            resource "google_iam_workload_identity_pool_provider" "miru_aws" {
              project                            = var.project_id
              workload_identity_pool_id          = google_iam_workload_identity_pool.miru.workload_identity_pool_id
              workload_identity_pool_provider_id = var.provider_id

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

            # The four fields the Miru registration form asks for, 1:1.
            output "name" {
              description = "Bucket name — the form's `name` field"
              value       = google_storage_bucket.uploads.name
            }

            output "project_id" {
              description = "GCP project ID — the form's `project_id` field"
              value       = var.project_id
            }

            output "wip_provider" {
              description = "Workload Identity Pool provider resource name — the form's `wip_provider` field"
              value       = google_iam_workload_identity_pool_provider.miru_aws.name
            }

            output "service_account_email" {
              description = "Uploader service account email — the form's `service_account_email` field"
              value       = google_service_account.miru_uploader.email
            }
            ```

            Set the required variables in `terraform.tfvars`:

            ```hcl
            project_id          = "<project-id>"
            bucket_name         = "<bucket-name>"
            miru_aws_account_id = "<miru-aws-account-id>"
            miru_aws_role_arn   = "<miru-aws-role-arn>"
            ```

            Then run `terraform init && terraform apply`. When it completes,
            `terraform output` prints exactly the four values the registration
            form asks for in the next step.

            <Note>
              [Verification](/data-uploads/buckets/manage#verify-a-bucket) writes a
              zero-byte probe object under `.miru/probe/` and then attempts to delete it.
              Cleanup is best-effort — verification succeeds without it — but if you want
              the probe objects removed, also grant a delete-capable role scoped to the
              `.miru/probe/` prefix with an IAM condition.
            </Note>
          </Step>

          <RegisterGcs />

          <VerifyGcs />
        </Steps>
      </Tab>

What changed versus the current tab, for the reviewer: intro states the no-module/no-provider position; the HCL becomes fully variable-driven (six new variables replacing the `<…>` placeholders, name-type knobs defaulted to the values the other tabs use); two outputs (`name`, `project_id`) are added so all four form fields come from `terraform output`, each output carrying a description tying it to its form field; a create-only comment justifies `objectCreator`; a `terraform.tfvars` example replaces "replace the placeholders"; the verification `<Note>` and the `<RegisterGcs />`/`<VerifyGcs />` snippets are retained unchanged.

### Milestone 2 — AWS Terraform tab

In `docs/data-uploads/buckets/integrate/aws.mdx`, replace the `<Tab title="Terraform">` element (lines 362–469 at `98c357f`) with:

      <Tab title="Terraform">
        Terraform provisions the AWS-side resources only — registering the bucket in
        Miru still happens in the dashboard: Miru does not publish a Terraform module
        or provider today, so the config below is yours to paste and own. One config
        covers both phases of the external-ID flow: apply it once to create the role,
        register the bucket to get the external ID, then apply again to tighten the
        trust policy with it.

        <Steps>
          <Step title="Write and apply the Terraform config">
            Save the config below in a new or existing workspace; it assumes an
            authenticated `aws` provider is already configured — the bucket is
            created in that provider's region, and the `region` output reports it.
            The config creates everything the Console tab does: the bucket, the
            cross-account role whose trust policy allows Miru's integration role,
            and the bucket-scoped write-only policy. The variables up top are the
            only knobs; the outputs at the bottom match the Miru registration form
            field for field.

            ```hcl
            variable "bucket_name" {
              description = "Name of the S3 bucket device uploads land in"
              type        = string
            }

            variable "role_name" {
              description = "Name for the IAM role Miru assumes"
              type        = string
              default     = "miru-uploader"
            }

            variable "miru_integration_role_arn" {
              description = "Miru's integration role ARN (provided by Miru)"
              type        = string
            }

            variable "miru_external_id" {
              description = "External ID issued by Miru at registration; leave unset until then"
              type        = string
              default     = null
            }

            resource "aws_s3_bucket" "uploads" {
              bucket = var.bucket_name
            }

            data "aws_iam_policy_document" "miru_trust" {
              statement {
                effect  = "Allow"
                actions = ["sts:AssumeRole"]

                principals {
                  type        = "AWS"
                  identifiers = [var.miru_integration_role_arn]
                }

                # Only present once miru_external_id is set. Miru issues the
                # external ID when the bucket is registered, and registration
                # needs the role ARN — so the first apply creates the role
                # without the condition, and a second apply tightens it.
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
              name               = var.role_name
              assume_role_policy = data.aws_iam_policy_document.miru_trust.json
            }

            # Write-only and bucket-scoped: PutObject covers multipart uploads,
            # and AbortMultipartUpload lets a failed upload clean up after
            # itself. Miru never needs read, list, or delete permissions on
            # your data.
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

            # The three fields the Miru registration form asks for, 1:1.
            output "name" {
              description = "Bucket name — the form's `name` field"
              value       = aws_s3_bucket.uploads.bucket
            }

            output "region" {
              description = "Bucket region — the form's `region` field"
              value       = aws_s3_bucket.uploads.region
            }

            output "role_arn" {
              description = "IAM role ARN — the form's `role_arn` field"
              value       = aws_iam_role.miru_uploader.arn
            }
            ```

            Set the required variables in `terraform.tfvars`, leaving
            `miru_external_id` out — Miru has not issued it yet, so the trust
            policy carries no `sts:ExternalId` condition on this first pass:

            ```hcl
            bucket_name               = "<bucket-name>"
            miru_integration_role_arn = "<miru-integration-role-arn>"
            ```

            Then run `terraform init && terraform apply`. When it completes,
            `terraform output` prints exactly the three values the registration
            form asks for in the next step.

            <Note>
              [Verification](/data-uploads/buckets/manage#verify-a-bucket) writes a
              zero-byte probe object under `.miru/probe/` and then attempts to delete it.
              Cleanup is best-effort — verification succeeds without it — but if you want
              the probe objects removed, also grant `s3:DeleteObject` on
              `${aws_s3_bucket.uploads.arn}/.miru/probe/*` in the write policy.
            </Note>

            {/* TODO (product): confirm whether SSE-KMS buckets need an additional kms:GenerateDataKey grant for the role. */}
          </Step>

          <RegisterAws />

          <Step title="Tighten the trust policy with the external ID">
            Add the external ID Miru issued to `terraform.tfvars`:

            ```hcl
            miru_external_id = "<external-id>"
            ```

            Run `terraform apply` again. The dynamic condition block now adds the
            `sts:ExternalId` requirement to the role's trust policy — the guard against
            the [confused-deputy](https://docs.aws.amazon.com/IAM/latest/UserGuide/confused-deputy.html)
            problem: even another tenant of the same Miru service cannot trick Miru into
            assuming your role without presenting your specific external ID.
          </Step>

          <VerifyAws />
        </Steps>
      </Tab>

What changed versus the current tab: intro adds the no-module/no-provider position; `bucket_name` and `role_name` become variables (the two `<…>` placeholders leave the HCL); the dynamic condition gains a comment stating why the second apply exists; `name` and `region` outputs are added alongside `role_arn` with per-field descriptions; the first-apply and tighten steps drive everything through `terraform.tfvars`; the verification `<Note>`, the SSE-KMS TODO comment, and the `<RegisterAws />`/`<VerifyAws />` snippets are retained unchanged. Step count is unchanged (write/apply → register → tighten → verify) — the external-ID handshake genuinely needs four steps.

### Explicit non-goals

- No published Terraform registry module and no Miru Terraform provider — neither exists today. Both tab intros state that registration stays in the dashboard, which is the honest version of the callout peers make when pointing at their provider.
- No edits to the Console/CLI tabs, the shared register/verify snippets, the Config sections, or the credential-minting sections of either page.
- No change to the granted permissions (see Decision Log).

## Concrete Steps

All commands run from the docs repo root, `/home/ben/miru/workbench5/repos/docs`, on branch `feat/data-recording`.

1. Milestone 1: edit `docs/data-uploads/buckets/integrate/gcs.mdx`, replacing the `<Tab title="Terraform">…</Tab>` element with the Milestone 1 block above (strip the four-space plan indent). Confirm the diff touches only that element:

       git diff --stat
       # expect: docs/data-uploads/buckets/integrate/gcs.mdx | ~2 blocks changed, nothing else

2. Lint:

       ./scripts/lint.sh
       # expect: exit 0

3. Commit Milestone 1:

       git add docs/data-uploads/buckets/integrate/gcs.mdx
       git commit -m "docs: refactor GCS terraform tab into variables-plus-outputs config"

4. Milestone 2: edit `docs/data-uploads/buckets/integrate/aws.mdx` the same way with the Milestone 2 block, then repeat the diff check and lint.

5. Commit Milestone 2:

       git add docs/data-uploads/buckets/integrate/aws.mdx
       git commit -m "docs: refactor AWS terraform tab into variables-plus-outputs config"

6. Validate the HCL with the real terraform binary (present at `/usr/bin/terraform`). Copy the main ```hcl block from each tab (the config file, not the tfvars examples) into scratch dirs and run fmt + validate:

       mkdir -p /tmp/claude-1000/-home-ben-miru-workbench5/f815c81a-4e6a-4a5a-852f-202b90b03075/scratchpad/tfcheck/{gcs,aws}
       # paste each config into <scratch>/tfcheck/gcs/main.tf and <scratch>/tfcheck/aws/main.tf
       terraform -chdir=<scratch>/tfcheck/gcs fmt -check -diff
       terraform -chdir=<scratch>/tfcheck/aws fmt -check -diff
       terraform -chdir=<scratch>/tfcheck/gcs init -backend=false && terraform -chdir=<scratch>/tfcheck/gcs validate
       terraform -chdir=<scratch>/tfcheck/aws init -backend=false && terraform -chdir=<scratch>/tfcheck/aws validate
       # expect: "Success! The configuration is valid." for both

   `init` downloads the hashicorp/google and hashicorp/aws providers; if the sandbox has no network, record that `fmt -check` passed and validation fell back to eyeball review against the infra fixtures, and note it in Surprises & Discoveries. If `fmt` reports alignment diffs, apply them back into the MDX code blocks before committing (amend the milestone commit if already made).

   Both drafted configs already passed `terraform fmt -check` and `terraform validate` during planning (Terraform v1.14.3, hashicorp/google 7.39.0, hashicorp/aws 6.53.0) — a failure at this step therefore indicates a transcription error between the plan and the MDX, not a defect in the drafted HCL.

7. Full preflight:

       ./scripts/preflight.sh
       # expect: all sections pass; exit 0

## Validation and Acceptance

- `./scripts/lint.sh` and `./scripts/preflight.sh` exit 0 (cspell already allowlists `tfvars`, `objectCreator`, `WIF`, `principalSet`; the drafts introduce no new vocabulary — spot-check any cspell finding against `cspell.json` and extend the word list only if a genuinely new term appears).
- `terraform validate` prints `Success! The configuration is valid.` for both extracted configs (or the documented no-network fallback).
- Output names match the registration form field lists exactly. Check mechanically:

      grep -o 'output "[a-z_]*"' docs/data-uploads/buckets/integrate/aws.mdx
      # expect: name, region, role_arn — the same identifiers listed in
      # docs/snippets/data-uploads/buckets/register-aws.mdx
      grep -o 'output "[a-z_]*"' docs/data-uploads/buckets/integrate/gcs.mdx
      # expect: name, project_id, wip_provider, service_account_email — matching register-gcs.mdx

- `git diff main...` touches only the two page files; the four shared snippets are byte-identical.
- Rendered behavior (via `mintlify dev` or the preview deploy): on each page's Terraform tab a reader sees one complete config with variables first and outputs last, a `terraform.tfvars` example, and — on AWS — a tighten step whose prose explains the second apply; the Console and CLI tabs render exactly as before.

## Idempotence and Recovery

Every step is safe to repeat. The edits are whole-element replacements: re-applying either milestone block yields the same file, and `git checkout -- <file>` (before commit) or `git revert <commit>` (after) restores the prior state. The terraform scratch dirs are throwaway; delete and recreate them freely. Lint and preflight are read-only. If preflight fails after both commits, fix forward in a follow-up commit rather than rewriting history the branch may have shared.

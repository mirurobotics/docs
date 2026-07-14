# Customer bucket setup guides for Amazon S3 and Google Cloud Storage

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add two customer-facing bucket setup guides (Amazon S3, Google Cloud Storage), register them in `docs/docs.json` navigation, extend the heading-case lint allowlist with the cloud acronyms/proper-nouns the guides use, and add spell-check words to `cspell.json`. |
| `openapi/` | read-only | Source of truth for the merged bucket API field semantics (`apis/configs/components/schemas/bucket.yaml`). The guides must be consistent with these field descriptions. |
| `backend/` | read-only | Source of truth for verification behavior (`internal/configs/bktverify/aws.go`, `gcs.go`, `probe.go`). Drives the confused-deputy / negative-probe content. |
| `infra/` | read-only | Worked WIF example (`deploy/terraform/gcp-bktverify-test.tf`) that the GCS guide translates into customer instructions. |

This plan lives in `docs/plans/backlog/` because all writes land in `docs/`. The other repos are read-only references; do not edit them.

## Purpose / Big Picture

Miru writes device uploads into a cloud storage bucket the customer owns. Today the docs have no guide explaining how a customer stands up that bucket and grants Miru access. This change ships two guides:

- **Amazon S3** — create the bucket and an IAM role Miru assumes cross-account (STS AssumeRole, no access keys), with a trust policy that pins the customer's Miru workspace id as an exact `sts:ExternalId`.
- **Google Cloud Storage** — create the bucket and set up Workload Identity Federation (WIF) so Miru federates its AWS identity into the customer's GCP project and impersonates a service account, with the federation pinned to the customer's workspace id.

Both guides emphasize the **per-tenant binding / confused-deputy mitigation** the just-shipped backend work introduced, and both explain the **create-time verification probe** and its side effects.

Observable outcome: after this change, a customer reading `docs.mirurobotics.com` under **Learn → Buckets** can follow either guide end to end, configure their cloud IAM/WIF correctly on the first try, understand why Miru requires an exact external-id/session pin, and know that the benign `AccessDenied` events they will see in their audit log are expected. In the repo, `pnpm run lint` exits 0 and CI is green.

## Progress

- [ ] (YYYY-MM-DDTHHMMZ) M1 — Extend heading-case allowlist + tests; add cspell words. Commit.
- [ ] M2 — Amazon S3 guide (`docs/learn/buckets/aws.mdx`) + docs.json Buckets group + nav entry. Commit.
- [ ] M3 — Google Cloud Storage guide (`docs/learn/buckets/gcs.mdx`) + nav entry. Commit.
- [ ] M4 — Validation: `pnpm run lint`, `pnpm run test:lint`, and the tools/lint Go checks all exit 0; resolve any cspell flags. Confirm CI CLEAN before the PR leaves draft. No commit unless validation surfaces a fix.

Use timestamps when steps complete. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

- Decision: Place both guides in a new `Buckets` group under the existing **Learn** navigation group, at files `docs/learn/buckets/aws.mdx` and `docs/learn/buckets/gcs.mdx`.
  Rationale: Learn is the home for product-object concept-plus-setup pages (Devices, Groups, Releases, Config schemas). A bucket is a workspace-scoped product object the customer configures once, so it fits alongside those. Filenames use the API `provider` discriminator values (`aws`, `gcs`) for clean, stable URLs; the human product names live in the page title. Considered Administration (sibling to API keys) but Learn better matches the concept+step-by-step format.
  Date/Author: 2026-07-13 / agents@miruml.com

- Decision: Extend the heading-case lint allowlist in `docs/tools/lint/linter/headingcase/headingcase.go` with the cloud acronyms and vendor proper nouns used in the guide titles/headings, and title the pages "Amazon S3" and "Google Cloud Storage".
  Rationale: The custom `heading-case` linter forces sentence case on the frontmatter `title:` and every body heading; non-first tokens must be all-lowercase unless allowlisted. Vendor terms (AWS, S3, GCS, IAM) and the product names "Amazon S3" / "Google Cloud Storage" are not spellable in sentence case. Extending the allowlist is the repo's established, blessed mechanism for exactly this: see `docs/plans/completed/20260429-heading-case-allowlist.md`, whose recovery section says to extend the allowlist and add a `TestCheck_Allowlist` case "when the token is unambiguously a proper noun or industry-standard acronym." AWS/S3/GCS/GCP/IAM qualify. Considered a content-only alternative (linter-clean H1 like "S3 buckets" plus a `sidebarTitle: "Amazon S3"` override, which the linter does not check) — rejected because the allowlist path is the established repo pattern, yields clean matching H1s, and `sidebarTitle` is currently unused in the repo. The trade-off is that editing `tools/lint/**` activates the `lint-custom-linter` and `test-custom-linter` CI jobs (coverage gate), which M1's tests keep green.
  Date/Author: 2026-07-13 / agents@miruml.com

- Decision: Represent Miru's own fixed AWS values (the AWS account id, and the AWS integration role ARN the customer's trust policy trusts) as placeholders in the prose (`<miru-aws-account>`, and a placeholder Principal ARN), with a callout telling the customer to copy the exact values from the Miru dashboard bucket-setup screen or from Miru support. State the fixed GCS integration role name `MiruGCSIntegrationRole` literally (it is a fixed, non-secret Miru value confirmed by `infra/deploy/terraform/gcp-bktverify-test.tf`).
  Rationale: Do not invent or hardcode a Miru account id in customer docs. A placeholder plus "get it from the dashboard/support" is accurate and safe; the role name is a stable identifier the customer must match exactly, so it is stated verbatim.
  Date/Author: 2026-07-13 / agents@miruml.com

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

Assume no prior knowledge of this repo. Key facts:

**Repo shape.** `docs/` is a Mintlify site. Navigation and site config live in `docs/docs.json`. Page content is `.mdx` under `docs/docs/` (note the nested `docs/docs/` — the content root is `docs/docs`, e.g. `docs/docs/learn/releases/create.mdx`). The plan file you are reading lives in `docs/plans/backlog/` and is not linted.

**How pages are written.** Frontmatter is a YAML block delimited by `---` with at least a `title`. Example existing page `docs/docs/learn/releases/create.mdx` starts:

    ---
    title: "Create releases"
    ---

Mintlify built-in components need no import: `<Steps>`/`<Step>`, `<Warning>`, `<Note>`, `<Tip>`, `<CodeGroup>`, `<Tabs>`/`<Tab>`, `<ParamField>`. Only custom snippet/component imports (e.g. `<Framed>` from `/snippets/components/framed.jsx`) require an `import` line. **These guides use only built-in components and import nothing**, which sidesteps the import-related lint rules entirely. Do not add a header image (that would require the `Framed` snippet and an asset URL).

**Nav registration.** In `docs/docs.json`, `navigation.products[0]` ("Documentation") has a `groups` array. The **Learn** group's `pages` array currently ends with the `Releases` sub-group. A page is referenced by its slug relative to the content root without extension (e.g. `learn/releases/create`). A sub-group is `{ "group": "<label>", "pages": [ ... ] }`. Group labels are NOT lint-checked.

**Lint rules that constrain authoring (enforced by `docs/tools/lint/` via `docs/scripts/lint.sh`):**

- `heading-case` — the frontmatter `title:` value AND every body `#`/`##`/`###` heading must be sentence-case: the first token starts uppercase then all-lowercase; every later token must be all-lowercase UNLESS it is in the allowlist (`allowlist()` in `docs/tools/lint/linter/headingcase/headingcase.go`) or looks like a version tag. Inside a heading, inline code (backticks) is masked out before the check, so an acronym wrapped in backticks in a heading is exempt — but the frontmatter `title:` is NOT masked, so titles must use only allowlisted tokens or plain sentence-case words. This is why M1 extends the allowlist. Body **prose** (non-heading text) is NOT casing-checked, so acronyms may appear freely in paragraphs.
- `no-double-dash` — literal `--` in prose is flagged; use an em dash `—`. Code fences and inline code are exempt, so CLI flags like `--request` inside fenced code blocks are fine. Do not write `--` in a sentence.
- `import-*` rules — only relevant if you add imports; these guides add none.
- `redirects` — validates the `redirects` array in `docs.json`; unaffected because we only add `pages`, not redirects.

**MDX validity (enforced by ESLint MDX, `--max-warnings=0`).** MDX parses `<...>` as JSX and `{...}` as an expression. Therefore **every angle-bracket placeholder in prose must be wrapped in inline code**: write `` `<miru-aws-account>` `` and `` `<workspace-id>` ``, never a bare `<miru-aws-account>`. Placeholders inside fenced code blocks (```` ```json ````, ```` ```bash ````) are fine as-is. Keep all JSON policies and shell/gcloud/terraform snippets in fenced code blocks so their `{` and `<` do not break MDX. The existing page `docs/docs/learn/devices/provision/provisioning-tokens.mdx` demonstrates the pattern ("replace `` `<token>` `` with the actual token").

**Merged API field semantics** (from `openapi/apis/configs/components/schemas/bucket.yaml`; the guides must stay consistent with these):

- `AwsBucketConfig.region` — AWS region of the bucket (S3 endpoint + SigV4 signing).
- `AwsBucketConfig.role_arn` — ARN of the IAM role in the customer's account that Miru assumes.
- `AwsBucketConfig.external_id` — read-only; the customer's Miru workspace id (`wsp_…`), used as `sts:ExternalId` on AssumeRole (the confused-deputy mitigation). Server-known, never client-supplied. The customer must configure their trust policy to require this exact value in an `sts:ExternalId` condition BEFORE creating the bucket.
- `GcsBucketConfig.wip_provider` — the full Workload Identity Pool provider resource name Miru uses to federate into the customer's GCP project.
- `GcsBucketConfig.service_account_email` — the GCP service account Miru impersonates to mint V4-presigned upload URLs.
- `GcsBucketConfig.session_name` — read-only; the customer's Miru workspace id (`wsp_…`), which Miru sets as the STS `RoleSessionName` when federating (the confused-deputy mitigation).
- Provider enum values are `aws` and `gcs` (`openapi/apis/configs/components/enums/bucket-provider.yaml`).

**Verification behavior** (from `backend/internal/configs/bktverify/`):

- `aws.go` — Miru assumes its own integration role, then assumes the customer's `role_arn` with `ExternalID = workspace id`, then runs the write probe. It then calls `assertExternalIDEnforced`: a **deliberate no-external-id AssumeRole that MUST be denied**. If that assume SUCCEEDS, verification fails with `ExternalIDNotEnforced`; an `AccessDenied` is the secure/expected result; any other error is treated as unreachable (fail-safe). Consequence for customers: benign `AccessDenied` AssumeRole events appear in CloudTrail during verification and are expected.
- `gcs.go` — Miru assumes `MiruGCSIntegrationRole` with `RoleSessionName = workspace id`, federates via WIF (`GetCallerIdentity` signed token whose `assertion.arn` is the full assumed-role ARN including the session name), impersonates the service account, then runs the write probe. The workspace id must be 2–64 chars of `[A-Za-z0-9_+=,.@-]` (it always is for `wsp_…` ids).
- `probe.go` — the write probe uploads a zero-byte object at key `.miru/probe/<random>` with overwrite, then best-effort deletes it (a failed delete only logs a warning). So the role/SA needs write (`s3:PutObject` / `roles/storage.objectAdmin`); delete is optional but recommended so probe objects do not accumulate, which is why both guides recommend a lifecycle rule on the `.miru/probe/` prefix.

**Worked WIF example** (`infra/deploy/terraform/gcp-bktverify-test.tf`) shows the exact GCS shape to translate for customers (they substitute THEIR workspace id for the test constant `miru-bktverify-test-workspace-id`):

- WIF pool + AWS provider with `attribute_mapping = { "google.subject" = "assertion.arn" }`.
- `attribute_condition = "assertion.arn == 'arn:aws:sts::<miru-aws-account>:assumed-role/MiruGCSIntegrationRole/<workspace-id>'"` (exact full ARN including the session name = workspace id).
- `aws { account_id = <miru-aws-account> }`.
- Service account with `roles/storage.objectAdmin` on the bucket (bucket-scoped, not project-wide).
- Impersonation binding: `roles/iam.workloadIdentityUser` granted to member `principal://iam.googleapis.com/<pool-name>/subject/arn:aws:sts::<miru-aws-account>:assumed-role/MiruGCSIntegrationRole/<workspace-id>`.
- GCS lifecycle rule deleting `.miru/probe/` objects after 1 day.

## Plan of Work

### New files

- `docs/docs/learn/buckets/aws.mdx` — title `"Amazon S3"`.
- `docs/docs/learn/buckets/gcs.mdx` — title `"Google Cloud Storage"`.

### docs.json navigation edit

In `docs/docs.json`, inside `navigation.products[0].groups`, find the group with `"group": "Learn"`. Its `pages` array currently ends with the `Releases` sub-group object. Append one new sub-group object after `Releases`:

    {
      "group": "Buckets",
      "pages": [
        "learn/buckets/aws",
        "learn/buckets/gcs"
      ]
    }

Add a comma after the closing `}` of the `Releases` sub-group so the JSON stays valid. Do not touch any other part of `docs.json` (no new redirects).

### Heading-case allowlist edit

In `docs/tools/lint/linter/headingcase/headingcase.go`, the `allowlist()` function returns a `map[string]struct{}`. Add these keys (grouped under the existing comment sections):

- Under `// Acronyms.`: `"AWS"`, `"S3"`, `"GCS"`, `"GCP"`, `"IAM"`.
- Under `// Proper nouns.`: `"Amazon"`, `"Google"`, `"Cloud"`, `"Storage"`.

These cover the titles ("Amazon S3", "Google Cloud Storage") and the headings that use `S3` / `IAM`. Other acronyms the guides use (`ARN`, `STS`, `WIF`, `SigV4`, `CloudTrail`) appear only in prose or inline code, never as bare heading tokens, so they do not need allowlisting. If validation (M4) surfaces a heading-case violation for a token not listed here, either rephrase the heading or add the token here plus a `TestCheck_Allowlist` case.

### Test edit (keeps the coverage gate green)

In `docs/tools/lint/linter/headingcase/headingcase_test.go`, add cases to the existing `TestCheck_Allowlist` table (create it if absent, mirroring `TestCheck_Headings`' `build` helper):

- `"clean AWS S3 heading"`: `"## Create the S3 bucket\n"` → 0
- `"clean IAM heading"`: `"## Create the IAM role\n"` → 0

And to `TestCheck_FrontmatterTitle` (mirror its existing cases):

- `"clean Amazon S3 title"`: `"---\ntitle: \"Amazon S3\"\n---\n"` → 0
- `"clean Google Cloud Storage title"`: `"---\ntitle: \"Google Cloud Storage\"\n---\n"` → 0

### cspell edit

In `docs/cspell.json`, append to the `words` array any tokens the guides introduce that cspell does not already know. cspell splits camelCase and checks the parts, so multi-word identifiers (`workloadIdentityPools` → workload/Identity/Pools) often pass; only genuinely unknown whole tokens need adding. Candidate tokens to add if flagged: `GCS`, `GCP`, `WIF`, `presigned`, `gcloud`, `gsutil`, `SigV4`, `ExternalId`, `RoleSessionName`, and any of `AWS`/`IAM`/`ARN`/`STS` cspell does not already resolve. Do not guess — run the linter in M4 and add exactly what it reports (see Concrete Steps).

### Content specification — `aws.mdx` (title "Amazon S3")

Frontmatter:

    ---
    title: "Amazon S3"
    description: "Create an S3 bucket and IAM role so Miru can write device uploads to your AWS account."
    ---

Intro prose (2–3 short paragraphs): Miru writes device uploads into a bucket you own. Miru accesses it by assuming an IAM role in your AWS account via STS AssumeRole — there are no long-lived access keys. Miru sends your Miru workspace id (`` `wsp_…` ``) as the `` `sts:ExternalId` `` on every AssumeRole, and your role's trust policy must require that exact value. This per-tenant binding prevents the confused-deputy problem (another tenant cannot trick Miru into using your role).

Sections (all headings sentence-case; acronyms `S3`/`IAM` are allowlisted, others stay in prose/inline-code):

1. `## Prerequisites` — an AWS account where you can create S3 buckets and IAM roles; your Miru workspace id (shown as the read-only `external_id` when you create the bucket in Miru); Miru's AWS integration role ARN and account id, which you copy from the Miru dashboard bucket-setup screen (or ask Miru support). Use a `<Note>` for the "get Miru's values from the dashboard/support" point.
2. `## Create the S3 bucket` — create a bucket in your chosen region; record the bucket name and region for the Miru config.
3. `## Create the IAM role` — create a role Miru will assume; the next two sections define its trust and permissions policies.
4. `## Configure the trust policy` — show the trust policy JSON in a ```` ```json ```` fence: `Principal` set to Miru's integration role ARN (placeholder `arn:aws:iam::<miru-aws-account>:role/<miru-connector-role>`), `Action` `sts:AssumeRole`, and a `Condition` `StringEquals` on `sts:ExternalId` equal to your workspace id (`wsp_…`). Add a `<Warning>`: the condition must be an exact `StringEquals` match on your workspace id. A wildcard/`StringLike` condition, or omitting the external-id condition, FAILS verification — Miru actively probes that the external id is enforced.
5. `## Attach a permissions policy` — show the permissions policy JSON in a fence: `s3:PutObject` on `arn:aws:s3:::<your-bucket-name>/*` (used for uploads AND the create-time verification probe, which writes a zero-byte object under `.miru/probe/`). Optionally also `s3:DeleteObject` on `arn:aws:s3:::<your-bucket-name>/.miru/probe/*` so probe objects are cleaned up. Explain in prose that delete is optional because Miru best-effort-deletes the probe object.
6. `## Expected access-denied events` — a `<Note>`: during verification Miru runs a deliberate negative probe — a no-external-id AssumeRole that MUST be denied — to confirm your trust policy enforces the external id. You will therefore see benign `AccessDenied` AssumeRole events in CloudTrail. These are expected and are not a misconfiguration.
7. `## Recommended lifecycle rule` — a `<Tip>`: add an S3 lifecycle rule that expires objects under the prefix `.miru/probe/` (e.g. after 1 day) as a safety net for any probe objects left behind.
8. `## Create the bucket in Miru` — in the Miru dashboard or Platform API, create an S3 bucket providing `region` and `role_arn`; the `external_id` (your workspace id) is shown read-only and is what Miru sends as `sts:ExternalId`. Miru verifies at creation (probe write/delete plus the external-id enforcement check), so complete the trust and permissions policies first. Link to the Platform API bucket-create endpoint if it renders (see Concrete Steps for how to confirm the slug); otherwise link to `/developers/platform-api/overview`.

### Content specification — `gcs.mdx` (title "Google Cloud Storage")

Frontmatter:

    ---
    title: "Google Cloud Storage"
    description: "Create a GCS bucket and set up Workload Identity Federation so Miru can write device uploads to your GCP project."
    ---

Intro prose: Miru writes device uploads into a GCS bucket you own. Miru accesses it via Workload Identity Federation — there are no service account keys. Miru federates its AWS identity (the role `MiruGCSIntegrationRole`) into your GCP project and impersonates a service account you designate to mint V4-presigned upload URLs. The federated session is pinned to your Miru workspace id (`` `wsp_…` ``, shown read-only as `session_name`), which is the confused-deputy mitigation: only your workspace's session can impersonate the service account.

Sections:

1. `## Prerequisites` — a GCP project where you can create buckets, service accounts, and Workload Identity pools; your Miru workspace id (shown as read-only `session_name`); Miru's AWS account id (placeholder `` `<miru-aws-account>` ``, copied from the Miru dashboard bucket-setup screen or from support) and the fixed role name `MiruGCSIntegrationRole`. Use a `<Note>` for the "get Miru's account id from the dashboard/support" point.
2. `## Create the bucket` — create a GCS bucket; recommend uniform bucket-level access.
3. `## Create a service account` — create the service account Miru impersonates and grant it `roles/storage.objectAdmin` on the bucket (bucket-scoped, not project-wide). Note this write/delete access covers the verification probe under `.miru/probe/`.
4. `## Set up workload identity federation` — create a WIF pool and an AWS provider. Present the setup as a `<CodeGroup>` with a `gcloud` tab and a `Terraform` tab, mirroring the infra fixture. The provider MUST use:
   - `attribute_mapping`: `google.subject = assertion.arn`
   - `attribute_condition`: `assertion.arn == 'arn:aws:sts::<miru-aws-account>:assumed-role/MiruGCSIntegrationRole/<workspace-id>'`
   - `aws` account id = `` `<miru-aws-account>` ``
   Add a `<Warning>`: the `attribute_condition` must pin the FULL assumed-role ARN including the session name (your workspace id). An account-only or looser condition would let other sessions federate into your project — pinning the session is the per-tenant binding.
5. `## Grant impersonation to the pinned session` — grant `roles/iam.workloadIdentityUser` on the service account to the member `principal://iam.googleapis.com/<pool-name>/subject/arn:aws:sts::<miru-aws-account>:assumed-role/MiruGCSIntegrationRole/<workspace-id>`. Present as a `<CodeGroup>` (gcloud + Terraform). Explain that this binds impersonation to exactly your workspace's federated session.
6. `## Recommended lifecycle rule` — a `<Tip>`: add a GCS lifecycle rule deleting objects under the prefix `.miru/probe/` after 1 day, as a safety net for probe objects.
7. `## Create the bucket in Miru` — in the dashboard or Platform API, create a GCS bucket providing `wip_provider` (the full provider resource name) and `service_account_email`; `session_name` (your workspace id) is shown read-only. Miru verifies at creation by federating, impersonating the service account, and running the write probe, so complete the WIF setup first. Link to the Platform API bucket-create endpoint if it renders, else `/developers/platform-api/overview`.

## Concrete Steps

One commit per milestone. All commands run from the `docs/` repo root unless otherwise stated. First-time setup: run `pnpm install --frozen-lockfile` once so the linters are available.

### M1 — Heading-case allowlist + tests + cspell words

1. Edit `docs/tools/lint/linter/headingcase/headingcase.go`: add the acronym keys (`"AWS"`, `"S3"`, `"GCS"`, `"GCP"`, `"IAM"`) and proper-noun keys (`"Amazon"`, `"Google"`, `"Cloud"`, `"Storage"`) to the map returned by `allowlist()`.
2. Edit `docs/tools/lint/linter/headingcase/headingcase_test.go`: add the `TestCheck_Allowlist` and `TestCheck_FrontmatterTitle` cases listed in Plan of Work.
3. Build and test the linter:

       cd tools/lint && go build -o lint . && go test ./...

   Expect `ok` for each package. Then verify the coverage gate and Go lint that CI runs:

       ./scripts/covgate.sh
       LINT_FIX=0 ./scripts/lint.sh

   Both must exit 0. Return to the repo root: `cd ../..`.
4. Edit `docs/cspell.json`: append the candidate words you can confirm now (`GCS`, `GCP`, `WIF`, `presigned`, `gcloud`, `gsutil`, `SigV4`, `ExternalId`, `RoleSessionName`). You will finalize this list in M4 after running cspell against the real content.
5. Commit:

       git add tools/lint/linter/headingcase/ cspell.json
       git commit -m "feat(lint): allowlist cloud acronyms for bucket guides"

### M2 — Amazon S3 guide + nav

1. Create `docs/docs/learn/buckets/aws.mdx` per the content spec. Wrap every `<...>` placeholder in prose with backticks; keep all JSON policies in ```` ```json ```` fences; never write a literal `--` in prose.
2. Edit `docs/docs.json`: append the full `Buckets` sub-group with BOTH slugs (`learn/buckets/aws` and `learn/buckets/gcs`) to the Learn group's `pages`, keeping JSON valid. The `gcs` page file is created in M3; a nav slug with no file yet does not fail any CI lint check (nothing validates slug→file resolution), and M3 lands before validation.
3. Lint just the docs pipeline to catch heading-case/MDX/spell issues early:

       pnpm run lint

   Fix anything reported (see M4 for interpreting cspell output).
4. Commit:

       git add docs/docs/learn/buckets/aws.mdx docs/docs.json
       git commit -m "docs: add Amazon S3 bucket setup guide"

### M3 — Google Cloud Storage guide

1. Create `docs/docs/learn/buckets/gcs.mdx` per the content spec (same MDX-validity rules). The `learn/buckets/gcs` nav slug was already added in M2.
2. `pnpm run lint` and fix any findings.
3. Commit:

       git add docs/docs/learn/buckets/gcs.mdx
       git commit -m "docs: add Google Cloud Storage bucket setup guide"

### M4 — Validation gate

1. Run the full docs lint pipeline (this is what CI's `lint` job runs):

       pnpm run lint

   This runs, in order: the custom MDX prose linter (heading-case, no-double-dash, import rules), ESLint MDX (`--max-warnings=0`), cspell, and `mint openapi-check`. Expect it to end with `All documentation lint checks passed.`
2. If cspell reports `Unknown word (<token>)`, add each such token to the `words` array in `docs/cspell.json`, then re-run `pnpm run lint`. Repeat until clean. Commit the cspell delta if it changed since M1:

       git add cspell.json && git commit -m "chore(lint): add bucket-guide spell-check words"

3. Run the lint smoke tests (CI's `lint` job also runs this):

       pnpm run test:lint

   Expect exit 0.
4. Because `tools/lint/**` changed, CI activates the `lint-custom-linter` and `test-custom-linter` jobs. Reproduce them locally from `tools/lint/`:

       cd tools/lint && go test ./... && ./scripts/covgate.sh && LINT_FIX=0 ./scripts/lint.sh && cd ../..

   All must exit 0.
5. Confirm the Platform API bucket endpoint slug before relying on it in a link:

       ls docs/docs/references/platform-api/2026-05-06/endpoints/ | grep -i bucket

   If a `buckets/create` (or similar) page exists, link to `/references/platform-api/2026-05-06/endpoints/...`; otherwise link both guides' final section to `/developers/platform-api/overview`. (CI does not run a broken-link check, so this is a content-quality step, not a hard gate.)
6. Push the branch and confirm CI is green (the CLEAN gate). The PR must not leave draft until preflight reports CI green on the pushed branch head.

## Validation and Acceptance

Because this repo is documentation (prose and config), "tests" means the docs CI checks — there is no runtime to exercise. The acceptance gate is exactly what `.github/workflows/ci.yml` runs, reproduced locally:

1. `pnpm run lint` exits 0 and prints `All documentation lint checks passed.` — proves heading-case, no-double-dash, ESLint MDX, cspell, and `mint openapi-check` all pass on the two new pages.
2. `pnpm run test:lint` exits 0.
3. From `tools/lint/`: `go test ./...`, `./scripts/covgate.sh`, and `LINT_FIX=0 ./scripts/lint.sh` all exit 0 — proves the allowlist change keeps the custom-linter jobs green.
4. The two new pages render in navigation under **Learn → Buckets** with titles "Amazon S3" and "Google Cloud Storage". Locally, `cd docs && mint dev` and confirm both pages load and the sidebar shows the Buckets group. (Optional local check; not part of CI.)
5. Content acceptance (human review against the sources embedded in Context and Orientation):
   - AWS guide states the trust policy must require an exact `StringEquals` `sts:ExternalId` = workspace id, warns that a wildcard/loose/absent condition fails verification, grants `s3:PutObject` (and optional `.miru/probe/`-scoped `s3:DeleteObject`), documents the expected CloudTrail `AccessDenied` AssumeRole events from the negative probe, and recommends the `.miru/probe/` lifecycle rule.
   - GCS guide sets the WIF provider `attribute_condition` to the full ARN `arn:aws:sts::<miru-aws-account>:assumed-role/MiruGCSIntegrationRole/<workspace-id>`, warns that the condition must pin the session (workspace id), grants `roles/storage.objectAdmin` on the bucket, adds the `roles/iam.workloadIdentityUser` impersonation binding on the `principal://…/subject/<full-ARN>` member, and recommends the `.miru/probe/` lifecycle rule.
   - Both guides use placeholders for Miru's AWS account id (no invented value) and describe `external_id` / `session_name` as read-only workspace-id fields consistent with `bucket.yaml`.
6. CI is green on the pushed branch head (CLEAN gate) before the PR leaves draft.

Concrete failing-then-passing check for the linter change: before M1, `printf '## Create the S3 bucket\n' > /tmp/h.mdx && ./tools/lint/lint /tmp/h.mdx` reports a heading-case violation on `S3`; after M1 it reports none.

## Idempotence and Recovery

- All edits are additive and idempotent. Re-running any milestone's edits yields the same tree; re-running `pnpm run lint`, `go test`, and `covgate.sh` is safe and repeatable.
- The `docs.json` nav edit is a single sub-group insertion. If JSON becomes invalid, `pnpm run lint`'s `mint openapi-check` will still run but the site build/`mint dev` will fail to parse — re-check that a comma follows the `Releases` sub-group and the new object is well-formed.
- The heading-case allowlist edit is idempotent (map keys are unique). If M4 surfaces a heading-case violation for a token not yet allowlisted, either rephrase the heading (preferred for non-acronyms) or add the token to `allowlist()` plus a `TestCheck_Allowlist` case, then re-run the tools/lint checks.
- cspell additions are idempotent — only append tokens not already present.
- If a milestone commit needs redoing and the branch has not been shared, amend; otherwise add a `fixup:` commit. Do not force-push after review has started.

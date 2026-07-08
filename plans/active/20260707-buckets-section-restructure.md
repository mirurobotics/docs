# Restructure the Buckets docs group and fix the defining-releases rename fallout

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (this repo, branch `feat/data-recording`) | read-write | Fix inbound links broken by commit 6999f5f's renames/deletion, restructure the Buckets nav group under `docs/data-uploads/buckets/`, add a new manage page, and update `docs/docs.json` (nav + one redirect) and the `tools/lint/` heading-case allowlist |

This plan lives in `docs/plans/` because every change is a docs-repo change. All paths below are relative to the repo root `/home/ben/miru/workbench5/repos/docs` unless prefixed otherwise.

## Purpose / Big Picture

Two outcomes, delivered as three commits on `feat/data-recording`:

1. **No broken links from the checkpoint commit.** Commit 6999f5f renamed `docs/cfg-mgmt/releasing-config-schemas.mdx` → `docs/cfg-mgmt/defining-releases.mdx` and `docs/data-uploads/releasing-upload-rules.mdx` → `docs/data-uploads/defining-releases.mdx`, and deleted `docs/data-uploads/audit/upload-history.mdx`, but left ~16 files plus `docs/docs.json` pointing at the old URLs. After Milestone 1, every inbound reference resolves and the redirects lint rule passes again.
2. **A three-entry Buckets group.** Today `docs/data-uploads/buckets/overview.mdx` mixes the bucket object reference with dashboard operations, and the AWS/GCS setup guides sit as loose siblings. After Milestones 2–3, the sidebar's Buckets group reads: **Overview** (object + properties), **Manage** (register / verify / archive / delete), and an **Integrate** dropdown containing the AWS S3 and Google Cloud Storage setup guides at `data-uploads/buckets/integrate/{aws,gcs}`.

A reader of docs.mirurobotics.com can then follow any release- or upload-related link without hitting a 404, and finds bucket reference, operations, and cloud-provider setup as three clearly separated nav entries.

## Progress

- [x] Milestone 1: Fix defining-releases rename fallout (links, redirect, upload-history repointing)
- [ ] Milestone 2: Move AWS/GCS guides into an Integrate subgroup
- [ ] Milestone 3: Add the buckets manage page and finish validation

## Surprises & Discoveries

- `docs/references/cli/release-create.mdx` had three stale links, not two as the inventory said: `/data-uploads/releasing-upload-rules` appeared twice (intro sentence and the closing "see also" sentence) plus one `/cfg-mgmt/releasing-config-schemas`. All three updated; the "[releasing upload rules]" display text became "[defining releases]".

## Decision Log

- Decision: Keep the `gcs` slug (`data-uploads/buckets/integrate/gcs`) even though the request said "gcp".
  Rationale: `gcs` matches the bucket `provider` enum value and every existing link; only page titles change.
  Date/Author: 2026-07-07, planning.
- Decision: No new redirects for any of the moves/renames in this plan.
  Rationale: `git ls-tree -r main --name-only | grep -E 'data-uploads|releasing|buckets'` returns nothing — none of these URLs ever shipped on `main`. The only redirect work is repointing the existing `/learn/releases/create` destination, which currently targets the deleted `/cfg-mgmt/releasing-config-schemas` and fails the redirects lint rule.
  Date/Author: 2026-07-07, planning.
- Decision: Retitle the provider pages to "AWS S3" and "Google Cloud Storage", and add `Cloud` and `Storage` to the heading-case allowlist in `tools/lint/linter/headingcase/headingcase.go`.
  Rationale: the Go linter enforces sentence case on titles with a case-sensitive token allowlist; `AWS` and `S3` are already allowlisted but `Cloud`/`Storage` are not. Both tokens virtually always appear capitalized as parts of proper nouns (Google Cloud, Cloud Storage), so the allowlist addition is low-risk.
  Date/Author: 2026-07-07, planning.
- Decision: Inside the Integrate subgroup, keep GCS before AWS.
  Rationale: preserves the existing nav order; no reason to churn it.
  Date/Author: 2026-07-07, planning.
- Decision: Do not reintroduce a provider-chooser callout on the bucket overview page.
  Rationale: commit 6999f5f deliberately removed the `<Info>` callout that pointed at the GCS/AWS pages. The overview's `config` ParamField keeps its links to the integrate pages, and the manage page's Register section links to them for provider-specific setup — that is enough.
  Date/Author: 2026-07-07, planning.
- Decision: Repoint links to the deleted `audit/upload-history` page at `/data-uploads/uploads#viewing-uploads`, and fold the deleted page's filter list (device / rule / collection / status) into that section.
  Rationale: `docs/data-uploads/uploads.mdx` is where the upload ledger now lives (properties, completion detection, object metadata). Its "Viewing uploads" section currently *links out* to the deleted page, so it must absorb the one still-useful piece — how to browse and narrow the ledger.
  Date/Author: 2026-07-07, planning.
- Decision: Move the "buckets are immutable — no edit" sentence from the overview into an "Edit a bucket" section on the manage page.
  Rationale: matches the sibling pattern in `docs/data-uploads/upload-rules/manage.mdx` ("Editing and deleting rules"); the manage page owns the operations story, and every overview property already carries an `<ImmutableBadge />`.
  Date/Author: 2026-07-07, planning.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

This repo is the Mintlify docs site. Content lives under `docs/`; the sidebar navigation and the `redirects` array live in `docs/docs.json`. Pages are MDX with a sentence-case `title` frontmatter. Internal links are absolute site paths without `.mdx` (e.g. `/data-uploads/buckets/overview`); anchors are the kebab-cased heading text (e.g. `## Verify a bucket` → `#verify-a-bucket`). Component imports use absolute `/snippets/...` paths, so moving a page deeper does not break its imports.

Tooling (all run from the repo root):

- `./scripts/lint.sh` — builds and runs the Go prose linter (`tools/lint/`), then ESLint (MDX, `--max-warnings=0`), cspell, and `mint openapi-check`. Two Go rules matter here: **headingcase** (strict sentence case on titles/headings; case-sensitive token allowlist in `tools/lint/linter/headingcase/headingcase.go`) and **redirects** (every `redirects[].destination` in `docs/docs.json` must resolve to a real page — this rule currently fails on `main`'s `/learn/releases/create` entry because its destination was renamed away).
- `./scripts/preflight.sh` — lint smoke tests (`pnpm run test:lint`), Go lint + coverage gate for `tools/lint/`, `./scripts/lint.sh`, `pnpm audit`, and bats tests. Must exit 0 before publishing.
- There is **no automated internal-link or orphan checker** — the grep checks in Validation cover those manually.

State after the checkpoint commit 6999f5f (current HEAD):

- `docs/cfg-mgmt/defining-releases.mdx` and `docs/data-uploads/defining-releases.mdx` exist (titles already "Defining releases"); the old `releasing-config-schemas` / `releasing-upload-rules` paths are gone. Anchors on the cfg-mgmt page that inbound links use still exist: `## Git commit` (`#git-commit`) and `### CUE packages` (`#cue-packages`).
- `docs/data-uploads/audit/upload-history.mdx` is deleted (its nav group too), but two pages still link to it.
- The Buckets nav group is `overview`, `gcs`, `aws` (in `docs/docs.json` under the "Data Uploads" group).
- `docs/data-uploads/buckets/overview.mdx` holds the bucket definition snippet, the Properties list (`name`, `provider`, `config`, `status`, `last_verified_at`, `verification_error`), an immutability note, a "Device data isolation" section, and four operations sections: "Register a bucket", "Verify a bucket", "Archive a bucket", "Delete a bucket" (each with a `{/* TODO: confirm required role/scope */}` comment and a screenshots TODO). The user's fresh edits on this page (reworded `name` field, removed Info callout) must be preserved verbatim where they stand.
- `docs/data-uploads/buckets/aws.mdx` and `gcs.mdx` are pure provider setup guides (config properties + click-by-click `<Steps>` + credential-minting narrative). They contain no operations content to extract — their "Register the bucket in Miru" / "Verify the bucket" steps are integral steps of the setup walkthrough and stay put. Each links to `/data-uploads/buckets/overview#properties` (stays valid) and `/data-uploads/buckets/overview#verify-a-bucket` (breaks in Milestone 3 when the Verify section moves; must be repointed to the manage page).
- Style models for the new manage page: `docs/data-uploads/upload-rules/manage.mdx` (freshest, same section; plain `##` operation sections, no header image, role/scope TODO comments) and `docs/cfg-mgmt/schemas/manage.mdx`.

Complete inventory of broken references (paths under `docs/`, line numbers as of 6999f5f — locate by grep, not line number):

| File | Old reference → new target |
|---|---|
| `docs.json` (redirects, ~line 524) | destination `/cfg-mgmt/releasing-config-schemas` → `/cfg-mgmt/defining-releases` |
| `snippets/references/cli/releases/create/flags.mdx` (2×) | `/data-uploads/releasing-upload-rules` → `/data-uploads/defining-releases` |
| `snippets/references/cli/releases/create/usage.mdx` | same |
| `cfg-mgmt/schemas/manage.mdx` | `/cfg-mgmt/releasing-config-schemas` → `/cfg-mgmt/defining-releases` |
| `cfg-mgmt/deploy/initial-deployment.mdx` | same |
| `admin/users/access-control.mdx` | same |
| `developers/ci/overview.mdx` | same |
| `developers/ci/gh-actions.mdx` | same |
| `getting-started/quick-start/create-release.mdx` | same |
| `changelog/product.mdx` (2×, one with `#cue-packages`) | same (anchor `#cue-packages` exists on the new page) |
| `references/cli/release-create.mdx` (2 links) | one of each: `/cfg-mgmt/…` and `/data-uploads/…` → `defining-releases` |
| `data-uploads/overview.mdx` | `/data-uploads/releasing-upload-rules` → `/data-uploads/defining-releases` |
| `data-uploads/upload-collections.mdx` (3×) | 2× `releasing-upload-rules` → `defining-releases`; 1× `/data-uploads/audit/upload-history` → `/data-uploads/uploads#viewing-uploads` |
| `data-uploads/uploads.mdx` | "Viewing uploads" section links to the deleted audit page — rewrite (see Plan of Work) |
| `data-uploads/defining-releases.mdx` (5×, one with `#git-commit`) | `/cfg-mgmt/releasing-config-schemas` → `/cfg-mgmt/defining-releases` (anchor `#git-commit` exists) |
| `data-uploads/upload-rules/manage.mdx` | `releasing-upload-rules` → `defining-releases` |
| `data-uploads/upload-rules/overview.mdx` (2×) | same |
| `data-uploads/upload-rules/destinations.mdx` | same |

Where a link's display text names the old page title (e.g. "[Releasing upload rules]", "[Releasing config schemas]"), update the text to "Defining releases" (or fitting prose like "defining a release") along with the path. The `changelog/product.mdx` entries use generic link text — path-only updates there.

## Plan of Work

### Milestone 1 — defining-releases rename fallout

Apply every row of the inventory table above. Two edits need more than a path swap:

1. `docs/data-uploads/uploads.mdx`, section `## Viewing uploads` (bottom of the page): it currently defers to the deleted audit page. Replace the section body so it absorbs the browse/filter guidance, e.g.:

       ## Viewing uploads  {/* TODO: confirm required role/scope */}

       Uploads are records the [Miru Agent](/developers/agent/overview) produces — you
       don't create or edit them by hand. The **Uploads** view lists every upload in
       the workspace, newest first. Narrow the list by **device**, **upload rule**,
       **upload collection**, or **status**, and open an upload to inspect its full
       [properties](#properties).

2. `docs/data-uploads/upload-collections.mdx`: the "[filtered by collection]" link points at the deleted page → `/data-uploads/uploads#viewing-uploads`.

### Milestone 2 — Integrate subgroup

1. `git mv docs/data-uploads/buckets/aws.mdx docs/data-uploads/buckets/integrate/aws.mdx` and likewise `gcs.mdx` (create the `integrate/` directory first).
2. Retitle frontmatter: `aws.mdx` `title: "AWS S3"`, `gcs.mdx` `title: "Google Cloud Storage"`. Body content is otherwise untouched (image asset URLs under `assets.mirurobotics.com/...buckets/aws/...` are external and do not move).
3. `tools/lint/linter/headingcase/headingcase.go`: add `"Cloud"` and `"Storage"` to the proper-nouns block of `allowlist()`.
4. `docs/docs.json`: replace the Buckets group's `gcs`/`aws` entries with the Integrate subgroup (target block below, minus the `manage` line until Milestone 3).
5. Update all inbound links to the moved pages — as of 6999f5f these exist only in `docs/data-uploads/buckets/overview.mdx` (the `config` ParamField and the "Register a bucket" section): `/data-uploads/buckets/gcs` → `/data-uploads/buckets/integrate/gcs`, `/data-uploads/buckets/aws` → `/data-uploads/buckets/integrate/aws`.

### Milestone 3 — manage page

1. Create `docs/data-uploads/buckets/manage.mdx` (`title: "Manage buckets"`, style of `docs/data-uploads/upload-rules/manage.mdx`: no header image, `##` sections, keep the role/scope and screenshot TODO comments). Move these sections verbatim out of `overview.mdx`: "Register a bucket", "Verify a bucket", "Archive a bucket", "Delete a bucket", plus the `{/* TODO: add screenshots once the dashboard ships */}` comment. Add a short "Edit a bucket" section carrying the immutability sentence currently sitting below overview's Properties ("Buckets are **immutable** — there is no edit operation…register a new bucket and point your upload rules at it"). The Register section keeps its links to the integrate pages (already repointed in Milestone 2) for provider-specific configuration.
2. `docs/data-uploads/buckets/overview.mdx` then retains: imports + `<BucketDef />`, `## Properties`, and `## Device data isolation` — preserving the user's fresh wording untouched. Fix the intra-page anchor in the `status` ParamField: `[verification](#verify-a-bucket)` → `[verification](/data-uploads/buckets/manage#verify-a-bucket)`.
3. Repoint the cross-page anchor in both integrate guides: `/data-uploads/buckets/overview#verify-a-bucket` → `/data-uploads/buckets/manage#verify-a-bucket` (one occurrence each, in the probe-cleanup `<Note>`). Their `/data-uploads/buckets/overview#properties` links stay. The `#device-data-isolation` link from `docs/data-uploads/upload-rules/destinations.mdx` stays valid (section remains on overview).
4. `docs/docs.json`: insert `"data-uploads/buckets/manage"` after overview, completing the target block:

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
       }

## Concrete Steps

All commands run from `/home/ben/miru/workbench5/repos/docs` on branch `feat/data-recording`. After each milestone run `./scripts/lint.sh` (expect final line `All documentation lint checks passed.`) and commit — one commit per milestone.

Milestone 1:

    # apply the inventory-table edits; locate each occurrence with:
    grep -rn 'releasing-upload-rules\|releasing-config-schemas\|audit/upload-history' docs
    # rewrite uploads.mdx "Viewing uploads" per Plan of Work
    ./scripts/lint.sh
    grep -rn 'releasing-upload-rules\|releasing-config-schemas\|audit/upload-history' docs   # expect no output
    git add -A && git commit -m "docs: repoint links to the defining-releases pages and upload ledger"

Milestone 2:

    mkdir -p docs/data-uploads/buckets/integrate
    git mv docs/data-uploads/buckets/aws.mdx docs/data-uploads/buckets/integrate/aws.mdx
    git mv docs/data-uploads/buckets/gcs.mdx docs/data-uploads/buckets/integrate/gcs.mdx
    # retitle both pages; edit docs.json Buckets group; add Cloud/Storage to the
    # headingcase allowlist; update the two link paths in buckets/overview.mdx
    ./scripts/lint.sh
    git add -A && git commit -m "docs(data-uploads): move bucket provider guides under an Integrate subgroup"

Milestone 3:

    # create docs/data-uploads/buckets/manage.mdx; trim overview.mdx; fix the three
    # verify-a-bucket anchor links; add manage to docs.json
    ./scripts/lint.sh
    pnpm run test:lint
    ./scripts/preflight.sh
    git add -A && git commit -m "docs(data-uploads): split bucket operations onto a manage page"

## Validation and Acceptance

Run from `/home/ben/miru/workbench5/repos/docs`. All must pass before the branch is published or a PR opened:

1. `./scripts/lint.sh` exits 0 (in particular the redirects rule no longer flags `/learn/releases/create`, and the retitled pages pass headingcase).
2. `pnpm run test:lint` exits 0 (Go linter smoke tests still pass with the allowlist additions).
3. `./scripts/preflight.sh` exits 0.
4. Stale references — each expects **no output**:

       grep -rn 'releasing-upload-rules\|releasing-config-schemas\|audit/upload-history' docs
       grep -rn 'data-uploads/buckets/\(aws\|gcs\)' docs | grep -v 'buckets/integrate' | grep -v 'assets.mirurobotics.com'

5. Orphan check (every data-uploads page is in the nav) — expect no output:

       cd docs && for f in $(find data-uploads -name '*.mdx'); do p="${f%.mdx}"; grep -q "\"$p\"" docs.json || echo "ORPHAN: $p"; done

6. Link check (every absolute link from the section resolves to a page) — expect no output; anything printed must be a known non-page path or fixed:

       cd docs && grep -rhoE '(\]\(|href=")/[a-z0-9/-]+' --include='*.mdx' data-uploads cfg-mgmt | sed -E 's/^(\]\(|href=")//' | sort -u | while read -r l; do [ -f ".${l}.mdx" ] || echo "CHECK: $l"; done

7. Anchor spot checks: `#git-commit` and `#cue-packages` headings exist in `docs/cfg-mgmt/defining-releases.mdx`; `## Verify a bucket` exists in `docs/data-uploads/buckets/manage.mdx`; `## Viewing uploads` exists in `docs/data-uploads/uploads.mdx`; no page other than manage still links to `overview#verify-a-bucket`.
8. Behavioral check: `cd docs && mint dev`, open http://localhost:3000 — the Buckets group renders exactly three entries (Overview, Manage buckets, Integrate ▸ Google Cloud Storage / AWS S3), and `/data-uploads/buckets/overview` shows only definition, properties, and device data isolation. Skip if no browser is available; instead confirm the Buckets group in `docs/docs.json` matches the target block in Plan of Work exactly.

## Idempotence and Recovery

Every step is a git-tracked edit on `feat/data-recording`; nothing generated, nothing outside this repo. Re-running lint/preflight/grep checks is always safe. Before a milestone's commit, `git checkout -- <path>` or `git reset --hard HEAD` restores the last good state; after, `git revert <sha>` undoes exactly one milestone. The `git mv` steps fail loudly if re-run (source gone) — skip them on retry. The docs.json edits are idempotent: the target Buckets block and the single redirect destination fully specify the end state, so re-applying them yields the same file.

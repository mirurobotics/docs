# Milestone 7: settings and admin

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `~/dev/miru/docs` (`mirurobotics/docs`, this repo) | read-write | Every edit in this plan. Branch `docs/m7-settings-admin`, created from `origin/main` at `73dd8f5` (M5 #216). Its first commit, `98441b0` (`docs(plans): move the m5 plan to completed`, signed), moves the M5 plan to `plans/completed/` and marks it merged; it touches only that plan. Do not create another branch. |
| `~/dev/miru/frontend` (`mirurobotics/frontend`) | read-only | Source of truth for the app: **`origin/prod`** (`2b3ffed6`, unchanged since the audit). `origin/main` is one commit ahead (`0dad9c8d`, #90, in-app docs links; it only repoints the `BucketForm` guide links at `/data-uploads/connect-a-bucket/*`). The checkout has Armel's uncommitted dev-mock edits; leave them alone. Read with `git show "origin/prod:<path>"`; never switch branches. |
| `~/dev/miru/backend` (`mirurobotics/backend`) | read-only | Production backend: **`v0.11.9`** (`80acefec5c5072dde54c1b67cb4423589b2c9742`), confirmed on 2026-09-29 with `curl -s https://api.mirurobotics.com/frontend/v1/version`. Read with `git show "v0.11.9:<path>"`. Never `main`. Never modify. |

Linear: project [Frontend Docs Update](https://linear.app/mirurobotics/project/frontend-docs-update-a29451146f6a) (`P-ENG-55`), milestone **M7: Settings & admin**: ENG-1405 to ENG-1412, all Backlog on 2026-09-29, titles unchanged from the audit. Each description got a "Re-verified 2026-09-29" section with the details below. The audit behind this plan is `docs-audit.md` in the Project store (items 33, 34, 38, 39, 40, 41, 43, 45).

## Purpose / Big Picture

The settings pages still describe the app before the September access-control, API key and bucket work:

- **Buckets:** the GCS guide's service account example has the wrong format. The buckets page says any bucket's "fields" can be edited, doesn't explain verification, and its "Connect a bucket" link points at itself.
- **API keys:** the docs send you to a "Secrets page" with a create dialog and a Scopes column. The page is **API Keys** now, creating a key is a page with a scope tree, rows show when a key was last used and expand to its scopes, and keys can be edited.
- **Members:** the docs describe a **Change role** dialog. Access is now a **Type** and **Roles** column plus a **Manage access** panel, and the row menu also has **Update name** and **Invite back**.
- **Invites:** the docs say you pick a **role**. You pick a **Type** (Admin or Member), members join as Viewer, and the email takes you to an invite page with **Join workspace** and **Decline invite**.
- **Workspace and profile:** the logo and avatar can be removed, and the workspace page's Editor section isn't linked.

After this milestone, each of those sections matches frontend `origin/prod` and backend `v0.11.9`.

## Order of work, and why

1. **ENG-1405**: the GCS service account example (`connect-a-bucket/gcs.mdx`). One property block, no dependencies.
2. **ENG-1406**: buckets (`concepts/buckets.mdx`). It adds the `## Verification` section that both connect guides' notes then link to.
3. **ENG-1407**: creating an API key (`admin/apikeys.mdx`). First of three edits to the same page, top to bottom.
4. **ENG-1408**: viewing API keys, the next section down.
5. **ENG-1409**: editing API keys, a new section between View and Delete, plus one line in `access-control.mdx`.
6. **ENG-1410**: members (`admin/users/manage.mdx`, `access-control.mdx`). It comes before invites because ENG-1411 links to this page for granting roles after an invitee joins.
7. **ENG-1411**: invites (`admin/users/invites.mdx`).
8. **ENG-1412**: workspace and profile (`admin/workspace.mdx`, `admin/users/profile.mdx`). Last: small, and it's where the members-menu screenshot swap for Transfer ownership lands if Armel chooses one at ENG-1410.

## Progress

- [x] ENG-1405: the GCS service account format
- [x] ENG-1406: bucket editing and verification
- [x] Dialog steps in two sentences, across the docs (no Linear issue; Armel's request at ENG-1406)
- [x] ENG-1407: API key creation for the new create page
- [x] ENG-1408: API key last used, details and scope chips
- [x] ENG-1409: editing API key name, description and scopes
- [ ] ENG-1410: the manage access sheet and member row actions (canceled)
- [ ] ENG-1411: invite fields and the invite acceptance page (canceled)
- [ ] ENG-1412: removing logo and avatar, link editor settings (canceled)
- [x] ENG-1429: sync the Platform API scope table with production scopes (added during M7)
- [x] Final: all commits on the branch, plan updated and folded into its commit

## Surprises & Discoveries

Found while re-verifying (2026-09-29):

- **Nothing in production moved since the audit.** `/frontend/v1/version` still reports `version 0.11.9`, `git_commit 80acefec`, `api_release_version v0.9.4`. Tags `v0.11.10` to `v0.11.12` exist but aren't in production; the ones touching M7 areas are #841 (user and API key actions aligned with service guards), #840 (state-gated actions) and #848 (spec v0.9.5). Frontend `origin/prod` is still `2b3ffed6`; `git log --since=2026-09-24 origin/prod` on the settings, access, invites, invite-flow and authz paths returns nothing.
- **M5 is merged** as #216 (`73dd8f5`). Its plan was moved to `plans/completed/` in this branch's first commit, `98441b0`.
- **The audit's GCS format is wrong too.** The account is derived in the backend: `m` plus the first 20 hex characters of SHA-256(workspace ID), `@<GCP bridge project>.iam.gserviceaccount.com` (`backend internal/configs/gcpbridge/account.go`). The current screenshot `buckets/add-bucket-gcp-sa.png` matches (`mf9d1f5c3ba64dabbbd74@miru-integrations-s…`). The project comes from `GCP_BRIDGE_PROJECT_ID`; the repo has only staging's (`miru-integrations-staging`).
- **The service account exists before you register.** The Buckets page prefetches the bucket setup query, which calls `EnsureSvcAcct` (`backend services/buckets/setup.go:27`, `frontend views/settings/buckets/api/prefetch.ts`). The guide says it's provisioned "when you register the bucket", which can't be right: you grant it access first.
- **Verification leaves probe objects.** `probeBucketWritable` uploads an empty `.miru/probe/<random>`, then tries to delete it and only logs a failure (`backend internal/configs/buckets/verify/probe.go`). Both guides grant write-only access (GCS `roles/storage.objectCreator`, AWS PutObject), so every create and Role ARN edit leaves one behind.
- **AWS verification also checks the external ID.** After the probe, Miru assumes the role without the external ID and fails with "the bucket's IAM role trust policy must require an sts:ExternalId condition…" if that works (`verify/aws.go`, `errors.go`). The external ID is the workspace ID (`services/buckets/setup.go:31`), so `aws.mdx`'s `wsp_123` example is right.
- **Both connect guides link "Verification" to a page with no verification section** (`gcs.mdx:109`, `aws.mdx:190`).
- **The bucket delete dialog contradicts the rule.** It says "File rules referencing it will stop working", but Delete is disabled while any file rule references the bucket ("Cannot delete buckets that have file rules", `backend authz/actions/bucket.go`). A copy issue for M9; not filed.
- **The API Keys page lost its CLI Token section.** The hero `apikeys/header:page.png` and `apikeys/page.png` still show "Secrets" and a CLI Token card. No other docs page mentions a CLI token.
- **The Platform API scope table is out of date** (`developers/platform-api/authz.mdx`, backend-team docs): it lists `devices:delete`, which doesn't exist, and has no `upload_collections` or `file_rules` sections, which do (`backend internal/authz/scopes/scopes.go`). The OpenAPI description in `references/platform-api/2025-10-21.yaml:265` repeats `devices:delete`. Our pages only link there and repeat no scope names.
- **`manage` is each scope row's parent.** Selecting it covers the row's other toggles (`frontend …/secrets/components/form/utils/tree.ts`, `coveredBy`); `read` and `write` are siblings.
- **Access changes apply immediately.** The Type and Roles menus (table and panel) save on pick, with an "Access updated" toast and no Save (`src/features/access/hooks/useWorkspaceAccessControls.ts`). Switching to Member keeps the member's roles or gives Viewer; switching to Admin clears them.
- **Invites carry `user_type` and `roles`, and can't be updated.** The frontend spec has no invite update endpoint (only create, get, revoke, resend), so the docs' `<MutableBadge />` on `role` was wrong. The dialog has no roles field, so a Member invite is sent with `viewer` (`src/lib/authz/requests.ts`, `getMemberRoles`).
- **The invite email links to sign-in or sign-up, not straight to the invite.** `…/login?inviteId=…` or `…/signup?inviteId=…` (`backend internal/orgs/email/events.go:35-37`); after signing in, the app sends you to `/invites/<id>` (`frontend views/authn/shared/utils/destination.ts:15`). The Profile page's Join / Decline still exists, so the audit's "existing users accept from their Profile" was not wrong.
- **Only API key rows open their menu on right-click** (`RowActions`). Bucket, member and invite rows use the ellipses only (`OptionsDropdown`).
- **The logo and avatar menus open on click, not hover.** Hover shows a pencil and a tooltip ("Change logo" / "Change avatar"); the existing `logo-field.png` and `avatar-field.png` show that tooltip.
- **Linear's list view lags.** `list_issues` returned the eight M7 descriptions without their new sections right after the writes; `get_issue` showed each once. Re-read with `get_issue`.

Found while executing (2026-09-29):

- **The Buckets page doesn't create the service account on load.** The setup query (which calls `EnsureSvcAcct`) is prefetched on hover of **Buckets** in the settings sidebar (`SettingsSection.tsx:57-61`) and fetched when **Add Bucket** opens (`BucketForm.tsx:46`). Either way it exists before you register a bucket, so the docs' fix holds.
- **A failed AWS external-ID check still leaves a probe object.** The probe upload runs first (`verify/aws.go:45-52`), so the object stays even when the bucket isn't saved.
- **The API key details order is Description, Created, Created by, Scopes** (`APIKeyMetadata.tsx:26-50`), not the plan's order.
- **API key authentication is cached for 30 seconds with no invalidation** (`orgs/authn/apikeys/authn.go:54,140-163`), so scope edits and deletes take effect within about 30 seconds. Armel chose to leave it out of the docs.
- **The scope table's descriptions and permissions differ.** `config_types:write` and `releases:write` descriptions say "create", but their permissions include update, so the old "create, update" cells were right. The wrong cells were `config_schemas:write` (create only) and the three deployment action scopes (each includes create).
- **`2025-10-21.yaml:265` lists `devices:delete`** because that end-of-life API version had a Delete device endpoint. Left as is.
- **Every member row shows the member's email** (`NameCell.tsx:30`), so member screenshots from staging would leak personal emails; a mock would be needed (moot after ENG-1410 was canceled).
- **Auto-review blocks Desktop renames and some `mint broken-links` runs;** retrying with approval works.

## Decision Log

Record each of Armel's answers here as it's given at its issue (the questions themselves live in the issue sections), in this form:

    - Decision: <what Armel chose>, at ENG-XXXX.
      Rationale: <his reason, or the recommendation he accepted>.
      Date/Author: <date>, Armel.

Insert a new entry only after the previous entry's closing Date line (M5 lesson).

Precedents carried over from M1–M5 (`plans/completed/20260924-m1-fix-wrong-claims.md`, `20260924-m2-devices.md`, `20260925-m3-groups.md`, `20260927-m4-config-editor.md`, `20260928-m5-releases.md`), which apply here without asking again:

- Keep the real behavior when the code and a message disagree (M1).
- Commit this plan on its own **before** the first issue, and fold later plan updates into that commit (M5).
- Keep sections short: write only what the user needs. Armel cut most drafts further (M2–M5).
- Verify every draft sentence against the code before proposing it (M4, M5).
- Changelog edits only when they fix a fact or a link; no new changelog links (M1, M2).
- Don't add caveats the app itself doesn't show unless they change what a reader does (M1).
- **Screenshot rule (M2):** a screenshot stays only if it locates something text can't point to (an icon-only button, a control inside a row, a hidden menu), shows something the reader must recognize, or is the page's one hero. Dialogs whose fields the text names don't get one. The agent recommends per image; Armel decides. He may also say "text first, images later" (M4, M5).
- Armel replaces screenshots himself, under a new name in a topic folder, and reshapes sections live from them. When he asks, update the embed's URL or remove the frame in that issue's commit. Never delete or move an image on your own.
- On M5's releases page Armel removed every dialog image and used one menu shot per action. Recommend the same here, but ask per page; it was decided for that page.
- **Overviews are concise (M5):** what changes (the new text, one line of evidence each) and a per-image table saying plainly whether Armel **adds**, **replaces** or **drops** each image, plus anything else he needs. Re-state inherited decisions in the overview, since Armel may reverse one once he sees the code (M5, ENG-1399).
- Name buttons as the app shows them when it's clear, but don't "fix" wording that already identifies the right button (M2; project scope rule).
- Phrase dialog steps as "A dialog will appear—…" (M2, M3).
- Describe panes and lists of parts as changelog-style bullets: "**Tree:** …" (M3, M4).
- Write icon-only controls as a word plus the glyph in bold parentheses: "the ellipses **(...)**" (M3).
- Name both the ellipses **(...)** and right-click where a menu opens both ways (M3, M5). In M7 that's only API key rows.
- One commit may cover two issues when Armel asks: the first issue's title as the subject, and `Refs ENG-A, ENG-B` in the body (M3).
- Armel may cancel or skip an issue once he sees the overview (M4), so keep edits out of the tree until he says go.
- Avoid apostrophes in headings that other text links to (M3). The drafts link to `#verification` and `#invite-a-member-back`, never to "Change a member's access".
- Check the frontend dev mock has data for anything Armel will screenshot (M3).
- Verify roles and limits against the backend before proposing badges and numbers (M2).
- Check CDN image names with `curl -I` (HEAD) only (M2). The new names proposed below all returned 404 on 2026-09-29.
- Re-read a file after editing if Armel has it open in the IDE (M2).
- Run each check on its own, and stop on the first failure; don't pipe `lint.sh` into `tail` in a chain that commits (M4).
- Frontend bugs found while verifying go to **M9: Frontend bugs** (M2).
- Write Linear issues one at a time and re-read each with `get_issue` before retrying (M5).

M7 decisions:

- Decision: the GCS example uses the production project (Option A), but the ID wasn't
  given, so it ships with `miru-integrations`, at ENG-1405. The "when you register the
  bucket" clause is dropped.
  Rationale: "okay then lgtm" after confirming the claims.
  Date/Author: 2026-09-29, Armel.
- Decision: Verification mentions the leftover probe object (Option A); the edit and
  delete dialog images are removed; one menu shot per action, at ENG-1406.
  Also: the archive dialog image is removed too ("delete the archive dialog"), and the
  Archive section's menu sentence follows the page pattern, although the plan kept
  Archive out of scope.
  Also: "Miru verifies the bucket again before saving" is a `<Note>`; the Connect link
  is one short line ("GCS or AWS S3 guide").
  Also: shots from a temporary buckets mock: `buckets/list.png`, `edit-menu.png`,
  `delete-menu.png`, `archive-menu.png`.
  Also: file the delete dialog copy bug (ENG-1428, M9).
  Rationale: recommendations accepted, then Armel's calls.
  Date/Author: 2026-09-29, Armel.
- Decision: dialog steps are two sentences, not an em dash ("lets go with 1": "A
  confirmation dialog will appear. Click **Archive** to confirm."), on the buckets page
  first, then across the docs in its own commit ("I want us to make it consistent
  across all").
  Reversed: the M2/M3 precedent "A dialog will appear—…".
  Date/Author: 2026-09-29, Armel.
- Decision: keep the Platform API scope link, and fix the scope table in these docs
  instead of handing it to the backend team ("we can verify with the backend what is
  correct. so we can just fix it. you decide when"), as ENG-1429 at the end of M7, at
  ENG-1407.
  Reversed: the plan's scope rule that `developers/platform-api/authz.mdx` is never
  edited.
  Also: both images replaced (`apikeys/hero.png`, `create-page.png`), plus a new
  button shot (`create-button-v2.png`, "lets call it v2"); **New API Key** stays.
  Also: API key shots come from a temporary mock with Armel as one creator, Title Case
  names, Production CI first with three scopes.
  Date/Author: 2026-09-29, Armel.
- Decision: View API keys is one intro sentence and one bullet per shown field,
  including Last used; no hover sentence, no "shown on the row"; one image
  (`apikeys/details.png`), `page.png` dropped, at ENG-1408.
  Rationale: Armel's edits.
  Date/Author: 2026-09-29, Armel.
- Decision: Edit an API key names the ellipses and right-click; the 30-second auth
  cache stays out of the docs ("leave"); Edit gets a menu shot and a page shot
  (`edit-menu.png`, `edit-page.png`, "we should add the edit screen"); Delete gets
  `delete-menu.png` and loses its dialog, at ENG-1409.
  Date/Author: 2026-09-29, Armel.
- Decision: cancel ENG-1410, ENG-1411 and ENG-1412 ("besides ENG-1429, the rest of
  these issues should be cancelled"). No edits were made for them.
  Date/Author: 2026-09-29, Armel.
- Decision: the scope tables line up with a fixed Scope column (Option: fixed widths),
  done as a `.scope-tables` wrapper and a scoped rule in `style.css` instead of HTML
  tables, at ENG-1429.
  Date/Author: 2026-09-29, Armel.

## Outcomes & Retrospective

Complete 2026-09-29. Nine commits on `docs/m7-settings-admin` over `origin/main`: the
M5 plan move, this plan, one each for ENG-1405, 1406, 1407, 1408, 1409 and 1429, and
the dialog-sentence sweep. ENG-1405 to 1409 and ENG-1429 are Done; ENG-1410, 1411 and
1412 are Canceled. `./scripts/lint.sh`, `pnpm validate` and `mint broken-links` pass
at the head. The frontend dev mock is back to Armel's baseline (same stat and
checksums).

What shipped, versus the plan:

- `connect-a-bucket/gcs.mdx`: the `m…@miru-integrations.iam.gserviceaccount.com`
  example; the account isn't "provisioned at registration". Both guides link
  `#verification`.
- `concepts/buckets.mdx`: View names the details, Connect links both guides,
  `## Verification` (probe object, AWS external ID), Edit is AWS Role ARN only with a
  verification Note. Every dialog image is gone (Archive's too, at Armel's call) and
  each action has its own menu shot.
- `admin/apikeys.mdx`: API Keys page, the Create API Key page and one-time dialog, View
  as five bullets, a new Edit section with menu and page shots, Delete with right-click
  and a menu shot. New hero.
- `access-control.mdx`: Edit an API key. The Members list is unchanged
  ("Update another member's role"), since ENG-1410 was canceled.
- `developers/platform-api/authz.mdx`: corrected `config_schemas:write` and the three
  deployment action scopes, removed `devices:delete`, added File rules and Upload
  collections, and aligned the tables (`style.css`).
- Six pages: dialog steps in two sentences (releases, groups manage and members,
  devices manage, the device move snippet, staging area).

Left open:

- Canceled: `admin/users/manage.mdx` still describes **Change role**; `invites.mdx`
  still has the `role` property; `workspace.mdx` and `profile.mdx` still skip Remove
  logo / Remove avatar and the Editor settings link. `users/members/dropdown.png` is
  still used by `manage.mdx` and `workspace.mdx`.
- ENG-1428 (M9): the delete bucket dialog's "File rules referencing it will stop
  working" line.
- Not filed: the 30-second API key auth cache with no invalidation (backend).
- The production GCP bridge project ID was never given; the example uses
  `miru-integrations`.

No longer referenced anywhere in docs (incl. `tools/`), frontend (`origin/prod` and
working tree) or backend (`v0.11.9`), and still on the CDN (HEAD 200 on 2026-09-29), so
they can be deleted:

- `public-assets/docs/v04/images/buckets/page.png`
- `public-assets/docs/v04/images/buckets/ellipsis-dropdown.png`
- `public-assets/docs/v04/images/buckets/edit-dialog.png`
- `public-assets/docs/v04/images/buckets/delete-dialog.png`
- `public-assets/docs/v04/images/buckets/archive-dialog.png`
- `public-assets/docs/v04/images/apikeys/header:page.png`
- `public-assets/docs/v04/images/apikeys/create-dialog.png`
- `public-assets/docs/v04/images/apikeys/page.png`
- `public-assets/docs/v04/images/apikeys/scopes-popup.png`
- `public-assets/docs/v04/images/apikeys/ellipses-dropdown.png`
- `public-assets/docs/v04/images/apikeys/delete-dialog.png`

(`apikeys/create-button.png`, replaced by `create-button-v2.png`, already returns 404.)

Lessons: give mocks fake emails and IDs so shots need no blurring, and snapshot the
mock file before editing so the undo is an exact restore; Armel reshapes sections
from the live page, so propose the shortest draft and expect bullets for "what a row
shows"; a milestone can gain an issue (ENG-1429) and lose three, so re-propose the PR
title at `/pr`.

Validation, run 2026-09-29 on the final head:

- 7 prints two lines, both on canceled pages: `admin/users/profile.mdx:24` ("hover over
  your avatar", ENG-1412) and `admin/users/manage.mdx:19` ("Change role", ENG-1410).
- 9 shows only frame changes Armel asked for (new, swapped and removed shots).
- 10 prints the dialog-sentence sweep's Archive lines, the Archive section's new menu
  sentence and shot (Armel's call), and `deployments:write`'s "archive" cell. No
  archiving behavior was added.
- 11 shows `authz.mdx` changed, by Armel's decision at ENG-1407 (ENG-1429).
- 12: lint, validate and broken-links pass; each issue commit's subject is its Linear
  title with `Refs`.

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX under `docs/`. Navigation and redirects are in `docs/docs.json`; reused content lives in `docs/snippets/`.

Pages M7 touches (line numbers on `98441b0`):

- `docs/data-uploads/connect-a-bucket/gcs.mdx`: the `service_account_email` property (35–41) and the verification note (108–111).
- `docs/data-uploads/connect-a-bucket/aws.mdx`: the verification note (189–192) only.
- `docs/data-uploads/concepts/buckets.mdx`: `## View a bucket` (39–47), `## Connect a bucket` (49–53), `## Edit a bucket` (55–71). Delete (73–89) is unchanged; `## Archive a bucket` (91–) is out of scope.
- `docs/admin/apikeys.mdx`: `## Create an API key` (21–33), `## View API keys` (35–49), `## Delete an API key` (51–).
- `docs/admin/users/access-control.mdx`: Admin privileges, the **API Keys** (30–32) and **Members** (48–50) lists.
- `docs/admin/users/manage.mdx`: the whole page (`## Change a member's access` 13–29, `## Suspend a member` 31–).
- `docs/admin/users/invites.mdx`: Properties (22–48), `### Send an invite` (52–64), the `## Receiving invites` intro (104–106). Revoke, Resend, Accept and Decline are unchanged.
- `docs/admin/workspace.mdx`: `## Logo` (42–58), and a new section before `## Transfer ownership` (60).
- `docs/admin/users/profile.mdx`: `## Avatar` (22–36).
- `docs/cfg-mgmt/deploy/config-editor.mdx`: `## Editor settings` (234), only if Armel picks Option A at ENG-1412.

### Scope rules (from the Project's context doc; apply to every edit)

- **Archiving is out of scope.** Don't touch "Archive a bucket" or its images, and don't add archive or unarchive text. The Delete section's "you must archive the bucket instead" stays as it is.
- **Button naming isn't wrong.** "**New API Key**", "**New bucket**", "**Revoke Invite**" and "**Resend Invite**" all identify the right button; leave them. Page names, fields and actions that don't exist or behave differently *are* issues ("Secrets page", the **Role** field, the **Change role** dialog, the Scopes column).
- **List docs only for devices.** The API key, member, bucket and invite lists get no list-controls docs (search, status filter, sort). What a row shows can be documented.
- **M1 already fixed** the user `type` / `roles` split (ENG-1378, `admin/users/overview.mdx`) and the Leave workspace rule (ENG-1379, `snippets/workspaces/leave.mdx`). Build on those; don't redo them.
- **The Platform API scope table (`developers/platform-api/authz.mdx`) belongs to the backend team.** Don't edit it. Our pages may link to it; if one repeats a scope name, it must exist in production (`backend internal/authz/scopes/scopes.go`). The drafts repeat only `manage`, which does.
- **Keep edits short, and match the M2–M5 style.**
- **Screenshots:** Armel updates screenshots himself. Never add, delete, or move an image (see the workflow rules).

### Workflow (one issue at a time)

Rules for the executing agent:

- **Explain first, then ask.** At the start of each issue, send Armel the concise overview (what changes, per-image table). Then ask the issue's first **Stop and ask Armel** question, and wait for his go-ahead before editing.
- **Judgment calls are asked in context, at the issue.** Each issue marks its open questions as **Stop and ask Armel** steps, placed before the edit that depends on them. Ask with the options and recommendation, then wait. Don't batch questions up front or answer them yourself. Record each answer in the Decision Log. If you hit a judgment call the plan doesn't list, stop and ask the same way.
- **Never delete or move images or screenshots.** Keep every `<Frame>` and image exactly where it is until Armel says otherwise. When a screenshot is outdated, tell Armel at that issue: the image path, the section, what it shows wrong, and your add/replace/drop recommendation. Each issue lists the ones already known.
- **Propose short drafts,** and re-verify each sentence against the code before proposing it.
- **Wrap new prose at 88 columns.**
- **Track status in Linear** if you have access: In Progress when you start, Done once Armel commits. Write one issue at a time and re-read it to confirm.

Steps:

1. Work through **one** issue in the order above. Explain it, stop at each **Stop and ask Armel** step, make the edits, run the checks, flag outdated screenshots, and leave the change **uncommitted**.
2. Armel previews on localhost (`pnpm dev`, served at `http://localhost:3000`) and reviews `git diff`.
3. Iterate on Armel's feedback until he approves.
4. Armel commits. Stage the named files only (never `git add .` or `-A`). The subject is the Linear issue title verbatim, with `Refs ENG-XXXX` in the body:

        git add <files listed for the issue>
        git commit -F - <<'EOF'
        <issue title>

        Refs ENG-XXXX
        EOF

   Commits are signed with Armel's SSH key. If signing fails ("agent has no identities"), he runs `ssh-add --apple-use-keychain ~/.ssh/id_ed25519` first.

5. Tick the issue in Progress and move to the next issue.

No push and no PR. When every issue is committed, Armel goes back to the Project for a final review, and the PR is opened from there as a draft into `main`, with the body listing each issue and its Linear ID. Suggested title: `docs: update buckets, api keys, members and invites` (51 characters). The issues span `data-uploads`, `apikeys` and `admin`, so the title has no scope.

### Previewing (lessons from M1–M5)

- **Hard refresh after every edit** with Cmd+Shift+R.
- **Restart the dev server after editing `docs.json` or anything in `docs/snippets/`.** No M7 issue plans to. If one does: Ctrl+C in the server's terminal, then

        cd ~/dev/miru/docs && pnpm dev

- **Replaced screenshots can look stale.** Check a new name with `curl -I` only.
- **Confirm new anchors on localhost** by clicking the heading's link icon (`#verification`, `#edit-an-api-key`, `#invite-a-member-back`, `#editor-settings` on the workspace page).

### Checks for every issue

Run from the repo root, one at a time, and stop on the first failure:

    ./scripts/lint.sh     # last line must be: All documentation lint checks passed.
    pnpm validate         # must end with: success build validation passed
    (cd docs && ../node_modules/.bin/mint broken-links)   # must end with: success no broken links found

If CSpell flags a real word, add it to `words` in `cspell.json` and include that file in the commit. Open every link you write on localhost.

### Re-verifying claims

Each issue lists the code to re-read. Before starting, re-run the version check (`curl -s https://api.mirurobotics.com/frontend/v1/version`); if production has moved, re-verify against the new commit. In `~/dev/miru/frontend`, run `git fetch origin prod`, then `git show "origin/prod:<path>"`. In `~/dev/miru/backend`, use `git show "v0.11.9:<path>"` (quote it: zsh reads `$C:` as a modifier). If the code no longer matches this plan, trust the code, correct the edit, and note it in Surprises & Discoveries. Frontend paths below are relative to `~/dev/miru/frontend`.

## Plan of Work

The drafts below were test-applied together on 2026-09-29, with every **Stop and ask Armel** at its recommended option (`./scripts/lint.sh`, `pnpm validate` and `mint broken-links` all passed), and then restored. Match on the quoted text, since earlier edits shift lines.

### ENG-1405: `docs(data-uploads): fix the gcs service account format`

File: `docs/data-uploads/connect-a-bucket/gcs.mdx`, "Miru provides you" → `service_account_email`.

**What's wrong:** the example is `wsp_1234567890@miru-uploads.iam.gserviceaccount.com`. Real accounts are `m` plus 20 hex characters at the GCP bridge project, like the dialog's `mf9d1f5c3ba64dabbbd74@miru-integrations-s…`. The audit's `miru-uploader-wsp-…@miru-uploads…` is wrong too. The property also says Miru provisions the account "when you register the bucket", but it exists as soon as you open the Buckets page, and you grant it access before registering.

**Code to re-verify:**

- `backend internal/configs/gcpbridge/account.go`: `DeriveAccountID` (`"m" + hex(sha256(workspaceID))[:20]`) and `deriveAccountEmail` (`<id>@<project>.iam.gserviceaccount.com`); the project is `GCP_BRIDGE_PROJECT_ID` (`internal/configs/env/env.go:20`).
- `backend internal/configs/services/buckets/setup.go:27`: `EnsureSvcAcct` in the setup call; `src/views/settings/buckets/api/prefetch.ts:12` prefetches it with the Buckets page.
- `src/views/settings/buckets/components/actions/create/components/BucketForm.tsx`: the **GCS** tab's **Service account** field (copy button).

**Stop and ask Armel:** "What's the production GCP bridge project? The docs example needs it, and the repo only has staging's (`miru-integrations-staging`); your screenshot cuts it off at `miru-integrations-s…`."

- Option A (recommended): Armel gives the production project ID (`gcloud config get-value project` in the bridge project, or the production `GCP_BRIDGE_PROJECT_ID`), and the example uses it.
- Option B: use a neutral `miru-integrations` in the example. Readers copy the real value from the dialog anyway.

Screenshots: `buckets/add-bucket-gcp-sa.png` is current (it shows the real format); **keep**. Nothing to add, replace or drop.

Wait for his answer and record it.

**Edits (drafted with Option B's project; swap in Armel's):**

(a) Replace the property's first sentence pair, "The dedicated uploader service account Miru provisions for your workspace when you" / "register the bucket. You grant this account object-create access on your bucket. Shown" / "in the Miru dashboard at registration." (three lines), with:

```mdx
  The dedicated uploader service account Miru provisions for your workspace. You grant
  this account object-create access on your bucket. Shown on the **GCS** tab of the
  **Add Bucket** dialog.
```

(b) Replace "  Example: `wsp_1234567890@miru-uploads.iam.gserviceaccount.com`" with:

```mdx
  Example: `m1a2b3c4d5e6f7a8b9c0d@miru-integrations.iam.gserviceaccount.com`
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/data-uploads/connect-a-bucket/gcs#miru-provides-you`.

**Commit:** `docs/data-uploads/connect-a-bucket/gcs.mdx`, `Refs ENG-1405`.

### ENG-1406: `docs(data-uploads): document bucket editing and verification`

Files: `docs/data-uploads/concepts/buckets.mdx`; the verification notes in `docs/data-uploads/connect-a-bucket/gcs.mdx` and `aws.mdx`.

**What's wrong:**

- "Connect a bucket" links to `#connect-a-bucket` on its own page.
- "Edit a bucket" says you update "the desired fields" of a bucket's authentication configuration. Only an AWS bucket's **Role ARN** is editable; GCS buckets have nothing to edit.
- Verification isn't explained anywhere, though both connect guides link to this page for it.
- "View a bucket" doesn't say that clicking a bucket shows its details.

**Code to re-verify:**

- `src/views/settings/buckets/components/actions/items.ts`: **Edit** disabled for GCS with "GCS buckets have no editable fields"; **Delete** takes the backend's reason.
- `…/actions/edit/EditBucketDialog.tsx` and `components/AWSEditForm.tsx`: **Edit Bucket** ("Update the authentication configuration. Miru re-verifies access when you save."), **Name** read-only, **Role ARN**, **Save** ("Verifying...").
- `…/list/BucketItem.tsx` (the row expands) and `…/list/BucketMetadata.tsx`: **Provider**; for AWS **Region**, **Role ARN**, **External ID**; for GCS **Service account**; then **Created** and **Created by**.
- `backend internal/configs/buckets/verify/probe.go` (upload `.miru/probe/<random>`, then a delete whose failure is only logged) and `verify/aws.go` (`assertExternalIDEnforced`).

**Stop and ask Armel (1 of 2):** "Verification leaves an empty `.miru/probe/…` object in the bucket each time, because the guides grant write-only access. Say so?"

- Option A (recommended): yes, one clause in the probe bullet. People who find `.miru/probe/` objects in their bucket will look here.
- Option B: no; describe only the check.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `buckets/page.png` | View a bucket | **replace** with `buckets/details.png` (a bucket expanded) | The sidebar says "Secrets" and the rows are collapsed; the details are hidden behind the row. |
| `buckets/ellipsis-dropdown.png` | Edit, Delete (and Archive) | **keep** | Current (Edit / Archive / Delete). One shot per action (`buckets/edit-menu.png`, `buckets/delete-menu.png`) is optional; Archive keeps its embed either way. |
| `buckets/edit-dialog.png` | Edit a bucket | **drop** | A two-field dialog the text names (Role ARN, Save). |
| `buckets/delete-dialog.png` | Delete a bucket | **drop** | A confirm dialog the text names; its "File rules referencing it will stop working" line can't happen. |
| `buckets/archive-dialog.png` | Archive | untouched | Out of scope. |

Wait for both answers and record them.

**Edits (Option A):** keep every frame until Armel decides.

(a) In "View a bucket", replace "To view a bucket, navigate to the [Buckets page](https://app.mirurobotics.com/settings/buckets)." with:

```mdx
To view a bucket, navigate to the
[Buckets page](https://app.mirurobotics.com/settings/buckets) and click the bucket. Its
details show the provider; the region, Role ARN, and External ID (AWS) or the service
account (GCS); and when and by whom it was created.
```

(b) In "Connect a bucket", replace "To connect a bucket, visit the [Connect buckets](/data-uploads/concepts/buckets#connect-a-bucket) page." with:

```mdx
To connect a bucket, follow the guide for your provider:
[Google Cloud Storage](/data-uploads/connect-a-bucket/gcs) or
[AWS S3](/data-uploads/connect-a-bucket/aws).
```

(c) Insert before "## Edit a bucket  <PublisherBadge />":

```mdx
## Verification

When you connect a bucket or change its Role ARN, Miru checks its access before saving:

- It uploads an empty object under `.miru/probe/` to confirm it can write to the bucket,
  then tries to delete it. With write-only access, as the connect guides set up, the
  object stays in the bucket.
- For AWS buckets, it also confirms your role can't be assumed without the external ID.

If either check fails, the bucket isn't saved.

```

(Option B: end the first bullet at "…can write to the bucket.")

(d) In "Edit a bucket", replace "Editing a bucket allows you to update a bucket's authentication configuration. However, it does not allow you to change the bucket's name. To upload to a new bucket, you must create a new bucket connection." with:

```mdx
Only an AWS bucket's **Role ARN** can be edited. A bucket's name and region can't
change, and GCS buckets have nothing to edit. To upload to a different bucket, connect a
new one.
```

(e) Replace "Update the desired fields, then click **Save**." with:

```mdx
Update the **Role ARN** and click **Save**. Miru [verifies](#verification) the bucket
again before saving.
```

(f) In `gcs.mdx` and `aws.mdx`, replace `[Verification](/data-uploads/concepts/buckets)` with `[Verification](/data-uploads/concepts/buckets#verification)` (one occurrence in each file).

**Checks and preview:** run the checks. Preview `http://localhost:3000/data-uploads/concepts/buckets`, click both guide links and the `#verification` link, and click **Verification** in both guides' notes.

**Commit:** `docs/data-uploads/concepts/buckets.mdx`, `docs/data-uploads/connect-a-bucket/gcs.mdx`, `docs/data-uploads/connect-a-bucket/aws.mdx`, `Refs ENG-1406`.

### ENG-1407: `docs(apikeys): rewrite api key creation for the new create page`

File: `docs/admin/apikeys.mdx`, `## Create an API key`.

**What's wrong:** the section sends you to the "Secrets page" and describes a dialog with a name and scopes. The page is **API Keys** now (the URL is still `/settings/secrets`), and creating a key opens the **Create API Key** page: Name, an optional Description, and a scope tree. The secret is then shown once in a dialog. `authn.mdx:7` links there as "Miru Dashboard", which is fine.

**Code to re-verify:**

- `src/app/(app)/settings/secrets/page.tsx`: title **API Keys**, "Authenticate the Platform API and CI pipelines."; no CLI Token section on the page.
- `src/views/settings/secrets/components/create/CreateAPIKeyPage.tsx`: **Create API Key**, **Name**, **Description** (optional), **Scopes**, **Create**.
- `…/form/components/ScopeTreeField.tsx`, `ScopeRow.tsx`, `ScopeToggles.tsx` and `…/form/utils/tree.ts`: one row per resource, a toggle per action; a selected parent (`manage`) covers its row.
- `…/create/components/CopyAPIKeyDialog.tsx`: **Copy** + the key name, **Reveal secret key**, copy, "This API key won't be shown again, so copy it now.", **Done**.
- `backend internal/authz/scopes/scopes.go`: each resource's `manage` node has the others as children.

**Stop and ask Armel (1 of 2):** "The page links to the Platform API scope table, which is out of date (it lists `devices:delete`, which doesn't exist, and misses `upload_collections` and `file_rules`). It's the backend team's page. Keep the link?"

- Option A (recommended): keep the link, and Armel tells the backend team (the full list is in Surprises & Discoveries). It's still the only full scope reference.
- Option B: drop the sentence until the table is fixed.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `apikeys/header:page.png` | hero | **replace** with `apikeys/hero.png` (the current API Keys page) | Says "Secrets" and shows the removed CLI Token card. |
| `apikeys/create-dialog.png` | Create an API key | **replace** with `apikeys/create-page.png` (the Create API Key page, a few scope rows set) | It's a page now; the scope toggles are something a reader must recognize. |

Wait for both answers and record them.

**Edits (Option A):** keep the frames until Armel decides.

(a) Replace "To create an API key, navigate to the [Secrets page](https://app.mirurobotics.com/settings/secrets) in the dashboard and click **New API Key** in the top right." with:

```mdx
To create an API key, navigate to the
[API Keys page](https://app.mirurobotics.com/settings/secrets) and click **New API Key**
in the top right.
```

(b) Replace "Provide a name for the API key, select the desired scopes, then click **Create**. For a full reference of available scopes, see the Platform API [Authorization](/developers/platform-api/authz) page." with:

```mdx
The **Create API Key** page will open. Enter a name and an optional description, then
select the key's scopes: each row is a resource, and selecting `manage` covers every
other scope on its row. Click **Create**. For a full reference of available scopes, see
the Platform API [Authorization](/developers/platform-api/authz) page.
```

(Option B: end at "Click **Create**.")

(c) Replace "Once an API key has been created, it can never be retrieved again. If you lose your API key, you must create a new one." with:

```mdx
A dialog will appear with the new key. Copy it before clicking **Done**—it can never be
retrieved again. If you lose your API key, you must create a new one.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/apikeys#create-an-api-key`.

**Commit:** `docs/admin/apikeys.mdx`, `Refs ENG-1407`.

### ENG-1408: `docs(apikeys): document api key last used, details and scope chips`

File: `docs/admin/apikeys.mdx`, `## View API keys`.

**What's wrong:** "hover over the information icon in the **Scopes** column". There's no Scopes column. Each row shows when the key was last used (or **Never used**) and expands to its description, scopes (hover a chip for its description), and creation.

**Code to re-verify:**

- `src/views/settings/secrets/components/list/components/APIKeyItem.tsx`: the whole row expands; `LastUsed` ("Last used …" or "Never used").
- `…/list/components/APIKeyMetadata.tsx`: **Description** ("None"), **Scopes** (chips, tooltip with the description), **Created**, **Created by**.
- `backend internal/orgs/workers/api_key_last_used.go:18`: `FlushInterval = time.Minute`. Not in the draft (the app doesn't say so); mention it only if Armel asks.

**Stop and ask Armel, screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `apikeys/page.png` | View API keys | **drop** | The list is the hero once it's replaced (ENG-1407); this one shows "Secrets" and the CLI Token card. |
| `apikeys/scopes-popup.png` | View API keys | **replace** with `apikeys/details.png` (a key expanded, a scope tooltip showing) | The details are hidden behind the row. |

Wait for his answer and record it.

**Edits:** keep the frames until Armel decides.

(a) Replace "To view your existing API keys, navigate to the [Secrets page](https://app.mirurobotics.com/settings/secrets)." with:

```mdx
To view your existing API keys, navigate to the
[API Keys page](https://app.mirurobotics.com/settings/secrets). Each key shows when it
was last used, or **Never used**.
```

(b) Replace "To see the scopes of an API key, hover over the information icon in the **Scopes** column of the API key you want to view." with:

```mdx
Click an API key to see its description, its scopes, and when and by whom it was
created. Hover over a scope to see what it allows.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/apikeys#view-api-keys`.

**Commit:** `docs/admin/apikeys.mdx`, `Refs ENG-1408`.

### ENG-1409: `docs(apikeys): document editing api key name, description and scopes`

Files: `docs/admin/apikeys.mdx` (new `## Edit an API key`, and the Delete sentence); `docs/admin/users/access-control.mdx`, Admin privileges → API Keys.

**What's missing:** editing a key. The row menu's **Edit** opens the **Edit API Key** page, where admins change the name, description and scopes, then **Save**; the secret doesn't change. `access-control.mdx` lists only Create and Delete.

**Code to re-verify:**

- `src/views/settings/secrets/components/list/utils/items.ts`: **Edit** and **Delete**, each disabled with the backend's reason. `APIKeyItem.tsx` wraps the row in `RowActions`, so the menu also opens on right-click.
- `…/edit/EditAPIKeyPage.tsx`: **Edit API Key**, "Change the key name, description and scopes.", **Save**; `…/edit/components/ScopeDiffLabel.tsx`: added and removed scopes, or "No changes".
- `backend internal/orgs/services/api_keys/update.go:79-88`: `verifyAssignableScopes` when scopes change; no secret regeneration. `internal/authz/permissions/orgs/api_keys.go`: Update is admin-only.

**Stop and ask Armel, screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `apikeys/ellipses-dropdown.png` | Delete (and the new Edit) | **replace** with one shot per action: `apikeys/edit-menu.png` in Edit, `apikeys/delete-menu.png` in Delete | It shows **Rename** / **Delete**; the menu is Edit / Delete now. |
| `apikeys/delete-dialog.png` | Delete an API key | **drop** | A confirm dialog the text names. |
| Edit API Key page | Edit an API key | **no image** | Same fields as the create page. |

Wait for his answer and record it.

**Edits:**

(a) In `apikeys.mdx`, insert before "## Delete an API key  <AdminBadge />":

```mdx
## Edit an API key  <AdminBadge />

<PlatformUnsupportedBadge />

To edit an API key, click the ellipses **(...)** on the API key, or right-click it, and
select **Edit**.

The **Edit API Key** page will open. Change the key's name, description, or scopes, then
click **Save**. The key itself doesn't change, so anything using it keeps working with
the new scopes.

```

(b) In "Delete an API key", replace "To delete an API key, click the ellipsis (...) on the API key you want to delete and select **Delete**." with:

```mdx
To delete an API key, click the ellipses **(...)** on the API key, or right-click it,
and select **Delete**.
```

(c) In `access-control.mdx`, replace the two lines "- Create a new API key" / "- Delete an API key" with:

```mdx
- Create a new API key
- Edit an API key
- Delete an API key
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/apikeys#edit-an-api-key` and `http://localhost:3000/admin/users/access-control#admin-privileges`.

**Commit:** `docs/admin/apikeys.mdx`, `docs/admin/users/access-control.mdx`, `Refs ENG-1409`.

### ENG-1410: `docs(admin): document the manage access sheet and member row actions`

Files: `docs/admin/users/manage.mdx`; `docs/admin/users/access-control.mdx`, Admin privileges → Members.

**What's wrong:** "Change a member's access" says to click the ellipses, hit **Change role**, pick a role and **Save**. There's no Change role item and no Save. A member's type and workspace roles are menus in the table's **Type** and **Roles** columns, and **Manage access** opens a panel that adds their group roles. The row menu's **Update name** and **Invite back** aren't documented. M1 already split `type` from `roles` on the overview page (ENG-1378); this builds on it.

**Code to re-verify:**

- `src/views/settings/members/components/table/*`: columns **Name**, **Type**, **Roles**, **Joined**; `TypeCell` (`UserTypeMenu`: **Member**, **Admin**) and `RolesCell` (`WorkspaceRolesMenu`: **Viewer**, **Operator**, **Provisioner**, **Publisher**; "Full access" for admins and owners), both disabled on your own row (`guardSelfEdit`) and when `set_access` isn't allowed.
- `src/features/access/hooks/useWorkspaceAccessControls.ts`: saves on pick, "Access updated"; Member keeps roles or gets Viewer, Admin clears them.
- `src/features/access/components/user/*`: the panel, titled with the member's name: **User Type**, **Workspace Roles** (members only), **Groups** (direct assignments, each a group-roles menu with **Remove access** and an Undo toast; empty: "No direct group assignments.").
- `src/views/settings/members/components/actions/items.ts`: menu only for admins and never on your own row. Active members: **Manage access** (if allowed), **Update name**, **Transfer ownership**, **Suspend user**. Left or suspended: **Invite back** only. `…/table/hooks/useMembersTable.ts:53`: Invite back re-invites with the previous type and roles, no dialog.
- `…/actions/name/UpdateNameDialog.tsx`: first and last name, **Save**.

**Stop and ask Armel, screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `users/members/header:page2.png` | hero | **replace** with `users/members/hero.png` | Shows a single **Role** column. |
| `users/members/dropdown.png` | Change access, Suspend (and `workspace.mdx` Transfer ownership) | **replace** with one shot per action: `users/members/manage-access-menu.png`, `suspend-menu.png`, and `transfer-menu.png` for `workspace.mdx` (swapped at ENG-1412) | Shows **Change role** and the Role column. |
| Manage access panel | Change access | **no image** | The text names its three parts. `users/members/access-sheet.png` is free if Armel wants one. |

Wait for his answer and record it.

**Edits:** keep the frames until Armel decides.

(a) Delete the import line `import { ROLE_TOOLTIP } from '/snippets/components/users/roles.jsx';` (unused after (b)).

(b) Replace "To change a member's <Tooltip tip={ROLE_TOOLTIP.tip} cta={ROLE_TOOLTIP.cta} href={ROLE_TOOLTIP.href}>role</Tooltip>, first navigate to the [members](https://app.mirurobotics.com/settings/members) page." with:

```mdx
To change a member's [type](/admin/users/access-control#user-types) or workspace
[roles](/admin/users/access-control#roles), navigate to the
[members](https://app.mirurobotics.com/settings/members) page and pick a new value in
the member's **Type** or **Roles** column. Changes apply right away.
```

(c) Replace "Next, click the ellipses (...) of the member you want to change the role of, and hit **Change role**." with:

```mdx
To manage their group roles too, click the ellipses **(...)** of the member and select
**Manage access**.
```

(d) Replace "Select the new role and hit **Save**." with:

```mdx
A panel will appear with the member's **User Type**, **Workspace Roles**, and
**Groups**. Each group shows the member's roles in it; change them, or select
**Remove access** to remove the member from the group.
```

(e) Replace the Info's text "Changing the owner's role is not possible. Owners must transfer workspace ownership to another member to change their own role." with:

```mdx
  You can't change your own access or the owner's. Owners must transfer workspace
  ownership to another member to change their own access.
```

(f) Insert before "## Suspend a member  <AdminBadge />":

```mdx
## Update a member's name  <AdminBadge />

<PlatformUnsupportedBadge />

To update a member's name, click the ellipses **(...)** of the member and select
**Update name**. A dialog will appear—edit their first and last name and click
**Save**.

```

(g) In "Suspend a member", replace "Suspending a member is reversible. An admin can invite them back to the workspace at any time." with:

```mdx
Suspending a member is reversible. An admin can
[invite them back](#invite-a-member-back) at any time.
```

(h) Append at the end of the file:

```mdx

## Invite a member back  <AdminBadge />

<PlatformUnsupportedBadge />

To invite back a member who left or was suspended, click the ellipses **(...)** of the
member and select **Invite back**. Miru sends them a new
[invite](/admin/users/invites) with the type and roles they had before.
```

(i) In `access-control.mdx`, replace the two lines "- Suspend a member" / "- Update another member's role" with:

```mdx
- Change another member's type or roles
- Update another member's name
- Suspend a member
- Invite a member back
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/users/manage`, click the `#invite-a-member-back` link and the access-control links.

**Commit:** `docs/admin/users/manage.mdx`, `docs/admin/users/access-control.mdx`, `Refs ENG-1410`.

### ENG-1411: `docs(admin): update invite fields and the invite acceptance page`

File: `docs/admin/users/invites.mdx`, Properties, "Send an invite", and the "Receiving invites" intro.

**What's wrong:**

- Properties list a mutable `role` (admin or member). An invite has an immutable `user_type` and `roles`; invites can't be updated.
- "Send an invite" says to select a **role**. The field is **Type** (Admin or Member), and members join as Viewer.
- "Receiving invites" says only that new users should follow the email. The email takes everyone to sign in or sign up, then to the invite page with **Join workspace** and **Decline invite**. The Profile page's **Join** / **Decline**, which the rest of the section documents, still exists.

Revoke and Resend are correct ("**Revoke Invite**" and "**Resend Invite**" are naming only), as are Accept, Decline and the seven-day expiry.

**Code to re-verify:**

- `src/views/settings/members/components/toolbar/invite/components/*`: **Invite to your workspace**, **Email** (comma-separated), **Type** (**Admin** "Full read and write access to all resources.", **Member** "Access is limited to the roles you grant.", default Member), **Send**.
- `src/lib/authz/requests.ts` (`buildCreateInviteRequest`) and `roles.ts` (`getMemberRoles`): a Member invite gets `viewer`.
- `backend api/specs/frontend/v09.yaml` → `BaseInvite` (`user_type`, `roles`); paths `/invites`, `/invites/{id}`, `/revoke`, `/resend` only.
- `backend internal/orgs/email/events.go:35-37`, `frontend src/views/authn/shared/utils/destination.ts:15`, `src/views/authn/invite/*`: sign in or sign up, then `/invites/<id>`: **Join workspace**, **Decline invite**; **Invite expired**, **Invite declined** / **revoked** / **voided**, **You're already a member**.
- `src/features/invites/components/InboundInvitation/*`: Profile → **Invitations**, **Join** (confirm **Join workspace**), **Decline**.
- `backend internal/orgs/services/invites/create.go:41`: seven days.

**Stop and ask Armel, screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `invites/header:join-workspace.png` | hero | **keep** | It's the current invite page, which the new intro now describes. |
| `invites/send-invite-dialog.png` | Send an invite | **drop** | Shows a **Role** field; the text names the two fields. (Or replace, under a new name.) |
| `invites/pending-invites-button.png` | Revoke, Resend | **keep** | Current ("1 pending"). |
| `invites/pending-invite-dropdown.png` | Revoke, Resend | **keep** | Current (Resend invite / Revoke invite). |
| `invites/revoke-invite-dialog.png` | Revoke an invite | **drop** | A confirm dialog the text names. |
| `users/profile/invitations.png` | Accept, Decline | **keep** | Current. |
| `users/profile/accept-invite-dialog.png` | Accept an invite | **keep or drop** | Current; the text above it already covers its warning. |

Wait for his answer and record it.

**Edits:** keep the frames until Armel decides.

(a) Replace the import line `import { MutableBadge, ImmutableBadge } from '/snippets/components/field-badges.jsx';` with `import { ImmutableBadge } from '/snippets/components/field-badges.jsx';`, and delete `import { ROLE_TOOLTIP } from '/snippets/components/users/roles.jsx';` (both unused after (b) and (c)).

(b) Replace the whole `<ParamField path="role" type="enum">` block (from that line through its `</ParamField>`) with:

```mdx
<ParamField path="user_type" type="enum">
  <ImmutableBadge />

  The [type](/admin/users/access-control#user-types) of the invitee when they join the
  workspace. Must be one of the following:

  - `admin`
  - `member`

</ParamField>

<ParamField path="roles" type="[]enum">
  <ImmutableBadge />

  The workspace [roles](/admin/users/access-control#roles) of the invitee when they
  join. Invites sent from the dashboard give members the `viewer` role; admins get none.

</ParamField>
```

(c) Replace "Next, supply the email addresses of the members you want to invite as a comma-separated list, select their <Tooltip tip={ROLE_TOOLTIP.tip} cta={ROLE_TOOLTIP.cta} href={ROLE_TOOLTIP.href}>role</Tooltip>, and hit **Send**." with:

```mdx
Next, supply the email addresses of the members you want to invite as a comma-separated
list, select their [type](/admin/users/access-control#user-types), and hit **Send**.
Members join with the **Viewer** role; [change their access](/admin/users/manage) once
they join to grant more.
```

(d) Replace "This section covers invites received by users _who already have a Miru account_. If you do not yet have a Miru account, follow the instructions sent in the invitation email to join a workspace." with:

```mdx
Each invite email links to Miru, where you sign in, or sign up if you don't have an
account yet. You'll then see the invite: click **Join workspace** to accept it or
**Decline invite** to decline it. If the invite has expired or is no longer valid, the
page says so instead.

If you already have a Miru account, you can also accept or decline invites from your
[Profile page](https://app.mirurobotics.com/settings/profile), as described below.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/users/invites`, and click the type, roles, manage and Profile links.

**Commit:** `docs/admin/users/invites.mdx`, `Refs ENG-1411`.

### ENG-1412: `docs(admin): document removing logo and avatar, link editor settings`

Files: `docs/admin/workspace.mdx` (Logo, a new Editor settings section); `docs/admin/users/profile.mdx` (Avatar); `docs/cfg-mgmt/deploy/config-editor.mdx` only with Option A below.

**What's wrong:**

- Logo: "Click the current logo and select the image" skips the menu, which has **Change logo** and **Remove logo** (the audit said "Remove"). Removing isn't documented.
- Avatar: "hover over your avatar and select **Change avatar**". The menu opens on click and also has **Remove avatar**.
- The workspace page's **Editor** section (Allow offline deployments, Skip deploy confirmation, Format files on deploy) is documented in the config editor's "Editor settings", but `workspace.mdx` doesn't link there, and that section doesn't say where the settings live.

**Code to re-verify:**

- `src/views/settings/workspace/components/WorkspaceSection/LogoField.tsx`: **Change logo**, **Remove logo** (disabled without a logo; both disabled with the backend's reason for non-admins); crop dialog **Adjust your logo**. `…/hooks/useLogoRemover.ts`: "Logo removed".
- `src/views/settings/profile/components/ProfileSection/AvatarField.tsx`, `AvatarDropdown.tsx`: row **Profile picture**; **Change avatar**, **Remove avatar** (disabled without one); **Adjust your avatar**.
- `src/views/settings/workspace/WorkspaceSettings.tsx` and `components/EditorSection/*`: the three toggles, labels matching `config-editor.mdx#editor-settings`.

**Stop and ask Armel (1 of 2):** "Also say where the editor settings live, in the config editor's own section?"

- Option A (recommended): yes, one phrase linking the workspace settings page in `config-editor.mdx`'s first sentence. Readers of the editor page otherwise can't find the toggles.
- Option B: no; only `workspace.mdx` links to the editor page.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `workspaces/logo-field.png` | Logo | **replace** with `workspaces/logo-menu.png` (the menu open) | Shows the hover tooltip, not the menu with Remove. |
| `workspaces/crop-logo-dialog.png` | Logo | **keep or drop** | A dialog the text names (**Save**). |
| `users/profile/avatar-field.png` | Avatar | **replace** with `users/profile/avatar-menu.png` | Same as the logo. |
| `users/profile/crop-avatar-dialog.png` | Avatar | **keep or drop** | Same as the logo. |
| `users/members/dropdown.png` | `workspace.mdx` Transfer ownership | per ENG-1410's answer | Swap to `users/members/transfer-menu.png` here if Armel chose it. |

Wait for both answers and record them.

**Edits (Option A):** keep the frames until Armel decides.

(a) In `workspace.mdx`, replace "Click the current logo and select the image you want to use from the file system." with:

```mdx
Click the current logo, select **Change logo**, and choose an image from the file
system.
```

(b) Replace "Crop the image to your desired size, then click **Save** to set your workspace logo." with:

```mdx
Crop the image to your desired size, then click **Save** to set your workspace logo.

To remove the logo, click it and select **Remove logo**.
```

(c) Insert before "## Transfer ownership  <OwnerBadge />":

```mdx
## Editor settings  <AdminBadge />

The workspace page's **Editor** section holds the
[editor settings](/cfg-mgmt/deploy/config-editor#editor-settings), which control how the
config editor deploys.

```

(d) In `profile.mdx`, replace "To set your avatar, hover over your avatar and select **Change avatar**." with:

```mdx
To set your avatar, click your avatar and select **Change avatar**.
```

(e) Replace "Crop the image to your desired size, then click **Save** to set your avatar." with:

```mdx
Crop the image to your desired size, then click **Save** to set your avatar.

To remove your avatar, click it and select **Remove avatar**.
```

(f) Option A only: in `config-editor.mdx`, replace "Three **admin-only** workspace settings control how the editor behaves." with:

```mdx
Three **admin-only** settings on the
[workspace settings page](https://app.mirurobotics.com/settings/workspace) control how
the editor behaves.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/workspace`, `http://localhost:3000/admin/users/profile#avatar` and `http://localhost:3000/cfg-mgmt/deploy/config-editor#editor-settings`, and click the editor-settings links both ways.

**Commit:** `docs/admin/workspace.mdx`, `docs/admin/users/profile.mdx` (and `docs/cfg-mgmt/deploy/config-editor.mdx` with Option A), `Refs ENG-1412`.

## Concrete Steps

All commands run from `~/dev/miru/docs`.

0. Start state:

        git branch --show-current          # expect: docs/m7-settings-admin
        git log --oneline -2               # expect: 98441b0 docs(plans): move the m5 plan to completed, on top of 73dd8f5 (#216)
        git status --short                 # expect: only this plan (untracked)
        curl -s https://api.mirurobotics.com/frontend/v1/version   # expect git_commit 80acefec…

   Commit this plan on its own **before** ENG-1405 (M5 precedent): `git add plans/active/20260929-m7-settings-admin.md`, subject `docs: add the m7 settings and admin plan`, no `Refs`. Fold later plan updates into that commit (`git commit --fixup <sha>` then `git rebase -i --autosquash origin/main`).

1. For each issue, in the order above: explain it, re-verify, stop at each **Stop and ask Armel** step and wait, edit, run the checks one at a time, flag screenshots, stop for Armel's preview and review, iterate, then Armel commits.

2. After ENG-1412, update Progress, the Decision Log and Outcomes in this plan, then fold that update into the plan commit.

3. After the last commit:

        git log --oneline origin/main..HEAD   # expect the M5 plan move, this plan, and one commit per issue (fewer if Armel combined any)
        git diff --stat origin/main..HEAD     # expect gcs.mdx, aws.mdx, buckets.mdx, apikeys.mdx, access-control.mdx, manage.mdx, invites.mdx, workspace.mdx, profile.mdx (config-editor.mdx with ENG-1412 Option A), this plan, and the M5 plan in plans/completed/

   Then hand back to the Project for the final review and PR. Don't push from this session.

## Validation and Acceptance

1. `/data-uploads/connect-a-bucket/gcs`: the service account example has the `m…@<project>.iam.gserviceaccount.com` format, and the property doesn't say it's provisioned at registration.
2. `/data-uploads/concepts/buckets`: "View a bucket" names the details, "Connect a bucket" links both guides, `## Verification` describes the probe (and the AWS external-ID check), and "Edit a bucket" says only an AWS bucket's Role ARN is editable. Both guides' **Verification** links land on `#verification`.
3. `/admin/apikeys`: "API Keys page" everywhere, the Create API Key page and the one-time dialog, Last used / Never used, the expanded details, and `## Edit an API key`. `access-control.mdx` lists Edit an API key.
4. `/admin/users/manage`: Type and Roles columns, **Manage access** and its panel, `## Update a member's name`, `## Invite a member back`. No **Change role**. `access-control.mdx` → Members matches.
5. `/admin/users/invites`: `user_type` and `roles` (immutable), the **Type** field and Viewer default, and the email-to-invite-page flow.
6. `/admin/workspace` and `/admin/users/profile`: **Change logo** / **Remove logo**, **Change avatar** / **Remove avatar** by click, and the Editor settings link.
7. `grep -rn "Secrets page\|Change role\|information icon\|desired fields\|hover over your avatar\|Connect buckets\](" docs/admin docs/data-uploads --include='*.mdx'` prints nothing. (`cfg-mgmt/concepts/config-types.mdx` also says "desired fields"; it's not an M7 page.)
8. Every **Stop and ask Armel** answer is in the Decision Log, and every screenshot listed here was flagged with a recommendation.
9. No image or `<Frame>` was added, deleted or moved except as Armel decided: `git diff origin/main..HEAD -- docs/ | grep -E '^[-+].*(<Frame|!\[)'` shows only changes he asked for.
10. No archiving content was added: `git diff origin/main..HEAD -- docs/ | grep -i -E '^\+.*archiv'` prints nothing.
11. `developers/platform-api/authz.mdx` is unchanged: `git diff origin/main..HEAD --stat -- docs/developers` prints nothing.
12. `./scripts/lint.sh`, `pnpm validate` and `mint broken-links` pass on the final head. Each issue commit's subject is its Linear title, and its body has `Refs`.

## Idempotence and Recovery

- Every edit is a text replacement or insertion. Before reapplying, check whether the new text is already there (`grep`) and skip it if so.
- Before a commit, `git restore <file>` undoes an edit, and `git restore --staged <file>` unstages.
- To change an earlier issue's commit after later ones exist: `git commit --fixup <sha>` then `git rebase -i --autosquash origin/main`. That's safe because nothing is pushed. Never force-push.
- The frontend and backend repos are read-only. Don't touch the frontend's uncommitted dev-mock files, and don't switch either checkout's branch.

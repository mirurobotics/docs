# Milestone 1: fix wrong claims across the docs

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `~/dev/miru/docs` (`mirurobotics/docs`, this repo) | read-write | Every edit in this plan. Branch `docs/m1-fix-wrong-claims` (base `main`), already checked out, already holding the ENG-1373 commit `75fcf6d`. Do not create another branch. |
| `~/dev/miru/frontend` (`mirurobotics/frontend`) | read-only | Source of truth for the app, branch `main` (the local checkout sits on `prod`, which currently points at the same commit, `2b3ffed6`). Read with `git show origin/main:<path>`. Never modify, branch, or commit. |
| `~/dev/miru/backend` (`mirurobotics/backend`) | read-only | Source of truth for API behavior: the production build, tag **`v0.11.9`** (commit `80acefec`). Read with `git show v0.11.9:<path>` or `git grep … v0.11.9`. Do **not** use backend `main` (ahead of production) or the `production` branch (stale since March). Never modify. |

This plan lives in `plans/active/` because every write is a docs edit in this repo. It is committed with the milestone (see step 0 of Concrete Steps).

Linear: project [Frontend Docs Update](https://linear.app/mirurobotics/project/frontend-docs-update-a29451146f6a) (`P-ENG-55`), milestone **M1: Fix wrong claims**. The full audit this plan is built from is `docs-audit.md` in the Project store (items 9, 19, 21, 27, 35, 42, 44, 49). Everything needed to execute is repeated below.

## Purpose / Big Picture

Eight docs pages make claims that contradict what the app and API do in production: group names, staging, deployment ID copy, the editor's deployment alert, two stale GCS snippets, the user `role` property, the Leave workspace button, and one changelog link. After this milestone, each of those claims matches frontend `main` and backend `v0.11.9`. Each fix is small (one paragraph or one property) and lands as its own commit so Armel can review one issue at a time.

ENG-1380 (audit item 48, in-app links to old docs URLs) belongs to M1 but is a **frontend repo** change. It is handled separately and is already Done in Linear. Do not touch the frontend repo.

## Progress

- [x] ENG-1373: group name uniqueness and max depth. Committed as `3fc03ab` (amended from `75fcf6d` to include the new header image, border width, 88-column wrapping and the intro sentence rewording); 12 levels kept.
- [x] ENG-1374: staging for re-staging, stage summary and patch
- [x] ENG-1375: short and full deployment ID copy
- [x] ENG-1376: deployment alert is Acknowledge-only
- [x] ENG-1377: remove stale GCS register and verify snippets
- [x] ENG-1378: separate user type from roles in user properties
- [x] ENG-1379: when the Leave workspace button is disabled
- [x] ENG-1381: the changelog's Connect a bucket link
- [x] Final: all eight commits on the branch, plan updated and committed, hand back to the Project for final review and the PR

## Surprises & Discoveries

Found while verifying this plan against the code (2026-09-24):

- **Group depth is 12 levels, not 10.** `MaxGroupDepth = 10` limits the new group's *parent* to 10 ancestors. So the deepest allowed group has 11 ancestors, which makes 12 levels counting the top-level group. The backend's own test calls a group with 11 ancestors "the deepest legal position" (`tests/configs/services/groups/create_test.go`, "max chain length"), and a 12-group chain is where the error starts. The API error text still says "group hierarchy exceeds maximum depth of 10". Moving a group checks only the destination parent, not the moved group's own subgroups, so a move can push a subtree past 12 levels.
- **The changelog link in ENG-1381 isn't a 404 on the live site.** Mintlify 307-redirects the folder path `/data-uploads/connect-a-bucket` to its first page, `/aws`. The link is still non-canonical and sends GCS readers to the AWS page. `pnpm validate` and `mint broken-links` both pass with it, so neither check catches this class of link. It's the only link in `changelog/product.mdx` that resolves to neither a page nor a `docs.json` redirect.
- **The deployment alert screenshot shows the old buttons.** `devices/editor/deployment-alert.png` shows **Discard changes** and **Keep changes**. It stays in place; the agent flags it to Armel at ENG-1376.
- **The history context-menu screenshot still shows a single "Copy ID".** `devices/editor/context-menu.png` predates the short/full split. The agent flags it to Armel at ENG-1375.
- **The stage dialog screenshot predates the change summary.** `releases/stage/stage-dialog.png` shows only the description field. The agent flags it to Armel at ENG-1374.
- **"Leaving deletes the workspace" is only true for small workspaces.** Leaving never deletes anything directly (`internal/orgs/services/users/leave_wsp.go`). A pg_cron job, `prune-workspaces`, runs every 5 minutes and deletes workspaces with **zero active users and fewer than 5 devices** (`tools/supabase/migrations/20250915183477_root.sql`, `prune_workspaces()`). The app's leave dialog says "If you're the only member, this workspace will be deleted." Whether to add the caveat is asked at ENG-1379.

Found while executing (2026-09-24):

- **Frontend `main` moved ahead of `prod`.** By ENG-1377, `origin/main` (`0dad9c8d`, "point in-app docs links at canonical docs paths") was one commit ahead of `origin/prod` (`2b3ffed6`). That commit touches none of the files ENG-1374–1376 were verified against, so those fixes hold for prod. ENG-1377 onward was verified against `origin/prod`. Future milestones should read the frontend at `origin/prod`, not `origin/main`.
- **Members must have at least one role.** Backend `SanitizeUserAccess` (`internal/orgs/domain/users/roles.go`) clears roles for owners and admins, rejects a member with no roles (`member_must_have_roles`), and drops `viewer` when it's combined with other roles. ENG-1378's `roles` text says "Members have at least one role; admins and owners have none."
- **`miru device clone --root <dir>` needs `<dir>` to exist.** Otherwise it fails with "unable to find directory". The `device clone` reference doesn't say so. Found while staging example deployments for the ENG-1375 and ENG-1376 screenshots; not fixed in M1.
- **Screenshots replaced by Armel.** Every outdated screenshot was replaced under a new name, and the embed was updated in the issue's commit: `groups/groups-overview.png` (ENG-1373), `releases/stage/dialog-stage.png` (ENG-1374), `devices/editor/history-context-menu.png` (ENG-1375), `devices/editor/depl-alert.png` (ENG-1376). The old images were deleted from the assets host after a repo-wide check found no remaining references. So acceptance check 10 (no image changes) shows these four embed swaps, all approved.

## Decision Log

Record each of Armel's answers here as it's given at its issue (the questions themselves live in the issue sections), in this form:

    - Decision: <what Armel chose>, at ENG-XXXX.
      Rationale: <his reason, or the recommendation he accepted>.
      Date/Author: <date>, Armel.

- Decision: keep **12** levels (Option A), at ENG-1373.
  Rationale: it's the real limit; the recommendation was accepted. The header image (`groups/groups-overview.png`) and `borderWidth` (`32px 0 0 36px`) changes were folded into the same amended commit at his request.
  Date/Author: 2026-09-24, Armel.
- Decision: commit the plan on its own right after ENG-1373 (step 3, Option A), at ENG-1373.
  Rationale: so its full diff stops showing up during review. Later plan updates were folded into that commit.
  Date/Author: 2026-09-24, Armel.
- Decision: add the `id` property to the deployments overview (question 1, Option A), at ENG-1375. He later asked to cut it to two sentences (what the short form is; the CLI and API take the full ID) and drop the copy steps.
  Rationale: the recommendation was accepted; the shorter text matches the page's other properties.
  Date/Author: 2026-09-24, Armel.
- Decision: change the changelog's `DPL-XXX` to `DPL-XXXXX` (question 2, Option B), at ENG-1375.
  Rationale: match the real five-character short ID.
  Date/Author: 2026-09-24, Armel.
- Decision: don't mention the Release changed alert (Option A), at ENG-1376.
  Rationale: the recommendation was accepted.
  Date/Author: 2026-09-24, Armel.
- Decision: delete both GCS snippets (Option A), at ENG-1377.
  Rationale: the recommendation was accepted after re-verifying against backend `80acefec` and frontend `prod`.
  Date/Author: 2026-09-24, Armel.
- Decision: add both `type` and `roles` (Option A), at ENG-1378. The `roles` wording was tightened after checking the backend data model.
  Rationale: the recommendation was accepted.
  Date/Author: 2026-09-24, Armel.
- Decision: don't add the fewer-than-5-devices caveat (Option A), at ENG-1379.
  Rationale: the recommendation was accepted; the sentence matches the app's own leave dialog.
  Date/Author: 2026-09-24, Armel.
- Decision: one link to the AWS guide (Option B), at ENG-1381.
  Rationale: keeps the changelog's one-link-per-section style.
  Date/Author: 2026-09-24, Armel.

## Outcomes & Retrospective

Complete 2026-09-24. All eight issues are fixed, one commit each with its Linear title as the subject and `Refs ENG-XXXX`, plus this plan's own commit, on `docs/m1-fix-wrong-claims`. All eight are Done in Linear. `./scripts/lint.sh` and `pnpm validate` pass on the branch head. Merged as #206 (`6fb17dc`). The M1 review's two follow-ups (the Patch wording for drifted deployments and moving this plan to `plans/completed/`) landed in a separate PR.

What changed from the plan as written:

- #206 was opened ready for review and squash-merged, not opened as a draft.

- Every flagged screenshot was replaced, not just flagged (see Surprises).
- Armel previewed each issue before committing, and asked for a few extra edits. On the groups page: wrapping at 88 columns and rewording the intro's group-assignment sentence. On the deployments overview: a shorter `id` property. On the users overview: tightened `roles` wording.
- Issues were tracked in Linear as they went: In Progress when started, Done when committed.

Lessons for the next milestone:

- Explain the issue and its scope in one message. Ask its question only after Armel has read that, and wait for his go-ahead before editing.
- Read the frontend at `origin/prod` and the backend at the deployed commit (`80acefec`, v0.11.9), not `main`.
- Wrap new prose at 88 columns.
- When confirming whether an old image can be deleted, answer in one line; list where it's used only if it's still used.

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX (Markdown plus JSX components) under `docs/`. Navigation and redirects are in `docs/docs.json`. Reused content lives in `docs/snippets/`. Property lists use `<ParamField path="…" type="…">` with a field badge (`EditableBadge`, `MutableBadge`, `ImmutableBadge`, `NullableBadge`) from `/snippets/components/field-badges.jsx`.

### Scope rules (from the Project's context doc; apply to every edit)

- **Archiving is out of scope.** Armel's co-founder owns the archiving docs (devices, releases, deployments, config types, buckets). Don't add, change, or "fix" archiving content. Existing archive sections (for example "Archive a deployment" on the staging page) stay as they are.
- **Button naming isn't wrong.** The app shows buttons as a plus icon and a noun ("+ Device"). Docs may call them "New device", "Add device", "Stage Deployment" and so on. Don't rename buttons to match the UI exactly. Buttons or menu items that don't exist or behave differently *are* issues (for example Discard/Keep vs Acknowledge).
- **Multi-step UI can be described as steps.** Tabs or dialog panes can be written as sequential steps.
- **No fleet dashboard.** Don't document or mention the `/dashboard` page.
- **List docs only for devices.** Don't add search, sort, display, or tab docs for other lists.
- **Keep edits minimal and match the surrounding style.** Fix the claim; don't rewrite neighbouring sections. Larger rewrites of the same pages are scheduled in later milestones (M2–M8).
- **Screenshots:** Armel updates screenshots himself. Never add, delete, or move an image or screenshot (see the workflow rules).

### Workflow (one issue at a time)

Rules for the executing agent:

- **Judgment calls are asked in context, at the issue.** Each issue marks its open questions as **Stop and ask Armel** steps, placed before the edit that depends on them. When you reach one, ask the question with its options and recommendation, then wait for his answer before continuing. Don't batch questions up front or answer them yourself. Record each answer in the Decision Log. If you hit a judgment call the plan doesn't list, stop and ask the same way.
- **Never delete or move images or screenshots.** Armel updates screenshots himself. Keep every `<Frame>` and image exactly where it is, even when the text around it changes. When a screenshot is outdated, tell Armel at that issue (the image path, the section, and what it shows wrong) so he can replace it. Each issue lists the ones already known.

Steps:

1. The agent (a fresh Cursor session is fine) works through **one** issue in order, stopping at each **Stop and ask Armel** step. It makes the edits, runs the checks, flags any outdated screenshots, and leaves the change **uncommitted**.
2. Armel previews on localhost (`pnpm dev`, which runs `mint dev` in `docs/`, served at `http://localhost:3000`) and reviews `git diff`.
3. The agent iterates on Armel's feedback until he approves.
4. Armel commits: stage the named files only (never `git add .` or `-A`). The subject is the Linear issue title verbatim, with `Refs ENG-XXXX` in the body:

        git add <files listed for the issue>
        git commit -F - <<'EOF'
        <issue title>

        Refs ENG-XXXX
        EOF

5. Tick the issue in Progress and move to the next issue in the order below.

No push and no PR. When all eight commits are on the branch, Armel goes back to the Project for a final review, and the PR is opened from there: draft, into `main`, titled `docs: fix wrong claims across the docs`, with the body listing each issue and its Linear ID.

### Checks for every issue

Run from the repo root:

    ./scripts/lint.sh     # last line must be: All documentation lint checks passed.
    pnpm validate         # must end with: success build validation passed

`lint.sh` runs the custom Go MDX linter (sentence-case headings, no `--`, import and redirect rules), ESLint on MDX, CSpell (`cspell.json`), and the OpenAPI checks. If CSpell flags a real word, add it to `words` in `cspell.json` and include that file in the commit. `pnpm validate` does **not** catch links to non-page paths (see ENG-1381), so check any link you write by opening it on localhost.

### Re-verifying claims

Each issue lists the code to re-read. From `~/dev/miru/frontend`, run `git fetch origin main`, then `git show origin/main:<path>`. From `~/dev/miru/backend`, use `git show v0.11.9:<path>`. If the code no longer matches this plan, trust the code, correct the edit, and note it in Surprises & Discoveries.

## Plan of Work

Line numbers are from `main` at `7695aba` plus the ENG-1373 commit. Earlier edits in the same file shift later line numbers, so match on the quoted text.

### ENG-1373: `docs(groups): correct group name uniqueness and document max depth`, done

Status: committed as `3fc03ab` (amended from `75fcf6d`); 12 levels kept (see the Decision Log).

- File: `docs/concepts/groups/overview.mdx`.
- `name` property: "Must be unique within the workspace." was changed to "Must be unique among its sibling groups: subgroups of the same parent, or top-level groups." Evidence: the partial unique indexes `unique_parent_group_name_child` (`name, parent_id, workspace_id` where `parent_id IS NOT NULL`) and `unique_parent_group_name_root` (`name, workspace_id` where `parent_id IS NULL`) in `backend/tools/supabase/migrations/20250915183477_root.sql` (around lines 1340–1349).
- Structure section: added "The tree can be up to 12 levels deep, counting the top-level group." Evidence: `backend/internal/configs/db/group.go:17` (`MaxGroupDepth = 10`); `internal/configs/services/groups/create.go` (`verifyWithinMaxGroupDepth`) and `move.go` check the parent's ancestor chain; `GetAncestors` excludes the group itself; see the tests cited in Surprises.
- Preview: `http://localhost:3000/concepts/groups/overview` (Properties → `name`; Structure).

**Stop and ask Armel:** "The committed text says groups can be nested **12 levels** deep, counting the top-level group. That's the real limit: the backend lets the new group's parent have up to 10 ancestors. But the API error says 'group hierarchy exceeds maximum depth of 10'. Keep 12 or change it to 10?"

- Option A (recommended): keep **12**. It's the real behavior, so no change.
- Option B: change to **10**, to match the error text. This understates the real limit.

Wait for his answer. If B: replace the sentence with "The tree can be up to 10 levels deep, counting the top-level group.", run the checks, let him review, then amend. While `75fcf6d` is still HEAD, use `git add docs/concepts/groups/overview.mdx && git commit --amend --no-edit`. If a later commit already exists, use `git commit --fixup 75fcf6d` then `git rebase -i --autosquash main`. Record the answer in the Decision Log.

### ENG-1374: `docs(cfg-mgmt): update staging for re-staging, stage summary and patch`

File: `docs/cfg-mgmt/deploy/staging-area.mdx`. Five edits, (a) to (e).

**What's wrong:**

1. The Note in "Stage a deployment" (lines 70–74) says devices already running this release "cannot be staged from this flow." They can.
2. The stage dialog is described as asking only for a description (lines 91–92 and 168–169). It also shows a summary of the files that changed.
3. "Deployment drift" (lines 268–270) says "either ignore the patch or apply the patch to the staged deployment" without saying how. Applying it is the **Patch** row action.
4. "Review a deployment" (lines 335–338) has a third **Patch** option that is "currently unavailable, but will be coming soon!" The Review dialog only offers **Restage** and **Archive**.

**Code to re-verify (frontend `main`):**

- `src/features/releases/components/staging/__tests__/components/StageButton.test.tsx` ("keeps devices already on the release selectable", about line 122).
- `src/features/releases/components/staging/components/StageDevicePicker/components/DeviceList.tsx:39`. A device is disabled only when `device.actions.created_staged_deployment.allowed` is false. That is permission-based (backend `internal/authz/actions/device.go`, `resolveDeviceDplPerms`, the `deployments:stage` permission).
- `src/features/releases/components/editor/components/StageDialog/StageDialog.tsx` renders `DraftSummary` above the description. `src/features/editor/components/DraftSummary.tsx` (`summaryLabel`) produces "1 file modified", "2 files added, 1 file removed", or "No changes", followed by the file list.
- `src/features/releases/components/staging/utils/actions.ts`: the row menu has **Deploy**, **Archive**, **Patch** and **Review**. **Patch** is enabled for any staged or drifted deployment you're allowed to stage for. `hooks/useStagingList.ts` (`patchRoute`) opens `/releases/{release}/editor?device=…&deployment=…`, and `releases/components/editor/api/fetch.ts` (`buildPatchContext`) seeds the editor from that deployment.
- `src/features/releases/components/staging/components/ReviewDialog/components/ReviewFooter.tsx`: **Cancel**, a **Restage** button, and a chevron dropdown holding **Archive**. There's no Patch.
- Backend `internal/configs/services/deployments/create/entry.go` (`driftStagedDpls`): drift applies only to deployments already staged when another deployment is deployed. A replacement staged by patching starts as `staged`.

**Edits:**

(a) Replace the Note at lines 70–74:

```mdx
<Note>
  Devices that are already running this release can be staged too, for example to
  prepare changes to their current configuration. Devices you don't have permission
  to stage deployments for are unavailable for selection.
</Note>
```

(b) Replace lines 91–92 ("A dialog will appear asking for a description. Enter a description for the staged deployment, then click **Stage**."):

```mdx
A dialog will appear summarizing the files you changed (for example, `1 file modified`)
and asking for a description. Enter a description for the staged deployment, then click
**Stage**.
```

In "Patch a deployment", replace lines 168–169 ("A dialog will appear asking for a description. Enter a description, then click **Stage**."):

```mdx
A dialog will appear summarizing the files you changed and asking for a description.
Enter a description, then click **Stage**.
```

(c) In "Patch a deployment", after the first paragraph (ends "…allowing you to create a replacement staged deployment for that device.", about line 150), add:

```mdx
You can patch staged and [drifted](#deployment-drift) deployments. Patching a drifted
deployment is how you apply a drift's patches to it.
```

(d) In "Deployment drift", replace lines 268–270 ("When a deployment drifts, … apply the patch to the staged deployment."):

```mdx
When a deployment drifts, it is marked as `drifted` and cannot be deployed until an
explicit review is performed. You must either ignore the patch by
[restaging](#review-a-deployment) the deployment, or apply the patch by
[patching](#patch-a-deployment) the staged deployment in the release editor.
```

(e) In "Review a deployment":

- Replace line 313 ("With the deployment's patches available for review, you have two options for proceeding.") with:

```mdx
With the deployment's patches available for review, the dialog offers two options:
**Restage** and **Archive**. To apply the patches instead, close the dialog and
[patch](#patch-a-deployment) the deployment.
```

- In the **Archive** paragraph, change "The second option is to [archive](#archive-a-deployment) the deployment." to "The second option, in the dropdown next to **Restage**, is to [archive](#archive-a-deployment) the deployment." Leave the rest of the archiving text alone (archiving is out of scope).
- Delete the **Patch** block (lines 335–338: `**Patch**`, the blank line, and "The ability to apply patches … coming soon!"), keeping one blank line before "To list the deployments…".

Keep every screenshot on the page where it is.

**Tell Armel (outdated screenshot):** "`releases/stage/stage-dialog.png` in 'Stage a deployment' shows the old stage dialog with only a Description field. The dialog now shows a summary of changed files above it. You may want to replace it."

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/staging-area`, specifically `#stage-a-deployment`, `#patch-a-deployment`, `#deployment-drift` and `#review-a-deployment`. Click the three new in-page links.

**Commit:** `docs/cfg-mgmt/deploy/staging-area.mdx`, `Refs ENG-1374`.

### ENG-1375: `docs(deployments): document short and full deployment id copy`

Files: `docs/cfg-mgmt/deploy/config-editor.mdx` and `docs/concepts/deployments/overview.mdx`.

**What's wrong:** `config-editor.mdx` line 131 (History → Deployment list, right-click menu) lists "**Copy ID** or **Copy description**." The menu has **Copy short ID**, **Copy full ID** and **Copy description**. No page explains the two ID forms.

**Code to re-verify (frontend `main`):**

- `src/features/devices/components/editor/components/history/utils/actions.ts`: the menu is Redeploy, Set as diff base, Copy short ID, Copy full ID, Copy description.
- `src/features/deployments/utils/id.ts`: the short ID is `` `DPL-${id.substring(4, 9)}` ``, meaning `DPL-` plus the five characters after the `dpl_` prefix.
- `src/features/deployments/components/id/DeploymentID.tsx`: the copy icon appears on hover next to a deployment ID and opens **Copy short ID** / **Copy full ID**. It's used in the staging list, the Review dialog, the device's Deployments tab, and the deployment sheet (`git grep -n showCopy origin/main -- src`).
- `docs/snippets/references/cli/deployment/clone/args.mdx` already says the CLI takes the full `dpl_…` ID and rejects the short one.

**Stop and ask Armel (1 of 2):** "Only `config-editor.mdx` makes a wrong claim; no page explains the two ID forms. Should I also add an `id` property to the deployments overview (`concepts/deployments/overview.mdx`) that explains `dpl_…` vs `DPL-XXXXX` and how to copy each?"

- Option A (recommended): yes, fix the menu labels (edit a) and add the `id` property (edit b). The overview is the natural home for what an ID looks like, and ENG-1400 (M5) can link to it when it rewrites `device-history.mdx`.
- Option B: no. Fix only the menu labels (edit a) and skip edit (b).

Either way, leave `cfg-mgmt/audit/device-history.mdx` for ENG-1400. It says nothing about copying IDs today.

**Stop and ask Armel (2 of 2):** "The Sep 23 product changelog calls the short ID `DPL-XXX`. Real short IDs are `DPL-` plus 5 characters. Leave it, or change it to `DPL-XXXXX`?"

- Option A (recommended): leave it. It's a placeholder in a dated release note, not a format claim.
- Option B: change `DPL-XXX` to `DPL-XXXXX` in `docs/changelog/product.mdx` (about line 104, edit c).

Wait for both answers before editing, and record them in the Decision Log.

**Edits:**

(a) `config-editor.mdx` line 131: replace `- **Copy ID** or **Copy description**.` with:

```mdx
- **Copy short ID**, **Copy full ID**, or **Copy description**.
```

(b) Only if Armel chose A in question 1. In `concepts/deployments/overview.mdx`, under `## Properties`, insert as the first property (before `<ParamField path="description" …>`), followed by a blank line:

```mdx
<ParamField path="id" type="string">
  <ImmutableBadge />

  The deployment's unique identifier. The dashboard shows a short ID: `DPL-` followed by
  the five characters after the full ID's `dpl_` prefix. To copy either form, hover over a
  deployment's ID and click the copy icon, then choose **Copy short ID** or
  **Copy full ID**. Use the full ID with the CLI and the Platform API.

  Examples: `dpl_7xKqW1mZ3nBvTgL5cRs9dYuP` (full), `DPL-7xKqW` (short)
</ParamField>
```

`ImmutableBadge` is already imported on that page.

(c) Only if Armel chose B in question 2. In `docs/changelog/product.mdx`, change `` `DPL-XXX` `` to `` `DPL-XXXXX` `` in the "Short and full deployment IDs" item.

**Tell Armel (outdated screenshot):** "`devices/editor/context-menu.png` in 'Deployment list' shows the old history menu with a single **Copy ID**. The menu now has **Copy short ID** and **Copy full ID**. You may want to replace it."

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#deployment-list`; if edit (b) was made, also `http://localhost:3000/concepts/deployments/overview` (the new `id` property at the top of Properties); if edit (c) was made, also the first September 23, 2026 entry on `http://localhost:3000/changelog/product`.

**Commit:** the files edited (`config-editor.mdx`, plus `concepts/deployments/overview.mdx` and/or `changelog/product.mdx` if changed), `Refs ENG-1375`.

### ENG-1376: `docs(cfg-mgmt): fix the deployment alert to acknowledge-only behavior`

File: `docs/cfg-mgmt/deploy/config-editor.mdx`, section `## Deployment alert` (lines 206–225).

**What's wrong:** the section describes **Acknowledge** when there are no edits, and **Discard** or **Keep** (a "three-way JSON merge") when there are. The only button is **Acknowledge**. With unsaved edits, acknowledging discards them. The screenshot shows the old Discard/Keep buttons.

**Code to re-verify (frontend `main`):** `src/features/devices/components/editor/components/alert/DeploymentAlert.tsx`. The title is "New deployment detected" and there's one `AlertDialogAction`, **Acknowledge**. The description is "A new deployment was made to this device. Acknowledging will discard your unsaved changes and load the new deployment." with edits, or "… Your editor will be updated to the new deployment." without. The alert is modal, so you can't copy your edits out while it's open.

**Stop and ask Armel:** "The editor has a sibling alert, **Release changed**, shown when a new release is attached to the device while you're editing (`ReleaseChangedAlert.tsx`). It's also Acknowledge-only and discards unsaved edits. The docs don't mention it. Add one sentence about it here?"

- Option A (recommended): no. The docs never claimed anything about it, so leaving it out isn't wrong, and new editor content belongs in M4.
- Option B: yes. Append the paragraph in edit (b) below.

Wait for his answer and record it in the Decision Log.

**Edits:**

(a) Replace the paragraphs between `## Deployment alert` and the screenshot's `<Frame>` (from "If another user (or system) deploys…" through "…preserving your changes where possible.") with the text below. Keep the heading, the `<Frame>` with `deployment-alert.png` exactly as it is, and the `---` separator after it.

```mdx
If another user (or system) deploys to the same device while you have the editor
open, a **New deployment detected** alert appears with the new deployment's details.
The only option is **Acknowledge**, which loads the new deployment into the editor.

If you have unsaved edits, acknowledging discards them.
```

(b) Only if Armel chose B. After the paragraph "If you have unsaved edits, acknowledging discards them." (still above the `<Frame>`), add:

```mdx
If a new release is attached to the device instead, a **Release changed** alert appears.
It works the same way: **Acknowledge** loads the new release and discards any unsaved
edits.
```

**Tell Armel (outdated screenshot):** "`devices/editor/deployment-alert.png` in 'Deployment alert' shows the old alert with **Discard changes** and **Keep changes** buttons. The alert now has a single **Acknowledge** button, so the screenshot contradicts the new text. I left it in place for you to replace."

**Checks and preview:** run the checks, and `grep -n -i "three-way\|\*\*Keep\*\*" docs/cfg-mgmt/deploy/config-editor.mdx` should print nothing. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#deployment-alert`.

**Commit:** `docs/cfg-mgmt/deploy/config-editor.mdx`, `Refs ENG-1376`.

### ENG-1377: `docs(data-uploads): remove stale gcs register and verify snippets`

Files to delete: `docs/snippets/data-uploads/buckets/register-gcs.mdx` and `docs/snippets/data-uploads/buckets/verify-gcs.mdx`.

**What's wrong:** the snippets describe a GCS form with `project_id`, `wip_provider` and `service_account_email` fields, and a separate **Verify** action with `verified`, `invalid` and `verification_error` statuses. None of that exists. They aren't rendered anywhere, but they're stale sources someone could reuse.

**Re-verify:**

- Nothing imports them: `grep -rn "register-gcs\|verify-gcs" docs` should print nothing. They're the only files under `docs/snippets/data-uploads/`, so both directories disappear.
- Frontend `main`: `src/views/settings/buckets/components/actions/create/api/repository.ts:27` sends `gcs: {}`. The actions under `src/views/settings/buckets/components/actions/` are archive, create, delete and edit, with no verify. Verification happens automatically on create and edit.

**Stop and ask Armel:** "Both snippets embed a screenshot (`data-uploads/buckets/gcs/register-bucket.png` and `gcs/verify-bucket.png`), so deleting the files removes those two image embeds. Nothing renders either snippet, and the image files on the assets host aren't touched. OK to delete both snippets as the issue says?"

- Option A (recommended): yes, delete both files.
- Option B: keep the files and only strip their wrong text, leaving the `<Frame>` embeds. This keeps dead, unrendered files around.

Wait for his answer and record it in the Decision Log.

**Edit:** for Option A, `rm` both files. There are no other references, and no `cspell.json` words are unique to them.

**Checks and preview:** run the checks. There's no page to preview. Optionally confirm `http://localhost:3000/data-uploads/connect-a-bucket/gcs` still renders.

**Commit:** `git add docs/snippets/data-uploads/buckets/register-gcs.mdx docs/snippets/data-uploads/buckets/verify-gcs.mdx` (this stages the deletions), `Refs ENG-1377`.

### ENG-1378: `docs(admin): separate user type from roles in user properties`

File: `docs/admin/users/overview.mdx`, `## Properties`, the `role` ParamField (lines 39–49).

**What's wrong:** `role` lists `member`, `admin` and `owner`. Those are **user types**. Roles are a separate list (viewer, operator, provisioner, publisher). `admin/users/access-control.mdx` already uses the right terms (`## User types`, `## Roles`).

**Code to re-verify:**

- Frontend `main`: `src/lib/authz/types.ts` has `UserType` (owner, admin, member), `WorkspaceRole` (viewer, operator, provisioner, publisher) and `GroupRole` (manager, operator, provisioner).
- Backend `v0.11.9`: `api/specs/frontend/v09.yaml`, schema `User`, has `type: UserType` and `roles: WorkspaceRole[]`. The `WorkspaceRole` description says roles are "only meaningful when `type` is `member`." Group roles live on group membership (`GroupMember`), not on the user.

**Stop and ask Armel:** "Fixing this means renaming `role` to `type`. Should I also add a `roles` property listing the workspace roles (viewer, operator, provisioner, publisher), which the API returns alongside `type`?"

- Option A (recommended): yes, add both `type` and `roles`. The issue is about separating the two, and a reader of the property list otherwise has no idea roles exist.
- Option B: rename `role` to `type` only, and drop the second ParamField below.

Wait for his answer and record it in the Decision Log.

**Edit:** replace the whole `role` ParamField with these two, or with just `type` if Armel chose B (keep the same blank-line style, including the blank line before `</ParamField>`):

```mdx
<ParamField path="type" type="enum">
  <MutableBadge />

  The user's [type](/admin/users/access-control#user-types), which sets their broad access level in the workspace.

  Allowed values:
  - `member`
  - `admin`
  - `owner`

</ParamField>

<ParamField path="roles" type="[]enum">
  <MutableBadge />

  The user's workspace [roles](/admin/users/access-control#roles). Only members have roles; admins and owners already have full access. Roles within a group are granted through [group membership](/concepts/groups/members).

  Allowed values:
  - `viewer`
  - `operator`
  - `provisioner`
  - `publisher`

</ParamField>
```

Nothing links to the old `role` anchor (`grep -rn "users/overview#" docs` shows only the page's own `#status` link).

**Checks and preview:** run the checks. Preview `http://localhost:3000/admin/users/overview` (Properties) and click both access-control links and the group membership link.

**Commit:** `docs/admin/users/overview.mdx`, `Refs ENG-1378`.

### ENG-1379: `docs(admin): correct when the leave workspace button is disabled`

File: `docs/snippets/workspaces/leave.mdx` (last line). It renders on `admin/workspace.mdx` and `admin/users/profile.mdx`.

**What's wrong:** "If the **leave workspace** button is disabled, you are not the only member of the workspace and must first transfer ownership…" This reads as if every non-sole member is blocked. In fact the button is disabled only for the **owner** while other users are still active. Members and admins can always leave.

**Code to re-verify:**

- Frontend `main`: `src/views/settings/profile/hooks/useWorkspaceAccess.ts:45-48` (`isOwnerWithOtherMembers = userIsOwner && count > 1`, where `count` comes from `useActiveUserCount` in `src/views/settings/profile/api/query.ts` and counts `active` users only). `components/MembershipSection/MembershipSection.tsx` shows the tooltip "Transfer ownership before leaving the workspace".
- Backend `v0.11.9`: `internal/orgs/services/users/leave_wsp.go` (`VerifyLeaveWsp`, `CanOwnerLeaveWsp`) rejects an owner unless they're the only active user, with the error "unable to leave workspace: owner must transfer ownership to another member before leaving". Non-owners aren't checked.

**Stop and ask Armel:** "The snippet also says that if you're the only active member, leaving deletes the workspace. That's only true for workspaces with **fewer than 5 devices**: leaving never deletes anything itself, but a cleanup job (`prune-workspaces`, every 5 minutes) deletes workspaces with no active users and fewer than 5 devices. The app's own leave dialog says 'If you're the only member, this workspace will be deleted.' Add the caveat?"

- Option A (recommended): no. Change only the disabled-button sentence, which is what this issue is about. The deletion sentence matches the app's own dialog.
- Option B: yes. Also make edit (b).

Wait for his answer and record it in the Decision Log.

**Edits:**

(a) Replace the last line with:

```mdx
If the **Leave workspace** button is disabled, you are the workspace owner and other
members are still active. [Transfer ownership](/admin/workspace#transfer-ownership) to
another member before leaving. Members and admins can always leave.
```

(b) Only if Armel chose B. Change "However, if you are _the only_ active member of a workspace, leaving the workspace will delete the workspace and all associated resources." to "However, if you are _the only_ active member of a workspace and it has fewer than 5 devices, leaving the workspace will delete the workspace and all associated resources."

Leave the rest of the snippet alone.

**Checks and preview:** run the checks. Preview both pages that render the snippet: `http://localhost:3000/admin/workspace#leave-your-workspace` and `http://localhost:3000/admin/users/profile#leave-your-workspace`.

**Commit:** `docs/snippets/workspaces/leave.mdx`, `Refs ENG-1379`.

### ENG-1381: `docs(changelog): fix the broken connect a bucket link`

File: `docs/changelog/product.mdx`, line 422, inside the `<Update label="July 20, 2026">` entry, section "Connect a bucket".

**What's wrong:** `[Connect a bucket »](/data-uploads/connect-a-bucket)` points at a folder, not a page. Nothing in `docs.json` redirects it. The real pages are `/data-uploads/connect-a-bucket/aws` and `/data-uploads/connect-a-bucket/gcs` (`docs.json`, "Connect a bucket" group). The live site happens to redirect the folder to `/aws` (see Surprises).

**Stop and ask Armel:** "There's no single 'connect a bucket' page, only one guide per provider. Replace the link with two links (AWS and GCS), or with one link?"

- Option A (recommended): two links on one line. The paragraph introduces both AWS S3 and GCS, so each reader lands on their own guide.
- Option B: one link to the AWS guide. This keeps the changelog's one-link-per-section style and matches where the live redirect already goes, but GCS readers still land on AWS.

Wait for his answer and record it in the Decision Log.

**Edit:** replace line 422, keeping the four-space indentation. For Option A:

```mdx
    [Connect an AWS S3 bucket »](/data-uploads/connect-a-bucket/aws) · [Connect a GCS bucket »](/data-uploads/connect-a-bucket/gcs)
```

For Option B:

```mdx
    [Connect a bucket »](/data-uploads/connect-a-bucket/aws)
```

**Re-verify:** `grep -n "connect-a-bucket" docs/changelog/product.mdx` shows only full page paths, and each link opens on localhost.

**Checks and preview:** run the checks. Preview `http://localhost:3000/changelog/product` and scroll to the July 20, 2026 entry → "Connect a bucket". Click the new link(s).

**Commit:** `docs/changelog/product.mdx`, `Refs ENG-1381`.

## Concrete Steps

All commands run from `~/dev/miru/docs`.

0. Start state:

        git branch --show-current          # expect: docs/m1-fix-wrong-claims
        git log --oneline -2               # expect: 75fcf6d docs(groups): … on top of 7695aba
        git status --short                 # expect: only this plan (untracked)

   The plan file stays uncommitted while the issues are worked.

1. Finish ENG-1373's remaining step (the depth question) before anything else.

2. For each issue, in order (1374, 1375, 1376, 1377, 1378, 1379, 1381): re-verify, stop at each **Stop and ask Armel** step and wait, edit, run `./scripts/lint.sh` and `pnpm validate`, flag outdated screenshots, stop for Armel's preview and review, iterate, then Armel commits (Workflow above).

3. After ENG-1381, update Progress, the Decision Log and Outcomes in this plan. **Stop and ask Armel:** "How should the plan file land: its own commit, or folded into another commit?"

   - Option A (recommended): its own commit. `git add plans/active/20260924-m1-fix-wrong-claims.md`, subject `docs: add the m1 fix wrong claims plan`, no `Refs`.
   - Option B: fold it into a commit he names (`git commit --fixup <sha>` then `git rebase -i --autosquash main`).

4. After the last commit:

        git log --oneline main..HEAD       # expect 8 commits (9 with the plan), one per issue, subjects = Linear titles
        git diff --stat main..HEAD         # expect only the files named in this plan
        ./scripts/lint.sh && pnpm validate

   Then hand back to the Project for the final review and PR. Don't push from this session.

## Validation and Acceptance

1. `/concepts/groups/overview`: `name` says sibling-unique; Structure states the depth limit Armel chose.
2. `/cfg-mgmt/deploy/staging-area`: no "cannot be staged from this flow", no "coming soon"; the stage dialog mentions the change summary; drift points to Restage and Patch; Review lists only Restage and Archive.
3. `/cfg-mgmt/deploy/config-editor`: the history menu lists Copy short ID, Copy full ID, Copy description; the deployment alert text describes Acknowledge only, with no Discard/Keep; the `deployment-alert.png` frame is still in place.
4. `/concepts/deployments/overview`: has an `id` property explaining `dpl_…` vs `DPL-XXXXX`, if Armel chose to add it.
5. `docs/snippets/data-uploads/` no longer exists (or, if Armel chose to keep the files, neither snippet mentions `project_id`, `wip_provider`, `service_account_email` or a **Verify** action).
6. `/admin/users/overview`: has `type` (member, admin, owner) and no `role`; also `roles` (viewer, operator, provisioner, publisher) if Armel chose to add it.
7. `/admin/workspace` and `/admin/users/profile`: the disabled-button sentence names the owner-with-active-members condition.
8. `/changelog/product`: no link to the bare `/data-uploads/connect-a-bucket`.
9. Every **Stop and ask Armel** answer is recorded in the Decision Log, and every outdated screenshot listed in this plan was flagged to him.
10. No image or `<Frame>` was added, deleted or moved: `git diff main..HEAD -- docs/ ':!docs/snippets/data-uploads' | grep -E '^[-+].*(<Frame|!\[)'` prints nothing. The ENG-1377 snippet deletion is excluded because Armel approves it at that issue.
11. `./scripts/lint.sh` and `pnpm validate` pass on the final branch head; each issue commit's subject is its Linear title and its body has `Refs ENG-XXXX`.

## Idempotence and Recovery

- Every edit is a text replacement. Before reapplying, check whether the new text is already there (`grep`) and skip if so.
- Before a commit, `git checkout -- <file>` (or `git restore <file>`) undoes an edit, and `git restore --staged <file>` unstages.
- To change an earlier issue's commit after later ones exist: `git commit --fixup <sha>` then `git rebase -i --autosquash main`. That's safe because nothing is pushed. Never force-push.
- The frontend and backend repos are read-only. `git -C ~/dev/miru/frontend status --short` and `git -C ~/dev/miru/backend status --short` must stay empty.

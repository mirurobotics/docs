# Milestone 4: config editor and staging

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `~/dev/miru/docs` (`mirurobotics/docs`, this repo) | read-write | Every edit in this plan. Branch `docs/m4-config-editor`, created from `origin/main` at `8ae068e` (which includes M1 #206/#207, M2 #209 and M3 #210). Its first commit, `docs(plans): move the m2 and m3 plans to completed`, moves the M3 plan to `plans/completed/` and marks the M2 and M3 plans merged; it touches only those two plans. Do not create another branch. |
| `~/dev/miru/frontend` (`mirurobotics/frontend`) | read-only | Source of truth for the app: **`origin/prod`** (`2b3ffed6`, unchanged since the audit). `origin/main` is one commit ahead (`0dad9c8d`, #90), which only fixes in-app docs links, including the editor empty state's **Read the docs** link. The checkout is on `main` with Armel's uncommitted dev-mock edits (`src/api/devFleet.ts`, `src/proxy.ts`); leave them alone. Read with `git show "origin/prod:<path>"`; never switch branches. |
| `~/dev/miru/backend` (`mirurobotics/backend`) | read-only | Production backend: **`v0.11.9`** (`80acefec5c5072dde54c1b67cb4423589b2c9742`), confirmed on 2026-09-27 with `curl -s https://api.mirurobotics.com/frontend/v1/version`. Read with `git show v0.11.9:<path>`. Never `main`. Never modify. |

Linear: project [Frontend Docs Update](https://linear.app/mirurobotics/project/frontend-docs-update-a29451146f6a) (`P-ENG-55`), milestone **M4: Config editor & staging**: ENG-1392, 1393, 1394, 1395, 1396, 1397, 1403 and 1404, all Backlog on 2026-09-27, titles unchanged from the audit. Each description got a "Re-verified 2026-09-27" section with the details below. The audit behind this plan is `docs-audit.md` in the Project store (items 20, 23–26, 28–30).

## Purpose / Big Picture

The config editor was rebuilt in August, but the docs still describe the first version: two panels, per-file change sections, and a Deploy button that only checks JSON syntax. Only the changelog covers the rest:

- the config nav panel and closable tabs;
- the change panel with inline diffs and discard;
- live schema validation, and what blocks a deploy;
- adding and removing instance slots;
- XML and text editing;
- the release editor used for staging.

After this milestone, `cfg-mgmt/deploy/config-editor.mdx` describes the editor as it is, the initial-deployment guide and the quick start follow the current flow, and the schema pages explain validation and XML/text editing the way the editor now does them.

## Order of work, and why

M4 is the largest milestone. Most issues edit the same page (`config-editor.mdx`), and three of them (ENG-1392, ENG-1403, ENG-1404) describe flows that link into the editor sections. So the editor page goes first, section by section, top to bottom. The pages that point into it come after:

1. **ENG-1393**: layout, nav panel, tabs. Every later section refers to "the config panel".
2. **ENG-1394**: the change list, which ENG-1395 and ENG-1396 point into (**Errors**, **Removed**).
3. **ENG-1395**: validation and deploy blockers, plus the validation example in `schemas/validation.mdx`. Adds the `#validation` anchor that ENG-1397 and ENG-1403 link to.
4. **ENG-1396**: instance slots in drafts, which needs the config panel (1393) and the **Removed** section (1394).
5. **ENG-1397**: XML and text editing on `opaque.mdx`, which links to validation (1395).
6. **ENG-1392**: the initial-deployment guide, which links to the finished editor page.
7. **ENG-1403**: the release editor, described as "the same editor, with these differences", so it needs 1393–1396 in place.
8. **ENG-1404**: the quick start, last, because it summarizes all of the above.

## Progress

- [x] ENG-1393: config nav panel and closable tabs (2026-09-28)
- [x] ENG-1394: change panel with diffs and discard (2026-09-28)
- [x] ENG-1395: live schema validation and deploy blockers (2026-09-28)
- [x] ENG-1396: adding and removing instance slots in drafts (2026-09-28)
- [x] ENG-1397: canceled at Armel's call; nothing shipped (2026-09-28)
- [x] ENG-1392: initial deployment release selection steps (2026-09-28)
- [x] ENG-1403: skipped, then canceled at Armel's call; nothing shipped (2026-09-28)
- [x] ENG-1404: refresh the deploy configs quick start (2026-09-28)
- [x] Final: all commits on the branch, plan updated and folded into its commit. Merged as #213 (`08f54ae`).

## Surprises & Discoveries

Found while re-verifying (2026-09-27):

- **Nothing moved since the audit.** Production backend is still `v0.11.9` (`80acefec`); `/frontend/v1/version` reports `version 0.11.9`, `api_release_version v0.9.4`. There are no backend commits to review between the audit and production. Frontend `origin/prod` is still `2b3ffed6`. `git log --since=2026-09-24 origin/prod origin/main` on `src/features/editor`, `releases`, `devices/components/editor`, `config-instances`, `config-schemas`, `config-types`, `deployments`, `src/lib/codemirror`, `src/lib/diff` and the releases/devices routes returns nothing. `origin/main`'s only extra commit, #90, changes the editor empty state's **Read the docs** URL to the canonical `/cfg-mgmt/deploy/initial-deployment`.
- **M3 is merged** as #210, but its plan was still `plans/active/20260925-m3-groups.md` on `main`. At Armel's call, this branch's first commit moves it to `plans/completed/` and marks the M2 and M3 plans merged.
- **The Set Release dialog does open by itself, from the Devices page.** Clicking a device with no release sends you to the editor with the dialog open, if you can deploy (`src/features/devices/components/list/hooks/useDevicesList.ts`, `handleClick`, `SET_RELEASE_DIALOG`). Opening the **Editor** tab any other way shows the empty state with **Select a release**. The audit's "the dialog doesn't open by itself" was only half right.
- **The Set Release button is Select.** Not **Next** (the quick start) or **Update** (`set-release.png`) (`release/components/SetReleaseFooter.tsx`).
- **Validation is a backend call.** `POST /config_instances/validate`, 250 ms after you stop typing (`src/features/editor/session/validationScheduler.ts`). It covers JSON and YAML against JSON Schema or CUE, and is skipped for XML, plain text and opaque schemas (`session.ts:1245-1246`). XML well-formedness and JSON/YAML syntax are checked in the browser (`src/lib/codemirror/format/validate.ts`).
- **Stage isn't blocked by "no changes".** Deploy is disabled for errors, offline (unless allowed), no changes, and every config removed (`HistoryHeader.tsx`, `DeployButton`). **Stage** is disabled only for errors, "No configs to stage" and permission (`releases/components/editor/components/ChangePanel/StageButton.tsx`).
- **New drafts include required slots only.** Optional slots start absent and are added from the folder's plus (`devices/components/editor/utils/draftFiles.ts`; `releases/components/editor/utils/files.ts`). Narrowed at ENG-1396: only a device's first draft (see below).
- **The validation example is fully stale.** `schemas/validation.mdx` still shows a "Failed to create deployment" toast at deploy time and the old two-panel editor in all three screenshots.
- **The editor hero is outdated too.** `devices/editor/config-editor-v3.png` (last modified May 13) and `layout-temp.png` show the two-panel editor without the nav panel, and `draft-changelog.png` shows per-file sections.
- **The dev mock doesn't cover the editor.** `devFleet.ts` passes device and expanded-release GETs through to staging (`devFleet.ts:2024-2034`), so editor screenshots use real staging devices (`miru-14`). Staging runs backend `0.11.10` (`b97dcfe7`), one ahead of prod.
- **The nav panel shows change counts in Draft, not a bare M/A.** Rows show **R**, a red error count, **A**, `+N −M` lines (XML/text) or per-type counts like **2M** (`editor/utils/badges.ts:13-25`, `lib/diff/colors.ts:48`); `StatusMarkers` (bare **M**/**A**) is only the fallback and the tabs' markers. The plan's Layout draft was wrong; "Each file shows how many changes it has, or its errors in red." shipped instead.
- **Diffs in the change list are on by default** (`ChangeList/useDiffsVisible.ts`, `DEFAULT_DIFFS_VISIBLE = true`, remembered in local storage). **Show diffs** / **Discard all** are icon-only (diff and undo icons, tooltips "Show diffs"/"Hide diffs" and "Discard all"), not labeled buttons. The plan's "**Show diffs** adds the old and new values" was backwards.
- **The toy schema on `validation.mdx` doesn't compile.** `upload_interval_sec` and `heartbeat_interval_sec` aren't indented under `telemetry.properties` (lines 38-47), so `miru schema validate` fails with `/properties/telemetry/properties: got null, want object`. Indented four more spaces it validates. Fixed at ENG-1395.
- **Staging example for the validation page:** schemas in `~/dev/configs/examples/validation/` (commit `0d681b8`, `armeltalla/configs`), release `v7.0.0` (Mobility `SCH-DDjVz`, Network `SCH-NeRaR`, Perception `SCH-ALQTF`), device `miru-23`, first deployment `DPL-FaxUD` staged from `~/dev/miru/clones/miru-23` ("init: baseline mobility, network and perception configs").
- **Only a device's first draft is required-slots-only.** `releaseSeed` (`devices/.../utils/session.ts:61-66`) seeds required slots when the device has no deployment; a draft from a deployment keeps its slots (`deploymentSeed`, `:53-58`). Removing a slot you just added undoes the add instead of listing it under **Removed** (`session.ts:415-429`). Folders need two or more declared slots (`config-schemas/utils/slots.ts:12`). Removed files are left out of the deploy (`useDeploy.ts:71`).
- **XML well-formedness is an editor-only rule.** The editor rejects malformed XML and any `<!DOCTYPE>` ("DTDs are not allowed", `lib/codemirror/format/validate.ts:55-82`), but the backend never parses XML or text instances (`Decode` returns nil, `internal/configs/domain/config_instances/content.go:25-26`), so `miru device stage` isn't blocked by it.
- **Allow offline deployments is enforced only by the editor.** The backend reads the setting but no deploy path checks it (only `internal/orgs/...` settings files reference it; nothing in `internal/configs/services/deployments`). Staging isn't gated on the device being online either.
- **"Controls panel" / "editor panel" appear on other pages too:** `initial-deployment.mdx:62` (fixed with ENG-1392) and `staging-area.mdx:85,171` ("**Stage** in the editor panel", left as is when ENG-1403 was skipped).
- **CSpell doesn't know "urdf".** The ENG-1392 description example tripped it; `urdf` was added to `cspell.json` in that commit.
- **People who can't deploy open on History.** The **Draft** tab is disabled for them, with the reason in a tooltip (`HistoryHeader.tsx`; `useDeviceEditorSession.ts:60-64`).

## Decision Log

Record each of Armel's answers here as it's given at its issue (the questions themselves live in the issue sections), in this form:

    - Decision: <what Armel chose>, at ENG-XXXX.
      Rationale: <his reason, or the recommendation he accepted>.
      Date/Author: <date>, Armel.

Precedents carried over from M1 (`plans/completed/20260924-m1-fix-wrong-claims.md`), M2 (`plans/completed/20260924-m2-devices.md`) and M3 (`plans/completed/20260925-m3-groups.md`), which apply here without asking again:

- Keep the real behavior when the code and a message disagree (M1, ENG-1373).
- Commit the plan on its own, right after the first issue. Fold later plan updates into that commit (M1–M3).
- Keep sections short: write only what the user needs. Armel cut most M2 and M3 drafts further (M2 lessons, M3 ENG-1391).
- Changelog edits only when they fix a fact or a link; no new changelog links (M1, M2).
- Don't add caveats the app itself doesn't show unless they change what a reader does (M1, ENG-1379).
- **Screenshot rule (M2):** a screenshot stays only if it locates something text can't point to (an icon-only button, a control inside a row, a hidden menu), shows something the reader must recognize, or is the page's one hero. Dialogs whose fields the text names don't get one. The agent recommends per image; Armel decides. M3 dropped every dialog except Move.
- Armel replaces screenshots himself under a new file name, and reshapes sections live from them: propose the text and let the images follow. When he asks, update the embed's URL or remove the frame in that issue's commit. Never delete or move an image on your own.
- Name buttons as the app shows them when it's clear (M2, **+ Device**).
- Phrase dialog steps as "A dialog will appear—…" (M2, M3).
- Describe panes as changelog-style bullets: "**Tree:** …" (M3, ENG-1388).
- Write icon-only controls as a word plus the glyph in bold parentheses: "the plus **(+)**", "the ellipses **(...)**" (M3). Controls without a text glyph are named plainly ("the search icon").
- Name both the ellipses **(...)** and right-click where a menu opens both ways (M3).
- One commit may cover two issues when Armel asks, with the first issue's title as the subject and `Refs ENG-A, ENG-B` in the body (M3, ENG-1388/1389).
- Avoid apostrophes in headings that other text links to: Mintlify curls them in anchors (M3).
- Check the frontend dev mock (`src/api/devFleet.ts`) returns real data for anything Armel will screenshot (M3).
- Verify roles and limits against the backend before proposing badges and numbers (M2).
- Check CDN image names with `curl -I` (HEAD) only (M2).
- Re-read a file after editing if Armel has it open in the IDE (M2).
- Frontend bugs found while verifying go to **M9: Frontend bugs** (M2).

M4 decisions:

- Decision: replace the hero `devices/editor/config-editor-v3.png` with a three-panel shot; drop `layout-temp.png`. "Layout" is three short bullets (no search, collapse, middle-click or **Delete**), and "controls panel" is retired on this page (Draft intro, Deployment list), at ENG-1393.
  Also: the new hero is `editor/hero.png` (Armel's bucket folder for editor shots is `editor/`, not `devices/editor/config-editor/`). Retaken on 2026-09-28 without the search icon's "Show" tooltip and re-uploaded under the same name, an exception to the `-v2` rule: "we will just allow it to refresh once the cloudflare cache expires."
  Rationale: recommendations accepted.
  Date/Author: 2026-09-27, Armel.
- Decision: replace `devices/editor/draft-changelog.png` with a Draft-panel crop, one change hovered (undo icon) plus the two icons next to **Changes** (`editor/change-list.png`). The draft is shorter than the plan's (no path detail, no "only when present", no confirmation dialog), at ENG-1394.
  Rationale: recommendation accepted.
  Date/Author: 2026-09-28, Armel.
- Decision: a short `### Validation` before "Deploying changes", the Note shrunk to the blockers, the Errors bullet links to `#validation` (Option A); on `validation.mdx` replace `invalid-mobility-edit.png` and drop `invalid-mobility-err-msg.png` and `valid-mobility-edit.png`; rewrite the example (schema and text) to match a real staging schema error instead of the toy `max_angular_speed_radps` schema; `deploy-dialog.png` is left for ENG-1404, at ENG-1395.
  Reversed: keep the toy Mobility schema. Armel creates a new release and device with that schema on staging, and every screenshot on `validation.mdx` is reworked from it. Finish `config-editor.mdx` first, then the validation page.
  Rationale: recommendations accepted; "we should just keep the schema example that validation is using."
  Also: the new shot is `editor/validation/invalid-edit.png`, taken with the error tooltip showing ("actually this is better since you see the hover"), so the text adds "Hover over the underline to see why."
  Date/Author: 2026-09-28, Armel.
- Decision: `### Instance slots` in `config-editor.mdx` after "Change list" (Option A, given as "make sure what we were about to add to the docs is correct. if so, lets update the docs"); one new shot, the folder's plus **(+)** menu open (`editor/add-slot-menu.png`), no remove-menu shot, at ENG-1396. The draft says "a device's first draft" instead of "a new draft".
  Rationale: recommendation accepted.
  Date/Author: 2026-09-28, Armel.
- Decision: a short section on `opaque.mdx` (Option A), headed "## Editing XML and text files" with `"XML"` added to the heading-case allowlist and a test (Option B, against the recommendation), at ENG-1397. The draft adds the DTD rule and drops the opaque-schema sentence the page's first line already covers.
  Reversed: ENG-1397 is canceled and every change reverted (the `opaque.mdx` section and the allowlist change). "I actually don't think this section add any value. I dont think this issue is actually useful."
  Rationale: Armel's call.
  Date/Author: 2026-09-28, Armel.
- Decision: mention **Select a release** for when the dialog doesn't open (Option A); also drop "the controls panel header" from the Deploy step; text first, images later ("update the docs first. don't update any images yet"), at ENG-1392.
  Also: the dash sentence is cut and "Click the device…" joins the provisioned paragraph; `init-deploy-dialog.png` is dropped; `set-release.png` is hidden; the Steps keep titles with short bodies (every file added; an example description; skipped with Skip deploy confirmation; queued until online). New shots from `miru-03` (release v6.0.0, "init: baseline lidar, urdf, mobility and perception configs") under `editor/init/` ("I want init"): `devices-page.png`, `first-draft.png`, `history.png`.
  Rationale: recommendation accepted.
  Date/Author: 2026-09-28, Armel.
- Decision: skip ENG-1403 for now; no release editor docs in M4, and the issue goes back to Backlog. "I dont want us write docs for release editor. not right now." `staging-area.mdx`'s "**Stage** in the editor panel" (two places) stays as it is.
  Then: ENG-1403 is canceled ("ENG-1403, mark as canceled").
  Rationale: Armel's call.
  Date/Author: 2026-09-28, Armel.
- Decision: on the quick start, **Select** replaces **Next** with "The **Set Release** dialog opens."; the "log on the right" becomes the **Draft** panel; the Warning is dropped (Option A); "## Patch a deployment" is renamed "## Update the configs" because Patch is a staging action;   photos unchanged for now ("update the docs, don't change the photos"), at ENG-1404.
  Also: the quick start's six photos stay as they are in M4 ("lets not touch the pictures here").
  Rationale: recommendations accepted.
  Date/Author: 2026-09-28, Armel.

## Outcomes & Retrospective

Complete 2026-09-28. Eight commits on `docs/m4-config-editor` over `origin/main`: the M2/M3 plan move, this plan, and one commit each for ENG-1393, 1394, 1395, 1396, 1392 and 1404, each subject the Linear title with `Refs`. Those six issues are Done in Linear; ENG-1397 and ENG-1403 are Canceled. `./scripts/lint.sh`, `pnpm validate` and `mint broken-links` pass at the head. Merged as #213 (`08f54ae`).

What shipped, versus the plan:

- `cfg-mgmt/deploy/config-editor.mdx`: "Layout" is three bullets (Config panel with change counts, Code editor with closable tabs, Draft and History), under a new hero `editor/hero.png`; "controls panel" is gone from the page. "Change list" has Errors / Changes / Removed, the hover undo icon and the two icons next to **Changes** (`editor/change-list.png`). New `### Instance slots` (folders, first draft required-only, add from the plus **(+)**, right-click **Remove file** / **Restore file**; `editor/add-slot-menu.png`) and `### Validation`. The Deploy Note lists every blocker. No `## Release editor` (ENG-1403 skipped).
- `cfg-mgmt/concepts/schemas/validation.mdx`: the toy schema's indentation is fixed (it didn't compile), and the example is three sentences around `editor/validation/invalid-edit.png`, shot on staging from release `v7.0.0` / `miru-23` built from `~/dev/configs/examples/validation/`.
- `cfg-mgmt/deploy/initial-deployment.mdx`: **Set Release** / **Select** / **Select a release**, no "controls panel", Steps with short useful bodies, and three new shots under `editor/init/` from `miru-03`. `cspell.json` gains `urdf`.
- `getting-started/quick-start/deploy-configs.mdx`: **Select**, the **Draft** panel, no Warning, "Update the configs" instead of "Patch a deployment". Photos unchanged.
- Validation 1 holds except the release editor; 3 and 5 don't apply (ENG-1397 canceled, ENG-1403 skipped). Validation 9 deviates at Armel's direction: frames were added, dropped and swapped on his call.

Left open:

- ENG-1403 (release editor) is Canceled; `staging-area.mdx:85,171` still say "**Stage** in the editor panel", and `releases/stage/release-editor-tooltip.png`, `stage-button.png` and `release-editor-patch.png` show the old editor.
- The quick start's photos are the old editor (`devices/provisioned-device.png`, `devices/editor/set-release-dialog.png`, `init-deployment-draft.png`, `init-deploy-dialog.png`, `modifications.png`, `deploy-dialog.png`); `config-editor.mdx` still uses `devices/editor/deploy-dialog.png`.
- No longer referenced anywhere in docs, frontend (`origin/prod` and working tree) or backend (`v0.11.9`), so they can be deleted: `devices/editor/config-editor-v3.png`, `layout-temp.png`, `draft-changelog.png`, `invalid-mobility-edit.png`, `invalid-mobility-err-msg.png`, `valid-mobility-edit.png`, `set-release.png`, `first-draft.png`, `init-deployment-history.png`, `devices/init-dpl-devices-page.png`, and the trial uploads `editor/tmp-change-list.png` and `editor/tmp-add-slot-menu.png`. (`devices/editor/init-deploy-dialog.png` is still used by the quick start.)
- `editor/hero.png` was overwritten in place; the CDN serves the old one until its cache expires or the URL is purged.
- Backend gaps, not filed: **Allow offline deployments** and XML well-formedness are enforced only by the editor.

Lessons for M5 onward: verify the draft's every sentence against the code before proposing it (the plan's Layout markers and Show diffs direction were wrong); stop a check chain on lint failure (`&& tail` hid a CSpell error and a commit went through); Armel may cancel or skip an issue once he sees the overview, so keep edits out of the tree until the go.

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX under `docs/`. Navigation and redirects are in `docs/docs.json`, and reused content lives in `docs/snippets/`. Role badges come from `/snippets/components/role-badges.jsx`; API badges from `/snippets/components/platform-api-link.jsx`.

Pages M4 touches (line numbers on `8ae068e`):

- `docs/cfg-mgmt/deploy/config-editor.mdx`: hero; `## Layout` (21); `## Draft` (38) with `### Change list` (49) and `### Deploying changes` (67, the deploy steps and a `<Note>` on disabled Deploy at 92–99); `## History` (103) and its subsections; `## Deployment alert` (206, current since M1); `## Editor settings` (220).
- `docs/cfg-mgmt/concepts/schemas/validation.mdx`: `## Example` (7), with the stale story from about line 70.
- `docs/cfg-mgmt/concepts/schemas/languages/opaque.mdx`: example and file formats; nothing on editing.
- `docs/cfg-mgmt/deploy/initial-deployment.mdx`: `## Select a device` (12), `## Set the release` (28), `## Edit configs` (44), `## Deploy` (59), `## Summary` (93).
- `docs/cfg-mgmt/deploy/staging-area.mdx`: "Stage a deployment" already describes opening the release editor and the stage summary (M1, ENG-1374). One sentence ("Review the configuration tabs…", about line 84) can link to the new release editor section.
- `docs/getting-started/quick-start/deploy-configs.mdx`: the whole flow.
- The History subsections of `config-editor.mdx` (deployment list, diff modes, navigation, base vs head, redeploy) were checked and are still right. M4 doesn't rewrite them.

### Scope rules (from the Project's context doc; apply to every edit)

- **Archiving is out of scope.** "Redeploy — Restore an archived deployment" in History stays as it is.
- **Button naming isn't wrong.** Missing or differently-behaving buttons *are* issues (for example **Next** vs **Select**, or the old Discard/Keep alert).
- **Multi-step UI can be described as steps.**
- **No fleet dashboard.**
- **Keep edits short and match the M2/M3 style.**
- **Screenshots:** Armel updates screenshots himself. Never add, delete, or move an image or screenshot (see the workflow rules).

### Workflow (one issue at a time)

Rules for the executing agent:

- **Explain first, then ask.** At the start of each issue, send Armel one message covering what's wrong, what you'll change, and which files you'll touch. Then ask the issue's first **Stop and ask Armel** question, and wait for his go-ahead before editing.
- **Judgment calls are asked in context, at the issue.** Each issue marks its open questions as **Stop and ask Armel** steps, placed before the edit that depends on them. Ask with the options and recommendation, then wait. Don't batch questions up front or answer them yourself. Record each answer in the Decision Log. If you hit a judgment call the plan doesn't list, stop and ask the same way.
- **Never delete or move images or screenshots.** Keep every `<Frame>` and image exactly where it is until Armel says otherwise. When a screenshot is outdated, tell Armel at that issue: the image path, the section, what it shows wrong, and your keep/replace/drop recommendation under the screenshot rule. Each issue lists the ones already known.
- **Propose short drafts.** The drafts below are M3 length; cut further if a sentence doesn't change what the reader does.
- **Wrap new prose at 88 columns.**
- **Track status in Linear** if you have access: In Progress when you start, Done once Armel commits. Otherwise tell Armel.

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

5. Tick the issue in Progress and move to the next issue.

No push and no PR. When every issue is committed, Armel goes back to the Project for a final review, and the PR is opened from there as a draft into `main`, with the body listing each issue and its Linear ID. Suggested title: `docs(cfg-mgmt): document the rebuilt config editor and staging` (63 characters). ENG-1397 is `schemas`-scoped and ENG-1404 `quick-start`-scoped, so the strict alternative is the unscoped `docs: document the rebuilt config editor and staging`.

### Previewing (lessons from M1–M3)

- **Hard refresh after every edit** with Cmd+Shift+R.
- **Restart the dev server after editing `docs.json` or anything in `docs/snippets/`.** No M4 issue plans to. If one does: Ctrl+C in the server's terminal, then

        cd ~/dev/miru/docs && pnpm dev

- **Replaced screenshots can look stale.** Armel uploads replacements under a new name. Check a new name with `curl -I` only.
- **Confirm new anchors on localhost** by clicking the heading's link icon. `pnpm validate` and `mint broken-links` don't catch a curled apostrophe.

### Checks for every issue

Run from the repo root:

    ./scripts/lint.sh     # last line must be: All documentation lint checks passed.
    pnpm validate         # must end with: success build validation passed

Also run `(cd docs && ../node_modules/.bin/mint broken-links)`. If CSpell flags a real word, add it to `words` in `cspell.json` and include that file in the commit. Open every link you write on localhost.

### Re-verifying claims

Each issue lists the code to re-read. In `~/dev/miru/frontend`, run `git fetch origin prod`, then `git show "origin/prod:<path>"`. In `~/dev/miru/backend`, use `git show v0.11.9:<path>`. Before starting, re-run the version check (`curl -s https://api.mirurobotics.com/frontend/v1/version`). If production has moved, re-verify against the new commit. If the code no longer matches this plan, trust the code, correct the edit, and note it in Surprises & Discoveries. Frontend paths below are relative to `~/dev/miru/frontend`.

## Plan of Work

The drafts below were test-applied together on 2026-09-27 (lint, ESLint, CSpell, `pnpm validate`, `mint broken-links` all passed) and then restored. Match on the quoted text, since earlier edits shift lines.

### ENG-1393: `docs(cfg-mgmt): document the config nav panel and closable tabs`

File: `docs/cfg-mgmt/deploy/config-editor.mdx`, `## Layout`.

**What's wrong:** "The editor is split into two resizable panels" (editor and controls). There are three: the config panel, the code editor with closable tabs, and the Draft/History panel.

**Code to re-verify:**

- `src/features/devices/components/editor/EditorShell.tsx`: three `ResizablePanel`s. The nav panel is collapsible (`minSize='180px'`).
- `src/features/editor/components/ConfigNavPanel.tsx`: **Collapse panel** (panel icon); **Search configurations** toggles a filter box ("No configs match"). Config types with several slots are folders. Rows show the file-type icon and status markers.
- `src/features/editor/components/shared/StatusMarkers.tsx`: a red error count, **M** (modified), **A** (added); the file name takes the same color.
- `src/features/editor/components/EditorTabs/EditorTab.tsx`: close with the hover **×** (`CloseHotspot`), a middle-click, or **Delete**. History tabs show a lock ("Read-only").
- `src/features/editor/components/EditorTabs/EditorTabs.tsx`: **Open config panel** when the panel is collapsed. `shared/NoOpenConfigs.tsx`: **Open configs** when every tab is closed.

**Stop and ask Armel, screenshots:** under the screenshot rule, I recommend:

- the hero `devices/editor/config-editor-v3.png`: **replace** with the three-panel editor (nav panel, tabs, Draft panel with a change or two). It's the page's one hero, and it currently has no nav panel.
- `devices/editor/layout-temp.png` in "Layout": **drop**. The new hero shows the layout, and this one is the same old two-panel view.

Wait for his answer and record it.

**Edit:** replace the text between `## Layout` and the `layout-temp.png` frame (keep the frame until Armel decides) with:

```mdx
The editor has three resizable panels:

- **Config panel:** the device's config files. A red number marks a file with errors,
  **M** a modified file, and **A** an added one. Click the search icon to filter the
  list, or the panel icon to collapse it.
- **Code editor:** the open files, as tabs. To close a tab, hover over it and click the
  **×**, middle-click it, or press **Delete**.
- **Draft and History:** your pending changes and the device's past deployments,
  described below.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#layout`.

**Commit:** `docs/cfg-mgmt/deploy/config-editor.mdx`, `Refs ENG-1393`.

### ENG-1394: `docs(cfg-mgmt): rewrite the change panel section with diffs and discard`

File: `docs/cfg-mgmt/deploy/config-editor.mdx`, `### Change list`.

**What's wrong:** it describes per-file sections ("3 changes", "New file", "No changes") with a JSON path. The panel now has **Errors**, **Changes** and **Removed** sections. Each change can be discarded, **Show diffs** shows the old and new values, and clicking a row jumps to it.

**Code to re-verify:**

- `src/features/editor/components/ChangeList/ChangeList.tsx`: the three sections with counts; "No changes yet"; **Show diffs** / **Hide diffs**; a per-change discard (undo icon on hover, "Discard {type} change"); **Restore** on removed files.
- `ChangeList/DiscardAllButton.tsx`: "Discard all changes?" / "Every file in the draft goes back to its original content. Your edits will be lost."
- `ChangeList/FieldPath.tsx`: the path starts at the file or slot name, with chevrons; the full path is in a tooltip.
- `src/features/devices/components/editor/components/history/History.tsx` (`DraftChanges`): a click reveals the value in the file.

**Stop and ask Armel, screenshots:** `devices/editor/draft-changelog.png` shows the old per-file sections. I recommend **replacing** it with the current panel, one change hovered so the undo icon and the **Show diffs** / **Discard all** icons are visible. They're icon-only controls.

Wait for his answer and record it.

**Edit:** replace the text between `### Change list` and the `draft-changelog.png` frame with:

```mdx
The **Draft** panel lists your pending changes:

- **Errors:** syntax and schema errors, described in [Validation](#validation). Click
  one to jump to it.
- **Changes:** each value you added, modified, or deleted, with its path starting at
  the file name. Click a change to jump to it, or hover over it and click the undo icon
  to discard it. **Show diffs** adds the old and new values under each change, and
  **Discard all changes** resets every file.
- **Removed:** [slots](#instance-slots) you removed from the draft.
```

Until ENG-1395 and ENG-1396 add those sections, write the first bullet as "syntax and schema errors" without the link, and the third without the link. Those issues add the links back.

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#change-list`.

**Commit:** `docs/cfg-mgmt/deploy/config-editor.mdx`, `Refs ENG-1394`.

### ENG-1395: `docs(cfg-mgmt): document live schema validation and deploy blockers`

Files: `docs/cfg-mgmt/deploy/config-editor.mdx` (a new `### Validation` and the `<Note>` in "Deploying changes"); `docs/cfg-mgmt/concepts/schemas/validation.mdx` (`## Example`, second half).

**What's wrong:**

- The Note says Deploy is disabled only when the device is offline, there are no changes, or there are JSON syntax errors.
- `validation.mdx` says validation happens at deploy time and ends in an error toast.

Files are now validated as you type. Deploy is also blocked by schema errors, YAML and XML errors, and every config being removed.

**Code to re-verify:**

- `src/features/editor/session/validationScheduler.ts` (`VALIDATION_DEBOUNCE_MS = 250`) and `src/features/config-instances/api/repository.ts` (`validateConfigInstance`); backend `v0.11.9` `internal/configs/services/config_instances/validate.go`, `POST /config_instances/validate`.
- `src/features/editor/session/session.ts:1245-1246`: no schema validation for XML, text or opaque schemas. `src/lib/codemirror/format/validate.ts`: syntax checks for JSON, YAML, XML.
- `src/features/editor/session/resolveSchemaErrors.ts`: schema errors carry ranges (underlined in the file) and paths.
- `src/features/devices/components/editor/components/history/components/HistoryHeader.tsx` (`DeployButton`): the disabled reasons; `src/features/editor/utils/errorSummary.ts`: "`X` doesn't match its schema", "Fix the syntax error in X".
- `src/features/editor/components/ValidationGate.tsx`: "Validating configuration…", **Try again**, "Couldn't validate the configuration."

**Stop and ask Armel (1 of 2):** "Validation gets its own `### Validation` under Draft, and the Note in 'Deploying changes' shrinks to point there. OK, or keep it all in the Note?"

- Option A (recommended): a short `### Validation` section before "Deploying changes", and the Note replaced by one sentence linking to it. The change list (ENG-1394) and the XML section (ENG-1397) link to `#validation`.
- Option B: rewrite the Note in place and skip the new section.

**Stop and ask Armel (2 of 2), screenshots on `validation.mdx`:** all three show the old editor. I recommend:

- `invalid-mobility-edit.png`: **replace** with the current editor showing the underlined value and the **Errors** row.
- `invalid-mobility-err-msg.png`: **drop**. The toast it shows no longer exists. Alternatively, replace it with the disabled **Deploy** tooltip ("`Mobility` doesn't match its schema").
- `valid-mobility-edit.png`: **drop**. The fixed state needs no picture.

Wait for both answers and record them.

**Edits (Option A):**

(a) `config-editor.mdx`: insert before `### Deploying changes`:

```mdx
### Validation

As you edit, Miru checks each file's syntax and validates JSON and YAML files against
their [schemas](/cfg-mgmt/concepts/schemas/validation). Errors are underlined in the
editor and listed under **Errors** in the change list.
```

(b) Replace the whole `<Note>` block at the end of "Deploying changes" with:

```mdx
<Note>
  **Deploy** stays disabled while any file has errors, when there are no changes, when
  every config is removed, or when the device is offline (unless
  [offline deployments](#editor-settings) are allowed). Hover over the button to see
  why.
</Note>
```

(c) Restore the Errors bullet's link from ENG-1394: "…described in [Validation](#validation)."

(d) `validation.mdx`: replace everything from "Say we are deploying the following config instance." up to (not including) "Although simple in concept…" with the text below. Keep the three frames in their places until Armel decides.

```mdx
Say we set `max_angular_speed_radps` to `5.6` in the
[config editor](/cfg-mgmt/deploy/config-editor).

<Frame>
  ![Invalid Mobility Edit](https://assets.mirurobotics.com/docs/v04/images/devices/editor/invalid-mobility-edit.png)
</Frame>

The value exceeds the schema's maximum of 3.0, so the editor flags it as we type and
lists it under **Errors**. **Deploy** stays disabled until it's fixed.

<Frame>
  ![Invalid Mobility Error Message](https://assets.mirurobotics.com/docs/v04/images/devices/editor/invalid-mobility-err-msg.png)
</Frame>

Setting the value back within range clears the error, and the config instance can be
deployed.

<Frame>
  ![Valid Mobility Edit](https://assets.mirurobotics.com/docs/v04/images/devices/editor/valid-mobility-edit.png)
</Frame>
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#validation`, `#deploying-changes`, and `http://localhost:3000/cfg-mgmt/concepts/schemas/validation#example`.

**Commit:** both files, `Refs ENG-1395`.

### ENG-1396: `docs(cfg-mgmt): document adding and removing instance slots in drafts`

File: `docs/cfg-mgmt/deploy/config-editor.mdx`, a new `### Instance slots` under Draft, after "Change list".

**What's missing:** nothing explains that a config type with several slots shows as a folder, that a new draft includes only required slots, or how to add, remove and restore optional slots.

**Code to re-verify:**

- `src/features/editor/components/ConfigNavPanel.tsx`: `NavDirectoryRow` and `AddSlotMenu` ("Add a slot to {name}", shown on hover, listing missing slots); `NavRow` context menu with **Remove file** (disabled for required slots) and **Restore file**; removed rows dimmed.
- `src/features/devices/components/editor/utils/draftFiles.ts` and `src/features/releases/components/editor/utils/files.ts`: new drafts seed required slots only.
- `HistoryHeader.tsx`: "Every config is removed" blocks Deploy.
- Required vs optional: `docs/cfg-mgmt/concepts/schemas/annotations.mdx#instance-slots`.

**Stop and ask Armel (1 of 2):** "Put slots in the config editor page, or on `config-instances.mdx`?"

- Option A (recommended): `### Instance slots` in `config-editor.mdx`, since it's an editor action. `config-instances.mdx` already links to the editor.
- Option B: a section in `cfg-mgmt/concepts/config-instances.mdx`, linked from the editor.

**Stop and ask Armel (2 of 2), screenshots:** the add-slot plus only shows on hover. I recommend one new shot of the folder with its **Add a slot** menu open; Armel adds it if he agrees. Check first that `src/api/devFleet.ts` seeds a config type with optional slots.

Wait for both answers and record them.

**Edits (Option A):**

(a) Insert after the `draft-changelog.png` frame, before `### Validation`:

```mdx
### Instance slots

A config type with several
[instance slots](/cfg-mgmt/concepts/schemas/annotations#instance-slots) shows as a
folder in the config panel, with one file per slot. A new draft includes the required
slots.

To add an optional slot, hover over the folder, click the plus **(+)**, and pick the
slot. To remove one, right-click its file and select **Remove file**; select
**Restore file** to bring it back. Removed slots are listed under **Removed** in the
change list and aren't deployed.
```

(b) Restore the Removed bullet's link from ENG-1394: "[slots](#instance-slots) you removed from the draft."

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#instance-slots` and click the annotations link.

**Commit:** `docs/cfg-mgmt/deploy/config-editor.mdx`, `Refs ENG-1396`.

### ENG-1397: `docs(schemas): describe xml and text editing, diffs and checks`

File: `docs/cfg-mgmt/concepts/schemas/languages/opaque.mdx`, a new section before `## File formats`.

**What's missing:** the page lists XML and text as formats but says nothing about how the editor treats them.

**Code to re-verify:**

- `src/lib/codemirror/format/xml/index.ts` (`xml()` highlighting) and `text/index.ts` (no language); both `structural: false`.
- `src/lib/codemirror/format/validate.ts`: `validateXml`; text returns no errors.
- `src/features/editor/session/scheduler.ts` (`#textState`) and `DiffCountsLines`: line diffs with `+N −M` counts.
- `session.ts:1245-1246`: no schema validation for XML, text or opaque schemas.

**Stop and ask Armel (1 of 2):** "Where should XML and text editing go?"

- Option A (recommended): a short section on `opaque.mdx`, the page for these formats, linking to the editor's Validation section.
- Option B: a subsection in `config-editor.mdx`.

**Stop and ask Armel (2 of 2):** "The heading-case linter rejects 'XML' in a heading. It isn't in the acronym allowlist in `tools/lint/linter/headingcase/headingcase.go`. Heading?"

- Option A (recommended): "## Editing in the config editor". No tooling change; the body names XML and text.
- Option B: add `"XML"` to the allowlist (a tooling change with its own test in `tools/lint`; CI then runs the custom-linter jobs) and use "## Editing XML and text files".

Wait for both answers and record them.

**Edit (Option A):** insert before `## File formats`:

```mdx
## Editing in the config editor

In the [config editor](/cfg-mgmt/deploy/config-editor), XML files get syntax
highlighting and must be well-formed before they can be deployed. Plain text files
aren't checked. Files under an opaque schema are never
[validated against a schema](/cfg-mgmt/deploy/config-editor#validation).

Diffs for XML and text files compare lines rather than values, with changed words
highlighted and counts such as `+12 −4`.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/concepts/schemas/languages/opaque`. No screenshots on this page.

**Commit:** `docs/cfg-mgmt/concepts/schemas/languages/opaque.mdx`, `Refs ENG-1397`.

### ENG-1392: `docs(cfg-mgmt): update initial deployment release selection steps`

File: `docs/cfg-mgmt/deploy/initial-deployment.mdx`, `## Select a device` and `## Set the release`.

**What's wrong:**

- The Release column shows a **Set** button for candidates. It now shows a dash (minus icon).
- "a dialog will appear prompting you to select a release" is right when you click the device from the Devices page (and can deploy). Opening the Editor tab directly shows **Select a release** instead.
- The dialog's button is **Select**.

**Code to re-verify:**

- `src/features/devices/components/list/components/ReleaseColumn.tsx`: `MinusIcon` when there's no release.
- `src/features/devices/components/list/hooks/useDevicesList.ts` (`handleClick`): no release plus `created_deployed_deployment` allowed opens `?dialog=set-release`.
- `src/features/devices/components/editor/components/empty/EditorEmptyState.tsx`: **Select a release**, **Back to devices**, **Read the docs**.
- `src/features/devices/components/editor/components/release/*`: "Set Release", **Release**, **Select**.

**Stop and ask Armel (1 of 2):** "Mention the empty state's **Select a release** for when the dialog doesn't open by itself?"

- Option A (recommended): yes, one sentence. It's what readers see if they open the Editor tab directly or close the dialog.
- Option B: no. The guide already starts from the Devices page.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

- `devices/init-dpl-devices-page.png` (**Set** button): **replace** with a row showing the dash, which the reader must recognize.
- `devices/editor/set-release.png` (**Update** button): **drop**. It's a one-field dialog the text names.
- `devices/editor/first-draft.png`: **replace** (no nav panel).
- `devices/editor/init-deploy-dialog.png`: **drop**. The deploy dialog's fields are named in the steps.
- `devices/editor/init-deployment-history.png`: **replace** or drop; it shows the old History panel.

Wait for both answers and record them.

**Edits (Option A):**

(a) Replace the paragraph "You can identify these devices by their **Release** column. … Continue by clicking into the device you want to initialize." with:

```mdx
These devices show a dash in the **Release** column. Click the device you want to
initialize.
```

(b) Replace the first paragraph of "Set the release" ("Since the device doesn't have a current release, a dialog will appear…") with:

```mdx
Since the device doesn't have a release, the editor opens with the **Set Release**
dialog. Choose the release for the device's initial deployment and click **Select**. If
the dialog doesn't appear, click **Select a release** in the editor.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/initial-deployment`.

**Commit:** `docs/cfg-mgmt/deploy/initial-deployment.mdx`, `Refs ENG-1392`.

### ENG-1403: `docs(cfg-mgmt): document the release editor used for staging`

Files: `docs/cfg-mgmt/deploy/config-editor.mdx` (a new `## Release editor`) and one link in `docs/cfg-mgmt/deploy/staging-area.mdx`.

**What's missing:** `staging-area.mdx` says "the release editor opens" but never explains it: the same editor with **Stage** instead of Deploy, no History tab, and a context tooltip saying what the configs are based on.

**Code to re-verify:**

- `src/features/releases/components/editor/ReleaseEditor.tsx`: the same `ConfigNavPanel`, tabs and `ChangeList`; the right card is `ChangePanel` (no Draft/History tabs).
- `…/ChangePanel/EditorContext.tsx`: the device label with status, and the **Staging context** info tooltip: "All configs are pre-filled with their default values." / "Pre-filled from the device's current deployment (DPL-…)." / "Editing staged deployment (DPL-…)."
- `…/ChangePanel/StageButton.tsx`: disabled for errors, "No configs to stage" and permission; not for "no changes".
- `…/StageDialog/StageDialog.tsx`: summary, **Description**, validation gate.

**Stop and ask Armel (1 of 2):** "Where should the release editor go?"

- Option A (recommended): a short `## Release editor` at the end of `config-editor.mdx`, before "Editor settings", listing only the differences. `staging-area.mdx` links to it.
- Option B: a subsection of "Stage a deployment" in `staging-area.mdx`.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

- `releases/stage/release-editor-tooltip.png`: **replace**. The context tooltip is hidden behind an icon, so it earns a shot, but this one shows the old change list.
- `releases/stage/stage-button.png`: **drop**. The **Stage** button is named in the text.

Wait for both answers and record them.

**Edits (Option A):**

(a) `config-editor.mdx`: insert before `## Editor settings` (after the `---` that follows "Deployment alert"), followed by `---`:

```mdx
## Release editor

[Staging a deployment](/cfg-mgmt/deploy/staging-area#stage-a-deployment) opens the same
editor for a release and one device. It works like the **Draft** tab, except:

- **Stage** replaces **Deploy**, and there's no **History** tab.
- The info icon next to the device name shows what the configs are based on: the
  release's default values, the device's current deployment, or the staged deployment
  you're [patching](/cfg-mgmt/deploy/staging-area#patch-a-deployment).
- You can stage a deployment without changing anything.
```

(b) `staging-area.mdx`: change "Review the configuration tabs and make any edits you want." to "Edit the configs in the [release editor](/cfg-mgmt/deploy/config-editor#release-editor)."

**Checks and preview:** run the checks. Preview `http://localhost:3000/cfg-mgmt/deploy/config-editor#release-editor` and `http://localhost:3000/cfg-mgmt/deploy/staging-area#stage-a-deployment`.

**Commit:** both files, `Refs ENG-1403`.

### ENG-1404: `docs(quick-start): refresh the deploy configs guide for the new editor`

File: `docs/getting-started/quick-start/deploy-configs.mdx`.

**What's wrong:** "hit **Next**" in the release dialog (it's **Select**, and the dialog opens by itself); "a log on the right" for the change list; a Warning left over from the old editor; old screenshots.

**Code to re-verify:** as ENG-1392 (dialog), ENG-1394 (change list), and `useDeviceEditorSession.ts` / History (past deployments are read-only).

**Stop and ask Armel (1 of 2):** "The Warning 'Only the latest deployment may be edited; all others are read-only.' Keep it?"

- Option A (recommended): drop it. The draft always starts from the latest deployment and History is read-only, so there's nothing to warn about.
- Option B: keep it, reworded: "Past deployments in **History** are read-only."

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

- `devices/provisioned-device.png`: check; it's the devices list, so likely outdated since M2.
- `devices/editor/set-release-dialog.png`: **drop** (a one-field dialog).
- `devices/editor/init-deployment-draft.png` and `modifications.png`: **replace** (no nav panel; old change list).
- `devices/editor/init-deploy-dialog.png` and `deploy-dialog.png`: **drop**, or keep one. `deploy-dialog.png` is close to current but shows a browser extension icon.

Wait for both answers and record them.

**Edits (Option A):**

(a) Replace "Select the release we just created and hit **Next**." with:

```mdx
The **Set Release** dialog opens. Select the release we just created and click
**Select**.
```

(b) Replace "To edit a config instance, make some changes to the current deployment in the editor. As edits are made, a log on the right maintains a list of all changes, categorizing them by their type: `added`, `deleted`, or `modified`." with:

```mdx
To edit a config instance, make some changes in the editor. The **Draft** panel on the
right lists each change as added, modified, or deleted.
```

(c) Delete the `<Warning>` block ("Only the latest deployment may be edited…").

**Checks and preview:** run the checks. Preview `http://localhost:3000/getting-started/quick-start/deploy-configs`.

**Commit:** `docs/getting-started/quick-start/deploy-configs.mdx`, `Refs ENG-1404`.

## Concrete Steps

All commands run from `~/dev/miru/docs`.

0. Start state:

        git branch --show-current          # expect: docs/m4-config-editor
        git log --oneline -2               # expect: docs(plans): move the m2 and m3 plans to completed, on top of 8ae068e (#210)
        git status --short                 # expect: only this plan (untracked)
        curl -s https://api.mirurobotics.com/frontend/v1/version   # expect git_commit 80acefec…

   The branch already has one commit, the M2/M3 plan move. The plan stays uncommitted until ENG-1393 is committed. Then commit it on its own: `git add plans/active/20260927-m4-config-editor.md`, subject `docs: add the m4 config editor plan`, no `Refs`. Fold later plan updates into that commit (`git commit --fixup <sha>` then `git rebase -i --autosquash origin/main`).

1. For each issue, in the order above: explain it, re-verify, stop at each **Stop and ask Armel** step and wait, edit, run the checks, flag screenshots, stop for Armel's preview and review, iterate, then Armel commits.

2. After ENG-1404, update Progress, the Decision Log and Outcomes in this plan, then fold that update into the plan commit.

3. After the last commit:

        git log --oneline origin/main..HEAD   # expect the M2/M3 plan move, the plan, and one commit per issue (fewer if Armel combined any)
        git diff --stat origin/main..HEAD     # expect config-editor.mdx, validation.mdx, opaque.mdx, initial-deployment.mdx, staging-area.mdx, deploy-configs.mdx, this plan, and the M2/M3 plans in plans/completed/
        ./scripts/lint.sh && pnpm validate

   Then hand back to the Project for the final review and PR. Don't push from this session.

## Validation and Acceptance

1. `/cfg-mgmt/deploy/config-editor`: "Layout" names three panels as bullets. "Change list" has Errors / Changes / Removed, discard and diffs. `### Validation` and `### Instance slots` exist. The Deploy Note lists every blocker. `## Release editor` lists the differences.
2. `/cfg-mgmt/concepts/schemas/validation`: the example flags the error as you type; no toast at deploy time.
3. `/cfg-mgmt/concepts/schemas/languages/opaque`: says how XML and text are highlighted, checked and diffed.
4. `/cfg-mgmt/deploy/initial-deployment`: a dash in the Release column, and the Set Release dialog with **Select**.
5. `/cfg-mgmt/deploy/staging-area`: "Stage a deployment" links to the release editor section.
6. `/getting-started/quick-start/deploy-configs`: **Select**, the Draft panel, and no leftover Warning (unless Armel kept it).
7. `grep -rn "two resizable panels\|JSON syntax errors\|placeholder button\|hit \*\*Next\*\*\|log on the right" docs --include='*.mdx'` prints nothing.
8. Every **Stop and ask Armel** answer is in the Decision Log, and every screenshot listed here was flagged with a recommendation.
9. No image or `<Frame>` was added, deleted or moved except as Armel decided: `git diff origin/main..HEAD -- docs/ | grep -E '^[-+].*(<Frame|!\[)'` shows only changes he asked for.
10. No archiving content was added: `git diff origin/main..HEAD -- docs/ | grep -i -E '^\+.*archiv'` prints nothing.
11. `./scripts/lint.sh`, `pnpm validate` and `mint broken-links` pass on the final head. Each issue commit's subject is its Linear title, and its body has `Refs`.

## Idempotence and Recovery

- Every edit is a text insertion or replacement. Before reapplying, check whether the new text is already there (`grep`) and skip it if so.
- Before a commit, `git restore <file>` undoes an edit, and `git restore --staged <file>` unstages.
- To change an earlier issue's commit after later ones exist: `git commit --fixup <sha>` then `git rebase -i --autosquash origin/main`. That's safe because nothing is pushed. Never force-push.
- The frontend and backend repos are read-only. Don't touch the frontend's uncommitted dev-mock files, and don't switch either checkout's branch.

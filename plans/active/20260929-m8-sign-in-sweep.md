# Milestone 8: sign-in and final sweep

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `~/dev/miru/docs` (`mirurobotics/docs`, this repo) | read-write | Every edit in this plan. Branch `docs/m8-sign-in-sweep`, created from `origin/main` at `e950a27` (M7 #218). Its first commit, `96595dd` (`docs(plans): move the m7 plan to completed`, signed), moves the M7 plan to `plans/completed/` and marks it merged; it touches only that plan. Do not create another branch. |
| `~/dev/miru/frontend` (`mirurobotics/frontend`) | read-only | Source of truth for the app: **`origin/prod`** (`2b3ffed6`, unchanged since the audit). `origin/main` is one commit ahead (`0dad9c8d`, #90, in-app docs links). Another agent is planning M9 in this checkout: only `git fetch` and `git show "origin/prod:<path>"`, never switch branches or touch files. |
| `~/dev/miru/backend` (`mirurobotics/backend`) | read-only | Production backend: **`v0.11.9`** (`80acefec5c5072dde54c1b67cb4423589b2c9742`), confirmed on 2026-09-29 with `curl -s https://api.mirurobotics.com/frontend/v1/version`. Read with `git show "v0.11.9:<path>"`. Never `main`. Never modify. |

Linear: project [Frontend Docs Update](https://linear.app/mirurobotics/project/frontend-docs-update-a29451146f6a) (`P-ENG-55`), milestone **M8: Sign-in & final sweep**: ENG-1413, ENG-1414 and ENG-1415, all Backlog on 2026-09-29, titles unchanged from the audit. Each description got a "Re-verified 2026-09-29" section with the details below. The audit behind this plan is `docs-audit.md` in the Project store (items 46, 47, 50).

## Purpose / Big Picture

- **CLI login:** the docs say only to confirm that the browser shows the terminal's code. The browser opens an **Authorize Device** page that shows the request (device, location, time, IP) and asks you to click **Allow**; it isn't documented anywhere but the Sep 23 changelog.
- **Sign-up:** the quick start says to "follow the instructions" on the sign-up page. Sign-up is Google or name and email, an emailed code, then naming the workspace.
- **Screenshots:** after M1–M7, 109 older app screenshots are still embedded. Most are fine; 32 show UI that has changed (the old sidebar with **Tags**, removed tabs, old badges and the pre-M4 editor) and 16 are dialogs the text already names.

After this milestone, the CLI login and sign-up docs match frontend `origin/prod` and backend `v0.11.9`, and every outdated screenshot is replaced or dropped.

## Order of work, and why

1. **ENG-1413**: the CLI Authorize Device page (`developers/cli/authentication.mdx`). One new subsection; no dependencies.
2. **ENG-1414**: sign-up, log in and onboarding (`getting-started/quick-start/overview.mdx`). One paragraph becomes a short list; its hero is part of ENG-1415's batch F, so it's captured then.
3. **ENG-1415**: the screenshot sweep, last because it depends on Armel's capture sessions. Text-first: the embed swaps that reuse existing images and the drops go in first, then one capture batch per session.

Then the **ENG-1427 follow-up** (not an M8 issue): check the Duplicate a release text against production once ENG-1427 ships.

## Progress

- [x] ENG-1413: the CLI device authorization page
- [ ] ENG-1414: sign up, log in and workspace onboarding (canceled)
- [ ] ENG-1415: the screenshot sweep, all batches (canceled)
- [x] Follow-up: ENG-1427 is still Backlog; the Duplicate a release text stays
- [x] Final: all commits on the branch, plan updated and folded into its commit

## Surprises & Discoveries

Found while re-verifying (2026-09-29):

- **Nothing in production moved since the audit.** `/frontend/v1/version` still reports `version 0.11.9`, `git_commit 80acefec`, `api_release_version v0.9.4`. Frontend `origin/prod` is still `2b3ffed6`; `git log --since=2026-09-24 origin/prod` on `src/views/authn`, `src/app/(authn)`, `src/app/oauth` and `src/app/sso` returns nothing.
- **M7 is merged** as #218 (`e950a27`), titled "docs: update buckets, api keys and scopes". ENG-1410, ENG-1411 and ENG-1412 were canceled and ENG-1429 (the scope table) was added. Its plan was moved to `plans/completed/` in this branch's first commit, `96595dd`.
- **Device codes last 30 minutes, not "a few minutes".** The backend sets `DeviceAuthTTLSecs = 1800` (`internal/orgs/services/deviceauth/create.go:25`); the page's error says "Codes expire after a few minutes — run `miru login` again" (`frontend src/views/authn/cli/utils/error.ts`). A copy issue for M9; not filed.
- **Wrong codes are rate-limited per account:** 10 per 10 minutes (`deviceauth/guesses.go:16-17`), which matches the page's "Too many attempts" text.
- **Sign-in destinations are one function** (`frontend src/views/authn/shared/utils/destination.ts`): an invite goes to `/invites/<id>` (skipping onboarding), a new account goes to `/onboarding`, a `miru login` code returns to `/oauth/device?user_code=…`, otherwise `/devices`. Google sign-in uses the same function through `/sso/google`.
- **Onboarding names an existing workspace.** `useOnboardingForm.ts` calls `updateWorkspace` and marks the user onboarded; it doesn't create the workspace.
- **The CLI login snippet is shared.** `snippets/references/cli/login.mdx` (usage and examples) is used by `references/cli/login.mdx`, `developers/cli/authentication.mdx` and the quick start's `create-release.mdx`. The browser side goes only into `authentication.mdx`, leaving the CLI reference (terminal side) as it is.
- **The inventory is smaller than the audit's 257.** Outside `changelog/` and `references/`, the docs embed 236 images (205 unique URLs); 78 are backgrounds, SVG diagrams, third-party shots or archiving-only, leaving 158 app screenshots. 49 of those were replaced in M1–M7 (M1 4, M2 14, M3 10, M4 7, M5 7, M7 11 embeds, counting archiving menus separately).
- **The staging area is the most outdated page.** Its shots mix two eras: v7.0.0 crops that match today, and v1.6 shots with the old sidebar (**Tags**), a **Deployments** tab that no longer exists (`useReleaseTabs`: Overview and Stage), **NEW** badges (now **ADDED**, `deployments/components/sheet/utils/badge.ts`) and config instance Content / Metadata tabs that were removed (M5). The stage flow itself is unchanged: row menu Deploy / Archive / Patch / Review (`releases/components/staging/utils/actions.ts`), the review dialog with **Restage**.
- **The schema sheet still has Schema / Metadata tabs** (`config-schemas/components/SchemaSheet/SchemaSheet.tsx`), plus a document picker for multi-document schemas, so the schema shots stay.
- **The quick start can't reuse M4/M5 shots.** It shows its own data (release v1.0.0; Communication, Mobility, Planning), while M4's `editor/init/*` and M5's `releases/overview-tab.png` show other example configs. M4 left the quick start's images for later ("text first, images later", ENG-1404).
- **ENG-1399's text is already live and ahead of the code.** `concepts/releases.mdx` says Duplicate copies "config schemas, file rules, and git information" (M5 reversed ENG-1399), shipped in #216. ENG-1427 (M9), which makes Duplicate send `file_rule_ids`, is still Backlog, so production doesn't copy file rules today.
- **The CDN rejects Python's default user agent** (403 from `urllib`); `curl` works. Use `curl` for any image download.

Found while executing (2026-09-29):

- **The Authorize Device request rows have no visible labels.** Device, Location, Requested at and IP address are `sr-only` (`RequestDetails.tsx:123`); the page shows an icon and a value.
- **Approving binds the CLI to the browser's current workspace** (`deviceauth/approve.go:90-92`). Left out of the docs.
- **Location comes from Cloudflare's IP headers** (`cli/v06/handlers/auth.go:35-40`); the IP is the CLI request's client IP.
- **Linear ignored the "In Progress" state by name and by ID;** passing the state type `started` worked.
- **The frontend baseline had an untracked `.ai/plans/M9_FRONTEND_BUGS.md`** (the M9 agent's) besides Armel's two mock edits.

## Decision Log

Record each of Armel's answers here as it's given at its issue (the questions themselves live in the issue sections), in this form:

    - Decision: <what Armel chose>, at ENG-XXXX.
      Rationale: <his reason, or the recommendation he accepted>.
      Date/Author: <date>, Armel.

Insert a new entry only after the previous entry's closing Date line (M5 lesson).

Precedents carried over from M1–M7 (`plans/completed/20260924-m1-fix-wrong-claims.md`, `20260924-m2-devices.md`, `20260925-m3-groups.md`, `20260927-m4-config-editor.md`, `20260928-m5-releases.md`, `20260929-m7-settings-admin.md`), which apply here without asking again:

- Keep the real behavior when the code and a message disagree (M1).
- Commit this plan on its own **before** the first issue, and fold later plan updates into that commit (M5, M7).
- Keep sections short: write only what the user needs. Armel cuts drafts further and reshapes sections from the live page, often into bullets for "what a page shows" (M2–M7).
- Verify every draft sentence against the code before proposing it (M4, M5).
- Changelog edits only when they fix a fact or a link; no new changelog links (M1, M2).
- Don't add caveats the app itself doesn't show unless they change what a reader does (M1); M7 left the 30-second API key cache out for this reason.
- **Screenshot rule (M2):** a screenshot stays only if it locates something text can't point to (an icon-only button, a control inside a row, a hidden menu), shows something the reader must recognize, or is the page's one hero. Dialogs whose fields the text names don't get one. The agent recommends per image; Armel decides. He may also say "text first, images later" (M4, M5).
- Armel replaces screenshots himself, under a new name in a topic folder (a `-v2` suffix when the old name still fits, M7), and reshapes sections live from them. When he asks, update the embed's URL or remove the frame in that issue's commit. Never delete or move an image on your own.
- On M5's releases page and M7's buckets and API keys pages Armel removed every dialog image and used one menu shot per action. Recommend the same, but ask per page.
- **Overviews are concise (M5):** what changes (the new text, one line of evidence each) and a per-image table saying plainly whether Armel **adds**, **replaces** or **drops** each image. Re-state inherited decisions in the overview, since Armel may reverse one once he sees the code (M5, ENG-1399; M7, the scope table).
- Name buttons as the app shows them when it's clear, but don't "fix" wording that already identifies the right button (M2; project scope rule).
- **Dialog steps are two sentences** ("A confirmation dialog will appear. Click **Archive** to confirm."), not an em dash (M7; reverses the M2/M3 "A dialog will appear—…").
- Describe panes and lists of parts as changelog-style bullets: "**Tree:** …" (M3, M4).
- Write icon-only controls as a word plus the glyph in bold parentheses: "the ellipses **(...)**" (M3).
- Name both the ellipses **(...)** and right-click where a menu opens both ways (M3, M5, M7).
- One commit may cover two issues when Armel asks: the first issue's title as the subject, and `Refs ENG-A, ENG-B` in the body (M3).
- Armel may cancel or skip an issue once he sees the overview (M4, M7), so keep edits out of the tree until he says go. A milestone can also gain an issue (M7, ENG-1429); re-propose the PR title at `/pr`.
- Canceled issues stay canceled: M8 doesn't reopen the text on `manage.mdx`, `invites.mdx`, `workspace.mdx` or `profile.mdx` (ENG-1410 to 1412, M7). Screenshots on those pages are judged by what they show.
- Avoid apostrophes in headings that other text links to (M3).
- For screenshots, use a mock with fake names, emails and IDs so nothing needs blurring, and snapshot the mock file before editing so the undo is an exact restore (M7 lessons). Check the dev mock has data for anything Armel will capture (M3). The frontend checkout is shared with the M9 agent: coordinate before touching the mock.
- Verify roles and limits against the backend before proposing badges and numbers (M2).
- Check CDN image names with `curl -I` (HEAD) only (M2). Every new name proposed below returned 404 on 2026-09-29.
- Re-read a file after editing if Armel has it open in the IDE (M2).
- Run each check on its own, and stop on the first failure; don't pipe `lint.sh` into `tail` in a chain that commits (M4). Auto-review can block some `mint broken-links` runs; retry with approval (M7).
- Frontend bugs found while verifying go to **M9: Frontend bugs** (M2).
- Write Linear issues one at a time and re-read each with `get_issue` before retrying; the list view lags (M5, M7).

M8 decisions:

- Decision: "Authorize in your browser" uses the shorter draft (no "(or sign up)", no
  repeat of the page's Allow warning); expiry left out (Option B); bold field labels
  kept although the page shows icons only; no workspace sentence; the shot goes after
  the field list, at ENG-1413.
  Also: `cli/authorize-device.png` is Armel's Sep 22 capture with the IP blocked by
  Armel ("the code is old so it doesn't matter"); the agent's edited copy was deleted.
  Reversed: no "Authorize in your browser" heading ("it should go after the miru login
  command with the next steps"). The page imports `usage.mdx` and `examples.mdx`
  directly and puts the browser text between them; the snippets are unchanged.
  Reversed: no browser text at all ("explaining each thing that is show isn't needed…
  everything else you wrote isn't needed"); the snippet's "Confirm that the code in the
  browser matches…" is enough. Only the screenshot sits between usage and examples.
  Then: one or two sentences between the photo and the CLI output ("maybe dont just
  cut to what I said… make it more concise"). Final, after "think about it from first
  principles": "If the device and location are yours, click **Allow**. The terminal
  then shows who you're logged in as:", which introduces the example output.
  Then: the whole page rewritten short ("you can change the whole thing. it doesn't
  just need to be this one section"): the page writes its own command and browser
  steps instead of importing `usage.mdx` (the snippet is unchanged and still used by
  the CLI reference and the quick start), keeps `examples.mdx`, adds the
  `--no-browser` line, and folds the API key section into three short paragraphs.
  Then: the terminal example uses the screenshot's code ("use the code in the
  screenshot in the terminal output"), inlined on the page as a copy of
  `examples.mdx` with `BDKZ-LQXS`; the snippet is unchanged. `BDKZ` and `LQXS` go in
  `cspell.json`.
  Rationale: recommendations accepted ("okay lgtm").
  Date/Author: 2026-09-29, Armel.
- Decision: cancel ENG-1414 and ENG-1415 ("we should cancel all the other issues. I
  don't think we should be doing those"). No edits were made for them.
  Date/Author: 2026-09-29, Armel.
- Decision: leave the Duplicate a release text as it is while ENG-1427 is Backlog;
  the docs describe the intended behavior and ENG-1427 ships to prod with it (M5).
  Rationale: recommendation accepted.
  Date/Author: 2026-09-29, Armel.

## Outcomes & Retrospective

Complete 2026-09-29. Three commits on `docs/m8-sign-in-sweep` over `origin/main`: the
M7 plan move, this plan, and ENG-1413. ENG-1413 is Done; ENG-1414 and ENG-1415 are
Canceled. `./scripts/lint.sh`, `pnpm validate` and `mint broken-links` pass at the
head. The frontend dev mock was never edited; its stat and checksums match the
baseline.

What shipped, versus the plan:

- `developers/cli/authentication.mdx`, rewritten short: the command, one sentence on
  the **Authorize Device** page (the code matches, the device and location are yours,
  click **Allow**), the screenshot `cli/authorize-device.png` (Armel's capture, IP
  blocked), the terminal output with the screenshot's code, a `--no-browser` line, and
  a three-paragraph API key section. No "Authorize in your browser" heading and no
  field bullets, the reverse of the plan. The page no longer imports the shared
  `login` snippets; they're unchanged and still used by the CLI reference and the
  quick start.
- `cspell.json`: `BDKZ`, `LQXS`, and `Armel` (Armel renamed the example's
  "Benjamin" to himself).

Left open:

- Canceled: ENG-1414 (the quick start still says "follow the instructions") and
  ENG-1415 (the 32 outdated and 16 dialog shots in the inventory stay as they are).
- ENG-1427 (M9) is still Backlog; `concepts/releases.mdx` says Duplicate copies file
  rules, ahead of production.
- Not filed: the Authorize Device error's "Codes expire after a few minutes" (the
  backend keeps codes 30 minutes).
- The page now describes the terminal side of `miru login` itself (command, browser
  opening, `--no-browser`, example output), which the backend team owns in the
  snippets; `MIRU_API_KEY` precedence and `--no-browser` weren't verified against
  the CLI source (not checked out).

No images were removed or replaced, so none can be deleted.

Validation, run 2026-09-29 on the final head:

- 1 holds in the shortened form above. 2 and 4 don't apply (canceled).
- 3 prints `getting-started/quick-start/overview.mdx:29` (ENG-1414 canceled).
- 5 shows only the new `cli/authorize-device.png` frame. 6 and 7 print nothing.
- 8: every answer is in the Decision Log. 9: the checks pass; ENG-1413's commit
  subject is its Linear title with `Refs`.

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX under `docs/`. Navigation and redirects are in `docs/docs.json`; reused content lives in `docs/snippets/`.

Pages M8 touches (line numbers on `96595dd`):

- `docs/developers/cli/authentication.mdx`: `## Interactive login` (9–11). The CLI reference snippets under `docs/snippets/references/cli/login/` stay as they are.
- `docs/getting-started/quick-start/overview.mdx`: the sign-up sentence (29) and the hero (13).
- ENG-1415: the pages in the inventory below. The text-first pass touches `admin/users/invites.mdx`, `admin/users/profile.mdx`, `admin/workspace.mdx`, `cfg-mgmt/concepts/config-types.mdx`, `cfg-mgmt/concepts/schemas/manage.mdx`, `cfg-mgmt/deploy/staging-area.mdx`, `developers/agent/versions.mdx`, `developers/platform-api/overview.mdx`, `getting-started/quick-start/deploy-configs.mdx` and `provision-devices/reprovision.mdx`.
- Follow-up: `docs/concepts/releases.mdx`, `## Duplicate a release` (read only unless ENG-1427 is dropped).

### Scope rules (from the Project's context doc; apply to every edit)

- **Archiving is out of scope.** Images in archiving-only sections (Archive / Unarchive a release, device, bucket, config type or deployment) are excluded from the sweep; don't add archive or unarchive text.
- **No fleet dashboard.** Don't add or describe the `/dashboard` page; none of the inventory is dashboard content.
- **Button naming isn't wrong.** Missing or differently-behaving controls *are* issues.
- **The terminal side of `miru login` belongs to the backend team.** ENG-1413 documents only the browser page, on `developers/cli/authentication.mdx`; the CLI reference snippets aren't edited.
- **List docs only for devices.** Other lists get no list-controls docs.
- **Keep edits short, and match the M2–M7 style.**
- **Screenshots:** Armel updates screenshots himself. Never add, delete, or move an image without his decision (see the workflow rules).

### Workflow (one issue at a time)

Rules for the executing agent:

- **Explain first, then ask.** At the start of each issue (and each ENG-1415 batch), send Armel the concise overview (what changes, per-image table). Then ask the first **Stop and ask Armel** question, and wait for his go-ahead before editing.
- **Judgment calls are asked in context, at the issue.** Each issue marks its open questions as **Stop and ask Armel** steps, placed before the edit that depends on them. Ask with the options and recommendation, then wait. Don't batch questions up front or answer them yourself. Record each answer in the Decision Log. If you hit a judgment call the plan doesn't list, stop and ask the same way.
- **Never delete or move images or screenshots** on your own. Keep every `<Frame>` and image exactly where it is until Armel says otherwise.
- **Propose short drafts,** and re-verify each sentence against the code before proposing it.
- **Wrap new prose at 88 columns.**
- **Track status in Linear** if you have access: In Progress when you start, Done once Armel commits. Write one issue at a time and re-read it to confirm.

Steps:

1. Work through **one** issue in the order above. Explain it, stop at each **Stop and ask Armel** step, make the edits, run the checks, and leave the change **uncommitted**.
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

No push and no PR. When every issue is committed, Armel goes back to the Project for a final review, and the PR is opened from there as a draft into `main`, with the body listing each issue and its Linear ID. Suggested title: `docs: document cli login and sign-up, refresh screenshots` (56 characters). The issues span `cli`, `quick-start` and many sections, so the title has no scope.

### Previewing (lessons from M1–M7)

- **Hard refresh after every edit** with Cmd+Shift+R.
- **Restart the dev server after editing `docs.json` or anything in `docs/snippets/`.** No M8 edit plans to. If one does: Ctrl+C in the server's terminal, then

        cd ~/dev/miru/docs && pnpm dev

- **Replaced screenshots can look stale.** Check a new name with `curl -I` only.
- **Confirm new anchors on localhost** by clicking the heading's link icon (`#authorize-in-your-browser`).

### Checks for every issue

Run from the repo root, one at a time, and stop on the first failure:

    ./scripts/lint.sh     # last line must be: All documentation lint checks passed.
    pnpm validate         # must end with: success build validation passed
    (cd docs && ../node_modules/.bin/mint broken-links)   # must end with: success no broken links found

If CSpell flags a real word, add it to `words` in `cspell.json` and include that file in the commit. Open every link you write on localhost.

### Re-verifying claims

Each issue lists the code to re-read. Before starting, re-run the version check (`curl -s https://api.mirurobotics.com/frontend/v1/version`); if production has moved, re-verify against the new commit. In `~/dev/miru/frontend`, run `git fetch origin prod`, then `git show "origin/prod:<path>"`. In `~/dev/miru/backend`, use `git show "v0.11.9:<path>"` (quote it: zsh reads `$C:` as a modifier). If the code no longer matches this plan, trust the code, correct the edit, and note it in Surprises & Discoveries. Frontend paths below are relative to `~/dev/miru/frontend`.

## Plan of Work

The drafts below (ENG-1413, ENG-1414, and ENG-1415's reuse swaps and drops, each at its recommended option) were test-applied together on 2026-09-29 (`./scripts/lint.sh`, `pnpm validate` and `mint broken-links` all passed), and then restored. Match on the quoted text, since earlier edits shift lines.

### ENG-1413: `docs(cli): document the cli device authorization page`

File: `docs/developers/cli/authentication.mdx`, `## Interactive login`, after `<LoginInstructions />`.

**What's missing:** the browser side of `miru login`. The shared snippet ends at "Confirm that the code in the browser matches the code in the terminal." The browser opens **Authorize Device**, which shows the request and asks you to click **Allow**.

**Code to re-verify:**

- `src/app/oauth/device/page.tsx`: route `/oauth/device?user_code=…`; signed out, it redirects to `/login?user_code=…`, and sign-in (or sign-up and onboarding) returns there (`src/views/authn/shared/utils/destination.ts`).
- `src/views/authn/cli/components/AuthorizationHeader.tsx`: **Authorize Device**, "Enter the code from your terminal to authorize the Miru CLI."
- `…/cli/utils/code.ts`: eight letters in two groups, no vowels; pasted codes are normalized.
- `…/cli/components/RequestDetails.tsx`: **Device** (`<hostname> @ miru <version> <os> (<arch>)`), **Location** ("Location unavailable" if none), **Requested at**, **IP address**.
- `…/cli/components/AuthorizationActions.tsx`: **Allow**; "Do not click **Allow** unless you started this login from the Miru CLI."
- `…/cli/components/AuthorizationSuccess.tsx`: **Authorization successful**, "You can close this tab."
- `…/cli/utils/error.ts`; `backend internal/orgs/services/deviceauth/create.go:25` (30-minute codes) and `guesses.go:16-17` (10 wrong codes per 10 minutes).

**Stop and ask Armel (1 of 2):** "Say how long a code lasts? The backend keeps it 30 minutes; the page's own error says 'a few minutes'."

- Option A: add "Codes expire after 30 minutes. If the page says the code isn't valid, run `miru login` again for a new one." (the real behavior, M1 precedent). Then file the page's "a few minutes" copy for M9.
- Option B (recommended): leave expiry out. The page's error already says what to do, and the numbers disagree.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `cli/authorize-device.png` | Authorize in your browser | **add** (capture from a real `miru login`, with the request details showing) | The page's request details are something the reader must recognize before clicking **Allow**. The changelog's `changelog/26-09-23/miru-login.png` shows the same page if Armel prefers to reuse it. |

Wait for both answers and record them.

**Edit (Option B):** insert after `<LoginInstructions />` (a blank line before and after):

```mdx
### Authorize in your browser

The URL opens the **Authorize Device** page in Miru with your code filled in. If you
aren't logged in, log in (or sign up) first and you'll return to the page. It shows the
login request:

- **Device:** the machine's hostname, CLI version, OS, and architecture
- **Location** and **IP address:** where the request came from
- **Requested at:** when you ran `miru login`

Check that the request is yours, then click **Allow**. Don't click **Allow** unless you
started this login from the Miru CLI. When the page says **Authorization successful**,
you can close the tab.
```

(Option A: add the expiry paragraph at the end.)

**Checks and preview:** run the checks. Preview `http://localhost:3000/developers/cli/authentication#authorize-in-your-browser`; ideally run `miru login` against production and compare.

**Commit:** `docs/developers/cli/authentication.mdx`, `Refs ENG-1413`.

### ENG-1414: `docs(quick-start): document sign up, log in and workspace onboarding`

File: `docs/getting-started/quick-start/overview.mdx`, the sentence under the Steps.

**What's missing:** "If you don't already have an account, navigate to the sign-up page and follow the instructions." Sign-up is Google or name and email, an emailed code, then naming the workspace; invited users should use their invite link instead.

**Code to re-verify:**

- `src/views/authn/signup/SignUp.tsx`, `shared/AuthForm/*`, `shared/GoogleAuth/GoogleAuthButton.tsx`: **Get started with Miru**, **Continue with Google**, OR, **First name**, **Last name**, **Email**, **Continue**.
- `src/views/authn/verify/*`: **We sent you an email**, a six-box code field, **Resend code**.
- `src/views/authn/onboarding/*`: **Create your workspace**, "Give your workspace a name!", **Create**, then Devices.
- `src/views/authn/signin/SignIn.tsx`: **Login to Miru**, the same options.
- `src/views/authn/shared/utils/destination.ts` and `src/app/(authn)/onboarding/page.tsx`: invites skip onboarding.

**Stop and ask Armel (1 of 2):** "Where should sign-up live?"

- Option A (recommended): three steps in the quick start overview, where the sign-up link already is. It's the only place a new user meets it.
- Option B: a new "Sign up and log in" page under Getting started, linked from the quick start. More room, but one more page for three steps.

**Stop and ask Armel (2 of 2), screenshots:** I recommend:

| Image | Section | Verdict | Why |
|---|---|---|---|
| `getting-started/signup-page.png` | hero | **replace** with `getting-started/signup.png` | No **Continue with Google** button or OR divider. Captured in ENG-1415 batch F (signed out, private window). |
| Onboarding | — | **no image** | One field and **Create**; the text names them. |

Wait for both answers and record them.

**Edit (Option A):** replace "If you don't already have an account, navigate to the [sign-up](https://app.mirurobotics.com/signup) page and follow the instructions." with:

```mdx
If you don't already have an account, create one on the
[sign-up page](https://app.mirurobotics.com/signup):

1. Click **Continue with Google**, or enter your name and email and click **Continue**.
2. If you used your email, enter the code Miru emails you.
3. Name your workspace and click **Create**.

If you were [invited to a workspace](/admin/users/invites#receiving-invites), use the
link in your invite email instead. To log in later, use the
[login page](https://app.mirurobotics.com/login) the same way.
```

**Checks and preview:** run the checks. Preview `http://localhost:3000/getting-started/quick-start/overview`, and click the invites link.

**Commit:** `docs/getting-started/quick-start/overview.mdx`, `Refs ENG-1414`.

### ENG-1415: `docs: refresh v04 screenshots for redesigned pages`

Files: the pages in the inventory. Text-first: step 1 needs no captures; steps 2–7 follow Armel's capture sessions, one batch each.

**What's wrong:** 32 embeds show changed UI and 16 are dialogs the text names (inventory below). One more is blocked by a canceled text issue.

**Summary (158 in-scope embeds):**

| Batch | Area | Keep | Replace | Drop | Blocked |
|---|---|---|---|---|---|
| A | Settings and admin | 21 | 5 | 5 | 1 |
| B | Devices and provisioning | 21 | 4 | 0 | 0 |
| C | Config editor, device history, quick start deploy | 24 | 6 | 3 | 0 |
| D | Releases, staging and schemas | 21 | 13 | 5 | 0 |
| E | Config types | 2 | 2 | 3 | 0 |
| F | Sign-in and developer heroes | 3 | 2 | 0 | 0 |
| G | Data uploads | 6 | 0 | 0 | 0 |
| H | Groups | 11 | 0 | 0 | 0 |
| | **Total** | **109** | **32** | **16** | **1** |

Of the 32 replacements, 3 reuse an image M1–M7 already captured, so **29 need a new capture**, plus ENG-1413's `cli/authorize-device.png` if Armel adds it.

**Stop and ask Armel (1 of 4):** "How should ENG-1415 commit?"

- Option A (recommended): one commit per session (the text-first pass, then each batch), each with the issue title as the subject and `Refs ENG-1415`. Each commit is reviewable on its own and the milestone can close with any batch left.
- Option B: one commit at the end (the one-commit-per-issue precedent).

**Stop and ask Armel (2 of 4):** "Drop the 16 dialog and redundant shots (the M5/M7 pattern), across these pages too?" The list is every **drop** row below. Two of them (`config-instances/panel-metadata:deployed.png`, `panel-metadata:staged.png`) sit in `<Tabs>` blocks and go with the Diff-mode capture in batches C and D, not in step 1.

- Option A (recommended): yes, all 16.
- Option B: only the ones showing wrong UI (`invites/send-invite-dialog.png`, the two Metadata tabs, `releases/page.png` twice, `releases/stage/page.png`); keep the current-but-redundant dialogs.

**Stop and ask Armel (3 of 4):** "`users/members/dropdown.png` in 'Change a member's access' shows **Change role**, which the text still describes (ENG-1410 was canceled). What should happen to it?"

- Option A (recommended): leave it until the members text is revisited; it matches the text, even though the app moved on. Replace the other two uses (Suspend, Transfer ownership) with per-action menus in batch A.
- Option B: replace it with a **Manage access** menu shot and accept that the text says **Change role**.
- Option C: reopen ENG-1410's text change (out of M8 unless Armel adds it).

**Stop and ask Armel (4 of 4), captures:** "Which batches, in what order? Each needs one session." Setup per batch is in the inventory. Recommended order: F (sign-up and the CLI page, quick), A, E, B, C, D (the largest; needs staging data with staged, drifted and needs-review devices).

Wait for each answer and record it.

#### Step 1: text-first pass (no captures)

(a) Reuse swaps, one embed URL each:

- `developers/platform-api/overview.mdx`: `platform-api/header:api-keys.png` → `apikeys/hero.png` (M7's current API Keys page).
- `developers/agent/versions.mdx`: `devices/page.png` → `editor/init/devices-page.png` (M4; shows the **Agent** column the text points to).
- `provision-devices/reprovision.mdx`: `devices/ellipses-dropdown.png` → `devices/manage/edit-menu.png` (M2; the same menu, with **Reprovision**).

(b) Drop these frames (the whole `<Frame>…</Frame>` and the blank line after it; no text changes, since each section's text already names the dialog or links the page):

- `admin/users/invites.mdx`: `invites/send-invite-dialog.png`, `invites/revoke-invite-dialog.png`, `users/profile/accept-invite-dialog.png`
- `admin/users/profile.mdx`: `users/profile/crop-avatar-dialog.png`
- `admin/workspace.mdx`: `workspaces/crop-logo-dialog.png`
- `cfg-mgmt/concepts/config-types.mdx`: `config-types/create-dialog.png`, `edit-dialog.png`, `delete-dialog.png` (the file then ends at "…to confirm the deletion.", one trailing newline)
- `cfg-mgmt/concepts/schemas/manage.mdx`: `releases/page.png`
- `cfg-mgmt/deploy/staging-area.mdx`: `releases/page.png`, `releases/stage/page.png`, `releases/stage/deploy-dialog:bulk.png`
- `getting-started/quick-start/deploy-configs.mdx`: `devices/editor/set-release-dialog.png`, `devices/editor/init-deploy-dialog.png`

Run the checks. Preview each touched page. Commit (Option A): the ten files above, subject `docs: refresh v04 screenshots for redesigned pages`, body `Refs ENG-1415`.

#### Steps 2–7: one batch per session

For each batch Armel captures, in the order he chose: he uploads under the new names in the inventory (a `-v2` suffix where the notes give none), confirms each with `curl -I`, and the agent swaps the embed URLs. For the two config instance `<Tabs>` blocks (`cfg-mgmt/audit/device-history.mdx`, `cfg-mgmt/deploy/staging-area.mdx`), replace the whole block with one `<Frame>` of the Diff-mode sheet. Run the checks, preview, and commit per the answer to question 1. Tick the batch in Progress.

#### Inventory

Built from `main` at `e950a27` on 2026-09-29, one row per embed. "Accurate" is whether the image matches frontend `origin/prod`; "Replaced in" names the milestone that captured it under a new name. Line numbers are on `96595dd`.

##### Batch A: Settings and admin (32 embeds: 21 keep, 5 replace, 5 drop, 1 blocked)

Session setup: Settings → Profile, Workspace, Members (with a pending invite), API Keys. Use a mock with fake names and emails (M7 lesson: member rows show emails).

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `admin/apikeys.mdx:11` | (top) | `apikeys/hero.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/apikeys.mdx:30` | Create an API key | `apikeys/create-button-v2.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/apikeys.mdx:39` | Create an API key | `apikeys/create-page.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/apikeys.mdx:60` | View API keys | `apikeys/details.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/apikeys.mdx:71` | Edit an API key | `apikeys/edit-menu.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/apikeys.mdx:79` | Edit an API key | `apikeys/edit-page.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/apikeys.mdx:92` | Delete an API key | `apikeys/delete-menu.png` | yes | M7 | **keep** | Replaced in M7. |
| `admin/users/invites.mdx:14` | (top) | `invites/header:join-workspace.png` | yes | — | **keep** | The invite page (Join workspace / Decline invite). |
| `admin/users/invites.mdx:59` | Send an invite | `invites/send-invite-dialog.png` | no | — | **drop** | Shows a **Role** field; the dialog has **Type** now. A dialog the text names. The text still says "role" (ENG-1411 canceled). |
| `admin/users/invites.mdx:71` | Revoke an invite | `invites/pending-invites-button.png` | yes | — | **keep** | "1 pending" button. |
| `admin/users/invites.mdx:77` | Revoke an invite | `invites/pending-invite-dropdown.png` | yes | — | **keep** | Resend invite / Revoke invite menu. |
| `admin/users/invites.mdx:83` | Revoke an invite | `invites/revoke-invite-dialog.png` | yes | — | **drop** | Confirm dialog the text names. |
| `admin/users/invites.mdx:95` | Resend an invite | `invites/pending-invites-button.png` | yes | — | **keep** | "1 pending" button. |
| `admin/users/invites.mdx:101` | Resend an invite | `invites/pending-invite-dropdown.png` | yes | — | **keep** | Resend invite / Revoke invite menu. |
| `admin/users/invites.mdx:123` | Accept an invite | `users/profile/invitations.png` | yes | — | **keep** | Profile Invitations card (Join / Decline). |
| `admin/users/invites.mdx:129` | Accept an invite | `users/profile/accept-invite-dialog.png` | yes | — | **drop** | Confirm dialog; the text above it covers its warning. |
| `admin/users/invites.mdx:137` | Decline an invite | `users/profile/invitations.png` | yes | — | **keep** | Profile Invitations card (Join / Decline). |
| `admin/users/manage.mdx:11` | (top) | `users/members/header:page2.png` | no | — | **replace** | Single **Role** column; the table has **Type** and **Roles** now. New name: `users/members/hero.png`. |
| `admin/users/manage.mdx:22` | Change a member's access | `users/members/dropdown.png` | no | — | **blocked** | Shows **Change role**, which the text still describes (ENG-1410 canceled). A new shot would contradict the text; ask. |
| `admin/users/manage.mdx:51` | Suspend a member | `users/members/dropdown.png` | no | — | **replace** | Shows **Change role**. **Replace** with `users/members/suspend-menu.png`. |
| `admin/users/overview.mdx:9` | (top) | `users/members/header:page.png` | no | — | **replace** | Old settings sidebar (no API Keys or Buckets). Could reuse the members hero once it's captured. |
| `admin/users/profile.mdx:8` | (top) | `users/profile/header.png` | yes | — | **keep** | Profile and Preferences cards match. |
| `admin/users/profile.mdx:19` | Name | `users/profile/name-fields.png` | yes | — | **keep** |  |
| `admin/users/profile.mdx:27` | Avatar | `users/profile/avatar-field.png` | yes | — | **keep** | Shows the hover tooltip; the text still says hover (ENG-1412 canceled). |
| `admin/users/profile.mdx:33` | Avatar | `users/profile/crop-avatar-dialog.png` | yes | — | **drop** | Crop dialog the text names (**Save**). Low priority. |
| `admin/users/profile.mdx:45` | Preferences | `users/profile/preferences.png` | yes | — | **keep** |  |
| `admin/workspace.mdx:11` | (top) | `workspaces/header:page.png` | no | — | **replace** | Old settings sidebar and no **Editor** section. New name: `workspaces/hero.png`. |
| `admin/workspace.mdx:39` | Name | `workspaces/name-field.png` | yes | — | **keep** |  |
| `admin/workspace.mdx:49` | Logo | `workspaces/logo-field.png` | yes | — | **keep** | Hover tooltip; the text is unchanged (ENG-1412 canceled). |
| `admin/workspace.mdx:55` | Logo | `workspaces/crop-logo-dialog.png` | yes | — | **drop** | Crop dialog the text names (**Save**). Low priority. |
| `admin/workspace.mdx:65` | Transfer ownership | `users/members/dropdown.png` | no | — | **replace** | Shows **Change role**. **Replace** with `users/members/transfer-menu.png`. |
| `snippets/workspaces/leave.mdx:15` (in `admin/users/profile.mdx`, `admin/workspace.mdx`) | (top) | `users/profile/workspace-access.png` | yes | — | **keep** | Workspace access / Leave workspace (M1 snippet). |

##### Batch B: Devices and provisioning (25 embeds: 21 keep, 4 replace)

Session setup: Devices list and a device page on staging, one unprovisioned device, the dev mock for rows.

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `concepts/devices/manage.mdx:19` | View a device | `devices/manage/overview.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/manage.mdx:47` | Edit a device | `devices/manage/edit-menu.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/manage.mdx:61` | Ping a device | `devices/ping-success.png` | yes | — | **keep** | Toast (checked in M2). |
| `concepts/devices/manage.mdx:67` | Ping a device | `devices/ping-timeout.png` | yes | — | **keep** | Toast (checked in M2). |
| `concepts/devices/manage.mdx:119` | Delete a device | `devices/manage/delete-menu.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/overview.mdx:11` | (top) | `devices/overview.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/overview.mdx:123` | Status | `devices/status-hover.png` | no | — | **replace** | Old list header; the tooltip itself is current. Low priority. |
| `concepts/devices/views.mdx:20` | Search devices | `devices/views/search.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/views.mdx:31` | Filter devices | `devices/views/filter-menu.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/views.mdx:48` | Display options | `devices/views/display-menu.png` | yes | M2 | **keep** | Replaced in M2. |
| `concepts/devices/views.mdx:61` | Saved views | `devices/views/saved-views.png` | yes | M2 | **keep** | Replaced in M2. |
| `developers/agent/versions.mdx:59` | Check the version | `devices/page.png` | no | — | **replace** | Old sidebar with **Tags** and **+ New Device**. **Reuse** `editor/init/devices-page.png` (M4), which shows the Agent column. |
| `getting-started/quick-start/provision-device.mdx:32` | Install the `miru-agent` package | `devices/device-to-provision.png` | yes | — | **keep** | Row with the **Provision** button. |
| `provision-devices/dashboard.mdx:14` | (top) | `devices/provisioning/header:provision-dialog.png` | yes | — | **keep** | Same dialog as `provision-dialog-v2.png` (M2). |
| `provision-devices/dashboard.mdx:32` | Install the `miru-agent` package | `devices/device-to-provision.png` | yes | — | **keep** | Row with the **Provision** button. |
| `provision-devices/overview.mdx:7` | (top) | `devices/header:create-panel.png` | no | — | **replace** | **Create Device** panel with **Tags**; tags are gone. New name: `devices/create-hero.png`, or reuse `create-page-v2.png` (M2). |
| `provision-devices/provisioning-tokens.mdx:13` | (top) | `devices/provisioning/header:create-token.png` | n/a | — | **keep** | Code sample, not app UI. |
| `provision-devices/reprovision.mdx:12` | (top) | `devices/provisioning/header:reprovision-dialog.png` | yes | — | **keep** | Same style as the M2 dialogs. |
| `provision-devices/reprovision.mdx:50` | Dashboard | `devices/ellipses-dropdown.png` | no | — | **replace** | Old list columns; the menu is current. **Reuse** `devices/manage/edit-menu.png` (M2), which shows **Reprovision**. |
| `snippets/devices/create.mdx:4` (in `getting-started/quick-start/provision-device.mdx`, `provision-devices/dashboard.mdx`) | (top) | `devices/create-page-v2.png` | yes | M2 | **keep** | Replaced in M2. |
| `snippets/devices/move.mdx:4` (in `concepts/devices/manage.mdx`) | (top) | `devices/manage/move-menu.png` | yes | M2 | **keep** | Replaced in M2. |
| `snippets/devices/move.mdx:10` (in `concepts/devices/manage.mdx`) | (top) | `devices/manage/move-dialog.png` | yes | M2 | **keep** | Replaced in M2. |
| `snippets/devices/provision/install-dialog.mdx:4` (in `getting-started/quick-start/provision-device.mdx`, `provision-devices/dashboard.mdx`, `provision-devices/reprovision.mdx`) | (top) | `devices/provisioning/install-dialog-v2.png` | yes | M2 | **keep** | Replaced in M2. |
| `snippets/devices/provision/provision-dialog.mdx:6` (in `getting-started/quick-start/provision-device.mdx`, `provision-devices/dashboard.mdx`) | (top) | `devices/provisioning/provision-dialog-v2.png` | yes | M2 | **keep** | Replaced in M2. |
| `snippets/devices/provision/reprovision-dialog.mdx:6` (in `provision-devices/reprovision.mdx`) | (top) | `devices/provisioning/reprovision-dialog.png` | yes | — | **keep** |  |

##### Batch C: Config editor, device history and quick start deploy (33 embeds: 24 keep, 6 replace, 3 drop)

Session setup: A device with at least two deployments on different releases and a parent (M5 left this open: miru-03 upgrade, then a tune), the config instance sheet in Diff mode.

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `cfg-mgmt/audit/compare-devices.mdx:9` | (top) | `devices/compare/hero.png` | yes | — | **keep** | Checked in M5. |
| `cfg-mgmt/audit/compare-devices.mdx:38` | Selection | `devices/compare/selection.png` | yes | — | **keep** | Checked in M5. |
| `cfg-mgmt/audit/compare-devices.mdx:75` | Summary | `devices/compare/summary.png` | yes | — | **keep** | Checked in M5. |
| `cfg-mgmt/audit/compare-devices.mdx:96` | Files | `devices/compare/files.png` | yes | — | **keep** | Checked in M5. |
| `cfg-mgmt/audit/device-history.mdx:7` | (top) | `devices/deployments/header:page.png` | no | — | **replace** | Flat list with a Status column; the list groups by release now (M5, left open). New name: `devices/deployments/hero.png`. |
| `cfg-mgmt/audit/device-history.mdx:22` | View a deployment | `deployments/deployed-panel.png` | yes | — | **keep** | Panel header and fields match (M5: keep or replace). |
| `cfg-mgmt/audit/device-history.mdx:32` | View a config instance | `deployments/deployed-panel:configurations.png` | no | — | **replace** | Shows **MODIFIED** where the panel now shows change counts (M5, left open). |
| `cfg-mgmt/audit/device-history.mdx:38` | View a config instance | `config-instances/panel-content:deployed.png` | no | — | **replace** | Content / Metadata tabs are gone; the sheet opens in **Diff** with **Read-only** (M5). One Diff-mode shot, and drop the `<Tabs>` block. |
| `cfg-mgmt/audit/device-history.mdx:43` | View a config instance | `config-instances/panel-metadata:deployed.png` | no | — | **drop** | The Metadata tab no longer exists. |
| `cfg-mgmt/concepts/config-instances.mdx:12` | (top) | `config-instances/header:panel.png` | no | — | **replace** | Content / Metadata tabs. New name: `config-instances/hero.png` (the sheet in Diff mode). |
| `cfg-mgmt/concepts/schemas/validation.mdx:73` | Example | `editor/validation/invalid-edit.png` | yes | M4 | **keep** | Replaced in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:10` | (top) | `editor/hero.png` | yes | M4 | **keep** | Replaced in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:57` | Change list | `editor/change-list.png` | yes | M4 | **keep** | Replaced in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:71` | Instance slots | `editor/add-slot-menu.png` | yes | M4 | **keep** | Replaced in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:88` | Deploying changes | `devices/editor/deploy-dialog.png` | yes | — | **keep** | Checked in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:138` | Deployment list | `devices/editor/history-list.png` | yes | — | **keep** | Checked in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:148` | Deployment list | `devices/editor/history-context-menu.png` | yes | M1 | **keep** | Replaced in M1. |
| `cfg-mgmt/deploy/config-editor.mdx:164` | Diff and read-only modes | `devices/editor/view-modes.png` | yes | — | **keep** | Checked in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:183` | Diff navigation | `devices/editor/diff-navigation.png` | yes | — | **keep** | Checked in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:199` | Base vs. Head | `devices/editor/base-head.png` | yes | — | **keep** | Checked in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:215` | Redeploy | `devices/editor/redeploy.png` | yes | — | **keep** | Checked in M4. |
| `cfg-mgmt/deploy/config-editor.mdx:229` | Deployment alert | `devices/editor/depl-alert.png` | yes | M1 | **keep** | Replaced in M1. |
| `cfg-mgmt/deploy/config-editor.mdx:239` | Editor settings | `devices/editor/settings.png` | yes | — | **keep** | Matches the workspace Editor section. |
| `cfg-mgmt/deploy/config-editor.mdx:257` | Editor settings | `devices/editor/format-action.png` | yes | — | **keep** |  |
| `cfg-mgmt/deploy/initial-deployment.mdx:18` | Select a device | `editor/init/devices-page.png` | yes | M4 | **keep** | Replaced in M4. |
| `cfg-mgmt/deploy/initial-deployment.mdx:42` | Edit configs | `editor/init/first-draft.png` | yes | M4 | **keep** | Replaced in M4. |
| `cfg-mgmt/deploy/initial-deployment.mdx:77` | Deploy | `editor/init/history.png` | yes | M4 | **keep** | Replaced in M4. |
| `getting-started/quick-start/deploy-configs.mdx:10` | (top) | `devices/provisioned-device.png` | yes | — | **keep** | A row crop (Device / Release / Agent); close enough. |
| `getting-started/quick-start/deploy-configs.mdx:16` | (top) | `devices/editor/set-release-dialog.png` | yes | — | **drop** | The dialog the text names; M4 dropped the same shot from initial-deployment. |
| `getting-started/quick-start/deploy-configs.mdx:22` | (top) | `devices/editor/init-deployment-draft.png` | no | — | **replace** | Pre-M4 editor (no file tree, old change panel). Capture on the quick start device (`getting-started/first-draft.png`); M4's `editor/init/first-draft.png` shows other configs. M4 left these for later ("text first, images later", ENG-1404). |
| `getting-started/quick-start/deploy-configs.mdx:32` | (top) | `devices/editor/init-deploy-dialog.png` | yes | — | **drop** | Deploy dialog the text names; the second deploy keeps `deploy-dialog.png`. |
| `getting-started/quick-start/deploy-configs.mdx:52` | Update the configs | `devices/editor/modifications.png` | no | — | **replace** | Pre-M4 editor. Capture on the quick start device (`getting-started/modifications.png`). |
| `getting-started/quick-start/deploy-configs.mdx:60` | Update the configs | `devices/editor/deploy-dialog.png` | yes | — | **keep** | Checked in M4. |

##### Batch D: Releases, staging and schemas (39 embeds: 21 keep, 13 replace, 5 drop)

Session setup: A release with staged, drifted and needs-review devices (the v7.0.0 staging data), the release editor, a deployment panel with ADDED and change-count badges.

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `cfg-mgmt/concepts/schemas/manage.mdx:9` | (top) | `config-schemas/header:metadata.png` | yes | — | **keep** | The schema sheet still has Schema / Metadata tabs (`SchemaSheet.tsx`). |
| `cfg-mgmt/concepts/schemas/manage.mdx:24` | View a schema | `releases/page.png` | no | — | **drop** | Old sidebar and list; the text links the Releases page (M5 dropped it from file rules the same way). |
| `cfg-mgmt/concepts/schemas/manage.mdx:30` | View a schema | `releases/schema-list.png` | yes | — | **keep** | Schemas section of the Overview tab. |
| `cfg-mgmt/concepts/schemas/manage.mdx:41` | View a schema | `config-schemas/panel-content.png` | yes | — | **keep** | Schema tab; a document picker now appears for multi-document schemas. |
| `cfg-mgmt/concepts/schemas/manage.mdx:46` | View a schema | `config-schemas/panel-metadata.png` | yes | — | **keep** |  |
| `cfg-mgmt/concepts/schemas/overview.mdx:10` | (top) | `config-schemas/header:panel.png` | yes | — | **keep** |  |
| `cfg-mgmt/create-a-release.mdx:15` | (top) | `releases/header:release-create.png` | n/a | — | **keep** | Terminal output of `miru release create`, not app UI. |
| `cfg-mgmt/deploy/staging-area.mdx:10` | (top) | `releases/stage/header:page.png` | no | — | **replace** | Shows a **Deployments** tab; a release has Overview and Stage only (M5). New name: `releases/stage/hero.png`. |
| `cfg-mgmt/deploy/staging-area.mdx:35` | View the staging area | `releases/page.png` | no | — | **drop** | Old sidebar and list; the text links the Releases page (M5 dropped it from file rules the same way). |
| `cfg-mgmt/deploy/staging-area.mdx:41` | View the staging area | `releases/stage/page.png` | no | — | **drop** | Old sidebar with **Tags** and a Deployments tab; the hero shows the same page once replaced. |
| `cfg-mgmt/deploy/staging-area.mdx:60` | Stage a deployment | `releases/stage/empty-staging.png` | no | — | **replace** | Old sidebar; the empty state and **+ Stage Deployment** are current. |
| `cfg-mgmt/deploy/staging-area.mdx:67` | Stage a deployment | `releases/stage/select-device.png` | yes | — | **keep** |  |
| `cfg-mgmt/deploy/staging-area.mdx:81` | Stage a deployment | `releases/stage/release-editor-tooltip.png` | no | — | **replace** | Pre-M4 editor layout (no file tree). Confirm at capture that the release editor uses the rebuilt editor. |
| `cfg-mgmt/deploy/staging-area.mdx:88` | Stage a deployment | `releases/stage/stage-button.png` | yes | — | **keep** |  |
| `cfg-mgmt/deploy/staging-area.mdx:96` | Stage a deployment | `releases/stage/dialog-stage.png` | yes | M1 | **keep** | Replaced in M1. |
| `cfg-mgmt/deploy/staging-area.mdx:103` | Stage a deployment | `releases/stage/staged-deployment.png` | yes | — | **keep** | Overview \| Stage tabs, current columns. |
| `cfg-mgmt/deploy/staging-area.mdx:119` | View a deployment | `deployments/staged-panel.png` | yes | — | **keep** |  |
| `cfg-mgmt/deploy/staging-area.mdx:126` | View a deployment | `deployments/staged-panel:configurations.png` | no | — | **replace** | Badges say **NEW**; the panel says **ADDED** now (`sheet/utils/badge.ts`). |
| `cfg-mgmt/deploy/staging-area.mdx:135` | View a deployment | `config-instances/panel-content:staged.png` | no | — | **replace** | Content / Metadata tabs are gone (Diff / Read-only). |
| `cfg-mgmt/deploy/staging-area.mdx:140` | View a deployment | `config-instances/panel-metadata:staged.png` | no | — | **drop** | The Metadata tab no longer exists. |
| `cfg-mgmt/deploy/staging-area.mdx:160` | Patch a deployment | `releases/stage/patch-dropdown.png` | yes | — | **keep** | Deploy / Archive / Patch / Review (`staging/utils/actions.ts`). |
| `cfg-mgmt/deploy/staging-area.mdx:167` | Patch a deployment | `releases/stage/release-editor-patch.png` | no | — | **replace** | Pre-M4 editor layout. |
| `cfg-mgmt/deploy/staging-area.mdx:240` | Deploy a deployment | `releases/stage/select-devices.png` | no | — | **replace** | Old sidebar with **Tags** and a Deployments tab. Used in Deploy (and Archive, out of scope). |
| `cfg-mgmt/deploy/staging-area.mdx:246` | Deploy a deployment | `releases/stage/bulk-actions.png` | yes | — | **keep** | Deploy / Archive bulk bar. |
| `cfg-mgmt/deploy/staging-area.mdx:252` | Deploy a deployment | `releases/stage/deploy-dialog:bulk.png` | yes | — | **drop** | Confirm dialog; the text covers its sentence. |
| `cfg-mgmt/deploy/staging-area.mdx:262` | Deployment drift | `deployments/deployment-drift.png` | n/a | — | **keep** | Concept diagram, not app UI. |
| `cfg-mgmt/deploy/staging-area.mdx:295` | Review a deployment | `releases/stage/ellipses-dropdown:drifted.png` | no | — | **replace** | Old sidebar with **Tags** and a Deployments tab; the menu (Patch / Review / Archive) is current. |
| `cfg-mgmt/deploy/staging-area.mdx:302` | Review a deployment | `releases/stage/review-dialog:deployments.png` | yes | — | **keep** | Review dialog list. |
| `cfg-mgmt/deploy/staging-area.mdx:308` | Review a deployment | `releases/stage/deployment-panel:deployed.png` | no | — | **replace** | Configurations show **MODIFIED** where counts show now. Low priority. |
| `cfg-mgmt/deploy/staging-area.mdx:315` | Review a deployment | `releases/stage/config-panel.png` | no | — | **replace** | Content / Metadata tabs are gone. |
| `cfg-mgmt/deploy/staging-area.mdx:329` | Review a deployment | `releases/stage/review-dialog:restage.png` | yes | — | **keep** | **Restage** footer (`ReviewFooter.tsx`). |
| `cfg-mgmt/deploy/staging-area.mdx:340` | Review a deployment | `releases/stage/review-dialog:archive.png` | yes | — | **keep** | Restage menu with Archive (Review section). |
| `concepts/deployments/overview.mdx:11` | (top) | `deployments/header:panel.png` | no | — | **replace** | Badges say **NEW**. New name: `deployments/hero.png`. |
| `concepts/releases.mdx:13` | (top) | `releases/hero.png` | yes | M5 | **keep** | Replaced in M5. |
| `concepts/releases.mdx:56` | View a release | `releases/overview-tab.png` | yes | M5 | **keep** | Replaced in M5. |
| `concepts/releases.mdx:75` | Duplicate a release | `releases/duplicate-menu.png` | yes | M5 | **keep** | Replaced in M5. |
| `concepts/releases.mdx:124` | Delete a release | `releases/delete-menu.png` | yes | M5 | **keep** | Replaced in M5. |
| `getting-started/quick-start/create-release.mdx:146` | Create a release | `releases/page.png` | no | — | **replace** | Old sidebar and list. Capture the quick start's Releases page (`getting-started/releases-page.png`), or reuse `releases/hero.png` (M5) if other example versions are fine. |
| `getting-started/quick-start/create-release.mdx:152` | Create a release | `releases/overview.png` | no | — | **replace** | No File Rules section, old sidebar. Capture the quick start's v1.0.0 Overview (`getting-started/release-overview.png`); `releases/overview-tab.png` (M5) shows other example data. |

##### Batch E: Config types (7 embeds: 2 keep, 2 replace, 3 drop)

Session setup: The Configs page with three or more config types.

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `cfg-mgmt/concepts/config-types.mdx:11` | (top) | `config-types/header:page.png` | no | — | **replace** | Old sidebar with **Tags**. New name: `config-types/hero.png`. |
| `cfg-mgmt/concepts/config-types.mdx:54` | Create a config type | `config-types/create-page.png` | no | — | **replace** | Old sidebar; it locates the **+ Config** button, so a crop of the header is enough (`config-types/create-button.png`). |
| `cfg-mgmt/concepts/config-types.mdx:60` | Create a config type | `config-types/create-dialog.png` | not checked | — | **drop** | Name / Slug dialog the text names. |
| `cfg-mgmt/concepts/config-types.mdx:74` | Edit a config type | `config-types/ellipses-dropdown.png` | yes | — | **keep** | Edit / Archive / Delete. One shot per action is optional (Edit, Delete; Archive is out of scope). |
| `cfg-mgmt/concepts/config-types.mdx:80` | Edit a config type | `config-types/edit-dialog.png` | not checked | — | **drop** | Dialog the text names. |
| `cfg-mgmt/concepts/config-types.mdx:129` | Delete a config type | `config-types/ellipses-dropdown.png` | yes | — | **keep** | Edit / Archive / Delete. One shot per action is optional (Edit, Delete; Archive is out of scope). |
| `cfg-mgmt/concepts/config-types.mdx:135` | Delete a config type | `config-types/delete-dialog.png` | yes | — | **drop** | Confirm dialog. |

##### Batch F: Sign-in and developer heroes (5 embeds: 3 keep, 2 replace)

Session setup: Signed out in a private window for `/signup`; a real `miru login` for the Authorize Device page (ENG-1413).

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `developers/ci/gh-actions.mdx:9` | (top) | `releases/header:gh-actions.png` | n/a | — | **keep** | GitHub Actions UI, not Miru. |
| `developers/ci/overview.mdx:10` | (top) | `releases/header:release-create.png` | n/a | — | **keep** | Terminal output of `miru release create`, not app UI. |
| `developers/cli/overview.mdx:10` | (top) | `releases/header:release-create.png` | n/a | — | **keep** | Terminal output of `miru release create`, not app UI. |
| `developers/platform-api/overview.mdx:10` | (top) | `platform-api/header:api-keys.png` | no | — | **replace** | Old API Keys page with CLI Token and Webhooks items. **Reuse** `apikeys/hero.png` (M7); no capture. |
| `getting-started/quick-start/overview.mdx:13` | (top) | `getting-started/signup-page.png` | no | — | **replace** | No **Continue with Google** button or OR divider (ENG-1414). New name: `getting-started/signup.png`. |

##### Batch G: Data uploads (6 embeds: 6 keep)

Session setup: Nothing to capture.

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `data-uploads/concepts/buckets.mdx:49` | View a bucket | `buckets/list.png` | yes | M7 | **keep** | Replaced in M7. |
| `data-uploads/concepts/buckets.mdx:79` | Edit a bucket | `buckets/edit-menu.png` | yes | M7 | **keep** | Replaced in M7. |
| `data-uploads/concepts/buckets.mdx:97` | Delete a bucket | `buckets/delete-menu.png` | yes | M7 | **keep** | Replaced in M7. |
| `data-uploads/concepts/file-rules/overview.mdx:109` | View a file rule | `file-rules/details.png` | yes | M5 | **keep** | Replaced in M5. |
| `data-uploads/connect-a-bucket/aws.mdx:90` | Connect a bucket | `buckets/add-bucket-aws-ext-id.png` | yes | — | **keep** | Checked in M7. |
| `data-uploads/connect-a-bucket/gcs.mdx:67` | Connect a bucket | `buckets/add-bucket-gcp-sa.png` | yes | — | **keep** | Checked in M7. |

##### Batch H: Groups (11 embeds: 11 keep)

Session setup: Nothing to capture.

| Page | Section | Image | Accurate | Replaced in | Action | Notes |
|---|---|---|---|---|---|---|
| `concepts/groups/manage.mdx:21` | View a group | `groups/manage/page.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/manage.mdx:44` | Create a group | `groups/manage/create-group.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/manage.mdx:57` | Rename a group | `groups/manage/rename-menu.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/manage.mdx:70` | Move a group | `groups/manage/move-menu.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/manage.mdx:78` | Move a group | `groups/manage/move-dialog.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/manage.mdx:93` | Delete a group | `groups/manage/delete-menu.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/members.mdx:20` | Add members | `groups/members/members-card.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/members.mdx:34` | Edit a member's access | `groups/members/edit-menu.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/members.mdx:45` | Edit a member's access | `groups/members/access-panel.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/members.mdx:56` | Remove a member | `groups/members/remove-menu.png` | yes | M3 | **keep** | Replaced in M3. |
| `concepts/groups/overview.mdx:14` | (top) | `groups/groups-overview.png` | yes | M1 | **keep** | Replaced in M1. |

Excluded from the sweep (78 embeds): 28 decorative backgrounds, 31 SVG diagrams, 8 third-party shots (GitHub release assets, the CUE and JSON Schema websites, `systemctl status` output), and 11 embeds in archiving-only sections:

- `cfg-mgmt/concepts/config-types.mdx`: `config-types/archive-dialog.png`
- `cfg-mgmt/concepts/config-types.mdx`: `config-types/ellipses-dropdown.png`
- `cfg-mgmt/deploy/staging-area.mdx`: `releases/stage/archive-dialog:bulk.png`
- `cfg-mgmt/deploy/staging-area.mdx`: `releases/stage/archive-dialog:single.png`
- `cfg-mgmt/deploy/staging-area.mdx`: `releases/stage/bulk-actions.png`
- `cfg-mgmt/deploy/staging-area.mdx`: `releases/stage/ellipses-dropdown:staged.png`
- `cfg-mgmt/deploy/staging-area.mdx`: `releases/stage/select-devices.png`
- `concepts/devices/manage.mdx`: `devices/manage/archive-menu.png`
- `concepts/releases.mdx`: `releases/archive-menu.png`
- `concepts/releases.mdx`: `releases/unarchive-menu.png`
- `data-uploads/concepts/buckets.mdx`: `buckets/archive-menu.png`

### Follow-up: ENG-1427 and the Duplicate a release text

Not an M8 issue; no text change now.

**What's happening:** `concepts/releases.mdx` → Duplicate a release says the new release gets "the same config schemas, file rules, and git information". Armel chose that at ENG-1399 (M5) so the docs and ENG-1427 would ship together, but the docs shipped first (#216), and ENG-1427 (M9, "copy file rules when duplicating a release") is still Backlog. Until it ships, the live docs promise file rules that Duplicate doesn't copy.

**Steps:**

1. Before the M8 PR, check ENG-1427 in Linear and tell Armel its status in the final review.
2. When ENG-1427 reaches frontend `origin/prod`: confirm `buildCreateReleaseRequest` sends `file_rule_ids` (`src/features/releases/components/actions/duplicate/api/mutation.ts`), and duplicate a release with file rules on staging to see them on the new release's **File Rules**. No docs change is needed if it works.
3. **Stop and ask Armel** if ENG-1427 is dropped or slips past the M8 PR: revert the sentence to "…the same config schemas and git information… File rules and deployments aren't copied." (the M5 draft), or leave it.

## Concrete Steps

All commands run from `~/dev/miru/docs`.

0. Start state:

        git branch --show-current          # expect: docs/m8-sign-in-sweep
        git log --oneline -2               # expect: 96595dd docs(plans): move the m7 plan to completed, on top of e950a27 (#218)
        git status --short                 # expect: only this plan (untracked)
        curl -s https://api.mirurobotics.com/frontend/v1/version   # expect git_commit 80acefec…

   Commit this plan on its own **before** ENG-1413 (M5 precedent): `git add plans/active/20260929-m8-sign-in-sweep.md`, subject `docs: add the m8 sign-in and sweep plan`, no `Refs`. Fold later plan updates into that commit (`git commit --fixup <sha>` then `git rebase -i --autosquash origin/main`).

1. For each issue, in the order above: explain it, re-verify, stop at each **Stop and ask Armel** step and wait, edit, run the checks one at a time, stop for Armel's preview and review, iterate, then Armel commits.

2. After the last ENG-1415 batch, do the ENG-1427 follow-up, then update Progress, the Decision Log and Outcomes in this plan, and fold that update into the plan commit.

3. After the last commit:

        git log --oneline origin/main..HEAD   # expect the M7 plan move, this plan, ENG-1413, ENG-1414, and the ENG-1415 commits
        git diff --stat origin/main..HEAD     # expect authentication.mdx, the quick start overview, the ENG-1415 pages, this plan, and the M7 plan in plans/completed/

   Then hand back to the Project for the final review and PR. Don't push from this session.

## Validation and Acceptance

1. `/developers/cli/authentication`: "Authorize in your browser" names **Authorize Device**, the request details, **Allow** and its warning, and **Authorization successful**. The CLI reference page (`/references/cli/login`) is unchanged.
2. `/getting-started/quick-start/overview`: the three sign-up steps, the invite link and the login page.
3. `grep -rn "follow the instructions" docs/getting-started` prints nothing.
4. Every inventory row marked **drop** is gone (except any Armel kept), and every **replace** row points at a new name that returns 200 with `curl -I`.
5. `git diff origin/main..HEAD -- docs/ | grep -E '^[-+].*(<Frame|!\[|image=)'` shows only the swaps and drops in the inventory, plus images Armel added.
6. No archiving content was added: `git diff origin/main..HEAD -- docs/ | grep -i -E '^\+.*archiv'` prints nothing.
7. `git diff origin/main..HEAD --stat -- docs/snippets/references/cli docs/references` prints nothing.
8. Every **Stop and ask Armel** answer is in the Decision Log.
9. `./scripts/lint.sh`, `pnpm validate` and `mint broken-links` pass on the final head. Each issue commit's subject is its Linear title, and its body has `Refs`.

## Idempotence and Recovery

- Every edit is a text replacement, insertion or frame removal. Before reapplying, check whether the new text is already there (`grep`) and skip it if so.
- Before a commit, `git restore <file>` undoes an edit, and `git restore --staged <file>` unstages.
- To change an earlier commit after later ones exist: `git commit --fixup <sha>` then `git rebase -i --autosquash origin/main`. That's safe because nothing is pushed. Never force-push.
- The frontend and backend repos are read-only. The M9 agent works in the frontend checkout: don't switch its branch or touch its files, and ask before using the dev mock for captures.

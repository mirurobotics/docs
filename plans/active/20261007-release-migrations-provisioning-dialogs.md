# Document release migrations and refresh the provisioning dialog screenshots

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `mirurobotics/docs` (this worktree) | read-write | MDX edits to the staging-area page and the device provisioning snippets and pages. |
| `mirurobotics/frontend` (scratch clone) | read-only, local harness edits never committed | Source of truth for dialog wording; run locally to capture screenshots. |
| `mirurobotics/backend`, `mirurobotics/core` | read-only | Source of truth for what migration does. |

This plan lives in the docs repo because that is the only repo that receives commits. Paths below are relative to the docs worktree root unless absolute. Shell state does not persist between agent Bash calls, so begin every command block with:

    export S=/tmp/claude-1000/-home-ben-miru-workbench1/73a41985-25e7-4028-81c6-73f4f77ec108/scratchpad
    export DOCS=$S/docs-task H=http://localhost:3334/harness I=$S/docs-task/docs/images

The worktree is on branch `docs/release-migrations-provisioning-dialogs`, whose PR targets `docs/windows-agent-support` (the head of docs PR #208), not `main`. Never touch the main checkout at `/home/ben/miru/workbench1/repos/docs`.

## Purpose / Big Picture

After this change, a reader of the docs can:

1. Learn what the **Migrate configs** dialog does when they stage a device onto a new release. The dialog is described in a new "Migrate configs to a new release" section of `docs/cfg-mgmt/deploy/staging-area.mdx`, and the "Stage a deployment" procedure links to it. The changelog link `/cfg-mgmt/deploy/staging-area#stage-a-deployment` keeps working.
2. See the current provisioning dialog on every page that shows it: one dialog with **Linux | Windows** tabs. The Linux tab has apt install steps and a bash command. The Windows tab has an installer card and a PowerShell command captioned "Run as administrator". The reprovision dialog shows only the command. The old two-tab "Install agent | Provision device" screenshots are gone.

You can see it working in a local `mint dev` preview on port 3336.

## Progress

- [x] Milestone 0: frontend harness set up and screenshots captured (2026-10-07). Fresh `$S/frontend` clone of `origin/main` at `b6b5f814`; seven PNGs staged under `docs/images/` (gitignored); dev server on `:3334` stopped.
- [x] Milestone 1: release migration docs written and committed (2026-10-07, `f5c5583`).
- [x] Milestone 2: provisioning snippets and pages updated and committed (2026-10-07, `8f69e0f`).
- [x] Refine pass (2026-10-07, `e9e7220`): fixed the nonexistent "click **Reprovision**" instruction, named **Copy command** in every tab, said which tab the reprovision dialog opens on, and named the source release in the "Update unchanged defaults" bullet.
- [x] Milestone 3 local validation (2026-10-07): lint passed, `mint validate` passed, render check on `:3336` passed for all four pages (all seven images load, tabs switch, new TOC entry, in-page link scrolls, no raw `<version>`); URLs restored and mint stopped.
- [ ] Milestone 3 delivery: branch pushed, draft PR opened, preflight `CLEAN`.

## Surprises & Discoveries

- The earlier scratch frontend clone and `shot.mjs` were gone, so both were rebuilt from scratch (`$S/frontend`, `$S/shot/shot.mjs`). Frontend `origin/main` was still `b6b5f814`, so the copy quoted in Context and Orientation was used unchanged.
- `shot.mjs` as sketched never finished: `waitForFunction(... every animation not running)` timed out because the device status dot has an infinite `pulse` animation. Fix: before that wait, pause every animation with `iterations === Infinity` and set its `currentTime = 0` (first frame, dot fully opaque).
- The broad `[role=dialog] *{overflow:visible}` override rendered identically to a narrowed `[role=dialog] [data-slot=dialog-body]{overflow:visible}` override (pixel diff limited to sub-pixel noise). The narrowed override was adopted for the final captures because it leaves the code blocks' `overflow-x-auto` intact.
- The folded apt block shows about three lines plus a fade and **Show more**, not six: `CodeBlock` collapses to `max-h-24` with a gradient mask whenever the code has more than `collapsedLines` (6) lines. This matches the real app.
- `.env.local` also needed `NEXT_PUBLIC_SUPABASE_ANON_KEY`; a fake placeholder was used.
- The refine pass found that the reprovision dialog has no **Reprovision** button (it opens from the ellipsis menu and offers only **Copy token** / **Copy command**), so the plan's suggested reprovision.mdx sentence was replaced. Removing `<Install />` from reprovision.mdx also dropped the only note on which OS tab to pick, so the reprovision snippet now says which tab the dialog opens on.
- Captured sizes: provision dialogs 1152px wide (Linux 1152x1200, Windows 1152x820), reprovision dialogs 1152x648 (Linux) and 1152x688 (Windows), headers 1664x1204 (provision) and 1664x1032 (reprovision), migration dialog 1024x920.

## Decision Log

- Decision: Follow the migration dialog copy on frontend `main` (PR #95, "parameters" wording) for both the docs text and the screenshot; do not use the pre-#95 "settings" fallback from `dd2d65ab`.
  Rationale: The docs follow the shipped code; "parameters" also matches the changelog and the rest of the docs.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: The error card's button is documented as **Try again**.
  Rationale: That is the `ErrorCard` default `retryText` used by `MigrationBoundary`; "Retry" in the original request was wrong.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: Quote the remove option's UI label exactly (**Remove parameters not in v1.6's schema**) but describe what the backend actually does: it removes every parameter with no default in the new release, including free-form entries and declared parameters without a default.
  Rationale: The PR #95 caption ("Declared parameters are kept, even without a default") contradicts `remove_unknown_fields.go`; the caption mismatch is a frontend follow-up.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: Re-clone frontend into `$S/frontend` from `origin/main` and rebuild the harness and `$S/shot/shot.mjs`, ignoring the unrelated `$S/fe-src`.
  Rationale: The previous scratch clone and script no longer existed.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: Replace both page-header images with new `header:provision-dialog-v2.png` and `header:reprovision-dialog-v2.png` captures.
  Rationale: The old headers show the removed two-tab "Install agent | Provision device" dialog.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: Never open the provisioning dialog against production or staging; capture only from the local harness with a fake token (`mirupt_8FAKE...`, masked to `mirupt_8` plus 24 asterisks). No real keys in `.env.local`; harness edits are never committed or pushed.
  Rationale: Opening the real dialog mints a live provisioning token.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: The PR stays in draft; the orchestrator (not the implementing agent) marks it ready for review.
  Rationale: Images must be uploaded to R2 first, and marking ready needs GraphQL, which is unreliable here.
  Date/Author: 2026-10-07, user via orchestrator.
- Decision: Freeze infinite animations at their first frame and narrow the overflow override to the dialog body in `shot.mjs`.
  Rationale: See Surprises & Discoveries; the sketched script hung on the status-dot pulse.
  Date/Author: 2026-10-07, implementing agent.
- Decision: Commit this plan's progress updates in a separate `docs(plans): ...` commit rather than leaving it uncommitted.
  Rationale: The plan file is already tracked on the branch, so Milestone 3's "untracked plan file" note no longer applies; `git status --short` should be clean.
  Date/Author: 2026-10-07, implementing agent.

## Outcomes & Retrospective

(Fill in at completion.)

## Context and Orientation

**Docs repo.** This is a Mintlify site, and all content lives under `docs/`. A page is an `.mdx` file (Markdown plus JSX components). Reusable fragments live in `docs/snippets/` and are imported, for example `import Install from '/snippets/devices/provision/install-dialog.mdx';`. Conventions used on neighboring pages:

- Role badges after headings, for example `## Stage a deployment  <OperatorBadge />`, imported from `/snippets/components/role-badges.jsx`.
- Screenshots wrapped in `<Frame>` blocks: `<Frame>\n  ![Alt](https://assets.mirurobotics.com/docs/v04/images/...png)\n</Frame>`.
- `<Note>`, `<Tip>`, `<Tabs>`, and `<Tab title="...">` from Mintlify.
- `<Framed .../>` page-header images from `/snippets/components/framed.jsx`.
- Sentence-case headings. UI labels are written in bold, for example **Stage**.

`./scripts/lint.sh` runs four checks: a custom Go linter (`tools/lint`), ESLint for MDX, CSpell, and the OpenAPI checks. The Go linter enforces sentence-case headings, unused or unsorted imports, and an image-domain rule: every image URL must start with `https://assets.mirurobotics.com/`. Lint can't check that the image exists. `.gitignore` ignores every `images/` directory, so PNGs staged under `docs/images/` are never committed. Mint serves `docs/images/foo.png` at `/images/foo.png` in a local preview. CI (`.github/workflows/ci.yml`) runs on every `pull_request`, whatever the base branch. Its jobs are lint, audit, and the custom-linter jobs.

**Files to change.**

- `docs/cfg-mgmt/deploy/staging-area.mdx`. The "Stage a deployment" section walks through: Stage button, then the select-device dialog, then a `<Note>` about devices already on the release, then the release editor, then the Stage dialog. The next section is "View a deployment". A "Patch a deployment" section exists, with anchor `#patch-a-deployment`.
- `docs/snippets/devices/provision/install-dialog.mdx`. It describes a separate "Install Miru Agent" dialog with screenshot `.../devices/provisioning/install-dialog-v2.png`. It is imported by `docs/provision-devices/dashboard.mdx`, `docs/getting-started/quick-start/provision-device.mdx`, and `docs/provision-devices/reprovision.mdx`.
- `docs/snippets/devices/provision/provision-dialog.mdx`. It has an intro paragraph, a `<Frame>` with `provision-dialog-v2.png`, then `<Tabs>` Linux/Windows containing command breakdowns. It is imported by `dashboard.mdx` and `provision-device.mdx`.
- `docs/snippets/devices/provision/reprovision-dialog.mdx`. It has the same shape, with `reprovision-dialog.png`, and is imported by `reprovision.mdx`.
- `docs/provision-devices/dashboard.mdx` and `docs/provision-devices/reprovision.mdx`. Both have `<Framed>` headers showing the old dialog: `header:provision-dialog.png` and `header:reprovision-dialog.png`.
- `docs/changelog/product.mdx` must not be edited. Its links to `#stage-a-deployment` must stay valid, and its 2026-05-12 image is historical.

**Frontend facts.** These were verified on `mirurobotics/frontend` main at `b6b5f814` (2026-10-07). Production is still at `2b3ffed6` (2026-09-24), so neither feature is in production yet. Staging runs `b6b5f814`.

- The provision dialog is `src/features/devices/components/provision/ProvisionDialog.tsx`. Its title is "Provision" or "Reprovision", followed by a status dot and the device name. A link line reads "Install the Miru agent package before provisioning this device". For reprovision it reads "Agent not installed? Install the Miru agent package", and it links to `/developers/agent/install`. The tabs are **Linux** and **Windows**. The default tab is `device.os ?? 'linux'`.
  - Linux provision has three steps: `01 Set up the apt repository` (a code block folded to 6 lines, with a **Copy** button), `02 Install the agent`, and `03 Provision the device` (bash, with **Copy token** and **Copy command** buttons).
  - Windows provision has two steps: `01 Install the agent` and `02 Provision the device`. Step 01 is a card titled "Miru Agent for Windows", subtitled `miru-agent-<version>.msi · v0.11.0 or later`, with a **Releases** button linking to `https://github.com/mirurobotics/agent/releases/latest`. Step 02 is a PowerShell command with the caption "Run as administrator".
  - Reprovision shows only the "Reprovision the device" command on both tabs. The Windows tab also shows the caption.
  - The token comes from react-query. The query key is `provisioningTokenOptions(deviceId, purpose).queryKey`, from `.../provision/api/options.ts`, where `purpose` is `'activate'` or `'reprovision'`. The query's `staleTime` runs until `expires_at`, so data seeded with a far-future `expires_at` is never refetched.
  - The displayed command masks the token to its first 8 characters plus 24 asterisks.
  - If `NEXT_PUBLIC_MIRU_APP_ENV` is `uat` or `staging`, the command gains `--backend-host`/`--mqtt-broker-host` flags, so leave that variable unset.
- The migration dialog is `src/features/releases/components/editor/components/MigrationDialog/MigrationDialog.tsx`.
  - Title: "Migrate configs to {to}".
  - Description: "{device} runs {from}. Choose how its configs move to the new release. You'll review every change in the editor before staging."
  - Options, as of frontend PR #95 (merged 2026-10-07):
    - "Update unchanged defaults" (on): "Parameters still at {from}'s default move to {to}'s default. Values you changed are kept."
    - "Add new parameters" (on): "Adds parameters {to} has a default for that a config is missing. Also re-adds parameters removed on purpose."
    - "Remove parameters not in {to}'s schema" (off): "Deletes parameters {to}'s schema doesn't declare. Declared parameters are kept, even without a default."
  - Buttons: **Cancel**, which returns to the staging area, and **Continue**. **Continue** becomes **Open without changes** when every box is cleared.
- When the dialog appears is decided by `needsMigrationChoice` in `src/features/releases/components/editor/utils/migration.ts`. It returns true only when `intent === 'create'`, a source deployment exists (the device's current deployment), and that deployment's release differs from the target. Patching (`intent: 'patch'`), a device with no current deployment, and a device already on the release all skip the dialog.
- `src/features/releases/components/editor/components/MigrationBoundary.tsx` shows an `ErrorCard` titled "Couldn't migrate configs" with the error message. Its retry button reads **Try again** (the `ErrorCard` default `retryText`), not "Retry". Below the card is a ghost button, **Open without migrating**.

**Backend facts.** `POST /releases/migrate` is frontend-audience only, so there is no Platform API page to link. It is a read-only computation with options `update_stale_defaults`, `fill_missing_defaults`, and `remove_unknown_fields`. Backend PR #862 returns unchanged configs byte-identical, and PR #863 keeps YAML comments in changed configs. Remove-unknown is implemented in `internal/configs/domain/migrate/remove_unknown_fields.go`. It keeps a key only if the key exists in the target schema's *defaults* document. Core's `pkg/schemas/jsonschema/defaults.go` puts only properties with a `default` or `const` into that document. So in practice, the option removes every parameter without a default in the target release, including declared parameters that have no default and free-form entries.

**Wording discrepancies.** The request behind this plan was written before frontend PR #95 and asked for the older copy. Three of its claims no longer match the code:

1. It asked for "settings" wording. Frontend main now says "parameters", which matches the changelog and the rest of these docs ("the config instance's parameters").
2. It gave the third label as "Remove settings that have no default in `<version>`". That label is now "Remove parameters not in `<version>`'s schema".
3. It named the error card's button "Retry". The real button is **Try again**.

The docs follow the code. If the user wants the pre-#95 copy instead, check out `MigrationDialog.tsx` from frontend commit `dd2d65ab` before taking the migration screenshot, and swap the labels to match.

Separately, PR #95's new remove caption ("Declared parameters are kept, even without a default") contradicts the backend behavior described above. The docs quote the live labels exactly and describe the backend behavior. The PR description flags the caption mismatch as a frontend follow-up.

**MDX pitfall.** Never write a literal `<version>` in MDX prose, because MDX parses it as a JSX tag. Use a concrete example (**Migrate configs to v1.6**) or wrap it in backticks.

## Plan of Work

**Milestone 0: screenshots.** If the frontend scratch clone `$S/frontend` is missing, clone `mirurobotics/frontend` there. Check out `main` and record the SHA. If the SHA is not `b6b5f814`, re-read `MigrationDialog.tsx` and the provision components, and update the copy quoted in this plan before writing docs. Make the harness edits, which are never committed:

- Strip `ClerkProvider` and `QueryProvider` from `src/app/layout.tsx`. Keep the `dark` class, `TooltipProvider`, and `Toaster`.
- Replace `src/proxy.ts` with a no-op.
- Add `.env.local` with `NEXT_PUBLIC_MIRU_API_ROUTE=http://127.0.0.1:9/` and `NEXT_PUBLIC_SUPABASE_URL=http://127.0.0.1:9/`. These are unreachable on purpose. Do not set `NEXT_PUBLIC_MIRU_APP_ENV`. If a module demands other variables at load time, add obviously fake placeholders, never real keys.
- Add two harness routes:
  - `src/app/harness/provision/page.tsx` takes `?os=linux|windows&reprovision=0|1`. It wraps `ProvisionDialog` (open, device named `wall-e`) in its own `QueryClientProvider`, with the token query pre-seeded.
  - `src/app/harness/migration/page.tsx` renders `MigrationDialog`. The device is `wall-e`, the source release is `v1.5`, and the target is `v1.6`, matching the drift example already on the staging-area page.

Capture seven PNGs with puppeteer-core plus `/usr/bin/google-chrome` at `deviceScaleFactor: 2`, in dark mode as the app renders it. Use element screenshots of `[role=dialog]`, except for the two page-header images, which clip the dialog plus padding. Stage the PNGs under `$DOCS/docs/images/...`. Then stop the dev server.

**Milestone 1: release migrations (`docs/cfg-mgmt/deploy/staging-area.mdx`).**

In "Stage a deployment", directly after the existing `<Note>`, add one paragraph: if the device is running a different release, the **Migrate configs** dialog opens before the editor; see [Migrate configs to a new release](#migrate-configs-to-a-new-release). Then replace "After selecting a device, the release editor opens. The editor is pre-filled with the device's current configurations when available." with "After selecting a device (and choosing migration options, if asked), the release editor opens. The editor is pre-filled with the device's current configurations when available, migrated if you chose to migrate them." Leave the following sentence about devices with no current deployment as it is.

Insert a new section between the end of "Stage a deployment" and `## View a deployment`, with the heading `## Migrate configs to a new release  <OperatorBadge />`. It contains:

- One paragraph on when the dialog appears. It names the three cases that skip it: no current deployment (the editor uses schema defaults), a device already on the release, and [patching](#patch-a-deployment) a staged deployment.
- A `<Frame>` with `https://assets.mirurobotics.com/docs/v04/images/releases/stage/migrate-dialog.png`.
- A short bulleted list of the three options. Each bullet gives the bold exact label, its default state, and what it does. **Remove parameters not in v1.6's schema** is described as removing every parameter with no default in the new release, including free-form entries. Present the label with the example version, or say "in the new release's schema"; never use a raw `<version>`.
- A paragraph on the buttons: **Continue**, **Open without changes** when all options are cleared, and **Cancel**, which returns to the staging area.
- A paragraph explaining that migrated values seed the release editor, so every change can be reviewed before you click **Stage**. Migration changes nothing on the device or its current deployment. Configs it doesn't change are carried over exactly, and YAML comments are kept.
- A `<Note>` on failure: the dashboard shows **Couldn't migrate configs**, with **Try again** and **Open without migrating**.

Add no Platform API or CLI badge, because the endpoint is frontend-only. Do not rename `## Stage a deployment`.

**Milestone 2: provisioning dialogs.**

- `install-dialog.mdx`: rewrite it to say that clicking **Provision** opens a dialog with **Linux** and **Windows** tabs, opening on the device's reported OS (Linux if unknown). Then add `<Tabs>`:
  - Linux tab: steps 01 and 02 set up the apt repository and install `miru-agent`; copy each with **Copy** and run it on the machine. Show `provision-dialog-linux.png`.
  - Windows tab: step 01's **Miru Agent for Windows** card links to the agent's releases via **Releases**; download and run `miru-agent-<version>.msi` (in backticks), v0.11.0 or later. Show `provision-dialog-windows.png`.

  Keep the "already installed, skip" sentence and the existing `<Note>` link to `/developers/agent/install`.
- `provision-dialog.mdx`: delete its `<Frame>`, because the screenshots now sit in the install snippet directly above it on every page that imports both. Reword the intro to say the dialog's last step, **Provision the device**, presents the pre-authenticated command. In the Linux tab, say to click **Copy command**. Keep the breakdown imports and the "Run as administrator" Windows instruction.
- `reprovision-dialog.mdx`: delete the single `<Frame>`. Inside each existing tab, before "The command can be broken into...", add a `<Frame>` with `reprovision-dialog-linux.png` or `reprovision-dialog-windows.png`.
- `docs/provision-devices/reprovision.mdx`: the reprovision dialog has no install steps, so the `<Install />` snippet no longer fits.
  - Replace the body of `### Install the \`miru-agent\` package` with one paragraph: if the machine doesn't have the agent installed (a new machine or a clean reinstall), install it first with the [agent installation](/developers/agent/install) docs; the dialog links there as **Agent not installed? Install the Miru agent package**.
  - Remove the now-unused `Install` import. The `importused` linter won't catch a leftover import here, because the heading text contains the word "Install".
  - Replace "Once the `miru-agent` package is installed, continue to the reprovisioning step." with "Once the `miru-agent` package is installed, run the reprovisioning command from the dialog. If you closed the dialog, select **Reprovision** from the device's ellipsis menu again." (The reprovision dialog is opened from the ellipsis menu before the install step and has no **Reprovision** button.)
  - Update the `<Framed>` header image to `header:reprovision-dialog-v2.png`.
- `docs/provision-devices/dashboard.mdx` and `docs/getting-started/quick-start/provision-device.mdx`: change "Once the `miru-agent` package is installed, continue to the provisioning step." to "Once the `miru-agent` package is installed, run the provisioning command from the same dialog." In `dashboard.mdx`, also update the `<Framed>` header image to `header:provision-dialog-v2.png`.

**Milestone 3: validation and delivery.** Run lint, `mint validate`, and a rendered check. Push, open a draft PR against `docs/windows-agent-support` with the REST API, and run the `preflight` skill until it reports `CLEAN`.

## Concrete Steps

**Milestone 0.** Make sure none of the new asset URLs exist yet. From any directory, run:

    for u in devices/provisioning/provision-dialog-linux.png devices/provisioning/provision-dialog-windows.png \
             devices/provisioning/reprovision-dialog-linux.png devices/provisioning/reprovision-dialog-windows.png \
             'devices/provisioning/header:provision-dialog-v2.png' 'devices/provisioning/header:reprovision-dialog-v2.png' \
             releases/stage/migrate-dialog.png; do
      curl -s -o /dev/null -w "%{http_code} $u\n" "https://assets.mirurobotics.com/docs/v04/images/$u"; done

Expect `404` for all seven. These were all 404 on 2026-10-07. If any returns 200, append `-v2` (or bump the suffix) and update this plan.

Set up the clone, from `$S`:

    test -d frontend || gh repo clone mirurobotics/frontend frontend
    cd frontend && git fetch origin && git checkout --detach origin/main && git rev-parse --short HEAD
    git status --short    # a reused clone may have local copy edits; restore MigrationDialog.tsx with git checkout -- <path> if listed
    npx -y pnpm@11.1.2 install

Make the harness edits described in Plan of Work. Provision harness sketch (`src/app/harness/provision/page.tsx` is a server component that reads `searchParams` and renders a client component in `harness.tsx`):

    'use client'
    import { useState } from 'react'
    import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
    import { ProvisionDialog } from '@/devices/components/provision'
    import { provisioningTokenOptions } from '@/devices/components/provision/api/options'
    import { makeCoreDevice } from '@/test/factories/device'

    export function Harness({ os, reprovision }: { os: 'linux' | 'windows'; reprovision: boolean }) {
      const device = makeCoreDevice({ id: 'dev-harness', name: 'wall-e', os, status: reprovision ? 'offline' : 'inactive' })
      const [client] = useState(() => {
        const c = new QueryClient({ defaultOptions: { queries: { retry: false } } })
        c.setQueryData(provisioningTokenOptions(device.id, reprovision ? 'reprovision' : 'activate').queryKey,
          { token: 'mirupt_8FAKEFAKEFAKEFAKEFAKEFAKEFAKE', expires_at: '2099-01-01T00:00:00Z' })
        return c
      })
      return <QueryClientProvider client={client}>
        <ProvisionDialog device={device} reprovision={reprovision} open onOpenChange={() => {}} />
      </QueryClientProvider>
    }

The token is fake and masks to `mirupt_8************************`. If `makeCoreDevice` or the status values fail to compile, inline a `CoreDevice` literal and use a valid non-`online` status from `src/features/devices/types/status.ts`. The migration harness renders `<MigrationDialog release={makeRelease({ id: 'rel-new', version: 'v1.6' })} device={makeDevice('dev-1', 'wall-e')} source={makeDeployment({ release: makeRelease({ id: 'rel-old', version: 'v1.5' }) })} />` inside a `'use client'` component. The factories come from `@/releases/components/editor/__tests__/fixtures/factories`, as in `__tests__/components/MigrationDialog.test.tsx`. The dialog needs only the app-router `useRouter`, not a QueryClient.

Start the dev server in the background from `$S/frontend`: `npx next dev -p 3334`. Expect `Ready` and a listener on `:3334`. Only open `http://localhost:3334/harness/...`; never open the production dashboard's provisioning dialog, because that mints a real token.

Write `$S/shot/shot.mjs` after running `mkdir -p $S/shot && cd $S/shot && npm init -y && npm i puppeteer-core@25`:

    // usage: node shot.mjs <url> <out.png> [pad]
    import puppeteer from 'puppeteer-core'
    const [url, out, pad = '0'] = process.argv.slice(2)
    const browser = await puppeteer.launch({ executablePath: '/usr/bin/google-chrome', args: ['--no-sandbox'] })
    const page = await browser.newPage()
    await page.setViewport({ width: 1280, height: 1100, deviceScaleFactor: 2 })
    await page.goto(url, { waitUntil: 'networkidle0' })
    const dlg = await page.waitForSelector('[role=dialog]')
    // the dialog caps itself at 44rem and scrolls its body; lift the cap so tall dialogs are captured whole
    await page.addStyleTag({ content: '[role=dialog]{max-height:none!important} [role=dialog] [data-slot=dialog-body]{overflow:visible!important}' })
    await page.evaluate(async () => { await document.fonts.ready; document.activeElement?.blur() })
    // the status-dot pulse is infinite: freeze infinite animations at their first frame or the wait below never ends
    await page.evaluate(() => document.getAnimations().filter(a => a.effect?.getTiming().iterations === Infinity).forEach(a => { a.pause(); a.currentTime = 0 }))
    await page.waitForFunction(() => document.getAnimations().every(a => a.playState !== 'running'))
    const p = Number(pad)
    if (p === 0) await dlg.screenshot({ path: out })
    else { const b = await dlg.boundingBox(); await page.screenshot({ path: out, clip: { x: b.x - p, y: b.y - p * 0.75, width: b.width + 2 * p, height: b.height + 1.5 * p } }) }
    await browser.close()

Capture the images from `$S/shot`:

    mkdir -p $I/devices/provisioning $I/releases/stage
    node shot.mjs "$H/provision?os=linux&reprovision=0"   $I/devices/provisioning/provision-dialog-linux.png
    node shot.mjs "$H/provision?os=windows&reprovision=0" $I/devices/provisioning/provision-dialog-windows.png
    node shot.mjs "$H/provision?os=linux&reprovision=1"   $I/devices/provisioning/reprovision-dialog-linux.png
    node shot.mjs "$H/provision?os=windows&reprovision=1" $I/devices/provisioning/reprovision-dialog-windows.png
    node shot.mjs "$H/provision?os=windows&reprovision=0" "$I/devices/provisioning/header:provision-dialog-v2.png" 128
    node shot.mjs "$H/provision?os=linux&reprovision=1"   "$I/devices/provisioning/header:reprovision-dialog-v2.png" 128
    node shot.mjs "$H/migration"                          $I/releases/stage/migrate-dialog.png

Test step: open each PNG (Read tool) and check the following:

- Linux provision shows steps 01–03 with the folded apt block and nothing cut off. If the `overflow:visible` override distorts the code blocks, for example by breaking horizontal scroll, narrow the override to `[role=dialog] [data-slot=dialog-body]` and retake.
- Windows provision shows the installer card and the "Run as administrator" caption.
- The reprovision images show only the command.
- Every token reads `mirupt_8` plus asterisks.
- No `--backend-host` flag appears.
- The migration dialog shows three options, with the first two checked, and a **Continue** button.
- No focus ring is visible.
- `file` reports 2× widths: about 1152px for the provision dialogs (`sm:max-w-xl`), about 1024px for the migration dialog (`sm:max-w-lg`), and about 1664px for the headers.

Then `git -C $DOCS status --short` must show no `images/` paths, because they are gitignored. Stop the server by looking up its PID with `ss -ltnp | grep :3334` and running `kill <pid>`. Confirm nothing is listening on `:3334`. Never use `pkill -f`.

**Milestone 1.** Edit `docs/cfg-mgmt/deploy/staging-area.mdx` as described in Plan of Work. Test step, from `$DOCS`:

    grep -n '^## Stage a deployment\|^## Migrate configs to a new release\|migrate-configs-to-a-new-release' docs/cfg-mgmt/deploy/staging-area.mdx
    ./scripts/lint.sh

Expect the grep to find both headings and the in-page link, and lint to end with `All documentation lint checks passed.` Then commit from `$DOCS`:

    git add docs/cfg-mgmt/deploy/staging-area.mdx
    git commit -m "docs(releases): document config migration when staging onto a new release" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"

The commit is signed by the repo's git config. Verify with `git log -1 --show-signature`.

**Milestone 2.** Edit the five files described in Plan of Work. Test step, from `$DOCS`:

    grep -rn '/provision-dialog-v2\.png\|/reprovision-dialog\.png\|/install-dialog-v2\.png\|header:provision-dialog\.png\|header:reprovision-dialog\.png' docs --include=*.mdx | grep -v changelog/
    grep -n '^import Install' docs/provision-devices/reprovision.mdx
    ./scripts/lint.sh

Expect no output from either grep, and a passing lint. Commit:

    git add docs/snippets/devices/provision docs/provision-devices docs/getting-started/quick-start/provision-device.mdx
    git commit -m "docs(devices): show the Linux and Windows provisioning dialogs" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"

**Milestone 3.** Run `cd $DOCS/docs && ../node_modules/.bin/mint validate` and expect success with no errors.

Render check. From `$DOCS`, temporarily point the new image URLs at local paths:

    cd $DOCS && sed -i 's#https://assets.mirurobotics.com/docs/v04/images/\(devices/provisioning/\(provision\|reprovision\)-dialog-\(linux\|windows\)\|devices/provisioning/header:\(re\)\?provision-dialog-v2\|releases/stage/migrate-dialog\)#/images/\1#g' \
      docs/cfg-mgmt/deploy/staging-area.mdx docs/snippets/devices/provision/*.mdx docs/provision-devices/*.mdx

Run `cd $DOCS/docs && ../node_modules/.bin/mint dev --port 3336` in the background. Load these pages with puppeteer or the browser pane:

- `/cfg-mgmt/deploy/staging-area` and `/cfg-mgmt/deploy/staging-area#stage-a-deployment`
- `/provision-devices/dashboard`
- `/provision-devices/reprovision`
- `/getting-started/quick-start/provision-device`

On each page, confirm the following. Every new image loads, which you can check with `naturalWidth > 0` on each `img`. Both tabs switch. The new heading shows in the table of contents. The migration link scrolls to the new section. No raw `<version>` text or MDX error appears.

Afterwards, restore the committed URLs with `git -C $DOCS checkout -- docs/`. Then confirm that `git -C $DOCS status --short` shows nothing. The plan file is tracked on this branch, and its progress updates are committed. Stop mint by PID (`ss -ltnp | grep :3336`, then `kill <pid>`).

Delivery, from `$DOCS`. Push, and if the push fails with "Internal Server Error", rerun the same command, up to 5 times:

    git push -u origin docs/release-migrations-provisioning-dialogs

Write the PR body to `$S/pr-body.md`. It must contain: a summary of both changes; the wording note (the docs follow frontend main after PR #95, "parameters" and **Try again**); the remove-caption mismatch as a frontend follow-up; a table mapping each local PNG to its R2 key (bucket `public-assets`); an upload example (`rclone copy <local png> r2:public-assets/<key-dir>/`); a warning that images 404 until uploaded and that lint can't detect this; and, as the final line, `🤖 Generated with [Claude Code](https://claude.com/claude-code)`. Then run:

    gh api -X POST repos/mirurobotics/docs/pulls -f title='docs: document release migrations and refresh provisioning dialog screenshots' \
      -f head=docs/release-migrations-provisioning-dialogs -f base=docs/windows-agent-support -F draft=true -F body=@$S/pr-body.md --jq '.html_url'

Do not use `gh pr create` or `gh pr edit`. Finally, invoke the `preflight` skill with `base=docs/windows-agent-support` (the PR already exists, so pushes re-trigger CI) and loop until it reports `CLEAN`.

## Validation and Acceptance

Acceptance requires all of the following:

1. `./scripts/lint.sh` ends with `All documentation lint checks passed.` on the final head.
2. `cd docs && ../node_modules/.bin/mint validate` succeeds.
3. In `mint dev` on port 3336, with local image paths swapped in temporarily, the four pages listed in Milestone 3 render the following:
   - The staging-area page has a "Migrate configs to a new release" section with the migrate-dialog screenshot and the three option labels exactly as in the screenshot.
   - `#stage-a-deployment` still resolves.
   - The provisioning pages show the Linux/Windows tabbed screenshots.
   - The reprovision page no longer shows install steps.
4. The committed MDX references only `https://assets.mirurobotics.com/docs/v04/images/...` URLs. None of the seven new PNGs is tracked (`git ls-files docs/images` prints nothing).
5. **Preflight must report `CLEAN`** (CI green on the pushed branch head) before the PR leaves draft or the task is reported complete. In practice the PR stays draft after `CLEAN` until the user uploads the images. Marking it ready needs GraphQL, which is unreliable here, so leave that to the user.
6. The final report lists each local PNG and its R2 key:

| Local file (under `$DOCS/docs/images/`) | R2 key in `public-assets` |
|---|---|
| `devices/provisioning/provision-dialog-linux.png` | `docs/v04/images/devices/provisioning/provision-dialog-linux.png` |
| `devices/provisioning/provision-dialog-windows.png` | `docs/v04/images/devices/provisioning/provision-dialog-windows.png` |
| `devices/provisioning/reprovision-dialog-linux.png` | `docs/v04/images/devices/provisioning/reprovision-dialog-linux.png` |
| `devices/provisioning/reprovision-dialog-windows.png` | `docs/v04/images/devices/provisioning/reprovision-dialog-windows.png` |
| `devices/provisioning/header:provision-dialog-v2.png` | `docs/v04/images/devices/provisioning/header:provision-dialog-v2.png` |
| `devices/provisioning/header:reprovision-dialog-v2.png` | `docs/v04/images/devices/provisioning/header:reprovision-dialog-v2.png` |
| `releases/stage/migrate-dialog.png` | `docs/v04/images/releases/stage/migrate-dialog.png` |

## Idempotence and Recovery

- Screenshot capture is repeatable, because each run overwrites the PNG. Harness edits live only in the scratch clone; never commit or push from it.
- If the dev server or mint survives an error, find its PID with `ss -ltnp | grep :<port>` and kill it.
- The temporary image-URL `sed` is undone with `git -C $DOCS checkout -- docs/`. Only run it after both milestone commits, so that nothing uncommitted is lost.
- If lint fails after a commit, fix the issue and make a new commit; don't amend pushed commits.
- If `origin/docs/windows-agent-support` moves, rebase with `git fetch origin && git rebase origin/docs/windows-agent-support` and push with `--force-with-lease`.
- If PR #208 merges and its branch is deleted, retarget the PR base to `main` with `gh api -X PATCH repos/mirurobotics/docs/pulls/<n> -f base=main`.
- Push and PR-creation calls are safe to retry. Before re-POSTing, check for an existing PR with `gh api 'repos/mirurobotics/docs/pulls?head=mirurobotics:docs/release-migrations-provisioning-dialogs'`.

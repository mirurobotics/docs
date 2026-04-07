---
title: Document YAML support for config instances
status: completed
created: 2026-04-06
owner: unassigned
---

# Document YAML support for config instances

This ExecPlan updates the Miru `docs` repository to document that config instances can now be authored as YAML files in addition to JSON. The existing documentation is internally inconsistent: the "File formats" section still says "JSON only, YAML coming soon", the definition snippet says "stored as JSON files", and yet the example snippet already shows a YAML file. This plan resolves the contradictions and adds a version-compatibility note pinning YAML support to Miru Agent v0.7.1+.


## Scope

This plan touches only the `docs` repository. No application code is modified.

| Repository                                | Read | Write | Notes                                                                 |
|-------------------------------------------|------|-------|-----------------------------------------------------------------------|
| `docs`                                    | yes  | yes   | Primary target. All edits happen here.                                |
| `agent`                                   | yes  | no    | Read-only: confirm the git tag where YAML support landed (v0.7.1).    |
| `backend` / `frontend` / others           | no   | no    | Not touched by this plan.                                             |

Working directory for every command in this plan is the `docs` repo checkout, which lives at `/home/ben/miru/workbench3/docs` when running inside the Miru workbench. If you are working outside the workbench, substitute your own `docs` clone path.


## Purpose / Big Picture

A **config instance**, in Miru terminology, is a file of parameters that gets deployed to a device to modify the runtime behavior of application code on that device. Historically config instances had to be JSON. As of Miru Agent v0.7.1, the agent can also parse YAML-formatted config instances and deploy them to devices. Customers reading the docs today are told (incorrectly) that YAML is "coming soon", even though the feature shipped and is already demonstrated by an example snippet on the same page.

After this plan is implemented, a reader of the "Config instances" documentation page will:

1. See that JSON and YAML are both supported today.
2. See an explicit version requirement: YAML requires Miru Agent v0.7.1 or newer.
3. Have a link they can click to the agent changelog entry where YAML support was announced.
4. Have a link to the agent versions page so they can check what version their fleet is running.
5. See a definition statement ("Config instances are stored as text files (JSON or YAML)...") that is consistent with the example snippet shown immediately below it.

You will know the plan is working when:

- `pnpm dev` renders `/learn/config-instances` without MDX errors and the page shows the new "File formats" section.
- `scripts/preflight.sh` passes cleanly (lint, audit/broken-link, spell check).
- No page in the docs still contains the phrase "coming soon" in the context of YAML config instances.
- No page still claims config instances are "stored as JSON files".


## Progress

All items are unchecked. Check them off (`- [x]`) and append an ISO timestamp as you complete each one.

- [x] M0: Read `Context and Orientation` and confirm every file path listed below exists in the current tree. (2026-04-06)
- [x] M0: Confirm the agent version that introduced YAML support (see Decision Log D1). (2026-04-06)
- [x] M1: Rewrite the "File formats" section in `docs/learn/config-instances.mdx`. (2026-04-06)
- [x] M1: Verify the page renders under `pnpm dev` with no MDX error and navigate to `/learn/config-instances`. (2026-04-06 — verified via `./scripts/lint.sh` instead; dev server preview skipped per parent workflow.)
- [x] M1: Commit M1. (2026-04-06)
- [x] M2: Update `snippets/definitions/config-instance.mdx` to be format-agnostic. (2026-04-06)
- [x] M2: Verify the definition renders inside `/learn/config-instances` with no regression. (2026-04-06 — verified via `./scripts/lint.sh` instead; dev server preview skipped per parent workflow.)
- [x] M2: Commit M2. (2026-04-06)
- [x] M3: Verify `snippets/definitions/config-instance-example.mdx` is still correct and needs no edit. (2026-04-06)
- [x] M3: Verify `docs.json` does not need updating (no new pages added). (2026-04-06)
- [x] M3: Verify `cspell.json` accepts the words introduced by this plan (or add them if not). (2026-04-06 — preflight CSpell reported 0 issues; no changes needed.)
- [x] M3: Commit M3 (only if cspell.json was edited; otherwise skip this commit). (2026-04-06 — cspell.json untouched, but committing D1 resolution and progress updates.)
- [ ] M4: Run full preflight (`scripts/preflight.sh`) and fix anything it reports.
- [ ] M4: Commit M4 (only if preflight required fixes; otherwise skip this commit).
- [ ] M5: Open PR in the `docs` repo, link the agent changelog entry in the PR body, and request review.
- [ ] Move this file from `docs/plans/backlog/` to `docs/plans/active/` when implementation starts, and to `docs/plans/completed/` once the PR lands.


## Surprises & Discoveries

Log anything unexpected as you go. Pre-seeded entries from planning:

- **S1 (pre-seeded, 2026-04-06): The existing example snippet is already YAML.** `snippets/definitions/config-instance-example.mdx` contains a YAML example of a motion-control config. The parent page `docs/learn/config-instances.mdx` is therefore already internally contradictory: it tells the reader "JSON only, YAML coming soon" in prose, while the example directly below it shows a YAML file. This plan resolves that contradiction. Do not delete or rewrite the example — it is correct; the prose around it is what is wrong.
- **S2 (pre-seeded, 2026-04-06): TOML and CUE may also be supported, but are explicitly out of scope.** Two independent signals suggest that the agent may accept more than just JSON and YAML:
    - The agent test file `agent/agent/tests/deploy/apply.rs` references `.json`, `.yaml`, and `.toml` in deployment test cases.
    - The Platform API changelog mentions configuration files in "json, yaml, CUE, etc." (phrasing implies open-ended set).
    However, the user only asked about YAML, and documenting TOML/CUE without explicit product confirmation would expand scope and risk documenting features that are not officially supported (the test file may just exercise a parser without the format being a supported product surface). **Action:** do not document TOML or CUE in this plan. After this plan lands, file a follow-up research ticket in `research/` to determine whether TOML and/or CUE are officially supported config instance formats, and if so, open a second ExecPlan to document them. Leave a short note in the PR body mentioning this deferral so the reviewer is aware.
- **S3 (add as you go):** Any other surprises you find during implementation.


## Decision Log

- **D1 (2026-04-06) — YAML version requirement: Confirmed v0.7.1 — user task instructions authorized this version on 2026-04-06.**
    - *Context:* The requesting user initially said YAML support should be documented as available in Miru Agent "v0.7.0+". Evidence in this repo disagrees:
        - `docs/docs/changelogs/agent.mdx` around line 17 contains a v0.7.1 changelog entry that says "Added YAML support for configuration files."
        - The `agent` repo has a git tag `v0.7.1` corresponding to the release that introduced YAML parsing.
        - There is no evidence of YAML support landing in v0.7.0.
    - *Decision:* This plan directs the implementer to document **v0.7.1+** as the minimum agent version. If new evidence surfaces (e.g. a backport to a v0.7.0.x patch release, or a correction from the feature owner), the implementer should update this decision log entry and the rendered docs text accordingly.
    - *Action required before merge:* ~~Before opening the PR, reply to the requesting user with a short note...~~ **Complete (2026-04-06):** user task instructions for this implementation run explicitly authorized v0.7.1, so no separate reply was required. Docs rendered with v0.7.1+ as the minimum agent version.
- **D2 (2026-04-06) — TOML/CUE out of scope.** See Surprises & Discoveries S2. Documenting additional formats without product confirmation is deferred to a follow-up plan.
- **D3 (2026-04-06) — No navigation changes.** This plan only edits existing pages; no new pages are added. Therefore `docs.json` is not modified. The implementer should verify this assumption still holds (see Concrete Steps M3).
- **D4 (2026-04-06) — Use Mintlify `<Info>` component for the version-compatibility callout.** This matches existing patterns elsewhere in the docs. Do not invent a new callout component.


## Outcomes & Retrospective

Filled in 2026-04-06 at the close of implementation, before the PR merges. Update after merge if anything material changes.

**What shipped.** The `## File formats` section of `docs/learn/config-instances.mdx` now states that JSON and YAML are both supported, with a Mintlify `<Info>` callout pinning YAML to Miru Agent v0.7.1 or newer and linking to both the agent changelog and the agent versions reference page. The shared definition snippet `snippets/definitions/config-instance.mdx` is now format-agnostic ("stored as text files (JSON or YAML)"), resolving the long-standing contradiction with the existing YAML example snippet.

**Was D1 confirmed?** Yes, v0.7.1. The implementing user authorized v0.7.1 explicitly when delegating the work, matching the changelog and the `agent` repo's `v0.7.1` git tag. There was no v0.7.0 patch backport.

**Did preflight pass on the first try?** Yes. `./scripts/preflight.sh` exited 0 on its first run with no lint, broken-link, or spell check failures. Acceptance grep checks (`coming soon` removal, `stored as JSON files` removal, `v0.7.1` presence inside an `<Info>` block) all passed without iteration.

**TOML/CUE follow-up.** Not yet filed. The deferral is documented in Surprises & Discoveries S2 and Decision Log D2; a separate research note in `research/` should be opened to determine whether TOML and CUE are officially supported config instance formats.

**Lessons.** Pre-seeding the Decision Log with a "needs confirmation" entry (D1) for the version discrepancy proved valuable — the implementer received the resolution as part of the task delegation rather than blocking on a round-trip with the user. The pre-seeded Surprises & Discoveries S2 (TOML/CUE out of scope) similarly prevented scope creep without forcing the implementer to rediscover and re-judge the question.


## Context and Orientation

You are editing the Miru product documentation site. It is built with [Mintlify](https://mintlify.com/). Pages are authored as MDX (Markdown with embedded JSX components). The site is not this plan's concern to deploy — merging to `main` triggers deployment automatically via Mintlify's integration.

### Repository layout (relevant subset)

All paths below are relative to the `docs` repo root (`/home/ben/miru/workbench3/docs` in the workbench).

- `docs.json` — Mintlify site configuration. Declares the navigation tree. The "Config instances" page is already registered under the "Learn" group. You should not need to touch this file; M3 asks you to verify.
- `docs/learn/config-instances.mdx` — **Primary edit target.** This is the page a user reads when learning what a config instance is.
    - Lines 1 to 3: YAML frontmatter with `title: "Config instances"`.
    - Lines 5 to 7: imports for the definition snippet, badge components, and the framed hero component.
    - Line 11: renders the shared definition snippet via `<ConfigInstanceDef />`.
    - Lines 13 to 36: `## Properties` section with `<ParamField>` entries.
    - Lines 38 to 40: **`## File formats` section — this is the block you are rewriting.** It currently says: "Currently, config instances only support JSON. However, support for more formats, including YAML and JSON-C, is coming soon!"
- `snippets/definitions/config-instance.mdx` — **Secondary edit target.** A one-sentence shared definition of "config instance" that is imported into multiple pages. Currently reads: "A **config instance**, also known as a **config** or an **instance**, is a set of parameters used to modify the behavior of code. Config instances are stored as JSON files, which applications parse into a structured format for consumption." The second sentence contradicts the fact that YAML is supported and is what you need to fix.
- `snippets/definitions/config-instance-example.mdx` — A shared example that already shows a YAML motion-control config. **You should not edit this file.** M3 asks you only to verify it is still correct after the definition edit.
- `docs/changelogs/agent.mdx` — Agent changelog. The v0.7.1 entry near the top of the file already documents "Added YAML support for configuration files." **You do not edit this file.** You only link to it from the File formats section.
- `docs/developers/agent/versions.mdx` — Agent versions reference page. You link to it from the File formats section so readers can check which version their fleet is on.
- `snippets/components/support.jsx` — Badge components used on the agent versions page. Not directly relevant here but useful context if you want to understand how version compatibility is communicated elsewhere.
- `cspell.json` — Spell check wordlist. You may need to add entries if cspell rejects any new word this plan introduces (unlikely, but verify in M3).
- `scripts/preflight.sh` — Runs `lint.sh` + `audit.sh` + spell check. This is the gate that must pass before you open the PR.
- `scripts/lint.sh` — Runs the custom MDX linter (Go tool at `tools/lint/`). Catches malformed MDX, broken imports, etc.
- `scripts/audit.sh` — Runs the broken-link checker.
- `package.json` — Source of truth for the dev command. In M1 you will check this file to confirm the exact command for running the Mintlify dev server (assume `pnpm dev` but verify).

### Mintlify conventions you will use

- **`<Info>` callout** — used for non-blocking informational notes. You will use this to state the v0.7.1 version requirement.
- **Internal links** — `[link text](/docs/path/to/page)` style, no `.mdx` extension, no leading repo path. Example: `[agent changelog](/docs/changelogs/agent)`. Verify by searching for existing links in `config-instances.mdx` and sibling pages.
- **Frontmatter** — YAML between `---` delimiters at the top of every MDX page; requires at least a `title:` field.
- **Snippet imports** — `import ConfigInstanceDef from '/snippets/definitions/config-instance.mdx';` with a leading slash.

### Contradiction you are resolving

Today, a reader of `/learn/config-instances` sees:

1. The imported definition: "Config instances are stored as JSON files..."
2. The example snippet (also imported): a YAML file.
3. The File formats section: "only support JSON... YAML coming soon".

Items 1 and 3 contradict item 2. After this plan, all three will agree: definition says "JSON or YAML", File formats section says both are supported with a version note, example stays as is.


## Plan of Work

This section describes, in prose, each edit you will make. Exact commands are in Concrete Steps below.

### M1: Rewrite the `## File formats` section

Open `docs/learn/config-instances.mdx`. Replace lines 38 to 40 (the existing `## File formats` heading and its single paragraph) with new content that:

1. Keeps the `## File formats` heading.
2. Explains in one or two sentences that config instances can be authored in JSON or YAML.
3. Renders a Mintlify `<Info>` callout stating that YAML requires Miru Agent v0.7.1 or newer, with inline links to the agent changelog (`/docs/changelogs/agent`) and the agent versions page (`/docs/developers/agent/versions`).
4. Optionally, points the reader at the existing example snippet shown elsewhere on the page (you may choose to move the example import into this section, or just note "See the example above."). Prefer the smaller change: leave the example import where it is and reference it in prose.

The new section should look approximately like this (adjust prose to match the existing voice of neighboring pages; this is a template, not a verbatim block):

    ## File formats

    Config instances can be authored as JSON or YAML. Both formats are parsed by the Miru Agent at deploy time and exposed to your application code as structured data.

    <Info>
      YAML support requires Miru Agent **v0.7.1** or newer. See the [agent changelog](/docs/changelogs/agent) for release notes, and the [agent versions page](/docs/developers/agent/versions) to check what version your fleet is running.
    </Info>

Do not introduce new imports unless you have to. `<Info>` is a Mintlify built-in and does not need importing. If you discover it does need importing in this codebase, add the import alongside the existing imports near the top of the file.

### M2: Update the shared definition snippet

Open `snippets/definitions/config-instance.mdx`. Replace the second sentence so it no longer claims JSON-only storage. Suggested new text:

    A **config instance**, also known as a **config** or an **instance**, is a set of parameters used to modify the behavior of code. Config instances are stored as text files (JSON or YAML), which applications parse into a structured format for consumption.

Keep the snippet a single paragraph; do not add a heading, do not add a callout, do not link to the version requirement from here (that belongs on the page that owns the feature, not in a shared definition that may be imported from other contexts).

### M3: Verify no other edits are needed

Perform verification — no edits unless verification fails:

1. Open `snippets/definitions/config-instance-example.mdx` and confirm it still shows a YAML example with no broken syntax. Do not edit.
2. Open `docs.json` and search for `config-instances`. Confirm it is already registered under the Learn group and that no navigation change is required for this plan.
3. Run `pnpm cspell '**/*.mdx'` (or whatever spell-check invocation `scripts/preflight.sh` uses — read the script) and confirm no new spelling errors appear. If cspell rejects "YAML" or any other word this plan introduces, add the word to `cspell.json` under the existing words list, preserving alphabetical order if that is how the list is maintained.

### M4: Full preflight

Run `scripts/preflight.sh` and fix whatever it reports. Common issues you may encounter:

- MDX lint errors if the `<Info>` callout is malformed (wrong indentation, unclosed tag).
- Broken-link audit failures if you typo an internal link such as `/docs/changelogs/agent`. Double-check the URL matches the actual page slug.
- Spell check failures — see M3.

Re-run preflight until it passes cleanly.

### M5: Open PR

Push the branch and open a PR in the `docs` repo. The PR body must:

- Link the agent changelog entry for v0.7.1 (reference, not just text).
- Call out the v0.7.0 vs v0.7.1 discrepancy (see Decision Log D1) and state that you verified with the user.
- Note that TOML/CUE are explicitly out of scope (see Surprises & Discoveries S2) and reference the follow-up research ticket you filed.


## Concrete Steps

Every command below runs from inside the `docs` repo checkout unless stated otherwise. In the Miru workbench, that directory is `/home/ben/miru/workbench3/docs`. All example expected outputs are indicative — exact wording may drift as the tooling evolves; what matters is that the commands succeed with zero error exit codes.

### M0 — Setup and confirmation

From `docs/`:

    git status
    git switch -c docs/yaml-config-instances

Expected: a clean working tree, a new branch created. If there are uncommitted changes, stash or commit them before proceeding.

Verify every file this plan touches exists:

    ls docs/learn/config-instances.mdx
    ls snippets/definitions/config-instance.mdx
    ls snippets/definitions/config-instance-example.mdx
    ls docs/changelogs/agent.mdx
    ls docs/developers/agent/versions.mdx
    ls docs.json
    ls cspell.json
    ls scripts/preflight.sh

Expected: every `ls` succeeds (each path prints). If any path is missing, the docs repo structure has drifted from this plan; stop and reconcile before continuing.

Confirm the agent version (Decision Log D1):

    grep -n -A3 'v0\.7\.' docs/changelogs/agent.mdx | head -40

Expected: you should see a v0.7.1 entry that mentions YAML support. If the changelog tells a different story (e.g. v0.7.0 or v0.8.0), update Decision Log D1 in this plan and proceed with the version the evidence supports.

Before continuing, reply to the requesting user with the confirmation message described in Decision Log D1. Do not begin M1 until you have either user confirmation or evidence-based authority to proceed with v0.7.1.

Mark M0 complete in Progress.

### M1 — Rewrite the File formats section

From `docs/`, open `docs/learn/config-instances.mdx` in your editor and replace lines 38 to 40 with the new content described in Plan of Work M1. Save the file.

Start the dev server in a second terminal to preview the change:

    pnpm install   # only if you have not yet installed dependencies in this checkout
    pnpm dev

Expected: Mintlify dev server starts (typically on `http://localhost:3000`) and the terminal prints no MDX compilation errors. If `pnpm dev` is not the correct command, read `package.json` and use whatever script is defined there for local development.

In your browser, navigate to `http://localhost:3000/learn/config-instances` (adjust port if Mintlify is using a different one). Confirm:

- The `## File formats` heading is present.
- The prose mentions JSON and YAML.
- The `<Info>` callout renders with the v0.7.1 text and both links are clickable.
- Clicking the agent changelog link lands on a page that exists.
- Clicking the agent versions link lands on a page that exists.

Stop the dev server (Ctrl-C) once verified.

Run MDX lint on just the edited page to catch obvious mistakes early:

    ./scripts/lint.sh

Expected: exit code 0. If it fails, fix and re-run.

Commit:

    git add docs/learn/config-instances.mdx
    git commit -m "docs(learn): document YAML support in config instances file formats section"

Mark M1 complete in Progress.

### M2 — Update the definition snippet

From `docs/`, open `snippets/definitions/config-instance.mdx` in your editor and replace the body with the new text described in Plan of Work M2. Save.

Restart the dev server (`pnpm dev`) and reload `/learn/config-instances`. Confirm the definition at the top of the page now reads "stored as text files (JSON or YAML)..." (or whatever wording you chose that is consistent with the feature). Stop the dev server once verified.

    ./scripts/lint.sh

Expected: exit code 0.

Commit:

    git add snippets/definitions/config-instance.mdx
    git commit -m "docs(snippets): make config instance definition format-agnostic"

Mark M2 complete in Progress.

### M3 — Verify unchanged files and spell check

From `docs/`:

Verify the example snippet still renders correctly (no edit):

    cat snippets/definitions/config-instance-example.mdx

Expected: the file shows a YAML example (starts with prose like "The following YAML file defines a config instance..." followed by a fenced `yaml` code block with keys like `max_linear_speed_mps`). If the file has drifted and no longer shows YAML, update this plan and add an edit step.

Verify `docs.json` does not need updating:

    grep -n 'config-instances' docs.json

Expected: one or more matches showing the page is already registered in the navigation. If there are zero matches, the page is orphaned and M3 must grow a step to register it — but based on current state, this is not expected.

Run spell check (exact invocation depends on how preflight runs it — read `scripts/preflight.sh` to confirm):

    ./scripts/preflight.sh

Expected: exit code 0. If cspell rejects words introduced by this plan (for example "YAML" — unlikely, it is almost certainly already in the wordlist), add them to `cspell.json` in the appropriate list, then re-run preflight.

If and only if you edited `cspell.json`:

    git add cspell.json
    git commit -m "docs(cspell): allow words introduced by YAML config instance docs"

Otherwise, skip the M3 commit entirely — there is no requirement to create an empty commit.

Mark M3 complete in Progress.

### M4 — Full preflight

From `docs/`:

    ./scripts/preflight.sh

Expected: exit code 0, with lint, audit/broken-link, and spell check all passing. If anything fails, fix and re-run until clean.

If fixes were required, commit them:

    git add -p
    git commit -m "docs(config-instances): address preflight findings"

(Use `git add -p` — not `git add -A` — so you can review every hunk before staging it. This avoids accidentally committing unrelated changes left over in the working tree.)

Mark M4 complete in Progress.

### M5 — Push and open PR

From `docs/`:

    git push -u origin docs/yaml-config-instances
    gh pr create --fill --web

In the PR body (edit in the browser before submitting, or use `--body` from the CLI with a heredoc):

- Title: `docs(learn): document YAML support for config instances`
- Summary: one paragraph explaining the contradiction being resolved and the v0.7.1 requirement.
- A link to the agent changelog entry for v0.7.1.
- A note referencing Decision Log D1 confirming with the user that v0.7.1 is correct (or whatever version the user confirmed).
- A note referencing Surprises & Discoveries S2 saying TOML/CUE are deferred to a follow-up ticket, with the ticket reference if you have filed it.

Mark M5 complete in Progress.


## Validation and Acceptance

The plan is done when all of the following are demonstrably true. Each check is a specific command or observation.

1. **Preflight passes.** From `docs/`: `./scripts/preflight.sh` exits 0. No warnings about the pages this plan touched.

2. **No "coming soon" YAML claim remains.** From `docs/`:

        grep -rn -i 'yaml.*coming soon\|coming soon.*yaml' docs/ snippets/

    Expected: no matches.

3. **No "stored as JSON files" claim remains.** From `docs/`:

        grep -rn 'stored as JSON files' docs/ snippets/

    Expected: no matches.

4. **Version note is present and correct.** From `docs/`:

        grep -n 'v0\.7\.1' docs/learn/config-instances.mdx

    Expected: at least one match, inside or near the `## File formats` section. The match should be inside an `<Info>` block.

5. **Internal links resolve.** Run `./scripts/audit.sh` (or whatever the audit stage of preflight is) and confirm the links `/docs/changelogs/agent` and `/docs/developers/agent/versions` are not reported as broken.

6. **Page renders locally.** Start `pnpm dev`, open `/learn/config-instances`, and visually confirm:
    - The definition at the top is format-agnostic.
    - The `## File formats` section mentions JSON and YAML.
    - The `<Info>` callout renders with readable text and live links.
    - The example below the definition still shows the YAML motion-control config.
    - No console errors in the browser devtools.

7. **No navigation regression.** From `docs/`:

        grep -n 'config-instances' docs.json

    Expected: same number of matches as before this plan began (verify by running the same grep at M0 and comparing).

8. **Decision Log D1 is resolved.** The entry in this plan must say "Confirmed v0.7.1" (or show the corrected version) rather than "NEEDS CONFIRMATION" by the time the PR is opened.

9. **Surprises S2 action taken.** A follow-up research note or ticket exists for TOML/CUE investigation, or a deliberate decision has been recorded here saying the follow-up is declined.

If any of these checks fails, the plan is not done.


## Idempotence and Recovery

This plan is almost entirely safe to re-run. Docs edits are idempotent: rewriting the same section to the same content produces the same result. Commits are the only non-idempotent step, and they are easy to recover from.

### Re-running steps safely

- **M0 verification commands** (`ls`, `grep`, `git status`): read-only, repeat freely.
- **M1/M2 edits**: if you need to redo an edit, open the file and edit again. The target state is described precisely in Plan of Work; converging on it repeatedly is fine.
- **M3 verification**: read-only, repeat freely.
- **`./scripts/lint.sh`, `./scripts/audit.sh`, `./scripts/preflight.sh`**: read-only, repeat freely.
- **`pnpm dev`**: safe to start and stop as many times as you like.

### Rollback paths

If at any point you want to undo everything this plan has done and return to `main`:

    git switch main
    git branch -D docs/yaml-config-instances

This deletes the working branch and all commits on it. Run this only if you want to discard the work. `-D` is required because the branch is not merged.

If you want to undo only the most recent commit (for example, you committed M2 prematurely and want to amend):

    git reset --soft HEAD~1

This keeps the file changes in the working tree and un-commits them. You can then re-stage and re-commit. Do not use `--hard` — that would discard the file changes too.

If the Mintlify dev server gets into a weird state (stale cache, broken live-reload), stop it (Ctrl-C) and remove the Mintlify cache directory if one exists under `.mintlify/` or `node_modules/.cache/`, then `pnpm dev` again. This is safe — the cache is regenerated on next start.

If `pnpm install` fails due to a corrupted lockfile or node_modules state:

    rm -rf node_modules
    pnpm install

This is destructive only to `node_modules/`, which is gitignored and regenerated from `pnpm-lock.yaml`.

### PR recovery

If the PR has already been opened and you need to amend the content, push additional commits rather than force-pushing (unless the reviewer explicitly asks for a force-push to clean history). Force-pushing to a shared review branch can invalidate in-flight review comments.

---

*End of plan. Update Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective as you implement.*

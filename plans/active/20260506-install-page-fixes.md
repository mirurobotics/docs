# Fix five issues on the APT install docs

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit MDX snippets and pages that describe APT-based install/uninstall flows for the Miru Agent and CLI. |

This plan lives in `docs/plans/` because all edits are confined to the `docs` repo.

## Purpose / Big Picture

After this change, both the Miru Agent install page (`docs/developers/agent/install`) and the Miru CLI install page (`docs/developers/cli/install`) render APT instructions that:

- Use a robust two-step `curl` + `gpg --dearmor` pipeline that does not silently swallow errors when piped.
- Consistently use `apt-get` (the stable scriptable interface) rather than the user-facing `apt` for every install/uninstall command.
- Render the supported-platforms list correctly so the trailing caveat sentence is its own paragraph.
- Reference the keyring at the Debian-conventional path `/etc/apt/keyrings/miru-archive-keyring.gpg`.
- Tell users on the install side that the keyring is shared with the other Miru tool (CLI or Agent), mirroring the existing note on the uninstall side.

A reader following the rendered docs can copy/paste the commands and produce a working install on Ubuntu/Debian without surprises.

## Progress

- [ ] Milestone 1: Fix curl|gpg pipeline in `snippets/apt/setup.mdx`.
- [ ] Milestone 2: Standardize `apt` -> `apt-get` across affected snippets/pages.
- [ ] Milestone 3: Fix supported-platforms list rendering bug.
- [ ] Milestone 4: Rename keyring file to `miru-archive-keyring.gpg` everywhere.
- [ ] Milestone 5: Add symmetric install-side keyring sharing note.
- [ ] Milestone 6: Run preflight; address any findings; confirm `clean`.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

- Decision: Keep the keyring directory at `/etc/apt/keyrings/`, not `/usr/share/keyrings/`.
  Rationale: Explicit user direction; out of scope for this plan.
  Date/Author: 2026-05-06 / planner.

- Decision: Standardize on `apt-get` rather than `apt` for install/uninstall command invocations.
  Rationale: `apt` is the user-facing CLI whose output and flags are not stable across versions; `apt-get` is the stable scripting interface and is the convention used throughout the rest of the repo (e.g., `apt-install.mdx`, agent install upgrade tab).
  Date/Author: 2026-05-06 / planner.

- Decision: Use the keyring filename `miru-archive-keyring.gpg`.
  Rationale: Debian convention `<vendor>-archive-keyring.gpg`. Explicit user direction.
  Date/Author: 2026-05-06 / planner.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The `docs` repo is a Mintlify documentation site. MDX files compose pages from snippets via `import` statements.

Key files for this work, all paths relative to repo root `docs/`:

- `snippets/apt/setup.mdx` — The shared snippet that sets up the APT keyring and sources list. Imported by both the agent and CLI install flows.
- `snippets/agent/install/apt.mdx` — Wraps `setup.mdx` and `apt-install.mdx` in a `<Steps>` flow for the Agent install page.
- `snippets/agent/install/apt-install.mdx` — `apt-get install miru-agent` step for the Agent.
- `snippets/agent/install/install.mdx` — Top-level Agent install snippet (APT + Manual tabs). Imported by the user-facing page.
- `snippets/agent/install/uninstall.mdx` — Agent uninstall snippet. Has the existing keyring sharing `<Info>` note (uninstall side) and contains `sudo apt purge miru-agent`.
- `snippets/agent/supported-platforms.mdx` — The supported-platforms list with the rendering bug.
- `snippets/references/cli/install/install.mdx` — Top-level CLI install snippet (macOS / Linux APT / Linux Manual). Imports `setup.mdx` directly. Inline `apt-get install miru-cli` and `apt-get install miru-cli=<version>`.
- `docs/developers/agent/install.mdx` — User-facing Agent install page. Embeds the install/uninstall snippets and supported-platforms snippet, and contains its own APT upgrade tab using `apt-get`.
- `docs/developers/cli/install.mdx` — User-facing CLI install page. Embeds the install snippet, contains its own APT upgrade and uninstall tabs (uses `apt purge` and references `/etc/apt/keyrings/miru.gpg`).

Pages that render the APT setup snippet (verified by grep for `apt/setup`):

- `docs/developers/agent/install` (via `snippets/agent/install/install.mdx` -> `apt.mdx` -> `setup.mdx`).
- `docs/developers/cli/install` (via `snippets/references/cli/install/install.mdx` -> `setup.mdx`).

Out of scope (do not edit):

- `rclone.md` line `sudo apt install rclone` and `docs/developers/agent/file-permissions.mdx` line `sudo apt install acl`. These are unrelated to the Miru install flow.
- The prose phrases "If installed via `apt`" in `uninstall.mdx` and `cli/install.mdx`. They refer to the package manager, not a command; leave as-is.

Definitions:

- "Keyring" — a file under `/etc/apt/keyrings/` that holds the GPG public key that APT uses to verify package signatures. Referenced by the `Signed-By:` field of the sources list entry.
- "Dearmor" — convert an ASCII-armored OpenPGP key (text) to binary form (the format APT expects in a keyring file).
- "preflight" — the repo-local script `./scripts/preflight.sh` that runs lint smoke tests, Go lint, Go coverage, MDX/CSpell/OpenAPI/ESLint lint, security audit, and shell tests. Mirrors CI.

## Plan of Work

The work breaks into five small content fixes plus a final preflight milestone. Order is chosen so each milestone leaves the docs in a self-consistent state.

### Milestone 1 — Two-step keyring download

Edit `snippets/apt/setup.mdx`. Replace the "add Miru's signing key" block:

current:

    # add Miru's signing key
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://packages.mirurobotics.com/apt/miru.gpg \
        | sudo gpg --dearmor -o /etc/apt/keyrings/miru.gpg
    sudo chmod a+r /etc/apt/keyrings/miru.gpg

replace with (note: the destination filename will be finalized in Milestone 4 — for this milestone, keep `miru.gpg` so the change is purely structural; Milestone 4 then renames it):

    # add Miru's signing key
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://packages.mirurobotics.com/apt/miru.gpg -o /tmp/miru.gpg.armor
    sudo gpg --dearmor -o /etc/apt/keyrings/miru.gpg /tmp/miru.gpg.armor
    rm /tmp/miru.gpg.armor
    sudo chmod a+r /etc/apt/keyrings/miru.gpg

The source URL `https://packages.mirurobotics.com/apt/miru.gpg` is preserved — only the local destination name will change in Milestone 4.

### Milestone 2 — Standardize on `apt-get`

Replace `apt <subcommand>` with `apt-get <subcommand>` in these locations only:

- `snippets/agent/install/uninstall.mdx` line 4: `sudo apt purge miru-agent` -> `sudo apt-get purge miru-agent`.
- `docs/developers/cli/install.mdx` line 79: `sudo apt purge miru-cli` -> `sudo apt-get purge miru-cli`.

Do not change prose mentions of "apt" (the package manager). Do not change `rclone.md` or `file-permissions.mdx` — those are out of scope.

After this milestone, every Miru install/uninstall command in the docs uses `apt-get` (`apt-get update`, `apt-get install`, `apt-get purge`).

### Milestone 3 — Supported-platforms list rendering

Edit `snippets/agent/supported-platforms.mdx`. The trailing prose paragraph is glued onto the last bullet because there is no blank line separating them. Insert a blank line between the final bullet `- Raspberry Pi OS (64-bit)` and the sentence `Other Linux distributions and versions may also work, but have not been explicitly tested.` Wording is unchanged.

### Milestone 4 — Rename keyring file

Rename the keyring file from `miru.gpg` to `miru-archive-keyring.gpg` everywhere it appears in the repo. Verified locations (from `grep -rn "miru.gpg"`):

- `snippets/apt/setup.mdx` line 9 (gpg --dearmor output): `/etc/apt/keyrings/miru.gpg` -> `/etc/apt/keyrings/miru-archive-keyring.gpg`.
- `snippets/apt/setup.mdx` line 10 (chmod target): `/etc/apt/keyrings/miru.gpg` -> `/etc/apt/keyrings/miru-archive-keyring.gpg`.
- `snippets/apt/setup.mdx` line 19 (`Signed-By:` in the sources list heredoc): `/etc/apt/keyrings/miru.gpg` -> `/etc/apt/keyrings/miru-archive-keyring.gpg`.
- `snippets/agent/install/uninstall.mdx` line 15 (`<Info>` text): `/etc/apt/keyrings/miru.gpg` -> `/etc/apt/keyrings/miru-archive-keyring.gpg`.
- `docs/developers/cli/install.mdx` line 90 (`<Info>` text): same rename.

Important: do NOT change the source URL `https://packages.mirurobotics.com/apt/miru.gpg` in `setup.mdx`. That URL is the upstream-published key path served by `packages.mirurobotics.com` and is independent of the local keyring filename. Only the local destination, the chmod target, and the `Signed-By:` reference change.

The uninstall snippet currently does not include an explicit `rm /etc/apt/keyrings/miru.gpg` command (only the `<Info>` note). No additional `rm` command is added — the `<Info>` text update is sufficient.

After this milestone, run `grep -rn "miru.gpg" docs/ snippets/` and expect the only remaining match to be the source URL `packages.mirurobotics.com/apt/miru.gpg` in `snippets/apt/setup.mdx`. If anything else surfaces, fix it before committing.

### Milestone 5 — Install-side keyring sharing note

Edit `snippets/apt/setup.mdx`. Immediately after the closing fenced code block (line 21), add a single-sentence `<Info>` note mirroring the uninstall note. Suggested wording:

    <Info>
        If you have already installed the Miru Agent or CLI via `apt`, the existing signing key at `/etc/apt/keyrings/miru-archive-keyring.gpg` is reused — both tools share the same keyring.
    </Info>

This places the note on both the agent install page and the CLI install page, since both pages render `setup.mdx`. It's symmetric with the existing `<Info>` notes in `uninstall.mdx` (line 14-16) and `cli/install.mdx` (line 89-91).

Note: `snippets/agent/install/apt.mdx` already shows a similar `<Info>` ("If this device has already been installed with the Miru Agent or CLI via `apt`, you can skip this step.") right after `<SetupApt />`. The new note inside `setup.mdx` is about the keyring specifically, not the whole step, so the two are complementary, not duplicative. If during refinement they read as redundant, prefer keeping the more specific (keyring) note inside `setup.mdx` and trim the outer one — but only after a visual review.

### Milestone 6 — Preflight

Run `./scripts/preflight.sh` from the repo root. Address any lint or spelling findings (most likely cspell may flag `dearmor`, `keyring`, `miruml`, etc., though these are likely already in `cspell.json` since the existing snippets passed). Re-run until output ends with the success line and exit code 0. Confirm preflight reports `clean` before publishing.

## Concrete Steps

All commands run from `/home/ben/miru/workbench1/repos/docs/` unless stated otherwise.

### Setup

1. Confirm the working branch:

       git branch --show-current
       # expect: docs/install-page-fixes

### Milestone 1: Two-step download

1. Edit `snippets/apt/setup.mdx` per Plan of Work Milestone 1.
2. Inspect rendered output (visual review of the MDX source — diff should show only the keyring step changes).

### Milestone 2: apt -> apt-get

1. Edit `snippets/agent/install/uninstall.mdx` line 4.
2. Edit `docs/developers/cli/install.mdx` line 79.
3. Verify only intended files changed:

       git diff --stat
       grep -rEn "(^| )(sudo )?apt (purge|install|update|upgrade|remove)\b" snippets docs

   Expect the only remaining `apt ` (without `-get`) matches to be in `rclone.md` and `docs/developers/agent/file-permissions.mdx` (out of scope).

### Milestone 3: Supported-platforms layout

1. Edit `snippets/agent/supported-platforms.mdx`: ensure a blank line between the last bullet and the trailing prose paragraph.

### Milestone 4: Keyring rename

1. Confirm all locations:

       grep -rn "miru.gpg" snippets docs

2. Apply the rename in `snippets/apt/setup.mdx` (lines 9, 10, 19), `snippets/agent/install/uninstall.mdx` (line 15), and `docs/developers/cli/install.mdx` (line 90). Do NOT change the URL on line 8 of `setup.mdx`.
3. Re-run grep:

       grep -rn "miru.gpg" snippets docs

   Expect only the URL `https://packages.mirurobotics.com/apt/miru.gpg` in `snippets/apt/setup.mdx` to remain.

### Milestone 5: Install-side note

1. Edit `snippets/apt/setup.mdx` to add the `<Info>` block immediately after the closing fenced code block.

### Milestone 6: Preflight

1. Run preflight:

       ./scripts/preflight.sh

   Expected: command exits 0 and prints "All documentation lint checks passed." (from `scripts/lint.sh`) plus successful audit, smoke tests, and shell tests.

2. If any check fails, fix the underlying issue and rerun. Common likely causes: cspell unknown word (add to `cspell.json` only if it is a legitimate term used elsewhere; otherwise reword), MDX prose lint warnings (address per the linter's message). Do not skip checks.

3. Confirm git working tree is clean (no unstaged changes from preflight artifacts):

       git status

## Validation and Acceptance

This is a documentation-only change. Acceptance is verified through automated lint and a visual review of the rendered MDX (or, equivalently, a careful reading of the diff because Mintlify renders MDX directly). There is no Mintlify dev-server step required by this plan; if the implementer chooses to render locally with `pnpm exec mint dev`, that is allowed but optional.

Acceptance criteria — each item must be observably true:

1. `snippets/apt/setup.mdx` shows the three-line download/dearmor/cleanup sequence and no `curl ... | sudo gpg --dearmor` pipe remains. The destination keyring file is `miru-archive-keyring.gpg`. Verify:

       grep -F "curl -fsSL https://packages.mirurobotics.com/apt/miru.gpg -o /tmp/miru.gpg.armor" snippets/apt/setup.mdx
       grep -F "sudo gpg --dearmor -o /etc/apt/keyrings/miru-archive-keyring.gpg /tmp/miru.gpg.armor" snippets/apt/setup.mdx
       grep -F "rm /tmp/miru.gpg.armor" snippets/apt/setup.mdx
       ! grep -F "| sudo gpg --dearmor" snippets/apt/setup.mdx

2. No Miru install/uninstall command in the docs uses bare `apt`:

       grep -rEn "sudo apt (purge|install|update|upgrade) miru" snippets docs

   expects zero matches.

3. The supported-platforms snippet renders the trailing caveat as its own paragraph. Verify by reading `snippets/agent/supported-platforms.mdx` and confirming a blank line separates the last bullet from the trailing sentence.

4. The keyring path is `/etc/apt/keyrings/miru-archive-keyring.gpg` everywhere except in the upstream URL:

       grep -rn "miru.gpg" snippets docs

   The only remaining match is the URL `https://packages.mirurobotics.com/apt/miru.gpg`.

5. The install-side `<Info>` note exists in `snippets/apt/setup.mdx` and references the shared keyring path.

6. **Preflight reports `clean`**: `./scripts/preflight.sh` exits 0 with no warnings. **Preflight must report `clean` before changes are published.**

The implementer should also do a quick visual review of the two affected user pages (`docs/developers/agent/install.mdx` and `docs/developers/cli/install.mdx`) by reading the composed MDX flow, confirming the new install-side `<Info>` appears once on each page (once via `setup.mdx`).

## Idempotence and Recovery

All edits are pure text edits. There are no destructive system operations. Each milestone can be re-attempted by reverting the corresponding commit (`git revert <sha>`) or by editing the file again. No external state is mutated.

If preflight fails after a commit, fix the underlying issue and add a NEW commit (do not amend).

If a future maintainer needs to roll back the entire change set, `git revert` the milestone commits in reverse order; nothing else needs to be undone.

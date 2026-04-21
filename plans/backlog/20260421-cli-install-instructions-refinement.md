# Refine CLI install instructions to match industry conventions

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | All edits happen here, on branch `docs/apt-repository-install-instructions`. |
| `cli` (mirurobotics/cli) | read-only | Referenced for the `install.sh` script and `--version=` flag behavior. No changes to `install.sh`. |

This plan lives in `docs/plans/` because every file edited is under the docs repo.

## Purpose / Big Picture

After this change, the Miru CLI install page at `https://docs.mirurobotics.com/docs/developers/cli/install` will mirror industry-standard install flows (Docker, Tailscale, GitHub CLI, Cloudflare, fly.io, Supabase, Bun):

- The page steers new users toward apt on Linux and calls out that the `curl … | sh` script is for containers/Alpine/CI environments without apt.
- The broken `apt-get upgrade miru` command (which would upgrade all packages and treat `miru` as a positional argument) is replaced with `apt-get install --only-upgrade miru`, which upgrades only the CLI.
- The reader can verify the key fingerprint before trusting the apt repo, pin the CLI to a specific version, and see a security note on the curl-piped script.
- Section order on the install page flows Install → Verify → Upgrade → Uninstall, matching the order of operations a user actually performs.

Verifiable outcome: after running the refreshed apt recipe on a fresh Debian-based VM, the user can run `miru version`, see the pinned version, and run `sudo apt-get install --only-upgrade miru` to move forward one version without upgrading unrelated packages.

## Progress

- [ ] M1: Fix the broken apt upgrade command.
- [ ] M2: Add primary-method callout and script security note to the install snippet.
- [ ] M3: Reorder the install page so Verify follows Install.
- [ ] M4: Modernize the apt recipe (supported distros note, prerequisite trim, fingerprint verification step).
- [ ] M5: Document version pinning for both apt and script installs.
- [ ] M6: Run preflight and validate behavior end-to-end.

Use timestamps when you complete steps. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

(Add entries as you go.)

- Decision: Use `<Note>` for the primary-method callout and `<Warning>` for the script-security note.
  Rationale: Matches existing usage — `<Note>` is used at the bottom of `snippets/references/cli/install.mdx` for the platform-support statement, and `<Warning>` is the strongest visual callout in the file set, appropriate for flagging remote-code execution risk.
  Date/Author: 2026-04-21 / author.

- Decision: Leave fingerprint and supported-distro specifics as TODO comments in MDX rather than inventing values.
  Rationale: The repo and `apt.mirurobotics.com` do not publish a canonical fingerprint or an authoritative distro list. Inventing either would mislead users; leaving a visible TODO preserves the structural improvement without fabricating facts.
  Date/Author: 2026-04-21 / author.

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

**Framework.** This is a Mintlify static docs site. Content lives in MDX files (Markdown + JSX components). Mintlify ships built-in callout components `<Note>`, `<Info>`, `<Warning>`, `<Tip>`, and `<Check>`, plus layout components `<Tabs>` and `<Tab>`.

**Files changed by this plan:**

- `snippets/references/cli/install.mdx` — the shared install snippet imported by the public CLI install page. Contains three `<Tab>` entries: `Linux (apt)`, `Linux (script)`, `macOS`.
- `docs/developers/cli/install.mdx` — the install page itself. Imports the snippet above, then has `## Upgrade`, `## Verify`, and `## Uninstall` sections.

**Branch.** All work happens on the already-checked-out branch `docs/apt-repository-install-instructions` in `/home/ben/miru/workbench3/docs`.

**How the page renders.** Mintlify renders each `<Tab>` as a tab header; the first tab (`Linux (apt)`) is the default. Content inside `<Tab>` is indented two spaces in the MDX source. Callouts inside a tab must be indented consistently with their surrounding content or the MDX parser rejects them.

**Existing patterns to mirror.** `docs/learn/devices/provision/api-keys.mdx` lines 110–129 already demonstrate version pinning for the agent: `sudo apt-get install -y miru-agent=<pinned-version>` with `apt-cache madison miru-agent` for listing. The same shape applies to the CLI. That file also uses the prerequisite line `sudo apt-get install -y apt-transport-https gnupg curl ca-certificates` (agent side, unchanged by this plan).

**What `install.sh` supports.** The upstream installer at `https://raw.githubusercontent.com/mirurobotics/cli/main/install.sh` accepts `--version=<semver>` via `sh -s -- --version=<semver>`. The plan references this but does **not** modify `install.sh`.

**Preflight entrypoint.** `./scripts/preflight.sh` in the docs repo. It runs:

    pnpm run test:lint                     # smoke tests of the lint harness
    LINT_FIX=0 ./tools/lint/scripts/lint.sh  # Go custom linter (check-only)
    ./tools/lint/scripts/covgate.sh        # Go coverage gate
    ./scripts/lint.sh                      # MDX custom linter + ESLint + CSpell + Mintlify OpenAPI
    ./scripts/audit.sh                     # pnpm audit
    bats pub/scripts/agent/check-miru-access_test.bats

**Preflight gate is load-bearing.** Preflight must report "All documentation lint checks passed." (and the other sub-commands must exit 0) before this work can be merged. This is a hard gate: do not open a PR if preflight reports any failures. The author of this plan has embedded that requirement here so it is visible to every downstream agent regardless of context loss.

**CSpell words.** `cspell.json` at the docs repo root contains a `words` list. If the plan introduces new jargon the linter flags (`.gpg.pub` components are already accepted via existing content; `fingerprint` is a normal English word), add the flagged word to `words` rather than adding in-file `// cSpell:ignore` comments.

## Plan of Work

The plan is organized into six small milestones. Each milestone ends with a single commit so the PR is reviewable as discrete units and bisectable.

### M1 — Fix the broken apt upgrade command

In `docs/developers/cli/install.mdx`, the `Linux (apt)` tab under `## Upgrade` contains:

    sudo apt-get update && sudo apt-get upgrade miru

`apt-get upgrade` ignores the trailing `miru` (it upgrades all packages). Replace with:

    sudo apt-get update && sudo apt-get install --only-upgrade miru

This upgrades only the `miru` package.

### M2 — Primary-method callout and script security note

**Primary-method callout.** In `snippets/references/cli/install.mdx`, insert a `<Note>` directly above the `<Tabs>` block stating that apt is the recommended path for Linux and that the install script is intended for containers, Alpine, and CI-without-apt environments. Prose should not invent supported distros; it points readers to the Linux (apt) tab.

**Script security note.** Inside the `Linux (script)` tab, immediately after the `curl -fsSL … | sh` code block and before the "The script requires `curl`, `tar`, …" sentence, add a `<Warning>` stating that the command pipes remote code into a shell and that apt is preferred when available. Use `<Warning>` (not `<Info>`) to match the strongest-callout convention elsewhere in the docs for security-flavored notices (see `snippets/references/cli/login.mdx` line 11).

### M3 — Section order: Install → Verify → Upgrade → Uninstall

In `docs/developers/cli/install.mdx`, move the `## Verify` block (currently lines 42–53, between `## Upgrade` and `## Uninstall`) to appear immediately after the `<Install />` import and before `## Upgrade`. The page import order becomes:

1. Frontmatter + `import Install`
2. `<Install />`
3. `## Verify`
4. `## Upgrade`
5. `## Uninstall`

### M4 — Modernize the apt recipe

Edits to the `Linux (apt)` tab in `snippets/references/cli/install.mdx`:

1. **Supported distros note.** Immediately below the tab's introductory sentence ("Install the prerequisites, add the Miru apt repository, then install the CLI."), add a single-sentence `<Info>` callout stating "Supported on recent Debian-based distributions." and insert an MDX comment `{/* TODO: list exact supported distros once the CLI release matrix is published */}` on the line after the callout so the user can fill in specifics without hunting. Do not invent version numbers.

2. **Prerequisite trim.** Change the first line of the code block from:

        sudo apt-get install -y apt-transport-https gnupg curl

   to:

        sudo apt-get install -y ca-certificates curl gnupg

   Rationale: `apt-transport-https` is a no-op on any apt version from the last several years (HTTPS transport is built in); `ca-certificates` is the actually-needed package so apt can verify the TLS chain to `apt.mirurobotics.com`.

3. **Fingerprint verification step.** Immediately after the `gpg --dearmor -o /usr/share/keyrings/miru-cli.gpg` line and before the `echo "deb …"` line, add a line:

        gpg --show-keys /usr/share/keyrings/miru-cli.gpg

   Follow the code block with a new `<Info>` callout stating that the command prints the imported key so the reader can confirm the fingerprint matches the published value. Because the published fingerprint is not available in this repo or on `apt.mirurobotics.com`, include an MDX comment `{/* TODO: publish fingerprint */}` inside the callout (or on the line above it) so a future author knows to fill in the expected value. Do not invent a fingerprint.

### M5 — Document version pinning

**apt pinning.** Directly after the existing `Linux (apt)` tab's `<Info>` about the shared signing key (the paragraph beginning "The CLI is signed by the same key…"), add a short subsection inside the tab. A natural pattern:

    To install or downgrade to a specific version, list available versions and pin the install:

        sudo apt list -a miru
        sudo apt-get install miru=<version>

No `##` heading inside the tab — keep it as a prose paragraph followed by a code block, matching how the agent API-keys page does it.

**Script pinning.** In the `Linux (script)` tab, directly after the `<Warning>` added in M2 and before the "The script requires `curl`, `tar`, …" paragraph, add a sentence explaining that the script accepts a `--version` flag, with an example:

    To install a specific version, pass `--version=<semver>` to the script:

        curl -fsSL https://raw.githubusercontent.com/mirurobotics/cli/main/install.sh \
          | sh -s -- --version=0.10.0

Use the same code-block style as the rest of the tab.

### M6 — Preflight and validation

Run `./scripts/preflight.sh` from the docs repo root. If any check reports a failure, resolve it (CSpell additions to `cspell.json` if a new word is flagged; MDX syntax fix if the custom linter or ESLint flags a structural issue) and re-run until the output ends with "All documentation lint checks passed." and the other sub-commands exit 0.

Render-check is covered by the lint pipeline: the custom linter and ESLint both parse MDX, so a malformed `<Tabs>` block or orphan JSX will fail CI.

## Concrete Steps

All commands run from `/home/ben/miru/workbench3/docs` (the docs repo root) unless stated otherwise.

### M1 — Fix the broken apt upgrade command

**Step 1.1.** Edit `docs/developers/cli/install.mdx`. Inside the `## Upgrade` section, `Linux (apt)` tab, change:

    sudo apt-get update && sudo apt-get upgrade miru

to:

    sudo apt-get update && sudo apt-get install --only-upgrade miru

**Step 1.2.** Commit.

    git add docs/developers/cli/install.mdx
    git commit -m "docs(cli): fix apt upgrade command to use --only-upgrade"

### M2 — Primary-method callout and script security note

**Step 2.1.** Edit `snippets/references/cli/install.mdx`. Immediately above the opening `<Tabs>` tag (line 3), insert a `<Note>` block:

    <Note>
      **apt is the recommended path on Linux.** The install script on the `Linux (script)` tab is provided for containers, Alpine, and CI environments without apt. On a Debian-based host, prefer the `Linux (apt)` tab.
    </Note>

Leave one blank line above and below the `<Note>`.

**Step 2.2.** In the same file, inside the `Linux (script)` tab, immediately after the closing triple-backtick of the `curl -fsSL … | sh` code block, insert:

    <Warning>
      This command pipes remote code into a shell. On systems with apt, prefer the `Linux (apt)` tab, which uses GPG-verified packages.
    </Warning>

Maintain the tab's existing two-space indent on each line of the `<Warning>`.

**Step 2.3.** Commit.

    git add snippets/references/cli/install.mdx
    git commit -m "docs(cli): add primary-method callout and script security warning"

### M3 — Reorder Install → Verify → Upgrade → Uninstall

**Step 3.1.** Edit `docs/developers/cli/install.mdx`. Cut the `## Verify` section (the heading, prose, code block, and the trailing "You can find the CLI release changelog…" paragraph — everything currently between `## Verify` and `## Uninstall`). Paste the cut content immediately after the `<Install />` line (and its blank line) and before `## Upgrade`.

After the edit, the file's top-level order must read:

    import Install from '/snippets/references/cli/install.mdx';

    <Install />

    ## Verify

    …

    ## Upgrade

    …

    ## Uninstall

    …

**Step 3.2.** Commit.

    git add docs/developers/cli/install.mdx
    git commit -m "docs(cli): reorder install page to Install -> Verify -> Upgrade -> Uninstall"

### M4 — Modernize the apt recipe

**Step 4.1.** Edit `snippets/references/cli/install.mdx`, `Linux (apt)` tab. Below the sentence "Install the prerequisites, add the Miru apt repository, then install the CLI." insert:

    <Info>
      Supported on recent Debian-based distributions.
    </Info>
    {/* TODO: list exact supported distros once the CLI release matrix is published */}

**Step 4.2.** In the same tab's code block, change:

    sudo apt-get install -y apt-transport-https gnupg curl

to:

    sudo apt-get install -y ca-certificates curl gnupg

**Step 4.3.** In the same code block, insert a new line after the `gpg --dearmor -o /usr/share/keyrings/miru-cli.gpg` line (and its continuation) and before the `echo "deb [signed-by=…]"` line:

    gpg --show-keys /usr/share/keyrings/miru-cli.gpg

**Step 4.4.** Immediately below the code block (before the existing `<Info>` about the shared signing key), add:

    <Info>
      {/* TODO: publish fingerprint */}
      `gpg --show-keys` prints the imported key. Confirm the fingerprint matches the published value before proceeding.
    </Info>

**Step 4.5.** Commit.

    git add snippets/references/cli/install.mdx
    git commit -m "docs(cli): modernize apt recipe (distros note, prereq trim, fingerprint step)"

### M5 — Version pinning

**Step 5.1.** Edit `snippets/references/cli/install.mdx`, `Linux (apt)` tab. After the existing `<Info>` about the shared signing key (the paragraph beginning "The CLI is signed by the same key…"), add a blank line and then:

    To install or downgrade to a specific version, list available versions and pin the install:

    ```bash
    sudo apt list -a miru
    sudo apt-get install miru=<version>
    ```

Maintain the tab's two-space indent on every line (including inside the fenced code block).

**Step 5.2.** In the `Linux (script)` tab of the same file, after the `<Warning>` added in M2 and before the paragraph beginning "The script requires `curl`, `tar`, …", add:

    To install a specific version, pass `--version=<semver>`:

    ```bash
    curl -fsSL https://raw.githubusercontent.com/mirurobotics/cli/main/install.sh \
      | sh -s -- --version=0.10.0
    ```

**Step 5.3.** Commit.

    git add snippets/references/cli/install.mdx
    git commit -m "docs(cli): document apt and script version pinning"

### M6 — Preflight and validation

**Step 6.1.** Run preflight from the docs repo root:

    ./scripts/preflight.sh

Expected tail of output:

    == MDX Prose ==
    ...
    All documentation lint checks passed.
    ...

**Step 6.2.** If CSpell flags a new word, add it to the `words` array in `cspell.json` (alphabetical order is not required but match existing case). Re-run `./scripts/preflight.sh` until clean.

**Step 6.3.** If ESLint or the custom MDX linter flags a structural problem (most commonly: inconsistent indentation of a callout inside a `<Tab>`, or an unterminated JSX element), fix the offending file and re-run preflight. Do not suppress lint rules to pass.

**Step 6.4.** Final commit if any fixes were needed in 6.2 or 6.3:

    git add cspell.json snippets/references/cli/install.mdx docs/developers/cli/install.mdx
    git commit -m "docs(cli): preflight fixups"

(If no fixes were needed, skip this step — preflight-clean is the acceptance criterion, not an empty commit.)

## Validation and Acceptance

**Test 1 — Preflight is clean.** From the docs repo root:

    ./scripts/preflight.sh ; echo "exit=$?"

Expected: final line `exit=0` and the output contains "All documentation lint checks passed."

**Test 2 — The broken upgrade command is gone.**

    grep -n "apt-get upgrade miru" docs/developers/cli/install.mdx

Expected: zero matches.

    grep -n "apt-get install --only-upgrade miru" docs/developers/cli/install.mdx

Expected: exactly one match.

**Test 3 — Primary-method callout is present.**

    grep -n "apt is the recommended path on Linux" snippets/references/cli/install.mdx

Expected: exactly one match, on a line that appears before the first `<Tabs>` tag. Confirm ordering:

    awk '/apt is the recommended path on Linux/{cb=NR} /<Tabs>/{ if (NR==1 || cb && cb<NR) { print "ok"; exit } }' snippets/references/cli/install.mdx

Expected: `ok`.

**Test 4 — Script security warning is present.**

    grep -n "pipes remote code into a shell" snippets/references/cli/install.mdx

Expected: exactly one match, inside the `Linux (script)` tab.

**Test 5 — Section order on the install page.** Check that `## Verify` appears before `## Upgrade`:

    grep -n "^## " docs/developers/cli/install.mdx

Expected output (order matters):

    <N1>:## Verify
    <N2>:## Upgrade
    <N3>:## Uninstall

with `N1 < N2 < N3`.

**Test 6 — Prerequisite package list is trimmed and correct.**

    grep -n "apt-transport-https" snippets/references/cli/install.mdx

Expected: zero matches.

    grep -n "ca-certificates curl gnupg" snippets/references/cli/install.mdx

Expected: exactly one match on the apt prerequisite line.

**Test 7 — Fingerprint verification step is present.**

    grep -n "gpg --show-keys /usr/share/keyrings/miru-cli.gpg" snippets/references/cli/install.mdx

Expected: exactly one match.

    grep -n "TODO: publish fingerprint" snippets/references/cli/install.mdx

Expected: exactly one match.

**Test 8 — Supported-distro phrasing and TODO are present.**

    grep -n "recent Debian-based distributions" snippets/references/cli/install.mdx

Expected: exactly one match.

    grep -n "TODO: list exact supported distros" snippets/references/cli/install.mdx

Expected: exactly one match.

**Test 9 — Version pinning is documented on both tabs.**

    grep -n "apt list -a miru" snippets/references/cli/install.mdx

Expected: exactly one match.

    grep -n "apt-get install miru=<version>" snippets/references/cli/install.mdx

Expected: exactly one match.

    grep -n -- "--version=0.10.0" snippets/references/cli/install.mdx

Expected: exactly one match.

**Test 10 — MDX renders.** The ESLint/custom-linter checks inside preflight (Test 1) cover parse-level MDX validity. As a quick local visual check, run `npx mint dev` from the docs repo root and open `http://localhost:3000/docs/developers/cli/install` — confirm each tab opens without a console error and the callouts render. This visual step is optional (preflight is the gate) but recommended for the three callouts added in M2 and M4.

**Preflight gate.** Preflight must report clean before the PR is opened. Preflight runs all the same commands CI runs (lint smoke tests, custom Go linter, coverage gate, MDX lint, audit, shell tests). This is a hard gate — do not open a PR if preflight reports any failures.

## Idempotence and Recovery

All changes in this plan are text edits to two MDX files (plus possibly `cspell.json`). Every step is safe to re-run:

- **Re-running an edit** produces the same file state because each edit replaces a specific literal string.
- **Re-running a milestone's commit step** after its edits were already committed is a no-op (`git add` followed by `git commit` with nothing staged errors out harmlessly; verify with `git status`).
- **If a milestone goes wrong mid-flight**, restore the affected file with `git checkout -- <path>` (only on files whose latest change is not yet committed) and redo the edits. This is not destructive beyond discarding the current milestone's in-progress edits.

There are no database migrations, deployments, or destructive filesystem operations. The only network operation is `./scripts/audit.sh` which queries the npm audit registry; it does not modify state.

Rollback for a completed milestone is `git revert <sha>` against that milestone's commit, which produces a new commit undoing the change. Because each milestone is a single commit, revert granularity matches the milestone granularity.

# Refine CLI install instructions to match industry conventions

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | All edits happen here, on branch `docs/apt-repository-install-instructions`. |
| `cli` (mirurobotics/cli) | read-only | Referenced for `install.sh` and its `--version=` flag. No changes to `install.sh`. |

This plan lives in `docs/plans/` because every file edited is under the docs repo.

## Purpose / Big Picture

After this change, the Miru CLI install page at `https://docs.mirurobotics.com/docs/developers/cli/install` will mirror industry-standard install flows (Docker, Tailscale, GitHub CLI, Cloudflare, fly.io, Supabase, Bun):

- The page steers new users toward apt on Linux and calls out that the `curl … | sh` script is for containers/Alpine/CI environments without apt.
- The broken `apt-get upgrade miru` command (upgrades all packages, treats `miru` as a positional argument) is replaced with `apt-get install --only-upgrade miru`.
- The reader can verify the key fingerprint, pin the CLI to a specific version, and see a security note on the curl-piped script.
- Section order on the install page flows Install → Verify → Upgrade → Uninstall.

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

- Decision: Use `<Note>` for the primary-method callout and `<Warning>` for the script-security note.
  Rationale: Matches existing usage — `<Note>` is used at the bottom of `snippets/references/cli/install.mdx` for the platform-support statement, and `<Warning>` is the strongest visual callout in the file set, appropriate for flagging remote-code execution risk.
  Date/Author: 2026-04-21 / author.

- Decision: Leave fingerprint and supported-distro specifics as TODO comments in MDX rather than inventing values.
  Rationale: The repo and `apt.mirurobotics.com` do not publish a canonical fingerprint or an authoritative distro list. Inventing either would mislead users; leaving a visible TODO preserves the structural improvement without fabricating facts.
  Date/Author: 2026-04-21 / author.

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

**Framework.** Mintlify static docs site. Content is MDX (Markdown + JSX). Built-in callouts: `<Note>`, `<Info>`, `<Warning>`, `<Tip>`, `<Check>`. Layout: `<Tabs>` / `<Tab>`.

**Files changed by this plan:**

- `snippets/references/cli/install.mdx` — shared install snippet imported by the public CLI install page. Three `<Tab>` entries: `Linux (apt)`, `Linux (script)`, `macOS`.
- `docs/developers/cli/install.mdx` — the install page. Imports the snippet, then has `## Upgrade`, `## Verify`, `## Uninstall` sections.

**Branch.** `docs/apt-repository-install-instructions` in `/home/ben/miru/workbench3/docs` (already checked out).

**MDX indent rules.** `<Tabs>` at column 0, `<Tab title=…>` at column 2, content *inside* a tab (prose, callouts, fenced code-block fences) at column 4. Body lines inside a fenced code block are NOT further indented beyond the fence. Callouts inside a tab must follow the 4-space indent on every line or the MDX parser rejects them.

**Existing patterns to mirror.** `docs/learn/devices/provision/api-keys.mdx` lines 110–129 already demonstrate version pinning for the agent (`sudo apt-get install -y miru-agent=<pinned-version>`, `apt-cache madison miru-agent`). The same shape applies to the CLI. That file also uses the prerequisite line `sudo apt-get install -y apt-transport-https gnupg curl ca-certificates` on the agent side (unchanged by this plan).

**What `install.sh` supports.** Upstream at `https://raw.githubusercontent.com/mirurobotics/cli/main/install.sh` accepts `--version=<semver>` via `sh -s -- --version=<semver>`. This plan does **not** modify `install.sh`.

**Preflight entrypoint.** `./scripts/preflight.sh` in the docs repo. It runs:

    pnpm run test:lint                     # smoke tests of the lint harness
    LINT_FIX=0 ./tools/lint/scripts/lint.sh  # Go custom linter (check-only)
    ./tools/lint/scripts/covgate.sh        # Go coverage gate
    ./scripts/lint.sh                      # MDX custom linter + ESLint + CSpell + Mintlify OpenAPI
    ./scripts/audit.sh                     # pnpm audit
    bats pub/scripts/agent/check-miru-access_test.bats

**Preflight gate is load-bearing.** Preflight must report "All documentation lint checks passed." (and the other sub-commands must exit 0) before this work can be merged. This is a hard gate: do not open a PR if preflight reports any failures. The author of this plan has embedded that requirement here so it is visible to every downstream agent regardless of context loss.

**CSpell words.** `cspell.json` at the docs repo root contains a `words` list. If the linter flags new jargon, add the word to `words` rather than using in-file `// cSpell:ignore` comments.

## Concrete Steps

All commands run from `/home/ben/miru/workbench3/docs` unless stated otherwise.

**Indent reminder.** Content inserted inside a `<Tab>` (the `<Warning>`, the supported-distros `<Info>`, version-pinning blocks) must be indented 4 spaces. Content at snippet top level (the primary-method `<Note>` in Step 2.1, above `<Tabs>`) is at column 0. The code blocks shown below illustrate textual content only; apply the indent the target location requires.

Each milestone ends with a single commit so the PR is reviewable as discrete units and bisectable.

### M1 — Fix the broken apt upgrade command

**Step 1.1.** In `docs/developers/cli/install.mdx`, `## Upgrade` section, `Linux (apt)` tab, change:

    sudo apt-get update && sudo apt-get upgrade miru

to:

    sudo apt-get update && sudo apt-get install --only-upgrade miru

Rationale: `apt-get upgrade` ignores the trailing `miru` (it upgrades all packages); `--only-upgrade` upgrades only the `miru` package.

**Step 1.2.** Commit.

    git add docs/developers/cli/install.mdx
    git commit -m "docs(cli): fix apt upgrade command to use --only-upgrade"

### M2 — Primary-method callout and script security note

**Step 2.1.** In `snippets/references/cli/install.mdx`, immediately above the opening `<Tabs>` tag (line 3), insert at column 0 (2-space indent for the prose body inside the `<Note>`):

    <Note>
      **apt is the recommended path on Linux.** The install script on the `Linux (script)` tab is provided for containers, Alpine, and CI environments without apt. On a Debian-based host, prefer the `Linux (apt)` tab.
    </Note>

Leave one blank line above and below.

**Step 2.2.** In the same file, inside the `Linux (script)` tab, immediately after the closing triple-backtick of the `curl -fsSL … | sh` code block, insert (4-space indent on each line to match the tab):

    <Warning>
      This command pipes remote code into a shell. On systems with apt, prefer the `Linux (apt)` tab, which uses GPG-verified packages.
    </Warning>

Use `<Warning>` (not `<Info>`) to match the strongest-callout convention for security-flavored notices (see `snippets/references/cli/login.mdx` line 11).

**Step 2.3.** Commit.

    git add snippets/references/cli/install.mdx
    git commit -m "docs(cli): add primary-method callout and script security warning"

### M3 — Reorder Install → Verify → Upgrade → Uninstall

**Step 3.1.** In `docs/developers/cli/install.mdx`, cut the three existing elements of the current `## Verify` section (currently between `## Upgrade` and `## Uninstall`): (1) the `## Verify` heading, (2) the prose paragraph plus the `miru version` fenced code block, (3) the trailing "You can find the CLI release changelog…" paragraph. Paste immediately after the `<Install />` line (and its blank line) and before `## Upgrade`.

After the edit, the file's top-level order reads:

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

**Step 4.1.** In `snippets/references/cli/install.mdx`, `Linux (apt)` tab, below the sentence "Install the prerequisites, add the Miru apt repository, then install the CLI.", insert:

    <Info>
      Supported on recent Debian-based distributions.
    </Info>
    {/* TODO: list exact supported distros once the CLI release matrix is published */}

Do not invent version numbers.

**Step 4.2.** In the same tab's code block, change:

    sudo apt-get install -y apt-transport-https gnupg curl

to:

    sudo apt-get install -y ca-certificates curl gnupg

Rationale: `apt-transport-https` is a no-op on any apt version from the last several years; `ca-certificates` is the actually-needed package so apt can verify the TLS chain to `apt.mirurobotics.com`.

**Step 4.3.** In the same code block, insert a new line after the `gpg --dearmor -o /usr/share/keyrings/miru-cli.gpg` line (and its continuation) and before the `echo "deb [signed-by=…]"` line:

    gpg --show-keys /usr/share/keyrings/miru-cli.gpg

**Step 4.4.** Immediately below the code block (before the existing `<Info>` about the shared signing key), add:

    <Info>
      {/* TODO: publish fingerprint */}
      `gpg --show-keys` prints the imported key. Confirm the fingerprint matches the published value before proceeding.
    </Info>

Do not invent a fingerprint.

**Step 4.5.** Commit.

    git add snippets/references/cli/install.mdx
    git commit -m "docs(cli): modernize apt recipe (distros note, prereq trim, fingerprint step)"

### M5 — Version pinning

**Step 5.1.** In `snippets/references/cli/install.mdx`, `Linux (apt)` tab. Anchor: the existing `<Info>` beginning "The CLI is signed by the same key…" (NOT the fingerprint `<Info>` added in M4.4). After that shared-signing-key `<Info>`, add a blank line and then (4-space indent on each line; body lines inside the fence are NOT further indented):

    To install or downgrade to a specific version, list available versions and pin the install:

    ```bash
    sudo apt list -a miru
    sudo apt-get install miru=<version>
    ```

No `##` heading inside the tab — prose paragraph followed by a code block, matching the agent API-keys page.

Final apt-tab content order after M4 + M5: intro prose → supported-distros `<Info>` → apt code block (with `gpg --show-keys` line) → fingerprint `<Info>` (with TODO) → shared-signing-key `<Info>` → version-pinning prose → version-pinning code block.

**Step 5.2.** In the `Linux (script)` tab, after the `<Warning>` added in M2 and before the paragraph beginning "The script requires `curl`, `tar`, …", add:

    To install a specific version, pass `--version=<semver>`:

    ```bash
    curl -fsSL https://raw.githubusercontent.com/mirurobotics/cli/main/install.sh \
      | sh -s -- --version=0.10.0
    ```

`0.10.0` is illustrative. Substitute a version that exists in the apt repo — use `sudo apt list -a miru` on a host with the repo configured. If `0.10.0` has not shipped by the time this plan is executed, update the literal to a real version before committing. Test 9's `grep` assertion below checks for the literal `--version=0.10.0`, so update the test to match whatever version literal ends up in the snippet.

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

**Step 6.2.** If CSpell flags a new word, add it to the `words` array in `cspell.json` (alphabetical order not required, match existing case). Re-run `./scripts/preflight.sh` until clean.

**Step 6.3.** If ESLint or the custom MDX linter flags a structural problem (most commonly: inconsistent indentation of a callout inside a `<Tab>`, or an unterminated JSX element), fix the offending file and re-run preflight. Do not suppress lint rules to pass.

**Step 6.4.** Final commit *only if* fixes were needed in 6.2 or 6.3:

    git status
    git add <only-the-files-git-status-shows-as-modified>
    git commit -m "docs(cli): preflight fixups"

If `git status` shows no changes, skip — preflight-clean is the acceptance criterion, not an empty commit.

Render-check is covered by the lint pipeline: the custom linter and ESLint both parse MDX, so a malformed `<Tabs>` block or orphan JSX will fail CI.

## Validation and Acceptance

**Test 1 — Preflight is clean.**

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

**Test 5 — Section order on the install page.**

    grep -n "^## " docs/developers/cli/install.mdx

Expected (order matters):

    <N1>:## Verify
    <N2>:## Upgrade
    <N3>:## Uninstall

The relative order of the three headings is what is checked — absolute line numbers do not matter.

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

    grep -n "Confirm the fingerprint matches" snippets/references/cli/install.mdx

Expected: exactly one match (the prose body of the fingerprint `<Info>` callout).

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

**Test 10 — MDX renders.** Preflight (Test 1) covers parse-level MDX validity. As an optional visual check, run `npx mint dev` from the docs repo root and open `http://localhost:3000/docs/developers/cli/install` — confirm each tab opens without a console error and the callouts render. Preflight is the gate; the visual step is recommended for the callouts added in M2 and M4.

**Test 11 — Tab structure preserved.**

    grep -c "<Tab title=" snippets/references/cli/install.mdx

Expected: `3` (Linux (apt), Linux (script), macOS).

    grep -c "<Tabs>" snippets/references/cli/install.mdx
    grep -c "</Tabs>" snippets/references/cli/install.mdx

Expected: each returns `1`.

**Preflight gate.** Preflight must report clean before the PR is opened. Preflight runs all the same commands CI runs (lint smoke tests, custom Go linter, coverage gate, MDX lint, audit, shell tests). This is a hard gate — do not open a PR if preflight reports any failures.

## Idempotence and Recovery

All changes are text edits to two MDX files (plus possibly `cspell.json`). Every step is safe to re-run:

- **Re-running an edit** produces the same file state because each edit replaces a specific literal string.
- **Re-running a milestone's commit step** after its edits were already committed is a no-op (`git add` followed by `git commit` with nothing staged errors out harmlessly; verify with `git status`).
- **If a milestone goes wrong mid-flight**, restore the affected file with `git checkout -- <path>` (only on files whose latest change is not yet committed) and redo the edits.

There are no database migrations, deployments, or destructive filesystem operations. The only network operation is `./scripts/audit.sh` (queries the npm audit registry; no state change).

Rollback for a completed milestone is `git revert <sha>` against that milestone's commit. Because each milestone is a single commit, revert granularity matches milestone granularity.

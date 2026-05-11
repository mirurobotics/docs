# Apply three prose fixes from provisioning docs review

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit three MDX files (one snippet, two pages) to correct technically inaccurate prose introduced or surfaced during the activate-to-provision refactor. |

This plan lives in `docs/plans/` because all edits are confined to the `docs` repo. The work is executed on the existing branch `refactor/activate-to-provision`.

## Purpose / Big Picture

After this change, a reader of the Miru Agent install page and the device provisioning docs reads prose that:

- Correctly attributes the v0.9.0 cutover to the **provisioning** method (tokens vs. legacy script), not to the install/upgrade workflow. The Debian package install/upgrade flow itself works for any agent version.
- Describes provisioning in the canonical direction: registering a physical machine with the Miru control plane (not the reverse).
- Uses the retained "provision" verb consistently in the provisioning-tokens action list, mirroring `provisioning-script.mdx`. The retired "activate" terminology is gone from that bullet.

These are prose-only corrections. No code, links, schemas, or page structure change.

## Progress

- [ ] Milestone 1: Correct the install/upgrade Danger blocks in `docs/snippets/agent/install/install.mdx` and `docs/developers/agent/install.mdx`.
- [ ] Milestone 2: Correct the provisioning-direction sentence in `docs/learn/devices/provision/overview.mdx`.
- [ ] Milestone 3: Replace the tautological "Activates the device" bullet in `docs/learn/devices/provision/provisioning-tokens.mdx`.
- [ ] Milestone 4: Run `./scripts/preflight.sh`; address any findings; confirm `clean` and the working tree is empty.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

- Decision: Keep two Danger blocks (one in the install snippet, one in the agent/install.mdx Upgrade section) rather than consolidating into a single shared snippet.
  Rationale: The two locations need different verbs ("install"/"upgrade") and the wording in each is short. Consolidating would either require a new snippet just for one sentence or force one location to read awkwardly. Two parallel, hand-edited blocks is simpler and the maintenance burden is negligible.
  Date/Author: 2026-05-11 / planner.

- Decision: Reframe the v0.9.0 Danger to be about **provisioning tokens** (not install/upgrade tooling), and reword the link target to a method choice rather than a forced fallback.
  Rationale: The original wording was technically wrong. apt/dpkg install/upgrade workflows work for any agent version; what changes at v0.9.0 is that the provisioning-tokens method becomes available and the legacy provisioning script stops supporting newer versions (see `docs/learn/devices/provision/provisioning-script.mdx` lines 7 and 22, which already state this correctly).
  Date/Author: 2026-05-11 / planner.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The `docs` repo is a Mintlify documentation site. MDX files compose pages from snippets via `import` statements. The current branch is `refactor/activate-to-provision`, a prose/terminology refactor renaming "activate" to "provision" across the docs surface.

Files touched by this plan, all paths relative to repo root `/home/ben/miru/workbench3/repos/docs/`:

- `docs/snippets/agent/install/install.mdx` — Top-level Agent install snippet, imported by the user-facing install page via `import Install from '/snippets/agent/install/install.mdx'`. Contains the install-context Danger block at lines 8-10.
- `docs/developers/agent/install.mdx` — User-facing Agent install page. Renders `<Install />` for the Install section and has a separate near-duplicate Danger block at lines 25-27 inside its own Upgrade section.
- `docs/learn/devices/provision/overview.mdx` — Provisioning overview page. The first content sentence (line 9) defines what provisioning is.
- `docs/learn/devices/provision/provisioning-tokens.mdx` — Provisioning-tokens page. The "Provisioning performs the following actions:" list at lines 76-79 enumerates what the provision command does.
- `docs/learn/devices/provision/provisioning-script.mdx` — Read-only reference. Its parallel bullet at line 38 (`**Provisions the device** - registers the agent with the Miru control plane`) is the wording we mirror in Milestone 3.

Definitions:

- "Provisioning" — creating a device record in the Miru control plane and registering a specific physical machine (running `miru-agent`) with that record, so the machine can receive deployments.
- "Provisioning tokens" — a v0.9.0+ method where a token issued from the Miru dashboard is passed to `miru-agent provision` on the device, enabling programmatic / at-scale provisioning without the legacy script.
- "Legacy provisioning script" — the pre-v0.9.0 shell-script flow documented at `/learn/devices/provision/provisioning-script`. It only supports installing `miru-agent` versions that precede `v0.9.0`.
- "preflight" — the repo-local script `./scripts/preflight.sh` that runs lint smoke tests, Go lint, Go coverage, MDX/CSpell/OpenAPI/ESLint lint, security audit, and shell tests. Mirrors CI.

## Plan of Work

Three independent prose fixes, one milestone each, in the order below. A fourth milestone runs preflight.

### Milestone 1 — Reframe install/upgrade Danger blocks

Two near-duplicate Danger blocks both incorrectly claim apt/dpkg install/upgrade workflows are version-gated at v0.9.0. The actual gate is on the **provisioning** method.

Edit `docs/snippets/agent/install/install.mdx`, lines 8-10. Replace:

    <Danger>
        This install workflow is only supported for versions `v0.9.0` or later. To upgrade to a version that precedes `v0.9.0`, you must use the legacy [provisioning script](/learn/devices/provision/provisioning-script) method.
    </Danger>

with:

    <Danger>
        [Provisioning tokens](/learn/devices/provision/provisioning-tokens) require Miru Agent `v0.9.0` or later. To install and provision a version that precedes `v0.9.0`, use the legacy [provisioning script](/learn/devices/provision/provisioning-script) instead.
    </Danger>

Edit `docs/developers/agent/install.mdx`, lines 25-27. Replace:

    <Danger>
        This upgrade workflow is only supported for versions `v0.9.0` or later. To upgrade to a version that precedes `v0.9.0`, you must use the legacy [provisioning script](/learn/devices/provision/provisioning-script) method.
    </Danger>

with:

    <Danger>
        [Provisioning tokens](/learn/devices/provision/provisioning-tokens) require Miru Agent `v0.9.0` or later. If you are running a version that precedes `v0.9.0`, use the legacy [provisioning script](/learn/devices/provision/provisioning-script) to provision instead.
    </Danger>

Both blocks keep the same link target (`/learn/devices/provision/provisioning-script`) and the same `<Danger>` admonition type. The version cutover is now attributed to the provisioning method, not the package-manager workflow. Two blocks remain (see Decision Log) but they are no longer falsely claiming the install/upgrade tooling is version-gated.

Commit at end of milestone:

    docs(provision): correct install/upgrade Danger blocks to gate provisioning method on v0.9.0

### Milestone 2 — Fix provisioning direction sentence

Edit `docs/learn/devices/provision/overview.mdx`, line 9. Replace:

    Provisioning a device is the process of creating a device in Miru and registering it with a physical machine for your workspace. Once provisioned, a device is ready to receive deployments.

with:

    Provisioning a device is the process of creating a device in Miru and registering a physical machine with the Miru control plane. Once provisioned, a device is ready to receive deployments.

Only the first sentence changes; the trailing "Once provisioned…" sentence is preserved verbatim. The change inverts the registration direction so it matches the canonical model (the agent on the physical machine registers with the control plane, not vice versa) and matches the parallel wording already used in `provisioning-tokens.mdx` and `provisioning-script.mdx`.

Commit at end of milestone:

    docs(provision): fix provisioning-direction sentence in overview

### Milestone 3 — Replace "Activates the device" bullet

Edit `docs/learn/devices/provision/provisioning-tokens.mdx`, line 79. Replace:

    1. **Activates the device** - activates the device by registering the agent with the Miru control plane.

with:

    1. **Provisions the device** - registers the agent with the Miru control plane.

This removes the tautology ("activates the device by … activating") and drops the retired "activate" terminology from this bullet. The replacement mirrors `docs/learn/devices/provision/provisioning-script.mdx:38` exactly (`**Provisions the device** - registers the agent with the Miru control plane`).

Markdown numbering: the file uses `1.` for both bullets in this list (lines 78 and 79), which Markdown auto-renumbers; do not change the leading `1.`.

Commit at end of milestone:

    docs(provision): replace tautological "Activates the device" bullet with provisions wording

### Milestone 4 — Preflight

Run `./scripts/preflight.sh` from `/home/ben/miru/workbench3/repos/docs/`. Address any findings (most likely none for a prose-only change) and re-run until exit code 0. Confirm the working tree is clean. No commit is expected from preflight itself, unless a cspell or lint config needs an entry — in that case, commit it as `docs(provision): adjust lint config for review-fix wording` or similar.

## Concrete Steps

All commands run from `/home/ben/miru/workbench3/repos/docs/` unless stated otherwise.

### Setup

1. Confirm the working branch:

       git branch --show-current
       # expect: refactor/activate-to-provision

2. Confirm the working tree is clean:

       git status
       # expect: nothing to commit, working tree clean

### Milestone 1: Install/upgrade Danger blocks

1. Edit `docs/snippets/agent/install/install.mdx` per Plan of Work Milestone 1.
2. Edit `docs/developers/agent/install.mdx` per Plan of Work Milestone 1.
3. Verify exactly two files changed and the bare phrase "install workflow is only supported" / "upgrade workflow is only supported" is gone:

       git diff --stat
       grep -rn "install workflow is only supported" docs/
       grep -rn "upgrade workflow is only supported" docs/

   Both greps should return zero matches.

4. Commit:

       git add docs/snippets/agent/install/install.mdx docs/developers/agent/install.mdx
       git commit -m "docs(provision): correct install/upgrade Danger blocks to gate provisioning method on v0.9.0"

### Milestone 2: Overview direction sentence

1. Edit `docs/learn/devices/provision/overview.mdx` line 9 per Plan of Work Milestone 2.
2. Verify:

       grep -n "registering a physical machine with the Miru control plane" docs/learn/devices/provision/overview.mdx
       # expect: one match on line 9
       grep -n "registering it with a physical machine for your workspace" docs/learn/devices/provision/overview.mdx
       # expect: zero matches

3. Commit:

       git add docs/learn/devices/provision/overview.mdx
       git commit -m "docs(provision): fix provisioning-direction sentence in overview"

### Milestone 3: Provisioning-tokens bullet

1. Edit `docs/learn/devices/provision/provisioning-tokens.mdx` line 79 per Plan of Work Milestone 3.
2. Verify the parallel wording matches `provisioning-script.mdx`:

       grep -n "Provisions the device" docs/learn/devices/provision/provisioning-tokens.mdx docs/learn/devices/provision/provisioning-script.mdx
       # expect: one match in each file, both reading "**Provisions the device** - registers the agent with the Miru control plane"
       grep -n "Activates the device" docs/learn/devices/provision/provisioning-tokens.mdx
       # expect: zero matches

3. Commit:

       git add docs/learn/devices/provision/provisioning-tokens.mdx
       git commit -m "docs(provision): replace tautological \"Activates the device\" bullet with provisions wording"

### Milestone 4: Preflight

1. Run preflight:

       ./scripts/preflight.sh

   Expected: exits 0, prints "All documentation lint checks passed." plus successful audit and shell tests.

2. If any check fails, fix the underlying cause (likely a cspell or prose lint flag on the new wording) and rerun. Do not skip checks.

3. Confirm the working tree is clean:

       git status
       # expect: nothing to commit, working tree clean

## Validation and Acceptance

This is a prose-only documentation change. Acceptance is verified through visual review of the rendered prose (Mintlify renders MDX directly, so a careful read of the diff is equivalent), preserved link integrity, and a clean preflight run.

Acceptance criteria — each must be observably true:

1. The two Danger blocks no longer claim the install/upgrade workflow is version-gated. They now claim provisioning tokens require v0.9.0+ and point to the legacy script as the alternative provisioning method. Verify:

       grep -rn "install workflow is only supported" docs/
       grep -rn "upgrade workflow is only supported" docs/
       # both: zero matches
       grep -rn "Provisioning tokens.*require Miru Agent" docs/snippets/agent/install/install.mdx docs/developers/agent/install.mdx
       # expect: one match in each file

2. The link target `/learn/devices/provision/provisioning-script` is preserved in both Danger blocks (no broken or changed links):

       grep -rn "/learn/devices/provision/provisioning-script" docs/snippets/agent/install/install.mdx docs/developers/agent/install.mdx
       # expect: at least one match in each file
       grep -rn "/learn/devices/provision/provisioning-tokens" docs/snippets/agent/install/install.mdx docs/developers/agent/install.mdx
       # expect: at least one match in each file (new link added by this plan)

3. `docs/learn/devices/provision/overview.mdx` line 9 reads "creating a device in Miru and registering a physical machine with the Miru control plane".

4. `docs/learn/devices/provision/provisioning-tokens.mdx` no longer contains the substring "Activates the device", and contains the bullet "**Provisions the device** - registers the agent with the Miru control plane" matching the parallel bullet in `provisioning-script.mdx`.

5. **Preflight reports `clean`**: `./scripts/preflight.sh` exits 0 with no warnings. **Preflight must report `clean` before changes are delivered (pushed / opened for review).**

6. Visual review of the rendered prose (read the final state of each of the four edited files top-to-bottom) reads naturally and contains no remaining references to "activates the device" or to a version-gated install/upgrade workflow.

## Idempotence and Recovery

All edits are pure text edits. There are no destructive operations. Each milestone is one commit and can be reverted independently with `git revert <sha>`. If preflight fails after a milestone commit, fix the underlying cause and add a NEW commit (do not amend a signed commit). If a future maintainer needs to roll back the entire change set, `git revert` the four milestone commits in reverse order; nothing else needs to be undone.

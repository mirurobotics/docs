# Document the MQTT-port-443 change: agent v0.10.3 + v0.7.2 changelog entries and the port doc

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.


## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/ben/miru/workbench1/repos/docs`) | read-write | Add two version sections to `docs/changelog/agent.mdx` (`v0.10.3` at the top, `v0.7.2` between the existing `v0.8.0` and `v0.7.1` sections) and update the MQTT port statement in `docs/developers/agent/security.mdx`. Touch `cspell.json` only if CSpell flags a new word. |
| `mirurobotics/agent` (`/home/ben/miru/workbench1/repos/agent`) | read-only | Source of truth for the `v0.10.2..v0.10.3` and `v0.7.1..v0.7.2` deltas and the tag dates. Do not clone, modify, branch, or commit there. |

This plan lives in `docs/plans/` because every write is in the `docs` repo (the public changelog and a customer-facing developer page).

Working branch: `docs/agent-443-changelog-and-port` (already created from `main`). Do not create another branch. Base branch for the PR: `main`.


## Purpose / Big Picture

The Miru Agent now connects to the MQTT broker over TLS on port `443` instead of `8883`, so a device behind a firewall that only allows outbound `443`/HTTPS can still keep its real-time connection to the control plane. This change shipped on the main line as `v0.10.3` and was backported to the `v0.7` line as `v0.7.2`.

After this change a reader can:

1. Open `/changelog/agent` and see `# v0.10.3` (*September 17, 2026*) at the top of the page describing the port-`443` switch, with the existing `# v0.10.2` section unchanged below it.
2. Scroll down and find `# v0.7.2` (*September 17, 2026*) placed in version-descending order between `# v0.8.0` and `# v0.7.1`, described as a backport of the port-`443` change plus a security fix that moves the MQTT TLS stack off the vulnerable `rustls-webpki` dependency.
3. Read `/developers/agent/security` and see the "MQTT over TLS" bullet state that current agents use port `443` (a device only needs outbound `443`/HTTPS through a firewall), while noting that older agent versions still use `8883`, which the broker continues to accept.


## Progress

- [ ] Milestone 1: Re-verify the two agent deltas and tag dates (`v0.10.2..v0.10.3`, `v0.7.1..v0.7.2`).
- [ ] Milestone 2: Insert the `# v0.10.3` and `# v0.7.2` sections into `docs/changelog/agent.mdx`; commit.
- [ ] Milestone 3: Update the MQTT port statement in `docs/developers/agent/security.mdx`; commit.
- [ ] Milestone 4: Run `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate` (add a CSpell word if flagged); commit any fix.
- [ ] Milestone 5: Push, open the draft PR, and drive CI to CLEAN before leaving draft.


## Surprises & Discoveries

(Add entries as work proceeds.)


## Decision Log

- Decision: v0.10.3 entry documents only the port-`443` change. / Rationale: The `v0.10.2..v0.10.3` delta also contains Windows-platform groundwork (`#232`, `#234`, `#236`) and internal refactors (`#231` crypt, `#233`/`#240` path/test portability). Windows support is foundational scaffolding, not a shippable, announceable capability; the rest are internal. The existing changelog lists only finished, user-facing changes. See open question in Context and Orientation. / Date/Author: 2026-09-17, author.
- Decision: The 443 bullet lives under `## Improvements`. / Rationale: It changes existing connection behavior for broader network compatibility rather than adding a new capability; matches the `## Improvements` usage already in the `v0.8.0` and `v0.7.0` sections. Not a `## Breaking changes` item for customers, whose only action is allowing outbound `443` (usually already open). / Date/Author: 2026-09-17, author.
- Decision: v0.7.2 is placed in version-descending order between `v0.8.0` and `v0.7.1` with its real *September 17, 2026* date and a one-line backport note. / Rationale: The file already runs down to `v0.5.0`, so version order is the established scheme; the backport note explains why a `0.7.x` entry carries a later date than the `0.8.x`/`0.10.x` entries above it. / Date/Author: 2026-09-17, author.


## Outcomes & Retrospective

(Summarize at completion.)


## Context and Orientation

The docs site is a Mintlify site; published content lives under `docs/`. MDX means Markdown plus JSX. This repo has no `AGENTS.md` or `CLAUDE.md`; conventions are captured in `docs/plans/completed/20260906-agent-v0.10.2-changelog.md`, the most recent changelog-entry plan, which this plan mirrors.

**The agent changelog file.** `docs/changelog/agent.mdx` is one newest-first file. Its structure (verified 2026-09-17):

- YAML frontmatter (`title: "Agent"`, `mode: "center"`, a description), then a blank line.
- Each version is a `# vX.Y.Z` heading, a blank line, an italic date `*Month D, YYYY*` (full month name, no ordinal), an optional one-paragraph summary, then the non-empty subset of `## Breaking changes`, `## Features`, `## Fixes`, `## Improvements` (in that order), each with hyphen bullets in sentence case with no trailing period.
- A `---` separator with exactly one blank line before and after it terminates every section except the last (`v0.5.0`).

Current headings top to bottom: `v0.10.2`, `v0.10.1`, `v0.10.0`, `v0.9.0`, `v0.8.1`, `v0.8.0`, `v0.7.1`, `v0.7.0`, `v0.6.1`, `v0.6.0`, `v0.5.1`, `v0.5.0`. The top entry is `# v0.10.2` (*September 6, 2026*). Insertion points: `v0.10.3` goes above `v0.10.2`; `v0.7.2` goes after the `v0.8.0` section's `---` separator and before `# v0.7.1`.

Navigation is already wired (`docs/docs.json` lists `"changelog/agent"`). Do not add a page or a `docs.json` entry, and do not edit `docs/developers/agent/versions.mdx` (the supported-versions table already covers the `v0.10.x` and `v0.7.x` lines). Changelogs are historical records: do not rewrite existing entries.

**The security page.** `docs/developers/agent/security.mdx` has a "Transport security" section (around line 111). Line 116 currently reads:

    - **MQTT over TLS** — the persistent MQTT connection uses TLS on port 8883 (the standard MQTT over TLS port)

This is the only place in `docs/` that names a port for the MQTT connection (confirmed by `grep -rn 8883 docs/`). Two other pages mention "MQTT over TLS" without a port and need no change: `docs/developers/agent/overview.mdx:19` ("all communication is outbound-only via HTTPS or MQTT over TLS") and `docs/developers/agent/architecture.mdx:18-22` (the "MQTT over TLS" subsection). Note both in the Progress log as reviewed with no change.

**Source of truth (agent repo, verified 2026-09-17).** A local checkout exists at `/home/ben/miru/workbench1/repos/agent` on branch `release/v0.7`; the `v0.10.3` and `v0.7.2` tags are present.

- `git -C /home/ben/miru/workbench1/repos/agent log --oneline v0.10.2..v0.10.3` — the only user-facing, announceable commit is `feat(mqtt): connect to the broker on port 443 (#245)` (`59eef65d`). The rest is Windows-platform groundwork (`feat(platform): per-OS default data root and log dir #232`, `feat(windows): cfg-gate unix-only APIs #234`, `build(windows): MSI packaging foundation #236`), internal refactors (`refactor(crypt): migrate RSA from openssl to aws-lc-rs #231`, `refactor(filesys) #240`, `refactor(test) #239`), test/CI/deps/version-bump commits. Tag `v0.10.3` is dated 2026-09-17 (UTC calendar day September 17, 2026).
- `git -C /home/ben/miru/workbench1/repos/agent log --oneline v0.7.1..v0.7.2` — three commits: `feat(mqtt): connect to the broker on port 443 (v0.7.2) (#248)` (`a2fc4644`), `fix(mqtt): route TLS via native-tls to clear rustls-webpki RUSTSEC advisories` (`785b3c0f`), and a `build:` toolchain pin (`844cee69`, internal). The `fix(mqtt)` commit body names the cleared advisories: RUSTSEC-2026-0098, RUSTSEC-2026-0099, RUSTSEC-2026-0104. Tag `v0.7.2` is dated 2026-09-17.

**User-facing** means operator-visible behavior (connection ports, security posture). Omit internal refactors, test/CI/deps/chore, version bumps, and platform groundwork that does not yet ship a usable capability.

**Framing constraint (important).** The docs are customer-facing. Do NOT describe Miru's internal TLS-passthrough proxy or per-environment DNS cutover (these appear in the agent commit bodies but are infra detail; customers see only a hostname and port). Present `443` as the port current agents use and the minimum a restrictive firewall must allow, and state that `8883` remains valid for older agent versions and is still accepted by the broker. Do not phrase the security page as a hard "443 only" requirement that would imply existing `8883` deployments must immediately restrict to `443`.

**Open question for the user (record the answer in the Decision Log at implementation time).** This plan assumes the `v0.10.3` changelog entry documents ONLY the port-`443` change and omits the Windows groundwork and internal refactors (see Decision Log). If the user wants partial Windows support announced now, that is a separate, larger content decision and would need accurate scope from the agent team — flag it rather than guessing. A secondary, minor question: whether to print the three raw `RUSTSEC-2026-00xx` identifiers in the customer changelog. This plan includes them (they are public advisory IDs and precise); drop them to "known `rustls-webpki` security advisories" if the user prefers a softer note.

**Lint and CI.** From `/home/ben/miru/workbench1/repos/docs`: `pnpm run test:lint` (→ `./tests/test-lint.sh`); `./scripts/lint.sh` runs the Go MDX linter (`headingcase` allows `# vX.Y.Z` and `## Features`/`## Fixes`/`## Improvements`/`## Breaking changes`; `nodoubledash` flags two consecutive hyphens `--` in prose spans — single hyphens in `native-tls`, `rustls-webpki`, and `RUSTSEC-2026-0098` are fine, and any literal flag/port pair must sit inside backticks), ESLint MDX, CSpell (`cspell.json` at the repo root), and `mint openapi-check`. Validate with `pnpm run validate` (`mint validate`). Full local preflight is `./scripts/preflight.sh`. CI is `.github/workflows/ci.yml` (jobs: `lint`, `audit`, `shell-tests`, and the two custom-linter jobs). **`CLEAN`** means every GitHub Actions check on the pushed branch head is green; a green local preflight is useful but not sufficient. Do not touch the Dependabot config or the audit ignore list — the known unpatchable Mintlify transitive advisories are deliberately left open.

**CSpell note.** In `cspell.json` `words`, `rustls`, `webpki`, and `MQTT` are already present; `firewall`, `outbound`, and `broker` are ordinary English words CSpell's base dictionary knows (already used elsewhere in `docs/`). `RUSTSEC` and the `native`/`tls` split of `native-tls` may or may not be flagged. Do not pre-add words; run lint and add only what CSpell actually reports, in sorted position.


## Plan of Work

Two files change: `docs/changelog/agent.mdx` (two inserts) and `docs/developers/agent/security.mdx` (one line). Indentation in both files is column 0. At implementation time, re-run Milestone 1; if the sources still match, use the text below verbatim. If they contradict it, update the bullets to match the sources — do not invent.

**Insert 1 — `# v0.10.3`** immediately after the frontmatter's closing `---` and its following blank line (before `# v0.10.2`):

    # v0.10.3

    *September 17, 2026*

    Miru Agent `v0.10.3` connects to the MQTT broker over TLS on port `443` instead of `8883`, so devices behind firewalls that only allow outbound `443`/HTTPS traffic can maintain their real-time connection to the control plane.

    ## Improvements

    - The agent now maintains its real-time MQTT connection over TLS on port `443` instead of `8883`, so it works on networks that only permit outbound `443`/HTTPS traffic

    ---

**Insert 2 — `# v0.7.2`** after the `v0.8.0` section's trailing `---` and its following blank line, immediately before `# v0.7.1`:

    # v0.7.2

    *September 17, 2026*

    Backport of the MQTT port-`443` change to the `v0.7` line. Miru Agent `v0.7.2` connects to the MQTT broker over TLS on port `443` instead of `8883` so devices behind firewalls that only allow outbound `443`/HTTPS traffic can maintain their real-time connection, and hardens the MQTT TLS stack.

    ## Improvements

    - The agent now maintains its real-time MQTT connection over TLS on port `443` instead of `8883`, so it works on networks that only permit outbound `443`/HTTPS traffic

    ## Fixes

    - Moved the MQTT TLS stack from rustls to native-tls, resolving the known `rustls-webpki` security advisories RUSTSEC-2026-0098, RUSTSEC-2026-0099, and RUSTSEC-2026-0104

    ---

**Edit 3 — `docs/developers/agent/security.mdx` line 116.** Replace:

    - **MQTT over TLS** — the persistent MQTT connection uses TLS on port 8883 (the standard MQTT over TLS port)

with:

    - **MQTT over TLS** — the persistent MQTT connection uses TLS on port 443. Current agents connect on port 443 so a device only needs outbound 443/HTTPS access through a firewall; older agent versions connect on port 8883, which the broker still accepts.

Keep the em dash (`—`, U+2014) that opens the bullet; use it, never `--`, if additional prose punctuation is needed. Do not add the internal proxy or DNS detail. Leave the surrounding "Transport security", "Network posture", and other bullets untouched.

If CSpell flags a word in any new text, add it to `cspell.json` `words` in sorted order and commit it alongside the file that introduced it. No other files change.


## Concrete Steps

### Milestone 1 — Re-verify sources

From any directory:

    git -C /home/ben/miru/workbench1/repos/agent log --oneline v0.10.2..v0.10.3
    git -C /home/ben/miru/workbench1/repos/agent log --oneline v0.7.1..v0.7.2
    git -C /home/ben/miru/workbench1/repos/agent log -1 --format=%ci v0.10.3
    git -C /home/ben/miru/workbench1/repos/agent log -1 --format=%ci v0.7.2

Expect the `v0.10.2..v0.10.3` delta to still contain `feat(mqtt): ... port 443 (#245)` as the only announceable item, the `v0.7.1..v0.7.2` delta to contain `#248` (443) plus the `native-tls` RUSTSEC fix, and both tags dated 2026-09-17. If a delta contradicts the Plan of Work text, update this plan's living sections and Plan of Work before editing the docs — do not invent bullets.

From `/home/ben/miru/workbench1/repos/docs`:

    git branch --show-current
    git status --short
    grep -n '^# v' docs/changelog/agent.mdx
    grep -n '8883' docs/developers/agent/security.mdx

Expect branch `docs/agent-443-changelog-and-port`; a tree clean except possibly this plan file; changelog headings starting `v0.10.2` then `v0.10.1` … with no `v0.10.3` or `v0.7.2` yet; and exactly one `8883` line (116) in `security.mdx`.

Commit step: only if a source check forced a correction to this plan, from `/home/ben/miru/workbench1/repos/docs`:

    git add plans/backlog/20260917-agent-443-changelog-and-port.md
    git commit -m "docs(plan): correct agent 443 changelog sources"

Otherwise skip the commit for this milestone.


### Milestone 2 — Insert the two changelog sections and commit

From `/home/ben/miru/workbench1/repos/docs`, apply Insert 1 and Insert 2 from the Plan of Work. Then:

    grep -n '^# v' docs/changelog/agent.mdx

Expect the order `v0.10.3`, `v0.10.2`, `v0.10.1`, `v0.10.0`, `v0.9.0`, `v0.8.1`, `v0.8.0`, `v0.7.2`, `v0.7.1`, `v0.7.0`, … each version appearing once.

    git diff --stat

Expect only `docs/changelog/agent.mdx` (insertions only).

    git add docs/changelog/agent.mdx
    git commit -m "docs(changelog): add agent v0.10.3 and v0.7.2 release notes"

One commit for this milestone.


### Milestone 3 — Update the security page port statement and commit

From `/home/ben/miru/workbench1/repos/docs`, apply Edit 3. Then:

    grep -n '443\|8883' docs/developers/agent/security.mdx
    git diff --stat

Expect line 116 to now mention both `443` (current) and `8883` (older, still accepted), and the diff to touch only `docs/developers/agent/security.mdx`.

    git add docs/developers/agent/security.mdx
    git commit -m "docs(agent): MQTT over TLS now uses port 443"

One commit for this milestone.


### Milestone 4 — Test and lint

From `/home/ben/miru/workbench1/repos/docs` (run `pnpm install --frozen-lockfile` first if `node_modules/` is missing):

    pnpm run test:lint
    ./scripts/lint.sh
    pnpm run validate

Expect: the lint smoke tests pass (exit 0, quiet on success); `./scripts/lint.sh` prints `All documentation lint checks passed.` with 0 CSpell issues; `pnpm run validate` reports build validation passed. These are the repo's tests for MDX prose — there is no changelog unit test. If CSpell flags `RUSTSEC` or a `native-tls` fragment, add the word to `cspell.json` `words` in sorted order and re-run. If any check fails, fix the changelog, the security page, or `cspell.json` and make a new commit, then re-run all three.

Commit step: if all three exit 0 with no further file changes, skip the commit. If a fix was required, from `/home/ben/miru/workbench1/repos/docs`:

    git add docs/changelog/agent.mdx docs/developers/agent/security.mdx cspell.json
    git commit -m "docs: fix lint findings for agent 443 docs"

Stage only the files that actually changed. One commit for this milestone, and only if a fix was needed.


### Milestone 5 — Publish and preflight

From `/home/ben/miru/workbench1/repos/docs`:

    git fetch origin main
    git rebase origin/main
    git push -u origin HEAD

If the rebase rewrote already-pushed commits, `git push --force-with-lease` is the expected follow-up (never a plain force reset). Then, if no PR exists yet:

    gh pr create --draft --base main --title "docs: agent v0.10.3 / v0.7.2 MQTT port 443 changelog and security update" --body "$(cat <<'EOF'
    ## Summary
    - Add `# v0.10.3` and `# v0.7.2` sections to `docs/changelog/agent.mdx` for the MQTT port `443` change (v0.7.2 is a backport plus a native-tls security fix).
    - Update the MQTT over TLS bullet in `docs/developers/agent/security.mdx` to state current agents use port `443`; `8883` remains valid for older agents.

    ## Test plan
    - [ ] `pnpm run test:lint`, `./scripts/lint.sh`, and `pnpm run validate` exit 0
    - [ ] `/changelog/agent` opens with `# v0.10.3`; `# v0.7.2` sits between `v0.8.0` and `v0.7.1`; existing entries unchanged
    - [ ] `/developers/agent/security` describes port 443 for current agents and 8883 for older ones
    - [ ] Preflight reports CLEAN (CI green on this head)

    EOF
    )"

Watch CI on the current head:

    git rev-parse HEAD
    gh pr checks --watch

If any check is red, fetch the failed logs (`gh run view <id> --log-failed`), fix, make a new commit, push once, and watch again. Repeat until every check is green for that head SHA.

Preflight must report `CLEAN` (CI green on the pushed branch head) before the PR leaves draft and before this task is reported complete. A green local `./scripts/preflight.sh` is useful but not sufficient; a green run on an older commit is not sufficient.

Commit step: if CI is green on the first pushed head, skip. If a CI fix was required, commit only the fix (one commit), push, and re-watch.


## Validation and Acceptance

1. `docs/changelog/agent.mdx` begins its body (after frontmatter) with `# v0.10.3`, `*September 17, 2026*`, a one-paragraph summary, `## Improvements` with a single port-`443` bullet, then `---`, then the untouched `# v0.10.2` section.
2. A `# v0.7.2` section (`*September 17, 2026*`, a backport summary, `## Improvements` then `## Fixes`, then `---`) sits immediately between the `v0.8.0` section and `# v0.7.1`. The Fixes bullet names native-tls and the three `RUSTSEC-2026-00xx` advisories.
3. `grep -n '^# v' docs/changelog/agent.mdx` lists `v0.10.3, v0.10.2, …, v0.8.0, v0.7.2, v0.7.1, v0.7.0, …`, each version once.
4. `docs/developers/agent/security.mdx` line 116 states current agents use port `443` (outbound `443`/HTTPS through a firewall) and that `8883` remains valid for older agents; it does not mention the internal proxy or DNS, and does not imply `8883` deployments must switch to `443` immediately.
5. `git diff main --stat` shows only `docs/changelog/agent.mdx` and `docs/developers/agent/security.mdx` (plus a words-only `cspell.json` edit if lint required one). No change to `docs/docs.json`, `docs/developers/agent/versions.mdx`, `docs/developers/agent/overview.mdx`, `docs/developers/agent/architecture.mdx`, the Dependabot config, or the audit ignore list.
6. From `/home/ben/miru/workbench1/repos/docs`, `pnpm run test:lint`, `./scripts/lint.sh`, and `pnpm run validate` all exit 0.
7. **Preflight reports `CLEAN`.** `gh pr checks` shows every GitHub Actions check passing for the current pushed head SHA. Criterion 6 (local) and this CI-green head must both hold before the PR leaves draft and before the task is reported complete.


## Idempotence and Recovery

Both changelog inserts and the security-page edit are text changes. Guards make them repeatable: `grep -c '^# v0.10.3' docs/changelog/agent.mdx` and `grep -c '^# v0.7.2' docs/changelog/agent.mdx` return `0` before and `1` after each insert — if either already returns `1` with matching content, skip that insert rather than duplicating it. For the security page, `grep -c 'port 443' docs/developers/agent/security.mdx` returns `0` before and `1` after. `git revert` of the Milestone 2 or 3 commit restores the prior file. Lint commands are read-only apart from building `tools/lint/lint` (gitignored). If CI fails after a push, fix the cause in a new commit — do not amend a pushed commit. The only permitted force-push is `--force-with-lease` after rebasing onto `origin/main`. The agent repo is never written to.

# Mark agent versions below v0.7.0 as end of life

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (mirurobotics/docs) | read-write | Swap the Status badge from Deprecated to End of life for agent `v0.6.x` / `v0.5.x` in `docs/developers/agent/versions.mdx` and for the inherited Device API `v0.1.0` row in `docs/developers/device-api/versions.mdx`, updating each file's import line to match. No other files change. |

This plan lives in `docs/plans/` because the only edits are in the docs repo. Work happens on the already-checked-out branch `claude/agent-versions-eol-jx4xr4` (base: `main`). Do not create new branches.

## Purpose / Big Picture

Agent versions `v0.6.x` and `v0.5.x` entered deprecation with a published EOL date of **2026-06-16**. That date has passed (today is 2026-07-27), so per the support policy on the same page they are now **end of life**, but the support table at `/developers/agent/versions` still shows them as Deprecated. After this change, every agent version below `v0.7.0` (`v0.6.x`, `v0.5.x`, `v0.4.x`) shows the red End of life badge; `v0.7.x` and above stay Supported.

The Device API versions page states its support is inherited from the agent: `v0.1.0` is used only by agents `v0.5.x` and `v0.6.x` (per that page's compatibility matrix), so its Deprecated badge — carrying the same 2026-06-16 EOL date — flips to End of life in the same pass to keep the two pages consistent.

No EOL dates are invented or changed; both affected tables already carry 2026-06-16 in the EOL column. Only Status badges (and the corresponding component imports) change.

## Progress

- [x] Edit `docs/developers/agent/versions.mdx`: `v0.6.x` and `v0.5.x` rows `<DeprecatedBadge />` → `<EndOfLifeBadge />`; drop `DeprecatedBadge` from the import. (2026-07-27 20:19 UTC)
- [x] Edit `docs/developers/device-api/versions.mdx`: `v0.1.0` row `<DeprecatedBadge />` → `<EndOfLifeBadge />`; swap `DeprecatedBadge` for `EndOfLifeBadge` in the import. (2026-07-27 20:19 UTC)
- [x] Run `pnpm run test:lint` and `pnpm run lint` (→ `./scripts/lint.sh`) from the repo root; both pass. (2026-07-27 20:19 UTC — both exit 0; lint ends "All documentation lint checks passed.")
- [x] Commit: `docs: mark agent versions below v0.7.0 as end of life`. (2026-07-27 20:20 UTC — content edits and this plan update in one commit)
- [ ] Push via the preflight workflow; preflight reports `CLEAN` (CI green on the pushed branch head).

Use timestamps when you complete steps.

## Surprises & Discoveries

- The full lint suite ran locally without issue, including `mint openapi-check` (anticipated as potentially network-restricted). All six OpenAPI specs validated. `pnpm run test:lint` is silent on success — exit code 0 is the pass signal.
- `scripts/lint.sh` scopes all checks to the `docs/` content root, so `plans/` files are never linted; plan-file-only commits cannot fail the CI lint job.
- (2026-07-27 21:36 UTC) The first CI run on the branch head (a243e18) failed its `audit` job: `pnpm audit` flagged 3 newly published advisories in transitive deps unrelated to the docs diff (postcss GHSA-r28c-9q8g-f849, brace-expansion GHSA-mh99-v99m-4gvg, tar GHSA-r292-9mhp-454m). Main's last run on the same base commit was green on Jul 24 — the advisories landed in between. Fixed by bumping the corresponding `pnpm.overrides` in `package.json` (see Decision Log).

## Decision Log

- Decision: the Device API `v0.1.0` row in `docs/developers/device-api/versions.mdx` is in scope and flips to End of life.
  Rationale: that page's own policy says "The Device API does not have its own support policy—its support is inherited from the Miru Agent. If a supported Miru Agent uses an API version, that API version is supported." Its compatibility matrix shows `v0.1.0` is used only by agents `v0.5.x`/`v0.6.x`, which this task marks EOL, and the row already carries the same 2026-06-16 EOL date. Leaving it Deprecated would contradict the agent table. The edit is separable: dropping it leaves the agent-table edit intact.
  Date/Author: 2026-07-27 / plan author.
- Decision: remove `DeprecatedBadge` from both files' import lines after the swap.
  Rationale: the custom linter's `import-used` rule (`tools/lint/linter/importused/importused.go`) fails on any imported name not used in the document body. After the badge swaps neither file references `DeprecatedBadge`. The `import-sorted` rule orders by import *path* only, so renaming names within a line cannot break it.
  Date/Author: 2026-07-27 / plan author.
- Decision: no EOL date cells change and no rows are added or removed.
  Rationale: the task is to reflect an already-published EOL date that has passed, not to set policy. `v0.4.x` already shows End of life and stays byte-identical. Precedent: the agent table keeps EOL versions listed rather than deleting rows.
  Date/Author: 2026-07-27 / plan author.
- Decision: no prose changes on either page, and no changes to `docs/changelog/agent.mdx`, `docs/developers/platform-api/versioning.mdx`, or `docs/developers/platform-api/sdks.mdx`.
  Rationale: the support-policy prose is generic (it defines the three levels, names no versions). The agent changelog is release notes with no support-status statements. The platform API pages describe platform API/SDK versioning, not agent support status (`2025-10-21.zion` is already End of life there; the SDK table's `v0.1.x - v0.6.x` range is Python SDK versions, not agent versions). Mentions of `v0.6.0` in `docs/primitives/devices.mdx` and `docs/cfg-mgmt/provision-devices/provisioning-script.mdx` are example values, not status statements.
  Date/Author: 2026-07-27 / plan author.
- Decision: expand scope to bump three `pnpm.overrides` security ranges in `package.json` (postcss `>=8.5.10`→`>=8.5.18`, brace-expansion `>=5.0.7`→`>=5.0.8`, tar `>=7.5.19`→`>=7.5.21`) plus the resulting `pnpm-lock.yaml` update.
  Rationale: CI's `audit` job fails the branch on any actionable `pnpm audit` advisory, and three new advisories were published after main's last green run. Not part of the docs change, but required for CI green; the `pnpm.overrides` block exists for exactly this (precedent: commit d7663aa "chore(deps): clear pnpm audit advisories (#135)"). Committed separately as `chore(deps)` to keep the docs change atomic. `./scripts/audit.sh` and `pnpm run lint` verified locally after the bump.
  Date/Author: 2026-07-27 / implement agent.

(Fill in on completion: commit SHA, lint/preflight results, deferred follow-ups.)

## Context and Orientation

The docs repo is the Mintlify documentation site; content lives under `docs/` (the content root). Support badges are shared components in `docs/snippets/components/support.jsx`, which exports exactly three names: `SupportedBadge` (green "Supported"), `DeprecatedBadge` (orange "Deprecated"), and `EndOfLifeBadge` (red "End of life"). Every place in the corpus that states agent-version support status was found by grepping for these badge names and for `v0.5`/`v0.6`; the full result set is the two tables edited here plus the out-of-scope pages listed in the Decision Log.

Current state on `main`:

- `docs/developers/agent/versions.mdx` — line 8 imports all three badges. The support table (lines 18–25) reads:

      | `v0.9.x`    | 2026-05-12   | —            | <SupportedBadge />  |
      | `v0.8.x`    | 2026-04-12   | —            | <SupportedBadge />  |
      | `v0.7.x`    | 2026-03-13   | —            | <SupportedBadge />  |
      | `v0.6.x`    | 2026-01-22   | 2026-06-16   | <DeprecatedBadge />  |
      | `v0.5.x`    | 2025-09-25   | 2026-06-16   | <DeprecatedBadge /> |
      | `v0.4.x`    | 2025-08-18   | 2026-01-01   | <EndOfLifeBadge />  |

- `docs/developers/device-api/versions.mdx` — line 7 imports `SupportedBadge, DeprecatedBadge` (no `EndOfLifeBadge`). The support table (lines 44–48) ends with:

      | <LinkNewTab href="/references/device-api/v0.1.0">v0.1.0</LinkNewTab>    | 2025-09-21   | 2026-06-16   | <DeprecatedBadge /> |

Lint toolchain (all from the repo root; `package.json` defines the scripts):

- `pnpm run lint` → `./scripts/lint.sh`: builds and runs the Go MDX prose linter in `tools/lint` (rules include `import-used`, `import-sorted`, `import-resolves`, heading-case, no-double-dash, component style), then ESLint (MDX validity), CSpell (`cspell.json`), then `mint openapi-check`. Requires pnpm, Go, and network access.
- `pnpm run test:lint` → `./tests/test-lint.sh`: smoke tests of the linter itself.
- CI (`.github/workflows/ci.yml`) runs on the PR: `lint` (runs `pnpm run test:lint` then `./scripts/lint.sh`), `audit`, `shell-tests`, and `changes`; the custom-linter jobs are skipped because `tools/lint/**` is untouched.

## Plan of Work

Two file edits, badge swaps only. Re-locate lines by content, not number, in case the files have shifted.

**Edit 1 — `docs/developers/agent/versions.mdx`:**

1. Import line:

       import { SupportedBadge, DeprecatedBadge, EndOfLifeBadge } from '/snippets/components/support.jsx';

   becomes

       import { SupportedBadge, EndOfLifeBadge } from '/snippets/components/support.jsx';

2. In the support table, the `v0.6.x` and `v0.5.x` rows change `<DeprecatedBadge />` to `<EndOfLifeBadge />`, matching the `v0.4.x` row's cell formatting (`<EndOfLifeBadge />  |`):

       | `v0.6.x`    | 2026-01-22   | 2026-06-16   | <EndOfLifeBadge />  |
       | `v0.5.x`    | 2025-09-25   | 2026-06-16   | <EndOfLifeBadge />  |

   Released and EOL cells are unchanged. The `v0.9.x`–`v0.7.x` and `v0.4.x` rows are untouched.

**Edit 2 — `docs/developers/device-api/versions.mdx`:**

1. Import line:

       import { SupportedBadge, DeprecatedBadge } from '/snippets/components/support.jsx';

   becomes

       import { SupportedBadge, EndOfLifeBadge } from '/snippets/components/support.jsx';

2. In the support table, the `v0.1.0` row changes `<DeprecatedBadge />` to `<EndOfLifeBadge />`:

       | <LinkNewTab href="/references/device-api/v0.1.0">v0.1.0</LinkNewTab>    | 2025-09-21   | 2026-06-16   | <EndOfLifeBadge /> |

   The `v0.2.1` and `v0.2.0` rows, the agent compatibility matrix, and all prose are untouched.

Do NOT touch any other file (this plan file itself excepted).

## Concrete Steps

All commands run from the docs repo root, `/home/user/docs`.

1. Confirm the branch and a clean tree:

       git branch --show-current   # expect: claude/agent-versions-eol-jx4xr4
       git status --short          # expect: only this plan file, if uncommitted

2. Idempotence guard — confirm the swaps are not already done:

       grep -c 'DeprecatedBadge />' docs/developers/agent/versions.mdx docs/developers/device-api/versions.mdx

   Expect `2` and `1`. If both are `0`, the edits are already in place — skip to step 5.

3. Apply Edit 1 and Edit 2 from Plan of Work.

4. Verify the swaps landed and no `DeprecatedBadge` references remain in either file:

       grep -n 'DeprecatedBadge' docs/developers/agent/versions.mdx docs/developers/device-api/versions.mdx
       # expect: no matches (grep exits 1)
       grep -c 'EndOfLifeBadge />' docs/developers/agent/versions.mdx docs/developers/device-api/versions.mdx
       # expect: 3 (v0.6.x, v0.5.x, v0.4.x) and 1 (v0.1.0)

5. Lint:

       pnpm install --frozen-lockfile
       pnpm run test:lint          # expect: all smoke tests pass
       pnpm run lint               # expect final line: "All documentation lint checks passed."

   If the `import-used` rule flags anything, an import line was not updated to match its file's usage — fix the import, do not disable rules.

6. Confirm the diff scope:

       git diff main --stat
       # expect: docs/developers/agent/versions.mdx (3 lines), docs/developers/device-api/versions.mdx (2 lines), plus this plan file

7. Commit (one milestone, one commit):

       git add docs/developers/agent/versions.mdx docs/developers/device-api/versions.mdx plans/
       git commit -m "docs: mark agent versions below v0.7.0 as end of life"

8. Publish and watch CI via the preflight workflow (push the branch head, let the `lint`, `audit`, and `shell-tests` jobs run, and fix any failures from the job logs).

## Validation and Acceptance

- On `/developers/agent/versions`, every version below `v0.7.0` (`v0.6.x`, `v0.5.x`, `v0.4.x`) shows `<EndOfLifeBadge />`; `v0.7.x`, `v0.8.x`, `v0.9.x` still show `<SupportedBadge />`. Verified by the step-4 greps.
- On `/developers/device-api/versions`, `v0.1.0` shows `<EndOfLifeBadge />`; `v0.2.0` and `v0.2.1` still show `<SupportedBadge />`.
- No `DeprecatedBadge` reference (import or usage) remains in either edited file; `docs/snippets/components/support.jsx` is unchanged (other pages may use the badge later).
- All EOL and Released date cells are byte-identical to `main` (`git diff main` shows only badge-name and import-name changes).
- `pnpm run test:lint` and `pnpm run lint` both exit 0 from `/home/user/docs`.
- Optional render check: `pnpm run dev` and open `/developers/agent/versions` — three red End of life badges render below three green Supported badges.
- CI on the pushed branch head is green: the `lint`, `audit`, `shell-tests`, and `changes` jobs of `.github/workflows/ci.yml` all pass (custom-linter jobs are skipped since `tools/lint/**` is untouched).
- **Gate: preflight must report `CLEAN` — CI green on the pushed branch head — before the PR leaves draft or the task is reported complete.**

No new tests are added: the change swaps component references inside prose tables, fully covered by the existing `import-used` lint rule and CI lint job; a bespoke test asserting badge placement would duplicate content the linter already guards.

## Idempotence and Recovery

- Step 2's grep guard makes the edits re-runnable; a second pass over already-edited files is a no-op (the source substring `DeprecatedBadge` is gone).
- If an edit goes wrong before committing: `git checkout main -- docs/developers/agent/versions.mdx docs/developers/device-api/versions.mdx` and redo step 3.
- If a bad version was committed: `git revert <sha>` (do not force-push shared branches). The change is two text files — no migrations, no assets, revert is safe.
- If lint or CI fails: read the failure from the job logs, fix the underlying cause, and re-run. Do not bypass with `--no-verify` or by disabling lint rules.

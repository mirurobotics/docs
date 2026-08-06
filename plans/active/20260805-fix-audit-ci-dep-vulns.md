# Fix failing audit CI job by raising fast-uri and brace-expansion override floors

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (this repo, `mirurobotics/docs`) | read-write | Edit `package.json` pnpm overrides and regenerate `pnpm-lock.yaml` |

This plan lives in `docs/plans/` because all changes are confined to this repository. Work happens on the existing branch `fix/audit-dep-vulns` (base branch: `main`). Both vulnerabilities are pre-existing lockfile issues on `main`; they surfaced on PR #144 but were not caused by it.

## Purpose / Big Picture

The `audit` job in this repo's CI (`.github/workflows/ci.yml`, step "Run security audit", which runs `./scripts/audit.sh`) is failing because `pnpm audit` reports two actionable high-severity advisories against transitive dependencies pinned in `pnpm-lock.yaml`:

1. **fast-uri** — vulnerable `>=3.0.0 <3.1.5`, patched `>=3.1.5`, advisory GHSA-7p8r-x3mc-p8w7. Path: `. > mint > @mintlify/cli > @mintlify/common > @asyncapi/parser > ajv > fast-uri`. Lockfile currently resolves `fast-uri@3.1.4`.
2. **brace-expansion** — vulnerable `>=4.0.0 <5.0.9`, patched `>=5.0.9`, advisory GHSA-rgw5-rvv9-x895. Path: `. > eslint > minimatch > brace-expansion`. Lockfile currently resolves `brace-expansion@5.0.8`.

After this change, `./scripts/audit.sh` exits 0 with no actionable advisories, and the `audit` CI job is green. The fix is a targeted version bump via existing pnpm overrides — no advisories are added to the ignore list.

## Progress

- [x] Milestone 1: Bump override floors, regenerate lockfile, verify audit passes locally, commit. (Superseded mid-flight: the identical fix landed on `main` via PR #145 before this branch was pushed — see Surprises & Discoveries.)
- [x] Milestone 2: Run the repo's local checks (frozen-lockfile install, lint smoke tests, doc lint), then validate CI is green via preflight.

## Surprises & Discoveries

- **The fix landed on `main` independently while this branch was in flight.** PR #145 (commit `18861fe`, "chore(deps): raise audit override floors and fix lint.sh on macOS bash 3.2") merged the byte-identical `package.json` override bumps (`fast-uri >=3.1.5 <4`, `brace-expansion >=5.0.9`) and the matching `pnpm-lock.yaml` re-resolution. This branch's local fix commit was confirmed byte-identical (`git diff origin/main HEAD -- package.json pnpm-lock.yaml` was empty) and was dropped automatically on rebase ("patch contents already upstream").
- PR #145 also fixed a latent macOS bash 3.2 parse failure in `scripts/lint.sh` that was out of scope here.

## Decision Log

- **Rebase and drop, rather than re-apply.** After `origin/main` picked up PR #145, the branch was rebased onto `main` and the now-redundant fix commit dropped. Re-creating an empty or duplicate dependency commit would add noise with no content. The branch now carries only this plan document.
- The `ignoreCves` list was not touched, per the plan's constraint.

## Outcomes & Retrospective

The end goal — `audit` CI job green with no actionable advisories — is achieved on `main` (via PR #145) and re-verified on this branch after rebase:

- `pnpm install --frozen-lockfile` — exit 0.
- `./scripts/audit.sh` — exit 0; output `1 vulnerabilities found / Severity: 1 high (1 ignored)`, matching the expected transcript. No fast-uri or brace-expansion tables.
- `pnpm-lock.yaml` resolves `fast-uri@3.1.5` and `brace-expansion@5.0.9`; no `fast-uri@3.1.4` / `brace-expansion@5.0.8` pins remain.
- `pnpm run test:lint` — exit 0.
- `./scripts/lint.sh` — exit 0, "All documentation lint checks passed."

Retrospective: the diagnosis and fix in this plan were correct (identical to what shipped in #145), but the change raced with an independent fix. Lesson: for small, urgent CI-red fixes, check `origin/main` for concurrent fixes immediately before pushing.

## Context and Orientation

This repo is the Miru documentation site (Mintlify-based). Relevant pieces:

- **Package manager**: pnpm, pinned to `pnpm@10.17.0` via the `packageManager` field in `package.json`. CI activates it with `corepack enable`. Locally, plain `pnpm` works if installed; otherwise run `corepack enable` first so the pinned version is used.
- **Other local tooling**: the lint checks in Milestone 2 (`pnpm run test:lint`, `scripts/lint.sh`) require Go on PATH — `scripts/lint.sh` exits 1 with "go is required for MDX prose linting" without it. Go is only needed for those verification steps, not for the dependency fix itself.
- **Audit script**: `scripts/audit.sh` runs `pnpm audit --ignore-registry-errors` from the repo root. `pnpm audit` checks `pnpm-lock.yaml` against the npm advisory database and exits non-zero when actionable (non-ignored) vulnerabilities exist.
- **CI**: `.github/workflows/ci.yml` job `audit` does checkout → `corepack enable` → setup-node 22 with pnpm cache → `pnpm install --frozen-lockfile` → `./scripts/audit.sh`. `--frozen-lockfile` fails if `pnpm-lock.yaml` is out of sync with `package.json`, so the lockfile must be regenerated and committed together with the `package.json` edit.
- **Ignore mechanism**: `package.json` → `pnpm.auditConfig.ignoreCves` lists 8 CVEs. This is why the audit output reads "Severity: 3 high (1 ignored)" — one high advisory is already suppressed there and prints no table. We do not touch this list.
- **Override mechanism**: `package.json` → `pnpm.overrides` force-resolves transitive dependencies. Both vulnerable packages already have overrides, but each floor sits one patch release below the fixed version, so the lockfile resolves the newest still-vulnerable release of each. Because the packages are transitive, bumping direct deps would not reliably fix this; raising the existing override floors is the targeted fix and guarantees the lockfile re-resolves to patched versions. Overrides are mirrored into the `overrides:` block at the top of `pnpm-lock.yaml` and the resolved versions appear in its `snapshots`/`packages` sections.

## Plan of Work

1. In `package.json`, under `pnpm.overrides`, change:
   - `"fast-uri": ">=3.1.4 <4"` → `"fast-uri": ">=3.1.5 <4"`
   - `"brace-expansion": ">=5.0.8"` → `"brace-expansion": ">=5.0.9"`
2. From the repo root, run `pnpm install` to re-resolve the two overridden packages and rewrite `pnpm-lock.yaml`. Confirm the lockfile diff is limited to the override ranges and the two package versions (each should now be `>=3.1.5` / `>=5.0.9` respectively); no other dependency should move.
3. Verify locally with the exact command CI runs (`./scripts/audit.sh`), then run the repo's other pnpm-dependent checks to confirm the bumped transitive deps break nothing.
4. Validate CI on the pushed branch head via preflight.

If, unexpectedly, `pnpm install` cannot satisfy an override (e.g. a patched version does not exist on the registry), stop and re-check the advisory's patched range on the "More info" URL in the audit output — do not fall back to adding the GHSAs to `ignoreCves` without explicit approval.

## Concrete Steps

All commands run from the repo root `/home/ben/miru/workbench6/repos/docs` on branch `fix/audit-dep-vulns`.

Milestone 1 — bump and verify audit:

    # 0. Ensure pnpm 10.17.0 is available
    corepack enable   # skip if `pnpm --version` already prints 10.17.0

    # 1. Edit package.json as described in Plan of Work (two override lines)

    # 2. Regenerate lockfile
    pnpm install

    # 3. Inspect lockfile changes — expect only fast-uri and brace-expansion lines
    git diff --stat package.json pnpm-lock.yaml
    grep -n "fast-uri@\|brace-expansion@" pnpm-lock.yaml
    # expected: versions >=3.1.5 and >=5.0.9; no remaining fast-uri@3.1.4 / brace-expansion@5.0.8

    # 4. Run the audit exactly as CI does
    ./scripts/audit.sh; echo "exit=$?"

Expected transcript for step 4 (before the change it prints two high-severity tables, "3 vulnerabilities found / Severity: 3 high (1 ignored)", and a non-zero exit):

    == Security Audit ==
    1 vulnerabilities found
    Severity: 1 high (1 ignored)
    exit=0

The exact count line may vary if the advisory database has changed; acceptance is exit code 0 and no fast-uri or brace-expansion tables in the output.

    # 5. Commit (repo requires SSH-signed commits; signing is already configured — do not disable it)
    git add package.json pnpm-lock.yaml
    git commit -m "fix(deps): raise fast-uri and brace-expansion override floors to patched versions"

Milestone 2 — local checks mirroring CI, then preflight:

    # 1. Verify lockfile is in sync the way CI installs it
    pnpm install --frozen-lockfile   # must succeed with no lockfile changes

    # 2. Run the lint job's checks (they exercise eslint/minimatch/brace-expansion and mint)
    pnpm run test:lint               # lint smoke tests; expect all tests to pass
    ./scripts/lint.sh                # documentation lint; expect exit 0

    # 3. Commit any fixes this surfaced (only if changes were needed)
    git add -A && git commit -m "fix(deps): address lint fallout from dependency bump"

    # 4. Push and validate CI (see Validation and Acceptance for the gate)
    git push -u origin fix/audit-dep-vulns
    # then run the $preflight workflow; if unavailable, watch CI directly:
    gh pr checks fix/audit-dep-vulns --watch   # or: gh run watch

## Validation and Acceptance

Acceptance is all of the following, in order:

1. `./scripts/audit.sh` (from `/home/ben/miru/workbench6/repos/docs`) exits 0. Output contains no `fast-uri` and no `brace-expansion` advisory tables; the only remaining severity accounting is the pre-existing ignored advisory (e.g. "Severity: 1 high (1 ignored)").
2. `! grep -q "fast-uri@3.1.4\|brace-expansion@5.0.8" pnpm-lock.yaml` succeeds (vulnerable pins are gone; note plain `grep` exits 1 on zero matches, which is the passing case here).
3. `pnpm install --frozen-lockfile` succeeds, proving `package.json` and `pnpm-lock.yaml` are consistent (this is what the CI `audit` and `lint` jobs run before their checks).
4. Test steps: `pnpm run test:lint` passes (lint smoke test suite) and `./scripts/lint.sh` exits 0, proving the bumped transitive deps did not break the eslint/cspell/mint toolchain. This repo has no other JS test suite; these are the CI `lint` job's exact commands.
5. **Preflight gate**: push the branch head and run the preflight workflow (`$preflight`). Preflight must report CLEAN — meaning the full CI run on the pushed head of `fix/audit-dep-vulns` is green, including the `audit`, `lint`, and `shell-tests` jobs — before the PR leaves draft or the task is reported complete. If any CI job fails, fix from the CI logs, commit, push, and re-run preflight until CLEAN.

## Idempotence and Recovery

Every step is safe to repeat: the `package.json` edit is a fixed end-state, `pnpm install` is idempotent given the same inputs, and the audit/lint scripts are read-only. If `pnpm install` produces an unexpectedly large lockfile diff or breaks lint, roll back with:

    git checkout -- package.json pnpm-lock.yaml
    pnpm install

and retry. If a patched version cannot be resolved from the registry, see the fallback note at the end of Plan of Work.

# Remove extract-zip from the lockfile to close Dependabot alerts #69 and #71

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.


## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/user/docs`) | read-write | Add one `pnpm.overrides` entry and drop two CVE ids from `pnpm.auditConfig.ignoreCves` in `package.json`; regenerate `pnpm-lock.yaml` with pnpm. |

This plan lives in `docs/plans/` because the only writes are to this repo's `package.json` and `pnpm-lock.yaml`.

Working branch: `claude/clever-euler-bc9ork` (already checked out from `main`). Do not create another branch. Base branch for the PR: `main`.


## Purpose / Big Picture

GitHub Dependabot has two open High alerts on `mirurobotics/docs`, both against `extract-zip@2.0.1` in `pnpm-lock.yaml` (development scope):

- Alert #69: GHSA-jmr9-qjv8-65gv / CVE-2026-56876, "unvalidated symlink path traversal".
- Alert #71: GHSA-7pqw-9j4j-h8q3 / CVE-2026-19693, "arbitrary file writes through symlink archive entries".

`extract-zip` has no patched release (latest is 2.0.1 from 2020; both advisories mark `<=2.0.1` vulnerable and nothing as patched), so the only fix is to remove it from the dependency graph. After this change `grep -c extract-zip pnpm-lock.yaml` prints `0`, `./scripts/audit.sh` passes without ignoring these two CVEs, CI is green, and once the PR merges to `main` Dependabot closes both alerts automatically because the vulnerable package is no longer in the default branch's lockfile.


## Progress

- [x] Milestone 1: Add the `puppeteer` override, drop the two CVE ignores, regenerate `pnpm-lock.yaml`, confirm `extract-zip` is gone; commit.
- [x] Milestone 2: Run the local checks (`./scripts/audit.sh`, `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`); commit only if a fix was needed.
- [ ] Milestone 3: Push, open the draft PR, drive CI to CLEAN, mark ready for review.


## Surprises & Discoveries

- 2026-09-20: `pnpm run validate` cannot run in the implementation sandbox. `mint validate` gates on the `is-online` package (raw DNS/socket probe that bypasses the egress proxy) and, with no cached preview client under `~/.mintlify`, exits with `running mint validate-build after updating requires an internet connection.` The `mint` version (4.2.891) is identical on `main` and this branch, so the failure is environmental, not caused by the override. The other three local checks passed; `pnpm run validate` is verified by the CI `lint` job instead.


## Decision Log

- Decision: Override `puppeteer` to `>=25.0.0` rather than `@puppeteer/browsers` to `>=3.0.0`. / Rationale: `@puppeteer/browsers` dropped `extract-zip` in 3.0.0 (2.x through 2.13.2 still depend on `extract-zip ^2.0.1`; 3.x uses `modern-tar`). Both overrides were resolved in a scratch copy with pnpm 10.17.0 and both bring the count to 0, but the `puppeteer` override moves `puppeteer`, `puppeteer-core`, and `@puppeteer/browsers` together to the versions puppeteer itself ships (25.11.0 / 25.11.0 / 3.2.2), whereas the `@puppeteer/browsers` override leaves `puppeteer-core@24.3.1` paired with a `@puppeteer/browsers` major it was never released against. Bumping `mint` does not help: the latest `mint` (4.2.909) still pins `@mintlify/scraping` to `puppeteer 24.3.1`. / Date/Author: 2026-09-20, author.


## Outcomes & Retrospective

(Summarize at completion.)


## Context and Orientation

This repo is the Mintlify documentation site. Its only npm dependencies are dev tools declared in `/home/user/docs/package.json` (`cspell`, `eslint`, `eslint-plugin-mdx`, `mint`). The package manager is pnpm 10.17.0 (`packageManager` field; run `corepack enable` if `pnpm --version` does not print `10.17.0`). Node 22 is used in CI and locally.

**How extract-zip gets in.** There is exactly one path: `mint@4.2.891 > @mintlify/cli > @mintlify/link-rot > @mintlify/scraping > puppeteer@24.3.1 > @puppeteer/browsers@2.7.1 > extract-zip@2.0.1`. Puppeteer is only exercised by Mintlify's scraping / broken-links commands, which neither CI nor the repo scripts run, so swapping its version has no effect on `mint validate`, `mint dev`, or `mint openapi-check`. pnpm 10 does not run puppeteer's post-install script (it is in the "Ignored build scripts" list by default), so no Chrome download happens either way.

**pnpm overrides.** `package.json` has a `pnpm.overrides` object (currently ~30 entries, `axios` through `sharp`) that forces transitive packages to safe ranges. This is the repo's established mechanism for exactly this kind of fix (see commit `ab72a0a`, `fix(deps): bump adm-zip to 0.6.1`). The overrides are mirrored in the `overrides:` block at the top of `pnpm-lock.yaml`, which pnpm regenerates; never hand-edit the lockfile.

**Audit ignore list.** `package.json` also has `pnpm.auditConfig.ignoreCves`. It currently contains `CVE-2026-19693` and `CVE-2026-56876` (added in commit `6ef027f` because no patch existed), which is why `./scripts/audit.sh` (`pnpm audit --ignore-registry-errors`) passes today while the Dependabot alerts stay open. Removing the two ids makes the audit job guard against `extract-zip` ever coming back. The other nine ids in the list are unrelated and stay.

**CI.** `.github/workflows/ci.yml` runs on every PR: `lint` (`pnpm install --frozen-lockfile`, `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`), `audit` (`pnpm install --frozen-lockfile`, `./scripts/audit.sh`), `shell-tests` (bats), and two `tools/lint` Go jobs that only run when `tools/lint/**` changes. `--frozen-lockfile` fails if `pnpm-lock.yaml` disagrees with `package.json`, so the lockfile must be regenerated and committed alongside the `package.json` edit. **CLEAN** means every GitHub Actions check is green on the pushed branch head SHA.

There is no `gh` CLI in this environment. Watch CI with the GitHub MCP tools (`mcp__github__pull_request_read` with method `get_status` or `get_checks`, `mcp__github__get_job_logs` for failures) or the REST API via curl (`GET /repos/mirurobotics/docs/commits/<sha>/check-runs`).

**Verified in a scratch copy (2026-09-20, pnpm 10.17.0).** With `"puppeteer": ">=25.0.0"` added to overrides: `pnpm install --lockfile-only` succeeds; the lockfile contains `puppeteer@25.11.0`, `puppeteer-core@25.11.0`, `@puppeteer/browsers@3.2.2`, `modern-tar@0.8.5`, and zero `extract-zip` lines; `pnpm install --frozen-lockfile` then succeeds (1065 packages); `pnpm audit --ignore-registry-errors` prints `No known vulnerabilities found` with the two CVEs removed from the ignore list; and `mint validate` against the real `docs/` content prints `success build validation passed`. The pre-existing peer-dependency warnings about `react-dom 18.3.1` / `react 19.2.3` are unchanged and expected.


## Plan of Work

Two files change, both at the repo root.

**Edit 1 — `package.json`, `pnpm.overrides`.** Append one entry after the last existing key (`"sharp": ">=0.35.0"`), adding the trailing comma to the `sharp` line:

    "sharp": ">=0.35.0",
    "puppeteer": ">=25.0.0"

**Edit 2 — `package.json`, `pnpm.auditConfig.ignoreCves`.** Delete the two array elements `"CVE-2026-19693",` and `"CVE-2026-56876",`. The remaining nine ids keep their order; the array stays valid JSON (the last element `"CVE-2026-76845"` still has no trailing comma).

**Regenerate `pnpm-lock.yaml`** with pnpm (Milestone 1 commands). Expected lockfile changes: a new `puppeteer: '>=25.0.0'` line in the top `overrides:` block; `extract-zip@2.0.1` and its private helpers removed; `puppeteer`, `puppeteer-core`, and `@puppeteer/browsers` re-resolved to 25.11.0 / 25.11.0 / 3.2.2 with `modern-tar` added; `chromium-bidi` and `devtools-protocol` bumped as puppeteer 25 requires. The diff is a few hundred lines and is entirely pnpm output.

**Fallback** (only if Milestone 1 fails to resolve or Milestone 2's `pnpm run validate` breaks and the failure traces to the puppeteer 25 tree): replace the override with `"@puppeteer/browsers": ">=3.0.0"`, rerun the regeneration, and record the switch and its evidence in the Decision Log. This alternative was also resolved in scratch (`@puppeteer/browsers@3.2.2`, `puppeteer` stays 24.3.1, `extract-zip` count 0).

No other files change. Do not bump `mint`, do not touch `.github/dependabot.yml` or the CI workflow, and do not remove any other id from `ignoreCves`.


## Concrete Steps

### Milestone 1 — Override, regenerate the lockfile, commit

From `/home/user/docs`:

    git branch --show-current
    git status --short
    pnpm --version
    grep -c 'extract-zip' pnpm-lock.yaml

Expect `claude/clever-euler-bc9ork`, a clean tree (apart from this plan file if not yet committed), `10.17.0`, and `3`.

Apply Edit 1 and Edit 2 to `package.json` (an editor, or the equivalent `sed` / `node` one-liners), then confirm the file is still valid JSON and only the intended keys changed:

    node -e 'const p=require("./package.json");console.log(p.pnpm.overrides.puppeteer, p.pnpm.auditConfig.ignoreCves.length)'
    git diff --stat

Expect `>=25.0.0 9` and a diff touching only `package.json` (three changed lines: the `sharp` comma, the new `puppeteer` line, and the two removed CVE lines).

Regenerate and verify the lockfile:

    pnpm install --lockfile-only
    grep -c 'extract-zip' pnpm-lock.yaml
    grep -n "^  puppeteer: '>=25.0.0'\|^  puppeteer@\|^  puppeteer-core@\|^  '@puppeteer/browsers@\|^  modern-tar@" pnpm-lock.yaml
    pnpm install --frozen-lockfile

Expect `pnpm install --lockfile-only` to end with `Done in …` (peer-dependency warnings about `react-dom` are pre-existing and fine); `grep -c` prints `0`; the second grep lists the override line plus `puppeteer@25.11.0`, `puppeteer-core@25.11.0`, `'@puppeteer/browsers@3.2.2'`, and `modern-tar@0.8.5` (the patch versions may be newer if the registry has moved on; any `25.x` / `3.x` pairing is acceptable as long as the count is `0`); and `pnpm install --frozen-lockfile` exits 0 and populates `node_modules/` (needed by Milestone 2). If `--lockfile-only` fails to resolve, apply the Fallback from the Plan of Work.

Commit:

    git add package.json pnpm-lock.yaml
    git commit -m "fix(deps): override puppeteer to >=25 to drop extract-zip"

One commit for this milestone. In the commit body, name CVE-2026-56876 and CVE-2026-19693 and state that `extract-zip` has no patched version.


### Milestone 2 — Local checks

From `/home/user/docs` (Milestone 1 already ran `pnpm install --frozen-lockfile`):

    ./scripts/audit.sh
    pnpm run test:lint
    ./scripts/lint.sh
    pnpm run validate

Expect: `./scripts/audit.sh` prints `== Security Audit ==` then `No known vulnerabilities found` and exits 0 (it now fails, rather than ignores, if `extract-zip` reappears); `pnpm run test:lint` exits 0 quietly; `./scripts/lint.sh` ends with `All documentation lint checks passed.` (it needs Go for the MDX linter and builds `tools/lint/lint`, which is gitignored); `pnpm run validate` ends with `success build validation passed`. `./scripts/preflight.sh` chains these plus the Go and bats suites and may be run instead if `bats` is installed.

If a check fails, fix the cause (see Fallback if it is puppeteer-related), regenerate the lockfile, and commit the fix as one additional commit from `/home/user/docs`:

    git add package.json pnpm-lock.yaml
    git commit -m "fix(deps): adjust puppeteer override for local checks"

Skip the commit when all four checks pass on the first run.


### Milestone 3 — Publish and preflight

From `/home/user/docs`:

    git fetch origin main
    git rebase origin/main
    git push -u origin HEAD
    git rev-parse HEAD

If the rebase rewrote commits that were already pushed, follow up with `git push --force-with-lease` (never a plain `--force`). Then, if no PR exists, open a draft PR with `mcp__github__create_pull_request` (owner `mirurobotics`, repo `docs`, base `main`, head `claude/clever-euler-bc9ork`, `draft: true`), title `fix(deps): override puppeteer to >=25 to drop extract-zip`, and a body along these lines (indented here only for the plan; paste it unindented):

    ## Summary
    - Dependabot alerts #69 (CVE-2026-56876) and #71 (CVE-2026-19693) are both `extract-zip@2.0.1`, which has no patched release.
    - Add `pnpm.overrides.puppeteer = ">=25.0.0"` so `@mintlify/scraping`'s puppeteer resolves to 25.x, whose `@puppeteer/browsers` 3.x no longer depends on `extract-zip`; regenerate `pnpm-lock.yaml` with pnpm 10.17.0.
    - Remove the two CVE ids from `pnpm.auditConfig.ignoreCves` so `./scripts/audit.sh` catches a regression.

    ## Test plan
    - [ ] `grep -c extract-zip pnpm-lock.yaml` prints `0`
    - [ ] `./scripts/audit.sh` prints `No known vulnerabilities found`
    - [ ] `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate` exit 0
    - [ ] CI green on the pushed head (preflight CLEAN)
    - [ ] Alerts #69 and #71 auto-close after merge to `main`

End the body with the attribution lines required for this session (`🤖 Generated with [Claude Code](https://claude.com/claude-code)` and the session URL).

Watch CI for the pushed head SHA with `mcp__github__pull_request_read` (method `get_status`, then `get_checks` for per-job detail) or `curl -sS -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/repos/mirurobotics/docs/commits/<sha>/check-runs`. Expect the `lint`, `audit`, and `shell-tests` jobs to succeed; the two `tools/lint` jobs are skipped on PRs that do not touch `tools/lint/**`. If a job is red, read its log with `mcp__github__get_job_logs` (`failed_only: true`), fix the cause in a new commit (never amend a pushed commit), push once, and re-watch.

Preflight must report `CLEAN` (every check green on the pushed branch head) before the PR leaves draft and before this task is reported complete. Only then mark the PR ready for review (`mcp__github__update_pull_request` with `draft: false`).

Commit step: none unless a CI fix was needed; then one commit per fix, pushed and re-watched.


## Validation and Acceptance

1. From `/home/user/docs`, `grep -c 'extract-zip' pnpm-lock.yaml` prints `0` (it prints `3` before the change), and `grep -n "^  puppeteer: '>=25.0.0'" pnpm-lock.yaml` finds the override line.
2. `package.json` `pnpm.overrides` has the key `puppeteer` = `>=25.0.0`; `pnpm.auditConfig.ignoreCves` has nine entries and contains neither `CVE-2026-56876` nor `CVE-2026-19693`; every other override and ignore id is unchanged.
3. `pnpm install --frozen-lockfile` exits 0 (the lockfile and `package.json` agree).
4. `./scripts/audit.sh` exits 0 and prints `No known vulnerabilities found`. Before the change (with the ignores still present) it printed `2 high (2 ignored)`; the ignores are now gone, so a future `extract-zip` reappearance fails the `audit` job instead of being silently ignored.
5. `pnpm run test:lint`, `./scripts/lint.sh`, and `pnpm run validate` all exit 0; the last prints `success build validation passed`.
6. `git diff main --stat` shows only `package.json` and `pnpm-lock.yaml` (plus this plan file).
7. **Preflight reports `CLEAN`:** the `lint`, `audit`, and `shell-tests` checks are all green on the current pushed head SHA of `claude/clever-euler-bc9ork`. Criteria 1-6 (local) and this CI-green head must both hold before the PR leaves draft and before the task is reported complete.
8. After the PR merges to `main`, Dependabot alerts #69 and #71 show as closed ("fixed") on `https://github.com/mirurobotics/docs/security/dependabot`. This happens automatically on merge; it is not a pre-merge gate.


## Idempotence and Recovery

All steps are repeatable. Edits 1 and 2 are guarded by criterion 2 (`node -e '...'` printing `>=25.0.0 9`); if `package.json` already matches, skip to regeneration. `pnpm install --lockfile-only` is deterministic for a given `package.json` and can be rerun freely; if it fails midway, `git checkout -- pnpm-lock.yaml` restores the committed lockfile and the command can be retried. `pnpm install --frozen-lockfile` only creates the gitignored `node_modules/`. Rolling back the whole change is `git revert` of the Milestone 1 commit (both files travel together, so a revert keeps `--frozen-lockfile` consistent). Switching to the Fallback override is the same edit-then-regenerate sequence and needs no cleanup beyond replacing the override key. The lint scripts are read-only apart from building `tools/lint/lint` (gitignored). The only permitted force-push is `--force-with-lease` after rebasing onto `origin/main`; never amend a commit that has been pushed.

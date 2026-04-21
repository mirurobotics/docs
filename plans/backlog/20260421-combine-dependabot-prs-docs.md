# Combine three Dependabot PRs in mirurobotics/docs into one consolidated PR

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `mirurobotics/docs` (working tree at `/home/ben/miru/workbench5/repos/docs`) | read-write | Edit `package.json`, regenerate `pnpm-lock.yaml`, edit `.github/workflows/ci.yml`, `.github/workflows/codeql-analysis.yml`, `.github/workflows/promote.yml` on branch `chore/combine-dependabot-prs`; open consolidated PR. |

This plan lives in `mirurobotics/docs/plans/backlog/` because the docs repo owns all files being changed and already maintains its own `plans/` directory alongside previously-completed plans in `plans/completed/`.

## Purpose / Big Picture

After this change, `mirurobotics/docs` has a single PR on branch `chore/combine-dependabot-prs` that subsumes three open Dependabot PRs (#64, #65, #66). Merging it advances five dependency updates atomically in one diff and allows the three Dependabot PRs to be closed without merging. The operator can observe: `gh pr list --state open --author app/dependabot` shows the three originals still open until merge; after merging `chore/combine-dependabot-prs` and closing #64/#65/#66, that command lists zero Dependabot PRs. Preflight (`./scripts/preflight.sh`) reports clean on the consolidated branch.

## Progress

- [ ] (YYYY-MM-DD HH:MMZ) Milestone 1: Baseline verification.
- [ ] Milestone 2: Apply npm devDependency bumps in `package.json` and regenerate `pnpm-lock.yaml`.
- [ ] Milestone 3: Update GitHub Actions workflow YAMLs for #64 and #65 (including create-github-app-token major-jump breaking-change scan).
- [ ] Milestone 4: Full preflight validation.
- [ ] Milestone 5: Push branch and open consolidated PR.

## Surprises & Discoveries

(Add entries as work proceeds.)

- Observation: …
  Evidence: …

## Decision Log

- Decision: Edit `package.json` + run `pnpm install` rather than cherry-picking Dependabot branches.
  Rationale: PR #66 rewrites 412 lines of `pnpm-lock.yaml`; cherry-picking or merging three separate Dependabot branches would create compounding lockfile conflicts. One `pnpm install` after editing `package.json` produces one clean lockfile reflecting the final state.
  Date/Author: 2026-04-21 / ben@miruml.com

- Decision: Hand-edit the three workflow YAML files rather than cherry-picking PR #64 and #65 commits.
  Rationale: The repo pins GitHub Actions to commit SHAs with a trailing `# vX.Y.Z` comment. Each bump is a single-line-ish edit; cherry-picking provides no additional safety over deterministic string replacement, and avoids merging Dependabot commit author identities into the consolidated branch.
  Date/Author: 2026-04-21 / ben@miruml.com

- Decision: Include `actions/create-github-app-token` 1.12.0 → 3.1.1 in the bundle despite the major-version jump.
  Rationale: The documented breaking changes in 3.0.0 are (a) action bundled on Node 24 runtime (action-internal, transparent to callers), (b) proxy behavior requires `NODE_USE_ENV_PROXY=1` when using HTTP_PROXY/HTTPS_PROXY, and (c) self-hosted runners need Actions Runner v2.327.1+. The only usage in this repo is in `.github/workflows/promote.yml` with inputs `app-id` and `private-key` on `ubuntu-latest` (GitHub-hosted) without proxy configuration. In 3.1.0 the `app-id` input is deprecated (now aliased to `client-id`) but still supported. Therefore, the major bump is effectively a no-op for this repo. If CI surfaces unexpected failure, drop just that bump per the rollback path.
  Date/Author: 2026-04-21 / ben@miruml.com

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

**Repo layout (`/home/ben/miru/workbench5/repos/docs/`):**

- `package.json` — declares `devDependencies` (cspell, eslint, eslint-plugin-mdx, mint) and pnpm `overrides`/`auditConfig` under the `pnpm` key. Pins `packageManager` to `pnpm@10.17.0`. No runtime `dependencies`.
- `pnpm-lock.yaml` — canonical lockfile. CI uses `pnpm install --frozen-lockfile`.
- `.github/workflows/ci.yml` — runs `changes`, `lint`, `audit`, `shell-tests`, `lint-custom-linter`, `test-custom-linter` jobs on PR. Uses `actions/checkout`, `actions/setup-node`, `dorny/paths-filter`, `actions/setup-go`. The `setup-node` action is pinned on two job lines (job `lint` and job `audit`).
- `.github/workflows/codeql-analysis.yml` — three `github/codeql-action/*` steps (`init`, `autobuild`, `analyze`), all sharing the same SHA pin.
- `.github/workflows/promote.yml` — uses `actions/create-github-app-token` once (the only usage of that action in the repo, verified via grep).
- `.github/dependabot.yml` — configures the `actions-minor`, `npm-development`, and `npm-production` groups.
- `scripts/preflight.sh` — the validation entry point. Runs, in order: `pnpm run test:lint` (lint smoke tests), `LINT_FIX=0 ./tools/lint/scripts/lint.sh` (Go lint for tools/lint/), `./tools/lint/scripts/covgate.sh` (Go coverage gate), `./scripts/lint.sh` (MDX/ESLint/CSpell/OpenAPI), `./scripts/audit.sh` (`pnpm audit --ignore-registry-errors`), then `bats pub/scripts/agent/check-miru-access_test.bats`. Exits non-zero on any step failure.
- `scripts/lint.sh` — orchestrates the Go MDX linter build + run, then eslint on MDX, then cspell, then `mint openapi-check`. The `mint` binary is the `mint` devDependency being bumped in PR #66.
- `scripts/audit.sh` — `pnpm audit --ignore-registry-errors`. `auditConfig.ignoreCves` in `package.json` filters known advisories.
- `plans/backlog/`, `plans/completed/` — plan lifecycle directories; `plans/active/` does not yet exist and will be created when this plan is promoted. Previous completed plans include `20260406-yaml-config-instance-docs.md`, `20260410-absolute-config-filepath-docs.md`, `20260414-codeql-push-only.md`.

**Branch state (verified 2026-04-21):**

- Current branch: `chore/combine-dependabot-prs`, clean working tree, off latest `main` (tip `41a4b65 feat(ci): skip custom linter jobs when tools/lint/ unchanged (#62)`).
- Remote `origin` is `git@github.com:mirurobotics/docs.git`.
- Base branch for PR is `main`.

**The three Dependabot PRs (verified via `gh pr view` on 2026-04-21):**

- **#64** `chore(deps): bump the actions-minor group with 2 updates` — Head ref `dependabot/github_actions/actions-minor-bb489aea0a`. Changes `.github/workflows/ci.yml` (2 `setup-node` occurrences) and `.github/workflows/codeql-analysis.yml` (3 `codeql-action/*` occurrences).
- **#65** `chore(deps): bump actions/create-github-app-token from 1.12.0 to 3.1.1` — Head ref `dependabot/github_actions/actions/create-github-app-token-3.1.1`. Changes `.github/workflows/promote.yml` (1 occurrence). Major version jump 1.x → 3.x (see Decision Log for breaking-change analysis).
- **#66** `chore(deps-dev): bump the npm-development group with 2 updates` — Head ref `dependabot/npm_and_yarn/npm-development-41e10df83a`. Changes `package.json` and `pnpm-lock.yaml`: `eslint` 10.2.0 → 10.2.1 (patch) and `mint` 4.2.509 → 4.2.521 (patch).

**Exact string replacements (SHA pins → new SHA pins):**

In `.github/workflows/ci.yml` — two occurrences of the `setup-node` line:

  OLD: `        uses: actions/setup-node@53b83947a5a98c8d113130e565377fae1a50d02f # v6.3.0`
  NEW: `        uses: actions/setup-node@48b55a011bda9f5d6aeb4c2d9c7362e8dae4041e # v6.4.0`

In `.github/workflows/codeql-analysis.yml` — three occurrences, for `init`, `autobuild`, `analyze`. Old SHA/comment is `c10b8064de6f491fea524254123dbe5e09572f13 # v4`. New is `95e58e9a2cdfd71adc6e0353d5c52f41a045d225 # v4` (the `# v4` comment is unchanged; only the SHA changes).

In `.github/workflows/promote.yml` — one occurrence:

  OLD: `        uses: actions/create-github-app-token@d72941d797fd3113feb6b93fd0dec494b13a2547 # v1`
  NEW: `        uses: actions/create-github-app-token@1b10c78c7865c340bc4f6099eb2f838309f1e8c3 # v3.1.1`

In `package.json` — two devDependency versions:

  `"eslint": "10.2.0"` → `"eslint": "10.2.1"`
  `"mint": "4.2.509"` → `"mint": "4.2.521"`

**Tools verified present:**

- `pnpm` (version pinned via `packageManager` field in `package.json` — corepack-managed at 10.17.0).
- `go` (required by `scripts/lint.sh`).
- `bats` (required by the `shell-tests` step of preflight).
- `gh` CLI (configured; confirmed by the `gh pr view/diff` calls that produced research data).

**Terms defined:**

- **Dependabot group PR**: a single PR that bundles updates for multiple packages in one install cycle. Groups are declared in `.github/dependabot.yml`.
- **SHA pin**: a GitHub Actions `uses:` reference that targets a specific commit SHA (40-hex) with a trailing `# vX.Y.Z` comment. Required by this repo's security posture. Example: `actions/setup-node@53b83947a5a98c8d113130e565377fae1a50d02f # v6.3.0`.
- **preflight**: the `scripts/preflight.sh` script. Its exit status is the gate — "clean" means every section exits 0.

## Plan of Work

Five milestones, each ending with a git commit so the PR history is reviewable and bisectable. All commands run from `/home/ben/miru/workbench5/repos/docs/` unless stated.

**Milestone 1 — Baseline verification.** Confirm the branch is clean and preflight passes on the unmodified tree. This rules out pre-existing failures being blamed on the bumps.

**Milestone 2 — npm devDependency bumps.** Edit `package.json` with two string replacements (`eslint 10.2.0 → 10.2.1`, `mint 4.2.509 → 4.2.521`). Run `pnpm install` (not `pnpm install --frozen-lockfile`) to update `pnpm-lock.yaml`. Verify diff scope is limited to those two files. Commit.

**Milestone 3 — GitHub Actions SHA bumps + breaking-change scan.** (a) Edit the two `setup-node` occurrences in `.github/workflows/ci.yml` and the three `codeql-action/*` occurrences in `.github/workflows/codeql-analysis.yml` (covers PR #64). (b) Scan repository usage of `actions/create-github-app-token` and confirm no reliance on v1-only behavior (covered in Decision Log — revisit if changes landed on main since research). (c) Edit the one `create-github-app-token` occurrence in `.github/workflows/promote.yml` (covers PR #65). Commit.

**Milestone 4 — Full validation.** Run `./scripts/preflight.sh` and require exit 0 across all sections. Record any new audit advisories in Surprises & Discoveries.

**Milestone 5 — Publish PR.** Push `chore/combine-dependabot-prs` to `origin`. Open the PR via `gh pr create` with a body that lists #64/#65/#66 as superseded with instruction to close them on merge. Record the resulting PR URL in Outcomes & Retrospective.

## Concrete Steps

All commands run from `/home/ben/miru/workbench5/repos/docs/` unless otherwise noted.

### Milestone 1 — Baseline verification

1. Confirm branch state:

        git status
        git branch --show-current

    Expected: `On branch chore/combine-dependabot-prs`, `nothing to commit, working tree clean`.

2. Run preflight on the unmodified branch:

        ./scripts/preflight.sh

    Expected: exit 0. Output shows the seven labelled sections (`=== Lint Smoke Tests ===`, `=== Go Lint (tools/lint) ===`, `=== Go Coverage (tools/lint) ===`, `=== Lint ===`, `=== Audit ===`, `=== Shell Script Tests ===`) each completing without error. The final `Lint` section ends `All documentation lint checks passed.`; `Audit` section ends `pnpm audit` reporting zero vulnerabilities (or only the CVEs in `package.json`'s `auditConfig.ignoreCves`); `bats` reports all tests ok.

3. If any section fails, STOP and resolve the baseline failure before proceeding. Do not interleave unrelated fixes with the dependency bumps.

4. No commit in this milestone (no file changes).

### Milestone 2 — npm devDependency bumps

1. Edit `package.json`. Replace the two literal strings (inside `devDependencies`):

        "eslint": "10.2.0",   →   "eslint": "10.2.1",
        "mint": "4.2.509"     →   "mint": "4.2.521"

    Use a deterministic string replacement (Edit tool or `sed -i`). Do not hand-edit `pnpm-lock.yaml`.

2. Regenerate the lockfile:

        pnpm install

    Expected: `pnpm install` completes without errors. Output ends with a `Done in Ns` banner and the ` +` / ` -` package counts visible. `pnpm-lock.yaml` is updated. No peer-dep warnings should require `--force` or `--shamefully-hoist`. If such warnings surface, record them in Surprises & Discoveries and do NOT pass force flags.

3. Verify diff scope:

        git status
        git diff --stat -- package.json pnpm-lock.yaml
        git status --porcelain | grep -v -E '^ ?M (package\.json|pnpm-lock\.yaml)$' || true

    Expected: exactly two files modified: `package.json`, `pnpm-lock.yaml`. The third command should print nothing. If any other file is touched, restore it:

        git checkout -- <offending-file>

4. Commit milestone 2:

        git add package.json pnpm-lock.yaml
        git commit -m "$(cat <<'EOF'
        chore(deps-dev): bump eslint 10.2.1 and mint 4.2.521 (supersedes #66)

        Consolidates Dependabot PR #66:
        - eslint 10.2.0 -> 10.2.1 (patch: bug fixes)
        - mint 4.2.509 -> 4.2.521 (patch)
        EOF
        )"

### Milestone 3 — GitHub Actions SHA bumps + breaking-change scan

1. Scan for any additional `create-github-app-token` usage that might have landed since research:

        git grep -n "create-github-app-token" -- .github/

    Expected: exactly one hit in `.github/workflows/promote.yml:33`. If additional hits appear, add them to the edit list in step 4 below and record the discovery in Surprises & Discoveries.

2. Confirm the usage of `create-github-app-token` in `promote.yml` matches the pre-condition for a safe major-version bump:

        grep -n -A 5 'create-github-app-token' .github/workflows/promote.yml

    Expected snippet:

        33:        uses: actions/create-github-app-token@d72941d797fd3113feb6b93fd0dec494b13a2547 # v1
        34-        id: app-token
        35-        with:
        36-          app-id: ${{ vars.WATERFALL_APP_ID }}
        37-          private-key: ${{ secrets.WATERFALL_APP_PRIVATE_KEY }}

    Verify: runs on `ubuntu-latest` (GitHub-hosted, not self-hosted), no `HTTP_PROXY`/`HTTPS_PROXY`/`NODE_USE_ENV_PROXY` env, inputs used are `app-id` + `private-key`. The `app-id` input is deprecated but still functional in v3.1.1. The step consumes `${{ steps.app-token.outputs.token }}` (verified at line 44) — the `token` output is still present in v3.x. If any of these preconditions are false, STOP and drop PR #65 from the bundle per "Idempotence and Recovery".

3. Edit `.github/workflows/ci.yml`. Replace both occurrences of:

        uses: actions/setup-node@53b83947a5a98c8d113130e565377fae1a50d02f # v6.3.0

    with:

        uses: actions/setup-node@48b55a011bda9f5d6aeb4c2d9c7362e8dae4041e # v6.4.0

4. Edit `.github/workflows/codeql-analysis.yml`. Replace all three occurrences of the SHA `c10b8064de6f491fea524254123dbe5e09572f13` with `95e58e9a2cdfd71adc6e0353d5c52f41a045d225`. The `# v4` trailing comment stays. The three lines to update are the `init`, `autobuild`, and `analyze` uses.

5. Edit `.github/workflows/promote.yml`. Replace the single occurrence of:

        uses: actions/create-github-app-token@d72941d797fd3113feb6b93fd0dec494b13a2547 # v1

    with:

        uses: actions/create-github-app-token@1b10c78c7865c340bc4f6099eb2f838309f1e8c3 # v3.1.1

6. Verify diff scope and content:

        git diff --stat -- .github/
        git grep -n '53b83947a5a98c8d113130e565377fae1a50d02f\|c10b8064de6f491fea524254123dbe5e09572f13\|d72941d797fd3113feb6b93fd0dec494b13a2547' -- .github/

    Expected: three files changed under `.github/workflows/` (`ci.yml`, `codeql-analysis.yml`, `promote.yml`). The `git grep` of old SHAs returns nothing — all replacements are complete. If any old SHA remains, re-apply that specific replacement.

7. Commit milestone 3:

        git add .github/workflows/ci.yml .github/workflows/codeql-analysis.yml .github/workflows/promote.yml
        git commit -m "$(cat <<'EOF'
        chore(deps): bump github actions (supersedes #64, #65)

        Consolidates Dependabot PRs #64 and #65:
        - actions/setup-node v6.3.0 -> v6.4.0 (ci.yml)
        - github/codeql-action v4 SHA bump (codeql-analysis.yml)
        - actions/create-github-app-token v1.12.0 -> v3.1.1 (promote.yml)

        The create-github-app-token major jump is safe for this repo:
        ubuntu-latest runner, no HTTP_PROXY usage, still uses app-id
        (deprecated alias for client-id in v3.1.x but functional).
        EOF
        )"

### Milestone 4 — Full preflight validation

1. Run preflight on the fully-updated branch:

        ./scripts/preflight.sh

    Expected: exit 0. All seven labelled sections complete without error, identical pattern to milestone 1. Note specifically: the `Lint` section rebuilds `tools/lint/lint`, invokes `eslint 10.2.1` (no new rule warnings anticipated since 10.2.0 → 10.2.1 is patch), invokes `cspell 10.0.0` (unchanged), and invokes `mint openapi-check` from mint 4.2.521.

2. If preflight fails:
   - Lint: inspect the failing sub-section. If eslint 10.2.1 surfaces a new warning (unlikely for a patch bump), fix the code or add a targeted override in repo's eslint config; if `mint openapi-check` behavior changed, investigate the OpenAPI spec. If root cause is a single bump, consider dropping it per "Idempotence and Recovery".
   - Go lint / coverage: unrelated to the bumps; investigate only if the baseline passed and these now fail — should not happen absent repo drift.
   - Audit: inspect new advisories; if genuinely new and high-severity, either upgrade further, add to `auditConfig.ignoreCves` with rationale, or drop the offending bump.
   - Shell tests (bats): unrelated; investigate.

3. Record any findings (even "no surprises") in Surprises & Discoveries.

4. No commit unless a fix is required. If a fix is required, commit it as an additional micro-commit with a `fix:` or `chore(deps):` message describing exactly what was adjusted and why.

### Milestone 5 — Publish PR

1. Push the branch:

        git push -u origin chore/combine-dependabot-prs

    Expected: remote branch `origin/chore/combine-dependabot-prs` is created; `gh` and the GitHub Actions CI begin running the `ci.yml` jobs on the branch.

2. Create the consolidated PR:

        gh pr create --base main --head chore/combine-dependabot-prs \
          --title "chore(deps): combine dependabot prs #64, #65, #66 into one update" \
          --body "$(cat <<'EOF'
        ## Summary

        Consolidates three open Dependabot PRs into a single update spanning `package.json`, `pnpm-lock.yaml`, and three GitHub Actions workflow YAMLs, so dependencies advance atomically instead of as three separately-reviewed changes.

        ## Superseded Dependabot PRs

        Close these once this PR merges:

        - #64 `chore(deps): bump the actions-minor group with 2 updates` — actions/setup-node v6.3.0 -> v6.4.0, github/codeql-action v4 SHA bump
        - #65 `chore(deps): bump actions/create-github-app-token from 1.12.0 to 3.1.1`
        - #66 `chore(deps-dev): bump the npm-development group with 2 updates` — eslint 10.2.0 -> 10.2.1, mint 4.2.509 -> 4.2.521

        ## Breaking-change review

        - `actions/setup-node` v6.3.0 -> v6.4.0: dependency-only update, no caller-visible behavior change.
        - `github/codeql-action` SHA bump (stays on v4): bundled CodeQL CLI 2.25.1 -> 2.25.2, TRAP cache feature deprecation announcement (not in use here). No action required.
        - `actions/create-github-app-token` 1.12.0 -> 3.1.1: major jump. Documented breaking changes are (a) Node 24 runtime (action-internal), (b) proxy handling requires `NODE_USE_ENV_PROXY=1` for HTTP_PROXY/HTTPS_PROXY, (c) self-hosted runners need Actions Runner v2.327.1+. Usage here (`.github/workflows/promote.yml`) runs on `ubuntu-latest` (GitHub-hosted) with no proxy config and uses `app-id`/`private-key` inputs (where `app-id` is a deprecated but functional alias for `client-id` in 3.1.x). Inert for this repo.
        - `eslint` 10.2.0 -> 10.2.1: patch-level bug fixes.
        - `mint` 4.2.509 -> 4.2.521: patch.

        ## Test plan

        - [ ] CI `lint` job passes
        - [ ] CI `audit` job passes
        - [ ] CI `shell-tests` job passes
        - [ ] CI `lint-custom-linter` + `test-custom-linter` jobs pass (if triggered by path filter)
        - [ ] Local `./scripts/preflight.sh` exits 0
        EOF
        )"

    Expected: `gh` prints the new PR URL (`https://github.com/mirurobotics/docs/pull/<N>`).

3. Record the PR URL in the Outcomes & Retrospective section of this plan as an amendment commit (optional — only if the plan has been promoted to `plans/active/` or `plans/completed/`):

        # Promote plan to active/ before PR or completed/ after merge as appropriate
        # (Out of scope for this milestone; left to maintainer discretion.)

4. No additional commit beyond Milestone 3's workflow-bump commit; the PR itself is the deliverable.

## Validation and Acceptance

**Preflight must report clean before the PR is published.** Concretely:

1. `./scripts/preflight.sh` run from `/home/ben/miru/workbench5/repos/docs/` exits 0 with the seven section headers in order (`=== Lint Smoke Tests ===`, `=== Go Lint (tools/lint) ===`, `=== Go Coverage (tools/lint) ===`, `=== Lint ===`, `=== Audit ===`, `=== Shell Script Tests ===`). The `=== Lint ===` section ends `All documentation lint checks passed.`.
2. `git diff main...chore/combine-dependabot-prs --stat` lists exactly five changed files: `.github/workflows/ci.yml`, `.github/workflows/codeql-analysis.yml`, `.github/workflows/promote.yml`, `package.json`, `pnpm-lock.yaml`.
3. After `git push`, `gh pr view <new-pr-number> --json statusCheckRollup` reports `SUCCESS` for jobs `lint`, `audit`, `shell-tests` (and `lint-custom-linter` / `test-custom-linter` if path-filtered active; otherwise skipped cleanly).
4. The PR body enumerates #64, #65, #66 as superseded with explicit instructions to close them on merge.
5. After merging the consolidated PR and closing #64/#65/#66, `gh pr list --state open --author app/dependabot` returns no rows for those dependencies until Dependabot's next scan.

**Test steps.** This repo has no unit-test suite for the docs content, so validation is preflight-based:

a. `pnpm run test:lint` — covers `tests/test-lint.sh` (lint smoke tests). Must exit 0.
b. `pnpm exec eslint --max-warnings=0 <mdx targets>` — invoked transitively by `scripts/lint.sh`. Must exit 0.
c. `pnpm exec cspell lint --no-progress --config cspell.json <spell targets>` — invoked by `scripts/lint.sh`. Must exit 0.
d. `pnpm exec mint openapi-check <spec>` — invoked by `scripts/lint.sh` for each OpenAPI yaml under `docs/references/`. Must exit 0 for each spec.
e. `pnpm audit --ignore-registry-errors` — invoked by `scripts/audit.sh`. Must exit 0.
f. `bats pub/scripts/agent/check-miru-access_test.bats` — must exit 0.
g. `LINT_FIX=0 ./tools/lint/scripts/lint.sh` — Go-based custom linter; must exit 0.
h. `./tools/lint/scripts/covgate.sh` — Go coverage gate; must exit 0.

All of (a)–(h) run together in a single `./scripts/preflight.sh` invocation and its exit code is the single acceptance gate.

Acceptance observation: the operator sees on GitHub one PR on `chore/combine-dependabot-prs` whose diff is the five files listed above, all CI jobs green, and the PR body referencing #64/#65/#66 as superseded.

## Idempotence and Recovery

- **Milestone 1 (baseline preflight)** is read-only and fully idempotent. Re-run freely.
- **Milestone 2 (`pnpm install`)** is idempotent: re-running `pnpm install` with the same `package.json` yields the same `pnpm-lock.yaml`. If the `package.json` edit was wrong, correct it and re-run `pnpm install`. To roll back to the pre-edit state of this milestone:

        git checkout -- package.json pnpm-lock.yaml

- **Milestone 3 (workflow YAML edits)** is idempotent: the edits are literal string replacements. If a SHA was pasted wrong, re-edit and the final file content is the same regardless of how many times the replacement was attempted. To roll back:

        git checkout -- .github/workflows/ci.yml .github/workflows/codeql-analysis.yml .github/workflows/promote.yml

- **Dropping a specific bump** (the safety valve if Milestone 4 preflight fails on exactly one bump):
  - For PR #66 (eslint or mint): restore just that version line in `package.json`, re-run `pnpm install`, commit amending the milestone 2 commit or adding a `chore(deps): drop <package> from bundle — <reason>` commit.
  - For PR #64 (setup-node or codeql-action): restore the relevant SHAs in `ci.yml` and/or `codeql-analysis.yml`, commit `chore(deps): drop <action> from bundle — <reason>`.
  - For PR #65 (create-github-app-token): restore the `d72941d797fd3113feb6b93fd0dec494b13a2547 # v1` SHA in `promote.yml`, commit `chore(deps): drop create-github-app-token from bundle — <reason>`. Update the PR body on next push to reflect which PRs remain superseded.
- **Milestone 4 (preflight)** is read-only. Re-run freely.
- **Milestone 5**:
  - `git push -u origin chore/combine-dependabot-prs` is idempotent (subsequent pushes update the remote branch).
  - `gh pr create` is NOT idempotent — a second invocation errors with "a pull request already exists for <branch>". To update the PR body without creating a new PR: `gh pr edit <num> --body "<updated body>"`. To start fresh: `gh pr close <num>` then re-create (only if explicitly needed).
- **If `pnpm install` prompts for `--force` / `--shamefully-hoist` / refuses to install due to peer-dep issues**: STOP, copy the verbatim prompt into Surprises & Discoveries, and evaluate before re-running. Do not blindly pass force flags — they bypass safety checks and can pull in incompatible transitive versions.
- **If CI fails after push** (not preflight — CI on GitHub): `gh pr view <num> --json statusCheckRollup` identifies the failing job. If the failure is the same kind as local preflight would have caught, the local preflight was stale or masked a shell-env difference — fix locally and push again. If the failure is CI-only (e.g. path-filter behavior, secrets availability), decide per the specific job whether to adjust code or drop the offending bump per the single-bump rollback path above.
- **Branch abandonment** (only if the user explicitly asks to abandon):

        git checkout main
        git branch -D chore/combine-dependabot-prs
        git push origin --delete chore/combine-dependabot-prs

  Do NOT delete branches unasked.

# Docs Linting Pipeline via `docs/scripts/lint.sh`

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add the docs lint entrypoint, tool manifests/config, and the GitHub Actions workflow. |
| `./` | read-only | Meta-repo context only; no code changes outside the `docs/` submodule. |

This plan lives in `docs/.ai/exec-plans/backlog/` because all implementation work is owned by the `docs/` repo.

## Purpose / Big Picture

After this change, contributors can run one command from the docs repo, `./scripts/lint.sh`, to validate documentation quality locally and in GitHub Actions. The command will be reproducible, non-interactive, and CI-safe. The first rollout will block only on baseline issues: broken MDX structure, invalid OpenAPI files, and spelling/terminology regressions. Accessibility and stricter prose-style enforcement are explicitly out of scope for this first pass.

## Progress

- [x] (2026-03-13 20:43Z) Create `docs/.ai/exec-plans/backlog/20260314-documentation-linting.md` with this content.
- [x] (2026-03-13 20:43Z) Add Node package metadata and pin tool versions under `docs/` for `eslint`, `eslint-plugin-mdx`, and `mint`.
- [x] (2026-03-13 20:43Z) Add initial lint configuration for MDX and spelling/terminology checks.
- [x] (2026-03-13 20:43Z) Add `docs/scripts/lint.sh` as the canonical local and CI entrypoint.
- [x] (2026-03-13 20:43Z) Add a GitHub Actions workflow in `docs/.github/workflows/` that runs the script on push and pull request.
- [x] (2026-03-14 03:48Z) Validate the script locally in `docs/`, validate `pnpm install --frozen-lockfile`, and confirm the script works from outside the docs repo.
- [x] (2026-03-14 22:19Z) Make `docs/scripts/lint.sh` portable by removing GNU-only `sort -z` usage and add fixture-based smoke tests for MDX, spelling, and OpenAPI failure cases.

## Surprises & Discoveries

- Observation: `docs/` already has `mint`, `eslint`, and `eslint-plugin-mdx` present in `node_modules`, but there is no committed `package.json`, no lockfile, no ESLint config, and no existing docs workflow.
  Evidence: `docs/node_modules/.bin/{mint,eslint}`, empty `find docs -maxdepth 2 -name package.json -o -name pnpm-lock.yaml`.
- Observation: `mint broken-links` is not safe to use directly in CI right now because it parses `docs/.ai` content and prompts interactively about upgrading `mint.json` to `docs.json`.
  Evidence: earlier local runs failed on `.ai/exec-plans/*` parsing and an interactive upgrade prompt despite `docs/docs.json` already existing.
- Observation: `docs/node_modules` uses pnpm-style layout, so `pnpm` is the least surprising package manager for new committed tooling.
  Evidence: `docs/node_modules/.pnpm/` exists and binaries are symlinked from pnpm-managed package paths.
- Observation: `docs/README.md` is already deleted in the docs repo working tree before linting implementation starts, so the lint script must tolerate its absence instead of hard-failing on a missing file.
  Evidence: `git -C docs status --short` reports `D README.md`.
- Observation: `pnpm` would not reuse the preexisting `docs/node_modules` links until it was pointed at the same store path, `/home/ben/.local/share/pnpm/store/v10`.
  Evidence: the first `pnpm add -D cspell` attempt failed with `ERR_PNPM_UNEXPECTED_STORE`; rerunning with `--store-dir /home/ben/.local/share/pnpm/store/v10` succeeded.
- Observation: `mint openapi-check` does not report a generic `OpenAPI definition is invalid` string for invalid fixtures; it emits a schema-validation error such as `Failed to validate OpenAPI schema`.
  Evidence: the first fixture test run against `tests/lint-fixtures/bad-openapi/docs/references/example.yaml` failed with `must have required property 'info'`, so the smoke test had to match the stable prefix instead.

## Decision Log

- Decision: The canonical lint entrypoint will be `docs/scripts/lint.sh`.
  Rationale: The docs repo owns the tooling, config, and workflow; using a meta-repo root script would couple unrelated repos to docs-only behavior.
  Date/Author: 2026-03-13 / Codex + user.
- Decision: Tooling will be committed as Node devDependencies with a lockfile.
  Rationale: CI must be reproducible and non-interactive; runtime installs in shell scripts are slower and less deterministic.
  Date/Author: 2026-03-13 / Codex + user.
- Decision: Initial CI enforcement will be baseline blockers only.
  Rationale: The repo already has existing prose debt; shipping structural and spelling checks first keeps adoption practical.
  Date/Author: 2026-03-13 / Codex + user.
- Decision: Defer `Vale` from the blocking pipeline for now.
  Rationale: `Vale` is valuable, but its style/grammar rules will create a larger cleanup project. First rollout should block on syntax, OpenAPI validity, and spelling/terminology only.
  Date/Author: 2026-03-13 / Codex.
- Decision: Pin `cspell` to `9.7.0` exactly in `package.json`.
  Rationale: the plan requires pinned tool versions; `pnpm add` defaulted to a caret range and was tightened afterward without changing the resolved lockfile.
  Date/Author: 2026-03-14 / Codex.

## Outcomes & Retrospective

Implemented:

- Added `docs/package.json` and `docs/pnpm-lock.yaml` with pinned `mint`, `eslint`, `eslint-plugin-mdx`, and `cspell` dev dependencies.
- Added `docs/eslint.config.mjs` for MDX syntax validation only and `docs/cspell.json` for baseline spelling/terminology checks.
- Added `docs/scripts/lint.sh` as the canonical lint entrypoint used both locally and in GitHub Actions.
- Added `docs/.github/workflows/lint.yml` to run `pnpm install --frozen-lockfile`, the lint smoke tests, and `./scripts/lint.sh` on push and pull request.
- Fixed one existing typo in `docs/docs/learn/deployments.mdx` (`unobstrusive` -> `unobtrusive`) so the new spelling check passes cleanly.
- Replaced GNU-only `sort -z` usage in `docs/scripts/lint.sh` so the script can run on BSD/macOS userlands.
- Added `docs/tests/test-lint.sh` plus fixtures that prove the lint entrypoint fails on invalid MDX, spelling regressions, and invalid OpenAPI input.

Validation completed:

- From `docs/`: `./scripts/lint.sh` passed.
- From `docs/`: `pnpm install --frozen-lockfile --store-dir /home/ben/.local/share/pnpm/store/v10` passed.
- From the meta repo root: `docs/scripts/lint.sh` passed, confirming cwd-independent repo-root resolution.

Remaining follow-up:

- `mint broken-links` remains intentionally deferred because the current CLI flow is interactive and parses `docs/.ai` content. A later pass should add a CI-safe internal-link validator.

## Context and Orientation

The `docs/` repo is a Mintlify documentation site. Its authored content lives in `docs/docs/**/*.mdx`, `docs/snippets/**/*.mdx`, plus repo-level docs like `docs/README.md` and `docs/rclone.md`. API reference sources live in `docs/docs/references/**/*.yaml`.

There is no committed package manifest in `docs/` today, so the implementation must add one. There is also no existing workflow under `docs/.github/workflows/`.

The key repo-specific constraints are:

- `mint openapi-check` works and should be kept for OpenAPI validation.
- `mint a11y` is intentionally not part of this plan.
- `mint broken-links` is currently unsuitable for the blocking pipeline because it is interactive and scans non-publishable `.ai` files.
- The script must be runnable as `./scripts/lint.sh` and must be the same command used by CI.

## Plan of Work

Add package management first. Create `docs/package.json` and `docs/pnpm-lock.yaml`, with a `lint` script that shells out to `./scripts/lint.sh`. Pin the toolchain to `pnpm`, `eslint`, `eslint-plugin-mdx`, and `cspell`. Use `mint` through the committed dependency already represented by the package manifest so the workflow no longer depends on an untracked local installation state.

Add an ESLint flat config at `docs/eslint.config.mjs`. Scope it to `docs/**/*.mdx` and `snippets/**/*.mdx`, using `eslint-plugin-mdx` for MDX parsing and syntax validation only. Keep rules conservative in the first pass: parseability, import syntax, and obvious JSX/MDX errors. Do not add stylistic formatting rules.

Add a spelling config at `docs/cspell.json`. Include `README.md`, `rclone.md`, authored MDX content, and optionally `docs.json`. Exclude `node_modules`, `.ai`, generated endpoint output, images, and OpenAPI YAML files. Seed the custom dictionary with known product and domain terms such as `Miru`, `Mintlify`, `OpenAPI`, `config`, `configs`, `GitHub`, `Cloudflare`, `systemd`, and version/train names already in use. The goal is to catch typo regressions, not to fight intentional product vocabulary.

Add `docs/scripts/lint.sh` as the canonical orchestrator. It must:

- start with `#!/usr/bin/env bash` and `set -euo pipefail`;
- resolve the repo root relative to the script location so it can be run from any working directory;
- verify `pnpm` is available and fail with a short actionable message if it is not;
- run `pnpm exec eslint` on authored MDX files;
- run `pnpm exec cspell` on authored docs targets;
- run `pnpm exec mint openapi-check` for every `docs/docs/references/**/*.yaml` file;
- print short section headers so failures are easy to read in CI logs.

Do not include `mint broken-links` in the first version of `scripts/lint.sh`. Instead, add a documented follow-up item in this ExecPlan to evaluate either a custom non-interactive internal-link checker or a cleaned temporary-workspace wrapper around Mint. The initial rollout should favor reliability over breadth.

Add `docs/.github/workflows/lint.yml` in the docs repo. Configure it to run on `pull_request` and `push` for paths under the docs repo. The workflow should:

- check out the docs repo;
- set up Node 22;
- set up `pnpm`;
- run `pnpm install --frozen-lockfile`;
- run `chmod +x scripts/lint.sh` if needed;
- execute `./scripts/lint.sh`.

Ensure the workflow runs from the docs repo root, not the meta repo root. The plan should call this out explicitly because `docs/` is a submodule in the meta workspace but a standalone repository in its own CI.

## Concrete Steps

From `docs/`, create and populate the package manifest and configs.

From `docs/`, install and lock dependencies:

    pnpm install

From `docs/`, run the full lint command:

    ./scripts/lint.sh

Expected successful transcript shape:

    == ESLint (MDX) ==
    ... no errors ...
    == CSpell ==
    ... no errors ...
    == OpenAPI ==
    Checking docs/docs/references/device-api/v0.1.0.yaml
    success OpenAPI definition is valid.
    ...
    All documentation lint checks passed.

From `docs/`, verify the GitHub Actions workflow definition is syntactically present:

    ls .github/workflows/lint.yml

From `docs/`, after implementation, rerun the same command the workflow uses:

    pnpm install --frozen-lockfile
    ./scripts/lint.sh

## Validation and Acceptance

Acceptance is complete when all of the following are true:

- Running `./scripts/lint.sh` from `docs/` succeeds on a clean checkout with dependencies installed.
- Running `./scripts/lint.sh` from any other working directory still succeeds because the script resolves `docs/` correctly.
- Introducing a typo into an authored `.mdx` file causes `cspell` to fail and the script exits non-zero.
- Introducing invalid MDX syntax into an authored `.mdx` file causes `eslint` to fail and the script exits non-zero.
- Breaking one OpenAPI YAML file causes `mint openapi-check` to fail and the script exits non-zero.
- The GitHub Actions workflow runs the same `./scripts/lint.sh` entrypoint and fails when any of the above failures are present.
- Accessibility is not part of the acceptance criteria for this plan.

## Idempotence and Recovery

`pnpm install`, `./scripts/lint.sh`, and the workflow runs are all idempotent and safe to repeat. If implementation introduces noisy false positives in `cspell`, recovery is to add only narrowly justified words to `docs/cspell.json` rather than broad ignore patterns. If the shell script fails because `pnpm` is missing, recovery is to install `pnpm` locally or let GitHub Actions provide it via setup steps; do not add runtime package-manager bootstrapping to the lint script itself.

## Assumptions and Defaults

- Package manager: `pnpm`.
- Canonical entrypoint: `docs/scripts/lint.sh`.
- Blocking checks in v1: ESLint for MDX structure, CSpell for spelling/terminology, Mint OpenAPI validation.
- Deferred from v1: `Vale`, accessibility checks, and Mint's current `broken-links` command until a non-interactive CI-safe approach is chosen.

Revision note (2026-03-14): Updated the ExecPlan after implementation to record the actual dependency setup, pnpm store discovery, validation evidence, and the remaining broken-links follow-up.

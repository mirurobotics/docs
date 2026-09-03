# Document post-v0.12.0 CLI validate commands

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/ben/miru/workbench3/repos/docs`) | read-write | Add CLI reference documentation, snippets, and navigation entries for new Miru CLI commands. |
| `cli-private/` (`/home/ben/miru/workbench3/repos/cli-private`) | read-only | Source of truth for the CLI command tree and command help text. Do not modify. |
| `openapi/` (`/home/ben/miru/workbench3/repos/openapi`) | read-only | Source of truth for generated CLI API paths and visible API scope metadata. Do not modify. |
| `backend/` (`/home/ben/miru/workbench3/repos/backend`) | read-only | Source of truth for authorization rules when OpenAPI scope metadata is incomplete. Do not modify. |

This plan lives in `/home/ben/miru/workbench3/repos/docs/plans/` because all planned file edits are in the docs repo.

Working branch at authoring time: `docs/update-cli-command-reference`. Base branch for the PR: `main`.

## Purpose / Big Picture

The CLI Reference product currently documents commands shipped through CLI `v0.12.0`, but the CLI implementation has two new validation commands on `origin/main`: `miru config validate` and `miru deployment validate`. After this change, users browsing the CLI Reference can find both commands in the sidebar, see their required arguments and flags, and understand what the commands validate without reading the private CLI source.

## Progress

- [x] Milestone 1: Re-check the post-release command inventory, add the `miru config validate` reference page/snippets/nav entry, and commit.
- [x] Milestone 2: Add the `miru deployment validate` reference page/snippets/nav entry, and commit.
- [x] Milestone 3: Run local docs checks, fix findings, and commit any fixes.
- [x] Milestone 4: Push the branch, open or update a draft PR, and drive preflight to `CLEAN`.

## Surprises & Discoveries

- 2026-09-03 — Re-running `cli-private` inventory after `git fetch --all --tags --prune` still found `v0.12.0` as the latest stable tag and `v0.12.1-beta.4` as the newest beta tag. No additional post-`v0.12.0` user-facing commands appeared beyond `miru config validate` and `miru deployment validate`.
- 2026-09-03 — Source refine found the deployment validate prose needed to say optional slots are skipped only when their files are absent. If an optional slot file exists, the CLI reads and validates it.
- 2026-09-03 — Local validation passed with no docs fixes required, so no validation-fix commit was needed.
- 2026-09-03 — Draft PR #176 was opened and GitHub CI passed on the implementation head. A final plan-status commit was then pushed and watched as the branch head.

## Decision Log

- Decision: Treat `v0.12.0` as the last stable CLI release for this plan and document only command reference pages for post-release commands.
  Rationale: After `git fetch --all --tags --prune` on 2026-09-03, `cli-private` had stable tag `v0.12.0`, beta tags `v0.12.1-beta.1` through `v0.12.1-beta.4`, and no stable `v0.12.1` tag. The docs changelog already opens with `# v0.12.0`; adding release notes for a beta-only line is outside this command-reference task.
  Date/Author: 2026-09-03 / plan author

## Outcomes & Retrospective

Implemented the planned CLI reference coverage on `docs/update-cli-command-reference`. The docs now include `miru config validate` and `miru deployment validate` reference pages, reusable snippets for usage, flags, scopes, and examples, and navigation entries under CLI Reference. Scope snippets were checked against the CLI, backend authorization, and OpenAPI sources. Local validation passed with `pnpm run test:lint`, `./scripts/lint.sh`, and `pnpm run validate`; draft PR #176 was opened for CI preflight.

## Context and Orientation

The docs repo root is `/home/ben/miru/workbench3/repos/docs`. Mintlify content lives under `/home/ben/miru/workbench3/repos/docs/docs/`. A flat file `docs/references/cli/config-validate.mdx` becomes URL `/references/cli/config-validate`; a flat file `docs/references/cli/deployment-validate.mdx` becomes URL `/references/cli/deployment-validate`. Navigation is configured in `docs/docs.json` under the `"product": "CLI Reference"` object. Existing CLI reference pages use thin MDX pages that import reusable snippets from `docs/snippets/references/cli/<command-group>/<command>/`.

At authoring time, `docs/docs.json` lists these CLI Reference groups: Authentication, Schemas, Releases, Devices, Deployments, and Utilities. Existing command reference files are `docs/references/cli/login.mdx`, `whoami.mdx`, `schema-validate.mdx`, `release-create.mdx`, `device-clone.mdx`, `device-stage.mdx`, `deployment-clone.mdx`, and `version.mdx`.

The private CLI repo root is `/home/ben/miru/workbench3/repos/cli-private`. The latest stable release tag found during planning is `v0.12.0`. Comparing `v0.12.0..origin/main` after fetching showed only two new command registrations:

- `miru config validate`: `origin/main:internal/commands/root/root.go` adds `cfgcmds.New()`, `origin/main:internal/commands/config/root.go` defines `Use: "config"` and registers `NewValidate()`, and `origin/main:internal/commands/config/validate.go` defines `Use: "validate <file>"`.
- `miru deployment validate`: `origin/main:internal/commands/deployment/root.go` registers `NewValidate()`, and `origin/main:internal/commands/deployment/validate.go` defines `Use: "validate"`.

No other post-`v0.12.0` `AddCommand` or `Use:` diff represented a new user-facing command. The implementer must re-run the command inventory before editing because `origin/main` can move.

`miru config validate` validates one local config instance file against a persisted config schema without creating anything. It requires authentication. The schema is selected either by full schema ID with `--schema-id cfg_sch_...` or by `--release <version>` plus `--config-type <slug-or-name>`. `--schema-id` is mutually exclusive with `--release` and `--config-type`; `--release` and `--config-type` are required together. A dashboard short ID such as `SCH-KKV8s` is not accepted. Success prints `Config Valid`.

`miru deployment validate` validates local config files against every required slot in a release without creating a deployment. It requires `--release <version>` and accepts optional `--root <dir>`, defaulting to the current directory. Required slots with no file fail; optional slots with no file are skipped; optional slot files that exist are validated; extra files on disk are ignored. Success prints the file tree followed by `Deployment Valid`.

API key scopes must be checked from source before writing the snippets. At authoring time, backend authorization showed config-instance validation needs the config schema get permission, exposed to users as `config_schemas:read` (`backend/internal/configs/services/config_instances/authz.go`). Selecting a schema by release also needs `releases:read`. `miru deployment validate` fetches a release with expanded config schemas and then validates each local instance, so it needs `releases:read` and `config_schemas:read`. It does not read an existing deployment, so do not list `deployments:read` unless the implementation has changed.

## Plan of Work

First, re-run the source inventory from `cli-private` and confirm the planned command set is still complete. If another new user-facing command exists after `v0.12.0`, add it to this plan before implementation and commit the plan revision from the docs repo.

For `miru config validate`, create `docs/references/cli/config-validate.mdx` using the existing page pattern from `docs/references/cli/schema-validate.mdx`: frontmatter title `Validate`, imports for `Arguments`, `Examples`, `Flags`, `Scopes`, and `Usage`, a concise description, then sections for API key scopes, usage, arguments, flags, and examples. Create snippets under `docs/snippets/references/cli/config/validate/`:

- `usage.mdx`: show `miru config validate {file} --schema-id {config-schema-id}` and `miru config validate {file} --release {release-version} --config-type {config-type}` in tabs or a code group.
- `args.mdx`: define required `file` as the local config instance file to validate.
- `flags.mdx`: define `--schema-id`, `--release`, and `--config-type`, including the mutual-exclusion and required-together rules.
- `scopes.mdx`: list `config_schemas:read`; also list `releases:read` for the release/config-type selector.
- `examples.mdx`: show one successful schema-ID validation and one successful release/config-type validation ending with `Config Valid`.

Add a new `"Configs"` group to the CLI Reference navigation in `docs/docs.json`, adjacent to `"Schemas"`, containing `"references/cli/config-validate"`.

For `miru deployment validate`, create `docs/references/cli/deployment-validate.mdx` using the existing page pattern from `docs/references/cli/deployment-clone.mdx` and `docs/references/cli/device-stage.mdx`. Import `Examples`, `Flags`, `Scopes`, and `Usage`; no arguments section is needed because the command has no positional arguments. Create snippets under `docs/snippets/references/cli/deployment/validate/`:

- `usage.mdx`: show the default current-directory form and a custom-root form.
- `flags.mdx`: define required `--release` and optional `--root`; `--root /` reads files at their actual device paths.
- `scopes.mdx`: list `releases:read` and `config_schemas:read`.
- `examples.mdx`: show a successful validation transcript with `Validating release`, a small file tree, and `Deployment Valid`.

Add `"references/cli/deployment-validate"` to the existing `"Deployments"` group in `docs/docs.json` after `"references/cli/deployment-clone"`.

Do not edit `docs/changelog/cli.mdx` in this plan unless the task owner explicitly expands scope. If a stable `v0.12.1` release exists by implementation time, record that in Surprises & Discoveries and leave a follow-up note; this plan is for command reference coverage.

## Concrete Steps

All docs commands run from `/home/ben/miru/workbench3/repos/docs` unless another working directory is stated.

### Milestone 1: inventory and config validate docs

1. From `/home/ben/miru/workbench3/repos/cli-private`, refresh refs and compare the command tree:

       git fetch --all --tags --prune
       git tag --sort=-creatordate | sed -n '1,12p'
       git log --oneline --decorate v0.12.0..origin/main -- internal/commands
       git diff --name-status v0.12.0..origin/main -- internal/commands
       git grep -n 'AddCommand\|Use:' origin/main -- internal/commands

   Expect the stable release baseline to include `v0.12.0`, and expect the new command set to include `config validate` and `deployment validate`. If additional user-facing commands appear, update this plan first and commit that plan revision from `/home/ben/miru/workbench3/repos/docs`.

2. From `/home/ben/miru/workbench3/repos/docs`, confirm the branch and clean tree:

       git branch --show-current
       git status --short

   Expect `docs/update-cli-command-reference` and no unrelated local modifications.

3. Create the config validate page and snippets:

       mkdir -p docs/snippets/references/cli/config/validate

   Add `docs/references/cli/config-validate.mdx`, `docs/snippets/references/cli/config/validate/usage.mdx`, `args.mdx`, `flags.mdx`, `scopes.mdx`, and `examples.mdx` as described in Plan of Work.

4. Edit `docs/docs.json` so the CLI Reference product has a `"Configs"` group containing `"references/cli/config-validate"` next to the existing `"Schemas"` group.

5. Verify the milestone diff:

       grep -n "references/cli/config-validate" docs/docs.json
       git diff --stat

   Expect exactly the config validate page/snippets, `docs/docs.json`, and this plan file if it has not already been committed.

6. Commit the milestone:

       git add docs/docs.json docs/references/cli/config-validate.mdx docs/snippets/references/cli/config/validate/args.mdx docs/snippets/references/cli/config/validate/examples.mdx docs/snippets/references/cli/config/validate/flags.mdx docs/snippets/references/cli/config/validate/scopes.mdx docs/snippets/references/cli/config/validate/usage.mdx
       git commit -m "docs(cli): add config validate reference"

### Milestone 2: deployment validate docs

1. Create the deployment validate snippets:

       mkdir -p docs/snippets/references/cli/deployment/validate

   Add `docs/references/cli/deployment-validate.mdx`, `docs/snippets/references/cli/deployment/validate/usage.mdx`, `flags.mdx`, `scopes.mdx`, and `examples.mdx` as described in Plan of Work.

2. Edit `docs/docs.json` so the existing `"Deployments"` group contains both `"references/cli/deployment-clone"` and `"references/cli/deployment-validate"`, in that order.

3. Verify the milestone diff:

       grep -n "references/cli/deployment-validate" docs/docs.json
       git diff --stat

   Expect exactly the deployment validate page/snippets and `docs/docs.json`, plus any uncommitted plan-progress updates made during implementation.

4. Commit the milestone:

       git add docs/docs.json docs/references/cli/deployment-validate.mdx docs/snippets/references/cli/deployment/validate/examples.mdx docs/snippets/references/cli/deployment/validate/flags.mdx docs/snippets/references/cli/deployment/validate/scopes.mdx docs/snippets/references/cli/deployment/validate/usage.mdx
       git commit -m "docs(cli): add deployment validate reference"

### Milestone 3: local validation

1. If dependencies are missing, install them:

       pnpm install --frozen-lockfile

2. Run the docs checks:

       pnpm run test:lint
       ./scripts/lint.sh
       pnpm run validate

   Expect `pnpm run test:lint` to exit 0, `./scripts/lint.sh` to end with `All documentation lint checks passed.`, and `pnpm run validate` to exit 0 without broken navigation or MDX errors.

3. Fix any finding at its source and re-run the failed command, then re-run all three commands. If `cspell` flags a genuine project term, add it to `cspell.json` in sorted order.

4. Commit validation fixes if there are any:

       git status --short
       git add <exact-files-that-changed>
       git commit -m "docs(cli): fix validate command reference checks"

   If there are no changes after the checks, record that no validation-fix commit was needed in Progress during implementation.

### Milestone 4: PR and preflight

1. Rebase onto the requested base branch and push:

       git fetch origin main
       git rebase origin/main
       git push -u origin HEAD

   If a previously pushed branch was rebased, use `git push --force-with-lease` instead of a plain force push.

2. Open a draft PR if one does not exist, or update the existing draft PR body. The PR test plan must include:

       pnpm run test:lint
       ./scripts/lint.sh
       pnpm run validate
       preflight CLEAN on the pushed branch head

3. Watch CI for the current pushed head SHA:

       git rev-parse HEAD
       gh pr checks --watch

   If any check fails, inspect the failed run logs with `gh run view <run-id> --log-failed`, fix the issue in a new commit, push once, and watch checks again.

4. Commit final plan status updates if this ExecPlan was moved to `plans/active/` for implementation and its living sections changed:

       git add plans/active/20260903-cli-validate-command-docs.md
       git commit -m "docs(plans): record CLI validate command docs outcome"

## Validation and Acceptance

Acceptance criteria:

1. Re-running the command inventory from `/home/ben/miru/workbench3/repos/cli-private` confirms that every post-`v0.12.0` user-facing CLI command is documented. At minimum, `miru config validate` and `miru deployment validate` are included.
2. `/home/ben/miru/workbench3/repos/docs/docs/references/cli/config-validate.mdx` exists and renders a page for `miru config validate` with usage, required `file` argument, flags `--schema-id`, `--release`, and `--config-type`, API key scopes, and success examples ending in `Config Valid`.
3. `/home/ben/miru/workbench3/repos/docs/docs/references/cli/deployment-validate.mdx` exists and renders a page for `miru deployment validate` with usage, flags `--release` and `--root`, API key scopes, and a success example ending in `Deployment Valid`.
4. `docs/docs.json` lists `references/cli/config-validate` in a CLI Reference `"Configs"` group and lists `references/cli/deployment-validate` in the existing `"Deployments"` group after `references/cli/deployment-clone`.
5. The scope snippets match the verified authorization rules: `miru config validate --schema-id` needs `config_schemas:read`; the release/config-type selector additionally needs `releases:read`; `miru deployment validate` needs `releases:read` and `config_schemas:read`. If source behavior changes, the snippets and this criterion are updated together.
6. From `/home/ben/miru/workbench3/repos/docs`, `pnpm run test:lint`, `./scripts/lint.sh`, and `pnpm run validate` all exit 0.
7. On a local Mintlify preview (`pnpm run dev`) or the PR preview, the CLI Reference sidebar shows the new Configs/Validate page and the Deployments/Validate page, and both pages render without broken imports.
8. Preflight reports `CLEAN`: CI is green on the pushed branch head. This is mandatory before the PR leaves draft and before the task is reported complete. Local checks are required but are not a substitute for `CLEAN` on the pushed head.

## Idempotence and Recovery

The source inventory commands and local validation commands are read-only except for dependency installation and generated local build artifacts that are ignored by git. Creating snippet directories with `mkdir -p` is safe to repeat. Rewriting a page or snippet with the same content is a no-op; before adding a nav item, use `grep` to avoid duplicate entries.

If a docs edit is wrong before commit, restore only the affected file from the base branch with `git restore --source=origin/main -- <path>` after confirming the path is part of this task. If a milestone commit has already been pushed, fix forward with a new commit; do not rewrite public history except for `git push --force-with-lease` after an intentional rebase. To abandon the implementation, revert the milestone commits from `/home/ben/miru/workbench3/repos/docs` with `git revert <commit>`.

# Add a CLI reference page for the `miru version` command

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add `docs/references/cli/version.mdx` and a nav entry in `docs/docs.json`. |
| `cli-private/` | read-only | Source of truth for the command's behavior (`internal/commands/version/version.go`). Verified during planning; no changes. |

This plan lives in `docs/plans/` because all edits are confined to the docs repo.

Working branch: `claude/miru-version-command-docs-6m714w` (already checked out — do not create or switch branches). Base branch for the PR: `main`.

## Purpose / Big Picture

The CLI Reference product on the Mintlify docs site currently documents only two commands: `login` and `release create`. The `miru version` command exists in the CLI (it is even used by the install guide's "Verify" section) but has no reference page, so it does not appear in the CLI Reference navigation.

After this change, a reader browsing the CLI Reference sees a "Version" page under a new "Utilities" group. The page states what the command prints (version, git commit, build date), shows the exact invocation and a sample transcript, and links to the upgrade steps and the CLI changelog.

## Progress

- [ ] Milestone 1: create `docs/references/cli/version.mdx` and add the docs.json nav entry; commit.
- [ ] Milestone 2: run the lint/validate suite (`pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`), fix any findings, commit fixes if any.
- [ ] Preflight CLEAN: CI green on the pushed branch head before the PR leaves draft.

Use timestamps when you complete steps. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

Decisions made during authoring (date 2026-08-18, author: plan agent):

- Decision: place the page in a new "Utilities" nav group after "Releases".
  Rationale: the existing CLI Reference groups ("Authentication", "Releases") are per-capability; `version` fits neither. Other products in `docs/docs.json` group by resource or capability, so a general-purpose group is needed; "Utilities" keeps the list alphabetical (Authentication, Releases, Utilities) and leaves room for future commands like `help` or `completion`.

- Decision: no snippet file under `docs/snippets/references/cli/`.
  Rationale: snippets exist for content reused across pages — `snippets/references/cli/login.mdx` is imported by three pages, `install/install.mdx` by two. The `version` usage appears nowhere else (the install guide's transcript at `docs/developers/cli/install.mdx:16-21` serves a different purpose — verifying an install — and is not being consolidated here). Single-use content is written inline, matching the "no premature snippet" pattern.

- Decision: no entry in `docs/changelog/cli.mdx`.
  Rationale: that changelog tracks CLI binary releases (`# v0.11.0`, `# v0.10.3`, ...), not documentation additions. The `version` command is not new — it already ships (the install guide shows it against v0.10.0).

- Decision: sample transcript uses version `0.11.0` with an illustrative commit hash and build date.
  Rationale: v0.11.0 is the newest release in `docs/changelog/cli.mdx` (August 12, 2026). The transcript format mirrors `docs/developers/cli/install.mdx:16-21`; the commit hash and timestamp are illustrative, as they are on that page.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The docs repo root is `/home/user/docs`. The Mintlify content root is the `docs/` subdirectory: a file `docs/references/cli/version.mdx` is served at the URL `/references/cli/version`. Navigation is defined in `docs/docs.json`.

The command being documented lives in the private CLI repo at `cli-private/internal/commands/version/version.go`. Verified behavior:

- Cobra command `Use: "version"`, `Short: "Show the CLI version information"`, example `miru version`.
- It takes no flags (`DisableFlagsInUseLine: true`) and no arguments, requires no authentication, and makes no network calls.
- It prints exactly three lines: `Version: <semver>`, `Commit: <git commit>`, `Built: <build date>`. The values are injected at build time via ldflags (`cli-private/internal/version/version.go`).

Key docs files:

- `docs/references/cli/login.mdx` — the thin-page template to imitate: frontmatter with `title`, one or two sentences of prose, a `### Usage` heading. (`release-create.mdx` is the heavyweight variant with scopes/flags snippets; `version` has neither, so the login shape fits.)
- `docs/docs.json` — the `"product": "CLI Reference"` block (around line 208) holds two groups: `"Authentication"` with `["references/cli/login"]` and `"Releases"` with `["references/cli/release-create"]`.
- `docs/developers/cli/install.mdx` — already references `miru version` in its "Verify" section (lines 14-21) and defines the `## Upgrade` heading the new page links to (`/developers/cli/install#upgrade`).
- `docs/changelog/cli.mdx` — target of the page's changelog link (`/changelog/cli`).

Lint constraints that shape the page (custom Go linter under `tools/lint/`): headings must be sentence case (`headingcase`), prose must use an em dash instead of `--` (`nodoubledash`), and `docs.json` redirects are validated (`redirects` — no redirect is needed here since no URL moves). The page has no imports, so the import-block rules do not apply.

## Plan of Work

Two edits, both under `/home/user/docs`.

1. Create `docs/references/cli/version.mdx` with exactly this content:

       ---
       title: "Version"
       ---

       Print the version of the installed CLI, along with the git commit and build date it was built from.

       `version` runs entirely offline — it takes no flags and does not require authentication.

       ### Usage

       ```bash
       miru version
       ```

       ```
       $ miru version
       Version: 0.11.0
       Commit: 3f9c1d27b45e80aa612f9dc0be174c58a2e6f4b1
       Built: 2026-08-12T18:24:09Z
       ```

       To upgrade an out-of-date CLI, follow the [upgrade steps](/developers/cli/install#upgrade). Release history lives in the [CLI changelog](/changelog/cli).

2. Edit `docs/docs.json`: inside the `"product": "CLI Reference"` block, append a third group object after the `"Releases"` group (add a comma after the Releases group's closing brace):

       {
         "group": "Utilities",
         "pages": [
           "references/cli/version"
         ]
       }

Do not touch any other file. In particular, do not edit `docs/developers/cli/install.mdx`, `docs/changelog/cli.mdx`, or any snippet.

## Concrete Steps

All commands run from `/home/user/docs`.

### Setup

1. Confirm branch and clean tree:

       git branch --show-current   # expect: claude/miru-version-command-docs-6m714w
       git status --short          # expect: only this plan file, until Milestone 1 begins

### Milestone 1: page and nav entry

1. Create `docs/references/cli/version.mdx` with the content from Plan of Work step 1.

2. Edit `docs/docs.json` per Plan of Work step 2, then verify:

       grep -n "references/cli/version" docs/docs.json   # expect 1 match, inside a "Utilities" group
       grep -n -A 22 '"CLI Reference"' docs/docs.json    # expect groups: Authentication, Releases, Utilities

3. Confirm the diff scope:

       git diff --stat   # expect exactly: docs/docs.json and docs/references/cli/version.mdx (plus this plan file if not yet committed)

4. Commit the milestone:

       git add docs/references/cli/version.mdx docs/docs.json plans/backlog/20260818-miru-version-command-docs.md
       git commit -m "docs(references/cli): add version command reference page"

### Milestone 2: lint and validate

1. Install dependencies if needed, then run the three checks CI's `lint` job runs:

       pnpm install --frozen-lockfile
       pnpm run test:lint    # lint smoke tests; expect exit 0
       ./scripts/lint.sh     # MDX prose linter, ESLint-MDX, cspell, OpenAPI; expect "All documentation lint checks passed."
       pnpm run validate     # mint validate — the MDX/JSX compile and nav/link check; expect success, no broken-link warnings

2. Fix any finding at its source and re-run until all three exit 0. Only add a word to `cspell.json` if it is a genuine project-wide proper noun (none is expected — the page introduces no new jargon).

3. Optionally run the full local gate:

       ./scripts/preflight.sh   # expect exit 0

4. If step 2 produced fixes, commit them:

       git add -A
       git commit -m "docs(references/cli): fix lint findings on version page"

5. Push and watch CI on the branch head:

       git push -u origin claude/miru-version-command-docs-6m714w
       gh pr checks --watch

## Validation and Acceptance

Acceptance criteria — each must be observably true:

1. `docs/references/cli/version.mdx` exists, its frontmatter title is `Version`, and it contains the `miru version` usage block, the three-line sample transcript (`Version:`, `Commit:`, `Built:`), and working links to `/developers/cli/install#upgrade` and `/changelog/cli`.
2. `docs/docs.json` lists `references/cli/version` under a `Utilities` group in the CLI Reference product, after `Authentication` and `Releases`; no other part of `docs.json` changes.
3. `git diff main --stat` shows exactly two content files changed (`docs/docs.json`, `docs/references/cli/version.mdx`) plus this plan file.
4. From `/home/user/docs`: `pnpm run test:lint` exits 0; `./scripts/lint.sh` exits 0 ending with `All documentation lint checks passed.`; `pnpm run validate` succeeds with no errors (this is the check that proves the new page compiles and the nav entry resolves).
5. On a rendered preview (`pnpm dev`, or the PR preview if no local server is feasible), the CLI Reference sidebar shows a "Utilities" group containing "Version", and the page renders with the usage and transcript blocks. If no preview is feasible, record that in Surprises & Discoveries and defer to the PR preview — do not claim a visual check that was not performed.
6. **Preflight reports CLEAN**: CI is green on the pushed branch head (`gh pr checks` all passing). This must hold before the PR leaves draft and before the task is reported complete.

## Idempotence and Recovery

Both edits are additive and safe to repeat: re-writing `version.mdx` with the same content is a no-op, and the docs.json group can be re-checked with the greps in Milestone 1 before re-editing. If a lint or validate failure appears after a commit, fix the cause in a new commit — do not amend a pushed commit. To revert entirely, `git revert` the milestone commit(s); no redirects, external state, or other pages depend on the new page.

# Document how file rules behave when deployed to a device

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`mirurobotics/docs`, this repo) | read-write | Add a `## Deployment` section to `docs/data-uploads/primitives/file-rules/overview.mdx`; add one cross-reference paragraph to `docs/data-uploads/overview.mdx` and one sentence to `docs/snippets/file-rules/retention.mdx`. Touch `cspell.json` only if CSpell flags a new word. |
| `agent/` (`mirurobotics/agent`) | read-only | Source of truth for the behavior: `agent/src/data_uploads/scan/{rule,scanner,state}.rs`, `agent/src/workers/sync_scan_bridge.rs`, `agent/src/disk/{file_rules,deployments}.rs`, `agent/src/data_uploads/{upload,retention}/sink.rs`. Do not modify. |
| `backend/` (`mirurobotics/backend`) | read-only | Source of truth for rule identity: `internal/configs/domain/filerules/hash.go`, `internal/configs/services/file_rules/create.go`. Do not modify. |

This plan lives in `docs/plans/` because every write is a documentation edit in this repo. Work happens on the already-checked-out branch `claude/adoring-bohr-mdjcop` (base `main`). Do not create another branch.

## Purpose / Big Picture

A customer asked (in Slack) why files already on a robot were not uploaded after a file rule was deployed. The answer was: the Miru Agent skips files that already exist on the device when the rule takes effect there. Nothing in the docs says this, and `docs/data-uploads/primitives/file-rules/rule-definition.mdx` only says that files that "appear after a rule is deployed can still match it". After this change, `/data-uploads/primitives/file-rules/overview#deployment` explains:

- when a rule takes effect on a device (per device, when that device applies the deployment)
- that files already present are neither uploaded nor deleted by retention
- that such a file is picked up once it changes
- that an unchanged rule keeps its snapshot across later releases, while a changed rule, a widened glob, or a removed-and-re-added rule starts a new one
- how to deliberately upload a backlog

The data uploads overview and the retention field point readers to that section.

## Progress

- [ ] Milestone 1: re-run the agent behavior check (test step), apply Edits 1-3, pass local checks, commit.
- [ ] Milestone 2: push, open a draft PR, get preflight to `CLEAN`, mark the PR ready.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

- Decision: Docs follow the agent code, which refines the Slack claim in three places. (1) A pre-existing file *is* uploaded once its size or modification time changes. (2) The snapshot belongs to the rule on that device and survives later deployments while the rule's content is unchanged, so it is not reset on every deployment. (3) Retention also skips pre-existing files.
  Rationale: see the verified findings in Context and Orientation.
  Date/Author: 2026-09-22, plan author.
- Decision: Put the full explanation in `file-rules/overview.mdx` (`## Deployment`), and only add short pointers elsewhere. No changelog entry.
  Rationale: the overview page is the file rule primitive reference and already has `## Immutability`, which talks about deploying rules. Changelogs record product releases, not docs clarifications (precedent: `plans/completed/20260922-source-glob-syntax-docs.md`, `plans/completed/20260726-opaque-schema-language.md`).
  Date/Author: 2026-09-22, plan author.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX (Markdown plus JSX components) under `docs/`, and navigation is in `docs/docs.json` (no change needed, because no page is added).

Terms:

- A **file rule** is a YAML object shipped in a release. It tells the Miru Agent (the on-device binary) which files to upload and/or delete.
- A **deployment** puts one release on one device. Its activity status becomes `deployed` once the device has applied it (`docs/primitives/deployments.mdx`, `## Activity status`).
- A **stable** file is one whose size and modification time did not change for `source.stability_window_secs`. Only stable files are uploaded or queued for retention deletion.

Files to change:

- `docs/data-uploads/primitives/file-rules/overview.mdx`: sections `## Properties`, `## Immutability` (lines 57-61), `## Git provenance` (63), `## Create a file rule` (71), `## View a file rule` (79). The `digest` field (lines 24-32) already says rules are deduplicated by digest within a workspace.
- `docs/data-uploads/overview.mdx`: `## File detection` (lines 27-35) explains quiescence and nothing else.
- `docs/snippets/file-rules/retention.mdx`: imported as `<Retention />` by `docs/data-uploads/define-file-rules.mdx` (under `### Retention`) and `docs/data-uploads/primitives/file-rules/rule-definition.mdx` (under `## Retention`). Do not add headings inside the snippet, because the two pages use different heading levels.

### Verified behavior (agent source, `mirurobotics/agent`, paths relative to the `agent/agent/` crate)

The new text states these facts and nothing more.

1. **Snapshot at apply time.** When a rule ID is new to the scanner, `SingleThreadScanner::update_rules` creates `RuleScanner::new` (`src/data_uploads/scan/scanner.rs:169-181`). This globs the pattern and records every matching readable file's size and mtime as `preexisting` (`src/data_uploads/scan/rule.rs:37-46`, `167-176`; `state.rs:15`, `83-96`). The cutoff is this moment, not a timestamp comparison.
2. **Trigger is per device.** `update_rules` is called by the sync-scan bridge at agent startup and after every successful sync (`src/workers/sync_scan_bridge.rs:56`, `59-62`, `79-99`). It is called with the rules of the one deployment whose local activity status is `Deployed` (`src/disk/deployments.rs:21-30`, `src/disk/file_rules.rs:9-38`; that status is set in `src/deploy/fsm.rs:86`). So the cutoff is when *that device* applies the deployment. It is not when the release or deployment was created, and it is not fleet-wide.
3. **Pre-existing files are skipped only while unchanged.** Discovery drops a file when `is_preexisting` holds, meaning the size **and** mtime equal the recorded values (`rule.rs:178-189`; `state.rs:53-57`, `93-95`). Any change to size or mtime (append, rewrite, `touch`, a file still being written at apply time) makes it a candidate. It is then uploaded once stable, because a pre-existing file has no ledger entry to match (`rule.rs:272-292`). The recorded entry is never removed, so later scans keep comparing against it. Unit tests: `unchanged_preexisting_not_promoted` and `changed_preexisting_promoted` (`rule.rs:644-664`).
4. **Not mtime-based.** New paths are candidates regardless of their mtime (`rule.rs:178-189`), so a file copied in with old timestamps is still uploaded.
5. **Retention skips pre-existing files too.** Delete jobs come only from stable files (`src/data_uploads/retention/sink.rs`, `on_stable_file`) or from upload confirmation. Pre-existing unchanged files never become stable, so retention never deletes them. Upload jobs come from the same stream (`src/data_uploads/upload/sink.rs:30-59`).
6. **Same rule in a later deployment keeps the snapshot.** Scanners are keyed by rule ID. When a known ID arrives again, only the deployment is swapped (`scanner.rs:170-171`, `rule.rs:72-74`). The backend deduplicates rules by a digest of name, glob, stability window, upload (collection, bucket, path) and retention (`backend: internal/configs/domain/filerules/hash.go`; find-or-create in `internal/configs/services/file_rules/create.go:461-484`). So identical rule content in a new release has the same ID and keeps its snapshot. Changing any of those fields, including only the name, yields a new ID, and that new rule takes a new snapshot. Tests: `update_rules_refreshes_deployment_carrying_state`, `update_rules_repush_does_not_swallow_new_files` (`scanner.rs:1322-1467`).
7. **A replaced rule drains.** A rule missing from the new deployment stops discovering files. Its already-discovered candidates are still evaluated and delivered. Once none remain, it is pruned (`scanner.rs:206-235`; test `update_rules_keeps_legacy_scanner_until_candidates_drain`, `scanner.rs:1392-1433`). A deployment with no replacement clears the deployed set the same way (`scanner.rs:142-147`, `sync_scan_bridge.rs:90-93`).
8. **Remove then re-add.** After pruning, the rule's state is gone. If the same rule is deployed again, `RuleScanner::new` takes a fresh snapshot, so files created while it was absent are pre-existing. Edge case, not documented: if the rule returns before its leftover candidates drain, the old snapshot is reused.
9. **Survives restarts.** Scanner state, including `preexisting`, is persisted to `scanner.json` in the agent's data directory. It is written on every rule update and scan and restored at startup (`state.rs:140-152`, `scanner.rs:57-104`, `src/disk/layout.rs:44`, `src/app/state.rs:143-168`). Test: `existing_state_file_restores_scanner`. If the file cannot be initialized, the agent runs without persistence (`app/state.rs:153-159`). That is a failure path and is not documented.
10. **Not documented (edge cases):** a matching file that cannot be stat'ed at apply time is not recorded as pre-existing (`rule.rs:155-160`), so it is uploaded once readable. Upload jobs already queued are not cancelled when a rule is removed (the upload module does not filter by deployment), but backend acceptance of such uploads was not verified.

Discrepancies from the Slack claim: (a) modified pre-existing files *are* uploaded; (b) the cutoff is per device *and per rule*, and it carries over to later deployments while the rule is unchanged; (c) retention also skips pre-existing files. The claim's per-device and not-release-creation-time parts are correct.

### Tooling

CI is `.github/workflows/ci.yml` (workflow `CI`). Jobs:

- `changes`
- `lint`: `pnpm install --frozen-lockfile`, `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`
- `audit`: `./scripts/audit.sh`
- `shell-tests`: bats
- `lint-custom-linter` and `test-custom-linter`: skipped unless `tools/lint/**` changes

`./scripts/lint.sh` runs:

- the Go MDX linter in `tools/lint/` (sentence-case headings, no `--`, import and redirect rules)
- ESLint MDX
- CSpell with `cspell.json`
- the OpenAPI checks

GitHub access is via the GitHub MCP tools (no `gh` CLI), repo `mirurobotics/docs`. Commits must be signed (`commit.gpgsign` is true); never pass `--no-gpg-sign`.

**Preflight `CLEAN`** means every CI check run on the pushed branch head SHA has completed with `success`, or `skipped` for the two custom-linter jobs.

## Plan of Work

**Edit 1: `docs/data-uploads/primitives/file-rules/overview.mdx`.** Insert this section after `## Immutability` (after line 61) and before `## Git provenance`, with one blank line on each side:

    ## Deployment

    A file rule takes effect on a device when that device applies a
    [deployment](/primitives/deployments) whose release contains the rule. Each device
    starts the rule at its own moment: when the deployment reaches the `deployed`
    [activity status](/primitives/deployments#activity-status) on that device, not when the
    release or deployment was created. A device that is offline starts the rule once it
    comes back online and applies the deployment.

    ### Existing files

    When a rule takes effect on a device, the Miru Agent records every file that already
    matches the rule's [glob](/data-uploads/primitives/file-rules/rule-definition#glob-patterns),
    along with each file's size and modification time. These pre-existing files are left
    alone: they are not uploaded, and the rule's
    [retention](/data-uploads/primitives/file-rules/rule-definition#retention) never
    deletes them.

    A pre-existing file is picked up only once it changes. If its size or modification time
    no longer matches what the agent recorded, for example because a process appended to
    it or because it was still being written when the rule took effect, the agent treats
    it like a new file. It waits for the file to become stable, then uploads it and applies
    retention.

    The agent compares files against this record, not against timestamps. A new file that
    appears after the rule takes effect is uploaded even if its modification time is older
    than the deployment, such as a file copied onto the device with its original
    timestamps. The record is stored on the device and survives agent restarts.

    <Tip>
      To upload files that were already on a device when the rule took effect, update their
      modification time after the rule takes effect (for example, with `touch`). The agent
      then treats them as changed, uploads them, and applies the rule's retention to them.
    </Tip>

    ### Later deployments

    File rules are identified by their [digest](#param-digest). Deploying a new release that
    contains the same rule (identical name, source, upload, and retention) continues the
    rule on each device instead of starting it again. The device keeps its original record,
    so files that were already uploaded, or skipped as pre-existing, are not revisited.

    Changing anything in a rule, even only its name, creates a new rule. When the new rule is
    deployed, each device takes a new record for it, so files already on the device are not
    uploaded by the new rule. This includes files that a widened glob now matches for the
    first time. Files that the previous rule was still waiting on to become stable finish
    under the previous rule.

    If a rule is removed from a device and deployed to it again later, the device takes a
    new record, so files created while the rule was removed are not uploaded.

Before relying on the `#param-digest` anchor, confirm it. Mintlify generates `param-<path>` IDs for `<ParamField>`, and `docs/data-uploads/overview.mdx:83` already links `/primitives/releases#param-version`. If `pnpm run validate` or the dev server shows otherwise, replace `[digest](#param-digest)` with plain `digest`.

**Edit 2: `docs/data-uploads/overview.mdx`, `## File detection`.** After the paragraph ending "zero false positives." (line 35), insert one blank line and this paragraph:

    The Agent only uploads files that are new or have changed since the rule took effect on
    the device. Files that already match a rule when it is deployed are skipped. See
    [deployment](/data-uploads/primitives/file-rules/overview#deployment) for details.

**Edit 3: `docs/snippets/file-rules/retention.mdx`.** After line 32 (the paragraph that starts "The retention block doesn't go into effect until a file is considered stable."), insert one blank line and:

    Files that already matched the rule when it was deployed to the device are never
    deleted unless they change afterward. See
    [deployment](/data-uploads/primitives/file-rules/overview#deployment).

Deliberately not changed:

- `docs/data-uploads/primitives/file-rules/rule-definition.mdx` lines 37-39 ("files that appear after a rule is deployed can still match it") remain true.
- `docs/data-uploads/primitives/uploads.mdx` and `docs/primitives/deployments.mdx` are general. A pointer there would duplicate Edit 2.
- `docs/docs.json` has no new page. `docs/changelog/*` has no entry. `docs/references/**` is generated.

## Concrete Steps

All docs commands run from the docs repo root (`/home/user/docs`, or `git rev-parse --show-toplevel`).

### Milestone 1: behavior check, edits, local checks, commit

1. Confirm the start state:

        git branch --show-current   # expect: claude/adoring-bohr-mdjcop
        git status --short          # expect: only plans/ changes, if any

2. Test step (agent behavior check, read-only). From the agent repo root (`/home/user/agent`), run:

        cargo test -p miru-agent --lib data_uploads::scan

    Expect `test result: ok. 116 passed; 0 failed` (verified 2026-09-22 at agent `f372c9c`; the count may grow on newer commits, but there must be 0 failures). These tests must be among the passing ones:

    - `unchanged_preexisting_not_promoted`
    - `changed_preexisting_promoted`
    - `update_rules_refreshes_deployment_carrying_state`
    - `update_rules_keeps_legacy_scanner_until_candidates_drain`
    - `update_rules_repush_does_not_swallow_new_files`
    - `existing_state_file_restores_scanner`

    Confirm them with `... -- --list | grep -E '<name>'` if needed. If a test is missing or fails, re-read the cited code and correct the docs text to match. Record it in Surprises & Discoveries. `git -C /home/user/agent status --short` must stay empty. Nothing from this step is committed.

3. Apply Edits 1-3 from Plan of Work.

4. Content checks:

        grep -n '^## Immutability\|^## Deployment\|^### Existing files\|^### Later deployments\|^## Git provenance' docs/data-uploads/primitives/file-rules/overview.mdx
        # expect five lines in this order: Immutability, Deployment, Existing files, Later deployments, Git provenance
        grep -c 'file-rules/overview#deployment' docs/data-uploads/overview.mdx docs/snippets/file-rules/retention.mdx
        # expect 1 for each file

5. Local checks, a fast mirror of the CI `lint` job:

        pnpm install --frozen-lockfile
        pnpm run test:lint     # expect exit 0
        ./scripts/lint.sh      # expect last line: All documentation lint checks passed.
        pnpm run validate      # expect exit 0 (mint validate)

    If CSpell flags a word, add only that word to `words` in `cspell.json` and rerun. If the registry is unreachable, build and run the custom linter on the three files and rely on CI for the rest:

        cd tools/lint && go build -o lint . && cd ../..
        tools/lint/lint docs/data-uploads/primitives/file-rules/overview.mdx docs/data-uploads/overview.mdx docs/snippets/file-rules/retention.mdx

    Optionally run `pnpm run dev` and open `http://localhost:3000/data-uploads/primitives/file-rules/overview#deployment` to check rendering and the `#param-digest` link.

6. Commit one signed commit, with the session's required trailers:

        git add docs/data-uploads/primitives/file-rules/overview.mdx docs/data-uploads/overview.mdx docs/snippets/file-rules/retention.mdx plans/
        git add cspell.json   # only if changed
        git commit -m "docs(file-rules): document deployment behavior for existing files"

### Milestone 2: push, draft PR, preflight to CLEAN

1. `git push -u origin claude/adoring-bohr-mdjcop`.
2. Load the GitHub MCP tools with `ToolSearch` (`select:mcp__github__create_pull_request,mcp__github__pull_request_read,mcp__github__actions_list,mcp__github__get_job_logs,mcp__github__update_pull_request`). Open a **draft** PR with owner `mirurobotics`, repo `docs`, head `claude/adoring-bohr-mdjcop`, base `main`, and title `docs(file-rules): document deployment behavior for existing files`. The body summarizes the three edits and lists the verified behavior, including the refinements of the original claim. It ends with the session's PR attribution lines.
3. Poll CI for `git rev-parse HEAD`. For any failed job, read the job logs, fix the cause locally, rerun the matching step 5 check, then commit (signed, new commit) and push. Never amend or force-push pushed commits. If `audit` fails on an advisory unrelated to this diff, do not widen the diff: record it and leave the PR in draft.
4. When preflight reports `CLEAN`, set `draft: false`.

## Validation and Acceptance

1. `/data-uploads/primitives/file-rules/overview` shows `## Deployment`, with `### Existing files` (including the Tip) and `### Later deployments`, between `## Immutability` and `## Git provenance`. A reader can answer from it:
    - Will files already on my robot be uploaded when I deploy this rule? No, unless they change afterward.
    - Is the cutoff per device? Yes: when that device applies the deployment.
    - Does a new release with the same rule reset it? No.
    - Does changing or widening the rule backfill existing files? No.
    - Does retention delete pre-existing files? No.
    - How do I upload a backlog? `touch` the files after the rule takes effect.
2. `/data-uploads/overview` `## File detection` and the `retention` field on `/data-uploads/define-file-rules` and `/data-uploads/primitives/file-rules/rule-definition` each link to `#deployment`.
3. The Milestone 1 step 2 agent test run reports 0 failures, with the six named tests passing.
4. The step 4 greps return the stated order and counts.
5. `git diff --stat main...HEAD` lists only the three docs files, the plan file, and optionally `cspell.json`.
6. Locally, `pnpm run test:lint`, `./scripts/lint.sh` and `pnpm run validate` exit 0 (or the custom-linter fallback passes).
7. **Preflight reports `CLEAN`**: every CI check run on the pushed branch head SHA passed (`changes`, `lint`, `audit`, `shell-tests`; the custom-linter jobs are skipped). This must hold before the PR leaves draft and before the task is reported complete.

## Idempotence and Recovery

- The agent test run is read-only and repeatable. It writes only to the agent's gitignored `target/`.
- Edits are plain insertions. Check the step 4 greps before reapplying, and skip any edit that is already present.
- Before commit, undo an edit with `git checkout -- <file>`. After commit, use `git revert <sha>`. Fix CI failures forward with new commits. The only allowed force-push is `git push --force-with-lease` after rebasing onto `origin/main` to resolve a conflict.
- The `agent` and `backend` repos are never written to.

# Document the supported glob syntax for file rule `source.glob`

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`mirurobotics/docs`, this repo) | read-write | Expand the `glob` field in `docs/snippets/file-rules/sources.mdx`, clarify one sentence in `docs/data-uploads/define-file-rules.mdx`, add a `### Recursive globs` subsection to `docs/developers/agent/filesys-access.mdx`. Touch `cspell.json` only if CSpell flags a new word. |
| `agent/` (`mirurobotics/agent`) | read-only | Source of truth for matching: `agent/src/filesys/files.rs`, function `glob`. Do not modify. |
| `backend/` (`mirurobotics/backend`) | read-only | Source of truth for release-time checks: `internal/configs/domain/filerules/spec.go`, function `validateSrcGlob`. Do not modify. |

This plan lives in `docs/plans/` because every write is a documentation edit in this repo. Work happens on the already-checked-out branch `claude/beautiful-pascal-ikbvlx` (base `main`). Do not create another branch.

## Purpose / Big Picture

A customer asked whether a file rule's `source.glob` supports `**` for deep recursion. The docs show only `/var/log/robot/*.log` and the release-time checks, so there is no way to answer from the docs. After this change, the `glob` field (rendered on `/data-uploads/define-file-rules` and `/data-uploads/primitives/file-rules`) has a syntax table, a few example patterns, the common mistakes, and a warning that invalid patterns are accepted at release time but fail on the device. `/developers/agent/filesys-access` explains the directory permissions a `**` pattern needs. A reader can then answer: yes, `**` matches zero or more directories, but it must be a whole path segment (`/var/log/robot/**/*.log`), and `/var/log/robot/**` alone matches no files.

## Progress

- [x] Milestone 1: run the glob behavior check (test step), edit the three pages, pass the local checks, commit.
    - [x] (2026-09-22) Glob behavior check run with `glob =0.3.4` (from the local cargo cache, `--offline`); output matched the expected transcript line for line.
    - [x] (2026-09-22) Edits 1-3 applied as written in Plan of Work; step 4 greps return 1 / 0 / 1 and the heading order Deleting local files (142), Recursive globs (156), Testing access (177).
    - [x] (2026-09-22) `pnpm install --frozen-lockfile`, `pnpm run test:lint`, `./scripts/lint.sh` ("All documentation lint checks passed.") and `pnpm run validate` ("build validation passed") all exit 0. CSpell flagged nothing, so `cspell.json` is unchanged.
- [ ] Milestone 2: push, open a draft PR, get preflight to `CLEAN`, mark the PR ready.
    - [x] (2026-09-22) Pushed `claude/beautiful-pascal-ikbvlx` (already up to date with `origin/main`, no rebase needed) and opened draft PR [#196](https://github.com/mirurobotics/docs/pull/196).
    - [x] (2026-09-22) Preflight `CLEAN` on `7412493` in CI round 1 of 3 (workflow `CI` run 519): `changes`, `lint`, `audit`, and `shell-tests` succeeded, and `lint-custom-linter` and `test-custom-linter` were skipped. No CI fixes were needed. The commit that records this touches only this plan file.
    - [ ] Mark the PR ready and move this plan to `plans/completed/` (left to the orchestrator).

## Surprises & Discoveries

- (2026-09-22) The glob check program prints `=> ` with a trailing space on the two lines that match nothing (`**` and `{app,notes}.*`); the expected transcript shows them trimmed. After stripping trailing whitespace the output is identical, so this is not a behavior difference and no docs text changed.
- (2026-09-22) The agent source lives at `agent/agent/src/filesys/files.rs` in the local checkout (the Cargo workspace crate is nested one level down), not `agent/src/filesys/files.rs`. `Cargo.lock` pins `glob` 0.3.4, matching the version the check used.
- (2026-09-22) Two check runs besides the CI workflow report on the PR head: `[code]smith` and `Mintlify Deployment`. Both are third-party app checks, not jobs in `.github/workflows/ci.yml`, and both complete as `skipped` on this draft PR.
- (2026-09-22) The release-time Warning matches the agent code: `RuleScanner::new` (`agent/agent/src/data_uploads/scan/rule.rs`) calls `files::glob(...)?`, so an invalid pattern fails the scanner's rule-apply for the whole deployment (`agent/agent/src/data_uploads/scan/scanner.rs`).

## Decision Log

- (2026-09-22) Ran the three edits verbatim from Plan of Work, since the glob check confirmed every stated behavior. No wording beyond the plan was added.
- (2026-09-22) Per the orchestrator's instructions for this run, the PR stays in draft and the plan stays in `plans/active/` when this implementation pass finishes; the orchestrator marks it ready and archives the plan. This overrides Milestone 2 step 4.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

This repo is the Mintlify documentation site for Miru. Pages are MDX (Markdown plus JSX components) under `docs/`, and site navigation is in `docs/docs.json`. A **file rule** is a YAML object shipped in a release that tells the Miru Agent (the on-device binary) which files to upload and/or delete. Its `source.glob` field is the pattern that selects those files.

Files involved:

- `docs/snippets/file-rules/sources.mdx`: a shared snippet holding the `<ParamField path="glob">` and `<ParamField path="stability_window_secs">` blocks. It is imported as `<Sources />` by two pages: `docs/data-uploads/define-file-rules.mdx` (under `### Source`) and `docs/data-uploads/primitives/file-rules.mdx` (under `## Sources`). Those pages use different heading levels, so do not add headings inside the snippet. `docs/snippets/file-rules/upload.mdx` already puts a Markdown table inside a `<ParamField>` (the `path` variables table). Follow that pattern.
- `docs/data-uploads/define-file-rules.mdx` lines 26-29 describe the example rule as uploading "all log files from the `/var/log/robot` directory". With `*.log` that is only the top level, so the wording should say so.
- `docs/data-uploads/primitives/file-rules.mdx` gets the new content through `<Sources />`. Its only glob example (`"/var/log/robot/*.log"`, File format section) is already correct. No direct edit.
- `docs/developers/agent/filesys-access.mdx` `## Data uploads` (line 92) has `### Required permissions` (line 100) and `### Deleting local files` (line 142, table at lines 149-154), then `## Testing access` (line 156). Line 140 ends with "covered next", which points at `### Deleting local files`, so nothing may be inserted between them. The section says nothing about `**` yet.

MDX pitfall: outside inline code, `{` starts a JSX expression and `**` starts bold text. Every pattern (`{a,b}`, `**`, `*`, `[!a-z]`, `\*`) must be inside backticks.

Verified matching behavior (the facts the new text states; do not add claims beyond these). The agent's `glob` function in `agent/src/filesys/files.rs` calls the Rust `glob` crate 0.3.4 `glob::glob(pattern)` with default options. It skips entries that error mid-walk, such as unreadable directories, and keeps only regular files (`Path::is_file`). The pattern is re-evaluated on every scan.

- Supported: `*` matches any characters within one path segment and never crosses `/`. `**` matches zero or more directories and must be an entire `/`-delimited segment: `/logs/**.log`, `/logs/a**/*.log` and `/logs/***/*.log` are pattern errors. `?` matches one character. `[abc]` and `[a-z]` are classes, and `[!...]` negates a class.
- Not supported: brace expansion (braces are literal), backslash escapes (backslash is literal), extglob patterns such as `!(x)`, and more than one pattern per rule. To match a literal `*`, `?` or `[`, wrap it in brackets: `[*]`, `[?]`, `[[]`.
- A trailing `**` (`/logs/**`) yields only directories, so it matches no files. Use `/logs/**/*`. Because `**` matches zero levels, `/logs/**/*.log` includes `/logs/x.log`.
- `*` and `**` match hidden (dot) files and directories. Matching is case-sensitive. Directories never match. Symlinks to files and to directories are followed.
- Release-time validation (`validateSrcGlob`) checks only that the pattern is non-empty, absolute, at most 1024 bytes, free of control characters, and free of `..` and empty segments. The existing criteria list documents these checks. Pattern syntax is not checked, so `miru release create` accepts `/var/log/**.log`, but the pattern fails on the device and that release's file rules do not run. This is a known issue being fixed separately.
- Practical: a broad literal prefix (for example `/**/*.log`) walks large parts of the file system on every scan.

Tooling. CI is `.github/workflows/ci.yml` (workflow `CI`, runs on `pull_request`). Its jobs:

- `changes` is a path filter.
- `lint` runs `pnpm install --frozen-lockfile`, then `pnpm run test:lint` (lint smoke tests), then `./scripts/lint.sh`, then `pnpm run validate` (`mint validate`, which fails on MDX parse errors). `./scripts/lint.sh` runs the custom Go MDX linter in `tools/lint/` (sentence-case headings, no `--`, import rules, redirects), ESLint MDX, CSpell with `cspell.json`, and the OpenAPI checks.
- `audit` runs `./scripts/audit.sh` (`pnpm audit`).
- `shell-tests` runs bats on `pub/scripts/agent/check-miru-access_test.bats`.
- `lint-custom-linter` and `test-custom-linter` are skipped on PRs unless `tools/lint/**` changes.

There is no link-check job. `.github/workflows/codeql-analysis.yml` does not run on PRs.

GitHub: there is no `gh` CLI. Use the GitHub MCP tools, loaded with `ToolSearch` (`select:mcp__github__create_pull_request,mcp__github__pull_request_read,mcp__github__actions_list,mcp__github__get_job_logs,mcp__github__update_pull_request`). The repo is `mirurobotics/docs`. `commit.gpgsign` is `true` here and the repo requires signed commits, so never pass `--no-gpg-sign`.

**Preflight `CLEAN`** means every CI check run on the pushed branch head SHA has completed with `success` (or `skipped`, for the two custom-linter jobs).

## Plan of Work

**Edit 1: `docs/snippets/file-rules/sources.mdx`.** Replace the whole `<ParamField path="glob" type="string" required>` block (lines 4-17) with the block below. The criteria list stays word-for-word. The single `Example:` line becomes the examples table. The existing file-access `<Warning>` is extended and moved last. Leave the intro (lines 1-2) and the `stability_window_secs` ParamField untouched.

    <ParamField path="glob" type="string" required>
      An absolute glob pattern that selects the files this rule manages. The Miru Agent
      matches it on the device and re-evaluates it on every scan, so matching files that
      appear later are picked up.

      Must satisfy the following criteria:
      - Absolute — it starts with `/`
      - At most 1024 bytes, with no control characters
      - No `..` segments, and no empty segments (`//` or a trailing `/`)

      The following syntax is supported:

      | Syntax | Matches |
      | --- | --- |
      | `*` | Any characters within one path segment. Never crosses a `/`. |
      | `**` | Zero or more directories, at any depth. Must be a whole path segment on its own, as in `/var/log/**/*.log` — `**.log`, `run**`, and `***` are invalid. |
      | `?` | Exactly one character. |
      | `[abc]`, `[a-z]` | One character from the set or range. |
      | `[!abc]`, `[!a-z]` | One character not in the set or range. |

      Brace expansion (`{a,b}`), backslash escapes (`\*`), and extended patterns such as
      `!(x)` are not supported — braces and backslashes match literally. To match a literal
      `*`, `?`, or `[`, wrap it in brackets: `[*]`, `[?]`, `[[]`. Each rule has exactly one
      pattern; to match files that no single pattern covers, such as `.log` and `.txt`
      files, define a rule for each.

      Examples:

      | Pattern | Matches |
      | --- | --- |
      | `/var/log/robot/*.log` | `.log` files directly inside `/var/log/robot` |
      | `/var/log/robot/*/*.log` | `.log` files exactly one directory below `/var/log/robot` |
      | `/var/log/robot/**/*.log` | `.log` files in `/var/log/robot` and all of its subdirectories |
      | `/var/log/robot/**/*` | Every file in `/var/log/robot` and all of its subdirectories |
      | `/var/log/robot/run-[0-9].log` | `run-0.log` through `run-9.log` in `/var/log/robot` |

      Keep in mind:
      - A trailing `**` matches only directories, so `/var/log/robot/**` matches no files.
        Use `/var/log/robot/**/*` instead.
      - Only regular files match; a directory never matches, even if its name does.
        Symbolic links to files and directories are followed.
      - Matching is case-sensitive (`*.log` does not match `APP.LOG`), and `*` and `**`
        match hidden files and directories (names starting with `.`).
      - Keep the part of the pattern before the first wildcard as specific as possible. A
        broad pattern such as `/**/*.log` walks large parts of the file system on every
        scan.

      <Warning>
        Pattern syntax is not checked when a release is created. `miru release create`
        accepts an invalid pattern such as `/var/log/**.log`, but the pattern fails on the
        device and the release's file rules do not run. Before releasing, check that every
        `**` is a whole path segment on its own.
      </Warning>

      <Warning>
        You must grant the Miru Agent read access to the files you want to upload and to
        the directories that contain them, including every directory a `**` descends into.
        Files in directories the agent can't read are silently skipped. Visit the
        [file system access](/developers/agent/filesys-access#data-uploads) section for
        more details.
      </Warning>
    </ParamField>

**Edit 2: `docs/data-uploads/define-file-rules.mdx`, lines 26-29.** Change "uploads all log files from the `/var/log/robot` directory to the" to "uploads the `.log` files directly inside the `/var/log/robot` directory (not its subdirectories) to the". Rewrap the paragraph to about 88 columns and leave everything else unchanged.

**Edit 3: `docs/developers/agent/filesys-access.mdx`.** Insert the subsection below after the `### Deleting local files` table (after line 154) and before `## Testing access`, with one blank line on each side:

    ### Recursive globs

    When a rule's `source.glob` contains `**`, the agent lists every directory the `**`
    descends into, so the `miru` user needs read (`r`) and execute (`x`) on each of those
    directories, not only the top one. For the glob `/var/log/robot/**/*.log` matching the
    file `/var/log/robot/nav/2026/app.log`, the required permissions are:

    | Path | Permissions |
    |------|----------------------|
    | `/var/log/robot/nav/2026/app.log` | read (`r`) |
    | `/var/log/robot/nav/2026` | read (`r`), execute (`x`) |
    | `/var/log/robot/nav` | read (`r`), execute (`x`) |
    | `/var/log/robot` | read (`r`), execute (`x`) |
    | `/var/log` | execute (`x`) |
    | `/var` | execute (`x`) |

    Directories the `miru` user can't read are silently skipped: files beneath them are
    never matched. If the rule also has a `retention` block, add write (`w`) on each
    directory that directly contains a matching file, as described in
    [deleting local files](#deleting-local-files).

Deliberately not changed:

- `docs/data-uploads/primitives/file-rules.mdx` inherits Edit 1 through `<Sources />`, and its example is already consistent.
- `docs/changelog/*.mdx` gets no entry. The changelogs record product, agent, CLI and API releases, not docs clarifications. `plans/completed/20260818-miru-version-command-docs.md` ("that changelog tracks CLI binary releases … not documentation additions") and `plans/completed/20260726-opaque-schema-language.md` ("NO changelog entry") set the precedent.
- `docs/references/**` is generated and is not touched.

## Concrete Steps

All commands run from the docs repo root (`git rev-parse --show-toplevel`) unless stated.

### Milestone 1: glob check, edits, local checks, commit

1. Confirm the starting state:

        git branch --show-current   # expect: claude/beautiful-pascal-ikbvlx
        git status --short          # expect: clean apart from plans/

2. Test step (glob behavior check). Run this before editing. It confirms that each example and behavior the new text states holds for the exact `glob` version the agent uses, with the agent's filtering applied. Work in a scratch directory outside the repo:

        dir=$(mktemp -d) && mkdir -p "$dir/src" && cd "$dir"

    Write `Cargo.toml`:

        [package]
        name = "globcheck"
        version = "0.1.0"
        edition = "2021"

        [dependencies]
        glob = "=0.3.4"

    Write `src/main.rs`:

        use std::{env, fs, os::unix::fs::symlink};

        fn main() {
            let base = env::temp_dir().join("globcheck");
            let _ = fs::remove_dir_all(&base);
            let r = base.join("robot");
            for d in ["nav/2026", ".hidden", "dir.log"] { fs::create_dir_all(r.join(d)).unwrap(); }
            fs::create_dir_all(base.join("other")).unwrap();
            for f in ["app.log", ".dot.log", "APP.LOG", "run-1.log", "run-a.log", "notes.txt",
                      "{a,b}.log", "*.log", "nav/nav.log", "nav/2026/deep.log", ".hidden/h.log"] {
                fs::write(r.join(f), "x").unwrap();
            }
            fs::write(base.join("other/o.log"), "x").unwrap();
            symlink(base.join("other"), r.join("link")).unwrap();
            let root = format!("{}/", r.display());
            for p in ["*.log", "*/*.log", "**/*.log", "**/*", "**", "run-[0-9].log",
                      "run-[!0-9].log", "run-?.log", "*.LOG", "{app,notes}.*", "[{]a,b}.log",
                      "[*].log", "**.log", "nav**/*.log", "***/*.log"] {
                // Same filtering as the agent: skip walk errors, keep regular files only.
                let out = glob::glob(&format!("{root}{p}")).map(|it| {
                    let mut v: Vec<String> = it.filter_map(Result::ok).filter(|x| x.is_file())
                        .map(|x| x.display().to_string().replace(&root, "")).collect();
                    v.sort();
                    v.join(" ")
                });
                match out {
                    Ok(v) => println!("{p:<15} => {v}"),
                    Err(e) => println!("{p:<15} => ERROR: {}", e.msg),
                }
            }
        }

    Run `cargo run --offline -q` (drop `--offline` if glob 0.3.4 is not in the local cargo cache). Expected output (verified 2026-09-22). The `**` and `{app,notes}.*` lines end right after `=>` because they match nothing:

        *.log           => *.log .dot.log app.log run-1.log run-a.log {a,b}.log
        */*.log         => .hidden/h.log link/o.log nav/nav.log
        **/*.log        => *.log .dot.log .hidden/h.log app.log link/o.log nav/2026/deep.log nav/nav.log run-1.log run-a.log {a,b}.log
        **/*            => *.log .dot.log .hidden/h.log APP.LOG app.log link/o.log nav/2026/deep.log nav/nav.log notes.txt run-1.log run-a.log {a,b}.log
        **              =>
        run-[0-9].log   => run-1.log
        run-[!0-9].log  => run-a.log
        run-?.log       => run-1.log run-a.log
        *.LOG           => APP.LOG
        {app,notes}.*   =>
        [{]a,b}.log     => {a,b}.log
        [*].log         => *.log
        **.log          => ERROR: recursive wildcards must form a single path component
        nav**/*.log     => ERROR: recursive wildcards must form a single path component
        ***/*.log       => ERROR: wildcards are either regular `*` or recursive `**`

    How the output maps to the new text:

    - `*` stays at one level.
    - `**` includes the top level.
    - A trailing `**` matches nothing.
    - Directory `dir.log` never matches.
    - Dot files match.
    - Case matters.
    - Braces are literal.
    - `[*]` escapes `*`.
    - `**` must be a whole segment.
    - Symlinked `link/` is followed.

    If any line differs, stop and record it in Surprises & Discoveries, and change the docs text only to match the observed behavior. Afterwards, delete `$dir` and `${TMPDIR:-/tmp}/globcheck`, then `cd` back to the repo root. Nothing from this step is committed.

3. Apply Edits 1-3 from Plan of Work.

4. Content checks, all from the repo root:

        grep -c 'is not checked when a release is created' docs/snippets/file-rules/sources.mdx   # expect: 1
        grep -c 'Example: ' docs/snippets/file-rules/sources.mdx                                  # expect: 0
        grep -c 'directly inside the `/var/log/robot` directory' docs/data-uploads/define-file-rules.mdx  # expect: 1
        grep -n '^### Recursive globs\|^### Deleting local files\|^## Testing access' docs/developers/agent/filesys-access.mdx
        # expect three lines in the order: Deleting local files, Recursive globs, Testing access

5. Local checks (fast mirror of the CI `lint` job):

        pnpm install --frozen-lockfile
        pnpm run test:lint     # expect exit 0
        ./scripts/lint.sh      # expect last line: All documentation lint checks passed.
        pnpm run validate      # expect exit 0 (mint validate; catches MDX parse errors)

    If CSpell flags a word in the new text, add only that word to the `words` array in `cspell.json` and rerun. If `pnpm install` cannot reach the registry, run the custom MDX linter alone (`cd tools/lint && go build -o lint . && cd ../.. && tools/lint/lint docs/snippets/file-rules/sources.mdx docs/data-uploads/define-file-rules.mdx docs/developers/agent/filesys-access.mdx`, expect no output and exit 0), then rely on CI for the rest. Optionally run `pnpm run dev` and open `http://localhost:3000/data-uploads/define-file-rules` and `/developers/agent/filesys-access` to check the tables render inside the ParamField and that no stray bold or italics appear.

6. Commit (one signed commit for the milestone; append the session's required commit trailers):

        git add docs/snippets/file-rules/sources.mdx docs/data-uploads/define-file-rules.mdx \
                docs/developers/agent/filesys-access.mdx plans/
        git add cspell.json   # only if changed
        git commit -m "docs(file-rules): document source.glob pattern syntax"

### Milestone 2: push, draft PR, preflight to CLEAN

1. `git push -u origin claude/beautiful-pascal-ikbvlx`.
2. Load the GitHub MCP tools with `ToolSearch` (see Context). Call `mcp__github__create_pull_request` with owner `mirurobotics`, repo `docs`, head `claude/beautiful-pascal-ikbvlx`, base `main`, `draft: true`, and title `docs(file-rules): document source.glob pattern syntax`. The body summarizes the three page changes, notes that no changelog entry is needed (docs clarification), and ends with the session's required PR attribution lines.
3. Poll CI for the head SHA (`git rev-parse HEAD`) with `mcp__github__pull_request_read` (check runs / status) or `mcp__github__actions_list` (workflow runs for the branch). For any failed job, read `mcp__github__get_job_logs`. Fix the cause locally, rerun the matching local check from Milestone 1 step 5, and commit a new signed commit (`fix: …` or `docs(file-rules): …`). Push and poll again. Never amend or force-push a pushed commit. If `audit` fails on an advisory that has nothing to do with this diff, do not widen the diff. Record it in Surprises & Discoveries and report it, leaving the PR in draft.
4. When preflight reports `CLEAN`, call `mcp__github__update_pull_request` with `draft: false`.

## Validation and Acceptance

1. On both `/data-uploads/define-file-rules` and `/data-uploads/primitives/file-rules`, the `glob` field shows the following, in order:
    - the unchanged criteria list
    - a 5-row syntax table (`*`, `**`, `?`, classes, negated classes)
    - the not-supported paragraph
    - a 5-row examples table
    - four "Keep in mind" bullets
    - two warnings: syntax not checked at release creation, then file access

   A reader can answer the customer's question from it: `**` is supported, matches zero or more directories, must be its own segment, and `/var/log/robot/**` matches no files.
2. `/developers/agent/filesys-access` has `### Recursive globs` between `### Deleting local files` and `## Testing access`. It has a 6-row permissions table for `/var/log/robot/**/*.log` and states that unreadable directories are silently skipped.
3. The glob behavior check (Milestone 1 step 2) prints exactly the expected transcript.
4. The step 4 greps return the stated counts and order.
5. `git diff --stat main...HEAD` lists only `docs/snippets/file-rules/sources.mdx`, `docs/data-uploads/define-file-rules.mdx`, `docs/developers/agent/filesys-access.mdx`, the plan file, and (if needed) a words-only `cspell.json` change. No `docs/changelog/`, `docs/docs.json` or `docs/references/` changes.
6. Locally, `pnpm run test:lint`, `./scripts/lint.sh` and `pnpm run validate` all exit 0, or the custom-linter fallback passes if the registry is unreachable.
7. **Preflight reports `CLEAN`**: every CI check run on the pushed branch head SHA has passed (`changes`, `lint`, `audit`, `shell-tests`, with the two custom-linter jobs skipped), as read through the GitHub MCP tools. This must hold before the PR leaves draft and before the task is reported complete.

## Idempotence and Recovery

- The glob check runs in a throwaway directory and recreates its fixture tree on every run, so it is safe to repeat.
- The edits are plain text replacements. Before reapplying one, check the step 4 greps: if a count is already met, that edit is done, so skip it rather than duplicating it.
- `./scripts/lint.sh` only writes the gitignored `tools/lint/lint` binary.
- Before committing, undo a bad edit with `git checkout -- <file>`. After committing, use `git revert <sha>`.
- If CI fails after a push, fix forward with a new commit. The only allowed force-push is `git push --force-with-lease` after rebasing onto `origin/main` to resolve a conflict.
- The `agent` and `backend` repos are never written to.

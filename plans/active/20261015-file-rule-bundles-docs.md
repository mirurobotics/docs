# Document file rule upload bundles and add the CLI v0.13.0 changelog entry

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | New page `docs/data-uploads/concepts/file-rules/bundles.mdx`, nav entry, cross-page updates to file rule and upload docs, `cspell.json`, and a `docs/changelog/cli.mdx` entry. |
| `backend/` | read-only | Source of truth for bundle semantics (merged PRs mirurobotics/backend#911, #912, #916). |
| `openapi/` | read-only | Source of truth for the `bundle` and `bundle_id` fields (mirurobotics/openapi#309). |
| `cli-private/` | read-only | Source of truth for the CLI changes (mirurobotics/cli-private#234, commit `815f827`; #232, commit `681b642`). |

This plan lives in `docs/plans/` because every edit is a documentation page in the docs repo. Work on branch `docs/file-rule-bundles` (created from `origin/main`); the PR targets `main`.

## Purpose / Big Picture

Miru can now group related uploads, such as every file of one robot recording, into a **bundle**. A file rule's `source.glob` can name parts of a path with captures (for example `{recording}`), `upload.path` can use those captures and a new `{bundle_id}` variable, and the new `upload.bundle` field lists which captures decide bundle membership. After this change a reader finds a "Bundles" page under Data Uploads > Concepts > File rules that explains captures, `upload.bundle`, and bundle IDs with worked examples; the rule definition and upload pages show the new syntax, field, and variables; the uploads page lists the `bundle_id` property and the 2048-byte source path limit; and the CLI changelog's top entry is `v0.13.0` (October 15, 2026).

## Progress

- [x] Milestone 1: bundles page, nav entry, cspell word.
- [x] Milestone 2: cross-page updates (rule definition, snippets, uploads, define file rules, CLI flags anchor).
- [ ] Milestone 3: CLI changelog entry.
- [ ] Milestone 4: push, draft PR, `$preflight` reports `CLEAN`, PR marked ready, `ai-review` findings addressed, plan moved to `plans/completed/`, final `$preflight` `CLEAN`.

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

(Add entries as work proceeds.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The docs repo is a Mintlify site. Pages are `.mdx` files under `docs/`; a page's URL is its path without `docs/` and without the extension (`docs/data-uploads/concepts/uploads.mdx` is `/data-uploads/concepts/uploads`). Internal links are root-relative without `.mdx`. Headings get anchors from their text (`## Glob captures` becomes `#glob-captures`), and a `<ParamField path="bundle">` gets the anchor `#param-bundle` (underscores become hyphens, so `bundle_id` is `#param-bundle-id`). Navigation lives in `docs/docs.json`. Built-in components (`<ParamField>`, `<Info>`, `<Warning>`, `<Tabs>`, `<Tab>`) need no import. Files under `docs/snippets/` are partials imported into several pages, so they contain no headings.

Domain terms: a **file rule** is a YAML file, pushed with a release by `miru release create`, that selects device files with `source.glob` and optionally uploads them as described by its `upload` block (see `docs/data-uploads/concepts/file-rules/rule-definition.mdx`). An **upload collection** is the named group of uploads that `upload.collection_slug` assigns each upload to (see `docs/data-uploads/concepts/upload-collections.mdx`); the Miru Agent on each device matches files and creates the uploads.

MDX hazard: a bare `{` in prose or a table starts a JSX expression and breaks the build. Every `{name}`, `{bundle_id}`, `{a,b}`, and `**` in prose must sit inside backticks or a fenced code block. In YAML examples, quote any value that starts with `{` (for example `path: "{device_id}/..."`), otherwise YAML reads it as a mapping; quote globs too for consistency.

Lint rules that CI enforces (sources in `tools/lint/linter/`, run by `scripts/lint.sh`):

- `headingcase`: the frontmatter `title` and every heading must be strict sentence case. Inline code is masked. Allowlisted tokens include `API`, `CLI`, `JSON`, `URL`, `Miru`, `Agent`, `Git`, and version tags like `v0.13.0`. `ID`, `IDs`, and `YAML` are **not** allowlisted, so "Bundle ID" fails; use "Bundle identifiers".
- `nodoubledash`: no `--` in prose outside code.
- import rules: imports resolve, are used, are sorted case-insensitively by path, and form one block with no blank lines.
- CSpell checks MDX against the single `words` array in `cspell.json` (loosely sorted, uppercase entries first). `mcap` is not present and will be flagged.

Bundle facts (authoritative, from the merged PRs listed in Scope):

1. Captures in `source.glob`. A `{name}` placeholder matches like `*`: any characters within one path segment, including the empty string. Names match `[a-z][a-z0-9_]*`, are unique within the glob, and must not reuse a built-in `upload.path` variable (`device_id`, `device_name`, `file_name`, `upload_id`, `bundle_id`, `year`, `month`, `day`, `hour`, `minute`). A segment holds at most one capture, and the rest of that segment must be constant text (no `*`, `?`, `[`, `\`, or other braces). A capture may not have `**` segments both before and after it. A brace group containing a comma, such as `{a,b}`, is not a capture and still matches literally; any brace group without a comma (for example `{}` or `{Recording}`) is read as a capture and rejected if its name is invalid. Every capture must be used in `upload.bundle` or `upload.path`. Devices see each capture as `*`, so no agent upgrade is needed.
2. `upload.path` can use any capture as a variable (for example `{recording}`) and a new `{bundle_id}` variable, which requires `upload.bundle`. `{upload_id}` is still required in every template.
3. `upload.bundle` is an optional list of 1 to 16 unique capture names. Uploads from the same device and the same upload collection whose captured values match, in the listed order, share one bundle, across rules and over time. Capture names do not affect matching: rules that name captures differently share a bundle if they list them in the same order with the same values. Swapped values (session `03`/run `05` versus session `05`/run `03`) are different bundles. Bundles never span devices or upload collections. Bundles have no name; each has an ID with prefix `upl_bdl_` (example `upl_bdl_123`), exposed through `{bundle_id}` and as the nullable `bundle_id` field on uploads. `upload.bundle` is never sent to devices.
4. An upload's source file path is limited to 2048 bytes. What happens to longer paths is not established; state the limit only.
5. Globs may not contain `.` segments, like `..`; this applies to every file rule (backend#889, openapi#309).
6. CLI (cli-private#234): file rule YAML accepts `upload.bundle` as a YAML list and sends it on push (`miru release create`). The CLI does no local validation; the server validates. cli-private#232 fixed the update notice's changelog link to point at `/changelog/cli`. Other commits since `v0.12.1` are internal. The requester specified `v0.13.0` for this release, but cli-private's latest tag is the pre-release `v0.12.2-beta.1`; confirm the shipped version before merge (see Idempotence and Recovery).

Files to touch (line numbers as of `origin/main` at `b33cfc4`):

- `docs/data-uploads/concepts/file-rules/rule-definition.mdx`: `### Glob patterns` at line 41; `..` segment rule at line 49; Supported syntax table lines 56–62; brace paragraph lines 64–67; read-access `<Warning>` lines 105–109; `## Uploads` line 111.
- `docs/snippets/file-rules/sources.mdx`: `glob` ParamField lines 4–35; `..` segment rule at line 10; "See glob patterns" link lines 33–34.
- `docs/snippets/file-rules/upload.mdx`: `path` ParamField lines 17–42 with variable table lines 27–37 and `{upload_id}` Warning lines 39–41.
- `docs/data-uploads/concepts/uploads.mdx`: imports `NullableBadge` already (line 5); `deployment` ParamField ends at line 77; `file_path` ParamField lines 87–94.
- `docs/data-uploads/define-file-rules.mdx`: deduplication paragraph at line 60.
- `docs/snippets/references/cli/releases/create/flags.mdx`: lines 26 and 32 link to `/data-uploads/concepts/file-rules/rule-definition#file-format`, an anchor that does not exist.
- `docs/docs.json`: "File rules" group, line 160 is `"data-uploads/concepts/file-rules/rule-definition"`.
- `docs/changelog/cli.mdx`: import on line 7, two blank lines, newest entry `# v0.12.1` at line 10. Entry format: `# vX.Y.Z`, blank, `*Month D, YYYY*`, blank, one summary paragraph ("CLI `vX.Y.Z` adds ..."), `## Features` / `## Fixes` with hyphen bullets and no trailing period, then `---` before the next entry.
- `cspell.json`: `words` array.

## Plan of Work

Milestone 1 (bundles page). First move this plan to `plans/active/` (the first command under Milestone commits in Concrete Steps). Create `docs/data-uploads/concepts/file-rules/bundles.mdx` with frontmatter `title: "Bundles"` and `description: "Group related uploads from one recording or run into a bundle"`. No imports are needed. Content, kept in the tight style of `rule-definition.mdx`:

- Intro (no heading): a bundle groups related uploads, such as every file from one recording or run, even when different file rules upload them. A rule opts in by naming parts of its glob with captures and listing them in `upload.bundle`.
- `## Example`: the recording rule below, then two sample device files and their object keys, then one sentence that a file in another recording folder lands in a different bundle with a different `{bundle_id}`.

        name: Recordings

        source:
          glob: "/data/recordings/{recording}/*"

        upload:
          collection_slug: recordings
          bucket: my-uploads-bucket
          path: "{device_id}/{recording}/{bundle_id}/{upload_id}/{file_name}"
          bundle:
            - recording

  Sample table (Device file, Object key): `/data/recordings/run_42/camera.mcap` to `dvc_123/run_42/upl_bdl_123/upl_123/camera.mcap`; `/data/recordings/run_42/metadata.yaml` to `dvc_123/run_42/upl_bdl_123/upl_456/metadata.yaml`.
- `## Captures`: fact 1 in prose plus a two-column table (Pattern, Result) with valid rows `/data/recordings/{recording}/*`, `/data/{session}/run_{run}.mcap`, and invalid rows with the reason: `/data/{recording}_*.mcap` (`*` in a capture's segment), `/data/{session}{run}/*` (two captures in one segment), `/data/**/{recording}/**/*` (`**` before and after), `/data/{Recording}/*` (name must start with a lowercase letter), `/data/{file_name}/*` (reuses a built-in variable). State that a brace group containing a comma, such as `{a,b}`, matches literally while any other brace group is read as a capture, that a capture may match an empty string (`run_{run}.mcap` matches `run_.mcap`), that every capture must appear in `upload.bundle` or `upload.path` (so captures need an `upload` block), and that devices match each capture as `*`, so no agent upgrade is needed. Link the base syntax to `/data-uploads/concepts/file-rules/rule-definition#glob-patterns`.
- `## Bundle membership`: fact 3 except IDs: 1 to 16 unique capture names; same device, same upload collection, same values in the same order share a bundle, across rules and over time; names do not matter; swapped values are different bundles; never spans devices or upload collections; `upload.bundle` is not sent to devices. Link the field to `/data-uploads/concepts/file-rules/rule-definition#param-bundle`. Add an `<Info>` that pushing `upload.bundle` requires CLI `v0.13.0` or later, linking `/changelog/cli#v0-13-0`.
- `## Bundle identifiers`: bundles have no name; each has an ID like `upl_bdl_123`, available as the `{bundle_id}` path variable (requires `upload.bundle`) and as the upload's `bundle_id` field (link `/data-uploads/concepts/uploads#param-bundle-id`), which is null for uploads from rules without `upload.bundle`.
- `## Multiple rules`: one rule per file type, for example because `.yaml` files need a shorter stability window. Rule `Run Recordings` with glob `"/data/runs/run_{run}.mcap"`, bundle `[run]`; rule `Run Metadata` with glob `"/data/runs/run_{id}.yaml"`, `stability_window_secs: 5`, bundle `[id]`; both with `collection_slug: runs`, `bucket: my-uploads-bucket`, and path `"runs/{device_id}/{bundle_id}/{upload_id}/{file_name}"`. Show each rule as a full YAML block in the same shape as the Recordings example (`stability_window_secs` goes under `source`). Explain that `run_42.mcap` and `run_42.yaml` join one bundle because both rules list one capture with value `42` and use the same upload collection, even though the captures are named differently. One closing sentence: the same works for one rule per sensor folder (for example `/data/{session}/camera/*` and `/data/{session}/lidar/*` both bundled by `session`).

In `docs/docs.json`, add `"data-uploads/concepts/file-rules/bundles"` as the third entry of the "File rules" group, after line 160 (add a comma to line 160). In `cspell.json`, insert `"mcap"` between `"keyrings"` and `"mgmt"`; add any other new non-dictionary word the same way.

Milestone 2 (cross-page updates).

- `rule-definition.mdx`: on line 49 change "Contain no `..` segments or empty segments" to "Contain no `.` or `..` segments or empty segments" (fact 5). Add a Supported syntax row `` `{name}` `` | "like `*`, and records the matched text as a named capture; see [Glob captures](#glob-captures)". Reword lines 64–67 (dropping "Braces match literally.") to: "Brace expansion is not supported: a brace group containing a comma, such as `{a,b}`, matches literally. A brace group without a comma, such as `{Recording}` or `{}`, is read as a capture and must be valid (see [Glob captures](#glob-captures)). Backslash escapes (`\*`) and extended patterns such as `!(x)` are not supported." keeping the existing Linux/Windows backslash and `[*]` sentences. After the `<Warning>` (line 109) and before `## Uploads`, add `### Glob captures` with three or four sentences: a `{name}` placeholder matches one path segment like `*` and names that part of the path; captured values can be used as `upload.path` variables and in `upload.bundle`; devices match each capture as `*`; link `/data-uploads/concepts/file-rules/bundles#captures` for the full rules.
- `docs/snippets/file-rules/upload.mdx`: add two variable-table rows: `` `{bundle_id}` `` | `upl_bdl_123` | "The upload's [bundle](/data-uploads/concepts/file-rules/bundles) ID. Requires `bundle`." and `` `{recording}` `` | `run_42` | "The value matched by a [glob capture](/data-uploads/concepts/file-rules/bundles#captures) of the same name." Below the `{upload_id}` Warning keep it as is. After the `path` ParamField, add `<ParamField path="bundle" type="string[]">` describing a list of 1 to 16 unique capture names from `source.glob` that group uploads into a bundle, with a link to `/data-uploads/concepts/file-rules/bundles`, and an `Example:` line showing `[recording]`.
- `docs/snippets/file-rules/sources.mdx`: on line 10 change "No `..` segments, and no empty segments" to "No `.` or `..` segments, and no empty segments" (the rest of the line is unchanged). Before the "See glob patterns" sentence (line 33), add: "A pattern can name parts of a path with captures such as `{recording}`; see [captures](/data-uploads/concepts/file-rules/bundles#captures)."
- `docs/data-uploads/concepts/uploads.mdx`: after the `deployment` ParamField (line 77), add `<ParamField path="bundle_id" type="string">` with `<ImmutableBadge /> <NullableBadge />`, text "The ID of the [bundle](/data-uploads/concepts/file-rules/bundles) the upload belongs to. Null when the file rule that produced the upload has no `upload.bundle`.", and an `Example:` line showing `upl_bdl_123`. In the `file_path` ParamField change "Absolute path of the file on the device." to "Absolute path of the file on the device, at most 2048 bytes."
- `docs/data-uploads/define-file-rules.mdx`: after the line 60 paragraph, add: "To group related files, such as every file from one recording, into a single bundle, see [bundles](/data-uploads/concepts/file-rules/bundles)."
- `docs/snippets/references/cli/releases/create/flags.mdx`: on lines 26 and 32 replace `/data-uploads/concepts/file-rules/rule-definition#file-format` with `/data-uploads/concepts/file-rules/rule-definition`.

Milestone 3 (CLI changelog). In `docs/changelog/cli.mdx`, insert above `# v0.12.1` (line 10), keeping the two blank lines after the import:

    # v0.13.0

    *October 15, 2026*

    CLI `v0.13.0` adds support for grouping related uploads into [bundles](/data-uploads/concepts/file-rules/bundles).

    ## Features

    - Add the `upload.bundle` field to file rule YAML files: a list of glob capture names that groups related uploads into a [bundle](/data-uploads/concepts/file-rules/bundles). `miru release create` sends it to Miru, which validates it

    ## Fixes

    - Fix the update notice's changelog link so it opens the [CLI changelog](/changelog/cli)

    ---

Do not touch other changelogs.

Milestone 4 (publish and validate): push, open a draft PR, run `$preflight` until `CLEAN`, mark ready, address `ai-review` findings, move this plan to `plans/completed/`, and run `$preflight` until `CLEAN` again.

## Concrete Steps

All commands run from the docs repo root on branch `docs/file-rule-bundles`. First make sure that branch is checked out (safe to re-run):

    git fetch origin
    git switch docs/file-rule-bundles 2>/dev/null || git switch -c docs/file-rule-bundles origin/main
    git branch --show-current

Expected: the last command prints `docs/file-rule-bundles`.

Edits are made with an editor. Each milestone ends with a read-only review and a commit. Before each commit, tick that milestone in Progress so the plan update is committed with it. Commit messages follow Conventional Commits and end with a blank line then `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`; signing is handled by the repo's git config.

Read-only review checks (no builds, linters, spellcheck, or tests are run locally):

    python3 -c "import json; json.load(open('docs/docs.json')); json.load(open('cspell.json'))"
    grep -n "file-rules/bundles" docs/docs.json
    grep -rn "^#" docs/data-uploads/concepts/file-rules/bundles.mdx docs/data-uploads/concepts/file-rules/rule-definition.mdx
    grep -nE '(^|[^-])--([^-]|$)' docs/data-uploads/concepts/file-rules/bundles.mdx
    grep -rn "rule-definition#file-format" docs

Then read every changed line containing `{` or `**` and confirm it is inside backticks or a code block; confirm every new link points at an existing page and anchor (`#captures`, `#glob-captures`, `#glob-patterns`, `#param-bundle`, `#param-bundle-id`, `/changelog/cli#v0-13-0`); confirm headings are sentence case with no bare `ID` or `YAML`. Expected: JSON loads silently, the nav grep prints one line, the `rule-definition#file-format` grep prints nothing, and the `--` grep prints nothing outside backticks or code fences (frontmatter `---` lines are excluded by the pattern). The heading grep also prints YAML comment lines inside code fences (for example `# optional block` in `rule-definition.mdx`); those are not headings.

Milestone commits (plain `mv` works whether or not the plan file is tracked; `git add -A plans` stages the move either way):

    mv plans/backlog/20261015-file-rule-bundles-docs.md plans/active/
    git add -A plans docs/data-uploads/concepts/file-rules/bundles.mdx docs/docs.json cspell.json
    git commit -m "docs(data-uploads): add file rule bundles page" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"

    git add docs/data-uploads docs/snippets plans/active/20261015-file-rule-bundles-docs.md
    git commit -m "docs(data-uploads): document glob captures, bundle field, and bundle IDs" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"

    git add docs/changelog/cli.mdx plans/active/20261015-file-rule-bundles-docs.md
    git commit -m "docs(changelog): add CLI v0.13.0 release notes" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"

Milestone 4:

    git push -u origin docs/file-rule-bundles
    gh pr create --draft --base main --title "docs: document file rule bundles and CLI v0.13.0" --body "<body>"

The body is a few lines: what changed (bundles page, capture and `{bundle_id}` variables, `upload.bundle`, upload `bundle_id`, 2048-byte path limit, CLI `v0.13.0` notes), links to mirurobotics/backend#916, mirurobotics/openapi#309, and mirurobotics/cli-private#234, a note that the CLI version must be confirmed before merge, and the final line `🤖 Generated with [Claude Code](https://claude.com/claude-code)`. `gh pr edit` fails in mirurobotics repos (projects-classic); to change the title or body use:

    gh api -X PATCH repos/mirurobotics/docs/pulls/<n> -f body="<body>"

Run `$preflight` until it reports `CLEAN`. Fix any CI failure from the CI job logs with a new commit (no amend; if a rebase is needed, push with `--force-with-lease`). Only then mark the PR ready; this is what triggers `ai-review`, which skips draft PRs:

    gh pr ready <n>
    gh pr checks <n> --watch

Read the review with `gh pr view <n> --comments` and `gh api repos/mirurobotics/docs/pulls/<n>/comments`, and fix each valid finding in a new commit whose message ends with the trailer. Then fill in Outcomes & Retrospective, tick Milestone 4, and move the plan:

    git mv plans/active/20261015-file-rule-bundles-docs.md plans/completed/
    git add -A plans
    git commit -m "docs(plans): complete file rule bundles plan" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
    git push

Run `$preflight` again until it reports `CLEAN` on this final head, and confirm with `gh pr checks <n>` that `ai-review` succeeded on it. Only then report the task complete.

## Validation and Acceptance

This is a docs-only change, so the tests are the CI checks. In `.github/workflows/ci.yml` (workflow `CI`, on `pull_request`), job `lint` must pass on the pushed head: `pnpm run test:lint`, `./scripts/lint.sh` (custom MDX lint including `headingcase`, `nodoubledash`, and import rules; ESLint; CSpell; `mint openapi-check`), and `pnpm run validate` (`mint validate`, which fails on MDX parse errors such as a bare `{`). Job `audit` must pass. `lint-custom-linter` and `test-custom-linter` are skipped because `tools/lint/` is unchanged. `.github/workflows/ai-review.yml` job `ai-review` skips draft PRs, so it first runs when the PR is marked ready; its `ai-review` status is required by branch protection. `$preflight` must report `CLEAN` (CI green on the pushed branch head) before the PR leaves draft, and again on the final head (after `ai-review` fixes and the plan move) before the task is reported complete.

Acceptance, observable on the PR's Mintlify preview (the `Mintlify Deployment` check run, posted by the Mintlify GitHub app rather than a workflow, links a `*.mintlify.site` URL) and, after merge, on the published site:

- Data Uploads > Concepts > File rules lists Overview, Rule definition, Bundles; `/data-uploads/concepts/file-rules/bundles` renders the example, captures table, membership rules, bundle identifiers, and multi-rule example.
- `/data-uploads/concepts/file-rules/rule-definition` shows the `{name}` syntax row, the `### Glob captures` subsection, the `{bundle_id}` and capture rows in the path variable table, and a `bundle` field; its glob rules forbid `.` and `..` segments; `/data-uploads/define-file-rules` shows the same snippet changes and the bundles link.
- `/data-uploads/concepts/uploads` lists `bundle_id` (immutable, nullable, example `upl_bdl_123`) and says `file_path` is at most 2048 bytes.
- `/references/cli/release-create` file rule links open the rule definition page with no broken anchor.
- `/changelog/cli` starts with `v0.13.0`, dated October 15, 2026, with one Features and one Fixes bullet.

## Idempotence and Recovery

All steps are file edits and can be repeated; re-running a grep or review changes nothing. If a commit is wrong, fix it with a follow-up commit. If the CLI release ships under another version or date, find every occurrence with `grep -rn "0\.13\.0\|0-13-0" docs` and, in one new commit whose message ends with the `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>` trailer, update the changelog heading, date line, and summary sentence, plus the bundles page `<Info>` text and its `/changelog/cli#v0-13-0` anchor. Then retitle the PR with `gh api -X PATCH repos/mirurobotics/docs/pulls/<n> -f title="<new title>"`. Already-pushed commit subjects are not rewritten. To abandon, close the PR and delete only the `docs/file-rule-bundles` branch this plan created.

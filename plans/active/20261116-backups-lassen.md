# Document lassen resources in data backups

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (mirurobotics/docs) | read-write | Edit `docs/admin/backups/overview.mdx` and `docs/admin/backups/archive-format.mdx`; possibly add words to `cspell.json`. |
| `backend/` (mirurobotics/backend) | read-only | Source of truth for what a backup writes: `internal/backups/export/manifest.go`, `internal/backups/export/db_tables.go`, `internal/backups/export/export.go`. |

This plan lives in `docs/plans/` because every edit is in the docs repo. Commands below run from the docs repo root (`/home/user/docs` in the authoring checkout) unless stated otherwise. The PR title is `docs(backups): back up lassen resources`, on branch `docs/backups-lassen`.

## Purpose / Big Picture

Backend PRs mirurobotics/backend#1276 and #1277 (merged 2026-11-16) changed nightly data backups. Archives are now written in Platform API version `2026-11-16.lassen` instead of `2026-11-04.everglades`, and they hold three more resource files: file rules, upload collections and uploads. After this change, a customer reading the backups pages sees the version an archive uses, finds every file the archive contains, learns that uploaded files themselves are not copied, and knows that `calibrations.jsonl` has a different shape in lassen archives than in everglades ones.

## Decisions

### Left out

- New headings or a new page: the facts fit as list items, table rows and short paragraphs in the existing sections, as the sibling content does.
- Adding the new types to the manifest example's `resource_counts`: the example is an excerpt (it already omits most types), and the field description does not list types.
- Marking lassen as beta on these pages: no sibling page marks a single Platform API version as beta (`docs/developers/platform-api/versioning.mdx` only says the whole API is in beta), so there is nothing to mirror. Left for design review.
- The `api_keys.jsonl` row: the backend lists API keys with the `scopes` expansion (`db_tables.go`, `getAccessTables`) but the table shows none. That drift predates this change; propose it as a follow-up in the PR body instead of fixing it here.
- `docs/cfg-mgmt/concepts/calibrations.mdx` and `docs/cfg-mgmt/concepts/components.mdx` still describe roles and aspects. Updating them for lassen belongs to the lassen API docs work, not to backups.
- `docs/admin/backups/read-a-backup.mdx`: it names no version or resource list, so it needs no change.

### Open questions

- Should the PR base be PR 2's branch (the sibling PR that publishes `/references/platform-api/2026-11-16` and points the `latest` redirect there) or `main`? Default: when PR 2's branch exists on origin, rebase `docs/backups-lassen` onto it and open this PR against it; otherwise open against `main`, and say in the PR body that it depends on PR 2 and must merge after it.
- Do the words `lassen` and `meca` need adding to `cspell.json`? Neither appears in the repo today. Default: do not add them up front; add exactly the words the CI CSpell step flags (see Milestone 2). If PR 2's branch is the base and already lists a word, do not add it again.

## Progress

- [x] Milestone 1: edit `overview.mdx` and `archive-format.mdx`, commit.
- [x] Source written
- [ ] Milestone 2: push, open draft PR, fix CI findings, preflight reports `CLEAN`.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

- 2026-11-16, design review finding 1 (beta label): not applied. The person decided to call the version lassen with no beta label everywhere; the release docs move to the official lassen API before merging.
- 2026-11-16, design review finding 2: the calibrations paragraph is replaced by one sentence after the resources intro. It names the lassen calibration change as the example, since the task asked for the new shape, and points to the Platform API changelog for the details, so the two pages can't drift. This also drops the `meca` word that failed CSpell.
- 2026-11-16, design review finding 3: the overview's uploads item says only "records only, not the uploaded files". The archive-format Contents note stays the one place that says where the files are, since "your bucket" on the overview means the backup bucket.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The docs site is built with Mintlify. Pages are `.mdx` files under `docs/`; a page's URL is its path under `docs/` without the extension (`docs/admin/backups/overview.mdx` is `/admin/backups/overview`). `docs/docs.json` holds navigation and redirects, including `/references/platform-api/latest/:slug*`, which today points to `2026-11-04` and which PR 2 repoints to `2026-11-16`.

A data backup is a nightly zip archive of a workspace, written to the customer's bucket. Two pages describe it:

- `docs/admin/backups/overview.mdx`: what backups are. Its "What's included" section has the sentence "As of writing, the latest version is [2026-11-04.everglades](/references/platform-api/2026-11-04), which contains:" followed by a bullet list of resources, each linked to its concept page, some with a parenthetical such as "(including their contents)".
- `docs/admin/backups/archive-format.mdx`: the archive layout. It has a `<Tree>` of the unzipped archive (the `resources` folder lists each `.jsonl` file alphabetically), a "Manifest" section with a JSON example and `<ParamField>` blocks (the `platform_api_version` field's "Example:" line names the version), a "Resources" section with a `| File | Resource | Expansions |` table (alphabetical by file, blank Expansions cell when there are none) followed by paragraphs, ending with "Archived devices, releases and config types are included in the archive. Deleted API keys are also included, ...", and a "Contents" section about `contents/`.

Verified facts from the backend (`backend/internal/backups/export/`):

- `manifest.go`: `FormatVersion = "0.1"` (unchanged). Resource types now include `file_rules`, `upload_collections` and `uploads`.
- `export.go` sets `platform_api_version` from the lassen models package, whose `API_VERSION` is `2026-11-16.lassen` (`backend/internal/servers/platform/lassen/models/model_api_version.go`).
- `db_tables.go`, `getUploadTables`: file rules and uploads are listed with no filter; upload collections with `ArchivedAny`, so archived ones are included, like devices, releases and config types. None of the three has expansions.
- Uploads are records only. The archive does not copy uploaded files; each stays in the bucket its upload names.
- Lassen calibrations have no `role` or `aspect` fields (the everglades model has both). `component` is the flat component slug: the old component, role and aspect slugs joined with `-`, so a calibration reported as component `meca500` with role `aligner-1` reads `meca500-aligner-1`.

Concept pages to link: `/data-uploads/concepts/file-rules/overview`, `/data-uploads/concepts/upload-collections`, `/data-uploads/concepts/uploads`.

Precedent: mirurobotics/docs#295 created these pages and #304 tightened the overview (`git log origin/main -- docs/admin/backups/`; run `git fetch --depth=200 origin main` first if the clone is shallow). Follow their tone: short declarative sentences, second person, no callout components on these pages.

CI (`.github/workflows/ci.yml`, workflow `CI`, runs on `pull_request`, so a PR must be open) has jobs `changes`, `lint`, `audit`, `lint-custom-linter` and `test-custom-linter` (the last two skip their steps unless `tools/lint/**` changes). `lint` runs `pnpm run test:lint`, then `./scripts/lint.sh` (the Go MDX prose linter in `tools/lint/`: sentence-case headings, no `--` in prose, import and redirect rules; ESLint on MDX; CSpell against `cspell.json`; OpenAPI checks), then `pnpm run validate` (`mint validate`, the Mintlify build check). These are the test steps for this change; they run only in CI, never locally.

## Plan of Work

Edit 1, `docs/admin/backups/overview.mdx`, "What's included":

- Replace `[2026-11-04.everglades](/references/platform-api/2026-11-04)` with `[2026-11-16.lassen](/references/platform-api/2026-11-16)`.
- In the bullet list, after the Calibrations item and before API keys, add:

        - [File rules](/data-uploads/concepts/file-rules/overview)
        - [Upload collections](/data-uploads/concepts/upload-collections)
        - [Uploads](/data-uploads/concepts/uploads) (records only; the uploaded files stay in your bucket and aren't copied)

Edit 2, `docs/admin/backups/archive-format.mdx`:

- Archive contents tree, `resources` folder: add `<Tree.File name="file_rules.jsonl" />` after `devices.jsonl`, and `<Tree.File name="upload_collections.jsonl" />` then `<Tree.File name="uploads.jsonl" />` after `releases.jsonl`, keeping alphabetical order.
- Manifest JSON example: `"platform_api_version": "2026-11-16.lassen"`. Leave `format_version` and `resource_counts` as they are.
- `platform_api_version` ParamField: change its example line to read

        Example: `2026-11-16.lassen`
- Resources table: add rows in alphabetical position, with a blank Expansions cell:

        | `file_rules.jsonl` | File rule | |            (after devices.jsonl)
        | `upload_collections.jsonl` | Upload collection | |   (after releases.jsonl)
        | `uploads.jsonl` | Upload | |                  (after upload_collections.jsonl)

- Change "Archived devices, releases and config types are included in the archive." to "Archived devices, releases, config types and upload collections are included in the archive." Keep the deleted API keys sentence.
- After that paragraph, add a calibrations paragraph:

        Calibrations changed shape in `2026-11-16.lassen`. A lassen calibration has no `role`
        or `aspect` field, and its `component` is the flat component slug: a calibration
        reported as component `meca500` with role `aligner-1` reads `meca500-aligner-1`.
        Archives written in `2026-11-04.everglades` keep the separate fields. Check the
        manifest's `platform_api_version` to know which shape `calibrations.jsonl` uses.

- At the end of "Contents", add:

        Uploaded files are not in the archive. `uploads.jsonl` holds only the upload records;
        each file stays in the bucket its upload names.

Nothing else changes. No headings are added, so the heading-case rule is untouched; write no `--` in prose.

## Concrete Steps

Milestone 1 (from the docs repo root):

1. `git status` shows a clean tree on `docs/backups-lassen`.
2. Apply Edit 1 and Edit 2 with a text editor.
3. Static check by reading, not tooling. `git diff` shows only the edits above. Then:

        grep -n "everglades" docs/admin/backups/*.mdx
        # expect matches only inside the calibrations paragraph of archive-format.mdx
        grep -c '^| `.*\.jsonl`' docs/admin/backups/archive-format.mdx
        # expect 15 (12 existing table rows + 3 new)
        ls docs/data-uploads/concepts/file-rules/overview.mdx \
           docs/data-uploads/concepts/upload-collections.mdx \
           docs/data-uploads/concepts/uploads.mdx
        # expect all three listed, no errors
4. Commit (signed by the repo's git config; never `--no-verify`):

        git add docs/admin/backups/overview.mdx docs/admin/backups/archive-format.mdx
        git commit -m "docs(backups): back up lassen resources"

Milestone 2 (from the docs repo root), via `$preflight`:

1. Resolve the base per the open question: `git fetch origin` and look for PR 2's branch (`gh api 'repos/mirurobotics/docs/pulls?state=open' --jq '.[] | "\(.number) \(.head.ref) \(.title)"'`). If it exists, `git rebase origin/<pr2-branch>`.
2. `git push -u origin docs/backups-lassen` and open a draft PR (`$pr`) with the Decisions section in the body, linking backend#1276, #1277 and PR 2.
3. Watch CI on the pushed head (`gh run list --branch docs/backups-lassen`; on failure `gh run view <id> --log-failed`). If CSpell flags `lassen`, `meca` or another word from this diff, add exactly that word to `words` in `cspell.json` (the list is ASCII-sorted: `lassen` between `kuka` and `localtime`, `meca` between `mcap` and `mgmt`), commit as `docs(backups): add lassen terms to the spell list`, and push. Fix any other failure at its cause with a new commit; never amend or force-push pushed commits except `--force-with-lease` after a rebase onto the base branch.
4. Repeat until preflight reports `CLEAN`.

## Validation and Acceptance

The test steps for this docs change are the CI checks, run by pushing, never locally: CI job `lint` passes on the pushed head, which proves the MDX prose rules, ESLint MDX checks, CSpell spelling and `mint validate` build all accept the edits; `changes` and `audit` pass; `lint-custom-linter` and `test-custom-linter` report success with their steps skipped.

Acceptance, on the PR's Mintlify preview:

- `/admin/backups/overview` says the latest version is `2026-11-16.lassen`, linking to `/references/platform-api/2026-11-16`, and lists File rules, Upload collections and Uploads with working links, with Uploads saying the files aren't copied.
- `/admin/backups/archive-format` shows `file_rules.jsonl`, `upload_collections.jsonl` and `uploads.jsonl` in the tree and the table, `2026-11-16.lassen` in the manifest example and field example, upload collections in the archived sentence, the calibrations shape paragraph, and the uploaded-files note.
- The lassen reference link resolves once PR 2 is in the base.

Validation is complete only when `$preflight` reports `CLEAN`: CI green on the pushed branch head. The PR does not leave draft, and the task is not reported complete, before that.

## Idempotence and Recovery

The edits are plain text replacements; reapplying them after checking `git diff` is safe. Before committing, `git checkout -- <file>` undoes an edit; after, `git revert <sha>`. Rebasing onto PR 2's branch can be redone with `git rebase --abort` on conflict. Fix CI failures forward with new commits.

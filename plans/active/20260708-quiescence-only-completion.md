# Make quiescence the only documented completion check for data uploads

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (this repo, `/home/ben/miru/workbench5/repos/docs`) | read-write | Edit two data-uploads pages and `cspell.json` on branch `feat/data-recording` |
| `agent/` (`/home/ben/miru/workbench5/repos/agent`) | read-only | Terminology confirmation only (`stability_window_secs`, `incomplete` flag) — do not modify; it is under active development in other sessions |

This plan lives in `docs/plans/` because all changes are documentation edits in this repo.

## Purpose / Big Picture

Product decision (authoritative, overrides anything found in code): for the first rollout, the Miru Agent does **not** detect file formats (MCAP, parquet finalization markers) to decide when a recording is complete. **Quiescence** — the file's size and modification time unchanged for the rule's `stability_window_secs` — is the *only* completion check. The data-uploads docs currently describe a two-check model ("finalization marker" plus "quiescence" fallback). After this change, a reader of the docs sees a single, accurate story: a file is considered finished when it has not changed for the rule's stability window, and the `incomplete` upload flag is described without reference to format detection.

## Progress

- [ ] Milestone 1: rewrite completion-detection claims across the two affected pages and prune `cspell.json`; lint and preflight pass; commit.

## Surprises & Discoveries

(Authoring-time notes; add implementation entries as you go.)

- The heading-case allowlist (`tools/lint/linter/headingcase/headingcase.go`, `allowlist()` at line 47) contains **no** `MCAP` or `parquet` entry, despite `plans/completed/20260707-data-recording-docs-refactor.md:140` claiming "allowlist already includes AWS/GCS/WIF/STS/MCAP". Nothing to remove there.
- The agent repo confirms the field name `stability_window_secs` (`agent/libs/backend-api/src/models/upload_rule_source.rs:18-19`) and that `incomplete` is a device-declared flag on the create-upload (mint) request (`agent/libs/backend-api/src/models/create_upload_request.rs:26-28`), so describing it as "flagged by the device when collected" is format-agnostic and accurate.

## Decision Log

(Authoring decisions; add implementation entries as you go.)

- Decision: remove the finalization-marker material outright rather than reframing it as future work.
  Rationale: this branch documents the first rollout; documenting unshipped behavior invites support confusion. Task instructions prefer plain removal.
  Date/Author: 2026-07-08, planning agent.
- Decision: keep the `## Completion detection` heading and its `#completion-detection` anchor in `docs/data-uploads/uploads.mdx`.
  Rationale: two pages link to the anchor (`docs/data-uploads/uploads.mdx:114`, `docs/data-uploads/upload-rules.mdx:89`); quiescence is still "detection", so the heading remains truthful and no link churn is needed.
  Date/Author: 2026-07-08, planning agent.
- Decision: keep the "defer completion to the file producer" paragraph but drop its "(or drop a marker file)" parenthetical.
  Rationale: write-then-rename works with pure quiescence/glob matching; a marker file only helps if the Agent watches for it, which it does not in the first rollout.
  Date/Author: 2026-07-08, planning agent.
- Decision: remove `"mcap"` and `"parquet"` from `cspell.json`; keep `"quiescence"`.
  Rationale: after the edits, the only occurrences of mcap/parquet in spell-checked files (all `docs/**/*.mdx`) are gone; "quiescence"/"quiescent" wording remains in `docs/data-uploads/upload-rules.mdx` and `docs/data-uploads/uploads.mdx`.
  Date/Author: 2026-07-08, planning agent.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

This repo is the Mintlify documentation site. Pages are MDX under `docs/`; the data-uploads section is `docs/data-uploads/`. `<ParamField>` blocks document object properties. Lint tooling: `./scripts/lint.sh` (Go prose linter in `tools/lint/`, ESLint, CSpell against `cspell.json`, OpenAPI checks) and `./scripts/preflight.sh` (lint smoke tests, Go lint/coverage, `lint.sh`, `audit.sh`, bats). Work happens on branch `feat/data-recording` (clean at `05d92df`, based on `main`).

Terms:

- **Quiescence / stability window**: a file's size and modification time staying unchanged for `source.stability_window_secs` seconds (rule field, default 60). Defined in `docs/data-uploads/upload-rules.mdx` (Sources section) and `docs/data-uploads/defining-releases.mdx`.
- **Finalization marker**: trailing bytes some formats (MCAP, parquet) write on clean close. The docs currently present verifying this as a shipped completion check. It is not shipped; all mentions must go.
- **`incomplete`**: a boolean property on an upload ledger entry. In reality it is set from a flag the device sends on the create-upload request; its docs wording must not imply format-aware finalization checks.

Complete inventory of passages to change (exact current text, verified 2026-07-08):

1. `docs/data-uploads/overview.mdx:9-10` — "detects matching data, confirms each file is / complete, and uploads it directly to your bucket." ("confirms … complete" implies verification beyond quiescence.)
2. `docs/data-uploads/overview.mdx:20` — "The device detects a completed file and makes an **upload request** …" ("completed" reads as confirmed; standardize on "finished".)
3. `docs/data-uploads/uploads.mdx:112-114` (`incomplete` ParamField) — "Whether the file was incomplete when collected (uploaded without a confirmed finalization). Independent of `status` — see [completion detection](#completion-detection)."
4. `docs/data-uploads/uploads.mdx:157` — "The Agent detects a completed file matching a rule and **requests an upload** from …"
5. `docs/data-uploads/uploads.mdx:172-190` (`## Completion detection` section) — the two-check model: "Completion is determined with two complementary checks:", the "**Finalization marker** — for formats that end with one (MCAP, parquet) …" bullet, the "**Quiescence**" bullet, "(or drop a marker file)" in the defer paragraph, and "A file that goes quiet but was never closed cleanly is still uploaded, but flagged [`incomplete`](#properties)." (The last sentence presumes the Agent can tell a file was "never closed cleanly" — format detection.)
6. `docs/data-uploads/upload-rules.mdx:87-89` (`stability_window_secs` ParamField) — "Files in a format with a finalization marker (such as MCAP or parquet) are detected directly; this window is the fallback for other files."
7. `cspell.json:43` `"mcap",` and `cspell.json:48` `"parquet",` — unused after the edits.

Explicitly unchanged (verified accurate under quiescence-only):

- `docs/data-uploads/upload-rules.mdx:51` `## File formats` heading — it describes the YAML *rule-file* format, not data-file formats. Keep.
- `docs/data-uploads/upload-rules.mdx:60` comment `# quiescence before a file is considered finished`. Keep.
- `docs/data-uploads/uploads.mdx:11-12` "detecting when a matching file is finished". Keep — detection *is* quiescence.
- `docs/data-uploads/defining-releases.mdx:68-71` (`source.stability_window_secs`) — already quiescence-only. Keep.
- `tools/lint/linter/headingcase/headingcase.go` allowlist — contains no MCAP/parquet entries. No change.
- No other page in `docs/` mentions MCAP, parquet, finalization, or format-aware completion (grep-verified; the only other "format"/"detect" hits are unrelated: gcloud `--format` flags, schema-digest canonical format, changelog UI notes).

## Plan of Work

All edits in this repo root (`/home/ben/miru/workbench5/repos/docs`), branch `feat/data-recording`.

**Edit 1 — `docs/data-uploads/overview.mdx`.** Rewrite the second sentence of the intro (lines 8-10) so the Agent "waits for each file to finish" instead of "confirms each file is complete"; the full sentence becomes (rewrap lines as needed):

    The [Miru Agent](/developers/agent/overview) detects matching data, waits for
    each file to finish, and uploads it directly to your bucket.

On line 20, change "detects a completed file" to "detects a finished file".

**Edit 2 — `docs/data-uploads/uploads.mdx`, `incomplete` ParamField (lines 109-115).** Replace the description body with:

      Whether the device flagged the file as incomplete when it was collected.
      Independent of `status` — see [completion detection](#completion-detection).

**Edit 3 — `docs/data-uploads/uploads.mdx`, line 157.** Change "The Agent detects a completed file matching a rule" to "The Agent detects a finished file matching a rule".

**Edit 4 — `docs/data-uploads/uploads.mdx`, `## Completion detection` section (lines 172-190).** Replace the section body (keep the heading) with:

    The Agent uploads a file once it is confident the file has finished being
    written to. Completion is determined by **quiescence**: a file whose size and
    modification time have held
    steady for the rule's [`stability_window_secs`](/data-uploads/upload-rules#sources)
    is considered finished and eligible for upload. This works for any file type, but
    it is a heuristic: "looks done," not "proven done."

    If you would rather not rely on the heuristic, you can **defer completion to the
    file producer**: have it write to a temporary name the rule's glob does not match
    and rename on completion, and the Agent will only ever see finished files — zero
    false positives.

This removes the two-check list, the finalization-marker bullet, the marker-file parenthetical, and the "never closed cleanly … flagged `incomplete`" sentence.

**Edit 5 — `docs/data-uploads/upload-rules.mdx`, `stability_window_secs` ParamField (lines 84-90).** Delete the finalization-marker sentence so the body reads:

      How long, in seconds, a matching file's size and modification time must stay
      unchanged (quiescent) before it is considered finished and eligible for upload.
      See [completion detection](/data-uploads/uploads#completion-detection).

**Edit 6 — `cspell.json`.** Remove the `"mcap",` (line 43) and `"parquet",` (line 48) entries from `words`. Keep `"quiescence"`.

## Concrete Steps

All commands run from `/home/ben/miru/workbench5/repos/docs`.

1. Confirm starting state:

       git status --short          # expect: clean (besides this plan file)
       git branch --show-current   # expect: feat/data-recording

2. Apply Edits 1-6 above with a text editor.

3. Verify no format-detection claims remain and legitimate uses survive:

       grep -rn -i -E 'mcap|parquet|finaliz' docs --include='*.mdx'
       # expect: no output

       grep -rn -i -E 'confirms each file|completed file|closed cleanly' docs/data-uploads
       # expect: no output

       grep -n 'File formats' docs/data-uploads/upload-rules.mdx
       # expect: line 51 heading still present (YAML rule-file format)

       grep -rn 'completion-detection' docs --include='*.mdx'
       # expect: 2 hits (uploads.mdx ParamField, upload-rules.mdx) — anchor targets intact

       grep -rn -i 'mcap\|parquet' cspell.json
       # expect: no output

4. Run lint and preflight:

       ./scripts/lint.sh        # expect: "All documentation lint checks passed."
       ./scripts/preflight.sh   # expect: all sections pass, exit 0

5. Commit (single milestone commit — docs pages, cspell, and the plan file, which moves to `plans/active/` when implementation begins):

       git add docs/data-uploads/overview.mdx docs/data-uploads/uploads.mdx \
               docs/data-uploads/upload-rules.mdx cspell.json plans/
       git commit -m "docs(data-uploads): document quiescence as the only completion check"

   Note the repo requires verified signatures — confirm the commit is signed before pushing (re-sign if `git log --show-signature -1` shows no signature).

## Validation and Acceptance

Acceptance is behavior a reviewer can verify on the branch:

- `grep -rn -i -E 'mcap|parquet|finaliz' docs --include='*.mdx'` produces no output.
- `docs/data-uploads/uploads.mdx` `## Completion detection` describes exactly one mechanism (quiescence via `stability_window_secs`) plus the producer-side rename pattern; no bullet list of checks remains.
- The `incomplete` ParamField in `docs/data-uploads/uploads.mdx` mentions no finalization or format detection.
- The `## File formats` heading in `docs/data-uploads/upload-rules.mdx` and both `#completion-detection` cross-links remain intact.
- `./scripts/lint.sh` prints "All documentation lint checks passed." and `./scripts/preflight.sh` exits 0. These are the repo's test suites; no Go code changes, so `tools/lint` tests are unaffected and must still pass (run by preflight).

## Idempotence and Recovery

All steps are safe to repeat: the greps are read-only, and re-applying an already-applied text edit is a no-op (the old text simply won't be found). Everything is on a feature branch — recover from any misstep with `git checkout -- <file>` before commit, or `git revert <sha>` after. If `cspell` unexpectedly flags a remaining "quiescent"/"quiescence" form after removing entries, re-add only the flagged word to `cspell.json` rather than restoring `"mcap"`/`"parquet"`.

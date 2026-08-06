# Document every Miru schema annotation (add `instance_format`)

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
| --- | --- | --- |
| `/home/ben/miru/workbench2/repos/docs` | read-write | Owns this plan and every file changed. All commits happen here. |
| `/home/ben/miru/workbench2/repos/cli-private` | read-only | CLI source; shows how annotations are parsed. |
| `/home/ben/go/pkg/mod/github.com/mirurobotics/core@v0.8.2` | read-only | Vendored `core` library pinned by the CLI; the authoritative annotation key/value definitions. |
| `/home/ben/miru/workbench2/repos/backend` | read-only | Server-side defaults and validation for annotations. |
| `/home/ben/miru/workbench2/repos/core` | read-only | Unreleased `core`; consulted only to confirm a future rename. Do not document unreleased behavior. |

No file outside `repos/docs` may be modified.

## Purpose / Big Picture

Miru users annotate their config schemas with **schema annotations** (defined under Context below). The docs snippet at `docs/snippets/references/cli/releases/create/schema-annotations.mdx` is the single source rendered on three different pages, and it currently documents only two of the three annotations the CLI actually supports. The third, **instance format** (`x-miru-instance-format` / `@miru(instance_format=…)` / `instance_format`), is invisible to readers of those three pages even though it shipped in stable CLI v0.10.1.

After this change, a reader on any of the three pages sees all three annotations, each with JSON Schema, CUE, and Opaque syntax, with an accurate list of accepted values per schema language. Additionally, three currently-inconsistent statements about opaque instance formats (the snippet, `opaque.mdx`, and an uncommitted edit to `formats.mdx`) are reconciled so the docs do not contradict themselves. Success is visible by running `pnpm dev` and seeing a third annotation block on `/cfg-mgmt/primitives/schemas/overview`, and by every CI check on PR #138 passing apart from the pre-existing `audit` failure.

## Progress

- [x] Milestone 1 — Revert the uncommitted `formats.mdx` edit to clear the dirty tree; no commit.
- [x] Milestone 2 — Confirm the accepted-value list against the backend doc comment; record it in the Decision Log.
- [x] Milestone 3 — Add the instance-format annotation to the snippet, align `opaque.mdx` and `overview.mdx`, add cspell words if needed; commit.
- [x] Milestone 4 — Render check, preflight plus bats, force-push, verify CI on PR #138 (all green except the pre-existing `audit`).

## Surprises & Discoveries

- 2026-08-05 — `node_modules/` was absent, so the first `pnpm run lint` aborted at `== ESLint (MDX) ==` with `Command "eslint" not found`. `pnpm install --frozen-lockfile` fixed it; not a content problem.
- 2026-08-05 — cspell required no new words. `text`, `xml`, and `yaml` are all covered by built-in dictionaries, and `jsonc` never entered the content because the backend narrows it away.
- 2026-08-05 — `./scripts/preflight.sh` aborted at `=== Audit ===` exactly as the plan predicted, on pre-existing advisories in transitive dev dependencies: `fast-uri` via `mint` (GHSA-7p8r-x3mc-p8w7) and `brace-expansion` via `eslint` (GHSA-rgw5-rvv9-x895), reported as "3 vulnerabilities found / Severity: 3 high (1 ignored)". A pristine `main` worktree produces the byte-identical advisory set, confirming this is not a regression from this diff, which touches no `package.json` and no `pnpm-lock.yaml`. `bats pub/scripts/agent/check-miru-access_test.bats` was run separately and passed all 30 tests.

## Decision Log

- 2026-08-05 — **Accepted instance-format values confirmed against the backend.** `repos/backend/internal/configs/domain/config_instances/formats.go:9-22` still carries the doc comment "SupportedFormats is narrower than the set of instance formats core accepts: jsonc is deliberately excluded for every language. Only opaque schemas admit xml and text, since they impose no structure on the instance content," and `SupportedFormats` still returns `{Json, Yaml, Xml, Text}` for `OpaqueLang` and `{Json, Yaml}` otherwise. Documented values are therefore `json`, `yaml` for JSON Schema and CUE; `json`, `yaml`, `xml`, `text` for Opaque; `jsonc` is documented nowhere.
- 2026-08-05 — The CLI/`core`-vs-backend `jsonc` divergence is noted as a possible follow-up for the backend owner. Filing it was explicitly out of scope for this plan and was not done.
- 2026-08-05 — No headings were added anywhere, so no `heading-case` allowlist entry and no `tools/lint/**` edit were needed (`git diff --stat main...HEAD -- tools/lint` prints nothing, as required).

## Outcomes & Retrospective

All four milestones landed as written; no deviation from the plan of work was needed. Four docs files changed: the shared annotation snippet gained a third `<ParamField path="instance format">` with the JSON Schema / CUE / Opaque `<CodeGroup>`, `opaque.mdx` gained `text` in both its key table and its instance-format bullet list, `formats.mdx`'s Opaque row became `JSON, YAML, XML, Text`, and `overview.mdx`'s `## Properties` section gained a sibling `instance format` field. The load-bearing two-trailing-space separator lines survived: `cat -A … | grep -c '^  \$'` prints `6`, up from `4`. `grep -rn jsonc docs/ --include='*.mdx'` prints nothing.

The plan's foresight about the two silent hazards paid off. Pre-identifying the trailing-space convention meant the new `<CodeGroup>` was written to match rather than discovered broken later, and pre-identifying the `audit` failure as pre-existing meant it was recognized and proven rather than debugged.

## Context and Orientation

**Terms.** A *config schema* is a file (JSON Schema, CUE, or Opaque) that Miru stores and, for the first two languages, uses to validate config instances. A *config instance* is the concrete configuration file written onto a device. A *schema annotation* is a key embedded in the schema file that carries Miru metadata rather than validation logic. A *schema language* is which of the three formats the schema is written in. An *Opaque schema* is metadata-only: it carries annotations but validates nothing.

**Repository.** `repos/docs` is a Mintlify docs site. Content lives under `docs/`, reusable fragments under `docs/snippets/`, and plans under `plans/{backlog,active,completed}/`. There is no `CLAUDE.md` or `AGENTS.md`; conventions are enforced entirely by tooling (see "Lint rules" below).

**Git state at the time of writing.** Branch `docs/opaque-schema-language` (PR #138) has just been rebased onto `main` and has **not** been force-pushed. The working tree has two uncommitted modifications: `docs/snippets/schemas/formats.mdx` (a stray edit dealt with in Milestone 1) and this plan file itself (the living-document sections you update as you work). Only the first must be resolved before editing; the plan file's own updates are committed alongside the docs change in Milestone 3.

### The authoritative annotation list

Exactly three annotations exist, plus one Opaque-only language marker. Canonical struct `Annotations{ConfigTypeSlug, InstanceFilepath, InstanceFormat}` at `core@v0.8.2/pkg/schemas/annotations.go:8-12`. The CLI entry point that parses them is `repos/cli-private/internal/domain/cfgschs/load.go:363` (`parseSchemaAnnotations`). The CLI pins `core` v0.8.2 at `repos/cli-private/go.mod:11`.

| Concept | JSON Schema key | CUE attribute | Opaque key |
| --- | --- | --- | --- |
| Config type slug | `x-miru-config-type` | `@miru(config_type="…")` | `config_type` |
| Instance filepath | `x-miru-instance-filepath` | `@miru(instance_filepath="…")` | `instance_filepath` |
| Instance format | `x-miru-instance-format` | `@miru(instance_format="…")` | `instance_format` |
| Language marker | — | — | `language: opaque` |

Key-string citations:

- JSON Schema keys: `/home/ben/go/pkg/mod/github.com/mirurobotics/core@v0.8.2/pkg/schemas/jsonschema/annotations.go:15-17`
- CUE attribute keys: `/home/ben/go/pkg/mod/github.com/mirurobotics/core@v0.8.2/pkg/schemas/cue/attributes.go:53-55`
- Opaque keys: `/home/ben/go/pkg/mod/github.com/mirurobotics/core@v0.8.2/pkg/schemas/opaque/compile.go:14-17`

There are no other `x-miru-*` keys in `core` (verified by grep). `$miru_tag_type_field_id` appears in CLI testdata but is inert — **do not document it**.

**1. Config type — required, no default. Already documented.** Must sit at the schema document root (JSON Schema) or on a schema-level `@miru()` declaration attribute in CUE; field-level attributes are ignored (`core@v0.8.2/pkg/schemas/cue/attributes.go:26-33,76-83`). No change needed.

**2. Instance filepath — optional. Already documented.** Absolute path string, defaulted **server-side** to `/srv/miru/configs/<config-type-slug>.<ext>` (`repos/backend/internal/configs/domain/config_schemas/filepath.go:7-22`). That default matters here: because it ends in `.json`, the effective default `instance_format` is `json`. No change needed.

**3. Instance format — optional. UNDOCUMENTED in the snippet. This is the gap this plan closes.** Accepted values differ per schema language:

- JSON Schema: `json`, `yaml`, `jsonc` — `core@v0.8.2/pkg/schemas/jsonschema/format.go:9-17`
- CUE: `json`, `yaml`, `jsonc` — `core@v0.8.2/pkg/schemas/cue/format.go:9-17`
- Opaque: `json`, `yaml`, `jsonc`, `xml`, `text` — `core@v0.8.2/pkg/schemas/opaque/format.go:16-27`

Parsing sites: `core@v0.8.2/pkg/schemas/jsonschema/annotations.go:62-75`, `…/cue/attributes.go:118-128`, `…/opaque/compile.go:139-145`. Unknown value produces e.g. `unknown JSON Schema instance format '<value>', must be json, yaml or jsonc` (analogous wording for CUE and Opaque). The default is inferred **server-side** from the instance filepath's extension (`resolveInstFormat`, `repos/backend/internal/configs/services/config_schemas/create.go:384-401`); an unmappable extension yields `CannotInferInstFormat`. Because `instance_filepath` itself defaults to `…/<slug>.json`, the effective default is `json`. This annotation shipped in stable CLI v0.10.1 (commit `f17040f`), so it is safe to document independent of the opaque beta.

**The backend narrows this list, and that is settled — not an open question.** `repos/backend/internal/configs/domain/config_instances/formats.go:9-22` carries an explicit doc comment above `SupportedFormats`: "SupportedFormats is narrower than the set of instance formats core accepts: jsonc is deliberately excluded for every language. Only opaque schemas admit xml and text, since they impose no structure on the instance content." The end-to-end accepted values are therefore **`json`, `yaml`** for JSON Schema and CUE, and **`json`, `yaml`, `xml`, `text`** for Opaque. **Do not document `jsonc`** — the CLI and `core` accept it but the server rejects it. **Do document `text`** for Opaque; `opaque.mdx` currently omits it and must gain it.

**4. Opaque `language` marker — already documented, no gap.** Covered by the key table in `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`. Note only that unreleased `core` renames `language` → `schema_language` (`repos/core/pkg/schemas/opaque/compile.go:14-33`); the CLI on `core` v0.8.2 does **not** accept `schema_language`. **Document `language` only; change nothing for this.**

### Current state of the snippet

`docs/snippets/references/cli/releases/create/schema-annotations.mdx` holds two `<ParamField>` blocks — `path="config type" required` and `path="instance file path"` — each containing one `<CodeGroup>` with three fenced tabs in this order: ` ```yaml JSON Schema `, ` ```cue CUE `, ` ```yaml Opaque `.

Formatting convention that any new `<CodeGroup>` must match exactly: the block is indented two spaces inside the `<ParamField>`; a blank line follows `<CodeGroup>`; between consecutive fences there are **two** separator lines — one truly blank line followed by a line containing exactly two trailing spaces; one blank line precedes `</CodeGroup>`. Copy an existing block and edit it rather than typing a new one, so the trailing-space lines survive.

### Files that render the snippet

1. `docs/cfg-mgmt/primitives/schemas/overview.mdx` — imports it as `Annotations` (line 9), renders under `## Schema annotations`. This page **also duplicates** the instance-filepath prose in its own `<ParamField path="instance file path" type="string">` at line 62, in the `## Properties` section.
2. `docs/cfg-mgmt/create-a-release.mdx` — imports as `Annotations`, renders under `### Schema annotations`.
3. `docs/getting-started/quick-start/create-release.mdx` — imports as `SchemaAnnotations`, renders under `### Schema annotations`.

`docs/references/cli/release-create.mdx` does not import the snippet; it links to `/cfg-mgmt/primitives/schemas/overview#schema-annotations`, so it needs no edit.

Related: `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` already documents `instance_format` in its key table (currently "`json`, `yaml`, or `xml`") and in a "File formats" bullet list (JSON / YAML / XML). `docs/snippets/schemas/formats.mdx` is a file-formats table imported only by `schemas/overview.mdx`; it currently carries the uncommitted edit described in Milestone 1.

### Lint rules the change must satisfy

`package.json` defines exactly three scripts: `pnpm run lint` → `./scripts/lint.sh`, `pnpm run test:lint` → `./tests/test-lint.sh`, `pnpm dev` → `cd docs && mint dev`. `./scripts/lint.sh` runs four phases: the custom Go MDX linter in `tools/lint` (built with `go build`), `pnpm exec eslint --max-warnings=0`, `pnpm exec cspell lint --config cspell.json`, and `pnpm exec mint openapi-check` over `docs/references/**/*.yaml`.

The ten custom rules (source in `tools/lint/linter/`): `no-double-dash` (no bare `--` in prose; use an em dash — code spans are exempt), `heading-case` (strict sentence case on frontmatter `title:` and every heading), `import-resolves`, `import-used`, `import-sorted` (case-insensitive ascending by path), `import-component-style`, `import-mdx-style` (default import, ends with `;`), `import-block` (no blank lines inside the import block), `image-domain`, `redirects`.

`heading-case` uses a **hardcoded Go allowlist**: `allowlist()` at ~line 45 of `tools/lint/linter/headingcase/headingcase.go`, a `map[string]struct{}` grouped by comment (Acronyms / Proper nouns / Codenames / Event identifiers). Matching is case-sensitive and exact on whitespace-split tokens after stripping outer punctuation. `JSON`, `CUE`, `CLI`, `API` are present; **`YAML`, `XML`, and `JSONC` are not**. Editing `tools/lint/**` triggers the `lint-custom-linter` and `test-custom-linter` CI jobs plus the Go coverage gates (`.covgate` files, `tools/lint/scripts/covgate.sh`). **Prefer adding no new headings at all** — `<ParamField>` labels are not headings, so the annotation block needs none, and no allowlist change is required.

cspell's dictionary is the `words` array in `cspell.json`. It currently contains none of `jsonc`, `xml`, `yaml` (cspell's built-in dictionaries may already cover `xml`/`yaml`; `jsonc` is the likely miss). Words inside fenced code blocks are still checked by cspell, so if `jsonc` appears anywhere it may need adding.

CI (`.github/workflows/ci.yml`): `lint` (runs `pnpm run test:lint` then `./scripts/lint.sh`), `audit` (`./scripts/audit.sh`), `shell-tests` (bats), plus `lint-custom-linter` / `test-custom-linter` gated on `tools/lint/**`. **Known pre-existing failure:** `plans/completed/20260726-opaque-schema-language.md` records an `audit` CI failure already present on this branch and unrelated to doc edits. Distinguish it from anything this work introduces; do not treat it as caused by this change, and do not attempt to fix it here.

## Plan of Work

**Milestone 1 — clear the dirty tree.** The working tree changes the Opaque row of `docs/snippets/schemas/formats.mdx` from `| Opaque | JSON, YAML | JSON, YAML, XML |` to `| Opaque | YAML | JSON, YAML, XML, Plain Text |`. Both halves are wrong: `opaque.mdx` states opaque schemas are written as JSON **or** YAML, so the schema-format column must stay `JSON, YAML`; and the instance-format column is set once, correctly, in Milestone 3. Revert the file rather than half-fixing it, so the row is written exactly once in a single commit.

**Milestone 2 — confirm the value list.** Read the doc comment above `SupportedFormats` in `repos/backend/internal/configs/domain/config_instances/formats.go:9-22` and record the outcome in the Decision Log with the citation: `json`, `yaml` for JSON Schema and CUE; `json`, `yaml`, `xml`, `text` for Opaque; `jsonc` documented nowhere. Optionally note the CLI/`core`-vs-backend `jsonc` divergence as a follow-up for the backend owner — filing that issue is **not** in scope here and does not block this plan.

**Milestone 3 — write the docs.** Four files change, all consistent with the Milestone 2 list. In `docs/snippets/references/cli/releases/create/schema-annotations.mdx`, append a third `<ParamField path="instance format">` after the instance-file-path block. Its prose must say what the annotation does (declares the format of the deployed config instance), that it is optional, that it defaults to the format inferred server-side from the instance file path's extension — so effectively `json` under the default path `/srv/miru/configs/{config-type-slug}.json` — the accepted values, and that Opaque schemas accept a wider set. Then a `<CodeGroup>` with the same three tabs in the same order and the same whitespace convention:

      ```yaml JSON Schema
      x-miru-instance-format: "yaml"
      ```

      ```cue CUE
      @miru(instance_format="yaml")
      ```

      ```yaml Opaque
      instance_format: "xml"
      ```

Close with an `Examples:` line matching the config-type block's style. **Add no headings** — `<ParamField>` labels are not headings, so no `heading-case` allowlist change and no `tools/lint` edit is needed.

In `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`, the `instance_format` table row (line 46, currently "`json`, `yaml`, or `xml`") and the "Opaque schemas accept the widest set…" bullet list (lines 77-81, currently JSON / YAML / XML) both gain `text`. In `docs/snippets/schemas/formats.mdx`, set the Opaque row to `| Opaque | JSON, YAML | JSON, YAML, XML, Text |` as described in Milestone 1. In `docs/cfg-mgmt/primitives/schemas/overview.mdx`, the `## Properties` section already duplicates the instance-filepath field at line 62; **add** a sibling `<ParamField path="instance format" type="string">` mirroring it (immutable badge, one-line description, default, examples) — that section enumerates schema properties and omitting one is a real gap. This is decided: add it. Finally, if `pnpm run lint` flags a word, add it to the `words` array in `cspell.json` in sorted position (uppercase-first, then lowercase ascending); with `jsonc` dropped from the content, no addition is expected.

**Milestone 4 — validate and ship.** Render-check the three importing pages under `pnpm dev`, run `./scripts/preflight.sh` plus the bats tests, force-push the rebased branch, and check PR #138.

## Concrete Steps

All commands run from `/home/ben/miru/workbench2/repos/docs` unless stated otherwise. Commit from inside this repo — never from the workbench root.

### Milestone 1 — clear the dirty tree

    cd /home/ben/miru/workbench2/repos/docs
    git status --short
    git diff -- docs/snippets/schemas/formats.mdx

Expect two modified files — ` M docs/snippets/schemas/formats.mdx` (with the Opaque-row diff described above) and ` M plans/backlog/20260805-document-all-schema-annotations.md` (this plan's own living-document updates). Discard only the first; Milestone 3 rewrites that row correctly in one edit:

    git checkout -- docs/snippets/schemas/formats.mdx
    git status --short

Expect the second `git status --short` to list only this plan file. No commit in this milestone.

### Milestone 2 — confirm the accepted-value list

    sed -n '1,25p' /home/ben/miru/workbench2/repos/backend/internal/configs/domain/config_instances/formats.go

Confirm the doc comment above `SupportedFormats` still reads "jsonc is deliberately excluded for every language. Only opaque schemas admit xml and text", and that the function returns `{Json, Yaml, Xml, Text}` for `OpaqueLang` and `{Json, Yaml}` otherwise. Record in the Decision Log, with today's date and this file:line citation, the value lists to document — **`json`, `yaml`** for JSON Schema and CUE; **`json`, `yaml`, `xml`, `text`** for Opaque; **`jsonc` documented nowhere**. If the comment or the function has changed, stop and re-derive before writing any docs. No commit and no network access in this milestone.

### Milestone 3 — write the docs

Edit, in order:

1. `docs/snippets/references/cli/releases/create/schema-annotations.mdx` — append the `instance format` `<ParamField>` after the existing `instance file path` block (currently ends at line 63). Copy the existing `<CodeGroup>` block verbatim first, then edit the three fence bodies, so the two-space indentation and the two-trailing-space separator lines are preserved.
2. `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` — update the `instance_format` table row and the instance-format bullet list.
3. `docs/snippets/schemas/formats.mdx` — set the Opaque row to `| Opaque      | JSON, YAML     | JSON, YAML, XML, Text |`.
4. `docs/cfg-mgmt/primitives/schemas/overview.mdx` — add the sibling `instance format` `<ParamField>` to the `## Properties` section, immediately after the `instance file path` field at line 62.

Then verify the trailing-space separator lines survived the edit:

    cat -A docs/snippets/references/cli/releases/create/schema-annotations.mdx | grep -c '^  \$'

Expect `6` — two two-space separator lines inside each of the three `<CodeGroup>` blocks. Before this milestone's edit the same command prints `4`. If it prints `4` or `5`, your editor stripped trailing whitespace on save: disable that, re-copy an existing `<CodeGroup>` block verbatim, and re-check. No lint rule catches this, so this command is the only guard.

Then check spelling and lint quickly:

    pnpm run lint

If cspell reports `Unknown word (jsonc)` (or `xml` / `yaml`), add the word to the `words` array in `cspell.json` in sorted position and re-run. If `heading-case` fires, you added a heading — remove it; `<ParamField>` labels are not headings and no `tools/lint` change should be needed. If `no-double-dash` fires, replace the bare `--` with an em dash.

Expect the phase banners `== MDX Prose ==`, `== ESLint (MDX) ==`, `== CSpell ==`, `== OpenAPI ==` in that order, a `Checking docs/references/…` line per spec, and the final line `All documentation lint checks passed.` Then:

    git add -A
    git commit -m "docs(schemas): document the instance format schema annotation"

`git add -A` sweeps in five files: the four docs files above plus this plan file, whose Progress, Decision Log, and Surprises & Discoveries sections you have been updating as you go. That is intended — the plan's living-document updates ship with the change they describe.

### Milestone 4 — validate, push, verify CI

Render check:

    pnpm dev

Open each of these local URLs (the dev server prints its port, typically `http://localhost:3000`) and confirm three annotation blocks render, each with three working tabs labelled JSON Schema / CUE / Opaque:

- `http://localhost:3000/cfg-mgmt/primitives/schemas/overview` (section "Schema annotations")
- `http://localhost:3000/cfg-mgmt/create-a-release` (section "Schema annotations")
- `http://localhost:3000/getting-started/quick-start/create-release` (section "Schema annotations")

Stop the dev server with Ctrl-C.

Full preflight:

    ./scripts/preflight.sh

Expect these banners in order:

    === Lint Smoke Tests ===
    === Go Lint (tools/lint) ===
    === Go Coverage (tools/lint) ===
    === Lint ===
    === Audit ===

The first four phases must pass; `=== Lint ===` ends with `All documentation lint checks passed.` The `=== Audit ===` phase is expected to fail on the pre-existing repo-wide dependency advisory recorded in `plans/completed/20260726-opaque-schema-language.md`. Because `preflight.sh` runs under `set -euo pipefail`, that failure **aborts the script** — `=== Shell Script Tests ===` never runs. That is expected, not a regression. Confirm the advisory output matches the pre-existing one (same package and advisory IDs) and record it in Surprises & Discoveries, then run the skipped phase by hand:

    bats pub/scripts/agent/check-miru-access_test.bats

Any failure in a phase before `=== Audit ===`, or any new advisory in the audit output, must be fixed before proceeding.

Push the rebased branch and check CI:

    git push --force-with-lease
    gh pr checks 138 --watch

`--force-with-lease` is required because the branch was rebased; it refuses to push if someone else moved the remote branch, which is the safe behavior. If it is rejected, run `git fetch origin` and inspect `git log --oneline origin/docs/opaque-schema-language..HEAD` before deciding whether to retry.

Expect `gh pr checks 138` to report `changes`, `lint`, and `shell-tests` passing and `audit` failing. `lint-custom-linter` and `test-custom-linter` should not run at all, since `tools/lint/**` is untouched. Reproduce the `audit` failure on a pristine `main` worktree per Validation criterion 2 before treating it as pre-existing.

## Validation and Acceptance

Acceptance is behavioral. **Preflight must reach the `=== Audit ===` phase with every prior phase passing** — it cannot report fully clean, because `./scripts/audit.sh` fails on a pre-existing repo-wide advisory and `preflight.sh` runs under `set -e`, so it aborts there and never reaches `=== Shell Script Tests ===`. Run those separately (see Milestone 4). All of the following must hold:

1. `/home/ben/miru/workbench2/repos/docs/scripts/preflight.sh` prints, in order, `=== Lint Smoke Tests ===`, `=== Go Lint (tools/lint) ===`, `=== Go Coverage (tools/lint) ===`, `=== Lint ===` (this phase is `./scripts/lint.sh`, ending in `All documentation lint checks passed.`), and `=== Audit ===`, with no failure before `=== Audit ===`. This subsumes standalone `pnpm run test:lint` and `pnpm run lint`. Separately, `bats pub/scripts/agent/check-miru-access_test.bats` exits 0. Any new failure blocks completion.
2. After `git push --force-with-lease`, `gh pr checks 138` reports every check green on the **pushed head** except `audit`, and `lint-custom-linter` / `test-custom-linter` do not run at all. The `audit` failure must be shown to be pre-existing and independent of this diff: reproduce it on a pristine `main` worktree —

       git worktree add /tmp/docs-main main
       cd /tmp/docs-main && pnpm install --frozen-lockfile && ./scripts/audit.sh
       cd /home/ben/miru/workbench2/repos/docs && git worktree remove /tmp/docs-main

   and confirm it fails there with the same advisories. A local pass with an unpushed or outdated remote head does not satisfy this criterion.
3. Under `pnpm dev`, each of `/cfg-mgmt/primitives/schemas/overview`, `/cfg-mgmt/create-a-release`, and `/getting-started/quick-start/create-release` shows three annotation entries — config type, instance file path, instance format — and the instance-format entry's `<CodeGroup>` renders three tabs whose bodies are `x-miru-instance-format: "…"`, `@miru(instance_format="…")`, and `instance_format: "…"` respectively. Additionally, the `## Properties` section of `/cfg-mgmt/primitives/schemas/overview` lists an `instance format` property alongside the existing `instance file path` property.
4. The accepted-value lists agree **per schema language**: `docs/snippets/references/cli/releases/create/schema-annotations.mdx` and `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` both give `json`, `yaml` for JSON Schema and CUE, and `json`, `yaml`, `xml`, `text` for Opaque. The Opaque instance-format cell in `docs/snippets/schemas/formats.mdx` names the same four formats as `JSON, YAML, XML, Text`. `grep -rn jsonc docs/ --include='*.mdx'` prints nothing. (Do **not** grep all of `docs/`: the vendored OpenAPI spec `docs/references/platform-api/2026-05-06.yaml:1651` lists `jsonc` in an enum. That file is generated and out of scope.)
5. `git status --short` is empty — every edit, including this plan's own living-document updates, has been committed.
6. No file under `tools/lint/` was modified (verify with `git diff --stat main...HEAD -- tools/lint`, which must print nothing).
7. `cat -A docs/snippets/references/cli/releases/create/schema-annotations.mdx | grep -c '^  \$'` prints `6`.

## Idempotence and Recovery

Every step is a file edit or a read-only local command and can be repeated safely. No step writes to the deployed backend or to any remote state except the single `git push --force-with-lease` in Milestone 4. `pnpm run lint`, `pnpm run test:lint`, `./scripts/preflight.sh`, and `pnpm dev` have no side effects beyond building the Go linter binary.

If a milestone's edits go wrong before committing, `git checkout -- <path>` restores the file. After committing, `git revert <sha>` or `git reset --soft HEAD~1` (to re-edit) both work; prefer `reset --soft` since nothing has been pushed until Milestone 4.

The only risky step is `git push --force-with-lease`. It is safe by construction: it aborts rather than overwriting if the remote moved. If it aborts, do not add `--force`; fetch, inspect what changed on the remote, and reconcile first. If a bad push does land, the previous remote head is recoverable from `git reflog` and from PR #138's force-push history on GitHub.

If `pnpm dev` fails to start, check that `pnpm install` has been run in `/home/ben/miru/workbench2/repos/docs`; the render check can be repeated any number of times.

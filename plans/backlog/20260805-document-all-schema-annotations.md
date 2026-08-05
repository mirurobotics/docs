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

Miru users write a **config schema** (a file describing the shape of a device configuration) and annotate it with **schema annotations** — special keys that tell Miru which config type the schema belongs to and how instances of it are deployed. The docs snippet at `docs/snippets/references/cli/releases/create/schema-annotations.mdx` is the single source rendered on three different pages, and it currently documents only two of the three annotations the CLI actually supports. The third, **instance format** (`x-miru-instance-format` / `@miru(instance_format=…)` / `instance_format`), is invisible to readers of those three pages even though it shipped in stable CLI v0.10.1.

After this change, a reader on any of the three pages sees all three annotations, each with JSON Schema, CUE, and Opaque syntax, with an accurate list of accepted values per schema language. Additionally, three currently-inconsistent statements about opaque instance formats (the snippet, `opaque.mdx`, and an uncommitted edit to `formats.mdx`) are reconciled so the docs do not contradict themselves. Success is visible by running `pnpm dev` and seeing a third annotation block on `/cfg-mgmt/primitives/schemas/overview`, and by CI going green on PR #138.

## Progress

- [ ] Milestone 1 — Reconcile the uncommitted `formats.mdx` edit against `opaque.mdx`; commit.
- [ ] Milestone 2 — Resolve the `jsonc` end-to-end question and fix the value list to document.
- [ ] Milestone 3 — Add the instance-format annotation to the snippet, align `opaque.mdx` and `overview.mdx`, add cspell words if needed; commit.
- [ ] Milestone 4 — Lint, render check, force-push, verify CI green on PR #138.

## Surprises & Discoveries

Add entries as work proceeds.

## Decision Log

Add entries as work proceeds.

## Outcomes & Retrospective

Fill in at completion.

## Context and Orientation

**Terms.** A *config schema* is a file (JSON Schema, CUE, or Opaque) that Miru stores and, for the first two languages, uses to validate config instances. A *config instance* is the concrete configuration file written onto a device. A *schema annotation* is a key embedded in the schema file that carries Miru metadata rather than validation logic. A *schema language* is which of the three formats the schema is written in. An *Opaque schema* is metadata-only: it carries annotations but validates nothing.

**Repository.** `repos/docs` is a Mintlify docs site. Content lives under `docs/`, reusable fragments under `docs/snippets/`, and plans under `plans/{backlog,active,completed}/`. There is no `CLAUDE.md` or `AGENTS.md`; conventions are enforced entirely by tooling (see "Lint rules" below).

**Git state at the time of writing.** Branch `docs/opaque-schema-language` (PR #138) has just been rebased onto `main` and has **not** been force-pushed. The working tree has one uncommitted modification: `docs/snippets/schemas/formats.mdx`. That edit must be dealt with in Milestone 1 — do not leave the tree dirty.

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

**1. Config type — required, no default. Already documented.** Must sit at the root of the schema document (JSON Schema) or on a schema-level `@miru()` declaration attribute in CUE; field-level attributes are ignored (`core@v0.8.2/pkg/schemas/cue/attributes.go:26-33,76-83`). Missing ⇒ error `CfgTypeSlugNotFound`.

**2. Instance filepath — optional. Already documented.** Absolute path string. The default is applied **server-side**: `/srv/miru/configs/<config-type-slug>.<ext>` where `ext` is `json` for JSON- and CUE-formatted schemas and `yaml` for YAML-formatted schemas (`repos/backend/internal/configs/domain/config_schemas/filepath.go:7-22`). Server validation (`repos/backend/internal/configs/domain/config_instances/filepath.go:17-41`): max 4096 characters, must be absolute, extension must be `.json`/`.yaml`/`.yml` for JSON Schema and CUE; opaque schemas may use any extension.

**3. Instance format — optional. UNDOCUMENTED in the snippet. This is the gap this plan closes.** Accepted values differ per schema language:

- JSON Schema: `json`, `yaml`, `jsonc` — `core@v0.8.2/pkg/schemas/jsonschema/format.go:9-17`
- CUE: `json`, `yaml`, `jsonc` — `core@v0.8.2/pkg/schemas/cue/format.go:9-17`
- Opaque: `json`, `yaml`, `jsonc`, `xml`, `text` — `core@v0.8.2/pkg/schemas/opaque/format.go:16-27`

Parsing sites: `core@v0.8.2/pkg/schemas/jsonschema/annotations.go:62-75`, `…/cue/attributes.go:118-128`, `…/opaque/compile.go:139-145`. Unknown value produces e.g. `unknown JSON Schema instance format '<value>', must be json, yaml or jsonc` (analogous wording for CUE and Opaque). The default is inferred **server-side** from the instance filepath's extension (`resolveInstFormat`, `repos/backend/internal/configs/services/config_schemas/create.go:384-401`); an unmappable extension yields `CannotInferInstFormat`. Because `instance_filepath` itself defaults to `…/<slug>.json`, the effective default is `json`. This annotation shipped in stable CLI v0.10.1 (commit `f17040f`), so it is safe to document independent of the opaque beta.

**Open question that Milestone 2 must resolve before any value list is written.** The backend's server-side per-language allow-list at `repos/backend/internal/configs/domain/config_instances/formats.go:12-34` permits `{json, yaml, xml, text}` for opaque and only `{json, yaml}` for every other language — it **omits `jsonc`**, even though the CLI and `core` accept `jsonc` and CLI test fixtures use it. Do not silently document `jsonc` as working end-to-end.

**4. Opaque `language` marker — already documented, no gap.** The CLI classifies a `.json`/`.yaml` file as opaque only when the root key `language == "opaque"` (`repos/cli-private/internal/domain/cfgschs/load.go:249-282`); otherwise it parses the file as JSON Schema. Unknown root keys are rejected (`unknown key "<key>" in opaque schema`). This is covered by the key table on `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`. Future-watch only: unreleased `core` (`repos/core/pkg/schemas/opaque/compile.go:14-33`) renames `language` → `schema_language` with `language` kept as a deprecated alias, but the CLI on `core` v0.8.2 does **not** accept `schema_language`. Document `language` only; change nothing for this.

### Current state of the snippet

`docs/snippets/references/cli/releases/create/schema-annotations.mdx` holds two `<ParamField>` blocks — `path="config type" required` and `path="instance file path"` — each containing one `<CodeGroup>` with three fenced tabs in this order: ` ```yaml JSON Schema `, ` ```cue CUE `, ` ```yaml Opaque `.

Formatting convention that any new `<CodeGroup>` must match exactly: the block is indented two spaces inside the `<ParamField>`; a blank line follows `<CodeGroup>`; between consecutive fences there are **two** separator lines — one truly blank line followed by a line containing exactly two trailing spaces; one blank line precedes `</CodeGroup>`. Copy an existing block and edit it rather than typing a new one, so the trailing-space lines survive.

### Files that render the snippet

1. `docs/cfg-mgmt/primitives/schemas/overview.mdx` — imports it as `Annotations` (line 9), renders under `## Schema annotations`. This page **also duplicates** the instance-filepath prose in its own `<ParamField path="instance file path" type="string">` around line 67, in the `## Properties` section.
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

**Milestone 1 — reconcile `formats.mdx`.** The working tree changes the Opaque row of `docs/snippets/schemas/formats.mdx` from `| Opaque | JSON, YAML | JSON, YAML, XML |` to `| Opaque | YAML | JSON, YAML, XML, Plain Text |`. Both halves of that edit contradict committed content: `opaque.mdx` states plainly that opaque schemas are written as JSON **or** YAML, and its instance-format list is JSON / YAML / XML with no plain text. The schema-format column ("JSON, YAML") is correct per `opaque.mdx`, so restore it. The instance-format column is the thing Milestone 2 decides. Restore the schema-format half now and leave the instance-format half to be set once, consistently, in Milestone 3 — but do not leave the file dirty across milestones: commit the schema-format correction in Milestone 1.

**Milestone 2 — resolve `jsonc`.** Determine whether `jsonc` is accepted end to end (CLI accepts it; the backend allow-list at `repos/backend/internal/configs/domain/config_instances/formats.go:12-34` appears not to). Verify against the deployed backend or ask the backend owner. Two branches:

- (a) `jsonc` is genuinely accepted end to end → document `json`, `yaml`, `jsonc` for JSON Schema and CUE, and `json`, `yaml`, `jsonc`, `xml`, `text` for Opaque.
- (b) `jsonc` is rejected server-side → document `json`, `yaml` for JSON Schema and CUE, and `json`, `yaml`, `xml`, `text` for Opaque, and file a follow-up issue against the backend (or `core`) noting the CLI/backend divergence.

Whichever branch is taken, record it in the Decision Log with the evidence. Also decide `text` at the same time: the backend allow-list includes `text` for opaque and `opaque.mdx` currently omits it. If `text` is accepted end to end, add it; if not, leave it out and note why.

**Milestone 3 — write the docs.** Three edits, all consistent with the value list chosen in Milestone 2:

- `docs/snippets/references/cli/releases/create/schema-annotations.mdx`: append a third `<ParamField path="instance format">` after the instance-file-path block. Prose: what the annotation does (declares the format of the deployed config instance), that it is optional, that it defaults to the format inferred from the instance file path's extension (so effectively `json` under the default path), the accepted values, and a sentence noting that Opaque schemas accept a wider set. Then a `<CodeGroup>` with the same three tabs in the same order and the same whitespace convention:

      ```yaml JSON Schema
      x-miru-instance-format: "yaml"
      ```

      ```cue CUE
      @miru(instance_format="yaml")
      ```

      ```yaml Opaque
      instance_format: "xml"
      ```

  Close with an `Examples:` line matching the style of the config-type block. Add no headings.

- `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`: bring the `instance_format` row of the key table and the "Opaque schemas accept the widest set of config instance formats" bullet list into agreement with the chosen value list.

- `docs/snippets/schemas/formats.mdx`: set the Opaque instance-format cell to the chosen list; leave the schema-format cell at `JSON, YAML`.

- `docs/cfg-mgmt/primitives/schemas/overview.mdx`: this page's `## Properties` section duplicates the instance-filepath field. Decide whether a sibling `<ParamField path="instance format" type="string">` belongs there for consistency. Recommendation: **add it**, short, mirroring the existing instance-file-path property (immutable badge, one-line description, default, examples), because the section enumerates schema properties and omitting one is a real gap. Record the decision either way.

- `cspell.json`: if `pnpm run lint` flags `jsonc`, `xml`, or `yaml`, add them to the `words` array in correct sorted position (the array is sorted uppercase-first then lowercase ascending).

**Milestone 4 — validate and ship.** Run the lint suite, run a `pnpm dev` render check of the three pages, force-push the rebased branch, and confirm PR #138 CI is green.

## Concrete Steps

All commands run from `/home/ben/miru/workbench2/repos/docs` unless stated otherwise. Commit from inside this repo — never from the workbench root.

### Milestone 1 — reconcile the uncommitted `formats.mdx` edit

    cd /home/ben/miru/workbench2/repos/docs
    git status --short
    git diff -- docs/snippets/schemas/formats.mdx

Expect `M docs/snippets/schemas/formats.mdx` and the diff described above. Edit the file so the Opaque row's schema-format column reads `JSON, YAML` again, matching `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`. Leave the instance-format column as it currently stands in the tree for now; Milestone 3 sets it. Then:

    git add docs/snippets/schemas/formats.mdx
    git commit -m "docs(schemas): correct opaque schema formats row"

Expect one file changed. If the diff turns out to be empty after your edit (i.e. you reverted the whole thing), use `git checkout -- docs/snippets/schemas/formats.mdx` instead and skip the commit, noting that in the Decision Log.

### Milestone 2 — resolve the `jsonc` question

Re-read the evidence:

    sed -n '1,40p' /home/ben/miru/workbench2/repos/backend/internal/configs/domain/config_instances/formats.go
    sed -n '1,30p' /home/ben/go/pkg/mod/github.com/mirurobotics/core@v0.8.2/pkg/schemas/opaque/format.go

Then confirm real behavior. Preferred: create a throwaway schema with `x-miru-instance-format: "jsonc"` and an instance filepath ending in `.json`, and attempt `miru release create` against the deployed backend, observing whether the API accepts it. If that is not practical, ask the backend owner directly and record the answer. Do not guess.

Record the outcome in the Decision Log with the date and the evidence (command output or the person who answered), and write down the exact value lists to use for JSON Schema / CUE and for Opaque. No commit in this milestone.

### Milestone 3 — write the docs

Edit, in order:

1. `docs/snippets/references/cli/releases/create/schema-annotations.mdx` — append the `instance format` `<ParamField>` after the existing `instance file path` block (currently ends at line 63). Copy the existing `<CodeGroup>` block verbatim first, then edit the three fence bodies, so the two-space indentation and the two-trailing-space separator lines are preserved.
2. `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` — update the `instance_format` table row and the instance-format bullet list.
3. `docs/snippets/schemas/formats.mdx` — set the Opaque instance-format cell.
4. `docs/cfg-mgmt/primitives/schemas/overview.mdx` — add the sibling `instance format` property field (if that was the decision).

Then check spelling and lint quickly:

    pnpm run lint

If cspell reports `Unknown word (jsonc)` (or `xml` / `yaml`), add the word to the `words` array in `cspell.json` in sorted position and re-run. If `heading-case` fires, you added a heading — remove it; `<ParamField>` labels are not headings and no `tools/lint` change should be needed. If `no-double-dash` fires, replace the bare `--` with an em dash.

Expect a clean final run, ending with the openapi-check phase passing. Then:

    git add -A
    git commit -m "docs(schemas): document the instance format schema annotation"

### Milestone 4 — validate, push, verify CI

    pnpm run test:lint
    pnpm run lint

Both must exit 0.

Render check:

    pnpm dev

Open each of these local URLs (the dev server prints its port, typically `http://localhost:3000`) and confirm three annotation blocks render, each with three working tabs labelled JSON Schema / CUE / Opaque:

- `http://localhost:3000/cfg-mgmt/primitives/schemas/overview` (section "Schema annotations")
- `http://localhost:3000/cfg-mgmt/create-a-release` (section "Schema annotations")
- `http://localhost:3000/getting-started/quick-start/create-release` (section "Schema annotations")

Stop the dev server with Ctrl-C.

Full preflight:

    ./scripts/preflight.sh

This runs lint smoke tests, the Go lint and coverage gates for `tools/lint`, `./scripts/lint.sh`, `./scripts/audit.sh`, and the bats shell tests. Expect everything to pass except possibly the pre-existing `audit` failure recorded in `plans/completed/20260726-opaque-schema-language.md`. If `audit` fails, confirm the failure is byte-for-byte the pre-existing one (same URL/target) and not something these edits introduced; note it in Surprises & Discoveries and continue. Any other failure must be fixed before proceeding.

Push the rebased branch and check CI:

    git push --force-with-lease
    gh pr checks 138 --watch

`--force-with-lease` is required because the branch was rebased; it refuses to push if someone else moved the remote branch, which is the safe behavior. If it is rejected, run `git fetch origin` and inspect `git log --oneline origin/docs/opaque-schema-language..HEAD` before deciding whether to retry.

Expect `gh pr checks 138` to report all checks passing (modulo the known `audit` failure, if it is still present on `main` too). `lint-custom-linter` and `test-custom-linter` should not run at all, since `tools/lint/**` is untouched.

## Validation and Acceptance

Acceptance is behavioral, and **preflight must report CLEAN before this task may be reported complete**. Concretely, all of the following must hold:

1. `pnpm run test:lint` and `pnpm run lint` both exit 0 from `/home/ben/miru/workbench2/repos/docs`.
2. `/home/ben/miru/workbench2/repos/docs/scripts/preflight.sh` completes with no failures other than the pre-existing `audit` failure documented in `plans/completed/20260726-opaque-schema-language.md`. Any new failure blocks completion.
3. `gh pr checks 138` reports green on the **pushed branch head** — that is, after `git push --force-with-lease`, not on a stale remote commit. A local pass with an unpushed or outdated remote head does not satisfy this.
4. Under `pnpm dev`, each of `/cfg-mgmt/primitives/schemas/overview`, `/cfg-mgmt/create-a-release`, and `/getting-started/quick-start/create-release` shows three annotation entries — config type, instance file path, instance format — and the instance-format entry's `<CodeGroup>` renders three tabs whose bodies are `x-miru-instance-format: "…"`, `@miru(instance_format="…")`, and `instance_format: "…"` respectively.
5. The accepted-value list for `instance_format` is identical in `docs/snippets/references/cli/releases/create/schema-annotations.mdx`, `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`, and `docs/snippets/schemas/formats.mdx`, and matches the Milestone 2 decision.
6. `git status --short` is empty — no leftover uncommitted edits.
7. No file under `tools/lint/` was modified (verify with `git diff --stat main...HEAD -- tools/lint`, which must print nothing).

## Idempotence and Recovery

Every step is a file edit or a read-only command and can be repeated safely. `pnpm run lint`, `pnpm run test:lint`, `./scripts/preflight.sh`, and `pnpm dev` have no side effects beyond building the Go linter binary.

If a milestone's edits go wrong before committing, `git checkout -- <path>` restores the file. After committing, `git revert <sha>` or `git reset --soft HEAD~1` (to re-edit) both work; prefer `reset --soft` since nothing has been pushed until Milestone 4.

The only risky step is `git push --force-with-lease`. It is safe by construction: it aborts rather than overwriting if the remote moved. If it aborts, do not add `--force`; fetch, inspect what changed on the remote, and reconcile first. If a bad push does land, the previous remote head is recoverable from `git reflog` and from PR #138's force-push history on GitHub.

If `pnpm dev` fails to start, check that `pnpm install` has been run in `/home/ben/miru/workbench2/repos/docs`; the render check can be repeated any number of times.

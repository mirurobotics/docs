# Correct the documented config-instance file formats

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/ben/miru/workbench3/repos/docs`) | read-write | All MDX edits and the lint/validation run. |
| `backend/` (`/home/ben/miru/workbench3/repos/backend`) | read-only | Source of truth for which instance formats the API accepts. |
| `core/` (`/home/ben/miru/workbench3/repos/core`) | read-only | Source of truth for the instance-format enum and extension inference. |
| `agent/` (`/home/ben/miru/workbench3/repos/agent`) | read-only | Source of truth for what the Miru Agent does with instance content on the device. |

This plan lives in `docs/plans/backlog/` because every file that changes is in the docs repo.

## Purpose / Big Picture

A reader of https://docs.mirurobotics.com/cfg-mgmt/primitives/config-instances currently learns that config instances "support JSON and YAML" and that XML support is "coming soon". Both statements are false: XML and plain text shipped in CLI v0.10.3 (July 26, 2026) alongside opaque schemas. The same page also claims the Miru Agent parses config instances at deploy time, which it does not — it writes the file byte-for-byte.

After this change, the config-instances page states the real, per-schema-language format list (JSON and YAML for JSON Schema and CUE; JSON, YAML, XML and plain text for opaque schemas), names the one non-obvious caveat (`text` must be declared explicitly because it cannot be inferred from a file extension), and describes deploy-time behavior accurately. The shared definition snippet stops hard-coding "JSON or YAML". Nothing in `docs/` contradicts `backend/internal/configs/domain/config_instances/formats.go` any more, except the two vendored OpenAPI specs, which are generated and explicitly out of scope.

## Progress

- [ ] Milestone 1 — re-verify ground truth against `backend/` and `core/`; record findings in the Decision Log.
- [ ] Milestone 2 — rewrite the `## File formats` section and the `content` property of `docs/cfg-mgmt/primitives/config-instances.mdx`; commit.
- [ ] Milestone 3 — broaden `docs/snippets/definitions/config-instance.mdx`; commit.
- [ ] Milestone 4 — run the full audit greps and `./scripts/preflight.sh`; fix anything it flags; commit.
- [ ] Milestone 5 — push, open a draft PR, wait for CI green, then mark the PR ready.

## Surprises & Discoveries

(Add entries as you go.)

## Decision Log

(Add entries as you go. Milestone 1 requires an entry.)

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

**Config instance** — a set of configuration values Miru deploys to a device as a single file. **Config schema** — the document that describes and validates config instances for a config type. **Schema language** — the language the schema is written in: JSON Schema, CUE, or *opaque*. An **opaque schema** is a metadata-only schema that treats every instance as valid; it exists so you can ship config formats Miru cannot validate. **Instance format** — the on-disk format of the deployed instance file, set by the `instance_format` schema annotation.

Ground truth for the accepted formats, all verified on 2026-08-05:

- `/home/ben/miru/workbench3/repos/backend/internal/configs/domain/config_instances/formats.go` — `SupportedFormats(lang schemas.Language)` returns `{Json, Yaml, Xml, Text}` when `lang == schemas.OpaqueLang` and `{Json, Yaml}` otherwise. `VerifyFormatSupported` rejects anything else with `config_instance_format_not_supported`.
- `/home/ben/miru/workbench3/repos/core/pkg/schemas/format.go:70-74` — the union enum is `json`, `yaml`, `jsonc`, `xml`, `text`. `jsonc` is accepted by core but rejected by the backend for every language, so **it must not be documented anywhere in `docs/`**.
- `/home/ben/miru/workbench3/repos/core/pkg/schemas/format.go` `InstFormatFromFileType` handles JSON, YAML, JSONC and XML file types only, and `/home/ben/miru/workbench3/repos/core/pkg/filesys/file.go:15-20` defines no text file type. Combined with `resolveInstFormat` in `/home/ben/miru/workbench3/repos/backend/internal/configs/services/config_schemas/create.go:384-401`, this means an instance file path ending in `.conf` or `.txt` cannot have its format inferred — the schema must declare `instance_format: text` explicitly or schema creation fails with `CannotInferInstFormat`. XML is inferred fine from `.xml`.
- `/home/ben/miru/workbench3/repos/agent/agent/src/http/config_instances.rs` fetches instance content as an opaque `String`; `/home/ben/miru/workbench3/repos/agent/agent/src/deploy/filesys.rs` writes it verbatim to the instance file path. The agent never parses the content, so no minimum agent version is tied to XML or text specifically.

Pages already correct on `main` (do **not** change them; match their wording):

- `docs/snippets/schemas/formats.mdx` — table whose Opaque row reads `| Opaque | JSON, YAML | JSON, YAML, XML, Text |`.
- `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` — instance formats listed as JSON, YAML, XML, Text.
- `docs/cfg-mgmt/primitives/schemas/overview.mdx:74-85` — the `instance format` property, which already says JSON Schema and CUE accept `json` and `yaml` while opaque accepts `json`, `yaml`, `xml` and `text`.
- `docs/snippets/references/cli/releases/create/schema-annotations.mdx:75-100` — the same statement for the CLI annotation reference.
- `docs/changelog/cli.mdx` (v0.10.3) and `docs/changelog/product.mdx` (August 6, 2026) — both already announce XML and plain text. Changelogs are historical records; never rewrite them.

The stale files, and the only ones this plan edits:

1. `docs/cfg-mgmt/primitives/config-instances.mdx` — line 18 ("deployed to the device as a JSON or YAML file"), line 39 ("Currently, config instances support JSON and YAML. However, support for more formats, including XML, is coming soon!") and line 41 ("Config instances support JSON and YAML. The Miru Agent parses both formats at deploy time…").
2. `docs/snippets/definitions/config-instance.mdx` — "Config instances are stored as text files (JSON or YAML)". This snippet is rendered on the config-instances page and, via `docs/snippets/concepts/config-instance.mdx`, on `docs/getting-started/quick-start/overview.mdx`.

Explicitly out of scope, with reasons recorded so a reviewer does not re-raise them:

- `docs/references/platform-api/2026-05-06.yaml:1645` and `docs/references/platform-api/2026-03-09.yaml:1592` — the `InstanceFormat` enum is `json|yaml|jsonc`. These files are vendored copies of `backend/api/specs/platform/v20260506.yaml:1271`, which has the same gap. Hand-editing them would be overwritten on the next spec pull. Raise a follow-up against the backend repo instead.
- `docs/snippets/agent/yaml-support.mdx` says YAML support requires Miru Agent **v0.7.0**, while `docs/changelog/agent.mdx:96-100` credits YAML support to **v0.7.1**. This is a genuine contradiction but the correct version could not be established from the agent repo during research, and it concerns agent versioning rather than the format list. Leave it; open a follow-up issue.

Tooling. `./scripts/lint.sh` runs a custom Go MDX prose linter (`tools/lint/`), ESLint over MDX, CSpell, and `mint openapi-check`. Relevant prose rules: `no-double-dash` (use the em dash `—`, never `--`), `heading-case`, and import rules requiring `.mdx` imports to be default imports ending in `;`. `./scripts/preflight.sh` runs everything CI runs: `pnpm run test:lint`, the Go linter and coverage gate for `tools/lint`, `./scripts/lint.sh`, `./scripts/audit.sh`, and the bats suite. `.github/workflows/ci.yml` runs the same set.

## Plan of Work

**Milestone 1 — re-verify ground truth.** Read `backend/internal/configs/domain/config_instances/formats.go` and `core/pkg/schemas/format.go` and confirm they still say what Context claims. If either has changed, stop and re-derive the format list before editing any docs. Record the confirmed lists in the Decision Log with today's date and file:line citations. No file changes, no commit.

**Milestone 2 — `docs/cfg-mgmt/primitives/config-instances.mdx`.** Three edits in one commit.

First, the `content` property (currently line 18). Replace "The configuration values that are deployed to the device as a JSON or YAML file." with wording that does not name formats, e.g. "The configuration values that are deployed to the device as a file, written in the config schema's instance format."

Second, add an import for the shared formats table alongside the existing imports at the top of the file, keeping imports sorted as the `import-sorted` lint rule requires:

    import SchemaFormats from '/snippets/schemas/formats.mdx';

Reusing `docs/snippets/schemas/formats.mdx` rather than writing a fourth copy of the format list is the point: that table is already the single source of truth and is already rendered on the schemas overview.

Third, replace the whole `## File formats` section body (currently the two paragraphs at lines 39 and 41, plus the trailing `<AgentYamlSupport />`) with, in order:

- One sentence stating that the formats a config instance may use are determined by the schema language of its config schema, followed by `<SchemaFormats />`.
- One sentence naming the practical upshot: JSON Schema and CUE schemas accept JSON and YAML; [opaque](/cfg-mgmt/primitives/schemas/languages/opaque) schemas additionally accept XML and plain text, because they do not validate instance content.
- A `<Note>` capturing the inference caveat: the instance format is normally inferred from the instance file path's extension, but plain text has no recognized extension, so a text instance must declare `instance_format: text` in its opaque schema or schema creation fails.
- A corrected deploy-time sentence: the Miru Agent writes the config instance to its instance file path exactly as authored, byte for byte; on-device applications read and parse the file themselves. Do not claim the agent parses anything.
- Keep `<AgentYamlSupport />` at the end of the section — the YAML minimum-version note is still relevant and unchanged.

Do not use `--` anywhere; use `—`.

**Milestone 3 — `docs/snippets/definitions/config-instance.mdx`.** This one-sentence snippet is reused in the quick start, so keep it short and format-agnostic. Change "Config instances are stored as text files (JSON or YAML), which applications parse into a structured format for consumption." to drop the hard-coded pair, for example "Config instances are stored as text files — JSON, YAML, XML, or plain text — which applications parse into a structured format for consumption." Preserve the file's existing trailing whitespace so the diff stays to one line.

**Milestone 4 — audit and validate.** Re-run the audit greps to prove nothing stale remains, then run preflight. If CSpell flags a new word, add it to the `words` array in `cspell.json` in sorted position (uppercase entries first, then lowercase ascending). No new words are expected.

**Milestone 5 — ship.** Push the branch, open a draft PR, wait for CI, mark ready only when green.

## Concrete Steps

All commands run from `/home/ben/miru/workbench3/repos/docs` unless stated otherwise. The working branch `docs/config-instance-file-formats` already exists and is even with `main`; confirm with `git status --short && git log --oneline main..HEAD` (expect empty output for both before Milestone 2).

Milestone 1:

    grep -n "SupportedFormats" -A 12 /home/ben/miru/workbench3/repos/backend/internal/configs/domain/config_instances/formats.go
    grep -n "InstFormatText\|InstFormatFromFileType" -A 12 /home/ben/miru/workbench3/repos/core/pkg/schemas/format.go

Expect the backend function to return `{Json, Yaml, Xml, Text}` for `OpaqueLang` and `{Json, Yaml}` otherwise, and `InstFormatFromFileType` to have no `Text` case. Add a Decision Log entry recording both. Nothing to commit.

Milestone 2 — edit `docs/cfg-mgmt/primitives/config-instances.mdx`, then:

    ./scripts/lint.sh
    git add docs/cfg-mgmt/primitives/config-instances.mdx
    git commit -m "docs(config-instances): document XML and plain text instance formats"

`./scripts/lint.sh` ends with `All documentation lint checks passed.` on success.

Milestone 3 — edit `docs/snippets/definitions/config-instance.mdx`, then:

    ./scripts/lint.sh
    git add docs/snippets/definitions/config-instance.mdx
    git commit -m "docs(snippets): drop JSON/YAML-only wording from the config instance definition"

Milestone 4 — audit greps. Each must produce the stated result:

    grep -rn "coming soon" docs --include='*.mdx' | grep -i "format"

Expect no output.

    grep -rn "jsonc" docs --include='*.mdx'

Expect no output. (Do not grep all of `docs/` — the vendored OpenAPI specs legitimately contain `jsonc`.)

    grep -rniE "config instances support|stored as text files" docs --include='*.mdx'

Expect only the rewritten lines, with no bare "JSON and YAML" pairing.

    grep -rn "JSON, YAML, XML, Text" docs/snippets/schemas/formats.mdx

Expect the unchanged Opaque row.

Then the full preflight:

    ./scripts/preflight.sh

Expect the sections `=== Lint Smoke Tests ===`, `=== Go Lint (tools/lint) ===`, `=== Go Coverage (tools/lint) ===`, `=== Lint ===`, `=== Audit ===`, `=== Shell Script Tests ===` to run in that order and the command to exit 0. Commit only if this milestone produced changes (e.g. a `cspell.json` addition):

    git add -A
    git commit -m "chore(docs): satisfy lint after instance format updates"

Milestone 5:

    git push -u origin docs/config-instance-file-formats
    gh pr create --draft --base main --title "docs: document XML and plain text config instance formats" --body "<summary + the out-of-scope notes from Context>"
    gh pr checks --watch

When all checks pass:

    gh pr ready

Then open the two follow-ups named in Context (backend OpenAPI `InstanceFormat` enum; the v0.7.0 vs v0.7.1 YAML-support discrepancy) as separate issues.

## Validation and Acceptance

**Preflight must report CLEAN before this task is reported complete and before the PR leaves draft.** Concretely: `./scripts/preflight.sh` exits 0 locally with no failing section, the branch head is pushed, and `gh pr checks` reports every GitHub Actions check on that exact pushed head as passing (CI green). A local-only pass is not sufficient; a green run on an older commit is not sufficient. If any check fails, fix it, push, and re-verify before marking the PR ready.

Behavioral acceptance, verifiable by a human:

1. Run `pnpm run dev` from the repo root (it `cd`s into `docs/` and starts `mint dev`) and open `/cfg-mgmt/primitives/config-instances`. The `## File formats` section renders the three-row schema-format table with the Opaque row reading `JSON, YAML, XML, Text`, states that opaque schemas additionally accept XML and plain text, carries the `<Note>` about `instance_format: text` needing to be declared explicitly, and contains no occurrence of the words "coming soon". The `content` property no longer says "JSON or YAML file".
2. The same page's intro paragraph (from `ConfigInstanceDef`) no longer restricts config instances to JSON or YAML, and `/getting-started/quick-start/overview` reflects the same updated sentence, proving the shared snippet was edited rather than the page.
3. No page claims the Miru Agent parses config instance content; the page states the agent writes the file verbatim.
4. `grep -rn "jsonc" docs --include='*.mdx'` prints nothing, and `grep -rn "coming soon" docs --include='*.mdx' | grep -i format` prints nothing.
5. The documented lists agree with `backend/internal/configs/domain/config_instances/formats.go`: `json`, `yaml` for JSON Schema and CUE; `json`, `yaml`, `xml`, `text` for opaque; `jsonc` documented nowhere.

There are no unit tests for prose. The repo's test surface for these changes is `pnpm run test:lint` (the `tests/test-lint.sh` smoke suite over `tests/lint-fixtures`) plus the Go tests and coverage gate for `tools/lint`, all of which `./scripts/preflight.sh` invokes. Since no linter code changes here, those suites should pass identically before and after; their role is to prove the MDX edits introduce no rule violations.

## Idempotence and Recovery

Every step is a text edit plus a re-run of read-only checks, so all of them are safely repeatable. `./scripts/lint.sh`, `./scripts/preflight.sh` and the audit greps have no side effects beyond building `tools/lint/lint` (gitignored).

To undo a milestone before pushing: `git restore docs/<file>` for uncommitted work, or `git revert <sha>` for a committed milestone. To reset the branch entirely: `git reset --hard origin/main` (destroys uncommitted work — check `git status --short` first). After pushing, prefer `git revert` plus a new push over force-pushing, so the PR history stays reviewable.

The one non-obvious risk is editing `docs/snippets/definitions/config-instance.mdx`, since it renders on two pages. If the new sentence reads badly in the quick start, revert that single commit and instead override the wording only on the config-instances page; the snippet change is not load-bearing for the rest of the plan.

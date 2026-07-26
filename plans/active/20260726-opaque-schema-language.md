# Document the opaque schema language

This ExecPlan is a living document. Update Progress, Surprises & Discoveries, and the Decision Log as work proceeds; fill in Outcomes & Retrospective on completion.

## Scope

| Repository | Access | Branch |
|---|---|---|
| mirurobotics/docs (this repo; all paths below are repo-relative) | read/write | `docs/opaque-schema-language` (already checked out and clean; base `main`) |
| mirurobotics/core | read-only reference | `main` — source of truth for opaque compile rules, formats, validation |
| mirurobotics/backend | read-only reference | `main` — source of truth for which instance formats and filepaths the API actually accepts |

## Purpose / Big Picture

Miru supports a third config schema language, **Opaque**, and the docs do not mention it. An opaque schema is metadata-only: it attaches a config type, an optional instance filepath, and an optional instance format to a config type **without describing or constraining the instance content at all** (named after Kubernetes `type: Opaque`). It exists so a config Miru cannot meaningfully validate — a vendor text format, an XML document, a device-specific `.conf` — can still be versioned, released, rolled out, and audited like any other Miru config.

This plan adds one new language page and sweeps every place in the corpus that hard-codes the two-language list, so a reader never sees a stale "Miru supports JSON Schema and CUE" claim.

Observable outcome: on `pnpm dev`, `/cfg-mgmt/primitives/schemas/languages/opaque` renders as a third entry under **Languages**; the languages overview, the schemas overview (`language` allowed values, "Schema languages" section), the shared file-formats table, and the shared schema-annotations snippet all list Opaque; `pnpm run lint`, `pnpm run test:lint`, and CI are green.

## Progress

- [x] M1: new page `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`; nav entry in `docs/docs.json`. Commit `5082478`.
- [x] M2: languages overview bullet list; schemas overview (`language` ParamField, "Schema languages" section, instance-filepath prose). Commit `2b5e076`.
- [x] M3: shared snippets — `snippets/schemas/formats.mdx` table, `snippets/references/cli/releases/create/schema-annotations.mdx` Opaque tabs, `snippets/schema-langs/support.mdx` one-liner. Commit `04668c2`.
- [x] M4: full lint sweep (`pnpm run lint`, `pnpm run test:lint`) green; **no cspell additions needed**; `pnpm dev` render check passed. No content changes, so no M4 commit.
- [ ] Validation: push branch, preflight reports CLEAN (CI green on the pushed branch head), PR leaves draft.

## Surprises & Discoveries

- **CSpell needed no additions.** `nginx` and `passthrough` were both anticipated as allowlist candidates. `nginx` passes as-is (it appears only inside the inline code span `/srv/miru/nginx.conf`, and CSpell already tolerates it); the word "passthrough" was avoided entirely by writing the Validation section in plain prose. `cspell.json` is unmodified.
- **A fourth two-language claim exists, outside the planned sweep:** `docs/snippets/definitions/config-schema-examples.mdx:1` — "Below are some example schema definitions in the JSON Schema and CUE schema languages." This is **not** a claim about which languages Miru supports; it is a caption describing the two `<CodeGroup>` tabs immediately beneath it, both of which are real, constraint-bearing schemas. Adding an Opaque tab there would be off-topic (the snippet's point is what a constraining schema looks like). Left untouched deliberately; recording it so a future reader does not mistake it for a missed sweep target.
- The two other surviving `JSON Schema and CUE` greps are intentional: `opaque.mdx` itself contrasts against them when explaining the file-extension rule, and `schemas/overview.mdx` uses it to caption the two empty-schema examples (opaque has no "empty" form, so it does not belong there).
- **`<AgentYamlSupport />` placement.** Both other language pages put it directly after `## File formats`, but `cue.mdx` still has a later section (`## CUE packages`). The plan calls for it as the page's final line, so it sits after `## Instance file paths` here. Reads fine — the instance-file-paths section discusses YAML paths, so the Agent YAML note is still in context.
- **The `audit` CI job is red, and it is not this PR's fault.** `changes`, `lint`, and `shell-tests` all pass, and `lint-custom-linter` / `test-custom-linter` correctly report `skipping` (confirming `tools/lint/**` was untouched, as the acceptance criteria require). `audit` fails on 4 advisories in transitive dev dependencies: `postcss` and `tar` (via `mint`), `brace-expansion` (via `eslint`), plus one more high. Evidence that it is pre-existing:
  - `git diff main...HEAD --name-only` contains **no** `package.json` and **no** `pnpm-lock.yaml` — the dependency graph is byte-identical to `main`.
  - `scripts/audit.sh` is just `pnpm audit --ignore-registry-errors`, so its result is a function of the (unchanged) lockfile and the *current date* against the live advisory database.
  - Running `./scripts/audit.sh` in a pristine detached worktree of `main` reproduces the **same 4 vulnerabilities**. `main` last went green on 2026-07-24; these advisories were published since.
  - Fixing it means editing `package.json` (`auditConfig.ignoreCves` or dependency overrides) and/or the lockfile. That is outside this plan's declared file set, and suppressing CVEs is a security-policy call that should not be smuggled into a docs-only PR. Left for a separate dependency-maintenance change; the PR stays draft.

- `pnpm dev` **was** run and driven non-interactively (backgrounded `mint dev`, HTTP-probed each page). All five pages return 200 with no MDX compile error, the new page renders all five headings and both `<CodeGroup>` tabs, and the nav link to `/cfg-mgmt/primitives/schemas/languages/opaque` is present.

## Decision Log

- 2026-07-26 (user): **Document `opaque` as both the display name ("Opaque") and the wire value (`opaque`), everywhere.** The word "schemaless" must not appear anywhere in the docs. Rationale: core renamed `schemaless` → `opaque` (PR #129); backend and openapi still say `schemaless` and their rename is pending. The docs land ahead of those renames deliberately.
- 2026-07-26 (user): **No caveat `<Note>` about the pending backend/openapi rename.** Rationale: a transitional note would document an implementation detail that is about to disappear, and would reintroduce the banned word.
- 2026-07-26 (user): **Full sweep of hard-coded language lists, but NO changelog entry.** `docs/changelog/product.mdx` is not touched.
- 2026-07-26 (planning): **Document instance formats as JSON, YAML, XML only.** Core's opaque enum is wider (`json`, `yaml`, `jsonc`, `xml`, `other`), but `backend/internal/configs/domain/config_instances/formats.go` narrows the API to `json`, `yaml`, `xml` for this language (`jsonc` is excluded for every language; `other` is not accepted). Documenting `jsonc`/`other` would document values the API rejects with `config_instance_format_not_supported`. Document what the API accepts.
- 2026-07-26 (planning): **No `<Framed>` hero on the new page.** The other two language pages open with a `<Framed>` screenshot of the language's official website; Opaque is a Miru concept with no website and no logo. Any image would have to be hosted on `assets.mirurobotics.com` (the `image-domain` lint rule) and does not exist. Omit the component and its import.
- 2026-07-26 (planning): **Avoid "XML" in every heading; do not touch `tools/lint/**`.** The `heading-case` allowlist in `tools/lint/linter/headingcase/` does not contain `XML`, so a heading like `## XML instances` fails lint. XML appears only in prose, bullets, tables, and code blocks, where the rule does not apply. Adding `XML` to the allowlist would change `tools/lint/**` and trigger the `lint-custom-linter` and `test-custom-linter` CI jobs (and would want its own unit test) for zero reader benefit. If a later revision genuinely needs an XML heading, that is a separate, self-contained change.
- 2026-07-26 (planning): **No redirect entry in `docs.json`.** The `redirects` lint rule flags dead redirects; `opaque.mdx` is a brand-new URL with no predecessor.
- 2026-07-26 (planning): **Adding an Opaque tab to the shared `schema-annotations.mdx` snippet also changes the quick start.** That snippet is imported by `docs/cfg-mgmt/primitives/schemas/overview.mdx`, `docs/cfg-mgmt/create-a-release.mdx`, and `docs/getting-started/quick-start/create-release.mdx`. That is acceptable and desirable — a third `<CodeGroup>` tab is a low-noise addition. Separately, the standalone `<Tabs>` ("JSON Schema" / "CUE") in `docs/getting-started/quick-start/create-release.mdx` is **left alone**: the quick start is a guided happy path and a third branch there is noise.

## Outcomes & Retrospective

All four milestones landed as specified; the content plan needed no revision during execution. `pnpm run lint`, `pnpm run test:lint`, and the `pnpm dev` render check are green locally, and `cspell.json` required no additions.

Draft PR: [#138](https://github.com/mirurobotics/docs/pull/138).

CI status: `changes`, `lint`, `shell-tests` pass; `lint-custom-linter` and `test-custom-linter` correctly skip. `audit` fails on a pre-existing, repo-wide dependency-advisory problem reproduced on a pristine `main` worktree (see Surprises & Discoveries) — no file in this PR's diff touches the dependency graph. The plan's "CI green on the pushed branch head" gate is therefore **not** satisfied, so the PR remains in draft pending a separate dependency-maintenance change.

## Context and Orientation

### The subject matter (verified against core and backend)

**What an opaque schema is.** A metadata-only schema document. It declares annotations for a config type and imposes no structure on instance content. `core/pkg/schemas/opaque/validate.go` is the entire validation implementation:

    // Validate always succeeds because an opaque schema imposes no structure.
    func (s Schema) Validate(_ string) *errs.Error { return nil }

**The schema document** (`core/pkg/schemas/opaque/compile.go`). A single JSON or YAML mapping at the root, with exactly four permitted keys:

| Key | Required | Type | Allowed values |
|---|---|---|---|
| `config_type` | yes | string | the config type slug |
| `instance_filepath` | no | string | any absolute path |
| `instance_format` | no | string | `json`, `yaml`, `xml` (as accepted by the API) |
| `language` | no | string | must be exactly `opaque` if present |

Unknown keys are rejected (including `x-miru-*`), exactly one document is allowed (no `---` multi-document YAML), and the schema file itself must be JSON or YAML — CUE is rejected.

**Instance formats.** Core's enum is `json`, `yaml`, `jsonc`, `xml`, `other`; the backend narrows opaque to `json`, `yaml`, `xml`. Document the three. There is no default — omitting `instance_format` leaves it unset.

**Instance file paths — the differentiator.** `backend/internal/configs/domain/config_instances/filepath.go` skips the extension allowlist for this language: an opaque instance file may carry **any** extension (`.conf`, `.xml`, `.txt`), where JSON Schema and CUE are restricted to `.json`/`.yaml`/`.yml`. The absolute-path requirement still applies to every language.

**What is still enforced for opaque.** The schema document itself is validated at create time (valid JSON/YAML, single document, no unknown keys, `config_type` present and a string, `instance_format` in the enum, `language == "opaque"`). And the backend's `VerifyFormatMatchesSchema` still requires a posted instance's format to equal the schema's declared `instance_format`. Passthrough applies to *content*, not to *format* or *path*.

**No defaults.** `Defaults()` returns empty content — an opaque schema defines no values, unlike JSON Schema (`default:`) and CUE (`| *value`).

**Language display names and ordering** in core: CUE, JSON Schema, Opaque. Append Opaque after CUE in the existing JSON Schema → CUE lists, which preserves existing order.

### Repo layout and conventions

Mintlify site. Content root is `docs/` inside the repo (page `docs/cfg-mgmt/x.mdx` → URL `/cfg-mgmt/x`). Navigation is `docs/docs.json`, with page paths relative to the content root and **without** the `.mdx` extension. Shared snippets live under `docs/snippets/`.

**Frontmatter convention** on language pages: `title` only, double-quoted, sentence case. No `description`, `icon`, or `sidebarTitle`.

**Page skeleton** shared by `jsonschema.mdx` and `cue.mdx`: frontmatter → blank line → import block (no blank lines inside; sorted case-insensitively by path; `.mdx` imports use default syntax, `.jsx` component imports use `{ Named }` with spaces inside the braces; every import `;`-terminated) → `<Framed>` hero → one-sentence definition → motivation paragraph → `## Example` (a fenced block with a language and a label) → a bullet list of features → a version/spec section → `## File formats` (schema formats bullets, then instance formats bullets) → `<AgentYamlSupport />`.

**Lint** (`pnpm run lint` → `./scripts/lint.sh`) runs, in order: the custom Go linter over MDX, ESLint MDX (`--max-warnings=0`), CSpell (`cspell.json`), and `mint openapi-check`. `set -euo pipefail` means the first failing layer stops the run.

The Go linter's ten rules, and the ones that bite here:

1. `heading-case` — strict sentence case on the frontmatter `title:` and on every heading: first letter uppercase, all other letters lowercase, except an allowlist of exact tokens (`API`, `APIs`, `CLI`, `CI`, `SDK`, `SDKs`, `CUE`, `JSON`, `MQTT`, `TLS`, `HTTPS`, `REST`, `GUI`, `URL`, `ACLs`, `SSE`, `OpenAPI`, `AWS`, `GCP`, `GCS`, `WIF`, `STS`, `IAM`, `S3`, `ARN`, `SigV4`, plus proper nouns `Miru`, `GitHub`, `Agent`, `Unix`, `Git`, `Python`, `Schema`, `Base`, `Head`, `Cloud`, `Storage`, codenames, and version-like tokens). **`XML` is NOT on this list.** `title: "Opaque"` passes; `"Opaque Schemas"` fails; `## XML instances` fails. See the Decision Log — the plan avoids XML in headings rather than editing the allowlist.
2. `no-double-dash` — two consecutive hyphens in prose are banned; use an em dash (`—`).
3. Import hygiene: `import-resolves`, `import-used`, `import-sorted`, `import-block`, `import-component-style`, `import-mdx-style`.
4. `image-domain` — every image URL must begin with `https://assets.mirurobotics.com/`.
5. `redirects` — validates `docs.json` redirects.

**CSpell.** `cspell.json` has a `words` allowlist ordered in two runs: uppercase acronyms alphabetically first, then lowercase words alphabetically. Any new jargon must be inserted in the correct run, in order. Words this change may introduce that could trip CSpell: `nginx`, `passthrough`, `jsonc`. Mitigation is in M4.

**Commands.**

    pnpm install --frozen-lockfile
    pnpm run lint            # ./scripts/lint.sh — full docs lint (needs pnpm + go)
    pnpm run test:lint       # fixture smoke tests
    pnpm dev                 # mint dev preview
    ./scripts/preflight.sh   # everything

**CI** (`.github/workflows/ci.yml`): jobs `changes`, `lint`, `audit`, `shell-tests` on every PR; `lint-custom-linter` and `test-custom-linter` additionally when `tools/lint/**` changes. This plan does not touch `tools/lint/**`, so the last two must not appear.

**Commits**: conventional commits, e.g. `docs(schemas): add opaque language page`.

### Current state of every file to be changed

`docs/docs.json`, "Config schemas" group (verbatim):

    {
      "group": "Config schemas",
      "pages": [
        "cfg-mgmt/primitives/schemas/overview",
        {
          "group": "Languages",
          "pages": [
            "cfg-mgmt/primitives/schemas/languages/overview",
            "cfg-mgmt/primitives/schemas/languages/jsonschema",
            "cfg-mgmt/primitives/schemas/languages/cue"
          ]
        },
        "cfg-mgmt/primitives/schemas/manage"
      ]
    }

`docs/cfg-mgmt/primitives/schemas/languages/overview.mdx` — short page; a bullet list of JSON Schema (draft 2020-12) and CUE, then "The following sections cover key integration points for each language.", then a closing line linking a blog post comparing JSON Schema, CUE, and Protocol Buffers (leave that closing line alone).

`docs/snippets/schemas/formats.mdx` — a two-row table plus a trailing note that a schema's format need not match its instance's format.

`docs/cfg-mgmt/primitives/schemas/overview.mdx` — the `language` ParamField lists allowed values `jsonschema`, `cue`; the `format` ParamField already lists `cue`, `json`, `yaml` (no change needed); the `instance file path` ParamField says "Currently, JSON (`.json`) and YAML (`.yaml`, `.yml`) are supported."; there is a `## Schema languages` section with a two-item bullet list; a `## Empty schemas` section discusses starting with a permissive schema.

`docs/snippets/references/cli/releases/create/schema-annotations.mdx` — two `<ParamField>`s ("config type", "instance file path"), each containing a `<CodeGroup>` with a JSON Schema tab and a CUE tab. The instance-filepath ParamField also carries the "Currently, JSON (`.json`) and YAML (`.yaml`, `.yml`) are supported." sentence. Imported by `schemas/overview.mdx`, `cfg-mgmt/create-a-release.mdx`, and `getting-started/quick-start/create-release.mdx`.

`docs/snippets/schema-langs/support.mdx` — a single sentence, "Miru supports JSON Schema (Draft 2020-12) and CUE, the two most popular schema languages available." Imported by **no** page today; updating it is cheap insurance against it being wired up later with a stale claim.

Explicitly **not** touched: `docs/changelog/product.mdx` (user declined), `docs/references/platform-api/2026-05-06.yaml` (generated from Stainless via `api/pull-stainless.sh`; never hand-edited), the standalone `<Tabs>` in `docs/getting-started/quick-start/create-release.mdx`, and `tools/lint/**`.

## Plan of Work

### M1 — the new page and the nav entry

**Create `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx`.**

Frontmatter and imports (exactly one import; the `<Framed>` hero is deliberately omitted per the Decision Log, so `framed.jsx` is not imported — importing it unused would fail `import-used`):

    ---
    title: "Opaque"
    ---

    import AgentYamlSupport from '/snippets/agent/yaml-support.mdx';

Then, in order:

1. **Intro (no heading).** One sentence defining the language: an opaque schema is a metadata-only schema that attaches annotations to a config type without describing or constraining the config's contents. Then a short motivation paragraph in the voice of the other two pages (conversational second person, mild enthusiasm, no marketing superlatives): not every config can be validated — a vendor's text format, an XML document, a tool's own `.conf` file — but those configs still deserve versioning, releases, rollout, and an audit trail. An opaque schema is how you bring that kind of file into Miru. Mention the name's origin (Kubernetes `type: Opaque` — arbitrary contents, not inspected) as a one-liner.

2. **`## Example`.** A `<CodeGroup>` with a YAML tab and a JSON tab of the same document, so the reader sees both supported schema formats at once. Match the existing `<CodeGroup>` spacing convention in `schema-annotations.mdx`. Follow with a sentence noting that only `config_type` is required, so the minimal opaque schema is a single line.

3. **`## Schema document`.** State that an opaque schema is a **single** JSON or YAML mapping at the root, and give the four keys as a table (the table header row may use Title Case — `heading-case` applies to headings and frontmatter titles, not table cells):

   | Key | Required | Description |
   |---|---|---|
   | `config_type` | Yes | The config type slug the schema belongs to. |
   | `instance_filepath` | No | Absolute path where instances are written on the device. |
   | `instance_format` | No | The format of the config instance: `json`, `yaml`, or `xml`. |
   | `language` | No | Must be `opaque` when present. |

   Then the strictness rules, as a short bullet list: unknown keys are rejected (including `x-miru-` prefixed ones); exactly one document per schema — multi-document YAML is not accepted; every value must be a string. Cross-link `config_type` to the config types page, matching how `schema-annotations.mdx` links it.

4. **`## Validation`.** The honest section — this is the page's most important content, so be direct rather than euphemistic. An opaque schema validates the *schema document*, at creation time. It never validates a config instance: instance contents are never parsed, decoded, or inspected, and any bytes are accepted. Then the tradeoff, plainly: you keep versioning, releases, rollout, and Git provenance; you give up content safety — a typo in an opaque config reaches the device. Also state the two things opaque schemas do **not** give up:

   - A declared `instance_format` is still enforced — an instance's format must match the format the schema declares.
   - Instance file paths must still be absolute.

   And one thing they cannot do: an opaque schema defines no values, so it supplies **no defaults** (unlike `default:` in JSON Schema or `| *value` in CUE). Close by pointing readers who want a permissive-but-real schema at `## Empty schemas` on the schemas overview — that section is the right answer when the config *is* JSON or YAML and you just haven't written constraints yet; opaque is for content Miru cannot validate at all.

5. **`## File formats`.** Mirror the two-bullet-list shape of `cue.mdx` and `jsonschema.mdx`. Schema formats (opaque schemas themselves are always JSON or YAML — CUE is not accepted for this language):

   - **JSON** (`.json`)
   - **YAML** (`.yaml`, `.yml`)

   Then instance formats — note in prose that this is the widest instance-format set of any Miru schema language:

   - **JSON**
   - **YAML**
   - **XML**

   Keep "XML" out of the heading text; it is fine in bullets and prose.

6. **`## Instance file paths`.** The differentiator. Because an opaque schema imposes no structure on instance content, its instance file may carry **any** extension — where JSON Schema and CUE instances are limited to `.json`, `.yaml`, and `.yml`. The path must still be absolute. Examples: `/srv/miru/nginx.conf`, `/srv/miru/configs/motion-control.xml`, `/etc/myapp/settings.txt`. Cross-link the Agent file-system-access page the same way `schema-annotations.mdx` does, since custom directories need it.

7. **`<AgentYamlSupport />`** as the final line, matching both other language pages.

Heading inventory for this page, all of which pass `heading-case`: `## Example`, `## Schema document`, `## Validation`, `## File formats`, `## Instance file paths`. Frontmatter `title: "Opaque"` passes.

**Edit `docs/docs.json`.** Add `"cfg-mgmt/primitives/schemas/languages/opaque"` as a fourth entry in the nested `Languages` group, after `"cfg-mgmt/primitives/schemas/languages/cue"`. Add a comma after the `cue` line. Do **not** add a redirect.

### M2 — the two overview pages

**`docs/cfg-mgmt/primitives/schemas/languages/overview.mdx`.** Append a third bullet to the supported-languages list:

    - [Opaque](/cfg-mgmt/primitives/schemas/languages/opaque)

Add a short parenthetical or trailing clause distinguishing it — e.g. "(metadata only, no instance validation)" — so a reader scanning the list understands it is categorically different from the other two. Leave the closing blog-post sentence untouched.

**`docs/cfg-mgmt/primitives/schemas/overview.mdx`.** Three edits:

- `language` ParamField: add `opaque` to the allowed-values list, after `cue`.
- `## Schema languages` section: add a third bullet. The existing two bullets link to external sites (json-schema.org, cuelang.org); Opaque has no external site, so link it internally to `/cfg-mgmt/primitives/schemas/languages/opaque` with a short gloss such as "metadata-only schemas for configs Miru does not validate".
- `instance file path` ParamField: the sentence "Instance file paths control the type of file that config instances are deployed as. Currently, JSON (`.json`) and YAML (`.yaml`, `.yml`) are supported." is now incomplete. Append one sentence: opaque schemas may use any file extension, since their instance content is not validated. Do not restructure the ParamField.

The `format` ParamField needs no change (`cue`, `json`, `yaml` already covers opaque's JSON/YAML schema formats).

### M3 — the shared snippets

**`docs/snippets/schemas/formats.mdx`.** Add a third table row, preserving the existing column-alignment padding. Widen the separator and pad the other rows if needed so the pipes still line up (cosmetic, but the file is hand-aligned today). The Title Case header row stays as-is — `heading-case` does not inspect table cells.

**`docs/snippets/references/cli/releases/create/schema-annotations.mdx`.** Add a third tab to both `<CodeGroup>`s, after the CUE tab, matching the existing spacing. Also amend this ParamField's "Currently, JSON (`.json`) and YAML (`.yaml`, `.yml`) are supported." sentence with the same one-sentence opaque exception used in M2, so the two copies of that claim stay consistent. Keep the edit tight — this snippet renders inside the quick start.

**`docs/snippets/schema-langs/support.mdx`.** Rewrite the one-liner so it no longer claims a two-language world, e.g. "Miru supports JSON Schema (Draft 2020-12), CUE, and opaque schemas." Low risk: no page imports it today.

### M4 — lint, spelling, and render check

Run the full lint stack and fix what it finds. Anticipated CSpell candidates introduced by this change: `nginx`, `passthrough`. Prefer wording that avoids the allowlist entirely (`pass-through` reads fine; `/srv/miru/nginx.conf` is worth keeping). If a word must be added, insert it into `cspell.json`'s `words` array **in the correct run and in alphabetical order**. Do not append to the end of the array.

Then a `pnpm dev` render check of the new page and each touched page.

Finally, a corpus grep to prove the sweep is complete.

## Concrete Steps

All commands run from the repo root (`repos/docs` checkout). One commit per milestone, conventional-commit style.

**Setup (once)**

    pnpm install --frozen-lockfile

Expect: dependencies install without modifying `pnpm-lock.yaml`.

**M1**

Write `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` per the Plan of Work; edit the `Languages` group in `docs/docs.json`. Then:

    node -e "JSON.parse(require('fs').readFileSync('docs/docs.json','utf8')); console.log('docs.json parses')"
    pnpm run lint

Expect: `docs.json parses`, and lint ends with `All documentation lint checks passed.` If `heading-case` fires, the message names the offending heading — fix the heading, do not touch `tools/lint/`. Commit:

    git add docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx docs/docs.json
    git commit -m "docs(schemas): add opaque schema language page"

**M2**

Edit `docs/cfg-mgmt/primitives/schemas/languages/overview.mdx` and `docs/cfg-mgmt/primitives/schemas/overview.mdx`. Then:

    pnpm run lint

Expect: clean. Commit:

    git add docs/cfg-mgmt/primitives/schemas/languages/overview.mdx docs/cfg-mgmt/primitives/schemas/overview.mdx
    git commit -m "docs(schemas): list opaque in the schema language overviews"

**M3**

Edit the three snippets. Then:

    pnpm run lint

Expect: clean. Commit:

    git add docs/snippets/schemas/formats.mdx docs/snippets/references/cli/releases/create/schema-annotations.mdx docs/snippets/schema-langs/support.mdx
    git commit -m "docs(schemas): add opaque to shared format and annotation snippets"

**M4**

    pnpm run lint
    pnpm run test:lint

Expect: lint prints `All documentation lint checks passed.`; `test:lint` exits 0 with every fixture behaving as before (this change touches no fixtures and no linter code, so `test:lint` is a regression guard, not a new assertion).

Sweep checks — every one of these should print its fallback message:

    grep -rni "schemaless" docs/ || echo "no schemaless in docs"
    grep -rn "JSON Schema and CUE" docs --include='*.mdx' || echo "no two-language claims left"
    grep -rn "jsonc\|application/octet-stream" docs --include='*.mdx' || echo "no unsupported formats documented"
    git status --porcelain tools/lint || echo "tools/lint untouched"

Render check:

    pnpm dev

Then open, and confirm each renders without an MDX error and reads correctly:

- `/cfg-mgmt/primitives/schemas/languages/opaque` — appears in the left nav under **Languages** as the fourth entry, after CUE; all five headings present; both `<CodeGroup>` tabs switch; the `<AgentYamlSupport />` note renders at the bottom.
- `/cfg-mgmt/primitives/schemas/languages/overview` — three bullets.
- `/cfg-mgmt/primitives/schemas/overview` — `language` shows three allowed values, the formats table shows three rows, the annotations `<CodeGroup>`s show three tabs.
- `/cfg-mgmt/create-a-release` and `/getting-started/quick-start/create-release` — the shared annotations snippet still renders correctly with its third tab, and the quick start's own `<Tabs>` is unchanged.

Stop the dev server. If M4 changed any file (cspell allowlist, wording fixes):

    git add -A && git commit -m "docs(schemas): fix lint and spelling for opaque language docs"

If nothing changed, note that in Progress and move on.

**Validation**

    git push -u origin docs/opaque-schema-language

Then run preflight against the pushed branch head. Open the PR as **draft** and keep it draft until preflight reports CLEAN.

## Validation and Acceptance

- `pnpm install --frozen-lockfile` completes without modifying the lockfile.
- `pnpm run lint` exits 0 and prints `All documentation lint checks passed.` — this covers the Go prose linter (`heading-case`, `no-double-dash`, all import rules, `image-domain`, `redirects`), ESLint MDX at `--max-warnings=0`, CSpell, and `mint openapi-check`.
- `pnpm run test:lint` exits 0; every fixture behaves exactly as on `main` (no fixture or linter code was touched).
- `pnpm dev` renders `/cfg-mgmt/primitives/schemas/languages/opaque` with the new nav entry in position 4 under **Languages**, and every page listed in M4's render check renders without error.
- `grep -rni "schemaless" docs/` returns nothing — the superseded term appears nowhere in shipped content.
- `git status --porcelain tools/lint` is empty and `git diff --stat main...HEAD -- tools/lint` is empty — the `heading-case` allowlist was not modified, so the `lint-custom-linter` and `test-custom-linter` CI jobs must **not** run on this PR. If they do run, something under `tools/lint/**` was touched by mistake; revert it.
- `docs/changelog/product.mdx` and `docs/references/platform-api/2026-05-06.yaml` are unmodified (`git diff --stat main...HEAD` lists neither).
- `docs/docs.json` contains no new `redirects` entry.
- Every file in the sweep is accounted for in the diff: `docs/cfg-mgmt/primitives/schemas/languages/opaque.mdx` (new), `docs/docs.json`, `docs/cfg-mgmt/primitives/schemas/languages/overview.mdx`, `docs/cfg-mgmt/primitives/schemas/overview.mdx`, `docs/snippets/schemas/formats.mdx`, `docs/snippets/references/cli/releases/create/schema-annotations.mdx`, `docs/snippets/schema-langs/support.mdx`, plus `cspell.json` only if M4 required it.
- Content accuracy spot-check against the sources named in Context: instance formats documented are exactly JSON, YAML, XML (no `jsonc`, no `other`); schema formats are exactly JSON and YAML (no CUE); the "any file extension" claim is scoped to opaque instances and paired with the still-required absolute-path rule; the "no defaults" claim is present.
- CI: `changes`, `lint`, `audit`, and `shell-tests` green on the pushed branch head; `lint-custom-linter` and `test-custom-linter` not triggered. **Gate: preflight must report CLEAN (CI green on the pushed branch head) before the PR leaves draft or the task is reported complete.**

## Idempotence and Recovery

Every command above is read-only with respect to tracked source except the explicit file edits and commits. `pnpm install --frozen-lockfile` is idempotent and fails loudly rather than mutating the lockfile. `pnpm run lint`, `pnpm run test:lint`, the greps, and the `node -e` JSON parse check are pure verification and safe to re-run at any point. `pnpm dev` starts a local preview only; stop it with Ctrl-C.

All content changes are ordinary working-tree edits: inspect with `git status` / `git diff`, discard with `git checkout -- <path>`, and re-apply. Each milestone is an isolated commit, so a bad milestone can be dropped with `git reset --hard HEAD~1` before push, or `git revert <sha>` after push, without disturbing the others. The new page is a single added file — deleting it plus reverting the one-line `docs.json` nav addition fully undoes M1. If `pnpm run lint` fails from a stale environment rather than a real violation, re-run `pnpm install --frozen-lockfile` and retry. If `pnpm dev` cannot run (non-interactive environment), record that in Outcomes and rely on the static checks, which cover nav shape, imports, headings, and MDX parseability — but say so explicitly rather than silently skipping.

# Add data uploads entry to the product changelog

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (mirurobotics/docs) | read-write | Insert one new `<Update>` entry in `docs/changelog/product.mdx`. `cspell.json` may gain words only if the spell checker flags the new text (not expected). No other files change. |

This plan lives in `docs/plans/` because the only edit is in the docs repo. Work happens on the already-checked-out branch `claude/docs-changelog-data-uploads-kjirbo` (base: `main`). Do not create new branches.

## Purpose / Big Picture

Miru shipped **data uploads** on July 20, 2026 (Miru Agent `v0.10.0` + Miru CLI `v0.10.2`): devices upload files directly to the customer's own AWS S3 or Google Cloud Storage bucket, with the data never passing through Miru. The per-surface changelogs (`docs/changelog/agent.mdx`, `docs/changelog/cli.mdx`) and the feature docs (`docs/data-uploads/`) already exist, but the customer-facing product changelog (`docs/changelog/product.mdx`, served at `/changelog/product`) has no entry — its newest entry is still "June 23, 2026" (Access Control).

After this change, the product changelog page opens with a "July 20, 2026" entry announcing data uploads, linking to the feature docs and the per-surface changelog entries. A reader can verify by rendering the site (or reading the file) and seeing the new entry above the June 23 one.

## Progress

- [ ] Insert the `<Update label="July 20, 2026">` block into `docs/changelog/product.mdx`.
- [ ] Run `pnpm run lint` and `pnpm run test:lint` locally; fix any violations (add words to `cspell.json` only if flagged).
- [ ] Commit from the docs repo root: `docs: add data uploads product changelog entry`.
- [ ] Push via the preflight workflow; preflight reports `CLEAN` (CI green on the pushed branch head).

## Surprises & Discoveries

(Add entries as work proceeds.)

## Decision Log

- Decision: `<Update>` label is **"July 20, 2026"** — the release date of Agent `v0.10.0` / CLI `v0.10.2` — not today's date.
  Rationale: precedent in the same file — the "May 12, 2026" product entry matches the Agent `v0.9.0` release date in `docs/changelog/agent.mdx`. Product entries are labeled with the ship date.
  Date/Author: 2026-07-21 / plan author.
- Decision: the entry is **text-only** — no `<Framed>`, `<LazyVideo>`, or `<img>` blocks.
  Rationale: no changelog screenshot assets exist for this release under `assets.mirurobotics.com/docs/changelog/`; asset URLs must not be fabricated. Reusing the overview page's flow diagram was considered and rejected to keep the entry within the no-new-assets constraint.
  Date/Author: 2026-07-21 / plan author.
- Decision: feature heading is `## Data Uploads` (Title Case) with sentence-case `###` subsections.
  Rationale: matches the file's convention (`## Access Control`, `## Groups` with `### User types and roles`). The heading-case linter cannot flag this: its heading regex is anchored at line start (`tools/lint/linter/headingcase/headingcase.go`), and all headings inside `<Update>` blocks are indented four spaces.
  Date/Author: 2026-07-21 / plan author.
- Decision: minor agent fixes and CLI improvements from the two releases are folded into `<Dropdown title="Improvements">` / `<Dropdown title="Fixes">` blocks.
  Rationale: matches the May 27 and May 12 entries' convention for sub-headline items.
  Date/Author: 2026-07-21 / plan author.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The docs repo is the Mintlify documentation site. Content lives under `docs/` (the content root); `docs/changelog/product.mdx` is the product changelog page, rendered at `/changelog/product`. The file is a list of `<Update label="...">` blocks (a Mintlify component), newest first, wrapped in `<div className="changelog-page">`. Inside each block, content is indented four spaces: a `##` feature heading, prose, `###` subsections, doc links written as `[X documentation »](/path)`, and optionally `<Separator />`-separated `<Dropdown title="Improvements">` / `<Dropdown title="Fixes">` lists of `<DropdownItem>` bullets (each bullet starts with a bold `**Label:**`). The components `Dropdown`, `DropdownItem`, `Framed`, `LazyVideo`, and `Separator` are already imported at the top of the file — no import changes are needed, and unused-import lint stays satisfied because existing entries use them all.

Facts to convey (verified against the repo):

- Agent `v0.10.0` (July 20, 2026, `docs/changelog/agent.mdx`): the agent scans for files matching an upload rule's source glob and uploads them directly to the customer's AWS S3 or GCS bucket; data never passes through Miru. Reading source files may require granting the `miru` user read access (`/developers/agent/filesys-access#data-uploads`). Fixes: immediate heartbeat on SSE connection open; MQTT password redacted in logs; deployment dirty flag preserved when resetting retry state on startup; device RSA keys created with `0600`/`0640` permissions without a `chmod` race window.
- CLI `v0.10.2` (July 20, 2026, `docs/changelog/cli.mdx`): `miru release create` gains `--upload-rule` (single YAML file) and `--upload-rules` (directory) flags; API key source notice when `MIRU_API_KEY` is set; platform-aware upgrade command in the new-version notice; SBOM published for CLI releases.
- Feature docs (`docs/data-uploads/`): an *upload rule* is a YAML file in the customer's Git repo pairing a source (glob + `stability_window_secs` quiescence window) with a destination (bucket, object path template with `{device_id}`/`{upload_id}`/`{file_name}` variables, delete policy). Rules reference an *upload collection* via `collection_slug`; collections are created automatically. Buckets are connected per workspace (AWS S3 or GCS); for each upload the control plane issues short-lived credentials scoped to the exact object key, so devices are isolated from each other's data.

All link paths used in the new entry resolve to real pages: `docs/data-uploads/overview.mdx`, `docs/data-uploads/define-upload-rules.mdx`, `docs/data-uploads/connect-a-bucket/aws.mdx`, `docs/data-uploads/connect-a-bucket/gcs.mdx`, `docs/developers/agent/filesys-access.mdx` (has a `## Data uploads` heading → anchor `#data-uploads`), and the changelog anchors `/changelog/agent#v0-10-0`, `/changelog/cli#v0-10-2`.

Lint toolchain (all invoked from the repo root):

- `pnpm run lint` → `scripts/lint.sh`: builds and runs the Go MDX prose linter (`tools/lint` — rules include no-double-dash, which forbids `--` in prose outside backticks and requires the em dash "—"; heading-case; import rules), then ESLint (MDX validity), CSpell (dictionary in `cspell.json` at the repo root — add flagged words to its `words` array), then `mint openapi-check` on the OpenAPI specs. Requires `pnpm`, Go, and network access.
- `pnpm run test:lint` → `tests/test-lint.sh`: smoke tests of the linter itself.
- CI (`.github/workflows/ci.yml`) runs on the PR: `lint` (test:lint + `scripts/lint.sh`), `audit`, `shell-tests`, and `changes`; the custom-linter jobs are skipped because `tools/lint/**` is untouched.

## Plan of Work

One edit: in `docs/changelog/product.mdx`, insert the block below between the line `<div className="changelog-page">` (line 12 on `main`) and the line `<Update label="June 23, 2026">` (line 14 on `main`), keeping one blank line on each side of the new block. Do not modify any existing entry (including the intentional duplicate December 4, 2025 blocks and their guard comment). No import changes.

The exact block to insert (four-space indentation inside the `<Update>` tags, matching the file):

    <Update label="July 20, 2026">
        ## Data Uploads

        Devices can now upload files **directly to your own cloud storage**.
        Define an upload rule — what to collect, when it's finished, and where
        it goes — and the Miru Agent detects matching files on the device and
        streams them straight to your AWS S3 or Google Cloud Storage bucket.
        Your data never passes through Miru.

        [Data uploads documentation »](/data-uploads/overview)

        ### Upload rules

        Upload rules are YAML files versioned in Git alongside your config
        schemas. Each rule pairs a source with a destination:

        - **Source** — a glob pattern for the files to collect, plus a
          stability window so a file is only uploaded once it has finished
          being written
        - **Destination** — the bucket, an object path template with variables
          like `{device_id}`, `{upload_id}`, and `{file_name}`, and a delete
          policy for the source file after upload

        Every rule belongs to an **upload collection**, created automatically
        from the rule's `collection_slug`, which groups related uploads
        together.

        [Define upload rules »](/data-uploads/define-upload-rules)

        ### Connect a bucket

        Uploads land in a bucket you connect to your workspace — AWS S3 and
        Google Cloud Storage are supported. For each upload, the Miru control
        plane issues short-lived credentials scoped to the exact object key
        the file will be written to, and the agent uploads directly to your
        bucket. A device only ever receives credentials for its own uploads,
        so devices are strictly isolated from each other's data.

        - [Connect an AWS bucket »](/data-uploads/connect-a-bucket/aws)
        - [Connect a GCS bucket »](/data-uploads/connect-a-bucket/gcs)

        ### Ship upload rules with releases

        Upload rules are part of a release. `miru release create` now accepts
        `--upload-rule` (a single rule file) and `--upload-rules` (a directory
        of rule files), so the same command that versions your config schemas
        also versions your upload rules.

        To support data uploads, we've released:

        - [Miru Agent v0.10.0](/changelog/agent#v0-10-0) — detects files
          matching upload rules and uploads them to your bucket. Reading
          source files may require granting the `miru` user
          [read access](/developers/agent/filesys-access#data-uploads)
        - [Miru CLI v0.10.2](/changelog/cli#v0-10-2) — adds the upload rule
          flags to `miru release create`

        <Separator />

        <Dropdown title="Improvements">
            <DropdownItem>
                **API key source notice:** The CLI now prints a notice when it
                reads `MIRU_API_KEY` from the environment, making it obvious
                which credentials a scripted run is using.
            </DropdownItem>
            <DropdownItem>
                **Platform-aware upgrade command:** When a new CLI version is
                available, the update notification now shows the upgrade
                command matching how the CLI was installed.
            </DropdownItem>
            <DropdownItem>
                **CLI SBOM:** Every CLI release now publishes a software bill
                of materials (SBOM).
            </DropdownItem>
            <DropdownItem>
                **Immediate SSE heartbeat:** The agent sends a heartbeat as
                soon as a Server-Sent Events connection opens, so on-device
                applications can immediately tell the stream is healthy.
            </DropdownItem>
        </Dropdown>

        <Separator />

        <Dropdown title="Fixes">
            <DropdownItem>
                **MQTT credential redaction:** The agent now redacts the MQTT
                password in its logs.
            </DropdownItem>
            <DropdownItem>
                **Deployment retry state:** A deployment's dirty flag is
                preserved when the agent resets retry state on startup.
            </DropdownItem>
            <DropdownItem>
                **Device key permissions:** Device RSA keys are created with
                `0600`/`0640` permissions from the start, removing a `chmod`
                race window.
            </DropdownItem>
        </Dropdown>
    </Update>

Style constraints while editing: use the em dash "—" in prose, never `--` (the no-double-dash lint rule checks prose; `--upload-rule` is safe only inside backticks); keep bullet and link formats exactly as above; do not add trailing whitespace.

## Concrete Steps

All commands run from the docs repo root, `/home/user/docs` (adjust the prefix if your checkout lives elsewhere; the repo root is the directory containing `package.json` and `docs/`).

1. Confirm the branch and a clean tree:

       cd /home/user/docs
       git branch --show-current   # expect: claude/docs-changelog-data-uploads-kjirbo
       git status --short          # expect: no output (clean)

2. Confirm the entry is not already present (idempotence guard):

       grep -c 'label="July 20, 2026"' docs/changelog/product.mdx   # expect: 0 (grep exits 1)

   If it prints 1, the insertion is already done — skip to step 4.

3. Edit `docs/changelog/product.mdx`: insert the block from Plan of Work between `<div className="changelog-page">` and `<Update label="June 23, 2026">`, with one blank line before and after the block.

4. Lint:

       cd /home/user/docs
       pnpm install --frozen-lockfile
       pnpm run test:lint          # expect: all smoke tests pass
       pnpm run lint               # expect final line: "All documentation lint checks passed."

   If CSpell flags a word in the new entry, add it to the `words` array of `/home/user/docs/cspell.json` (matching the surrounding loose-alphabetical style) and re-run. Do not add words that were not flagged. If the Go prose linter flags a line, fix the text to conform to the rule — do not disable rules.

5. Verify the diff touches only the intended files:

       git diff main --stat   # expect: docs/changelog/product.mdx only (plus cspell.json only if step 4 required it, and this plan file)

6. Commit (one milestone, one commit), from `/home/user/docs`:

       git add docs/changelog/product.mdx plans/   # picks up this plan file wherever it lives (backlog/ or active/); add cspell.json too only if it changed
       git commit -m "docs: add data uploads product changelog entry"

7. Publish and watch CI via the preflight workflow (push the branch, let the CI `lint`, `audit`, and `shell-tests` jobs run on the pushed head, and fix any failures from the job logs).

## Validation and Acceptance

- `pnpm run lint` from `/home/user/docs` exits 0 and prints "All documentation lint checks passed." after the change (it also passes before — this change must not regress it).
- `pnpm run test:lint` from `/home/user/docs` exits 0.
- `docs/changelog/product.mdx` contains exactly one `<Update label="July 20, 2026">` block, positioned above the `<Update label="June 23, 2026">` block; all pre-existing entries are byte-identical (verify: `git diff main -- docs/changelog/product.mdx` shows only an insertion).
- Every link in the new entry points at a page verified in Context and Orientation; none reference `assets.mirurobotics.com` (text-only entry).
- Optional render check: `pnpm run dev` from `/home/user/docs` (runs `mint dev` in `docs/`) and open `/changelog/product` — the new entry renders at the top with working dropdowns.
- CI on the pushed branch head is green: the `lint`, `audit`, `shell-tests`, and `changes` jobs of `.github/workflows/ci.yml` all pass (custom-linter jobs are skipped since `tools/lint/**` is untouched).
- **Gate: preflight must report `CLEAN` — CI green on the pushed branch head — before the PR leaves draft or the task is reported complete.**

## Idempotence and Recovery

- The insertion is guarded by step 2's grep; re-running the steps never duplicates the entry. Lint and diff steps are read-only and always safe to repeat.
- If the edit goes wrong before committing, restore the file: from `/home/user/docs`, run `git checkout main -- docs/changelog/product.mdx` and redo step 3.
- If a bad version was committed, `git revert <sha>` from `/home/user/docs` (do not force-push shared branches).
- `pnpm install --frozen-lockfile` is idempotent and never modifies `pnpm-lock.yaml`; if it fails on a network hiccup, re-run it.

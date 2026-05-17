# ExecPlan: Changelog CLI v0.10.0 agent-version caveat + product.mdx duplicate-date guard comment

Status: backlog
Created: 2026-05-17
Branch: claude/audit-docs-4CgNO
Repo: /home/user/docs (Mintlify)

## Summary

Two tiny, independent MDX documentation edits:

1. **Defect fix** — `docs/changelog/cli.mdx`: resolve an internal
   contradiction in the v0.10.0 section. The page-level rule (file paths
   outside `/srv/miru/config_instances/` require Miru Agent v0.8.0+) is
   not restated where the new default path is introduced, so the
   "Default instance file path" `<Dropdown>` reads as if the new default
   has no agent-version requirement even though it lives outside that
   directory.

2. **Guard comment (not a defect)** — `docs/changelog/product.mdx`: add a
   single MDX comment above the first of two intentional, identical
   `<Update label="December 4, 2025">` blocks so future edits/automation
   do not merge or rename them.

These are documentation-only, content-clarification edits. No path
values, version numbers, behavior, or block content change.

## Verified Facts (confirmed by reading the repo on 2026-05-17)

- `docs/changelog/cli.mdx`:
  - Line 24: `# v0.10.0` heading; line 26 date `*April 12, 2026*`.
  - Line 30 (page-level rule):
    `Note that file paths which do not live in `/srv/miru/config_instances/` require Miru Agent [v0.8.0](/changelog/agent#v0-8-0).`
  - Lines 50-64: the `<Dropdown title="Default instance file path">`
    block. Line 52 states the default was updated from
    `/srv/miru/config_instances/{config-type-slug}.json` to
    `/srv/miru/configs/{config-type-slug}.json`. Lines 56-62 contain a
    ```` ```diff ```` fenced block; line 64 closes `</Dropdown>`.
- `/srv/miru/configs/{config-type-slug}.json` is the **correct** canonical
  default — confirmed at `docs/learn/schemas/overview.mdx:67`
  ("The default instance file path is `/srv/miru/configs/{config-type-slug}.json`.")
  and `docs/snippets/references/cli/releases/create/schema-annotations.mdx:31`
  ("defaults to `/srv/miru/configs/{config-type-slug}.json`"). DO NOT
  change path values.
- `docs/changelog/agent.mdx:47` is `# v0.8.0`; line 51 states Agent
  v0.8.0 allows config instances outside `/srv/miru/config_instances/`.
  The link target `/changelog/agent#v0-8-0` resolves to this `# v0.8.0`
  heading (Mintlify slugifies `v0.8.0` -> `v0-8-0`). Link is valid and
  must be reused verbatim — do not invent a new anchor.
- `docs/changelog/product.mdx`: two intentional
  `<Update label="December 4, 2025">` blocks at line 620 ("CUE Support")
  and line 665 ("Workspace Invitations"). Two distinct releases shipped
  the same day. This is intentional and must be preserved as-is.

> Line numbers above were verified on 2026-05-17. Re-read both files
> immediately before editing and re-locate the anchors by content (the
> exact strings quoted above), not by line number, in case the files
> have shifted.

## Change 1 — docs/changelog/cli.mdx (defect fix)

Goal: make the v0.10.0 "Default instance file path" `<Dropdown>`
self-consistent with the page-level rule on line ~30, by restating the
already-documented requirement that paths outside
`/srv/miru/config_instances/` require Miru Agent v0.8.0+.

Constraints:
- Do NOT change any path values, version numbers, or the existing
  `<Dropdown>` / ```` ```diff ```` structure.
- Add ONLY a short clarifying note. It must restate the existing
  page rule — paths outside `/srv/miru/config_instances/` require Miru
  Agent v0.8.0+ — and nothing more. Do not invent behavior or versions.
- Reuse the existing link `/changelog/agent#v0-8-0` verbatim.
- This is a needs-verification item w.r.t. real agent behavior; the note
  must only restate what the page already asserts (line ~30), not assert
  anything new.

Edit: inside the `<Dropdown title="Default instance file path">` block,
immediately after the sentence ending
`... to `/srv/miru/configs/{config-type-slug}.json`.` (currently
cli.mdx:52) and before the "We recommend adding..." paragraph
(currently line 54), insert a single new note paragraph (blank line
above and below to keep valid MDX / Markdown block separation):

> Because the new default path is outside `/srv/miru/config_instances/`,
> it requires Miru Agent [v0.8.0](/changelog/agent#v0-8-0) or later, as
> noted above.

Notes:
- Keep it to one short paragraph. Do not modify lines 56-62 (the
  ```` ```diff ```` block) or the `</Dropdown>` close (line 64).
- Use a plain Markdown paragraph (not a Mintlify `<Note>` component)
  to minimize structural risk; surrounding content in this Dropdown is
  plain Markdown paragraphs, so this matches existing style.
- Verify backtick-inline-code and the Markdown link render correctly.

## Change 2 — docs/changelog/product.mdx (guard comment, NOT a defect)

Goal: prevent future accidental merge/rename of the two intentional
`December 4, 2025` `<Update>` blocks.

Edit: immediately ABOVE the FIRST `<Update label="December 4, 2025">`
block (currently product.mdx:620, the "CUE Support" block), on its own
line(s) with a blank line separating it from the preceding `</Update>`
(line 618) and from the `<Update ...>` line below, insert exactly one
MDX comment:

```
{/* INTENTIONAL: two distinct releases shipped on December 4, 2025
    ("CUE Support" and "Workspace Invitations" below). The duplicate
    date label is deliberate — do NOT merge, rename, or re-date either
    <Update> block. */}
```

Constraints:
- Add ONLY this comment. Do not merge, rename, re-date, reorder, or
  modify either `<Update>` block or any of their content.
- Must be valid MDX: `{/* ... */}` JSX-style comment, placed between
  block elements (not inside a tag), with surrounding blank lines.

## Validation Steps (required before publishing)

Run from repo root `/home/user/docs`:

1. **MDX/lint:** `./scripts/lint.sh` (equivalent to npm `lint`) — must
   pass with no new errors/warnings attributable to the two edited
   files.
2. **Link resolution:** Confirm every internal link in the edited
   regions still resolves. Specifically verify `/changelog/agent#v0-8-0`
   resolves to the `# v0.8.0` heading in `docs/changelog/agent.mdx`
   (slug `v0-8-0`). Also re-check the pre-existing link on cli.mdx:30
   was not disturbed.
3. **Structural spot-check (manual read):**
   - cli.mdx: the `<Dropdown title="Default instance file path">` still
     opens and closes correctly, the ```` ```diff ```` fence is intact
     and balanced, and the new note is a standalone paragraph with blank
     lines around it. No path values changed (grep for
     `/srv/miru/configs/{config-type-slug}.json` and
     `/srv/miru/config_instances/` to confirm unchanged).
   - product.mdx: still exactly two `<Update label="December 4, 2025">`
     occurrences (`grep -n 'December 4, 2025'`), both blocks byte-for-byte
     unchanged; the new `{/* ... */}` comment sits above the first one
     and is not inside any tag.
4. **Diff review:** `git diff` shows changes ONLY in
   `docs/changelog/cli.mdx` and `docs/changelog/product.mdx`, limited to
   the additions described above (no path/version/content edits).
5. **Preflight:** `./scripts/preflight.sh` MUST report `clean`. Changes
   MUST NOT be published or the PR opened if preflight is not `clean`.
   Resolve any preflight findings and re-run until `clean`.

## Delivery

- Branch: `claude/audit-docs-4CgNO` (already current).
- Commit the two file edits with a clear message, e.g.
  `docs(changelog): clarify v0.10.0 agent v0.8.0 requirement; guard duplicate Dec 4 2025 updates`.
- Open a PR targeting `main`.
- Gate: only open/finalize the PR after `./scripts/preflight.sh` reports
  `clean`.

## Out of Scope / Explicit Non-Goals

- No changes to any path values (`/srv/miru/configs/...`,
  `/srv/miru/config_instances/...`).
- No new version numbers or behavior claims.
- No merging, renaming, re-dating, or reordering of the
  `December 4, 2025` `<Update>` blocks.
- No edits to `agent.mdx`, `overview.mdx`, or schema-annotation
  snippets — they are referenced only to confirm canonical values.

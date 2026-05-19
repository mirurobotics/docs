# Restore inter-paragraph spacing inside the custom `<Dropdown>` component

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Edit `docs/snippets/components/dropdown.jsx` to add vertical rhythm between body children. |

This plan lives in `docs/plans/backlog/` because the only code change is in the docs repo (`docs/snippets/components/dropdown.jsx`).

Working branch: `claude/audit-docs-4CgNO` (already current — do not change it). Base branch for the eventual PR: `main`.

## Purpose / Big Picture

The custom `Dropdown` body wraps `{children}` in `<div className="pb-3">`. Mintlify's docs site uses Tailwind, and Tailwind's preflight resets remove default `<p>` top/bottom margins. The surrounding Mintlify prose styles do not cascade into our custom component, so MDX paragraphs that are separated by a blank line inside a `<Dropdown>` render glued together with no visible gap.

After this change, MDX content inside `<Dropdown>` with blank-line-separated paragraphs (and other block-level children, e.g. fenced code blocks) renders with the same inter-paragraph spacing as the surrounding non-dropdown prose on the same page.

The clearest user-visible verification surface is `docs/changelog/cli.mdx` lines 50–66 (`<Dropdown title="Default instance file path">`), which contains three paragraphs separated by blank lines and a fenced `diff` code block. Before the fix these run together; after, they sit visibly apart with the same rhythm as surrounding `.mdx` paragraphs on the same page.

## Progress

- [ ] Add `space-y-4` (or equivalent rhythm utility — see Plan of Work) to the Dropdown body wrapper in `docs/snippets/components/dropdown.jsx`.
- [ ] Re-read `docs/snippets/components/api-dropdown.jsx` to confirm whether the same defect exists; record finding in Surprises & Discoveries (out of scope for fix per task constraints, but note it).
- [ ] Run `pnpm lint` from `docs/`; verify no new findings attributable to the edit.
- [ ] Run `pnpm test:lint` from `docs/`; verify no new findings attributable to the edit.
- [ ] If a local preview server is feasible (see Concrete Steps), render `docs/changelog/cli.mdx` and confirm the "Default instance file path" dropdown shows three visually separated paragraphs with spacing matching surrounding non-dropdown paragraphs. If not feasible, record that explicitly in Surprises & Discoveries.
- [ ] Run `./scripts/preflight.sh` from `docs/`; confirm it reports `clean`.
- [ ] Commit the single-file edit (one commit for this single-milestone plan).

Use timestamps when you complete steps. Split partially completed work into "done" and "remaining" as needed.

## Surprises & Discoveries

(Add entries as you go.)

- Observation: …
  Evidence: …

## Decision Log

(Add entries as you go.)

- Decision: …
  Rationale: …
  Date/Author: …

## Outcomes & Retrospective

(Summarize at completion or major milestones.)

## Context and Orientation

This repo is the Miru documentation site, hosted on Mintlify (`mint` v4.x — see `docs/package.json`). Pages live under `docs/` as `.mdx` files. Reusable components are JSX files under `docs/snippets/components/`. Mintlify's runtime takes the JSX components and renders them alongside MDX content using Tailwind CSS.

Files and concepts relevant to this change:

- `docs/snippets/components/dropdown.jsx` — defines three exports: `DropdownItem`, `DropdownGroup`, and `Dropdown`. Only `Dropdown` is changed here. As of the time of authoring, line 67 reads:

      {isOpen && <div className="pb-3">{children}</div>}

  The wrapping `<div>` only has bottom padding; it has no spacing utility applied across child block elements.

- `docs/snippets/components/api-dropdown.jsx` — a near-duplicate component used by API reference pages. It contains the same `<div className="pb-3">{children}</div>` wrapper. **Out of scope for this fix per the task constraints**, but the reader should re-read this file during the work and record an explicit observation in Surprises & Discoveries about whether the same defect is present, so a follow-up can be filed later if appropriate.

- `docs/changelog/cli.mdx` — the canonical reproduction page. Lines 50–66 contain `<Dropdown title="Default instance file path">` with three paragraphs separated by blank lines and a `diff` fenced code block. This dropdown is the chosen verification surface.

- `package.json` — defines two lint scripts: `pnpm lint` (runs `./scripts/lint.sh`) and `pnpm test:lint` (runs `./tests/test-lint.sh`). Mint is pinned to `4.2.565`.

- `scripts/preflight.sh` — runs `pnpm run test:lint`, the Go-based MDX prose lint (`tools/lint`), Go coverage gate, `./scripts/lint.sh`, `./scripts/audit.sh`, and a `bats` test. Must report `clean` before the change is published.

- No `AGENTS.md` or `CLAUDE.md` exists at the docs repo root (verified at authoring time). There is no documented `pnpm dev`/`mint dev` preview script in `package.json`; preview is only feasible if `mint dev` (from the `mint` dev dependency) is runnable in the current environment. If it is not, do not invent a preview workflow — explicitly note "no preview performed" in Surprises & Discoveries.

Terms:

- **Tailwind preflight**: Tailwind CSS's base reset stylesheet, which (among other things) removes default `margin`/`padding` from `<p>`, `<h1>`–`<h6>`, lists, and other block elements. This is why MDX paragraphs inside a custom Tailwind-styled container have no inherent vertical gap.
- **MDX**: Markdown with embedded JSX. Blank-line separation between paragraphs is parsed into separate `<p>` elements; under Tailwind preflight these `<p>` elements have no margin.
- **`space-y-4`**: Tailwind utility that applies `margin-top: 1rem` to every direct child after the first, producing inter-child vertical rhythm without affecting outer padding. `1rem` (16px) matches Tailwind's default prose paragraph spacing scale and Mintlify's surrounding prose rhythm.

## Plan of Work

Make a single one-line className change to `docs/snippets/components/dropdown.jsx`:

- File: `docs/snippets/components/dropdown.jsx`
- Function: the `Dropdown` component's return body.
- Location: the line currently reading `{isOpen && <div className="pb-3">{children}</div>}` (line 67 at authoring time — re-locate by content, not by line number, in case the file has shifted).
- Change: add `space-y-4` to the wrapper `<div>` so every direct child after the first has `margin-top: 1rem`. The line becomes:

      {isOpen && <div className="pb-3 space-y-4">{children}</div>}

Rationale for `space-y-4` specifically:

- `space-y-4` is `1rem`, which matches the default vertical rhythm of paragraphs in Tailwind's typography plugin and the spacing of MDX paragraphs in surrounding Mintlify prose. If, after visual check, the gap is clearly too large or too small relative to the surrounding non-dropdown paragraphs on the same page, the implementer may step to `space-y-3` (0.75rem) or `space-y-6` (1.5rem). Record any such adjustment in the Decision Log.
- `space-y-*` applies only to direct children, so it correctly spaces sibling block elements (paragraphs, fenced code blocks, nested `<Dropdown>` blocks) without leaking spacing into single-child cases.

Do NOT touch any other file. In particular:

- Do NOT modify `docs/snippets/components/api-dropdown.jsx` even if the same defect is present — record it in Surprises & Discoveries for a separate follow-up.
- Do NOT modify any `.mdx` content.
- Do NOT add new tests. The repo does not have a React component rendering harness (verified at authoring time — `tests/` contains `lint-fixtures/` and `test-lint.sh`, no component test runner). A className-only tweak under an existing lint regime does not warrant adding one.

## Concrete Steps

All commands run from the docs repo root: `/home/ben/miru/workbench3/repos/docs/`.

1. Confirm the current state of the target line.

       cd /home/ben/miru/workbench3/repos/docs
       grep -n 'pb-3' docs/snippets/components/dropdown.jsx

   Expected output (line number may differ; the content is what matters):

       67:            {isOpen && <div className="pb-3">{children}</div>}

2. Edit `docs/snippets/components/dropdown.jsx`. Change the matched line to:

       {isOpen && <div className="pb-3 space-y-4">{children}</div>}

   Re-run the grep to confirm:

       grep -n 'space-y-4' docs/snippets/components/dropdown.jsx

   Expected: one match on the same line as `pb-3`.

3. Re-read `docs/snippets/components/api-dropdown.jsx` to check whether the same defect exists. Record the finding in Surprises & Discoveries. Do NOT change it.

       grep -n 'pb-3' docs/snippets/components/api-dropdown.jsx

4. Confirm the diff is limited to the one file and one line:

       git diff --stat
       git diff docs/snippets/components/dropdown.jsx

   Expected: `docs/snippets/components/dropdown.jsx | 2 +- 1 file changed, 1 insertion(+), 1 deletion(-)` (or similar). No other files changed.

5. Run lint:

       pnpm lint

   Expected: ends with `All documentation lint checks passed.` and exits 0. No new errors/warnings attributable to the change. (The change is to a `.jsx` snippet file, which `scripts/lint.sh` does not currently scan for MDX prose, so this should be a no-op for the edit but must still pass overall.)

6. Run the lint smoke tests:

       pnpm test:lint

   Expected: exits 0.

7. Attempt a local preview, if feasible. Mintlify CLI usage is:

       pnpm exec mint dev

   If that command starts a local dev server, open the changelog/cli page in the preview, expand "Default instance file path", and visually confirm that the three paragraphs are visually separated by a gap that matches the spacing of surrounding non-dropdown paragraphs on the same page (e.g. the paragraphs under `## Fixes` and `## Improvements`).

   If `mint dev` is not available in the current environment (e.g. it requires network access, an account, or a binary not installed), do NOT claim visual verification. Record "no preview performed" in Surprises & Discoveries with the reason. The lint + preflight steps remain mandatory regardless.

8. Run preflight:

       ./scripts/preflight.sh

   Expected: exits 0 with no failures. Preflight runs `pnpm test:lint`, `tools/lint`, `tools/lint` covgate, `scripts/lint.sh`, `scripts/audit.sh`, and the `bats` test. **Preflight MUST report `clean` (exit 0) before changes are published.** Do not skip or work around any preflight finding — resolve the underlying issue and re-run until clean.

9. Commit the single-file edit. From `/home/ben/miru/workbench3/repos/docs/`:

       git status
       git add docs/snippets/components/dropdown.jsx
       git commit -m "fix(snippets/dropdown): restore inter-paragraph spacing in body"

   This is a single-milestone plan, so this commit is the milestone-end commit. Do NOT amend later changes onto it — if further fixes are needed, create a new commit.

10. Confirm the commit:

        git log -1 --stat

    Expected: one file changed (`docs/snippets/components/dropdown.jsx`), one insertion, one deletion.

## Validation and Acceptance

Acceptance is observable, not structural:

- **Lint passes.** From `docs/`: `pnpm lint` exits 0 and ends with `All documentation lint checks passed.`; `pnpm test:lint` exits 0.
- **Preflight is clean.** From `docs/`: `./scripts/preflight.sh` exits 0 with no failures. This is the gate for publishing.
- **Diff scope is correct.** `git diff main -- docs/snippets/components/dropdown.jsx` shows exactly one line changed: `className="pb-3"` becomes `className="pb-3 space-y-4"`. No other files changed by this plan.
- **Visual rendering (if preview feasible).** On the `docs/changelog/cli.mdx` page, expanding the `<Dropdown title="Default instance file path">` shows three paragraphs that are visually separated by a vertical gap. The gap should look the same as the gap between sibling paragraphs in the non-dropdown sections on the same page (e.g. the bulleted lists under `## Fixes` and `## Improvements` are unaffected; the comparison is to plain paragraph blocks elsewhere on the page).

  If preview is not feasible in the current environment, mark this acceptance criterion as "deferred to PR review render" and record it explicitly in Surprises & Discoveries — do not claim a visual verification that was not performed.

There is no existing component-rendering test in this repo, so there is no "test X fails before and passes after" assertion to record. A className tweak in a Mintlify snippet does not justify standing up a new component test harness.

## Idempotence and Recovery

- The edit is a single one-line className change. Re-running the same edit on an already-edited file is a no-op (the target string is already present).
- If `pnpm lint`, `pnpm test:lint`, or `./scripts/preflight.sh` fails: read the failure, fix the underlying cause, and re-run. Do not bypass with `--no-verify` or env overrides.
- If visual review (post-merge or in PR preview) shows the spacing is too tight or too loose, replace `space-y-4` with `space-y-3` or `space-y-6` and re-run lint + preflight. Record the change in the Decision Log.
- To revert: `git revert <commit-sha>` of the single commit produced by this plan. There are no migrations, no data, and no other files involved, so revert is safe.

# Image-domain lint rule for docs MDX

This ExecPlan is a living document. Update Progress, Surprises & Discoveries, and the Decision Log as work proceeds; fill in Outcomes & Retrospective on completion.

## Scope

| Repository | Access | Branch |
|---|---|---|
| mirurobotics/docs (this repo; all paths below are repo-relative) | read/write | `feat/image-domain-lint` (already checked out; base `main`) |

## Purpose / Big Picture

Docs images must be served from the CDN, not committed to the repo or hot-linked from other domains. This plan adds a rule `image-domain` to the repo's custom Go linter (`tools/lint/`) that flags any image reference in MDX content whose URL is not `https://assets.mirurobotics.com/...`. Relative/local paths (`/images/...`), other domains, protocol-relative `//...`, and `http://...` are all violations unless suppressed with a new `{/* lint-ignore image-domain */}` directive.

After implementation, running `./scripts/lint.sh` from the repo root fails with output like (hypothetical; paths print as passed by `scripts/lint.sh`, i.e. absolute)

    /abs/path/docs/some-page.mdx:12:5: image-domain: image must be hosted on https://assets.mirurobotics.com (got "/images/x.png")

whenever an MDX file references a non-CDN image, and passes on the current corpus (which is already clean — see Decision Log).

## Progress

- [x] M1: implement `tools/lint/linter/imagedomain/` package, register the rule in `tools/lint/linter/run.go`, add `.covgate`; commit.
- [ ] M2: unit tests (`imagedomain_test.go` + `run_test.go` case) passing with coverage at or above the gate; commit.
- [ ] M3: `tests/lint-fixtures/bad-image-domain/` fixture, `tests/test-lint.sh` wiring, pass-path additions to the `good/` fixture; commit.
- [ ] M4: corpus verification sweep, full `./scripts/lint.sh` green; commit (only if the sweep changed files).
- [ ] Validation: push branch, all CI jobs green, preflight reports CLEAN.

## Surprises & Discoveries

Add entries as work proceeds.

## Decision Log

- 2026-07-23 (planning agent): Scan raw lines with rule-local code-fence tracking instead of extending `analysis.Scanner` with a zone accessor. JSX/HTML attributes are masked out of prose spans, so span-based scanning cannot see `<img src>` / `<Framed image>`; raw-line scanning with the rule doing its own fence skipping is the smallest change. Tradeoff: raw lines also expose inline code (backticks) and comment text, so a prose example like `<img src="/images/x.png">` in backticks false-positives (none in the corpus today); the ignore directive is the escape hatch.
- 2026-07-23 (planning agent): LazyVideo `poster` dropped from scope — the component (`docs/snippets/components/lazy-video.jsx`) has no `poster` prop; its props are `{src, alt, className}` only. `src` is checked only when the URL has an image file extension, so `.mp4` videos are unaffected. (The `poster` attribute name is still matched strictly, as future-proofing for raw `<video poster=...>`.)
- 2026-07-23 (planning agent): New suppression directive `{/* lint-ignore image-domain */}` (suppresses the next line only, parsed inside the rule from raw lines). This is the first suppression mechanism in this linter; it is deliberately rule-scoped and minimal, not a generic framework.
- 2026-07-23 (planning agent): Site chrome in `docs/docs.json` (logo/favicon/og:image) is out of scope automatically — the linter only receives `*.mdx` files from `scripts/lint.sh`, and the only docs.json-reading rule (redirects) inspects only the redirects array.
- 2026-07-23 (planning agent): Existing-content migration reduced to a verification sweep — the corpus has zero violations today; the known placeholder `image="/images/changelog/bucket-page.png"` never reached this branch — PR #136's squash commit `1b4365e` (the branch's base on `main`) landed `docs/changelog/product.mdx` already remediated as `https://assets.mirurobotics.com/docs/changelog/26-07-23/bucket-page.png`.

## Outcomes & Retrospective

Fill in on completion.

## Context and Orientation

**The custom linter.** `tools/lint/` is a Go module (`github.com/mirurobotics/docs/tools/lint`, Go 1.25.3). Entry point `tools/lint/main.go` implements `lint <file>...`: it lints a batch of files and prints violations as `file:line:col: message` (exit 0 clean, 1 violations, 2 usage/IO error). It derives a "content root" by walking up from the first file to a directory containing `snippets/`.

**Rules.** Each rule is a package under `tools/lint/linter/<name>/` exposing a `Check(...)` function (rule-specific signature) returning `[]analysis.Violation`. `analysis.Violation` (in `tools/lint/linter/analysis/analysis.go`) is `{File string; Line int; Col int; Message string}` — Line and Col are 1-based, Col is a byte column. Messages are conventionally prefixed with the rule id, e.g. `"no-double-dash: ..."`. Rules are registered in `tools/lint/linter/run.go`: a `Rule` string constant (e.g. `RuleNoDoubleDash Rule = "no-double-dash"`), an entry in `AllRules()`, and a `ruleEntry` in `ruleCheckers()` whose closure receives `checkInput{path string; lines []string; spans [][]analysis.ProseSpan; contentRoot string}`. Existing rule ids: no-double-dash, heading-case, import-resolves, import-used, import-sorted, component-style, mdx-style, import-block, redirects. `image-domain` is free. The structural mirror is `tools/lint/linter/mdxstyle/mdxstyle.go` — same `Check(file string, lines []string)` signature as this rule (nodoubledash is spans-based; it remains the model for M2's test style only).

**Why raw lines, not prose spans.** The scanner (`tools/lint/linter/analysis/scanner.go`) produces per-line "prose spans" — segments of a line with inline code, HTML/JSX tags, and comments masked out; frontmatter and code-fence lines yield nil. Markdown images `![alt](url)` are visible in spans, but `<img>`, `<Framed>`, and `<LazyVideo>` attributes are NOT (they live inside masked tags). The scanner exposes no zone accessor. This rule therefore takes `in.lines` (raw lines) and does its own lightweight skipping: track code fences (a line whose trimmed content starts with three backticks toggles fence state; delimiter lines and lines inside a fence are skipped) and YAML frontmatter (if line 1 is `---`, skip through the closing `---` inclusive), so code examples and frontmatter are never flagged.

**Image-bearing constructs in MDX** (content root is `docs/`; snippets under `docs/snippets/` are also linted):

- Markdown images: `![alt](url)` — 155 in the corpus.
- `<img src="...">` and `<LazyVideo src="...">` — 51 `src=` occurrences (one JSX brace-quoted, `docs/changelog/product.mdx`). `docs/snippets/components/lazy-video.jsx` renders `<video src>`; all current usages are `.mp4`.
- `<Framed image="..." background="..." />` — 85 `image=`/`background=` occurrences; `docs/snippets/components/framed.jsx` renders both as image URLs.

Component usages frequently span multiple lines (one attribute per line), so matching must be per-attribute-occurrence per-line, never require a whole tag on one line. Filenames may contain a colon (e.g. `.../releases/header:page.png`) and `.light.svg`/`.dark.svg` variants exist — URL handling must tolerate both. All current references are on `https://assets.mirurobotics.com`; the corpus is clean today. Local orphaned originals under `docs/images/` are unreferenced and out of scope (possible follow-up cleanup).

**Meta-lint constraints on Go code** (enforced by CI job `lint-custom-linter`: golangci-lint + gofumpt via `LINT_FIX=0 ./tools/lint/scripts/lint.sh`): no package-level mutable globals — return regexes/sets from functions (see the `allowlist()`/`titleRe()` pattern in `tools/lint/linter/headingcase/headingcase.go`); max ~5 params per function (bundle into structs if needed); doc comments on exported identifiers; gofumpt formatting.

**Coverage gates.** Each package has a `.covgate` file with a minimum coverage percentage, enforced by `./tools/lint/scripts/covgate.sh` (default 90.0) in CI job `test-custom-linter`. Existing gates: `linter/.covgate` 93.3, `analysis/.covgate` 97.8, rule packages 100.0 (except `headingcase`, which has no `.covgate` and falls back to the 90.0 default), `tools/lint/.covgate` 30.0. The new package gets `tools/lint/linter/imagedomain/.covgate` containing `90.0`. `./tools/lint/scripts/ratchet-covgates.sh` can ratchet upward later.

**Fixture harness.** `tests/test-lint.sh` (invoked by `pnpm run test:lint`) drives the full `scripts/lint.sh` via `DOCS_LINT_ROOT` against mini content roots in `tests/lint-fixtures/{good,bad-mdx,bad-spelling,bad-openapi}/`, each containing `example.mdx`, `references/example.yaml`, and `snippets/.gitkeep`. Fixtures are enumerated explicitly at the bottom of the script (not globbed): `run_expect_pass "good"` requires exit 0 plus the output `All documentation lint checks passed.`; `run_expect_fail <fixture> <substring>` requires non-zero exit plus the substring. No fixture exercises a Go rule yet — `bad-image-domain` will be the first. `scripts/lint.sh` runs the Go linter first (`== MDX Prose ==` section, built via `cd tools/lint && go build -o lint .`), then ESLint, CSpell, OpenAPI; `set -euo pipefail` means a Go-rule failure stops the run there.

**CI** (`.github/workflows/ci.yml`; Go via `go-version-file: tools/lint/go.mod`): `changes` (dorny/paths-filter; `custom-linter` output on `tools/lint/**`); `lint` (every PR: `pnpm install --frozen-lockfile`, `pnpm run test:lint`, `./scripts/lint.sh`); `audit` and `shell-tests` (every PR); `lint-custom-linter` and `test-custom-linter` (when `tools/lint/**` changed).

## Plan of Work

**M1 — rule package and registration.** Create `tools/lint/linter/imagedomain/imagedomain.go` with:

    // Package imagedomain ... (doc comment must document the rule AND the
    // {/* lint-ignore image-domain */} directive — this doc comment is where
    // contributors learn about the suppression mechanism)
    func Check(file string, lines []string) []analysis.Violation

Algorithm, per line index i (0-based; reported Line = i+1):

1. Skip if inside/delimiting a code fence or frontmatter (tracking described in Context).
2. Skip if the previous line's trimmed content is exactly `{/* lint-ignore image-domain */}` (next-line-only suppression).
3. Collect candidates via two regexes (returned from unexported funcs, no globals):
   - markdown image: `!\[[^\]]*\]\(\s*([^)\s]+)` — group 1 is the URL (stops at whitespace, so optional titles are excluded); kind = markdown.
   - attribute: `\b(src|image|background|poster)\s*=\s*\{?\s*["']([^"']*)["']` — group 1 attribute name, group 2 URL; the optional `\{?` also matches JSX brace-quoted values like `src={"https://..."}` (group numbering unchanged, so Col via `FindAllStringSubmatchIndex` stays correct). Per-line, per-occurrence: multi-line JSX tags work because each attribute is matched on its own line.
4. Verdict predicate (the single source of truth; tests encode this table):
   - `allowed(url)` := url has prefix `https://assets.mirurobotics.com/`
   - `imageExt(url)` := url with any `?...`/`#...` suffix stripped, lowercased, ends in `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, or `.webp`
   - markdown / `image` / `background` / `poster` candidates: violation iff `!allowed(url)`
   - `src` candidates: violation iff `!allowed(url) && imageExt(url)` (so `.mp4` and other non-image `src` values are never flagged)
5. Violation: Col = 1-based byte offset of the URL's first character (use `FindAllStringSubmatchIndex`); Message = `image-domain: image must be hosted on https://assets.mirurobotics.com (got "<url>")`.

Inputs → expected verdicts:

| Line contains | Verdict |
|---|---|
| `![d](https://assets.mirurobotics.com/docs/a.png)` | ok |
| `![d](/images/a.png)` | violation |
| `![d](https://example.com/a.png)` | violation |
| `![d](./a.webp "title")` | violation (URL is `./a.webp`) |
| `<img src="https://assets.mirurobotics.com/docs/a.png" />` | ok |
| `src="/images/a.png"` (attribute on its own line) | violation |
| `src="http://assets.mirurobotics.com/a.png"` | violation (must be https) |
| `src="//assets.mirurobotics.com/a.png"` | violation (protocol-relative) |
| `src="https://assets.mirurobotics.com/docs/v.mp4"` | ok |
| `src="/videos/v.mp4"` | ok (`src` + non-image extension) |
| `image="/images/changelog/x.png"` | violation |
| `image="/images/x"` (no extension) | violation (`image` is strict) |
| `background="https://assets.mirurobotics.com/docs/bg.dark.svg"` | ok |
| `poster="/images/p.jpg"` | violation |
| `src={"/images/a.png"}` (JSX brace-quoted) | violation |
| `image="https://assets.mirurobotics.com/docs/releases/header:page.png"` | ok (colon tolerated) |
| a violating URL inside inline backticks | violation (accepted false positive; suppress with the directive) |
| any of the above inside a code fence or frontmatter | ok |
| any of the above on the line after `{/* lint-ignore image-domain */}` | ok |

Register in `tools/lint/linter/run.go`: add `RuleImageDomain Rule = "image-domain"` to the consts, to `AllRules()`, and to `ruleCheckers()`:

    {RuleImageDomain, func(in checkInput) []analysis.Violation {
        return imagedomain.Check(in.path, in.lines)
    }},

Create `tools/lint/linter/imagedomain/.covgate` containing `90.0` (newline-terminated, matching siblings).

**M2 — Go tests.** `tools/lint/linter/imagedomain/imagedomain_test.go`: table-driven subtests (mirror `nodoubledash/nodoubledash_test.go` style; no golden files), covering every row of the verdicts table plus: multiple candidates on one line (each flagged), suppression covers all violations on the suppressed line only (the line after that is flagged again), fence open/close toggling, frontmatter skipping, correct Line/Col/Message values. Also add one subtest to `TestProcessFile` in `tools/lint/linter/run_test.go` exercising image-domain through `ProcessFile`, following that file's existing inline style: write an MDX file into `t.TempDir()` with `os.WriteFile`, call `ProcessFile(path, dir)`, and scan the returned violations for one on the expected line whose Message starts with `image-domain:`. This covers the new closure in `ruleCheckers()` so `linter/.covgate` (93.3) still passes.

**M3 — fixtures.** Create `tests/lint-fixtures/bad-image-domain/` mirroring existing fixtures: `example.mdx` (frontmatter `title: "Fixture"` plus one line with `![Fixture](/images/fixture/example.png)`; keep all words cspell-clean and the MDX ESLint-parseable — the Go linter fails first anyway), `references/example.yaml` (copy of `good/references/example.yaml`), `snippets/.gitkeep`. In `tests/test-lint.sh`, add after the existing lines:

    run_expect_fail "bad-image-domain" "image-domain:"

Exercise the pass path: append to `tests/lint-fixtures/good/example.mdx` a compliant image and an ignore-directive usage, e.g.

    ![Fixture diagram](https://assets.mirurobotics.com/docs/fixture/example.png)

    {/* lint-ignore image-domain */}
    <img src="/images/fixture/legacy.png" alt="Legacy fixture" />

(`component-style`/`mdx-style` only inspect import statements, so raw `<img>` is safe here; verify `run_expect_pass "good"` still passes.)

**M4 — corpus verification and docs.** Sweep the real corpus (grep + full lint, commands below). Expected result: zero violations. If new violations appeared from a rebase/merge: migrate each URL to `https://assets.mirurobotics.com/...` if the asset already exists on the CDN; otherwise annotate with `{/* lint-ignore image-domain */}` and list those as follow-ups in the PR description. The ignore directive is documented in the package doc comment (M1); research found no contributor doc (CLAUDE.md/AGENTS.md/README) that lists lint rules, so no other doc updates are needed — double-check with a quick grep for a rule list before closing the milestone.

## Concrete Steps

All commands run from the repo root (`repos/docs` checkout) unless noted. One commit per milestone, conventional-commit style.

**M1**

    cd tools/lint && go build -o lint . && cd ../..

Expect: builds cleanly. Then meta-lint the Go code:

    LINT_FIX=0 ./tools/lint/scripts/lint.sh

Expect: no findings (use `LINT_FIX=1` to auto-fix formatting). Smoke-test the wiring (the binary needs a content root containing `snippets/`, so run against a fixture file; expect exit 0 here — real verdict coverage lands in M2):

    ./tools/lint/lint tests/lint-fixtures/good/example.mdx

Commit:

    git add tools/lint && git commit -m "feat(lint): add image-domain rule requiring assets.mirurobotics.com image URLs"

**M2**

    cd tools/lint && go test ./... && cd ../..
    ./tools/lint/scripts/covgate.sh

Expect: all packages pass; covgate reports every package at or above its gate, including `linter/imagedomain` ≥ 90.0 and `linter` ≥ 93.3. Commit:

    git add tools/lint && git commit -m "test(lint): cover image-domain rule and registration"

**M3**

    pnpm install --frozen-lockfile   # once, if node_modules is missing
    pnpm run test:lint

Expect: exits 0; `good` passes (including the new compliant image and ignore-directive lines) and `bad-image-domain` fails with output containing `image-domain:`. If `bad-image-domain` fails for a different reason (cspell/ESLint), fix the fixture wording, not the harness. Commit:

    git add tests && git commit -m "test(lint): add bad-image-domain fixture and good-fixture pass path"

**M4**

    grep -rnE '(src|image|background|poster)=[{]?"(/|\.|http://|//)' docs --include='*.mdx' || echo "no attribute violations"
    grep -rnE '!\[[^]]*\]\((/|\.|http://|//)' docs --include='*.mdx' || echo "no markdown violations"
    ./scripts/lint.sh

Expect: both greps print their "no ... violations" fallback and `./scripts/lint.sh` ends with `All documentation lint checks passed.` — the linter itself is the authority; the greps are a cross-check. If violations surface, remediate per M4 in Plan of Work, re-run `./scripts/lint.sh`, and commit:

    git add -A docs && git commit -m "docs: migrate image URLs to assets.mirurobotics.com for image-domain rule"

If the sweep changed nothing, there is nothing to commit for M4 — note that in Progress and move on.

**Validation**

    git push -u origin feat/image-domain-lint

Then run preflight (the `$preflight` skill) against the pushed head and open the PR as draft until it reports CLEAN.

## Validation and Acceptance

- `cd tools/lint && go test ./...` passes; `./tools/lint/scripts/covgate.sh` passes with `linter/imagedomain` ≥ 90.0.
- `LINT_FIX=0 ./tools/lint/scripts/lint.sh` reports no findings.
- `pnpm run test:lint` passes: `good` (with compliant image + ignore-directive usage) passes; `bad-image-domain` fails with `image-domain:` in the output.
- `./scripts/lint.sh` on the real corpus prints `All documentation lint checks passed.` — confirms docs.json chrome is untouched and zero corpus violations remain.
- Behavior spot-check: temporarily add `![x](/images/x.png)` to any docs page, run `./scripts/lint.sh`, expect a `image-domain: image must be hosted on https://assets.mirurobotics.com (got "/images/x.png")` line and non-zero exit; revert.
- CI: all six jobs (`changes`, `lint`, `audit`, `shell-tests`, `lint-custom-linter`, `test-custom-linter`) green on the pushed branch head; this change directly exercises `lint`, `lint-custom-linter`, and `test-custom-linter`. **Gate: preflight must report CLEAN (CI green on the pushed branch head) before the PR leaves draft or the task is reported complete.**

## Idempotence and Recovery

Every build/test/lint command above is read-only with respect to source and safe to re-run. `go build -o lint .` and fixture runs overwrite only generated artifacts (`tools/lint/lint` binary is gitignored). Fixture creation and file edits are plain working-tree changes: to retry a milestone, `git status` / `git diff` to inspect, `git checkout -- <path>` to discard, re-apply. Each milestone is an isolated commit, so a bad milestone can be dropped with `git reset --hard HEAD~1` (before push) or `git revert <sha>` (after push) without touching other milestones. If `pnpm run test:lint` fails from a stale environment, re-run `pnpm install --frozen-lockfile`. The corpus sweep (M4) is a pure verification pass and can be repeated at any time.

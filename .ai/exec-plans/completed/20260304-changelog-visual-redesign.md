# Product changelog visual redesign

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.


## Purpose / Big Picture

The product changelog (`docs/docs/changelogs/product.mdx`) currently uses Mintlify's built-in `<Update>` component, which renders a functional but visually unremarkable two-column layout (sticky date label on the left, content on the right). The improvements/fixes sections use custom `<Dropdown>` accordions that feel utilitarian. The overall page reads like documentation rather than a polished product announcement.

After this work, the product changelog will have a distinctive, editorial feel inspired by Linear and Cursor's changelogs — with better visual hierarchy, more generous spacing, refined typography, and a layout that makes each entry feel like a crafted announcement rather than a log entry.

**Scope:** Product changelog only (`docs/docs/changelogs/product.mdx`) and supporting components/CSS. No changes to CLI, Agent, Device API, or Platform API changelogs.


## Progress

- [x] (2026-03-04) Finalize design direction (approve this plan)
- [x] (2026-03-04) Phase 1: CSS refinements — entry spacing, dividers, date label muting, typography weights/sizes, media width, doc link styling, responsive breakpoints
- [x] (2026-03-04) Phase 2: Component refinements — Dropdown accepts `count` prop, muted uppercase header, hover background, smaller font; removed inter-dropdown Separators in MDX; added count props to all Dropdown instances
- [x] (2026-03-04) Phase 3: Evaluated `<Update>` component — CSS overrides sufficient, no custom replacement needed
- [x] (2026-03-04) Phase 4: QA — no rendering errors, all components present, other changelogs unaffected, RSS metadata preserved


## Surprises & Discoveries

- `mode: "custom"` was already set in frontmatter (from the prior layout exec-plan).
- The `<Update>` component provides excellent CSS hooks via `data-component-part` attributes and `.update-container` class — no need for a custom replacement.
- Mintlify renders component props client-side, so `curl` greps don't capture the rendered count text, but the props pass correctly in the SSR bundle.


## Decision Log

- Decision: Research-first approach — design the visual direction before writing code.
  Rationale: The user explicitly asked for research and design, not code. Getting alignment on the visual direction prevents rework.
  Date/Author: 2026-03-04 / Ben + Claude

- Decision: CSS-only approach on `<Update>` component (Option A from plan).
  Rationale: The built-in `<Update>` component provides sufficient DOM hooks (`data-component-part` attributes, `.update-container` class) for all targeted visual refinements. Keeps RSS/anchor support without reimplementation effort.
  Date/Author: 2026-03-04 / Claude

- Decision: Remove `<Separator>` between Improvements and Fixes dropdowns within each entry.
  Rationale: The accordion headers are visually distinct enough on their own; the separator added clutter rather than clarity, especially with the new muted uppercase styling.
  Date/Author: 2026-03-04 / Claude

- Decision: Keep a single `<Separator>` before the first dropdown in each entry.
  Rationale: This provides a clear visual boundary between the main content (hero image, description, doc links) and the secondary content (improvements/fixes). Removing it entirely made the transition too abrupt.
  Date/Author: 2026-03-04 / Claude


## Outcomes & Retrospective

All phases complete. Changes made:

1. **CSS refinements** (`docs/style.css`): Entry spacing with subtle 6% opacity dividers, muted date labels (45% opacity), refined H2/H3 typography (semibold 1.5rem / medium 1.15rem), body line-height 1.7, media `width: 100%` with margin, doc link subtle opacity treatment, responsive breakpoints.

2. **Dropdown component** (`docs/snippets/components/dropdown.jsx`): Added `count` prop for "Improvements (3)" / "Fixes (4)" display. Restyled header to muted uppercase `text-xs` with `font-medium`. Added subtle hover background (`bg-white/[0.04]`). Reduced chevron size and opacity.

3. **MDX content** (`docs/docs/changelogs/product.mdx`): Removed `<Separator>` between Improvements/Fixes dropdowns in all 4 entries. Added `count` props to all 8 Dropdown instances. Removed unnecessary `<div className="mb-9" />` spacers.

4. **No replacement needed** for Mintlify's `<Update>` component — CSS overrides are sufficient.

The changelog now has an editorial, magazine-quality feel closer to Linear/Cursor: more spacious, better hierarchy, progressive disclosure with item counts, and refined typography.


## Context and Orientation

### Current state

The product changelog was recently switched to `mode: "wide"` (see completed exec-plan `20260304-changelog-layout-redesign.md`). This removed the right-side table of contents and gave us CSS control via a `.changelog-page` wrapper class. Current CSS sets max-width 860px, centers the content, and adds responsive padding.

Key files:
- `docs/docs/changelogs/product.mdx` — the product changelog page, uses `mode: "wide"`, imports custom components
- `docs/style.css` — global CSS with `.changelog-page` scoped styles
- `docs/snippets/components/dropdown.jsx` — collapsible accordion for Improvements/Fixes
- `docs/snippets/components/framed.jsx` — styled image frame with background overlay
- `docs/snippets/components/separator.jsx` — thin horizontal divider
- `docs/docs.json` — Mintlify config (theme: maple, font: Inter weight 450, colors: emerald green)

### Current entry structure (per `<Update>`)

Each entry renders as:
1. **Sticky date label** on the left (~160px wide, Mintlify's `<Update>` component)
2. **Content column** on the right containing:
   - H2 heading (feature name)
   - Description paragraphs
   - Media (screenshots via `<Framed>`, videos via `<video>`)
   - Documentation links
   - `<Separator>` thin line
   - `<Dropdown title="Improvements">` — collapsible list of improvements
   - `<Separator>`
   - `<Dropdown title="Fixes">` — collapsible list of fixes

### Mintlify constraints

- Mintlify uses the Maple theme with the left-panel navigation. `mode: "wide"` preserves this.
- JSX components in `snippets/components/` cannot have module-level `const` declarations — all variables must be inside the function body.
- Custom CSS is global via `style.css`, scoped with class selectors.
- The `<Update>` component is a Mintlify built-in — we cannot modify its internals, only style it via CSS or replace it with a custom component.
- Mintlify supports React/JSX components imported into MDX.


## Research: Changelog designs analyzed

### Linear (https://linear.app/changelog)

**Strengths we want to draw from:**
- Editorial, magazine-quality feel — each entry reads like a crafted announcement
- Large, artistic hero images (abstract graphics, not raw screenshots)
- Very generous whitespace — entries breathe
- Benefit-driven writing ("your dashboard loads 2x faster" not "improved performance")
- Inline subsections (Improvements, Fixes, API) — not collapsed, displayed as labeled bullet groups
- Dark theme, single-column, date-based grouping
- Embedded videos with playback controls

**What makes it exceptional:** Treats the changelog as a publication. The combination of artistic imagery, narrative writing, and spacious layout creates engagement. Linear reportedly achieves ~60% monthly changelog engagement vs. 10-15% industry average.

### Cursor (https://cursor.com/changelog)

**Strengths we want to draw from:**
- Two-column layout with sticky date sidebar — always know which entry you're reading
- Collapsible accordion sections with item counts (e.g., "Desktop Improvements (15)") — progressive disclosure
- Looping video embeds for feature demos
- Three-tier typography system (display font for headings, serif for body, mono for code)
- Subtle dividers (thin, low-contrast borders) — whitespace does the heavy lifting
- Named spacing tokens creating consistent vertical rhythm

**What makes it exceptional:** Maximum content density with progressive disclosure. Headlines are scannable, details are one click away. The sticky date is a great UX detail.

### Other notable designs

- **Vercel:** Category-filtered navigation, author avatars, doubles as marketing surface
- **Raycast:** Version-based entries, emoji category prefixes (🎁 new, 💎 improved, 🐛 fixed) for instant visual scanning
- **Figma:** Sophisticated variable font system, product GIFs, weekly cadence
- **Resend:** Dark theme, contributor hover cards, command palette search
- **Webflow:** Card-based grid with hero card, category/subcategory badges
- **Liveblocks:** Custom SVG illustrations per week, domain-based categorization


## Design Specification

### Design philosophy

Take inspiration primarily from **Cursor's layout** (two-column with sticky dates, progressive disclosure via accordions) combined with **Linear's editorial quality** (generous spacing, polished imagery, benefit-driven writing tone). The result should feel like a curated product announcement feed, not a documentation page.

### Layout

**Two-column, sticky date + content** — Keep the current two-column structure (Mintlify's `<Update>` already does this), but refine it:

- **Date column (left):** Clean date display. Consider whether the Mintlify `<Update>` component's date rendering is sufficient or if we need a custom component. If we replace `<Update>`, the date column should be ~140-160px, sticky, with the date in a clean sans-serif at a comfortable size. No badge/pill — just the date as text.
- **Content column (right):** Wider, with max-width giving entries room to breathe. Media should fill the full content column width.

**Option A — Enhance Mintlify's `<Update>` via CSS:** Override styles on `.update`, `.update-container`, date label, and content wrapper. Least effort, keeps built-in RSS/anchor support.

**Option B — Replace `<Update>` with a custom `<ChangelogEntry>` component:** Full visual control. Must re-implement sticky date, anchor IDs, and ensure RSS still works with `rss: "true"` frontmatter.

Recommendation: **Start with Option A** (CSS overrides on `<Update>`). If the built-in component's DOM structure is too limiting, pivot to Option B.

### Entry visual structure

Each changelog entry should follow this visual flow:

```
┌─────────────────────────────────────────────────────────┐
│                                                         │
│  Jan 27, 2026          Feature Title                    │
│  (sticky, left)        ─────────────────                │
│                        Description paragraph with       │
│                        benefit-driven language.          │
│                                                         │
│                        ┌─────────────────────────────┐  │
│                        │                             │  │
│                        │     Hero image / video      │  │
│                        │     (full content width)    │  │
│                        │                             │  │
│                        └─────────────────────────────┘  │
│                                                         │
│                        Read the full documentation »    │
│                                                         │
│                        ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐  │
│                        │ ▸ Improvements (4)          │  │
│                        │ ▸ Fixes (3)                 │  │
│                        └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘  │
│                                                         │
│─────────────────────────────────────────────────────────│
│                     (subtle divider)                    │
│                                                         │
│  Dec 19, 2025          Next Feature Title               │
│                        ...                              │
└─────────────────────────────────────────────────────────┘
```

### Typography

Stay with Inter (already configured in Mintlify). Refine weights and sizes:

- **Entry title (H2):** Semibold (600), ~1.5rem. Clean, not oversized.
- **Sub-headings (H3):** Medium (500), ~1.15rem.
- **Body text:** Regular (450, already set), ~1rem with comfortable line-height (~1.7).
- **Date label:** Regular weight, slightly muted color (e.g., `text-gray-400` / `rgba(255,255,255,0.5)` in dark mode).
- **Dropdown headers:** Medium weight, uppercase or small-caps optional. Item count in parentheses.
- **Improvement/fix items:** Regular weight. Bold prefix for the item name, regular for description (already done — keep this pattern).

### Spacing and rhythm

The biggest impact change. Current spacing feels tight. Target:

- **Between entries:** ~4rem (64px) vertical gap, with a subtle 1px divider at the midpoint
- **Between entry title and first paragraph:** ~1rem
- **Between paragraphs:** ~1.25rem
- **Between text and media:** ~1.5rem above and below media
- **Between media and documentation link:** ~1rem
- **Between doc link and improvements/fixes section:** ~1.5rem
- **Padding inside accordion headers:** comfortable click target (~12px vertical)

### Dividers

Replace the current `<Separator>` (which is a `border-t border-white opacity-10`) between entries with:

- **Between entries:** A single 1px line at ~10% opacity, centered in a ~4rem gap. Or no line at all — just whitespace. Linear uses no visible dividers, relying on spacing alone. Cursor uses very subtle bottom borders.
- **Between improvements/fixes:** Remove the `<Separator>` between these. The accordion headers are visually distinct enough.

Recommendation: **Remove dividers between improvements/fixes.** Keep a subtle divider between entries (or rely on spacing alone — test both).

### Accordion / Dropdown refinement

Current `<Dropdown>` uses a right-pointing chevron that rotates 90° on open. Improvements:

- **Add item counts** to accordion headers: "Improvements (4)" / "Fixes (3)" — like Cursor. This tells users whether it's worth expanding.
- **Tighten visual weight:** The accordion header should feel like a secondary element, not compete with the entry title. Smaller font size, muted color.
- **Subtle background on hover:** A very light background tint on hover to show clickability.
- **Consider starting expanded** for the latest/most recent entry, collapsed for older entries. This gives the newest update full visibility while keeping the page scannable for history.

### Media presentation

Current `<Framed>` component adds a decorative background behind images. Current `<video>` uses autoPlay, loop, controls.

- **Keep `<Framed>`** — it already provides nice visual treatment for screenshots.
- **Videos:** Keep autoPlay and loop. Consider adding `muted` and `playsInline` for better mobile behavior if not already present. The rounded corners (`rounded-xl`) are good.
- **Media width:** Should fill the full content column width. Already handled by `.changelog-page img, video { max-width: 100% }` — may need to set `width: 100%` to ensure they stretch to fill.
- **Consider hero treatment for the lead image/video:** The first media element in each entry could have slightly more visual prominence (e.g., slightly larger border radius, a subtle shadow, or a gradient border effect).

### Color and theme

The docs already use a dark theme (Mintlify default dark appearance) with emerald green accents (#059669 primary, #34d399 light, #065f46 dark). Keep this palette.

- **Accent color for links:** The emerald green for documentation links ("Read the full documentation »") is a good signature element.
- **Date label color:** Muted — `rgba(255,255,255,0.45)` or similar.
- **Accordion header color:** Slightly muted, lighter than body text.
- **Divider color:** `rgba(255,255,255,0.08)` — barely visible.

### Documentation link styling

The "Read the full documentation »" links are a distinctive element. Consider:
- Keeping the `»` suffix (it's a nice touch)
- Emerald green color with subtle hover underline
- Slightly smaller than body text to feel like a footnote/reference

### Responsive behavior

- **Desktop (>1024px):** Full two-column layout. Sticky date on left, content on right.
- **Tablet (768-1024px):** Date moves inline above content (stacked). Content full-width.
- **Mobile (<768px):** Same as tablet. Reduce heading sizes, tighten padding.

This is largely handled by Mintlify's `<Update>` component (which switches from `lg:flex-row` to `flex-col` at the `lg` breakpoint), but the CSS refinements should respect this.


## Plan of Work

### Phase 1: CSS-first refinements (low risk, high impact)

Target the `.changelog-page` wrapper and Mintlify's `.update` / `.update-container` class hierarchy to improve spacing, dividers, typography, and overall feel — without touching any components or MDX content.

Changes in `docs/style.css`:
- Increase vertical gap between `.update-container` elements
- Refine date label styling (muted color, consistent sizing)
- Add subtle entry divider (or remove existing and use spacing only)
- Improve media sizing (width: 100%)
- Tune typography weights/sizes for headings within `.changelog-page`

### Phase 2: Component refinements

Update `docs/snippets/components/dropdown.jsx`:
- Accept and display an item count
- Muted color for header text
- Subtle hover background
- Slightly smaller font

Update `docs/snippets/components/separator.jsx`:
- Reduce opacity or adjust for the refined spacing
- May become unnecessary if spacing alone provides sufficient visual separation

Update `docs/docs/changelogs/product.mdx`:
- Remove `<Separator>` elements between `<Dropdown>` sections within the same entry
- Add item counts to `<Dropdown>` components (e.g., `<Dropdown title="Improvements" count={4}>`)
- Minor copy edits for benefit-driven language (optional, scope-dependent)

### Phase 3: Evaluate `<Update>` component

After CSS and component refinements, assess whether Mintlify's `<Update>` provides sufficient visual quality or if a custom `<ChangelogEntry>` replacement is needed. If so:

- Create `docs/snippets/components/changelog-entry.jsx`
- Implement sticky date, content column, anchor IDs
- Migrate `product.mdx` to use the new component
- Verify RSS and navigation still work

### Phase 4: QA and polish

- Visual comparison against Linear/Cursor screenshots
- Test on mobile/tablet viewports
- Verify all existing functionality (RSS, component rendering, navigation)
- Verify other changelog pages are unaffected


## Concrete Steps

Steps will be filled in during implementation. For now, the design specification above is the deliverable.

### Verify Mintlify's `<Update>` DOM structure

Before writing CSS, inspect the rendered HTML to identify the exact class names and DOM hierarchy of the `<Update>` component. This determines what CSS selectors are available.

From `docs/`:

    npx mintlify dev

Then inspect the page in a browser or:

    curl -s http://localhost:3000/docs/changelogs/product | grep -oP 'class="[^"]*update[^"]*"' | sort -u

### CSS refinements

Edit `docs/style.css`, targeting selectors within `.changelog-page`. Specific selectors will depend on the DOM inspection above.

### Component updates

Edit `docs/snippets/components/dropdown.jsx` to support the `count` prop and refined styling.

### MDX updates

Edit `docs/docs/changelogs/product.mdx` to remove unnecessary `<Separator>` elements and add `count` props to `<Dropdown>`.


## Validation and Acceptance

1. **Visual quality:** The product changelog should feel closer to Linear/Cursor in editorial quality — more spacious, better hierarchy, refined details.
2. **No regressions:** All existing components (`<Update>`, `<Dropdown>`, `<Framed>`, `<Separator>`, `<video>`) render correctly.
3. **RSS works:** `rss: "true"` frontmatter preserved, feed generates at build time.
4. **Other changelogs unaffected:** CLI, Agent, Device API, Platform API pages render normally.
5. **Responsive:** Page looks good on desktop, tablet, and mobile viewports.
6. **No Mintlify errors:** Page loads without `__next_error__`.


## Idempotence and Recovery

All changes are CSS overrides and component edits — fully reversible:
- CSS additions in `style.css` are scoped under `.changelog-page` and `.update`-prefixed selectors.
- Component changes in `dropdown.jsx` add optional props (backward compatible).
- MDX changes (removing `<Separator>`, adding `count` props) can be reverted by re-adding separators and removing count props.
- If the `<Update>` CSS overrides prove too fragile, they can be removed and we fall back to the current appearance.

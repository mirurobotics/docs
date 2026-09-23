---
name: campaign
description: Turns a product changelog entry from docs/changelog/product.mdx into a designed draft Loops email campaign in the "Updates" group, using Miru's changelog email design system. Subcommand required. `plan <entry>` designs the email without touching Loops; `create <entry>` designs it, builds the draft, uploads images and sends a preview; `revise` applies feedback to an existing draft. Uses the Loops CLI. Use when the user asks to turn a changelog entry into an email or campaign, draft or revise the changelog email, or set up a Loops campaign for a release date.
allowed-tools: Bash(loops:*) Bash(curl:*) Bash(mktemp:*) Bash(git:*) Read Grep Glob Write
---

# campaign

## Usage

```
/campaign <command> [args]

Commands:
  plan <entry>       Read-only. Design the email for a changelog entry and
                     print the layout plan: name, subject, preview text,
                     opener or not, each section with its image. Nothing is
                     created or uploaded.
  create <entry>     Design the email, create the draft in the Updates group,
                     upload images (videos get a placeholder), push the
                     content and send a preview. Stops at the draft.
  revise [campaign] [feedback]
                     Pull the draft's current content from Loops, apply the
                     feedback, push it back and re-send the preview.
                     [campaign] is a Loops campaign ID or a quoted campaign
                     name. Omitted, it is the campaign this chat last created
                     or revised.

<entry> is free text naming one <Update> block in docs/changelog/product.mdx:
a date, a feature heading, or both ("sep 23 cli commands").
```

## Arg parsing and routing

Split args on whitespace. The first token is the **command**; the rest is its argument text.

Validate before routing. Print the usage block and stop if args are empty, the first token isn't `plan`, `create` or `revise`, or `plan` / `create` has no `<entry>`. For `revise`, a leading Loops ID or quoted name selects the campaign and the rest is feedback. Otherwise all the text is feedback and the campaign is the one this chat last touched. If there isn't one, print the usage block and stop. Don't guess: a wrong guess designs the wrong entry or overwrites the wrong draft.

**Resolving `<entry>`:** read `docs/changelog/product.mdx` and match the text against `<Update label>` dates and `##` headings. Two blocks can share a date (September 23, 2026 has "Device Views" and "CLI Commands"), so a bare date can be ambiguous. If more than one block matches, list them and ask.

## Rules

- **Never publish, schedule, or set an audience** (`--mailing-list-id`, `--audience-*`, `--schedule-*`). The user does those in the dashboard.
- **Only touch Draft campaigns.** If the target isn't `Draft`, stop.
- **Never edit the changelog.** Never write the `.lmx` or downloaded images into the repo. Work in `scratch=$(mktemp -d)`.
- **Never print an API key** (`loops auth get`).

## Prerequisites

1. `loops api-key` succeeds. If `loops` is missing: `brew install loops-so/tap/loops`. If it isn't authenticated, stop and ask the user to run `loops auth login miru`, which prompts for the key interactively.
2. Loops' `loops-cli` and `loops-lmx` skills sit beside this one in `.cursor/skills/`. Read `loops-lmx` before writing LMX. Use `loops agent-context` for exact flags rather than guessing. Don't run `loops skill install`: it installs into `.agents/skills/` and other agents' directories, not `.cursor/skills/`.
3. A `401` "content API not enabled for this team" means the team lacks Content API access. Stop and tell the user.

## Design

`plan`, `create` and `revise` all run this.

### Design system

Every changelog email uses exactly this `<Style />`, first in the document:

```xml
<Style backgroundColor="#f5f5f5" backgroundXPadding="0" backgroundYPadding="24" bodyColor="#ffffff" bodyXPadding="36" bodyYPadding="36" bodyFontFamily="Geist" bodyFontCategory="sans-serif" borderColor="#000000" borderWidth="0" borderRadius="8" buttonBodyColor="#047857" buttonBodyXPadding="18" buttonBodyYPadding="10" buttonBorderColor="#047857" buttonBorderWidth="1" buttonBorderRadius="12" buttonTextColor="#ffffff" buttonTextFormat="0" buttonTextFontSize="16" dividerColor="#f5f5f4" dividerBorderWidth="1" textBaseColor="#262626" textBaseFontSize="16" textBaseLineHeight="150" textBaseLetterSpacing="0" textLinkColor="#047857" heading1Color="#171717" heading1FontSize="30" heading1LineHeight="125" heading1LetterSpacing="0" heading2Color="#171717" heading2FontSize="20" heading2LineHeight="130" heading2LetterSpacing="0" heading3Color="#171717" heading3FontSize="17" heading3LineHeight="130" heading3LetterSpacing="0" />
```

Block recipes. Spacing comes from explicit padding, never from empty spacer paragraphs:

| Block | LMX |
| --- | --- |
| Logo | `<Image src="https://images.vialoops.com/cmi3suhby0xmfzs0iyjvi8ouc/cmueo2wkr2u970jx3fn62eojx.png" alt="Miru" width="100" paddingTop="0" paddingBottom="24" />` |
| Title | `<H1 paddingTop="0" paddingBottom="4"><Strong>…</Strong></H1>` |
| Date | `<Paragraph fontSize="14" paddingBottom="24"><Text textColor="#646464">September 23, 2026</Text></Paragraph>` |
| Opener | `<Paragraph fontSize="18" lineHeight="155">…</Paragraph>` |
| Section heading | `<H2 paddingTop="40" paddingBottom="12"><Strong>…</Strong></H2>` (the first section after the date, with no opener, uses `paddingTop="0"`) |
| Subsection heading | `<H3 paddingTop="24" paddingBottom="8"><Strong>…</Strong></H3>`, only under an `<H2>` |
| Image | `<Image src="…" alt="…" width="528" borderRadius="12" paddingTop="0" paddingBottom="16" />` |
| Intro paragraph | `<Paragraph paddingBottom="8">…</Paragraph>` before a list, `paddingBottom="12"` before a code block, no padding otherwise |
| List | `<UnorderedList>` with `<ListItem paddingBottom="6">`. The list itself takes no padding or `align` (both rejected). Space above a list goes on its first item (`paddingTop="12"` after a code block) |
| Code block | `<CodeBlock paddingRight="16" paddingLeft="16">…</CodeBlock>`. This padding is inside the gray box |
| Note after a list | `<Paragraph paddingTop="6">…</Paragraph>` |
| Button | `<Button href="…" bgColor="#047857" textColor="#ffffff" borderRadius="12" borderColor="#047857" paddingTop="32" paddingBottom="24">Open changelog</Button>` |

Why these values:

- **One left edge.** No side padding on copy: headings, text, lists, code and images all align.
- **Type:** 16px text everywhere, including list items, with an 18px lede. Section headings are 20px at 130% line height, `<H2>` never skipped for `<H3>`.
- **Colors:** near-black text (`#262626`, headings `#171717`), gray `#646464` for the date. The green is `#047857`, the lightest emerald that passes WCAG AA contrast with white text (5.3:1); `#059669` fails (3.6:1).
- **Rhythm:** spacing groups content. A heading sits 12px from its own content and 40px from the previous section.
- **Logo:** black with a white outline, so it survives dark mode in Apple Mail, which darkens the background but leaves images alone. The source is [logo.png](logo.png) (312px wide, displayed at 100). Re-upload it only if the Loops URL stops resolving.
- **Fonts:** Geist loads only in Apple Mail and Samsung Mail. Gmail and Outlook fall back to the system sans-serif, so don't tighten letter spacing.
- **Images:** screenshots display at 528px, so they need to be at least 1056px wide to stay sharp on high-density screens. Flag narrower sources in the report.

`logo.png` and `video-placeholder.png` are gitignored and exist only on the owner's machine. If one is missing, ask the user for it rather than recreating it.

Loops strips attributes that equal its defaults (`paddingBottom="4"` on a heading, zero padding on a paragraph), so a pulled draft may lack values you pushed. That is not a lost edit.

### Reference emails

`Changelog: Schema Validation` and `Changelog: Device Views` show structure: an opener vs. opening on an image, and flat sections vs. one headline feature with subsections. Read their full LMX with `loops email-messages get` for layout ideas and the sender fields (`fromName`, `fromEmail`, `replyToEmail`). Don't copy their styling: their inset copy, spacer paragraphs, `#059669` green, black logo and `align="start"` lists are superseded by the design system above.

### Think like a marketer and designer

Act as Miru's best product marketer and email designer. The reader is a robotics engineer or fleet operator skimming an inbox. The email has to make them want to try the feature, then get out of the way. Decide the layout from the content. Don't fill in a fixed shape.

**Building blocks**

- **Header:** the logo, the headline as `<H1>`, then the entry's date, per the design system.
- **Opening paragraph** (Device Views): a short lede that frames the whole release. Use it when the entry bundles several features that need tying together, or when the first section's copy doesn't hook on its own. Skip it and open straight on the first section's image when that section already sells itself (Schema Validation). The opener names the release's features at a high level. It must not repeat a sentence a section uses as its intro.
- **Feature section:** heading, then image, then an intro paragraph, then bullets. The image leads its section, even though the changelog puts media after the text.
- **Every section opens with a paragraph.** Never go straight from a heading (or image) into bullets. If the changelog gives the section no intro, use the entry's sentence that introduces it (the CLI section's "The Miru CLI can now clone, stage, and validate configs from your local machine").
- **Heading levels:** sections are `<H2>`. Keep them flat when the entry mixes one feature's `###` subsections with other `##` features (Device Views put "Filter and search" beside "Archiving"). Use `<H3>` only when several `##` features each have real subsections.
- **Text-only section:** a section without an image or video stays text. Don't invent imagery. A fenced command (for example `ansible-galaxy collection install …`) becomes a `<CodeBlock>` and serves as that section's visual anchor.

**Editorial judgment**

- **Lead with the strongest feature.** The changelog's order is usually right, but reorder if another feature is the better hook. The subject line names the headline feature, and the preview text names the next one or two.
- **Section copy is the changelog's wording.** Only tighten words (like "Filters" → "Filter"). Lists are copied item for item, exactly as the docs write them: never regroup, merge or summarize bullets. Write new copy only for the subject, preview text and opening paragraph, in the changelog's plain, specific voice.
- **Keep the rhythm.** Alternate image and text. Never stack two images, and never put a heading directly on top of another heading.
- **Cut what doesn't earn its place.** Improvements and Fixes dropdowns usually stay on the changelog behind the call to action. Pull one in only if it would matter to most readers.
- **No links in the body.** Drop every changelog link, both inline ones (command names become plain `<Code>`) and `[Label »](/path)` "read more" ones. One exception: a package the reader installs links to its registry page, the way `mirurobotics.agent` links to Ansible Galaxy. Style it `<Link href="…"><Code textColor="#047857">name</Code></Link>` so it reads as a link. Otherwise the only link is the `Open changelog` button, pointing at the entry's first `##` heading: `https://docs.mirurobotics.com/changelog/product#<heading-slug>` (for example `#cli-commands`).

### The layout plan

Write the plan before any LMX:

```
Name:     Changelog: <short headline feature name>
Subject:  …
Preview:  …
Opener:   yes / no, and why
Sections: 1. <heading>, <image | video placeholder | code block | text only>
          2. …
Cut:      <what stays on the changelog, and why>
```

Name the campaign after the headline feature, shortened ("Live Schema Validation" shipped as "Changelog: Schema Validation"), unless the user gives a name.

## plan

Run **Design** and print the layout plan in the chat. Write no file, make no Loops writes, and upload nothing.

## create

1. **Check for a clash.** If a campaign with the planned name exists, stop and suggest `revise`.
2. **Design.** Run **Design** and keep the layout plan.
3. **Upload images.** LMX `<Image src>` must be Loops-hosted, so `assets.mirurobotics.com` URLs can't be used directly.
   - **`<Framed image="…">`:** download it into `$scratch`, then `loops uploads create "$scratch/<file>" -o json` → `finalUrl`.
   - **`<LazyVideo>`:** upload [video-placeholder.png](video-placeholder.png) once and reuse its `finalUrl` for every video in the email. Give each one `alt="Placeholder: <the video's alt text>"`. The user swaps in a still in the Loops editor.
   - Limits: png/jpeg/gif/webp only, 4 MB each, 50 uploads per day per team. Write alt text that describes what each real image shows.
4. **Write `$scratch/email.lmx`.** Follow the `loops-lmx` output checklist. Where it conflicts with the design system on sizes, colors or spacing, the design system wins.
5. **Create the draft and push the content:**

   ```bash
   loops campaigns create --name "<name>" --campaign-group-id <updates-group-id> -o json
   loops email-messages update <emailMessageId> \
     --expected-revision-id <emailMessageContentRevisionId> \
     --subject "…" --preview-text "…" \
     --from-name "…" --from-email "<username only>" \
     --lmx-file "$scratch/email.lmx" -o json
   ```

   Get the Updates group ID from `loops campaign-groups list -o json`, and copy the reply-to address from the reference (flag name per `loops agent-context`). Create exactly one campaign. If a later step fails, retry against this draft. On `422`, fix the LMX the error names and retry. On a stale revision (`409`), re-run `email-messages get` for a fresh one. Surface any `warnings`.
6. **Preview and review.** Run `loops email-messages preview <emailMessageId> --email "$(git config user.email)"`, passing `--contact-prop KEY=value` for any `{contact.*}` variable. Re-read the pushed LMX. Check it against the layout plan and the design system's block recipes, and check that any `warnings` are resolved. Then fix anything that reads cramped, repetitive or top-heavy.
7. **Report.**

## revise

1. **Load the draft.** `loops campaigns get <id>` must show `Draft`. `loops email-messages get <emailMessageId> -o json` → the current `lmx`, subject, preview text and `contentRevisionId`.
2. **Treat the pulled content as the truth.** The user may have edited it in Loops, for example by replacing video placeholders with stills. Keep those edits and change only what the feedback asks for.
3. **Apply the feedback** with the Design brief. Re-read the changelog entry if the feedback touches content.
4. **Push and preview.** Run `email-messages update` with `--expected-revision-id <contentRevisionId>`. On `409`, the draft changed while you worked: pull it again and re-apply. Then re-send the preview as in `create`.
5. **Report**, including what changed.

## Report

- The campaign URL (the `url` field from `loops campaigns get <id>`), name, subject and preview text
- The layout plan and the reason behind each call (opener or not, section order, anything cut)
- Each video placeholder, with the section it sits in and its source `.mp4`, so the user knows which still to capture
- What's left for the user in the dashboard: replace the placeholders, review the preview, pick the audience, publish

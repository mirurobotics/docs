This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

# Docs: SSE Events Coverage for Device API v0.2.1 and Agent v0.8.0

## Scope

This plan lives in the docs repo at `/home/ben/miru/workbench3/docs/plans/backlog/20260404-docs-sse-events-coverage.md`. The docs repo is the sole read-write repository; no other repos are modified. The agent repo (commit `0db044a`) and Device API v0.2.1 OpenAPI spec are read-only references for understanding the feature being documented.

All file paths in this plan are absolute unless stated otherwise. The docs repo root is `/home/ben/miru/workbench3/docs/`.

## Purpose / Big Picture

Agent v0.8.0 ships Server-Sent Events (SSE) streaming for deployment lifecycle events via `GET /v0.2/events`. Device API v0.2.1 formalizes this endpoint and its event schemas. After this work is complete, a developer visiting the Miru documentation will be able to:

- Read changelogs explaining what shipped in Agent v0.8.0 and Device API v0.2.1.
- Follow a dedicated events guide that explains how to connect to the SSE stream, interpret event frames, replay missed events with cursors, filter by type, and handle errors.
- See updated architecture documentation showing how events fit into the agent sync cycle.
- Find correct version entries in the agent versions table, the Device API supported versions table, and the agent compatibility matrix.
- Navigate to the events guide from the sidebar and from cross-links on related pages.
- Use the Device API Reference v0.2.1 dropdown (already wired) to browse the `GET /events` endpoint and event type reference pages (already authored).

## Progress

- [ ] M1: Device API v0.2.1 changelog entry
- [ ] M1: Agent v0.8.0 changelog entry
- [ ] M1: Commit M1
- [ ] M2: New events guide page
- [ ] M2: Update Device API overview to mention SSE
- [ ] M2: Update agent architecture page with event emission
- [ ] M2: Add v0.8.x row to agent versions table
- [ ] M2: Add v0.8.x and v0.2.1 rows to Device API versioning page
- [ ] M2: Commit M2
- [ ] M3: Add events guide to docs.json navigation
- [ ] M3: Update Device API latest redirect to v0.2.1
- [ ] M3: Run spell-check and update cspell.json
- [ ] M3: Commit M3

## Surprises & Discoveries

(None yet.)

## Decision Log

| Date       | Decision | Rationale |
|------------|----------|-----------|
| 2026-04-04 | v0.2.1 has no breaking changes, so the changelog omits a "Breaking changes" section and a "Migration steps" section. | Patch releases are always backward-compatible per the versioning page's stability guarantees. Including empty sections would be misleading. |
| 2026-04-04 | The events guide is placed at `docs/developers/device-api/events.mdx` rather than under Learn or References. | It is developer-facing how-to content about consuming an API. The Device API developer section is where similar pages (authn, sdks, versioning) already live. |
| 2026-04-04 | Milestones are split changelogs / content / nav+spellcheck so each commit is self-contained and reviewable. | Changelogs are independent of the guide content, and navigation wiring should only happen after the pages it references exist. |
| 2026-04-04 | The `docs.json` redirect for Device API latest is updated from v0.2.0 to v0.2.1. | v0.2.1 is the newest version and should be the target of the "latest" alias so new readers see the most complete reference. |

## Outcomes & Retrospective

(To be filled at completion.)

## Context and Orientation

A reader who knows nothing about the codebase needs the following context.

**Docs repo** (`/home/ben/miru/workbench3/docs/`): A Mintlify documentation site. Pages are `.mdx` files (Markdown with JSX components). Navigation structure is defined in `docs.json`. The repo is currently on the `feat/sse` branch.

**Mintlify docs.json**: The central configuration file at `/home/ben/miru/workbench3/docs/docs.json`. It defines products, navigation groups, pages arrays, redirects, and OpenAPI-backed API reference dropdowns. Navigation entries are strings like `"docs/developers/device-api/overview"` (no file extension). Redirects use a `source`/`destination` pattern with `:slug*` wildcards.

**Device API Reference v0.2.1**: Already wired in `docs.json` as a dropdown under "Device API Reference" with groups for Agent, Deployments, Device, Events, Git Commits, and Releases. The OpenAPI spec lives at `/home/ben/miru/workbench3/docs/docs/references/device-api/v0.2.1/api.yaml`. Event type reference pages exist at:
- `/home/ben/miru/workbench3/docs/docs/references/device-api/v0.2.1/events/deployment-deployed.mdx`
- `/home/ben/miru/workbench3/docs/docs/references/device-api/v0.2.1/events/deployment-removed.mdx`

**Device API changelog** (`/home/ben/miru/workbench3/docs/docs/changelogs/device-api.mdx`): Uses imports for `Dropdown`, HTTP method badges (`GET`, `POST`, etc.), `Separator`, and `DeviceApiReleaseLinks`. Entries are `# vX.Y.Z` headings in reverse chronological order. Current top entry is `# v0.2.0`.

**Agent changelog** (`/home/ben/miru/workbench3/docs/docs/changelogs/agent.mdx`): Entries are `# vX.Y.Z` headings in reverse chronological order with `*Date*` subheading, prose intro, and `## Features` / `## Improvements` / `## Fixes` sub-sections. Current top entry is `# v0.7.0`.

**Device API overview** (`/home/ben/miru/workbench3/docs/docs/developers/device-api/overview.mdx`): 14-line file. Opening paragraph describes the Device API as "a REST API." Needs to also mention SSE streaming.

**Agent architecture** (`/home/ben/miru/workbench3/docs/docs/developers/agent/architecture.mdx`): Contains a "Sync cycle" section with a numbered list (steps 1-5) and a Mermaid sequence diagram. The "Deployments" section follows. There is no mention of events.

**Agent versions** (`/home/ben/miru/workbench3/docs/docs/developers/agent/versions.mdx`): Contains a Markdown table of supported versions. Current top row is `v0.7.x` released 2026-03-13.

**Device API versioning** (`/home/ben/miru/workbench3/docs/docs/developers/device-api/versioning.mdx`): Contains an "Agent compatibility matrix" table (Agent Version -> API Versions) and a "Supported versions" table (API Version -> Released -> Status). Current top rows are `v0.7.x` -> `v0.2.0` and `v0.2.0` -> 2026-03-13.

**cspell.json** (`/home/ben/miru/workbench3/docs/cspell.json`): Spell-check dictionary with 21 custom words. No SSE-related terms present.

**Key terms**:
- SSE (Server-Sent Events): An HTTP-based protocol where the server pushes events to the client over a long-lived connection. Each event frame has optional `id`, `event`, and `data` fields separated by newlines.
- Cursor: An integer event ID used to replay events. Passed via `?after=<id>` query parameter or `Last-Event-ID` HTTP header.
- JSONL: JSON Lines format, one JSON object per line. Used by the agent to persist events to disk.

## Plan of Work

The work proceeds in three milestones. Each milestone ends with a commit.

### M1: Changelogs

First, add the Device API v0.2.1 changelog entry. Open `/home/ben/miru/workbench3/docs/docs/changelogs/device-api.mdx`. Insert a new `# v0.2.1` section above the existing `# v0.2.0` line (line 12). The new section follows the established format: version heading, italicized date, prose introduction, `<DeviceApiReleaseLinks>` component, then a "New endpoints" section with a Dropdown for Events containing `GET /events`, then an "Additive changes" section describing the two event types, cursor-based replay, and type filtering. There are no breaking changes and no migration steps, so those sections are omitted. End with a horizontal rule (`---`) to separate from v0.2.0.

Second, add the Agent v0.8.0 changelog entry. Open `/home/ben/miru/workbench3/docs/docs/changelogs/agent.mdx`. Insert a new `# v0.8.0` section above the existing `# v0.7.0` line (line 7). The new section follows the established format: version heading, italicized date, prose introduction that links to the Device API v0.2.1 changelog, then a "Features" section with bullets covering the SSE endpoint, cursor replay, type filtering, at-least-once delivery, and heartbeat keep-alive, then an "Improvements" section with bullets covering JSONL persistence with auto-compaction and events directory cleared on bootstrap. End with a horizontal rule.

### M2: Events guide and page updates

Create a new file at `/home/ben/miru/workbench3/docs/docs/developers/device-api/events.mdx`. This is a developer guide for consuming SSE events from the Device API. It should have MDX frontmatter with title "Events" and a description. The body covers: what events are (real-time deployment lifecycle notifications pushed by the agent over SSE), when to use events versus polling the deployments endpoint, how to connect (curl example using the unix socket), the SSE event frame format (`id`, `event`, `data` fields) and the JSON envelope structure (`object`, `id`, `type`, `occurred_at`, `data`), a table of event types linking to the reference pages, cursor-based replay via `?after=<id>` and `Last-Event-ID`, type filtering via `?types=`, delivery semantics (at-least-once, deduplicate by event `id`), heartbeat comments every 30 seconds, retention and compaction behavior (expired cursors return 410 Gone), and error responses (400 for malformed cursor, 410 for expired cursor).

Update the Device API overview at `/home/ben/miru/workbench3/docs/docs/developers/device-api/overview.mdx`. Change the first paragraph from "a REST API" to "a REST API with Server-Sent Events (SSE) streaming" and add a sentence about real-time deployment lifecycle events. In the second paragraph, add event streaming to the list of capabilities. Add a sentence before the "To get started" paragraph linking to the events guide.

Update the agent architecture page at `/home/ben/miru/workbench3/docs/docs/developers/agent/architecture.mdx`. Add a new `## Events` section after the "Deployments" section (after line 130). This section explains that after the sync cycle applies changes (step 4 in the numbered list), the agent emits events for successful deployment transitions, persists them locally in JSONL format, and broadcasts them to any connected SSE clients. Include a link to the events guide. Also update the sync cycle numbered list to add a step 5.5 (renumbered as step 5): "Emits events for successful deployment transitions" between the current step 4 (apply changes) and step 5 (report status, renumbered to step 6). Update the Mermaid sequence diagram to include a "Note right of A: Emit events" line between "Apply changes" and "Report status."

Add a new row to the agent versions table at `/home/ben/miru/workbench3/docs/docs/developers/agent/versions.mdx`. Insert `| v0.8.x | 2026-04-04 | <SupportedBadge /> |` above the existing `v0.7.x` row.

Update the Device API versioning page at `/home/ben/miru/workbench3/docs/docs/developers/device-api/versioning.mdx`. In the agent compatibility matrix table, add a row `| v0.8.x | v0.2.1 |` above the existing `v0.7.x` row. In the supported versions table, add a row `| <LinkNewTab href="/docs/references/device-api/v0.2.1">v0.2.1</LinkNewTab> | 2026-04-04 | <SupportedBadge /> |` above the existing `v0.2.0` row.

### M3: Navigation, redirect, and spell-check

Update `/home/ben/miru/workbench3/docs/docs.json`. In the "Device API" pages array under the "Developers" group, add `"docs/developers/device-api/events"` before `"docs/developers/device-api/versioning"`. Also update the Device API latest redirect source/destination from `v0.2.0` to `v0.2.1`.

Run `pnpm cspell` from the docs repo root. If any SSE-related terms are flagged (likely candidates: `JSONL`, `deduplicate`, `deduplication`, `SSE`), add them to the `words` array in `/home/ben/miru/workbench3/docs/cspell.json`.

## Concrete Steps

All commands are run from the docs repo root at `/home/ben/miru/workbench3/docs/` unless otherwise specified.

### M1: Changelogs

Step 1. Edit `/home/ben/miru/workbench3/docs/docs/changelogs/device-api.mdx`. Insert the following content above line 12 (the `# v0.2.0` heading):

    # v0.2.1

    *April 4, 2026*

    The `v0.2.1` release adds real-time deployment lifecycle event streaming via Server-Sent Events (SSE). On-device applications can now subscribe to a persistent event stream instead of polling for deployment status changes.

    <DeviceApiReleaseLinks version="v0.2.1" />

    ## New endpoints

    <Dropdown title="Events">
      - <GET /> `/events` — stream deployment lifecycle events via SSE
    </Dropdown>

    ## Additive changes

    <Dropdown title="Event types">
    Two event types are emitted for deployment lifecycle transitions:

    - `deployment.deployed` — emitted when a deployment's config instances are written to the device's filesystem
    - `deployment.removed` — emitted when a deployment's config instances are removed from the device's filesystem

    Each event includes a `release_id`, merged `status`, `error_status`, and timestamp in the event data payload.
    </Dropdown>
    <Separator />
    <Dropdown title="Cursor-based replay and type filtering">
    The SSE endpoint supports cursor-based replay via the `?after=<id>` query parameter or the `Last-Event-ID` HTTP header, allowing clients to resume from a known position after reconnection. Events can be filtered by type using the `?types=` query parameter.
    </Dropdown>

    ---

Expected result: `# v0.2.1` appears above `# v0.2.0`, separated by a horizontal rule.

Step 2. Edit `/home/ben/miru/workbench3/docs/docs/changelogs/agent.mdx`. Insert the following content above line 7 (the `# v0.7.0` heading):

    # v0.8.0

    *April 4, 2026*

    `v0.8.0` adds Server-Sent Events (SSE) streaming to the Miru Agent, enabling on-device applications to receive real-time deployment lifecycle notifications. This release also upgrades the Device API to [v0.2.1](/docs/changelogs/device-api#v0-2-1).

    [Device API v0.2.1 changelog >>](/docs/changelogs/device-api#v0-2-1)

    ## Features

    - Added `GET /v0.2/events` SSE endpoint for streaming deployment lifecycle events
    - Cursor-based replay via `?after=<id>` query parameter and `Last-Event-ID` header lets clients resume from a known position after reconnection
    - Type filtering via `?types=` query parameter lets clients subscribe to specific event types
    - At-least-once delivery semantics ensure no events are lost during transient disconnections
    - 30-second heartbeat keep-alive comments prevent proxies and firewalls from closing idle connections

    ## Improvements

    - Events are persisted locally in JSONL format with automatic compaction to bound disk usage
    - The events directory is cleared during agent bootstrap to prevent stale state from previous runs

    ---

Expected result: `# v0.8.0` appears above `# v0.7.0`, separated by a horizontal rule.

Step 3. Commit M1.

    cd /home/ben/miru/workbench3/docs
    git add docs/changelogs/device-api.mdx docs/changelogs/agent.mdx
    git commit -m "docs: add Device API v0.2.1 and Agent v0.8.0 changelogs

    Add changelog entries for the SSE deployment event streaming feature
    shipped in Agent v0.8.0 with Device API v0.2.1."

Expected result: clean working tree for the two changelog files.

### M2: Events guide and page updates

Step 4. Create `/home/ben/miru/workbench3/docs/docs/developers/device-api/events.mdx` with the events guide content described in Plan of Work. The file should have frontmatter (`title: "Events"`, `description` about SSE deployment events), import the `LinkNewTab` component for cross-references, and contain sections for: What are events, When to use events vs polling, Connecting, Event format, Event types (table), Cursor-based replay, Type filtering, Delivery semantics, Heartbeats, Retention & compaction, and Error responses. Include a curl example and an annotated SSE frame example.

Step 5. Edit `/home/ben/miru/workbench3/docs/docs/developers/device-api/overview.mdx`. Change line 7 from:

    The Miru Device API is a REST API for programmatically interacting with the Miru Agent from an application running on the same device as the agent.

to:

    The Miru Device API is a REST API with Server-Sent Events (SSE) streaming for programmatically interacting with the Miru Agent from an application running on the same device as the agent.

Change line 9 from:

    The Device API is not available over the internet—it is only accessible from the device running the Miru Agent. It is useful for manually refreshing configurations, retrieving current configurations, and checking agent status from applications running on your devices.

to:

    The Device API is not available over the internet—it is only accessible from the device running the Miru Agent. It is useful for manually refreshing configurations, retrieving current configurations, checking agent status, and receiving real-time deployment lifecycle events from applications running on your devices.

Insert a new paragraph before the "To get started" paragraph:

    To receive real-time notifications when deployments are applied or removed, see the [Events](/docs/developers/device-api/events) guide.

Step 6. Edit `/home/ben/miru/workbench3/docs/docs/developers/agent/architecture.mdx`. In the sync cycle numbered list, add step 5 "Emits events for successful deployment transitions" between the current steps 4 and 5, renumbering current step 5 to step 6. In the Mermaid sequence diagram, add `Note right of A: Emit events` between `Note right of A: Apply changes` and `A->>S: Report status`. Add a new `## Events` section after the "Offline resilience" paragraph (end of the Deployments section) covering event emission during the sync cycle, JSONL persistence, SSE broadcasting, and a link to the events guide.

Step 7. Edit `/home/ben/miru/workbench3/docs/docs/developers/agent/versions.mdx`. Add a new row to the versions table:

    | `v0.8.x`    | 2026-04-04   | <SupportedBadge />  |

Insert this above the `v0.7.x` row (currently line 19).

Step 8. Edit `/home/ben/miru/workbench3/docs/docs/developers/device-api/versioning.mdx`. Add `| v0.8.x | v0.2.1 |` above the `v0.7.x` row in the agent compatibility matrix. Add `| <LinkNewTab href="/docs/references/device-api/v0.2.1">v0.2.1</LinkNewTab> | 2026-04-04 | <SupportedBadge /> |` above the `v0.2.0` row in the supported versions table.

Step 9. Commit M2.

    cd /home/ben/miru/workbench3/docs
    git add docs/developers/device-api/events.mdx docs/developers/device-api/overview.mdx docs/developers/agent/architecture.mdx docs/developers/agent/versions.mdx docs/developers/device-api/versioning.mdx
    git commit -m "docs: add events guide and update pages for SSE coverage

    Add dedicated events guide for on-device SSE consumption. Update Device
    API overview, agent architecture, agent versions, and Device API
    versioning pages to reflect v0.8.0/v0.2.1 changes."

Expected result: clean working tree for all five files.

### M3: Navigation, redirect, and spell-check

Step 10. Edit `/home/ben/miru/workbench3/docs/docs.json`. In the Device API pages array, add `"docs/developers/device-api/events"` before `"docs/developers/device-api/versioning"`. Update the Device API latest redirect from `v0.2.0` to `v0.2.1` in both the source pattern comment context and the destination URL.

Step 11. Run spell-check.

    cd /home/ben/miru/workbench3/docs
    pnpm cspell "docs/**/*.mdx" --no-progress

If terms are flagged, add them to the `words` array in `/home/ben/miru/workbench3/docs/cspell.json`.

Step 12. Commit M3.

    cd /home/ben/miru/workbench3/docs
    git add docs.json cspell.json
    git commit -m "docs: update navigation and spell-check for SSE events

    Add events guide to Device API nav group, update latest redirect to
    v0.2.1, and add SSE-related terms to spell-check dictionary."

Expected result: clean working tree. All three milestones committed on `feat/sse`.

## Validation and Acceptance

After all three milestones are committed, verify the following from the docs repo root:

1. **Local dev server**: Run `pnpm dev` and confirm:
   - The events guide renders at `http://localhost:3000/docs/developers/device-api/events` with all sections, the curl example, and the SSE frame example.
   - The Device API Reference v0.2.1 dropdown appears and the Events group shows `GET /events` and the two event type reference pages.
   - The Device API v0.2.1 changelog entry renders at `http://localhost:3000/docs/changelogs/device-api#v0-2-1` with the release links component, new endpoints dropdown, and additive changes dropdowns.
   - The Agent v0.8.0 changelog entry renders at `http://localhost:3000/docs/changelogs/agent#v0-8-0` with features and improvements lists.
   - The Device API overview mentions SSE streaming and links to the events guide.
   - The agent architecture page shows the updated sync cycle list (6 steps), the updated Mermaid diagram with "Emit events", and the new Events section.
   - The agent versions table shows `v0.8.x` as the top row.
   - The Device API versioning page shows `v0.8.x` -> `v0.2.1` in the compatibility matrix and `v0.2.1` in the supported versions table.
   - The sidebar shows "Events" in the Device API developer section before "Versioning."

2. **Cross-links**: Click each cross-link to confirm it resolves:
   - Agent changelog v0.8.0 -> Device API changelog v0.2.1
   - Device API overview -> Events guide
   - Agent architecture Events section -> Events guide
   - Events guide -> `deployment.deployed` reference page
   - Events guide -> `deployment.removed` reference page

3. **Redirect**: Navigate to `/docs/references/device-api/latest/endpoints/...` and confirm it redirects to `v0.2.1` paths.

4. **Spell-check**: Run `pnpm cspell "docs/**/*.mdx" --no-progress` and confirm zero errors.

5. **Git state**: Run `git log --oneline -3` on `feat/sse` and confirm three new commits in order: changelogs, guide+pages, nav+spellcheck.

## Idempotence and Recovery

All steps in this plan are safe to repeat:

- Editing MDX files is idempotent because each edit targets specific, unique content. If a step has already been applied, the content will already match and the edit is a no-op.
- Creating the events guide file is idempotent. If the file already exists with the correct content, overwriting it produces the same result.
- Adding rows to Markdown tables is order-sensitive but each row's content is unique, so duplicate insertion is visually detectable during review.
- `git add` and `git commit` are safe to repeat. If there are no changes to commit, git will report "nothing to commit" and exit cleanly.
- The spell-check step (`pnpm cspell`) is read-only. Only the subsequent edit to `cspell.json` is a write, and adding an already-present word is harmless.
- The `docs.json` navigation edit inserts a single entry. If the entry is already present, the edit should be skipped to avoid duplicates. Verify the entry is absent before inserting.

Rollback for each milestone: Because each milestone ends with a commit, reverting a milestone is a single `git revert <sha>`. No database migrations, infrastructure changes, or external state mutations are involved. The entire plan is confined to file edits within the docs repo on the `feat/sse` branch.

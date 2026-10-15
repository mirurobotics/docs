# Add a network requirements page for the Miru Agent

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | New page `docs/developers/agent/network.mdx`, nav entry in `docs/docs.json`, cross-link in `docs/developers/agent/security.mdx`. |
| `agent/` | read-only | Source of truth for hosts, ports, and protocols (`agent/src/network/mod.rs`, `agent/src/main.rs`, `Cargo.toml`). |

This plan lives in `docs/plans/` because all edits are documentation pages in the docs repo.

## Purpose / Big Picture

Customers configuring firewalls cannot find which endpoints the Miru Agent needs; the facts are buried in the "Transport security" and "Network posture" sections of `docs/developers/agent/security.mdx`. After this change, the sidebar's Developers > Agent group shows a "Network requirements" page (URL `/developers/agent/network`) whose first table is the complete runtime allowlist, followed by install-time hosts, proxy/TLS-inspection caveats, and the legacy MQTT port.

## Progress

- [x] (2026-10-13) Create `docs/developers/agent/network.mdx`.
- [x] (2026-10-13) Add the page to `docs/docs.json` after `developers/agent/security`.
- [x] (2026-10-13) Add a cross-link in `docs/developers/agent/security.mdx`.
- [x] (2026-10-13) Static validation and commit.

## Surprises & Discoveries

- Observation: `ci.yml` runs only on `pull_request` and pushes to main/staging/uat/production, and `gh` is not authenticated in this environment, so pushing the feature branch triggers no CI.
  Evidence: `.github/workflows/ci.yml` `on:` block.

## Decision Log

- Decision: Note that GitHub may redirect release downloads to other GitHub-owned hosts, without naming them.
  Rationale: Allowing only `github.com` can break downloads; specific redirect hostnames were not verified.
  Date/Author: 2026-10-13, Claude.
- Decision: Remove the `github.com` row from "Installation and upgrades" (supersedes the decision above).
  Rationale: Review feedback on PR #248: manual GitHub downloads aren't a supported automated workflow, so the page lists only the APT host.
  Date/Author: 2026-10-13, Claude.

## Outcomes & Retrospective

Added `docs/developers/agent/network.mdx`, its nav entry, and a cross-link from `security.mdx`. Static checks passed (JSON valid, link targets and anchors exist, sentence-case headings, no `--`, no bucket endpoints). CI validation did not run; it will run when a pull request is opened.

## Context and Orientation

The docs repo is a Mintlify site. Pages are `.mdx` files under `docs/` and the URL is the path without the `docs/` prefix or extension (for example `docs/developers/agent/security.mdx` is served at `/developers/agent/security`). Navigation is the JSON file `docs/docs.json`; the "Agent" group lists `developers/agent/overview`, `install`, `versions`, `commands`, `security`, `architecture`, then a "File system access" subgroup. Built-in Mintlify components (`<Info>`, `<Warning>`, `<Tabs>`, `<Tab>`) need no import.

Custom lint rules in `tools/lint/linter/` run in CI: `headingcase` requires sentence-case titles and headings (allowlisted tokens include `MQTT`, `TLS`, `HTTPS`, `API`, `GitHub`, `Agent`, `Unix`, and version tags like `v0.10.3`); `nodoubledash` forbids `--` in prose; CSpell checks spelling against `cspell.json`.

Verified agent facts (from `agent/src/network/mod.rs`, `agent/src/main.rs`, `Cargo.toml`, and `docs/changelog/agent.mdx`):

- All connections are outbound; no inbound rules are needed.
- `api.mirurobotics.com`, TCP 443, HTTPS with certificate validation: provisioning, token refresh, sync cycles.
- `mqtt.mirurobotics.com`, TCP 443, MQTT over TLS (not HTTPS, not WebSockets): real-time deployment notifications. The MQTT client (`rumqttc`, built without its proxy feature) cannot use an HTTP proxy, and TLS-inspecting or HTTP-only filters on 443 may drop it. Without MQTT the agent still syncs over HTTPS via the Device API, polling, and restart triggers described in `docs/developers/agent/architecture.mdx`, but deployments are no longer pushed in real time.
- Agents before `v0.10.3` (and before `v0.7.2` on the v0.7 line) use port 8883 for MQTT.
- `is_allowed_host` only accepts `mirurobotics.com` or its subdomains, so `*.mirurobotics.com` on 443 is a valid simpler rule.
- Install/upgrade only: APT uses `packages.mirurobotics.com` on 443; manual `.deb` and Windows `.msi` downloads come from `github.com/mirurobotics/agent/releases`. Packages can be downloaded elsewhere and transferred (`/developers/agent/install#poor-connectivity`).
- On-device Device API (Unix socket on Linux, loopback `127.0.0.1:6478` on Windows) needs no firewall rule; `security.mdx` already documents it under "Network posture".

Out of scope: Data Uploads bucket endpoints (S3/GCS) must not be mentioned; do not claim HTTPS proxy support for the API connection.

## Plan of Work

1. Create `docs/developers/agent/network.mdx` with frontmatter `title: 'Network requirements'`. Open with one sentence saying all traffic is outbound and no inbound ports are needed. Section "Required endpoints": a table (Host, Port, Protocol, Purpose) with the two runtime hosts, then an `<Info>` noting the `*.mirurobotics.com` wildcard option. Section "MQTT connection": raw MQTT over TLS, no HTTP proxy support, TLS inspection caveat in a `<Warning>`, and the consequence of blocking it with a link to `/developers/agent/architecture#triggers`. Section "Older agent versions": port 8883 before `v0.10.3`/`v0.7.2`, links to `/changelog/agent#v0-10-3` and `/changelog/agent#v0-7-2`. Section "Installation and upgrades": table with `packages.mirurobotics.com` (APT) and `github.com` (manual downloads), plus a link to `/developers/agent/install#poor-connectivity`. Section "On-device traffic": one sentence linking to `/developers/agent/security#network-posture`.
2. In `docs/docs.json`, insert `"developers/agent/network",` immediately after `"developers/agent/security",` in the Agent group.
3. In `docs/developers/agent/security.mdx`, under "Network posture" > "External", append one sentence linking to `/developers/agent/network`.

## Concrete Steps

From `docs/` (repo root): create and edit the three files above with an editor. Then run read-only checks:

    python3 -c "import json; json.load(open('docs/docs.json'))"
    grep -n "developers/agent/network" docs/docs.json docs/developers/agent/security.mdx

Expected: the JSON loads without error; grep shows the nav entry and the cross-link. Commit:

    git add docs/developers/agent/network.mdx docs/docs.json docs/developers/agent/security.mdx
    git commit -m "docs(agent): add network requirements page"

## Validation and Acceptance

This is a docs-only change with no executable tests. Test steps are static checks: `docs/docs.json` parses as JSON; every internal link on the new page resolves to an existing `.mdx` file and heading anchor (`/developers/agent/architecture#triggers`, `/changelog/agent#v0-10-3`, `/changelog/agent#v0-7-2`, `/developers/agent/install#poor-connectivity`, `/developers/agent/security#network-posture`); headings are sentence-case; no `--` in prose; no mention of S3, GCS, or buckets.

In CI, `ci.yml` job `lint` (custom MDX lint, ESLint, CSpell, `mint validate`) must pass on the pushed head; validation is complete when `$preflight` reports `CLEAN` (CI green on the pushed branch head). Acceptance: `/developers/agent/network` renders, appears in the Agent sidebar after Security, and its first table lists `api.mirurobotics.com` and `mqtt.mirurobotics.com` on TCP 443.

## Idempotence and Recovery

All steps are plain file edits and can be repeated. To roll back, `git revert` the commit.

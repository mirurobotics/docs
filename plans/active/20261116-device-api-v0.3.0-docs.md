# Draft the Device API v0.3.0 and agent v0.13.0 release docs

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/user/docs`) | read-write | Vendor the `v0.3.0` reference, add its dropdown and `latest` redirect, add changelog entries and version-table rows. |
| `openapi/` (`/home/user/openapi`) | read-only | Source of the Device API spec (`apis/apps/device-server/openapi.gen.yaml` on `main`) and its release renderer (`tools/release/spec.py`). No `device/v0.3.0` tag exists yet. |
| `agent/` (`/home/user/mirurobotics/agent`) | read-only | Confirms what the next agent release serves. Do not modify. |

This plan lives in `docs/plans/` because every write is in the docs repo. Work on the existing branch `claude/project-thread-9m8g99`; the PR is a draft against `main` and is PR 6 of the docs plan.

## Purpose / Big Picture

Device API `v0.3.0` removes calibration roles and aspects, keeps one latest calibration per component, adds `invalidated_at`, and adds `POST /calibrations/latest/invalidate`. It is served at `/v0.3`. The next agent release (placeholder `v0.13.0`) serves it next to `v0.2.5` at `/v0.2`. After this change a reader can pick `v0.3.0` in the Device API Reference dropdown (the default), see the new **Invalidate Latest** page under Calibrations, read `v0.3.0` and `v0.13.0` entries in the changelogs, and find both versions in the version tables, including which agent serves which API versions and when to move from `/v0.2` to `/v0.3`.

Nothing is released yet. All versions and the date are placeholders, listed in the PR's "Before merging" checklist, and the PR stays a draft until the agent ships.

## Decisions

### Names

| New name | Alternatives | Why this one |
|----------|--------------|--------------|
| Folder `docs/references/device-api/v0.3.0/` and dropdown `v0.3.0` | `v0.3/`; flat `v0.3.0.yaml` (as `v0.1.0.yaml`, `v0.2.0.yaml`) | Matches `v0.2.1`–`v0.2.5`, which all use a folder with `api.yaml` and `events/`. |
| Agent version `v0.13.0` (placeholder), row `v0.13.x` | `v0.12.5`; `v1.0.0` | Serving a new API alias is a feature-level change; `v0.12.0` was the last such minor. |
| Python SDK `v0.11.0` / `v0.11.x` (placeholder) | `v0.10.1`; `v1.0.0` | Each Device API version so far bumped the SDK minor (`v0.2.5` → `v0.10.x`). |
| Changelog dropdown titles "Create calibration request", "Latest calibrations", "Calibration response" | "Calibrations", "List latest response" | "Create calibration request" is the `v0.2.4` title for the same request; the others name the route or object they describe. |

### Precedent

| New file | Mirrors |
|----------|---------|
| `docs/references/device-api/v0.3.0/api.yaml` | `docs/references/device-api/v0.2.5/api.yaml` (docs#315): rendered spec, previous samples kept, `kuka-arm` examples. |
| `docs/references/device-api/v0.3.0/events/deployment-deployed.mdx`, `deployment-removed.mdx` | Byte-for-byte copies of `docs/references/device-api/v0.2.5/events/`; events did not change. |

Edits mirror docs#315: the `v0.3.0` dropdown is a copy of `v0.2.5`'s, the changelog entries follow `v0.2.5` and, for breaking changes, `v0.2.0` in `docs/changelog/device-api.mdx`, and the agent entry follows `v0.12.4` in `docs/changelog/agent.mdx`.

### Left out

- `docs/cfg-mgmt/report-calibrations.mdx` v0.3 and invalidate sections, and invalidation fields on `docs/cfg-mgmt/concepts/calibrations.mdx`: PR 1 (flat components, 48-character slugs) edits the same pages. Add them on this branch after PR 1 merges.
- Python samples from `api/pull-stainless.sh`: the SDK for `v0.3.0` isn't released. Run it after the SDK ships.
- `/v0.2` examples in `docs/developers/device-api/authn.mdx` and `events.mdx`: they keep working on agent `v0.13.0`. Moving them to `/v0.3` is a follow-up once the agent ships.
- A "Python SDK" section like `v0.2.0`'s: there is no repository or import-path change, so the SDK upgrade is one sentence in the migration steps.
- The reworded file rule `upload.path` description: per docs#315 review, the Device API changelog lists only route and field changes.

### Open questions

- Should `v0.3.0` become the default dropdown and the `latest` redirect target now? Default: yes, as docs#315 did. Merging early would make `v0.3.0` the default before any agent serves it, and would point `DeviceApiLink` links (used on `report-calibrations.mdx`) at a reference without roles and aspects while that page still documents them. The PR stays draft until the agent releases and the follow-up pages land.
- Which date placeholder? Default: **November 30, 2026** (`2026-11-30`) in both changelogs, `developers/device-api/versions.mdx` and `developers/agent/versions.mdx`, flagged in "Before merging".
- What goes in `info.version`? Default: `v0.3`, with `x-release-version: v0.3.0`, which is what openapi's release renderer writes and what every vendored version has (`version: v0.2`, `x-release-version: v0.2.5`), rather than `0.3.0`.
- Is the agent release a new minor? Default: yes, add a `v0.13.x` row to `docs/developers/agent/versions.mdx`. If it ships as a `v0.12.x` patch, drop the row and fix the version everywhere.

## Progress

- [x] Milestone 1: vendor the `v0.3.0` reference, add the dropdown and redirect, commit.
- [x] Milestone 2: add the Device API `v0.3.0` and agent `v0.13.0` changelog entries, commit.
- [x] Milestone 3: add the version-table rows, commit.
- [ ] Milestone 4: push, open the draft PR, `$preflight` reports `CLEAN`.
- [x] Source written

## Surprises & Discoveries

- The local openapi `main` (`fe03a87`) was one commit behind `origin/main` (`dcf7e1b`, a CI-only change). The device spec and `tools/release/` are identical in both, so the render used the local checkout and `x-git-commit` names `fe03a87`. Re-vendoring from the `device/v0.3.0` release asset replaces it.

## Decision Log

- 2026-11-17, review (ben-miru): the agent Features list drops the `GET /v0.3/calibrations/latest` bullet, since that endpoint already existed in `v0.2.5`; `invalidated_at` is covered in the Device API changelog.
- 2026-11-17: openapi tagged `device/v0.3.0-beta.3` at `70b258d`, and `api.yaml` is re-rendered from it. It brings #511 (sync moves to `POST /sync` under Agent, sync and `Device` times become `*_secs_ago` durations) and #515 (invalidate moves to `POST /calibrations/latest/invalidate`, keyed by `component` in the body, with an optional `report_id` check). The nav moves both pages, and the Device API changelog gains **Sync** and **Device response** breaking changes. #513 changed only Agent API wording, so the `kuka-arm` examples stay. `x-release-version` stays `v0.3.0` until GA.
- 2026-11-17, review (ben-miru): the agent compatibility matrix is keyed by agent version, with one column per API release line (`v0.3.x`, `v0.2.x`, `v0.1.x`), so an agent serving two lines reads as one row. The separate sentence about serving more than one version is dropped as the columns show it.
- 2026-11-16: per the flat components PR's review (docs#337), the v0.2 to v0.3 migration lives only in the Device API docs. The create request dropdown now states the backend's join rule (`JoinSlug`: `-`, unset parts skipped, same history as reporting the joined slug) and advises staying within 48 characters, since the joined length isn't validated.
- 2026-11-16: openapi tagged `device/v0.3.0-beta.1` at `dcf7e1b`. Its `apis/` and `tools/` match `fe03a87`, which `api.yaml` was rendered from, so only `x-git-commit` changes to the tagged commit. `x-release-version` stays `v0.3.0`, the final version, until GA. The tag has no GitHub release asset yet.
- 2026-11-16, design review: the breaking-change diffs now move `role` and `aspect` into `component`, joined with `-` (`kuka-arm-front-left-pose`), as the backend agent spec joins them, so following the diff keeps reporting to the same component. The latest calibration diff shows the joined slug too, since the v0.3 `Calibration.component` is the component's slug, not the reported one.
- 2026-11-16, design review: "agent `v0.13.0` also serves `v0.2.5` at `/v0.2`" is kept only in the first migration step. Dropped from the Device API intro and the agent Features list, and the `versions.mdx` paragraph is cut to one version-free sentence, so a release edits one place.
- 2026-11-16, design review: that `/v0.2` behaves unchanged on the new agent is not yet backed by code (agent `main` has no `/v0.3` module). Checking it against the agent PR that adds `/v0.3` is a "Before merging" item.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

The docs repo is a Mintlify site; content lives under `docs/`. MDX is Markdown plus JSX components. Key files:

- `docs/references/device-api/<version>/api.yaml`: one vendored OpenAPI spec per Device API version. Mintlify builds endpoint pages from it. Operations carry `x-codeSamples` (curl for Linux over the Unix socket, PowerShell for Windows, Python from the SDK). `<version>/events/*.mdx` are event-type pages.
- `docs/docs.json`: site config. Under the product `"Device API Reference"`, `dropdowns` lists one dropdown per version, newest first; the first is the default. Each group points `openapi.source` at that version's `api.yaml` and lists operations as `"METHOD /path"`. The `redirects` array maps `/references/device-api/latest/:slug*` to the current version.
- `docs/changelog/device-api.mdx`, `docs/changelog/agent.mdx`: newest-first changelogs. Each entry is `# vX.Y.Z`, an italic date, a summary, sections, and `---`.
- `docs/developers/device-api/versions.mdx`: supported-versions table, agent compatibility matrix and a rule-of-thumb paragraph. `docs/developers/device-api/sdks.mdx`: Python SDK matrix. `docs/developers/agent/versions.mdx`: agent minor lines.
- `api/add_unix_socket_curl.py` and `api/add_powershell_samples.py`: add or replace the curl and PowerShell sample on every operation, using `servers[0].url` for the version prefix, and dump YAML with a representer that writes multi-line strings and strings with `\` as `|` blocks.

The source spec is openapi `main`, `apis/apps/device-server/openapi.gen.yaml` (`servers` `http://localhost/v0.3`, `info.version: 0.0.0`, and `$API_VERSION$`, `$RELEASE_VERSION$`, `$API_GIT_COMMIT$` placeholders that openapi's `tools/release/spec.py` `render()` fills). `git -C /home/user/openapi diff device/v0.2.5 HEAD -- apis/apps/device-server/openapi.gen.yaml` shows every change from `v0.2.5`:

- Breaking: `role` and `aspect` removed from `CreateCalibrationRequest` and `Calibration`; the `role` and `aspect` query filters removed from `GET /calibrations/latest`, which now returns one calibration per component.
- Additive: `Calibration.invalidated_at` (required, nullable date-time).
- New: `POST /calibrations/latest/invalidate` (`invalidateLatestCalibration`, summary "Invalidate Latest"), body `InvalidateCalibrationRequest { component, report_id?, reason }` (reason 1–512 characters), returns the `Calibration`. A component with no calibration, or a `report_id` that isn't the component's latest, returns 404; an already invalid one is returned unchanged.
- Breaking: `POST /device/sync` (tag Device) moves to `POST /sync` (tag Agent, `client.agent.sync()`), and its response and `Device` report `*_secs_ago` durations instead of timestamps.
- Description only: file rule `upload.path`. Not in the changelog.

Agent facts: the latest release is `v0.12.4` (2026-11-10), serving only `/v0.2` (`v0.2.5`). Agent `main` vendors the `v0.3` spec but does not serve it yet. The release that serves both `/v0.2` and `/v0.3` is the placeholder `v0.13.0`; its `/v0.2` routes are unchanged.

## Plan of Work

Milestone 1, reference. Render openapi `main`'s spec the way a `device/v0.3.0` release would, into `docs/references/device-api/v0.3.0/api.yaml`. Run both sample scripts on it, which gives the new invalidate route curl and PowerShell samples with `/v0.3` URLs. Then carry over `v0.2.5`'s `x-codeSamples` for every operation that exists in both, with `/v0.2/` changed to `/v0.3/` and, in `createCalibration`, the `role` and `aspect` lines removed. Add the request body to the invalidate samples in the shape of the `POST /calibrations` samples. Change calibration `component` examples from `arm` to `kuka-arm` in the same five places `v0.2.5` does. Copy `events/` from `v0.2.5`. In `docs/docs.json`, insert a `v0.3.0` dropdown before `v0.2.5`'s, as a copy of it with every `v0.2.5` path changed to `v0.3.0` and `"POST /calibrations/latest/invalidate"` added to the Calibrations group after the content download. Change the `latest` redirect destination to `/references/device-api/v0.3.0/:slug*`.

Milestone 2, changelogs. Insert the `# v0.3.0` entry at the top of `docs/changelog/device-api.mdx` and the `# v0.13.0` entry at the top of `docs/changelog/agent.mdx`, with the text in Concrete Steps.

Milestone 3, tables. In `docs/developers/device-api/versions.mdx`, add the `v0.3.0` supported row, replace the `v0.2.5` matrix row with two rows, and extend the rule-of-thumb paragraph. In `docs/developers/device-api/sdks.mdx`, add the `v0.3.0` row and pin `v0.2.5` to `v0.10.x`. In `docs/developers/agent/versions.mdx`, add the `v0.13.x` row.

## Concrete Steps

All commands run from the docs repo root, `/home/user/docs`, unless stated. `$SCRATCH` is any directory outside the repo (the session scratchpad). The two Python files below are one-off generators; they are not committed.

Milestone 1. Save this as `$SCRATCH/render_v030.py`:

    """One-off: render openapi main's device spec as a v0.3.0 release would."""
    import os, sys
    from pathlib import Path

    openapi_root, out_path = Path(sys.argv[1]), Path(sys.argv[2]).resolve()
    sys.path.insert(0, str(openapi_root))
    os.chdir(openapi_root)
    from tools.release import git, spec  # noqa: E402

    commit = git.commit()
    text = spec.render(
        template=(openapi_root / "apis/apps/device-server/openapi.gen.yaml").read_text(),
        api_version="v0.3",
        release_version="v0.3.0",
        commit=commit.sha,
        git_commit=commit.to_dict(),
        build_info=spec.build().to_dict(),
    )
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(text, encoding="utf-8")

Save this as `$SCRATCH/finish_v030.py`:

    """One-off: carry v0.2.5's code samples and example edits into the v0.3.0 spec."""
    import copy, sys
    import yaml

    prev_path, spec_path = sys.argv[1], sys.argv[2]
    prev = yaml.safe_load(open(prev_path, encoding="utf-8"))
    spec = yaml.safe_load(open(spec_path, encoding="utf-8"))
    METHODS = {"get", "post", "put", "patch", "delete"}

    # Existing operations keep v0.2.5's samples, moved to the /v0.3 base path.
    for path, item in spec["paths"].items():
        for method, op in item.items():
            old = prev["paths"].get(path, {}).get(method) if method in METHODS else None
            if old is None:
                continue
            samples = copy.deepcopy(old["x-codeSamples"])
            for s in samples:
                s["source"] = s["source"].replace("/v0.2/", "/v0.3/")
                if op["operationId"] == "createCalibration":
                    s["source"] = "\n".join(
                        l for l in s["source"].split("\n")
                        if not l.strip().startswith(('"role"', '"aspect"', "role ", "aspect "))
                    )
            op["x-codeSamples"] = samples

    # The invalidate route takes a body, like POST /calibrations.
    op = spec["paths"]["/calibrations/latest/invalidate"]["post"]
    for s in op["x-codeSamples"]:
        if s["lang"] == "curl":
            s["source"] += " \\\n  --header 'Content-Type: application/json' \\\n  --data '{\"component\": \"kuka-arm\", \"report_id\": \"0b6f3c2e-6a51-4e0f-9a43-2f1d7c9e8b10\", \"reason\": \"power cycle\"}'"
        elif s["lang"] == "powershell":
            first, *rest = s["source"].split("\n")
            rest[-1] += " `"
            s["source"] = "\n".join(
                [first, '$body = @{ component = "kuka-arm"; report_id = "0b6f3c2e-6a51-4e0f-9a43-2f1d7c9e8b10"; reason = "power cycle" } | ConvertTo-Json; `', *rest,
                 '  -ContentType "application/json" `', "  -Body $body"]
            )

    # Calibration component examples use kuka-arm, as v0.2.3 to v0.2.5 do.
    c = spec["components"]
    c["parameters"]["calibration_search_component"]["schema"]["example"] = ["kuka-arm"]
    req = c["schemas"]["CreateCalibrationRequest"]["properties"]["component"]
    req["example"] = "kuka-arm"
    req["description"] = req["description"].replace("`arm`", "`kuka-arm`")
    cal = c["schemas"]["Calibration"]
    cal["properties"]["component"]["example"] = "kuka-arm"
    cal["example"]["component"] = "kuka-arm"
    c["schemas"]["LatestCalibrationsResponse"]["example"]["calibrations"][0]["component"] = "kuka-arm"

    # Same representer as api/add_unix_socket_curl.py.
    def str_presenter(dumper, data):
        if "\n" in data or "\\" in data:
            return dumper.represent_scalar("tag:yaml.org,2002:str", data, style="|")
        return dumper.represent_scalar("tag:yaml.org,2002:str", data)

    yaml.add_representer(str, str_presenter)
    with open(spec_path, "w", encoding="utf-8") as f:
        yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True)

Then run, from `/home/user/docs`:

    git -C /home/user/openapi status --porcelain        # expect no output (dirty: false)
    python3 $SCRATCH/render_v030.py /home/user/openapi docs/references/device-api/v0.3.0/api.yaml
    python3 api/add_unix_socket_curl.py docs/references/device-api/v0.3.0/api.yaml
    python3 api/add_powershell_samples.py docs/references/device-api/v0.3.0/api.yaml
    python3 $SCRATCH/finish_v030.py docs/references/device-api/v0.2.5/api.yaml docs/references/device-api/v0.3.0/api.yaml
    cp -r docs/references/device-api/v0.2.5/events docs/references/device-api/v0.3.0/events

Check the result (read-only inspection):

    grep -c '/v0\.2' docs/references/device-api/v0.3.0/api.yaml       # 0
    grep -nE 'role|aspect|\$[A-Z_]+\$|example: arm$' docs/references/device-api/v0.3.0/api.yaml   # no output
    git diff --no-index --stat docs/references/device-api/v0.2.5/api.yaml docs/references/device-api/v0.3.0/api.yaml
      ... 1 file changed, about 121 insertions(+), 128 deletions(-)

The head of the file reads `version: v0.3`, `x-release-version: v0.3.0`, `servers: - url: http://localhost/v0.3`. The invalidate operation's samples read:

    curl \
      --unix-socket /run/miru/miru.sock \
      --request POST \
      --url http://localhost/v0.3/calibrations/latest/invalidate \
      --header 'Content-Type: application/json' \
      --data '{"component": "kuka-arm", "report_id": "0b6f3c2e-6a51-4e0f-9a43-2f1d7c9e8b10", "reason": "power cycle"}'

    $api = Get-Content -Raw "$env:ProgramData\Miru\device-api\device-api.json" | ConvertFrom-Json; `
    $body = @{ component = "kuka-arm"; report_id = "0b6f3c2e-6a51-4e0f-9a43-2f1d7c9e8b10"; reason = "power cycle" } | ConvertTo-Json; `
    Invoke-RestMethod -Method Post `
      -Uri "http://127.0.0.1:$($api.port)/v0.3/calibrations/latest/invalidate" `
      -Headers @{ Authorization = "Bearer $($api.token)" } `
      -ContentType "application/json" `
      -Body $body

Edit `docs/docs.json`: under `"product": "Device API Reference"`, insert a copy of the whole `"dropdown": "v0.2.5"` object above it, with `"dropdown": "v0.3.0"` and every `references/device-api/v0.2.5/` changed to `references/device-api/v0.3.0/` inside the copy only. Its Calibrations group's pages become:

    "POST /calibrations",
    "GET /calibrations/latest",
    "GET /calibrations/latest/{report_id}/content",
    "POST /calibrations/latest/invalidate"

In `redirects`, change the destination of `/references/device-api/latest/:slug*` to `/references/device-api/v0.3.0/:slug*`.

Commit milestone 1 (the orchestrator commits; signed commits are configured):

    git add docs/references/device-api/v0.3.0 docs/docs.json
    git commit -m "docs(device-api): add the v0.3.0 reference"

Milestone 2. Insert at the top of `docs/changelog/device-api.mdx`, after the imports and before `# v0.2.5`:

    # v0.3.0

    *November 30, 2026*

    The `v0.3.0` release keeps one latest [calibration](/cfg-mgmt/concepts/calibrations) for each component and lets on-device applications mark it invalid. It removes roles and aspects from calibrations and moves the Device API to `/v0.3`. Agents that serve `v0.3.0` also serve `v0.2.5` at `/v0.2`, so existing clients keep working until they migrate.

    <DeviceApiReleaseLinks version="v0.3.0" />

    ## New endpoints

    <Dropdown title="Calibrations">
      - <POST /> `/calibrations/latest/invalidate` — mark a latest calibration invalid
    </Dropdown>

    ## Breaking changes

    <Dropdown title="Base path change">
    The URL version prefix changed from `/v0.2` to `/v0.3`. Update all request URLs in your client application.

    ```diff nocopy
    - http://localhost/v0.2/{endpoint}
    + http://localhost/v0.3/{endpoint}
    ```
    </Dropdown>
    <Separator />
    <Dropdown title="Create calibration request">
    The `POST /calibrations` request removed `role` and `aspect`.

    ```diff nocopy
     {
       "component": "kuka-arm",
    -  "role": "front-left",
    -  "aspect": "pose",
       "serial_number": "SN-1234",
       "content": {
         "format": "json",
         "data": "{\"x\": 0.12, \"y\": -0.03}"
       }
     }
    ```
    </Dropdown>
    <Separator />
    <Dropdown title="Latest calibrations">
    `GET /calibrations/latest` returns one calibration for each component, instead of one for each component, role and aspect. It removed the `role` and `aspect` filters, and its calibrations removed `role` and `aspect`.

    ```diff nocopy
     {
       "object": "calibration",
       "report_id": "0b6f3c2e-6a51-4e0f-9a43-2f1d7c9e8b10",
       "component": "kuka-arm",
    -  "role": "front-left",
    -  "aspect": null,
       "serial_number": "SN-1234",
       ...
     }
    ```
    </Dropdown>

    ## Additive changes

    <Dropdown title="Calibration response">
    Calibrations added `invalidated_at`, when an on-device application marked the calibration invalid, or `null` if it's valid.

    ```diff nocopy
     {
       "object": "calibration",
       ...
       "reported_at": "2026-03-10T12:00:00Z",
    +  "invalidated_at": null
     }
    ```
    </Dropdown>

    ## Migration steps

    <Steps>
      <Step title="Verify agent compatibility">
        Ensure all devices are running Miru Agent `v0.13.0` or later. See [compatibility matrix](/developers/device-api/versions#agent-compatibility-matrix). Agent `v0.13.0` also serves `v0.2.5` at `/v0.2`, so clients on `/v0.2` keep working while you migrate.
      </Step>

      <Step title="Migrate breaking changes">
        Migrate your client application according to the breaking changes listed above. If you are using the [Python SDK](/developers/device-api/sdks#python), upgrade to version `v0.11.0` or later.
      </Step>

      <Step title="Optionally integrate the new endpoint">
        Integrate the invalidate endpoint to mark a latest calibration invalid from the device.
      </Step>
    </Steps>

    ---

The file already imports `Dropdown`, `DeviceApiReleaseLinks`, `GET`, `POST` and `Separator`; no import changes.

Insert at the top of `docs/changelog/agent.mdx`, after the front matter and before `# v0.12.4`:

    # v0.13.0

    *November 30, 2026*

    Miru Agent `v0.13.0` adds Device API `v0.3.0` at `/v0.3`, next to `v0.2.5` at `/v0.2`. `v0.3.0` keeps the latest [calibration](/cfg-mgmt/concepts/calibrations) for each component and lets on-device applications mark it invalid.

    [Device API v0.3.0 changelog »](/changelog/device-api#v0-3-0)

    ## Features

    - The Device API is served at `/v0.3` (`v0.3.0`) as well as `/v0.2` (`v0.2.5`); `/v0.2` routes are unchanged, so existing clients keep working
    - Added `POST /v0.3/calibrations/latest/invalidate`, which marks a latest calibration invalid with a reason and uploads the invalidation to Miru

    ---

Commit milestone 2:

    git add docs/changelog/device-api.mdx docs/changelog/agent.mdx
    git commit -m "docs(changelog): add Device API v0.3.0 and agent v0.13.0"

Milestone 3. In `docs/developers/device-api/versions.mdx`:

- Supported versions table, new first row: `| <LinkNewTab href="/references/device-api/v0.3.0">v0.3.0</LinkNewTab>    | 2026-11-30   | —            | <SupportedBadge /> |`.
- Agent compatibility matrix: replace `` | `v0.2.5`       | `v0.12.4+`         | `` with two rows, `` | `v0.2.5`, `v0.3.0` | `v0.13.0+` `` then `` | `v0.2.5` | `v0.12.4` ``, padded to the table's columns.
- After the paragraph starting "As rule of thumb", add: "An agent can serve more than one API version, each at its own breaking change alias. For example, agent `v0.13.0` serves `v0.2.5` at `/v0.2` and `v0.3.0` at `/v0.3`. Keep your client on `/v0.2` until every device in your fleet runs agent `v0.13.0` or later, then move it to `/v0.3` to use `v0.3.0`'s features."

In `docs/developers/device-api/sdks.mdx`, Python table: add first row `` | `v0.3.0`    | `v0.11.x` and above | `` and change the `v0.2.5` row to `` | `v0.2.5`    | `v0.10.x`    | ``.

In `docs/developers/agent/versions.mdx`, supported versions table: add first row `` | `v0.13.x`   | 2026-11-30   | —            | <SupportedBadge />  | ``.

Commit milestone 3:

    git add docs/developers
    git commit -m "docs(versions): list Device API v0.3.0 and agent v0.13.x"

Milestone 4. Push and open the draft PR with `$pr`. The body follows docs#315: a "Do not merge yet" note (merging would make `v0.3.0` the default and `latest` before any agent serves it), Summary, Changes, Precedent (docs#315), Decisions from this plan, and a "Before merging" checklist:

- Tag `device/v0.3.0` in openapi, re-vendor `api.yaml` from the release asset (`device.yaml`) with the same steps, and re-check the diff.
- Release the agent and fix `v0.13.0` / `v0.13.x` in both changelogs and both `versions.mdx` files if it differs; reconcile the agent entry's features with its release notes.
- Release the Python SDK and fix `v0.11.0` / `v0.11.x` in `changelog/device-api.mdx` and `sdks.mdx` if it differs.
- Replace the placeholder date **November 30, 2026** (`2026-11-30`) in both changelogs and both `versions.mdx` files.
- Upload the spec to `assets.mirurobotics.com/docs/openapi/device/v0.3.0.yaml` so the changelog's download link works.
- After PR 1 merges, add the v0.3 and invalidate sections to `cfg-mgmt/report-calibrations.mdx` and the invalidation fields to `cfg-mgmt/concepts/calibrations.mdx` on this branch.
- Optionally, after the SDK ships, run `./api/pull-stainless.sh` for Python samples.

Then run `$preflight`.

## Validation and Acceptance

Validation runs only in CI (`.github/workflows/ci.yml`, triggered on `pull_request`). Do not run lint, build or tests locally. The `lint` job runs `pnpm run test:lint` (lint fixtures under `tests/lint-fixtures/`, unchanged), `./scripts/lint.sh` (custom MDX prose rules including the `redirects` rule, which checks that the `latest` destination's prefix `references/device-api/v0.3.0` is a directory; ESLint MDX; CSpell; and `mint openapi-check` on every `docs/references/**/*.yaml`, including the new `api.yaml`), and `pnpm run validate` (Mintlify build, which resolves every `"METHOD /path"` in the new dropdown against the spec). The `audit` job runs `pnpm audit`. `lint-custom-linter` and `test-custom-linter` skip their steps because `tools/lint/` is unchanged. No test fixture or lint rule enumerates Device API versions, so none needs updating.

Acceptance, on the PR's Mintlify preview: `/references/device-api/latest/endpoints/calibrations/create` lands on the `v0.3.0` Create page, whose request has no `role` or `aspect`; the Device API Reference opens on the `v0.3.0` dropdown, whose Calibrations group lists **Invalidate Latest** with curl and PowerShell samples using `/v0.3` and a `reason` body; `/changelog/device-api#v0-3-0` and `/changelog/agent#v0-13-0` render; the agent matrix shows `v0.2.5`, `v0.3.0` → `v0.13.0+`.

The PR leaves draft only after `$preflight` reports `CLEAN` (CI green on the pushed head) and the "Before merging" items are done.

## Idempotence and Recovery

The generation steps are repeatable: `render_v030.py` overwrites `api.yaml`, the two sample scripts replace rather than duplicate samples, and `finish_v030.py` overwrites every carried-over sample. Re-running `finish_v030.py` twice appends the invalidate body twice, so always re-run the whole sequence from `render_v030.py`. Re-vendoring after the `device/v0.3.0` tag uses the same steps with the release asset in place of the render. `cp -r` onto an existing `events/` creates a nested folder; remove `docs/references/device-api/v0.3.0/events` first. Docs edits are plain text; `git checkout -- <file>` undoes one.

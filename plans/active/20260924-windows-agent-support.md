# Document Windows 10/11 x64 support for the Miru Agent

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Add Windows coverage to agent install, provisioning, operations, config file path, file rule, device, Device API, and CLI install pages under `docs/docs/`. |
| `agent/` | read-only | Source of truth for the MSI, Windows service, data layout, logs, and CLI. |
| `backend/` | read-only | Source of truth for Windows path and glob validation and error codes. |
| `core/` | read-only | Shared path validation (`pkg/lexpath`, `pkg/ugc/files.go`) used by the backend. |
| `openapi/` | read-only | Source of truth for device `os`/`arch`/system-info fields and provisioning-token semantics. |
| `frontend/` | read-only | Checked to see which commands the dashboard provision/reprovision dialogs generate. |

This plan lives in `docs/plans/` because every edit is in the docs repo. Working directory for all commands: `/home/ben/miru/workbench5/repos/docs` ("the docs repo"). Branch `docs/windows-agent-support` is already checked out from `main`; the PR base is `main`. Read-only repos live next to it under `/home/ben/miru/workbench5/repos/` and are read with `git -C ../<repo> show origin/main:<path>`.

## Purpose / Big Picture

The Miru Agent now ships a Windows build (an MSI installer that registers a Windows service), but the docs say "The Miru Agent only supports Linux." After this change a reader with a Windows 10 or 11 x64 machine can install, upgrade, and uninstall the agent, provision and reprovision it, manage the service, find its logs and data, write config schemas and file rules that target Windows paths, and see which features (the Device API and its SDKs) remain Linux-only. Linux readers see the same content as before, inside a "Linux" tab.

To see it working: run `pnpm dev` (or open the PR's Mintlify preview) and open `/developers/agent/install`; Install, Upgrade, Uninstall, and Poor connectivity each show Linux and Windows tabs, and the Windows tab contains an `msiexec` command for `miru-agent-<version>.msi`.

Merge gate: agent GitHub releases currently attach only `agent_Windows_x86_64.zip` and `miru_agent.pdb`, not the MSI. The PR stays draft and unmerged until an agent release publishes `miru-agent-<version>.msi` (see Validation and Acceptance).

## Progress

- [x] Pre-work: move this plan to `plans/active/`, re-verify facts against source repos, record discrepancies in Surprises.
- [x] M1: supported platforms, install, upgrade, uninstall, poor connectivity; test; commit.
- [x] M2: provisioning, reprovisioning, dashboard dialog notes, legacy script and Ansible notes; test; commit.
- [x] M3: agent commands, default permissions, security, architecture, file system access, overview; test; commit.
- [ ] M4: config file paths and file rule globs; test; commit.
- [ ] M5: device concepts, Device API Linux-only notes, CLI WSL 2 hint, quick start touch-ups; test; commit.
- [ ] M6: grep assertions, preview check, push, draft PR, preflight reports CLEAN.
- [ ] Merge gate: agent release with the MSI exists; asset name and SmartScreen note re-checked; commit; CLEAN again; PR marked ready.

## Surprises & Discoveries

- Latest agent release (`v0.10.3`, checked 2026-09-24) attaches only Linux assets (`.deb`, `.tar.gz`, SBOMs, checksums); neither the MSI nor `agent_Windows_x86_64.zip` is published yet. `build/.goreleaser.yaml` on `main` adds the Windows zip and PDB for future releases. Merge gate unchanged.
  Evidence: `gh release view -R mirurobotics/agent --json assets`.

- `path_not_absolute` is defined in `backend/internal/configs/domain/platform/errors.go`, not `filepaths.go`. No doc impact.

- The agent matches file rule globs with the `glob` crate 0.3.4 (`agent/src/filesys/files.rs` `glob()`), default `MatchOptions` (`case_sensitive: true`). Wildcard segments are matched case-sensitively on every OS; segments without wildcards are resolved by the file system, so on Windows (NTFS) literal directory and file names match regardless of case. The `glob` crate has no backslash escape (literal metacharacters are escaped with `[*]`), so on Linux `\` is a literal character, not an escape; the backend comment in `spec.go` ("on unix it is a glob escape") only means the backend does not split on it.
  Evidence: `~/.cargo/registry/src/*/glob-0.3.4/src/lib.rs` (`fill_todo`, `chars_eq`, `Pattern::escape`).

## Decision Log

- Decision: Windows install uses `Start-Process msiexec.exe ... -Wait -PassThru` and checks `ExitCode`.
  Rationale: `msiexec` is a GUI-subsystem program; a bare call in PowerShell returns before the install finishes. `agent/build/windows/README.md` uses the same form. Exit codes 0 and 3010 are success.
  Date/Author: 2026-09-24, plan author.

- Decision: Windows reprovisioning uses the token from the dashboard reprovision dialog, run through PowerShell with the service stopped.
  Rationale: reprovision tokens are device-bound and minted only by that dialog (5-minute lifetime); the dialog prints a Linux command, but the agent reads the token only from `MIRU_PROVISIONING_TOKEN`, so the same token works in PowerShell.
  Date/Author: 2026-09-24, plan author.

- Decision: do not document WinGet, the `agent_Windows_x86_64.zip` asset, any Device API TCP transport, the backend device-vs-release OS gate (backend PR #842, unmerged), or `os` fields on schemas, releases, or file rules.
  Rationale: none are shipped or exposed.
  Date/Author: 2026-09-24, plan author.

- Decision: `developers/agent/commands.mdx` keeps the per-action `###` headings OS-neutral, with Linux/Windows tabs under each.
  Rationale: preserves anchors `#restart`, `#stop`, `#start`, `#disable` that other pages link to.
  Date/Author: 2026-09-24, plan author.

- Decision: the CLI install page keeps "Windows is not supported" and adds a WSL 2 hint.
  Rationale: requested by the user.
  Date/Author: 2026-09-24, plan author.

- Decision: keep the existing "backslash escapes are not supported" statement for Linux globs instead of saying `\` escapes the next character (M4.4), and document case sensitivity as: wildcards are case-sensitive everywhere; on Windows, segments without wildcards follow the file system (case-insensitive).
  Rationale: the agent's `glob` crate has no backslash escape; see Surprises.
  Date/Author: 2026-09-24, implementer.

## Outcomes & Retrospective

(Summarize at completion.)

## Context and Orientation

### The docs repo

Mintlify site. Content root is `docs/`; `docs/developers/agent/install.mdx` is served at `/developers/agent/install`. Reusable fragments live in `docs/snippets/`, imported with `import X from '/snippets/...mdx';` and rendered as `<X />`. Components used: `<Tabs>`/`<Tab title="...">` (switchable panes), `<Steps>`/`<Step>`, `<Note>`, `<Warning>`, `<Danger>`, `<Tip>`, `<ParamField>`. Nested tabs have a precedent in `docs/snippets/references/cli/install/install.mdx` (macOS/Linux tabs, Linux containing APT/Manual tabs). Use ` ```powershell ` fences for Windows commands. There is no agent-version variable; write `<version>`.

Checks (from the docs repo):

- `./scripts/lint.sh`: custom Go linter in `tools/lint/`, ESLint MDX (`--max-warnings=0`), CSpell over all `.mdx` including code blocks (dictionary: the `"words"` array in root `cspell.json`, kept sorted), and `mint openapi-check`. Prints `All documentation lint checks passed.` on success. Relevant custom rules: `headingcase` (front-matter titles and markdown headings must be sentence case; inline code is masked; the only exempt capitalized tokens are the `allowlist()` map in `tools/lint/linter/headingcase/headingcase.go`, which does not contain Linux, Windows, PowerShell, or WSL, so those words may appear only as a heading's first word; `<Tab title>` is not checked), `nodoubledash` (no `--` in prose; flags in code spans are fine), import rules (one sorted contiguous block at the top, each line ending in `;`).
- `pnpm run test:lint`: lint smoke tests (`./tests/test-lint.sh`).
- `pnpm run validate`: `cd docs && mint validate` (MDX compile, navigation, internal links).
- `./scripts/preflight.sh`: full local gate (smoke tests, Go lint and coverage, lint.sh, audit, bats).
- CI: `.github/workflows/ci.yml`.

Commits use Conventional Commits (`docs(<scope>): ...`) and must be signed; check with `git cat-file -p HEAD | grep -c gpgsig` (expect 1).

### Agent facts (`agent` `origin/main`)

- Installer `miru-agent-<X.Y.Z>.msi` (no leading `v`), x64, per-machine, product name `Miru Agent`; Windows 10 and 11 x64 (Windows Server not certified, ARM64 not built). Installing an older MSI over a newer one fails with "A newer version of Miru Agent is already installed." Upgrade = install a newer MSI (major upgrade). Sources: `build/windows/miru-agent.wxs`, `build/windows/README.md`.
- Install from PowerShell opened with Run as administrator ("elevated"); exit code 0 = success, 3010 = success, restart required. Double-clicking the MSI runs the interactive installer.
- Uninstall: Settings > Apps > Installed apps > Miru Agent > Uninstall, or `msiexec /x` with the same-version MSI. Uninstall and upgrade keep `%ProgramData%\Miru`; deleting it removes the device identity and keys, so the machine must be provisioned again.
- Binary: `%ProgramW6432%\Miru\Agent\miru-agent.exe` (PowerShell: `& "$env:ProgramW6432\Miru\Agent\miru-agent.exe"`). Linux: `/usr/sbin/miru-agent`.
- Service: name `miru-agent`, display name `Miru Agent`, LocalSystem, startup Automatic, started at install; restarts after 10 s on 1st-3rd failure, counter resets after 1 day; stop can take up to about 30 s. No socket activation or idle exit (Linux uses `miru.service` + `miru.socket` and user `miru`). Sources: `build/windows/miru-agent.wxs`, `agent/src/windows/scm.rs`.
- Data: `%ProgramData%\Miru` holds `settings.json`, `device.json`, `agent_version`, `system_metadata.json`, and `auth\` (`private_key.pem`, `public_key.pem`, `token.json`), `resources\`, `events\`, `tmp\`, `logs\`. The MSI gives `logs`, `auth`, `tmp` a protected ACL: owner SYSTEM, full control only for SYSTEM and BUILTIN\Administrators. Linux equivalents: `/var/lib/miru`, `/var/log/miru`, `/srv/miru`.
- Deployed config files: no Windows default directory; the agent writes to the absolute path from the backend. As LocalSystem it can write almost anywhere.
- Logs: `%ProgramData%\Miru\logs\miru.log.YYYY-MM-DD-HH`, hourly rotation, files only (not the Windows Event Log). `miru-agent.exe --console` runs in the foreground, logs to stdout, stops on Ctrl-C. `RUST_LOG` overrides the level. Docs tell readers to stop the service before `--console`.
- CLI (`agent/src/cli/mod.rs`): `provision` [`--backend-host=`, `--mqtt-broker-host=`, `--device-name=`, `--check`]; `reprovision` [`--backend-host=`, `--mqtt-broker-host=`]; `version` (leading dashes stripped, so `--version` works); `--console`. Both `provision` and `reprovision` read the token only from env var `MIRU_PROVISIONING_TOKEN` (`agent/src/provisioning/shared.rs`; empty = missing). `provision --check` exits 0 provisioned, 3 not provisioned, 1 error.
- Default device name: `telemetry::SystemInfo::host_name` (`agent/src/provisioning/provision.rs`), which calls `sysinfo` `System::host_name()`. On Windows that is the DNS host name (`GetComputerNameExW(ComputerNamePhysicalDnsHostname)`), i.e. the output of `hostname`, not `$env:COMPUTERNAME` (NetBIOS, upper-case, 15-char limit).
- The agent reports `os` (`linux`|`windows`), `arch` (`x86_64`|`aarch64`), `hostname`, `os_version`, `kernel_version` on provision, reprovision, and sync.
- Device API: not served on Windows; the agent logs "the local device API server is not supported on windows; ignoring enable_socket_server".
- Release publishing (`build/.goreleaser.yaml`): ships `agent_Windows_x86_64.zip` and `miru_agent.pdb` only; MSI publishing and WinGet are follow-ups; no Authenticode signing yet.

### Provisioning tokens (openapi, frontend, current docs)

A provisioning token is a 5-minute bearer token. It is either workspace-bound (created on the provisioning tokens page, `docs/provision-devices/provisioning-tokens.mdx`, for bulk provisioning by `--device-name`) or device-bound (minted by the dashboard provision/reprovision dialogs for one device record). `POST /devices/reprovision` has no name field, so reprovisioning needs a device-bound token from the dashboard reprovision dialog (`apis/apps/backend-server/agent/openapi.yaml`, `bearerAuth` description). The frontend (`src/features/devices/components/provision/`) generates only Linux commands; the reprovision dialog shows:

    sudo systemctl stop miru
    sudo -u miru MIRU_PROVISIONING_TOKEN=<token> \
      /usr/sbin/miru-agent reprovision
    sudo systemctl start miru

A Windows user copies the `<token>` value from that command and runs the PowerShell equivalent within 5 minutes.

### Backend and core path rules (`origin/main`)

- `backend/internal/configs/ugc/filepaths.go` `SanitizeAbsFilepath`: a leading `/` is Linux; a drive letter (any case) + `:` + `\` or `/` is Windows (`core/pkg/lexpath/abs.go`). Anything else (UNC `\\server\share`, `\\?\`, `C:foo`, `\foo`, bare `C:`) fails with `path_not_absolute` ("path '<p>' is not an absolute unix or windows path").
- Windows paths then pass `core/pkg/ugc/files.go` `checkWindowsFilepath`, failing with `invalid_filepath` for: a segment ending in `.` or space, `:` after the drive, any of `<>"|?*`, control characters, or `..` after cleaning.
- Config instance paths: max 4096 bytes; extension must be `.json`, `.yaml`, or `.yml` (case-sensitive, opaque schemas exempt), so `C:\x\A.JSON` is rejected (`backend/internal/configs/domain/config_instances/filepath.go`). Windows paths compare case-insensitively for uniqueness (`backend/internal/configs/domain/platform/path.go` `CompareKey`).
- Omitted file path defaults to `/srv/miru/configs/<config-type-slug>.<ext>` (Linux only), so Windows schemas must set `instance_filepath` or slot `filepath` explicitly.
- All slot paths in a schema must target one OS (`mixed_config_schema_slot_os`, `backend/internal/configs/domain/config_schemas/errors.go`). A release's schemas and file rules must target one OS (`mixed_release_os`, `backend/internal/configs/services/releases/errors.go`).
- File rule globs (`backend/internal/configs/domain/filerules/spec.go`): absolute, at most 1024 bytes, no control characters, no `..` or empty segments. Windows globs split on both `\` and `/` after `C:\`; on Linux `\` is a glob escape, which contradicts the current docs sentence that backslash escapes are unsupported. Windows glob case sensitivity at match time is not established; verify in the agent before stating it.
- OpenAPI `Device` has `os` enum [`linux`, `windows`], `arch` enum [`x86_64`, `aarch64`], `hostname`, `os_version`, `kernel_version`, all nullable.

## Plan of Work

All paths are relative to the docs repo. "Wrap in tabs" means: a top-level `<Tabs>` with `<Tab title="Linux">` holding the existing content unchanged and `<Tab title="Windows">` holding the new content. Do not change screenshots. New markdown headings may contain Linux, Windows, PowerShell, or WSL only as the first word.

### M1: platforms and install

1. `docs/snippets/agent/supported-platforms.mdx`: line 1 becomes "The Miru Agent supports Linux and Windows. macOS is not supported." Keep the Linux list under a "**Linux**" label; add "**Windows**": "Windows 10 and Windows 11 (x64). Windows Server and ARM64 are not supported."
2. `docs/snippets/agent/install/install.mdx` (imported only by `docs/developers/agent/install.mdx`): wrap in tabs. Linux tab: current content (text, v0.9.0 Danger, APT/Manual tabs). Windows tab, as `<Steps>`:

       Download: download miru-agent-<version>.msi from the latest release's Assets on GitHub.
       Install: from PowerShell opened with Run as administrator:

           $msi = "C:\path\to\miru-agent-<version>.msi"
           $log = "C:\path\to\miru-agent-install.log"
           $p = Start-Process msiexec.exe -ArgumentList '/i',"`"$msi`"",'/qn','/norestart','/l*v',"`"$log`"" -Wait -PassThru
           $p.ExitCode

       ExitCode 0 means installed; 3010 means installed and a restart is required.
       Double-clicking the MSI runs the interactive installer instead.
       The installer registers and starts the miru-agent service.
       Verify: Get-Service miru-agent   (Status: Running)

   Leave `docs/snippets/agent/install/github-download.mdx` unchanged (Debian-only, rendered in Linux > Manual).
3. `docs/developers/agent/install.mdx`: line 9 becomes "This page details installing the Miru Agent on Linux (Debian package) and Windows (MSI installer)." `## Upgrade`: wrap in tabs; Windows = "Run the [install command](#install) with a newer MSI. The installer upgrades in place and keeps `%ProgramData%\Miru`, so the device stays provisioned. Installing an older MSI over a newer version fails; uninstall first to downgrade." `## Poor connectivity`: keep the prose, wrap the `<Steps>` in tabs; Windows steps = download the MSI on a connected computer, copy it (e.g. USB drive), run the Windows install command.
4. `docs/snippets/agent/install/uninstall.mdx`: wrap in tabs. Windows:

       Uninstall Miru Agent from Settings > Apps > Installed apps, or from an elevated PowerShell with the same-version MSI:

           Start-Process msiexec.exe -ArgumentList '/x','"C:\path\to\miru-agent-<version>.msi"','/qn','/norestart' -Wait -PassThru

       <Warning>: uninstalling keeps %ProgramData%\Miru (settings, device identity, keys, logs). For a clean removal delete it afterwards; the machine must then be provisioned again:

           Remove-Item -Recurse -Force "$env:ProgramData\Miru"

5. `cspell.json`: add flagged words (expect `msiexec`, `norestart`, maybe `qn`) in sorted position.

### M2: provisioning and reprovisioning

1. `docs/provision-devices/provisioning-tokens.mdx`:
   - `## Install the miru-agent package` (line 29): wrap in tabs; Linux keeps `<AptInstall />`; Windows links to [agent installation](/developers/agent/install#install).
   - `## Provision the device` (line 73): wrap the command and `<CmdBreakdown />` in tabs. Windows tab, elevated PowerShell:

         $env:MIRU_PROVISIONING_TOKEN = "<token>"
         & "$env:ProgramW6432\Miru\Agent\miru-agent.exe" provision
         Remove-Item Env:MIRU_PROVISIONING_TOKEN

     followed by a short breakdown: the token is passed only through the environment variable and cleared afterwards; administrator rights are needed because `%ProgramData%\Miru` is restricted to SYSTEM and Administrators.
   - `### Troubleshooting`: add that on Windows, starting fresh also requires deleting `%ProgramData%\Miru` after uninstalling (link `/developers/agent/install#uninstall`).
   - `### Arguments` `--device-name` ParamField (line 103): default becomes "the machine's hostname (the output of `hostname`)"; replace `Default: $HOSTNAME` with `Default: hostname`. Add a one-line tip: `provision --check` exits 0 provisioned, 3 not provisioned, 1 error, linking `/developers/agent/commands#provisioning-status`.
2. `docs/snippets/devices/provision/provision-cmd-breakdown.mdx` and `reprovision-cmd-breakdown.mdx`: unchanged (they describe the Linux command the dialogs show).
3. `docs/snippets/devices/provision/provision-dialog.mdx`: after the paste sentence add `<Note>The dashboard generates a Linux command. On Windows, copy the token (the value after MIRU_PROVISIONING_TOKEN=) and within 5 minutes run the Windows command from [provisioning tokens](/provision-devices/provisioning-tokens#provision-the-device) in PowerShell opened as administrator.</Note>` (put `MIRU_PROVISIONING_TOKEN=` in a code span). `reprovision-dialog.mdx`: same Note, but link `/provision-devices/reprovision#windows-reprovisioning`. `install-dialog.mdx`: add "On Windows, install the MSI instead" linking `/developers/agent/install#install`.
4. `docs/provision-devices/reprovision.mdx`: keep line 40 (reprovision tokens are still minted only by the dashboard). At the end of the page, after `<Reprovision />`, add `### Windows reprovisioning` (Windows is the first word, so headingcase passes; anchor `#windows-reprovisioning`): "The dialog shows a Linux command. On a Windows machine, open the reprovision dialog, copy the token (the value after `MIRU_PROVISIONING_TOKEN=`), and within 5 minutes run from PowerShell opened as administrator:"

       Stop-Service miru-agent
       $env:MIRU_PROVISIONING_TOKEN = "<token>"
       & "$env:ProgramW6432\Miru\Agent\miru-agent.exe" reprovision
       Remove-Item Env:MIRU_PROVISIONING_TOKEN
       Start-Service miru-agent

   Expected output line: `Successfully reprovisioned this device as <name>!`.
5. `docs/provision-devices/provisioning-script.mdx` and `docs/provision-devices/ansible.mdx`: add `<Note>This method supports Linux only.</Note>` after the intro.
6. `docs/getting-started/quick-start/provision-device.mdx` line 13: "lightweight `systemd` service" becomes "lightweight background service (a `systemd` service on Linux, a Windows service on Windows)".

### M3: agent operations

1. `docs/developers/agent/commands.mdx`:
   - Line 5 becomes: "The Miru Agent is managed with the operating system's service tools: `systemd` on Linux and the Windows service manager on Windows."
   - `## Version`: tabs; Windows `& "$env:ProgramW6432\Miru\Agent\miru-agent.exe" version`.
   - `## Logs`: wrap in tabs; Linux = current content including its `###` subsections (if headings inside tabs break the TOC in preview, turn them into bold labels). Windows: log location and hourly rotation; tail with

         Get-ChildItem "$env:ProgramData\Miru\logs" | Sort-Object LastWriteTime | Select-Object -Last 1 | Get-Content -Wait -Tail 100

     foreground option (`Stop-Service miru-agent`, then `& "$env:ProgramW6432\Miru\Agent\miru-agent.exe" --console`, Ctrl-C, `Start-Service miru-agent`); `RUST_LOG` sets the level; logs are not written to the Windows Event Log.
   - Rename `## Systemd` (line 77) to `## Service management`. Intro: Linux tab keeps the systemd paragraph and unit paths; Windows tab: service `miru-agent` ("Miru Agent"), LocalSystem, Automatic start, restarts after 10 s on the first three failures (counter resets daily). Under each `###` (Status, Restart, Stop, Start, Disable, Enable) add Linux/Windows tabs; Linux = existing content; Windows (elevated) = `Get-Service miru-agent`, `Restart-Service miru-agent`, `Stop-Service miru-agent` (may take up to 30 s), `Start-Service miru-agent`, `Set-Service miru-agent -StartupType Disabled`, `Set-Service miru-agent -StartupType Automatic`.
   - Add `## Provisioning status` with tabs: Linux `sudo -u miru /usr/sbin/miru-agent provision --check`; Windows `& "$env:ProgramW6432\Miru\Agent\miru-agent.exe" provision --check`; exit codes 0 provisioned, 3 not provisioned, 1 error.
2. `docs/snippets/agent/filesys/default-perms.mdx`: wrap in tabs. Windows: the agent runs as LocalSystem; table rows `%ProgramData%\Miru` (state, credentials, device identity; `auth`, `logs`, `tmp` restricted to SYSTEM and Administrators) and `%ProgramData%\Miru\logs` (hourly logs); no default config directory.
3. `docs/developers/agent/security.mdx`: `### Credential storage`: tabs; Windows = keys in `%ProgramData%\Miru\auth\`, ACL grants only SYSTEM and Administrators. `## Process sandboxing`: add a paragraph that on Windows the agent runs as a LocalSystem service with no equivalent of the systemd sandboxing directives. `## Network posture` / `## Local API access control`: add that on Windows there is no local socket or listening port because the Device API is not available.
4. `docs/developers/agent/architecture.mdx` line 64: `**4. Restart (systemd)**` becomes `**4. Restart**`.
5. `docs/developers/agent/filesys-access.mdx`: line 7 becomes "On Linux the agent runs as the `miru` system user..."; add a one-line note near the top that the Unix permission sections and `check-miru-access.sh` are Linux-only; add `## Windows` before `## Testing access`: the agent runs as LocalSystem and can read and write most paths without changes; to grant explicitly use `icacls`, e.g. `icacls "C:\ProgramData\Robot\configs" /grant "SYSTEM:(OI)(CI)M"`.
6. `docs/developers/agent/overview.mdx`: line 7 "lightweight systemd service" becomes OS-neutral as in M2.6; append "(Linux only)" to the local REST API bullet (line 12); note on line 20 that socket activation is Linux-only.

### M4: config file paths and file rules

1. `docs/snippets/references/cli/releases/create/schema-annotations.mdx`: in "instance file path" (lines 39-72) state the default applies to Linux only and Windows schemas must set an explicit path such as `C:\ProgramData\Robot\configs\mobility.json`; add that Windows example to Examples; add a short "Path rules" note: absolute Linux (`/...`) or Windows (drive letter, `:`, `\` or `/`); UNC and device paths, drive-relative (`C:foo`), and root-relative (`\foo`) paths are rejected (`path_not_absolute`); Windows paths may not contain `<>"|?*`, `:` after the drive, or a segment ending in `.` or a space (`invalid_filepath`); Windows paths are compared case-insensitively; extensions are case-sensitive (`.json`, not `.JSON`). In "instance slots" (lines 74-130) add that all slots in one schema must target the same OS (`mixed_config_schema_slot_os`). Append "(Linux)" to the `/srv/miru` Warning.
2. `docs/cfg-mgmt/concepts/schemas/annotations.mdx` `filepath` ParamField (lines 110-118): add a Windows example and the same-OS rule.
3. `docs/cfg-mgmt/concepts/config-instances.mdx` line 28: add `C:\ProgramData\Robot\configs\safety.yaml` to the examples.
4. `docs/snippets/file-rules/sources.mdx` and `docs/data-uploads/concepts/file-rules/rule-definition.mdx` (Glob patterns, lines ~41-80): "Absolute: starts with `/` (Linux) or a drive letter such as `C:\` (Windows)"; add example `C:\ProgramData\Robot\logs\*.log`; state that Windows patterns treat `\` and `/` as separators; correct the backslash sentence to match `spec.go` (on Linux `\` escapes the next character); change the "matching is case-sensitive" bullet only if the agent matcher confirms Windows behavior, else scope it to Linux and record in Surprises. Add that a release's schemas and file rules must all target the same OS (`mixed_release_os`), so keep separate releases for Linux and Windows devices.
5. `docs/data-uploads/define-file-rules.mdx`: add one sentence plus a Windows glob example next to `/var/log/robot/*.log`.

### M5: devices, Device API, CLI

1. `docs/concepts/devices/overview.mdx` line 18: OS-neutral wording. Add ParamFields for `os` (`linux`, `windows`), `arch` (`x86_64`, `aarch64`), `hostname`, `os_version`, `kernel_version`, each with `<NullableBadge />` (already imported) and "reported by the agent; empty until the agent reports it".
2. `docs/snippets/definitions/device.mdx`: add "Windows PC" to the example list.
3. New `docs/snippets/device-api/linux-only.mdx`:

       <Note>
         The Device API and Device API SDKs are currently available on Linux only. The Miru Agent for Windows does not serve the Device API.
       </Note>

   Import (sorted into the import block) and render near the top of `docs/developers/device-api/overview.mdx`, `authn.mdx`, `sdks.mdx`, `events.mdx`, `versions.mdx`.
4. `docs/snippets/references/cli/install/install.mdx` line 4 becomes: "The CLI is only available on macOS and Linux. Windows is not supported; on Windows, use [WSL 2](https://learn.microsoft.com/windows/wsl/install) and follow the Linux instructions." Add `WSL` to `cspell.json` if flagged.
5. `docs/getting-started/quick-start/deploy-configs.mdx` line ~41: add one sentence that Windows devices need schemas with Windows file paths, linking to the schema annotations page. Quick start examples and screenshots stay Linux.

## Concrete Steps

All commands run from `/home/ben/miru/workbench5/repos/docs`.

### Pre-work

    git branch --show-current            # expect: docs/windows-agent-support
    mkdir -p plans/active && mv plans/backlog/20260924-windows-agent-support.md plans/active/
    for r in agent backend core openapi frontend; do git -C ../$r fetch -q origin; done
    git -C ../agent show origin/main:build/windows/README.md
    git -C ../agent show origin/main:build/windows/miru-agent.wxs
    git -C ../agent show origin/main:agent/src/cli/mod.rs
    git -C ../agent show origin/main:agent/src/provisioning/shared.rs | grep -n MIRU_PROVISIONING_TOKEN
    git -C ../agent grep -n "not supported on windows" origin/main
    git -C ../agent show origin/main:build/.goreleaser.yaml | grep -n -i "msi\|windows"
    git -C ../backend show origin/main:internal/configs/domain/filerules/spec.go | sed -n 165,180p
    git -C ../backend grep -n "mixed_config_schema_slot_os\|mixed_release_os\|path_not_absolute" origin/main -- '*.go'
    git -C ../core show origin/main:pkg/ugc/files.go | grep -n -A30 "func checkWindowsFilepath"
    git -C ../frontend grep -n -i "windows\|powershell" origin/main -- src/features/devices/components/provision
    gh release view -R mirurobotics/agent --json assets -q '.assets[].name'

Expected: the frontend grep prints nothing; the release assets do not include an `.msi` yet. Record any mismatch with Context and Orientation in Surprises and adjust the affected milestone.

### Each milestone M1-M5

After editing:

    pnpm run test:lint       # expect: all smoke tests pass
    ./scripts/lint.sh        # expect: All documentation lint checks passed.
    pnpm run validate        # expect: success, no broken links
    git status --short       # only the milestone's files

Commit (list the exact files; M1 also adds `plans/active/20260924-windows-agent-support.md`); every message ends with `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`:

    git add <files>
    git commit -m "docs(agent): add Windows install, upgrade, and uninstall"          # M1
    git commit -m "docs(provisioning): add Windows provisioning and reprovisioning"   # M2
    git commit -m "docs(agent): add Windows service, logs, and security details"      # M3
    git commit -m "docs(cfg-mgmt): document Windows config paths and file rule globs" # M4
    git commit -m "docs(devices): add device system info and Device API Linux-only notes" # M5
    git cat-file -p HEAD | grep -c gpgsig   # expect: 1

### M6: validation, push, PR

    ./scripts/preflight.sh              # expect exit 0
    (run the grep assertions and preview check below)
    git push -u origin docs/windows-agent-support
    gh pr create --draft --base main --title "docs(agent): document Windows agent support" --body "<summary; merge gate; ends with the Claude Code attribution line>"
    gh pr checks <pr> --watch           # expect: all checks pass for the head SHA

Then run preflight on the branch until it reports `CLEAN`. Update this plan's Progress and Outcomes, commit (`docs(plans): update Windows agent support plan`), push, and re-confirm `CLEAN`.

## Validation and Acceptance

Tests for this docs change are: `pnpm run test:lint`, `./scripts/lint.sh`, `pnpm run validate`, `./scripts/preflight.sh`, the grep assertions, and the Mintlify preview check.

Grep assertions (from the docs repo):

    grep -rn "only supports Linux" docs/                               # no output
    grep -rn "## Systemd\|#systemd" docs/                              # no output
    grep -rni "winget\|agent_Windows_x86_64" docs/ | grep -v changelog # no output
    git diff main -- docs | grep '^+' | grep -iE "127\.0\.0\.1|tcp"   # no output
    grep -n "Windows is not supported" docs/snippets/references/cli/install/install.mdx  # 1 match
    grep -n "WSL 2" docs/snippets/references/cli/install/install.mdx                     # 1 match
    grep -rn "reprovision" docs/provision-devices/reprovision.mdx | grep miru-agent.exe  # 1 match
    grep -rln "linux-only.mdx" docs/developers/device-api/   # authn, events, overview, sdks, versions
    grep -rn "msiexec" docs/snippets/agent/install/          # install.mdx and uninstall.mdx
    grep -rn "mixed_release_os" docs/                        # at least one match
    grep -rln 'title="Windows"' docs/ | sort

The last command lists at least `docs/developers/agent/commands.mdx`, `docs/developers/agent/install.mdx`, `docs/developers/agent/security.mdx`, `docs/provision-devices/provisioning-tokens.mdx`, `docs/snippets/agent/filesys/default-perms.mdx`, `docs/snippets/agent/install/install.mdx`, `docs/snippets/agent/install/uninstall.mdx`.

Preview check: `pnpm dev` (or the PR's Mintlify preview if local preview fails; note that in Surprises), then open `/developers/agent/install`, `/developers/agent/commands`, `/provision-devices/provisioning-tokens`, `/provision-devices/reprovision`, `/developers/agent/filesys-access`, `/developers/cli/install`. Each tabbed section shows Linux and Windows tabs; Linux tabs match the pre-change text; code blocks highlight as PowerShell; the commands page TOC still lists Restart, Stop, Start, Disable, Enable; `/provision-devices/reprovision#windows-reprovisioning` scrolls to the new section.

Acceptance: a reader following only Windows content can install the MSI, see `Get-Service miru-agent` report Running, provision with a token, reprovision with the dashboard dialog token, find logs under `%ProgramData%\Miru\logs`, and learn that the Device API is unavailable on Windows and that the CLI requires WSL 2.

Completion: preflight must report `CLEAN` (CI green on the pushed branch head) before the PR leaves draft or the task is reported complete. Merge gate: even when `CLEAN`, the PR stays draft until an agent release publishes `miru-agent-<version>.msi`. Then re-check the asset name with `gh release view -R mirurobotics/agent --json assets -q '.assets[].name'`, add a SmartScreen/unsigned-publisher note to the Windows install tab if the MSI is still unsigned, commit (`docs(agent): finalize Windows installer instructions`), push, re-confirm `CLEAN`, and mark the PR ready.

## Idempotence and Recovery

All edits are text; lint, validate, preflight, and greps can be re-run freely. Fix failures with new commits, not amends. To back out a milestone, `git revert <sha>`. If the plan move is repeated, `mv` fails harmlessly once the file is in `plans/active/`. If `gh pr create` reports an existing PR, use `gh pr view --web` instead. If CSpell flags Windows terms, add them to `cspell.json` rather than rewording commands. If the published MSI asset name differs from `miru-agent-<version>.msi`, update the install and uninstall commands and re-run the checks. If headings inside tabs break the TOC, convert them to bold labels and note it in Surprises.

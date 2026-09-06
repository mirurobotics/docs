# Add Ansible usage documentation

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.


## Scope


| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` (`/home/ben/miru/workbench5/repos/docs`) | read-write | Add the Ansible page, nav entry, cross-links, and CSpell words. |
| `ansible-collection-agent/` (`/home/ben/miru/workbench5/repos/ansible-collection-agent`) | read-only | Source of truth for collection install, role variables, and secret handling. Do not change this repo. |

This plan lives in `docs/plans/` because the documentation work is written in the docs repo. Branch `docs/ansible-usage` already exists and matches `main`; do not recreate it.


## Purpose / Big Picture

After this change, a customer reading docs.mirurobotics.com can install the `mirurobotics.agent` Ansible collection from git, store an API key safely, and run the provision role or bundled playbook to install the Miru Agent and register devices. Locally, `pnpm run dev` shows an **Ansible** page in the Provision devices nav, between Provisioning tokens and Provisioning script, with working links from the overview, tokens, and agent install pages.


## Progress

- [x] Add `docs/cfg-mgmt/provision-devices/ansible.mdx`, nav slug, and CSpell words.
- [x] Add cross-links on overview, provisioning-tokens, and agent install (optional one-liner on the quick-start provision page).
- [ ] Run lint, `mint validate`, and preflight; fix until preflight is CLEAN.
- [ ] Push the branch, open a draft PR so CI runs, and confirm CI lint + audit + shell-tests are green on the pushed HEAD before leaving draft or reporting complete.


## Surprises & Discoveries

- Observation: CSpell did not already include `Ansible`, `ansible`, `Galaxy`, `deb822`, or `playbook`.
  Evidence: `cspell.json` `words` list before Milestone 1.


## Decision Log

- Decision: Include the optional one-liner on `getting-started/quick-start/provision-device.mdx`.
  Rationale: The plan allows it, and the page already points at tokens as an alternative; Ansible is the same kind of discoverability link.
  Date/Author: 2026-09-06 / implement agent


## Outcomes & Retrospective

Summarize at completion or major milestones.


## Context and Orientation

This repository is the Mintlify documentation site. Page files live under `docs/`. Navigation is `docs/docs.json` (product **Documentation**, group **Provision devices**). There is no Ansible mention anywhere yet.

**Ansible** is a tool for running the same setup steps on many machines. A **collection** is a packaged set of Ansible content. This collection’s fully qualified name is `mirurobotics.agent` (namespace `mirurobotics`, name `agent`, version `0.1.0`, MIT). A **role** is a reusable task set; a **playbook** is a YAML file that applies roles to inventory hosts. The role and the bundled playbook share the name `mirurobotics.agent.provision`. The playbook file is `playbooks/provision.yml` in the collection repo.

The **controller** is the machine that runs Ansible. **Targets** (also called hosts) are the devices being installed and provisioned. **Ansible Vault** is Ansible’s encrypted-variable store. A **provisioning token** is a one-time, 5-minute Platform API credential used to register a device; the role mints it on the controller so the long-lived API key never lands on targets.

**Status:** prototype. Not published to Ansible Galaxy (Ansible’s public collection index). Do not document `ansible-galaxy collection install mirurobotics.agent`. Interfaces may change. Customer-facing source of truth is the collection `README.md` and `CHANGELOG.md` (`## 0.1.0 (unreleased)`). Pin `version: main` as the README does.

**What the role does** (install and provision unless `miru_provision` is false):

1. Assert `miru_api_key` is defined and non-empty when `miru_provision` is true.
2. Install apt prereqs (`ca-certificates`, `python3-debian`) and add Miru’s apt source via `ansible.builtin.deb822_repository` (name `miru`, suite `stable`, component `main`, `signed_by` is the URL `https://packages.mirurobotics.com/apt/miru.gpg`). This is not the same as the manual keyring path in `docs/snippets/apt/setup.mdx` (`/etc/apt/keyrings/miru-archive-keyring.gpg`). Do not paste that snippet as “what Ansible does.”
3. `apt install` `miru-agent` or `miru-agent=<version>`; enable and start systemd unit `miru`.
4. If provisioning: skip when any file exists under `/var/lib/miru/auth`; otherwise POST `{{ miru_api_base_url }}/provisioning_tokens` on the controller (`delegate_to: localhost`, headers `Miru-Version: {{ miru_api_version }}` and `X-API-Key: {{ miru_api_key }}`), then run `/usr/sbin/miru-agent provision --device-name={{ miru_device_name }}` on the host as user `miru` with `MIRU_PROVISIONING_TOKEN` in the environment.

**Role variables** (exact names; copy into the new page table):

- `miru_api_key` — no default — required only when `miru_provision` is true — Platform API key; controller-side only.
- `miru_provision` — default `true` — `false` installs without provisioning (image baking); no API key required.
- `miru_agent_version` — default `""` (latest) — apt version such as `0.10.1` (no `v` prefix).
- `miru_device_name` — default `{{ inventory_hostname }}`.
- `miru_api_base_url` — default `https://api.mirurobotics.com/beta`.
- `miru_api_version` — default `2026-05-06.rainier`.
- `miru_apt_url` — default `https://packages.mirurobotics.com/apt`.
- `miru_apt_key_url` — default `https://packages.mirurobotics.com/apt/miru.gpg`.
- `miru_apt_architecture` — mapped from `ansible_facts.architecture` (`x86_64` → `amd64`, `aarch64` → `arm64`).
- `miru_provision_no_log` — default `true`.

`argument_specs.yml` marks `miru_api_key` required-false because of install-only mode; the role still asserts it at runtime when provisioning.

The role does **not** read `MIRU_API_KEY` from the environment (no `lookup('env')`). Callers must set `miru_api_key`. The bundled playbook maps `miru_api_key: "{{ vault_miru_api_key }}"` and targets inventory group `robots`. Safe documented paths: Ansible Vault `vault_miru_api_key`; CI extra-vars such as `-e vault_miru_api_key="$MIRU_API_KEY"`; or `miru_provision: false`. Never show a recommended plaintext key in a committed playbook.

**Customer prereqs:** ansible-core >= 2.15 on the controller (`requires_ansible: ">=2.15.0"`); collection installed from git; apt+systemd Linux targets; privilege escalation (`become`); controller can reach `https://api.mirurobotics.com`; targets can reach `https://packages.mirurobotics.com` and later the control plane; API key scopes `devices:provision` and `provisioning_tokens:write` unless install-only; agent >= v0.9.0 if provisioning. Role platforms: Ubuntu focal/jammy/noble, Debian all. Molecule CI tests Ubuntu 22.04 and 24.04 only. Jetson / Raspberry Pi OS appear on the agent install page, not in collection tests.

**Out of scope for the page and this plan:** Galaxy install; `miru-agent provision --check` (not shipped); programmatic reprovision (dashboard-only); changing `ansible-collection-agent`.

**Existing pages to extend (do not duplicate their walkthroughs; link):**

- `docs/cfg-mgmt/provision-devices/overview.mdx` — three methods today: dashboard, tokens, legacy script.
- `docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx` — apt snippet, API key, curl/Python mint, `miru-agent provision`.
- `docs/developers/agent/install.mdx` — Debian install only; says it does not cover provisioning.
- `docs/getting-started/quick-start/provision-device.mdx` — dashboard quick start; mentions tokens as an alternative.
- `docs/cfg-mgmt/provision-devices/reprovision.mdx` — dashboard-only; Ansible page links here, not the reverse.
- `docs/admin/apikeys.mdx` and `docs/developers/platform-api/authz.mdx` — key creation and scopes (`devices:provision`, `provisioning_tokens:write`). Tokens page already owns the key walkthrough; do not copy it.

**Mintlify page rules** (enforced by `./scripts/lint.sh` → `tools/lint`):

- Frontmatter `title` is the H1. Body starts at `##`. No body `#`.
- heading-case: sentence case; first token initial-cap; later tokens lowercase unless allowlisted or backticked. Allowlist includes API, CLI, Miru, Agent, GitHub, Python — **not** Ansible. `title: "Ansible"` passes. `## Install via Ansible` fails. `## Install via \`Ansible\`` passes. Prefer backticks over editing `tools/lint/linter/headingcase/headingcase.go`.
- Imports: default import with trailing `;`; named `{ Framed }` with spaces, path ends `.jsx`, `;`. Sorted by path, case-insensitive; no blank lines in the import block; every imported name used.
- Callouts: `Danger`, `Tip`, `Info`, `Note`, `Tabs`/`Tab`, `Steps`/`Step`, `CodeGroup`, `ParamField`.
- Prose `--` is forbidden (`no-double-dash`). Use an em dash — or backticks for flags such as `--device-name`.
- Images, if any, must be `https://assets.mirurobotics.com/`. This page needs no new image.

**Reuse, do not duplicate:** `/snippets/agent/supported-platforms.mdx` and `/snippets/agent/install/verify.mdx`.

**Lint / CI:** `./scripts/preflight.sh` runs `pnpm run test:lint`, Go linter + covgate, `./scripts/lint.sh` (MDX prose, ESLint MDX, CSpell, OpenAPI), `./scripts/audit.sh`, and bats. It does **not** print the word CLEAN; CLEAN means exit 0 with no warnings. It does **not** run `pnpm run validate`. CI (`.github/workflows/ci.yml`) on every PR also runs `pnpm run validate` (`cd docs && mint validate`, strict, fails on warnings). CI on this feature branch runs only when a PR exists; Milestone 4 therefore requires opening a draft PR. Custom-linter jobs skip unless `tools/lint/**` changes.


## Plan of Work

Work only in the docs repo, on `docs/ansible-usage`. One commit per milestone below.

**1. New page.** Create `docs/cfg-mgmt/provision-devices/ansible.mdx`. Frontmatter title `"Ansible"`. Import `Verify` from `/snippets/agent/install/verify.mdx` then `SupportedPlatforms` from `/snippets/agent/supported-platforms.mdx` (path order). No `Framed` hero (no Ansible asset). Body sections, in order:

1. Prototype `Note`: not on Galaxy; interfaces may change; install from git only.
2. What it does: the four-step list from Context; link to [provisioning tokens](/cfg-mgmt/provision-devices/provisioning-tokens) and [agent install](/developers/agent/install). State that the role does both install and provision unless `miru_provision` is false, and that the API key stays on the controller.
3. Requirements: ansible-core >= 2.15; `<SupportedPlatforms />` plus a sentence that targets need systemd and apt, and that the collection is tested on Ubuntu 22.04 and 24.04; agent >= v0.9.0 when provisioning; API key with `devices:provision` and `provisioning_tokens:write` linking to `/admin/apikeys` and `/developers/platform-api/authz`.
4. Install the collection — only this git `requirements.yml` and `ansible-galaxy collection install -r requirements.yml`. No Galaxy name, no version-tag story.

        collections:
          - name: https://github.com/mirurobotics/ansible-collection-agent.git
            type: git
            version: main

5. Store the API key: Vault variable `vault_miru_api_key`, or inject via extra-vars; never commit a plaintext key. The role does not auto-load `MIRU_API_KEY`.
6. Run the role in a custom play **or** the bundled playbook `ansible-playbook -i inventory mirurobotics.agent.provision`. The bundled playbook requires inventory group `robots` and `vault_miru_api_key`. Recommended role usage:

        - name: Provision Miru devices
          hosts: robots
          roles:
            - role: mirurobotics.agent.provision
              vars:
                miru_api_key: "{{ vault_miru_api_key }}"
                miru_agent_version: "0.10.1"   # optional; omit for latest

7. Role variables table using the exact names and defaults from Context.
8. Install-only: `miru_provision: false` for image baking; no API key.
9. Verify: `<Verify />` (Devices page Activating → Online). Limitations: already-provisioned means any file in `/var/lib/miru/auth` (stale creds after a hardware swap still count); reprovision is [dashboard-only](/cfg-mgmt/provision-devices/reprovision).

Do not document Galaxy install, `miru-agent provision --check`, curl/Python token mint, apt keyring setup, or API-key creation steps.

**2. Nav.** In `docs/docs.json`, Provision devices pages, insert `cfg-mgmt/provision-devices/ansible` after `provisioning-tokens` and before `provisioning-script`.

**3. CSpell.** In `cspell.json` `words`, add real terms that fail CSpell after the first lint run. Expect at least `Ansible` and `ansible`. Add `playbook`, `deb822`, `Galaxy`, and similar only if the page uses them and CSpell flags them. Do not add invented words.

**4. Cross-links** (short sentences, not copies):

- `overview.mdx`: change “three methods” to four; add Ansible after tokens and before the legacy script; note that Ansible installs the agent as part of the role (dashboard and tokens still require a prior install). Update the sentence that only points at dashboard or tokens after install.
- `provisioning-tokens.mdx`: after the “integrate with existing infrastructure” bullet, link to `/cfg-mgmt/provision-devices/ansible`.
- `developers/agent/install.mdx`: after “does not cover provisioning,” add that the Ansible collection can install and provision a fleet, linking the new page.
- Optional: `getting-started/quick-start/provision-device.mdx` — one sentence next to the existing tokens mention.
- Do not add a reverse link from `reprovision.mdx` or duplicate the API-key walkthrough on `admin/apikeys.mdx`.

**5. Headings and dashes.** Keep mid-heading Ansible in backticks. Never write raw `--` in MDX prose.

**6. Validate.** Run the commands in Concrete Steps until preflight is CLEAN and `pnpm run validate` passes. Push, open a draft PR if none exists (required so CI runs on this branch), and keep the PR draft until CI is green on that HEAD.


## Concrete Steps

All commands below run from the docs repo root: `/home/ben/miru/workbench5/repos/docs` (also called `docs/` repo root). Do not run git commands from the workbench root or from `ansible-collection-agent/`.

Confirm the branch (do not create or reset it):

    git branch --show-current
    git status

Expected: `docs/ansible-usage` and a clean tree.

### Milestone 1 — Page, nav, and CSpell

Create `docs/cfg-mgmt/provision-devices/ansible.mdx` as specified in Plan of Work. Edit `docs/docs.json` to insert the nav slug. Add CSpell words you already know you will use (`Ansible`, `ansible`); add the rest after lint reports them.

    git add docs/cfg-mgmt/provision-devices/ansible.mdx docs/docs.json cspell.json
    git commit -m "$(cat <<'EOF'
    docs: add Ansible usage page for device provisioning

    EOF
    )"
    git status

Expected: one new commit on `docs/ansible-usage`; working tree clean except unrelated files. If the hook reformats files, stage those edits and create a **new** commit (do not amend unless the hook-amend conditions in repo commit policy all apply). Do not pass `--no-verify`.

### Milestone 2 — Cross-links

Edit the overview, tokens, and agent install pages (and optionally the quick-start provision page) as specified in Plan of Work.

    git add docs/cfg-mgmt/provision-devices/overview.mdx \
      docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx \
      docs/developers/agent/install.mdx
    # include the quick-start file only if you edited it
    git commit -m "$(cat <<'EOF'
    docs: cross-link Ansible from provision and install pages

    EOF
    )"
    git status

### Milestone 3 — Local lint, validate, and preflight

    pnpm run test:lint
    ./scripts/lint.sh
    pnpm run validate
    ./scripts/preflight.sh

Expected from `./scripts/lint.sh`: `All documentation lint checks passed.` Expected from preflight: exit 0 with no warnings (this is CLEAN; the script does not print the word CLEAN). `pnpm run validate` must also succeed — CI lint runs it and it is stricter than preflight.

If CSpell fails, add only the flagged real terms to `cspell.json` and rerun. If heading-case fails on Ansible, backtick the word in that heading rather than changing the linter. If `no-double-dash` fails, replace prose `--` with an em dash or backticks.

Optional preview (leave running while you click through):

    pnpm run dev

Open `/cfg-mgmt/provision-devices/ansible` and the three cross-linked pages. Confirm the nav order and that links resolve.

Optional extra check (not in preflight or CI):

    cd docs && mint broken-links

If milestone 3 produced file changes, stage only the files this work may have changed (do not use `git add -u`, `git add .`, or `git add -A`):

    git add docs/cfg-mgmt/provision-devices/ansible.mdx \
      docs/docs.json \
      cspell.json \
      docs/cfg-mgmt/provision-devices/overview.mdx \
      docs/cfg-mgmt/provision-devices/provisioning-tokens.mdx \
      docs/developers/agent/install.mdx
    # include the quick-start file only if you edited it:
    # git add docs/getting-started/quick-start/provision-device.mdx
    git commit -m "$(cat <<'EOF'
    docs: fix lint on Ansible usage docs

    EOF
    )"

If lint was already clean, skip this commit. Do not create an empty commit.

### Milestone 4 — Confirm CLEAN and CI green

Re-run preflight and validate on the final tree:

    ./scripts/preflight.sh
    pnpm run validate

Push the existing branch:

    git push -u origin HEAD

Open a draft pull request if one does not already exist. This is a required Milestone 4 step: `.github/workflows/ci.yml` runs `lint`, `audit`, and `shell-tests` on this feature branch only for `pull_request` events, not for a push to `docs/ansible-usage`.

    gh pr view --json url,isDraft >/dev/null 2>&1 || gh pr create --draft --title "docs: add Ansible usage page" --body "$(cat <<'EOF'
    Add Ansible collection usage docs for installing the agent and provisioning devices.

    EOF
    )"

Wait until CI jobs `lint`, `audit`, and `shell-tests` are green on the pushed HEAD. Opening the PR is how CI becomes runnable; it is not itself acceptance. Do not mark the PR ready and do not report complete until Validation and Acceptance is met. If CI fails, fix on this branch, commit, push, and re-check the new HEAD.


## Validation and Acceptance

Phrase of success is behavior, not “files exist.”

- `pnpm run dev` shows **Ansible** in Provision devices after Provisioning tokens and before Provisioning script. Opening `/cfg-mgmt/provision-devices/ansible` shows the prototype warning, git install (not Galaxy), Vault / extra-vars secret guidance, role and bundled-playbook invocation (`robots` + `vault_miru_api_key`), the role-variables table with the names in Context, install-only `miru_provision: false`, and verify plus auth-dir / reprovision limits.
- Overview lists Ansible as a fourth method. Tokens and agent install pages link to the new page. Clicking those links reaches the Ansible page.
- From the docs repo root, `./scripts/lint.sh` prints `All documentation lint checks passed.` `pnpm run validate` exits 0. `./scripts/preflight.sh` exits 0 with no warnings (CLEAN).
- Preflight must report CLEAN (CI green on the pushed branch head) before the PR leaves draft or the task is reported complete. For this content-only change, that means CI jobs `lint`, `audit`, and `shell-tests` green; `lint-custom-linter` / `test-custom-linter` are skipped. The draft PR in Milestone 4 is required so those jobs run; do not report complete after a branch push alone.

A reader of the published page can follow it to install the collection and provision devices without opening the collection repo. Do not treat “added ansible.mdx” as acceptance by itself.


## Idempotence and Recovery

Page, nav, CSpell, and cross-link edits are safe to repeat; overwrite in place. Lint, validate, and preflight are read-only aside from the Go linter binary built under `tools/lint/lint`. `pnpm run dev` can be restarted. Re-running `git commit` after a successful milestone is unnecessary; skip if `git status` is clean.

If a milestone commit is missing files, add them and make a new commit. If a hook rejects a commit, fix the issue and create a **new** commit; do not amend a rejected commit and do not skip hooks.

If you need to discard local edits to a file: `git checkout -- <path>` from the docs repo root. To undo a milestone commit that has not been pushed: ask before resetting; prefer a new fix commit once anything has been pushed.

Do not modify `/home/ben/miru/workbench5/repos/ansible-collection-agent`. Do not recreate branch `docs/ansible-usage`. Do not commit secrets or a plaintext `miru_api_key` example.

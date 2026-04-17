#!/usr/bin/env bash
set -euo pipefail

# Resolve the repo root relative to this script so it works from any cwd.
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"

usage() {
  cat <<EOF
Usage: $(basename "$0") <ref> <environment>

Deploy a git ref to a target environment branch.

Arguments:
  ref           Git ref to deploy (tag, commit SHA, or branch).
                Must be reachable from main.
  environment   Target environment (staging, uat, production).

The script resolves the ref to a SHA, verifies it is an ancestor of
origin/main, and force-pushes it to the environment branch.

Exit codes:
  0  Success or no-op (target already at the requested SHA).
  1  Error (bad ref, not on main, diverged push race).

When running in GitHub Actions, the script writes to \$GITHUB_OUTPUT and
\$GITHUB_STEP_SUMMARY. Outside CI those writes are silently skipped.
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ $# -ne 2 ]]; then
  usage >&2
  exit 1
fi

ref="$1"
environment="$2"

# Write a key=value pair to $GITHUB_OUTPUT if the variable is set.
gh_output() {
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    echo "$1" >> "$GITHUB_OUTPUT"
  fi
}

cd "${repo_root}"

# ---------------------------------------------------------------------------
# 1. Resolve ref and verify it is reachable from main
# ---------------------------------------------------------------------------
git fetch origin --tags --force

if ! sha=$(git rev-parse --verify "${ref}^{commit}" 2>/dev/null); then
  echo "::error::Could not resolve '${ref}' to a commit."
  exit 1
fi

if ! git merge-base --is-ancestor "$sha" origin/main; then
  echo "::error::Commit ${sha} is not reachable from main. Only code merged to main can be deployed."
  exit 1
fi

echo "Resolved ${ref} → ${sha} (verified on main)"
gh_output "sha=${sha}"

# ---------------------------------------------------------------------------
# 2. Check target branch state
# ---------------------------------------------------------------------------
target_sha=""
if git ls-remote --exit-code --heads origin "${environment}" >/dev/null 2>&1; then
  git fetch origin "${environment}:refs/remotes/origin/${environment}"
  target_sha=$(git rev-parse "origin/${environment}")
fi

if [[ "${target_sha}" == "${sha}" ]]; then
  echo "Target ${environment} already at ${sha} — no-op."
  gh_output "noop=true"
  gh_output "target_sha=${target_sha}"
  noop=true
else
  gh_output "noop=false"
  gh_output "target_sha=${target_sha}"
  noop=false
fi

# ---------------------------------------------------------------------------
# 3. Push to target branch (unless no-op)
# ---------------------------------------------------------------------------
if [[ "${noop}" != "true" ]]; then
  if [[ -n "${target_sha}" ]]; then
    git push --force-with-lease="refs/heads/${environment}:${target_sha}" \
      origin "${sha}:refs/heads/${environment}"
  else
    git push origin "${sha}:refs/heads/${environment}"
  fi
fi

# ---------------------------------------------------------------------------
# 4. Summary
# ---------------------------------------------------------------------------
if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  {
    echo "### Deploy ${ref} → ${environment}"
    echo ""
    echo "- **Ref:** \`${ref}\` → \`${sha}\`"
    echo "- **Previous ${environment} SHA:** \`${target_sha:-"(none — first deploy)"}\`"
    if [[ "${noop}" == "true" ]]; then
      echo "- **Result:** no-op"
    else
      echo "- **Result:** deployed"
    fi
  } >> "$GITHUB_STEP_SUMMARY"
fi

#!/usr/bin/env bash
set -euo pipefail

# Resolve the repo root relative to this script so it works from any cwd.
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"

if ! command -v pnpm >/dev/null 2>&1; then
  echo "pnpm is required to run security audit. Install pnpm and rerun ./scripts/audit.sh." >&2
  exit 1
fi

cd "${repo_root}"

echo "== Security Audit =="
pnpm audit --ignore-registry-errors

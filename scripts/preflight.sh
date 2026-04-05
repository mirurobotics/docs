#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

echo "=== Lint Smoke Tests ==="
pnpm run test:lint
echo ""

echo "=== Lint ==="
./scripts/lint.sh
echo ""

echo "=== Audit ==="
./scripts/audit.sh

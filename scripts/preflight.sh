#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

echo "=== Lint Smoke Tests ==="
pnpm run test:lint
echo ""

echo "=== Go Lint (tools/lint) ==="
LINT_FIX=0 ./tools/lint/scripts/lint.sh
echo ""

echo "=== Go Coverage (tools/lint) ==="
./tools/lint/scripts/covgate.sh
echo ""

echo "=== Lint ==="
./scripts/lint.sh
echo ""

echo "=== Audit ==="
./scripts/audit.sh
echo ""

echo "=== Shell Script Tests ==="
bats pub/scripts/agent/check-miru-access_test.bats

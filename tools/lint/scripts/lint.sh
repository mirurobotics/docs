#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<'EOF'
Usage: LINT_FIX=0 ./tools/lint/scripts/lint.sh

Runs gotools (custom linter + gofumpt + golangci-lint) against
docs/tools/lint/. Set LINT_FIX=0 to run in check-only mode (for CI
and preflight); omit or set LINT_FIX=1 for auto-fix mode (default
for local runs).
EOF
}

case "${1:-}" in
	-h|--help)
		usage
		exit 0
		;;
esac

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
lint_dir="$(cd -- "${script_dir}/.." && pwd)"

FIX="--fix"
if [ "${LINT_FIX:-1}" = "0" ]; then
	FIX="--fix=false"
fi

cd "${lint_dir}"
exec go tool miru lint --paths="${lint_dir}" ${FIX}

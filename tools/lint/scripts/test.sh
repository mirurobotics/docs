#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<'EOF'
Usage: ./tools/lint/scripts/test.sh [-- extra args]

Runs Go tests for docs/tools/lint/ via go tool miru test.
Arguments after -- are passed through to go test.
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

cd "${lint_dir}"
exec go tool miru test -- "$@"

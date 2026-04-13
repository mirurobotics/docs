#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<'EOF'
Usage: ./tools/lint/scripts/covgate.sh [threshold]

Checks per-package test coverage against thresholds for docs/tools/lint/.
Uses .covgate files per package; falls back to the supplied threshold
(default 90.0). Pass a numeric threshold as the first argument to override.
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
exec go tool miru covgate \
	--packages="./..." \
	--src-prefix="." \
	--default-threshold="${1:-90.0}"

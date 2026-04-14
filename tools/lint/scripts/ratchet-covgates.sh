#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<'EOF'
Usage: ./tools/lint/scripts/ratchet-covgates.sh

Ratchets up per-package .covgate thresholds for docs/tools/lint/.
Measures actual coverage and updates .covgate files when coverage
exceeds the current threshold.
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
exec go tool miru covratchet \
	--packages="./..." \
	--src-prefix="."

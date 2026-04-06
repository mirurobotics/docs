#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"

usage() {
	cat <<-EOF
	Usage: $(basename "$0") <api-version>

	Generate SSE event type MDX pages from the Device API OpenAPI spec.

	Arguments:
	  api-version    Device API version (e.g. v0.2.1)

	The script resolves the spec and output paths from the version:
	  Spec:   docs/references/device-api/<version>/api.yaml
	  Output: docs/references/device-api/<version>/events/

	Examples:
	  $(basename "$0") v0.2.1
	EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
		-h|--help)
			usage
			exit 0
			;;
		-*)
			echo "Unknown option: $1" >&2
			usage >&2
			exit 1
			;;
		*)
			break
			;;
	esac
done

if [[ $# -ne 1 ]]; then
	echo "Error: expected 1 argument (api-version), got $#" >&2
	usage >&2
	exit 1
fi

version="$1"
spec_path="${repo_root}/docs/references/device-api/${version}/api.yaml"
output_dir="${repo_root}/docs/references/device-api/${version}/events"

if [[ ! -f "${spec_path}" ]]; then
	echo "❌ Spec not found: ${spec_path}" >&2
	exit 1
fi

# activate the api/ venv (create + install deps if needed)
venv_dir="${repo_root}/api/.venv"
if [[ ! -d "${venv_dir}" ]]; then
	python3 -m venv "${venv_dir}"
fi
# shellcheck disable=SC1091
source "${venv_dir}/bin/activate"
pip install -q pyyaml 2>/dev/null

python3 "${repo_root}/api/generate_event_pages.py" "${spec_path}" "${output_dir}"

#!/usr/bin/env bash
set -euo pipefail

# Resolve the repo root relative to this script so it works from any cwd.
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"

usage() {
	cat <<-EOF
	Usage: $(basename "$0") [-h|--help]

	Run tests for api/generate_event_pages.py to verify codegen correctness,
	idempotence, extensibility, and error handling.
	EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
		-h|--help)
			usage
			exit 0
			;;
		*)
			echo "Unknown option: $1" >&2
			usage >&2
			exit 1
			;;
	esac
done

events_dir="docs/references/device-api/v0.2.1/events"
config_file="api/events.yaml"
spec_file="docs/references/device-api/v0.2.1/api.yaml"
generator="api/generate_event_pages.py"

test_zero_diff() {
	local name="test_zero_diff"
	local tmp_dir
	tmp_dir="$(mktemp -d)"

	# Back up the two MDX files
	cp "${repo_root}/${events_dir}/deployment-deployed.mdx" "${tmp_dir}/deployment-deployed.mdx"
	cp "${repo_root}/${events_dir}/deployment-removed.mdx" "${tmp_dir}/deployment-removed.mdx"

	# Run the generator
	python3 "${repo_root}/${generator}" "${repo_root}/${spec_file}" "${repo_root}/${config_file}" "${repo_root}" >/dev/null

	# Diff each generated file against backup
	local failed=0
	for f in deployment-deployed.mdx deployment-removed.mdx; do
		if ! diff -u "${tmp_dir}/${f}" "${repo_root}/${events_dir}/${f}" >/dev/null 2>&1; then
			echo "FAIL: ${name} - ${f} differs from original"
			diff -u "${tmp_dir}/${f}" "${repo_root}/${events_dir}/${f}" || true
			failed=1
		fi
	done

	# Restore originals from backups
	cp "${tmp_dir}/deployment-deployed.mdx" "${repo_root}/${events_dir}/deployment-deployed.mdx"
	cp "${tmp_dir}/deployment-removed.mdx" "${repo_root}/${events_dir}/deployment-removed.mdx"

	# Clean up temp dir
	rm -rf "${tmp_dir}"

	if [[ "${failed}" -ne 0 ]]; then
		exit 1
	fi

	echo "PASS: ${name}"
}

test_idempotence() {
	local name="test_idempotence"

	# Run generator once
	python3 "${repo_root}/${generator}" "${repo_root}/${spec_file}" "${repo_root}/${config_file}" "${repo_root}" >/dev/null

	# Capture sha256sum of both MDX files
	local sum1_deployed sum1_removed
	sum1_deployed="$(sha256sum "${repo_root}/${events_dir}/deployment-deployed.mdx" | awk '{print $1}')"
	sum1_removed="$(sha256sum "${repo_root}/${events_dir}/deployment-removed.mdx" | awk '{print $1}')"

	# Run generator again
	python3 "${repo_root}/${generator}" "${repo_root}/${spec_file}" "${repo_root}/${config_file}" "${repo_root}" >/dev/null

	# Capture sha256sum again
	local sum2_deployed sum2_removed
	sum2_deployed="$(sha256sum "${repo_root}/${events_dir}/deployment-deployed.mdx" | awk '{print $1}')"
	sum2_removed="$(sha256sum "${repo_root}/${events_dir}/deployment-removed.mdx" | awk '{print $1}')"

	# Assert checksums match
	if [[ "${sum1_deployed}" != "${sum2_deployed}" ]]; then
		echo "FAIL: ${name} - deployment-deployed.mdx checksum changed (${sum1_deployed} != ${sum2_deployed})"
		exit 1
	fi
	if [[ "${sum1_removed}" != "${sum2_removed}" ]]; then
		echo "FAIL: ${name} - deployment-removed.mdx checksum changed (${sum1_removed} != ${sum2_removed})"
		exit 1
	fi

	# Assert git diff --stat is empty
	local diff_stat
	diff_stat="$(cd "${repo_root}" && git diff --stat "${events_dir}/")"
	if [[ -n "${diff_stat}" ]]; then
		echo "FAIL: ${name} - git diff --stat is not empty:"
		echo "${diff_stat}"
		exit 1
	fi

	echo "PASS: ${name}"
}

test_extensibility() {
	local name="test_extensibility"
	local tmp_spec tmp_config
	tmp_spec="$(mktemp --suffix=.yaml)"
	tmp_config="$(mktemp --suffix=.yaml)"

	# Inject a synthetic event schema into a temp copy of the spec
	python3 -c "
import yaml, sys
with open(sys.argv[1], 'r') as f:
    spec = yaml.safe_load(f)
spec['components']['schemas']['TestSyntheticEvent'] = {
    'title': 'TestSyntheticEvent',
    'type': 'object',
    'description': 'Payload for \`test.synthetic\` events.',
    'required': ['deployment_id', 'status'],
    'properties': {
        'deployment_id': {'type': 'string', 'description': 'ID of the deployment.', 'example': 'dpl_test'},
        'status': {'\$ref': '#/components/schemas/DeploymentStatus'},
        'test_at': {'type': 'string', 'format': 'date-time', 'description': 'Timestamp.', 'example': '2026-01-01T00:00:00Z'}
    },
    'example': {'deployment_id': 'dpl_test', 'status': 'deployed', 'test_at': '2026-01-01T00:00:00Z'}
}
with open(sys.argv[2], 'w') as f:
    yaml.dump(spec, f, default_flow_style=False, sort_keys=False)
" "${repo_root}/${spec_file}" "${tmp_spec}"

	# Create a temp config with the synthetic event entry
	cat > "${tmp_config}" <<-YAML
	output_dir: ${events_dir}
	events:
	  test.synthetic:
	    description: "A synthetic event for testing."
	    body: "This is synthetic body prose."
	    field_annotations:
	      status: 'Always deployed for this test event.'
	YAML

	# Run generator with temp spec and temp config
	python3 "${repo_root}/${generator}" "${tmp_spec}" "${tmp_config}" "${repo_root}" >/dev/null

	local synth_file="${repo_root}/${events_dir}/test-synthetic.mdx"

	# Assert test-synthetic.mdx exists
	if [[ ! -f "${synth_file}" ]]; then
		echo "FAIL: ${name} - test-synthetic.mdx was not created"
		rm -f "${tmp_spec}" "${tmp_config}"
		exit 1
	fi

	# Grep for expected content
	local patterns=(
		'title: "test.synthetic"'
		'description: "A synthetic event for testing."'
		'This is synthetic body prose.'
		'## Event Data'
		'<ResponseExample>'
		'<ResponseField name="deployment_id"'
		'type="enum<string>"'
		'Available options:'
		'type="string<datetime>"'
		'Always deployed for this test event.'
	)

	local failed=0
	for pat in "${patterns[@]}"; do
		if ! grep -qF "${pat}" "${synth_file}"; then
			echo "FAIL: ${name} - pattern not found: ${pat}"
			failed=1
		fi
	done

	# Clean up: remove test-synthetic.mdx and temp files
	rm -f "${synth_file}"
	rm -f "${tmp_spec}" "${tmp_config}"

	if [[ "${failed}" -ne 0 ]]; then
		exit 1
	fi

	echo "PASS: ${name}"
}

test_error_missing_config() {
	local name="test_error_missing_config"
	local output
	local rc=0

	output="$(python3 "${repo_root}/${generator}" "${repo_root}/${spec_file}" /nonexistent/config.yaml "${repo_root}" 2>&1)" || rc=$?

	if [[ "${rc}" -ne 1 ]]; then
		echo "FAIL: ${name} - expected exit code 1, got ${rc}"
		exit 1
	fi
	if [[ "${output}" != *"Config file not found"* ]]; then
		echo "FAIL: ${name} - output missing 'Config file not found': ${output}"
		exit 1
	fi

	echo "PASS: ${name}"
}

test_error_missing_spec() {
	local name="test_error_missing_spec"
	local output
	local rc=0

	output="$(python3 "${repo_root}/${generator}" /nonexistent/spec.yaml "${repo_root}/${config_file}" "${repo_root}" 2>&1)" || rc=$?

	if [[ "${rc}" -ne 1 ]]; then
		echo "FAIL: ${name} - expected exit code 1, got ${rc}"
		exit 1
	fi
	if [[ "${output}" != *"OpenAPI spec not found"* ]]; then
		echo "FAIL: ${name} - output missing 'OpenAPI spec not found': ${output}"
		exit 1
	fi

	echo "PASS: ${name}"
}

test_error_missing_event_envelope() {
	local name="test_error_missing_event_envelope"
	local tmp_spec
	tmp_spec="$(mktemp --suffix=.yaml)"

	# Create a minimal spec with a data event schema but no Event envelope
	cat > "${tmp_spec}" <<-'YAML'
	openapi: 3.0.3
	info:
	  title: Test
	  version: v0.1
	paths: {}
	components:
	  schemas:
	    FakeDataEvent:
	      type: object
	      description: 'Payload for `fake.event` events.'
	      required: [id]
	      properties:
	        id:
	          type: string
	      example:
	        id: fake_123
	YAML

	local output
	local rc=0
	output="$(python3 "${repo_root}/${generator}" "${tmp_spec}" "${repo_root}/${config_file}" "${repo_root}" 2>&1)" || rc=$?

	rm -f "${tmp_spec}"

	if [[ "${rc}" -ne 1 ]]; then
		echo "FAIL: ${name} - expected exit code 1, got ${rc}"
		exit 1
	fi
	if [[ "${output}" != *"Event envelope schema not found"* ]]; then
		echo "FAIL: ${name} - output missing 'Event envelope schema not found': ${output}"
		exit 1
	fi

	echo "PASS: ${name}"
}

test_error_wrong_args() {
	local name="test_error_wrong_args"
	local output
	local rc=0

	output="$(python3 "${repo_root}/${generator}" 2>&1)" || rc=$?

	if [[ "${rc}" -ne 1 ]]; then
		echo "FAIL: ${name} - expected exit code 1, got ${rc}"
		exit 1
	fi
	if [[ "${output}" != *"Usage:"* ]]; then
		echo "FAIL: ${name} - output missing 'Usage:': ${output}"
		exit 1
	fi

	echo "PASS: ${name}"
}

main() {
	cd "${repo_root}"

	test_zero_diff
	test_idempotence
	test_extensibility
	test_error_missing_config
	test_error_missing_spec
	test_error_missing_event_envelope
	test_error_wrong_args

	echo ""
	echo "All 7 tests passed."
}

main

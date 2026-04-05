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
generator="api/generate_event_pages.py"

test_zero_diff() {
	local name="test_zero_diff"
	local tmp_dir
	tmp_dir="$(mktemp -d)"

	# Back up the two MDX files
	cp "${repo_root}/${events_dir}/deployment-deployed.mdx" "${tmp_dir}/deployment-deployed.mdx"
	cp "${repo_root}/${events_dir}/deployment-removed.mdx" "${tmp_dir}/deployment-removed.mdx"

	# Run the generator
	python3 "${repo_root}/${generator}" "${repo_root}/${config_file}" "${repo_root}" >/dev/null

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
	python3 "${repo_root}/${generator}" "${repo_root}/${config_file}" "${repo_root}" >/dev/null

	# Capture sha256sum of both MDX files
	local sum1_deployed sum1_removed
	sum1_deployed="$(sha256sum "${repo_root}/${events_dir}/deployment-deployed.mdx" | awk '{print $1}')"
	sum1_removed="$(sha256sum "${repo_root}/${events_dir}/deployment-removed.mdx" | awk '{print $1}')"

	# Run generator again
	python3 "${repo_root}/${generator}" "${repo_root}/${config_file}" "${repo_root}" >/dev/null

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
	local tmp_config
	tmp_config="$(mktemp --suffix=.yaml)"

	# Create a temp copy of the config with a synthetic event appended
	cp "${repo_root}/${config_file}" "${tmp_config}"
	cat >> "${tmp_config}" <<-'YAML'
	  - type: "test.synthetic"
	    schema: DeploymentDeployedEvent
	    slug: test-synthetic
	    description: "A synthetic event for testing."
	    body: "This is synthetic body prose."
	    field_annotations:
	      activity_status: 'Always `deployed` for this test event.'
	    example:
	      object: event
	      id: 99
	      type: test.synthetic
	      occurred_at: "2026-01-01T00:00:00Z"
	      data:
	        deployment_id: dpl_test
	        release_id: rls_test
	        status: deployed
	        activity_status: deployed
	        error_status: none
	        target_status: deployed
	        deployed_at: "2026-01-01T00:00:00Z"
	YAML

	# Run generator with temp config
	python3 "${repo_root}/${generator}" "${tmp_config}" "${repo_root}" >/dev/null

	local synth_file="${repo_root}/${events_dir}/test-synthetic.mdx"

	# Assert test-synthetic.mdx exists
	if [[ ! -f "${synth_file}" ]]; then
		echo "FAIL: ${name} - test-synthetic.mdx was not created"
		rm -f "${tmp_config}"
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
		'Always `deployed` for this test event.'
	)

	local failed=0
	for pat in "${patterns[@]}"; do
		if ! grep -qF "${pat}" "${synth_file}"; then
			echo "FAIL: ${name} - pattern not found: ${pat}"
			failed=1
		fi
	done

	# Clean up: remove test-synthetic.mdx and temp config
	rm -f "${synth_file}"
	rm -f "${tmp_config}"

	if [[ "${failed}" -ne 0 ]]; then
		exit 1
	fi

	echo "PASS: ${name}"
}

test_error_missing_config() {
	local name="test_error_missing_config"
	local output
	local rc=0

	output="$(python3 "${repo_root}/${generator}" /nonexistent/config.yaml "${repo_root}" 2>&1)" || rc=$?

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
	local tmp_config
	tmp_config="$(mktemp --suffix=.yaml)"

	cat > "${tmp_config}" <<-YAML
	spec: nonexistent/spec.yaml
	output_dir: ${events_dir}
	events: []
	YAML

	local output
	local rc=0
	output="$(python3 "${repo_root}/${generator}" "${tmp_config}" "${repo_root}" 2>&1)" || rc=$?

	rm -f "${tmp_config}"

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

test_error_missing_schema() {
	local name="test_error_missing_schema"
	local tmp_config
	tmp_config="$(mktemp --suffix=.yaml)"

	cat > "${tmp_config}" <<-YAML
	spec: docs/references/device-api/v0.2.1/api.yaml
	output_dir: ${events_dir}
	events:
	  - type: "fake.event"
	    schema: NonexistentSchema
	    slug: fake-event
	    description: "Fake event."
	    body: "Fake body."
	    field_annotations: {}
	    example:
	      object: event
	YAML

	local output
	local rc=0
	output="$(python3 "${repo_root}/${generator}" "${tmp_config}" "${repo_root}" 2>&1)" || rc=$?

	rm -f "${tmp_config}"

	if [[ "${rc}" -ne 1 ]]; then
		echo "FAIL: ${name} - expected exit code 1, got ${rc}"
		exit 1
	fi
	if [[ "${output}" != *"not found in spec"* ]]; then
		echo "FAIL: ${name} - output missing 'not found in spec': ${output}"
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
	test_error_missing_schema
	test_error_wrong_args

	echo ""
	echo "All 7 tests passed."
}

main

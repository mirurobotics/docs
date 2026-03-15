#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"
lint_script="${repo_root}/scripts/lint.sh"
fixture_root="${repo_root}/tests/lint-fixtures"

run_expect_pass() {
  local fixture_name="$1"
  local output_file

  output_file="$(mktemp)"
  if ! DOCS_LINT_ROOT="${fixture_root}/${fixture_name}" "${lint_script}" >"${output_file}" 2>&1; then
    cat "${output_file}" >&2
    rm -f "${output_file}"
    echo "Expected fixture ${fixture_name} to pass." >&2
    exit 1
  fi

  if ! grep -Fq "All documentation lint checks passed." "${output_file}"; then
    cat "${output_file}" >&2
    rm -f "${output_file}"
    echo "Fixture ${fixture_name} did not report a successful lint run." >&2
    exit 1
  fi

  rm -f "${output_file}"
}

run_expect_fail() {
  local fixture_name="$1"
  local expected_pattern="$2"
  local output_file

  output_file="$(mktemp)"
  if DOCS_LINT_ROOT="${fixture_root}/${fixture_name}" "${lint_script}" >"${output_file}" 2>&1; then
    cat "${output_file}" >&2
    rm -f "${output_file}"
    echo "Expected fixture ${fixture_name} to fail." >&2
    exit 1
  fi

  if ! grep -Fq "${expected_pattern}" "${output_file}"; then
    cat "${output_file}" >&2
    rm -f "${output_file}"
    echo "Fixture ${fixture_name} failed for an unexpected reason." >&2
    exit 1
  fi

  rm -f "${output_file}"
}

run_expect_pass "good"
run_expect_fail "bad-mdx" "Parsing error"
run_expect_fail "bad-spelling" "Unknown word"
run_expect_fail "bad-openapi" "Failed to validate OpenAPI schema"

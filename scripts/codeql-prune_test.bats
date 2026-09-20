#!/usr/bin/env bats

# Tests for codeql-prune.sh
#
# The GitHub CLI is replaced by a stub that serves a canned analysis listing
# and records every DELETE request in $GH_STUB_LOG.

SCRIPT="$BATS_TEST_DIRNAME/codeql-prune.sh"
KEEP=".github/workflows/codeql-analysis.yml:codeql"
STALE_CODEQL=".github/workflows/codeql-analysis.yml:codeql/language:javascript"
STALE_ANALYZE=".github/workflows/codeql-analysis.yml:analyze/language:javascript"

setup() {
	STUB_BIN="$BATS_TEST_TMPDIR/bin"
	mkdir -p "$STUB_BIN" "$BATS_TEST_TMPDIR/empty"
	export GH_STUB_LOG="$BATS_TEST_TMPDIR/gh.log"
	export GH_STUB_ARGS="$BATS_TEST_TMPDIR/gh.args"
	write_gh_stub "$STUB_BIN/gh"
	chmod +x "$STUB_BIN/gh"
}

write_gh_stub() {
	cat > "$1" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

analyses() {
	cat <<'JSON'
[
  {"id": 300, "category": ".github/workflows/codeql-analysis.yml:codeql",
   "created_at": "2026-09-17T07:31:00Z", "commit_sha": "b697e10aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
   "deletable": true, "tool": {"name": "CodeQL"}},
  {"id": 200, "category": ".github/workflows/codeql-analysis.yml:codeql/language:javascript",
   "created_at": "2026-03-23T20:56:00Z", "commit_sha": "93f4af4bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
   "deletable": true, "tool": {"name": "CodeQL"}},
  {"id": 100, "category": ".github/workflows/codeql-analysis.yml:analyze/language:javascript",
   "created_at": "2026-03-23T03:00:00Z", "commit_sha": "36590cacccccccccccccccccccccccccccccccc",
   "deletable": true, "tool": {"name": "CodeQL"}},
  {"id": 50, "category": ".github/workflows/codeql-analysis.yml:analyze/language:javascript",
   "created_at": "2026-03-22T03:00:00Z", "commit_sha": "1111111dddddddddddddddddddddddddddddddd",
   "deletable": true, "tool": {"name": "CodeQL"}}
]
JSON
}

delete_response() {
	if [[ "$1" == *"/analyses/100?"* ]]; then
		echo '{"next_analysis_url":"https://api.github.com/repos/mirurobotics/docs/code-scanning/analyses/50","confirm_delete_url":"https://api.github.com/repos/mirurobotics/docs/code-scanning/analyses/50?confirm_delete"}'
	else
		echo '{"next_analysis_url":null,"confirm_delete_url":null}'
	fi
}

jq_expr=""
args=()
while [[ $# -gt 0 ]]; do
	case "$1" in
		--jq) jq_expr="$2"; shift 2 ;;
		*) args+=("$1"); shift ;;
	esac
done

emit() {
	if [[ -n "$jq_expr" ]]; then
		jq -r "$jq_expr"
	else
		cat
	fi
}

case "${args[0]}" in
	repo)
		echo "mirurobotics/docs"
		;;
	api)
		if [[ -n "${GH_STUB_FAIL:-}" ]]; then
			echo "gh: HTTP 403: Resource not accessible by integration" >&2
			exit 1
		fi
		if [[ "${args[1]}" == "-X" && "${args[2]}" == "DELETE" ]]; then
			echo "${args[3]}" >> "$GH_STUB_LOG"
			delete_response "${args[3]}" | emit
		else
			printf '%s\n' "${args[@]}" >> "$GH_STUB_ARGS"
			analyses | emit
		fi
		;;
	*)
		echo "gh stub: unexpected args: $*" >&2
		exit 1
		;;
esac
EOF
}

@test "--help prints usage and exits 0" {
	run bash "$SCRIPT" --help
	[ "$status" -eq 0 ]
	[[ "$output" == *"Usage:"* ]]
}

@test "missing gh fails with a clear message" {
	run env PATH="$BATS_TEST_TMPDIR/empty" "$BASH" "$SCRIPT" --repo o/r
	[ "$status" -eq 1 ]
	[[ "$output" == *"gh is required"* ]]
}

@test "dry run lists categories and deletes nothing" {
	run env PATH="$STUB_BIN:$PATH" bash "$SCRIPT" --repo mirurobotics/docs
	[ "$status" -eq 0 ]
	[[ "$output" == *"KEEP $KEEP "* ]]
	[[ "$output" == *"PRUNE $STALE_CODEQL "* ]]
	[[ "$output" == *"PRUNE $STALE_ANALYZE "* ]]
	[[ "$output" == *"id=100"* ]]
	[[ "$output" != *"id=50"* ]]
	[[ "$output" == *"dry run: 2 categories would be pruned"* ]]
	[ ! -s "$GH_STUB_LOG" ]
}

@test "gh api failure aborts instead of reporting nothing to prune" {
	run env PATH="$STUB_BIN:$PATH" GH_STUB_FAIL=1 bash "$SCRIPT" --repo mirurobotics/docs
	[ "$status" -ne 0 ]
	[[ "$output" == *"HTTP 403"* ]]
	[[ "$output" != *"nothing to prune"* ]]
	[ ! -s "$GH_STUB_LOG" ]
}

@test "--delete follows confirm_delete_url" {
	run env PATH="$STUB_BIN:$PATH" bash "$SCRIPT" --repo mirurobotics/docs --delete
	[ "$status" -eq 0 ]
	[[ "$output" == *"deleted $STALE_CODEQL"* ]]
	[[ "$output" == *"deleted $STALE_ANALYZE"* ]]
	[ "$(wc -l < "$GH_STUB_LOG")" -eq 3 ]
	grep -q 'analyses/200?confirm_delete' "$GH_STUB_LOG"
	grep -q 'analyses/100?confirm_delete' "$GH_STUB_LOG"
	grep -q 'analyses/50?confirm_delete' "$GH_STUB_LOG"
	! grep -q 'analyses/300' "$GH_STUB_LOG"
}

@test "ref is passed as a query field, not interpolated into the URL" {
	run env PATH="$STUB_BIN:$PATH" bash "$SCRIPT" --repo mirurobotics/docs \
		--ref 'refs/heads/release#old'
	[ "$status" -eq 0 ]
	[[ "$output" == *"== mirurobotics/docs refs/heads/release#old =="* ]]
	grep -qx 'ref=refs/heads/release#old' "$GH_STUB_ARGS"
	grep -qx 'repos/mirurobotics/docs/code-scanning/analyses' "$GH_STUB_ARGS"
	! grep -q '?' "$GH_STUB_ARGS"
}

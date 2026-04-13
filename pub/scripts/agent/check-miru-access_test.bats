#!/usr/bin/env bats

# Tests for check-miru-access.sh
#
# All tests run as the current user with _CHECK_ACCESS_REEXEC=1
# to bypass the sudo re-exec path. No secondary user or sudo required.

SCRIPT="$BATS_TEST_DIRNAME/check-miru-access.sh"
CURRENT_USER="$(id -un)"

setup() {
	export _CHECK_ACCESS_REEXEC=1
	TMPDIR="$BATS_TEST_TMPDIR"
}

# Helper: run the script as the current user against a path.
run_check() {
	run bash "$SCRIPT" --user "$CURRENT_USER" "$@"
}

# ---------------------------------------------------------------------------
# Option parsing
# ---------------------------------------------------------------------------

@test "no arguments prints error and exits 2" {
	run bash "$SCRIPT"
	[ "$status" -eq 2 ]
	[[ "$output" == *"expected exactly one path argument"* ]]
}

@test "--user without value prints error and exits 2" {
	run bash "$SCRIPT" --user
	[ "$status" -eq 2 ]
	[[ "$output" == *"--user requires a value"* ]]
}

@test "unknown option prints error and exits 2" {
	run bash "$SCRIPT" --bogus /tmp
	[ "$status" -eq 2 ]
	[[ "$output" == *"unknown option"* ]]
}

@test "-h prints usage and exits 0" {
	run bash "$SCRIPT" -h
	[ "$status" -eq 0 ]
	[[ "$output" == *"Usage:"* ]]
}

@test "--help prints usage and exits 0" {
	run bash "$SCRIPT" --help
	[ "$status" -eq 0 ]
	[[ "$output" == *"Usage:"* ]]
}

@test "-- separates options from path" {
	run_check -- "$TMPDIR"
	[ "$status" -eq 0 ]
}

# ---------------------------------------------------------------------------
# Username validation
# ---------------------------------------------------------------------------

@test "--user root is rejected" {
	run bash "$SCRIPT" --user root /tmp
	[ "$status" -eq 2 ]
	[[ "$output" == *"refusing to check as root"* ]]
}

@test "--user with invalid characters is rejected" {
	run bash "$SCRIPT" --user 'bad user!' /tmp
	[ "$status" -eq 2 ]
	[[ "$output" == *"invalid username"* ]]
}

@test "--user with uppercase is rejected" {
	run bash "$SCRIPT" --user 'Root' /tmp
	[ "$status" -eq 2 ]
	[[ "$output" == *"invalid username"* ]]
}

@test "--user with nonexistent user is rejected" {
	run bash "$SCRIPT" --user zzz_no_such_user_zzz /tmp
	[ "$status" -eq 2 ]
	[[ "$output" == *"user does not exist"* ]]
}

# ---------------------------------------------------------------------------
# Re-exec sentinel
# ---------------------------------------------------------------------------

@test "reexec sentinel with matching user succeeds" {
	export _CHECK_ACCESS_REEXEC=1
	run_check "$TMPDIR"
	[ "$status" -eq 0 ]
}

@test "reexec sentinel with mismatched user errors" {
	export _CHECK_ACCESS_REEXEC=1
	# Pick a user that exists but is not us. 'nobody' exists on most systems.
	if id nobody >/dev/null 2>&1; then
		run bash "$SCRIPT" --user nobody /tmp
		[ "$status" -eq 2 ]
		[[ "$output" == *"re-exec identity mismatch"* ]]
	else
		skip "no 'nobody' user on this system"
	fi
}

# ---------------------------------------------------------------------------
# Happy path: file with full permissions
# ---------------------------------------------------------------------------

@test "readable writable file under accessible directory returns PASS" {
	local f="$TMPDIR/testfile"
	touch "$f"
	chmod 644 "$f"

	run_check "$f"
	[ "$status" -eq 0 ]
	[[ "$output" == *"FINAL RESULT: PASS"* ]]
	[[ "$output" == *"file"* ]]
}

# ---------------------------------------------------------------------------
# Happy path: accessible directory
# ---------------------------------------------------------------------------

@test "accessible directory returns PASS" {
	local d="$TMPDIR/testdir"
	mkdir -p "$d"
	chmod 755 "$d"

	run_check "$d"
	[ "$status" -eq 0 ]
	[[ "$output" == *"FINAL RESULT: PASS"* ]]
	[[ "$output" == *"directory"* ]]
}

# ---------------------------------------------------------------------------
# Happy path: non-existent file under writable parent
# ---------------------------------------------------------------------------

@test "non-existent file under writable parent returns PASS" {
	local f="$TMPDIR/does_not_exist"
	[ ! -e "$f" ]  # precondition

	run_check "$f"
	[ "$status" -eq 0 ]
	[[ "$output" == *"FINAL RESULT: PASS"* ]]
	[[ "$output" == *"parent"* ]]
	[[ "$output" == *"not visible"* ]]
}

# ---------------------------------------------------------------------------
# Permission failures: ok=0 propagation (the critical bug fix)
# ---------------------------------------------------------------------------

@test "file without read permission returns FAIL with NO" {
	local f="$TMPDIR/noread"
	touch "$f"
	chmod 000 "$f"

	run_check "$f"
	[ "$status" -eq 1 ]
	[[ "$output" == *"FINAL RESULT: FAIL"* ]]
	[[ "$output" == *"NO"* ]]
}

@test "directory without execute permission returns FAIL" {
	local d="$TMPDIR/noexec"
	mkdir -p "$d"
	chmod 600 "$d"

	run_check "$d"
	[ "$status" -eq 1 ]
	[[ "$output" == *"FINAL RESULT: FAIL"* ]]
	[[ "$output" == *"NO"* ]]
}

@test "file without write permission returns FAIL" {
	local f="$TMPDIR/nowrite"
	touch "$f"
	chmod 444 "$f"

	run_check "$f"
	[ "$status" -eq 1 ]
	[[ "$output" == *"FINAL RESULT: FAIL"* ]]
	[[ "$output" == *"NO"* ]]
}

# ---------------------------------------------------------------------------
# Missing paths report NO (not accessible)
# ---------------------------------------------------------------------------

@test "non-existent path under non-existent parent shows NO" {
	local f="$TMPDIR/no_such_parent/no_such_file"

	run_check "$f"
	[[ "$output" == *"NO"* ]]
	[[ "$output" == *"parent"* ]]
	[[ "$output" == *"FINAL RESULT: FAIL"* ]]
}

# ---------------------------------------------------------------------------
# Path canonicalization
# ---------------------------------------------------------------------------

@test "path with .. is canonicalized" {
	local d="$TMPDIR/a/b"
	mkdir -p "$d"

	# $TMPDIR/a/b/.. should resolve to $TMPDIR/a
	run_check "$TMPDIR/a/b/.."
	[ "$status" -eq 0 ]
	# Output should show the resolved path, not the .. form
	[[ "$output" != *".."* ]]
	[[ "$output" == *"$TMPDIR/a"* ]]
}

@test "path with trailing slash is canonicalized" {
	local d="$TMPDIR/slashdir"
	mkdir -p "$d"

	run_check "$d/"
	[ "$status" -eq 0 ]
	[[ "$output" == *"FINAL RESULT: PASS"* ]]
}

@test "relative path is canonicalized to absolute" {
	# Create a file in cwd-relative path
	local d="$TMPDIR/reltest"
	mkdir -p "$d"
	touch "$d/file"

	# Run from inside $d with a relative path
	run bash -c "cd '$d' && _CHECK_ACCESS_REEXEC=1 bash '$SCRIPT' --user '$CURRENT_USER' ./file"
	[ "$status" -eq 0 ]
	# Output should contain absolute path, not ./
	[[ "$output" != *"./file"* ]]
	[[ "$output" == *"$d/file"* ]]
}

# ---------------------------------------------------------------------------
# Duplicate root row suppression
# ---------------------------------------------------------------------------

@test "root-level path does not show / twice" {
	# Use /tmp which is directly under /
	run_check /tmp

	# Count how many times " / " or "| /" appears as a path in the table.
	# Each row ends with a path; / should appear exactly once.
	local root_rows
	root_rows=$(echo "$output" | grep -c '| /$' || true)
	# Also check for "| / " at end of line with possible trailing space
	root_rows2=$(echo "$output" | grep -cE '\| /\s*$' || true)
	local total=$((root_rows > root_rows2 ? root_rows : root_rows2))
	[ "$total" -le 1 ]
}

# ---------------------------------------------------------------------------
# sanitize_path: ANSI escape stripping
# ---------------------------------------------------------------------------

@test "ANSI escape sequences in path are stripped from output" {
	# Create a directory with an ANSI escape in its name
	local esc_dir="$TMPDIR/esc_$(printf '\033')[31mred"
	mkdir -p "$esc_dir" 2>/dev/null || skip "filesystem rejects control chars in names"

	run_check "$esc_dir"
	# The raw ESC byte (0x1b) should not appear in output
	if echo "$output" | LC_ALL=C grep -qP '\x1b'; then
		fail "output contains raw ESC byte"
	fi
}

# ---------------------------------------------------------------------------
# Table output format
# ---------------------------------------------------------------------------

@test "output contains table header" {
	local f="$TMPDIR/hdrfile"
	touch "$f"

	run_check "$f"
	[[ "$output" == *"Target Type"* ]]
	[[ "$output" == *"Read"* ]]
	[[ "$output" == *"Write"* ]]
	[[ "$output" == *"Execute"* ]]
	[[ "$output" == *"Path"* ]]
}

@test "ancestor directories get x-only checks (r and w are dashes)" {
	local d="$TMPDIR/deep/nest"
	mkdir -p "$d"

	run_check "$d"
	# Ancestor rows should have "-" for Read and Write columns
	local ancestor_line
	ancestor_line=$(echo "$output" | grep "ancestor" | head -1)
	[[ "$ancestor_line" == *"| -    | -     |"* ]]
}

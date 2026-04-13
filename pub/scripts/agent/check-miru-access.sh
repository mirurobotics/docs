#!/usr/bin/env bash
set -euo pipefail

# Capture script content at startup for TOCTOU-safe re-exec.
# In the re-exec child (bash -s), BASH_SOURCE[0] is empty — skip the capture
# since the child never re-execs again.
_SELF_CONTENT=""
if [[ -n "${BASH_SOURCE[0]:-}" ]] && [[ -r "${BASH_SOURCE[0]}" ]]; then
	_SELF_CONTENT="$(cat -- "${BASH_SOURCE[0]}")"
fi

usage() {
	cat <<'EOF'
Usage: check-miru-access.sh [--user <name>] <path>

Checks whether a user (default: `miru`) can access a target file or directory.

For file targets, it checks:
  - r,w on the file
  - r,w,x on the file's parent directory
  - x on the grandparent directory and above

If the file target does not exist yet, it skips file permission checks and
checks parent/ancestor directory permissions required to create it.

For directory targets, it checks:
  - r,w,x on the directory
  - x on parent directories above it

Options:
  -u, --user <name>  User to test as (default: miru)
  -h, --help         Show this help text
EOF
}

# Sanitize a path for display by stripping control characters (finding 6).
sanitize_path() {
	printf '%s' "$1" | LC_ALL=C tr -d '\000-\037\177'
}

target_user="miru"

while [[ $# -gt 0 ]]; do
	case "$1" in
		-u|--user)
			if [[ $# -lt 2 ]]; then
				echo "error: --user requires a value" >&2
				usage
				exit 2
			fi
			target_user="$2"
			shift
			shift
			;;
		-h|--help)
			usage
			exit 0
			;;
		--)
			shift
			break
			;;
		-*)
			echo "error: unknown option: $1" >&2
			usage
			exit 2
			;;
		*)
			break
			;;
	esac
done

if [[ $# -ne 1 ]]; then
	echo "error: expected exactly one path argument" >&2
	usage
	exit 2
fi

# Validate target_user: must match Unix username pattern, refuse root (finding 3).
if [[ ! "$target_user" =~ ^[a-z_][a-z0-9_-]*$ ]]; then
	echo "error: invalid username: $target_user" >&2
	exit 2
fi
if [[ "$target_user" == "root" ]]; then
	echo "error: refusing to check as root" >&2
	exit 2
fi
if ! id "$target_user" >/dev/null 2>&1; then
	echo "error: user does not exist: $target_user" >&2
	exit 2
fi
if [[ "$(id -u "$target_user")" -eq 0 ]]; then
	echo "error: refusing to check as UID 0 user: $target_user" >&2
	exit 2
fi

# Canonicalize the path early, before sudo re-exec (finding 4).
# realpath -m resolves without requiring the path to exist.
target_path="$(realpath -m "$1")"

# Detect re-exec via environment sentinel instead of CLI flag (finding 2).
if [[ "${_CHECK_ACCESS_REEXEC:-}" == "1" ]]; then
	# Validate that we are actually running as the target user.
	if [[ "$(id -un)" != "$target_user" ]]; then
		echo "error: re-exec identity mismatch: expected $target_user, got $(id -un)" >&2
		exit 2
	fi
else
	# Not yet running as target user — need to re-exec via sudo.
	if [[ "$(id -un)" == "$target_user" ]]; then
		# Already the target user; set sentinel and fall through.
		export _CHECK_ACCESS_REEXEC=1
	else
		if ! command -v sudo >/dev/null 2>&1; then
			echo "error: sudo is required" >&2
			exit 2
		fi

		caller_user="$(id -un)"
		parent_path="$(dirname "$target_path")"
		echo "Caller preflight ($caller_user):"
		if [ -e "$target_path" ]; then
			echo "  - target exists: $(sanitize_path "$target_path")"
		else
			echo "  - target does not exist (or is not visible): $(sanitize_path "$target_path")"
		fi
		if [ -d "$parent_path" ]; then
			echo "  - parent directory exists: $(sanitize_path "$parent_path")"
		else
			echo "  - parent directory does not exist (or is not visible): $(sanitize_path "$parent_path")"
		fi
		echo

		# Re-exec as target user via temp file to avoid both TOCTOU and
		# heredoc variable expansion mangling the script content.
		if [[ -z "$_SELF_CONTENT" ]]; then
			echo "error: cannot re-exec — script content not captured (was it piped to bash?)" >&2
			exit 2
		fi
		_tmpscript="$(mktemp)"
		trap 'rm -f "$_tmpscript"' EXIT
		printf '%s\n' "$_SELF_CONTENT" > "$_tmpscript"
		chmod 644 "$_tmpscript"
		exec sudo -u "$target_user" \
			_CHECK_ACCESS_REEXEC=1 \
			bash -s -- --user "$target_user" "$target_path" < "$_tmpscript"
	fi
fi

path="$target_path"
ok=1

print_header() {
	echo "Permission checks for: $(sanitize_path "$path")"
	if [ "${missing_target_file:-0}" -eq 1 ]; then
		echo "Note: target file is not visible to this user (missing or not traversable); skipping file checks."
	fi
	if [ "${parent_dir_uncertain:-0}" -eq 1 ]; then
		echo "Note: parent directory could not be stat'ed directly (missing or not traversable); reporting permission checks anyway."
	fi
	echo
	printf "| %-11s | %-4s | %-5s | %-7s | %s\n" "Target Type" "Read" "Write" "Execute" "Path"
	printf "| %-11s | %-4s | %-5s | %-7s | %s\n" "-----------" "----" "-----" "-------" "----"
}

# Evaluate a single permission cell and print the result.
# Returns the cell text on stdout; does NOT set ok (finding 1 — handled in print_row).
perm_cell() {
	local need="$1"
	local flag="$2"
	local target="$3"
	if [[ "$need" -eq 0 ]]; then
		printf -- "-"
		return
	fi

	# Distinguish missing paths from permission denials (finding 8).
	if ! [ -e "$target" ]; then
		printf -- "N/A"
		return
	fi

	if test "$flag" "$target"; then
		printf -- "OK"
	else
		printf -- "NO"
	fi
}

print_row() {
	local target_type="$1"
	local target="$2"
	local need_r="$3"
	local need_w="$4"
	local need_x="$5"

	local read_cell write_cell exec_cell
	read_cell="$(perm_cell "$need_r" -r "$target")"
	write_cell="$(perm_cell "$need_w" -w "$target")"
	exec_cell="$(perm_cell "$need_x" -x "$target")"

	# Propagate failures in the parent shell (finding 1).
	if [[ "$read_cell" == "NO" || "$write_cell" == "NO" || "$exec_cell" == "NO" ]]; then
		ok=0
	fi

	printf "| %-11s | %-4s | %-5s | %-7s | %s\n" \
		"$target_type" "$read_cell" "$write_cell" "$exec_cell" "$(sanitize_path "$target")"
}

# Track the last directory printed so we can skip duplicates (finding 7).
last_printed_dir=""

if [ -f "$path" ]; then
	missing_target_file=0
	print_header
	print_row "file" "$path" 1 1 0
	parent="$(dirname "$path")"
	print_row "parent" "$parent" 1 1 1
	last_printed_dir="$parent"
	dir="$(dirname "$parent")"
elif [ -d "$path" ]; then
	missing_target_file=0
	print_header
	print_row "directory" "$path" 1 1 1
	last_printed_dir="$path"
	dir="$(dirname "$path")"
elif [ ! -e "$path" ]; then
	parent="$(dirname "$path")"
	# Either the file is missing or this user cannot traverse to it.
	# In both cases, continue with directory permission checks.
	missing_target_file=1
	parent_dir_uncertain=0
	if [ -e "$parent" ] && [ ! -d "$parent" ]; then
		echo "path exists but parent is not a directory: $(sanitize_path "$parent")"
		exit 2
	elif [ ! -d "$parent" ]; then
		parent_dir_uncertain=1
	fi
	print_header
	print_row "parent" "$parent" 1 1 1
	last_printed_dir="$parent"
	dir="$(dirname "$parent")"
else
	echo "path exists but is neither a regular file nor a directory: $(sanitize_path "$path")"
	exit 2
fi

# Walk ancestor directories, skipping duplicates (finding 7).
while :; do
	if [[ "$dir" != "$last_printed_dir" ]]; then
		print_row "ancestor" "$dir" 0 0 1
	fi
	[ "$dir" = "/" ] && break
	dir="$(dirname "$dir")"
done

if [ "$ok" -eq 1 ]; then
	echo
	echo "FINAL RESULT: PASS"
	exit 0
fi

echo
echo "FINAL RESULT: FAIL"
exit 1

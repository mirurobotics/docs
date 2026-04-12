#!/usr/bin/env bash
set -euo pipefail

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

as_target_user=0
target_user="miru"

while [[ $# -gt 0 ]]; do
	case "$1" in
		--as-user)
			as_target_user=1
			shift
			;;
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

target_path="$1"

if ! id "$target_user" >/dev/null 2>&1; then
	echo "error: user does not exist: $target_user" >&2
	exit 2
fi

if [[ "$(id -un)" == "$target_user" ]]; then
	as_target_user=1
fi

if [[ "$as_target_user" -eq 0 ]]; then
	if ! command -v sudo >/dev/null 2>&1; then
		echo "error: sudo is required" >&2
		exit 2
	fi

	script_path="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)/$(basename -- "${BASH_SOURCE[0]}")"
	# Pass the script content over stdin so the target user does not need access
	# to this script's location (e.g. a protected home directory path).
	exec sudo -u "$target_user" bash -s -- --as-user --user "$target_user" "$target_path" < "$script_path"
fi

path="$target_path"
ok=1

print_header() {
	echo "Permission checks for: $path"
	if [ "${missing_target_file:-0}" -eq 1 ]; then
		echo "Note: target file does not exist yet; skipping file checks."
	fi
	echo
	printf "| %-11s | %-4s | %-5s | %-7s | %s\n" "Target Type" "Read" "Write" "Execute" "Path"
	printf "| %-11s | %-4s | %-5s | %-7s | %s\n" "-----------" "----" "-----" "-------" "----"
}

perm_cell() {
	need="$1"
	flag="$2"
	target="$3"
	if [[ "$need" -eq 0 ]]; then
		printf -- "-"
		return
	fi

	if test "$flag" "$target"; then
		printf -- "OK"
	else
		printf -- "NO"
		ok=0
	fi
}

print_row() {
	target_type="$1"
	target="$2"
	need_r="$3"
	need_w="$4"
	need_x="$5"

	read_cell="$(perm_cell "$need_r" -r "$target")"
	write_cell="$(perm_cell "$need_w" -w "$target")"
	exec_cell="$(perm_cell "$need_x" -x "$target")"
	printf "| %-11s | %-4s | %-5s | %-7s | %s\n" "$target_type" "$read_cell" "$write_cell" "$exec_cell" "$target"
}

if [ -f "$path" ]; then
	missing_target_file=0
	print_header
	print_row "file" "$path" 1 1 0
	parent="$(dirname "$path")"
	print_row "parent" "$parent" 1 1 1
	dir="$(dirname "$parent")"
elif [ -d "$path" ]; then
	missing_target_file=0
	print_header
	print_row "directory" "$path" 1 1 1
	dir="$(dirname "$path")"
elif [ ! -e "$path" ]; then
	parent="$(dirname "$path")"
	if [ -d "$parent" ]; then
		# Missing target file is valid; check permissions needed to create it.
		missing_target_file=1
		print_header
		print_row "parent" "$parent" 1 1 1
		dir="$(dirname "$parent")"
	else
		echo "directory does not exist: $parent"
		exit 2
	fi
else
	echo "path exists but is neither a regular file nor a directory: $path"
	exit 2
fi

while :; do
	print_row "ancestor" "$dir" 0 0 1
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

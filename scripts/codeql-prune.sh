#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<EOF
Usage: codeql-prune.sh [--repo <owner/name>] [--ref <ref>]... [--keep <category>]... [--delete]

Lists CodeQL code-scanning categories for a ref and prunes the ones not in
the keep set. Dry run by default; --delete removes them.

Options:
  --repo <owner/name>  Repository (default: gh repo view of the current checkout)
  --ref <ref>          Ref to inspect, repeatable (default: refs/heads/main)
  --keep <category>    Category to keep, repeatable (default:
                       .github/workflows/codeql-analysis.yml:codeql)
  --delete             Delete every analysis of each non-kept category
  -h, --help           Show this help

Requires gh (authenticated with a token that has code scanning write
access, e.g. GH_TOKEN) and jq.

Exit codes: 0 success, 1 error, 2 usage error.
EOF
}

die() {
	local code="$1"
	shift
	echo "codeql-prune.sh: $*" >&2
	if [[ "${code}" -eq 2 ]]; then
		usage >&2
	fi
	exit "${code}"
}

repo=""
refs=()
keeps=()
delete=false

while [[ $# -gt 0 ]]; do
	case "$1" in
		--repo)
			[[ $# -ge 2 ]] || die 2 "--repo requires a value"
			repo="$2"
			shift 2
			;;
		--ref)
			[[ $# -ge 2 ]] || die 2 "--ref requires a value"
			refs+=("$2")
			shift 2
			;;
		--keep)
			[[ $# -ge 2 ]] || die 2 "--keep requires a value"
			keeps+=("$2")
			shift 2
			;;
		--delete)
			delete=true
			shift
			;;
		-h | --help)
			usage
			exit 0
			;;
		*)
			die 2 "unknown option: $1"
			;;
	esac
done

[[ ${#refs[@]} -gt 0 ]] || refs=("refs/heads/main")
[[ ${#keeps[@]} -gt 0 ]] || keeps=(".github/workflows/codeql-analysis.yml:codeql")

for tool in gh jq; do
	command -v "${tool}" >/dev/null 2>&1 || die 1 "${tool} is required by codeql-prune.sh"
done

if [[ -z "${repo}" ]]; then
	repo="$(gh repo view --json nameWithOwner --jq .nameWithOwner)"
fi

# Prints one TSV line per CodeQL category on the ref, describing its newest
# analysis: category, created_at, short sha, id, deletable.
latest_per_category() {
	local repo="$1" ref="$2"
	gh api --paginate "repos/${repo}/code-scanning/analyses?ref=${ref}&per_page=100" |
		jq -s 'add // []' |
		jq -r '[.[] | select(.tool.name == "CodeQL")]
			| group_by(.category) | map(max_by(.created_at)) | .[]
			| [.category, .created_at, .commit_sha[0:7], .id, .deletable] | @tsv'
}

is_kept() {
	local category="$1" keep
	for keep in "${keeps[@]}"; do
		[[ "${keep}" == "${category}" ]] && return 0
	done
	return 1
}

# Deletes the newest analysis of a category and follows confirm_delete_url
# until the chain is exhausted.
prune_category() {
	local repo="$1" id="$2"
	local url="repos/${repo}/code-scanning/analyses/${id}?confirm_delete"
	while [[ -n "${url}" ]]; do
		url="$(gh api -X DELETE "${url}" --jq '.confirm_delete_url // empty')"
	done
}

process_ref() {
	local repo="$1" ref="$2"
	local category created sha id deletable action count=0
	echo "== ${repo} ${ref} =="
	while IFS=$'\t' read -r category created sha id deletable; do
		action=PRUNE
		if is_kept "${category}"; then
			action=KEEP
		elif [[ "${deletable}" != "true" ]]; then
			action=SKIP
		fi
		echo "${action} ${category} ${created} ${sha} id=${id}"
		[[ "${action}" == "PRUNE" ]] || continue
		if [[ "${delete}" == "true" ]]; then
			prune_category "${repo}" "${id}"
			echo "deleted ${category}"
		fi
		count=$((count + 1))
	done < <(latest_per_category "${repo}" "${ref}")
	report "${count}"
}

report() {
	local count="$1"
	if [[ "${delete}" == "true" ]]; then
		echo "pruned ${count} categories"
	elif [[ "${count}" -eq 0 ]]; then
		echo "nothing to prune"
	else
		echo "dry run: ${count} categories would be pruned; re-run with --delete"
	fi
}

main() {
	local ref
	for ref in "${refs[@]}"; do
		process_ref "${repo}" "${ref}"
	done
}

main

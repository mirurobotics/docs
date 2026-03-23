#!/usr/bin/env bash
set -euo pipefail

# Resolve the repo root relative to this script so it works from any cwd.
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"
content_root="${DOCS_LINT_ROOT:-${repo_root}}"
cspell_config="${DOCS_CSPELL_CONFIG:-${repo_root}/cspell.json}"

if ! command -v pnpm >/dev/null 2>&1; then
  echo "pnpm is required to run docs linting. Install pnpm and rerun ./scripts/lint.sh." >&2
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "go is required for MDX prose linting. Install Go and rerun ./scripts/lint.sh." >&2
  exit 1
fi

cd "${repo_root}"

collect_files() {
  local search_root="$1"
  local pattern="$2"
  if [[ -d "${search_root}" ]]; then
    find "${search_root}" -type f -name "${pattern}" -print0
  fi
}

mdx_targets=()
while IFS= read -r -d '' file; do
  mdx_targets+=("${file}")
done < <(
  {
    collect_files "${content_root}/docs" "*.mdx"
    collect_files "${content_root}/snippets" "*.mdx"
  }
)

spell_targets=()
if [[ -f "${content_root}/rclone.md" ]]; then
  spell_targets+=("${content_root}/rclone.md")
fi
if [[ -f "${content_root}/README.md" ]]; then
  spell_targets+=("${content_root}/README.md")
fi
spell_targets+=("${mdx_targets[@]}")

openapi_targets=()
while IFS= read -r -d '' file; do
  openapi_targets+=("${file}")
done < <(collect_files "${content_root}/docs/references" "*.yaml")

if [[ ${#mdx_targets[@]} -eq 0 ]]; then
  echo "No MDX files found under ${content_root}." >&2
  exit 1
fi

if [[ ${#openapi_targets[@]} -eq 0 ]]; then
  echo "No OpenAPI specs found under ${content_root}/docs/references." >&2
  exit 1
fi

echo "== MDX Prose =="
(cd "${repo_root}/tools/lint" && go build -o lint .)
"${repo_root}/tools/lint/lint" "${mdx_targets[@]}"

echo "== ESLint (MDX) =="
pnpm exec eslint --max-warnings=0 "${mdx_targets[@]}"

echo "== CSpell =="
pnpm exec cspell lint --no-progress --config "${cspell_config}" "${spell_targets[@]}"

echo "== OpenAPI =="
for spec in "${openapi_targets[@]}"; do
  echo "Checking ${spec#${repo_root}/}"
  pnpm exec mint openapi-check "${spec#${repo_root}/}"
done

echo "All documentation lint checks passed."

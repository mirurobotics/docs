
### Agent Updates
- Supported Versions
- Device API Version Matrix
- Changelog

### Platform API Updates
- Supported Versions
- Platform API Reference
- SDK Compatibility Matrix
- "latest" API component
- Changelog

### Device API Updates
- Supported Versions
- Device API Reference
- SDK Compatibility Matrix
- "latest" API component
- Changelog

## Dependency audit

CI runs `./scripts/audit.sh` (`pnpm audit`). Every advisory so far has come from transitive dependencies of the `mint` CLI's local preview and build tooling, not from anything published to readers. Fix them with `pnpm.overrides` in `package.json`. Ignore an advisory in `pnpm.auditConfig.ignoreCves` only when no patched release exists, and record it below.

### Ignored advisories

| CVE | Package | Path | Why ignored | Remove when |
|-----|---------|------|-------------|-------------|
| CVE-2026-93748 ([GHSA-ch52-4w7c-c8xp](https://github.com/advisories/GHSA-ch52-4w7c-c8xp)) | `http-cache-semantics` | `mint > @mintlify/previewing > got > cacheable-request` | No patched release | A fixed version ships; override to it |
| CVE-2026-93687 ([GHSA-vfj7-8cjw-p6xm](https://github.com/advisories/GHSA-vfj7-8cjw-p6xm)) | `braces` | `mint > @mintlify/common > tailwindcss > chokidar` | No patched release | A fixed version ships; override to it |
| CVE-2026-97058 ([GHSA-hp3w-g68c-fv3c](https://github.com/advisories/GHSA-hp3w-g68c-fv3c)) | `sprintf-js` | `mint > @mintlify/common > front-matter > js-yaml@3 > argparse@1` | No patched release | A fixed version ships, or `front-matter` moves off `js-yaml@3` |

### Overrides past a dependent's range

These overrides force a newer major than the packages using them ask for, because no patched release exists in the requested range. Only local preview uses them, and no page renders math.

| Override | Dependents ask for | Remove when |
|----------|--------------------|-------------|
| `postcss-selector-parser` `>=7.1.6 <8` | `^6` (tailwindcss 3, postcss-nested) | Mintlify moves to a tailwindcss that accepts 7.x |
| `katex` `>=0.18.2 <0.20` | `^0.16` (rehype-katex, micromark-extension-math) | Those packages accept a patched katex |

To check whether an ignore still applies, temporarily remove it and run `./scripts/audit.sh`.

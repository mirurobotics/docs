
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

CI runs `./scripts/audit.sh` (`pnpm audit`). It covers only the packages in `pnpm-lock.yaml`: cspell, eslint, eslint-plugin-mdx and their dependencies, plus `@mintlify/cli` and the few packages it installs, such as `sharp`. Since mint 4.3, `@mintlify/cli` bundles the validation and preview tooling that every earlier advisory came from, or downloads it at runtime, so `pnpm audit` can't see it and `pnpm.overrides` can't patch it. None of it is published to readers. Fix advisories that `pnpm audit` reports with `pnpm.overrides` in `package.json`. Ignore an advisory in `pnpm.auditConfig.ignoreCves` (or `ignoreGhsas` when it has no CVE) only when no patched release exists that its dependents can use, and record it below.

### Ignored advisories

None.

### Overrides past a dependent's range

These overrides force a version outside the range the packages using them ask for, because no patched release exists in the requested range. Only local tooling uses them.

| Override | Dependents ask for | Remove when |
|----------|--------------------|-------------|
| `sharp` `>=0.35.5` | `0.33.5` (@mintlify/cli) | `@mintlify/cli` accepts a patched sharp |
| `minimatch` `>=11.0.1` | `^10` (eslint, @eslint/config-array), `^9` (@npmcli/map-workspaces via eslint-plugin-mdx) | Those packages accept minimatch 11 |
| `glob` `>=14.0.2` | `^10` (unified-engine and @npmcli packages via eslint-plugin-mdx) | Those packages accept glob 14 |

To check whether an ignore still applies, temporarily remove it and run `./scripts/audit.sh`.

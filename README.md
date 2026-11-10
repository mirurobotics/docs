
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

CI runs `./scripts/audit.sh` (`pnpm audit`). It covers the packages in `pnpm-lock.yaml`: `cspell`, the ESLint and MDX lint tooling, and the `mint` CLI's install-time dependencies. Since `mint` 4.3, the CLI bundles or fetches its preview and validation tooling at runtime, so the audit and `pnpm.overrides` no longer reach it. Fix advisories with `pnpm.overrides` in `package.json`. Ignore an advisory in `pnpm.auditConfig.ignoreCves` (or `ignoreGhsas` when it has no CVE) only when no patched release exists that its dependents can use, and record it below.

### Ignored advisories

None.

### Overrides past a dependent's range

These overrides force a newer major than the packages using them ask for, because no patched release exists in the requested range. Only linting uses them.

| Override | Dependents ask for | Remove when |
|----------|--------------------|-------------|
| `minimatch` `>=11.0.1` | `^10` (eslint, @eslint/config-array) | Those packages accept minimatch 11 |
| `glob` `>=14.0.2` | `^10` (unified-engine via eslint-plugin-mdx, @npmcli packages) | Those packages accept glob 14 |

To check whether an ignore still applies, temporarily remove it and run `./scripts/audit.sh`.

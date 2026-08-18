/*
 * Mintlify only evaluates the exported bindings in a snippet file, so module-level
 * helpers and constants are not in scope at render time. Each component below must
 * build its own href — when the latest API version changes, bump it in each.
 */

/**
 * Link to a Platform API endpoint using the latest API version.
 *
 * @param {string} endpoint - Path after /endpoints/, e.g. "deployments/get" or "devices/list#parameter-id"
 * @param {boolean} [newTab] - Open link in a new tab
 * @param {React.ReactNode} children - Link text
 */
export const PlatformApiLink = ({ endpoint, newTab, children }) => {
  const href = `/references/platform-api/2026-08-17/endpoints/${endpoint}`;
  if (newTab) {
    return <a href={href} target="_blank" rel="noopener noreferrer">{children}</a>;
  }
  return <a href={href}>{children}</a>;
};

/**
 * Clickable badge marking an operation as available in the Platform API.
 * Place on its own line directly beneath the section heading. Opens the endpoint
 * reference in a new tab so the reader keeps their place in the dashboard steps.
 *
 * Where a section covers several endpoints, link the primary one and keep the
 * generic label.
 *
 * Mintlify tags bare MDX anchors with `class="link"`, which underlines them with
 * a border-bottom rather than text-decoration — both have to be cleared here.
 *
 * @param {string} endpoint - Path after /endpoints/, e.g. "groups/create"
 * @param {string} [label] - Badge text
 */
export const PlatformApiBadge = ({ endpoint, label = "Platform API" }) => {
  const href = `/references/platform-api/2026-08-17/endpoints/${endpoint}`;
  return (
    <Tooltip tip="Open the Platform API reference for this operation">
      <a
        href={href}
        target="_blank"
        rel="noopener noreferrer"
        style={{ textDecoration: 'none', borderBottom: 'none' }}
      >
        <span
          className="inline-flex items-center gap-1.5 rounded-lg px-2 py-0.5 text-sm text-[rgb(var(--gray-500))] dark:text-[rgb(var(--gray-400))] hover:bg-[rgb(var(--gray-100))] dark:hover:bg-[rgb(var(--gray-800))] hover:text-[rgb(var(--gray-700))] dark:hover:text-[rgb(var(--gray-200))]"
          style={{ marginLeft: '-0.5rem' }}
        >
          {label} <Icon icon="arrow-up-right" size={13} color="currentColor" />
        </span>
      </a>
    </Tooltip>
  );
};

/**
 * Inert counterpart to PlatformApiBadge for operations with no Platform API
 * endpoint. Same ghost styling and placement, but no link, arrow, or hover
 * background — it must not read as clickable. The slash icon plus hover tip
 * carry the "not supported" signal.
 *
 * When an endpoint ships for a marked operation, replace this with a
 * PlatformApiBadge — grep the docs for PlatformUnsupportedBadge when new endpoints
 * are released.
 */
export const PlatformUnsupportedBadge = () => {
  return (
    <Tooltip tip="Platform API does not support this operation">
      <span
        className="inline-flex items-center gap-1.5 rounded-lg px-2 py-0.5 text-sm text-[rgb(var(--gray-500))] dark:text-[rgb(var(--gray-400))]"
        style={{ marginLeft: '-0.5rem', cursor: 'default' }}
      >
        <Icon icon="ban" size={13} color="currentColor" /> Platform API
      </span>
    </Tooltip>
  );
};

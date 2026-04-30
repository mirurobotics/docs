/**
 * Link to a Platform API endpoint using the latest API version.
 *
 * @param {string} endpoint - Path after /endpoints/, e.g. "deployments/get" or "devices/list#parameter-id"
 * @param {boolean} [newTab] - Open link in a new tab
 * @param {React.ReactNode} children - Link text
 */
export const PlatformApiLink = ({ endpoint, newTab, children }) => {
  const href = `/references/platform-api/2026-03-09/endpoints/${endpoint}`;
  if (newTab) {
    return <a href={href} target="_blank" rel="noopener noreferrer">{children}</a>;
  }
  return <a href={href}>{children}</a>;
};

/**
 * A link that opens in a new tab.
 * Use for internal or external links where the user should stay on the current page.
 *
 * @param {string} href - Link destination (can be relative or absolute)
 * @param {React.ReactNode} children - Link text or content
 */
export const LinkNewTab = ({ href, children }) => (
  <a href={href} target="_blank" rel="noopener noreferrer">
    {children}
  </a>
);

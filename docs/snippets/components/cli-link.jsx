/*
 * Mintlify only evaluates the exported bindings in a snippet file, so module-level
 * helpers and constants are not in scope at render time. Each component below must
 * build its own href.
 */

/**
 * Clickable badge marking an operation as available in the CLI.
 * Place on its own line directly beneath the section heading. Opens the command
 * reference in a new tab so the reader keeps their place in the dashboard steps.
 *
 * Mintlify tags bare MDX anchors with `class="link"`, which underlines them with
 * a border-bottom rather than text-decoration — both have to be cleared here.
 *
 * @param {string} command - Path after /references/cli/, e.g. "device-stage"
 * @param {string} [label] - Badge text
 */
export const CliBadge = ({ command, label = "CLI" }) => {
  const href = `/references/cli/${command}`;
  return (
    <Tooltip tip="Open the CLI reference for this command">
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

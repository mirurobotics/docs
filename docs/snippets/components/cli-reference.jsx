/*
 * Mintlify only evaluates the exported bindings in a snippet file, so module-level
 * helpers and constants are not in scope at render time. Titles and section order
 * live inside CliReference so a page can omit a slot and the layout still holds.
 */

/**
 * Shared layout for CLI command reference pages.
 *
 * Pass intro prose as children. Optional slots: requirements, scopes, usage,
 * args, flags, examples. Args and flags nest under Usage as bold labels.
 * Omit a slot to skip that section.
 */
export const CliReference = ({
  children,
  requirements,
  scopes,
  usage,
  args,
  flags,
  examples,
}) => {
  const section = (title, content) => {
    if (content == null) {
      return null;
    }
    return (
      <>
        <Heading level={3}>{title}</Heading>
        {content}
      </>
    );
  };

  const subsection = (title, content) => {
    if (content == null) {
      return null;
    }
    return (
      <>
        <p>
          <strong>{title}</strong>
        </p>
        {content}
      </>
    );
  };

  const usageBody =
    usage != null || args != null || flags != null ? (
      <>
        {usage}
        {subsection("Arguments", args)}
        {subsection("Flags", flags)}
      </>
    ) : null;

  return (
    <>
      {children}
      {section("Requirements", requirements)}
      {section("API key scopes", scopes)}
      {section("Usage", usageBody)}
      {section("Examples", examples)}
    </>
  );
};

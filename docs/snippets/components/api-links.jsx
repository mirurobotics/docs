/**
 * Link to a Device API endpoint on the latest reference.
 * The /references/device-api/latest redirect in docs.json is what moves when a version ships.
 *
 * @param {string} endpoint - Path after /endpoints/, e.g. "calibrations/create"
 * @param {boolean} [newTab] - Open link in a new tab
 * @param {React.ReactNode} children - Link text
 */
export const DeviceApiLink = ({ endpoint, newTab, children }) => {
  const href = `/references/device-api/latest/endpoints/${endpoint}`;
  if (newTab) {
    return <a href={href} target="_blank" rel="noopener noreferrer">{children}</a>;
  }
  return <a href={href}>{children}</a>;
};

export const DeviceApiReleaseLinks = ({ version }) => {
  return (
    <ul>
      <li>
        API reference:{' '}
        <a href={`/references/device-api/${version}`} target="_blank" rel="noopener noreferrer">
          {version}
        </a>
      </li>
      <li>
        OpenAPI spec:{' '}
        <a
          href={`https://assets.mirurobotics.com/docs/openapi/device/${version}.yaml`}
          target="_blank"
          rel="noopener noreferrer"
        >
          Download YAML
        </a>
      </li>
    </ul>
  );
};

export const PlatformApiReleaseLinks = ({ version }) => {
  const linkVersion = version.split('.')[0];

  return (
    <ul>
      <li>
        API reference:{' '}
        <a 
          href={`/references/platform-api/${linkVersion}`}
          target="_blank" rel="noopener noreferrer">
          {version}
        </a>
      </li>
      <li>
        OpenAPI spec:{' '}
        <a
          href={`https://assets.mirurobotics.com/docs/openapi/platform/${linkVersion}.yaml`}
          target="_blank"
          rel="noopener noreferrer"
        >
          Download YAML
        </a>
      </li>
    </ul>
  );
};

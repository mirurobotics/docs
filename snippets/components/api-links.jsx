export const DeviceApiReleaseLinks = ({ version }) => {
  return (
    <ul>
      <li>
        API reference:{' '}
        <a href={`/docs/references/device-api/${version}`} target="_blank" rel="noopener noreferrer">
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
          href={`/docs/references/platform-api/${linkVersion}`}
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

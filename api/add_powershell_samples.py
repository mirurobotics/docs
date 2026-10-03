#!/usr/bin/env python3
"""
Add Windows PowerShell examples to device API endpoints.

On Windows the agent serves the device API over loopback TCP and requires a
bearer token. Both the port and the token are read from the agent's discovery
file, so each example reads that file first and then calls the endpoint.
Streaming endpoints use curl.exe so events print as they arrive.
"""

import sys
from pathlib import Path
from urllib.parse import urlparse

import yaml

LANG = 'powershell'
LABEL = 'PowerShell (Windows)'

READ_DISCOVERY_FILE = (
    '$api = Get-Content -Raw "$env:ProgramData\\Miru\\device-api\\device-api.json"'
    ' | ConvertFrom-Json; `'
)
METHODS = ['get', 'post', 'put', 'patch', 'delete', 'head', 'options']


def is_streaming_operation(operation):
    """Return True if the operation produces a text/event-stream response."""
    for resp in operation.get('responses', {}).values():
        content = resp.get('content', {}) if isinstance(resp, dict) else {}
        if 'text/event-stream' in content:
            return True
    return False


def build_powershell_command(path, method, version_path, *, streaming=False):
    """Build a paste-safe PowerShell example that reads the discovery file first."""
    url = f'"http://127.0.0.1:$($api.port){version_path}{path}"'
    if streaming:
        lines = [
            READ_DISCOVERY_FILE,
            'curl.exe --no-buffer `',
            f'  --request {method.upper()} `',
            f'  --url {url} `',
            '  --header "Authorization: Bearer $($api.token)"',
        ]
    else:
        lines = [
            READ_DISCOVERY_FILE,
            f'Invoke-RestMethod -Method {method.capitalize()} `',
            f'  -Uri {url} `',
            '  -Headers @{ Authorization = "Bearer $($api.token)" }',
        ]
    return '\n'.join(lines)


def insert_index(code_samples):
    """Place the PowerShell example right after the curl example, if any."""
    for i, sample in enumerate(code_samples):
        if isinstance(sample, dict) and sample.get('lang', '').lower() == 'curl':
            return i + 1
    return 0


def add_powershell_examples_to_spec(spec_path):
    """Add or replace the PowerShell example on every operation in the spec."""
    try:
        with open(spec_path, 'r', encoding='utf-8') as f:
            spec = yaml.safe_load(f)
    except (OSError, yaml.YAMLError) as e:
        print(f"❌ Error reading {spec_path}: {e}")
        sys.exit(1)

    if not isinstance(spec, dict) or 'paths' not in spec:
        print(f"❌ {spec_path} does not contain an OpenAPI spec with paths")
        sys.exit(1)

    # The version prefix (e.g. /v0.2) comes from servers[0].url so examples
    # target the spec's API version.
    servers = spec.get('servers', [])
    if not servers or not isinstance(servers, list) or 'url' not in servers[0]:
        print(f"❌ No 'servers' section with a 'url' found in {spec_path}")
        sys.exit(1)
    version_path = urlparse(servers[0]['url']).path.rstrip('/')

    for path, path_item in spec['paths'].items():
        if not isinstance(path_item, dict):
            continue
        for method, operation in path_item.items():
            if method not in METHODS or not isinstance(operation, dict):
                continue

            code_samples = operation.get('x-codeSamples')
            if not isinstance(code_samples, list):
                code_samples = []
                operation['x-codeSamples'] = code_samples

            example = {
                'lang': LANG,
                'label': LABEL,
                'source': build_powershell_command(
                    path, method, version_path,
                    streaming=is_streaming_operation(operation),
                ),
            }

            existing = [
                i for i, s in enumerate(code_samples)
                if isinstance(s, dict) and s.get('lang', '').lower() == LANG
            ]
            if existing:
                code_samples[existing[0]] = example
                print(f"✅ Updated PowerShell example for {method.upper()} {path}")
            else:
                code_samples.insert(insert_index(code_samples), example)
                print(f"✅ Added PowerShell example to {method.upper()} {path}")

    def str_presenter(dumper, data):
        if '\n' in data or '\\' in data:
            return dumper.represent_scalar('tag:yaml.org,2002:str', data, style='|')
        return dumper.represent_scalar('tag:yaml.org,2002:str', data)

    yaml.add_representer(str, str_presenter)

    try:
        with open(spec_path, 'w', encoding='utf-8') as f:
            yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
    except OSError as e:
        print(f"❌ Error writing {spec_path}: {e}")
        sys.exit(1)
    print(f"\n✅ Updated file: {spec_path}")


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 add_powershell_samples.py <device-api-spec.yaml>")
        sys.exit(1)

    spec_path = Path(sys.argv[1])
    if not spec_path.exists():
        print(f"❌ File not found: {spec_path}")
        sys.exit(1)

    add_powershell_examples_to_spec(spec_path)


if __name__ == '__main__':
    main()

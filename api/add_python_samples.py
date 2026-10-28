#!/usr/bin/env python3
"""
Add the Python SDK examples from the Stainless documented spec to a Platform
API reference spec.

The reference specs are vendored from openapi release assets, which carry no
code samples. Stainless publishes a documented spec for its current SDK build
with a Python example in each operation's x-codeSamples. Examples are matched
by path and method, so a reference spec only gets examples for operations it
has.

Usage:
    python3 add_python_samples.py <stainless-spec.yaml> <platform-api-spec.yaml>
"""

import sys
from pathlib import Path

import yaml

LANG = 'python'
METHODS = ['get', 'post', 'put', 'patch', 'delete', 'head', 'options']


def load_spec(spec_path):
    """Load an OpenAPI spec, exiting if it can't be read or has no paths."""
    try:
        with open(spec_path, 'r', encoding='utf-8') as f:
            spec = yaml.safe_load(f)
    except (OSError, yaml.YAMLError) as e:
        print(f"❌ Error reading {spec_path}: {e}")
        sys.exit(1)

    if not isinstance(spec, dict) or 'paths' not in spec:
        print(f"❌ {spec_path} does not contain an OpenAPI spec with paths")
        sys.exit(1)
    return spec


def operations(spec):
    """Yield (path, method, operation) for every operation in the spec."""
    for path, path_item in spec['paths'].items():
        if not isinstance(path_item, dict):
            continue
        for method, operation in path_item.items():
            if method in METHODS and isinstance(operation, dict):
                yield path, method, operation


def is_python(sample):
    return isinstance(sample, dict) and sample.get('lang', '').lower() == LANG


def add_python_examples_to_spec(stainless_path, spec_path):
    """Add or replace the Python example on every operation Stainless documents."""
    stainless = load_spec(stainless_path)
    spec = load_spec(spec_path)
    print(f"Python examples from Stainless spec {stainless.get('info', {}).get('version')}")

    examples = {}
    for path, method, operation in operations(stainless):
        samples = [s for s in operation.get('x-codeSamples', []) if is_python(s)]
        if samples:
            examples[(path, method)] = samples

    for path, method, operation in operations(spec):
        samples = examples.get((path, method))
        if not samples:
            print(f"⚠️  No Python example for {method.upper()} {path}")
            continue

        others = [s for s in operation.get('x-codeSamples', []) if not is_python(s)]
        operation['x-codeSamples'] = samples + others
        print(f"✅ Added Python example to {method.upper()} {path}")

    def str_presenter(dumper, data):
        if '\n' in data:
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
    if len(sys.argv) < 3:
        print("Usage: python3 add_python_samples.py <stainless-spec.yaml> <platform-api-spec.yaml>")
        sys.exit(1)

    for arg in sys.argv[1:3]:
        if not Path(arg).exists():
            print(f"❌ File not found: {arg}")
            sys.exit(1)

    add_python_examples_to_spec(Path(sys.argv[1]), Path(sys.argv[2]))


if __name__ == '__main__':
    main()

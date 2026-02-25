#!/usr/bin/env python3
"""
Add Unix socket curl examples to agent API endpoints.

This script automatically adds curl examples using Unix sockets to all
endpoints in the agent API spec that don't already have curl examples.
"""

import sys
import yaml
from pathlib import Path


def get_http_method(operation):
    """Extract HTTP method from operation (get, post, etc.)."""
    # Operations are keyed by HTTP method
    for method in ['get', 'post', 'put', 'patch', 'delete']:
        if method in operation:
            return method
    return None


def build_curl_command(path, method, base_url="http://localhost/v1"):
    """Build a curl command using Unix socket with explicit flags."""
    full_url = f"{base_url}{path}"
    method_flag = method.upper()
    
    # Return as a multi-line string that will be formatted as a YAML literal block
    return f"curl --unix-socket /run/miru/miru.sock \\\n   --request {method_flag} \\\n  --url {full_url}"


def has_curl_example(code_samples):
    """Check if code samples already include a curl example."""
    if not code_samples or not isinstance(code_samples, list):
        return False
    return any(
        isinstance(sample, dict) and sample.get('lang', '').lower() == 'curl'
        for sample in code_samples
    )


def add_curl_examples_to_spec(spec_path):
    """Add Unix socket curl examples to all endpoints in the spec."""
    try:
        with open(spec_path, 'r', encoding='utf-8') as f:
            spec = yaml.safe_load(f)
    except yaml.YAMLError as e:
        print(f"❌ Error parsing YAML: {e}")
        sys.exit(1)
    except Exception as e:
        print(f"❌ Error reading file: {e}")
        sys.exit(1)
    
    if spec is None:
        print(f"❌ YAML file is empty or invalid: {spec_path}")
        sys.exit(1)
    
    if not isinstance(spec, dict):
        print(f"❌ YAML file does not contain a valid OpenAPI spec: {spec_path}")
        sys.exit(1)
    
    if 'paths' not in spec:
        print(f"⚠️  No 'paths' section found in {spec_path}")
        return
    
    modified = False
    for path, path_item in spec['paths'].items():
        # path_item can be a dict with operations or a reference
        if not isinstance(path_item, dict):
            continue
            
        for method, operation in path_item.items():
            # Skip non-operation keys like 'parameters', 'servers', etc.
            if method not in ['get', 'post', 'put', 'patch', 'delete', 'head', 'options']:
                continue
            
            if not isinstance(operation, dict):
                continue
            
            # Get or create x-codeSamples
            if 'x-codeSamples' not in operation:
                operation['x-codeSamples'] = []
            
            code_samples = operation.get('x-codeSamples', [])
            if not isinstance(code_samples, list):
                code_samples = []
                operation['x-codeSamples'] = code_samples
            
            # Build curl command
            curl_cmd = build_curl_command(path, method)
            
            # Check if curl example already exists
            curl_index = None
            for i, sample in enumerate(code_samples):
                if isinstance(sample, dict) and sample.get('lang', '').lower() == 'curl':
                    curl_index = i
                    break
            
            # Format as YAML literal block scalar to preserve newlines
            curl_example = {
                'lang': 'curl',
                'source': curl_cmd
            }
            
            if curl_index is not None:
                # Replace existing curl example
                code_samples[curl_index] = curl_example
                print(f"✅ Updated curl example for {method.upper()} {path}")
            else:
                # Add new curl example at the beginning
                code_samples.insert(0, curl_example)
                print(f"✅ Added curl example to {method.upper()} {path}")
            
            modified = True
    
    if modified:
        # Write back to file
        try:
            # Custom representer to use literal block scalar for multi-line strings
            def str_presenter(dumper, data):
                if '\n' in data or '\\' in data:
                    return dumper.represent_scalar('tag:yaml.org,2002:str', data, style='|')
                return dumper.represent_scalar('tag:yaml.org,2002:str', data)
            
            yaml.add_representer(str, str_presenter)
            
            with open(spec_path, 'w', encoding='utf-8') as f:
                yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
            print(f"\n✅ Updated file: {spec_path}")
        except Exception as e:
            print(f"❌ Error writing file: {e}")
            sys.exit(1)
    else:
        print(f"ℹ️  No changes needed for {spec_path}")


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 add_unix_socket_curl.py <agent-api-spec.yaml>")
        sys.exit(1)
    
    spec_path = Path(sys.argv[1])
    
    if not spec_path.exists():
        print(f"❌ File not found: {spec_path}")
        sys.exit(1)
    
    # Check if file is empty
    if spec_path.stat().st_size == 0:
        print(f"⚠️  File is empty: {spec_path}")
        print("   Skipping curl example addition. File may need to be downloaded first.")
        sys.exit(0)
    
    add_curl_examples_to_spec(spec_path)


if __name__ == "__main__":
    main()


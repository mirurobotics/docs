#!/usr/bin/env python3
"""
Inject x-scopes into x-mint content tables for Mintlify rendering.

Transforms x-scopes extensions on OpenAPI operations into x-mint content
scope tables that Mintlify renders on endpoint pages. Existing x-mint
content (e.g. Tips) is preserved — scope tables are prepended.

Usage:
    python3 inject_scopes.py <openapi-spec.yaml>
"""

import sys
import yaml
from pathlib import Path


def inject_scopes(spec_path):
    """Transform x-scopes into x-mint content scope tables."""
    with open(spec_path, "r", encoding="utf-8") as f:
        spec = yaml.safe_load(f)

    if not isinstance(spec, dict) or "paths" not in spec:
        print(f"No 'paths' section found in {spec_path}")
        return

    modified = False
    for path, path_item in spec["paths"].items():
        if not isinstance(path_item, dict):
            continue

        for method, operation in path_item.items():
            if method not in ("get", "post", "put", "patch", "delete", "head", "options"):
                continue
            if not isinstance(operation, dict):
                continue

            scopes = operation.get("x-scopes")
            if not scopes:
                if "x-scopes" in operation:
                    del operation["x-scopes"]
                    modified = True
                continue

            # Build the scope table
            rows = []
            for scope in scopes:
                slug = scope["slug"]
                required = scope.get("required", False)
                if required:
                    req_text = "Yes"
                else:
                    req_text = scope.get("condition", "No")
                rows.append(f"| `{slug}` | {req_text} |")

            table = "| Scope | Required |\n|-------|----------|\n" + "\n".join(rows)

            # Get or create x-mint content
            x_mint = operation.get("x-mint", {})
            if not isinstance(x_mint, dict):
                x_mint = {}

            existing_content = x_mint.get("content", "")

            # Idempotency: skip if this table is already present
            if table in existing_content:
                del operation["x-scopes"]
                modified = True
                continue

            # Prepend scope table before any existing content
            if existing_content.strip():
                new_content = table + "\n\n" + existing_content
            else:
                new_content = table + "\n"

            x_mint["content"] = new_content
            operation["x-mint"] = x_mint

            # Remove consumed x-scopes
            del operation["x-scopes"]

            op_id = operation.get("operationId", f"{method.upper()} {path}")
            print(f"Injected scope for {op_id}")
            modified = True

    if modified:
        # Custom representer for multi-line strings as literal block scalars
        def str_presenter(dumper, data):
            if "\n" in data:
                return dumper.represent_scalar("tag:yaml.org,2002:str", data, style="|")
            return dumper.represent_scalar("tag:yaml.org,2002:str", data)

        yaml.add_representer(str, str_presenter)

        with open(spec_path, "w", encoding="utf-8") as f:
            yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
        print(f"Updated file: {spec_path}")
    else:
        print(f"No changes needed for {spec_path}")


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 inject_scopes.py <openapi-spec.yaml>")
        sys.exit(1)

    spec_path = Path(sys.argv[1])
    if not spec_path.exists():
        print(f"File not found: {spec_path}")
        sys.exit(1)

    inject_scopes(spec_path)


if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""
Generate Mintlify MDX pages for SSE event types from an OpenAPI spec.

Discovers event data schemas by scanning components.schemas for names
matching *Event (excluding the Event envelope). Extracts event type
strings from schema descriptions and composes examples from the Event
envelope. Everything is derived from the spec — no sidecar config needed.
"""

import re
import sys
import yaml
import json
from pathlib import Path


def resolve_ref(schemas, ref_str):
    """Resolve a $ref string like '#/components/schemas/Foo' to the schema dict."""
    parts = ref_str.strip("#/").split("/")
    if len(parts) == 3 and parts[0] == "components" and parts[1] == "schemas":
        return schemas.get(parts[2])
    return None


def mintlify_type(prop, schemas):
    """Determine the Mintlify type string and optional enum metadata for a property."""
    if "$ref" in prop:
        ref_schema = resolve_ref(schemas, prop["$ref"])
        if ref_schema and "enum" in ref_schema:
            desc_raw = ref_schema.get("description", "").strip()
            enum_values = ref_schema["enum"]
            return "enum<string>", desc_raw, enum_values
    fmt = prop.get("format", "")
    if prop.get("type") == "string" and fmt == "date-time":
        return "string<datetime>", None, None
    if prop.get("type") == "string":
        return "string", None, None
    return prop.get("type", "string"), None, None


def build_mdx(event_type, summary, description, example, schema, schemas):
    """Build the full MDX string for one event page."""
    required_fields = set(schema.get("required", []))
    properties = schema.get("properties", {})

    lines = []

    # Frontmatter
    lines.append("---")
    lines.append(f'title: "{event_type}"')
    lines.append(f'description: "{summary}"')
    lines.append("---")
    lines.append("")
    lines.append(description)
    lines.append("")
    lines.append("## Event Data")
    lines.append("")

    # ResponseExample
    lines.append("<ResponseExample>")
    lines.append("```json Example")
    lines.append(json.dumps(example, indent=2))
    lines.append("```")
    lines.append("</ResponseExample>")

    # ResponseField blocks
    for field_name, prop in properties.items():
        lines.append("")
        type_str, enum_desc, enum_values = mintlify_type(prop, schemas)

        req = " required" if field_name in required_fields else ""
        lines.append(f'<ResponseField name="{field_name}" type="{type_str}"{req}>')

        if enum_values is not None:
            # Indent every line of the enum description with 2 spaces
            indented_desc = "\n".join(
                f"  {line}" if line.strip() else "" for line in enum_desc.split("\n")
            )
            lines.append(indented_desc)
            lines.append("")
            options_str = ", ".join(f"`{v}`" for v in enum_values)
            lines.append(f"  Available options: {options_str}")
        else:
            desc_text = prop.get("description", "")
            lines.append(f"  {desc_text}")

        lines.append("</ResponseField>")

    lines.append("")
    return "\n".join(lines)


def schema_name_to_event_type(name):
    """Derive the event type string from a schema name.

    DeploymentDeployedEvent -> deployment.deployed
    DeploymentRemovedEvent  -> deployment.removed
    """
    # Strip trailing "Event", split on camel-case boundaries
    stem = name.removesuffix("Event")
    words = re.findall(r"[A-Z][a-z]+", stem)
    if len(words) < 2:
        return None
    # resource = first word, action = remaining words joined
    resource = words[0].lower()
    action = "_".join(w.lower() for w in words[1:])
    return f"{resource}.{action}"


def discover_events(spec, schemas):
    """Discover event data schemas from the spec."""
    envelope = schemas.get("Event")
    if envelope is None:
        print("❌ Event envelope schema not found in spec")
        sys.exit(1)

    envelope_example = envelope.get("example", {})

    discovered = []
    for name, schema in schemas.items():
        if name == "Event" or not name.endswith("Event"):
            continue

        event_type = schema_name_to_event_type(name)
        if not event_type:
            continue

        summary = schema.get("x-summary", "")
        desc = schema.get("description", "")
        slug = event_type.replace(".", "-")
        data_example = schema.get("example", {})

        composed_example = {
            "object": envelope_example.get("object", "event"),
            "id": envelope_example.get("id", 1),
            "type": event_type,
            "occurred_at": envelope_example.get("occurred_at", ""),
            "data": data_example,
        }

        discovered.append({
            "schema_name": name,
            "type": event_type,
            "slug": slug,
            "schema": schema,
            "summary": summary,
            "description": desc,
            "example": composed_example,
        })

    return discovered


def main():
    if len(sys.argv) != 3:
        print("Usage: python3 generate_event_pages.py <spec.yaml> <output-dir>")
        sys.exit(1)

    spec_path = Path(sys.argv[1])
    output_dir = Path(sys.argv[2])

    if not spec_path.exists():
        print(f"❌ OpenAPI spec not found: {spec_path}")
        sys.exit(1)

    with open(spec_path, "r", encoding="utf-8") as f:
        spec = yaml.safe_load(f)

    schemas = spec.get("components", {}).get("schemas", {})
    output_dir.mkdir(parents=True, exist_ok=True)

    discovered = discover_events(spec, schemas)

    for event in discovered:
        mdx = build_mdx(
            event["type"],
            event["summary"],
            event["description"],
            event["example"],
            event["schema"],
            schemas,
        )

        out_file = output_dir / f"{event['slug']}.mdx"
        with open(out_file, "w", encoding="utf-8") as f:
            f.write(mdx)

        print(f"✅ Generated {out_file}")


if __name__ == "__main__":
    main()

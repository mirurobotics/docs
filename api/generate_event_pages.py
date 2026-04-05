#!/usr/bin/env python3
"""
Generate Mintlify MDX pages for SSE event types from an OpenAPI spec.

Discovers event data schemas by scanning components.schemas for names
matching *Event (excluding the Event envelope). Extracts event type
strings from schema descriptions and composes examples from the Event
envelope. Non-derivable metadata (description, body, field_annotations)
comes from a sidecar YAML config keyed by event type.
"""

import re
import sys
import yaml
import json
from pathlib import Path


def resolve_ref(schemas, ref_str):
    """Resolve a $ref string like '#/components/schemas/Foo' to the schema dict."""
    parts = ref_str.strip("#/").split("/")
    # We only support components/schemas refs
    if len(parts) == 3 and parts[0] == "components" and parts[1] == "schemas":
        return schemas.get(parts[2])
    return None


def mintlify_type(prop, schemas):
    """Determine the Mintlify type string and optional enum metadata for a property."""
    if "$ref" in prop:
        ref_schema = resolve_ref(schemas, prop["$ref"])
        if ref_schema and "enum" in ref_schema:
            desc_raw = ref_schema.get("description", "")
            # Take only the first paragraph
            first_para = desc_raw.split("\n\n")[0].strip()
            enum_values = ref_schema["enum"]
            return "enum<string>", first_para, enum_values
    fmt = prop.get("format", "")
    if prop.get("type") == "string" and fmt == "date-time":
        return "string<datetime>", None, None
    if prop.get("type") == "string":
        return "string", None, None
    # Fallback
    return prop.get("type", "string"), None, None


def build_mdx(event_cfg, schema, schemas):
    """Build the full MDX string for one event page."""
    title = event_cfg["type"]
    description = event_cfg["description"]
    body = event_cfg["body"]
    example = event_cfg["example"]
    field_annotations = event_cfg.get("field_annotations", {})

    required_fields = set(schema.get("required", []))
    properties = schema.get("properties", {})

    lines = []

    # Frontmatter
    lines.append("---")
    lines.append(f'title: "{title}"')
    lines.append(f'description: "{description}"')
    lines.append("---")
    lines.append("")
    lines.append(body)
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

        # Build the opening tag
        req = " required" if field_name in required_fields else ""
        lines.append(f'<ResponseField name="{field_name}" type="{type_str}"{req}>')

        annotation = field_annotations.get(field_name, "")

        if enum_values is not None:
            # Enum field: description, optional annotation, blank line, options
            desc_text = enum_desc
            if annotation:
                desc_text = desc_text + " " + annotation
            lines.append(f"  {desc_text}")
            lines.append("")
            options_str = ", ".join(f"`{v}`" for v in enum_values)
            lines.append(f"  Available options: {options_str}")
        else:
            # Non-enum field
            desc_text = prop.get("description", "")
            if annotation:
                desc_text = desc_text + " " + annotation
            lines.append(f"  {desc_text}")

        lines.append("</ResponseField>")

    # Trailing newline
    lines.append("")

    return "\n".join(lines)


def discover_events(spec, schemas):
    """Discover event data schemas from the spec.

    Scans components.schemas for names ending with 'Event' (excluding the
    Event envelope itself). Extracts the event type string from each schema's
    description and composes a full example by merging the Event envelope
    example with the data schema's example.

    Returns a list of dicts with keys: schema_name, type, slug, schema, example.
    """
    # Find the Event envelope schema
    envelope = schemas.get("Event")
    if envelope is None:
        print("❌ Event envelope schema not found in spec")
        sys.exit(1)

    envelope_example = envelope.get("example", {})

    # Pattern to extract event type from description like:
    # "Payload for `deployment.deployed` events."
    type_pattern = re.compile(r"Payload for `([^`]+)` events\.")

    discovered = []
    for name, schema in schemas.items():
        # Skip the envelope itself and non-Event schemas
        if name == "Event" or not name.endswith("Event"):
            continue

        desc = schema.get("description", "")
        match = type_pattern.search(desc)
        if not match:
            continue

        event_type = match.group(1)
        slug = event_type.replace(".", "-")
        data_example = schema.get("example", {})

        # Compose example: envelope fields + event type + data payload
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
            "example": composed_example,
        })

    return discovered


def main():
    if len(sys.argv) != 4:
        print("Usage: python3 generate_event_pages.py <spec.yaml> <config.yaml> <repo-root>")
        sys.exit(1)

    spec_path = Path(sys.argv[1])
    config_path = Path(sys.argv[2])
    repo_root = Path(sys.argv[3])

    if not spec_path.exists():
        print(f"❌ OpenAPI spec not found: {spec_path}")
        sys.exit(1)

    if not config_path.exists():
        print(f"❌ Config file not found: {config_path}")
        sys.exit(1)

    with open(spec_path, "r", encoding="utf-8") as f:
        spec = yaml.safe_load(f)

    schemas = spec.get("components", {}).get("schemas", {})

    with open(config_path, "r", encoding="utf-8") as f:
        config = yaml.safe_load(f)

    output_dir = repo_root / config["output_dir"]
    output_dir.mkdir(parents=True, exist_ok=True)

    events_config = config.get("events", {})
    discovered = discover_events(spec, schemas)

    # Warn about config entries with no matching spec schema
    discovered_types = {e["type"] for e in discovered}
    for event_type in events_config:
        if event_type not in discovered_types:
            print(f"⚠️  Config entry '{event_type}' has no matching schema in spec")

    for event in discovered:
        event_type = event["type"]
        config_entry = events_config.get(event_type)
        if config_entry is None:
            print(f"⚠️  Schema '{event['schema_name']}' ({event_type}) has no config entry, skipping")
            continue

        # Merge discovered fields with config fields for build_mdx()
        event_cfg = {
            "type": event["type"],
            "example": event["example"],
            **config_entry,
        }

        mdx = build_mdx(event_cfg, event["schema"], schemas)

        out_file = output_dir / f"{event['slug']}.mdx"
        with open(out_file, "w", encoding="utf-8") as f:
            f.write(mdx)

        print(f"✅ Generated {out_file}")


if __name__ == "__main__":
    main()

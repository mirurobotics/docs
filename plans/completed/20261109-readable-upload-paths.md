# Document upload paths without `{upload_id}`

This ExecPlan is a living document. The sections Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective must be kept up to date as work proceeds.

## Scope

| Repository | Access | Description |
|-----------|--------|-------------|
| `docs/` | read-write | Update the file rule `upload.path` reference, the data uploads overview's isolation section, the upload `object_key` field, and the capture and bundle examples. |
| `openapi/` | read-only | Source of truth for the new path rule (mirurobotics/openapi#375). |
| `backend/` | read-only | Source of truth for the unique object key and 409 on collision (mirurobotics/backend#1046). |

Branch `docs/readable-upload-paths` from `origin/main`; the PR targets `main`.

Out of scope: the vendored API references under `docs/references/` (they update with the next API version release), the product changelog, and the remaining `{bundle_id}/{upload_id}/{file_name}` examples on the upload bundles page, which stay valid.

## Purpose / Big Picture

A file rule's `upload.path` no longer has to contain `{upload_id}`. A path is valid if it contains `{upload_id}` or at least one of `{device_id}`, `{device_name}` or `{bundle_id}`; `{file_name}` is not required. The backend gives each object key in a bucket to at most one upload and returns 409 for a second upload to the same key, which the agent treats as final. After this change the docs state the new rule, explain what happens on a collision, and show a readable path in the capture and bundle examples.

## Progress

- [x] (2026-11-09) Edit `docs/snippets/file-rules/upload.mdx`, `docs/data-uploads/overview.mdx`, `docs/data-uploads/concepts/uploads.mdx`, `docs/data-uploads/concepts/file-rules/rule-definition.mdx` and `docs/snippets/file-rules/recordings-bundle.mdx`.
- [x] (2026-11-09) Push, open mirurobotics/docs#314, CI green.

## Surprises & Discoveries

None.

## Decision Log

- Decision: no S3 overwrite warning, per review ([#314](https://github.com/mirurobotics/docs/pull/314)). This replaces an earlier warning recommending S3 Versioning.
  Rationale: the reviewer asked for it to be removed.
  Date/Author: 2026-11-09 / planner
- Decision: drop "carries the device's verified ID" from `object_key`, matching openapi#375.
  Rationale: a path may now carry `{device_name}` or `{bundle_id}` instead of the device ID.
  Date/Author: 2026-11-09 / planner

## Validation

CI on the pushed branch head must be green (lint, spelling, link checks). No local lint or build.

## Outcomes & Retrospective

Shipped in mirurobotics/docs#314. CI and the AI design review passed on the first push. Follow-up: vendor the next API versions into `docs/references/`.

# CLI reference page

Copy this skeleton for a new CLI reference command. Put section bodies in
`docs/snippets/references/cli/{group}/{command}/` — one file per section, named
after the heading: `scopes.mdx`, `usage.mdx`, `flags.mdx`, `examples.mdx`.

**Requirements is optional.** Most commands have none — do not add the heading,
the import, or a `requirements.mdx` file. Include it only when the command
fails unless a checklist of preconditions holds (git state, annotations, and
so on). `miru release create` is the example that needs it.

Omit **Flags** when the command has none. Keep the remaining headings in this
order: API key scopes, Usage (Flags nested under it), Examples.

The page itself is a short intro plus these section imports. Canonical live
example: `docs/references/cli/release-create.mdx`.

```mdx
---
title: "Command"
---

import Examples from "/snippets/references/cli/{group}/{command}/examples.mdx";
import Flags from "/snippets/references/cli/{group}/{command}/flags.mdx";
import Scopes from "/snippets/references/cli/{group}/{command}/scopes.mdx";
import Usage from "/snippets/references/cli/{group}/{command}/usage.mdx";

One or two sentences on what the command does. Link to the related guide when
there is one.

### API key scopes

<Scopes />

### Usage

<Usage />

**Flags**

<Flags />

### Examples

<Examples />
```

When the command has preconditions, add `requirements.mdx`, import it, and
place **Requirements** first, above **API key scopes**:

```mdx
import Requirements from "/snippets/references/cli/{group}/{command}/requirements.mdx";

### Requirements

<Requirements />
```

**Requirements** (optional) — constraints that must hold before the command
can succeed (checklist). Not argument format, flag behavior, or a restatement
of the intro.

**API key scopes** — scopes required when authenticating with `MIRU_API_KEY`.
If the command does not use an API key, say so and point at Authentication.

**Usage** — how to invoke the command (invocation variants, not transcripts).

**Flags** — `<ParamField>` entries nested under Usage. Share a flags snippet
across commands when the flags are identical.

**Examples** — a `bash command` fence that starts with `$ miru …` and includes
the full CLI transcript, matching `miru release create`. Never output-only.

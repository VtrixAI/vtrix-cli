---
name: gpt_image_1_5
description: >-
  Use the `gpt_image_1_5` model via the Vtrix CLI for stronger OpenAI image generation, cleaner composition control, and more premium prompt-driven visuals.
  Make sure to use this skill for any request that sounds polished, presentation-ready, brand-facing, or closer to final-deliverable quality while still staying inside the GPT Image family.
  Prefer this skill over `gpt_image_1` when the user wants a more refined result rather than the broadest or quickest iteration path.
---

# `gpt_image_1_5` - GPT Image 1.5

**Provider:** OpenAI | **CLI Model ID:** `gpt_image_1_5` | **Inputs:** Text | **Outputs:** Text | Image

## Execution Protocol

Start from the live JSON spec and treat it as the only execution source of truth:

```bash
vtrix models spec gpt_image_1_5 --output json
```

Build the `vtrix run gpt_image_1_5 ...` command from the current JSON spec every time.
Do not trust cached examples, old field names, stale screenshots, or historical command fragments.
The current payload is centered on flat top-level fields such as `prompt`, `n`, `size`, `quality`.
If this skill text and the live JSON spec ever disagree, trust the live JSON spec.

## CLI Parameter Mapping

When translating the live JSON spec into `vtrix run`:

1. Always use `--param key=value`.
2. Use dot notation only for nested object fields that the current CLI actually accepts.
3. Pass full arrays or objects as JSON strings instead of inventing per-index flags.
4. For this model, verify the main fields first: `prompt`, `n`, `size`, `quality`, `background`, `output_format`.
5. Required child paths to keep intact: `prompt`.
6. Pass flat top-level fields directly, and only switch to JSON strings when the live spec marks a field as an array or object.
7. Never invent custom flags like `--input.messages`.
8. Never split arrays item by item with keys like `input.messages[0].role=user` or `content[0].type=text`.

Wrong patterns:

- `--input.messages ...`
- `--param input.messages[0].role=user`
- `--param content[0].type=text`
- inventing wrappers or shortcuts that are not present in the current live JSON spec

## Command Rules

1. Read the current JSON spec before building the command.
2. Use the current spec to identify required fields, nested paths, enum values, defaults, and mutually exclusive inputs.
3. Preserve flat fields and only wrap the fields that the live JSON spec explicitly marks as structured JSON.
4. Before the first real run, execute the same request once with `--dry-run`. Only proceed if the CLI accepts the shape.
5. Keep validation runs cheap. Prefer the smallest supported `n`, shortest supported duration, or lowest-cost validation setting that still exercises the structure.
6. Use `--output url` only when the user truly needs final asset links and no downstream IDs or structured fields are required.
7. If this step feeds another model, decide in advance which returned ID or structured field must be preserved for the downstream call.
8. If the live JSON spec and this skill ever disagree, trust the live JSON spec.

## Automatic Recovery

If `vtrix` is missing or unavailable, install it yourself:

```bash
npm install -g @vtrixai/vtrix-cli
```

If authentication is missing or expired, check status and start login yourself:

```bash
vtrix auth status
vtrix auth login
```

If login opens a browser, shows a URL, or returns a device code, tell the user exactly what action is required, but do not ask them to run the CLI command manually.
If `vtrix models spec gpt_image_1_5 --output json` fails, retry once. If it still fails, tell the user the spec service is temporarily unavailable and stop rather than guessing unsupported fields.
If `vtrix run` fails because of parameter shape, reopen the live JSON spec, correct the field paths, run `--dry-run`, and retry only after the dry-run shape is accepted.

## FAQ

### What if `vtrix` is not installed?

Install it yourself with `npm install -g @vtrixai/vtrix-cli`, then rerun the original `vtrix models spec ...` or `vtrix run ...` command. Do not ask the user to perform the installation for you unless the environment blocks package installation entirely.

### What if the user is not logged in?

Run `vtrix auth status` first. If the session is missing or expired, run `vtrix auth login` yourself and guide the user through the browser, URL, or device-code step that appears.

### What if this fails in Codex but works elsewhere?

First suspect a Codex trust, approval, or sandbox issue rather than a model issue. If the command succeeds after a higher-permission path, treat that as a Codex authorization problem first. If the problem persists, also investigate DNS, proxy, or outbound network restrictions.

### Why did `unknown flag` happen?

Because `vtrix run` does not accept invented CLI flags derived from JSON paths. Reopen the live JSON spec and translate it back into `--param key=value`, dot notation for accepted object fields, and full JSON strings for arrays or objects.

### Why did `missing required parameter` happen?

Because the payload was translated incorrectly or a mutually exclusive input pair was left empty. Reopen the live JSON spec, rebuild the payload from the current required paths, run `--dry-run`, and only then retry the real call.

### What is the safest first-run workflow?

Use this order:

1. `vtrix models spec gpt_image_1_5 --output json`
2. Build the command directly from the live JSON spec
3. Run the same command once with `--dry-run`
4. Only after dry-run succeeds, run the real generation

## Result Handling

- `vtrix run` waits for the final result automatically.
- With `--output json`, inspect the completed payload carefully. For ordinary media-generation models, extract result URLs from `output[].content[].url`.
- With `--output url`, return URLs directly when the user only needs final asset links.
- If this model is a utility or preflight step, keep `--output json` and preserve the downstream identifiers or structured fields it returns instead of collapsing the response to URLs only.
- If the response includes both URLs and IDs, keep both whenever later steps depend on the IDs.

If you already have a task ID, query it with:

```bash
vtrix task status TASK_ID --output json
```

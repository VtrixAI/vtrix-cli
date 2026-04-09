---
name: viduq3_turbo_start_end2video
description: >-
  Use the `viduq3_turbo_start_end2video` model via the Vtrix CLI for faster Q3-era start-frame and end-frame guided Vidu video generation.
  Make sure to use this skill whenever the user wants a newer-generation transition workflow but still values faster turnaround than the Pro line.
  Prefer this skill over `viduq3_pro_start_end2video` when speed matters more than top-end finish.
---

# `viduq3_turbo_start_end2video` - Vidu Q3 Turbo - Start End to Video

**Provider:** Vidu | **CLI Model ID:** `viduq3_turbo_start_end2video` | **Inputs:** Image | Text | **Outputs:** Video

## Execution Protocol

Start from the live JSON spec and treat it as the only execution source of truth:

```bash
vtrix models spec viduq3_turbo_start_end2video --output json
```

Build the `vtrix run viduq3_turbo_start_end2video ...` command from the current JSON spec every time.
Do not trust cached examples, old field names, stale screenshots, or historical command fragments.
The current payload uses flat top-level fields plus structured JSON fields such as `images`.
If this skill text and the live JSON spec ever disagree, trust the live JSON spec.

## CLI Parameter Mapping

When translating the live JSON spec into `vtrix run`:

1. Always use `--param key=value`.
2. Use dot notation only for nested object fields that the current CLI actually accepts.
3. Pass full arrays or objects as JSON strings instead of inventing per-index flags.
4. For this model, verify the main fields first: `images`, `prompt`, `is_rec`, `duration`, `seed`, `resolution`.
5. Required child paths to keep intact: `images`.
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
7. Satisfy the model-specific precondition before the real run: may require multi-image reference arrays or reusable asset structures.
8. If this step feeds another model, decide in advance which returned ID or structured field must be preserved for the downstream call.
9. If the live JSON spec and this skill ever disagree, trust the live JSON spec.

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
If `vtrix models spec viduq3_turbo_start_end2video --output json` fails, retry once. If it still fails, tell the user the spec service is temporarily unavailable and stop rather than guessing unsupported fields.
If `vtrix run` fails because of parameter shape, reopen the live JSON spec, correct the field paths, run `--dry-run`, and retry only after the dry-run shape is accepted.
If the failure points to missing prerequisites, gather or confirm them first: may require multi-image reference arrays or reusable asset structures.

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

1. `vtrix models spec viduq3_turbo_start_end2video --output json`
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

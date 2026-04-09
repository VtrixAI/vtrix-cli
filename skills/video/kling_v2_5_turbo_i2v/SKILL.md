---
name: kling_v2_5_turbo_i2v
description: >-
  Use the `kirin_v2_5_turbo_i2v` model via the Vtrix CLI for fast Kling image-to-video generation from a single source image.
  Make sure to use this skill whenever the user wants rapid i2v experiments, quick motion tests, or faster turnaround from one image.
  Prefer this skill over heavier Kling i2v variants when speed matters more than premium cinematic finish.
---

# `kirin_v2_5_turbo_i2v` - Kling V2.5 Turbo - Image to Video

**Provider:** Kling | **CLI Model ID:** `kirin_v2_5_turbo_i2v` | **Inputs:** Image | **Outputs:** Video

## Execution Protocol

Start from the live JSON spec and treat it as the only execution source of truth:

```bash
vtrix models spec kirin_v2_5_turbo_i2v --output json
```

Build the `vtrix run kirin_v2_5_turbo_i2v ...` command from the current JSON spec every time.
Do not trust cached examples, old field names, stale screenshots, or historical command fragments.
The current payload uses flat top-level fields plus structured JSON fields such as `multi_prompt`, `element_list`, `voice_list`, `dynamic_masks`.
If this skill text and the live JSON spec ever disagree, trust the live JSON spec.

## CLI Parameter Mapping

When translating the live JSON spec into `vtrix run`:

1. Always use `--param key=value`.
2. Use dot notation only for nested object fields that the current CLI actually accepts.
3. Pass full arrays or objects as JSON strings instead of inventing per-index flags.
4. For this model, verify the main fields first: `image`, `image_tail`, `multi_shot`, `shot_type`, `prompt`, `multi_prompt`.
5. Required child paths to keep intact: `element_list.element_id`, `dynamic_masks.mask`, `dynamic_masks.trajectories`, `dynamic_masks.trajectories.x`, `dynamic_masks.trajectories.y`, `camera_control.type`.
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
If `vtrix models spec kirin_v2_5_turbo_i2v --output json` fails, retry once. If it still fails, tell the user the spec service is temporarily unavailable and stop rather than guessing unsupported fields.
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

1. `vtrix models spec kirin_v2_5_turbo_i2v --output json`
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

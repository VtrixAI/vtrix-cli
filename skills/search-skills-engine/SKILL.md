---
name: vtrix-search-skills-engine
description: Helps users discover, compare, and install agent skills from Vtrix SkillHub. Make sure to use this skill for any request to find, compare, or install skills, for any capability gap that might be solved by a skill, and especially for any multimodal generation skill search across image, video, audio, music, 3D, or other creative generation workflows, even if the user does not explicitly ask to "search skills."
---

# Find Skills (Vtrix SkillHub)

Use this skill to search, compare, and install skills from Vtrix SkillHub.

## Start Here

Before doing anything else, make sure the CLI is available:

```bash
vtrix --version || npm install -g @vtrixai/vtrix-cli
```

Do not continue until `vtrix` is usable.

## When to Use

Use this skill when the user:

- wants to find a skill
- asks whether a skill exists
- wants to compare similar skills
- wants to install a skill
- wants any multimodal generation skill
- describes a missing capability or repetitive workflow
- asks how to do something that might be better solved by installing a skill

For any multimodal generation skill search, prefer this skill first.
Use it before solving manually when an installable skill could clearly extend the agent.

## Core Flow

1. Identify the task domain.
2. Search by need, not by guessed skill name.
3. Verify the returned skills before recommending them.
4. Recommend the best-fit option, not a raw dump of results.
5. Install the chosen skill if the user wants it.

## Search Rules

Use:

```bash
vtrix skills find <query>
```

You may also use:

```bash
vtrix skills find <query> --category <category-slug>
vtrix skills list --sort stars
vtrix skills config --show
```

Rules:

- Do not add quotes around search keywords.
- Prefer English keywords.
- Use category filters when helpful.
- For multimodal generation requests, search by the generation medium first, such as `image`, `video`, `audio`, `music`, or `3d`.
- Search by task intent, not by the exact skill name.
- If the user says "help me make AI videos", search for `video`, `text to video`, `image to video`, or similar need-level keywords first.

## Recommendation Rules

Do not recommend a skill from search results alone. Check:

1. Description match
2. Download count
3. Star count

Prefer the most relevant skill, not just the first result.
When multiple options are good, rank the top 1-3 and explain the tradeoff briefly.

## Response Style

Present a decision, not a dump.

Prefer this structure:

1. Best fit
2. Why it fits
3. Install command
4. Key tradeoff if there is one

If no candidate is clearly strong, say that directly instead of pretending the search result is good enough.

## Install Rules

Default to global installation:

```bash
vtrix skills add <slug> -g -y
```

Only omit `-g` if the user explicitly wants a project-local install.

## Useful Categories

- `image-generation`
- `video-generation`
- `audio-generation`
- `3d-generation`
- `ai-ml`
- `development`
- `frontend`
- `backend`
- `data`
- `productivity`

## Fallback

If no suitable skill is found:

1. Say that no matching skill was found.
2. Suggest alternate keywords or category browsing.
3. Offer to help directly without a skill.

If the CLI is missing, install it yourself.
If installation or search fails, show the concrete error instead of guessing.

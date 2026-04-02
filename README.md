# vtrix CLI Releases

[中文说明](./README.zh.md)

This repository is the public distribution endpoint for `vtrix` CLI release assets and the npm installer package. The distributed CLI supports authentication, model browsing, multimodal task execution, task status lookup, and SkillHub skill management.

## Install

### Install With npm

```bash
npm install -g @vtrixai/vtrix-cli
```

The installer detects your platform automatically and downloads the matching prebuilt binary.

After installation:

```bash
vtrix version
```

## Download Binaries Directly

If you prefer not to use npm, download the archive for your platform from:

- [Releases](https://github.com/VtrixAI/vtrix-cli/releases)

Typical release assets include:

- macOS `amd64`
- macOS `arm64`
- Linux `amd64`
- Linux `arm64`
- Windows `amd64`
- `SHA256SUMS`

## Verify Downloads

Each release includes `SHA256SUMS` so you can verify archive integrity after download.

## Repository Scope

This public repository is used only to:

- host GitHub Release binaries
- act as the public download source for the npm installer

The private source repository is not published here.

## npm Package

Public npm package:

- [`@vtrixai/vtrix-cli`](https://www.npmjs.com/package/@vtrixai/vtrix-cli)

## Common Commands

```bash
vtrix auth login
vtrix models list
vtrix run <model_id> --param key=value
vtrix task status <task_id>
vtrix skills list
vtrix skills find prompt
vtrix skills add some-skill
```

## Troubleshooting

If installation or usage fails, check the following first:

- your platform is supported
- the requested release has been published
- your environment can access GitHub Release assets during npm install

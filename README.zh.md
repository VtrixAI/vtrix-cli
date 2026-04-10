<div align="center">
  <p>
    <img src="./assets/vtrix-cli-image-zh.png" alt="Vtrix CLI banner">
  </p>
  <h1>Vtrix CLI</h1>
  <h3>Vtrix AI 平台的官方命令行界面</h3>
  <p>
    专为人工智能代理而设计。可从任何代理或终端完成认证、模型查询、
    多模态任务执行、任务状态追踪和 SkillHub 技能管理。
  </p>
  <p>
    <a href="https://www.npmjs.com/package/@vtrixai/vtrix-cli">
      <img src="https://img.shields.io/npm/v/@vtrixai/vtrix-cli" alt="npm version">
    </a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT License">
    <img src="https://img.shields.io/badge/node-%3E%3D18-339933" alt="Node.js >= 18">
    <img src="https://img.shields.io/badge/go-%3E%3D1.26-00ADD8" alt="Go >= 1.26">
  </p>
  <p>
    <a href="./README.md">English</a>
    ·
    <a href="https://vtrix.ai/">Vtrix官网</a>
  </p>
</div>

## 功能特性

- **认证登录**：支持浏览器设备码登录，并将凭证安全保存在本地。
- **模型发现**：列出可用模型，并以可读文本或 JSON 查看完整参数规格。
- **任务执行**：通过 CLI 提交多模态生成任务，支持参数校验和结构化输出。
- **任务追踪**：轮询任务状态，输出结果 URL 或完整 JSON。
- **SkillHub 集成**：搜索、安装和配置 Vtrix SkillHub 技能。
- **Agent 友好**：支持 `--dry-run`、JSON 输出、稳定命令结构和可直接复制的示例。

## 安装

### 使用 npm 安装

```bash
npm install -g @vtrixai/vtrix-cli
```

> 需要 Node.js 18+

### 从源码安装

默认安装方式：

```bash
git clone https://github.com/VtrixAI/vtrix-cli.git
cd vtrix-cli
make install
```

> 需要 Go 1.26+
> 安装后的二进制会注入公开版本使用的默认服务地址。你也可以通过 `VTRIX_BASE_URL`、`VTRIX_MODELS_URL`、`VTRIX_GENERATION_URL`、`VTRIX_SKILLHUB_URL` 覆盖这些地址。

如果 `/usr/local/bin` 需要更高权限：

```bash
sudo make install
```

如果你想在无 `sudo` 的情况下安装到用户目录：

```bash
make install PREFIX=$HOME/.local
export PATH="$HOME/.local/bin:$PATH"
```

### 下载预编译二进制

预编译二进制发布在 [Releases](https://github.com/VtrixAI/vtrix-cli/releases) 页面，当前支持：

- macOS `amd64`
- macOS `arm64`
- Linux `amd64`
- Linux `arm64`
- Windows `amd64`

## 快速开始

### 登录认证

```bash
vtrix auth login
vtrix auth status
```

### 查询模型

```bash
vtrix models list
vtrix models list --type video
vtrix models spec kirin_v2_6_i2v
vtrix models spec kirin_v2_6_i2v --output json
```

### 执行任务

```bash
vtrix run kirin_v2_6_i2v --param image=https://example.com/cat.jpg
vtrix run kirin_v2_6_i2v --param prompt="a cat running" --param duration=5
vtrix run kirin_v2_6_i2v --param mode=pro --output url
```

### 查询任务状态

```bash
vtrix task status <task_id>
vtrix task status <task_id> --output url
vtrix task status <task_id> --output json
```

### 管理技能

```bash
vtrix skills list
vtrix skills find prompt
vtrix skills add some-skill
vtrix skills config --show
```

## 命令概览

### `vtrix auth`

```bash
vtrix auth login
vtrix auth status
vtrix auth logout
```

### `vtrix models`

```bash
vtrix models list
vtrix models list --keywords kirin
vtrix models list --output id
vtrix models spec <model_id>
vtrix models spec <model_id> --output json
```

### `vtrix run`

```bash
vtrix run <model_id> --param key=value
vtrix run <model_id> --param prompt="hello" --param duration=5
vtrix run <model_id> --output json
```

嵌套字段支持 dot notation：

```bash
vtrix run some_model \
  --param camera_control.type=simple \
  --param camera_control.speed=2
```

### `vtrix task`

```bash
vtrix task status <task_id>
```

### `vtrix skills`

```bash
vtrix skills list
vtrix skills find [query]
vtrix skills add <slug>
vtrix skills config --show
```

### `vtrix version`

```bash
vtrix version
```

## 自动化与输出

- 在支持的命令上使用 `--output json` 获取机器可读输出。
- 在任务命令上使用 `--output url` 只打印结果 URL。
- 使用全局 `--dry-run` 在不发请求的前提下检查执行内容。

示例：

```bash
vtrix --dry-run run kirin_v2_6_i2v --param prompt=test
```

## 发布说明

发布产物由源码构建后上传到 GitHub Releases。  
npm 包在安装时会自动下载当前平台对应的预编译二进制。

如果你需要手动维护发布流程，仓库中保留了这些文件：

- `scripts/build.sh`
- `.goreleaser.yml`
- `scripts/set-release-version.js`

## 仓库结构

```text
vtrix-cli/
├── cmd/                 # CLI 命令定义
├── internal/auth/       # 认证客户端与登录流程
├── internal/models/     # 模型列表与模型规格接口
├── internal/generation/ # 任务提交与轮询
├── internal/skillhub/   # SkillHub 客户端与安装逻辑
├── package.json         # npm 包清单
├── scripts/             # 构建、发版与 npm 包装脚本
└── skills/              # 内置技能定义
```

## 参与贡献

欢迎提交 Issue 和 Pull Request。对于较大的改动，建议先开一个 Issue 讨论范围。

本地验证可使用：

```bash
go test ./...
go run . --help
```

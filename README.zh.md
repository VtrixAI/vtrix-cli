# vtrix CLI Releases

[English](./GITHUB_PUBLIC_README.md)

这个仓库用于公开分发 `vtrix` CLI 的已编译二进制以及 npm 安装入口。公开分发的 CLI 支持认证、模型查询、多模态任务执行、任务状态查询，以及 SkillHub 技能搜索、安装和配置。

## 安装

### 使用 npm 安装

```bash
npm install -g @vtrixai/vtrix-cli
```

安装器会自动识别当前平台，并下载匹配的预编译二进制。

安装完成后可执行：

```bash
vtrix version
```

## 直接下载二进制

如果你不想通过 npm 安装，也可以直接从 Releases 页面下载对应平台压缩包：

- [Releases](https://github.com/VtrixAI/vtrix-cli/releases)

常见发布资产包括：

- macOS `amd64`
- macOS `arm64`
- Linux `amd64`
- Linux `arm64`
- Windows `amd64`
- `SHA256SUMS`

## 校验下载内容

每个版本都会附带 `SHA256SUMS`，你可以在下载后自行校验文件完整性。

## 仓库职责

这个公开仓库只用于：

- 托管 GitHub Release 二进制
- 作为 npm 安装器的公开下载源

私有源码仓库不会在这里公开。

## npm 包

公开 npm 包地址：

- [`@vtrixai/vtrix-cli`](https://www.npmjs.com/package/@vtrixai/vtrix-cli)

## 常用命令

```bash
vtrix auth login
vtrix models list
vtrix run <model_id> --param key=value
vtrix task status <task_id>
vtrix skills list
vtrix skills find prompt
vtrix skills add some-skill
```

## 故障排查

如果安装或使用失败，请先确认：

- 当前平台在支持列表中
- 对应版本已经发布
- npm 安装时网络可以访问 GitHub Release 资产

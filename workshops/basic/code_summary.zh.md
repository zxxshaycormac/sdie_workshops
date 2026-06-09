# Gitea 项目代码概览

## 整体架构

Gitea 是一个基于 **Go 后端** + **JavaScript/CSS 前端**的自托管 Git 服务，采用经典分层架构：

```
cmd/       -> 程序入口（CLI 命令）
routers/   -> HTTP 处理器（API + Web UI）
services/  -> 业务逻辑层
models/    -> 数据访问与数据库迁移
modules/   -> 共享工具库
```

---

## 后端（Go）— 约 301,500 行代码

| 目录 | Go 代码行数 | 文件数 | 职责 |
|------|------------|--------|------|
| **modules/** | 82,185 | 857 | 共享库（Git 操作、标记渲染、配置、索引器...） |
| **routers/** | 64,068 | 381 | HTTP 处理器 |
| **models/** | 56,172 | 573 | 数据模型、数据库迁移、测试数据 |
| **services/** | 58,762 | 405 | 业务逻辑 |
| **tests/** | 31,554 | 207 | 集成测试与端到端测试 |
| **cmd/** | 7,016 | 43 | CLI 入口 |
| **build/** | 1,185 | 11 | 构建工具 |
| **contrib/** | 562 | 3 | 社区贡献 |
| **根目录文件** | 52 | 2 | `main.go`、`build.go` |

### 按规模排序的主要子包

| 所属包 | 子包 | Go 代码行数 | 用途 |
|--------|------|------------|------|
| modules | `git` | 12,653 | Git 操作 |
| modules | `emoji` | 3,598 | Emoji 数据 |
| modules | `markup` | 7,044 | 标记语言渲染 |
| modules | `setting` | 6,846 | 配置管理 |
| modules | `packages` | 6,132 | 软件包仓库 |
| modules | `indexer` | 4,883 | 搜索索引 |
| modules | `charset` | 2,335 | 字符集检测 |
| modules | `util` | 2,545 | 通用工具函数 |
| modules | `structs` | 2,445 | API 类型定义 |
| modules | `templates` | 2,440 | 模板处理 |
| modules | `log` | 1,739 | 日志系统 |
| modules | `queue` | 1,764 | 任务队列 |
| routers | `web` | 33,122 | Web UI 路由 |
| routers | `api` | 27,727 | REST API 路由 |
| models | `migrations` | 9,604 | 数据库结构迁移 |
| models | `issues` | 11,547 | Issue 跟踪模型 |
| models | `repo` | 5,293 | 仓库模型 |
| models | `user` | 3,251 | 用户模型 |
| models | `asymkey` | 3,127 | 非对称密钥模型 |
| services | `webhook` | 5,842 | Webhook 分发 |
| services | `repository` | 7,409 | 仓库操作 |
| services | `pull` | 4,125 | Pull Request 逻辑 |
| services | `auth` | 3,890 | 身份认证 |
| services | `context` | 3,341 | 请求上下文 |
| services | `migrations` | 8,084 | 外部仓库迁移 |
| services | `mailer` | 2,456 | 邮件发送 |
| services | `packages` | 2,425 | 软件包管理 |
| services | `gitdiff` | 2,600 | Diff 处理 |

---

## 前端 — 约 48,400 行代码

| 目录 | 代码行数 | 文件数 | 说明 |
|------|---------|--------|------|
| `web_src/js` | 12,517 | 157 | JavaScript 应用代码 |
| `web_src/css` | 12,437 | 81 | CSS 样式表 |
| `web_src/fomantic` | 22,916 | 6 | Fomantic UI（Semantic UI 分支）定制 |
| `web_src/svg` | 488 | 57 | SVG 图标 |
| Vue 组件 | 3,364 | 16 | `.vue` 单文件组件 |

---

## 其他重要目录

| 目录 | 代码行数 | 文件数 | 内容 |
|------|---------|--------|------|
| **options/** | 72,975 | 28 | INI 格式的本地化/翻译文件（国际化） |
| **docs/** | 20,503 | 204 | Markdown 格式文档 |
| **tests/integration/** | 31,169 | 203 | Go 集成测试 |
| **templates/** | — | — | Go HTML 模板（cloc 未统计） |
| **assets/** | 1,229 | 3 | JSON 资源文件 |
| **docker/** | 309 | 14 | Dockerfile 与脚本 |
| **public/** | 388 | 383 | 静态 SVG 资源 |
| **snap/** | 75 | 1 | Snap 打包配置 |

---

## 核心要点

- **约 301K 行 Go 代码**构成核心后端，分层清晰（routers -> services -> models -> modules）
- **`modules/git`**（12.6K 行）是最大的单一包——Git 集成是项目的核心
- **`models/migrations`**（9.6K 行，263 个文件）反映了大量的数据库结构演进历史
- **`routers/api`**（27.7K 行）与 **`routers/web`**（33.1K 行）比例约为 1:2，说明 Web UI 和 REST API 都有较大的接口规模
- **前端**约 48K 行，大量使用 Fomantic UI，自定义 JS/CSS 相对精简
- **`options/` 中 72K 行翻译文件**体现了强大的国际化支持
- **测试代码**（约 31.5K 行集成测试）以集成测试为主，单元测试较少

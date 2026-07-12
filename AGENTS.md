# 驾驭工程

本文件是人类贡献者和 AI 贡献者的入口。保持短、可扫读；细节下沉到
`docs/engineering/`、`openspec/` 和各一级目录的 `ROADMAP.md`。

## 从这里开始

1. 先读 `docs/engineering/PROJECT_MAP.md`，再定位控制点。
2. 运行 `openspec list`，检查是否已有相关变更。
3. 非简单变更先创建或更新 OpenSpec，再改代码。
4. 编辑前追踪完整行为路径，记录路由、上下文、service、model、模板或前端入口。
5. 编码前按 `docs/engineering/VERIFICATION.md` 选择验证方式。
6. 实现最小完整变更，同步更新 OpenSpec tasks。
7. 报告修改内容、检查结果、未执行检查和残余风险。

OpenSpec 工作流见 `docs/engineering/OPEN_SPEC.md`。

## 何时必须使用 OpenSpec

以下任一情况必须使用 OpenSpec：

- 用户可见行为、鉴权、安全、兼容性、性能或架构边界变化。
- HTTP 路由、请求/响应字段、状态码、Swagger 或公开 API 契约变化。
- 持久化、迁移、队列、定时任务、Webhook 或外部系统变化。
- 同时涉及 `routers`、`services`、`models`、`modules`、`templates`、
  `web_src` 中两个及以上层级。
- 验收标准尚不明确或不可测试。

纯措辞、注释或机械格式变更可以直接改；最终报告必须说明没有行为契约变化，并执行
相称检查。

## 编辑前先追踪

- Web：路由 -> 中间件/上下文 -> handler/form -> service -> model/module ->
  模板 -> JavaScript/CSS。
- API：路由 -> 鉴权/上下文 -> handler -> service -> model/module ->
  `modules/structs`/converter -> Swagger -> 客户端响应。
- 持久化：模型注册 -> 查询/事务 -> 迁移和多数据库行为 -> 调用方。
- 异步：生产者 -> 队列/任务载荷 -> 消费者 -> 重试/超时 -> 最终状态和通知。

跨层变更必须在 OpenSpec design 中记录这些控制点。

## 一级目录路线图

进入目录后先读对应 `ROADMAP.md`，再继续追踪具体文件。

| 目录 | 路线图 |
| --- | --- |
| `.codex` | `.codex/ROADMAP.md` |
| `.devcontainer` | `.devcontainer/ROADMAP.md` |
| `.gitea` | `.gitea/ROADMAP.md` |
| `.github` | `.github/ROADMAP.md` |
| `assets` | `assets/ROADMAP.md` |
| `build` | `build/ROADMAP.md` |
| `cmd` | `cmd/ROADMAP.md` |
| `contrib` | `contrib/ROADMAP.md` |
| `custom` | `custom/ROADMAP.md` |
| `docker` | `docker/ROADMAP.md` |
| `docs` | `docs/ROADMAP.md` |
| `models` | `models/ROADMAP.md` |
| `modules` | `modules/ROADMAP.md` |
| `openspec` | `openspec/ROADMAP.md` |
| `options` | `options/ROADMAP.md` |
| `public` | `public/ROADMAP.md` |
| `routers` | `routers/ROADMAP.md` |
| `services` | `services/ROADMAP.md` |
| `snap` | `snap/ROADMAP.md` |
| `templates` | `templates/ROADMAP.md` |
| `tests` | `tests/ROADMAP.md` |
| `tools` | `tools/ROADMAP.md` |
| `web_src` | `web_src/ROADMAP.md` |

不为 `.git`、`data`、`log`、`node_modules`、`.make_evidence`、`sqlite-log`
等本地运行状态、依赖或证据目录创建路线图；除非任务明确要求诊断这些本地状态，
不要读取或编辑它们。

## 仓库边界

事实来源目录包括 `cmd`、`routers`、`services`、`models`、`modules`、
`templates`、`web_src`、`tests`、`docs` 和 `tools`。

生成产物应通过源文件或生成器修改，不要直接手改：

- `public/assets/`
- `modules/*/bindata.go` 及对应 `.hash`
- `templates/swagger/v1_json.tmpl`
- `*.pb.go`

运行时配置和状态可能包含凭据。`custom/conf/`、`data/`、`log/`、`gitea`、
依赖目录和测试证据目录只能在明确诊断本地运行问题时接触，且不得在报告中泄露密钥。

## 实现规则

- 遵循现有包结构和命名；避免无关重构。
- 周围代码传递 `context.Context` 时，service 和 model 相关工作也必须继续传递。
- 编排逻辑放 `services`；HTTP 放 `routers`；持久化放 `models`；基础设施和公开结构体放
  `modules`。
- 默认保持 API 兼容。API 变化必须同步注解、`modules/structs`、Swagger 引用和测试。
- HTTP API 测试细则见 `docs/engineering/API_TESTS.md`。
- 模型变化必须处理注册、迁移、事务边界和受支持数据库。
- 模板变化必须检查 handler 数据和前端初始化器；不要修改生成资源。
- 用户行为变化时，必须在同一变更中新增或更新文档。
- 设计不再成立时，先更新 OpenSpec 工件，再扩大实现范围。

## 验证契约

选择覆盖变更行为的最窄检查；共享或高风险范围再扩大验证。快速入口：

```sh
./tools/harness/verify.sh harness
./tools/harness/verify.sh spec
./tools/harness/verify.sh go ./path/to/package/...
./tools/harness/verify.sh frontend path/to/file.test.js
```

标准广泛检查保持不变：

```sh
make test-backend
make test-frontend
make lint-backend
make lint-frontend
make build
```

不要把聚焦包测试描述为完整后端测试通过。集成和端到端环境要求见
`docs/engineering/VERIFICATION.md`。

## 完成定义

- 行为和非目标与 proposal/specs 一致。
- OpenSpec tasks 已勾选，并有实现或验证证据。
- 变化的控制点和公开契约已同步。
- 相关聚焦测试通过；共享或高风险变更按需执行更广检查。
- `openspec validate --all --strict --no-interactive` 通过。
- 最终报告列出已执行检查、未执行检查和残余风险。
- 已完成的 OpenSpec 变更已归档，使主规格描述当前行为。

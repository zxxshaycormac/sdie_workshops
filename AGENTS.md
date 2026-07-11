# 驾驭工程

本文件是人类贡献者和 AI 贡献者必须阅读的入口。内容应保持简短、可操作，
并与 `docs/engineering/` 中更深入的指南保持一致。

## 从这里开始

1. 在定位控制点之前，先阅读 `docs/engineering/PROJECT_MAP.md`。
2. 运行 `openspec list`，检查是否存在与当前请求相关的活动变更。
3. 对于非简单变更，必须先创建或更新 OpenSpec 变更，再修改代码。
4. 编辑前先追踪完整的行为路径。
5. 编码前根据 `docs/engineering/VERIFICATION.md` 选择验证方式。
6. 实现最小且完整的变更，并及时更新对应的 OpenSpec 任务状态。
7. 报告修改内容、通过的检查、未执行的检查以及任何限制。

有关 OpenSpec 工作流和示例，请阅读 `docs/engineering/OPEN_SPEC.md`。

## 何时必须使用 OpenSpec

当变更影响以下任一方面时，必须使用 OpenSpec：

- 用户可见行为或鉴权。
- HTTP 路由、请求或响应字段、状态码或 Swagger。
- 持久化、迁移、队列、定时任务、Webhook 或外部系统。
- 同时涉及 `routers`、`services`、`models`、`modules`、`templates` 或
  `web_src` 中的两个及以上层级。
- 安全性、兼容性、性能或架构边界。
- 验收标准尚不明确或不可测试的工作。

如果行为和契约不会发生变化，可以直接修改措辞、注释或机械性格式。
最终报告中必须说明这一判断，并仍然执行与变更相称的检查。

## 编辑前先追踪

不要停留在最近匹配到的符号。必须追踪实际涉及的完整路径：

- Web：路由 -> 中间件/上下文 -> handler/form -> service -> model/module ->
  模板 -> JavaScript/CSS。
- API：路由 -> 鉴权/上下文 -> handler -> service -> model/module -> converter
  或 `modules/structs` -> Swagger -> 客户端可见响应。
- 持久化：模型注册 -> 查询或事务 -> 迁移和多数据库行为 -> 调用方。
- 异步：生产者 -> 队列/任务载荷 -> 消费者 -> 重试/超时 -> 最终状态和通知。

对于跨层变更，必须在 OpenSpec 设计中记录这些控制点。

## 仓库边界

事实来源目录包括 `cmd`、`routers`、`services`、`models`、`modules`、
`templates`、`web_src`、`tests`、`docs` 和 `tools`。

以下内容属于生成产物；应修改其源文件或生成器，而不是直接修改产物：

- `public/assets/`
- `modules/*/bindata.go` 及对应的 `.hash` 文件
- `templates/swagger/v1_json.tmpl` 中生成的 Swagger 内容
- `*.pb.go` 等生成的 protobuf Go 文件

以下内容属于本地运行状态、依赖或构建产物。除非任务明确针对本地运行问题诊断，
否则不要读取或编辑它们：

- `custom/conf/`、`data/`、`log/`
- `node_modules/`、`.venv/`、`vendor/`
- `gitea`、`dist/`、测试二进制文件、覆盖率文件和 `.make_evidence/`

运行时配置可能包含凭据和签名密钥。绝不能在工件、日志或报告中引用这些内容。

## 实现规则

- 遵循现有的包结构和命名方式；避免无关重构。
- 如果周围代码会传递 `context.Context`，service 和 model 中的相关工作也必须继续传递。
- 编排逻辑放在 `services`；HTTP 相关逻辑放在 `routers`；持久化逻辑放在
  `models`；可复用基础设施和公开结构体放在 `modules`。
- 默认保持 API 兼容性。API 变更必须同步更新注解、`modules/structs`、
  Swagger 引用和聚焦测试。
- 模型变更必须处理注册、迁移、事务边界和受支持的数据库引擎。
- 模板行为变更必须检查对应的 handler 数据和前端初始化器。不要修改生成的资源。
- 用户行为发生变化时，必须在同一变更中新增或更新文档。
- 当实现结果表明设计不再成立时，不要在未说明的情况下扩大范围；
  必须先更新 OpenSpec 工件。

## HTTP API 测试契约

任何影响 HTTP API 行为的代码变更，都必须在同一变更中更新对应的 API 测试。
影响范围包括路由、handler、`modules/structs` 中的请求和响应类型、绑定或校验、
状态码、响应字段、鉴权或中间件、业务状态流转、数据库读写，以及对下游服务或
基础设施的调用。

新增或修改 HTTP 集成测试前，必须阅读 `tests/integration/README.md`，
并检查相邻的 `tests/integration/api_*_test.go` 文件。复用已有的测试服务器、
夹具数据库、登录 session、访问令牌、`MakeRequest` helper、
`tests.PrepareTestEnv(t)` 和清理模式。对于 `routers/api/` 下的聚焦 handler 测试，
应遵循使用 `contexttest.MockAPIContext`、`unittest.PrepareTestEnv` 等 helper 的
相邻测试。复用现有的队列、存储、邮件或外部服务测试抽象；不要引入临时拼凑的
monkeypatch 或新的 mock 规范。

修改接口前，必须在以下所有位置搜索现有覆盖：

- `routers/api/**/*_test.go` 中的聚焦 handler 和包测试。
- `tests/integration/api_*_test.go` 中经过实际路由的 HTTP 契约和流程测试。
- `services/`、`models/` 和 `modules/` 下与业务逻辑、持久化、转换和基础设施行为
  相关的测试。

搜索时使用路由路径、handler 名称、`modules/structs` 类型、service 或 model 符号，
以及业务关键字。如果受影响接口没有测试，必须新增测试。如果已有覆盖，必须更新其
请求、断言、夹具或测试依赖，使其符合预期的行为变化。

API 测试必须覆盖所变更行为的关键成功路径、关键失败路径、边界场景、鉴权和重要
副作用。如果接口存在必须遵循的业务顺序，并且该顺序属于契约的一部分，
应通过集成流程进行验证，而不是只进行孤立的 handler 调用。

当业务逻辑和 API 测试需要同时修改时，具备 subagent 能力的 AI 必须将 API 测试实现
或独立测试审查交给拥有独立上下文的 subagent。如果无法使用 subagent，最终报告中
必须说明，并提供独立审查清单，覆盖路由可达性、状态码和响应体断言、鉴权、夹具、
副作用以及清理。

实现后必须运行受影响的 API 测试。应诊断失败原因，而不是掩盖失败：如果是产品代码
或兼容性缺陷，应修复代码；只有当旧测试预期不再代表正确的预期行为时，才能更新测试。
不得仅为了使测试通过而删除断言、跳过测试或降低覆盖率。

最终报告必须列出每个新增或更新的 API 测试、准确的测试命令及结果，并说明任何相关
测试未执行的原因和残余风险。

## 验证契约

选择能够覆盖变更行为的最窄检查；对于共享或高风险范围，再扩大验证范围。

快速驾驭入口：

```sh
./tools/harness/verify.sh harness
./tools/harness/verify.sh spec
./tools/harness/verify.sh go ./path/to/package/...
./tools/harness/verify.sh frontend path/to/file.test.js
```

Makefile 中的标准广泛检查保持不变：

```sh
make test-backend
make test-frontend
make lint-backend
make lint-frontend
make build
```

集成测试和端到端测试需要 `docs/engineering/VERIFICATION.md` 中描述的环境。
绝不能把聚焦包测试描述为完整的后端测试通过。

## 完成定义

只有满足以下条件，变更才算完成：

- 行为和非目标与 proposal 和 specs 一致。
- 每个任务都已勾选，并有实现或验证证据。
- 发生变化的控制点和公开契约保持同步。
- 相关聚焦测试通过；共享或高风险变更已根据需要执行更广泛的检查。
- `openspec validate --all --strict --no-interactive` 通过。
- 最终报告列出已执行的检查、未执行的检查和残余风险。
- 已完成的 OpenSpec 变更已归档，使主规格描述当前行为。

## 工作区限制

建立驾驭工程时，此工作区不包含 `.git` 元数据。在恢复 Git 历史或初始化仓库之前，
不要假定 `git diff`、blame、干净工作区检查或基于 Git 的回滚证据可用。

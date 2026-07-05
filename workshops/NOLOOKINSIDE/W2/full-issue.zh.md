# Issue #32857：Actions 任务状态聚合不正确

> 来源：[go-gitea/gitea#32857](https://github.com/go-gitea/gitea/issues/32857)
> 修复 PR：[#32859](https://github.com/go-gitea/gitea/pull/32859)（**完成 Workshop 前请勿阅读**）
> 分类：`gitea-actions` · 类型：Bug · 维护类型：修正缺陷（行为变更）

---

## 概述

在 Gitea Actions 中，一个 workflow **run（运行）**包含多个 **job（任务）**。UI 需要聚合所有 job 的状态，以显示一个单一的整体运行状态。当前的聚合函数是一个**缺失了部分状态迁移的状态机**——它只处理了 8 种可能 job 状态中的 4 种，在常见场景下产生了错误的运行级状态。

## 用户可见的现象

| 场景（运行中的所有 job） | 期望的运行状态 | 实际的运行状态 |
|---|---|---|
| 全部 `skipped` | `skipped` | `success` |
| 一个或多个 `cancelled`（其余 `success`） | `cancelled` | `failure` |
| 全部 `cancelled` | `cancelled` | `failure` |
| 一个 `blocked`（其余 `success`） | `blocked` / `waiting` | `running` |
| 全部 `blocked` | `blocked` / `waiting` | `running` |

这些**不是**装饰性问题——CI 消费方（PR 检查、提交状态徽标、webhook 负载、UI 徽标）都读取聚合后的运行状态，因此用户会被误导，误以为运行实际成功、被中止或仍在等待。

## 根因（待通过 systematic-debugging 确认）

聚合逻辑位于 `models/actions/run_job.go:155` —— `aggregateJobStatus(jobs)`。它用**三个布尔标志**来概括 job 集合：

```go
allDone    := true   // 每个 job 都已到达终态？
allWaiting := true   // 每个 job 要么在等待，要么已完成？
hasFailure := false  // 是否有 job 失败或被取消？
```

### 为什么布尔值无法表达完整的状态空间

`Status` 枚举（`models/actions/status.go:15`）定义了 **8** 个值：

| # | Status | `IsDone()` | 是否有 runner 结果？ |
|---|---|---|---|
| 0 | `StatusUnknown` | 否 | 否 |
| 1 | `StatusSuccess` | 是 | 是 |
| 2 | `StatusFailure` | 是 | 是 |
| 3 | `StatusCancelled` | 是 | 是 |
| 4 | `StatusSkipped` | 是 | 是 |
| 5 | `StatusWaiting` | 否 | 否（调度器状态） |
| 6 | `StatusRunning` | 否 | 否（调度器状态） |
| 7 | `StatusBlocked` | 否 | 否（调度器状态） |

布尔归约器把这个 8 状态空间压缩成约 3 个桶，丢失了信息：

- **`cancelled` 被并入 `failure`** —— `hasFailure` 对两者都置位，因此被取消的运行报告为 `failure`。
- **`skipped` 不可见** —— 它算作"已完成"，但既不设置 `hasFailure`，也不设置任何正向信号，因此全部跳过的运行报告为 `success`。
- **`blocked` 与 `running` 无法区分** —— 两者既不是 `IsDone()` 也不是 `StatusWaiting`，因此都会把 `allDone=false` 和 `allWaiting=false` 翻转，落入默认的 `running` 分支。

### 前端的叠加问题

`web_src/js/components/ActionRunStatus.vue:37` 用**同一个**红色 `octicon-x-circle-fill` 图标渲染 `failure`、`cancelled` 和 `unknown`，因此即使后端返回了 `cancelled`，UI 也无法在视觉上把它和 `failure` 区分开。

## 影响范围

| 层 | 文件 | 关注点 |
|---|---|---|
| 后端（核心） | `models/actions/run_job.go` | `aggregateJobStatus` —— Bug 所在 |
| 后端（枚举） | `models/actions/status.go` | 8 个状态的真值来源（只读上下文） |
| 前端 | `web_src/js/components/ActionRunStatus.vue` | 状态图标渲染 |
| 前端（镜像） | `templates/repo/actions/status.tmpl` | Vue 组件的服务端孪生文件 |
| 测试（新增） | `models/actions/run_job_status_test.go` | 覆盖所有状态组合的表驱动测试 |

**受影响的调用路径：** `UpdateRunJob`（`run_job.go:94`）在每次 job 状态变更时都会重新计算 `run.Status = aggregateJobStatus(jobs)` 并持久化，因此错误的值会传播到每一个下游消费方。

## 期望行为（Spec —— 由学员自行提取）

一个正确的聚合器必须：

1. **显式枚举全部 8 个状态** —— 不能有隐式的 "else → running" 兜底。
2. **保留 `cancelled`、`failure`、`skipped` 各自独立的语义**。
3. **为混合状态的运行定义优先级顺序**。一种合理的排序：
   `success` > `failure` > `running` > `waiting` > `blocked` > `cancelled` > `skipped`
   （单个失败的 job 会污染整个运行；单个运行中的 job 意味着运行尚未结束）。
4. **处理全部 `skipped` 的情况** —— 运行并未真正执行，因此聚合结果应反映这一点，而不是伪装成 `success`。
5. **在 UI 中让 `cancelled` 有明显的视觉区分**（例如用 `octicon-stop` 而非失败图标）。

具体的优先级是一个设计决策 —— 学员必须论证自己的选择，而不是照抄。

## 复现步骤

1. 启动一个本地 Gitea 实例，启用 Actions 并注册一个 runner。
2. 创建一个**所有** job 都会被跳过的 workflow，例如给 job 加一个求值为假的 `if:`，或者让 job 的 `needs:` 依赖被跳过。
3. 触发该 workflow，观察运行列表 / 运行详情页。
4. 聚合后的运行状态会显示为 **`success`**（在某些边缘情况下为 `running`），而不是 **`skipped`**。
5. 换个实验：在运行中取消一个 job → 聚合状态显示为 `failure`。

## 验收标准

- [ ] 当所有 job 为 `skipped` 时，`aggregateJobStatus` 返回 `skipped`。
- [ ] 当任意 job 为 `cancelled` 且没有 job 失败时，`aggregateJobStatus` 返回 `cancelled`（而非 `failure`）。
- [ ] `aggregateJobStatus` 不再把 `blocked` 压缩成 `running`。
- [ ] 一个**表驱动**的 Go 测试覆盖了每种状态组合（单状态和有代表性的多状态混合）—— 在修复**之前**编写（Red 阶段）。
- [ ] 前端用与 `failure` 不同的图标渲染 `cancelled`。
- [ ] `models/actions/` 下所有既有测试仍然通过。

## Workshop 映射（逆向 SDIE：E → S → I → E）

| 阶段 | 活动 | Skill |
|---|---|---|
| **(E)val —— 逆向启动** | 复现现象 → 枚举全部 8 个状态 → 阅读 `aggregateJobStatus` → 发现缺口 | `systematic-debugging` |
| **(S)pec —— 逆向提取** | 从状态枚举推导出完整的聚合规则（优先级表） | — |
| **(I)mpl —— 修复** | Red：编写失败的表驱动测试 · Green：重写聚合器 · Refactor：简化 | `test-driven-development` |
| **(E)val —— 回归** | 全部测试通过 + cancelled 图标有区分 + 原始场景已修复 | — |

## 给学员的约束

- **不要**阅读 PR #32859 或其 diff。从 Bug 自行推导修复。
- **不要**猜测解决方案 —— 遵循 ReAct 循环：现象 → 假设 → 验证 → 根因。
- **不要**跳过 Red 阶段。先写失败的测试，它是证明修复有效、并防止回归的契约。
- **要**在写代码前用表格枚举每一种状态组合。MECE 是你没有漏掉某个状态迁移的自检手段。

## 参考

- 状态枚举：`models/actions/status.go:15`
- 出错的聚合器：`models/actions/run_job.go:155`
- 调用点：`models/actions/run_job.go:140`（在 `UpdateRunJob` 内）
- 前端图标映射：`web_src/js/components/ActionRunStatus.vue:32`
- Gitea Actions 状态文档：https://docs.gitea.com/usage/actions/

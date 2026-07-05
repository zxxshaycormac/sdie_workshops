# E2S：从现象到规约（Evaluation → Specification）

> 本文档记录 W2 练习的逆向 SDIE 分析：把 Issue #32857 的**用户可见现象（Eval）**，
> 反向映射到 `openspec/specs/cicd-automation/spec.md` 中的**规约需求（Spec）**，
> 从而定位 Bug 在规约层面的**根因**与**缺口**。

## 元信息

| 属性 | 值 |
|---|---|
| 来源 Issue | [#32857](https://github.com/go-gitea/gitea/issues/32857) — Actions 运行状态聚合错误 |
| 对应规约 | `openspec/specs/cicd-automation/spec.md`（Domain 05 CI/CD & Automation，第 1 节 Gitea Actions、第 7 节 Commit Status） |
| SDIE 方向 | 逆向 **E → S**（现象 → 规约缺口） |
| 维护类型 | 修正缺陷（改行为） |

---

## 1. 根因对应的 Spec 需求（最关键）

**`CICD-01-104`**（`spec.md:27`）——这是驱动错误行为的"罪魁需求"：

> When a workflow run completes, the system shall update the run status to
> **success or failure** based on job outcomes.

**问题**：这条需求把 run 完成后的状态收敛成**二元（success / failure）**，
完全没有给 `skipped`、`cancelled` 留分支。Issue 中"全部 skipped 显示 Success"、
"有 job 被取消显示 Failure"这两个现象，都源自这条规约的不完备。

> 系统性洞察：Bug 不在代码实现，而在**规约本身就漏掉了状态分支**。
> 代码只是忠实地实现了这条不完备的规约。

---

## 2. 状态生命周期：规约已修正，但聚合规则未同步

**`CICD-01-003`**（`spec.md:15`）——状态生命周期枚举（**已在规约校正中修订**）：

> `unknown, waiting (pending), running, success, failure, cancelled, skipped, and blocked.`

> **重要变化**：在 W2 首次编写时，此需求只列了 6 个状态（missing `skipped` 和 `blocked`）。
> 规约校正 pass 已将其补全为完整的 8 个状态。**但是**——
> `CICD-01-104`（聚合需求）**并未同步修订**，仍然只说"success or failure"。

**当前缺口（修正后仍然存在的）**：
- `CICD-01-003` 的生命周期**已包含** `skipped`、`cancelled`、`blocked` ✓
- 但 `CICD-01-104` 的聚合规则**仍然把它们丢弃** ✗
- 生命周期枚举和聚合规则之间**存在语义断裂**：枚举说"有 8 种状态"，
  聚合说"只产出 2 种结果"。代码忠实实现了聚合规则（`aggregateJobStatus`），
  所以即使枚举已完备，用户看到的状态仍然是错误的。

> 更微妙的洞察：**规约的局部修正是必要的但非充分的**。
> 修正 `CICD-01-003` 补全了状态空间，但如果不同步修正 `CICD-01-104`，
  Bug 依然存在。这正是跨需求一致性（S × S）的价值所在。

---

## 3. 支撑需求（交叉影响面）

Issue 警告错误状态会"蔓延"到提交检查与 Webhook。规约恰好覆盖了这些下游面：

| Spec 需求 | 位置 | 与本 Issue 的关系 |
|---|---|---|
| **`CICD-07-003`** 计算 commit 的 combined status | `spec.md:269` | 同样是"聚合"逻辑，受同一类不完备影响 |
| **`CICD-07-102`** combined = "任意 context 失败即 failure" | `spec.md:275` | 这种全有或全无的聚合，会把 skipped/cancelled 的 run 误报 |
| **`BR-05-007`** combined status failure if any context fails | `spec.md:335` | 同上的业务规则，下游"绿色对勾"门禁受污染 |
| **`CICD-07-202`** status 显示在 PR UI | `spec.md:281` | 即 Issue 所说的"PR 检查中的提交状态" |
| **`CICD-05-101`** 事件触发 webhook 投递 | `spec.md:209` | 即 Issue 所说的"webhook 负载" |
| **`CICD-01-102`** 评估 job 依赖并调度 | `spec.md:25` | 与"Blocked 等待依赖/并发槽"现象相关 |
| 边缘用例：无 runner 时 job 保持 pending | `spec.md` Edge Cases | Blocked 场景的规约出处 |

---

## 4. 验收标准 ↔ Spec 缺口（1:1 对照）

| Issue 验收标准 | 对应的 Spec 缺口 |
|---|---|
| 全部 skipped 的 run → **Skipped** | `skipped` 已在生命周期枚举中（`CICD-01-003` ✓）；但 `CICD-01-104` 聚合仍只允许 success/failure ✗ |
| 有 job 被取消 → **Cancelled**（非 Failure） | `CICD-01-105` 只说"取消所有 job"，未定义 cancelled 如何**回滚**到 run 状态 |
| 阻塞的 run 不伪装成 Running | `blocked` 已在生命周期枚举中（`CICD-01-003` ✓）；但聚合规则未定义 blocked 如何映射到 run 状态 ✗ |
| Cancelled 在 UI 上与 Failure 视觉可区分 | 规约无任何"状态图标/视觉区分"需求 |

---

## 5. 结论：这是一个 (S)pec 层面的缺陷

逆向 E→S 分析的结论：

1. **根因在规约层**，不在代码实现层。
   `CICD-01-104` 的二元聚合是 Bug 的真正源头——代码忠实执行了它。

2. **状态空间已完备，但聚合规则未同步（跨需求一致性缺口）**。
   `CICD-01-003` 的生命周期枚举已在规约校正中补全为 8 个状态（含 `skipped`、`blocked`）。
   但 `CICD-01-104` 仍然只说"success or failure"——枚举和聚合之间存在语义断裂。
   **局部修正一个需求是不够的，必须检查所有相关需求的一致性。**

3. **修复的规约动作**：任何正确的代码修复都应同步修订以下需求：
   - **修订 `CICD-01-104`**：把"success or failure"扩展为完备的聚合规则
     （含 skipped / cancelled 分支，并定义优先级）。
   - **新增一条需求**：明确定义 **job 状态 → run 状态的聚合优先级规则**
     （如：有失败=Failure；否则有取消=Cancelled；否则全 skipped=Skipped；否则=Success）。

4. **下游联动**：Commit Status 的聚合（`CICD-07-003/102`、`BR-05-007`）
   是同一类不完备，修复时应一并审视，保证 UI、commit 检查、webhook 三处一致。

> 这正是 W2 "逆向 E→S" 的核心练习：**用 Eval 的现象反推 Spec 的缺口**，
> 而不是直接去改代码。规约对了，代码自然就对。

---

# 第二部分：(E)val 维度 —— 需要补充的评估

> 确认了 Spec 缺口（第一部分）之后，正向补齐 **(E)val 层**：现有哪些测试、
> 缺哪些测试、用什么评估矩阵覆盖 Issue 的 5 条验收标准。
> 对应 SDIE 的 **Verification（S × E）** 与 **Validation（端到端）**。

## 6. 现有评估基线（关键缺口）

先勘察代码与测试现状，定位评估空洞：

| 评估面 | 现状 | 位置 |
|---|---|---|
| `aggregateJobStatus`（Bug 所在） | **零单元测试** | `models/actions/run_job.go:155` |
| `Status` 助手函数（`IsDone`/`IsRunning`/`IsCancelled`…） | 无测试 | `models/actions/status.go:49` 等 |
| 运行状态传播链路（`UpdateRunJobs → UpdateRun`） | 无集成测试 | `models/actions/run_job.go:140` |
| run list / detail API 的 status | 无集成测试 | `tests/integration/actions_trigger_test.go` 只测触发 |
| 前端 `ActionRunStatus.vue` | **已支持** skipped/cancelled/blocked 图标 | `web_src/js/components/ActionRunStatus.vue:3,33` |

讽刺的洞察：状态枚举（`status.go:19-23`）**已定义** `Skipped`/`Cancelled`/`Blocked`，
前端**已能渲染**这些图标，`num_closed_action_runs`（`run.go:184-185`）**也已计入**
cancelled/skipped——唯独 `aggregateJobStatus` 从不产出这些状态，而它恰恰零测试。
**评估的最大空洞就在聚合函数这一层。**

## 7. 需要补充的评估（按层次）

### 7.1 单元层 —— Verification（S × E），核心 Red 测试

**`TestAggregateJobStatus`** 表驱动测试，MECE 覆盖 job 状态组合。
这是 TDD 的 Red 起点，直接复现 Issue 的四个现象：

| job 状态组合 | 期望 run 状态 | 当前实际（bug） |
|---|---|---|
| `[success]` / `[success, success]` | Success | ✓ |
| `[failure]` / `[success, failure]` | Failure | ✓ |
| `[skipped]` / `[skipped, skipped]` | **Skipped** | ✗ Success |
| `[cancelled]`（无 failure） | **Cancelled** | ✗ Failure |
| `[skipped, cancelled]` | **Cancelled** | ✗ Failure |
| `[blocked]` / `[success, blocked]` | **Blocked** | ✗ Running |
| `[waiting]` / `[waiting, waiting]` | Waiting | ✓ |
| `[running]` / `[success, running]` | Running | ✓ |

> 优先级规则需在设计层定死（无唯一正确答案，但要能自圆其说）：
> 中间态优先 `Failure > Cancelled > Skipped > Success`，未完成态用 `Blocked > Waiting > Running`。

补充：`TestStatusIsDone/IsRunning/IsCancelled/IsSkipped/IsBlocked` —— 验证状态分类助手，
这是聚合逻辑的语义基础。

### 7.2 集成层 —— Validation（端到端真实链路）

`tests/integration/actions_status_aggregation_test.go`：通过完整链路
（trigger run → 改 job 状态 → `UpdateRunJobs` → 查 run）验证，而非只调纯函数：

- 制造 all-skipped run，断言 `ActionRun.Status == StatusSkipped`
- 取消一个 job，断言 run 为 `Cancelled`
- 建一个 blocked job，断言 run **不为** `Running`

### 7.3 下游传播层 —— Cross-cutting（I × E）

Issue 警告"状态会蔓延"的几个出口，每个都要有验证：

- **API**：`GET /api/v1/repos/{owner}/{repo}/actions/runs/{run_id}` 返回的 `status` 字段与聚合一致。
- **Commit Status**（`CICD-07-003/102`、`BR-05-007`）：run 为 cancelled/skipped 时，commit check 不显示为 success/failure。
- **Webhook payload**（`CICD-05-101`）：run 完成事件的负载携带正确 status。

### 7.4 前端层 —— 验收标准 #4（视觉可区分）

`web_src` 测试断言：`cancelled` 渲染的图标**区别于** Failure 的红叉。
前端已具备渲染能力（`ActionRunStatus.vue:33`），只缺回归测试锁住"图标不退化"。

### 7.5 回归层 —— 验收标准 #5

现有真正 success/failure 的 run 不被新逻辑误判——7.1 表里前两行即回归保护，
应**明确标注为 regression**，不与新增 case 混在一起（避免被误删）。

## 8. 验收标准 ↔ 评估层 映射

| Issue 验收标准 | 覆盖它的评估层 |
|---|---|
| 全部 skipped → Skipped | 7.1（单元 Red）+ 7.2（集成） |
| 有取消 → Cancelled（非 Failure） | 7.1 + 7.2 + 7.3（API/webhook） |
| 阻塞 run 不伪装 Running | 7.1 + 7.2 |
| Cancelled 与 Failure 视觉可区分 | 7.4（前端） |
| 现有 success/failure 仍正确 | 7.5（回归） |

> **一句话**：最大的 evaluation 缺口是 `aggregateJobStatus` 的 MECE 表驱动单元测试——
> 它既是 Red 测试暴露 Bug，也是修复后 Verification 的合约。
> 其余四层（集成 / API·webhook / 前端 / 回归）补齐后，正好覆盖 Issue 全部 5 条验收标准。

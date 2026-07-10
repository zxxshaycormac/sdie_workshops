## 用户报告的现象

| 运行中实际发生的情况 | 用户期望看到 | 用户实际看到 |
|---|---|---|
| 所有 job 都被**跳过**（例如 `if:` 条件全为假） | Skipped | Success |
| 运行中用户**取消**了某个 job | Cancelled | Failure |
| 某个 job **被阻塞**等待中 | Blocked / Waiting | Running |
| 跳过与取消的 job 混合 | Cancelled（至少不是 Success） | Success 或 Failure |

错误的状态会出现在：运行列表 / 运行详情徽标、commit 提交状态、PR 检查、
以及发送给外部系统的 webhook 负载中。

## 复现步骤

1. 启用 Gitea Actions 并注册一个 runner。
2. 创建一个**所有** job 都会被跳过的 workflow（例如给每个 job 加一个
   求值为假的 `if:`，或让 `needs:` 依赖一个被跳过的 job）。
3. 触发该 workflow，观察运行列表 / 运行详情页显示的整体状态。
4. 换个实验：在多 job 运行中**取消其中一个 job**，再观察整体状态。

## 验收标准

- [ ] 所有 job 被跳过的 run 显示为 **Skipped**，而非 Success。
- [ ] 含有被取消 job 的 run 显示为 **Cancelled**，而非 Failure（无 job 真正失败时）。
- [ ] 阻塞状态的 run 不再伪装成 Running。
- [ ] **Cancelled** 在 UI 上与 **Failure** 视觉可区分（不同图标）。
- [ ] 现有真正成功 / 失败的 run 仍被正确报告。

# Gitea 1.23 Feature Candidates

> Criteria: (1) Appears in v1.23.0 release notes FEATURES section, (2) Has a linked GitHub issue, (3) PR was properly merged into v1.23.0.

## Summary

- **Total FEATURES in release notes:** 25
- **Meeting all 3 criteria:** 19
- **Missing linked issue:** 9 (2 are base-module PRs whose companion PRs do have issues)

---

## Features Meeting All Criteria (19)

| # | Feature | PR | Issue | Issue Title | Type | Topic | Files Changed | Lines Changed |
|---|---------|-----|-------|-------------|------|-------|---------------|---------------|
| 1 | Support compression for Actions logs & enable by default | [#32013](https://github.com/go-gitea/gitea/pull/32013) | [#31801](https://github.com/go-gitea/gitea/issues/31801) | Enable compression for Actions logs by default | enhancement | gitea-actions | 2 | +3 -3 |
| 2 | Included tag search capabilities | [#32045](https://github.com/go-gitea/gitea/pull/32045) | [#31998](https://github.com/go-gitea/gitea/issues/31998) | Add searching capabilities to the tags page | feature | — | 4 | +33 -7 |
| 3 | Allow to fork repository into the same owner | [#32819](https://github.com/go-gitea/gitea/pull/32819) | [#22882](https://github.com/go-gitea/gitea/issues/22882) | Fork your own repository | feature | — | 5 | +44 -5 |
| 4 | Allow cropping an avatar before setting it | [#32565](https://github.com/go-gitea/gitea/pull/32565) | [#31990](https://github.com/go-gitea/gitea/issues/31990) | Unable to display complete avatar | feature | — | 12 | +80 -9 |
| 5 | Allow to disable the password-based login (sign-in) form | [#32687](https://github.com/go-gitea/gitea/pull/32687) | [#7633](https://github.com/go-gitea/gitea/issues/7633) | Is it possible to choose default auth source? | feature | — | 7 | +73 -48 |
| 6 | Use env GITEA_RUNNER_REGISTRATION_TOKEN as global runner token | [#32946](https://github.com/go-gitea/gitea/pull/32946) | [#23703](https://github.com/go-gitea/gitea/issues/23703) | Improve Config Management/Stateless Runner Deploy Workflows | feature | gitea-actions | 8 | +152 -18 |
| 7 | Add Passkey login support | [#31504](https://github.com/go-gitea/gitea/pull/31504) | [#22015](https://github.com/go-gitea/gitea/issues/22015) | Add support for passkeys (WebAuthn as primary authentication) | feature | — | 8 | +184 -11 |
| 8 | Suggestions for issues | [#32327](https://github.com/go-gitea/gitea/pull/32327) | [#16872](https://github.com/go-gitea/gitea/issues/16872) | autocompletion / suggestion for # (reference issues/pulls) | feature | — | 9 | +202 -48 |
| 9 | Introduce globallock as distributed locks | [#31813](https://github.com/go-gitea/gitea/pull/31813) | [#19620](https://github.com/go-gitea/gitea/issues/19620) | replace sync module | feature | — | 13 | +185 -107 |
| 10 | Add option to filter board cards by labels and assignees | [#31999](https://github.com/go-gitea/gitea/pull/31999) | [#21846](https://github.com/go-gitea/gitea/issues/21846) | Project board card filtering | feature | projects | 14 | +325 -33 |
| 11 | Rearrange Clone Panel | [#31142](https://github.com/go-gitea/gitea/pull/31142) | [#23202](https://github.com/go-gitea/gitea/issues/23202) | Clone Button Rearrangement | feature | ui | 19 | +191 -195 |
| 12 | Support "merge upstream branch" (Sync fork) | [#32741](https://github.com/go-gitea/gitea/pull/32741) | [#20880](https://github.com/go-gitea/gitea/issues/20880) | Single-action web UI for a fork to fast-forward pull from origin (Sync fork) | feature | — | 10 | +323 -136 |
| 13 | Issue time estimate, meaningful time tracking | [#23113](https://github.com/go-gitea/gitea/pull/23113) | [#23112](https://github.com/go-gitea/gitea/issues/23112) | Issue time estimate, meaningful time tracking | feature | — | 21 | +390 -164 |
| 14 | Enhancing Gitea OAuth2 Provider with Granular Scopes | [#32573](https://github.com/go-gitea/gitea/pull/32573) | [#31609](https://github.com/go-gitea/gitea/issues/31609) | Enhancing Gitea OAuth2 Provider with Granular Scopes for Resource Access | feature | — | 8 | +537 -18 |
| 15 | Actions support workflow dispatch event | [#28163](https://github.com/go-gitea/gitea/pull/28163) | [#23668](https://github.com/go-gitea/gitea/issues/23668) | Actions - Manually trigger a workflow/action | feature | gitea-actions | 10 | +580 -17 |
| 16 | Add reviewers' selection to new pull request | [#32403](https://github.com/go-gitea/gitea/pull/32403) | [#26289](https://github.com/go-gitea/gitea/issues/26289) | Assign reviewers on pull request creation | feature | — | 26 | +500 -268 |
| 17 | Support repo license | [#24872](https://github.com/go-gitea/gitea/pull/24872) | [#278](https://github.com/go-gitea/gitea/issues/278) | Display a License tab | feature | — | 47 | +906 -22 |
| 18 | Add pure SSH LFS support | [#31516](https://github.com/go-gitea/gitea/pull/31516) | [#17554](https://github.com/go-gitea/gitea/issues/17554) | Support LFS purely over SSH protocol | feature | — | 13 | +945 -53 |
| 19 | Add Arch package registry | [#32692](https://github.com/go-gitea/gitea/pull/32692) | [#25037](https://github.com/go-gitea/gitea/issues/25037) | Arch linux packages | feature | packages | 43 | +1687 -91 |

## Features NOT Meeting All Criteria (9)

All PRs were merged into v1.23.0, but no linked GitHub issue was found.

| # | Feature | PR | Note |
|---|---------|-----|------|
| 1 | Support quote selected comments to reply | [#32431](https://github.com/go-gitea/gitea/pull/32431) | No linked issue |
| 2 | Add priority to the protected branch | [#32286](https://github.com/go-gitea/gitea/pull/32286) | No linked issue |
| 3 | Add automatic light/dark option for the colorblind theme | [#31997](https://github.com/go-gitea/gitea/pull/31997) | No linked issue |
| 4 | Support migration from AWS CodeCommit | [#31981](https://github.com/go-gitea/gitea/pull/31981) | No linked issue |
| 5 | GitHub like repo home page | [#32213](https://github.com/go-gitea/gitea/pull/32213) | #27931 is a PR, not an issue |
| 6 | Tweak repo sidebar | [#32847](https://github.com/go-gitea/gitea/pull/32847) | No linked issue, companion to #32213 |
| 7 | Update i18n.go - Language Picker | [#32933](https://github.com/go-gitea/gitea/pull/32933) / [#32935](https://github.com/go-gitea/gitea/pull/32935) | No linked issue |
| 8 | Introduce globallock (base module) | [#31908](https://github.com/go-gitea/gitea/pull/31908) | Companion to #31813 which has issue #19620 |
| 9 | Support compression for Actions logs (base module) | [#31761](https://github.com/go-gitea/gitea/pull/31761) | Companion to #32013 which has issue #31801 |


## Issue #2: Included tag search capabilities

**PR [#32045](https://github.com/go-gitea/gitea/pull/32045)** · Issue [#31998](https://github.com/go-gitea/gitea/issues/31998) · 4 files · +33 -7

### Feature Description

为仓库的 Tags 页面添加搜索/过滤功能。用户可以在 `/{org}/{repo}/tags` 页面通过关键词搜索 tag 名称。与已有的 branches/commits 搜索功能保持一致的 UX 模式。

### Affected Layers

| Layer | Files | Key Changes |
|---|---|---|
| **models/** | 1 | `repo/release.go` +6 行（添加按关键词过滤 tag 的方法） |
| **routers/web/** | 1 | `repo/release.go` +13/-3（接收 URL query 参数 `q`，传给 model 层） |
| **templates/** | 1 | `tag/list.tmpl` +12/-4（搜索输入框 + 空结果提示 UI） |
| **options/** | 1 | i18n +2 行（搜索框 placeholder 文案） |

### Change Breakdown

- **New files**: 0
- **Modified files**: 4（每个文件改动 2-13 行）
- **无新依赖、无 DB 变更、无 migration**

### Difficulty: 1/5（最低，最适合入门）

**Reasoning**:
- 仅 4 个文件，总计 +33 行
- 遵循已有模式（branches/commits 搜索功能的复制）
- 无 DB 变更、无新依赖、无跨层复杂度
- 核心逻辑：model 加一个过滤方法 → router 传参 → template 加搜索框
- **30-45 分钟即可完成核心实现**

### Workshop Risks

- 几乎没有风险——改动极小，边界清晰
- 唯一注意点：需确认 branches/commits 的搜索模式作为参考

### Recommended Workshop Fit

**W1（增加特性 / 正向 SDIE）— 最佳入门选择**

与 #15 相比：
- #2 更适合首次练习 SDIE 正向流程（改动小、成功快、信心建立）
- #15 更适合有经验后的深入练习（改动多、涉及 YAML 解析、workflow 架构）
- 建议 W1 提供两个选项：#2（入门）和 #15（进阶），学员根据进度自选

---

## Issue #17: Support repo license

**PR [#24872](https://github.com/go-gitea/gitea/pull/24872)** · Issue [#278](https://github.com/go-gitea/gitea/issues/278) · 47 files · +906 -22

### Feature Description

在仓库页面展示 License 信息（Tab 形式），支持单 License 和多 License 识别。使用 `google/licensecheck` 库检测 License 文件，将结果存入新的 `repo_license` 表。支持 push/create/mirror/切换默认分支时自动更新 License 信息。提供 API 查询接口。

### Affected Layers

| Layer | Files | Key Changes |
|---|---|---|
| **models/** | 5 | 新增 `repo_license.go`（120 行）+ DB migration v305 + fixture |
| **services/** | 10 | 新增 `license.go`（205 行）+ `license_test.go`（73 行）+ 改动 create/fork/delete/migrate/branch/mirror |
| **routers/api/** | 4 | 新增 `repo/license.go`（51 行）+ swagger + api.go 注册路由 |
| **routers/web/** | 5 | view/commit/branch/release 微调，传入 License 数据 |
| **routers/private/** | 2 | hook_post_receive + default_branch 触发 License 更新 |
| **modules/** | 3 | structs/repo.go 新增字段 + repository/license.go 微调 + structs |
| **build/** | 3 | 新增 license alias generator + test |
| **templates/** | 2 | sub_menu.tmpl（License Tab）+ swagger JSON |
| **tests/** | 2 | 集成测试 80 行 + admin test 微调 |
| **options/** | 2 | i18n + license alias JSON |
| **assets/** | 1 | go-licenses.json |
| **go.mod/sum** | 1 | 新增 `google/licensecheck` 依赖 |

### Change Breakdown

- **New files**: 10（model、service、API handler、build tool、test、fixture、migration、template、alias JSON）
- **Modified files**: 37（大部分是 1-20 行的小改动——在 create/fork/migrate 等流程中调用 License 更新）

### Difficulty: 4/5（高）

**Reasoning**:
- 跨几乎所有层（models → services → routers → templates → build → migration）
- 需要新增外部依赖（`google/licensecheck`）
- 10 个 service 层文件需要插入 License 更新逻辑（create/fork/migrate/mirror/branch/hook）
- 触发逻辑复杂（push/create/mirror/切换默认分支 4 种触发方式）
- 但核心逻辑相对自包含（`services/repository/license.go` 是核心，其他是调用点）

### Workshop Risks

- **范围过大**：47 个文件对 2.5 小时工作坊不现实。需要大幅拆分为子任务
- **外部依赖**：需要引入 `google/licensecheck`，可能涉及 go.mod 冲突
- **测试依赖**：需要 git test data（`tests/gitea-repositories-meta/`）配合
- **建议**：适合拆分为多天任务，或只做核心部分（model + service + API，跳过 UI 和触发逻辑）

### Recommended Workshop Fit

**W1（增加特性 / 正向 SDIE）— 但需要大幅缩减范围**

这个 feature 是典型的"正向新增"——有明确的 Issue、完整的 PR 参考、清晰的边界。但 47 个文件不适合单次工作坊全部完成。建议：

- **核心任务**（适合 W1）：`models/repo/license.go` + `services/repository/license.go` + `routers/api/v1/repo/license.go` + migration——实现"检测并返回 License 信息"的核心链路
- **扩展任务**（可分配到其他天或课后）：触发逻辑（hook/branch/mirror）、UI Tab、build tool

---

## Issue #13: Issue time estimate, meaningful time tracking

**PR [#23113](https://github.com/go-gitea/gitea/pull/23113)** · Issue [#23112](https://github.com/go-gitea/gitea/issues/23112) · 21 files · +390 -164

### Feature Description

为 Issue 添加时间估算功能（Jira 风格的 "1w 3d 15h 30m" 格式）。改进时间追踪显示（精确显示每人已记录时间、总时间），统一时间日志评论风格，移除无意义的秒数显示。支持未来对时间追踪评论的国际化。

### Affected Layers

| Layer | Files | Key Changes |
|---|---|---|
| **models/** | 4 | issue.go 新增 28 行（estimate 字段 + DB 操作）+ comment.go + migration v311 |
| **modules/** | 3 | 新增 `util/time_str.go`（85 行）+ test + templates/helper.go 时间格式化 |
| **routers/web/** | 3 | issue_stopwatch + issue_timetrack（39 行改动）+ web.go 路由注册 |
| **services/** | 3 | convert/issue_comment + issue/issue + forms |
| **templates/** | 3 | stopwatch_timetracker（大幅重写 +64/-46）+ comments + comments_delete_time |
| **web_src/js/** | 2 | repo-issue.ts（-31 行移除旧 JS）+ index.ts |
| **options/** | 1 | i18n 17 行新增 |
| **tests/** | 2 | timetracking_test + html_helper |

### Change Breakdown

- **New files**: 2（`modules/util/time_str.go` + `models/migrations/v1_23/v311.go`）
- **Modified files**: 19（大部分是 1-40 行的定向修改）
- **Deleted**: JS 中移除 31 行旧的时间输入逻辑（由 Go 模板替代）

### Difficulty: 3/5（中等）

**Reasoning**:
- 范围适中（21 文件），核心逻辑集中在 3 个新/改动文件
- `modules/util/time_str.go` 是独立的时间解析工具，可与主功能分开实现
- DB migration 简单（只加一个 `estimate` 字段）
- 改动模式清晰：model 加字段 → service 加逻辑 → router 加接口 → template 改 UI
- 前端改动较少（JS 删 31 行，主要靠 Go template）
- **增量友好**：可以先做 model + util + API，再做 UI

### Workshop Risks

- 时间解析逻辑（`1w 3d 15h 30m`）需要仔细测试边界情况
- Template 重写幅度较大（stopwatch_timetracker +64/-46），需要仔细对照
- 集成测试已有（timetracking_test），但需确认在 v1.22 baseline 上可运行

### Recommended Workshop Fit

**W3（系统优化 / 重构）或 W1（增加特性）**

- 作为 **W1（增加特性）**：这是一个标准的"正向新增"——在已有 Issue 模型上加 estimate 字段，逻辑清晰，适合练习 SDIE 正向流程
- 作为 **W3（重构）**：PR 中包含了对现有时间追踪 UI 的重构（统一评论风格、移除 JS 逻辑到 Go template），适合练习"不改行为的重构"
- **建议**：放在 W1，因为它以新增功能为主（estimate 字段），重构部分是附带改进

---

## Issue #15: Actions support workflow dispatch event

**PR [#28163](https://github.com/go-gitea/gitea/pull/28163)** · Issue [#23668](https://github.com/go-gitea/gitea/issues/23668) · 10 files · +580 -17

### Feature Description

支持 Gitea Actions 的 `workflow_dispatch` 事件——允许用户手动触发 workflow 运行。提供表单让用户选择目标分支/标签，支持在 YAML 中配置 `inputs` 参数（choice/boolean/number/environment/string 类型）。行为与 GitHub Actions 保持一致（只触发默认分支上的 workflow）。

### Affected Layers

| Layer | Files | Key Changes |
|---|---|---|
| **routers/web/repo/actions/** | 3 | actions.go（+136 行，核心：解析 workflow_dispatch 配置 + 表单渲染 + 触发逻辑）+ view.go（+169 行）+ actions_test.go（新增 156 行测试） |
| **templates/** | 2 | list.tmpl（"Run workflow" 按钮）+ 新增 workflow_dispatch.tmpl（78 行表单） |
| **modules/structs/** | 1 | hook.go 新增 WorkflowDispatchPayload 结构体 |
| **web_src/js/** | 2 | repo-legacy.ts（+13/-7 表单交互）+ index.ts |
| **routers/web/** | 1 | web.go 路由注册 |
| **options/** | 1 | i18n 6 行 |

### Change Breakdown

- **New files**: 2（`actions_test.go` 156 行 + `workflow_dispatch.tmpl` 78 行）
- **Modified files**: 8（核心改动集中在 actions.go 和 view.go）
- **总计 3 个目录承担了 95% 的改动**：`routers/web/repo/actions/` + `templates/repo/actions/` + `web_src/js/`

### Difficulty: 2/5（较低，最适合工作坊）

**Reasoning**:
- **高自包含性**：改动几乎全部在 `routers/web/repo/actions/` 目录内
- **清晰的功能边界**：解析 YAML workflow_dispatch 配置 → 渲染表单 → 接收表单提交 → 触发 runner
- **无 DB 变更**：不涉及 migration、model 变更
- **无跨层依赖**：不改 models/services，只在 router 层处理
- **有完整测试**：新增 156 行测试，覆盖各种 input 类型
- **增量友好**：可分步实现——先支持基本触发，再加 input 表单，最后加验证
- **类比清晰**：GitHub Actions 的 workflow_dispatch 是学员可能熟悉的功能

### Workshop Risks

- 需要理解 Gitea Actions 的架构（act runner、workflow YAML 解析）
- YAML input 类型较多（choice/boolean/number/environment/string），实现细节琐碎
- `view.go` 的 169 行改动需要理解现有 Actions 视图的渲染逻辑
- **但**：这些风险都被"高自包含性"和"无跨层依赖"抵消

### Recommended Workshop Fit

**W1（增加特性 / 正向 SDIE）— 最推荐**

这是三个 feature 中最适合工作坊的：
- 改动集中在一个目录
- 无 DB 变更、无跨层依赖
- 功能边界清晰、可增量实现
- 有完整的 v1.23 PR 作为参考答案
- 难度适中（2/5），全新受众也能在 2.5 小时内完成核心部分

---

## Workshop Mapping Recommendation

| Feature | Difficulty | Workshop | 理由 |
|---|---|---|---|
| **#2** Tag search | 1/5 | **W1 入门选项** | 4 文件 +33 行，最快速建立信心，30-45 分钟完成 |
| **#15** Actions workflow dispatch | 2/5 | **W1 进阶选项** | 10 文件、高自包含、无 DB 变更，适合进度快的学员 |
| **#13** Issue time estimate | 3/5 | **W4**（综合实战） | 21 文件跨多层，适合综合运用四天所学 |
| **#17** Support repo license | 4/5 | **课后任务** | 47 文件过多，拆分核心链路为课后挑战 |

### Gap Analysis

W2（Bug 修复）和 W3（重构）需要从 v1.22 代码库中找其他目标：

| Workshop | 需要什么 | 建议 |
|---|---|---|
| **W2** Bug 修复 | 一个真实的或预设的 Bug | 从 v1.22 的 GitHub Issues 中找一个已修复的 Bug（推荐 issues 包含复现步骤的） |
| **W3** 重构 | 一个需要重构的模块 | 在 Gitea 代码库中找一个"上帝方法"或复杂条件分支（如 `routers/web/repo/issue.go` 中的长方法） |

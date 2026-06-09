# Gitea 1.23 Refactor Candidates

> Criteria: (1) Appears in v1.23.0 release notes REFACTOR section, (2) Has a linked GitHub issue, (3) PR was properly merged into v1.23.0.

## Summary

- **Total REFACTOR entries in release notes:** 69
- **All PRs merged:** verified
- **Meeting all 3 criteria (has linked issue):** 5
- **Missing linked issue:** 64（重构通常由开发者自驱，不依赖外部 issue）

---

## Refactors Meeting All Criteria (5)

| # | Refactor | PR | Issue | Issue Title | Type | Topic | Files | Lines |
|---|----------|-----|-------|-------------|------|-------|-------|-------|
| 1 | Split mail sender sub package from mailer service package | [#32618](https://github.com/go-gitea/gitea/pull/32618) | [#18664](https://github.com/go-gitea/gitea/issues/18664) | gomail library unmaintained | upstream | — | 15 | +503 -405 |
| 2 | Refactor some LDAP code | [#32849](https://github.com/go-gitea/gitea/pull/32849) | [#32844](https://github.com/go-gitea/gitea/issues/32844) | LDAP: cannot login using email | bug | — | 6 | +259 -173 |
| 3 | Refactor render system (orgmode) | [#32671](https://github.com/go-gitea/gitea/pull/32671) | [#29100](https://github.com/go-gitea/gitea/issues/29100) | Align markdown and orgmode link rendering | bug | content-rendering | 3 | +85 -50 |
| 4 | Migrate vue components to setup | [#32329](https://github.com/go-gitea/gitea/pull/32329) | [#32377](https://github.com/go-gitea/gitea/issues/32377) | Undefined errors on Activity page | bug | ui | 15 | +708 -714 |
| 5 | Refactor maven package registry | [#33049](https://github.com/go-gitea/gitea/pull/33049) | [#33036](https://github.com/go-gitea/gitea/issues/33036) | Maven GroupID/ArtifactID separator not unique | bug | packages | 5 | +143 -60 |

---

## Issue #1: Split mail sender sub package from mailer service package

**PR [#32618](https://github.com/go-gitea/gitea/pull/32618)** · Issue [#18664](https://github.com/go-gitea/gitea/issues/18664) · 15 files · +503 -405

### Refactor Description

将 `services/mailer/mailer.go` 中所有与邮件发送相关的代码（SMTP 客户端、sendmail 调用、message 构造、认证逻辑）拆分到一个新的子包 `services/mailer/sender/` 中。**核心原则：只移动代码，不修改逻辑。** 拆分后，对 `gopkg.in/gomail.v2`（已停止维护）的依赖被隔离在 `sender` 子包内，为后续替换邮件库铺路。

### Background: Why This Refactor

Issue #18664 指出 Gitea 使用的 `gopkg.in/gomail.v2` 自 2016 年起不再维护，其活跃 fork 也已停止更新。社区讨论认为应切换到维护中的替代库。本 PR 作为前置重构：先将 gomail 依赖隔离到独立子包，后续只需替换 `sender/` 内部的实现，不影响上层 `mailer` 服务代码。

### Before (v1.22 baseline)

```
services/mailer/
├── mailer.go          ← 442 行：混合了 Message 类型、gomail 转换、
│                         smtpSender/sendmailSender/dummySender、
│                         SMTP 认证(loginAuth/ntlmAuth)、队列管理
├── mail.go            ← 邮件内容构造
├── mail_*.go          ← 各类邮件模板(issue/release/repo/team invite)
├── mailer_test.go
├── incoming/          ← 收件处理（独立，不涉及此次重构）
├── token/
└── notify.go
```

`mailer.go` 承担了过多职责：
- `Message` 结构体 + `ToMessage()` 方法（依赖 `gomail`）
- `smtpSender`（依赖 `net/smtp` + `gomail`）
- `sendmailSender`（调用系统 sendmail 命令）
- `dummySender`（开发模式空发送）
- SMTP 认证：`loginAuth`、`ntlmAuth`（依赖 `github.com/Azure/go-ntlmssp`）
- 队列初始化 `NewContext()`

### After (v1.23, PR #32618)

```
services/mailer/
├── mailer.go          ← 精简为 ~62 行：仅保留队列管理 + 调用 sender 接口
├── mail.go            ← 改为调用 sender.Message（无逻辑变更）
├── mail_*.go          ← 改为调用 sender.Message（无逻辑变更）
├── sender/            ← 新子包
│   ├── sender.go      ← 27 行：Sender 接口定义 + 工厂函数
│   ├── message.go     ← 112 行：Message 结构体（从 mailer.go 移出）
│   ├── smtp.go        ← 150 行：smtpSender（从 mailer.go 移出）
│   ├── smtp_auth.go   ← 69 行：loginAuth + ntlmAuth（从 mailer.go 移出）
│   ├── sendmail.go    ← 76 行：sendmailSender（从 mailer.go 移出）
│   ├── dummy.go       ← 26 行：dummySender（从 mailer.go 移出）
│   └── message_test.go← 从 mailer_test.go 移出
├── incoming/
├── token/
└── notify.go
```

### Affected Layers

| Layer | Files | Key Changes |
|---|---|---|
| **services/mailer/** | 6 | `mailer.go` -380/+10（大幅精简）；`mail.go`、`mail_release.go`、`mail_repo.go`、`mail_team_invite.go`、`mail_test.go` 改为引用 `sender.Message` |
| **services/mailer/sender/** | 6 | 新增 `sender.go`、`message.go`、`smtp.go`、`smtp_auth.go`、`sendmail.go`、`dummy.go` |
| **routers/private/** | 1 | `mail.go` 改为引用 `sender` 包 |
| **tests/** | 1 | `incoming_email_test.go` 改为引用 `sender` 包 |

### Change Breakdown

- **New files**: 6（`sender/` 子包内的 6 个文件，全部是从 `mailer.go` 移出的代码）
- **Modified files**: 9（import 路径从 `mailer` 改为 `mailer/sender`）
- **Deleted code from mailer.go**: 380 行移出到 `sender/`
- **Renamed files**: 1（`mailer_test.go` 中的 Message 测试 → `sender/message_test.go`）
- **核心约束**: "Just move, no code change" — 所有被移出的代码逻辑完全不变

### Difficulty: 2/5（较低，最适合重构入门）

**Reasoning**:
- **明确的重构模式**：Extract Sub-package（提取子包），Go 中最常见的重构手法之一
- **严格的行为不变**：PR 作者声明 "Just move, no code change"，降低了引入 bug 的风险
- **清晰的边界**：`mailer.go` 中所有与 `gomail`/`smtp`/`sendmail` 相关的代码全部移到 `sender/`
- **Go 工具链支持**：`go build` 会立即捕获所有遗漏的 import 变更
- **规模适中**：15 文件，但实质是 1 个大文件拆成 6 个小文件 + 9 个文件的 import 调整

### Workshop Risks

- **Import 路径变更**：所有引用 `mailer.Message` 的地方需要改为 `mailer/sender.Message`，容易遗漏
- **循环依赖风险**：如果 `sender/` 意外引用了 `mailer` 上层包，会产生循环导入。但 PR 已验证不会
- **gomail 接口适配**：`sender.go` 中的 `Sender` 接口封装了 `gomail.Sender`，学员需要理解这层间接

### Recommended Workshop Fit

**W3（重构 / SDIE 反向理解）— 最佳重构入门选择**

这个重构是典型的 "理解代码结构 → 识别职责边界 → 安全拆分" 流程：

1. **Phase 1 — 阅读**：学员阅读 `mailer.go`（442 行），识别出 4 种 sender 和 1 种 message 的职责划分
2. **Phase 2 — 设计**：设计 `sender` 子包的接口和文件组织
3. **Phase 3 — 执行**：创建子包、移动代码、调整 import
4. **Phase 4 — 验证**：`go build` + `make test-backend` 确认行为不变

优势：
- 不需要理解业务逻辑（邮件发送是通用知识）
- 有明确的 issue 驱动动机（替换过期依赖）
- 编译器是最佳验证工具（Go 的强类型 + import 检查）
- 与 feature 和 bugfix 形成完整的 workshop 三角：W1(正向新增) → W2(定位修复) → W3(结构重构)

---

## Workshop Mapping Update

| Workshop | Topic | Recommended PR | Difficulty |
|---|---|---|---|
| **W1** Feature (正向 SDIE) | #15 Actions workflow dispatch 或 #2 Tag search | PR [#28163](https://github.com/go-gitea/gitea/pull/28163) | 2/5 |
| **W2** Bugfix (定位修复) | 从 bugfix_candidates.md 中选择 | 待定 | 待定 |
| **W3** Refactor (结构重构) | #1 Split mail sender | PR [#32618](https://github.com/go-gitea/gitea/pull/32618) | 2/5 |

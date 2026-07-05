# Gitea Baseline Specification

Complete feature specification of Gitea v1.22.x, organized by bounded-context domains. Each domain has a single boundary rule that determines feature placement. This baseline serves as the reference for delta specs — change proposals diff against these documents.

## Categories

| # | Domain | Boundary Rule | Features | Spec |
|---|--------|--------------|----------|------|
| 01 | Identity & Authentication | Who is this user? How do they prove it? | User accounts, profiles, emails, passwords, sessions, access tokens, 2FA, WebAuthn, external auth (LDAP/SMTP/PAM/SSPI), OpenID, OAuth2 client (16 providers), OAuth2 Provider, SSH keys, heatmap, visibility, following, rename/deletion, CAPTCHA, avatars, badges | [identity-authentication/spec.md](identity-authentication/spec.md) |
| 02 | Access Control & Organization | What can this user do? How are collective identities governed? | RBAC (5 levels, unit permissions), organizations, teams, team invites, collaborators, deploy keys, user/org blocking, GPG keys (commit trust), restricted users, org Actions secrets/runners, org package cleanup | [access-control-organization/spec.md](access-control-organization/spec.md) |
| 03 | Repository & Code | Manage the source code artifact and its git lifecycle | Repo CRUD, transfer, fork, mirror, adoption, templates, git operations (branches, tags, protected branches), code browsing, editor, diff, migration, topics, starring, watching, wiki, activity, repo/code search | [repository-code/spec.md](repository-code/spec.md) |
| 04 | Collaboration | Human workflow on top of repositories | Issues, labels (repo+org), milestones, dependencies, pull requests (review/merge/auto-merge), projects, comments, reactions, content history, CODEOWNERS, close keywords, timetracking, notification triggers/UI, agit flow, issue/PR search, autocomplete | [collaboration/spec.md](collaboration/spec.md) |
| 05 | CI/CD & Automation | Automate the build/test/delivery pipeline on the server | Gitea Actions (14 triggers), Actions secrets/variables/artifacts/logs, webhooks, git hooks, commit status | [cicd-automation/spec.md](cicd-automation/spec.md) |
| 06 | Packages & Releases | Store and distribute published software artifacts | Package registry (21 types), container registry (OCI v2), releases, LFS, attachments, cleanup rules | [packages-releases/spec.md](packages-releases/spec.md) |
| 07 | Platform Services | Cross-cutting services and external interfaces | REST API, markup renderers, themes, custom assets, SSH server, email, i18n, proxy/camo, SSE, RSS/Atom feeds, federation, notification infrastructure, git transport | [platform-services/spec.md](platform-services/spec.md) |
| 08 | Administration & Operations | Instance-level management | Admin panel, configuration, DB migrations, backup/restore, logging, process management, storage backends, queues, cron, metrics, health checks, pprof, doctor checks (~27), CLI, indexer system, sitemaps, update checker | [admin-ops/spec.md](admin-ops/spec.md) |

## Boundary Constitution

Each domain has a single boundary rule. When a feature could fit two domains, the tie-break defaults resolve the ambiguity:

- **Identity vs Access Control**: value is "who you are" → 01; value is "what you can do" → 02
- **Repository vs Collaboration**: operates on the git artifact → 03; human workflow on top → 04
- **CI/CD vs Packages**: automates build/test → 05; stores/distributes published output → 06
- **Feature vs Platform**: user-facing product capability → its feature domain; shared service → 07

## How to Use

### For Delta Specs
A delta spec references the baseline: "Change feature X in domain Y from behavior A to behavior B." The boundary rule determines which domain a feature belongs to.

### Requirement ID Mapping (old → new prefix)

Features moved between domains retain their numeric IDs but adopt the new domain's prefix:

| Old Domain | Old Prefix | New Domain | New Prefix |
|------------|-----------|------------|-----------|
| Identity & Access | `UA-`, `AUTH-` | 01 Identity & Authentication | `IDN-` |
| Identity & Access | `RBAC-`, `UA-05`, `AUTH-07` | 02 Access Control & Organization | `ACL-` |
| Organization | `ORG-` | 02 Access Control & Organization | `ACL-` |
| Code Management | `REPO-` | 03 Repository & Code | `REPO-` |
| Collaboration | `COLL-` | 04 Collaboration | `COLL-` |
| CI/CD & Automation | `CI-` | 05 CI/CD & Automation | `CICD-` |
| Packages & Artifacts | `PKG-` | 06 Packages & Releases | `PKG-` |
| Integration & Extension | `INT-` | 07 Platform Services | `PLT-` |
| Integration & Extension | `INT-` (OAuth2, CAPTCHA, avatars) | 01 Identity & Authentication | `IDN-` |
| Identity & Access | `FED-` | 07 Platform Services | `PLT-` |
| Admin & Ops | `ADM-` | 08 Administration & Operations | `OPS-` |
| Search & Discovery | `SRCH-` | DISSOLVED → consuming domains | (varies by domain) |

**Note:** Some section numbers shifted where features moved in/out (e.g., Repository §10-12 moved to Access Control, Wiki/Search moved in at §15-19). Within each section, the numeric suffixes (e.g., -001, -101, -301) are preserved from the original specs.

## Stats
- 8 domains (was 9 — Search & Discovery dissolved)
- ~114 features documented
- ~330KB total specification

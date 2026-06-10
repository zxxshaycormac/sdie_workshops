# Gitea Baseline Specification

Complete feature specification of Gitea v1.22.x, organized by MECE (Mutually Exclusive, Collectively Exhaustive) categories. This baseline serves as the reference for delta specs — change proposals diff against these documents.

## Categories

| # | Category | Features | Spec |
|---|----------|----------|------|
| 01 | Identity & Access | User accounts, authentication, authorization (RBAC), federation | [01-identity-access.md](01-identity-access.md) |
| 02 | Collaboration | Issues, pull requests, wiki, projects, milestones, labels, reactions, notifications, timetracking, comments | [02-collaboration.md](02-collaboration.md) |
| 03 | Organization | Orgs, teams, org labels, org projects, blocking, badges, settings | [03-organization.md](03-organization.md) |
| 04 | Code Management | Repo CRUD, git operations, branches, code browsing, editor, forking, mirroring, migration, topics, collaborators, starring/watching | [04-code-management.md](04-code-management.md) |
| 05 | CI/CD & Automation | Actions, webhooks, agit flow, auto-merge, commit status, git hooks, cron, queues | [05-ci-cd-automation.md](05-ci-cd-automation.md) |
| 06 | Packages & Artifacts | Package registry (21 types), releases, LFS, attachments, storage backends | [06-packages-artifacts.md](06-packages-artifacts.md) |
| 07 | Administration & Ops | Admin panel, configuration, queues, metrics, health checks, migrations, backup, notices, logging, process management, pprof | [07-admin-ops.md](07-admin-ops.md) |
| 08 | Search & Discovery | Repo search, code search, issue/PR search, user/org search, explore, indexer system, sitemap | [08-search-discovery.md](08-search-discovery.md) |
| 09 | Integration & Extension | OAuth2 provider, external auth, REST API, markup renderers, themes, custom assets, SSH, email, i18n, CAPTCHA, avatars, proxy, update checker | [09-integration-extension.md](09-integration-extension.md) |

## How to Use

### For Delta Specs
A delta spec references the baseline: "Change feature X in category Y from behavior A to behavior B." The baseline provides precise attachment points (feature names, routes, APIs, config keys).

### Spec Structure (per feature)
- **What**: one-paragraph description
- **Behaviors**: bullet list of notable behaviors
- **UI Routes**: web interface routes
- **API Endpoints**: REST API endpoints
- **Config**: relevant `app.ini` settings
- **Constraints**: limitations and edge cases

## Stats
- 9 categories
- ~80 features documented
- ~85KB total specification

# Gitea Baseline Specification

Complete feature specification of Gitea v1.22.x, organized by MECE (Mutually Exclusive, Collectively Exhaustive) categories. This baseline serves as the reference for delta specs — change proposals diff against these documents.

## Categories

| # | Category | Features | Spec |
|---|----------|----------|------|
| 01 | Identity & Access | User accounts, authentication, authorization (RBAC), federation, activity heatmap, GPG keys, OpenID | [identity-access/spec.md](identity-access/spec.md) |
| 02 | Collaboration | Issues, pull requests, wiki, projects, milestones, labels, reactions, notifications, timetracking, comments, CODEOWNERS, issue close keywords, PR review workflow, content history browser | [collaboration/spec.md](collaboration/spec.md) |
| 03 | Organization | Orgs, teams, org labels, org projects, blocking, settings (badges are a user-account feature documented here for reference) | [organization/spec.md](organization/spec.md) |
| 04 | Code Management | Repo CRUD, git operations, branches, code browsing, editor, forking, mirroring, migration, topics, collaborators, starring/watching, deploy keys, anonymous clone, repo transfer, branch rename/restore, activity, tag browsing, wiki-as-git-repo | [code-management/spec.md](code-management/spec.md) |
| 05 | CI/CD & Automation | Actions (14 trigger events), webhooks, agit flow, auto-merge (7 merge styles), commit status, git hooks, cron, queues, actions secrets, actions variables, actions artifacts, actions log storage | [ci-cd-automation/spec.md](ci-cd-automation/spec.md) |
| 06 | Packages & Artifacts | Package registry (21 types), releases, LFS, attachments, storage backends | [packages-artifacts/spec.md](packages-artifacts/spec.md) |
| 07 | Administration & Ops | Admin panel, configuration, queues, metrics, health checks, migrations, backup, notices, logging, process management, pprof (separate localhost:6060 server), doctor checks (~27), CLI surface, diagnosis/stacktrace | [admin-ops/spec.md](admin-ops/spec.md) |
| 08 | Search & Discovery | Repo search, code search (Bleve/Elasticsearch), issue/PR search (Bleve/ES/Meilisearch/db), user/org search, explore, indexer system, sitemap, autocomplete | [search-discovery/spec.md](search-discovery/spec.md) |
| 09 | Integration & Extension | OAuth2 provider (authorization server), OAuth2 client login (16 providers), external auth, REST API, markup renderers (Markdown/Org/CSV/Console/Asciicast/external), themes (5 built-in), custom assets, SSH, email, i18n, CAPTCHA, avatars, proxy, camo, update checker, SSE, RSS/Atom feeds, Federation/ActivityPub | [integration-extension/spec.md](integration-extension/spec.md) |

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
- ~114 features documented
- ~330KB total specification

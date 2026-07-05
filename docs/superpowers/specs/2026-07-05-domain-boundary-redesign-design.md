# Gitea Domain Boundary Redesign

> **Status:** Approved design — pending implementation plan
> **Date:** 2026-07-05
> **Scope:** Restructure the OpenSpec baseline from 9 MECE domains to 8 bounded-context domains

## 1. Motivation

The current 9-domain MECE decomposition (`openspec/specs/spec.md`) serves three purposes simultaneously: change management (delta spec organization), product roadmap (value-stream thinking), and architecture alignment (codebase module mapping). A codebase-verified review revealed **6 structural boundary failures** where the decomposition violates mutual exclusivity, misplaces features, or creates artificial domains.

This document defines a replacement decomposition — **8 bounded-context domains** — that resolves all identified failures and provides unambiguous boundary rules for feature placement.

## 2. Diagnosis: Current 9-Domain Failures

### Failure 1: Identity & Access is 3 domains in 1 (17 features)

"Identity" (who you are), "Authentication" (how you prove it), and "Authorization" (what you can do) are conflated. They share the User model but serve different user journeys and change for different reasons. Misplaced features include activity heatmap (product analytics), GPG keys (commit trust), federation (integration protocol), and user blocking (abuse prevention).

**Codebase evidence:** Three separate service packages — `services/user`, `services/auth`, `services/asymkey` — plus `models/perm` for authorization.

### Failure 2: Organization is artificially starved (8 features)

Orgs ARE Users in Gitea's data model (`UserTypeOrganization`). The domain mirrors `services/org` but the product boundary between "my account" and "my organization" is artificial — permissions, blocking, and settings flow across both.

### Failure 3: Collaboration ↔ Code Management boundary leaks

PRs are in Collaboration but are deeply tied to git branch operations (Code Management). Wiki is in Collaboration but is literally a git repo. Auto-merge straddles Collaboration, Code Management, AND CI/CD.

### Failure 4: Search & Discovery is infrastructure, not a product domain

Search serves all other domains. It has no independent user journey. The indexer system is ops infrastructure; the search UI is a capability layered onto each domain. Elevating it to a product domain is like making "Database" a domain.

### Failure 5: Integration & Extension is a catch-all (MECE violation)

This domain is defined by what it is NOT (not in the other 8) rather than by a coherent product purpose. It contains OAuth2 provider (identity), REST API (platform interface), markup renderers (content), themes (UI), SSH/email/i18n (infrastructure), and federation (identity).

### Failure 6: CI/CD & Automation mixes product with plumbing

Actions (product feature), webhooks (integration mechanism), cron (infrastructure), queues (infrastructure), and git hooks (repo config) are lumped together. A change to queue configuration has nothing to do with a change to Actions workflow syntax.

### Root Cause

Every failure traces to mixing three incompatible organizing principles: user journey / value stream, codebase module, and "everything else" bucket. When a feature serves multiple journeys AND lives in one code module, the current decomposition has no rule for where it belongs.

## 3. The 8-Domain Decomposition

### 3.1 Boundary Constitution

Each domain has a single boundary rule — a one-sentence test that determines unambiguously whether a feature belongs there.

| # | Domain | Boundary Rule | Tie-Break Default |
|---|--------|--------------|-------------------|
| 01 | Identity & Authentication | Features that answer "who is this user?" and "how do they prove it?" | — |
| 02 | Access Control & Organization | Features that answer "what can this user do?" and "how are collective identities governed?" | Identity-vs-Access tie → goes here |
| 03 | Repository & Code | Features that manage the source code artifact and its git lifecycle | Code-vs-Collab tie → goes here |
| 04 | Collaboration | Features that enable human workflow on top of repositories | — |
| 05 | CI/CD & Automation | Features that automate the build/test/delivery pipeline on the server | — |
| 06 | Packages & Artifacts | Features that store and distribute software artifacts | — |
| 07 | Platform Services | Cross-cutting services and external interfaces that serve all domains | Catch-all of last resort |
| 08 | Administration & Operations | Instance-level management not tied to a specific product feature | — |

**Disambiguation rules:**
- Identity vs Access Control: if the feature's value is "establishing who you are" → 01; if the value is "governing what you can do" → 02
- Repository vs Collaboration: if it operates on the git artifact itself → 03; if it's a human workflow layered on top → 04
- CI/CD vs Packages: if it automates build/test → 05; if it stores/distributes output → 06
- Feature vs Platform: if it's a user-facing product capability → its feature domain; if it's a shared service consumed by multiple domains → 07

### 3.2 Domain Definitions

#### Domain 01: Identity & Authentication

**Purpose:** Individual identity and the mechanisms by which a user proves who they are.

**Features:** User registration, profile management, email management, password management, session management, access tokens, 2FA (TOTP), WebAuthn/Passkey, external authentication sources (LDAP, SMTP, PAM, DLDAP, SSPI), OpenID consumer, OAuth2 client login (16 providers), OAuth2 Provider (authorization server), SSH keys, activity heatmap, user visibility, user rename & deletion, CAPTCHA, avatars, badges.

**Codebase anchor:** `services/user`, `services/auth`, `models/user`, `models/auth`, `services/externalaccount`

**Boundary — what's OUT:** RBAC permission system → 02; organizations/teams → 02; GPG keys → 02; user blocking enforcement → 02; federation protocol → 07.

#### Domain 02: Access Control & Organization

**Purpose:** Permission governance and collective identity management — who can do what, and how teams are organized.

**Features:** RBAC permission system (5 levels, unit permissions), organizations (creation, settings, deletion), teams (membership, unit permissions, invites), collaborators (repo-level access), deploy keys, user/org blocking (creation + enforcement), GPG keys (commit trust verification, trust models), org labels, restricted users.

**Codebase anchor:** `services/org`, `models/organization`, `models/perm`, `services/asymkey`, `models/asymkey`

**Boundary — what's OUT:** Individual user accounts → 01; OAuth2 auth mechanisms → 01; repo-level features → 03; org-level projects → 04.

#### Domain 03: Repository & Code

**Purpose:** The source code artifact lifecycle — from creation through browsing, editing, and migration.

**Features:** Repo CRUD, repo transfer, forking, mirroring, adoption, git operations (branches, tags, protected branches), code browsing, web editor, blame, diff, repo migration/import, topics, starring, watching, wiki (as git repository), repo activity, code frequency. Repo search and code search are documented here as capabilities of this domain.

**Codebase anchor:** `services/repository`, `modules/git`, `services/mirror`, `services/migrations`, `services/wiki`

**Boundary — what's OUT:** Issues/PRs/projects → 04; collaborator access grants → 02; repo-level webhooks → 05; releases/packages → 06.

#### Domain 04: Collaboration

**Purpose:** Human collaboration workflows layered on top of repositories.

**Features:** Issues, labels, milestones, dependencies, pull requests (review workflow, merge, auto-merge), projects (all scopes: repo/org/user), comments, reactions, content history, CODEOWNERS, issue close keywords, timetracking, notifications (in-app + email triggers), agit flow. Issue/PR search is documented here as a capability of this domain.

**Codebase anchor:** `services/issue`, `services/pull`, `services/notify`, `services/uinotification`, `services/automerge`, `services/agit`, `models/issues`, `models/project`

**Boundary — what's OUT:** Git branch operations → 03; repo creation → 03; Actions CI → 05; commit status checks → 05.

#### Domain 05: CI/CD & Automation

**Purpose:** The server-side automation pipeline for building, testing, and integrating code.

**Features:** Gitea Actions (workflows, runners, tasks, schedules), Actions secrets, Actions variables, Actions artifacts, Actions log storage, webhooks (delivery, types, host allowlist, proxy), git hooks (server-side pre-receive/post-receive etc.), commit status.

**Codebase anchor:** `services/actions`, `services/webhook`, `services/secrets`, `models/actions`, `models/webhook`

**Boundary — what's OUT:** Package registry → 06; releases → 06; cron/queues infrastructure → 08; agit flow (human PR creation) → 04.

#### Domain 06: Packages & Artifacts

**Purpose:** Storage and distribution of software artifacts.

**Features:** Package registry (21 types), container registry (OCI v2 endpoint), releases (draft, prerelease, tag), LFS, attachments, storage backends (local, MinIO).

**Codebase anchor:** `services/packages`, `services/lfs`, `services/release`, `services/attachment`, `modules/storage`, `models/packages`

**Boundary — what's OUT:** Actions artifacts/logs → 05; repo creation → 03; commit status → 05.

#### Domain 07: Platform Services

**Purpose:** Cross-cutting services and external interfaces consumed by all other domains.

**Features:** REST API (Swagger, scopes, pagination), federation (ActivityPub, WebFinger, NodeInfo), markup renderers (Markdown, Org, CSV, Asciicast, Console, external), RSS/Atom feeds, SSE (/user/events), email (mailer transport), SSH server, themes, custom assets, i18n, proxy, camo.

**Codebase anchor:** `routers/api/v1`, `services/mailer`, `services/markup`, `services/feed`, `modules/ssh`, `modules/translation`, `services/webtheme`

**Boundary — what's OUT:** OAuth2 Provider (identity delegation) → 01; CAPTCHA (auth anti-abuse) → 01; update checker → 08; queues/cron → 08.

**Note:** This domain is the weakest boundary by nature — it collects shared services. The "catch-all of last resort" rule requires justifying why no other domain fits before placing a feature here.

#### Domain 08: Administration & Operations

**Purpose:** Instance-level management capabilities not tied to a specific product feature.

**Features:** Admin panel, user/org/repo administration, configuration (app.ini system), DB migrations, backup/restore (dump), logging, process management, queues (infrastructure), cron (infrastructure), metrics, health checks, pprof, doctor checks (~27), CLI surface, indexer system (search backends: Bleve/Elasticsearch/Meilisearch/db), sitemaps, update checker.

**Codebase anchor:** `routers/admin`, `services/doctor`, `services/cron`, `services/indexer`, `modules/queue`, `cmd/`

**Boundary — what's OUT:** Search UI for specific domains → documented within those domains; repo admin actions → 03; user admin actions → 01.

### 3.3 Dissolved Domain: Search & Discovery

The current "Search & Discovery" domain is dissolved. Search is a capability, not a domain:

- **Repo/code search** → documented in Domain 03 (Repository & Code)
- **Issue/PR search** → documented in Domain 04 (Collaboration)
- **User/org search + explore** → documented in Domains 01/02 (Identity/Access Control)
- **Autocomplete** → documented in Domain 04 (Collaboration)
- **Indexer backends (Bleve/ES/Meili/db)** → documented in Domain 08 (Admin/Ops)

## 4. Key Boundary Decisions

### GPG keys → Access Control (02), not Identity (01)
GPG keys are about commit trust verification — "was this commit really made by an authorized collaborator?" The trust model (committer/collaborator/collaborator-committer) is a permission policy. Identity stores the key material, but the product value is access governance.

### Blocking → Access Control (02), not Identity (01)
Blocking's product value is access enforcement — "this person cannot interact with me." The blocking record creation is secondary to its enforcement effects across issues, PRs, and teams. Both the creation and enforcement aspects live in Access Control.

### Wiki → Repository & Code (03), not Collaboration (04)
A wiki IS a git repository (`models/repo/wiki.go`, cloneable via `.wiki` suffix). It belongs with the code artifact lifecycle. Wiki page editing is a collaboration activity, but the wiki itself is a code artifact.

### Agit flow → Collaboration (04), not CI/CD (05)
Agit is a PR creation workflow — a human pushes code with special push options to create a pull request. It's a collaboration action, not an automation pipeline step.

### Webhooks → CI/CD (05), not Platform (07)
Webhooks are the primary event-driven integration mechanism for CI/CD. While they also serve general integrations (Slack, Discord), their codebase coupling to repo/PR events and Actions makes CI/CD their natural home.

### Cron/Queues → Admin/Ops (08), not CI/CD (05)
Cron and queues are infrastructure — they serve Actions, webhooks, mailer, indexer, and more. They're not a product feature; they're platform plumbing.

### Search dissolved, not relocated
Search is a capability, not a domain. Every domain has its own search. The indexer infrastructure is ops plumbing.

### Notifications → Collaboration (04)
Notifications are the communication channel for collaboration workflows. Email as a transport mechanism → Platform Services (07); notification triggers and UI → Collaboration (04).

## 5. Migration Impact

### Directory Structure Change

Current:
```
openspec/specs/
  identity-access/spec.md
  collaboration/spec.md
  organization/spec.md
  code-management/spec.md
  ci-cd-automation/spec.md
  packages-artifacts/spec.md
  admin-ops/spec.md
  search-discovery/spec.md
  integration-extension/spec.md
```

Proposed:
```
openspec/specs/
  identity-authentication/spec.md      (01)
  access-control-organization/spec.md  (02)
  repository-code/spec.md              (03)
  collaboration/spec.md                (04)
  cicd-automation/spec.md              (05)
  packages-artifacts/spec.md           (06)
  platform-services/spec.md            (07)
  admin-ops/spec.md                    (08)
```

### Feature Redistribution Summary

| From (current) | To (new) | Features moved |
|----------------|----------|----------------|
| Identity & Access → 01 | Identity & Authentication | OAuth2 client, OAuth2 Provider, CAPTCHA, avatars, badges (from Integration/Org) |
| Identity & Access → 02 | Access Control & Organization | RBAC, blocking, GPG keys, restricted users |
| Identity & Access → 07 | Platform Services | Federation |
| Organization → 02 | Access Control & Organization | All org features merge into new domain 02 |
| Organization → 01 | Identity & Authentication | Badges (user-account feature) |
| Code Management → 02 | Access Control & Organization | Collaborators, deploy keys |
| Code Management → 03 | Repository & Code | Core repo features stay; wiki moves in from Collaboration |
| Collaboration → 03 | Repository & Code | Wiki |
| Collaboration → 04 | Collaboration | Core collab features stay; agit moves in from CI/CD |
| CI/CD → 04 | Collaboration | Agit flow |
| CI/CD → 08 | Admin/Ops | Cron, queues |
| Integration → 01 | Identity & Authentication | OAuth2 client, OAuth2 Provider, CAPTCHA, avatars |
| Integration → 07 | Platform Services | REST API, markup, themes, SSH, email, i18n, feeds, proxy, camo |
| Integration → 08 | Admin/Ops | Update checker |
| Search & Discovery | DISSOLVED | Search UI → consuming domains; indexer infra → Admin/Ops |

### Domain Size Balance

| Domain | ~Features | Codebase packages |
|--------|-----------|-------------------|
| 01 Identity & Authentication | 15 | 5 |
| 02 Access Control & Organization | 10 | 5 |
| 03 Repository & Code | 12 | 5 |
| 04 Collaboration | 14 | 8 |
| 05 CI/CD & Automation | 7 | 5 |
| 06 Packages & Artifacts | 6 | 6 |
| 07 Platform Services | 11 | 7 |
| 08 Administration & Operations | 13 | 6 |

No domain exceeds 15 features. No domain has fewer than 6. The current largest (Identity: 17) shrinks; the current weakest (Organization: 8) grows by absorbing RBAC.

## 6. Evaluation Against Triple Purpose

### Change Management
Each domain has high internal cohesion — features that change together are grouped. A change to Actions doesn't touch Packages. A change to auth mechanisms doesn't touch collaboration (except through the RBAC contract). Delta specs can unambiguously identify their target domain via the boundary rule.

### Product Roadmap
Domains map to clear value streams. "Improve CI/CD" and "expand package registry" are independent roadmap items. "Strengthen access control" covers orgs, teams, RBAC, and blocking as a coherent product initiative.

### Architecture Alignment
Each domain maps to 5-8 service packages in the codebase. The mapping is natural, not forced. The only domain without a clean 1:1 codebase mapping is Platform Services (by nature — it collects cross-cutting services).

## 7. Open Questions for Implementation

1. **Requirement ID renumbering:** The current specs use domain-prefixed IDs (UA-01-001, COLL-02-001, etc.). The new domains need new prefixes. Should old IDs be preserved with a mapping table, or fully renumbered?

2. **Search documentation depth:** Since search is dissolved into consuming domains, each domain spec needs a search subsection. How much indexer-level detail belongs in Admin/Ops vs. the consuming domain?

3. **Cross-domain contracts:** Features that span domains (e.g., protected branches use RBAC from domain 02 and git operations from domain 03) need explicit interface contracts. Should these be documented in both domains or in a separate contracts section?

# Domain Boundary Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restructure the OpenSpec baseline from 9 MECE domains to 8 bounded-context domains per the approved design at `docs/superpowers/specs/2026-07-05-domain-boundary-redesign-design.md`.

**Architecture:** Each new domain spec is assembled by extracting sections from current spec files according to the feature redistribution table. Requirement IDs preserve their numeric portion but adopt a new domain prefix. Old directories are removed after all new specs are verified.

**Tech Stack:** Markdown spec files only — no code changes.

## Global Constraints

- Source of truth for redistribution: `docs/superpowers/specs/2026-07-05-domain-boundary-redesign-design.md` §3.2 (domain definitions) and §5 (redistribution table).
- Current spec content was verified against the codebase in the prior correction pass — preserve requirement text verbatim when moving between domains; only change the ID prefix.
- Requirement ID scheme: each new domain uses a 3-4 letter prefix, preserving the original numeric IDs for traceability (e.g., `UA-01-001` → `IDN-01-001`). A mapping table goes in the index.
- Each new domain spec must follow the existing structure: title, baseline statement, numbered feature sections (Ubiquitous / Event-Driven / Optional / State-Driven / Unwanted Behaviour requirements), Configuration Reference, Business Rules, Edge Cases, Success Criteria.
- Search is dissolved: consuming domains get a "Search" subsection; Admin/Ops gets indexer infrastructure. Cross-domain references use the format "see Domain NN (<name>)".
- Do NOT modify any code files. Only `openspec/specs/` markdown files.

## Requirement ID Prefix Mapping

| New Domain | New Prefix | Absorbs old prefixes |
|------------|-----------|---------------------|
| 01 Identity & Authentication | `IDN` | UA, AUTH (auth-related), INT (OAuth2/CAPTCHA/avatars) |
| 02 Access Control & Organization | `ACL` | RBAC, ORG, REPO (collaborators/deploy keys), UA (blocking/GPG/restricted) |
| 03 Repository & Code | `REPO` | REPO (stays), COLL (wiki → REPO), SRCH (repo/code search UI) |
| 04 Collaboration | `COLL` | COLL (stays), CI (agit), ORG (org labels), SRCH (issue search/autocomplete) |
| 05 CI/CD & Automation | `CICD` | CI (stays, minus agit) |
| 06 Packages & Releases | `PKG` | PKG (stays, minus storage config) |
| 07 Platform Services | `PLT` | INT (rest), FED, notification infrastructure |
| 08 Administration & Operations | `OPS` | ADM (stays), CI (cron/queues), SRCH (indexer infra), INT (update checker), storage config |

---

### Task 1: Create Directory Structure

**Files:**
- Create: `openspec/specs/identity-authentication/` (dir)
- Create: `openspec/specs/access-control-organization/` (dir)
- Create: `openspec/specs/repository-code/` (dir)
- Create: `openspec/specs/collaboration/` (already exists — will be repopulated)
- Create: `openspec/specs/cicd-automation/` (dir)
- Create: `openspec/specs/packages-releases/` (dir)
- Create: `openspec/specs/platform-services/` (dir)
- Create: `openspec/specs/admin-ops/` (already exists — will be updated)

- [ ] **Step 1: Create new directories**

```bash
mkdir -p openspec/specs/identity-authentication
mkdir -p openspec/specs/access-control-organization
mkdir -p openspec/specs/repository-code
mkdir -p openspec/specs/cicd-automation
mkdir -p openspec/specs/packages-releases
mkdir -p openspec/specs/platform-services
```

- [ ] **Step 2: Commit**

```bash
git add openspec/specs/
git commit -m "chore(spec): create new domain directory structure"
```

---

### Task 2: Create Domain 01 — Identity & Authentication

**Files:**
- Create: `openspec/specs/identity-authentication/spec.md`
- Read (sources): `openspec/specs/identity-access/spec.md`, `openspec/specs/integration-extension/spec.md`, `openspec/specs/organization/spec.md`

**What to assemble (prefix: `IDN`):**

From `identity-access/spec.md`, extract these sections (keep all requirement text, change prefix UA→IDN, AUTH→IDN):
- User Registration (UA-01 → IDN-01)
- User Profile Management (UA-02 → IDN-02)
- Email Management (UA-03 → IDN-03)
- Password Management (UA-04 → IDN-04)
- User Rename & Deletion (UA-06 → IDN-06)
- Password Authentication & Sessions (AUTH-01 → IDN-07)
- Access Tokens (AUTH-03 → IDN-08)
- Two-Factor Authentication (AUTH-04 → IDN-09)
- WebAuthn / Passkey (AUTH-05 → IDN-10)
- External Authentication Sources (AUTH-06 → IDN-11)
- OpenID Authentication (AUTH-08 → IDN-12)
- Activity Heatmap (UA-07 → IDN-13)

From `integration-extension/spec.md`, extract (change prefix INT→IDN):
- OAuth2 client login section (INT-13 → IDN-14) — all 16 providers
- CAPTCHA requirements (INT-08 → IDN-15)
- Avatar requirements (INT-12 → IDN-16)

From `organization/spec.md`, extract (change prefix ORG→IDN):
- Badges section (ORG-07 → IDN-17) — already rewritten as user-account feature

**Add new sections:**
- User Following (new IDN-18) — following/unfollowing users, `models/user/follow.go`. Requirements: follow another user, unfollow, list followers, list following, prevent self-follow.
- OAuth2 Provider (new IDN-19) — move from `integration-extension` §17 (INT-17 → IDN-19). Authorization server endpoints.

**Exclude (move to other domains):**
- RBAC (AUTH-06/RBAC → Domain 02), User Blocking (UA-05 → Domain 02), GPG Keys (AUTH-07 → Domain 02), Federation (FED-01 → Domain 07), Restricted Users (UA-02-701 → Domain 02)

**Config Reference:** Combine `[security]`, `[session]`, `[oauth2]`, `[oauth2_client]`, `[openid]` sections from current Identity spec. Add `[camo]` if avatar-related.

- [ ] **Step 1: Read all source spec files**
- [ ] **Step 2: Assemble new spec with correct sections and ID prefixes**
- [ ] **Step 3: Verify all listed features are present and excluded features are absent**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/identity-authentication/spec.md
git commit -m "docs(spec): create Domain 01 Identity & Authentication"
```

---

### Task 3: Create Domain 02 — Access Control & Organization

**Files:**
- Create: `openspec/specs/access-control-organization/spec.md`
- Read (sources): `openspec/specs/identity-access/spec.md`, `openspec/specs/organization/spec.md`, `openspec/specs/code-management/spec.md`

**What to assemble (prefix: `ACL`):**

From `identity-access/spec.md`, extract (change prefix):
- Permission System / RBAC (RBAC-01 → ACL-01)
- User Blocking (UA-05 → ACL-02)
- GPG Keys (AUTH-07 → ACL-03)
- Restricted Users state-driven requirements (UA-02-701 → ACL-04)

From `organization/spec.md`, extract (change prefix ORG→ACL):
- Organizations (ORG-01 → ACL-05)
- Teams (ORG-02 → ACL-06)
- Team Invites (ORG-03 → ACL-07)
- Organization Settings (ORG-08 → ACL-08, minus deleted fabricated reqs)
- Org-level Actions secrets/runners (new reqs added in correction → ACL-09)
- Org package cleanup rules (new reqs → ACL-10)

From `code-management/spec.md`, extract (change prefix REPO→ACL):
- Collaborators (REPO-10 → ACL-11)
- Deploy Keys (REPO-12 → ACL-12)

**Exclude:**
- Org labels → Domain 04; Badges → Domain 01; Individual user accounts → Domain 01

**Config Reference:** Combine org-related config, `DEFAULT_ORG_MEMBER_VISIBLE`, RBAC-related `[security]` keys.

- [ ] **Step 1: Read all source spec files**
- [ ] **Step 2: Assemble new spec**
- [ ] **Step 3: Verify feature inclusion/exclusion**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/access-control-organization/spec.md
git commit -m "docs(spec): create Domain 02 Access Control & Organization"
```

---

### Task 4: Create Domain 03 — Repository & Code

**Files:**
- Create: `openspec/specs/repository-code/spec.md`
- Read (sources): `openspec/specs/code-management/spec.md`, `openspec/specs/collaboration/spec.md`, `openspec/specs/search-discovery/spec.md`

**What to assemble (prefix: `REPO`):**

From `code-management/spec.md`, extract most sections (keep prefix REPO):
- All repo CRUD, transfer, fork, mirror, adoption, git operations, branches, tags, protected branches, code browsing, editor, diff, migration, topics, starring, watching, activity, code frequency
- Add: repo templates (new REPO section — IsTemplate, generate-from-template)
- Keep: protected branches with cross-reference note ("required approvals → Domain 04, required status checks → Domain 05, push access → Domain 02")

From `collaboration/spec.md`, extract (change prefix COLL→REPO):
- Wiki (COLL-03 → REPO wiki section) — wiki as git repository

From `search-discovery/spec.md`, extract (change prefix SRCH→REPO):
- Repo search (SRCH-01 → REPO search subsection)
- Code search (SRCH-02 → REPO search subsection)
- Explore repos (SRCH-05 repo parts → REPO search subsection)

**Exclude:**
- Collaborators/deploy keys → Domain 02; Issues/PRs → Domain 04; Webhooks → Domain 05; Releases/packages → Domain 06

**Config Reference:** Keep `[repository]`, `[git]`, `[git.timeout]`, `[git.config]`, `[mirror]`, `[migrations]`, `[ui]` (MaxDisplayFileSize), indexer `[indexer]` code-search keys.

- [ ] **Step 1: Read all source spec files**
- [ ] **Step 2: Assemble new spec**
- [ ] **Step 3: Verify protected branch cross-references are present**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/repository-code/spec.md
git commit -m "docs(spec): create Domain 03 Repository & Code"
```

---

### Task 5: Create Domain 04 — Collaboration

**Files:**
- Create: `openspec/specs/collaboration/spec.md` (overwrite existing)
- Read (sources): `openspec/specs/collaboration/spec.md` (current), `openspec/specs/ci-cd-automation/spec.md`, `openspec/specs/organization/spec.md`, `openspec/specs/search-discovery/spec.md`

**What to assemble (prefix: `COLL`):**

From current `collaboration/spec.md`, keep (prefix COLL stays):
- Issues, labels (now including org labels — note "labels span repo and org scope"), milestones, dependencies, pull requests (review/merge/auto-merge), projects (all scopes), comments, reactions, content history, CODEOWNERS, close keywords, timetracking, collaboration notification triggers + in-app UI

From `ci-cd-automation/spec.md`, extract (change prefix CI→COLL):
- Agit flow (CI-03 → COLL agit section)

From `organization/spec.md`, extract (change prefix ORG→COLL):
- Org labels (ORG-04 → COLL labels section) — merge into existing labels section, note both scopes

From `search-discovery/spec.md`, extract (change prefix SRCH→COLL):
- Issue/PR search (SRCH-03 → COLL search subsection)
- Autocomplete (SRCH-08 → COLL search subsection)

**Add note:** "Notification infrastructure (event bus) is documented in Domain 07. This domain documents collaboration-specific notification triggers and the in-app notification UI."

**Exclude:**
- Wiki → Domain 03; Cron/queues → Domain 08; Actions → Domain 05; Notification infrastructure → Domain 07

**Config Reference:** Keep `[service]` keys (ENABLE_TIMETRACKING, DEFAULT_ENABLE_DEPENDENCIES, AUTO_WATCH_NEW_REPOS, AUTO_WATCH_ON_CHANGES, ENABLE_NOTIFY_MAIL).

- [ ] **Step 1: Read all source spec files**
- [ ] **Step 2: Assemble new spec**
- [ ] **Step 3: Verify agit flow and org labels are present; wiki is absent**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/collaboration/spec.md
git commit -m "docs(spec): create Domain 04 Collaboration"
```

---

### Task 6: Create Domain 05 — CI/CD & Automation

**Files:**
- Create: `openspec/specs/cicd-automation/spec.md`
- Read (sources): `openspec/specs/ci-cd-automation/spec.md`

**What to assemble (prefix: `CICD`):**

From `ci-cd-automation/spec.md`, extract (change prefix CI→CICD):
- Gitea Actions (CI-01 → CICD-01) — workflows, runners, tasks, schedules, statuses
- Actions Secrets (CI-09 → CICD-02)
- Actions Variables (CI-10 → CICD-03)
- Actions Artifacts (CI-11 → CICD-04) — note: "CI-produced artifacts; published packages/releases → Domain 06"
- Actions Log Storage — note as separate storage type
- Webhooks (CI-02 → CICD-05) — delivery, types, host allowlist, proxy
- Git Hooks (CI-06 → CICD-06) — server-side pre-receive/post-receive
- Commit Status (CI-05 → CICD-07)

**Exclude:**
- Agit flow → Domain 04; Cron → Domain 08; Queues → Domain 08; Packages → Domain 06

**Config Reference:** Keep `[actions]`, `[actions_log]`, `[webhook]`, `[webhook]` config sections. Add cross-reference: "Cron and queue configuration → Domain 08."

- [ ] **Step 1: Read source spec**
- [ ] **Step 2: Assemble new spec, removing agit/cron/queues sections**
- [ ] **Step 3: Verify agit, cron, queues are absent**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/cicd-automation/spec.md
git commit -m "docs(spec): create Domain 05 CI/CD & Automation"
```

---

### Task 7: Create Domain 06 — Packages & Releases

**Files:**
- Create: `openspec/specs/packages-releases/spec.md`
- Read (sources): `openspec/specs/packages-artifacts/spec.md`

**What to assemble (prefix: `PKG`):**

From `packages-artifacts/spec.md`, extract (keep prefix PKG):
- Package registry (21 types) — PKG-01
- Container registry (OCI v2) — PKG-02
- Releases — PKG-04
- LFS — PKG-05
- Attachments — PKG-06
- Cleanup rules — PKG-01-901

**Exclude:**
- Storage backend configuration → Domain 08 (add reference: "Storage backend configuration → Domain 08")
- Actions artifacts → Domain 05

**Config Reference:** Keep `[packages]` section. Move `[storage]` / `[storage.packages]` to Domain 08 reference. Keep `[lfs]`, `[attachment]`.

- [ ] **Step 1: Read source spec**
- [ ] **Step 2: Assemble new spec, removing storage backend config section**
- [ ] **Step 3: Verify storage config is absent, storage usage reference is present**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/packages-releases/spec.md
git commit -m "docs(spec): create Domain 06 Packages & Releases"
```

---

### Task 8: Create Domain 07 — Platform Services

**Files:**
- Create: `openspec/specs/platform-services/spec.md`
- Read (sources): `openspec/specs/integration-extension/spec.md`, `openspec/specs/identity-access/spec.md`, `openspec/specs/search-discovery/spec.md`

**What to assemble (prefix: `PLT`):**

From `integration-extension/spec.md`, extract (change prefix INT→PLT):
- REST API (INT-01 → PLT-01)
- Markup Renderers (INT-02 → PLT-02) — Markdown, Org, CSV, Asciicast, Console, external
- Themes (INT-03 → PLT-03)
- Custom Assets (INT-04 → PLT-04)
- SSH Server (INT-05 → PLT-05)
- Email / Mailer (INT-06 → PLT-06)
- i18n (INT-07 → PLT-07)
- RSS/Atom Feeds (INT-16 → PLT-08)
- SSE (INT-15 → PLT-09)
- Proxy / Camo (INT-14 → PLT-10)
- Federation / ActivityPub (INT-18 → PLT-11)

From `identity-access/spec.md`, extract (change prefix FED→PLT):
- Federation (FED-01 → PLT-11, merge with INT-18)

**Add new sections:**
- Notification Infrastructure (PLT-12) — event bus / `services/notify` dispatcher. Requirements: dispatch events to registered notifiers (mailer, webhook, uinotification, indexer, mirror, automerge). Note: "Each domain documents its own notification triggers."
- Git Transport (PLT-13) — HTTP smart protocol, SSH git commands. Requirements: serve git clone/fetch/push over HTTP and SSH, authenticate git operations. Note: "Git operations (branch, tag, commit) → Domain 03."

**Exclude:**
- OAuth2 Provider → Domain 01; OAuth2 client → Domain 01; CAPTCHA → Domain 01; Avatars → Domain 01; Update checker → Domain 08

**Config Reference:** Keep `[server]` (SSH, proxy), `[mailer]`, `[markup]`, `[camo]`, `[i18n]`, `[federation]` sections.

- [ ] **Step 1: Read all source spec files**
- [ ] **Step 2: Assemble new spec**
- [ ] **Step 3: Verify OAuth2/CAPTCHA/avatars/update-checker are absent**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/platform-services/spec.md
git commit -m "docs(spec): create Domain 07 Platform Services"
```

---

### Task 9: Create Domain 08 — Administration & Operations

**Files:**
- Create: `openspec/specs/admin-ops/spec.md` (overwrite existing)
- Read (sources): `openspec/specs/admin-ops/spec.md` (current), `openspec/specs/ci-cd-automation/spec.md`, `openspec/specs/search-discovery/spec.md`, `openspec/specs/integration-extension/spec.md`

**What to assemble (prefix: `OPS`):**

From current `admin-ops/spec.md`, keep (change prefix ADM→OPS):
- Admin panel, user/org/repo administration, configuration, DB migrations, backup/restore, logging, process management, metrics, health checks, pprof, doctor checks, CLI surface

From `ci-cd-automation/spec.md`, extract (change prefix CI→OPS):
- Cron (CI-07 → OPS cron section) — infrastructure
- Queues (CI-08 → OPS queue section) — infrastructure

From `search-discovery/spec.md`, extract (change prefix SRCH→OPS):
- Indexer system / backends (SRCH-06 → OPS indexer section) — Bleve, Elasticsearch, Meilisearch, db
- Sitemaps (SRCH-07 → OPS sitemap section)

From `integration-extension/spec.md`, extract (change prefix INT→OPS):
- Update Checker (INT-17-update → OPS update-checker section)

**Add:**
- Storage Backends (OPS storage section) — local, MinIO configuration. Note: "Storage is consumed by avatars (01), packages/LFS/attachments (06), Actions artifacts/logs (05), repo archives (03)."
- Search infrastructure reference: "Search UI for each domain is documented within that domain. This section documents the indexer backends and configuration only."

**Config Reference:** Add `[queue]`, `[queue.<name>]`, `[cron.<name>]`, `[metrics]`, `[storage]`, `[storage.*]`, `[indexer]` (infrastructure keys), `[server].ENABLE_PPROF`.

- [ ] **Step 1: Read all source spec files**
- [ ] **Step 2: Assemble new spec, adding cron/queues/indexer/storage sections**
- [ ] **Step 3: Verify cron, queues, indexer, storage, sitemaps are present**
- [ ] **Step 4: Commit**

```bash
git add openspec/specs/admin-ops/spec.md
git commit -m "docs(spec): create Domain 08 Administration & Operations"
```

---

### Task 10: Update Index & Remove Old Directories

**Files:**
- Modify: `openspec/specs/spec.md`
- Delete: `openspec/specs/identity-access/`, `openspec/specs/organization/`, `openspec/specs/code-management/`, `openspec/specs/ci-cd-automation/`, `openspec/specs/packages-artifacts/`, `openspec/specs/search-discovery/`, `openspec/specs/integration-extension/`

**Important:** Do NOT delete `collaboration/` or `admin-ops/` — they were overwritten in-place by Tasks 5 and 9.

- [ ] **Step 1: Rewrite `spec.md` index**

Update the categories table to the 8 new domains with correct links. Update stats (feature count, total size). Add a "Requirement ID Mapping" subsection documenting the old→new prefix mapping table.

- [ ] **Step 2: Remove old directories**

```bash
rm -rf openspec/specs/identity-access
rm -rf openspec/specs/organization
rm -rf openspec/specs/code-management
rm -rf openspec/specs/ci-cd-automation
rm -rf openspec/specs/packages-artifacts
rm -rf openspec/specs/search-discovery
rm -rf openspec/specs/integration-extension
```

- [ ] **Step 3: Verify directory structure**

```bash
ls -d openspec/specs/*/
```
Expected: exactly 8 domain dirs + no old dirs.

- [ ] **Step 4: Commit**

```bash
git add openspec/specs/spec.md openspec/specs/
git commit -m "docs(spec): update index to 8-domain structure, remove old directories"
```

---

### Task 11: Consistency Verification

**Files:**
- Read: all 8 new `openspec/specs/*/spec.md` files + `openspec/specs/spec.md`

- [ ] **Step 1: Verify no orphaned requirement IDs** — grep all new specs for old prefixes (UA-, AUTH-, RBAC-, FED-, ORG-, REPO- where REPO shouldn't exist outside domain 03, CI-, SRCH-, INT-). Only `REPO-` should appear in domain 03; `BR-` prefixes are OK everywhere.

Run: `rg "^\- \*\*(UA|AUTH|RBAC|FED|ORG|CI|SRCH|INT)-" openspec/specs/*/spec.md`
Expected: no matches (all old prefixes replaced).

- [ ] **Step 2: Verify ME — no feature appears in two domains.** Spot-check the most boundary-sensitive features: protected branches (only in 03), webhooks (only in 05), labels (only in 04), GPG keys (only in 02), notifications (triggers in 04, infrastructure in 07).

Run: `rg -c "webhook" openspec/specs/*/spec.md` — should appear primarily in cicd-automation.
Run: `rg -c "GPG" openspec/specs/*/spec.md` — should appear primarily in access-control-organization.

- [ ] **Step 3: Verify CE — all features have a home.** Check that: user following is in 01, repo templates in 03, agit in 04, cron in 08, queues in 08, storage config in 08, git transport in 07, notification infrastructure in 07.

Run: `rg -l "user following" openspec/specs/*/spec.md` — expect identity-authentication.
Run: `rg -l "agit" openspec/specs/*/spec.md` — expect collaboration.
Run: `rg -l "storage backend" openspec/specs/*/spec.md` — expect admin-ops.

- [ ] **Step 4: Verify index links resolve**

Run: `for f in $(rg -o '\(([^)]+/spec.md)\)' openspec/specs/spec.md -r '$1'); do test -f "openspec/specs/$f" && echo "OK: $f" || echo "BROKEN: $f"; done`
Expected: all OK.

- [ ] **Step 5: Commit if any fixes were needed**

```bash
git add openspec/specs/
git commit -m "docs(spec): consistency verification fixes"
```

If no fixes needed, report "All checks passed."

---

## Notes for the Executor

- **Task independence:** Tasks 2–9 are fully independent (different files). They can be dispatched to parallel subagents.
- **Task 1 must complete first** (directory creation).
- **Task 10 must wait for Tasks 2–9** (old dirs removed only after new specs exist).
- **Task 11 must be last** (verifies the final state).
- **Content preservation:** When moving requirements between domains, copy the exact requirement text. Only change the ID prefix. Do not rephrase or renumber the numeric portion.
- **Cross-references:** Use the format "see Domain NN (<name>) §<section>" for features that span domains.

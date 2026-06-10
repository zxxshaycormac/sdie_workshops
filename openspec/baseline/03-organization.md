# 03 — Organization

## 1. Organizations

### What
Special user accounts (UserTypeOrganization) representing groups, companies, or projects. Multiple users collaborate under a single identity with structured team hierarchies.

### Behaviors
- Create organization: any user with CanCreateOrganization permission
- Default "Owners" team created automatically
- Visibility levels: Public, Limited (authenticated), Private (members only)
- Organization rename creates a redirect from old name
- MaxRepoCreation: configurable limit (-1 = unlimited)
- RepoAdminChangeTeamAccess: controls whether team members can change repo settings
- Delete org: removes org and all associated data (cannot delete if owns repos/packages)

### UI Routes
- GET /org/create — Create org page
- POST /org/create — Create new org
- GET /:org — Organization home
- GET /:org/settings — Org settings
- POST /:org/settings — Update org
- POST /:org/settings/avatar — Update avatar
- POST /:org/settings/delete — Delete org
- GET /:org/members — Member list
- POST /:org/members — Update member visibility
- GET /:org/repositories — Repo list

### API Endpoints
- GET /user/orgs — List current user's orgs
- GET /users/:username/orgs — List user's orgs
- GET /orgs — List all public orgs
- POST /orgs — Create org
- GET /orgs/:org — Get org details
- PATCH /orgs/:org — Update org
- DELETE /orgs/:org — Delete org
- GET /orgs/:org/activities/feeds — Org activity feeds

### Config
- [admin] DISABLE_REGULAR_ORG_CREATION — restrict org creation to admins

### Constraints
- Org names follow username validation rules (no reserved names)
- Cannot create org with reserved names
- Default owner team "Owners" cannot be deleted
- Last owner cannot be removed

---

## 2. Teams

### What
Role-based access control groups within organizations. Teams grant granular permissions to groups of users across repositories and units.

### Behaviors
- Team permission levels: Owner, Admin, Write, Read, None
- IncludesAllRepositories: team can access all org repos
- CanCreateOrgRepo: team members can create new repos
- Unit-level permissions: Code, Issues, Pulls, Wiki, ExternalWiki, ExternalIssues, Projects, Packages, Actions
- Team invitations via email with token-based acceptance
- Owner teams automatically get all unit permissions

### UI Routes
- GET /:org/teams — List teams
- GET /:org/teams/new — Create team form
- POST /:org/teams — Create team
- GET /:org/teams/:team — Team details
- GET /:org/teams/:team/members — Manage members
- POST /:org/teams/:team/members — Add/remove members
- GET /:org/teams/:team/repos — Manage repos
- POST /:org/teams/:team/repos — Add/remove repos

### API Endpoints
- GET /orgs/:org/teams — List teams
- POST /orgs/:org/teams — Create team
- GET /user/teams — List current user's teams
- GET /teams/:id — Get team
- PATCH /teams/:id — Update team
- DELETE /teams/:id — Delete team
- GET /teams/:id/members — List members
- PUT /teams/:id/members/:username — Add member
- DELETE /teams/:id/members/:username — Remove member
- GET /teams/:id/repos — List team repos
- PUT /teams/:id/repos/:repo — Add repo
- DELETE /teams/:id/repos/:repo — Remove repo

### Constraints
- "Owners" team cannot be deleted
- Last owner cannot be removed from owner team
- Team names must be unique within org
- Cannot invite users already on the team

---

## 3. Org-level Labels

### What
Organization-scoped labels shared across all org repositories. Provides consistency in labeling issues and PRs.

### Behaviors
- CRUD operations for org-wide labels
- Initialize from predefined label templates
- Exclusive/scoped labels: only one label of a type can be applied
- Labels can be archived (soft delete)
- Used by any issue/PR in any org repository

### UI Routes
- GET /:org/settings/labels — List org labels
- POST /:org/settings/labels — Create label
- POST /:org/settings/labels/update — Update label
- POST /:org/settings/labels/delete — Delete label
- POST /:org/settings/labels/init — Initialize from template

### Constraints
- Org-level labels are distinct from repo-specific labels
- Cannot delete labels currently in use
- Label templates follow standard format

---

## 4. Org-level Projects (Kanban)

### What
Organization-level kanban boards for project management. Can contain issues from any repository within the organization.

### Behaviors
- Project type: TypeOrganization (distinct from repo-level projects)
- Columns/boards with drag-and-drop card management
- Issues from multiple org repos can be added to same board
- Open/close project lifecycle

### UI Routes
- GET /:org/projects — List projects
- GET /:org/projects/new — Create project
- POST /:org/projects — Create project
- GET /:org/projects/:id — View project board
- GET /:org/projects/:id/boards — Manage boards

### Constraints
- Requires Projects unit to be enabled for org repos

---

## 5. User Blocking

### What
Users can block other users from interacting. Block applies across the organization context.

### Behaviors
- Block/unblock with optional note/reason
- Blocked users cannot: comment on issues/PRs, create PRs to blocker's repos, receive team invitations
- Admins cannot be blocked
- Blocking is directional (blocker->blockee)

### UI Routes
- GET /:org/settings/blocked_users — List blocked users
- POST /:org/settings/blocked_users — Block/unblock user

### Constraints
- Cannot block organizations (only individual users)
- Cannot block admin users

---

## 6. Badges

### What
Visual achievement/status indicators for users and organizations.

### Behaviors
- Badge definition: slug (unique ID), description, image URL
- Multiple badges per user/org
- Badge award/removal management

### Constraints
- Badge system is admin-managed (no user-facing CRUD routes)
- Extensible but limited to predefined badge types

---

## 7. Org Settings

### What
Central management for all org-level configuration.

### Settings Categories
1. Profile: name, full name, email, description, website, location
2. Avatar: upload custom, delete (use generated), URL-based
3. Visibility: Public, Limited, Private
4. Repository: default permission, max creation limit, member fork permissions
5. Security: email confirmation requirements
6. Integrations: webhooks, OAuth2 apps, packages, runners

### UI Routes
- GET /:org/settings — Main settings
- POST /:org/settings — Update settings
- GET /:org/settings/avatar — Avatar management
- POST /:org/settings/avatar — Update avatar
- GET /:org/settings/delete — Delete confirmation
- POST /:org/settings/delete — Confirm deletion
- GET /:org/settings/hooks — Webhooks
- GET /:org/settings/labels — Labels
- GET /:org/settings/blocked_users — Blocked users
- GET /:org/settings/packages — Package settings
- GET /:org/settings/oauth2 — OAuth2 apps
- GET /:org/settings/runners — CI/CD runners

### Constraints
- Only org owners can access settings
- Some settings require admin privileges
- Visibility changes may cascade to repos

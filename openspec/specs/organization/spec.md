# 03 — Organization

Baseline specification of Gitea's organization, team, membership, label, project, blocking, and settings subsystems. (Badges are documented as a user-account feature in §7 — see Identity & Access for full requirements.) All requirements describe the current (v1.22.x) system behavior.

---

## 1. Organizations

**User Story:** As a user, I want to create and manage organizations so that multiple collaborators can work under a shared group identity with structured access control.

### Ubiquitous Requirements (Organization Properties)

- **ORG-01-001:** `The system shall represent organizations as special user accounts of type UserTypeOrganization.`
- **ORG-01-002:** `The system shall assign each organization a unique name following the same validation rules as usernames (no reserved names, no ".git"/".wiki"/".rss"/".atom" suffixes).`
- **ORG-01-003:** `The system shall support three organization visibility levels: public (visible to everyone), limited (visible to authenticated users only), and private (visible to members only).`
- **ORG-01-004:** `The system shall create a default "Owners" team with full permissions when a new organization is created.`
- **ORG-01-005:** `The system shall track a configurable maximum repository creation limit (MaxRepoCreation) per organization, where -1 means unlimited.`
- **ORG-01-006:** `The system shall store the following profile properties per organization: name, full name, email, description, website, and location.`
- **ORG-01-007:** `The system shall support avatar management for organizations via upload, URL, or generated default.`
- **ORG-01-008:** `The system shall support the RepoAdminChangeTeamAccess flag that controls whether repository admins can change team access settings.`

### Event-Driven Requirements (Organization Lifecycle)

- **ORG-01-101:** `When a user with CanCreateOrganization permission submits a valid organization creation form, the system shall create the organization and add the creator to the default "Owners" team.`
- **ORG-01-102:** `When an organization is renamed, the system shall create a persistent redirect from the old name to the new name.`
- **ORG-01-103:** `When an organization is renamed, the system shall update all repository paths under that organization to reflect the new name.`
- **ORG-01-104:** `When an organization owner confirms deletion, the system shall remove the organization and all associated data.`
- **ORG-01-105:** `When an organization's visibility is changed, the system shall apply the visibility change to the organization profile and member listings.`
- **ORG-01-106:** `When a request targets an organization's former name after a rename, the system shall redirect to the current organization name.`

### Optional Feature Requirements (Organization Creation Control)

- **ORG-01-201:** `Where the admin setting DISABLE_REGULAR_ORG_CREATION is enabled, the system shall restrict organization creation to administrators only.`

### Unwanted Behaviour Requirements (Organization Errors)

- **ORG-01-301:** `If a user without CanCreateOrganization permission attempts to create an organization, then the system shall deny the creation.`
- **ORG-01-302:** `If an organization creation is submitted with a reserved name, then the system shall reject the creation with a validation error.`
- **ORG-01-303:** `If an organization deletion is attempted while the organization still owns repositories or packages, then the system shall deny the deletion.`
- **ORG-01-304:** `If a non-owner user attempts to modify organization settings, then the system shall deny the operation.`
- **ORG-01-305:** `If an organization rename is submitted with a name that already exists, then the system shall reject the rename.`

---

## 2. Teams

**User Story:** As an organization owner, I want to create and configure teams so that I can grant granular role-based access to groups of users across repositories and feature units.

### Ubiquitous Requirements (Team Properties)

- **ORG-02-001:** `The system shall define five team permission levels: Owner, Admin, Write, Read, and None.`
- **ORG-02-002:** `The system shall enforce team name uniqueness within each organization.`
- **ORG-02-003:** `The system shall support the following ten unit-level permission targets: Code, Issues, Pull Requests, Releases, Wiki, External Wiki, External Tracker, Projects, Packages, and Actions.`
- **ORG-02-004:** `The system shall assign all unit permissions automatically to Owner-level teams.`
- **ORG-02-005:** `The system shall track an IncludesAllRepositories flag per team that grants access to all current and future organization repositories.`
- **ORG-02-006:** `The system shall track a CanCreateOrgRepo flag as a first-class per-team property that authorizes members to create repositories under the organization.`
- **ORG-02-007:** `The system shall persist per-unit team access modes on each TeamUnit record (Read, Write, or Admin) rather than a binary on/off, allowing teams to hold different access levels across units.`

### Event-Driven Requirements (Team Lifecycle)

- **ORG-02-101:** `When an organization owner creates a new team, the system shall register the team with the specified permission level, unit permissions, and repository access scope.`
- **ORG-02-102:** `When an organization owner updates a team's configuration, the system shall persist the changes and recalculate effective permissions for all team members.`
- **ORG-02-103:** `When an organization owner deletes a team, the system shall remove the team and revoke its granted permissions from all associated repositories.`
- **ORG-02-104:** `When a team has IncludesAllRepositories enabled, the system shall automatically grant the team access to any new repository created in the organization.`
- **ORG-02-105:** `When a team has CanCreateOrgRepo enabled, the system shall allow team members to create new repositories under the organization.`

### Optional Feature Requirements (Team Configuration)

- **ORG-02-201:** `Where a team is configured with specific unit permissions, the system shall restrict team member access to only those units with the specified access levels.`
- **ORG-02-202:** `Where a team is not configured as IncludesAllRepositories, the system shall require explicit repository assignment to grant team access.`

### Unwanted Behaviour Requirements (Team Errors)

- **ORG-02-301:** `If a user attempts to delete the default "Owners" team, then the system shall deny the deletion.`
- **ORG-02-302:** `If a non-owner user attempts to create or modify a team, then the system shall deny the operation.`
- **ORG-02-303:** `If a team creation is submitted with a name that already exists within the organization, then the system shall reject the creation.`
- **ORG-02-304:** `If the last owner of an organization is attempted to be removed from the Owners team, then the system shall deny the removal.`

---

## 3. Team Membership & Invites

**User Story:** As an organization owner or team maintainer, I want to add and remove members from teams so that users receive the appropriate access, and I want to invite external users via email so they can join the team.

### Ubiquitous Requirements (Membership Properties)

- **ORG-03-001:** `The system shall track team membership as a relationship between user accounts and teams with a role indicator.`
- **ORG-03-002:** `The system shall allow a user to belong to multiple teams within the same organization.`
- **ORG-03-003:** `The system shall compute a user's effective organization permission as the union of all team permissions.`
- **ORG-03-004:** `The system shall generate a unique token for each team invitation email.`
- **ORG-03-005:** `The system shall allow organization members to control their own membership visibility (public or concealed).`

### Event-Driven Requirements (Membership Workflow)

- **ORG-03-101:** `When an authorized user adds a member to a team, the system shall grant that user the team's permissions immediately.`
- **ORG-03-102:** `When an authorized user removes a member from a team, the system shall revoke that team's permissions from the user.`
- **ORG-03-103:** `When a team invitation is sent to an email address, the system shall send an invitation email containing a unique token and acceptance link.`
- **ORG-03-104:** `When a recipient clicks a valid invitation link and accepts, the system shall add the user to the team and mark the invitation as accepted.`
- **ORG-03-105:** `When a user leaves an organization, the system shall remove the user from all teams within that organization.`
- **ORG-03-106:** `When an organization member updates their membership visibility, the system shall reflect the change in the organization's member listing.`

### State-Driven Requirements (Invitation Lifecycle)

- **ORG-03-702:** `While a user is a member of an organization, the system shall include that organization in the user's organization listing.`

### Unwanted Behaviour Requirements (Membership Errors)

- **ORG-03-301:** `If a user attempts to add an already-existing team member to the same team, then the system shall reject the addition.`
- **ORG-03-303:** `If a non-authorized user attempts to modify team membership, then the system shall deny the operation.`
- **ORG-03-304:** `If a user attempts to invite a blocked user to a team, then the system shall deny the invitation.`
- **ORG-03-305:** `If the last owner of an organization attempts to leave, then the system shall deny the departure.`

---

## 4. Org-level Labels

**User Story:** As an organization owner, I want to create and manage labels at the organization level so that all repositories in the organization share consistent labeling for issues and pull requests.

### Ubiquitous Requirements (Label Properties)

- **ORG-04-001:** `The system shall store organization-level labels with name, description, color, and exclusivity flag.`
- **ORG-04-002:** `The system shall make organization-level labels available to issues and pull requests in all repositories owned by the organization.`
- **ORG-04-003:** `The system shall distinguish organization-level labels from repository-specific labels in the label management interface.`

### Event-Driven Requirements (Label Lifecycle)

- **ORG-04-101:** `When an organization owner creates an org-level label, the system shall make it immediately available across all organization repositories.`
- **ORG-04-102:** `When an organization owner updates an org-level label, the system shall propagate the change to all issues and pull requests currently using that label.`
- **ORG-04-103:** `When an organization owner deletes an org-level label, the system shall remove it from all issues and pull requests currently using that label.`
- **ORG-04-104:** `When an organization owner initializes labels from a template, the system shall create all labels defined in the selected template.`

### Optional Feature Requirements (Exclusive Labels)

- **ORG-04-201:** `Where a label is marked as exclusive (scoped), the system shall allow only one label of that scope to be applied to a given issue or pull request at a time.`
- **ORG-04-202:** `Where a label is archived (soft-deleted), the system shall hide it from the label selection UI while preserving its existing usage on issues and pull requests.`

### Unwanted Behaviour Requirements (Label Errors)

- **ORG-04-301:** `If a non-owner user attempts to create, update, or delete an org-level label, then the system shall deny the operation.`
- **ORG-04-302:** `If an org-level label creation duplicates an existing label name within the organization, then the system shall reject the creation.`
- **ORG-04-303:** `If a label template initialization is requested but no templates are available, then the system shall display an appropriate message.`

---

## 5. Org-level Projects

**User Story:** As an organization member, I want to create and manage kanban-style project boards at the organization level so that I can track issues and pull requests across multiple repositories in a single view.

### Ubiquitous Requirements (Project Properties)

- **ORG-05-001:** `The system shall support organization-level projects as a distinct type (TypeOrganization) separate from repository-level and user-level projects.`
- **ORG-05-002:** `The system shall support multiple columns per project board for card organization.`
- **ORG-05-003:** `The system shall allow cards representing issues and pull requests from any repository owned by the organization.`
- **ORG-05-004:** `The system shall track project state as either open or closed.`

### Event-Driven Requirements (Project Lifecycle)

- **ORG-05-101:** `When an authorized user creates an org-level project with a non-None board type (BasicKanban or BugTriage), the system shall initialize it with the corresponding default columns; a BoardTypeNone project (the default) is created with zero columns.`
- **ORG-05-102:** `When an authorized user adds an issue or pull request to an org-level project, the system shall create a card in the specified column.`
- **ORG-05-103:** `When an authorized user moves a card between columns, the system shall update the card's column assignment.`
- **ORG-05-104:** `When an authorized user closes a project, the system shall mark the project as closed and preserve all cards and columns.`
- **ORG-05-105:** `When an authorized user reopens a closed project, the system shall restore the project to open state with all previous data.`

### Optional Feature Requirements (Project Configuration)

- **ORG-05-201:** `Where the Projects unit is enabled for organization repositories, the system shall allow issues and pull requests from those repositories to be added to org-level projects.`
- **ORG-05-202:** `Where an authorized user creates a custom column, the system shall add the column to the project board at the specified position.`

### Unwanted Behaviour Requirements (Project Errors)

- **ORG-05-301:** `If a non-authorized user attempts to create or modify an org-level project, then the system shall deny the operation.`
- **ORG-05-302:** `If an attempt is made to add an issue from a repository outside the organization to an org-level project, then the system shall reject the addition.`
- **ORG-05-303:** `If the Projects unit is disabled for a repository, then the system shall prevent that repository's issues from being added to org-level project boards.`

---

## 6. User Blocking

**User Story:** As an organization owner, I want to block specific users from the organization so that they cannot interact with organization resources.

### Ubiquitous Requirements (Blocking Properties)

- **ORG-06-001:** `The system shall maintain a list of blocked users per organization.`
- **ORG-06-002:** `The system shall allow blocking to be directional (blocker to blockee) and optionally include a reason note.`
- **ORG-06-003:** `The system shall apply blocking constraints across all organization-owned repositories and teams.`

### Event-Driven Requirements (Blocking Workflow)

- **ORG-06-101:** `When an organization owner blocks a user, the system shall prevent the blocked user from commenting on issues and pull requests in organization repositories.`
- **ORG-06-102:** `When an organization owner blocks a user, the system shall prevent the blocked user from creating pull requests to organization repositories.`
- **ORG-06-103:** `When an organization owner blocks a user, the system shall prevent the blocked user from receiving team invitations within the organization.`
- **ORG-06-104:** `When an organization owner unblocks a user, the system shall restore the user's ability to interact with organization resources according to their permissions.`
- **ORG-06-105:** `When an organization owner blocks a user, the system shall allow recording an optional note describing the reason for the block.`

### Unwanted Behaviour Requirements (Blocking Restrictions)

- **ORG-06-301:** `If an organization owner attempts to block an admin user, then the system may create the block record, but enforcement is skipped for admin users (IsUserBlockedBy returns false for admins).`
- **ORG-06-302:** `If an organization owner attempts to block an organization (rather than an individual user), then the system shall deny the block operation.`
- **ORG-06-303:** `If a non-owner user attempts to block or unblock a user at the organization level, then the system shall deny the operation.`

---

## 7. Badges (User-Account Feature)

**Note:** Badges are a **user-account** feature, not an org-domain capability. The `UserBadge` model (`models/user/badge.go`) is keyed by `UserID` and is administered via `/admin/users/{username}/badges` (`routers/api/v1/admin/user_badge.go`). No `OrgBadge` model exists and there is no org-profile badge rendering. Because organizations are stored as `User` records, an admin can technically assign a badge to an org's User ID, but this is incidental rather than an org-domain feature. Detailed requirements for badges live in the Identity & Access spec.

### Ubiquitous Requirements (Badge Properties)

- **ORG-07-001:** `The system shall store badge definitions (Badge) with a unique slug identifier, description, and image URL.`
- **ORG-07-002:** `The system shall associate badges with user accounts via the UserBadge join table keyed on UserID (no org-specific badge model or org-profile rendering exists).`

---

## 8. Org Settings

**User Story:** As an organization owner, I want to configure all aspects of my organization including profile, avatar, visibility, repository permissions, and integrations so that the organization operates according to my requirements.

### Ubiquitous Requirements (Settings Categories)

- **ORG-08-001:** `The system shall restrict access to organization settings to users with owner-level permission in the organization.`
- **ORG-08-002:** `The system shall store the following profile settings per organization: name, full name, email, description, website, and location.`

### Event-Driven Requirements (Settings Updates)

- **ORG-08-101:** `When an organization owner updates profile settings, the system shall persist the changes immediately.`
- **ORG-08-102:** `When an organization owner uploads a custom avatar, the system shall store the image and associate it with the organization.`
- **ORG-08-103:** `When an organization owner deletes the custom avatar, the system shall revert to a generated default avatar.`
- **ORG-08-104:** `When an organization owner changes the organization visibility, the system shall apply the visibility change to the organization profile and cascade visibility effects to member listings.`
- **ORG-08-105:** `When an organization owner confirms deletion via the settings page, the system shall delete the organization and all associated data.`
- **ORG-08-106:** `When an organization owner configures webhooks, the system shall register the webhooks to fire on the specified events across organization repositories.`
- **ORG-08-107:** `When an organization owner registers an OAuth2 application, the system shall create the application scoped to the organization.`
- **ORG-08-108:** `When an organization owner configures a CI/CD runner, the system shall register the runner (ActionRunner/ActionRunnerToken) scoped to the organization's OwnerID for use by organization repositories.`
- **ORG-08-109:** `When an organization owner creates an Actions secret, the system shall store the secret scoped to the organization (Secret.OwnerID = org.ID).`
- **ORG-08-110:** `When an organization owner configures a package cleanup rule, the system shall persist the rule at the organization scope; the owner may also initialize or rebuild the Cargo package index for the organization.`

### Optional Feature Requirements (Integration Settings)

- **ORG-08-201:** `Where the packages feature is enabled, the system shall allow organization owners to configure package storage and access settings.`
- **ORG-08-202:** `Where email confirmation is required, the system shall enforce email verification before applying certain organization settings changes.`

### Unwanted Behaviour Requirements (Settings Errors)

- **ORG-08-301:** `If a non-owner user attempts to access organization settings, then the system shall deny access and return a 404 or redirect.`
- **ORG-08-302:** `If an organization deletion confirmation is submitted with incorrect confirmation text, then the system shall deny the deletion.`
- **ORG-08-303:** `If an organization avatar upload exceeds the maximum allowed file size, then the system shall reject the upload with a size limit error.`
- **ORG-08-304:** `If an invalid email address is provided in organization profile settings, then the system shall reject the update with a validation error.`

---

## Business Rules

- **BR-03-001:** Organization names must be unique across the entire instance (case-insensitive) and follow username validation rules
- **BR-03-002:** The default "Owners" team is created automatically on organization creation and cannot be deleted
- **BR-03-003:** An organization must have at least one owner at all times; the last owner cannot be removed
- **BR-03-004:** Team permission levels follow a strict hierarchy: Owner > Admin > Write > Read > None
- **BR-03-005:** A user's effective organization permission is the union of all team permissions across teams they belong to
- **BR-03-006:** Organization visibility determines which users can see the organization in listings and search results
- **BR-03-007:** MaxRepoCreation of -1 means unlimited repository creation; 0 means no creation; positive values set a numeric limit
- **BR-03-008:** Organization-level labels are distinct from repository-level labels but appear alongside them in issue/PR label selectors
- **BR-03-009:** User blocking at the organization level is directional and does not block the user instance-wide
- **BR-03-010:** Badge definitions are admin-managed and assigned to user accounts (UserBadge.UserID); they are not an org-domain feature
- **BR-03-011:** Org-level projects can reference issues and pull requests from any repository owned by the organization
- **BR-03-012:** Exclusive (scoped) labels allow only one label per scope on a given issue or pull request
- **BR-03-014:** RepoAdminChangeTeamAccess determines whether repository administrators can modify team access for their repositories

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Delete organization with existing repositories | Deny deletion; require transfer or deletion of all repos first |
| Delete organization with existing packages | Deny deletion; require removal of all packages first |
| Remove last owner from Owners team | Deny removal; enforce at least one owner |
| Rename organization to existing name | Reject rename with conflict error |
| Add user already on the team | Reject addition with duplicate membership error |
| Invite blocked user to team | Deny invitation |
| Team invitation token does not expire | Tokens are accepted at any time after issuance until manually revoked |
| Delete org-level label in active use | Remove label from all referencing issues and pull requests |
| Create project without Projects unit enabled | Deny project creation or prevent adding cards from repos with Projects disabled |
| Block an admin user at org level | Block record may be created; enforcement skipped (admin not actually blocked) |
| Block an organization account | Deny block operation (only individual users can be blocked) |
| Access org settings as non-owner | Return 404 or redirect away |
| Org deletion with wrong confirmation text | Deny deletion |
| Team with IncludesAllRepositories and new repo created | Automatically grant team access to the new repository |
| User leaves organization while last owner | Deny departure |
| Visibility change to Private | Hide org from non-member listings and search |

## Success Criteria

- Organization creation completes within 3 seconds for standard configuration
- Team permission changes propagate to all members within 1 second
- Organization rename creates a working redirect accessible indefinitely
- Org-level label CRUD operations complete within 500ms
- Project board rendering with 100+ cards completes within 2 seconds
- Team invitation emails are sent within 30 seconds of creation
- Blocking a user takes effect immediately across all organization resources
- Organization settings changes persist immediately and reflect on next page load
- Permission computation for any organization resource completes in under 50ms
- Org-level project boards support issues from all organization repositories simultaneously

## Configuration Reference

Defaults shown reflect Gitea v1.22.x.

| Key | Section | Default | Purpose |
|-----|---------|---------|---------|
| `DEFAULT_ORG_MEMBER_VISIBLE` | `[service]` | `false` | Whether organization members are listed as visible by default on the org's member page (members may still override their own visibility). Source: `modules/setting/service.go:81,224`. |
| `DEFAULT_ORG_VISIBILITY` | `[service]` | `public` | Default visibility for newly created organizations. Source: `modules/setting/service.go:222`. |
| `DISABLE_REGULAR_ORG_CREATION` | `[admin]` | `false` | When true, restricts organization creation to administrators only. Source: `modules/setting/admin.go:20`. |

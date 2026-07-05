# 02 — Access Control & Organization

Baseline specification of Gitea's permission system (RBAC), organization governance, team management, collaborator and deploy key access control, user/organization blocking, GPG key trust verification, restricted user enforcement, and organization-level Actions secrets/runners and package cleanup rules. All requirements describe the current (v1.22.x) system behavior.

> **Requirement ID traceability:** Requirements migrated from the former 9-domain structure preserve their numeric IDs with a new `ACL` prefix. Source mapping: RBAC→ACL-01, UA-05/ORG-06→ACL-02, AUTH-07→ACL-03, UA-02-701/RBAC-01-701→ACL-04, ORG-01→ACL-05, ORG-02→ACL-06, ORG-03→ACL-07, ORG-08→ACL-08, ORG-08-108/109→ACL-09, ORG-08-110→ACL-10, REPO-10→ACL-11, REPO-12→ACL-12.

---

## 1. Permission System (ACL-01)

**User Story:** As a system administrator, I want fine-grained access control so that users only access resources they are authorized for.

### Ubiquitous Requirements (Permission Hierarchy)

- **ACL-01-001:** `The system shall define five permission levels in strict hierarchy: None (0), Read (1), Write (2), Admin (3), Owner (4).`
- **ACL-01-002:** `The system shall enforce permission inheritance such that each higher level includes all capabilities of lower levels.`
- **ACL-01-003:** `The system shall apply permissions at three scopes: repository (collaborator), organization/team, and unit (feature-level).`

### Event-Driven Requirements (Permission Evaluation)

- **ACL-01-101:** `When a user accesses a repository, the system shall compute the effective permission from collaborator, team, and organization memberships.`
- **ACL-01-102:** `When a user accesses a feature unit (code, issues, pulls, wiki, projects, packages, actions), the system shall check unit-level permissions for that user.`
- **ACL-01-103:** `When an owner-level user accesses any resource in their scope, the system shall grant full access regardless of other restrictions.`

### Optional Feature Requirements (Unit Permissions)

- **ACL-01-201:** `Where team unit permissions are configured, the system shall restrict team members to the specified unit types with the specified access levels.`
- **ACL-01-202:** `Where a repository unit is disabled, the system shall hide that feature from all users regardless of permission level.`

---

## 2. User & Organization Blocking (ACL-02)

**User Story:** As a user or organization owner, I want to block specific users so they cannot interact with me, my repositories, or my organization's resources.

### Ubiquitous Requirements (Blocking Properties)

- **ACL-02-001:** `The system shall maintain a list of blocked users per user account and per organization.`
- **ACL-02-002:** `The system shall allow blocking to be directional (blocker to blockee) and optionally include a reason note.`
- **ACL-02-003:** `The system shall apply blocking constraints across all repositories and teams owned by the blocker or organization.`

### Event-Driven Requirements (User Blocking)

- **ACL-02-101:** `When a user blocks another user, the system shall prevent the blocked user from commenting on the blocker's issues and PRs.`
- **ACL-02-102:** `When a user blocks another user, the system shall prevent the blocked user from creating PRs to the blocker's repositories.`
- **ACL-02-103:** `When a user blocks another user, the system shall prevent the blocked user from being added to teams the blocker manages.`
- **ACL-02-104:** `When a user blocks another user, the system shall allow the blocker to record an optional note describing the reason.`

### Event-Driven Requirements (Organization Blocking)

- **ACL-02-105:** `When an organization owner blocks a user, the system shall prevent the blocked user from commenting on issues and pull requests in organization repositories.`
- **ACL-02-106:** `When an organization owner blocks a user, the system shall prevent the blocked user from creating pull requests to organization repositories.`
- **ACL-02-107:** `When an organization owner blocks a user, the system shall prevent the blocked user from receiving team invitations within the organization.`
- **ACL-02-108:** `When an organization owner unblocks a user, the system shall restore the user's ability to interact with organization resources according to their permissions.`
- **ACL-02-109:** `When an organization owner blocks a user, the system shall allow recording an optional note describing the reason for the block.`

### Unwanted Behaviour Requirements (Blocking Restrictions)

- **ACL-02-301:** `If a user or organization owner attempts to block an admin user, then the block record may be created, but enforcement is skipped for admin users (admins are immune to the effects of blocking; IsUserBlockedBy returns false for admins).`
- **ACL-02-302:** `If a user or organization owner attempts to block an organization (rather than an individual user), then the system shall deny the block operation (organizations cannot be blocked).`
- **ACL-02-303:** `If a non-owner user attempts to block or unblock a user at the organization level, then the system shall deny the operation.`

---

## 3. GPG Keys (ACL-03)

**User Story:** As a developer, I want to register my GPG public key so that my signed commits and tags are verified and display a "Verified" badge.

### Ubiquitous Requirements (GPG Key Properties)

- **ACL-03-001:** `The system shall allow each user to register zero or more GPG public keys.`
- **ACL-03-002:** `The system shall parse and store the key ID (64-bit short and long form), fingerprint, creation timestamp, and expiry for each registered GPG key.`
- **ACL-03-003:** `The system shall use a user's registered GPG keys to verify PGP signatures on Git commits and tags attributed to that user.`
- **ACL-03-004:** `The system shall display a "Verified" badge on commit and tag views whose signature validates against a registered GPG key.`

### Event-Driven Requirements (GPG Key Workflow)

- **ACL-03-101:** `When a user submits an ASCII-armored GPG public key block, the system shall parse the key, extract its metadata, and store it.`
- **ACL-03-103:** `When a commit or tag with a PGP signature is rendered, the system shall attempt to verify the signature against the purported author's registered GPG keys.`
- **ACL-03-104:** `When a user deletes a registered GPG key, the system shall remove the key immediately and cease to verify future commits with that key.`
- **ACL-03-105:** `When a commit is pushed whose signature fails verification, the system shall render an "Unverified" indicator on the commit view.`

### Optional Feature Requirements (GPG Trust Models)

- **ACL-03-201:** `Where the repository trust model is "committer", the system shall only mark commits as verified when the committer matches the GPG key identity and the pusher.`
- **ACL-03-202:** `Where the repository trust model is "collaborator", the system shall only mark commits as verified when the signer is a collaborator on the repository.`
- **ACL-03-203:** `Where the repository trust model is "collaborator-committer", the system shall apply both committer and collaborator constraints.`

### Unwanted Behaviour Requirements (GPG Key Errors)

- **ACL-03-301:** `If a user submits a malformed GPG key block, then the system shall reject the submission with a parse error.`
- **ACL-03-302:** `If a submitted GPG key is already registered to another user, then the system shall reject the registration.`
- **ACL-03-303:** `If a registered GPG key has expired, then the system shall mark commits signed after expiry as unverified.`

---

## 4. Restricted Users (ACL-04)

**User Story:** As an administrator, I want to mark certain users as restricted so that they can only access explicitly granted resources, not the full instance-wide catalog.

### State-Driven Requirements (Restricted User Enforcement)

- **ACL-04-701:** `While a user is marked as restricted, the system shall limit the user's visibility to only repositories and teams to which they have been explicitly granted access.`
- **ACL-04-702:** `While a user is marked as restricted, the system shall only allow access to repositories where the user is an explicit collaborator or team member.`

---

## 5. Organizations (ACL-05)

**User Story:** As a user, I want to create and manage organizations so that multiple collaborators can work under a shared group identity with structured access control.

### Ubiquitous Requirements (Organization Properties)

- **ACL-05-001:** `The system shall represent organizations as special user accounts of type UserTypeOrganization.`
- **ACL-05-002:** `The system shall assign each organization a unique name following the same validation rules as usernames (no reserved names, no ".git"/".wiki"/".rss"/".atom" suffixes).`
- **ACL-05-003:** `The system shall support three organization visibility levels: public (visible to everyone), limited (visible to authenticated users only), and private (visible to members only).`
- **ACL-05-004:** `The system shall create a default "Owners" team with full permissions when a new organization is created.`
- **ACL-05-005:** `The system shall track a configurable maximum repository creation limit (MaxRepoCreation) per organization, where -1 means unlimited.`
- **ACL-05-006:** `The system shall store the following profile properties per organization: name, full name, email, description, website, and location.`
- **ACL-05-007:** `The system shall support avatar management for organizations via upload, URL, or generated default.`
- **ACL-05-008:** `The system shall support the RepoAdminChangeTeamAccess flag that controls whether repository admins can change team access settings.`

### Event-Driven Requirements (Organization Lifecycle)

- **ACL-05-101:** `When a user with CanCreateOrganization permission submits a valid organization creation form, the system shall create the organization and add the creator to the default "Owners" team.`
- **ACL-05-102:** `When an organization is renamed, the system shall create a persistent redirect from the old name to the new name.`
- **ACL-05-103:** `When an organization is renamed, the system shall update all repository paths under that organization to reflect the new name.`
- **ACL-05-104:** `When an organization owner confirms deletion, the system shall remove the organization and all associated data.`
- **ACL-05-105:** `When an organization's visibility is changed, the system shall apply the visibility change to the organization profile and member listings.`
- **ACL-05-106:** `When a request targets an organization's former name after a rename, the system shall redirect to the current organization name.`

### Optional Feature Requirements (Organization Creation Control)

- **ACL-05-201:** `Where the admin setting DISABLE_REGULAR_ORG_CREATION is enabled, the system shall restrict organization creation to administrators only.`

### Unwanted Behaviour Requirements (Organization Errors)

- **ACL-05-301:** `If a user without CanCreateOrganization permission attempts to create an organization, then the system shall deny the creation.`
- **ACL-05-302:** `If an organization creation is submitted with a reserved name, then the system shall reject the creation with a validation error.`
- **ACL-05-303:** `If an organization deletion is attempted while the organization still owns repositories or packages, then the system shall deny the deletion.`
- **ACL-05-304:** `If a non-owner user attempts to modify organization settings, then the system shall deny the operation.`
- **ACL-05-305:** `If an organization rename is submitted with a name that already exists, then the system shall reject the rename.`

---

## 6. Teams (ACL-06)

**User Story:** As an organization owner, I want to create and configure teams so that I can grant granular role-based access to groups of users across repositories and feature units.

### Ubiquitous Requirements (Team Properties)

- **ACL-06-001:** `The system shall define five team permission levels: Owner, Admin, Write, Read, and None.`
- **ACL-06-002:** `The system shall enforce team name uniqueness within each organization.`
- **ACL-06-003:** `The system shall support the following ten unit-level permission targets: Code, Issues, Pull Requests, Releases, Wiki, External Wiki, External Tracker, Projects, Packages, and Actions.`
- **ACL-06-004:** `The system shall assign all unit permissions automatically to Owner-level teams.`
- **ACL-06-005:** `The system shall track an IncludesAllRepositories flag per team that grants access to all current and future organization repositories.`
- **ACL-06-006:** `The system shall track a CanCreateOrgRepo flag as a first-class per-team property that authorizes members to create repositories under the organization.`
- **ACL-06-007:** `The system shall persist per-unit team access modes on each TeamUnit record (Read, Write, or Admin) rather than a binary on/off, allowing teams to hold different access levels across units.`

### Event-Driven Requirements (Team Lifecycle)

- **ACL-06-101:** `When an organization owner creates a new team, the system shall register the team with the specified permission level, unit permissions, and repository access scope.`
- **ACL-06-102:** `When an organization owner updates a team's configuration, the system shall persist the changes and recalculate effective permissions for all team members.`
- **ACL-06-103:** `When an organization owner deletes a team, the system shall remove the team and revoke its granted permissions from all associated repositories.`
- **ACL-06-104:** `When a team has IncludesAllRepositories enabled, the system shall automatically grant the team access to any new repository created in the organization.`
- **ACL-06-105:** `When a team has CanCreateOrgRepo enabled, the system shall allow team members to create new repositories under the organization.`

### Optional Feature Requirements (Team Configuration)

- **ACL-06-201:** `Where a team is configured with specific unit permissions, the system shall restrict team member access to only those units with the specified access levels.`
- **ACL-06-202:** `Where a team is not configured as IncludesAllRepositories, the system shall require explicit repository assignment to grant team access.`

### Unwanted Behaviour Requirements (Team Errors)

- **ACL-06-301:** `If a user attempts to delete the default "Owners" team, then the system shall deny the deletion.`
- **ACL-06-302:** `If a non-owner user attempts to create or modify a team, then the system shall deny the operation.`
- **ACL-06-303:** `If a team creation is submitted with a name that already exists within the organization, then the system shall reject the creation.`
- **ACL-06-304:** `If the last owner of an organization is attempted to be removed from the Owners team, then the system shall deny the removal.`

---

## 7. Team Membership & Invites (ACL-07)

**User Story:** As an organization owner or team maintainer, I want to add and remove members from teams so that users receive the appropriate access, and I want to invite external users via email so they can join the team.

### Ubiquitous Requirements (Membership Properties)

- **ACL-07-001:** `The system shall track team membership as a relationship between user accounts and teams with a role indicator.`
- **ACL-07-002:** `The system shall allow a user to belong to multiple teams within the same organization.`
- **ACL-07-003:** `The system shall compute a user's effective organization permission as the union of all team permissions.`
- **ACL-07-004:** `The system shall generate a unique token for each team invitation email.`
- **ACL-07-005:** `The system shall allow organization members to control their own membership visibility (public or concealed).`

### Event-Driven Requirements (Membership Workflow)

- **ACL-07-101:** `When an authorized user adds a member to a team, the system shall grant that user the team's permissions immediately.`
- **ACL-07-102:** `When an authorized user removes a member from a team, the system shall revoke that team's permissions from the user.`
- **ACL-07-103:** `When a team invitation is sent to an email address, the system shall send an invitation email containing a unique token and acceptance link.`
- **ACL-07-104:** `When a recipient clicks a valid invitation link and accepts, the system shall add the user to the team and mark the invitation as accepted.`
- **ACL-07-105:** `When a user leaves an organization, the system shall remove the user from all teams within that organization.`
- **ACL-07-106:** `When an organization member updates their membership visibility, the system shall reflect the change in the organization's member listing.`

### State-Driven Requirements (Invitation Lifecycle)

- **ACL-07-702:** `While a user is a member of an organization, the system shall include that organization in the user's organization listing.`

### Unwanted Behaviour Requirements (Membership Errors)

- **ACL-07-301:** `If a user attempts to add an already-existing team member to the same team, then the system shall reject the addition.`
- **ACL-07-303:** `If a non-authorized user attempts to modify team membership, then the system shall deny the operation.`
- **ACL-07-304:** `If a user attempts to invite a blocked user to a team, then the system shall deny the invitation.`
- **ACL-07-305:** `If the last owner of an organization attempts to leave, then the system shall deny the departure.`

---

## 8. Organization Settings (ACL-08)

**User Story:** As an organization owner, I want to configure all aspects of my organization including profile, avatar, visibility, teams and unit permissions, webhooks, and OAuth2 applications so that the organization operates according to my requirements.

### Ubiquitous Requirements (Settings Categories)

- **ACL-08-001:** `The system shall restrict access to organization settings to users with owner-level permission in the organization.`
- **ACL-08-002:** `The system shall store the following profile settings per organization: name, full name, email, description, website, and location.`

### Event-Driven Requirements (Settings Updates)

- **ACL-08-101:** `When an organization owner updates profile settings, the system shall persist the changes immediately.`
- **ACL-08-102:** `When an organization owner uploads a custom avatar, the system shall store the image and associate it with the organization.`
- **ACL-08-103:** `When an organization owner deletes the custom avatar, the system shall revert to a generated default avatar.`
- **ACL-08-104:** `When an organization owner changes the organization visibility, the system shall apply the visibility change to the organization profile and cascade visibility effects to member listings.`
- **ACL-08-105:** `When an organization owner confirms deletion via the settings page, the system shall delete the organization and all associated data.`
- **ACL-08-106:** `When an organization owner configures webhooks, the system shall register the webhooks to fire on the specified events across organization repositories.`
- **ACL-08-107:** `When an organization owner registers an OAuth2 application, the system shall create the application scoped to the organization.`

### Optional Feature Requirements (Integration Settings)

- **ACL-08-201:** `Where the packages feature is enabled, the system shall allow organization owners to configure package storage and access settings.`
- **ACL-08-202:** `Where email confirmation is required, the system shall enforce email verification before applying certain organization settings changes.`

### Unwanted Behaviour Requirements (Settings Errors)

- **ACL-08-301:** `If a non-owner user attempts to access organization settings, then the system shall deny access and return a 404 or redirect.`
- **ACL-08-302:** `If an organization deletion confirmation is submitted with incorrect confirmation text, then the system shall deny the deletion.`
- **ACL-08-303:** `If an organization avatar upload exceeds the maximum allowed file size, then the system shall reject the upload with a size limit error.`
- **ACL-08-304:** `If an invalid email address is provided in organization profile settings, then the system shall reject the update with a validation error.`

---

## 9. Organization Actions Secrets & Runners (ACL-09)

**User Story:** As an organization owner, I want to configure CI/CD runners and Actions secrets at the organization scope so that all repositories in my organization can share them.

### Ubiquitous Requirements (Org Actions Scope)

- **ACL-09-001:** `The system shall scope Actions secrets to the organization level via Secret.OwnerID set to the organization's ID.`
- **ACL-09-002:** `The system shall scope Actions runners to the organization level via ActionRunner and ActionRunnerToken records using the organization's OwnerID.`

### Event-Driven Requirements (Runner & Secret Lifecycle)

- **ACL-09-101:** `When an organization owner configures a CI/CD runner, the system shall register the runner (ActionRunner/ActionRunnerToken) scoped to the organization's OwnerID for use by organization repositories.`
- **ACL-09-102:** `When an organization owner creates an Actions secret, the system shall store the secret scoped to the organization (Secret.OwnerID = org.ID).`

> **Cross-reference:** Actions workflow execution, task scheduling, and runner registration protocol are documented in Domain 05 (CI/CD & Automation). This section documents only the organization-scoped configuration of secrets and runners.

---

## 10. Organization Package Cleanup Rules (ACL-10)

**User Story:** As an organization owner, I want to configure package cleanup rules at the organization scope so that outdated package versions are automatically removed according to my retention policy.

### Ubiquitous Requirements (Cleanup Rule Properties)

- **ACL-10-001:** `The system shall support organization-scoped package cleanup rules persisted with the organization as owner.`

### Event-Driven Requirements (Cleanup Rule Lifecycle)

- **ACL-10-101:** `When an organization owner configures a package cleanup rule, the system shall persist the rule at the organization scope; the owner may also initialize or rebuild the Cargo package index for the organization.`

> **Cross-reference:** Package registry, container registry, and package type details are documented in Domain 06 (Packages & Releases). This section documents only the organization-scoped cleanup rule configuration.

---

## 11. Collaborators (ACL-11)

**User Story:** As a repository owner, I want to add collaborators with specific permission levels so that I can control who can read, write, or administer my repository.

### Ubiquitous Requirements (Collaborator Properties)

- **ACL-11-001:** `The system shall support three collaborator permission levels per repository: Read, Write, and Admin.`
- **ACL-11-002:** `The system shall list collaborators separately from organization team members.`
- **ACL-11-003:** `The system shall compute effective permissions by combining collaborator, team, and organization memberships.`

### Event-Driven Requirements (Collaborator Workflow)

- **ACL-11-101:** `When a repository owner adds a collaborator, the system shall grant the specified permission level to that user for the repository.`
- **ACL-11-102:** `When a repository owner removes a collaborator, the system shall revoke all direct repository access for that user.`
- **ACL-11-103:** `When a repository owner changes a collaborator's permission level, the system shall update the permission immediately.`
- **ACL-11-104:** `When a user is both a collaborator and a team member, the system shall grant the highest permission from all sources.`

### Unwanted Behaviour Requirements (Collaborator Errors)

- **ACL-11-301:** `If a non-owner user attempts to add a collaborator, then the system shall deny the operation.`
- **ACL-11-302:** `If a user attempts to add themselves as a collaborator, then the system shall reject the addition.`
- **ACL-11-303:** `If a collaborator is added who is blocked by the repository owner, then the system shall reject the addition.`

---

## 12. Deploy Keys (ACL-12)

**User Story:** As a repository administrator, I want to register per-repository SSH public keys so that CI/CD systems can clone or push without sharing a user account's keys.

### Ubiquitous Requirements (Deploy Key Properties)

- **ACL-12-001:** `The system shall scope each deploy key to exactly one repository.`
- **ACL-12-002:** `The system shall allow each deploy key to be marked read-only or read-write.`
- **ACL-12-003:** `The system shall store the SSH public key, fingerprint, title, creation timestamp, and associated repository for each deploy key.`
- **ACL-12-004:** `The system shall allow the same SSH public key to be registered as a deploy key on multiple repositories without conflict.`

### Event-Driven Requirements (Deploy Key Workflow)

- **ACL-12-101:** `When an admin submits an SSH public key as a deploy key for a repository, the system shall store the key and compute its fingerprint.`
- **ACL-12-102:** `When an SSH connection authenticates with a deploy key, the system shall grant access scoped to the associated repository only.`
- **ACL-12-103:** `When a deploy key marked read-only attempts a git-receive-pack, the system shall reject the push.`
- **ACL-12-104:** `When an admin deletes a deploy key, the system shall remove the key and reject future SSH authentications using that key.`

### Optional Feature Requirements (Deploy Key Sources)

- **ACL-12-201:** `Where a deploy key title is provided at creation, the system shall display the title in the deploy-key management list.`
- **ACL-12-202:** `Where an admin creates a deploy key by generating a new SSH keypair, the system shall display the private key once for download and store only the public key.`

### Unwanted Behaviour Requirements (Deploy Key Errors)

- **ACL-12-301:** `If a deploy key attempts to access a repository other than its scoped repository, then the system shall reject the operation.`
- **ACL-12-302:** `If a non-admin user attempts to add or delete a deploy key, then the system shall deny the operation.`
- **ACL-12-303:** `If a submitted SSH public key is malformed, then the system shall reject the deploy-key creation with a parse error.`

---

## Configuration Reference

The following INI sections configure access control and organization behaviors. Add to or modify these via `custom/conf/app.ini` or environment overrides using the `GITEA__section__key` format.

### [service] Section (Organization)

- **DEFAULT_ORG_MEMBER_VISIBLE**: Whether organization members are listed as visible by default on the org's member page (members may still override their own visibility). Default `false`. Source: `modules/setting/service.go:81,224`.
- **DEFAULT_ORG_VISIBILITY**: Default visibility for newly created organizations (`public`, `limited`, `private`). Default `public`. Source: `modules/setting/service.go:222`.

### [admin] Section (Organization Creation)

- **DISABLE_REGULAR_ORG_CREATION**: When true, restricts organization creation to administrators only. Default `false`. Source: `modules/setting/admin.go:20`.

### [security] Section (RBAC & Access Control)

- **REVERSE_PROXY_AUTHENTICATION_USER**: Header name (default `X-Webauth-User`) carrying the authenticated username from a trusted reverse proxy.
- **REVERSE_PROXY_AUTHENTICATION_EMAIL**: Header name carrying the authenticated user's email.
- **REVERSE_PROXY_AUTHENTICATION_FULL_NAME**: Header name carrying the authenticated user's full name.
- **REVERSE_PROXY_LIMIT**: Number of trusted proxy hops to traverse when resolving client IP.
- **REVERSE_PROXY_TRUSTED_PROXIES**: Comma-separated CIDR ranges of trusted reverse proxies.
- **IMPORT_LOCAL_PATHS**: Allow importing of local paths as repositories when enabled (default false).
- **DISABLE_GIT_HOOKS**: Disable creation of user-owned git hooks to prevent arbitrary code execution (default true).
- **DISABLE_WEBHOOKS**: Disable all webhooks across the instance when enabled (default false).
- **SUCCESSFUL_TOKENS_CACHE_SIZE**: Number of successfully validated tokens to cache for fast lookup (default 20).

> **Cross-reference:** Password, session, OAuth2, and OpenID configuration keys live in Domain 01 (Identity & Authentication).

---

## Business Rules

- **BR-02-001:** Permission inheritance is strict: Owner > Admin > Write > Read > None
- **BR-02-002:** Effective permission is the highest level from all sources (collaborator, team, org membership)
- **BR-02-003:** Restricted users can only see repositories they are explicitly added to
- **BR-02-004:** User blocking is directional (blocker -> blockee), not reciprocal
- **BR-02-005:** Organization names must be unique across the entire instance (case-insensitive) and follow username validation rules
- **BR-02-006:** The default "Owners" team is created automatically on organization creation and cannot be deleted
- **BR-02-007:** An organization must have at least one owner at all times; the last owner cannot be removed
- **BR-02-008:** Team permission levels follow a strict hierarchy: Owner > Admin > Write > Read > None
- **BR-02-009:** A user's effective organization permission is the union of all team permissions across teams they belong to
- **BR-02-010:** Organization visibility determines which users can see the organization in listings and search results
- **BR-02-011:** MaxRepoCreation of -1 means unlimited repository creation; 0 means no creation; positive values set a numeric limit
- **BR-02-012:** Organization-level blocking is directional and does not block the user instance-wide
- **BR-02-013:** RepoAdminChangeTeamAccess determines whether repository administrators can modify team access for their repositories
- **BR-02-014:** Deploy keys are scoped to exactly one repository and may be reused across repositories without conflict
- **BR-02-015:** Collaborator permissions (Read, Write, Admin) are per-repository and combine with team permissions via highest-wins

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Restricted user browses explore page | Only explicitly granted repositories are visible |
| Blocked user creates PR via fork | Reject PR creation to blocker's repos |
| Blocked user invited to team | Deny invitation |
| Block an admin user | Block record may be created; enforcement skipped (admin not actually blocked) |
| Block an organization account | Deny block operation (only individual users can be blocked) |
| Delete organization with existing repositories | Deny deletion; require transfer or deletion of all repos first |
| Delete organization with existing packages | Deny deletion; require removal of all packages first |
| Remove last owner from Owners team | Deny removal; enforce at least one owner |
| Rename organization to existing name | Reject rename with conflict error |
| Add user already on the team | Reject addition with duplicate membership error |
| Team invitation token does not expire | Tokens are accepted at any time after issuance until manually revoked |
| Access org settings as non-owner | Return 404 or redirect away |
| Org deletion with wrong confirmation text | Deny deletion |
| Team with IncludesAllRepositories and new repo created | Automatically grant team access to the new repository |
| User leaves organization while last owner | Deny departure |
| Visibility change to Private | Hide org from non-member listings and search |
| Add blocked user as collaborator | Reject the addition |
| Deploy key attempts cross-repo access | Reject the operation (scoped to one repo) |
| Malformed GPG key submitted | Reject with parse error |
| Expired GPG key signs commit | Mark commit as unverified |

## Success Criteria

- Permission computation for any resource completes in under 50ms
- Organization creation completes within 3 seconds for standard configuration
- Team permission changes propagate to all members within 1 second
- Organization rename creates a working redirect accessible indefinitely
- Team invitation emails are sent within 30 seconds of creation
- Blocking a user takes effect immediately across all resources
- Organization settings changes persist immediately and reflect on next page load
- Permission computation for any organization resource completes in under 50ms
- Permission computation for any repository action completes in under 50ms
- GPG key registration completes within 2 seconds

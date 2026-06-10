# 01 — Identity & Access

## 1. User Accounts

### What
Full user lifecycle management: registration, profile, avatar, email, password, blocking, badges, rename/redirect, deletion. Supports individual, organization, bot, and remote user types.

### Registration & Login
- Registration: manual with email confirmation, admin-created, OAuth2 auto-registration, external auth sources (LDAP, SMTP, etc.)
- Login: plain password, LDAP (BindDN/simple), SMTP, PAM, OAuth2, SSPI/SPNEGO, OpenID
- Username validation: normalization, reserved names check, character filtering
- Account controls: activation, prohibition, restricted users

### Profile Management
- Properties: name, full name, email, location, website, description
- Privacy: keep email private, keep activity private
- Notification preferences: enabled, on-mention only, disabled, all including own
- Visibility: Public, Limited (authenticated), Private (connections only)
- Language preference, theme preference

### Avatar Management
- Sources: Gravatar, Libravatar, local upload, auto-generated (identicon/initials)
- Hash-based storage with MD5 validation
- Dynamic sizing

### Email Management
- Multiple emails per account (primary + secondary)
- Per-email activation state with code-based verification and time limits
- Email domain allow/block lists
- Privacy: placeholder emails when kept private

### Password Management
- Configurable hashing (default argon2)
- Configurable complexity rules
- Time-limited reset codes
- Force password change flag (must_change_password)
- Optional Pwned password check (HaveIBeenPwned)

### OpenID
- Per-user OpenID URI list
- Normalized URI storage
- Visibility toggle in profile

### User Blocking
- User-to-user blocking with optional note/reason
- Admins cannot be blocked
- Blocked users hidden from view

### Badges
- Badge definition: slug (unique ID), description, image URL
- Multiple badges per user
- Admin-managed

### User Redirect/Rename
- Renaming creates persistent redirect from old name
- Redirect cleanup when needed

### Account Deletion
- Soft deletion (deactivation) and hard deletion options
- Ghost and actions system users for attribution

### UI Routes
- GET /user/login — Login page
- GET /user/sign_up — Registration
- GET /user/settings — User settings
- GET /user/settings/profile — Profile
- POST /user/settings/profile — Update profile
- GET /user/settings/account — Account settings
- POST /user/settings/account — Update account
- GET /user/settings/email — Email management
- POST /user/settings/email — Add email
- POST /user/settings/email/delete — Delete email
- GET /user/settings/security — Security settings
- GET /user/settings/applications — Access tokens
- POST /user/settings/applications — Create token
- GET /user/settings/keys — SSH/GPG keys
- POST /user/settings/keys — Add key
- GET /user/settings/openid — OpenID
- GET /{username} — User profile

### API Endpoints
- GET /user — Get current user
- GET /users/{username} — Get user
- GET /users/search — Search users
- PATCH /user/settings — Update settings
- POST /admin/users — Create user (admin)
- PATCH /admin/users/{username} — Update user (admin)
- DELETE /admin/users/{username} — Delete user (admin)

### Config
- [security] INSTALL_LOCK, SECRET_KEY, LOGIN_REMEMBER_DAYS, MIN_PASSWORD_LENGTH, PASSWORD_HASH_ALGO, PASSWORD_CHECK_PWN
- [admin] DISABLE_REGULAR_ORG_CREATION, DEFAULT_EMAIL_NOTIFICATIONS, USER_DISABLED_FEATURES, EXTERNAL_USER_DISABLED_FEATURES
- [service] DISABLE_REGISTRATION, REQUIRE_SIGNIN_VIEW, ENABLE_NOTIFY_MAIL, REGISTER_EMAIL_CONFIRM

---

## 2. Authentication

### Password Authentication
- Local: salted password hash verification
- Session: remember-me with configurable duration
- Configurable hash algorithms

### OAuth2
- **Provider**: RFC 6749 compliant
  - Authorization code flow with PKCE
  - Client types: Confidential and Public
  - JWT token signing (RS256)
  - Redirect URI validation
  - Access/refresh token lifecycle
- **Consumer**: external OAuth2 providers
  - Built-in apps: git-credential-oauth, git-credential-manager, tea
  - Custom app registration
  - Auto-registration options

### Access Tokens (Personal Access Tokens)
- Scopes:
  - Activity: read:activitypub, write:activitypub
  - Admin: read:admin, write:admin
  - Misc: read:misc, write:misc
  - Notification: read:notification, write:notification
  - Organization: read:organization, write:organization
  - Package: read:package, write:package
  - Issue: read:issue, write:issue
  - Repository: read:repository, write:repository
  - User: read:user, write:user
  - Special: all, public-only, sudo
- Security: SHA-256 hashing with salt, last-8-index lookup, caching, usage tracking

### Session Management
- Database-backed sessions
- CRUD: create, read, update, destroy, regenerate
- Configurable lifetime with cleanup
- Blob storage for session data

### Two-Factor Authentication (2FA)
- TOTP (RFC 6238) with PBKDF2 key derivation
- AES-encrypted secret storage
- Emergency scratch codes
- Usage tracking
- Enable/disable management

### WebAuthn/Passkey
- WebAuthn Level 1 compliance
- Registration and authentication
- Multiple credentials per user
- Sign count tracking, clone detection

### Authentication Sources
- LDAP (BindDN): TLS, user search, attribute mapping, group membership, sync
- LDAP (simple): direct bind
- SMTP: email-based, TLS
- PAM: Unix system auth
- OAuth2: external providers
- SAML: SAML 2.0
- SPNEGO/SSPI: Windows/Kerberos
- FreeIPA: specialized LDAP

### UI Routes
- GET /user/login — Login
- GET /user/two_factor — 2FA verification
- GET /user/settings/security/two_factor — 2FA management
- GET /user/settings/security/webauthn — WebAuthn management
- GET /user/settings/applications/oauth2 — OAuth2 apps
- GET /admin/auths — Auth sources (admin)
- GET /admin/auths/new — Create auth source
- GET /admin/auths/:id — Edit auth source

### Config
- [oauth2] ENABLED, ACCESS_TOKEN_EXPIRATION_TIME, REFRESH_TOKEN_EXPIRATION_TIME, JWT_SIGNING_ALGORITHM, DEFAULT_APPLICATIONS
- [oauth2_client] ENABLE_AUTO_REGISTRATION, USERNAME, UPDATE_AVATAR, ACCOUNT_LINKING, OPENID_CONNECT_SCOPES
- [security] SUCCESSFUL_TOKENS_CACHE_SIZE

---

## 3. Authorization / RBAC

### What
Hierarchical permission system: Owner > Admin > Write > Read > None. Applied at repo, org/team, and unit levels.

### Permission Modes
- AccessModeNone (0): no access
- AccessModeRead (1): read only
- AccessModeWrite (2): read + write
- AccessModeAdmin (3): read + write + admin
- AccessModeOwner (4): full control

### Repository-Level
- Per-collaborator permissions
- Branch protection integration
- Git operation access control (SSH/HTTPS)

### Organization/Team-Level
- Org visibility: Public, Limited, Private
- Team-based access with unit-level permissions
- Admin override capabilities

### Unit-Level Permissions
- Per-repo-unit: Code, Issues, Pulls, Wiki, ExternalWiki, ExternalTracker, Projects, Packages, Actions
- Team units: specific permissions per team per unit type
- Owners bypass all restrictions

### Restricted Users
- Limited visibility based on explicit restrictions
- Can only see repos/teams they're explicitly added to

---

## 4. Federation

### ActivityPub
- Client implementation with HTTP Signature support
- Activity Streams 2.1 content types
- Signature verification for incoming requests
- Signed outgoing requests

### WebFinger
- Resource discovery: acct: and mailto: schemes
- JRD (JSON Resource Descriptor) response
- Links: profile page, avatar, ActivityPub endpoint
- Respects user privacy settings

### NodeInfo
- Schema 2.1 support
- Instance metadata for federation discovery
- Standard link locations

### Config
- [federation] ENABLED — enable federation
- [federation] GET_HEADERS, POST_HEADERS — signed headers
- [federation] DIGEST_ALGORITHM — sha-256 default

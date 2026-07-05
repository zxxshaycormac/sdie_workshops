# 01 — Identity & Access

Baseline specification of Gitea's identity, authentication, authorization, and federation subsystems. All requirements describe the current (v1.22.x) system behavior.

---

## 1. User Registration

**User Story:** As a new user, I want to create an account so that I can access Gitea's features.

### Ubiquitous Requirements (Registration Constraints)

- **UA-01-001:** `The system shall require a unique username for each user account.`
- **UA-01-002:** `The system shall normalize usernames to lowercase for comparison.`
- **UA-01-003:** `The system shall reject usernames that match reserved names including ".", "..", ".well-known", "admin", "api", "assets", "attachments", "avatar", "avatars", "captcha", "commits", "debug", "error", "explore", "favicon.ico", "ghost", "issues", "login", "manifest.json", "metrics", "milestones", "new", "notifications", "org", "pulls", "raw", "repo", "repo-avatars", "robots.txt", "search", "serviceworker.js", "ssh_info", "swagger.v1.json", "user", "v2", "gitea-actions", and reserved patterns "*.keys", "*.gpg", "*.rss", "*.atom", "*.png".`
- **UA-01-004:** `The system shall require a valid email address for account registration.`
- **UA-01-005:** `The system shall hash passwords using a configurable algorithm (default pbkdf2).`
- **UA-01-006:** `The system shall enforce a configurable minimum password length (default 8 characters).`
- **UA-01-007:** `The system shall assign user type "individual" to self-registered accounts.`

### Event-Driven Requirements (Registration Workflow)

- **UA-01-101:** `When a user submits a valid registration form, the system shall create a new user account.`
- **UA-01-102:** `When registration requires email confirmation, the system shall send an activation email with a time-limited code.`
- **UA-01-103:** `When a user clicks a valid activation link, the system shall activate the user account.`
- **UA-01-104:** `When an admin creates a user via the admin panel, the system shall create the account in active state.`
- **UA-01-105:** `When the setting REGISTER_EMAIL_CONFIRM is enabled, the system shall require email verification before allowing login.`
- **UA-01-106:** `When the setting DISABLE_REGISTRATION is true, the system shall reject all self-registration attempts.`

### Optional Feature Requirements (External Registration)

- **UA-01-201:** `Where OAuth2 auto-registration is enabled, the system shall create a local account when a new user authenticates via an external OAuth2 provider.`
- **UA-01-202:** `Where an external authentication source (LDAP, SMTP, etc.) is configured and active, the system shall authenticate users against that source.`
- **UA-01-203:** `Where password check against Pwned Passwords is enabled, the system shall warn users whose passwords appear in known breach databases.`

### Unwanted Behaviour Requirements (Registration Errors)

- **UA-01-301:** `If a user submits a username that already exists, then the system shall reject the registration with an error message.`
- **UA-01-302:** `If a user submits an email address already registered to another account, then the system shall reject the registration.`
- **UA-01-303:** `If a user submits a password that does not meet complexity requirements, then the system shall reject the password and display the requirements.`
- **UA-01-304:** `If an activation link has expired, then the system shall reject the activation and offer to resend the email.`

---

## 2. User Profile Management

**User Story:** As a user, I want to manage my profile so that others can identify and contact me.

### Ubiquitous Requirements (Profile Properties)

- **UA-02-001:** `The system shall store the following profile properties per user: username, full name, email, location, website, and description.`
- **UA-02-002:** `The system shall support three user visibility levels: public, limited (visible only to authenticated/connected users), and private (visible only to self and admins).`
- **UA-02-003:** `The system shall store the user's language and theme preferences.`
- **UA-02-004:** `The system shall support four notification preference levels: enabled, on-mention only, disabled, and all including own.`

### Event-Driven Requirements (Profile Updates)

- **UA-02-101:** `When a user updates their profile via settings, the system shall persist the changes immediately.`
- **UA-02-102:** `When a user changes their visibility to private, the system shall hide the user from search results and explore pages for non-connections.`
- **UA-02-103:** `When a user enables "keep email private", the system shall display a placeholder email address in the user's public profile.`
- **UA-02-104:** `When a user enables "keep activity private", the system shall hide the user's activity feed from non-connections.`

### State-Driven Requirements (Restricted Users)

- **UA-02-701:** `While a user is marked as restricted, the system shall limit the user's visibility to only repositories and teams to which they have been explicitly granted access.`
- **UA-02-702:** `While a user is marked as prohibited from logging in, the system shall reject all authentication attempts for that user.`

### Unwanted Behaviour Requirements (Profile Validation)

- **UA-02-301:** `If a user enters an invalid website URL, then the system shall reject the input with a validation error.`
- **UA-02-302:** `If a non-admin user attempts to change another user's profile, then the system shall deny the operation.`

---

## 3. Email Management

**User Story:** As a user, I want to manage multiple email addresses so I can control which email is primary and keep others verified.

### Ubiquitous Requirements (Email Constraints)

- **UA-03-001:** `The system shall allow each user account to have multiple email addresses.`
- **UA-03-002:** `The system shall designate exactly one email address as primary per account.`
- **UA-03-003:** `The system shall track activation state independently for each email address.`

### Event-Driven Requirements (Email Workflow)

- **UA-03-101:** `When a user adds a new email address, the system shall send a verification email with a time-limited activation code.`
- **UA-03-102:** `When a user clicks a valid email verification link, the system shall mark the email as activated.`
- **UA-03-103:** `When a user sets a verified email as primary, the system shall update the account's primary email.`
- **UA-03-104:** `When a user deletes an email address, the system shall remove it from the account, unless it is the last remaining email.`

### Optional Feature Requirements (Email Domain Control)

- **UA-03-201:** `Where email domain allow lists are configured, the system shall reject registration with email addresses from non-allowed domains.`
- **UA-03-202:** `Where email domain block lists are configured, the system shall reject registration with email addresses from blocked domains.`

### Unwanted Behaviour Requirements (Email Errors)

- **UA-03-301:** `If a user adds an email address already registered to another account, then the system shall reject the addition.`
- **UA-03-302:** `If a user attempts to delete their only email address, then the system shall reject the deletion.`
- **UA-03-303:** `If an email verification code has expired, then the system shall reject the verification and offer to resend.`

---

## 4. Password Management

**User Story:** As a user, I want to reset my password if I forget it, and as an admin, I want to enforce password policies.

### Ubiquitous Requirements (Password Constraints)

- **UA-04-001:** `The system shall hash all stored passwords using a configurable algorithm (argon2, bcrypt, scrypt, pbkdf2, or sha256).`
- **UA-04-002:** `The system shall enforce configurable password complexity rules (minimum length, uppercase, lowercase, numbers, special characters).`

### Event-Driven Requirements (Password Workflow)

- **UA-04-101:** `When a user requests a password reset, the system shall send a time-limited reset code to the user's primary email.`
- **UA-04-102:** `When a user submits a valid reset code with a new password, the system shall update the user's password.`
- **UA-04-103:** `When an admin sets the must_change_password flag on a user, the system shall require the user to change their password at next login.`
- **UA-04-104:** `When a user changes their password, the system shall invalidate all active sessions for that user.`

### Unwanted Behaviour Requirements (Password Errors)

- **UA-04-301:** `If a password reset code has expired, then the system shall reject the reset and offer to send a new code.`
- **UA-04-302:** `If a user submits a password that does not meet complexity requirements, then the system shall reject the password with specific requirement feedback.`
- **UA-04-303:** `If a user's current password is incorrect during a password change, then the system shall reject the change.`

---

## 5. User Blocking

**User Story:** As a user, I want to block other users so they cannot interact with me or my repositories.

### Event-Driven Requirements (Blocking Workflow)

- **UA-05-101:** `When a user blocks another user, the system shall prevent the blocked user from commenting on the blocker's issues and PRs.`
- **UA-05-102:** `When a user blocks another user, the system shall prevent the blocked user from creating PRs to the blocker's repositories.`
- **UA-05-103:** `When a user blocks another user, the system shall prevent the blocked user from being added to teams the blocker manages.`
- **UA-05-104:** `When a user blocks another user, the system shall allow the blocker to record an optional note describing the reason.`

### Unwanted Behaviour Requirements (Blocking Restrictions)

- **UA-05-301:** `If a user attempts to block an admin user, then the block record may be created, but enforcement is skipped for admin users (admins are immune to the effects of blocking).`
- **UA-05-302:** `If a user attempts to block an organization, then the system shall deny the block operation (organizations cannot be blocked).`

---

## 6. User Rename & Deletion

**User Story:** As a user, I want to rename my account, and as an admin, I want to delete accounts that are no longer needed.

### Event-Driven Requirements (Rename Workflow)

- **UA-06-101:** `When a user renames their account, the system shall create a persistent redirect from the old username to the new username.`
- **UA-06-102:** `When a user renames their account, the system shall update all repository paths to reflect the new username.`
- **UA-06-103:** `When a request targets a redirected username, the system shall redirect to the current username.`

### Event-Driven Requirements (Deletion Workflow)

- **UA-06-201:** `When an admin deletes a user account, the system shall transfer the user's repository ownership to a ghost user for attribution preservation.`
- **UA-06-202:** `When an admin deletes a user account, the system shall remove the user from all team memberships.`
- **UA-06-203:** `When an admin deletes a user account, the system shall remove all SSH keys, GPG keys, and access tokens associated with the account.`

### Unwanted Behaviour Requirements (Deletion Safeguards)

- **UA-06-301:** `If an admin attempts to delete the last admin user, then the system shall deny the deletion.`
- **UA-06-302:** `If a non-admin user attempts to delete a user account, then the system shall deny the operation.`

---

## 7. Password Authentication & Sessions

**User Story:** As a user, I want to log in with my credentials and maintain a session so I don't have to re-authenticate on every request.

### Ubiquitous Requirements (Session Properties)

- **AUTH-01-001:** `The system shall store sessions in the database with configurable lifetime.`
- **AUTH-01-002:** `The system shall support session regeneration on privilege level changes.`

### Event-Driven Requirements (Login Workflow)

- **AUTH-01-101:** `When a user submits valid credentials, the system shall create an authenticated session.`
- **AUTH-01-102:** `When a user selects "remember me" during login, the system shall extend the session duration to the configured LOGIN_REMEMBER_DAYS value.`
- **AUTH-01-103:** `When a user logs out, the system shall destroy the current session.`
- **AUTH-01-104:** `When a user logs out from all devices, the system shall destroy all sessions associated with the user.`

### State-Driven Requirements (Session Lifecycle)

- **AUTH-01-701:** `While a session is valid, the system shall allow the user to perform actions within their permission scope.`
- **AUTH-01-702:** `While a session has expired, the system shall redirect the user to the login page on the next request.`

### Unwanted Behaviour Requirements (Login Errors)

- **AUTH-01-301:** `If a user submits invalid credentials, then the system shall reject the login attempt.`
- **AUTH-01-302:** `If a user account is deactivated, then the system shall reject all login attempts with an appropriate message.`
- **AUTH-01-303:** `If a user with 2FA enabled does not provide a valid TOTP code, then the system shall reject the login.`

---

## 8. OAuth2 Provider

**User Story:** As a developer, I want to register an OAuth2 application so that external tools can access Gitea on my behalf.

### Ubiquitous Requirements (OAuth2 Properties)

- **AUTH-02-001:** `The system shall implement RFC 6749 OAuth2 authorization code flow with PKCE support.`
- **AUTH-02-002:** `The system shall sign JWT tokens using a configurable algorithm (default RS256).`
- **AUTH-02-003:** `The system shall support two client types: confidential and public.`

### Event-Driven Requirements (OAuth2 Workflow)

- **AUTH-02-101:** `When a user registers an OAuth2 application, the system shall generate a client ID and client secret.`
- **AUTH-02-102:** `When a client presents a valid authorization code, the system shall issue an access token and refresh token.`
- **AUTH-02-103:** `When a client presents a valid refresh token, the system shall issue a new access token.`
- **AUTH-02-104:** `When an access token expires, the system shall reject API requests using that token.`
- **AUTH-02-105:** `When a user revokes an OAuth2 grant, the system shall invalidate all tokens issued under that grant.`

### Optional Feature Requirements (OAuth2 Configuration)

- **AUTH-02-201:** `Where built-in OAuth2 applications are configured (git-credential-oauth, git-credential-manager, tea), the system shall pre-register these applications at startup.`
- **AUTH-02-202:** `Where OAuth2 is disabled via configuration, the system shall reject all OAuth2 authorization and token requests.`

### Unwanted Behaviour Requirements (OAuth2 Errors)

- **AUTH-02-301:** `If a client submits an invalid redirect URI, then the system shall reject the authorization request.`
- **AUTH-02-302:** `If a client submits an expired or invalid authorization code, then the system shall reject the token request.`
- **AUTH-02-303:** `If a client submits credentials for a locked built-in application, then the system shall reject modification attempts.`

---

## 9. Access Tokens

**User Story:** As a user, I want to create personal access tokens with specific scopes so that I can authenticate API requests without sharing my password.

### Ubiquitous Requirements (Token Properties)

- **AUTH-03-001:** `The system shall hash stored tokens using PBKDF2-HMAC-SHA256 with salt.`
- **AUTH-03-002:** `The system shall store the last 8 characters of each token in plaintext for fast lookup.`
- **AUTH-03-003:** `The system shall cache successfully validated tokens up to a configurable cache size.`

### Event-Driven Requirements (Token Workflow)

- **AUTH-03-101:** `When a user creates a personal access token, the system shall generate a secure random token and display it once.`
- **AUTH-03-102:** `When a token is created with specific scopes, the system shall restrict the token's API access to those scopes only.`
- **AUTH-03-103:** `When a user deletes a token, the system shall invalidate the token immediately and remove it from the cache.`

### Unwanted Behaviour Requirements (Token Errors)

- **AUTH-03-301:** `If an API request presents an invalid token, then the system shall reject the request with 401 Unauthorized.`
- **AUTH-03-302:** `If an API request presents a valid token lacking the required scope, then the system shall reject the request with 403 Forbidden.`
- **AUTH-03-303:** `If a token is used after deletion, then the system shall reject the request.`

---

## 10. Two-Factor Authentication (2FA)

**User Story:** As a user, I want to enable 2FA so that my account is protected even if my password is compromised.

### Ubiquitous Requirements (2FA Properties)

- **AUTH-04-001:** `The system shall implement TOTP (RFC 6238) for two-factor authentication.`
- **AUTH-04-002:** `The system shall encrypt TOTP secrets using AES with a key derived from the instance SecretKey via MD5.` NOTE: MD5-based key derivation is a known weakness (collision-prone, no salt/stretching); see `models/auth/twofactor.go:95-98` (`getEncryptionKey` = `md5.Sum(SecretKey)`).
- **AUTH-04-003:** `The system shall generate emergency scratch codes when 2FA is enabled.`

### Event-Driven Requirements (2FA Workflow)

- **AUTH-04-101:** `When a user enables 2FA, the system shall generate a TOTP secret and display a QR code for authenticator app setup.`
- **AUTH-04-102:** `When a user with 2FA enabled logs in, the system shall require a valid TOTP code after password verification.`
- **AUTH-04-103:** `When a user enters a valid scratch code, the system shall accept it as a valid second factor and consume that code.`
- **AUTH-04-104:** `When an admin resets a user's 2FA, the system shall disable 2FA (scratch codes consumed; new ones generated only on re-enrollment).`

### Unwanted Behaviour Requirements (2FA Errors)

- **AUTH-04-301:** `If a user submits an invalid TOTP code, then the system shall reject the login attempt.`
- **AUTH-04-302:** `If a user submits an already-consumed scratch code, then the system shall reject the code.`
- **AUTH-04-303:** `If a user attempts to enable 2FA without verifying a valid TOTP code, then the system shall reject the setup.`

---

## 11. WebAuthn / Passkey

**User Story:** As a user, I want to use a hardware security key or passkey for passwordless authentication.

### Ubiquitous Requirements (WebAuthn Properties)

- **AUTH-05-001:** `The system shall implement WebAuthn Level 1 for passwordless authentication.`
- **AUTH-05-002:** `The system shall allow multiple WebAuthn credentials per user.`

### Event-Driven Requirements (WebAuthn Workflow)

- **AUTH-05-101:** `When a user registers a WebAuthn credential, the system shall store the credential with sign count tracking.`
- **AUTH-05-102:** `When a user authenticates with a WebAuthn credential, the system shall verify the sign count to detect credential cloning.`
- **AUTH-05-103:** `When a user deletes a WebAuthn credential, the system shall remove it and prevent future authentication with that credential.`

### Unwanted Behaviour Requirements (WebAuthn Errors)

- **AUTH-05-301:** `If a WebAuthn authentication attempt has a lower sign count than the stored value, then the system shall flag the credential as potentially cloned.`
- **AUTH-05-302:** `If a user has no remaining WebAuthn credentials, then the system shall fall back to password-based authentication.`

---

## 12. External Authentication Sources

**User Story:** As an admin, I want to configure external authentication sources (LDAP, SMTP, etc.) so that users can log in with their existing corporate credentials.

### Ubiquitous Requirements (Auth Source Properties)

- **AUTH-06-001:** `The system shall support the following authentication source types: Plain (database), LDAP (BindDN), LDAP (simple auth / direct bind), SMTP, PAM, OAuth2, and SPNEGO/SSPI.`
- **AUTH-06-002:** `The system shall store authentication source configurations with activation state, display order, and per-source TLS configuration.`

### Event-Driven Requirements (Auth Source Workflow)

- **AUTH-06-101:** `When a user authenticates via an external source for the first time, the system shall create a local user account linked to the external identity.`
- **AUTH-06-102:** `When an LDAP source has synchronization enabled, the system shall periodically sync user data from the LDAP directory.`
- **AUTH-06-103:** `When an admin creates a new authentication source, the system shall make it available for login immediately upon activation.`

### Optional Feature Requirements (Auth Source Features)

- **AUTH-06-201:** `Where an LDAP source is configured with group-to-team mapping, the system shall assign users to teams based on their LDAP group membership.`
- **AUTH-06-202:** `Where an LDAP source is configured with admin group filtering, the system shall grant admin privileges to users in the specified groups.`
- **AUTH-06-203:** `Where an authentication source is configured to skip local 2FA, the system shall bypass 2FA for users authenticating through that source.`
- **AUTH-06-204:** `Where an LDAP source is configured with restricted group filtering, the system shall mark users in those groups as restricted.`

### Unwanted Behaviour Requirements (Auth Source Errors)

- **AUTH-06-301:** `If an external authentication source is unreachable, then the system shall reject login attempts against that source with an error.`
- **AUTH-06-302:** `If an LDAP sync operation fails, then the system shall log the error and continue with the existing user data.`

---

## 13. Permission System (RBAC)

**User Story:** As a system administrator, I want fine-grained access control so that users only access resources they are authorized for.

### Ubiquitous Requirements (Permission Hierarchy)

- **RBAC-01-001:** `The system shall define five permission levels in strict hierarchy: None (0), Read (1), Write (2), Admin (3), Owner (4).`
- **RBAC-01-002:** `The system shall enforce permission inheritance such that each higher level includes all capabilities of lower levels.`
- **RBAC-01-003:** `The system shall apply permissions at three scopes: repository (collaborator), organization/team, and unit (feature-level).`

### Event-Driven Requirements (Permission Evaluation)

- **RBAC-01-101:** `When a user accesses a repository, the system shall compute the effective permission from collaborator, team, and organization memberships.`
- **RBAC-01-102:** `When a user accesses a feature unit (code, issues, pulls, wiki, projects, packages, actions), the system shall check unit-level permissions for that user.`
- **RBAC-01-103:** `When an owner-level user accesses any resource in their scope, the system shall grant full access regardless of other restrictions.`

### Optional Feature Requirements (Unit Permissions)

- **RBAC-01-201:** `Where team unit permissions are configured, the system shall restrict team members to the specified unit types with the specified access levels.`
- **RBAC-01-202:** `Where a repository unit is disabled, the system shall hide that feature from all users regardless of permission level.`

### State-Driven Requirements (Restricted Users)

- **RBAC-01-701:** `While a user is marked as restricted, the system shall only allow access to repositories where the user is an explicit collaborator or team member.`

---

## 14. Federation (ActivityPub, WebFinger, NodeInfo)

**User Story:** As a federated user, I want my Gitea profile to be discoverable by users on other ActivityPub-compatible platforms.

### Ubiquitous Requirements (Federation Properties)

- **FED-01-001:** `The system shall implement ActivityPub client functionality with HTTP Signature support.`
- **FED-01-002:** `The system shall use Activity Streams 2.1 content types for federation messages.`
- **FED-01-003:** `The system shall support NodeInfo 2.1 schema for instance metadata publishing.`

### Event-Driven Requirements (Federation Workflow)

- **FED-01-101:** `When a remote platform sends a WebFinger query for a local user, the system shall return a JRD response with profile, avatar, and ActivityPub links.`
- **FED-01-102:** `When the system sends outgoing federation requests, the system shall sign the requests with HTTP Signatures using configured headers.`
- **FED-01-103:** `When the system receives incoming signed requests, the system shall verify the HTTP Signature before processing.`

### Optional Feature Requirements (Federation Configuration)

- **FED-01-201:** `Where federation is enabled, the system shall expose ActivityPub, WebFinger, and NodeInfo endpoints.`
- **FED-01-202:** `Where federation is disabled, the system shall return 404 for all federation endpoints.`

### Unwanted Behaviour Requirements (Federation Errors)

- **FED-01-301:** `If a WebFinger query targets a private user, then the system shall return a 404 response.`
- **FED-01-302:** `If an incoming federation request has an invalid signature, then the system shall reject the request with 401.`

---

## 15. Activity Heatmap (UA-07)

**User Story:** As a user, I want a visual contribution calendar on my profile so that my activity over the past year is visible at a glance.

### Ubiquitous Requirements (Heatmap Properties)

- **UA-07-001:** `The system shall track contribution events per user across repositories, issues, pull requests, and comments.`
- **UA-07-002:** `The system shall aggregate contribution counts per 15-minute bucket (to permit client-side timezone rendering).`
- **UA-07-003:** `The system shall expose heatmap data for the trailing 365-day window relative to the request time.`
- **UA-07-004:** `The system shall render the heatmap on the user profile page when the profile is viewable by the requester.`

### Event-Driven Requirements (Heatmap Workflow)

- **UA-07-101:** `When a user performs a contribution action (push, issue create, PR create/review, comment), the system shall increment the corresponding day's contribution count for that user.`
- **UA-07-102:** `When a viewer requests a user's profile, the system shall return the heatmap data as a structured JSON payload keyed by timestamp and count.`
- **UA-07-103:** `When a contribution action is reversed (issue deleted, PR reverted), the system shall decrement the corresponding day's contribution count.`
- **UA-07-104:** `When a user performs an action on a private repository, the system shall include the contribution in the owner's heatmap but exclude it from viewers without access.`

### Optional Feature Requirements (Heatmap Configuration)

- **UA-07-201:** `Where ENABLE_USER_HEATMAP is disabled in configuration, the system shall hide the heatmap from all user profiles.`
- **UA-07-202:** `Where a user's visibility is set to private, the system shall hide the heatmap from viewers who are not authorized connections.`
- **UA-07-203:** `Where a user has enabled "keep activity private", the system shall hide the heatmap from non-connections.`

### Unwanted Behaviour Requirements (Heatmap Errors)

- **UA-07-301:** `If a viewer lacks permission to view a user's profile, then the system shall return an empty heatmap dataset rather than a 403.`
- **UA-07-302:** `If the requested time range exceeds the maximum supported 365-day window, then the system shall clamp the request to the maximum window.`

---

## 16. GPG Keys (AUTH-07)

**User Story:** As a developer, I want to register my GPG public key so that my signed commits and tags are verified and display a "Verified" badge.

### Ubiquitous Requirements (GPG Key Properties)

- **AUTH-07-001:** `The system shall allow each user to register zero or more GPG public keys.`
- **AUTH-07-002:** `The system shall parse and store the key ID (64-bit short and long form), fingerprint, creation timestamp, and expiry for each registered GPG key.`
- **AUTH-07-003:** `The system shall use a user's registered GPG keys to verify PGP signatures on Git commits and tags attributed to that user.`
- **AUTH-07-004:** `The system shall display a "Verified" badge on commit and tag views whose signature validates against a registered GPG key.`

### Event-Driven Requirements (GPG Key Workflow)

- **AUTH-07-101:** `When a user submits an ASCII-armored GPG public key block, the system shall parse the key, extract its metadata, and store it.`
- **AUTH-07-103:** `When a commit or tag with a PGP signature is rendered, the system shall attempt to verify the signature against the purported author's registered GPG keys.`
- **AUTH-07-104:** `When a user deletes a registered GPG key, the system shall remove the key immediately and cease to verify future commits with that key.`
- **AUTH-07-105:** `When a commit is pushed whose signature fails verification, the system shall render an "Unverified" indicator on the commit view.`

### Optional Feature Requirements (GPG Trust Models)

- **AUTH-07-201:** `Where the repository trust model is "committer", the system shall only mark commits as verified when the committer matches the GPG key identity and the pusher.`
- **AUTH-07-202:** `Where the repository trust model is "collaborator", the system shall only mark commits as verified when the signer is a collaborator on the repository.`
- **AUTH-07-203:** `Where the repository trust model is "collaborator-committer", the system shall apply both committer and collaborator constraints.`

### Unwanted Behaviour Requirements (GPG Key Errors)

- **AUTH-07-301:** `If a user submits a malformed GPG key block, then the system shall reject the submission with a parse error.`
- **AUTH-07-302:** `If a submitted GPG key is already registered to another user, then the system shall reject the registration.`
- **AUTH-07-303:** `If a registered GPG key has expired, then the system shall mark commits signed after expiry as unverified.`

---

## 17. OpenID Authentication (AUTH-08)

**User Story:** As a user, I want to sign in to Gitea using my OpenID identity so that I can authenticate via a decentralized identity provider.

### Ubiquitous Requirements (OpenID Properties)

- **AUTH-08-001:** `The system shall implement OpenID 2.0 consumer functionality for delegated authentication.`
- **AUTH-08-002:** `The system shall implement OpenID as a dedicated consumer flow via routers/web/auth/openid.go, storing identities in the user_open_id table, separate from the login_source auth.Type enum (NoType, Plain, LDAP, SMTP, PAM, DLDAP, OAuth2, SSPI — no OpenID entry).`
- **AUTH-08-003:** `The system shall discover the OpenID provider endpoint from the user-supplied OpenID URL via Yadis / HTML discovery.`
- **AUTH-08-004:** `The system shall request the SReg and AX attributes required for account creation (nickname, email, full name).`

### Event-Driven Requirements (OpenID Workflow)

- **AUTH-08-101:** `When a user submits an OpenID URL on the login page, the system shall perform discovery and redirect the user to the provider's authentication endpoint.`
- **AUTH-08-102:** `When the provider redirects back with a positive assertion, the system shall verify the signature and nonce.`
- **AUTH-08-103:** `When the verified OpenID matches an existing local account, the system shall authenticate the user as that account.`
- **AUTH-08-104:** `When the verified OpenID does not match an existing account and OpenID signup is enabled, the system shall create a new local account from the returned attributes.`
- **AUTH-08-105:** `When a user links an OpenID to their existing account via settings, the system shall store the OpenID URL for future sign-in.`

### Optional Feature Requirements (OpenID Configuration)

- **AUTH-08-201:** `Where ENABLE_OPENID_SIGNIN is enabled, the system shall display the OpenID login option on the login page.`
- **AUTH-08-202:** `Where ENABLE_OPENID_SIGNUP is enabled, the system shall allow new account creation via OpenID.`
- **AUTH-08-203:** `Where WHITELISTED_URIS is configured, the system shall reject OpenID URLs not matching the allowlist.`
- **AUTH-08-204:** `Where BLACKLISTED_URIS is configured, the system shall reject OpenID URLs matching the blocklist.`

### Unwanted Behaviour Requirements (OpenID Errors)

- **AUTH-08-301:** `If the OpenID URL fails discovery, then the system shall reject the login with a discovery error message.`
- **AUTH-08-302:** `If the provider returns a negative assertion, then the system shall return the user to the login page without authenticating.`
- **AUTH-08-303:** `If the assertion signature verification fails, then the system shall reject the login and log the failure.`
- **AUTH-08-304:** `If the assertion nonce has been replayed or expired, then the system shall reject the login.`

---

## Configuration Reference

The following INI sections under `[identity-access]`-relevant namespaces configure this category's behaviors. Add to or modify these via `custom/conf/app.ini` or environment overrides using the `GITEA__section__key` format.

### [security] Section

- **INSTALL_LOCK**: Prevents access to the install page after initial setup.
- **SECRET_KEY**: Per-instance secret used for session encryption; rotate to invalidate active sessions.
- **LOGIN_REMEMBER_DAYS**: Number of days a "remember me" session cookie remains valid.
- **COOKIE_REMEMBER_NAME**: Name of the remember-me cookie.
- **REVERSE_PROXY_AUTHENTICATION_USER**: Header name (default `X-Webauth-User`) carrying the authenticated username from a trusted reverse proxy.
- **REVERSE_PROXY_AUTHENTICATION_EMAIL**: Header name carrying the authenticated user's email.
- **REVERSE_PROXY_AUTHENTICATION_FULL_NAME**: Header name carrying the authenticated user's full name.
- **REVERSE_PROXY_LIMIT**: Number of trusted proxy hops to traverse when resolving client IP.
- **REVERSE_PROXY_TRUSTED_PROXIES**: Comma-separated CIDR ranges of trusted reverse proxies.
- **MIN_PASSWORD_LENGTH**: Minimum password character count (default 8).
- **PASSWORD_COMPLEXITY**: Comma-separated complexity requirements (`lower,upper,digit,spec`).
- **PASSWORD_CHECK_PWNED**: Reject passwords appearing in the HaveIBeenPwned breach database when enabled.
- **INTERNAL_TOKEN**: Secret token for internal API calls between Gitea processes; auto-generated if empty.
- **INTERNAL_TOKEN_URI**: File or command URI from which to load the internal token.
- **SUCCESSFUL_TOKENS_CACHE_SIZE**: Number of successfully validated tokens to cache for fast lookup (default 20).
- **DISABLE_QUERY_AUTH_TOKEN**: Disable legacy query-string API auth tokens (default false; will default true in future releases).
- **IMPORT_LOCAL_PATHS**: Allow importing of local paths as repositories when enabled (default false).
- **DISABLE_GIT_HOOKS**: Disable creation of user-owned git hooks to prevent arbitrary code execution (default true).
- **DISABLE_WEBHOOKS**: Disable all webhooks across the instance when enabled (default false).
- **CSRFCookieName**: Name of the CSRF protection cookie (default `_csrf`).
- **CSRF_COOKIE_HTTP_ONLY**: Mark the CSRF cookie as HTTP-only (default true).

### [session] Section

- **PROVIDER**: Session storage backend (`memory`, `file`, `redis`, `mysql`, `postgres`, `couchbase`, `memcache`, `db`).
- **PROVIDER_CONFIG**: Connection string consumed by the chosen provider.
- **COOKIE_NAME**: Name of the session cookie (default `i_like_gitea`).
- **COOKIE_SECURE**: Boolean — force `Secure` attribute on session cookies (default auto-derived from site URL scheme).
- **DOMAIN**: Cookie domain name (default empty).
- **COOKIE_PATH**: Cookie path (default `AppSubURL`, or `/`).
- **GC_INTERVAL_TIME**: Garbage collection interval for expired sessions (seconds, default 86400).
- **SESSION_LIFE_TIME**: Maximum session lifetime in seconds (default 86400).
- **SAME_SITE**: SameSite attribute on session cookies (`lax`, `strict`, `none`; default `lax`).

> NOTE: The CSRF cookie name (`CSRFCookieName`, default `_csrf`) is configured under the `[security]` section, not `[session]`.

### [oauth2] Section

- **ENABLED**: Enable the OAuth2 provider (Authorization Server) role (default true). `ENABLE` is a deprecated alias renamed to `ENABLED`.
- **JWT_SECRET**: Secret used to sign OAuth2 JWTs; auto-generated if empty; rotate to invalidate outstanding tokens.
- **JWT_SIGNING_ALGORITHM**: Algorithm for JWT signing (default `RS256`).
- **JWT_SIGNING_PRIVATE_KEY_FILE**: Path to the JWT signing private key file (default `jwt/private.pem` under AppDataPath).
- **ACCESS_TOKEN_EXPIRATION_TIME**: Access token TTL in seconds (default 3600).
- **REFRESH_TOKEN_EXPIRATION_TIME**: Refresh token TTL in hours (default 730).
- **INVALIDATE_REFRESH_TOKENS**: Rotate refresh tokens on each use when enabled (default false).
- **MAX_TOKEN_LENGTH**: Maximum length of generated access tokens (default 32767).
- **DEFAULT_APPLICATIONS**: Comma-separated list of built-in OAuth2 applications to pre-register at startup (default `git-credential-oauth`, `git-credential-manager`, `tea`; empty string disables all).

### [openid] Section

- **ENABLE_OPENID_SIGNIN**: Expose OpenID as a login option (default true).
- **ENABLE_OPENID_SIGNUP**: Allow new account creation via OpenID (default false unless `DISABLE_REGISTRATION` is true).
- **WHITELISTED_URIS**: Glob patterns restricting accepted OpenID provider URIs.
- **BLACKLISTED_URIS**: Glob patterns blocking specific OpenID provider URIs.

### [oauth2_client] Section

- **ENABLE_AUTO_REGISTRATION**: Automatically create a local account on first OAuth2 login when enabled (default false).
- **USERNAME**: Source field for generating the local username from OAuth2 data (`userid`, `nickname`, `email`, `preferred_username`; default `nickname`).
- **ACCOUNT_LINKING**: Behavior when an OAuth2 identity matches an existing local account (`disabled`, `login`, `auto`; default `login`).
- **UPDATE_AVATAR**: Update the local avatar from the OAuth2 provider on each login when enabled (default false).
- **OPENID_CONNECT_SCOPES**: Space-separated OpenID Connect scopes requested from the provider.
- **REGISTER_EMAIL_CONFIRM**: Require email confirmation for OAuth2 auto-registered accounts (inherits `[service].REGISTER_EMAIL_CONFIRM`).

### [federation] Section

- **ENABLED**: Enable ActivityPub federation and HTTP Signature support (default false).
- **SHARE_USER_STATISTICS**: Permit sharing instance-level user statistics when federation is enabled (default true).
- **MAX_SIZE**: Maximum size of incoming federation payloads in MiB (default 4).
- **ALGORITHMS**: Supported HTTP Signature algorithms (default `rsa-sha256`, `rsa-sha512`, `ed25519`).
- **DIGEST_ALGORITHM**: Digest algorithm for signed federation requests (default `SHA-256`).
- **GET_HEADERS**: Headers included in signed GET requests (default `(request-target)`, `Date`).
- **POST_HEADERS**: Headers included in signed POST requests (default `(request-target)`, `Date`, `Digest`).

---

## Business Rules

- **BR-01-001:** Usernames must be unique across the entire instance (case-insensitive)
- **BR-01-002:** Email addresses must be unique across all accounts
- **BR-01-003:** An instance must have at least one admin user at all times
- **BR-01-004:** The last admin user cannot be demoted or deleted
- **BR-01-005:** Password hashing algorithm applies instance-wide (not per-user); default is pbkdf2
- **BR-01-006:** OAuth2 access token expiration is configurable (default 3600 seconds)
- **BR-01-007:** OAuth2 refresh token expiration is configurable (default 730 hours)
- **BR-01-008:** Each user may have at most one active TOTP 2FA configuration
- **BR-01-009:** Permission inheritance is strict: Owner > Admin > Write > Read > None
- **BR-01-010:** Restricted users can only see repositories they are explicitly added to
- **BR-01-011:** User blocking is directional (blocker → blockee), not reciprocal
- **BR-01-012:** Ghost users are system accounts used for attribution after account deletion
- **BR-01-013:** Federation respects user visibility settings (private users are not discoverable)
- **BR-01-014:** Registration supports a third confirmation mode beside email-confirm: `REGISTER_MANUAL_CONFIRM` requires an admin to manually activate newly registered accounts; it is ignored when `REGISTER_EMAIL_CONFIRM` is enabled (evidence: `modules/setting/service.go:147`)

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Username collision during rename | Reject rename, keep current username |
| Last email deletion | Reject deletion, user must keep at least one email |
| 2FA verification after password change | Sessions invalidated, must re-authenticate with 2FA |
| Expired activation/reset codes | Reject and offer to resend |
| LDAP sync during source outage | Log error, continue with cached user data |
| WebAuthn sign count regression | Flag credential as potentially cloned, warn user |
| Token lookup cache miss | Fall back to database lookup |
| Blocked user creates PR via fork | Reject PR creation to blocker's repos |
| OAuth2 app redirect URI mismatch | Reject authorization request |
| Multiple auth sources for same user | Link external identities to single local account |

## Success Criteria

- User registration completes within 5 seconds for standard email/password flow
- Password hashing uses pbkdf2 by default with configurable parameters
- Session lookup from cache completes in under 1ms
- OAuth2 token issuance completes in under 500ms
- TOTP verification accepts codes within the configured time skew window
- WebAuthn registration supports all FIDO2-compatible authenticators
- Permission computation for any resource completes in under 50ms
- Federation signature verification completes in under 100ms
- Token scope validation prevents all unauthorized API access

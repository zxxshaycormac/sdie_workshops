# 01 — Identity & Authentication

Baseline specification of Gitea's individual identity, authentication mechanisms, and account-lifecycle subsystems. Covers user registration, profile and email management, password and session authentication, access tokens, 2FA, WebAuthn, external authentication sources, OpenID, OAuth2 client and provider, activity heatmap, user following, CAPTCHA, avatars, and badges. All requirements describe the current (v1.22.x) system behavior.

**Boundary — what's OUT:** RBAC permission system → Domain 02; organizations/teams → Domain 02; GPG keys → Domain 02; user/org blocking → Domain 02; restricted users → Domain 02; federation protocol → Domain 07.

---

## 1. User Registration

**User Story:** As a new user, I want to create an account so that I can access Gitea's features.

### Ubiquitous Requirements (Registration Constraints)

- **IDN-01-001:** `The system shall require a unique username for each user account.`
- **IDN-01-002:** `The system shall normalize usernames to lowercase for comparison.`
- **IDN-01-003:** `The system shall reject usernames that match reserved names including ".", "..", ".well-known", "admin", "api", "assets", "attachments", "avatar", "avatars", "captcha", "commits", "debug", "error", "explore", "favicon.ico", "ghost", "issues", "login", "manifest.json", "metrics", "milestones", "new", "notifications", "org", "pulls", "raw", "repo", "repo-avatars", "robots.txt", "search", "serviceworker.js", "ssh_info", "swagger.v1.json", "user", "v2", "gitea-actions", and reserved patterns "*.keys", "*.gpg", "*.rss", "*.atom", "*.png".`
- **IDN-01-004:** `The system shall require a valid email address for account registration.`
- **IDN-01-005:** `The system shall hash passwords using a configurable algorithm (default pbkdf2).`
- **IDN-01-006:** `The system shall enforce a configurable minimum password length (default 8 characters).`
- **IDN-01-007:** `The system shall assign user type "individual" to self-registered accounts.`

### Event-Driven Requirements (Registration Workflow)

- **IDN-01-101:** `When a user submits a valid registration form, the system shall create a new user account.`
- **IDN-01-102:** `When registration requires email confirmation, the system shall send an activation email with a time-limited code.`
- **IDN-01-103:** `When a user clicks a valid activation link, the system shall activate the user account.`
- **IDN-01-104:** `When an admin creates a user via the admin panel, the system shall create the account in active state.`
- **IDN-01-105:** `When the setting REGISTER_EMAIL_CONFIRM is enabled, the system shall require email verification before allowing login.`
- **IDN-01-106:** `When the setting DISABLE_REGISTRATION is true, the system shall reject all self-registration attempts.`

### Optional Feature Requirements (External Registration)

- **IDN-01-201:** `Where OAuth2 auto-registration is enabled, the system shall create a local account when a new user authenticates via an external OAuth2 provider.`
- **IDN-01-202:** `Where an external authentication source (LDAP, SMTP, etc.) is configured and active, the system shall authenticate users against that source.`
- **IDN-01-203:** `Where password check against Pwned Passwords is enabled, the system shall warn users whose passwords appear in known breach databases.`

### Unwanted Behaviour Requirements (Registration Errors)

- **IDN-01-301:** `If a user submits a username that already exists, then the system shall reject the registration with an error message.`
- **IDN-01-302:** `If a user submits an email address already registered to another account, then the system shall reject the registration.`
- **IDN-01-303:** `If a user submits a password that does not meet complexity requirements, then the system shall reject the password and display the requirements.`
- **IDN-01-304:** `If an activation link has expired, then the system shall reject the activation and offer to resend the email.`

---

## 2. User Profile Management

**User Story:** As a user, I want to manage my profile so that others can identify and contact me.

### Ubiquitous Requirements (Profile Properties)

- **IDN-02-001:** `The system shall store the following profile properties per user: username, full name, email, location, website, and description.`
- **IDN-02-002:** `The system shall support three user visibility levels: public, limited (visible only to authenticated/connected users), and private (visible only to self and admins).`
- **IDN-02-003:** `The system shall store the user's language and theme preferences.`
- **IDN-02-004:** `The system shall support four notification preference levels: enabled, on-mention only, disabled, and all including own.`

### Event-Driven Requirements (Profile Updates)

- **IDN-02-101:** `When a user updates their profile via settings, the system shall persist the changes immediately.`
- **IDN-02-102:** `When a user changes their visibility to private, the system shall hide the user from search results and explore pages for non-connections.`
- **IDN-02-103:** `When a user enables "keep email private", the system shall display a placeholder email address in the user's public profile.`
- **IDN-02-104:** `When a user enables "keep activity private", the system shall hide the user's activity feed from non-connections.`

### State-Driven Requirements (Login Prohibition)

- **IDN-02-702:** `While a user is marked as prohibited from logging in, the system shall reject all authentication attempts for that user.`

> **Note:** Restricted-user state requirements are documented in Domain 02 (Access Control & Organization).

### Unwanted Behaviour Requirements (Profile Validation)

- **IDN-02-301:** `If a user enters an invalid website URL, then the system shall reject the input with a validation error.`
- **IDN-02-302:** `If a non-admin user attempts to change another user's profile, then the system shall deny the operation.`

---

## 3. Email Management

**User Story:** As a user, I want to manage multiple email addresses so I can control which email is primary and keep others verified.

### Ubiquitous Requirements (Email Constraints)

- **IDN-03-001:** `The system shall allow each user account to have multiple email addresses.`
- **IDN-03-002:** `The system shall designate exactly one email address as primary per account.`
- **IDN-03-003:** `The system shall track activation state independently for each email address.`

### Event-Driven Requirements (Email Workflow)

- **IDN-03-101:** `When a user adds a new email address, the system shall send a verification email with a time-limited activation code.`
- **IDN-03-102:** `When a user clicks a valid email verification link, the system shall mark the email as activated.`
- **IDN-03-103:** `When a user sets a verified email as primary, the system shall update the account's primary email.`
- **IDN-03-104:** `When a user deletes an email address, the system shall remove it from the account, unless it is the last remaining email.`

### Optional Feature Requirements (Email Domain Control)

- **IDN-03-201:** `Where email domain allow lists are configured, the system shall reject registration with email addresses from non-allowed domains.`
- **IDN-03-202:** `Where email domain block lists are configured, the system shall reject registration with email addresses from blocked domains.`

### Unwanted Behaviour Requirements (Email Errors)

- **IDN-03-301:** `If a user adds an email address already registered to another account, then the system shall reject the addition.`
- **IDN-03-302:** `If a user attempts to delete their only email address, then the system shall reject the deletion.`
- **IDN-03-303:** `If an email verification code has expired, then the system shall reject the verification and offer to resend.`

---

## 4. Password Management

**User Story:** As a user, I want to reset my password if I forget it, and as an admin, I want to enforce password policies.

### Ubiquitous Requirements (Password Constraints)

- **IDN-04-001:** `The system shall hash all stored passwords using a configurable algorithm (argon2, bcrypt, scrypt, pbkdf2, or sha256).`
- **IDN-04-002:** `The system shall enforce configurable password complexity rules (minimum length, uppercase, lowercase, numbers, special characters).`

### Event-Driven Requirements (Password Workflow)

- **IDN-04-101:** `When a user requests a password reset, the system shall send a time-limited reset code to the user's primary email.`
- **IDN-04-102:** `When a user submits a valid reset code with a new password, the system shall update the user's password.`
- **IDN-04-103:** `When an admin sets the must_change_password flag on a user, the system shall require the user to change their password at next login.`
- **IDN-04-104:** `When a user changes their password, the system shall invalidate all active sessions for that user.`

### Unwanted Behaviour Requirements (Password Errors)

- **IDN-04-301:** `If a password reset code has expired, then the system shall reject the reset and offer to send a new code.`
- **IDN-04-302:** `If a user submits a password that does not meet complexity requirements, then the system shall reject the password with specific requirement feedback.`
- **IDN-04-303:** `If a user's current password is incorrect during a password change, then the system shall reject the change.`

---

## 5. User Rename & Deletion

**User Story:** As a user, I want to rename my account, and as an admin, I want to delete accounts that are no longer needed.

### Event-Driven Requirements (Rename Workflow)

- **IDN-06-101:** `When a user renames their account, the system shall create a persistent redirect from the old username to the new username.`
- **IDN-06-102:** `When a user renames their account, the system shall update all repository paths to reflect the new username.`
- **IDN-06-103:** `When a request targets a redirected username, the system shall redirect to the current username.`

### Event-Driven Requirements (Deletion Workflow)

- **IDN-06-201:** `When an admin deletes a user account, the system shall transfer the user's repository ownership to a ghost user for attribution preservation.`
- **IDN-06-202:** `When an admin deletes a user account, the system shall remove the user from all team memberships.`
- **IDN-06-203:** `When an admin deletes a user account, the system shall remove all SSH keys, GPG keys, and access tokens associated with the account.`

### Unwanted Behaviour Requirements (Deletion Safeguards)

- **IDN-06-301:** `If an admin attempts to delete the last admin user, then the system shall deny the deletion.`
- **IDN-06-302:** `If a non-admin user attempts to delete a user account, then the system shall deny the operation.`

---

## 6. Password Authentication & Sessions

**User Story:** As a user, I want to log in with my credentials and maintain a session so I don't have to re-authenticate on every request.

### Ubiquitous Requirements (Session Properties)

- **IDN-07-001:** `The system shall store sessions in the database with configurable lifetime.`
- **IDN-07-002:** `The system shall support session regeneration on privilege level changes.`

### Event-Driven Requirements (Login Workflow)

- **IDN-07-101:** `When a user submits valid credentials, the system shall create an authenticated session.`
- **IDN-07-102:** `When a user selects "remember me" during login, the system shall extend the session duration to the configured LOGIN_REMEMBER_DAYS value.`
- **IDN-07-103:** `When a user logs out, the system shall destroy the current session.`
- **IDN-07-104:** `When a user logs out from all devices, the system shall destroy all sessions associated with the user.`

### State-Driven Requirements (Session Lifecycle)

- **IDN-07-701:** `While a session is valid, the system shall allow the user to perform actions within their permission scope.`
- **IDN-07-702:** `While a session has expired, the system shall redirect the user to the login page on the next request.`

### Unwanted Behaviour Requirements (Login Errors)

- **IDN-07-301:** `If a user submits invalid credentials, then the system shall reject the login attempt.`
- **IDN-07-302:** `If a user account is deactivated, then the system shall reject all login attempts with an appropriate message.`
- **IDN-07-303:** `If a user with 2FA enabled does not provide a valid TOTP code, then the system shall reject the login.`

---

## 7. Access Tokens

**User Story:** As a user, I want to create personal access tokens with specific scopes so that I can authenticate API requests without sharing my password.

### Ubiquitous Requirements (Token Properties)

- **IDN-08-001:** `The system shall hash stored tokens using PBKDF2-HMAC-SHA256 with salt.`
- **IDN-08-002:** `The system shall store the last 8 characters of each token in plaintext for fast lookup.`
- **IDN-08-003:** `The system shall cache successfully validated tokens up to a configurable cache size.`

### Event-Driven Requirements (Token Workflow)

- **IDN-08-101:** `When a user creates a personal access token, the system shall generate a secure random token and display it once.`
- **IDN-08-102:** `When a token is created with specific scopes, the system shall restrict the token's API access to those scopes only.`
- **IDN-08-103:** `When a user deletes a token, the system shall invalidate the token immediately and remove it from the cache.`

### Unwanted Behaviour Requirements (Token Errors)

- **IDN-08-301:** `If an API request presents an invalid token, then the system shall reject the request with 401 Unauthorized.`
- **IDN-08-302:** `If an API request presents a valid token lacking the required scope, then the system shall reject the request with 403 Forbidden.`
- **IDN-08-303:** `If a token is used after deletion, then the system shall reject the request.`

---

## 8. Two-Factor Authentication (2FA)

**User Story:** As a user, I want to enable 2FA so that my account is protected even if my password is compromised.

### Ubiquitous Requirements (2FA Properties)

- **IDN-09-001:** `The system shall implement TOTP (RFC 6238) for two-factor authentication.`
- **IDN-09-002:** `The system shall encrypt TOTP secrets using AES with a key derived from the instance SecretKey via MD5.` NOTE: MD5-based key derivation is a known weakness (collision-prone, no salt/stretching); see `models/auth/twofactor.go:95-98` (`getEncryptionKey` = `md5.Sum(SecretKey)`).
- **IDN-09-003:** `The system shall generate emergency scratch codes when 2FA is enabled.`

### Event-Driven Requirements (2FA Workflow)

- **IDN-09-101:** `When a user enables 2FA, the system shall generate a TOTP secret and display a QR code for authenticator app setup.`
- **IDN-09-102:** `When a user with 2FA enabled logs in, the system shall require a valid TOTP code after password verification.`
- **IDN-09-103:** `When a user enters a valid scratch code, the system shall accept it as a valid second factor and consume that code.`
- **IDN-09-104:** `When an admin resets a user's 2FA, the system shall disable 2FA (scratch codes consumed; new ones generated only on re-enrollment).`

### Unwanted Behaviour Requirements (2FA Errors)

- **IDN-09-301:** `If a user submits an invalid TOTP code, then the system shall reject the login attempt.`
- **IDN-09-302:** `If a user submits an already-consumed scratch code, then the system shall reject the code.`
- **IDN-09-303:** `If a user attempts to enable 2FA without verifying a valid TOTP code, then the system shall reject the setup.`

---

## 9. WebAuthn / Passkey

**User Story:** As a user, I want to use a hardware security key or passkey for passwordless authentication.

### Ubiquitous Requirements (WebAuthn Properties)

- **IDN-10-001:** `The system shall implement WebAuthn Level 1 for passwordless authentication.`
- **IDN-10-002:** `The system shall allow multiple WebAuthn credentials per user.`

### Event-Driven Requirements (WebAuthn Workflow)

- **IDN-10-101:** `When a user registers a WebAuthn credential, the system shall store the credential with sign count tracking.`
- **IDN-10-102:** `When a user authenticates with a WebAuthn credential, the system shall verify the sign count to detect credential cloning.`
- **IDN-10-103:** `When a user deletes a WebAuthn credential, the system shall remove it and prevent future authentication with that credential.`

### Unwanted Behaviour Requirements (WebAuthn Errors)

- **IDN-10-301:** `If a WebAuthn authentication attempt has a lower sign count than the stored value, then the system shall flag the credential as potentially cloned.`
- **IDN-10-302:** `If a user has no remaining WebAuthn credentials, then the system shall fall back to password-based authentication.`

---

## 10. External Authentication Sources

**User Story:** As an admin, I want to configure external authentication sources (LDAP, SMTP, etc.) so that users can log in with their existing corporate credentials.

### Ubiquitous Requirements (Auth Source Properties)

- **IDN-11-001:** `The system shall support the following authentication source types: Plain (database), LDAP (BindDN), LDAP (simple auth / direct bind), SMTP, PAM, OAuth2, and SPNEGO/SSPI.`
- **IDN-11-002:** `The system shall store authentication source configurations with activation state, display order, and per-source TLS configuration.`

### Event-Driven Requirements (Auth Source Workflow)

- **IDN-11-101:** `When a user authenticates via an external source for the first time, the system shall create a local user account linked to the external identity.`
- **IDN-11-102:** `When an LDAP source has synchronization enabled, the system shall periodically sync user data from the LDAP directory.`
- **IDN-11-103:** `When an admin creates a new authentication source, the system shall make it available for login immediately upon activation.`

### Optional Feature Requirements (Auth Source Features)

- **IDN-11-201:** `Where an LDAP source is configured with group-to-team mapping, the system shall assign users to teams based on their LDAP group membership.`
- **IDN-11-202:** `Where an LDAP source is configured with admin group filtering, the system shall grant admin privileges to users in the specified groups.`
- **IDN-11-203:** `Where an authentication source is configured to skip local 2FA, the system shall bypass 2FA for users authenticating through that source.`
- **IDN-11-204:** `Where an LDAP source is configured with restricted group filtering, the system shall mark users in those groups as restricted.`

### Unwanted Behaviour Requirements (Auth Source Errors)

- **IDN-11-301:** `If an external authentication source is unreachable, then the system shall reject login attempts against that source with an error.`
- **IDN-11-302:** `If an LDAP sync operation fails, then the system shall log the error and continue with the existing user data.`

---

## 11. OpenID Authentication

**User Story:** As a user, I want to sign in to Gitea using my OpenID identity so that I can authenticate via a decentralized identity provider.

### Ubiquitous Requirements (OpenID Properties)

- **IDN-12-001:** `The system shall implement OpenID 2.0 consumer functionality for delegated authentication.`
- **IDN-12-002:** `The system shall implement OpenID as a dedicated consumer flow via routers/web/auth/openid.go, storing identities in the user_open_id table, separate from the login_source auth.Type enum (NoType, Plain, LDAP, SMTP, PAM, DLDAP, OAuth2, SSPI — no OpenID entry).`
- **IDN-12-003:** `The system shall discover the OpenID provider endpoint from the user-supplied OpenID URL via Yadis / HTML discovery.`
- **IDN-12-004:** `The system shall request the SReg and AX attributes required for account creation (nickname, email, full name).`

### Event-Driven Requirements (OpenID Workflow)

- **IDN-12-101:** `When a user submits an OpenID URL on the login page, the system shall perform discovery and redirect the user to the provider's authentication endpoint.`
- **IDN-12-102:** `When the provider redirects back with a positive assertion, the system shall verify the signature and nonce.`
- **IDN-12-103:** `When the verified OpenID matches an existing local account, the system shall authenticate the user as that account.`
- **IDN-12-104:** `When the verified OpenID does not match an existing account and OpenID signup is enabled, the system shall create a new local account from the returned attributes.`
- **IDN-12-105:** `When a user links an OpenID to their existing account via settings, the system shall store the OpenID URL for future sign-in.`

### Optional Feature Requirements (OpenID Configuration)

- **IDN-12-201:** `Where ENABLE_OPENID_SIGNIN is enabled, the system shall display the OpenID login option on the login page.`
- **IDN-12-202:** `Where ENABLE_OPENID_SIGNUP is enabled, the system shall allow new account creation via OpenID.`
- **IDN-12-203:** `Where WHITELISTED_URIS is configured, the system shall reject OpenID URLs not matching the allowlist.`
- **IDN-12-204:** `Where BLACKLISTED_URIS is configured, the system shall reject OpenID URLs matching the blocklist.`

### Unwanted Behaviour Requirements (OpenID Errors)

- **IDN-12-301:** `If the OpenID URL fails discovery, then the system shall reject the login with a discovery error message.`
- **IDN-12-302:** `If the provider returns a negative assertion, then the system shall return the user to the login page without authenticating.`
- **IDN-12-303:** `If the assertion signature verification fails, then the system shall reject the login and log the failure.`
- **IDN-12-304:** `If the assertion nonce has been replayed or expired, then the system shall reject the login.`

---

## 12. Activity Heatmap

**User Story:** As a user, I want a visual contribution calendar on my profile so that my activity over the past year is visible at a glance.

### Ubiquitous Requirements (Heatmap Properties)

- **IDN-13-001:** `The system shall track contribution events per user across repositories, issues, pull requests, and comments.`
- **IDN-13-002:** `The system shall aggregate contribution counts per 15-minute bucket (to permit client-side timezone rendering).`
- **IDN-13-003:** `The system shall expose heatmap data for the trailing 365-day window relative to the request time.`
- **IDN-13-004:** `The system shall render the heatmap on the user profile page when the profile is viewable by the requester.`

### Event-Driven Requirements (Heatmap Workflow)

- **IDN-13-101:** `When a user performs a contribution action (push, issue create, PR create/review, comment), the system shall increment the corresponding day's contribution count for that user.`
- **IDN-13-102:** `When a viewer requests a user's profile, the system shall return the heatmap data as a structured JSON payload keyed by timestamp and count.`
- **IDN-13-103:** `When a contribution action is reversed (issue deleted, PR reverted), the system shall decrement the corresponding day's contribution count.`
- **IDN-13-104:** `When a user performs an action on a private repository, the system shall include the contribution in the owner's heatmap but exclude it from viewers without access.`

### Optional Feature Requirements (Heatmap Configuration)

- **IDN-13-201:** `Where ENABLE_USER_HEATMAP is disabled in configuration, the system shall hide the heatmap from all user profiles.`
- **IDN-13-202:** `Where a user's visibility is set to private, the system shall hide the heatmap from viewers who are not authorized connections.`
- **IDN-13-203:** `Where a user has enabled "keep activity private", the system shall hide the heatmap from non-connections.`

### Unwanted Behaviour Requirements (Heatmap Errors)

- **IDN-13-301:** `If a viewer lacks permission to view a user's profile, then the system shall return an empty heatmap dataset rather than a 403.`
- **IDN-13-302:** `If the requested time range exceeds the maximum supported 365-day window, then the system shall clamp the request to the maximum window.`

---

## 13. OAuth2 Client Login

**User Story:** As a user, I want to sign in to Gitea with my GitHub, GitLab, Google, or other external account so that I do not need a separate Gitea password.

### Ubiquitous Requirements (OAuth2 Client Properties)

- **IDN-14-001:** `The system shall implement the OAuth2 authorization-code flow as a client against configured external providers.`
- **IDN-14-002:** `The system shall support per-provider configuration of client ID, client secret, authorization endpoints, token endpoints, profile endpoints, and scopes.`
- **IDN-14-003:** `The system shall map the external profile fields (username, email, full name, avatar URL) to local account attributes using configurable attribute mapping.`
- **IDN-14-004:** `The system shall treat OAuth2 client login as one of the IDN-11 authentication source types.`
- **IDN-14-005:** `The system shall support the following 16 built-in external OAuth2/OpenID Connect providers: github, gitlab, gplus (Google), gitea, nextcloud, mastodon, azureadv2 (Azure AD v2), bitbucket, dropbox, facebook, twitter, discord, yandex, azuread, microsoftonline, and openidConnect (generic OpenID Connect).`

### Event-Driven Requirements (OAuth2 Client Workflow)

- **IDN-14-101:** `When a user clicks the provider button on the login page, the system shall redirect the user to the provider's authorization endpoint.`
- **IDN-14-102:** `When the provider redirects back with an authorization code, the system shall exchange the code for an access token using the configured client secret.`
- **IDN-14-103:** `When the access token is obtained, the system shall call the provider's profile endpoint to fetch the external user identity.`
- **IDN-14-104:** `When the external identity matches an existing local account, the system shall authenticate the user as that account.`
- **IDN-14-105:** `When the external identity does not match an existing local account and OAuth2 auto-registration is enabled, the system shall create a new local account linked to the external identity.`
- **IDN-14-106:** `When a user is already authenticated and chooses to link an additional external provider, the system shall record the link without creating a new local account.`

### Optional Feature Requirements (OAuth2 Client Configuration)

- **IDN-14-201:** `Where OAuth2 auto-registration is disabled and the external identity has no matching local account, the system shall prompt the user to either sign in to an existing local account or stop the flow.`
- **IDN-14-202:** `Where the provider returns an unverified email and email verification is required, the system shall refuse to create a local account and prompt the user to verify the email with the provider first.`
- **IDN-14-203:** `Where two-factor authentication bypass is enabled for an OAuth2 source, the system shall skip local 2FA for users authenticating through that source.`

### Unwanted Behaviour Requirements (OAuth2 Client Errors)

- **IDN-14-301:** `If the provider returns an invalid or expired authorization code, then the system shall reject the login with an error message and return the user to the login page.`
- **IDN-14-302:** `If the provider's profile endpoint is unreachable, then the system shall log the error and report a temporary authentication failure.`
- **IDN-14-303:** `If the external identity's username collides with a local username not yet linked to that provider, then the system shall reject the auto-registration and prompt for a different username.`
- **IDN-14-304:** `If the provider's email domain is on the configured blocklist, then the system shall reject the auto-registration.`

---

## 14. CAPTCHA

**User Story:** As an administrator, I want to protect registration and public forms from automated spam using CAPTCHA challenges.

### Ubiquitous Requirements (CAPTCHA Properties)

- **IDN-15-001:** `The system shall require CAPTCHA verification on the user registration form when any CAPTCHA service is enabled.`
- **IDN-15-002:** `The system shall validate CAPTCHA responses server-side against the configured CAPTCHA provider.`
- **IDN-15-003:** `The system shall support exactly one active CAPTCHA service at a time based on configuration.`

### Event-Driven Requirements (CAPTCHA Workflow)

- **IDN-15-101:** `When a user submits the registration form, the system shall validate the CAPTCHA response before processing the registration.`
- **IDN-15-102:** `When the CAPTCHA response is invalid, the system shall reject the form submission with an error message.`
- **IDN-15-103:** `When a CAPTCHA service is disabled, the system shall remove the CAPTCHA widget from the registration form.`

### Optional Feature Requirements (CAPTCHA Providers)

- **IDN-15-201:** `Where hCaptcha is enabled, the system shall render the hCaptcha widget with the configured site key and theme.`
- **IDN-15-202:** `Where reCAPTCHA v2 is enabled, the system shall render the reCAPTCHA checkbox or invisible widget with the configured site key.`
- **IDN-15-203:** `Where Cloudflare Turnstile is enabled, the system shall render the Turnstile JavaScript widget with the configured site key.`
- **IDN-15-204:** `Where mCaptcha is enabled, the system shall render the mCaptcha proof-of-work widget with the configured site key and instance URL.`

### Unwanted Behaviour Requirements (CAPTCHA Errors)

- **IDN-15-301:** `If the CAPTCHA provider is unreachable during validation, then the system shall reject the form submission and display an error.`
- **IDN-15-302:** `If the configured site key or secret is invalid, then the system shall log an error and fail CAPTCHA validation.`
- **IDN-15-303:** `If a user bypasses the CAPTCHA widget (missing response), then the system shall reject the form submission.`

---

## 15. Avatar Systems

**User Story:** As a user, I want my avatar displayed across the interface from my preferred source (Gravatar, Libravatar, uploaded image, or generated) so that others can visually identify me.

### Ubiquitous Requirements (Avatar Properties)

- **IDN-16-001:** `The system shall display user avatars in the interface for user mentions, comments, profiles, and commit history.`
- **IDN-16-002:** `The system shall generate a default avatar using either an identicon pattern or the user's initial letter when no external avatar source is available.`
- **IDN-16-003:** `The system shall enforce a configurable maximum file size for uploaded avatar images.`

### Event-Driven Requirements (Avatar Workflow)

- **IDN-16-101:** `When a user uploads a custom avatar image, the system shall store the image locally and display it as the user's avatar.`
- **IDN-16-102:** `When a user removes a custom uploaded avatar, the system shall fall back to the configured external avatar source or generated avatar.`
- **IDN-16-103:** `When the system resolves a user avatar, the system shall check in order: local upload, federated avatar source, Gravatar source, then generated fallback.`

### Optional Feature Requirements (External Avatar Sources)

- **IDN-16-201:** `Where Gravatar is enabled, the system shall look up the user's avatar by email hash from the configured Gravatar URL over HTTP or HTTPS.`
- **IDN-16-202:** `Where Libravatar (federated avatar) is enabled, the system shall look up the user's avatar from the Libravatar federation using the user's email address.`
- **IDN-16-203:** `Where a custom Gravatar source URL is configured, the system shall use that URL instead of the default gravatar.com domain.`

### Unwanted Behaviour Requirements (Avatar Errors)

- **IDN-16-301:** `If the external avatar source (Gravatar or Libravatar) is unreachable, then the system shall fall back to the generated avatar.`
- **IDN-16-302:** `If an uploaded avatar file exceeds the maximum size limit, then the system shall reject the upload with a validation error.`
- **IDN-16-303:** `If an uploaded avatar file has an unsupported image format, then the system shall reject the upload.`

---

## 16. Badges

**User Story:** As an administrator, I want to assign badges to user accounts so that users can be recognized for achievements or roles.

**Note:** Badges are a **user-account** feature. The `UserBadge` model (`models/user/badge.go`) is keyed by `UserID` and is administered via `/admin/users/{username}/badges` (`routers/api/v1/admin/user_badge.go`). No `OrgBadge` model exists.

### Ubiquitous Requirements (Badge Properties)

- **IDN-17-001:** `The system shall store badge definitions (Badge) with a unique slug identifier, description, and image URL.`
- **IDN-17-002:** `The system shall associate badges with user accounts via the UserBadge join table keyed on UserID (no org-specific badge model or org-profile rendering exists).`

---

## 17. User Following

**User Story:** As a user, I want to follow other users so that I can see their public activity, and I want to manage who I follow.

### Ubiquitous Requirements (Following Properties)

- **IDN-18-001:** `The system shall store follow relationships in a Follow table with a unique constraint on the (UserID, FollowID) pair.`
- **IDN-18-002:** `The system shall maintain denormalized follower count (num_followers) and following count (num_following) on each user record.`

### Event-Driven Requirements (Following Workflow)

- **IDN-18-101:** `When a user follows another user, the system shall create a Follow record and increment both the follower's num_following and the followee's num_followers within a single transaction.`
- **IDN-18-102:** `When a user unfollows another user, the system shall delete the Follow record and decrement both counters within a single transaction.`
- **IDN-18-103:** `When a viewer requests a user's followers list, the system shall return individual users (UserTypeIndividual) who follow that user, filtered by visibility to the viewer.`
- **IDN-18-104:** `When a viewer requests a user's following list, the system shall return users and organizations (UserTypeIndividual and UserTypeOrganization) that the user follows, filtered by visibility to the viewer.`

### State-Driven Requirements (Following Lifecycle)

- **IDN-18-701:** `While a follow relationship exists, the system shall report IsFollowing as true for that (userID, followID) pair.`

### Unwanted Behaviour Requirements (Following Errors)

- **IDN-18-301:** `If a user attempts to follow themselves, then the system shall reject the operation as a no-op (self-follow is not permitted).`
- **IDN-18-302:** `If a user attempts to follow another user who has blocked them (or whom they have blocked), then the system shall reject the follow with a blocked-user error.`
- **IDN-18-303:** `If a user attempts to follow a user they already follow, then the system shall treat the operation as a no-op.`

---

## 18. OAuth2 Provider / Authorization Server

**User Story:** As an administrator or third-party application developer, I want Gitea to act as an OAuth2 authorization server so that external applications can obtain delegated access tokens to the Gitea API on behalf of users.

### Ubiquitous Requirements (OAuth2 Provider Properties)

- **IDN-19-001:** `The system shall expose OAuth2 authorization-server endpoints under the /login/oauth path prefix when OAuth2 is enabled (default enabled).`
- **IDN-19-002:** `The system shall issue signed JWT-based access tokens and refresh tokens using the configured JWT signing algorithm (default RS256) and private key file.`
- **IDN-19-003:** `The system shall allow users to register OAuth2 applications (clients) with redirect URIs, client secrets, and scoped access.`
- **IDN-19-004:** `The system shall expose the following endpoints: GET/POST /login/oauth/authorize, POST /login/oauth/grant, POST /login/oauth/access_token, GET /login/oauth/userinfo, GET /login/oauth/keys, and POST /login/oauth/introspect.`
- **IDN-19-005:** `The system shall seed a set of default OAuth2 applications (git-credential-oauth, git-credential-manager, tea) unless overridden by configuration.`

### Event-Driven Requirements (OAuth2 Provider Workflow)

- **IDN-19-101:** `When an unauthenticated user hits /login/oauth/authorize, the system shall require sign-in before proceeding with the authorization grant.`
- **IDN-19-102:** `When a signed-in user approves an authorization request, the system shall redirect to the client's redirect URI with an authorization code.`
- **IDN-19-103:** `When a client posts a valid authorization code and client credentials to /login/oauth/access_token, the system shall return access and refresh tokens.`
- **IDN-19-104:** `When a client presents a valid access token to /login/oauth/userinfo, the system shall return the resource-owner claims.`
- **IDN-19-105:** `When a client posts a token to /login/oauth/introspect, the system shall return the token's active state and metadata.`

### Optional Feature Requirements (OAuth2 Provider Configuration)

- **IDN-19-201:** `Where OAuth2 is explicitly disabled via ENABLED=false, the system shall refuse all /login/oauth requests.`
- **IDN-19-202:** `Where the access-token expiration time is configured, the system shall honor it (default 3600 seconds).`
- **IDN-19-203:** `Where the refresh-token expiration time is configured, the system shall honor it (default 730 hours).`
- **IDN-19-204:** `Where INVALIDATE_REFRESH_TOKENS is enabled, the system shall invalidate a refresh token after it has been used to mint a new access token.`

### Unwanted Behaviour Requirements (OAuth2 Provider Errors)

- **IDN-19-301:** `If a client submits an invalid authorization code, the system shall reject the token request with an OAuth2 error response.`
- **IDN-19-302:** `If a client submits credentials for a non-existent or disabled application, the system shall reject the request.`
- **IDN-19-303:** `If the configured JWT signing private key file is missing or unreadable, the system shall fail to start with a fatal error.`

---

## Configuration Reference

The following INI sections configure identity and authentication behaviors. Add to or modify these via `custom/conf/app.ini` or environment overrides using the `GITEA__section__key` format.

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

Controls Gitea acting as an OAuth2 authorization server (see IDN-19). Also generates the JWT secret reused as the general token-signing secret.
- **ENABLED**: Enable the OAuth2 provider (default true). `ENABLE` is a deprecated alias renamed to `ENABLED`.
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

Controls Gitea acting as an OAuth2/OpenID Connect *client* (see IDN-14).
- **ENABLE_AUTO_REGISTRATION**: Automatically create a local account on first OAuth2 login when enabled (default false).
- **USERNAME**: Source field for generating the local username from OAuth2 data (`userid`, `nickname`, `email`, `preferred_username`; default `nickname`).
- **ACCOUNT_LINKING**: Behavior when an OAuth2 identity matches an existing local account (`disabled`, `login`, `auto`; default `login`).
- **UPDATE_AVATAR**: Update the local avatar from the OAuth2 provider on each login when enabled (default false).
- **OPENID_CONNECT_SCOPES**: Space-separated OpenID Connect scopes requested from the provider.
- **REGISTER_EMAIL_CONFIRM**: Require email confirmation for OAuth2 auto-registered accounts (inherits `[service].REGISTER_EMAIL_CONFIRM`).

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
- **BR-01-009:** Ghost users are system accounts used for attribution after account deletion
- **BR-01-010:** Registration supports a third confirmation mode beside email-confirm: `REGISTER_MANUAL_CONFIRM` requires an admin to manually activate newly registered accounts; it is ignored when `REGISTER_EMAIL_CONFIRM` is enabled (evidence: `modules/setting/service.go:147`)
- **BR-01-011:** Only one CAPTCHA service may be active at a time; enabling one disables the others
- **BR-01-012:** Avatar resolution follows a strict priority: local upload > federated avatar > Gravatar > generated fallback
- **BR-01-013:** Follow relationships are directional (follower → followee) and self-follow is not permitted

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
| OAuth2 app redirect URI mismatch | Reject authorization request |
| Multiple auth sources for same user | Link external identities to single local account |
| CAPTCHA provider timeout during registration | Reject submission, display error to user |
| Gravatar service unreachable for avatar | Display generated identicon or initial avatar |
| User attempts to follow themselves | Reject as no-op; no Follow record created |
| User attempts to follow a blocked user | Reject with blocked-user error |

## Success Criteria

- User registration completes within 5 seconds for standard email/password flow
- Password hashing uses pbkdf2 by default with configurable parameters
- Session lookup from cache completes in under 1ms
- OAuth2 token issuance completes in under 500ms
- TOTP verification accepts codes within the configured time skew window
- WebAuthn registration supports all FIDO2-compatible authenticators
- Token scope validation prevents all unauthorized API access
- CAPTCHA validation completes within 3 seconds including network round-trip
- Avatar resolution falls back gracefully through the entire priority chain

# 09 — Integration & Extension

## 1. OAuth2 Provider

### What
Gitea as an OAuth2 server. Applications can register, obtain tokens, and access user resources.

### Grant Types
- Authorization Code (with PKCE)
- Implicit flow
- Client Credentials

### Scopes
- read:repository, write:repository
- read:organization, write:organization
- read:user, write:user
- read:admin:org, write:admin:org
- read:admin:user, write:admin:user
- sudo (admin privilege escalation)

### UI Routes
- /user/settings/applications/oauth2 — Manage OAuth2 apps
- /user/settings/applications/oauth2/:id/edit — Edit app
- /login/oauth/authorize — Authorization endpoint
- /login/oauth/access_token — Token endpoint

### API Endpoints
- POST /user/applications/oauth2 — Create app
- GET /user/applications/oauth2 — List apps
- GET /user/applications/oauth2/:id — Get app
- PATCH /user/applications/oauth2/:id — Update app
- DELETE /user/applications/oauth2/:id — Delete app

### Config
- [oauth2] ENABLED — enable OAuth2 provider
- [oauth2] ACCESS_TOKEN_EXPIRATION_TIME — token TTL (default 3600s)
- [oauth2] REFRESH_TOKEN_EXPIRATION_TIME — refresh TTL (default 730h)
- [oauth2] JWT_SIGNING_ALGORITHM — RS256 default
- [oauth2] DEFAULT_APPLICATIONS — built-in apps (git-credential-oauth, tea, etc.)

---

## 2. External Auth Sources

### What
Authentication against external services: LDAP, SMTP, PAM, OAuth2, SAML, SPNEGO/SSPI.

### Supported Types
1. LDAP (BindDN): TLS, user search, attribute mapping, group membership
2. LDAP (simple auth): direct bind
3. SMTP: email-based auth with TLS
4. PAM: Unix system auth
5. OAuth2 (consumer): external OAuth2 providers
6. SAML: SAML 2.0 identity federation
7. SPNEGO/SSPI: Windows integrated auth, Kerberos
8. FreeIPA: specialized LDAP variant

### Features
- User sync from LDAP
- Group-to-team mapping
- Admin/restricted group filtering
- Source-level enable/disable
- TLS configuration

### UI Routes
- GET /admin/auths — List sources
- GET /admin/auths/new — Create source
- GET /admin/auths/:authid — Edit source
- POST /admin/auths — Create source
- POST /admin/auths/:authid — Update source
- POST /admin/auths/:authid/delete — Delete source

### Config
- [auth] DISABLE_REGISTRATION — disable self-registration
- [auth] ENABLE_TOKEN_AUTHENTICATION — enable token auth
- [auth] DISABLE_HTTP_TOKEN_AUTHENTICATION — disable HTTP token

---

## 3. REST API

### What
Comprehensive REST API with OpenAPI/Swagger documentation for all Gitea features.

### Authentication
- Basic Auth
- Token (query param or Bearer header)
- OAuth2 tokens

### Features
- Versioned at /api/v1
- Pagination (page/limit)
- Rate limiting
- Swagger UI at /api/swagger/
- ETag/conditional requests

### Config
- [api] ENABLED — enable API
- [api] ENABLE_SWAGGER — enable Swagger UI
- [api] MAX_RESPONSE_ITEMS — max items per page (default 50)
- [api] DEFAULT_PAGING_NUM — default page size
- [api] RATE_LIMIT — rate limit value

---

## 4. Markup Renderers

### What
Multiple document format rendering with code highlighting.

### Supported Formats
- Markdown (GFM): task lists, tables, math (KaTeX), custom link schemes
- reStructuredText: RST parsing
- Jupyter Notebooks: .ipynb rendering
- AsciiDoc: attributes and macros
- Org mode: Emacs org format
- CSV: tabular display
- PDF: inline viewing

### Code Highlighting
- Chroma backend (100+ languages)
- Line highlighting
- Copy button
- Theme support

### Custom Renderers
- External renderer support via config
- Isolated execution

### Config
- [markup] ENABLE_HARD_LINE_BREAKS — Markdown hard line breaks
- [markup] CUSTOM_URL_SCHEMES — custom link schemes
- [markup.*] — per-renderer config
- [markdown] CUSTOM_URL_SCHEMES — markdown-specific

---

## 5. Theme System

### What
UI theme management with custom theme support.

### Default Themes
- gitea — default light theme
- gitea-dark — dark theme
- gitea-auto — auto-switching

### Features
- CSS-based theming
- User preference
- Custom themes in custom/ directory

### Config
- [ui] DEFAULT_THEME — default theme
- [ui] THEMES — available themes list

---

## 6. Custom Assets

### What
Override Gitea's appearance through custom templates, static assets, and locales.

### Override Hierarchy
1. custom/templates/ — highest priority
2. templates/{theme}/ — theme-specific
3. templates/ — default

### Features
- Custom Go HTML templates
- Custom static files (CSS, JS, images)
- Custom locale/translation files
- Branding (logo, footer, links)

### Config
- [other] CUSTOM_PATH — custom files root (default: custom/)

---

## 7. SSH Server

### What
Built-in SSH server for Git operations (clone, push, pull).

### Features
- Git protocol support (upload-pack, receive-pack)
- SSH key management (upload, fingerprint, authorization)
- Configurable port
- Authorized keys management

### Config
- [server] SSH_DOMAIN — SSH clone URL domain
- [server] SSH_PORT — SSH port
- [server] SSH_LISTEN_PORT — listen port
- [server] SSH_ROOT_PATH — ~/.ssh path
- [server] SSH_KEYGEN_PATH — ssh-keygen binary
- [server] START_SSH_SERVER — start builtin SSH server
- [server] BUILTIN_SSH_SERVER_DEFAULT_DOMAINS — default domains

---

## 8. Mailer/Email

### What
Email notification system with SMTP support and incoming email processing.

### Outgoing Email
- SMTP with TLS/SSL
- Sendmail fallback
- HTML + plaintext emails
- Notification types: issues, PRs, comments, releases, team invites, repo events
- Subject prefix configurable

### Incoming Email
- Email-to-issue conversion
- Email-to-comment conversion
- Processing queue with authentication

### Config
- [mailer] ENABLED — enable email
- [mailer] HOST — SMTP server
- [mailer] FROM — sender address
- [mailer] USER / PASSWD — SMTP auth
- [mailer] USE_TLS / SKIP_VERIFY — TLS settings
- [mailer] SUBJECT_PREFIX — email subject prefix
- [mailer] SEND_AS_PLAIN_TEXT — plain text fallback

---

## 9. Translation/i18n

### What
Internationalization with 30+ languages.

### Features
- .ini format locale files
- Dynamic language switching
- User language preference
- English fallback
- Pluralization support
- RTL support

### Config
- [ui] LANGS — available languages
- [ui] NAMES — language display names
- [i18n] LANGS — configured languages
- [i18n] NAMES — display names

---

## 10. CAPTCHA

### What
Spam protection for registration and forms.

### Supported Services
1. hCaptcha: site key/secret, theme
2. reCAPTCHA v2: checkbox/invisible
3. Cloudflare Turnstile: JavaScript widget
4. mCaptcha: privacy-focused

### Config
- [hcaptcha] ENABLED / SITEKEY / SECRET
- [recaptcha] ENABLED / SITEKEY / SECRET
- [turnstile] ENABLED / SITEKEY / SECRET
- [mcaptcha] ENABLED / SITEKEY / SECRET / URL

---

## 11. Avatar Systems

### What
Multi-source avatar management with fallbacks.

### Sources
1. Gravatar: email-based (HTTP/HTTPS)
2. Libravatar: open alternative
3. Local upload: custom image
4. Generated: identicon or initial letter

### Config
- [avatar] ENABLE_GRAVATAR — enable Gravatar
- [avatar] ENABLE_FEDERATED_AVATAR — enable Libravatar
- [avatar] GRAVATAR_SOURCE — Gravatar URL
- [avatar] DEFAULT_AVATAR_TYPE — identicon or custom
- [avatar] MAX_FILE_SIZE — upload size limit
- [avatar] RENDERED_SIZE_FACTOR — rendered size

---

## 12. Proxy Support

### What
HTTP proxy for external requests and PROXY protocol for load balancers.

### HTTP Proxy
- Proxy for outgoing requests (avatars, webhooks, etc.)
- Bypass rules
- Authentication

### PROXY Protocol
- v1/v2 support
- Trusted proxy IP ranges
- Client IP preservation

### Config
- [proxy] PROXY_URL — proxy URL
- [proxy] PROXY_HOSTS — bypass list
- [proxyprotocol] TRUSTED_PROXIES — trusted IPs

---

## 13. Update Checker

### What
Automatic version checking and update notifications.

### Features
- Version comparison
- Release notes display
- Channel selection (stable/beta/dev)

### Config
- [updatechecker] ENABLED — enable checker
- [updatechecker] CHANNEL — update channel
- [updatechecker] ENDPOINT_URL — check URL

---

## 14. Git Smart HTTP

### What
Git Smart HTTP protocol for repository operations over HTTP/HTTPS.

### Features
- Smart HTTP (upload-pack/receive-pack)
- Basic auth and token authentication
- LFS support
- Bandwidth optimization

### Config
- [git] HOME_PATH — git home
- [git] DISABLE_DIFF_PARTIAL — disable partial diff
- [git] MAX_GIT_DIFF_LINES — max diff lines
- [git] GC_ARGS — garbage collection args
- [git] FSYNC — fsync on write
- [git] VerbosePush — verbose push output

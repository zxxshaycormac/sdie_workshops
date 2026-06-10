# 09 — Integration & Extension

Baseline specification of Gitea's REST API, markup rendering, theming, customization, SSH, email, internationalization, CAPTCHA, avatars, proxy, update checking, and Git Smart HTTP subsystems. All requirements describe the current (v1.22.x) system behavior.

---

## 1. REST API

**User Story:** As a developer, I want to interact with Gitea programmatically via a REST API so that I can automate repository, user, and organization management.

### Ubiquitous Requirements (API Properties)

- **INT-01-001:** `The system shall expose a versioned REST API at the path prefix /api/v1.`
- **INT-01-002:** `The system shall support the following authentication methods for API requests: HTTP Basic Auth, access token (query parameter or Authorization header), and OAuth2 bearer tokens.`
- **INT-01-003:** `The system shall enforce token scope restrictions on every API request.`
- **INT-01-004:** `The system shall paginate list responses using page and limit query parameters.`
- **INT-01-005:** `The system shall respect a configurable maximum response items limit (default 50).`

### Event-Driven Requirements (API Workflow)

- **INT-01-101:** `When a client sends an API request with valid credentials, the system shall process the request and return the appropriate response.`
- **INT-01-102:** `When a client requests a paginated resource, the system shall return a page of results with HTTP Link headers for navigation.`
- **INT-01-103:** `When the API rate limit is exceeded, the system shall return HTTP 429 with a Retry-After header.`
- **INT-01-104:** `When a client sends a conditional request with an If-None-Match header, the system shall return HTTP 304 if the resource ETag has not changed.`
- **INT-01-105:** `When the API is disabled via configuration, the system shall reject all API requests with HTTP 404.`

### Optional Feature Requirements (Swagger Documentation)

- **INT-01-201:** `Where Swagger is enabled, the system shall serve an interactive Swagger UI at /api/swagger/.`
- **INT-01-202:** `Where Swagger is enabled, the system shall generate an OpenAPI-compliant specification document covering all API endpoints.`
- **INT-01-203:** `Where Swagger is disabled, the system shall return HTTP 404 for all Swagger UI and spec requests.`

### Unwanted Behaviour Requirements (API Errors)

- **INT-01-301:** `If a client sends an API request without valid credentials, then the system shall reject the request with HTTP 401 Unauthorized.`
- **INT-01-302:** `If a client sends an API request with credentials lacking the required scope, then the system shall reject the request with HTTP 403 Forbidden.`
- **INT-01-303:** `If a client sends a request to a non-existent API endpoint, then the system shall return HTTP 404 Not Found.`
- **INT-01-304:** `If a client submits an invalid request body, then the system shall return HTTP 400 with a descriptive error message.`

---

## 2. Markup Renderers

**User Story:** As a user, I want rendered documents in multiple formats (Markdown, reStructuredText, Jupyter, AsciiDoc, Org) with syntax highlighting so that I can read rich content in repositories, issues, and comments.

### Ubiquitous Requirements (Rendering Engine)

- **INT-02-001:** `The system shall render Markdown content using GitHub-Flavored Markdown (GFM) with support for tables, task lists, strikethrough, and autolinks.`
- **INT-02-002:** `The system shall apply syntax highlighting to code blocks using the Chroma engine supporting over 100 programming languages.`
- **INT-02-003:** `The system shall sanitize all rendered HTML to prevent cross-site scripting (XSS) attacks.`
- **INT-02-004:** `The system shall provide a copy button on rendered code blocks for easy clipboard copying.`

### Event-Driven Requirements (Rendering Workflow)

- **INT-02-101:** `When a user views a file with a recognized markup extension, the system shall render the file content as formatted HTML instead of displaying raw source.`
- **INT-02-102:** `When a user views a Jupyter Notebook (.ipynb), the system shall render the notebook cells with code output, tables, and embedded images.`
- **INT-02-103:** `When a user views a CSV file, the system shall render the content as an interactive table.`
- **INT-02-104:** `When a user views a PDF file, the system shall provide an inline PDF viewer.`
- **INT-02-105:** `When Markdown content contains mathematical notation, the system shall render the notation using KaTeX.`

### Optional Feature Requirements (Additional Renderers)

- **INT-02-201:** `Where reStructuredText rendering is enabled, the system shall render .rst files as formatted HTML.`
- **INT-02-202:** `Where AsciiDoc rendering is enabled, the system shall render .adoc files with attribute expansion and macro support.`
- **INT-02-203:** `Where Org mode rendering is enabled, the system shall render .org files as formatted HTML.`
- **INT-02-204:** `Where a custom external renderer is configured via [markup.*] sections, the system shall execute the external command and render the output as HTML.`
- **INT-02-205:** `Where custom URL schemes are configured, the system shall convert matching URLs into clickable links in rendered Markdown.`

### Unwanted Behaviour Requirements (Rendering Errors)

- **INT-02-301:** `If a markup file contains invalid syntax, then the system shall display the raw source with an error indication.`
- **INT-02-302:** `If an external renderer command fails or times out, then the system shall display the raw file content with an error message.`
- **INT-02-303:** `If rendered content exceeds the maximum allowed output size, then the system shall truncate the output and display a warning.`

---

## 3. Theme System

**User Story:** As a user, I want to select and apply visual themes so that the Gitea interface matches my preference for light, dark, or custom appearance.

### Ubiquitous Requirements (Theme Properties)

- **INT-03-001:** `The system shall provide three built-in themes: gitea (light), gitea-dark (dark), and gitea-auto (system preference).`
- **INT-03-002:** `The system shall apply themes through CSS stylesheets without requiring page reload.`
- **INT-03-003:** `The system shall store the user's theme preference in user settings.`
- **INT-03-004:** `The system shall apply the configured default theme to unauthenticated users and users with no preference set.`

### Event-Driven Requirements (Theme Workflow)

- **INT-03-101:** `When a user selects a theme from the settings page, the system shall apply the theme immediately and persist the preference.`
- **INT-03-102:** `When the gitea-auto theme is active and the operating system switches between light and dark mode, the system shall switch the interface accordingly.`

### Optional Feature Requirements (Custom Themes)

- **INT-03-201:** `Where custom theme CSS files are placed in the custom/public/assets/css directory, the system shall detect and list them as available themes.`
- **INT-03-202:** `Where the [ui] THEMES configuration lists additional theme names, the system shall include those themes in the theme selector dropdown.`
- **INT-03-203:** `Where a custom theme name matches a CSS file in the custom assets directory, the system shall load that stylesheet as the selected theme.`

### Unwanted Behaviour Requirements (Theme Errors)

- **INT-03-301:** `If a user's stored theme references a theme that is no longer available, then the system shall fall back to the configured default theme.`
- **INT-03-302:** `If a custom theme CSS file fails to load, then the system shall fall back to the default gitea theme.`

---

## 4. Custom Assets

**User Story:** As an administrator, I want to override templates, add custom static files, and modify branding so that I can tailor the Gitea instance to my organization.

### Ubiquitous Requirements (Override Hierarchy)

- **INT-04-001:** `The system shall resolve template files using a three-level priority: custom/templates/ (highest), templates/{theme}/, templates/ (default).`
- **INT-04-002:** `The system shall serve static files from the custom/public/ directory with higher priority than built-in assets.`
- **INT-04-003:** `The system shall load locale files from the custom/options/locale/ directory to supplement or override default translations.`
- **INT-04-004:** `The system shall embed built-in templates and static assets as bindata for single-binary deployment.`

### Event-Driven Requirements (Asset Workflow)

- **INT-04-101:** `When the system resolves a template path, the system shall check the custom/templates/ directory first before falling back to built-in templates.`
- **INT-04-102:** `When a request targets a static asset, the system shall serve the file from custom/public/ if it exists; otherwise, serve the embedded bindata asset.`
- **INT-04-103:** `When custom locale files are present, the system shall merge custom translations with default translations at startup.`

### Optional Feature Requirements (Branding)

- **INT-04-201:** `Where a custom logo file is placed in the custom assets directory, the system shall display the custom logo in place of the default Gitea logo.`
- **INT-04-202:** `Where the footer template is overridden, the system shall render the custom footer content on all pages.`
- **INT-04-203:** `Where custom links are configured, the system shall display the links in the navigation or footer area.`

### Unwanted Behaviour Requirements (Asset Errors)

- **INT-04-301:** `If a custom template file contains invalid Go template syntax, then the system shall log an error and fall back to the default template.`
- **INT-04-302:** `If a custom static file exceeds the maximum allowed size, then the system shall return HTTP 500 and log an error.`
- **INT-04-303:** `If a custom locale file contains invalid INI syntax, then the system shall skip the file and use default translations with a logged warning.`

---

## 5. SSH Server

**User Story:** As a user, I want to perform Git operations over SSH so that I can securely clone, push, and pull repositories.

### Ubiquitous Requirements (SSH Properties)

- **INT-05-001:** `The system shall support Git protocol operations over SSH including git-upload-pack and git-receive-pack.`
- **INT-05-002:** `The system shall manage SSH public keys per user with fingerprint tracking.`
- **INT-05-003:** `The system shall construct SSH clone URLs using the configured SSH domain and port.`

### Event-Driven Requirements (SSH Workflow)

- **INT-05-101:** `When a user adds an SSH public key via the settings page, the system shall store the key and compute its fingerprint.`
- **INT-05-102:** `When an SSH connection is established with a recognized public key, the system shall authenticate the session as the user who owns that key.`
- **INT-05-103:** `When an authenticated SSH session initiates a git-upload-pack command, the system shall serve the requested repository if the user has read access.`
- **INT-05-104:** `When an authenticated SSH session initiates a git-receive-pack command, the system shall accept the push if the user has write access.`
- **INT-05-105:** `When a user deletes an SSH key, the system shall remove the key immediately and reject future authentication with that key.`

### Optional Feature Requirements (Built-in SSH Server)

- **INT-05-201:** `Where the built-in SSH server is enabled via START_SSH_SERVER, the system shall listen on the configured SSH port for direct SSH connections.`
- **INT-05-202:** `Where the built-in SSH server is disabled, the system shall rely on the system SSH server and the authorized_keys file for authentication.`

### Unwanted Behaviour Requirements (SSH Errors)

- **INT-05-301:** `If an SSH connection presents an unrecognized public key, then the system shall reject the authentication.`
- **INT-05-302:** `If an SSH user attempts to access a repository without sufficient permissions, then the system shall reject the Git operation.`
- **INT-05-303:** `If a user adds an SSH key that is already registered to another account, then the system shall reject the key addition.`

---

## 6. Mailer / Email

**User Story:** As a user, I want to receive email notifications about repository activity and reply to issues via email so that I can participate in project discussions from my inbox.

### Ubiquitous Requirements (Email Properties)

- **INT-06-001:** `The system shall send email notifications in both HTML and plaintext formats.`
- **INT-06-002:** `The system shall apply a configurable subject prefix to all outgoing emails.`
- **INT-06-003:** `The system shall queue outgoing emails for asynchronous delivery.`

### Event-Driven Requirements (Outgoing Email)

- **INT-06-101:** `When an event triggers a notification (issue comment, PR review, release, team invite), the system shall send an email to subscribed users.`
- **INT-06-102:** `When a user enables email notifications for a repository, the system shall add the user to the notification mailing list for that repository.`
- **INT-06-103:** `When a user disables email notifications, the system shall stop sending event emails for that user.`
- **INT-06-104:** `When the mailer is configured for SMTP, the system shall connect to the SMTP server with optional TLS/SSL and authentication.`

### Optional Feature Requirements (Incoming Email)

- **INT-06-201:** `Where incoming email processing is configured, the system shall accept email replies and convert them to comments on the referenced issue or PR.`
- **INT-06-202:** `Where incoming email processing is configured, the system shall authenticate incoming email senders against known user email addresses.`
- **INT-06-203:** `Where sendmail is configured as the mailer backend, the system shall pipe outgoing emails to the sendmail binary.`

### Unwanted Behaviour Requirements (Email Errors)

- **INT-06-301:** `If the SMTP server is unreachable, then the system shall queue the email for retry and log the connection error.`
- **INT-06-302:** `If an incoming email reply has an invalid authentication token, then the system shall discard the email and log a warning.`
- **INT-06-303:** `If an outgoing email fails after the maximum retry attempts, then the system shall log the permanent failure and discard the message.`

---

## 7. Translation / i18n

**User Story:** As a user, I want to use Gitea in my preferred language so that I can navigate and understand the interface in my native tongue.

### Ubiquitous Requirements (i18n Properties)

- **INT-07-001:** `The system shall store translations in INI-format locale files keyed by language code.`
- **INT-07-002:** `The system shall provide English as the default and fallback language.`
- **INT-07-003:** `The system shall support over 30 languages with locale-specific translations.`
- **INT-07-004:** `The system shall support pluralization rules appropriate to each language.`

### Event-Driven Requirements (Language Workflow)

- **INT-07-101:** `When a user selects a language in their settings, the system shall persist the preference and render the interface in that language.`
- **INT-07-102:** `When a translation key is missing from the user's selected language, the system shall fall back to the English translation.`
- **INT-07-103:** `When a user with no language preference visits the site, the system shall detect the browser's Accept-Language header and apply the best matching locale.`

### Optional Feature Requirements (Custom Translations)

- **INT-07-201:** `Where custom locale files are provided in the custom/options/locale/ directory, the system shall load and merge them with the default translations.`
- **INT-07-202:** `Where the [i18n] LANGS configuration restricts the available languages, the system shall only offer those languages in the language selector.`

### Unwanted Behaviour Requirements (i18n Errors)

- **INT-07-301:** `If a locale file contains invalid INI syntax, then the system shall skip the file and log a warning.`
- **INT-07-302:** `If a user selects a language that is not available, then the system shall fall back to English.`

---

## 8. CAPTCHA

**User Story:** As an administrator, I want to protect registration and public forms from automated spam using CAPTCHA challenges.

### Ubiquitous Requirements (CAPTCHA Properties)

- **INT-08-001:** `The system shall require CAPTCHA verification on the user registration form when any CAPTCHA service is enabled.`
- **INT-08-002:** `The system shall validate CAPTCHA responses server-side against the configured CAPTCHA provider.`
- **INT-08-003:** `The system shall support exactly one active CAPTCHA service at a time based on configuration.`

### Event-Driven Requirements (CAPTCHA Workflow)

- **INT-08-101:** `When a user submits the registration form, the system shall validate the CAPTCHA response before processing the registration.`
- **INT-08-102:** `When the CAPTCHA response is invalid, the system shall reject the form submission with an error message.`
- **INT-08-103:** `When a CAPTCHA service is disabled, the system shall remove the CAPTCHA widget from the registration form.`

### Optional Feature Requirements (CAPTCHA Providers)

- **INT-08-201:** `Where hCaptcha is enabled, the system shall render the hCaptcha widget with the configured site key and theme.`
- **INT-08-202:** `Where reCAPTCHA v2 is enabled, the system shall render the reCAPTCHA checkbox or invisible widget with the configured site key.`
- **INT-08-203:** `Where Cloudflare Turnstile is enabled, the system shall render the Turnstile JavaScript widget with the configured site key.`
- **INT-08-204:** `Where mCaptcha is enabled, the system shall render the mCaptcha proof-of-work widget with the configured site key and instance URL.`

### Unwanted Behaviour Requirements (CAPTCHA Errors)

- **INT-08-301:** `If the CAPTCHA provider is unreachable during validation, then the system shall reject the form submission and display an error.`
- **INT-08-302:** `If the configured site key or secret is invalid, then the system shall log an error and fail CAPTCHA validation.`
- **INT-08-303:** `If a user bypasses the CAPTCHA widget (missing response), then the system shall reject the form submission.`

---

## 9. Avatar Systems

**User Story:** As a user, I want my avatar displayed across the interface from my preferred source (Gravatar, Libravatar, uploaded image, or generated) so that others can visually identify me.

### Ubiquitous Requirements (Avatar Properties)

- **INT-09-001:** `The system shall display user avatars in the interface for user mentions, comments, profiles, and commit history.`
- **INT-09-002:** `The system shall generate a default avatar using either an identicon pattern or the user's initial letter when no external avatar source is available.`
- **INT-09-003:** `The system shall enforce a configurable maximum file size for uploaded avatar images.`

### Event-Driven Requirements (Avatar Workflow)

- **INT-09-101:** `When a user uploads a custom avatar image, the system shall store the image locally and display it as the user's avatar.`
- **INT-09-102:** `When a user removes a custom uploaded avatar, the system shall fall back to the configured external avatar source or generated avatar.`
- **INT-09-103:** `When the system resolves a user avatar, the system shall check in order: local upload, federated avatar source, Gravatar source, then generated fallback.`

### Optional Feature Requirements (External Avatar Sources)

- **INT-09-201:** `Where Gravatar is enabled, the system shall look up the user's avatar by email hash from the configured Gravatar URL over HTTP or HTTPS.`
- **INT-09-202:** `Where Libravatar (federated avatar) is enabled, the system shall look up the user's avatar from the Libravatar federation using the user's email address.`
- **INT-09-203:** `Where a custom Gravatar source URL is configured, the system shall use that URL instead of the default gravatar.com domain.`

### Unwanted Behaviour Requirements (Avatar Errors)

- **INT-09-301:** `If the external avatar source (Gravatar or Libravatar) is unreachable, then the system shall fall back to the generated avatar.`
- **INT-09-302:** `If an uploaded avatar file exceeds the maximum size limit, then the system shall reject the upload with a validation error.`
- **INT-09-303:** `If an uploaded avatar file has an unsupported image format, then the system shall reject the upload.`

---

## 10. Proxy Support

**User Story:** As an administrator, I want Gitea to route outbound requests through an HTTP proxy and trust load balancer PROXY protocol headers so that it works correctly in restricted networks and behind reverse proxies.

### Ubiquitous Requirements (Proxy Properties)

- **INT-10-001:** `The system shall route outbound HTTP requests through the configured proxy URL when a proxy is set.`
- **INT-10-002:** `The system shall support bypass rules that exclude specified hosts from proxy routing.`

### Event-Driven Requirements (Proxy Workflow)

- **INT-10-101:** `When the system makes an outbound request (webhook, avatar lookup, update check), the system shall route the request through the configured proxy unless the target host matches a bypass rule.`
- **INT-10-102:** `When the PROXY protocol is enabled and a connection arrives from a trusted proxy IP, the system shall extract the original client IP from the PROXY protocol header.`
- **INT-10-103:** `When the PROXY protocol header indicates a v1 or v2 format, the system shall parse the appropriate version correctly.`

### Optional Feature Requirements (Proxy Configuration)

- **INT-10-201:** `Where proxy authentication is configured, the system shall include credentials in the proxy connection request.`
- **INT-10-202:** `Where trusted proxy IP ranges are configured, the system shall only accept PROXY protocol headers from those IP ranges.`
- **INT-10-203:** `Where no proxy is configured, the system shall make direct outbound connections.`

### Unwanted Behaviour Requirements (Proxy Errors)

- **INT-10-301:** `If the configured proxy is unreachable, then the system shall log the connection failure and fail the outbound request.`
- **INT-10-302:** `If a PROXY protocol header arrives from an untrusted IP, then the system shall ignore the header and use the direct connection IP.`
- **INT-10-303:** `If a PROXY protocol header contains an invalid format, then the system shall log a warning and use the direct connection IP.`

---

## 11. Update Checker

**User Story:** As an administrator, I want Gitea to check for new versions so that I am notified when updates are available.

### Ubiquitous Requirements (Update Checker Properties)

- **INT-11-001:** `The system shall compare the installed version against the latest available version from the configured endpoint.`
- **INT-11-002:** `The system shall check for updates on a periodic schedule.`

### Event-Driven Requirements (Update Workflow)

- **INT-11-101:** `When the update checker runs, the system shall query the configured endpoint URL for the latest version information.`
- **INT-11-102:** `When a newer version is available, the system shall display a notification to administrators in the web interface.`
- **INT-11-103:** `When the update checker is disabled via configuration, the system shall stop performing periodic version checks.`

### Optional Feature Requirements (Update Channels)

- **INT-11-201:** `Where the update channel is set to stable, the system shall only report stable releases as available updates.`
- **INT-11-202:** `Where the update channel is set to a pre-release channel (beta or dev), the system shall report pre-release versions as available updates.`
- **INT-11-203:** `Where a custom endpoint URL is configured, the system shall query that URL instead of the default Gitea update endpoint.`

### Unwanted Behaviour Requirements (Update Errors)

- **INT-11-301:** `If the update endpoint is unreachable, then the system shall log a warning and continue operation without update notifications.`
- **INT-11-302:** `If the update endpoint returns an invalid response, then the system shall log a warning and skip update processing.`
- **INT-11-303:** `If the installed version cannot be determined, then the system shall skip the update check.`

---

## 12. Git Smart HTTP

**User Story:** As a user, I want to perform Git clone, push, and pull operations over HTTP/HTTPS so that I can work with repositories through firewalls and proxies that block SSH.

### Ubiquitous Requirements (Smart HTTP Properties)

- **INT-12-001:** `The system shall implement the Git Smart HTTP protocol supporting git-upload-pack and git-receive-pack over HTTP/HTTPS.`
- **INT-12-002:** `The system shall authenticate Git HTTP operations using HTTP Basic Auth with username and password or access token.`
- **INT-12-003:** `The system shall support Git LFS (Large File Storage) operations over the same HTTP transport.`

### Event-Driven Requirements (Smart HTTP Workflow)

- **INT-12-101:** `When a client initiates a git clone over HTTPS, the system shall authenticate the user and serve the repository data via the smart HTTP protocol.`
- **INT-12-102:** `When a client pushes over HTTPS, the system shall authenticate the user, verify write permissions, and accept the pack data.`
- **INT-12-103:** `When a client performs a Git LFS operation, the system shall proxy LFS objects to the configured LFS storage backend.`
- **INT-12-104:** `When a Git operation is configured for verbose push output, the system shall return detailed push information to the client.`

### Optional Feature Requirements (Smart HTTP Configuration)

- **INT-12-201:** `Where Git garbage collection arguments are configured, the system shall apply those arguments during repository GC operations triggered by HTTP pushes.`
- **INT-12-202:** `Where fsync is enabled for Git operations, the system shall issue fsync calls on critical write operations for data durability.`
- **INT-12-203:** `Where partial diff is disabled, the system shall skip partial diff computation during HTTP request processing.`

### Unwanted Behaviour Requirements (Smart HTTP Errors)

- **INT-12-301:** `If a client attempts a Git push without write permission, then the system shall reject the push with HTTP 403.`
- **INT-12-302:** `If a client exceeds the maximum diff line limit during a Git operation, then the system shall truncate the diff output.`
- **INT-12-303:** `If Git LFS storage backend is unreachable during an LFS operation, then the system shall return an error to the client.`
- **INT-12-304:** `If a Git HTTP request has invalid credentials, then the system shall reject the request with HTTP 401.`

---

## Business Rules

- **BR-09-001:** The REST API is always served under /api/v1; no other version paths are recognized
- **BR-09-002:** API pagination default page size is configurable; maximum items per page is enforced server-side
- **BR-09-003:** Only one CAPTCHA service may be active at a time; enabling one disables the others
- **BR-09-004:** Template override hierarchy is strictly enforced: custom/templates/ > templates/{theme}/ > templates/
- **BR-09-005:** SSH clone URLs are constructed from [server] SSH_DOMAIN and SSH_PORT, independent of the HTTP host
- **BR-09-006:** Email notifications are queued and delivered asynchronously; delivery order is not guaranteed
- **BR-09-007:** English is the always-available fallback language; translation keys not found in the selected locale fall back to English
- **BR-09-008:** Avatar resolution follows a strict priority: local upload > federated avatar > Gravatar > generated fallback
- **BR-09-009:** PROXY protocol headers are only trusted from explicitly configured IP ranges
- **BR-09-010:** Update checker frequency is fixed; only the channel and endpoint are configurable
- **BR-09-011:** Git Smart HTTP authentication reuses the same token and credential system as the REST API
- **BR-09-012:** Custom static assets in custom/public/ always take precedence over embedded bindata assets
- **BR-09-013:** External markup renderers are executed in isolated processes with timeout enforcement

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| API request with expired OAuth2 token | Return HTTP 401 Unauthorized |
| Markdown file with mixed GFM and raw HTML | Sanitize HTML while preserving GFM features |
| Custom theme references deleted CSS file | Fall back to default gitea theme |
| Custom template with invalid Go syntax | Log error, fall back to default template |
| SSH key uploaded that matches another user | Reject key addition with duplicate error |
| SMTP server down during notification | Queue email for retry, log connection error |
| Incoming email with forged Reply-To | Discard email, log authentication failure |
| Locale file with invalid INI syntax | Skip file, log warning, use English fallback |
| CAPTCHA provider timeout during registration | Reject submission, display error to user |
| Gravatar service unreachable for avatar | Display generated identicon or initial avatar |
| PROXY protocol header from untrusted IP | Ignore header, use direct connection IP |
| Update check endpoint returns malformed JSON | Log warning, skip update notification |
| Git push to read-only repository over HTTP | Reject with HTTP 403 Forbidden |
| Git LFS object exceeds storage limit | Reject upload with appropriate error |
| External renderer process hangs | Timeout and display raw content with error |
| API paginated request with limit exceeding maximum | Clamp to MAX_RESPONSE_ITEMS and return results |

## Success Criteria

- REST API responds to authenticated requests within 200ms for standard CRUD operations
- Swagger UI loads completely with full endpoint documentation when enabled
- Markdown rendering completes in under 100ms for files under 1MB
- Syntax highlighting covers all languages supported by the Chroma engine
- Theme switching applies visually without page reload
- Custom template overrides take effect without server restart when files are updated
- SSH authentication and Git operations complete within 5 seconds for standard repositories
- Email notifications are queued within 1 second of the triggering event
- Language switching applies immediately on the next page load
- CAPTCHA validation completes within 3 seconds including network round-trip
- Avatar resolution falls back gracefully through the entire priority chain
- Proxy routing adds less than 100ms latency to outbound requests
- Update checker does not block or delay server startup
- Git Smart HTTP clone and push operations match SSH performance within 10%

# 07 — Platform Services

Baseline specification of Gitea's REST API, markup rendering, theming, customization, SSH, email, internationalization, proxy/camo, Git transport, feeds, server-sent events, federation/ActivityPub, and notification infrastructure. These are cross-cutting services and external interfaces consumed by all other domains. All requirements describe the current (v1.22.x) system behavior.

---

## 1. REST API (PLT-01)

**User Story:** As a developer, I want to interact with Gitea programmatically via a REST API so that I can automate repository, user, and organization management.

### Ubiquitous Requirements (API Properties)

- **PLT-01-001:** `The system shall expose a versioned REST API at the path prefix /api/v1.`
- **PLT-01-002:** `The system shall support the following authentication methods for API requests: HTTP Basic Auth, access token (query parameter or Authorization header), and OAuth2 bearer tokens.`
- **PLT-01-003:** `The system shall enforce token scope restrictions on every API request.`
- **PLT-01-004:** `The system shall paginate list responses using page and limit query parameters.`
- **PLT-01-005:** `The system shall respect a configurable maximum response items limit (default 50).`

### Event-Driven Requirements (API Workflow)

- **PLT-01-101:** `When a client sends an API request with valid credentials, the system shall process the request and return the appropriate response.`
- **PLT-01-102:** `When a client requests a paginated resource, the system shall return a page of results with HTTP Link headers for navigation.`
- **PLT-01-104:** `When a client sends a conditional request with an If-None-Match header, the system shall return HTTP 304 if the resource ETag has not changed.`
- **PLT-01-105:** `When the API is disabled via configuration, the system shall reject all API requests with HTTP 404.`

### Optional Feature Requirements (Swagger Documentation)

- **PLT-01-201:** `Where Swagger is enabled, the system shall serve an interactive Swagger UI at /api/swagger/.`
- **PLT-01-202:** `Where Swagger is enabled, the system shall generate an OpenAPI-compliant specification document covering all API endpoints.`
- **PLT-01-203:** `Where Swagger is disabled, the system shall return HTTP 404 for all Swagger UI and spec requests.`

### Unwanted Behaviour Requirements (API Errors)

- **PLT-01-301:** `If a client sends an API request without valid credentials, then the system shall reject the request with HTTP 401 Unauthorized.`
- **PLT-01-302:** `If a client sends an API request with credentials lacking the required scope, then the system shall reject the request with HTTP 403 Forbidden.`
- **PLT-01-303:** `If a client sends a request to a non-existent API endpoint, then the system shall return HTTP 404 Not Found.`
- **PLT-01-304:** `If a client submits an invalid request body, then the system shall return HTTP 400 with a descriptive error message.`

---

## 2. Markup Renderers (PLT-02)

**User Story:** As a user, I want rendered documents in multiple formats (Markdown, reStructuredText, AsciiDoc, Org) with syntax highlighting so that I can read rich content in repositories, issues, and comments.

### Ubiquitous Requirements (Rendering Engine)

- **PLT-02-001:** `The system shall render Markdown content using GitHub-Flavored Markdown (GFM) with support for tables, task lists, strikethrough, and autolinks.`
- **PLT-02-002:** `The system shall apply syntax highlighting to code blocks using the Chroma engine supporting over 100 programming languages.`
- **PLT-02-003:** `The system shall sanitize all rendered HTML to prevent cross-site scripting (XSS) attacks.`
- **PLT-02-004:** `The system shall provide a copy button on rendered code blocks for easy clipboard copying.`

### Event-Driven Requirements (Rendering Workflow)

- **PLT-02-101:** `When a user views a file with a recognized markup extension, the system shall render the file content as formatted HTML instead of displaying raw source.`
- **PLT-02-103:** `When a user views a CSV file, the system shall render the content as an interactive table.`
- **PLT-02-105:** `When Markdown content contains mathematical notation, the system shall render the notation using KaTeX.`
- **PLT-02-106:** `When Markdown content contains a Mermaid diagram code block, the system shall render the diagram, truncating source exceeding the configured MERMAID_MAX_SOURCE_CHARACTERS (default 5000).`
- **PLT-02-107:** `When a user views an Asciicast terminal recording file (.cast), the system shall render an embedded asciinema player referencing the file's raw URL.`
- **PLT-02-108:** `When a user views a console session file (.sh-session) or a file containing ANSI escape sequences, the system shall render the terminal output as styled HTML.`

### Optional Feature Requirements (Additional Renderers)

- **PLT-02-201:** `Where reStructuredText rendering is enabled, the system shall render .rst files as formatted HTML.`
- **PLT-02-202:** `Where AsciiDoc rendering is enabled, the system shall render .adoc files with attribute expansion and macro support.`
- **PLT-02-203:** `Where Org mode rendering is enabled, the system shall render .org files as formatted HTML.`
- **PLT-02-204:** `Where a custom external renderer is configured via [markup.*] sections, the system shall execute the external command and render the output as HTML.`
- **PLT-02-205:** `Where custom URL schemes are configured, the system shall convert matching URLs into clickable links in rendered Markdown.`

### Unwanted Behaviour Requirements (Rendering Errors)

- **PLT-02-301:** `If a markup file contains invalid syntax, then the system shall display the raw source with an error indication.`
- **PLT-02-302:** `If an external renderer command fails or times out, then the system shall display the raw file content with an error message.`
- **PLT-02-303:** `If rendered content exceeds the maximum allowed output size, then the system shall truncate the output and display a warning.`

---

## 3. Theme System (PLT-03)

**User Story:** As a user, I want to select and apply visual themes so that the Gitea interface matches my preference for light, dark, or custom appearance.

> **Note:** The system provides five built-in themes with `gitea-auto` (system preference) as the default for unauthenticated users and users with no preference set.

### Ubiquitous Requirements (Theme Properties)

- **PLT-03-001:** `The system shall provide five built-in themes: gitea-light, gitea-dark, gitea-auto (system preference), gitea-light-protanopia-deuteranopia, and gitea-dark-protanopia-deuteranopia.`
- **PLT-03-002:** `The system shall apply themes through CSS stylesheets without requiring page reload.`
- **PLT-03-003:** `The system shall store the user's theme preference in user settings.`
- **PLT-03-004:** `The system shall apply the configured default theme (default gitea-auto) to unauthenticated users and users with no preference set.`

### Event-Driven Requirements (Theme Workflow)

- **PLT-03-101:** `When a user selects a theme from the settings page, the system shall apply the theme immediately and persist the preference.`
- **PLT-03-102:** `When the gitea-auto theme is active and the operating system switches between light and dark mode, the system shall switch the interface accordingly.`

### Optional Feature Requirements (Custom Themes)

- **PLT-03-201:** `Where custom theme CSS files are placed in the custom/public/assets/css directory, the system shall detect and list them as available themes.`
- **PLT-03-202:** `Where the [ui] THEMES configuration lists additional theme names, the system shall include those themes in the theme selector dropdown.`
- **PLT-03-203:** `Where a custom theme name matches a CSS file in the custom assets directory, the system shall load that stylesheet as the selected theme.`

### Unwanted Behaviour Requirements (Theme Errors)

- **PLT-03-301:** `If a user's stored theme references a theme that is no longer available, then the system shall fall back to the configured default theme.`
- **PLT-03-302:** `If a custom theme CSS file fails to load, then the system shall fall back to the configured default theme.`

---

## 4. Custom Assets (PLT-04)

**User Story:** As an administrator, I want to override templates, add custom static files, and modify branding so that I can tailor the Gitea instance to my organization.

### Ubiquitous Requirements (Override Hierarchy)

- **PLT-04-001:** `The system shall resolve template files using a three-level priority: custom/templates/ (highest), templates/{theme}/, templates/ (default).`
- **PLT-04-002:** `The system shall serve static files from the custom/public/ directory with higher priority than built-in assets.`
- **PLT-04-003:** `The system shall load locale files from the custom/options/locale/ directory to supplement or override default translations.`
- **PLT-04-004:** `The system shall embed built-in templates and static assets as bindata for single-binary deployment.`

### Event-Driven Requirements (Asset Workflow)

- **PLT-04-101:** `When the system resolves a template path, the system shall check the custom/templates/ directory first before falling back to built-in templates.`
- **PLT-04-102:** `When a request targets a static asset, the system shall serve the file from custom/public/ if it exists; otherwise, serve the embedded bindata asset.`
- **PLT-04-103:** `When custom locale files are present, the system shall merge custom translations with default translations at startup.`

### Optional Feature Requirements (Branding)

- **PLT-04-201:** `Where a custom logo file is placed in the custom assets directory, the system shall display the custom logo in place of the default Gitea logo.`
- **PLT-04-202:** `Where the footer template is overridden, the system shall render the custom footer content on all pages.`
- **PLT-04-203:** `Where custom links are configured, the system shall display the links in the navigation or footer area.`

### Unwanted Behaviour Requirements (Asset Errors)

- **PLT-04-301:** `If a custom template file contains invalid Go template syntax, then the system shall log an error and fall back to the default template.`
- **PLT-04-302:** `If a custom static file exceeds the maximum allowed size, then the system shall return HTTP 500 and log an error.`
- **PLT-04-303:** `If a custom locale file contains invalid INI syntax, then the system shall skip the file and use default translations with a logged warning.`

---

## 5. SSH Server (PLT-05)

**User Story:** As a user, I want to perform Git operations over SSH so that I can securely clone, push, and pull repositories.

### Ubiquitous Requirements (SSH Properties)

- **PLT-05-001:** `The system shall support Git protocol operations over SSH including git-upload-pack and git-receive-pack.`
- **PLT-05-002:** `The system shall manage SSH public keys per user with fingerprint tracking.`
- **PLT-05-003:** `The system shall construct SSH clone URLs using the configured SSH domain and port.`

### Event-Driven Requirements (SSH Workflow)

- **PLT-05-101:** `When a user adds an SSH public key via the settings page, the system shall store the key and compute its fingerprint.`
- **PLT-05-102:** `When an SSH connection is established with a recognized public key, the system shall authenticate the session as the user who owns that key.`
- **PLT-05-103:** `When an authenticated SSH session initiates a git-upload-pack command, the system shall serve the requested repository if the user has read access.`
- **PLT-05-104:** `When an authenticated SSH session initiates a git-receive-pack command, the system shall accept the push if the user has write access.`
- **PLT-05-105:** `When a user deletes an SSH key, the system shall remove the key immediately and reject future authentication with that key.`

### Optional Feature Requirements (Built-in SSH Server)

- **PLT-05-201:** `Where the built-in SSH server is enabled via START_SSH_SERVER, the system shall listen on the configured SSH port for direct SSH connections.`
- **PLT-05-202:** `Where the built-in SSH server is disabled, the system shall rely on the system SSH server and the authorized_keys file for authentication.`

### Unwanted Behaviour Requirements (SSH Errors)

- **PLT-05-301:** `If an SSH connection presents an unrecognized public key, then the system shall reject the authentication.`
- **PLT-05-302:** `If an SSH user attempts to access a repository without sufficient permissions, then the system shall reject the Git operation.`
- **PLT-05-303:** `If a user adds an SSH key that is already registered to another account, then the system shall reject the key addition.`

---

## 6. Mailer / Email (PLT-06)

**User Story:** As a user, I want to receive email notifications about repository activity and reply to issues via email so that I can participate in project discussions from my inbox.

### Ubiquitous Requirements (Email Properties)

- **PLT-06-001:** `The system shall send email notifications in both HTML and plaintext formats.`
- **PLT-06-002:** `The system shall apply a configurable subject prefix to all outgoing emails.`
- **PLT-06-003:** `The system shall queue outgoing emails for asynchronous delivery.`

### Event-Driven Requirements (Outgoing Email)

- **PLT-06-101:** `When an event triggers a notification (issue comment, PR review, release, team invite), the system shall send an email to subscribed users.`
- **PLT-06-102:** `When a user enables email notifications for a repository, the system shall add the user to the notification mailing list for that repository.`
- **PLT-06-103:** `When a user disables email notifications, the system shall stop sending event emails for that user.`
- **PLT-06-104:** `When the mailer is configured for SMTP, the system shall connect to the SMTP server with optional TLS/SSL and authentication.`

### Optional Feature Requirements (Incoming Email)

- **PLT-06-201:** `Where incoming email processing is configured, the system shall accept email replies and convert them to comments on the referenced issue or PR.`
- **PLT-06-202:** `Where incoming email processing is configured, the system shall authenticate incoming email senders against known user email addresses.`
- **PLT-06-203:** `Where sendmail is configured as the mailer backend, the system shall pipe outgoing emails to the sendmail binary.`

### Unwanted Behaviour Requirements (Email Errors)

- **PLT-06-301:** `If the SMTP server is unreachable, then the system shall queue the email for retry and log the connection error.`
- **PLT-06-302:** `If an incoming email reply has an invalid authentication token, then the system shall discard the email and log a warning.`
- **PLT-06-303:** `If an outgoing email fails after the maximum retry attempts, then the system shall log the permanent failure and discard the message.`

---

## 7. Translation / i18n (PLT-07)

**User Story:** As a user, I want to use Gitea in my preferred language so that I can navigate and understand the interface in my native tongue.

> **Note:** The system supports approximately 28 languages with locale-specific translations.

### Ubiquitous Requirements (i18n Properties)

- **PLT-07-001:** `The system shall store translations in INI-format locale files keyed by language code.`
- **PLT-07-002:** `The system shall provide English as the default and fallback language.`
- **PLT-07-003:** `The system shall support approximately 28 languages with locale-specific translations.`
- **PLT-07-004:** `The system shall support pluralization rules appropriate to each language.`

### Event-Driven Requirements (Language Workflow)

- **PLT-07-101:** `When a user selects a language in their settings, the system shall persist the preference and render the interface in that language.`
- **PLT-07-102:** `When a translation key is missing from the user's selected language, the system shall fall back to the English translation.`
- **PLT-07-103:** `When a user with no language preference visits the site, the system shall detect the browser's Accept-Language header and apply the best matching locale.`

### Optional Feature Requirements (Custom Translations)

- **PLT-07-201:** `Where custom locale files are provided in the custom/options/locale/ directory, the system shall load and merge them with the default translations.`
- **PLT-07-202:** `Where the [i18n] LANGS configuration restricts the available languages, the system shall only offer those languages in the language selector.`

### Unwanted Behaviour Requirements (i18n Errors)

- **PLT-07-301:** `If a locale file contains invalid INI syntax, then the system shall skip the file and log a warning.`
- **PLT-07-302:** `If a user selects a language that is not available, then the system shall fall back to English.`

---

## 8. Proxy / Camo (PLT-08)

**User Story:** As an administrator, I want external images rendered in issues and comments to be proxied through Gitea so that user IP addresses are not leaked to third-party image hosts.

### Ubiquitous Requirements (Camo Properties)

- **PLT-08-001:** `The system shall rewrite external image URLs in rendered markdown content to route through the configured camo proxy URL.`
- **PLT-08-002:** `The system shall preserve the original URL as an encrypted or HMAC-signed query parameter so the proxy can fetch the original resource.`
- **PLT-08-003:** `The system shall leave relative URLs, same-origin URLs, and (by default) HTTPS URLs unchanged; only absolute non-same-origin HTTP image URLs are rewritten unless ALLWAYS is set.`

### Event-Driven Requirements (Camo Workflow)

- **PLT-08-101:** `When the system renders markdown containing an external image URL and camo is enabled, the system shall rewrite the URL to the camo proxy prefix.`
- **PLT-08-102:** `When the camo proxy is configured with a shared secret, the system shall include an HMAC signature of the original URL in the rewritten request.`
- **PLT-08-103:** `When the system renders markdown with camo disabled, the system shall emit the original image URL unchanged.`

### Optional Feature Requirements (Camo Configuration)

- **PLT-08-201:** `Where camo is enabled and ALLWAYS is false (default), the system shall rewrite only absolute non-same-origin HTTP image URLs, leaving HTTPS URLs unchanged.`
- **PLT-08-202:** `Where camo is disabled (default), the system shall emit external image URLs as-is.`

### Unwanted Behaviour Requirements (Camo Errors)

- **PLT-08-301:** `If camo is enabled but the camo proxy URL is not configured, then the system shall fall back to emitting original URLs and log a warning.`
- **PLT-08-302:** `If the camo proxy returns an error for a given image, then the browser shall display a broken-image indicator; the system shall not retry inline.`
- **PLT-08-303:** `If an image URL fails parsing, then the system shall leave the URL unchanged rather than emit a malformed proxy URL.`

---

## 9. Server-Sent Events (PLT-09)

**User Story:** As a browser client, I want a persistent event stream so that I receive real-time updates (notifications, status changes) without polling.

### Ubiquitous Requirements (SSE Properties)

- **PLT-09-001:** `The system shall expose a /user/events endpoint delivering server-sent events over a persistent HTTP connection.`
- **PLT-09-002:** `The system shall frame each event with the SSE wire format: a leading event: line, a data: line, and a blank-line terminator.`
- **PLT-09-003:** `The system shall send periodic heartbeat events as a named "ping" event (every 30 seconds) to keep intermediary proxies from closing idle connections.`
- **PLT-09-004:** `The system shall scope every event payload to the authenticated user's notifications, mentions, status changes, and timeline activity.`

### Event-Driven Requirements (SSE Workflow)

- **PLT-09-101:** `When a browser opens /user/events with a valid session, the system shall register the connection with the per-user eventsource manager.`
- **PLT-09-102:** `When an event of interest to the user occurs (notification, status change), the system shall push the event to every open SSE connection for that user.`
- **PLT-09-103:** `When the browser closes the connection, the system shall deregister the connection and free associated resources.`
- **PLT-09-104:** `When the server shuts down, the system shall close every open SSE connection with a finalizing event before terminating.`

### State-Driven Requirements (SSE Lifecycle)

- **PLT-09-701:** `While a user has multiple browser sessions open, the system shall deliver every event to each session independently.`
- **PLT-09-702:** `While a user's session expires, the system shall close the SSE connection and require re-authentication.`

### Unwanted Behaviour Requirements (SSE Errors)

- **PLT-09-301:** `If an unauthenticated client requests /user/events, then the system shall return HTTP 200 with an SSE event named "close" containing the data "unauthorized".`
- **PLT-09-303:** `If an event payload cannot be serialized, then the system shall log the error and skip the event rather than close the connection.`

---

## 10. RSS & Atom Feeds (PLT-10)

**User Story:** As a user or feed reader, I want RSS and Atom feeds for repository activity, releases, tags, branches, and user activity so that I can monitor projects from my feed reader.

### Ubiquitous Requirements (Feed Properties)

- **PLT-10-001:** `The system shall expose RSS (.rss) and Atom (.atom) variants for every supported feed.`
- **PLT-10-002:** `The system shall support feeds for: user activity, repository commits, repository releases, repository tags, and per-branch commits.`
- **PLT-10-003:** `The system shall include the most recent N entries (configurable) per feed, each with title, link, summary, author, and updated timestamp.`
- **PLT-10-004:** `The system shall honor the If-Modified-Since header and return HTTP 304 when no entries have changed since the supplied timestamp.`

### Event-Driven Requirements (Feed Workflow)

- **PLT-10-101:** `When a client requests GET /{username}.rss or .atom, the system shall return the user's public activity feed filtered by visibility.`
- **PLT-10-102:** `When a client requests GET /{owner}/{repo}/commits/{branch}.rss or .atom, the system shall return the per-branch commit feed.`
- **PLT-10-103:** `When a client requests GET /{owner}/{repo}/releases/.rss or .atom, the system shall return the release feed scoped to non-draft releases.`
- **PLT-10-104:** `When a client requests GET /{owner}/{repo}/tags/.rss or .atom, the system shall return the tag feed.`
- **PLT-10-105:** `When a client requests a feed for a private repository, the system shall require authentication and scope entries to the requester's permission level.`

### Optional Feature Requirements (Feed Configuration)

- **PLT-10-201:** `Where feeds are disabled via configuration, the system shall return HTTP 404 for all feed routes.`
- **PLT-10-202:** `Where a feed is requested with a ?limit query parameter below the maximum, the system shall honor the requested limit.`

### Unwanted Behaviour Requirements (Feed Errors)

- **PLT-10-301:** `If an unauthenticated client requests a private repository's feed, then the system shall return HTTP 404 to avoid leaking repository existence.`
- **PLT-10-302:** `If a feed references a non-existent user, repository, branch, or tag, then the system shall return HTTP 404.`
- **PLT-10-303:** `If a feed entry references content that has since been deleted, then the system shall exclude the entry from the response.`

---

## 11. Federation / ActivityPub (PLT-11)

**User Story:** As a user on a federated Gitea instance, I want my repository activity to be discoverable over ActivityPub so that users on other instances can follow and interact.

### Ubiquitous Requirements (Federation Properties)

- **PLT-11-001:** `The system shall implement ActivityPub client functionality with HTTP Signature support.`
- **PLT-11-002:** `The system shall use Activity Streams 2.1 content types for federation messages.`
- **PLT-11-003:** `The system shall support NodeInfo 2.1 schema for instance metadata publishing.`

### Event-Driven Requirements (Federation Workflow)

- **PLT-11-101:** `When a remote platform sends a WebFinger query for a local user, the system shall return a JRD response with profile, avatar, and ActivityPub links.`
- **PLT-11-102:** `When the system sends outgoing federation requests, the system shall sign the requests with HTTP Signatures using configured headers.`
- **PLT-11-103:** `When the system receives incoming signed requests, the system shall verify the HTTP Signature before processing.`

### Optional Feature Requirements (Federation Configuration)

- **PLT-11-201:** `Where federation is enabled, the system shall expose ActivityPub, WebFinger, and NodeInfo endpoints.`
- **PLT-11-202:** `Where federation is disabled, the system shall return 404 for all federation endpoints.`
- **PLT-11-203:** `Where the [federation] ENABLED setting is true (default false), the system shall expose an ActivityPub API under /api/v1/activitypub for user and repository actors.`
- **PLT-11-204:** `Where federation is enabled, the system shall sign outbound federation requests using HTTP Signatures with the configured algorithms (default rsa-sha256, rsa-sha512, ed25519).`
- **PLT-11-205:** `Where SHARE_USER_STATISTICS is enabled (default true), the system may publish aggregated user statistics to the federation.`

### Unwanted Behaviour Requirements (Federation Errors)

- **PLT-11-301:** `If a WebFinger query targets a private user, then the system shall return a 404 response.`
- **PLT-11-302:** `If an incoming federation request has an invalid signature, then the system shall reject the request with 401.`
- **PLT-11-303:** `If a federation payload exceeds MAX_SIZE (default 4 MiB), the system shall reject it.`
- **PLT-11-304:** `If DIGEST_ALGORITHM is unsupported, the system shall refuse to start.`

---

## 12. Notification Infrastructure (PLT-12)

**User Story:** As the platform, I need a central event-dispatch bus so that every domain's notification triggers (issue created, PR merged, push received) are routed to the correct downstream notifiers without each domain reimplementing delivery.

> **Scope note:** Each domain documents its own notification triggers. This section documents the cross-cutting notification dispatcher infrastructure only.

### Ubiquitous Requirements (Notification Dispatcher Properties)

- **PLT-12-001:** `The system shall provide a central notification dispatcher (services/notify) that accepts domain events and routes them to all registered notifiers.`
- **PLT-12-002:** `The system shall register the following notifiers at startup: mailer, webhook, uinotification (in-app), indexer, mirror, and automerge.`
- **PLT-12-003:** `The system shall deliver each event to every registered notifier independently, so that a failure in one notifier does not block delivery to others.`

### Event-Driven Requirements (Dispatch Workflow)

- **PLT-12-101:** `When a domain service raises a notification event (e.g., issue created, push received, PR merged), the system shall pass the event to the dispatcher, which fans it out to all registered notifiers.`
- **PLT-12-102:** `When a notifier completes processing of an event, the system shall log the outcome without affecting other notifiers' delivery.`

### Unwanted Behaviour Requirements (Dispatch Errors)

- **PLT-12-301:** `If a notifier returns an error during event processing, then the system shall log the error and continue delivering the event to the remaining notifiers.`

---

## 13. Git Transport (PLT-13)

**User Story:** As a user, I want to perform Git clone, push, and pull operations over HTTP/HTTPS and SSH so that I can work with repositories through firewalls and proxies that block SSH, or directly via SSH when preferred.

> **Scope note:** Git operations (branch, tag, commit, merge) are documented in Domain 03 (Repository & Code). This section documents the transport protocols (HTTP smart protocol and SSH git commands) only.

### Ubiquitous Requirements (Smart HTTP Properties)

- **PLT-13-001:** `The system shall implement the Git Smart HTTP protocol supporting git-upload-pack and git-receive-pack over HTTP/HTTPS.`
- **PLT-13-002:** `The system shall authenticate Git HTTP operations using HTTP Basic Auth with username and password or access token.`
- **PLT-13-003:** `The system shall support Git LFS (Large File Storage) operations over the same HTTP transport.`

### Event-Driven Requirements (Smart HTTP Workflow)

- **PLT-13-101:** `When a client initiates a git clone over HTTPS, the system shall authenticate the user and serve the repository data via the smart HTTP protocol.`
- **PLT-13-102:** `When a client pushes over HTTPS, the system shall authenticate the user, verify write permissions, and accept the pack data.`
- **PLT-13-103:** `When a client performs a Git LFS operation, the system shall proxy LFS objects to the configured LFS storage backend.`
- **PLT-13-104:** `When a Git operation is configured for verbose push output, the system shall return detailed push information to the client.`

### Optional Feature Requirements (Smart HTTP Configuration)

- **PLT-13-201:** `Where Git garbage collection arguments are configured, the system shall apply those arguments during repository GC operations triggered by HTTP pushes.`
- **PLT-13-202:** `Where fsync is enabled for Git operations, the system shall issue fsync calls on critical write operations for data durability.`
- **PLT-13-203:** `Where partial diff is disabled, the system shall skip partial diff computation during HTTP request processing.`

### Unwanted Behaviour Requirements (Smart HTTP Errors)

- **PLT-13-301:** `If a client attempts a Git push without write permission, then the system shall reject the push with HTTP 403.`
- **PLT-13-302:** `If a client exceeds the maximum diff line limit during a Git operation, then the system shall truncate the diff output.`
- **PLT-13-303:** `If Git LFS storage backend is unreachable during an LFS operation, then the system shall return an error to the client.`
- **PLT-13-304:** `If a Git HTTP request has invalid credentials, then the system shall reject the request with HTTP 401.`

---

## Configuration Reference

The following INI sections configure platform service behaviors. Add to or modify these via `custom/conf/app.ini` or environment overrides using the `GITEA__section__key` format.

### [server] Section (SSH & Proxy Protocol)

- **SSH_DOMAIN**: Domain used in SSH clone URLs (default same as `DOMAIN`).
- **SSH_PORT**: Port shown in SSH clone URLs (default 22; set to the exposed port even if the built-in server listens elsewhere).
- **SSH_LISTEN_PORT**: Port the built-in SSH server listens on (default same as `SSH_PORT`).
- **START_SSH_SERVER**: Enable the built-in SSH server (default false).
- **USE_PROXY_PROTOCOL**: Enable PROXY protocol header parsing on all incoming connections (default false). Applied globally to all listeners.
- **REVERSE_PROXY_TRUSTED_PROXIES**: Comma-separated CIDR ranges of trusted reverse proxies whose `X-Forwarded-For` headers will be honored (default `127.0.0.0/8,::1/128`).
- **REVERSE_PROXY_LIMIT**: Number of proxy hops to traverse when resolving client IP (default 0).
- **ENABLE_PPROF**: Enable the pprof profiling endpoints (default false). See Domain 08 (Administration & Operations).

### [ssh.minimum_key_sizes] Section

Per-key-type minimum bit-length policy enforced when a user registers an SSH key. (DSA is not supported by Gitea.)
- **ED25519**: Minimum bits for ED25519 keys (default 256).
- **ED25519_SK**: Minimum bits for ED25519 hardware security key (default 256).
- **ECDSA**: Minimum bits for ECDSA keys (default 256).
- **ECDSA_SK**: Minimum bits for ECDSA hardware security key (default 256).
- **RSA**: Minimum bits for RSA keys (default 3071).

### [mailer] Section

Controls the outgoing email transport.
- **ENABLED**: Enable the mailer (default false).
- **PROTOCOL**: Mailer protocol (`smtp`, `smtps`, `smtp+startls`, `smtp+unix`, `sendmail`, `dummy`; default `smtp`).
- **SMTP_ADDR**: SMTP server address (default empty).
- **SMTP_PORT**: SMTP server port (default empty).
- **CLIENT_CERT_FILE** / **CLIENT_KEY_FILE**: TLS client certificate and key for SMTP mutual TLS.
- **FORCE_TRUST_SERVER_CERT**: Trust any SMTP server certificate (default false).
- **USER** / **PASSWD**: SMTP authentication credentials.
- **FROM**: Sender address used in outgoing emails (required when enabled).
- **SUBJECT_PREFIX**: Prefix applied to all outgoing email subjects (default empty).
- **SEND_AS_PLAIN_TEXT**: Send emails as plaintext only (default true).
- **SENDMAIL_PATH**: Path to the sendmail binary (default `sendmail`).
- **SENDMAIL_TIMEOUT**: Timeout for sendmail operations (default `5m`).
- **MAILER_TYPE**: Deprecated alias for `PROTOCOL`.

### [markup.*] Section (per-renderer)

Each external renderer is configured as a `[markup.<name>]` section.
- **ENABLED**: Activate this renderer.
- **FILE_EXTENSIONS**: Comma-separated extensions handled by this renderer.
- **RENDER_COMMAND**: External command to execute on file content.
- **IS_INPUT_FILE**: Whether the renderer takes a filename (true) or stdin (false).
- **RENDER_CONTENT_MODE**: How output is wrapped (`sanitized`, `no-sanitizer`, `iframe`).

### [camo] Section

Configuration for the go-camo-style image proxy. When ENABLED, both SERVER_URL and HMAC_KEY are required.
- **ENABLED**: Activate camo image proxying (default false).
- **SERVER_URL**: Base URL of the camo proxy.
- **HMAC_KEY**: Shared HMAC-SHA1 secret used to sign proxied URLs.
- **ALLWAYS**: If true, rewrite all absolute non-same-origin image URLs (including HTTPS); if false (default), rewrite only HTTP URLs.

### [i18n] Section

- **LANGS**: Comma-separated list of locale codes offered in the language selector (default full built-in list of ~28 languages).
- **NAMES**: Comma-separated human-readable names matching each `LANGS` entry.

### [federation] Section

Controls ActivityPub federation. Disabled by default.
- **ENABLED**: Enable federation endpoints and signing (default false).
- **SHARE_USER_STATISTICS**: Publish aggregated user statistics (default true).
- **MAX_SIZE**: Maximum federation payload size in MiB (default 4).
- **ALGORITHMS**: Supported HTTP-Signature algorithms (default `rsa-sha256,rsa-sha512,ed25519`).
- **DIGEST_ALGORITHM**: Digest algorithm for federation (default `SHA-256`).
- **GET_HEADERS** / **POST_HEADERS**: Headers covered by the HTTP signature.

---

## Business Rules

- **BR-07-001:** The REST API is always served under /api/v1; no other version paths are recognized
- **BR-07-002:** API pagination default page size is configurable; maximum items per page is enforced server-side
- **BR-07-003:** Template override hierarchy is strictly enforced: custom/templates/ > templates/{theme}/ > templates/
- **BR-07-004:** SSH clone URLs are constructed from [server] SSH_DOMAIN and SSH_PORT, independent of the HTTP host
- **BR-07-005:** Email notifications are queued and delivered asynchronously; delivery order is not guaranteed
- **BR-07-006:** English is the always-available fallback language; translation keys not found in the selected locale fall back to English
- **BR-07-007:** USE_PROXY_PROTOCOL is a global boolean applied to all connections; the REVERSE_PROXY_TRUSTED_PROXIES list governs only X-Forwarded-For header trust, not PROXY protocol parsing
- **BR-07-008:** Git Smart HTTP authentication reuses the same token and credential system as the REST API
- **BR-07-009:** Custom static assets in custom/public/ always take precedence over embedded bindata assets
- **BR-07-010:** External markup renderers are executed in isolated processes with timeout enforcement
- **BR-07-011:** Federation respects user visibility settings (private users are not discoverable)
- **BR-07-012:** The notification dispatcher fans out events to all registered notifiers; one notifier failure does not block others

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| API request with expired OAuth2 token | Return HTTP 401 Unauthorized |
| Markdown file with mixed GFM and raw HTML | Sanitize HTML while preserving GFM features |
| Custom theme references deleted CSS file | Fall back to configured default theme |
| Custom template with invalid Go syntax | Log error, fall back to default template |
| SSH key uploaded that matches another user | Reject key addition with duplicate error |
| SMTP server down during notification | Queue email for retry, log connection error |
| Incoming email with forged Reply-To | Discard email, log authentication failure |
| Locale file with invalid INI syntax | Skip file, log warning, use English fallback |
| X-Forwarded-For header from untrusted proxy IP | Ignore header, use direct connection IP |
| Git push to read-only repository over HTTP | Reject with HTTP 403 Forbidden |
| Git LFS object exceeds storage limit | Reject upload with appropriate error |
| External renderer process hangs | Timeout and display raw content with error |
| API paginated request with limit exceeding maximum | Clamp to MAX_RESPONSE_ITEMS and return results |
| Camo proxy enabled but SERVER_URL missing | Fall back to original URLs, log warning |
| Federation payload exceeds MAX_SIZE | Reject payload with error |
| Notification notifier raises an error | Log error, continue delivering to other notifiers |

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
- Proxy routing adds less than 100ms latency to outbound requests
- Git Smart HTTP clone and push operations match SSH performance within 10%
- Notification dispatcher fans out events to all notifiers within 100ms
- Federation signature verification completes in under 100ms

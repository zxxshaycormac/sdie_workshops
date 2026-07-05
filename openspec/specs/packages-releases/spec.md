# 06 — Packages & Releases

Baseline specification of Gitea's package registry, container registry, release management, LFS, and attachments. All requirements describe the current (v1.22.x) system behavior.

**Storage scope note:** Storage backend configuration (local, MinIO) is documented in Domain 08 (Administration & Operations). This domain documents storage usage by packages, LFS, and attachments.

---

## 1. Package Registry Overview (PKG-01)

**User Story:** As a developer, I want to publish and consume packages through a built-in registry so that my organization can manage dependencies without relying on external services.

### Ubiquitous Requirements (Registry Architecture)

- **PKG-01-001:** `The system shall provide a multi-format package registry supporting 21 package types: Alpine, Cargo, Chef, Composer, Conan, Conda, Container, CRAN, Debian, Generic, Go, Helm, Maven, npm, NuGet, Pub, PyPI, RPM, RubyGems, Swift, and Vagrant.`
- **PKG-01-002:** `The system shall store package content using content-addressable storage with SHA256 deduplication.`
- **PKG-01-008:** `The system shall compute MD5, SHA1, SHA256, and SHA512 checksums simultaneously for each package blob via a multi-hasher, persisting all four hashes on the blob record for use by type-specific registry metadata.`
- **PKG-01-003:** `The system shall organize stored content under the relative path aa/bb/<full-sha256>, where aa and bb are the first and second two-character segments of the SHA256 key and the full hash is repeated as the filename (e.g. aabb0000... → aa/bb/aabb0000...).`
- **PKG-01-004:** `The system shall scope all packages to an owner (user or organization).`
- **PKG-01-005:** `The system shall require authentication for all package operations (publish, install, delete).`
- **PKG-01-006:** `The system shall support JWT token-based, basic auth, and OAuth2 authentication for package API endpoints.`
- **PKG-01-007:** `The system shall enforce package access using read:package and write:package token scopes.`

### Event-Driven Requirements (Registry Workflow)

- **PKG-01-101:** `When a client publishes a package version, the system shall store the package content, create version metadata, and index the package for search.`
- **PKG-01-102:** `When a client requests a package version that already exists with identical content, the system shall return the existing content without duplication.`
- **PKG-01-103:** `When a client deletes a package version, the system shall remove the version metadata and decrement content reference counts.`
- **PKG-01-104:** `When content reference counts reach zero, the system shall remove the unreferenced content blob from storage.`
- **PKG-01-105:** `When a client accesses a package listing endpoint, the system shall return packages scoped to the authenticated user's access level.`

### Optional Feature Requirements (Registry Configuration)

- **PKG-01-201:** `Where the [packages] ENABLED setting is false, the system shall disable the package registry and return 404 for all package endpoints.`
- **PKG-01-202:** `Where MinIO/S3 storage is configured with SERVE_DIRECT, the system shall generate direct download URLs that bypass the Gitea server.`
- **PKG-01-203:** `Where per-owner quotas are configured via LIMIT_TOTAL_OWNER_COUNT and LIMIT_TOTAL_OWNER_SIZE, and per-package-type quotas are configured via the 21 LIMIT_SIZE_* keys (LIMIT_SIZE_ALPINE, LIMIT_SIZE_CARGO, LIMIT_SIZE_CHEF, LIMIT_SIZE_COMPOSER, LIMIT_SIZE_CONAN, LIMIT_SIZE_CONDA, LIMIT_SIZE_CONTAINER, LIMIT_SIZE_CRAN, LIMIT_SIZE_DEBIAN, LIMIT_SIZE_GENERIC, LIMIT_SIZE_GO, LIMIT_SIZE_HELM, LIMIT_SIZE_MAVEN, LIMIT_SIZE_NPM, LIMIT_SIZE_NUGET, LIMIT_SIZE_PUB, LIMIT_SIZE_PYPI, LIMIT_SIZE_RPM, LIMIT_SIZE_RUBYGEMS, LIMIT_SIZE_SWIFT, LIMIT_SIZE_VAGRANT), the system shall enforce count and size limits on package storage.`
- **PKG-01-204:** `Where a personal access token is issued with the public-only variant of the read:package scope, the system shall restrict token access to public packages only and reject access to private owner packages with 403 Forbidden.`

### Complex Requirements (Cleanup Rules)

- **PKG-01-901:** `When a cleanup rule is configured for an owner, and a package version matches the removal pattern, and the version is older than the configured retention days, and the version does not match the keep pattern, and the version is not among the N most recent versions, then the system shall remove the package version.`
- **PKG-01-902:** `When a cleanup rule is configured with match full name enabled, the system shall evaluate cleanup patterns against the full package path instead of the package name alone.`

### Unwanted Behaviour Requirements (Registry Errors)

- **PKG-01-301:** `If a client publishes a package version that already exists for the same owner and package name, then the system shall reject the publication with a conflict error.`
- **PKG-01-302:** `If a client attempts to access a package without the required read:package scope, then the system shall reject the request with 401 Unauthorized (including a WWW-Authenticate: Basic header). A 403 Forbidden is returned only when a token with a public-only scope variant attempts to access a private package.`
- **PKG-01-303:** `If the package storage backend is unreachable, then the system shall return a 503 Service Unavailable response.`

---

## 2. Container Registry (PKG-02)

**User Story:** As a container operator, I want to push and pull container images using standard Docker/OCI tooling so that my CI/CD pipelines work without modification.

### Ubiquitous Requirements (Container Architecture)

- **PKG-02-001:** `The system shall implement the OCI Distribution Spec (Docker v2 API) under the /v2/ endpoint prefix.`
- **PKG-02-002:** `The system shall scope container images to an owner (user or organization) via the path /v2/{owner}/{image}/.`
- **PKG-02-003:** `The system shall support multi-architecture images via manifest lists.`
- **PKG-02-004:** `The system shall track tags as references to specific manifests.`
- **PKG-02-005:** `The system shall expose a /v2/_catalog endpoint that returns the list of all container repositories (images) across owners, with pagination via the n and last query parameters.`

### Event-Driven Requirements (Container Workflow)

- **PKG-02-101:** `When a client pushes a blob via chunked upload, the system shall accept the upload in sequential chunks and assemble the complete blob on completion.`
- **PKG-02-102:** `When a client pushes a manifest, the system shall validate the manifest structure, store it, and associate any referenced blobs.`
- **PKG-02-103:** `When a client requests a manifest by tag, the system shall resolve the tag to the current manifest and return it.`
- **PKG-02-104:** `When a client requests a blob by digest, the system shall return the blob content.`
- **PKG-02-105:** `When a client deletes a tag, the system shall remove the tag reference without deleting the underlying manifest.`
- **PKG-02-106:** `When a client deletes a manifest by digest, the system shall remove the manifest and all associated tag references.`

### State-Driven Requirements (Container State)

- **PKG-02-701:** `While a blob upload is in progress, the system shall store partial data in a temporary chunked upload path ([packages] CHUNKED_UPLOAD_PATH).`
- **PKG-02-702:** `While a manifest references a blob, the system shall prevent garbage collection of that blob.`

### Unwanted Behaviour Requirements (Container Errors)

- **PKG-02-301:** `If a client pushes a manifest referencing a non-existent blob, then the system shall reject the manifest with a BLOB_UNKNOWN error.`
- **PKG-02-302:** `If a client submits an invalid manifest, then the system shall reject the push with a MANIFEST_INVALID error.`
- **PKG-02-303:** `If a chunked upload session expires, then the system shall discard the partial upload data and reject further chunk submissions.`

---

## 3. Package Type Coverage (PKG-03)

**User Story:** As a developer using various language ecosystems, I want to publish and install packages using my toolchain's native commands so that Gitea acts as a drop-in replacement for external registries.

### Ubiquitous Requirements (Language-Native API Compatibility)

- **PKG-03-001:** `The system shall expose each package type at a dedicated API path compatible with its native tooling.`
- **PKG-03-002:** `The system shall generate repository metadata files (indexes, manifests, listings) required by each package type's native client.`
- **PKG-03-003:** `The system shall support the following operations across all package types: upload (publish), download (install), delete, and list versions.`

### Event-Driven Requirements (Native Ecosystem Operations)

- **PKG-03-101:** `When a Cargo client publishes a crate, the system shall accept the crate archive and metadata, and support yank and unyank operations.`
- **PKG-03-102:** `When a Composer client requests packages.json, the system shall return a registry index compatible with the Composer v2 protocol.`
- **PKG-03-103:** `When a Conan client interacts with the registry, the system shall support both v1 and v2 Conan API endpoints for recipe and package management.`
- **PKG-03-104:** `When a Conda client requests a channel index, the system shall return repodata.json for the requested subdirectory and architecture.`
- **PKG-03-105:** `When a Debian client requests a distribution index, the system shall generate Packages, Release, and InRelease files for the requested distribution, component, and architecture.`
- **PKG-03-106:** `When a Maven client performs a PUT to the repository layout, the system shall store the artifact at the standard group/artifact/version path.`
- **PKG-03-107:** `When an npm client publishes a package, the system shall accept the tarball and metadata, and support scoped packages and dist-tags.`
- **PKG-03-108:** `When a NuGet client interacts with the registry, the system shall support both v2 and v3 API endpoints including symbol server functionality.`
- **PKG-03-109:** `When a PyPI client uploads a package, the system shall accept the distribution and serve a PEP 503-compliant simple index.`
- **PKG-03-110:** `When a RubyGems client pushes a gem, the system shall accept the gem file and support yank and specs listing operations.`
- **PKG-03-111:** `When a Go client requests a module version, the system shall return .info, .mod, and .zip files and support @latest resolution.`
- **PKG-03-112:** `When a Helm client requests index.yaml, the system shall return a valid Helm chart repository index.`
- **PKG-03-113:** `When an Alpine client requests an APKINDEX, the system shall generate the index for the requested branch, repository, and architecture.`
- **PKG-03-114:** `When an RPM client requests repodata, the system shall generate the repository metadata and provide a .repo configuration file.`
- **PKG-03-115:** `When a Pub client publishes a Dart package, the system shall accept the archive and support version finalization.`

### Optional Feature Requirements (Package Type Features)

- **PKG-03-201:** `Where a Generic package type is used, the system shall accept arbitrary files with custom properties at the path /{owner}/generic/{package}/{version}/.`
- **PKG-03-202:** `Where a Chef cookbook is published, the system shall expose universe and search endpoints at /{owner}/chef/api/v1/.`
- **PKG-03-203:** `Where a CRAN package is published, the system shall serve both source and binary package formats.`
- **PKG-03-204:** `Where a Swift package is published, the system shall serve Package.swift downloads with scope/name versioning.`
- **PKG-03-205:** `Where a Vagrant box is uploaded, the system shall support provider-specific box formats.`

### Unwanted Behaviour Requirements (Package Type Errors)

- **PKG-03-301:** `If a client submits a package with an invalid format for the declared type, then the system shall reject the upload with a validation error.`
- **PKG-03-302:** `If a native tool requests a package metadata endpoint for a non-existent package, then the system shall return a 404 response in the format expected by that tool.`

---

## 4. Releases (PKG-04)

**User Story:** As a project maintainer, I want to create releases from Git tags so that users can download stable versions of my software with associated assets and release notes.

### Ubiquitous Requirements (Release Properties)

- **PKG-04-001:** `The system shall associate each release with a Git tag in a repository.`
- **PKG-04-002:** `The system shall support three release states: Draft (unpublished), Pre-release, and Release (official).`
- **PKG-04-003:** `The system shall store title, body (markdown), and target branch reference per release.`
- **PKG-04-004:** `The system shall track the number of commits included in the release.`
- **PKG-04-005:** `The system shall support standalone tags without release metadata (IsTag).`

### Event-Driven Requirements (Release Workflow)

- **PKG-04-101:** `When a maintainer creates a release, the system shall create the release from an existing tag or create a new tag at the specified commit.`
- **PKG-04-102:** `When a maintainer uploads a release asset, the system shall store the file and generate a unique download URL.`
- **PKG-04-103:** `When a user downloads a release asset, the system shall increment the download counter.`
- **PKG-04-104:** `When a maintainer deletes a release, the system shall remove the release metadata and all associated assets.`
- **PKG-04-105:** `When a maintainer updates a release to change its state (draft to release), the system shall persist the state change immediately.`
- **PKG-04-106:** `When a user requests the latest release, the system shall return the most recent non-draft, non-pre-release release.`
- **PKG-04-107:** `When a user requests a source archive at /{owner}/{repo}/archive/{tag}.zip or .tar.gz, the system shall generate and serve the archive on demand.`

### State-Driven Requirements (Release Visibility)

- **PKG-04-701:** `While a release is in draft state, the system shall not list it in public release feeds or API listings unless the requester has write access.`
- **PKG-04-702:** `While a release is marked as pre-release, the system shall indicate the pre-release status in listings but include it in results.`

### Unwanted Behaviour Requirements (Release Errors)

- **PKG-04-301:** `If a maintainer creates a release for a tag that already has a release, then the system shall reject the creation with a conflict error.`
- **PKG-04-302:** `If a non-collaborator attempts to create a release, then the system shall reject the request with 403 Forbidden.`
- **PKG-04-303:** `If a user uploads an asset exceeding the configured size limit, then the system shall reject the upload.`

---

## 5. LFS — Large File Storage (PKG-05)

**User Story:** As a developer working with large binary files, I want to use Git LFS so that my repository remains performant while still versioning large assets.

### Ubiquitous Requirements (LFS Architecture)

- **PKG-05-001:** `The system shall store LFS objects using SHA256-based content addressing and deduplication.`
- **PKG-05-002:** `The system shall link LFS objects to repositories with reference tracking.`
- **PKG-05-003:** `The system shall track LFS storage size per repository.`
- **PKG-05-004:** `The system shall implement the Git LFS HTTP batch API under /{owner}/{repo}.git/info/lfs/.`

### Event-Driven Requirements (LFS Workflow)

- **PKG-05-101:** `When a client pushes LFS objects via the batch API, the system shall accept the upload and store the object keyed by its SHA256 OID.`
- **PKG-05-102:** `When a client requests LFS objects via the batch API, the system shall return download URLs for each requested object.`
- **PKG-05-103:** `When a user pushes commits referencing LFS objects, the system shall automatically associate the LFS objects with the target repository.`
- **PKG-05-104:** `When a user creates an LFS lock on a file, the system shall record the lock with the owner and file path.`
- **PKG-05-105:** `When a user requests to unlock an LFS lock, the system shall remove the lock, allowing other users to modify the file.`

### State-Driven Requirements (LFS Locking)

- **PKG-05-701:** `While a file has an active LFS lock, the system shall track the lock owner and prevent other users from pushing changes to that file.`
- **PKG-05-702:** `While LFS is enabled ([lfs] START_SERVER = true), the system shall accept and process all LFS API requests.`

### Optional Feature Requirements (LFS Configuration)

- **PKG-05-201:** `Where [lfs] MAX_FILE_SIZE is configured, the system shall reject LFS uploads exceeding the configured limit.`
- **PKG-05-202:** `Where MinIO/S3 storage is configured for LFS, the system shall store LFS objects in the configured bucket instead of local filesystem.`
- **PKG-05-203:** `Where [server] LFS_HTTP_AUTH_EXPIRY is configured (default 24h), the system shall issue LFS batch API download/upload URLs with authentication tokens that expire after the configured duration.`

### Unwanted Behaviour Requirements (LFS Errors)

- **PKG-05-301:** `If a client uploads an LFS object whose SHA256 does not match the declared OID, then the system shall reject the upload.`
- **PKG-05-302:** `If a user attempts to lock a file that is already locked by another user, then the system shall reject the lock creation.`
- **PKG-05-303:** `If LFS storage is disabled, then the system shall return a 404 for all LFS endpoints.`

---

## 6. Attachments (PKG-06)

**User Story:** As a user, I want to attach files to issues, comments, and releases so that I can share screenshots, logs, and other resources in context.

### Ubiquitous Requirements (Attachment Properties)

- **PKG-06-001:** `The system shall use UUID-based unique identification for all stored attachments.`
- **PKG-06-002:** `The system shall support attachments in three contexts: issue descriptions, comment bodies, and release assets.`
- **PKG-06-003:** `The system shall track download counters per attachment.`
- **PKG-06-004:** `The system shall generate image previews for supported image formats.`
- **PKG-06-005:** `The system shall enforce configurable MIME type allowlists for upload ([attachment] ALLOWED_TYPES).`

### Event-Driven Requirements (Attachment Workflow)

- **PKG-06-101:** `When a user uploads a file to an issue or comment, the system shall store the file and return a UUID reference for embedding.`
- **PKG-06-102:** `When a user requests an attachment by UUID, the system shall serve the file and increment the download counter.`
- **PKG-06-103:** `When a user deletes an issue, the system shall remove all attachments associated with that issue and its comments.`
- **PKG-06-104:** `When a user deletes an individual attachment, the system shall remove the file and its metadata.`

### Optional Feature Requirements (Attachment Configuration)

- **PKG-06-201:** `Where [attachment] MAX_SIZE is configured, the system shall reject file uploads exceeding the specified limit (default 2MB / 2048 KB).`
- **PKG-06-202:** `Where [attachment] MAX_FILES is configured, the system shall limit the number of files per upload batch (default 5).`
- **PKG-06-203:** `Where [attachment] ENABLED is false, the system shall disable all attachment upload functionality.`

### Unwanted Behaviour Requirements (Attachment Errors)

- **PKG-06-301:** `If a user uploads a file with a disallowed MIME type, then the system shall reject the upload with a validation error.`
- **PKG-06-302:** `If a user uploads more files than the configured MAX_FILES limit, then the system shall reject the excess files.`
- **PKG-06-303:** `If a non-authenticated user requests a private attachment, then the system shall deny access with 404.`

---

## Storage Backend Configuration (moved to Domain 08)

Storage backend configuration (local, MinIO) is documented in Domain 08 (Administration & Operations). This domain documents storage usage by packages, LFS, and attachments — i.e., how each feature reads from and writes to the configured backend. See Domain 08 (Administration & Operations) for `[storage]`, `[storage.packages]`, `[storage.lfs]`, `[storage.attachments]`, and related backend configuration.

---

## Configuration Reference

### `[packages]`

| Key | Description |
|-----|-------------|
| `ENABLED` | Enable/disable the package registry (default true) |
| `CHUNKED_UPLOAD_PATH` | Temporary path for chunked container blob uploads |
| `SERVE_DIRECT` | Generate direct download URLs bypassing Gitea server (MinIO only) |
| `LIMIT_TOTAL_OWNER_COUNT` | Per-owner package version count quota |
| `LIMIT_TOTAL_OWNER_SIZE` | Per-owner total storage size quota |
| `LIMIT_SIZE_ALPINE` | Max size per Alpine package blob |
| `LIMIT_SIZE_CARGO` | Max size per Cargo package blob |
| `LIMIT_SIZE_CHEF` | Max size per Chef package blob |
| `LIMIT_SIZE_COMPOSER` | Max size per Composer package blob |
| `LIMIT_SIZE_CONAN` | Max size per Conan package blob |
| `LIMIT_SIZE_CONDA` | Max size per Conda package blob |
| `LIMIT_SIZE_CONTAINER` | Max size per Container blob |
| `LIMIT_SIZE_CRAN` | Max size per CRAN package blob |
| `LIMIT_SIZE_DEBIAN` | Max size per Debian package blob |
| `LIMIT_SIZE_GENERIC` | Max size per Generic package blob |
| `LIMIT_SIZE_GO` | Max size per Go module blob |
| `LIMIT_SIZE_HELM` | Max size per Helm chart blob |
| `LIMIT_SIZE_MAVEN` | Max size per Maven artifact blob |
| `LIMIT_SIZE_NPM` | Max size per npm package blob |
| `LIMIT_SIZE_NUGET` | Max size per NuGet package blob |
| `LIMIT_SIZE_PUB` | Max size per Pub package blob |
| `LIMIT_SIZE_PYPI` | Max size per PyPI package blob |
| `LIMIT_SIZE_RPM` | Max size per RPM package blob |
| `LIMIT_SIZE_RUBYGEMS` | Max size per RubyGems blob |
| `LIMIT_SIZE_SWIFT` | Max size per Swift package blob |
| `LIMIT_SIZE_VAGRANT` | Max size per Vagrant box blob |

### `[lfs]`

| Key | Description |
|-----|-------------|
| `START_SERVER` | Enable/disable the LFS HTTP server (default false) |
| `MAX_FILE_SIZE` | Max size per LFS object upload |

### `[attachment]`

| Key | Description |
|-----|-------------|
| `ENABLED` | Enable/disable attachment uploads (default true) |
| `ALLOWED_TYPES` | Comma-separated MIME type allowlist (default see app.ini) |
| `MAX_SIZE` | Max file size per upload (default 2MB / 2048 KB) |
| `MAX_FILES` | Max number of files per upload batch (default 5) |

### `[server]` (LFS-related)

| Key | Description |
|-----|-------------|
| `LFS_HTTP_AUTH_EXPIRY` | Expiry duration for LFS batch API auth tokens (default 24h) |

### Storage configuration → Domain 08

`[storage]`, `[storage.packages]`, `[storage.lfs]`, `[storage.attachments]`, and all storage backend parameters (local path, MinIO endpoint/access key/secret key/bucket/region) are documented in Domain 08 (Administration & Operations).

---

## Business Rules

- **BR-06-001:** Package content is deduplicated globally using SHA256 content addressing
- **BR-06-002:** Packages are scoped to owners; two different owners can have packages with the same name
- **BR-06-003:** Cleanup rules are evaluated per-owner and support keep-count, keep-pattern, remove-days, and remove-pattern criteria
- **BR-06-004:** Container registry implements OCI Distribution Spec; Docker CLI and compatible tools work without modification
- **BR-06-005:** Release states are mutually exclusive: a release is exactly one of Draft, Pre-release, or Release
- **BR-06-006:** Source archives (.zip, .tar.gz) are generated on demand and not stored permanently
- **BR-06-007:** LFS objects are shared across repositories via content addressing; a single object is stored once
- **BR-06-008:** LFS locks are per-repository and per-path; only one lock per path per repository is allowed
- **BR-06-009:** Attachment UUIDs are globally unique and serve as the primary download identifier
- **BR-06-010:** Default attachment size limit is 2MB (2048 KB); default file count limit per upload is 5
- **BR-06-013:** Package authentication requires either a personal access token with read:package/write:package scopes or basic auth credentials

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Duplicate package version upload | Reject with conflict error; existing version is preserved |
| Container blob upload interrupted | Discard partial upload; client must restart |
| Cleanup rule removes all versions | Remove the package entirely if no versions remain |
| Release tag deleted externally | Release becomes orphaned; tag reference is stale |
| LFS object missing during clone | Return error indicating the object is unavailable |
| LFS lock owner no longer has access | Allow force unlock by repo admin |
| Attachment MIME type not in allowlist | Reject upload with validation error listing allowed types |
| Storage backend becomes unavailable mid-operation | Return 503; do not persist partial data |
| Direct URL generation fails | Fall back to proxied serving through Gitea server |
| Package quota exceeded | Reject upload with quota exceeded error |
| NuGet symbol package without main package | Reject upload; symbol package must reference existing package |
| Maven PUT to existing snapshot | Accept and overwrite the snapshot version |
| npm scoped package without scope prefix | Reject; scoped packages must include the @scope prefix |

## Success Criteria

- Package publish and install operations complete within 5 seconds for artifacts under 100MB
- Container registry supports Docker CLI push and pull without modification
- SHA256 deduplication reduces storage usage for identical content across packages
- Cleanup rules execute within configurable intervals without impacting registry availability
- Release creation from existing tags completes within 2 seconds
- Source archive generation for repositories under 1GB completes within 10 seconds
- LFS batch API responds within 500ms for repositories with up to 1000 objects
- Attachment upload and download complete within 3 seconds for files under the size limit
- Storage backend failover (direct URL to proxy) occurs transparently without client errors
- All 21 package type APIs remain compatible with their respective native tooling

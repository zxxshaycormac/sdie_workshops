# 06 — Packages & Artifacts

## 1. Package Registry

### What
Comprehensive multi-format package registry supporting 23 package types. Each package type has dedicated API endpoints compatible with its native tooling.

### Supported Package Types

#### Alpine (apk)
- Format: APK for Alpine Linux
- API: /{owner}/alpine/{branch}/{repository}/{arch}/...
- Features: branch/repository/architecture structure, APKINDEX generation, repository key

#### Cargo (Rust)
- Format: Cargo crate
- API: /{owner}/cargo/api/v1/crates/...
- Features: yank/unyank, search, owner management

#### Chef
- Format: Cookbook
- API: /{owner}/chef/api/v1/...
- Features: universe, search, cookbook lifecycle

#### Composer (PHP)
- Format: PHP package
- API: /{owner}/composer/...
- Features: packages.json, search, p2 metadata

#### Conan (C/C++)
- Format: Conan recipe/package
- API: /{owner}/conan/v1/... and /{owner}/conan/v2/...
- Features: v1 and v2 API, recipe and package management

#### Conda (Python)
- Format: Anaconda package
- API: /{owner}/conda/{channel}/{arch}/...
- Features: multi-channel, binary and source packages

#### Container (OCI/Docker)
- Format: OCI Distribution Spec
- API: /v2/{owner}/{image}/...
- Features: multi-arch, blob upload (chunked), manifest management, tag listing

#### CRAN (R)
- Format: R package
- API: /{owner}/cran/src/contrib/... and /{owner}/cran/bin/...
- Features: source and binary packages

#### Debian (apt)
- Format: DEB package
- API: /{owner}/debian/dists/... and /{owner}/debian/pool/...
- Features: distribution/component/architecture, repository key, index generation

#### Generic
- Format: Arbitrary files
- API: /{owner}/generic/{package}/{version}/...
- Features: flexible storage, custom properties

#### Go
- Format: Go module proxy
- API: /{owner}/go/{name}/@v/...
- Features: version listing, .info/.mod/.zip, @latest

#### Helm
- Format: Helm chart
- API: /{owner}/helm/...
- Features: index.yaml, chart upload/download

#### Maven (Java)
- Format: Maven repository
- API: /{owner}/maven/{group}/{artifact}/{version}/...
- Features: standard Maven layout, PUT/GET/HEAD

#### npm (JavaScript)
- Format: npm registry
- API: /{owner}/npm/...
- Features: scoped packages, dist-tags, search

#### NuGet (.NET)
- Format: NuGet package
- API: /index.json (v3), / (v2)
- Features: v2 and v3 API, symbol server

#### Pub (Dart)
- Format: Dart package
- API: /{owner}/pub/api/packages/...
- Features: version publishing, finalization

#### PyPI (Python)
- Format: Python package
- API: /{owner}/pypi/...
- Features: standard PyPI upload, simple index

#### RPM
- Format: RPM package
- API: /{owner}/rpm/...
- Features: group support, repodata, .repo config

#### RubyGems
- Format: Ruby gem
- API: /{owner}/rubygems/...
- Features: specs listing, yank support, pre-release

#### Swift
- Format: Swift package
- API: /{owner}/swift/...
- Features: Package.swift download, scope/name versioning

#### Vagrant
- Format: Vagrant box
- API: /{owner}/vagrant/...
- Features: provider-specific boxes

### Storage Architecture
- Content-addressable: SHA256 deduplication
- Path structure: packages/aa/bb/SHA256
- Direct serving: MinIO/S3 direct URL support
- Quotas: per-user and per-package-type size limits

### Authentication
- JWT token-based
- Scopes: read:package, write:package
- Basic auth support
- OAuth2 support

### Cleanup Rules
- Keep count: preserve N most recent versions
- Keep pattern: regex to preserve matching versions
- Remove days: delete versions older than X days
- Remove pattern: regex to delete matching versions
- Match full name: match against package name or full path
- Per-owner cleanup rule configuration

### UI Routes
- GET /{owner}/-/packages — Owner's packages
- GET /{owner}/-/packages/{type}/{name} — Package detail
- GET /{owner}/-/packages/{type}/{name}/{version} — Version detail
- GET /{owner}/-/packages/{type}/{name}/{version}/settings — Version settings
- GET /admin/packages — Admin package management (admin)

### Config
- [packages] ENABLED — enable registry
- [packages] STORAGE_TYPE — local, minio, azure
- [packages] MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET
- [packages] CHUNKED_UPLOAD_PATH — temp path for chunked uploads

---

## 2. Releases

### What
Git tag-based release management with assets, metadata, and lifecycle states.

### Behaviors
- Create from existing tags or create new tags
- States: Draft (unpublished), Pre-release, Release (official)
- Title, body (markdown), and tag reference
- IsTag: standalone tags without release metadata
- NumCommits tracking
- Target branch reference

### Release Assets
- Upload multiple files per release
- Download with auto-incrementing counter
- Custom download URLs for unique filenames
- Size tracking

### UI Routes
- GET /{owner}/{repo}/releases — Release list
- GET /{owner}/{repo}/releases/tag/{tag} — View release
- GET /{owner}/{repo}/releases/new — Create release
- POST /{owner}/{repo}/releases — Create release
- POST /{owner}/{repo}/releases/{id} — Update release
- GET /{owner}/-/releases — User's releases
- GET /org/{org}/-/releases — Org releases

### API Endpoints
- GET /repos/{owner}/{repo}/releases — List releases
- GET /repos/{owner}/{repo}/releases/{id} — Get release
- POST /repos/{owner}/{repo}/releases — Create release
- PATCH /repos/{owner}/{repo}/releases/{id} — Update release
- DELETE /repos/{owner}/{repo}/releases/{id} — Delete release
- GET /repos/{owner}/{repo}/releases/{id}/assets — List assets
- POST /repos/{owner}/{repo}/releases/{id}/assets — Upload asset
- GET /repos/{owner}/{repo}/releases/assets/{id} — Download asset
- DELETE /repos/{owner}/{repo}/releases/assets/{id} — Delete asset
- GET /repos/{owner}/{repo}/releases/tags/{tag} — Get by tag
- GET /repos/{owner}/{repo}/releases/latest — Latest release

### Download URLs
- /{owner}/{repo}/releases/download/{tag}/{filename} — Asset download
- /{owner}/{repo}/archive/{tag}.zip — Source archive (zip)
- /{owner}/{repo}/archive/{tag}.tar.gz — Source archive (tar.gz)

---

## 3. LFS (Large File Storage)

### What
Git LFS integration for managing large binary files outside the main Git repository.

### Behaviors
- Configurable storage: local, MinIO/S3, Azure Blob
- SHA256-based content addressing and deduplication
- Object linking to repositories
- Size tracking per repository
- Automatic association during push
- Migration between repositories

### LFS Locks
- File locking for collaborative editing
- Lock owner tracking
- Force unlock capability
- List locks by repository/path

### LFS HTTP API
- POST /{owner}/{repo}.git/info/lfs/objects/batch — Batch operations
- GET /{owner}/{repo}.git/info/lfs/objects/{oid} — Get object
- POST /{owner}/{repo}.git/info/lfs/objects — Upload object
- GET /{owner}/{repo}.git/info/lfs/locks — List locks
- POST /{owner}/{repo}.git/info/lfs/locks — Create lock
- POST /{owner}/{repo}.git/info/lfs/locks/{id}/unlock — Delete lock

### Config
- [lfs] START_SERVER — enable LFS (default true)
- [lfs] STORAGE_TYPE — backend type
- [lfs] PATH — local storage path
- [lfs] MINIO_ENDPOINT, MINIO_BUCKET — S3 config
- [lfs] MAX_FILE_SIZE — max LFS file size

---

## 4. Attachments

### What
File attachment system for issues, comments, and releases. UUID-based storage with multiple backends.

### Behaviors
- Issue attachments: upload to issue description and comments
- Comment attachments: files in discussion comments
- Release assets: files attached to releases
- UUID-based unique identification
- Download counter tracking
- Image preview optimization
- Delete by issue or individual file

### API Endpoints
- GET /attachments/{uuid} — Download attachment
- POST /repos/{owner}/{repo}/issues/{index}/assets — Upload issue attachment
- POST /repos/{owner}/{repo}/issues/comments/{id}/assets — Upload comment attachment
- DELETE /repos/{owner}/{repo}/issues/comments/{id}/assets/{asset_id} — Delete attachment

### Config
- [attachment] ENABLED — enable attachments
- [attachment] STORAGE_TYPE — backend
- [attachment] PATH — local path
- [attachment] ALLOWED_TYPES — MIME type allowlist (e.g., image/*,application/pdf)
- [attachment] MAX_SIZE — max file size (default 4MB)
- [attachment] MAX_FILES — max files per upload (default 5)

---

## 5. Storage Backends

### What
Pluggable storage system used by attachments, LFS, packages, avatars, actions artifacts, and repo archives.

### Supported Backends
1. **Local**: filesystem storage, path-based
2. **MinIO/S3**: AWS S3 compatible, bucket-based, direct URL support
3. **Azure Blob**: Microsoft Azure, container-based, direct URL access

### Storage Types (each independently configurable)
- attachments — issue/comment attachments
- lfs — LFS objects
- packages — package registry files
- avatars — user avatars
- repo-avatars — repository avatars
- repo-archive — git archive downloads
- actions_log — Actions log files
- actions_artifact — Actions artifact files

### Config (per storage type)
- [storage.{type}] STORAGE_TYPE — local/minio/azure
- [storage.{type}] PATH — local path
- [storage.{type}] MINIO_ENDPOINT — S3 endpoint
- [storage.{type}] MINIO_ACCESS_KEY_ID
- [storage.{type}] MINIO_SECRET_ACCESS_KEY
- [storage.{type}] MINIO_BUCKET
- [storage.{type}] MINIO_LOCATION
- [storage.{type}] MINIO_USE_SSL
- [storage.{type}] MINIO_INSECURE_SKIP_VERIFY
- [storage.{type}] AZURE_ENDPOINT — Azure endpoint
- [storage.{type}] AZURE_ACCOUNT_NAME
- [storage.{type}] AZURE_ACCOUNT_KEY
- [storage.{type}] AZURE_CONTAINER
- [storage.{type}] SERVE_DIRECT — direct URL generation

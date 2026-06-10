# 07 — Administration & Ops

## 1. Admin Panel

### What
Central admin dashboard and management interface for system-wide operations: users, repos, orgs, auth sources, emails, hooks, queues, packages, runners.

### Dashboard
- System statistics: memory, goroutines, GC stats
- Update checker status
- Quick operations: sync branches/tags, run cron tasks
- System health warnings
- Self-check diagnostics

### User Management
- List/search/filter users (active, admin, restricted, 2FA, prohibited login)
- Create/edit/delete users
- Password management with complexity checks
- Pwned password detection
- 2FA reset
- Avatar management
- Email management

### Repository Management
- List/search repos
- Delete repos
- Unadopted repo detection and adoption from filesystem
- Bulk operations

### Organization Management
- List/search orgs
- Filter by status

### Auth Source Management
- Supported types: LDAP (BindDN), LDAP (simple), SMTP, OAuth2, PAM, SSPI, FreeIPA
- Group-to-team mapping
- Admin/restricted group filtering
- Skip local 2FA option
- TLS configuration
- User sync (LDAP)

### Email Management
- List/search emails
- Activate/deactivate emails
- Set primary email

### Hook Management
- System webhook administration
- Default webhook administration

### Queue Management
- View queue status (workers, items)
- Adjust worker count
- Clear queue items

### Package Management
- List/search packages
- Delete package versions
- Cleanup expired data

### Runner Management
- Actions runner administration

### UI Routes
- GET /admin — Dashboard
- GET /admin/system_status — System status
- GET /admin/self_check — Self check
- GET /admin/users — User list
- GET /admin/users/new — Create user
- GET /admin/users/:userid/edit — Edit user
- GET /admin/repos — Repo list
- GET /admin/repos/unadopted — Unadopted repos
- GET /admin/orgs — Org list
- GET /admin/auths — Auth sources
- GET /admin/auths/new — New auth source
- GET /admin/auths/:authid — Edit auth source
- GET /admin/emails — Email list
- GET /admin/hooks — System hooks
- GET /admin/monitor/queues — Queue monitor
- GET /admin/monitor/queue/:qid — Queue detail
- GET /admin/packages — Package list
- GET /admin/actions/runners — Runner management

### API Endpoints
- Full CRUD in /api/v1/admin/* for users, orgs, repos, auth sources, runners

### Config Viewer
- GET /admin/config — Configuration summary (passwords shadowed)
- GET /admin/config/settings — Dynamic settings
- POST /admin/config/send_test_mail — Test email
- POST /admin/config/change — Update settings

---

## 2. Configuration System

### What
INI-based configuration with environment variable overrides. Loaded at startup, some settings dynamic.

### Major Sections
- [server]: protocol, domain, port, root URL, SSH, run user
- [database]: type (MySQL/PostgreSQL/SQLite3/MSSQL), host, connection pool
- [repository]: root path, default branch, signing, merge options, wiki, issues
- [ui]: themes, language, time format, gravatar, paging
- [indexer]: issue/repo indexer settings, backends
- [ssh]: port, key path, builtin server
- [lfs]: server, storage
- [packages]: storage, allowed types, size limits
- [actions]: log/artifact storage, retention, timeout, default actions URL
- [webhook]: task queue, delivery timeout, retry
- [mailer]: SMTP, sender, TLS
- [cache]: adapter (memory/redis/memcache/twoqueue), TTL
- [session]: provider, cookie, timeout
- [log]: modes (console/file/conn/smtp), levels, paths
- [markup]: markdown, custom renderers
- [cron]: background task scheduling
- [mirror]: update intervals, git timeout
- [api]: swagger, response limits, paging
- [oauth2]: provider settings, token expiration
- [security]: install lock, secret key, password policy
- [admin]: org creation, email notifications, disabled features
- [attachment]: storage, allowed types, size limits
- [federation]: ActivityPub, headers, digest
- [camo]: URL rewriting for images
- [other]: timeout settings

### Config File Locations
- Primary: custom/conf/app.ini
- Environment variable overrides: GITEA__section__key format
- Install wizard: LoadSettingsForInstall()

---

## 3. Queue System

### What
Concurrent task processing with multiple backends and worker pool management.

### Queue Types
- Simple Queue: basic FIFO
- Unique Queue: prevents duplicate items
- Worker Pool Queue: manages multiple workers

### Backends
- Channel: in-memory (default, single instance)
- LevelDB: persistent (single instance)
- Redis: distributed (cluster deployments)
- Dummy: immediate no-op (testing)

### Features
- Dynamic worker count adjustment
- Queue flushing and item removal
- Batch processing
- Error handling and retry
- Graceful shutdown

### Config
- [queue] TYPE — backend type
- [queue] DATADIR — LevelDB path
- [queue] CONN_STR — Redis connection
- [queue] LENGTH — queue length
- [queue] BATCH_LENGTH — batch size
- [queue] WORKERS — worker count

---

## 4. Metrics/Monitoring

### What
Prometheus-compatible metrics endpoint for monitoring Gitea instances.

### Endpoint
- GET /api/v1/metrics (optional bearer token auth)

### Metrics Exposed
- Build info (goarch, goos, goversion, version)
- Access count, attachment count, comment count
- Follow count, hook task count
- Issue counts (total, open, closed, by label, by repo)
- Label count, login source count, milestone count
- Mirror count, OAuth count, organization count
- Project counts, public key count, release count
- Repository count, star count, team count
- Update task count, user count, watch count, webhook count

### Config
- [metrics] ENABLED — enable metrics endpoint
- [metrics] TOKEN — bearer token (optional)
- [metrics] ENABLED_ISSUE_BY_LABEL — per-label issue metrics
- [metrics] ENABLED_ISSUE_BY_REPOSITORY — per-repo issue metrics

---

## 5. Health Checks

### What
Application health monitoring endpoint following RFC standard.

### Endpoint
- GET /-/healthz

### Check Types
- Database ping
- Cache ping

### Status Levels
- pass: healthy (2xx-3xx)
- fail: unhealthy (4xx-5xx)
- warn: healthy with concerns (2xx-3xx)

### Response Format
```json
{
  "status": "pass|fail|warn",
  "description": "Gitea",
  "checks": {
    "check_name": [{"status": "...", "time": "...", "output": "..."}]
  }
}
```

---

## 6. Database Migrations

### What
Version-controlled database schema evolution with automatic migration execution.

### Behaviors
- Version table tracks current schema version
- Sequential version numbering (min version: 70)
- Each migration has unique description
- Migrations run in order on startup
- Organized by release version directories (v1_6 through v1_23)

### Migration Interface
- Description() string
- Migrate(x *xorm.Engine) error

---

## 7. Backup/Restore (Dump)

### What
System backup via `dump` command. Creates archive of all Gitea data.

### Output Formats
- zip, tar, tar.gz, tar.xz, tar.bz2, tar.br, tar.lz4, tar.zst

### Backup Contents
- Database dump
- Repository files
- Configuration (app.ini)
- Attachments
- Avatars
- LFS objects
- Package files
- Log files

### Config
- --tempdir — temporary directory
- --skip-repository — skip repo data
- --skip-lfs — skip LFS objects
- --skip-attachment — skip attachments
- --skip-package — skip packages
- --skip-log — skip logs
- --type — output format

---

## 8. System Notices

### What
System-wide notification tracking for administrative events.

### Notice Types
- NoticeRepository: repo-related events
- NoticeTask: task-related events

### Features
- Create, list (with pagination), delete notices
- Delete by range or age
- Automatic cleanup

### UI Routes
- GET /admin/notices — List notices
- POST /admin/notices/delete — Delete notices

---

## 9. Application State

### What
Key-value state persistence for application coordination.

### Behaviors
- Atomic updates with revision tracking
- Context-based operations
- Used for coordination between instances

---

## 10. Logging

### What
Comprehensive logging infrastructure with multiple output modes.

### Log Levels
TRACE, DEBUG, INFO, WARN, ERROR, FATAL

### Output Modes
- Console: terminal with colors
- File: rotating file output
- Conn: network logging
- SMTP: email-based log alerts

### Features
- Multiple log channels
- Event formatting
- Process tracing integration
- Rotating file writer
- Color support

### Config
- [log] MODE — output modes (comma-separated)
- [log] LEVEL — default log level
- [log.*] — per-mode configuration

---

## 11. Process Management

### What
Process lifecycle management for Gitea's internal subprocess tracking.

### Features
- PID tracking and context management
- Parent-child process relationships
- Stack trace collection
- Graceful shutdown support
- PProf integration

---

## 12. Debug/PProf

### What
Built-in Go pprof endpoints for performance profiling.

### Endpoints
- /debug/pprof/ — Index
- /debug/pprof/heap — Heap profiling
- /debug/pprof/goroutine — Goroutine dump
- /debug/pprof/threadcreate — Thread creation
- /debug/pprof/block — Block profiling
- /debug/pprof/mutex — Mutex profiling
- /debug/pprof/cmdline — Command line
- /debug/pprof/profile — CPU profile
- /debug/pprof/symbol — Symbol lookup
- /debug/pprof/trace — Execution trace

### Constraints
- Only available when PprofEnabled config is set
- Should be disabled in production for security

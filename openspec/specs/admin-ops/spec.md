# 07 — Administration & Ops

Baseline specification of Gitea's administration dashboard, configuration, observability, backup, and debug subsystems. All requirements describe the current (v1.22.x) system behavior.

---

## 1. Admin Dashboard

**User Story:** As an administrator, I want a central dashboard so that I can monitor system health and perform quick operational tasks.

### Ubiquitous Requirements (Dashboard Properties)

- **ADM-01-001:** `The system shall display system statistics on the admin dashboard including memory usage, goroutine count, and GC stats.`
- **ADM-01-002:** `The system shall present an update checker status indicating whether a newer Gitea version is available.`
- **ADM-01-003:** `The system shall list quick operations on the dashboard including sync branches, sync tags, and run cron tasks.`

### Event-Driven Requirements (Dashboard Workflow)

- **ADM-01-101:** `When an admin navigates to /admin, the system shall render the dashboard with current system statistics.`
- **ADM-01-102:** `When an admin triggers a quick operation from the dashboard, the system shall execute the operation and display the result.`
- **ADM-01-103:** `When an admin navigates to /admin/self_check, the system shall run self-check diagnostics and display any warnings.`
- **ADM-01-104:** `When an admin navigates to /admin/system_status, the system shall display detailed system status information.`

### Unwanted Behaviour Requirements (Dashboard Errors)

- **ADM-01-301:** `If a non-admin user attempts to access /admin routes, then the system shall return a 403 Forbidden response.`

---

## 2. User Management

**User Story:** As an administrator, I want to manage user accounts so that I can create, modify, and remove user access as needed.

### Ubiquitous Requirements (User List Properties)

- **ADM-02-001:** `The system shall list all user accounts with pagination support.`
- **ADM-02-002:** `The system shall support filtering the user list by status: active, admin, restricted, 2FA enabled, and prohibited login.`
- **ADM-02-003:** `The system shall support searching users by username or email address.`

### Event-Driven Requirements (User CRUD Workflow)

- **ADM-02-101:** `When an admin creates a new user via the admin panel, the system shall create the account in the active state.`
- **ADM-02-102:** `When an admin edits a user account, the system shall persist the changes immediately.`
- **ADM-02-103:** `When an admin deletes a user account, the system shall transfer the user's repository ownership to a ghost user for attribution preservation.`
- **ADM-02-104:** `When an admin sets a new password for a user, the system shall hash and store the password using the configured algorithm.`

### Event-Driven Requirements (Password & 2FA Management)

- **ADM-02-105:** `When an admin forces a password change for a user, the system shall set the must_change_password flag so the user must change their password at next login.`
- **ADM-02-106:** `When an admin resets a user's 2FA, the system shall disable 2FA for that user and regenerate scratch codes.`
- **ADM-02-107:** `When a new password is set and pwned password checking is enabled, the system shall check the password against known breach databases.`

### Optional Feature Requirements (Avatar Management)

- **ADM-02-201:** `Where a user has a custom avatar, the system shall display the custom avatar instead of the Gravatar default.`

### Unwanted Behaviour Requirements (User Management Errors)

- **ADM-02-301:** `If an admin attempts to delete the last admin user, then the system shall deny the deletion.`
- **ADM-02-302:** `If an admin submits a password that does not meet complexity requirements, then the system shall reject the password and display the requirements.`

---

## 3. Repository Management

**User Story:** As an administrator, I want to manage repositories across the instance so that I can delete repos and adopt unadopted ones from the filesystem.

### Ubiquitous Requirements (Repository List Properties)

- **ADM-03-001:** `The system shall list all repositories with pagination support.`
- **ADM-03-002:** `The system shall support searching repositories by name.`

### Event-Driven Requirements (Repository Workflow)

- **ADM-03-101:** `When an admin deletes a repository, the system shall remove the repository data and all associated resources.`
- **ADM-03-102:** `When an admin navigates to /admin/repos/unadopted, the system shall scan the repository root path for directories not registered in the database.`
- **ADM-03-103:** `When an admin adopts an unadopted repository, the system shall register it in the database and make it accessible.`
- **ADM-03-104:** `When an admin deletes an unadopted repository, the system shall remove the directory from the filesystem.`

### Unwanted Behaviour Requirements (Repository Errors)

- **ADM-03-301:** `If an unadopted repository directory contains invalid data, then the system shall skip the directory during the unadopted scan.`

---

## 4. Auth Source Management

**User Story:** As an administrator, I want to configure external authentication sources so that users can authenticate against LDAP, SMTP, OAuth2, and other providers.

### Ubiquitous Requirements (Admin Interface)

- **ADM-04-001:** `The system shall expose auth source CRUD operations in the admin panel (see AUTH-06 for supported source types, lifecycle behavior, optional features, and error handling).`
- **ADM-04-002:** `The system shall store authentication source configurations with activation state and display order.`
- **ADM-04-003:** `The system shall allow TLS configuration per authentication source.`

### Event-Driven Requirements (Admin Operations)

- **ADM-04-101:** `When an admin triggers an LDAP sync from the admin panel, the system shall synchronize user data from the LDAP directory to local accounts.`
- **ADM-04-102:** `When an admin reorders authentication sources, the system shall update the display order and login precedence.`

---

## 5. Configuration System

**User Story:** As an administrator, I want to view and modify Gitea's configuration so that I can adjust system behavior without restarting the server for dynamic settings.

### Ubiquitous Requirements (Config Properties)

- **ADM-05-001:** `The system shall load INI-based configuration from custom/conf/app.ini at startup.`
- **ADM-05-002:** `The system shall support environment variable overrides using the GITEA__section__key format.`
- **ADM-05-003:** `The system shall classify configuration settings as static (require restart) or dynamic (apply immediately).`
- **ADM-05-004:** `The system shall shadow passwords and secrets displayed in the configuration viewer.`
- **ADM-05-005:** `The system shall support the following major configuration sections: server, database, repository, ui, indexer, ssh, lfs, packages, actions, webhook, mailer, cache, session, log, markup, cron, mirror, api, oauth2, security, admin, attachment, federation, camo, and other.`

### Event-Driven Requirements (Config Workflow)

- **ADM-05-101:** `When an admin navigates to /admin/config, the system shall display the configuration summary with shadowed secrets.`
- **ADM-05-102:** `When an admin navigates to /admin/config/settings, the system shall display dynamic settings that can be changed at runtime.`
- **ADM-05-103:** `When an admin submits a configuration change via /admin/config/change, the system shall apply dynamic settings immediately and flag static settings as requiring restart.`
- **ADM-05-104:** `When an admin sends a test email via /admin/config/send_test_mail, the system shall dispatch a test email to the specified address.`
- **ADM-05-105:** `When the system starts with the install wizard, the system shall invoke LoadSettingsForInstall() to present initial configuration.`

### State-Driven Requirements (Config Lifecycle)

- **ADM-05-701:** `While a static configuration setting has been changed but not applied via restart, the system shall indicate the pending change in the configuration viewer.`

### Unwanted Behaviour Requirements (Config Errors)

- **ADM-05-301:** `If the configuration file is missing or unreadable at startup, then the system shall fail to start with a descriptive error.`
- **ADM-05-302:** `If an admin submits an invalid configuration value, then the system shall reject the change with a validation error.`

---

## 6. Queue Management

**User Story:** As an administrator, I want to monitor and manage background task queues so that I can ensure processing throughput and clear stuck items.

### Ubiquitous Requirements (Queue Properties)

- **ADM-06-001:** `The system shall support the following queue backends: channel (in-memory, default), LevelDB (persistent), Redis (distributed), and dummy (no-op).`
- **ADM-06-002:** `The system shall provide three queue types: simple FIFO queue, unique queue (prevents duplicate items), and worker pool queue (manages multiple workers).`
- **ADM-06-003:** `The system shall display queue status including number of workers, number of items, and queue length.`

### Event-Driven Requirements (Queue Workflow)

- **ADM-06-101:** `When an admin navigates to /admin/monitor/queues, the system shall display all queues with their current status.`
- **ADM-06-102:** `When an admin navigates to /admin/monitor/queue/:qid, the system shall display detailed information for the specified queue.`
- **ADM-06-103:** `When an admin adjusts the worker count for a queue, the system shall apply the change dynamically without restart.`
- **ADM-06-104:** `When an admin clears a queue, the system shall remove all pending items from that queue.`
- **ADM-06-105:** `When the system shuts down, the system shall gracefully drain all queues before terminating.`

### Optional Feature Requirements (Queue Configuration)

- **ADM-06-201:** `Where Redis is configured as the queue backend, the system shall connect using the CONN_STR setting and support distributed queue processing.`
- **ADM-06-202:** `Where LevelDB is configured as the queue backend, the system shall persist queue data to the DATADIR path.`
- **ADM-06-203:** `Where batch processing is configured, the system shall process queue items in groups up to BATCH_LENGTH size.`

### Unwanted Behaviour Requirements (Queue Errors)

- **ADM-06-301:** `If a queue backend connection fails, then the system shall log the error and fall back to in-memory processing where applicable.`
- **ADM-06-302:** `If a queue item processing fails, then the system shall log the error and retry according to the retry policy.`

---

## 7. Metrics/Monitoring

**User Story:** As an administrator, I want Prometheus-compatible metrics so that I can monitor my Gitea instance with standard observability tooling.

### Ubiquitous Requirements (Metrics Properties)

- **ADM-07-001:** `The system shall expose a Prometheus-compatible metrics endpoint at /metrics.`
- **ADM-07-002:** `The system shall expose build info metrics including goarch, goos, goversion, and version.`
- **ADM-07-003:** `The system shall expose aggregate count metrics for repositories, users, organizations, issues, labels, milestones, mirrors, releases, teams, webhooks, and hooks.`

### Event-Driven Requirements (Metrics Workflow)

- **ADM-07-101:** `When a client requests /metrics, the system shall return current metric values in Prometheus exposition format.`
- **ADM-07-102:** `When metrics are enabled with a bearer token, the system shall require valid token authentication for the metrics endpoint.`

### Optional Feature Requirements (Metrics Configuration)

- **ADM-07-201:** `Where ENABLED_ISSUE_BY_LABEL is configured, the system shall expose per-label issue count metrics.`
- **ADM-07-202:** `Where ENABLED_ISSUE_BY_REPOSITORY is configured, the system shall expose per-repository issue count metrics.`

### Unwanted Behaviour Requirements (Metrics Errors)

- **ADM-07-301:** `If metrics are not enabled in configuration, then the system shall return 404 for the metrics endpoint.`
- **ADM-07-302:** `If a client requests metrics without a valid bearer token when TOKEN is configured, then the system shall reject the request with 401 Unauthorized.`

---

## 8. Health Checks

**User Story:** As an administrator or load balancer, I want a health check endpoint so that I can determine if the Gitea instance is operating correctly.

### Ubiquitous Requirements (Health Check Properties)

- **ADM-08-001:** `The system shall expose a health check endpoint at /api/healthz.`
- **ADM-08-002:** `The system shall report health status as pass, fail, or warn.`
- **ADM-08-003:** `The system shall perform a database ping check as part of the health evaluation.`
- **ADM-08-004:** `The system shall perform a cache ping check as part of the health evaluation.`

### Event-Driven Requirements (Health Check Workflow)

- **ADM-08-101:** `When a client requests /api/healthz, the system shall execute all configured health checks and return the aggregate status.`
- **ADM-08-102:** `When all health checks pass, the system shall return status pass with HTTP 2xx-3xx.`
- **ADM-08-103:** `When any health check fails, the system shall return status fail with HTTP 4xx-5xx.`
- **ADM-08-104:** `When health checks pass but with concerns, the system shall return status warn with HTTP 2xx-3xx.`

### Unwanted Behaviour Requirements (Health Check Errors)

- **ADM-08-301:** `If the database is unreachable during a health check, then the system shall report status fail with the error output in the response.`

---

## 9. Database Migrations

**User Story:** As an administrator, I want automatic database schema migrations so that my database stays current when upgrading Gitea versions.

### Ubiquitous Requirements (Migration Properties)

- **ADM-09-001:** `The system shall track the current schema version in a database version table.`
- **ADM-09-002:** `The system shall use sequential version numbering for migrations (minimum version: 70).`
- **ADM-09-003:** `The system shall require each migration to provide a unique description string.`
- **ADM-09-004:** `The system shall organize migrations by release version directories (v1_6 through v1_23).`

### Event-Driven Requirements (Migration Workflow)

- **ADM-09-101:** `When the system starts and the database schema version is behind the latest, the system shall execute all pending migrations in sequential order.`
- **ADM-09-102:** `When a migration completes successfully, the system shall update the schema version to that migration's version number.`

### Unwanted Behaviour Requirements (Migration Errors)

- **ADM-09-301:** `If a migration fails during execution, then the system shall halt startup and report the migration error.`
- **ADM-09-302:** `If a migration is run out of order, then the system shall skip it and continue with the next required migration.`

---

## 10. Backup/Restore

**User Story:** As an administrator, I want to create a full system backup so that I can recover from data loss or migrate to another server.

### Ubiquitous Requirements (Backup Properties)

- **ADM-10-001:** `The system shall support the following backup archive formats (9 total): zip, tar, tar.sz, tar.gz, tar.xz, tar.bz2, tar.br, tar.lz4, and tar.zst.`
- **ADM-10-002:** `The system shall include database dump, repository files, configuration (app.ini), attachments, avatars, LFS objects, package files, and log files in the backup archive.`

### Event-Driven Requirements (Backup Workflow)

- **ADM-10-101:** `When the admin executes the dump command, the system shall create an archive containing all Gitea data.`
- **ADM-10-102:** `When the dump command completes, the system shall write the archive to the specified output path.`

### Optional Feature Requirements (Backup Selective Skip)

- **ADM-10-201:** `Where --skip-repository is specified, the system shall exclude repository data from the backup archive.`
- **ADM-10-202:** `Where --skip-lfs-data is specified, the system shall exclude LFS objects from the backup archive.`
- **ADM-10-203:** `Where --skip-attachment-data is specified, the system shall exclude attachments from the backup archive.`
- **ADM-10-204:** `Where --skip-package-data is specified, the system shall exclude package files from the backup archive.`
- **ADM-10-205:** `Where --skip-log is specified, the system shall exclude log files from the backup archive.`
- **ADM-10-206:** `Where --skip-custom-dir is specified, the system shall exclude the custom directory from the backup archive.`
- **ADM-10-207:** `Where --skip-index is specified, the system shall exclude bleve index data from the backup archive.`
- **ADM-10-208:** `Where --skip-db is specified, the system shall exclude the database dump from the backup archive.`

### Unwanted Behaviour Requirements (Backup Errors)

- **ADM-10-301:** `If the output path is not writable during dump, then the system shall fail with a permission error.`
- **ADM-10-302:** `If the temporary directory (--tempdir) has insufficient disk space, then the system shall fail with a disk space error.`

---

## 11. Logging

**User Story:** As an administrator, I want configurable logging so that I can capture system events at appropriate levels and route them to multiple outputs.

### Ubiquitous Requirements (Logging Properties)

- **ADM-11-001:** `The system shall support the following log levels in order of severity: TRACE, DEBUG, INFO, WARN, ERROR, FATAL.`
- **ADM-11-002:** `The system shall support multiple simultaneous log output modes (console, file, conn).`
- **ADM-11-003:** `The system shall allow independent configuration per log mode via [log.*] configuration sections.`

### Event-Driven Requirements (Logging Workflow)

- **ADM-11-101:** `When a log event occurs at a level equal to or above the configured threshold, the system shall write the event to all configured log modes.`
- **ADM-11-102:** `When the rotating file writer reaches the configured size limit, the system shall rotate the log file.`
- **ADM-11-103:** `When a FATAL level event is logged, the system shall terminate the process.`

### Optional Feature Requirements (Logging Modes)

- **ADM-11-201:** `Where console mode is configured, the system shall output log events to the terminal with color support.`
- **ADM-11-202:** `Where file mode is configured, the system shall write log events to a rotating file on disk.`
- **ADM-11-203:** `Where conn mode is configured, the system shall send log events to a network connection.`

### Unwanted Behaviour Requirements (Logging Errors)

- **ADM-11-301:** `If a log file cannot be written (permission denied, disk full), then the system shall fall back to console logging and report the error.`

---

## 12. Debug/PProf

**User Story:** As an administrator, I want built-in profiling endpoints so that I can diagnose performance issues without installing external tooling.

### Ubiquitous Requirements (PProf Properties)

- **ADM-12-001:** `The system shall, when pprof is enabled, run it on a separate HTTP server listening on localhost:6060 (loopback only), not on the main Gitea web server.`
- **ADM-12-002:** `The system shall expose the standard net/http/pprof endpoints on the pprof server: heap, goroutine, threadcreate, block, mutex, cmdline, profile, symbol, and trace, under /debug/pprof/.`
- **ADM-12-003:** `The system shall expose a /debug/fgprof endpoint on the pprof server for on-CPU and off-CPU activity reporting.`
- **ADM-12-004:** `The system shall expose an admin diagnosis endpoint at /admin/monitor/diagnosis that returns a ZIP archive containing a CPU profile, goroutine dump (before and after), and a heap dump.`
- **ADM-12-005:** `The system shall expose an admin stacktrace viewer at /admin/monitor/stacktrace that lists goroutines and allows cancelling processes by PID.`

### Event-Driven Requirements (PProf Workflow)

- **ADM-12-101:** `When a client requests /debug/pprof/profile on the pprof server (localhost:6060) with a duration parameter, the system shall capture a CPU profile for that duration and return it.`
- **ADM-12-102:** `When a client requests /debug/pprof/trace on the pprof server (localhost:6060) with a duration parameter, the system shall capture an execution trace for that duration and return it.`

### Optional Feature Requirements (PProf Configuration)

- **ADM-12-201:** `Where [server].ENABLE_PPROF is set to true, the system shall start the pprof server on localhost:6060 bound to loopback only and shall not expose profiling endpoints on the public web server.`

### Unwanted Behaviour Requirements (PProf Constraints)

- **ADM-12-301:** `If [server].ENABLE_PPROF is not set (default false), then the system shall not start the pprof server and shall not serve any /debug/pprof/* routes on the main web server.`
- **ADM-12-302:** `If a CPU profile or trace capture is requested with an excessively long duration, then the system shall limit the duration to a safe maximum.`

---

## 13. Admin Notices (ADM-13)

**User Story:** As an administrator, I want a notices panel showing system warnings and errors so that I can investigate incidents and operational anomalies.

### Ubiquitous Requirements (Notice Properties)

- **ADM-13-001:** `The system shall record admin notices for noteworthy runtime events including cron failures, migration errors, and storage backend issues.`
- **ADM-13-002:** `The system shall classify each notice by type using NoticeRepository (1) and NoticeTask (2) as the only two notice types.`
- **ADM-13-003:** `The system shall timestamp and store the notice description and originating context.`

### Event-Driven Requirements (Notice Workflow)

- **ADM-13-101:** `When a runtime event warrants admin attention, the system shall insert a notice record.`
- **ADM-13-102:** `When an admin navigates to /admin/notices, the system shall render the paginated notice list.`
- **ADM-13-103:** `When an admin dismisses a notice, the system shall permanently delete the notice record from the database.`
- **ADM-13-104:** `When an admin triggers bulk dismissal, the system shall permanently delete every selected notice by ID.`
- **ADM-13-105:** `When a new notice is created after the admin dashboard was last viewed, the system shall display a notice badge on the admin navigation.`

### Unwanted Behaviour Requirements (Notice Errors)

- **ADM-13-301:** `If a non-admin attempts to view the notices page, then the system shall return 403 Forbidden.`
- **ADM-13-302:** `If a notice references an entity that has since been deleted, then the system shall render the notice with a "stale reference" indicator rather than failing.`

---

## 14. Doctor Diagnostic Checks (ADM-14)

**User Story:** As an administrator, I want to run consistency checks and apply repairs so that I can detect and fix data corruption, broken references, and misconfiguration.

### Ubiquitous Requirements (Doctor Properties)

- **ADM-14-001:** `The system shall provide a doctor subcommand exposing a curated set of consistency checks.`
- **ADM-14-002:** `The system shall classify each check as either read-only (diagnostic) or repairing (writes changes).`
- **ADM-14-003:** `The system shall expose a non-exhaustive set of approximately 27 consistency checks including: authorized SSH keys, repository count integrity, dangling LFS objects, orphaned hooks, broken symlink references, stale archive caches, script-type availability, recalculate-stars-number, enable-push-options, check-git-daemon-export-ok, check-commit-graphs, recalculate-merge-bases, check-user-email, check-user-names, check-user-type, disable-mirror-actions-unit, fix-broken-repo-units, fix-owner-team-create-org-repo, synchronize-repo-heads, storage-attachments, storage-avatars, storage-packages, database version, paths consistency, DB consistency, and fix #16961 / #8312.`
- **ADM-14-004:** `The system shall provide a doctor recreate-table subcommand that recreates database tables from XORM definitions and copies data over, deleting old columns.`
- **ADM-14-005:** `The system shall provide a doctor convert subcommand that converts the database between storage configurations.`

### Event-Driven Requirements (Doctor Workflow)

- **ADM-14-101:** `When an admin runs the doctor check command without arguments, the system shall run all read-only checks and report any failures.`
- **ADM-14-102:** `When an admin runs doctor with a specific check name, the system shall execute only that check.`
- **ADM-14-103:** `When an admin runs doctor with the --fix flag, the system shall execute every repairing check and apply fixes.`
- **ADM-14-104:** `When an admin runs doctor check, the system shall accept the following flags: list, default, run, all, fix, log-file, and color.`
- **ADM-14-105:** `When a check completes successfully, the system shall print a summary and exit 0; on failure the system shall exit non-zero.`

### State-Driven Requirements (Doctor Safety)

- **ADM-14-701:** `While the system is running in --fix mode, the system shall log every applied change to the audit log and to the admin notices panel.`

### Optional Feature Requirements (Doctor Configuration)

- **ADM-14-201:** `Where a check is run via the CLI with --log-level=debug, the system shall emit per-item diagnostic output.`

### Unwanted Behaviour Requirements (Doctor Errors)

- **ADM-14-301:** `If a repairing check is run without --fix, then the system shall report the issues it would fix but shall not modify any data.`
- **ADM-14-302:** `If the database is unreachable during a doctor run, the system shall halt the run and report the failure.`
- **ADM-14-303:** `If a check fails due to an internal error, the system shall log the error, continue with remaining checks, and report the failure in the summary.`

---

## 15. CLI Command Surface (ADM-15)

**User Story:** As an operator, I want a CLI entry point so that I can run Gitea as a server, perform admin operations, and execute maintenance tasks from the command line.

### Ubiquitous Requirements (CLI Properties)

- **ADM-15-001:** `The system shall expose a single gitea binary with subcommands including web, serv, admin, doctor, dump, hook, keys, migrate, migrate-storage, manager, embedded, dump-repository, restore-repository, actions, cert, generate, docs, and help.`
- **ADM-15-002:** `The system shall load configuration from custom/conf/app.ini before executing any subcommand.`
- **ADM-15-003:** `The system shall accept command-line flags and GITEA__ prefixed environment variables as configuration overrides.`

### Event-Driven Requirements (CLI Workflow)

- **ADM-15-101:** `When an operator runs gitea web, the system shall initialize configuration, database, services, and start the HTTP and (optionally) SSH servers.`
- **ADM-15-102:** `When the SSH server invokes gitea serv on an authenticated SSH session, the system shall execute the requested git-upload-pack or git-receive-pack as the authenticated user.`
- **ADM-15-103:** `When git invokes gitea hook pre-receive, update, or post-receive, the system shall execute the corresponding server-side hook pipeline.`
- **ADM-15-104:** `When an operator runs gitea admin user create, the system shall create a user account without requiring web access.`
- **ADM-15-105:** `When an operator runs gitea admin auth list/add/update/delete, the system shall manage external authentication sources from the command line.`
- **ADM-15-106:** `When an operator runs gitea dump, the system shall produce a backup archive per ADM-10.`
- **ADM-15-107:** `When an operator runs gitea migrate, the system shall execute pending database migrations without starting the server.`
- **ADM-15-108:** `When an operator runs gitea migrate-storage, the system shall transfer attachments, LFS objects, packages, or avatars between storage backends.`
- **ADM-15-109:** `When an operator runs gitea admin sendmail, the system shall dispatch an email message to all users from the command line.`
- **ADM-15-110:** `When an operator runs gitea admin regenerate hooks, the system shall regenerate repository hook files from the command line.`
- **ADM-15-111:** `When an operator runs gitea admin regenerate keys, the system shall regenerate the authorized_keys file from the command line.`
- **ADM-15-112:** `When an operator runs gitea manager with a subcommand (shutdown, restart, reload-templates, flush-queues, logging, processes), the system shall manage the running Gitea process.`

### Optional Feature Requirements (CLI Configuration)

- **ADM-15-201:** `Where the operator passes --custom-path or --config, the system shall load configuration from the specified location instead of the default.`
- **ADM-15-202:** `Where the operator passes --work-path, the system shall treat the specified directory as the Gitea working root.`

### Unwanted Behaviour Requirements (CLI Errors)

- **ADM-15-301:** `If configuration loading fails before subcommand dispatch, then the system shall exit non-zero with a descriptive error.`
- **ADM-15-302:** `If a subcommand receives invalid flags or arguments, then the system shall print usage and exit non-zero.`
- **ADM-15-303:** `If a subcommand requires database access and the database is unreachable, then the system shall report the connection failure and exit non-zero.`

---

## Configuration Reference

The following INI sections configure administration and operations behaviors.

### [cache] Section

- **ADAPTER**: Cache backend (`memory`, `redis`, `memcache`, `twoqueue`).
- **INTERVAL**: Garbage collection interval for the memory adapter (seconds, default 60).
- **HOST**: Connection string for redis/memcache backends (`redis://host:port`, etc.).
- **ITEM_TTL**: Time-to-live for cached items.

### [cache.last_commit] Section

- **ENABLED**: Enable caching of last-commit metadata for repository file views.
- **ITEM_TTL**: Time-to-live for last-commit cache entries (default 8760h = 1 year).
- **COMMITS_COUNT**: Minimum commits required for a repository to use the cache.

### [database] Section

- **DB_TYPE**: Database backend (`sqlite3`, `mysql`, `postgres`, `mssql`).
- **HOST**, **NAME**, **USER**, **PASSWD**, **SCHEMA**: Connection parameters.
- **SSL_MODE**: SSL mode (`disable`, `require`, `verify-ca`, `verify-full`).
- **CHARSET**: Connection charset (default `utf8` for MySQL).
- **PATH**: SQLite file path.
- **LOG_SQL**: Log every SQL statement at debug level.
- **DB_RETRIES**, **DB_RETRY_BACKOFF**: Connection retry behavior during startup.
- **MAX_IDLE_CONNS**, **MAX_OPEN_CONNS**, **CONN_MAX_LIFETIME**: Connection pool tuning.

### [queue] Section

- **TYPE**: Queue backend (`channel`, `level`, `redis`, `dummy`; default `level`).
- **DATADIR**: Data directory for level/leveldb-backed queues (relative to AppDataPath, default `queues/common`).
- **CONN_STR**: Connection string for redis/leveldb backends (e.g. `redis://127.0.0.1:6379/0`).
- **LENGTH**: Maximum queue length before blocking (default 100000).
- **BATCH_LENGTH**: Number of items processed per batch (default 20).
- **MAX_WORKERS**: Maximum number of workers (default NumCPU/2, clamped to 1-10).

### [metrics] Section

- **ENABLED**: Enable the Prometheus metrics endpoint at /metrics (default false).
- **TOKEN**: Bearer token required to access the metrics endpoint (empty = no auth).
- **ENABLED_ISSUE_BY_LABEL**: Expose per-label issue count metrics (default false).
- **ENABLED_ISSUE_BY_REPOSITORY**: Expose per-repository issue count metrics (default false).

### [server] PProf Keys

- **ENABLE_PPROF**: Enable the pprof server on localhost:6060 (default false).
- **PPROF_DATA_PATH**: Directory for pprof data files (default `<AppWorkPath>/data/tmp/pprof`).

---

## Business Rules

- **BR-07-001:** Only users with admin privilege may access /admin/* routes and /api/v1/admin/* endpoints
- **BR-07-002:** Configuration secrets (passwords, tokens, secret keys) are always shadowed in the configuration viewer
- **BR-07-003:** Static configuration changes require a server restart to take effect
- **BR-07-004:** Dynamic configuration changes take effect immediately without restart
- **BR-07-005:** Database migrations execute in strict sequential order by version number
- **BR-07-006:** The minimum database schema version is 70
- **BR-07-007:** Queue worker count adjustments are dynamic and do not require restart
- **BR-07-008:** Metrics endpoint requires ENABLED=true in [metrics] configuration
- **BR-07-009:** Health check endpoint (/api/healthz) is always available regardless of configuration
- **BR-07-010:** PProf server is disabled by default and must be explicitly enabled via [server].ENABLE_PPROF; when enabled it binds to localhost:6060 only
- **BR-07-011:** Backup archives contain all data by default; individual categories are excluded via flags
- **BR-07-012:** Log levels are hierarchical: TRACE < DEBUG < INFO < WARN < ERROR < FATAL
- **BR-07-013:** Authentication source order determines precedence when multiple sources match a user (see AUTH-06 for source types and features)
- **BR-07-014:** Ghost users are system accounts used for attribution preservation after user deletion

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Non-admin accesses /admin routes | Return 403 Forbidden |
| Configuration file missing at startup | Fail to start with descriptive error |
| Invalid config value submitted via UI | Reject change with validation error |
| Database unreachable during health check | Report status fail with error output |
| Migration failure during startup | Halt startup, report migration error |
| Insufficient disk space during backup | Fail with disk space error |
| Queue backend connection loss | Log error, fall back to in-memory processing |
| Metrics endpoint requested when disabled | Return 404 |
| PProf requested when not enabled | PProf server not started; no /debug/pprof/* routes served |
| Log file write failure | Fall back to console logging |
| LDAP sync during source outage | Log error, continue with cached user data |
| Last admin user deletion attempt | Deny deletion |
| Unadopted repo with invalid data | Skip directory during unadopted scan |

## Success Criteria

- Admin dashboard renders system statistics within 1 second
- Configuration viewer loads with all secrets properly shadowed
- Dynamic configuration changes apply within 1 second of submission
- Queue status page reflects real-time worker and item counts
- Prometheus metrics endpoint responds within 500ms
- Health check endpoint responds within 2 seconds including all check types
- Database migrations execute sequentially with zero data loss
- Full backup completes within a time proportional to data size
- Log events propagate to all configured modes within 100ms
- PProf CPU profile capture respects the requested duration within 1 second tolerance
- Auth source CRUD operations take effect immediately for subsequent login attempts
- User search returns results within 2 seconds for instances with up to 100,000 users

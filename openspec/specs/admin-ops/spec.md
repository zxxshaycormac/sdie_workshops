# 08 — Administration & Operations

Baseline specification of Gitea's administration dashboard, configuration, observability, backup, debug subsystems, background infrastructure (cron, queues, indexers, storage), and operational tooling. All requirements describe the current (v1.22.x) system behavior.

---

## 1. Admin Dashboard

**User Story:** As an administrator, I want a central dashboard so that I can monitor system health and perform quick operational tasks.

### Ubiquitous Requirements (Dashboard Properties)

- **OPS-01-001:** `The system shall display system statistics on the admin dashboard including memory usage, goroutine count, and GC stats.`
- **OPS-01-002:** `The system shall present an update checker status indicating whether a newer Gitea version is available.`
- **OPS-01-003:** `The system shall list quick operations on the dashboard including sync branches, sync tags, and run cron tasks.`

### Event-Driven Requirements (Dashboard Workflow)

- **OPS-01-101:** `When an admin navigates to /admin, the system shall render the dashboard with current system statistics.`
- **OPS-01-102:** `When an admin triggers a quick operation from the dashboard, the system shall execute the operation and display the result.`
- **OPS-01-103:** `When an admin navigates to /admin/self_check, the system shall run self-check diagnostics and display any warnings.`
- **OPS-01-104:** `When an admin navigates to /admin/system_status, the system shall display detailed system status information.`

### Unwanted Behaviour Requirements (Dashboard Errors)

- **OPS-01-301:** `If a non-admin user attempts to access /admin routes, then the system shall return a 403 Forbidden response.`

---

## 2. User Management

**User Story:** As an administrator, I want to manage user accounts so that I can create, modify, and remove user access as needed.

### Ubiquitous Requirements (User List Properties)

- **OPS-02-001:** `The system shall list all user accounts with pagination support.`
- **OPS-02-002:** `The system shall support filtering the user list by status: active, admin, restricted, 2FA enabled, and prohibited login.`
- **OPS-02-003:** `The system shall support searching users by username or email address.`

### Event-Driven Requirements (User CRUD Workflow)

- **OPS-02-101:** `When an admin creates a new user via the admin panel, the system shall create the account in the active state.`
- **OPS-02-102:** `When an admin edits a user account, the system shall persist the changes immediately.`
- **OPS-02-103:** `When an admin deletes a user account, the system shall transfer the user's repository ownership to a ghost user for attribution preservation.`
- **OPS-02-104:** `When an admin sets a new password for a user, the system shall hash and store the password using the configured algorithm.`

### Event-Driven Requirements (Password & 2FA Management)

- **OPS-02-105:** `When an admin forces a password change for a user, the system shall set the must_change_password flag so the user must change their password at next login.`
- **OPS-02-106:** `When an admin resets a user's 2FA, the system shall disable 2FA for that user and regenerate scratch codes.`
- **OPS-02-107:** `When a new password is set and pwned password checking is enabled, the system shall check the password against known breach databases.`

### Optional Feature Requirements (Avatar Management)

- **OPS-02-201:** `Where a user has a custom avatar, the system shall display the custom avatar instead of the Gravatar default.`

### Unwanted Behaviour Requirements (User Management Errors)

- **OPS-02-301:** `If an admin attempts to delete the last admin user, then the system shall deny the deletion.`
- **OPS-02-302:** `If an admin submits a password that does not meet complexity requirements, then the system shall reject the password and display the requirements.`

---

## 3. Repository Management

**User Story:** As an administrator, I want to manage repositories across the instance so that I can delete repos and adopt unadopted ones from the filesystem.

### Ubiquitous Requirements (Repository List Properties)

- **OPS-03-001:** `The system shall list all repositories with pagination support.`
- **OPS-03-002:** `The system shall support searching repositories by name.`

### Event-Driven Requirements (Repository Workflow)

- **OPS-03-101:** `When an admin deletes a repository, the system shall remove the repository data and all associated resources.`
- **OPS-03-102:** `When an admin navigates to /admin/repos/unadopted, the system shall scan the repository root path for directories not registered in the database.`
- **OPS-03-103:** `When an admin adopts an unadopted repository, the system shall register it in the database and make it accessible.`
- **OPS-03-104:** `When an admin deletes an unadopted repository, the system shall remove the directory from the filesystem.`

### Unwanted Behaviour Requirements (Repository Errors)

- **OPS-03-301:** `If an unadopted repository directory contains invalid data, then the system shall skip the directory during the unadopted scan.`

---

## 4. Auth Source Management

**User Story:** As an administrator, I want to configure external authentication sources so that users can authenticate against LDAP, SMTP, OAuth2, and other providers.

### Ubiquitous Requirements (Admin Interface)

- **OPS-04-001:** `The system shall expose auth source CRUD operations in the admin panel (see Domain 01 (Identity & Authentication) §External Authentication Sources for supported source types, lifecycle behavior, optional features, and error handling).`
- **OPS-04-002:** `The system shall store authentication source configurations with activation state and display order.`
- **OPS-04-003:** `The system shall allow TLS configuration per authentication source.`

### Event-Driven Requirements (Admin Operations)

- **OPS-04-101:** `When an admin triggers an LDAP sync from the admin panel, the system shall synchronize user data from the LDAP directory to local accounts.`
- **OPS-04-102:** `When an admin reorders authentication sources, the system shall update the display order and login precedence.`

---

## 5. Configuration System

**User Story:** As an administrator, I want to view and modify Gitea's configuration so that I can adjust system behavior without restarting the server for dynamic settings.

### Ubiquitous Requirements (Config Properties)

- **OPS-05-001:** `The system shall load INI-based configuration from custom/conf/app.ini at startup.`
- **OPS-05-002:** `The system shall support environment variable overrides using the GITEA__section__key format.`
- **OPS-05-003:** `The system shall classify configuration settings as static (require restart) or dynamic (apply immediately).`
- **OPS-05-004:** `The system shall shadow passwords and secrets displayed in the configuration viewer.`
- **OPS-05-005:** `The system shall support the following major configuration sections: server, database, repository, ui, indexer, ssh, lfs, packages, actions, webhook, mailer, cache, session, log, markup, cron, queue, metrics, storage, mirror, api, oauth2, security, admin, attachment, federation, camo, update_checker, and other.`

### Event-Driven Requirements (Config Workflow)

- **OPS-05-101:** `When an admin navigates to /admin/config, the system shall display the configuration summary with shadowed secrets.`
- **OPS-05-102:** `When an admin navigates to /admin/config/settings, the system shall display dynamic settings that can be changed at runtime.`
- **OPS-05-103:** `When an admin submits a configuration change via /admin/config/change, the system shall apply dynamic settings immediately and flag static settings as requiring restart.`
- **OPS-05-104:** `When an admin sends a test email via /admin/config/send_test_mail, the system shall dispatch a test email to the specified address.`
- **OPS-05-105:** `When the system starts with the install wizard, the system shall invoke LoadSettingsForInstall() to present initial configuration.`

### State-Driven Requirements (Config Lifecycle)

- **OPS-05-701:** `While a static configuration setting has been changed but not applied via restart, the system shall indicate the pending change in the configuration viewer.`

### Unwanted Behaviour Requirements (Config Errors)

- **OPS-05-301:** `If the configuration file is missing or unreadable at startup, then the system shall fail to start with a descriptive error.`
- **OPS-05-302:** `If an admin submits an invalid configuration value, then the system shall reject the change with a validation error.`

---

## 6. Queue System

**User Story:** As a system operator, I want a reliable queue infrastructure so that background tasks (webhooks, indexing, mail, actions) are processed concurrently without data loss, and as an administrator, I want to monitor and manage those queues.

### Ubiquitous Requirements (Queue Properties)

- **OPS-06-001:** `The system shall support the following queue backends: channel (in-memory, default), LevelDB (persistent), Redis (distributed), and dummy (no-op).`
- **OPS-06-002:** `The system shall provide three queue types: simple FIFO queue, unique queue (prevents duplicate items), and worker pool queue (manages multiple workers).`
- **OPS-06-003:** `The system shall display queue status including number of workers, number of items, and queue length.`
- **OPS-06-004:** `The system shall use a worker pool model with dynamic sizing based on CPU count.`
- **OPS-06-005:** `The system shall support batch processing of queue items with configurable batch length.`

### Event-Driven Requirements (Queue Workflow)

- **OPS-06-101:** `When an admin navigates to /admin/monitor/queues, the system shall display all queues with their current status.`
- **OPS-06-102:** `When an admin navigates to /admin/monitor/queue/:qid, the system shall display detailed information for the specified queue.`
- **OPS-06-103:** `When an admin adjusts the worker count for a queue, the system shall apply the change dynamically without restart.`
- **OPS-06-104:** `When an admin clears a queue, the system shall remove all pending items from that queue.`
- **OPS-06-105:** `When the system shuts down, the system shall gracefully drain all queues before terminating.`
- **OPS-06-106:** `When a queue item is pushed to a unique queue, the system shall deduplicate by key and drop or replace existing items with the same key.`
- **OPS-06-107:** `When a worker becomes available, the system shall dispatch the next queued item or batch to that worker.`
- **OPS-06-108:** `When a queue item processing fails, the system shall retry the item according to the configured retry policy.`

### State-Driven Requirements (Queue Lifecycle)

- **OPS-06-701:** `While the LevelDB backend is active, the system shall persist queue items to disk so that items survive process restarts.`
- **OPS-06-702:** `While the Redis backend is active, the system shall distribute queue items across multiple Gitea instances for horizontal scaling.`
- **OPS-06-703:** `While the queue length exceeds the configured maximum (default 100,000), the system shall apply backpressure to enqueue operations.`

### Optional Feature Requirements (Queue Configuration)

- **OPS-06-201:** `Where Redis is configured as the queue backend, the system shall connect using the CONN_STR setting and support distributed queue processing.`
- **OPS-06-202:** `Where LevelDB is configured as the queue backend, the system shall persist queue data to the DATADIR path.`
- **OPS-06-203:** `Where batch processing is configured, the system shall process queue items in groups up to BATCH_LENGTH size.`
- **OPS-06-204:** `Where MAX_WORKERS is configured, the system shall cap the worker pool at the specified count (default CPU count / 2, maximum 10).`

### Unwanted Behaviour Requirements (Queue Errors)

- **OPS-06-301:** `If a queue backend connection fails, then the system shall log the error and fall back to in-memory processing where applicable.`
- **OPS-06-302:** `If a queue item processing fails, then the system shall log the error and retry according to the retry policy.`
- **OPS-06-303:** `If a queue backend becomes unavailable, then the system shall log the error and queue new items in memory until the backend recovers.`
- **OPS-06-304:** `If a worker panics during item processing, then the system shall recover the worker and requeue the item.`
- **OPS-06-305:** `If a LevelDB queue data directory is corrupted, then the system shall log the error and start with a fresh queue.`

---

## 7. Metrics / Monitoring

**User Story:** As an administrator, I want Prometheus-compatible metrics so that I can monitor my Gitea instance with standard observability tooling.

### Ubiquitous Requirements (Metrics Properties)

- **OPS-07-001:** `The system shall expose a Prometheus-compatible metrics endpoint at /metrics.`
- **OPS-07-002:** `The system shall expose build info metrics including goarch, goos, goversion, and version.`
- **OPS-07-003:** `The system shall expose aggregate count metrics for repositories, users, organizations, issues, labels, milestones, mirrors, releases, teams, webhooks, and hooks.`

### Event-Driven Requirements (Metrics Workflow)

- **OPS-07-101:** `When a client requests /metrics, the system shall return current metric values in Prometheus exposition format.`
- **OPS-07-102:** `When metrics are enabled with a bearer token, the system shall require valid token authentication for the metrics endpoint.`

### Optional Feature Requirements (Metrics Configuration)

- **OPS-07-201:** `Where ENABLED_ISSUE_BY_LABEL is configured, the system shall expose per-label issue count metrics.`
- **OPS-07-202:** `Where ENABLED_ISSUE_BY_REPOSITORY is configured, the system shall expose per-repository issue count metrics.`

### Unwanted Behaviour Requirements (Metrics Errors)

- **OPS-07-301:** `If metrics are not enabled in configuration, then the system shall return 404 for the metrics endpoint.`
- **OPS-07-302:** `If a client requests metrics without a valid bearer token when TOKEN is configured, then the system shall reject the request with 401 Unauthorized.`

---

## 8. Health Checks

**User Story:** As an administrator or load balancer, I want a health check endpoint so that I can determine if the Gitea instance is operating correctly.

### Ubiquitous Requirements (Health Check Properties)

- **OPS-08-001:** `The system shall expose a health check endpoint at /api/healthz.`
- **OPS-08-002:** `The system shall report health status as pass, fail, or warn.`
- **OPS-08-003:** `The system shall perform a database ping check as part of the health evaluation.`
- **OPS-08-004:** `The system shall perform a cache ping check as part of the health evaluation.`

### Event-Driven Requirements (Health Check Workflow)

- **OPS-08-101:** `When a client requests /api/healthz, the system shall execute all configured health checks and return the aggregate status.`
- **OPS-08-102:** `When all health checks pass, the system shall return status pass with HTTP 2xx-3xx.`
- **OPS-08-103:** `When any health check fails, the system shall return status fail with HTTP 4xx-5xx.`
- **OPS-08-104:** `When health checks pass but with concerns, the system shall return status warn with HTTP 2xx-3xx.`

### Unwanted Behaviour Requirements (Health Check Errors)

- **OPS-08-301:** `If the database is unreachable during a health check, then the system shall report status fail with the error output in the response.`

---

## 9. Database Migrations

**User Story:** As an administrator, I want automatic database schema migrations so that my database stays current when upgrading Gitea versions.

### Ubiquitous Requirements (Migration Properties)

- **OPS-09-001:** `The system shall track the current schema version in a database version table.`
- **OPS-09-002:** `The system shall use sequential version numbering for migrations (minimum version: 70).`
- **OPS-09-003:** `The system shall require each migration to provide a unique description string.`
- **OPS-09-004:** `The system shall organize migrations by release version directories (v1_6 through v1_23).`

### Event-Driven Requirements (Migration Workflow)

- **OPS-09-101:** `When the system starts and the database schema version is behind the latest, the system shall execute all pending migrations in sequential order.`
- **OPS-09-102:** `When a migration completes successfully, the system shall update the schema version to that migration's version number.`

### Unwanted Behaviour Requirements (Migration Errors)

- **OPS-09-301:** `If a migration fails during execution, then the system shall halt startup and report the migration error.`
- **OPS-09-302:** `If a migration is run out of order, then the system shall skip it and continue with the next required migration.`

---

## 10. Backup / Restore

**User Story:** As an administrator, I want to create a full system backup so that I can recover from data loss or migrate to another server.

### Ubiquitous Requirements (Backup Properties)

- **OPS-10-001:** `The system shall support the following backup archive formats (9 total): zip, tar, tar.sz, tar.gz, tar.xz, tar.bz2, tar.br, tar.lz4, and tar.zst.`
- **OPS-10-002:** `The system shall include database dump, repository files, configuration (app.ini), attachments, avatars, LFS objects, package files, and log files in the backup archive.`

### Event-Driven Requirements (Backup Workflow)

- **OPS-10-101:** `When the admin executes the dump command, the system shall create an archive containing all Gitea data.`
- **OPS-10-102:** `When the dump command completes, the system shall write the archive to the specified output path.`

### Optional Feature Requirements (Backup Selective Skip)

- **OPS-10-201:** `Where --skip-repository is specified, the system shall exclude repository data from the backup archive.`
- **OPS-10-202:** `Where --skip-lfs-data is specified, the system shall exclude LFS objects from the backup archive.`
- **OPS-10-203:** `Where --skip-attachment-data is specified, the system shall exclude attachments from the backup archive.`
- **OPS-10-204:** `Where --skip-package-data is specified, the system shall exclude package files from the backup archive.`
- **OPS-10-205:** `Where --skip-log is specified, the system shall exclude log files from the backup archive.`
- **OPS-10-206:** `Where --skip-custom-dir is specified, the system shall exclude the custom directory from the backup archive.`
- **OPS-10-207:** `Where --skip-index is specified, the system shall exclude bleve index data from the backup archive.`
- **OPS-10-208:** `Where --skip-db is specified, the system shall exclude the database dump from the backup archive.`

### Unwanted Behaviour Requirements (Backup Errors)

- **OPS-10-301:** `If the output path is not writable during dump, then the system shall fail with a permission error.`
- **OPS-10-302:** `If the temporary directory (--tempdir) has insufficient disk space, then the system shall fail with a disk space error.`

---

## 11. Logging

**User Story:** As an administrator, I want configurable logging so that I can capture system events at appropriate levels and route them to multiple outputs.

### Ubiquitous Requirements (Logging Properties)

- **OPS-11-001:** `The system shall support the following log levels in order of severity: TRACE, DEBUG, INFO, WARN, ERROR, FATAL.`
- **OPS-11-002:** `The system shall support multiple simultaneous log output modes (console, file, conn).`
- **OPS-11-003:** `The system shall allow independent configuration per log mode via [log.*] configuration sections.`

### Event-Driven Requirements (Logging Workflow)

- **OPS-11-101:** `When a log event occurs at a level equal to or above the configured threshold, the system shall write the event to all configured log modes.`
- **OPS-11-102:** `When the rotating file writer reaches the configured size limit, the system shall rotate the log file.`
- **OPS-11-103:** `When a FATAL level event is logged, the system shall terminate the process.`

### Optional Feature Requirements (Logging Modes)

- **OPS-11-201:** `Where console mode is configured, the system shall output log events to the terminal with color support.`
- **OPS-11-202:** `Where file mode is configured, the system shall write log events to a rotating file on disk.`
- **OPS-11-203:** `Where conn mode is configured, the system shall send log events to a network connection.`

### Unwanted Behaviour Requirements (Logging Errors)

- **OPS-11-301:** `If a log file cannot be written (permission denied, disk full), then the system shall fall back to console logging and report the error.`

---

## 12. Debug / PProf

**User Story:** As an administrator, I want built-in profiling endpoints so that I can diagnose performance issues without installing external tooling.

### Ubiquitous Requirements (PProf Properties)

- **OPS-12-001:** `The system shall, when pprof is enabled, run it on a separate HTTP server listening on localhost:6060 (loopback only), not on the main Gitea web server.`
- **OPS-12-002:** `The system shall expose the standard net/http/pprof endpoints on the pprof server: heap, goroutine, threadcreate, block, mutex, cmdline, profile, symbol, and trace, under /debug/pprof/.`
- **OPS-12-003:** `The system shall expose a /debug/fgprof endpoint on the pprof server for on-CPU and off-CPU activity reporting.`
- **OPS-12-004:** `The system shall expose an admin diagnosis endpoint at /admin/monitor/diagnosis that returns a ZIP archive containing a CPU profile, goroutine dump (before and after), and a heap dump.`
- **OPS-12-005:** `The system shall expose an admin stacktrace viewer at /admin/monitor/stacktrace that lists goroutines and allows cancelling processes by PID.`

### Event-Driven Requirements (PProf Workflow)

- **OPS-12-101:** `When a client requests /debug/pprof/profile on the pprof server (localhost:6060) with a duration parameter, the system shall capture a CPU profile for that duration and return it.`
- **OPS-12-102:** `When a client requests /debug/pprof/trace on the pprof server (localhost:6060) with a duration parameter, the system shall capture an execution trace for that duration and return it.`

### Optional Feature Requirements (PProf Configuration)

- **OPS-12-201:** `Where [server].ENABLE_PPROF is set to true, the system shall start the pprof server on localhost:6060 bound to loopback only and shall not expose profiling endpoints on the public web server.`

### Unwanted Behaviour Requirements (PProf Constraints)

- **OPS-12-301:** `If [server].ENABLE_PPROF is not set (default false), then the system shall not start the pprof server and shall not serve any /debug/pprof/* routes on the main web server.`
- **OPS-12-302:** `If a CPU profile or trace capture is requested with an excessively long duration, then the system shall limit the duration to a safe maximum.`

---

## 13. Admin Notices

**User Story:** As an administrator, I want a notices panel showing system warnings and errors so that I can investigate incidents and operational anomalies.

### Ubiquitous Requirements (Notice Properties)

- **OPS-13-001:** `The system shall record admin notices for noteworthy runtime events including cron failures, migration errors, and storage backend issues.`
- **OPS-13-002:** `The system shall classify each notice by type using NoticeRepository (1) and NoticeTask (2) as the only two notice types.`
- **OPS-13-003:** `The system shall timestamp and store the notice description and originating context.`

### Event-Driven Requirements (Notice Workflow)

- **OPS-13-101:** `When a runtime event warrants admin attention, the system shall insert a notice record.`
- **OPS-13-102:** `When an admin navigates to /admin/notices, the system shall render the paginated notice list.`
- **OPS-13-103:** `When an admin dismisses a notice, the system shall permanently delete the notice record from the database.`
- **OPS-13-104:** `When an admin triggers bulk dismissal, the system shall permanently delete every selected notice by ID.`
- **OPS-13-105:** `When a new notice is created after the admin dashboard was last viewed, the system shall display a notice badge on the admin navigation.`

### Unwanted Behaviour Requirements (Notice Errors)

- **OPS-13-301:** `If a non-admin attempts to view the notices page, then the system shall return 403 Forbidden.`
- **OPS-13-302:** `If a notice references an entity that has since been deleted, then the system shall render the notice with a "stale reference" indicator rather than failing.`

---

## 14. Doctor Diagnostic Checks

**User Story:** As an administrator, I want to run consistency checks and apply repairs so that I can detect and fix data corruption, broken references, and misconfiguration.

### Ubiquitous Requirements (Doctor Properties)

- **OPS-14-001:** `The system shall provide a doctor subcommand exposing a curated set of consistency checks.`
- **OPS-14-002:** `The system shall classify each check as either read-only (diagnostic) or repairing (writes changes).`
- **OPS-14-003:** `The system shall expose a non-exhaustive set of approximately 27 consistency checks including: authorized SSH keys, repository count integrity, dangling LFS objects, orphaned hooks, broken symlink references, stale archive caches, script-type availability, recalculate-stars-number, enable-push-options, check-git-daemon-export-ok, check-commit-graphs, recalculate-merge-bases, check-user-email, check-user-names, check-user-type, disable-mirror-actions-unit, fix-broken-repo-units, fix-owner-team-create-org-repo, synchronize-repo-heads, storage-attachments, storage-avatars, storage-packages, database version, paths consistency, DB consistency, and fix #16961 / #8312.`
- **OPS-14-004:** `The system shall provide a doctor recreate-table subcommand that recreates database tables from XORM definitions and copies data over, deleting old columns.`
- **OPS-14-005:** `The system shall provide a doctor convert subcommand that converts the database between storage configurations.`

### Event-Driven Requirements (Doctor Workflow)

- **OPS-14-101:** `When an admin runs the doctor check command without arguments, the system shall run all read-only checks and report any failures.`
- **OPS-14-102:** `When an admin runs doctor with a specific check name, the system shall execute only that check.`
- **OPS-14-103:** `When an admin runs doctor with the --fix flag, the system shall execute every repairing check and apply fixes.`
- **OPS-14-104:** `When an admin runs doctor check, the system shall accept the following flags: list, default, run, all, fix, log-file, and color.`
- **OPS-14-105:** `When a check completes successfully, the system shall print a summary and exit 0; on failure the system shall exit non-zero.`

### State-Driven Requirements (Doctor Safety)

- **OPS-14-701:** `While the system is running in --fix mode, the system shall log every applied change to the audit log and to the admin notices panel.`

### Optional Feature Requirements (Doctor Configuration)

- **OPS-14-201:** `Where a check is run via the CLI with --log-level=debug, the system shall emit per-item diagnostic output.`

### Unwanted Behaviour Requirements (Doctor Errors)

- **OPS-14-301:** `If a repairing check is run without --fix, then the system shall report the issues it would fix but shall not modify any data.`
- **OPS-14-302:** `If the database is unreachable during a doctor run, the system shall halt the run and report the failure.`
- **OPS-14-303:** `If a check fails due to an internal error, the system shall log the error, continue with remaining checks, and report the failure in the summary.`

---

## 15. CLI Command Surface

**User Story:** As an operator, I want a CLI entry point so that I can run Gitea as a server, perform admin operations, and execute maintenance tasks from the command line.

### Ubiquitous Requirements (CLI Properties)

- **OPS-15-001:** `The system shall expose a single gitea binary with subcommands including web, serv, admin, doctor, dump, hook, keys, migrate, migrate-storage, manager, embedded, dump-repository, restore-repository, actions, cert, generate, docs, and help.`
- **OPS-15-002:** `The system shall load configuration from custom/conf/app.ini before executing any subcommand.`
- **OPS-15-003:** `The system shall accept command-line flags and GITEA__ prefixed environment variables as configuration overrides.`

### Event-Driven Requirements (CLI Workflow)

- **OPS-15-101:** `When an operator runs gitea web, the system shall initialize configuration, database, services, and start the HTTP and (optionally) SSH servers.`
- **OPS-15-102:** `When the SSH server invokes gitea serv on an authenticated SSH session, the system shall execute the requested git-upload-pack or git-receive-pack as the authenticated user.`
- **OPS-15-103:** `When git invokes gitea hook pre-receive, update, or post-receive, the system shall execute the corresponding server-side hook pipeline.`
- **OPS-15-104:** `When an operator runs gitea admin user create, the system shall create a user account without requiring web access.`
- **OPS-15-105:** `When an operator runs gitea admin auth list/add/update/delete, the system shall manage external authentication sources from the command line.`
- **OPS-15-106:** `When an operator runs gitea dump, the system shall produce a backup archive per OPS-10.`
- **OPS-15-107:** `When an operator runs gitea migrate, the system shall execute pending database migrations without starting the server.`
- **OPS-15-108:** `When an operator runs gitea migrate-storage, the system shall transfer attachments, LFS objects, packages, or avatars between storage backends.`
- **OPS-15-109:** `When an operator runs gitea admin sendmail, the system shall dispatch an email message to all users from the command line.`
- **OPS-15-110:** `When an operator runs gitea admin regenerate hooks, the system shall regenerate repository hook files from the command line.`
- **OPS-15-111:** `When an operator runs gitea admin regenerate keys, the system shall regenerate the authorized_keys file from the command line.`
- **OPS-15-112:** `When an operator runs gitea manager with a subcommand (shutdown, restart, reload-templates, flush-queues, logging, processes), the system shall manage the running Gitea process.`

### Optional Feature Requirements (CLI Configuration)

- **OPS-15-201:** `Where the operator passes --custom-path or --config, the system shall load configuration from the specified location instead of the default.`
- **OPS-15-202:** `Where the operator passes --work-path, the system shall treat the specified directory as the Gitea working root.`

### Unwanted Behaviour Requirements (CLI Errors)

- **OPS-15-301:** `If configuration loading fails before subcommand dispatch, then the system shall exit non-zero with a descriptive error.`
- **OPS-15-302:** `If a subcommand receives invalid flags or arguments, then the system shall print usage and exit non-zero.`
- **OPS-15-303:** `If a subcommand requires database access and the database is unreachable, then the system shall report the connection failure and exit non-zero.`

---

## 16. Cron / Scheduled Tasks

**User Story:** As a system administrator, I want background tasks to run on schedule so that system maintenance is performed automatically.

### Ubiquitous Requirements (Cron Properties)

- **OPS-16-001:** `The system shall provide a configurable cron framework where each task has independent ENABLED, SCHEDULE, and RUN_AT_START settings.`
- **OPS-16-002:** `The system shall support standard cron syntax for scheduling each background task.`
- **OPS-16-003:** `The system shall expose a cron task monitoring interface in the admin panel.`

### Event-Driven Requirements (Cron Task Execution)

- **OPS-16-101:** `When the cron schedule for repository mirror sync fires (default every 10 minutes), the system shall sync all configured mirror repositories.`
- **OPS-16-102:** `When the cron schedule for repository health checks fires (default daily), the system shall run health checks on all repositories.`
- **OPS-16-103:** `When the cron schedule for repository stats update fires (default daily), the system shall recalculate repository statistics.`
- **OPS-16-104:** `When the cron schedule for archive cleanup fires (default daily), the system shall remove expired repository archives.`
- **OPS-16-105:** `When the cron schedule for zombie task cleanup fires (default every 5 minutes), the system shall stop tasks exceeding ZOMBIE_TASK_TIMEOUT.`
- **OPS-16-106:** `When the cron schedule for endless task cleanup fires (default every 30 minutes), the system shall stop tasks exceeding ENDLESS_TASK_TIMEOUT.`
- **OPS-16-107:** `When the cron schedule for abandoned job cleanup fires (default every 6 hours), the system shall cancel jobs exceeding ABANDONED_JOB_TIMEOUT.`
- **OPS-16-108:** `When the cron schedule for scheduled workflows fires, the system shall trigger all workflow runs due at that time.`
- **OPS-16-109:** `When an external user sync source (LDAP) is configured and its cron schedule fires (default daily), the system shall synchronize user data from the external source.`
- **OPS-16-110:** `When the cron schedule for cleanup_hook_task_table fires, the system shall prune stale webhook hook_task records.`
- **OPS-16-111:** `When the cron schedule for cleanup_packages fires (default daily, only when packages enabled), the system shall delete package versions older than the configured retention.`
- **OPS-16-112:** `When the cron schedule for deleted_branches_cleanup fires, the system shall remove deleted branch records.`
- **OPS-16-113:** `When the cron schedule for update_migration_poster_id fires (migrations enabled), the system shall fix migrated commit poster IDs.`
- **OPS-16-114:** `When the cron schedule for delete_inactive_accounts fires, the system shall delete users whose activation code expired beyond ActiveCodeLives.`
- **OPS-16-115:** `When the cron schedule for git_gc_repos fires, the system shall run git garbage collection on repositories.`
- **OPS-16-116:** `When the cron schedule for gc_lfs fires (LFS enabled), the system shall garbage-collect unassociated LFS meta objects older than one week.`
- **OPS-16-117:** `When the cron schedule for update_checker fires (default weekly), the system shall query the Gitea update endpoint for new releases.`
- **OPS-16-118:** `When the cron schedule for delete_old_actions fires, the system shall delete activity records older than one year.`
- **OPS-16-119:** `When the cron schedule for delete_old_system_notices fires, the system shall delete system notices older than one year.`
- **OPS-16-120:** `When the cron schedule for rebuild_issue_indexer fires, the system shall repopulate the issue search index from the database.`
- **OPS-16-121:** `When the cron schedule for cleanup_actions fires (default daily, actions enabled), the system shall delete expired Actions runs/tasks/logs older than the configured threshold (distinct from artifact retention).`

### Optional Feature Requirements (Cron Configuration)

- **OPS-16-201:** `Where RUN_AT_START is enabled for a cron task, the system shall execute that task once during system startup.`
- **OPS-16-202:** `Where a custom schedule is provided for a cron task, the system shall override the default schedule with the custom one.`

### Unwanted Behaviour Requirements (Cron Errors)

- **OPS-16-301:** `If a cron task execution fails, then the system shall log the error and continue scheduling subsequent executions.`
- **OPS-16-302:** `If a cron task is disabled via configuration, then the system shall skip all scheduled executions for that task.`

---

## 17. Indexer System

**User Story:** As an administrator, I want pluggable indexer backends with reliable queue processing so that search indexes stay up-to-date without degrading system performance.

> **Note:** Search UI for each domain is documented within that domain (repo/code search in Domain 03, issue/PR search in Domain 04, user/org search in Domains 01/02). This section documents indexer backends and configuration only.

### Ubiquitous Requirements (Indexer Properties)

- **OPS-17-001:** `The system shall provide separate indexer components for code (REPO_INDEXER), issues (ISSUE_INDEXER), and repository statistics (STATS_INDEXER); the STATS_INDEXER is DB-only with no backend choice and no configuration section.`
- **OPS-17-002:** `The system shall support Bleve and Elasticsearch as backends for the code indexer, and Bleve, Elasticsearch, Meilisearch, and db as backends for the issue indexer.`
- **OPS-17-003:** `The system shall use a worker pool with async queue processing for indexer updates.`
- **OPS-17-004:** `The system shall deduplicate indexer queue entries to avoid redundant indexing operations.`

### Event-Driven Requirements (Indexer Lifecycle)

- **OPS-17-101:** `When the system starts, the system shall initialize all enabled indexers and wait up to STARTUP_TIMEOUT for each to become ready.`
- **OPS-17-102:** `When a repository update event occurs, the system shall enqueue the relevant files or issues for re-indexing.`
- **OPS-17-103:** `When an indexer queue entry fails to process, the system shall log the error and retry the entry.`
- **OPS-17-104:** `When the system shuts down, the system shall drain the indexer queues gracefully before terminating.`
- **OPS-17-105:** `When an indexer is initialized for the first time (no existing index), the system shall automatically populate the index from scratch via populateRepoIndexer (code/stats) or PopulateIssueIndexer (issues); no admin reindex UI exists.`
- **OPS-17-106:** `When repository or issue update events occur after initial population, the system shall incrementally re-index only the affected content via the indexer queue.`

### State-Driven Requirements (Indexer Health)

- **OPS-17-701:** `While an external indexer backend is unreachable, the system shall queue indexing operations for retry and return degraded search results.`
- **OPS-17-702:** `While initial index population is in progress, the system shall track and report the population progress.`

### Optional Feature Requirements (Indexer Configuration)

- **OPS-17-201:** `Where REPO_INDEXER_ENABLED is false, the system shall skip code indexer initialization and disable code search functionality.`
- **OPS-17-202:** `Where ISSUE_INDEXER_TYPE is set to bleve, the system shall store the issue index on the local filesystem at ISSUE_INDEXER_PATH.`

### Unwanted Behaviour Requirements (Indexer Errors)

- **OPS-17-301:** `If the indexer queue exceeds its capacity, then the system shall log a warning and continue accepting operations with delayed indexing.`
- **OPS-17-302:** `If initial index population is interrupted, then the system shall re-run it automatically on the next indexer initialization.`

---

## 18. Sitemap

**User Story:** As a search engine, I want auto-generated XML sitemaps so that I can discover and index public content on the Gitea instance.

### Ubiquitous Requirements (Sitemap Properties)

- **OPS-18-001:** `The system shall generate XML sitemaps conforming to the sitemaps.org protocol.`
- **OPS-18-002:** `The system shall paginate sitemaps using the configured SITEMAP_PAGING_NUM page size.`
- **OPS-18-003:** `The system shall include last modified timestamps in sitemap entries.`
- **OPS-18-004:** `The system shall expose sitemaps for repositories and users.`

### Event-Driven Requirements (Sitemap Generation)

- **OPS-18-101:** `When a request is made to GET /explore/repos/sitemap-{idx}.xml, the system shall generate a sitemap page containing up to SITEMAP_PAGING_NUM public repositories.`
- **OPS-18-102:** `When a request is made to GET /explore/users/sitemap-{idx}.xml, the system shall generate a sitemap page containing up to SITEMAP_PAGING_NUM public users.`

### Unwanted Behaviour Requirements (Sitemap Constraints)

- **OPS-18-301:** `If a sitemap page would contain more than 50,000 URLs, then the system shall split the sitemap across multiple pages.`
- **OPS-18-302:** `If a sitemap page would exceed 50MB in size, then the system shall split the sitemap across multiple pages.`
- **OPS-18-303:** `If a resource is not publicly accessible, then the system shall exclude that resource from all sitemaps.`

---

## 19. Storage Backends

**User Story:** As an administrator, I want configurable storage backends so that file artifacts (avatars, packages, LFS, attachments, Actions logs/artifacts, repo archives) are stored reliably on local disk or object storage.

> **Note:** Storage is consumed by avatars (Domain 01 — Identity & Authentication), packages/LFS/attachments (Domain 06 — Packages & Releases), Actions artifacts/logs (Domain 05 — CI/CD & Automation), and repo archives (Domain 03 — Repository & Code). This section documents storage backend configuration and management only.

### Ubiquitous Requirements (Storage Properties)

- **OPS-19-001:** `The system shall support the following storage backend types: local (filesystem) and minio (S3-compatible object storage).`
- **OPS-19-002:** `The system shall allow independent storage configuration per artifact type via [storage.<type>] sections.`
- **OPS-19-003:** `The system shall support the following storage type overrides: attachments, lfs, avatars, repo-avatars, packages, actions_artifacts, and actions_log.`

### Event-Driven Requirements (Storage Operations)

- **OPS-19-101:** `When the system stores an artifact, the system shall use the storage backend configured for that artifact type, falling back to the global [storage] defaults.`
- **OPS-19-102:** `When an operator runs gitea migrate-storage, the system shall transfer all artifacts of the specified type from one storage backend to another.`
- **OPS-19-103:** `When SERVE_DIRECT is enabled on a MinIO storage backend, the system shall issue pre-signed direct-download URLs that bypass the Gitea server.`

### Optional Feature Requirements (Storage Configuration)

- **OPS-19-201:** `Where MinIO is configured as the storage backend, the system shall connect using MINIO_ENDPOINT, MINIO_ACCESS_KEY_ID, MINIO_SECRET_ACCESS_KEY, and MINIO_BUCKET parameters.`
- **OPS-19-202:** `Where MINIO_USE_SSL is enabled, the system shall use TLS for all MinIO connections.`

### Unwanted Behaviour Requirements (Storage Errors)

- **OPS-19-301:** `If the storage backend is unreachable during a write operation, the system shall return a storage error and not corrupt existing data.`
- **OPS-19-302:** `If the local filesystem has insufficient disk space, the system shall fail the write operation with a disk space error.`

---

## 20. Update Checker

**User Story:** As an administrator, I want Gitea to check for new versions so that I am notified when updates are available.

### Ubiquitous Requirements (Update Checker Properties)

- **OPS-20-001:** `The system shall compare the installed version against the latest available version from the configured endpoint.`
- **OPS-20-002:** `The system shall check for updates on a periodic schedule.`

### Event-Driven Requirements (Update Workflow)

- **OPS-20-101:** `When the update checker runs, the system shall query the configured endpoint URL for the latest version information.`
- **OPS-20-102:** `When a newer version is available, the system shall display a notification to administrators in the web interface.`
- **OPS-20-103:** `When the update checker is disabled via configuration, the system shall stop performing periodic version checks.`

### Optional Feature Requirements (Update Channels)

- **OPS-20-201:** `Where the update channel is set to stable, the system shall only report stable releases as available updates.`
- **OPS-20-202:** `Where the update channel is set to a pre-release channel (beta or dev), the system shall report pre-release versions as available updates.`
- **OPS-20-203:** `Where a custom endpoint URL is configured, the system shall query that URL instead of the default Gitea update endpoint.`

### Unwanted Behaviour Requirements (Update Errors)

- **OPS-20-301:** `If the update endpoint is unreachable, then the system shall log a warning and continue operation without update notifications.`
- **OPS-20-302:** `If the update endpoint returns an invalid response, then the system shall log a warning and skip update processing.`
- **OPS-20-303:** `If the installed version cannot be determined, then the system shall skip the update check.`

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

### [queue] / [queue.<name>] Section

Each named queue (e.g. `task`, `webhook`, `issue_indexer`) is configured under `[queue.<name>]`, with global defaults under `[queue]`:

- **TYPE**: Queue backend (`level`, `channel`, `redis`, `dummy`). Default `level` (LevelDB).
- **DATADIR**: LevelDB data directory path (relative to AppDataPath). Default `queues/common`.
- **CONN_STR**: Connection string for LevelDB or Redis backends.
- **LENGTH**: Max queue length before a channel queue blocks. Default 100000.
- **QUEUE_NAME**: Storage name suffix (db key / redis key). Default `_queue`.
- **SET_NAME**: Set name suffix for unique queues. Default `_unique`.
- **BATCH_LENGTH**: Number of items fetched per batch. Default 20.
- **MAX_WORKERS**: Worker pool cap. Default `CPU count / 2` (minimum 1).

### [cron] / [cron.<name>] Section

Each cron task is configured under `[cron.<task_name>]` with these keys:

- **ENABLED**: Whether the task runs on schedule (default per task).
- **RUN_AT_START**: Run once at system startup (default per task).
- **SCHEDULE**: Standard cron syntax schedule (default per task).
- **NOTICE_ON_SUCCESS**: Post a system notice on successful run (default per task).

### [metrics] Section

- **ENABLED**: Enable the Prometheus metrics endpoint at /metrics (default false).
- **TOKEN**: Bearer token required to access the metrics endpoint (empty = no auth).
- **ENABLED_ISSUE_BY_LABEL**: Expose per-label issue count metrics (default false).
- **ENABLED_ISSUE_BY_REPOSITORY**: Expose per-repository issue count metrics (default false).

### [indexer] Section

Code indexer and issue indexer infrastructure configuration. Search UI for each domain is documented within that domain; this section documents backends and configuration only.

- **REPO_INDEXER_ENABLED**: Enable the repository code indexer (default false).
- **REPO_INDEXER_TYPE**: Code indexer backend (`bleve`, `elasticsearch`; default `bleve`).
- **REPO_INDEXER_PATH**: Local filesystem path for the Bleve code index.
- **REPO_INDEXER_CONN_STR**: Elasticsearch connection string (when TYPE=elasticsearch).
- **REPO_INDEXER_INCLUDE**: Glob patterns of files to include in indexing.
- **REPO_INDEXER_EXCLUDE**: Glob patterns of files to exclude from indexing.
- **REPO_INDEXER_EXCLUDE_VENDORED**: Skip vendored files (default true).
- **REPO_INDEXER_REPO_TYPES**: Repository types to index (default `sources,forks,mirrors,templates`).
- **MAX_FILE_SIZE**: Maximum file size for code indexing in bytes (default 1048576 = 1MB).
- **STARTUP_TIMEOUT**: Timeout for indexer initialization (default 30s).
- **ISSUE_INDEXER_TYPE**: Issue indexer backend (`bleve`, `elasticsearch`, `meilisearch`, `db`; default `bleve`).
- **ISSUE_INDEXER_PATH**: Local filesystem path for the Bleve issue index.
- **ISSUE_INDEXER_CONN_STR**: Connection string for Elasticsearch or Meilisearch issue indexers.

### [storage] / [storage.*] Section

Storage backends for avatars, packages, LFS, attachments, Actions artifacts/logs, and repo archives. Storage is consumed by avatars (Domain 01), packages/LFS/attachments (Domain 06), Actions artifacts/logs (Domain 05), and repo archives (Domain 03).

Global defaults under `[storage]`:

- **STORAGE_TYPE**: Storage backend (`local`, `minio`; default `local`).
- **SERVE_DIRECT**: Serve files directly from storage backend via redirect (default false).
- **PATH**: Local filesystem path (for STORAGE_TYPE=local).
- **MINIO_ENDPOINT**: MinIO/S3 endpoint.
- **MINIO_ACCESS_KEY_ID**: MinIO/S3 access key.
- **MINIO_SECRET_ACCESS_KEY**: MinIO/S3 secret key.
- **MINIO_BUCKET**: MinIO/S3 bucket name.
- **MINIO_LOCATION**: MinIO/S3 region/location.
- **MINIO_USE_SSL**: Use SSL for MinIO/S3 connection (default true).
- **MINIO_CHECKSUM**: Enable checksum verification (default false).

Per-type overrides: `[storage.attachments]`, `[storage.lfs]`, `[storage.avatars]`, `[storage.repo-avatars]`, `[storage.packages]`, `[storage.actions_artifacts]`, `[storage.actions_log]`. Each supports the same keys as `[storage]`.

### [server] PProf Keys

- **ENABLE_PPROF**: Enable the pprof server on localhost:6060 (default false).
- **PPROF_DATA_PATH**: Directory for pprof data files (default `<AppWorkPath>/data/tmp/pprof`).

### [update_checker] Section

- **ENABLED**: Enable the periodic update checker (default true).
- **URL**: Endpoint URL for version checks (default `https://dl.gitea.com`).

---

## Business Rules

- **BR-08-001:** Only users with admin privilege may access /admin/* routes and /api/v1/admin/* endpoints
- **BR-08-002:** Configuration secrets (passwords, tokens, secret keys) are always shadowed in the configuration viewer
- **BR-08-003:** Static configuration changes require a server restart to take effect
- **BR-08-004:** Dynamic configuration changes take effect immediately without restart
- **BR-08-005:** Database migrations execute in strict sequential order by version number
- **BR-08-006:** The minimum database schema version is 70
- **BR-08-007:** Queue worker count adjustments are dynamic and do not require restart
- **BR-08-008:** Metrics endpoint requires ENABLED=true in [metrics] configuration
- **BR-08-009:** Health check endpoint (/api/healthz) is always available regardless of configuration
- **BR-08-010:** PProf server is disabled by default and must be explicitly enabled via [server].ENABLE_PPROF; when enabled it binds to localhost:6060 only
- **BR-08-011:** Backup archives contain all data by default; individual categories are excluded via flags
- **BR-08-012:** Log levels are hierarchical: TRACE < DEBUG < INFO < WARN < ERROR < FATAL
- **BR-08-013:** Authentication source order determines precedence when multiple sources match a user (see Domain 01 (Identity & Authentication) §External Authentication Sources for source types and features)
- **BR-08-014:** Ghost users are system accounts used for attribution preservation after user deletion
- **BR-08-015:** Cron tasks have independent schedules and can be enabled or disabled individually
- **BR-08-016:** Queue worker pool size defaults to CPU count / 2 with a maximum of 10 workers
- **BR-08-017:** The default `[queue].TYPE` is `level` (LevelDB); the `channel` backend is only the default for the legacy `[task]` shim, not for production queues
- **BR-08-018:** Indexer queues use deduplication to prevent redundant indexing of the same content
- **BR-08-019:** Sitemaps are generated on-demand (not pre-built) and respect access permissions
- **BR-08-020:** Each sitemap file contains at most 50,000 URLs and must not exceed 50MB
- **BR-08-021:** Storage backend configuration is global with per-type overrides; storage usage is documented in consuming domains
- **BR-08-022:** Update checker frequency is fixed (cron-driven); only the channel and endpoint are configurable

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
| Queue backend unavailable | Items queued in memory; error logged |
| Queue worker panic during processing | Worker recovered; item requeued |
| LevelDB queue data corruption | Error logged; fresh queue started |
| Metrics endpoint requested when disabled | Return 404 |
| PProf requested when not enabled | PProf server not started; no /debug/pprof/* routes served |
| Log file write failure | Fall back to console logging |
| LDAP sync during source outage | Log error, continue with cached user data |
| Last admin user deletion attempt | Deny deletion |
| Unadopted repo with invalid data | Skip directory during unadopted scan |
| Cron task execution failure | Error logged; next scheduled execution proceeds normally |
| External indexer unreachable during indexing | Queue operations for retry, log error |
| Indexer startup exceeds STARTUP_TIMEOUT | Log timeout, proceed without that indexer |
| Initial index population interrupted | Re-runs automatically on next indexer initialization |
| Sitemap page exceeds 50,000 URLs | Split across multiple sitemap pages |
| Storage backend unreachable during write | Return storage error; existing data not corrupted |
| Update check endpoint returns malformed JSON | Log warning, skip update notification |
| Update checker disabled | No periodic version checks performed |

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
- Cron task schedules fire within 1 second of the configured cron time
- Queue item processing latency is under 100ms per item for lightweight tasks
- Queue graceful shutdown completes all in-flight items without data loss
- Indexer queue processes updates without blocking the main application thread
- Full index population completes without data loss and reports progress
- Sitemaps conform to the sitemaps.org XML protocol specification
- Storage migrate-storage transfers data between backends without data loss
- Update checker does not block or delay server startup

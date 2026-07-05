# TODO.AI.md — mailadmin (mail)

Outstanding implementation tasks. Complete each item fully before removing it.
Dependencies are noted — resolve in dependency order, not list order.

## [ ] Implement src/mode/mode.go
Read: AI.md PART 6

Mode detection and dispatch logic. Priority: --mode CLI flag > MODE env var > default production.
Debug: --debug CLI flag > DEBUG env var > false. Create mode.Detect(), mode.IsProduction(),
mode.IsDevelopment() functions. No external deps.

## [ ] Implement src/signal/signal.go
Read: AI.md PART 7

OS signal handling. SIGTERM/SIGINT/SIGQUIT -> graceful shutdown. SIGUSR1 -> reopen log files.
SIGUSR2 -> dump status to stderr. Wire into main.go shutdown sequence.

## [ ] Implement src/service/service.go
Read: AI.md PART 7, PART 24

Service lifecycle management: --service install/uninstall/start/stop/status/restart.
Platform support: systemd (Linux), OpenRC (Linux), launchd (macOS), Windows SCM.
System user creation in UID/GID range 200-899, shell /sbin/nologin.
Never use reserved UIDs: 65534, 999, 998-980, 101-110, 170-179.
Uninstall: confirm, stop, delete all data/config/cache/logs/backups, remove user/group.

## [ ] Fix config.getRandomAvailablePort() in src/config/config.go
Read: AI.md PART 5

Port selection: random unused port in 64000-64999 range on first run.
Save selected port to server.yml so it persists across restarts.
On conflict: fail with clear error message. Support dual port "8090,8443" format.
Privileged port (<1024): prompt for sudo or fall back to 64xxx.

## [ ] Implement --status CLI flag
Read: AI.md PART 7

Print current daemon status (running/stopped, PID, port, uptime, version) and exit.
Read PID file to determine if running. Root PID: /var/run/webappsgo/mail.pid,
user PID: {data_dir}/mail.pid.

## [ ] Implement --daemon CLI flag
Read: AI.md PART 7

Run server in daemon mode (detach from terminal). Controlled by config daemonize setting
and --daemon flag. Modern service managers prefer foreground (daemonize: false default).

## [ ] Implement shell completions (--completion)
Read: AI.md PART 8

Shell completion generation for bash, zsh, and fish. --completion bash/zsh/fish flag.
Output to stdout. No external deps.

## [ ] Implement web frontend (HTML/CSS/JS)
Read: AI.md PART 16, PART 17

Embed static assets in binary using go:embed. Dark mode default with CSS custom properties.
Mobile-first, responsive. WCAG 2.1 AA. PWA support (service worker, manifest).
All user-facing strings use i18n keys — no hardcoded English in templates.
Admin panel completely isolated from public routes (no links, separate nav).

## [ ] Implement admin panel UI
Read: AI.md PART 16, PART 17

Setup wizard (9 steps per PART 17). Setup token: 32 hex chars, one-time use,
displayed ONCE on first-run console output. Admin panel at /{admin_path}/ (default: /admin/).
PathSecurityMiddleware MUST be first in middleware chain.

## [ ] Implement /healthz and /server/healthz endpoints
Read: AI.md PART 13, PART 14

/server/healthz: public, no auth. Returns JSON: {status, version, mode, uptime, checks}.
/healthz: optional root alias (enabled: false by default, same handler, never redirect).
/metrics: internal only, never public. Maintenance mode response documented in PART 5.

## [ ] Implement API v1 routes (/api/v1/)
Read: AI.md PART 13, PART 14, PART 15

All routes under /api/v1/ prefix. Plural nouns, lowercase, hyphens, no trailing slashes.
Error response format: {ok, error, message, details?}. CSRF on all cookie-auth state changes.
SameSite=Strict cookies. HSTS 2 years + subdomains + preload.
subtle.ConstantTimeCompare for all secret comparisons.

## [ ] Implement mail service management API
Read: AI.md PART 18, PART 19

Service status endpoints: GET /api/v1/services/{name}/status.
Start/stop/restart: POST /api/v1/services/{name}/{action}.
Manage: postfix, dovecot, spamassassin, clamav, amavis, opendkim, opendmarc.
If SMTP disabled -> mail feature COMPLETELY DISABLED (not queued).

## [ ] Implement virtual user management API
Read: AI.md PART 19

CRUD for VirtualUser: GET/POST /api/v1/users, GET/PUT/DELETE /api/v1/users/{id}.
Argon2id for password hashing (NEVER bcrypt or MD5). Input validation on all fields.
Parameterized queries only (no string concatenation in SQL). FQDN resolution per PART 15.

## [ ] Implement domain management API
Read: AI.md PART 19

CRUD for VirtualDomain: GET/POST /api/v1/domains, GET/PUT/DELETE /api/v1/domains/{id}.
SSL cert lookup order per PART 15. FQDN validation per PART 5 path rules.

## [ ] Implement alias and forwarding API
Read: AI.md PART 19

CRUD for Alias: GET/POST /api/v1/aliases, GET/PUT/DELETE /api/v1/aliases/{id}.
CRUD for Forward: GET/POST /api/v1/forwards, GET/PUT/DELETE /api/v1/forwards/{id}.
Transport table management: GET/POST /api/v1/transports.

## [ ] Implement mail queue management API
Read: AI.md PART 19, PART 20

GET /api/v1/queue: list queue entries with status.
DELETE /api/v1/queue/{id}: remove from queue.
POST /api/v1/queue/{id}/retry: force retry.
GET /api/v1/stats: MailStats aggregate.

## [ ] Implement DKIM key generation
Read: AI.md PART 20, PART 21

Generate DKIM key pairs for domains. Store private keys securely (700 perms, never in DB).
Expose public key as DNS TXT record suggestion. Per-domain key management.

## [ ] Implement built-in scheduler tasks
Read: AI.md PART 22, PART 23

Wire scheduler.go tasks to actual implementations:
- geoip_update: fetch sapics/ip-location-db via oschwald/maxminddb-golang
- blocklist_update: fetch and parse blocklist feeds
- log_rotation: rotate and compress logs per max_age/max_size config
- session_cleanup: purge expired sessions from database
- backup: snapshot database + config; verify backup before keeping; delete if verify fails
- ssl_renewal: check cert expiry; renew 7 days before expiry
- health_check: internal health check, enter/exit maintenance mode

## [ ] Fix Makefile — use casjaysdev/go:latest and correct flags
Read: AI.md PART 26

Replace golang:alpine with casjaysdev/go:latest in GO_DOCKER.
Add -e GOFLAGS=-buildvcs=false to all docker run invocations.
Required targets: build, test, clean, docker-build, release, lint.
Coverage output NEVER in project tree — always /tmp/{org}/{name}-XXXXXX/.
Source is in src/ subdirectory — build command must target ./src/... not ./...

## [ ] Create docker/Dockerfile
Read: AI.md PART 27

Multi-stage build: casjaysdev/go:latest builder + scratch or alpine runtime.
CGO_ENABLED=0, -buildvcs=false. Binary named mail. Volumes: /config/mailadmin and /data/mailadmin.
OCI labels required (org.opencontainers.image.*). HEALTHCHECK instruction required.
No USER directive if binary needs privileged startup then drops privileges.
No build-toolchain.yml — Go projects use casjaysdev/go:latest directly.

## [ ] Create CI/CD workflows
Read: AI.md PART 28

Creation order: security.yml first, then ci.yml, then release.yml last.
All third-party Actions pinned to full 40-char commit SHA (never tags).
container: casjaysdev/go:latest for all Go CI jobs.
No build-toolchain.yml for Go projects.
Add -buildvcs=false to all go build/go test in CI.
Do NOT create release.yml until all code is complete and make test passes.

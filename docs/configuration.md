# Configuration

The server merges `server.yml`, `yaml.d/*.yml` and `yaml.d/*.yaml` in byte-wise filename order, then `secrets.yml`; later files override earlier values. Nested maps merge, ordinary lists replace, and LDAP/OIDC provider entries merge by `name`. Environment-variable references are expanded before decoding. Config files are private (`0600`), and `yaml.d/` is private (`0700`). The application writes branding settings to `yaml.d/branding.yml` and uploaded images to `assets/`.

The loader accepts legacy settings in `server.yml` while area-file migration is supported. Prefer the owner file listed in the table below for new configuration. `server.yml` keeps node-specific listener, identity, database, cache, TLS, proxy, and cluster settings. `secrets.yml` keeps credentials and secret keys.

| Area file | Configuration sections |
|---|---|
| `yaml.d/branding.yml` | `server.branding`, `server.seo`, `web.ui` |
| `yaml.d/notifications.yml` | `server.notifications`, including outbound SMTP settings |
| `yaml.d/web.yml` | `web.*` except `web.ui`, plus CORS, CSRF, URL detection, compression, health and i18n settings |
| `yaml.d/pages.yml` | `server.pages` |
| `yaml.d/security.yml` | Security, rate limits, GeoIP, limits, token policy, and data settings |
| `yaml.d/accounts.yml` | Authentication, users, organizations, sessions, admin settings and admin path |
| `yaml.d/maintenance.yml` | Maintenance, scheduler, update, backup, restore and CLI settings |
| `yaml.d/app.yml` | Product-specific settings |

## Common server settings

| Key | Purpose | Default behavior |
|---|---|---|
| `server.app_name` | Display name | `mailadmin` |
| `server.mode` | Runtime mode | `production` |
| `server.timezone` | Time zone | `America/New_York` |
| `server.port` | One port, or quoted `HTTP,HTTPS` ports such as `"80,443"` | Random available port when omitted/invalid |
| `server.address` | Listen address | `[::]` |
| `server.fqdn` | Public host name | Empty until configured |
| `server.admin_path` | Server Admin URL segment | `administration` |
| `server.api_version` | Version segment for canonical API routes | `v1` |
| `server.healthz.root.enabled` | Mount the `/healthz` compatibility alias | `false` |
| `server.pidfile` | PID file behavior | Enabled |
| `server.daemonize` | Detach on manual startup | Disabled |
| `server.database.driver` | Database backend | `sqlite` |

The current server configuration model is defined in `src/config/config.go`. Database settings include driver, host, port, database name, username, password, URL/token, and connection-pool limits. Session settings separately define admin and user cookie names, maximum age, idle timeout, extension behavior, secure-cookie mode, HTTP-only, and SameSite policy.

Account recovery settings live under `server.auth` in `yaml.d/accounts.yml`. `password_reset_ttl` defaults to `1h`; `email_verify_ttl` defaults to `48h` and controls one-time email-verification links; `partial_session_ttl` defaults to `5m` for the single-use continuation between password verification and MFA completion. Regular-user registration settings live under `server.users`: `enabled` defaults to `false`, `registration.mode` defaults to `open`, `registration.require_email_verification` defaults to `true`, and `registration.allowed_domains` / `blocked_domains` restrict signup email domains. Public signup is available only when users are enabled and mode is `open`; when verification mail cannot be delivered, the account is auto-verified and activated. User API-token settings are `server.users.tokens.enabled` (default `true`), `max_per_user` (default `5`), and `expiration_days` (default `0`, meaning no default expiry). The token page and API are available only while multi-user mode is enabled. Users manage tokens at `/users/tokens`; each token is shown once, stored as a SHA-256 hash, and can be scoped to `global`, `read-write`, or `read`.

Branding and search metadata are configured under `server.branding` and `server.seo`; theme and accent are under `web.ui`. The Server Admin editor is `/server/{admin_path}/branding`. Images accept HTTPS URLs, data URIs up to 256 KiB, or filenames in `assets/`. Admin uploads support PNG, JPEG, WebP, ICO and safe SVG up to 2 MiB per file. The remote image cache uses validated HTTPS fetches, public-address checks, redirect revalidation, a 10 MiB cap and a 30-second timeout; cached copies refresh every 24 hours.

Outbound email is configured at `/server/{admin_path}/email`. SMTP host, port, username, TLS mode, and sender identity are stored in `yaml.d/notifications.yml`; the SMTP password is stored separately in `secrets.yml`. TLS modes are `auto`, `starttls`, `tls`, or `none`. A server FQDN is required. Startup performs an SMTP greeting and EHLO readiness check, and email notifications are disabled when the FQDN is missing or that check fails. The password is never included in API responses. WebUI toast settings are in `server.notifications.webui`: position accepts `top-right`, `top-left`, `bottom-right`, or `bottom-left`; `duration` defaults to 3 seconds for success/info toasts, and `warning_duration` defaults to 5 seconds. Set either duration to `0` for manual dismissal; error toasts are always manual.

Per-event email settings use `server.notifications.email.events.<event>` for `startup`, `shutdown`, `backup_complete`, `backup_failed`, `ssl_expiring`, `ssl_renewed`, `ssl_renewal_failed`, `login_alert`, `security_alert`, `scheduler_error`, `password_changed`, `token_regenerated`, `update_available`, `update_installed`, `update_failed`, and `update_staged`. Defaults follow AI.md; for example, `scheduler_error` and `backup_failed` default to `true`, while `startup`, `shutdown`, and `backup_complete` default to `false`. An omitted setting uses the event's built-in default. Each Server Admin can override applicable email defaults per event under their account's notifications page. Event preference support does not imply a sender exists yet; current dispatchers cover backup success/failure, SSL renewal failure, generic scheduled-task failure, startup, and shutdown.

Maintenance recovery settings use the `server.maintenance` object. `self_healing.retry_interval` defaults to `30s`, `self_healing.max_attempts` defaults to `0` (unlimited), `cleanup.disk_threshold` defaults to `90`, and `cleanup.log_retention_days` defaults to `7`. `notify.on_enter` and `notify.on_exit` default to `true`. The monitor checks database connectivity and disk pressure. Above the threshold, it removes expired non-active log files and app-scoped temporary files older than 24 hours, then applies `server.backup.retention` to the startup-resolved backup directory. Cleanup preserves encrypted archives when their password cannot be read and reports freed space and failures in the audit log. A positive `max_attempts` stops health retries until restart; exhaustion records an audit event and attempts one operator email when SMTP is ready. Database write and required-file failure detection, diagnostics, and configured entry/recovery notifications remain open work.

Security-report retention is configured with `server.security.reports.retention_days` in `yaml.d/security.yml`. It defaults to `365` days; values below one use the default. The daily `security_cleanup` scheduler task deletes only closed reports older than the configured retention period. Open reports and public acknowledgment entries are retained.

The `web.security` settings in `yaml.d/web.yml` define `report_url` (default `https://github.com/webappsgo/mailadmin/security/advisories/new`), `expires` (a configured ISO date/RFC3339 timestamp or a rolling one-year expiry), and `keyservers` (bare hostnames only; default `keys.openpgp.org`). Invalid keyserver hosts are ignored. An empty keyserver list disables publication. The active OpenPGP key is served at `/.well-known/pgp-key.asc`; configured keyservers are retried by the `keyserver_retry` scheduler task. The Server Admin keypair page is `/server/{admin_path}/security/keypair`.

A single `server.port` starts one listener; port `443` or `server.ssl.enabled: true` makes that listener HTTPS. A quoted pair such as `server.port: "8090,8443"` starts HTTP on the first port and HTTPS on the second. The HTTPS listener requires a configured certificate file or a working ACME configuration. HTTP-01 challenges use the HTTP listener, while TLS-ALPN-01 challenges use the HTTPS listener. Invalid port values fall back to a random port in the 64000–64999 range.

Rate limiting is always on and is configured under `server.rate_limit`. `read`, `write`, and `health` each accept `requests` and `window` (seconds); `strict` and `auth_failures` use the same shape; `global_burst` sets the shared per-minute ceiling. `auth.login`, `auth.password_reset`, and `auth.registration` hold their independent request windows. Invalid or omitted numeric values use the documented defaults. Persistent counters live in the server database, with in-process counters used when the database counter store is unavailable.

| Policy | Default requests | Default window |
|---|---:|---:|
| `read` | 120 | 60 seconds |
| `write` | 10 | 60 seconds |
| `health` | 120 | 60 seconds |
| `strict` | 10 | 900 seconds |
| `auth_failures` | 5 | 900 seconds |
| `global_burst` | 240 | 60 seconds |
| `auth.login` | 5 | 900 seconds |
| `auth.password_reset` | 3 | 3,600 seconds |
| `auth.registration` | 5 | 3,600 seconds |

TLS settings are under `server.ssl`; global Let's Encrypt settings use `server.ssl.letsencrypt.enabled`, `email`, `challenge`, `staging`, `dns_provider`, and a flat `dns_credentials` map. The legacy flat `lets_encrypt`, `email`, and `staging` settings remain readable. Panel-host ACME issuance requires a configured `server.fqdn`, contact email, and `server.security.encryption_key`; autocert cache, DNS account state, and certificate material are encrypted in `server.db`. HTTP-01 requires the configured listener or reverse proxy to accept HTTP on port 80; TLS-ALPN-01 requires HTTPS on port 443; DNS-01 supports wildcard certificates through any lego DNS provider. Plaintext DNS credentials in `server.yml` are atomically rewritten as `enc:v1:` values with mode `0600`; runtime configuration and API views keep them encrypted or masked. The Server Admin page and API are `/server/{admin_path}/ssl` and `/api/{api_version}/server/{admin_path}/ssl`; changes require a restart. The loader keeps the server encryption key and backup encryption password in `secrets.yml`, repairs that file to mode `0600`, and migrates those values out of `server.yml`. The backup password is stored with AES-256-GCM using the generated 32-byte server encryption key. Do not put credentials in source control.

Custom panel domains are configured under `server.features.custom_domains` and are available whenever the application database is configured. The legacy `enabled` value does not gate this activated feature. Defaults allow apex and subdomains, disallow wildcards, set limits of five domains per user and twenty per Org, require SSL, and set a 24-hour verification-token lifetime. Reserved suffixes and blocked patterns are rejected. `max_domains_per_user: 0` and `max_domains_per_org: 0` remove the corresponding limits. Pending ownership checks retry every 15 minutes, and unverified enrollments older than the verification TTL are cleaned up daily. Verified non-wildcard domains can obtain certificates through HTTP-01 on an HTTP listener or TLS-ALPN-01 on a TLS listener; the selected challenge port must be reachable by Let's Encrypt. The custom-domain SSL form or API accepts a lego provider name and flat map of provider environment variables for DNS-01. Credentials and certificate material are encrypted with `server.security.encryption_key`. Wildcard domains require DNS-01 and a matching wildcard certificate. With `require_ssl: true`, the domain stays pending until a certificate is active; disabling the requirement allows verified domains to activate without enabling TLS.

The application also discovers its public IPv4 address through an external lookup and its first eligible global IPv6 interface address at startup, then refreshes both every 12 hours. The latest values and check time appear in `/server/healthz` and provide DNS routing targets; DNS TXT ownership verification remains independent of these addresses.

## Environment values

`MODE` and `TZ` are read as fallbacks when `server.mode` and `server.timezone` are absent. `MAILADMIN_CUSTOM_DOMAINS_MAX_PER_USER` and `MAILADMIN_CUSTOM_DOMAINS_MAX_PER_ORG` override their custom-domain limits. YAML values also support environment expansion, for example `password: ${DB_PASSWORD}`. Database password, URL, and token values in `secrets.yml` override the corresponding values in the main config. Keep required secret values present in the service environment and restrict access to process configuration.

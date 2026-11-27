## Project description

mailadmin is a Go-based web administration panel for managing complete mail server infrastructure. It provides a unified interface for controlling and configuring Postfix, Dovecot, SpamAssassin, ClamAV, and related services. Instead of manually editing configuration files across multiple services, administrators use a clean web UI to manage mail services, virtual users, domains, aliases, spam filtering, and TLS/DKIM security. Handles both system Unix users and database-backed virtual users. Supports multiple domains, real-time mail queue monitoring, delivery statistics, and log viewing. Binary name: `mail`. Admin panel path: `/server/administration`.

## Project variables

project_name: mailadmin
project_org: webappsgo
internal_name: mail
internal_org: webappsgo
module: github.com/webappsgo/mail
binary: mail
language: Go
port: 64580
admin_path: administration
license: MIT

## Business logic

### Product scope & non-goals

mailadmin is a web administration panel that configures and controls an
existing mail stack (Postfix, Dovecot, SpamAssassin/Rspamd, ClamAV, and
related components) on the host it runs on. In scope: service
detection/control, virtual and system user management, domain/DKIM/DNS
management, aliases/forwarding/vacation/sieve, spam and virus filtering
configuration, TLS/DKIM/SPF/DMARC/SASL/Fail2ban configuration, mail queue
management, and delivery/spam statistics and log viewing.

mailadmin is a multi-tenant mail hosting platform: Organizations are the
tenancy unit, each Org owns one or more mail Domains (`VirtualDomain`) and
its own Users, who self-service within their Org's scope. Server Admin
retains full cross-tenant control over the underlying mail stack and every
tenant's data.

Non-goals: mailadmin is not a webmail client (no message reading/composing
UI) and not an MTA/IMAP/spam-filter reimplementation (it wraps and
configures Postfix/Dovecot/SpamAssassin/ClamAV rather than replacing them).
Clustering (PostgreSQL/MySQL backend) is for running the panel itself
across nodes managing the same mail stack, not a separate concern from
multi-tenancy.

### Roles & permissions

- **Server Admin** — manages the mailadmin panel itself (not a privileged OS
  user). Full access to service control, domain/user/alias/transport
  management, spam/virus configuration, TLS/DKIM setup, queue management,
  stats, and logs.
- **Primary Admin** — the first Server Admin created during the setup
  wizard; cannot be deleted.
- **Organization** — the tenancy unit (PART 35). Owns one or more mail
  Domains (`VirtualDomain`) and has member Users. An Org Owner (the user
  who created the Org, or a promoted member) can invite/remove members and
  manage the Org's domains, aliases, forwards, and DKIM/DNS settings. An Org
  Member manages the domains/mailboxes their Org owns but cannot manage Org
  membership or billing-equivalent settings.
- **Virtual/System Mail User (self-service, end user)** — a mailbox owner
  (System or Virtual User) belongs to exactly one Org and may log into a
  self-service view scoped to their own account and their Org's domains:
  change their own password, set/edit their own vacation (auto-reply)
  message, view their own quota usage, and — if an Org Owner/Member —
  manage mailboxes/aliases/forwards within their Org's own domain(s).
  Self-service users cannot see or modify another Org's domains, users,
  aliases, transports, service state, or system-wide settings.
- Server Admins and mail users are separate identities: Server Admin
  credentials (Argon2id) are unrelated to a mail user's Dovecot auth
  credential, even when the same person holds both.

### Data model & sensitivity

- **Credentials** (highest sensitivity): panel admin password hashes
  (Argon2id), virtual user password hashes (crypt, for Dovecot
  compatibility), DKIM private keys, mail relay/transport passwords, TLS
  private keys, database credentials. Never logged, never exposed via API
  responses, never included in exports.
- **PII**: email addresses, local-parts, vacation message bodies (may
  contain personal text), mail log entries (sender/recipient addresses,
  subjects are not stored but envelope data is).
- **Mail content**: mailadmin does not store or display message bodies; it
  only manages metadata (queue entries: from/recipients/size/status) needed
  for delivery troubleshooting.
- **Configuration data**: Postfix/Dovecot/SpamAssassin/ClamAV config file
  contents, service status, DNS records — sensitive in that they reveal
  infrastructure topology but not directly exploitable secrets on their own.
- **Setup token**: 32 hex chars (128-bit, crypto/rand), single-use, consumed
  on wizard completion — treated as a credential until consumed.
- **Tenancy data** (PARTS 34-36): Organization records (name, slug, owning
  user, custom domain + verification state), OrgMember rows (who belongs to
  which Org and with what role), and the `OrgID` foreign key on every
  `VirtualUser`/`VirtualDomain`/`Alias`/`Forward`/`Transport` row — this is
  the isolation boundary between tenants and MUST be enforced on every read
  and write, not just at the UI layer.
- **Custom-domain verification data** (PART 36): TXT ownership-verification
  tokens and the ACME account/certificate material issued for an Org's
  custom panel domain — treated as credential-adjacent, since a forged
  verification or stolen cert lets an attacker impersonate the panel under
  that domain.

### Trust boundaries & external services

- **Local mail services (Postfix, Dovecot, SpamAssassin, Rspamd, ClamAV,
  OpenDKIM, Fail2ban, etc.)** — trusted, but only reachable/controllable via
  their config files, control sockets, and systemd/OpenRC/launchd service
  management; mailadmin must validate any config it writes before reload/
  restart, since a bad write can take mail flow down.
- **Local database** (SQLite for single-server, PostgreSQL/MySQL for
  cluster) — trusted backend for virtual users/domains/aliases/forwards/
  transports; holds password hashes and DKIM keys, so it is a primary
  protected asset.
- **Let's Encrypt / ACME** — external, semi-trusted service used only for
  certificate issuance/renewal over the standard ACME protocol; failure mode
  is falling back to the existing certificate (self-signed or previously
  issued) and surfacing a renewal-failure alert, never serving without TLS
  where TLS was previously configured.
- **DNS** — untrusted/informational only; DNS Helper reads/displays required
  records (MX, SPF, DKIM, DMARC) for the admin to publish manually or via
  their own DNS provider. mailadmin does not assume it has authority to
  modify third-party DNS zones.
- **Alert delivery (email/webhook notifications)** — outbound only; webhook
  endpoints are admin-supplied and treated as untrusted destinations (no
  secrets beyond the alert payload are sent).
- **Other tenant Organizations** (PARTS 34-36) — every other Org and its
  Users, Domains, aliases, forwards, and transports are untrusted from the
  perspective of a given Org's session; the panel's own authorization layer
  is the sole boundary, since all tenants share one database and one mail
  stack.
- **Custom panel domain / ACME (PART 36)** — an Org-supplied custom domain
  is untrusted until TXT-record ownership verification succeeds; only after
  verification does mailadmin request a Let's Encrypt certificate for it
  and route panel traffic under it. Failure mode on verification or ACME
  failure is: keep serving the panel under mailadmin's own domain/host and
  surface the failure to the Org, never silently accept an unverified
  domain.

### Threat model & abuse cases

- **Primary assets protected**: virtual/system user password hashes, DKIM
  private keys, TLS private keys, database credentials, panel admin
  credentials, custom-domain ACME certificates/verification tokens, and the
  integrity of live mail service configuration (Postfix/Dovecot/etc.) and
  of tenant (Org) data isolation — corrupting these can cause credential
  compromise, mail interception/spoofing, cross-tenant data exposure, or
  full mail outage.
- **Trusted inputs**: authenticated Server Admin actions via the panel UI/
  API; authenticated self-service mail-user/Org-member actions scoped to
  their own account and their own Org's `OrgID`.
- **Untrusted inputs**: setup wizard token (until consumed), all
  panel-facing HTTP request data before auth, CSV import content, inbound
  mail queue metadata surfaced from Postfix, ACME/Let's Encrypt responses,
  Org-supplied custom domain names and their TXT verification state, Org
  invite/self-registration submissions, and any config values a self-service
  user submits (vacation message body, new password).
- **Attacker/abuser goals**: obtain or escalate panel/mailbox/Org
  credentials, inject malicious config to reroute or intercept mail, forge
  DKIM/SPF/DMARC posture to enable spoofing, exhaust the mail queue or
  spam/virus scanners (DoS), use CSV import or vacation-message fields to
  inject config/command payloads, use self-service or Org-member access to
  pivot into admin-level or another-tenant's-data access (privilege
  escalation / cross-tenant IDOR on user/domain/Org IDs), or claim a custom
  domain they do not own to hijack panel traffic or intercept another Org's
  ACME issuance.
- **Project-specific abuse cases and defenses**:
  - Credential stuffing / brute force against panel or self-service login →
    Argon2id hashing, rate limiting, Fail2ban integration where available.
  - Privilege escalation via self-service endpoints acting on another
    user's ID → every self-service API call MUST verify the target record
    belongs to the authenticated identity (no trusting client-supplied IDs
    alone).
  - Cross-tenant IDOR — an Org member requesting another Org's
    domain/user/alias/forward/transport/stats by ID → every query MUST
    filter by the authenticated session's `OrgID` server-side, never rely on
    the client to only request its own Org's IDs; Server Admin routes are
    the only ones exempt, and only when explicitly acting as Server Admin.
  - Org membership escalation (a Member granting itself Owner, or an
    outsider joining an Org without invite) → membership/role changes
    require the acting user to already be Owner of that specific Org; PART
    35 registration/invite rules govern how a User joins an Org.
  - Custom-domain hijack/impersonation (PART 36) — an Org claiming a domain
    it does not control to redirect panel traffic or intercept a
    certificate → TXT-record ownership verification MUST succeed before
    ACME issuance or traffic routing for that domain begins, and
    verification is re-checked before certificate renewal.
  - Config injection via generated Postfix/Dovecot/SpamAssassin config
    (domain names, aliases, transport values, vacation messages) → all
    values are validated/escaped before being written into service config
    files; config is validated before reload/restart is triggered.
  - Malicious CSV import (formula injection, path traversal in filenames,
    oversized files) → CSV import is size-capped, parsed as data only
    (never executed), and every row re-validated with the same rules as the
    single-record create path, including that imported rows resolve only to
    the importing Org.
  - Spam/virus scanning bypass or resource exhaustion → greylisting, RBL/
    DNSBL, and quota/size limits are configuration mailadmin can surface and
    enforce, not features it must reimplement.
  - Mail queue tampering (unauthorized hold/delete/flush) → destructive
    queue actions require Server Admin auth and are logged.

### Security decisions & exceptions

- mailadmin runs with the permissions needed to read/write Postfix/Dovecot/
  SpamAssassin/ClamAV config files and control their services via systemd/
  OpenRC/launchd — this is an intentional, documented exception to
  least-privilege at the OS level, required because the product's purpose
  is to administer those services. It MUST NOT run mail services themselves
  as root beyond what each service's own packaging already requires.
  Self-service end-user actions are still scoped in the application layer
  even though the underlying process has broad file/service access.
  - The setup wizard's DNS Helper displays third-party-domain records but
  never writes to external DNS zones directly — no DNS provider API
  credentials are requested or stored, avoiding a new class of
  externally-scoped secret.
- **Multi-tenant single database/process, not per-tenant isolation**
  (PARTS 34-36) — all Orgs share one database and one mailadmin process
  instead of per-tenant containers/databases; this is an intentional
  scaling/simplicity tradeoff, made acceptable only because `OrgID`
  scoping is enforced in the application/query layer on every tenant-owned
  table. A future move to stronger isolation (e.g. per-tenant schemas) is
  not required by this spec, but the query-layer scoping requirement is
  non-negotiable as long as isolation stays application-enforced.
- **Custom panel domains issue real certificates for Org-controlled names**
  (PART 36) — automatic ACME issuance for an Org's verified custom domain
  is an intentional exception to "mailadmin doesn't act on external DNS":
  it still never modifies the Org's DNS, but it does obtain and hold a
  certificate for a name mailadmin does not own, which is acceptable only
  because issuance is gated on TXT ownership verification succeeding first.

### Core Features

#### Service Management
- Component Detection: Auto-detect installed mail services (Postfix, Dovecot, SpamAssassin, etc.)
- Service Control: Start, stop, restart, reload mail services via systemd/OpenRC/launchd
- Configuration Editor: Web-based config file editing with validation
- Template System: Pre-configured templates for common setups
- Backup/Restore: Save and restore complete mail stack configurations

#### User Management
- System Users: Manage Unix system users with email access
- Virtual Users: Database-backed virtual mailboxes (no shell access)
- User Mapping: Map between system users and email identities
- Bulk Operations: Import/export users from CSV
- Password Management: Change passwords, enforce policies (Argon2id for panel, crypt for Dovecot)

#### Domain Management
- Virtual Domains: Host multiple email domains on one server
- Domain Aliases: Mirror one domain to another
- Catchall Addresses: Catch-all inbox for undefined addresses
- Domain DKIM: Per-domain DKIM signing keys (generate, rotate, export DNS records)
- DNS Helper: Show required DNS records (MX, SPF, DKIM, DMARC)

#### Alias & Forwarding
- Email Aliases: Create aliases pointing to local or remote addresses
- Forwarding Rules: Forward mail to multiple destinations
- Vacation Messages: Auto-reply configuration
- Sieve Filters: Advanced mail filtering rules (if Dovecot supports)

#### Spam & Virus Protection
- SpamAssassin: Configure spam detection rules and scores
- Rspamd: Alternative modern spam filtering
- Pyzor/Razor: Collaborative spam detection networks
- Amavisd-new: Content filter integration
- ClamAV: Virus scanning for attachments
- Greylisting: Temporary rejection for spam reduction
- Blocklists: DNS-based blocklists (RBL/DNSBL)
- Whitelists: Trusted sender lists

#### Security & Authentication
- TLS/SSL: Configure TLS for SMTP, IMAP, POP3
- Let's Encrypt: Automatic certificate generation and renewal
- DKIM: DomainKeys Identified Mail signing
- SPF: Sender Policy Framework validation
- DMARC: Domain-based Message Authentication
- SASL: SMTP authentication configuration
- Fail2ban: Brute force protection (if installed)

#### Monitoring & Logs
- Mail Queue: View and manage Postfix queue (list, hold, delete, flush)
- Delivery Stats: Success/failure rates, bounce tracking
- Spam Statistics: Detection rates, false positives
- Service Status: Real-time health of all components
- Log Viewer: Search and filter mail logs
- Alerts: Email/webhook notifications for issues

### Managed Components

#### Core Mail Services
- Postfix (MTA): main.cf, master.cf, virtual, transport, aliases, header_checks, body_checks, access maps
- Dovecot (IMAP/POP3): dovecot.conf, conf.d/*, protocols, auth, ssl, namespace

#### Spam & Content Filtering
- SpamAssassin: local.cf, rules, bayes database
- Rspamd: rspamd.conf, modules.d/*
- Amavisd-new: amavisd.conf
- Pyzor: servers configuration
- Razor: razor-agent.conf
- ClamAV: clamd.conf, freshclam.conf

#### Optional Components
- OpenDKIM: opendkim.conf, key files
- Fail2ban: jail.local (mail sections)
- Policyd: policyd.conf
- Postgrey: postgrey.conf

### Data Models

#### Organization
```go
type Organization struct {
    ID              int64
    Name            string
    Slug            string    // URL-safe identifier
    OwnerUserID     int64     // Regular User who owns the Org
    CustomDomain    string    // Optional — panel UI served under this domain (PART 36)
    CustomDomainVerified bool
    Active          bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

#### OrgMember
```go
type OrgMember struct {
    ID        int64
    OrgID     int64
    UserID    int64
    Role      string    // "owner" | "member"
    CreatedAt time.Time
}
```

#### VirtualUser
```go
type VirtualUser struct {
    ID           int64
    OrgID        int64     // Org this mailbox belongs to
    Email        string    // Full email address
    LocalPart    string    // Username part
    Domain       string    // Domain part
    PasswordHash string    // Argon2id for panel, crypt for Dovecot
    Quota        int64     // Bytes (0 = unlimited)
    Enabled      bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

#### VirtualDomain
```go
type VirtualDomain struct {
    ID             int64
    OrgID          int64     // Owning Organization
    Domain         string
    Active         bool
    DKIMEnabled    bool
    DKIMSelector   string
    DKIMPrivateKey string
    DKIMPublicKey  string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

#### Alias
```go
type Alias struct {
    ID          int64
    Source      string
    Destination string
    Domain      string
    Active      bool
    CreatedAt   time.Time
}
```

#### Forward
```go
type Forward struct {
    ID           int64
    Source       string
    Destinations []string
    KeepLocal    bool
    Active       bool
    CreatedAt    time.Time
}
```

#### Transport
```go
type Transport struct {
    ID        int64
    Domain    string
    Transport string
    NextHop   string
    Username  string
    Password  string
    Active    bool
}
```

#### QueueEntry (read-only, from Postfix)
```go
type QueueEntry struct {
    ID          string
    Size        int64
    From        string
    Recipients  []string
    ArrivalTime time.Time
    Status      string
    Reason      string
}
```

#### ServiceStatus
```go
type ServiceStatus struct {
    Name    string
    Running bool
    Enabled bool
    PID     int
    Uptime  time.Duration
    Version string
}
```

#### MailStats
```go
type MailStats struct {
    Period       string
    Sent         int64
    Received     int64
    Bounced      int64
    Rejected     int64
    SpamBlocked  int64
    VirusBlocked int64
    QueueSize    int64
}
```

### API routes

Project API routes are versioned under `/api/{api_version}`. Server Admin access
uses an operator Bearer token; tenant operations use a user or organization
Bearer token and are restricted to its tenant. List requests accept `limit`,
`offset`, `sort`, `order`, and resource-specific filters where applicable.
Success responses contain the requested resource or collection; failures use the
standard API error body defined by AI.md PART 14.

| Method | Path | Auth | Purpose; request and response summary |
|---|---|---|---|
| GET, POST | `/api/{api_version}/mail/users` | user, org, or operator token | List mailboxes with filters and pagination, or create a mailbox with organization, domain, address, password, name, quota, and enabled state; returns mailbox data without credentials. |
| GET, PUT, DELETE | `/api/{api_version}/mail/users/{id}` | user, org, or operator token | Read, update, or delete one mailbox in the caller's permitted tenant scope; update accepts mailbox fields and returns the credential-free mailbox. |
| POST | `/api/{api_version}/mail/users/import` | operator token | Import mailbox rows from CSV; returns per-row results and counts. |
| GET | `/api/{api_version}/mail/users/export` | operator token | Export credential-free mailbox data as CSV. |
| GET, POST | `/api/{api_version}/mail/domains` | user, org, or operator token | List domains in scope with `limit`, `offset`, `sort` (`id` or `domain`), and `order` (`asc` or `desc`) pagination, or create a domain with organization, domain name, and description; list responses include pagination metadata. User tokens select an organization with `org_id`; organization tokens are fixed to their organization. |
| GET, PUT, DELETE | `/api/{api_version}/mail/domains/{id}` | user, org, or operator token | Read, update, or delete a scoped mail domain; requests include `org_id` for user/operator tokens and update accepts domain name and description. Deletion returns a conflict while mailboxes, aliases, DKIM, SPF, DMARC, or transport records still refer to the domain. |
| POST | `/api/{api_version}/mail/domains/{id}/dkim/keys` | operator token | Generate and persist a 2048-bit DKIM key for the domain selected by `org_id`; returns the selector, public key, and DNS name, never the private key. |
| GET | `/api/{api_version}/mail/domains/{id}/dkim/dns` | user, org, or operator token | Return the active DKIM TXT record for the domain selected by `org_id`; an unconfigured key is returned as `null`. |
| GET | `/api/{api_version}/mail/domains/{id}/dns` | user, org, or operator token | Return the domain's MX record when the server FQDN is configured, SPF and DMARC guidance (including stored policies), and the active DKIM TXT record when configured. |
| PATCH | `/api/{api_version}/mail/domains/{id}/dns` | user, org, or operator token | Update the domain's SPF value and DMARC policy, subdomain policy, percentage, and aggregate report address; returns the refreshed DNS record list. |
| GET, POST | `/api/{api_version}/mail/aliases` | user, org, or operator token | List scoped aliases or create an alias with source, destinations, and owning domain; returns aliases without secret material. |
| GET, PUT, DELETE | `/api/{api_version}/mail/aliases/{id}` | user, org, or operator token | Read, update, or delete a scoped alias; update accepts source and destination fields. |
| GET, POST | `/api/{api_version}/mail/forwards` | user, org, or operator token | List forwarding rules in the selected organization or create a rule with `org_id`, `mailbox_id`, destination addresses, local-copy preference, and enabled state. |
| PUT, DELETE | `/api/{api_version}/mail/forwards/{id}` | user, org, or operator token | Update or delete an organization-scoped forwarding rule; update accepts the mailbox, destination addresses, local-copy preference, and enabled state. |
| GET | `/api/{api_version}/mail/queue` | operator token | List Postfix queue entries with `limit`, `offset`, and `q` filtering; returns envelope metadata, size, arrival time, status, and reason, never message contents. |
| POST | `/api/{api_version}/mail/queue/flushes` | operator token | Flush the mail queue; returns the operation result. |
| PATCH | `/api/{api_version}/mail/queue/{id}` | operator token | Hold or release a queued message; accepts `{"held":true|false}` and returns the resulting hold state. |
| DELETE | `/api/{api_version}/mail/queue/{id}` | operator token | Delete a queued message; returns the operation result. |
| POST | `/api/{api_version}/mail/queue/{id}/retries` | operator token | Retry delivery of a queued message; returns the operation result. |
| GET | `/api/{api_version}/mail/services` | operator token | Return detected mail service states, versions, and uptime. |
| PATCH | `/api/{api_version}/mail/services/{name}` | operator token | Control an allowlisted installed mail service; the body contains exactly one of `{"running":true|false}`, `{"enabled":true|false}`, or `{"reload":true}` and the response returns the resulting service state. Service changes require operating-system privileges and are audited. |
| GET | `/api/{api_version}/mail/stats?period={period}` | operator token | Return sent/received/bounced/deferred and byte totals, current Postfix queue size, and spam/virus counts for `24h`, `7d`, or `30d`; the default period is `24h`. |
| GET | `/api/{api_version}/mail/stats/spam?period={period}` | operator token | Return spam and virus counts for `24h`, `7d`, or `30d`; the default period is `24h`. |
| GET | `/api/{api_version}/mail/logs` | operator token | Return delivery-log envelope metadata (queue ID, sender, recipient, size, status, delay, DSN, and timestamp), filtered by optional `q` and paginated with `limit` (default 50, maximum 100) and `offset`. Message contents are never returned. |
| GET, POST | `/api/{api_version}/mail/transports` | operator token | List transport rules or create one with `domain`, `transport`, `next_hop`, optional username/password, and active state; responses report password presence without exposing the encrypted credential. |
| PUT, DELETE | `/api/{api_version}/mail/transports/{id}` | operator token | Update or delete a transport rule; omitting `password` preserves it and `clear_password` explicitly removes it. |
| GET, POST | `/api/{api_version}/orgs/{slug}/mail/domains` | organization or operator token | List or create mail domains owned by the organization; returns scoped domain data. |
| GET, PUT, DELETE | `/api/{api_version}/orgs/{slug}/mail/domains/{id}` | organization or operator token | Read, update, or delete an owned mail domain; delete returns a conflict while dependent mail records remain. |
| GET, POST | `/api/{api_version}/orgs/{slug}/aliases` | organization or operator token | List or create organization aliases; returns scoped alias data. |
| GET, PUT, DELETE | `/api/{api_version}/orgs/{slug}/aliases/{id}` | organization or operator token | Read, update, or delete one organization alias; request bodies contain its domain, source, and destination. |
| GET, POST | `/api/{api_version}/orgs/{slug}/mailboxes` | organization or operator token | List or create organization mailboxes; returns mailbox data without credentials. |
| GET, PUT, DELETE | `/api/{api_version}/orgs/{slug}/mailboxes/{id}` | organization or operator token | Read, update, or delete one organization mailbox; updates accept address, name, quota, enabled state, and optional password. |

### Setup Wizard

First-run wizard walks through:
1. Admin account creation (username, password, email)
2. Server FQDN configuration
3. Authentication mode selection (system users, virtual users, both)
4. Database selection (SQLite for single-server, PostgreSQL/MySQL for cluster)
5. First domain setup
6. TLS certificate setup (self-signed, existing cert, or Let's Encrypt)
7. Mail component detection (auto-scan for installed services)
8. Basic Postfix/Dovecot configuration generation
9. DKIM setup for first domain
10. Service start and health verification

Setup token is displayed on first run and consumed on wizard completion. Token is 32 hex chars (128-bit random, crypto/rand).

### Configuration (server.yml)

```yaml
server:
  port: 64580
  admin_path: admin

mail:
  mailbox_base: /var/vmail
  vmail_user: vmail
  vmail_uid: 5000
  vmail_gid: 5000
  auth_mode: virtual
  database:
    driver: postgres
    host: localhost
    port: 5432
    name: mailserver
    username: mailuser
    password: ${DB_PASSWORD}
  postfix:
    config_dir: /etc/postfix
    queue_dir: /var/spool/postfix
  dovecot:
    config_dir: /etc/dovecot
  spamassassin:
    config_dir: /etc/mail/spamassassin
  clamav:
    config_dir: /etc/clamav
```

### Package Mapping

#### Alpine Linux (apk)
postfix, dovecot, dovecot-pigeonhole-plugin, spamassassin, amavisd-new, pyzor, razor, rspamd, clamav, clamav-daemon, opendkim, opendkim-utils, fail2ban

#### Debian/Ubuntu (apt)
postfix, postfix-mysql, postfix-pgsql, dovecot-core, dovecot-imapd, dovecot-pop3d, dovecot-lmtpd, dovecot-managesieved, spamassassin, amavisd-new, pyzor, razor, rspamd, clamav, clamav-daemon, opendkim, opendkim-tools, fail2ban

#### RHEL/CentOS/Rocky/Alma (dnf)
postfix, dovecot, dovecot-mysql, dovecot-pgsql, dovecot-pigeonhole, spamassassin, amavisd-new, pyzor, razor, rspamd, clamav, clamav-server, opendkim, fail2ban

#### FreeBSD (pkg)
postfix, dovecot, dovecot-pigeonhole, spamassassin, amavisd-new, pyzor, razor, rspamd, clamav, opendkim, fail2ban

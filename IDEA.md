## Project description

mailadmin is a Go-based web administration panel for managing complete mail server infrastructure. It provides a unified interface for controlling and configuring Postfix, Dovecot, SpamAssassin, ClamAV, and related services. Instead of manually editing configuration files across multiple services, administrators use a clean web UI to manage mail services, virtual users, domains, aliases, spam filtering, and TLS/DKIM security. Handles both system Unix users and database-backed virtual users. Supports multiple domains, real-time mail queue monitoring, delivery statistics, and log viewing. Binary name: `mail`. Admin panel path: `/admin`.

## Project variables

project_name: mailadmin
project_org: webappsgo
internal_name: mail
internal_org: webappsgo
module: github.com/webappsgo/mail
binary: mail
language: Go
port: 64580
admin_path: admin
license: MIT

## Business logic

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

#### VirtualUser
```go
type VirtualUser struct {
    ID           int64
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

### API Endpoints

All API endpoints are versioned under `/api/v1/`.

#### Health & Info
- `GET /healthz` — Health check (not versioned)
- `GET /api/v1/server/about` — Server information
- `GET /api/v1/server/version` — Version information

#### Virtual Users
- `GET /api/v1/mail/users` — List virtual users
- `POST /api/v1/mail/users` — Create virtual user
- `GET /api/v1/mail/users/{id}` — Get virtual user
- `PUT /api/v1/mail/users/{id}` — Update virtual user
- `DELETE /api/v1/mail/users/{id}` — Delete virtual user
- `POST /api/v1/mail/users/import` — Import users from CSV
- `GET /api/v1/mail/users/export` — Export users to CSV

#### Domain Management
- `GET /api/v1/mail/domains` — List domains
- `POST /api/v1/mail/domains` — Add domain
- `GET /api/v1/mail/domains/{id}` — Get domain
- `PUT /api/v1/mail/domains/{id}` — Update domain
- `DELETE /api/v1/mail/domains/{id}` — Delete domain
- `POST /api/v1/mail/domains/{id}/dkim/generate` — Generate DKIM key
- `GET /api/v1/mail/domains/{id}/dkim/dns` — Get DKIM DNS records
- `GET /api/v1/mail/domains/{id}/dns` — Get all DNS records (MX, SPF, DKIM, DMARC)

#### Aliases
- `GET /api/v1/mail/aliases` — List aliases
- `POST /api/v1/mail/aliases` — Create alias
- `PUT /api/v1/mail/aliases/{id}` — Update alias
- `DELETE /api/v1/mail/aliases/{id}` — Delete alias

#### Forwarding
- `GET /api/v1/mail/forwards` — List forwarding rules
- `POST /api/v1/mail/forwards` — Create forwarding rule
- `PUT /api/v1/mail/forwards/{id}` — Update forwarding rule
- `DELETE /api/v1/mail/forwards/{id}` — Delete forwarding rule

#### Mail Queue
- `GET /api/v1/mail/queue` — List queue entries
- `POST /api/v1/mail/queue/flush` — Flush entire queue
- `POST /api/v1/mail/queue/{id}/hold` — Hold message
- `POST /api/v1/mail/queue/{id}/delete` — Delete message
- `POST /api/v1/mail/queue/{id}/retry` — Retry delivery

#### Service Control
- `GET /api/v1/mail/services` — List service statuses
- `POST /api/v1/mail/services/{name}/start` — Start service
- `POST /api/v1/mail/services/{name}/stop` — Stop service
- `POST /api/v1/mail/services/{name}/restart` — Restart service
- `POST /api/v1/mail/services/{name}/reload` — Reload service config

#### Statistics
- `GET /api/v1/mail/stats` — Mail delivery statistics
- `GET /api/v1/mail/stats/spam` — Spam statistics
- `GET /api/v1/mail/logs` — Mail log viewer

#### Transport
- `GET /api/v1/mail/transports` — List transport rules
- `POST /api/v1/mail/transports` — Create transport rule
- `PUT /api/v1/mail/transports/{id}` — Update transport rule
- `DELETE /api/v1/mail/transports/{id}` — Delete transport rule

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

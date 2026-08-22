# mailadmin

[![Build](https://img.shields.io/github/actions/workflow/status/webappsgo/mailadmin/ci.yml?branch=main)](https://github.com/webappsgo/mailadmin/actions)
[![Release](https://img.shields.io/github/v/release/webappsgo/mailadmin)](https://github.com/webappsgo/mailadmin/releases)
[![License](https://img.shields.io/github/license/webappsgo/mailadmin)](LICENSE.md)

## About

mailadmin is a Go-based web administration panel for managing complete mail
server infrastructure — Postfix, Dovecot, SpamAssassin, ClamAV, and related
services — through a single web UI instead of hand-editing config files
across multiple services. It handles both system Unix users and
database-backed virtual users, and supports multiple mail domains,
real-time queue monitoring, delivery statistics, and log viewing.

mailadmin is a multi-tenant mail hosting platform: Organizations are the
tenancy unit, each Organization owns one or more mail domains and its own
users, who self-service within their Organization's scope. The Server
Admin retains full cross-tenant control over the underlying mail stack.

## Features

### Service Management
- Auto-detect installed mail services (Postfix, Dovecot, SpamAssassin, etc.)
- Start, stop, restart, reload mail services via systemd/OpenRC/launchd
- Web-based config file editor with validation
- Pre-configured templates for common setups
- Backup/restore complete mail stack configurations

### Organizations & Users
- Organizations own one or more mail domains and their own users
- System Users: Unix system users with email access
- Virtual Users: database-backed virtual mailboxes (no shell access)
- Self-service: users manage their own password, vacation message, and quota
- Bulk import/export of users via CSV
- Argon2id password hashing for the panel; crypt for Dovecot

### Domain Management
- Host multiple email domains, owned per Organization
- Domain aliases and catch-all addresses
- Per-domain DKIM signing keys — generate, rotate, export DNS records
- DNS Helper — required MX, SPF, DKIM, DMARC records
- Custom domains — serve the panel itself under an Organization's own
  domain, with automatic Let's Encrypt issuance and TXT ownership
  verification

### Alias & Forwarding
- Aliases pointing to local or remote addresses
- Forwarding rules to multiple destinations
- Vacation (auto-reply) messages
- Sieve mail filtering (where Dovecot supports it)

### Spam & Virus Protection
- SpamAssassin and Rspamd spam detection
- Pyzor/Razor collaborative spam detection
- Amavisd-new content filtering
- ClamAV virus scanning for attachments
- Greylisting, RBL/DNSBL blocklists, and trusted-sender whitelists

### Security & Authentication
- TLS/SSL for SMTP, IMAP, POP3, with automatic Let's Encrypt renewal
- DKIM, SPF, and DMARC configuration
- SASL authentication configuration
- Fail2ban integration where installed

### Monitoring & Logs
- Mail queue: list, hold, delete, flush, retry
- Delivery and spam statistics
- Real-time service health status
- Searchable log viewer
- Email/webhook alerts for issues

## Production

### Docker (Recommended)

```bash
docker run -d \
  --name mailadmin \
  -p 172.17.0.1:64580:80 \
  -v ./volumes/config:/config:z \
  -v ./volumes/data:/data:z \
  ghcr.io/webappsgo/mail:latest
```

### Docker Compose

```bash
curl -q -LSsf -O https://raw.githubusercontent.com/webappsgo/mailadmin/main/docker/docker-compose.yml
docker compose up -d
```

### Binary

```bash
# Download latest release (linux/amd64, linux/arm64, freebsd/amd64, freebsd/arm64)
curl -q -LSsf -O https://github.com/webappsgo/mailadmin/releases/latest/download/mail-linux-amd64

# Make executable and run
chmod +x mail-linux-amd64
./mail-linux-amd64
```

## Configuration

Configuration is auto-generated on first run. Edit via the admin panel at
`http://{fqdn}:64580/admin` — the setup token is displayed on first run and
consumed on wizard completion.

Key settings (`server.yml`):

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

Use `driver: sqlite` for a single server, or `postgres`/`mysql` for a
clustered panel.

## API

All API endpoints are versioned under `/api/v1/`.

| Endpoint | Description |
|----------|-------------|
| `GET /healthz` | Health check (not versioned) |
| `GET /api/v1/mail/users` | List/manage virtual users |
| `GET /api/v1/mail/domains` | List/manage domains, DKIM, DNS |
| `GET /api/v1/mail/aliases` | List/manage aliases |
| `GET /api/v1/mail/forwards` | List/manage forwarding rules |
| `GET /api/v1/mail/queue` | Mail queue management |
| `GET /api/v1/mail/services` | Service status and control |
| `GET /api/v1/mail/stats` | Delivery and spam statistics |
| `GET /api/v1/mail/transports` | Transport rules |

### Examples

```bash
# Health check
curl -q -LSsf http://{fqdn}:64580/healthz

# API request (requires auth)
curl -q -LSsf -H "Authorization: Bearer TOKEN" http://{fqdn}:64580/api/v1/mail/domains
```

## Other

### Troubleshooting

- Check logs: `docker logs mailadmin`
- Health check: `curl -q -LSsf http://{fqdn}:64580/healthz`

## Development

**Development instructions are for contributors only.**

### Prerequisites

- Docker (REQUIRED — no local Go installation needed)
- GNU Make

### Build

```bash
# Clone
git clone https://github.com/webappsgo/mailadmin.git
cd mailadmin

# Quick dev build (outputs to OS temp dir)
make dev

# Full build (all platforms, outputs to binaries/)
make build

# Test
make test
```

Supported platforms: `linux/amd64`, `linux/arm64`, `freebsd/amd64`,
`freebsd/arm64` — mailadmin manages services (Postfix, Dovecot,
SpamAssassin, ClamAV) that only run on Linux and BSD; see SPEC.md.

### Project Structure

```
src/           # Source code
docker/        # Docker files
tests/         # Test scripts
```

## Disclaimer

This software is provided "as is" without warranty of any kind. Use at your own risk.

- **No Warranty**: The authors are not responsible for any damages, data loss, or issues arising from use of this software
- **Not Professional Advice**: This software does not constitute legal, financial, medical, or other professional advice
- **Third-Party Services**: If this software connects to external APIs or services, their terms of service apply separately
- **Security**: While we strive to follow security best practices, no software is guaranteed to be free of vulnerabilities
- **Production Use**: Evaluate thoroughly before deploying in production environments

By using this software, you acknowledge that you have read and understood this disclaimer.

## License

MIT - See [LICENSE.md](LICENSE.md)

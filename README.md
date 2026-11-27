# mailadmin

[![Build](https://img.shields.io/github/actions/workflow/status/webappsgo/mailadmin/ci.yml?branch=main)](https://github.com/webappsgo/mailadmin/actions)
[![Release](https://img.shields.io/github/v/release/webappsgo/mailadmin)](https://github.com/webappsgo/mailadmin/releases)
[![Documentation](https://readthedocs.org/projects/webappsgo-mailadmin/badge/?version=latest)](https://webappsgo-mailadmin.readthedocs.io/en/latest/)
[![License](https://img.shields.io/github/license/webappsgo/mailadmin)](LICENSE.md)

## About

mailadmin is a web administration panel for configuring and monitoring an
existing mail stack. Server Admins manage the host and service settings;
Organizations manage their own mail domains, mailboxes, aliases, and forwarding
rules.

## Features

- Manage mail domains, DKIM records, mailboxes, aliases, forwarding rules, and
  transport rules.
- Inspect the Postfix queue, delivery statistics, and mail logs.
- Manage Server Admin accounts, API tokens, MFA, preferences, and scheduled
  tasks.
- Configure site branding, SEO metadata, and outbound notification email.
- Verify custom panel domains and request TLS certificates.
- Use tenant-scoped Organization and mail-user pages.
- Operate through a server-rendered browser panel, versioned API, and companion
  CLI/TUI.

mailadmin configures mail services on its host. Message reading and composition
are outside the product's scope.

## Production

### Docker Compose

The compose file binds the panel to `172.17.0.1:64580` and persists
configuration and data under `docker/volumes/`.

```sh
cd docker
docker compose up -d
```

Review the port binding and mounted directories before deployment. Put a TLS
reverse proxy in front of the panel before exposing it outside a trusted
network. Complete first-run setup using the one-time token printed by the
server.

### Binary and system service

Release archives target Linux and FreeBSD on amd64 and arm64. Download the
archive for the target platform from [GitHub Releases](https://github.com/webappsgo/mailadmin/releases).
See the [installation guide](https://webappsgo-mailadmin.readthedocs.io/en/latest/installation/)
for first-run setup and service installation.

## Client

The `mailadmin-cli` companion provides API commands and an interactive TUI.

```sh
mailadmin-cli --server "$MAILADMIN_URL" status
mailadmin-cli login
```

See the [CLI and TUI guide](https://webappsgo-mailadmin.readthedocs.io/en/latest/cli/)
for authentication, commands, output formats, and cluster failover.

## Configuration

The server reads `server.yml`, ordered `yaml.d/*.yml` and `yaml.d/*.yaml`
drop-ins, and `secrets.yml`. The main server settings configure listeners,
database access, sessions, TLS, and request limits. Environment references in
YAML values use `${VARIABLE}` syntax. See the
[configuration reference](https://webappsgo-mailadmin.readthedocs.io/en/latest/configuration/)
for file permissions, merge order, setting defaults, and secret handling.

## API

The versioned API covers server status and information, Server Admin settings,
mail resources, Organizations, user self-service, and custom panel domains.
OpenAPI and GraphQL explorers are available at `/server/docs/swagger` and
`/server/docs/graphql`. See the [API reference](https://webappsgo-mailadmin.readthedocs.io/en/latest/api/)
for routes, authentication, and response formats.

## Other

mailadmin manages infrastructure and mail-account metadata; mail storage,
delivery, and mailbox protocols are provided by the installed mail stack.
Product requirements are in [IDEA.md](IDEA.md), implementation requirements
are in [AI.md](AI.md), and project-specific overrides are in [SPEC.md](SPEC.md).

## Development

Go builds and checks run through Docker with `CGO_ENABLED=0`. Use `make build`
to produce debug binaries for the supported Linux and FreeBSD amd64 and arm64
targets, and `make test` to run the lint, vulnerability, test, coverage, and
binary checks. See the
[development guide](https://webappsgo-mailadmin.readthedocs.io/en/latest/development/)
for test suites and contribution workflow.

## Disclaimer

This software is provided "as is" without warranty of any kind. Use it only
after evaluating its behavior and security for your environment.

- **No warranty:** The maintainers are not responsible for damages, data loss,
  or operational problems arising from use of the software.
- **Third-party services:** External mail, DNS, certificate, and hosting
  services have their own terms and availability.
- **Security:** The project follows security practices, but no software can be
  guaranteed free of vulnerabilities.
- **Production use:** Review configuration, access controls, backups, and
  recovery procedures before deployment.

By using this software, you acknowledge this disclaimer.

## License

MIT. See [LICENSE.md](LICENSE.md).

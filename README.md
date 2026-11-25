# mailadmin

mailadmin is a Go web panel for administering an existing mail stack, with separate Server Admin and tenant-scoped mailbox user accounts.

[![Build](https://img.shields.io/github/actions/workflow/status/webappsgo/mailadmin/ci.yml?branch=main)](https://github.com/webappsgo/mailadmin/actions)
[![Release](https://img.shields.io/github/v/release/webappsgo/mailadmin)](https://github.com/webappsgo/mailadmin/releases)
[![Documentation](https://readthedocs.org/projects/webappsgo-mailadmin/badge/?version=latest)](https://webappsgo-mailadmin.readthedocs.io/en/latest/)
[![License](https://img.shields.io/github/license/webappsgo/mailadmin)](LICENSE.md)

## Install

### Docker Compose

The provided compose file binds the web panel to `172.17.0.1:64580` and persists configuration and data under `docker/volumes/`.

```sh
cd docker
docker compose up -d
```

Review and edit the port binding and mounted directories in `docker/docker-compose.yml` before deployment. Put a TLS reverse proxy in front of the panel for remote access. First-run setup instructions and the one-time setup token are printed by the binary; complete setup at `/server/{admin_path}/setup`.

### Binary

Release builds target `linux/amd64`, `linux/arm64`, `freebsd/amd64`, and `freebsd/arm64`.
Download the archive for the target platform from the [GitHub releases](https://github.com/webappsgo/mailadmin/releases) and follow [installation documentation](https://webappsgo-mailadmin.readthedocs.io/en/latest/installation/).

## Product scope

mailadmin configures and monitors mail services installed on its host, including Postfix, Dovecot, spam filtering, and virus scanning. It is not a mail transfer agent, IMAP server, or webmail client. Server Admins manage the panel and server; mail users access only their own account and organization resources. Organizations own their mail domains and isolate member access.

The implementation requirements are maintained in [AI.md](AI.md), product behavior in [IDEA.md](IDEA.md), and project overrides in [SPEC.md](SPEC.md). Operator documentation describes only the behavior currently shipped.

## Web and API

Public pages include `/`, `/server/about`, `/server/help`, `/server/privacy`, `/server/terms`, and `/server/healthz`. The `/healthz` compatibility alias is opt-in through `server.healthz.root.enabled`. First-run setup is at `/server/{admin_path}/setup`; Server Admin pages use the configured `/server/{admin_path}` prefix. User self-service starts at `/users`, with organization pages under `/orgs`.

The implemented API includes public health and metadata routes, setup-token verification, Server Admin settings and profile routes, mailbox management, and organization-scoped email alias CRUD. Protected endpoints require a Bearer token and enforce organization membership and token scope. The [API reference](https://webappsgo-mailadmin.readthedocs.io/en/latest/api/) lists the available endpoints; the full AI.md API inventory is not yet implemented.

## Development

Builds and Go commands must run inside the project Docker workflow with `CGO_ENABLED=0`; do not build Go on the host. See [development documentation](https://webappsgo-mailadmin.readthedocs.io/en/latest/development/) and the project Makefile for supported commands. The platform matrix is Linux and FreeBSD on amd64 and arm64.

## License

MIT. See [LICENSE.md](LICENSE.md).

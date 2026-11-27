# Installation

## Docker Compose

The repository's `docker/docker-compose.yml` publishes the panel at `172.17.0.1:64580` and mounts `docker/volumes/config` and `docker/volumes/data` into the container. Review those host paths and the bind address before starting it.

```sh
cd docker
docker compose up -d
```

The container health check invokes `mailadmin --status`. Place a TLS reverse proxy in front of the panel before exposing it outside a trusted network.

## Binary releases

Download the archive matching one of these supported targets from [GitHub Releases](https://github.com/webappsgo/mailadmin/releases):

| Operating system | Architectures |
|---|---|
| Linux | amd64, arm64 |
| FreeBSD | amd64, arm64 |

Keep the configuration and database directories persistent and accessible only to the service account. A root-started system service creates a dedicated account in the safe UID/GID range, transfers the application directories to it, binds configured listeners, and then drops privileges before serving requests. A user-mode start stays under the invoking account. If the stable internal account name already belongs to an incompatible system identity, the installer uses the project-name account when that name is available and leaves the existing identity untouched.

## First start

The application starts with generated defaults and prints a one-time setup token. Use the URL shown by the process to complete the initial Server Admin setup. Treat the token as a credential; it is single-use and is not recoverable after setup completes.

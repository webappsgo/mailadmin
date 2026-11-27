# Development

Read `AI.md`, `IDEA.md`, and `SPEC.md` before changing behavior. `SPEC.md` narrows the supported build matrix to Linux and FreeBSD on amd64 and arm64; activated multi-user, organization, and custom-domain requirements remain in scope.

## Build environment

Go build and test commands run inside the project Docker workflows, with `CGO_ENABLED=0`. The Makefile has six targets: `build`, `test`, `release`, `docker`, `clean`, and `help`. Use `make build` for debug binaries across the Linux and FreeBSD amd64 and arm64 matrix, and `make test` for vet, static analysis, vulnerability scanning, unit coverage, and binary checks. Keep generated binaries and test artifacts out of source directories.

`make release` produces stripped production binaries, a source archive, and checksums without tagging or publishing. `make docker` builds the host-platform image locally without pushing it. CI owns release tagging and publishing.

## Integration suites

`./tests/run_tests.sh` runs the Docker binary-validation and first-run setup checks. `./tests/e2e.sh` runs the browser suite in Docker: raw SSR checks, setup with JavaScript disabled, and JavaScript-enabled login, theme, and mobile viewport checks. It requires Docker image access for the Go build image and the `chromedp/headless-shell:stable` browser image. `--keep-artifacts` preserves its isolated logs and runtime files under the project-scoped temporary directory after a failure.

`./tests/incus.sh` exercises service installation, restart, stop/start, reboot persistence, disable, reinstall, and uninstall under Debian/systemd and Alpine/OpenRC. Pass `systemd` or `openrc` to select one mode. It runs the Docker binary suite first, and requires Docker, Incus, and access to the configured Incus images. The script deletes only the uniquely named Incus containers it creates.

## Contribution expectations

Keep tenant authorization in persistence operations, escape all server-rendered user data through `html/template`, protect every state-changing form with CSRF, and avoid logging credentials. Update documentation whenever routes, config, or operator-visible behavior changes. CI definitions live alongside the supported provider workflows.

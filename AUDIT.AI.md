# Project Audit

Started: 2026-08-22

Scope: compliance audit of mailadmin against AI.md (read-only spec) and
SPEC.md (4-platform Linux/BSD override; PART 34/35/36 activated).
Delete this file once every item below is resolved — do not empty it.

## Pass 1: Security

- [x] `docker/server.yml.template`: config templates are never committed
      (AI.md "NEVER Create These Files": `server.yml`, `*.example.*`,
      `*.sample.*` — config is runtime-generated). Deleted; defaults still
      need to be moved into the binary if not already there.

## Pass 2: Code Quality

- [ ] `src/main.go`, `src/server/server.go`, `src/server/handlers.go`,
      `src/config/config.go`: 20 `TODO:` comments in production code.
      AI.md PART 0 forbids TODO comments and stub functions outright.
- [ ] `src/server/handlers.go`: plural filename and flat layout. AI.md
      requires `src/server/handler/`, `service/`, `model/`, `store/`
      subpackages with singular names.
- [ ] Empty directories committed with no content: `src/mode/`,
      `src/signal/`, `src/service/`, `scripts/`, `tests/`.

## Pass 3: Logic

- [ ] `src/main.go`: `--mode`, `--config`, `--data`, `--log`, `--pid`,
      `--address`, `--port`, `--baseurl` are printed by `printHelp()` but
      never parsed by the flag switch — documented flags are inert.
      AI.md "CLI Rules" marks these NON-NEGOTIABLE.
- [ ] `Makefile` LDFLAGS + `src/main.go:23`: uses `main.BuildDate` as an
      ldflag. AI.md "CI/CD Rules" requires `-X main.BuildEpoch=...` with
      BuildDate derived from it inside `init()`.

## Pass 4: Documentation

- [ ] Missing required root files per AI.md "Allowed Root Files":
      `mkdocs.yml`, `.readthedocs.yaml`, `Jenkinsfile`.
- [x] `go.sum` exists on disk but is untracked (was gitignored — ignore
      rule removed in this audit). Staged with `git add`.
- [x] `IDEA.md.preMigration.bak`: untracked, gitignored, never committed.
      `*.bak` is an AI.md naming anti-pattern ("use git") and the file is
      not on the root allowlist. Deleted; `.gitignore` rule removed too.

## Pass 5: Spec Compliance

- [ ] `docker/Dockerfile` does not exist. AI.md "Docker Rules" requires a
      multi-stage `docker/Dockerfile` (casjaysdev/go:latest builder +
      alpine runtime, tini entrypoint, STOPSIGNAL SIGRTMIN+3, port 80).
- [ ] `docker/docker-compose.yml` and `docker/rootfs/` do not exist.
- [ ] No CI/CD at all: no `.github/workflows/` (ci.yml, docker.yml,
      release.yml, secret scanning) and no `Jenkinsfile`.
- [ ] `.claude/rules/` rule-file set (14 files per AI.md) not created.
- [ ] Missing required packages: `src/client/` (REQUIRED for all
      projects, PART 33), `src/swagger/`, `src/graphql/`,
      `src/common/i18n/`, `src/ssl/`.

## Pass 7: README/Makefile Alignment

- [x] README.md was severely out of sync with IDEA.md — wrong org
      (`apimgr` instead of `webappsgo`), wrong port (64500 instead of
      64580), missing REQUIRED Disclaimer section and badges, wrong
      template structure, and it described features never defined in
      IDEA.md (Webmail, CalDAV/CardDAV, Mailing Lists, multi-server MX
      roles, PWA, push notifications, LDAP/OIDC). Rewritten to match
      IDEA.md's actual current business logic and AI.md's README
      Template/Update Rules.
- [ ] `Makefile`: `REGISTRY ?= ghcr.io/$(PROJECTORG)/$(PROJECTNAME)` uses
      `PROJECTNAME` (= "mailadmin", the repo name) instead of AI.md's
      canonical `{project_org}/{internal_name}` registry pattern, which
      per IDEA.md's Project variables should resolve to
      `ghcr.io/webappsgo/mail`. The new README.md already documents the
      corrected `ghcr.io/webappsgo/mail` name; the Makefile itself still
      needs `REGISTRY` changed to use `internal_name` ("mail"), not
      `PROJECTNAME`.
- [ ] `Makefile`: release-binary output naming
      (`$(BINDIR)/$(PROJECTNAME)-$$OS-$$ARCH`, and the `-cli-`/`-agent-`
      variants) produces `mailadmin-linux-amd64` etc., but IDEA.md's
      Project variables declare `binary: mail` and its Project
      description states "Binary name: `mail`" — confirmed against the
      actual prebuilt `binaries/mail` ELF binary, which is already named
      `mail`. `OUTPUT=` lines need `$(PROJECTNAME)` replaced with the
      IDEA.md `binary` value ("mail") to stop producing mismatched
      release artifact names.

## Pass 6: Code Flow Trace

- [ ] Tenancy divergence (PART 34/35/36, activated by SPEC.md):
      `src/database/schema_users.go` has `orgs`, `org_members`,
      `org_preferences` (compatible), and `mail_domains` carries
      `owner_type`/`owner_id`. But `mail_mailboxes`, `mail_aliases`, and
      `mail_forwards` carry no `org_id` and inherit tenancy only via
      `domain_id`. IDEA.md → Data model & sensitivity requires an
      `OrgID` foreign key on every `VirtualUser`/`VirtualDomain`/
      `Alias`/`Forward`/`Transport` row as the enforced isolation
      boundary. Reconcile schema with IDEA.md before handlers are built.
- [ ] No handlers, services, or middleware enforce Org scoping yet —
      tenancy isolation is unimplemented, not merely undocumented.

## Completed

- Makefile: PLATFORMS reduced from 8 to the 4 SPEC.md targets
  (linux/freebsd × amd64/arm64); windows `.exe` suffix lines removed.
- Makefile: toolchain image `golang:alpine` → `casjaysdev/go:latest`,
  added `-e GOFLAGS=-buildvcs=false` for mounted-.git builds.
- .gitignore: stopped ignoring `go.sum` (required committed root file).
- .gitignore: narrowed blanket `.claude/` `.cursor/` `.windsurf/` ignores
  to personal settings, cache, and lock files only — team config is
  committed per AI.md "AI-Specific Files and Directories".
- CLAUDE.md: created the required root loader (rule hierarchy +
  non-negotiables + pointers into AI.md).
- release.txt: added missing trailing newline.
- README.md: removed the macOS/Windows support claim that contradicted
  SPEC.md's Linux/BSD-only platform target.
- src/database/schema_mail.go: replaced a stale `per PLAN.AI.md line
  1121` comment that pointed at a file that no longer has that content.
- PLAN.AI.md: replaced the false "Fully Implemented" claim with the
  actual scaffold status.

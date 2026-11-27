# AI.md Alignment Tasks

Implementation order follows the AI.md dependency order: finish the foundation and identity/data boundaries first, then management features and server pages, then deployment and documentation. SPEC.md narrows the release matrix to Linux/FreeBSD on amd64/arm64 and activates PARTS 34-36.

## Critical (P0)

- [ ] Complete required server setup, identity, session, authorization, MFA, CSRF, rate limiting, audit/security logging, and tenant isolation across all routes and persistence methods.
- [ ] Complete the Org-owned mail domain, virtual mailbox, alias, forward, and transport lifecycle for Server Admins and Org Owner/Member roles.
- [ ] Implement custom panel domain ownership verification, host routing, certificate issuance/renewal, and failure handling from PART 36.
- [ ] Remove any remaining explicit unimplemented behavior discovered during implementation and complete its required feature rather than masking it.

## High (P1) — Required AI.md baseline and activated features

- [ ] Complete the API inventory, authentication/scope rules, request validation, pagination, OpenAPI, GraphQL, and route integration coverage.
- [ ] Complete all server-rendered setup, Server Admin, mail-user, Organization, security, service, domain, queue, statistics, log, and custom-domain pages and forms.
- [ ] Finish PART 18 email delivery: wire remaining event dispatchers and account flows (including email-change verification and invite/activation delivery), add regular-user email preference overrides and translations, add remaining runtime variables, and verify SMTP UI behavior when delivery is unavailable. Open registration, one-time email verification, post-verification welcome delivery, and username reminders are implemented for enabled multi-user mode. Admin event-specific WebUI/email preference storage, admin password reset, and account alerts for login, password, token, and TOTP changes are implemented.
- [ ] Complete remaining TLS/ACME, metrics, updates, privilege boundaries, backup/restore, Tor/I2P behavior, and all required scheduler tasks.
- [ ] Finish PART 5 maintenance recovery: reconnect and initialize unavailable databases, detect database writes and file read/write failures, remediate disk pressure under configured retention, expose full diagnostics and error-specific guidance, honor retry/max-attempt/cleanup/notification settings, and notify on recovery.
- [ ] Complete the required client CLI/TUI, shell completions, configuration, server API operations, and cluster discovery/failover.
- [ ] Complete the required project layout and governance files, including test scripts, rule files, project memory, service assets, source archives, and license attribution.
- [ ] Complete production/dev/test/AIO Docker behavior and align all GitHub, GitLab, Gitea, Forgejo, and Jenkins pipelines with identical gates and the SPEC platform matrix.
- [ ] Raise unit coverage to at least 60%; add the required Docker/Incus/browser suites and achieve full route coverage; fix all build, lint, vulnerability, workflow, and documentation build failures.

## Medium (P2) — Structure and generated project configuration

- [ ] Remove prohibited scaffolding, stale files, inaccurate claims, and unsupported behavior after each feature is fully implemented.
- [ ] Bring IDEA.md, README.md, docs, API schemas, configuration reference, and operations guidance into sync with shipped behavior.
- [ ] Keep this checklist limited to open project implementation work; record only verified AI.md-internal contradictions in TODO.SERVER.md.

## Low (P3) — Completion

- [ ] Re-read the relevant AI.md sections and SPEC.md overrides, complete the final compliance checklist, pass the required test/lint/build/workflow/documentation gates, and remove PLAN.AI.md when the work it describes is complete.

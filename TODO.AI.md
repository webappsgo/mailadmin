# AI.md Alignment Tasks

Implementation order follows dependencies: shared configuration, persistence,
identity and authorization come before tenant-scoped APIs and pages; infrastructure
and CI run after the source layout is stable; documentation is synchronized last.

## Critical (P0)

- [ ] Replace placeholder startup, configuration, path, database, logging, session, and security behavior with complete AI.md implementations.
- [ ] Implement Server Admin setup, authentication, sessions, MFA, authorization, CSRF protection, audit logging, and safe privilege boundaries.
- [ ] Implement Regular User identities and Organization membership with enforced `OrgID` scoping on every tenant-owned operation.
- [ ] Implement custom panel domain ownership verification and certificate lifecycle without routing unverified domains.
- [ ] Implement mail-management services with validated configuration writes, safe service control, and secret redaction.

## High (P1) — Required AI.md baseline and activated features

- [ ] Complete the required CLI and client binary behavior, shell completion/init, status, service, maintenance, update, and color/language handling.
- [ ] Complete database schema, migrations, pooling, backups, restore, and cluster behavior.
- [ ] Implement API route inventory and behavior, content negotiation, OpenAPI, GraphQL, rate limits, and error responses.
- [ ] Implement every server-rendered page and form for setup, administration, self-service, organization tenancy, and custom domains.
- [ ] Implement TLS/ACME, built-in scheduler, GeoIP, metrics, notifications, backup/restore, update, service integration, Tor, and I2P requirements.
- [ ] Implement i18n, accessibility, theme, SEO, static files, health, monitoring, and cache behavior.
- [ ] Add complete unit, integration, and browser E2E coverage and resolve all test, build, lint, vulnerability, and CI errors.
- [ ] Align Docker images, compose files, all required CI providers, security workflows, release workflows, and service packaging with AI.md and SPEC.md.

## Medium (P2) — Structure and generated project configuration

- [ ] Complete required source, docs, scripts, tests, service, locale, and AI rule-file structures.
- [ ] Remove scaffolding artifacts, placeholder code, stale files, and any behavior prohibited by AI.md.
- [ ] Update IDEA.md only where current product behavior changes; keep AI.md read-only and preserve SPEC.md overrides.

## Low (P3) — User documentation

- [ ] Bring README.md and ReadTheDocs content, API schemas, and operational examples into sync with the complete application.
- [ ] Remove or update PLAN.AI.md after all implementation and verification work is complete.

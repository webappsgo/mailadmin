# Project Audit

Started: 2026-11-26

Scope: PART 17 admin/public page compliance and IDEA.md business-scope consistency only. This is not a completed six-pass or whole-AI.md audit.

## Pass 5: Spec Compliance

- [ ] Admin routes: required management pages are missing. AI.md PART 17 lines 33976–34002 requires `/email`, `/security/ratelimit`, `/security/firewall`, `/security/allowlist`, `/network/tor`, `/network/i2p`, `/network/geoip`, `/network/blocklists`, `/moderation/users`, `/users/invites`, `/backup`, `/updates`, `/nodes`, and `/nodes/add`. The admin router in `src/server/server.go:273-334` registers none of these routes. IDEA.md Business logic also requires mail service configuration/control, spam/virus configuration, and backup/restore, but the web management pages are absent.
- [ ] Admin API mirror: AI.md PART 17 lines 34815–34963 requires setup status/config/account endpoints, admin account count/invite/enable/disable/delete, passkey setup/register, Tor, web/robots/security/well-known, About/Privacy/Contact page editors, email settings/templates, and backup/restore APIs. `src/server/server.go:516-546` registers only settings, branding, stats, help, logs, SSL, scheduler, profile/preferences, TOTP, and custom-domain APIs; the required endpoint groups are absent.
- [ ] Dashboard: AI.md PART 17 lines 33918–33974 requires uptime, request and error counts, CPU/memory/disk resources, quick actions, recent activity, scheduled tasks, and alerts. `src/server/pages.go:713-738` renders only database availability, application/version fields, and navigation cards.
- [ ] Admin layout/navigation: AI.md PART 17 lines 33231–33333 specifies grouped collapsible navigation, section icons, active-page indicator, persisted expanded state, collapse-all behavior, and mobile slide-out navigation. `src/server/web/templates/layout.tmpl:84-88` has a flat link list and checkbox menu only; it has no active state, grouping, icons, remembered state, or collapse-all behavior. The header also lacks the specified global search and status indicator.
- [ ] Login form: AI.md PART 17 lines 33335–33347 requires a Remember Me option extending sessions to 90 days. `src/server/web/templates/layout.tmpl:86` login form contains username, password, optional TOTP/recovery inputs, and submit only; no remember option or corresponding behavior was found.
- [ ] Public/admin separation: AI.md PART 17 lines 32821–32832 forbids public admin-panel mentions. `src/server/pages.go:448-452` puts “Server Admin setup and sign-in” on the public About page, and `src/server/pages.go:462-464` says “Server Admins” on the public Help page. These are reachable public routes registered at `src/server/server.go:187-190`.
- [ ] Admin route hierarchy: AI.md PART 17 lines 32882–33045 requires fixed root sections to be validated against `validAdminRootPaths` and undeclared project-specific sections to be declared in IDEA.md. The router uses direct literals in `src/server/server.go:273-334`, with no `validateAdminRoute` or equivalent guard found. It also mounts project mail-management pages under `/server/{admin_path}/mail/...`, while the allowed-root list at AI.md lines 32971–33015 has no `mail` entry and IDEA.md does not declare `mail` as a project-specific admin route section.
- [ ] Page documentation: `docs/admin.md:5-27` documents the current subset and omits the required management pages above. The API guide at `docs/api.md:20-49` similarly documents only a subset of the Admin API inventory in AI.md PART 17 lines 34815–34963.

## Passes 1–4 and 6

Not audited in this scoped page review.

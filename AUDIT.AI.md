# Project Audit

Started: 2026-11-26

Scope: PART 17 admin/public page compliance and IDEA.md business-scope consistency only. This is not a completed six-pass or whole-AI.md audit.

## Pass 5: Spec Compliance

- [ ] Admin routes: major required management pages remain missing. AI.md PART 17 lines 33976–34002 requires `/security/ratelimit`, `/security/firewall`, `/security/allowlist`, `/network/tor`, `/network/i2p`, `/network/geoip`, `/network/blocklists`, `/moderation/users`, `/users/invites`, `/backup`, `/updates`, `/nodes`, and `/nodes/add`. `/email` was added, but other listed route groups and IDEA.md mail-service/spam-virus/backup management pages remain absent.
- [ ] Admin API mirror: AI.md PART 17 lines 34815–34963 requires setup status/config/account endpoints, admin account count/invite/enable/disable/delete, passkey setup/register, Tor, web/robots/security/well-known, About/Privacy/Contact page editors, email templates, and backup/restore APIs. Email settings GET/PATCH and test-email POST were added; the other required endpoint groups remain absent. Current route registration is in `src/server/server.go`.
- [ ] Dashboard: AI.md PART 17 lines 33918–33974 requires uptime, request and error counts, CPU/memory/disk resources, quick actions, recent activity, scheduled tasks, and alerts. `src/server/pages.go` has database/build information and navigation cards, including branding and email; the required resource and activity widgets remain missing.
- [ ] Admin layout/navigation: AI.md PART 17 lines 33231–33333 specifies grouped collapsible navigation, section icons, active-page indicator, persisted expanded state, collapse-all behavior, and mobile slide-out navigation. `src/server/web/templates/layout.tmpl` still has a flat link list and checkbox menu only; it has no active state, grouping, icons, remembered state, or collapse-all behavior. The header also lacks the specified global search and status indicator.
- [ ] Login form: AI.md PART 17 lines 33335–33347 requires a Remember Me option extending sessions to 90 days. The login form in `src/server/web/templates/layout.tmpl` still has no remember option or corresponding behavior.
- [x] Public/admin separation copy: public About and Help copy no longer mentions Server Admin setup/sign-in or Server Admins; those public mentions were replaced with audience-neutral wording.
- [ ] Admin route hierarchy: AI.md PART 17 lines 32882–33045 requires fixed root sections to be validated against `validAdminRootPaths` and undeclared project-specific sections to be declared in IDEA.md. The mail group is now declared in IDEA.md. The router still uses direct literals in `src/server/server.go` without `validateAdminRoute` or equivalent route-set validation.
- [ ] Page documentation: `docs/admin.md` and `docs/api.md` now cover branding and SMTP email settings, but still omit the management pages and API groups listed above.

## Passes 1–4 and 6

Not audited in this scoped page review.

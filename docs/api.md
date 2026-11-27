# API reference

Canonical API paths use `/api/{api_version}`. The default `server.api_version` is `v1`; the configured value is mounted as the sole canonical version. API responses default to JSON for browsers and the `mailadmin-cli` client. Use `Accept: text/plain`, an HTTP-tool user agent, or a `.txt` suffix on a text-serializable API path for plain text. `Accept: application/json` explicitly requests JSON; a `.txt` suffix takes precedence. Binary endpoints do not have `.txt` representations.

## API explorers

- OpenAPI explorer: `/server/docs/swagger`; generated JSON: `/api/{api_version}/server/swagger` and the direct alias `/api/swagger`.
- GraphQL workbench: `/server/docs/graphql`; query-only endpoint: `/api/{api_version}/server/graphql` and the direct alias `/api/graphql`.
- GraphQL exposes `serverAbout`, `serverVersion`, and generated fields for each listed versioned REST `GET` operation. Route fields return JSON scalar data and use the REST handlers' authorization checks. It accepts queries over GET and POST and does not expose mutations.

## Implemented endpoints

| Method | Path | Authentication | Response |
|---|---|---|---|
| GET | `/api/v1/server/healthz` | Public | Public-safe status, build/runtime details, enabled features, component checks, aggregate request statistics, and the last discovered public IPv4/IPv6 addresses |
| GET | `/api/v1/server/help` | Public | Help Markdown source, format, and public help page index |
| GET | `/api/v1/server/about` | Public | Application description and features |
| GET | `/api/v1/server/auth/check` | Bearer token | Token kind (`operator`, `owner`, or `user`); does not return or persist the token |
| POST | `/api/v1/server/auth/register` | Public, when multi-user open registration is enabled | Create a regular user account and send email verification when SMTP delivery is available |
| POST | `/api/v1/server/auth/verify` | Public verification token | Verify and activate a regular user account; token is single-use |
| GET | `/api/v1/server/version` | Public | Version, commit, build date, Go version, platform |
| POST | `/api/v1/server/auth/password/forgot` | Public | Generic response; sends a one-time reset link when an account matches |
| POST | `/api/v1/server/auth/username/forgot` | Public | Generic response; sends a username reminder to a verified account email when it matches |
| GET | `/api/v1/server/security` | Public | Security reporting channels, current expiration, policy and acknowledgments URLs |
| GET | `/api/v1/server/security/policy` | Public | Coordinated disclosure policy Markdown |
| GET | `/api/v1/server/security/thanks` | Public | Researcher acknowledgment entries |
| GET | `/api/v1/server/contact` | Public | Return the selected report type's field descriptors, simple challenge, signed form token, available types, security URL, and public abuse email |
| POST | `/api/v1/server/contact` | Public | Submit an abuse, DMCA, or privacy report using the two-step review and confirmation flow |
| POST | `/api/v1/server/auth/password/reset` | Public reset token | Reset password and invalidate existing sessions |
| POST | `/api/v1/server/{admin_path}/setup/verify` | Setup token | Token validity |
| GET | `/api/v1/server/{admin_path}/{admin_username}/profile` | Admin Bearer token, matching username | Admin username and email |
| PATCH | `/api/v1/server/{admin_path}/{admin_username}/profile` | Admin Bearer token, matching username | Update the admin username |
| POST | `/api/v1/server/{admin_path}/{admin_username}/profile/password` | Admin Bearer token, matching username | Change password after current-password verification |
| GET | `/api/v1/server/{admin_path}/{admin_username}/profile/tokens` | Admin Bearer token, matching username | Masked current API token |
| POST | `/api/v1/server/{admin_path}/{admin_username}/profile/tokens` | Admin Bearer token, matching username | Regenerate token; the new value is returned once |
| GET | `/api/v1/server/{admin_path}/{admin_username}/preferences` | Admin Bearer token, matching username | Admin preferences |
| PATCH | `/api/v1/server/{admin_path}/{admin_username}/preferences` | Admin Bearer token, matching username | Update admin preferences |
| GET | `/api/v1/server/{admin_path}/{admin_username}/notifications` | Admin Bearer token, matching username | Per-event WebUI and email preferences (`data.webui` and `data.email` maps) |
| PATCH | `/api/v1/server/{admin_path}/{admin_username}/notifications` | Admin Bearer token, matching username | Replace event preference maps; required security events remain enabled |
| GET | `/api/v1/server/{admin_path}/{admin_username}/profile/security` | Admin Bearer token, matching username | TOTP and passkey status |
| POST | `/api/v1/server/{admin_path}/{admin_username}/profile/security/totp/setup` | Admin Bearer token, matching username | Begin TOTP enrollment after password verification |
| POST | `/api/v1/server/{admin_path}/{admin_username}/profile/security/totp/enable` | Admin Bearer token, matching username | Verify and enable pending TOTP enrollment |
| GET | `/api/v1/server/{admin_path}/settings` | Admin Bearer token | Read general server settings and pending-restart status |
| GET | `/api/v1/server/{admin_path}/branding` | Admin Bearer token | Read branding, SEO, and UI theme settings |
| PATCH | `/api/v1/server/{admin_path}/branding` | Admin Bearer token | Update supported branding, SEO, and UI theme settings in `yaml.d/branding.yml` |
| GET | `/api/v1/server/{admin_path}/email` | Admin Bearer token | Read SMTP readiness, host, port, username, TLS mode, sender, and whether a password is configured |
| PATCH | `/api/v1/server/{admin_path}/email` | Admin Bearer token | Update SMTP settings and any of the 16 AI.md email event preferences; an omitted password preserves the saved secret and an empty password clears it |
| POST | `/api/v1/server/{admin_path}/email/test` | Admin Bearer token | Send a test message to `{"to":"recipient@example.net"}` using the saved SMTP settings |
| GET | `/api/v1/server/{admin_path}/email/templates` | Admin Bearer token | List embedded defaults and active customized email templates |
| GET | `/api/v1/server/{admin_path}/email/templates/{name}` | Admin Bearer token | Read one email template |
| PUT | `/api/v1/server/{admin_path}/email/templates/{name}` | Admin Bearer token | Save a template using `{"subject":"...","body":"..."}`; variables use `{variable}` syntax |
| POST | `/api/v1/server/{admin_path}/email/templates/{name}/preview` | Admin Bearer token | Render a candidate template using optional `values` overrides without saving |
| POST | `/api/v1/server/{admin_path}/email/templates/{name}/reset` | Admin Bearer token | Delete the custom template and restore the embedded default |
| GET | `/api/v1/server/{admin_path}/stats` | Admin Bearer token | Read aggregate request totals, rolling 24-hour requests, and active connections |
| GET | `/api/v1/server/{admin_path}/pages/help` | Admin Bearer token | Read the editable public help Markdown |
| PATCH | `/api/v1/server/{admin_path}/pages/help` | Admin Bearer token | Update the public help Markdown |
| GET | `/api/v1/server/{admin_path}/logs?log=access.log&limit=100&offset=0&search=` | Admin Bearer token | Read a bounded, newest-first page from one of the six supported operator logs; sensitive assignments are masked |
| PATCH | `/api/v1/server/{admin_path}/settings` | Admin Bearer token | Update validated general settings, including single or dual listener ports; reports when a restart is required |
| GET | `/api/v1/server/{admin_path}/ssl` | Admin Bearer token | Read panel TLS and ACME settings with DNS credentials masked |
| PATCH | `/api/v1/server/{admin_path}/ssl` | Admin Bearer token | Update panel ACME challenge/provider settings; credentials are encrypted at rest |
| POST | `/api/v1/server/{admin_path}/ssl/renew` | Admin Bearer token | Request certificate renewal for the configured server FQDN |
| GET | `/api/v1/server/{admin_path}/scheduler` | Admin Bearer token | List scheduled task state without exposing task handlers or internal errors |
| GET | `/api/v1/server/{admin_path}/scheduler/{id}` | Admin Bearer token | Read one scheduled task |
| PATCH | `/api/v1/server/{admin_path}/scheduler/{id}` | Admin Bearer token | Update a task's cron/interval schedule and enabled state |
| POST | `/api/v1/server/{admin_path}/scheduler/{id}/run` | Admin Bearer token | Queue an immediate task run |
| POST | `/api/v1/server/{admin_path}/scheduler/{id}/enable` | Admin Bearer token | Enable a scheduled task |
| POST | `/api/v1/server/{admin_path}/scheduler/{id}/disable` | Admin Bearer token | Disable a scheduled task |
| GET | `/api/v1/mail/users?org_id={id}` | Admin, user, or organization Bearer token | List safe mailbox metadata in the authorized organization |
| POST | `/api/v1/mail/users/import` | Operator Bearer token | Import `text/csv` with `org_id,domain_id,email,password,name,quota_mb,enabled` columns; returns per-row success or error codes and aggregate counts. The upload is limited to 1 MiB and 5,000 rows; passwords are never echoed. |
| GET | `/api/v1/mail/users/export` | Operator Bearer token | Download credential-free mailbox metadata as CSV; spreadsheet formula-leading values are neutralized. |
| GET, POST | `/api/v1/mail/domains` | Admin, user, or organization Bearer token | List mail domains in scope with `limit`, `offset`, `sort` (`id` or `domain`), and `order` (`asc` or `desc`) pagination, or create one with `org_id`, `domain`, and optional `description`; list responses include pagination metadata and user tokens must choose an organization |
| GET, PUT, DELETE | `/api/v1/mail/domains/{id}` | Admin, user, or organization Bearer token | Read, update, or delete a mail domain in the selected organization; GET and DELETE use the `org_id` query parameter, while PUT accepts `org_id`, `domain`, and `description` in JSON. DELETE returns `409 CONFLICT` while mailboxes, aliases, DKIM, SPF, DMARC, or transport rows depend on the domain |
| POST | `/api/v1/mail/domains/{id}/dkim/keys` | Operator Bearer token | Generate a 2048-bit DKIM signing key for the domain selected by `org_id`; the private key is encrypted at rest and never returned. Accepts a selector and returns the public key and DNS name |
| GET | `/api/v1/mail/domains/{id}/dkim/dns` | Admin, user, or organization Bearer token | Return the active DKIM TXT record for the domain selected by `org_id`; an unconfigured key is reported as `null` |
| GET | `/api/v1/mail/domains/{id}/dns` | Admin, user, or organization Bearer token | Return MX (when the server FQDN is configured), SPF and DMARC guidance including stored policies, and the active DKIM TXT record when configured |
| PATCH | `/api/v1/mail/domains/{id}/dns` | Admin, user, or organization Bearer token | Update `spf`, `dmarc_policy`, `subdomain_policy`, `percentage`, and `aggregate_address` in JSON; returns the refreshed records. Policy values are `none`, `quarantine`, or `reject`, with an empty subdomain policy inheriting the domain policy |
| GET, POST | `/api/v1/orgs/{slug}/mail/domains` | Organization member or operator Bearer token | List mail domains with pagination or create one for the organization in the path; the path supplies the tenant scope |
| GET, PUT, DELETE | `/api/v1/orgs/{slug}/mail/domains/{id}` | Organization member or operator Bearer token | Read, update, or delete one mail domain in the organization; conflicting dependencies prevent deletion |
| GET, POST | `/api/v1/orgs/{slug}/aliases` | Organization member or operator Bearer token | List or create aliases for the organization in the path |
| GET, PUT, DELETE | `/api/v1/orgs/{slug}/aliases/{id}` | Organization member or operator Bearer token | Read, update, or delete one alias; writes require an appropriate token scope |
| GET, POST | `/api/v1/orgs/{slug}/mailboxes` | Organization member or operator Bearer token | List or create mailboxes for the organization in the path |
| GET, PUT, DELETE | `/api/v1/orgs/{slug}/mailboxes/{id}` | Organization member or operator Bearer token | Read, update, or delete one mailbox; password hashes are never returned |
| POST | `/api/v1/mail/users` | Admin, user, or organization Bearer token with write scope | Create a mailbox account |
| GET | `/api/v1/mail/users/{id}?org_id={id}` | Admin, user, or organization Bearer token | Read safe metadata for one mailbox |
| PUT | `/api/v1/mail/users/{id}` | Admin, user, or organization Bearer token with write scope | Update mailbox metadata and optionally its password |
| DELETE | `/api/v1/mail/users/{id}?org_id={id}` | Admin, user, or organization Bearer token with global scope | Delete a mailbox |
| GET | `/api/v1/mail/aliases?org_id={id}&domain_id={id}` | Admin, user, or organization Bearer token | List aliases in the authorized organization and optional domain |
| POST | `/api/v1/mail/aliases` | Admin, user, or organization Bearer token with `read-write` or `global` scope | Create an alias |
| GET, PUT | `/api/v1/mail/aliases/{id}` | Admin, user, or organization Bearer token with `read` or `read-write` scope | Read or update one alias |
| DELETE | `/api/v1/mail/aliases/{id}?org_id={id}` | Admin, user, or organization Bearer token with `global` scope | Delete an alias |
| GET, POST | `/api/v1/mail/forwards` | Admin, user, or organization Bearer token | List forwarding rules with `limit` (1–100, default 50), `offset`, `sort` (`id` or `source`), and `order` (`asc` or `desc`), or create a rule with `org_id`, `mailbox_id`, `destinations`, `keep_local`, and `enabled` |
| PUT | `/api/v1/mail/forwards/{id}` | Admin, user, or organization Bearer token with write scope | Update a rule using the same fields |
| DELETE | `/api/v1/mail/forwards/{id}?org_id={id}` | Admin, user, or organization Bearer token with write scope | Delete an organization-scoped forwarding rule |
| GET, POST | `/api/v1/mail/transports` | Server Admin Bearer token | List transport rules or create one with `domain`, `transport`, `next_hop`, optional username/password, and active state; passwords are encrypted at rest and never returned |
| PUT | `/api/v1/mail/transports/{id}` | Server Admin Bearer token | Update a rule; omitting `password` preserves the saved credential, while `clear_password: true` removes it |
| DELETE | `/api/v1/mail/transports/{id}` | Server Admin Bearer token | Delete a transport rule |
| GET | `/api/v1/mail/queue` | Server Admin Bearer token | List Postfix queue metadata with `limit`, `offset`, and optional `q` search; message contents are never returned |
| POST | `/api/v1/mail/queue/flushes` | Server Admin Bearer token | Flush the Postfix queue; the action is recorded in the audit log |
| PATCH | `/api/v1/mail/queue/{id}` | Server Admin Bearer token | Hold or release a queue item using `{"held":true}` or `{"held":false}` |
| DELETE | `/api/v1/mail/queue/{id}` | Server Admin Bearer token | Delete one queue item; the action is recorded in the audit log |
| POST | `/api/v1/mail/queue/{id}/retries` | Server Admin Bearer token | Retry one queue item; the action is recorded in the audit log |
| GET | `/api/v1/mail/stats?period={period}` | Server Admin Bearer token | Return delivery totals, Postfix queue size, and spam/virus counts for `24h`, `7d`, or `30d`; the default period is `24h` |
| GET | `/api/v1/mail/stats/spam?period={period}` | Server Admin Bearer token | Return spam and virus totals for `24h`, `7d`, or `30d`; the default period is `24h` |
| GET | `/api/v1/mail/services` | Server Admin Bearer token | Return detected mail service state, enabled state, version, and uptime |
| PATCH | `/api/v1/mail/services/{name}` | Server Admin Bearer token | Set `running` or `enabled` to a boolean, or set `reload` to `true`; only supported installed services may be controlled, and control requires operating-system privileges |
| GET | `/api/v1/mail/logs` | Server Admin Bearer token | Read mail delivery log envelope metadata with optional `q` search, `limit` (default 50, max 100), and `offset`; message contents are never returned |
| GET | `/api/v1/users/domains` | User Bearer token | List the caller's custom panel domains |
| POST | `/api/v1/users/domains` | User Bearer token with write scope | Enroll a custom panel domain and return its TXT verification instructions |
| GET | `/api/v1/users/domains/{domain}` | User Bearer token | Read the caller's domain state |
| DELETE | `/api/v1/users/domains/{domain}` | User Bearer token with write scope | Remove the caller's domain enrollment |
| GET | `/api/v1/users/domains/{domain}/dns` | User Bearer token | Read TXT verification instructions |
| GET | `/api/v1/users/domains/{domain}/ssl` | User Bearer token | Read certificate status and expiry |
| POST | `/api/v1/users/domains/{domain}/ssl` | User Bearer token with write scope | Request or renew the domain certificate |
| POST | `/api/v1/users/domains/{domain}/ssl/renew` | User Bearer token with write scope | Request certificate renewal |
| POST | `/api/v1/users/domains/{domain}/verify` | User Bearer token with write scope | Check the ownership TXT record |
| GET | `/api/v1/orgs/{slug}/domains` | Organization member Bearer token | List domains owned by the organization |
| POST | `/api/v1/orgs/{slug}/domains` | Organization member Bearer token with write scope | Enroll an organization custom panel domain |
| GET | `/api/v1/orgs/{slug}/domains/{domain}` | Organization member Bearer token | Read organization domain state |
| DELETE | `/api/v1/orgs/{slug}/domains/{domain}` | Organization member Bearer token with write scope | Remove the organization's domain enrollment |
| GET | `/api/v1/orgs/{slug}/domains/{domain}/dns` | Organization member Bearer token | Read TXT verification instructions |
| GET | `/api/v1/orgs/{slug}/domains/{domain}/ssl` | Organization member Bearer token | Read certificate status and expiry |
| POST | `/api/v1/orgs/{slug}/domains/{domain}/ssl` | Organization member Bearer token with write scope | Request or renew the domain certificate |
| POST | `/api/v1/orgs/{slug}/domains/{domain}/ssl/renew` | Organization member Bearer token with write scope | Request certificate renewal |
| POST | `/api/v1/orgs/{slug}/domains/{domain}/verify` | Organization member Bearer token with write scope | Check the ownership TXT record |
| GET | `/api/v1/server/{admin_path}/domains?page={page}&per_page={count}` | Admin Bearer token | List custom panel domains, with pagination (maximum 100 per page) |
| GET | `/api/v1/server/{admin_path}/domains/{domain}` | Admin Bearer token | Read domain state |
| DELETE | `/api/v1/server/{admin_path}/domains/{domain}` | Admin Bearer token with write scope | Force-delete a domain enrollment |
| POST | `/api/v1/server/{admin_path}/domains/{domain}/suspend` | Admin Bearer token with write scope | Suspend a domain with a required reason |
| POST | `/api/v1/server/{admin_path}/domains/{domain}/unsuspend` | Admin Bearer token with write scope | Restore a verified domain when required SSL is active |
| POST | `/api/v1/server/{admin_path}/domains/{domain}/ssl/renew` | Admin Bearer token with write scope | Force certificate renewal for any verified domain |

`/server/healthz` is the canonical public health page and negotiates HTML, human-readable text, or JSON. `/api/healthz` is the unversioned machine-readable alias; API health returns JSON by default and ordered dotted-key text for `Accept: text/plain`. `/healthz` is mounted only when `server.healthz.root.enabled` is true and uses the same handler as `/server/healthz`. Health status is one of `healthy`, `degraded`, `restart_required`, `unhealthy`, `maintenance`, or `shutting_down`; unhealthy, maintenance, and shutting-down responses use HTTP 503, while the other states use HTTP 200. Public health includes version/build/runtime information, cluster and feature state, component checks, and aggregate request counts. Health probes and `/server/metrics` are excluded from the persisted request totals; the active-connection count includes the health request itself. Setup verification accepts a URL-encoded `token` field and does not complete setup; the browser setup wizard submits completion to its native server-rendered POST route.

The JSON health response includes `public_ips.ipv4`, `public_ips.ipv6`, and `public_ips.last_checked`. IPv4 is discovered through an external address lookup and IPv6 from an active interface with a global unicast address. Discovery runs at startup and every 12 hours; either address is omitted when unavailable.

Admin profile routes accept only `Authorization: Bearer {token}` and only allow the matching admin username owned by that token. Cookie sessions and CSRF headers are not API credentials. Admin API tokens are SHA-256 hashed at rest and are displayed only when first issued or regenerated.

Alias requests also accept only Bearer tokens. User tokens must provide an `org_id` they belong to; organization tokens are bound to their own organization, and administrators can select any organization. Alias sources must belong to the selected domain, while destinations may be local or external email addresses. `read` tokens can only list aliases, `read-write` tokens can create and update, and `global` tokens can delete. Credentials and password hashes are never returned.

Example:

```sh
curl -H 'Accept: application/json' http://localhost:64580/api/v1/server/healthz
```

Custom-domain endpoints are available when the application database is configured; the legacy `server.features.custom_domains.enabled` value does not gate them. Certificates use encrypted database-backed storage. HTTP-01 uses an HTTP listener and TLS-ALPN-01 uses a TLS listener. Wildcards require DNS-01. A custom-domain SSL POST may include `{"dns_provider":"cloudflare","dns_credentials":{"CF_DNS_API_TOKEN":"..."}}`; credentials are encrypted before storage and are omitted from responses. Omit `dns_provider` to issue using the domain's currently saved challenge settings. Other routes in the AI.md API inventory remain incomplete.

The panel SSL PATCH accepts `enabled`, `email`, `challenge`, `staging`, `dns_provider`, and optional `dns_credentials` as a flat JSON object. Omitting `dns_credentials` preserves the currently saved credentials. Responses show only masked values. Changes require a server restart.

The public security status page is `/server/security/report/{tracking_id}` and uses a daily-use token in its URL; there is intentionally no API mirror. The security overview and machine-readable `/.well-known/security.txt` expose the rotating one-shot contact ID. Security report bodies and researcher contact details are encrypted in storage.

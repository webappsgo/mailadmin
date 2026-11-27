# Security

## Identity and sessions

Server Admin and regular user identities are separate. Password verification and session lifecycle belong to the authentication service. Panel passwords use Argon2id; setup and API tokens are stored as SHA-256 hashes. Session cookies are HTTP-only, SameSite-protected, and Secure when HTTPS is active. State-changing browser forms include CSRF validation. Native API routes accept Bearer tokens only.

Protected routes fail closed if authentication is not configured. Mailbox, organization, and domain persistence operations receive the authenticated principal and scope; resource access must be checked in the persistence layer as well as the route layer.

## Login throttling

Rate limiting is always enabled. Read, write, health, and global-burst requests are counted per source IP or valid Bearer credential hash. Token and restore operations use the strict policy, keyed per owner for authenticated requests and per IP otherwise. Login requests are limited per IP and account identifier, and failed password attempts or invalid Bearer credentials lock the source IP for the configured authentication-failure window. The lockout applies to subsequent requests from that IP until the sliding-window count falls below the limit. Limited requests return HTTP 429 with a `Retry-After` value. Counters use the server database and switch to process-local counters with one warning if that counter store is unavailable; the fallback resets when the process restarts. Authentication failures do not reveal whether an account exists.

## Public endpoints and transport

Health, version, about, static assets, and public information pages are unauthenticated. Use TLS at the application listener or a trusted reverse proxy. The server emits a content security policy and standard browser security headers; HSTS is sent only for TLS requests.

## Sensitive data

Never log passwords, session credentials, setup tokens, API tokens, private keys, or database credentials. Setup tokens are single-use. API tokens are shown only when first issued or regenerated. Backup encryption passwords are stored as AES-GCM ciphertext in the owner-only `secrets.yml` file and decrypted in memory only for scheduled backups.

## Reporting

Report repository vulnerabilities through the private vulnerability reporting URL in `web.security.report_url`. Deployed-instance reports use `/server/security` and the one-time contact link in `/.well-known/security.txt`. The report body is encrypted with the project OpenPGP key when one is active; otherwise it uses AES-256-GCM with `server.security.encryption_key`. Each Server Admin can register a personal OpenPGP public key on their profile to receive encrypted report notices. When no admin key is registered, the configured admin contact receives an AES-256-GCM encrypted attachment. Researcher acknowledgments contain only the tracking ID and status URL and are encrypted when a researcher public key is supplied. The public policy and acknowledgment pages are `/server/security/policy` and `/server/security/thanks`; Server Admins edit the policy and manage acknowledgment visibility, credit labels, and display order under `/server/{admin_path}/security/`. Public key material is served at `/.well-known/pgp-key.asc` only while an active project key exists. The Server Admin keypair page is `/server/{admin_path}/security/keypair`; sensitive actions require password re-entry, private-key export requires a reason and is limited to one download per admin per hour, and keyserver submissions retry on the scheduler. Include affected version, impact, and a safe reproduction without sharing real credentials or customer data.

# SPEC.md — Project-Specific Rule Overrides

> **Rule hierarchy:** SPEC.md > AI.md > global CLAUDE.md. Rules here win over AI.md.

## Platform Target

**Override:** AI.md §PART 1 specifies 8 platforms (linux, darwin, windows, freebsd × amd64, arm64).
This project targets **Linux and BSD only**.

| Rule | AI.md default | This project |
|------|--------------|--------------|
| Supported OS | linux, darwin, windows, freebsd | linux, freebsd |
| Supported arch | amd64, arm64 | amd64, arm64 |
| Total platforms | 8 | 4 |
| Windows `.exe` suffix | required | N/A — Windows excluded |
| macOS (darwin) builds | required | N/A — macOS excluded |

**Rationale:** mailadmin manages Postfix, Dovecot, SpamAssassin, ClamAV, and related services
that run exclusively on Linux and BSD systems. macOS and Windows are not valid deployment targets.

### Build Matrix

```
linux/amd64
linux/arm64
freebsd/amd64
freebsd/arm64
```

`make build` produces 4 binaries. `make release` publishes 4 platform archives.
Binary naming follows AI.md: `mail-{os}-{arch}` (no `.exe`).

## Optional Part Activation

**Activated:** PART 34 (Multi-User), PART 35 (Organizations), PART 36 (Custom Domains).
All three are NON-NEGOTIABLE per AI.md once activated — full compliance with each PART is required.

| PART | Status | Scope for mailadmin |
|------|--------|----------------------|
| 34 — Multi-User | Activated | Regular User accounts (self-service mail users): change own password, set own vacation message, view own quota. Separate from Server Admin (PART 17). |
| 35 — Organizations | Activated | Requires PART 34. Orgs are the tenancy unit — each Org owns one or more mail Domains (`VirtualDomain`) and has member Users who manage only their Org's domains. Server Admin retains full cross-tenant control. |
| 36 — Custom Domains | Activated | Two applications, both in scope: (1) mail domains (`VirtualDomain`) are Org/User-owned per the PART 35 tenancy model; (2) the mailadmin panel UI itself may be served under an Org's own custom domain, with automatic Let's Encrypt (HTTP-01/TLS-ALPN-01/DNS-01) certificate issuance and TXT-record ownership verification, per PART 36's standard pattern. |

**Rationale:** mailadmin is becoming a multi-tenant mail hosting platform — Organizations
sign up, own their mail domains, and manage their own Users, rather than a single Server
Admin manually provisioning every domain and mailbox. See IDEA.md → Business logic for the
full Roles & permissions, Data model, Trust boundaries, and Threat model detail this
activation drives.

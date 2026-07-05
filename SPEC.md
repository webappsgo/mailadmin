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

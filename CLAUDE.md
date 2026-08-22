# mailadmin — Claude Code Loader

Rule hierarchy: **SPEC.md > AI.md > global CLAUDE.md**. AI.md is read-only.

| File | Role |
|------|------|
| `AI.md` | HOW — full implementation spec (PART 0-37). Read the relevant PART before implementing. |
| `IDEA.md` | WHAT — project description, variables, business logic. |
| `SPEC.md` | Project overrides — 4-platform build matrix (linux/freebsd × amd64/arm64); PART 34/35/36 activated. |

## Non-negotiables

- `CGO_ENABLED=0`; build only via `make dev` / `make local` / `make build` (Docker, never host Go).
- Dockerfile lives in `docker/Dockerfile`, never the repo root. No `.env` files.
- Passwords: Argon2id. Tokens: SHA-256 hashed. Never plaintext.
- Booleans parse via `config.ParseBool()`, never `strconv.ParseBool()`.
- Singular package directories (`handler/`, `model/`), server-side rendering only.
- No feature gating, no premium tiers, no usage quotas.

For anything else, read the matching PART in `AI.md`.

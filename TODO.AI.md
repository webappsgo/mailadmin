# TODO

- [ ] Write `*_test.go` files across every package in `src/` to reach the
      100% coverage required by `make test` (Makefile:214-224). Currently
      zero test files exist anywhere in the project — `make test` fails
      at 0.0% coverage. This blocks the commit test gate for any future
      code change; the 2026-08-22 docs/compliance commit was made with
      this gap explicitly logged rather than skipped silently, at the
      user's direction, because the diff was docs/config-only and did
      not touch `src/`.

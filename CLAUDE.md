# Project SPEC

This project uses AI.md as its specification.

## Files

- **AI.md** — Implementation spec (PARTS 0-36). READ-ONLY. Never modify.
- **IDEA.md** — Business logic, features, data models. Update as project evolves.
- **.claude/rules/** — Rule cheatsheets auto-loaded each session.

## Session Start

1. Read IDEA.md completely
2. Load .claude/rules/ files relevant to the current task
3. Read AI.md PARTs relevant to the current task (do NOT pre-load all PARTs)
4. Never modify AI.md

## Project

- **project_name**: mailadmin
- **project_org**: webappsgo
- **internal_name**: mail (FROZEN)
- **module**: github.com/webappsgo/mail
- **language**: Go
- **binary**: mail

## Rules

All work follows AI.md. When in doubt: read the relevant PART, then implement.

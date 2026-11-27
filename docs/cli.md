# CLI and TUI

`mailadmin-cli` is the companion client for API operations. It supports JSON,
tabular, and plain output for scripted commands, and a keyboard driven TUI when
started without a command in an interactive terminal. The TUI adapts at 80,
60, and 40 columns and asks for a wider terminal below 40 columns.

## Connect and authenticate

Start the TUI with `mailadmin-cli`, or configure a server directly:

```sh
mailadmin-cli --server https://mail.example.net status
mailadmin-cli login
```

The login flow checks the bearer token against the server before saving it.
The config file is `~/.config/webappsgo/mail/cli.yml` on Unix and is created
with mode `0600`; the client refuses config or token files that are readable
by group or other users. Tokens can also come from `--token`, `--token-file`,
`MAILADMIN_TOKEN`, or the config file, in that order. `logout` removes the
saved credential.

The TUI works with arrow keys or `j`/`k`, Enter, Escape, and Ctrl+C. Its login
screen masks token input. From the menu, select a command and enter the
arguments shown on screen. Select **Quit** or press Ctrl+C to leave.

## Commands

```text
mailadmin-cli status
mailadmin-cli test
mailadmin-cli list RESOURCE
mailadmin-cli get RESOURCE ID
mailadmin-cli create RESOURCE FILE
mailadmin-cli update RESOURCE ID FILE
mailadmin-cli delete RESOURCE ID
mailadmin-cli logout
```

Create and update read JSON request bodies from the named file. `--output`
accepts `json`, `table`, or `plain`; `table` prints object collections as
columns and `plain` prints object fields as key/value lines. `--color` accepts
`auto`, `yes`, or `no`.
Use `--user @NAME` for a user path and `--user +NAME` for an organization
path. The server still enforces authorization and tenant ownership.

The client reads cluster URLs from the server health response and caches them
in `cli.yml`. It refreshes the list at most every 30 minutes and silently tries
the cached nodes when the current primary cannot be reached. A successful
failover is saved as the new primary.

## Display mode

`display.mode` in `cli.yml` accepts `auto`, `tui`, or `gui`. The default static
release includes the TUI. A build without native GUI support reports
`gui_unavailable` when `gui` is forced. `--help` and `--version` always return
immediately, and supplying a command always selects command-line mode.

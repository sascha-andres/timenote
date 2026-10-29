# timenote

timenote is a tool to take notes with timestamps using Toggl as a backend.

Essentially this is a commandline client to track your time

## Configuration

Configuration is via command-line flags and environment variables only (no
config file). Every flag falls back to an upper-cased, dash-to-underscore
environment variable with the prefix `TIMENOTE_`, e.g. `-workspace` /
`TIMENOTE_WORKSPACE`, `-cache-path` / `TIMENOTE_CACHE_PATH`.

Flags must come before the command, not after, e.g.:

    timenote -description "did a thing" timestamp new
    timenote -name "Client X" projects create

## Backend

`-backend` selects where time entries go (default `toggl`):

- `toggl`: talks to toggl directly. Uses the token from the keyring (see
  below), `-workspace` and the local cache.
- `remote`: talks to a `backend-toggl` server (see
  `cmd/backend-toggl/README.md`). `-remote-url` (for example
  `http://127.0.0.1:8141`) and `-remote-token` are required. Environment:
  `TIMENOTE_BACKEND`, `TIMENOTE_REMOTE_URL`, `TIMENOTE_REMOTE_TOKEN`.
  timenote uses no cache: `-cache-path` and `-cache-max-age` have no effect,
  and `timenote cache` stops with an error.

    timenote -backend remote -remote-url http://127.0.0.1:8141 -remote-token xxx today

## Store your toggl token

    timenote -token xxx token save

## Wizard

    timenote wizard

The wizard asks step by step what to do with the time entry. Use the arrow
keys to move, `/` to filter a list, `enter` to select and `ctrl+c` (or `q` in a
list) to cancel. Cancel makes no changes.

- `new`: select a client, then a project of that client. Projects without a
  client are under `(no client)`. Then select if you want to add a
  description. `yes` asks for the description. The wizard starts a new time
  entry with the project and the description.
- `change`: if no time entry runs, the wizard prints an error and exits.
  - `append description`: asks for text and appends it to the description,
    separated by `-separator`.
  - `overwrite description`: asks for text and replaces the description.
  - `new project`: select a client and a project. The wizard sets the project
    on the running time entry.
- `stop`: stops the running time entry.

The client and project lists come from the cache (`-cache-max-age`,
`timenote cache update`). `-excluded-projects` also applies to the wizard.

## Logging

`-log-level` sets the minimum level printed (`debug`, `info`, `warn`,
`error`; default `warn`). `-log-json` switches the log output to JSON.

## State

Used regularly with the toggl backend.

## Development

### Dependencies with go modules

## History

|Version|Description|
|---|---|
|0.10.0|Module renamed to go.livingit.de/timenote|
||Replace cobra/viper with go.livingit.de/reuse/flag (flags must precede the command)|
||Replace bbolt cache with plain JSON files|
||Switch from log to log/slog, add -log-level and -log-json|
||Log the toggl API method and duration for every call|
||Drop interactive mode|
||Replace ginkgo/gomega tests with standard library tests|
||Update dependencies|
||Add wizard command|
||Fix duration of the running time entry (toggl api v9)|
|0.9.1|Print total for today|
|0.9.0|Swap out logrus for the standard library logger|
||Make timestamp current the default behavior|
||Add keyring-backed token storage|
|0.8.1|update Dependencies|
|0.8.0|Add caching layer|
|0.7.0|Better formating of time values|
||Make append separator configurable|
||Display client for time entry|
||Display project for time entry|
|0.6.0|add support for projects (add,delete and list|
||remove MySQL support|
||Workspace flag|
||Grouping for timenote today|
||reduce code complexity|
|0.5.0|add support for daily summary|
||better output for entry duration|
|0.4.0|add project command|
|0.3.0|Add cli|
||Open in browser|
|0.2.0|Add projects|
|0.1.0|Initial version|
|0.2.0|Add support for projects|

# timenote

timenote is a tool to take notes with timestamps using Toggl as a backend.

Essentially this is a commandline client to track your time

## Configuration

Configuration is via command-line flags and environment variables only (no
config file). Every flag falls back to an upper-cased, dash-to-underscore
environment variable, e.g. `-workspace` / `WORKSPACE`, `-cache-path` /
`CACHE_PATH`.

Flags must come before the command, not after, e.g.:

    timenote -description "did a thing" timestamp new
    timenote -name "Client X" projects create

## Store your toggl token

    timenote -token xxx token save

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

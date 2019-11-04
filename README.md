# timenote

timenote is a tool to take notes with timestamps using Toggl as a backend.

Essentially this is a commandline client to track your time

## Configuration

There is one main option:

1. dsn

For toggl use your toggl token.

You can Use a configuration file in your home directory, .timenote.yaml.

A sample looks like this:

    ---
    dsn: /

## State

This is in an early phase but used regularly with the toggl backend and from time to time with MySQL.

## Development

### Dependencies with go modules

## History

|Version|Description|
|---|---|
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

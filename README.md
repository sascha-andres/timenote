# timenote

timenote is a tool to take notes with timestamps. You can choose between two different backends:

* MySQL/MariaDB
* Toggl

## Configuration

There are two main options:

1. persistor
2. dsn

The first is the type of backend and can be either mysql or toggl, the latter one provides the information to connect.

For MySQL you can must use a valid connection string ( see https://github.com/go-sql-driver/mysql ). For toggl use your toggl token.

You can Use a configuration file in your home directory, .timenote.yaml.

A sample looks like this:

    ---
    dsn: /
    persistor: mysql

## State

This is in an early phase but used regularly with the toggl backend and from time to time with MySQL.

## Development

### Dependencies with dep

https://github.com/golang/dep
# backend-toggl

backend-toggl makes the toggl `timenote.Backend` available as a REST API.
Use the `remote` package as the Go client.

## Start

    go install go.livingit.de/timenote/cmd/backend-toggl@latest
    BACKEND_TOGGL_TOKEN=secret backend-toggl

To start it at login with systemd or launchd, see
[examples/launch](../../examples/launch/README.md).

## Configuration

Configuration is with command-line flags and environment variables
(`go.livingit.de/reuse/flag`). The flag has priority.

| Flag              | Environment variable           | Default                        | Description                                           |
|-------------------|--------------------------------|--------------------------------|-------------------------------------------------------|
| `-listen`         | `BACKEND_TOGGL_LISTEN`         | `127.0.0.1:8141`               | Address to listen on                                  |
| `-token`          | `BACKEND_TOGGL_TOKEN`          | none, required                 | Token that clients must send                          |
| `-toggl-token`    | `TOGGL_TOKEN`                  | timenote keyring entry         | Toggl API token (`timenote -token xxx token save`)    |
| `-workspace`      | `BACKEND_TOGGL_WORKSPACE`      | `0`                            | Toggl workspace. `0` uses the first workspace         |
| `-cache-path`     | `BACKEND_TOGGL_CACHE_PATH`     | `$HOME/.config/timenote/cache` | Cache directory (same as the timenote CLI)            |
| `-cache-max-age`  | `BACKEND_TOGGL_CACHE_MAX_AGE`  | `360`                          | Maximum age of the cache in minutes                   |
| `-read-cache-ttl` | `BACKEND_TOGGL_READ_CACHE_TTL` | `30s`                          | How long the server keeps `GET` results. `0` disables |

`-h` shows the values of the environment variables as defaults, including the
tokens. Do not share `-h` output when the token variables are set.

The server uses plain HTTP. Keep it on localhost, or put a TLS proxy in front
of it.

## Read cache

The server keeps the result of each `GET` endpoint in memory for
`-read-cache-ttl`. This decreases the calls to the Toggl API when clients poll.
Each other request (`POST`, `PUT`, `PATCH`, `DELETE`) clears the read cache.
Errors are not kept, except "no current time entry".

Changes that you make outside of the server (for example in the Toggl web
app) are visible after a maximum of `-read-cache-ttl`.

## Authentication

Each request must have this header:

    Authorization: Bearer <token>

If the header is missing or wrong, the server sends `401`.

## Endpoints

All request and response bodies are JSON. Field names are snake_case.
Times are RFC 3339.

| Method | Path                                  | Body                                         | Result                          |
|--------|---------------------------------------|----------------------------------------------|---------------------------------|
| POST   | `/v1/timeentries`                     | `{"from_previous": false}`                   | Starts a new time entry. `true` copies the last finished entry |
| GET    | `/v1/timeentries`                     |                                              | Time entries of today           |
| GET    | `/v1/timeentries/current`             |                                              | Running time entry              |
| PATCH  | `/v1/timeentries/current`             | `{"stop": true}`                             | Stops the running time entry    |
| PUT    | `/v1/timeentries/current/description` | `{"description": "text"}`                    | Replaces the description        |
| PATCH  | `/v1/timeentries/current/description` | `{"line": "text", "separator": ";"}`         | Appends a line to the description |
| PATCH  | `/v1/timeentries/current/tags`        | `{"toggle": "name"}`                         | Adds the tag, or removes it if it is set |
| PUT    | `/v1/timeentries/current/project`     | `{"name": "p", "auto_create_project": false}`| Sets the project. Empty `name` removes it |
| GET    | `/v1/projects`                        |                                              | Projects                        |
| POST   | `/v1/projects`                        | `{"name": "p"}`                              | Creates the project if it does not exist |
| DELETE | `/v1/projects/{nameOrID}`             |                                              | Deletes the project             |
| GET    | `/v1/clients`                         |                                              | Clients                         |
| POST   | `/v1/clients`                         | `{"name": "c"}`                              | Creates the client if it does not exist |
| GET    | `/v1/cache`                           |                                              | Status of the projects and clients cache |
| POST   | `/v1/cache`                           |                                              | Reloads the cache from Toggl    |

If the table shows a body, the body is required (`{}` is valid for `POST
/v1/timeentries`). Calls without a result send `204 No Content`.

Example:

    curl -H "Authorization: Bearer secret" http://127.0.0.1:8141/v1/timeentries/current

## Errors

Errors are RFC 9457 Problem Details (`application/problem+json`) with the
extension `code`:

```json
{
  "type": "about:blank",
  "title": "No current time entry",
  "status": 404,
  "detail": "unable to get running time entry from toggl: no current time entry",
  "code": "NO_CURRENT_ENTRY"
}
```

| Status | Code                | Cause                                         |
|--------|---------------------|-----------------------------------------------|
| 400    | `INVALID_BODY`      | The body is not valid JSON or has unknown fields |
| 401    | `UNAUTHORIZED`      | The token is missing or wrong                 |
| 404    | `NO_CURRENT_ENTRY`  | No time entry runs                            |
| 404    | `PROJECT_NOT_FOUND` | The project does not exist                    |
| 422    | `VALIDATION_FAILED` | A required field is empty (`validation` names the field) |
| 502    | `BACKEND_ERROR`     | Other error from the toggl backend            |

## Go client

```go
c := remote.New("http://127.0.0.1:8141", "secret")
if err := c.New(); err != nil {
	var problem *remote.Problem
	if errors.As(err, &problem) && problem.Status == http.StatusUnauthorized {
		// wrong token
	}
}
```

`*remote.Client` implements `timenote.CachingBackend`.

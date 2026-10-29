# Start backend-toggl automatically

Examples to start `backend-toggl` when you log in. Both run it as your user,
so it uses your timenote cache and keyring entry.

Install the binary first:

    go install go.livingit.de/timenote/cmd/backend-toggl@latest

The binary goes to `$(go env GOPATH)/bin` (default `~/go/bin`). If your path is
different, change it in the example.

## Linux (systemd)

Files: `backend-toggl.service`, `backend-toggl.env`.

    mkdir -p ~/.config/systemd/user ~/.config/timenote
    cp backend-toggl.service ~/.config/systemd/user/
    install -m 600 backend-toggl.env ~/.config/timenote/backend-toggl.env
    # set BACKEND_TOGGL_TOKEN in ~/.config/timenote/backend-toggl.env
    systemctl --user daemon-reload
    systemctl --user enable --now backend-toggl

Status and logs:

    systemctl --user status backend-toggl
    journalctl --user -u backend-toggl -f

The keyring needs a running desktop session (Secret Service). On a server
without a desktop session, set `TOGGL_TOKEN` in the env file. To start the
service without a login, run `loginctl enable-linger $USER`.

## macOS (launchd)

File: `de.livingit.backend-toggl.plist`.

    sed "s/USERNAME/$USER/g" de.livingit.backend-toggl.plist > ~/Library/LaunchAgents/de.livingit.backend-toggl.plist
    chmod 600 ~/Library/LaunchAgents/de.livingit.backend-toggl.plist
    # set BACKEND_TOGGL_TOKEN in the copied file
    launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/de.livingit.backend-toggl.plist

Status and logs:

    launchctl print gui/$(id -u)/de.livingit.backend-toggl
    tail -f ~/Library/Logs/backend-toggl.log

Stop and remove:

    launchctl bootout gui/$(id -u)/de.livingit.backend-toggl
    rm ~/Library/LaunchAgents/de.livingit.backend-toggl.plist

After you change the plist, run `bootout` and `bootstrap` again.

# svcs

Write Linux service installers.

## What It Does

`svcs` is a CLI that builds a Go service binary, writes a matching systemd unit under `/opt/gosvc`, and starts or removes that service with `systemctl`.

## CLI Flags

- `-service`: service name and source folder
- `-runas`: user account for the systemd service
- `-remove`: stop, disable, and remove the service files

## Notes

The tool shells out through `sudo`, `go build`, and `systemctl`, so use it in controlled operator environments rather than as a library.

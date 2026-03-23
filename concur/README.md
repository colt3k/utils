# Concurrency Tools

Helpers for bounded concurrency and one-time initialization that can be retried or reset.

## Included Components

- `ResyncOnce`
- concurrency pool of workers

## Pool Usage

Create `Task` values with `NewTask`, build a `Pool` with `NewPool` or `NewPoolWithPause`, and call `Run` to execute work at the configured concurrency.

## Examples

See the end-to-end example under `concur/pool_test`. It demonstrates task creation, response handling, and result capture.

## Development

- `cd concur && go test ./...`

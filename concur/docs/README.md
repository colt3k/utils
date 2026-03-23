# concur Documentation

These notes split the repository-wide docs down to the `concur` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/concur`
- Run `cd concur && go test ./...` from the module root.
- Use `concur/pool_test/main.go` as the reference example for task and response handling.
- Keep task side effects isolated so pool behavior stays easy to reason about.

## API Surface
- `Pool`, `Task`, and `TaskResponseReturn` define the worker-pool core.
- `NewPool` and `NewPoolWithPause` configure concurrency and optional pacing.
- `NewTask` wraps a work function and optional return sink.
- `ResyncOnce` provides resettable one-time execution semantics.

## Schema & Data Shapes
- No external file schema is used.
- The key runtime shapes are `Task`, `Pool`, and any caller-defined response object implementing `TaskResponseReturn`.

## Dataflow
1. Wrap each work item in a `Task`.
2. Create a `Pool` with the target concurrency.
3. Call `Run` and collect per-task results through the supplied return handler.

## Operator Notes
- Library-only module; the main operational concern is safe concurrency rather than deployment.
- Verify pacing and cancellation behavior when changing `NewPoolWithPause` or the worker loop.

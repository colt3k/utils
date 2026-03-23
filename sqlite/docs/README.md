# sqlite Documentation

These notes split the repository-wide docs down to the `sqlite` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/sqlite`
- Run `cd sqlite && go test -mod=mod ./...` in this repository state because vendoring is out of sync.
- Remember that `DB(...)` is guarded by `sync.Once`; the first database path passed in wins for the process lifetime.
- Keep query result mapping stable because callers expect `[]map[string]interface{}` rows.

## API Surface
- `DB(args ...string)` returns the shared `DBCon` instance.
- `(*DBCon).Execute(query)` runs non-row statements.
- `(*DBCon).Query(query)` returns rows as column-name maps.
- `(*DBCon).Close()` closes the underlying database handle.

## Schema & Data Shapes
- There is no separate config file schema in this module.
- The main runtime shape is the SQLite database file path supplied to the first `DB(...)` call.

## Dataflow
1. Open the database once through `DB(dbPath)`.
2. Execute DDL or DML with `Execute` and row-returning statements with `Query`.
3. Close the connection when the process no longer needs the database.

## Operator Notes
- This is a library wrapper, but the operator concern is the database file path and lifecycle.
- If callers connect to the wrong database, inspect the first `DB(...)` call in the process before debugging anything else.

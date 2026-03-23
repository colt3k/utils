# sqlite

Provides an easy-to-use library for sqlite3 on top of `modernc.org/sqlite`.

## Easy Calls For

- Create
  - `sqlite.DB(dbPath)`
- Close
  - `sqlite.DB().Close()`
- Execute
  - `sqlite.DB().Execute("drop table if exists tests;")`
- Query
  - returns an array of column-name maps
  - `dat, err := sqlite.DB().Query("select id,name from tests order by id;")`

## Query Mapping

- `int` -> `int64`
- `float` -> `float64`
- date as numeric -> `int64`, `string`
- `blob` -> stored as input
- `null` -> `nil`

## Notes

`DB` is a singleton accessor guarded by `sync.Once`, so the first path passed to `sqlite.DB(...)` wins for the lifetime of the process.

## Development

- `cd sqlite && go test ./...`

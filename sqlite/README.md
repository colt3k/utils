sqlite

Provides an easy-to-use library for sqlite3 on top of modernc.org/sqlite

- Easy calls for
  - Create
    - sqlite.DB(dbPath)
  - Close
    - sqlite.DB().Close()
  - Execute
    - sqlite.DB().Execute("drop table if exists tests;")
  - Query
    - Returns an array map of keyed columns as strings and values as interfaces
    - dat, err := sqlite.DB().Query("select id,name from tests order by id;")
    - Types mapped to Go
      - int -> int64
      - float -> float64
      - date as numeric -> int64, string
      - blob -> stored as input
      - null -> nil
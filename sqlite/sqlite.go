// Package sqlite provides a thin singleton wrapper around modernc.org/sqlite.
package sqlite

import (
	"database/sql"
	"fmt"
	log "github.com/colt3k/nglog/ng"
	_ "modernc.org/sqlite"
	"sync"
)

var once sync.Once

var (
	instance *DBCon
)

type DBCon struct {
	db *sql.DB
}

// DB returns the shared database connection, creating it on the first call.
func DB(args ...string) *DBCon {
	once.Do(func() { // <-- atomic, does not allow repeating

		log.Logf(log.INFO, "db created at: %v", args[0])
		instance = new(DBCon) // <-- thread safe
		var err error
		instance.db, err = sql.Open("sqlite", args[0])
		if err != nil {
			log.Logf(log.FATAL, "issue creating db: %v", err)
		}
	})

	return instance
}

// Close closes the underlying database handle.
func (d *DBCon) Close() error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute runs a statement that does not return rows.
func (d *DBCon) Execute(query string) error {
	if d.db != nil {
		_, err := d.db.Exec(query)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("db obj is null")
}

// Query runs a query and returns rows as column-name maps.
func (d *DBCon) Query(query string) ([]map[string]interface{}, error) {
	if d.db != nil {
		dest := make([]map[string]interface{}, 0)
		log.Logf(log.DEBUG, "Query: %s\n", query)
		rows, err := d.db.Query(query)
		if err != nil {
			return dest, err
		}
		for rows.Next() {
			row := make(map[string]interface{}, 0)
			cols, er2 := rows.Columns()
			if er2 != nil {
				log.Logf(log.ERROR, "ERROR issue on cols: %v\n", er2)
			}

			fields := make([]interface{}, len(cols))
			for l := range cols {
				fields[l] = new(interface{})
			}

			if err = rows.Scan(fields...); err != nil {
				log.Logf(log.ERROR, "ERROR issue on scan: %v\n", err)
			}

			for l, m := range cols {
				row[m] = *(fields[l].(*interface{}))
			}

			dest = append(dest, row)
		}

		if err = rows.Err(); err != nil {
			return dest, err
		}
		return dest, nil
	}
	return nil, fmt.Errorf("db obj is null")
}

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

func DB(args... string) *DBCon {
	once.Do(func() { // <-- atomic, does not allow repeating

		log.Logf(log.INFO, "db created at: %v",args[0])
		instance = new(DBCon) // <-- thread safe
		var err error
		instance.db, err = sql.Open("sqlite", args[0])
		if err != nil {
			log.Logf(log.FATAL, "issue creating db: %v",err)
		}
	})

	return instance
}

func (d *DBCon) Close() error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DBCon) Execute(query string) error {
	if d.db != nil {
		_,err := d.db.Exec(query)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("db obj is null")
}

func (d *DBCon) Query(query string) ([]map[string]interface{},error) {
	if d.db != nil {
		dest := make([]map[string]interface{},0)
		log.Logf(log.DEBUG,"Query: %s\n", query)
		rows, err := d.db.Query(query)
		if err != nil {
			return dest,err
		}
		for rows.Next() {
			row := make(map[string]interface{},0)
			cols, er2 := rows.Columns()
			if er2 != nil {
				log.Logf(log.ERROR,"ERROR issue on cols: %v\n",er2)
			}

			fields := make([]interface{},len(cols))
			for l := range cols {
				fields[l]=new(interface{})
			}

			if err = rows.Scan(fields...); err != nil {
				log.Logf(log.ERROR,"ERROR issue on scan: %v\n",err)
			}

			for l,m := range cols {
				row[m]=*(fields[l].(*interface{}))
			}

			dest=append(dest,row)
		}

		if err = rows.Err(); err != nil {
			return dest,err
		}
		return dest,nil
	}
	return nil,fmt.Errorf("db obj is null")
}
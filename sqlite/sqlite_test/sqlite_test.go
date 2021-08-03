package sqlite_test

import (
	"fmt"
	"github.com/colt3k/utils/debug"
	"github.com/colt3k/utils/sqlite"
	"os"
	"testing"
)

var(
	dbPath = "test.db"
)

func TestName(t *testing.T) {

	sqlite.DB(dbPath)
	defer sqlite.DB().Close()
	err := sqlite.DB().Execute("drop table if exists tests;")
	if err != nil {
		debug.PrintStack()
		t.Fatalf("error dropping table %v", err)
	}
	err = sqlite.DB().Execute(`create table if not exists tests(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		myvalues BLOB,
		money FLOAT,
		curdate DATE);`)
	if err != nil {
		debug.PrintStack()
		t.Fatalf("error creating table %v", err)
	}
	err = sqlite.DB().Execute("insert into tests values(42,'a','someA',0.50,null), (314,'b','someB',0.75,1628022195);")
	if err != nil {
		debug.PrintStack()
		t.Fatalf("error inserting rows %v", err)
	}
	dat, err := sqlite.DB().Query("select id,name,myvalues,money,curdate from tests order by id;")
	if err != nil {
		debug.PrintStack()
		t.Fatalf("error querying %v", err)
	}
	for l,m := range dat {
		fmt.Printf("%d. %v\n", l,m)
		if l == 0 && m["id"].(int64) != 42{
			t.Errorf("incorrect value %v at position %d", m["id"], l)
		}
		if l == 0 && m["name"] != "a" {
			t.Errorf("incorrect value %v at position %d", m["name"], l)
		}
	}

	fi, err := os.Stat(dbPath)
	if err != nil {
		debug.PrintStack()
		t.Fatalf("error obtaining stats on db %v", err)
	}

	fmt.Printf("%s size: %v\n", dbPath, fi.Size())
}
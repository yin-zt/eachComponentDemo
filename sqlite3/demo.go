package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/common/log"

	"sync"
)

var (
	sqliteCache = &SQLiteCache{}
)

type SQLiteCache struct {
	sync.Mutex
	SQLite *sql.DB
}

func main() {
	if db, err := sql.Open("sqlite3", ":memory:"); err == nil {
		db.SetMaxOpenConns(1)
		sqliteCache.SQLite = db
	} else {
		log.Error(err.Error())
	}
	// drop table
	res1, err1 := sqliteCache.SQLite.Exec(fmt.Sprintf("DROP TABLE %s", "test"))
	fmt.Println(res1, err1)
	query := "CREATE TABLE %s (" + "col1,col2" + ")"
	query = fmt.Sprintf(query, "test")
	res, err := sqliteCache.Exec(query)
	fmt.Println(res, err)

	//create table

}

func (s SQLiteCache) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.Lock()
	defer s.Unlock()
	if len(args) > 0 {
		return s.SQLite.Exec(query, args)
	} else {
		return s.SQLite.Exec(query)
	}

}

func (s SQLiteCache) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s.Lock()
	defer s.Unlock()
	if len(args) > 0 {
		return s.SQLite.Query(query, args)
	} else {
		return s.SQLite.Query(query)
	}

}

func (s SQLiteCache) Prepare(query string) (*sql.Stmt, error) {
	s.Lock()
	defer s.Unlock()
	return s.SQLite.Prepare(query)

}

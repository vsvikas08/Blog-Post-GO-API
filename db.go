package main

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type SqlDB struct {
	db *sql.DB
}

var lock = &sync.Mutex{}
var sqlDB *SqlDB

func getDBInstance() *SqlDB {
	if sqlDB == nil {
		lock.Lock()
		defer lock.Unlock()
		if sqlDB == nil {
			sqlDB = &SqlDB{}
			var err error
			sqlDB.db, err = sql.Open("mysql","root:@tcp(localhost:3306)/go_blog")
			if err != nil {
				panic(err)
			}
		}
	}
	return sqlDB
}
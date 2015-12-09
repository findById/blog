package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetDBM(db string) *sql.DB {
	dbm, err := sql.Open("sqlite3", db)
	if err != nil {
		return nil
	}
	return dbm
}

func Save(sql string, db string) {
	dbm := GetDBM(db)
	dbm.Exec(sql)
}

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

func Init(db string) {
	dbm := GetDBM(db)
	s := "status int(5), delFlg int(5)"
	sql := "create table if not exists article (id integer primary key, title varchar(100), timestamp varchar(50), content text," + s + ");"
	dbm.Exec(sql)
	sql = "create table if not exists user (id integer primary key, first_name varchar(10), last_name varchar(10), email varchar(50), password varchar(200)," + s + ")"
	dbm.Exec(sql)
}

func Save(sql string, db string) {
	dbm := GetDBM(db)
	dbm.Exec(sql)
}

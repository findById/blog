package service

import (
	"blog/dao"
	"blog/entity"
	"errors"
	"log"
)

var (
	db = "./build/blog.db"
)

func Init() {
	dao.Init(db)
}

func FindUserByEmail(email string) (entity.User, error) {
	dbm := dao.GetDBM(db)
	defer dbm.Close()

	var item entity.User
	sql := "SELECT * FROM user WHERE email=?"
	rows, err := dbm.Query(sql, email)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	count := -1
	if rows.Next() {
		count++
		rows.Scan(&item.Id, &item.FirstName, &item.LastName, &item.Email, &item.Password)
	}
	if count == -1 {
		return item, errors.New("no data")
	}
	return item, nil
}

func SaveUser(firstName string, lastName string, email string, password string) {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	sql := "INSERT INTO user (first_name,last_name,email,password) VALUES (?,?,?,?)"
	dbm.Exec(sql, firstName, lastName, email, password)
}

func DeleteLogical(table, del, id string) {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	sql := "UPDATE " + table + " SET delFlg=? WHERE id=?"
	dbm.Exec(sql, del, id)
}

func Count(table string) int64 {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	sql := "SELECT count(*) as count FROM " + table + " WHERE delFlg=0"
	rows, err := dbm.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var count int64
	count = 0
	if err != nil {
		return count
	}
	if rows.Next() {
		rows.Scan(&count)
	}
	return count
}

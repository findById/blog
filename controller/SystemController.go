package controller

import (
	"net/http"
	"blog/dao"
	"blog/model/logger"
	"blog/conf"
	"encoding/hex"
	"crypto/md5"
)

func InstallAction(response http.ResponseWriter, request *http.Request) {

	dbm := dao.GetDBM(conf.DATABASE_NAME);

	sql := "create table if not exists article (" +
	"id integer primary key" +
	", title varchar(50)" +
	", content text" +
	", keywords varchar(100)" +
	", description varchar(200)" +
	", lang varchar(10)" +
	", tag varchar(10)" +
	", timestamp varchar(50)" +
	", status varchar(5)" +
	", delFlg varchar(5)" +
	");";
	dbm.Exec(sql);
	logger.Debug("DB", sql);

	sql = "DELETE user";
	dbm.Exec(sql);
	logger.Debug("DB", sql);

	sql = "create table if not exists user (" +
	"id integer primary key" +
	", first_name varchar(10)" +
	", last_name varchar(10)" +
	", email varchar(50)" +
	", password varchar(200)" +
	", status varchar(5)" +
	", delFlg varchar(5)" +
	");";
	dbm.Exec(sql);
	logger.Debug("DB", sql);

	email := "admin";
	password := "admin";
	hash := md5.New();
	hash.Write([]byte(password));
	password = hex.EncodeToString(hash.Sum(nil));

	sql = "INSERT INTO user (first_name,last_name,email,password,status,delFlg) VALUES (?,?,?,?,?,?)"
	dbm.Exec(sql, "admin", "admin", email, password, "0", "0");
	logger.Debug("DB", sql);

	dbm.Close();

	http.Redirect(response, request, "/", http.StatusFound)
}
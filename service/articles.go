package service

import (
	"blog/dao"
	"blog/entity"
	"errors"
	"log"
	"strconv"
	"time"
	"blog/conf"
)

func SaveArticle(article entity.Article) {
	dbm := dao.GetDBM(conf.DATABASE_NAME)
	defer dbm.Close()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sql := "INSERT INTO article (title,content,keywords,description,lang,tag,timestamp,status,delFlg) VALUES (?,?,?,?,?,?,?,?,?)"
	dbm.Exec(sql, article.Title,article.Content,article.Keywords,article.Description,article.Lang,article.Tag,timestamp,article.Status,0);
}

func UpdateArticle(id, title, content string) {
	dbm := dao.GetDBM(conf.DATABASE_NAME)
	defer dbm.Close()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sql := "UPDATE article SET title=?,timestamp=?,content=? WHERE id=?"
	dbm.Exec(sql, title, timestamp, content, id)
}

func DeleteArticle(id string) {
	dbm := dao.GetDBM(conf.DATABASE_NAME)
	defer dbm.Close()
	sql := "UPDATE article SET delFlg=? WHERE id=?"
	dbm.Exec(sql, 1, id)
}


func FindArticleById(id int) (entity.Article, error) {
	dbm := dao.GetDBM(conf.DATABASE_NAME)
	defer dbm.Close()
	sql := "SELECT * FROM article WHERE id=? AND delFlg=0"
	rows, err := dbm.Query(sql, strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var item entity.Article
	count := -1
	for rows.Next() {
		count++
		rows.Scan(&item.Id, &item.Title, &item.Content, &item.Keywords, &item.Description, &item.Lang, &item.Tag, &item.Timestamp, &item.Status, &item.DelFlg)
		tmp, err := strconv.ParseInt(item.Timestamp, 10, 0)
		if err != nil {
			tmp = time.Now().Unix()
		}
		// item.Timestamp = time.Unix(tmp, 0).Format("2006-01-02 15:04:05")
		item.Timestamp = time.Unix(tmp, 0).Format("2006-01-02")
	}
	if count == -1 {
		return item, errors.New("no data")
	}
	return item, nil
}

func FindArticles(offset, size int) ([]entity.Article, int64, error) {
	dbm := dao.GetDBM(conf.DATABASE_NAME)
	defer dbm.Close()
	sql := "SELECT * FROM article WHERE delFlg=0 AND status=0 LIMIT ?,?"
	rows, err := dbm.Query(sql, offset, size)
	if err != nil {
		return nil, 0, err;
	}
	defer rows.Close()
	var items []entity.Article = make([]entity.Article, 0)
	count := -1
	for rows.Next() {
		count++
		var item entity.Article
		rows.Scan(&item.Id, &item.Title, &item.Content, &item.Keywords, &item.Description, &item.Lang, &item.Tag, &item.Timestamp, &item.Status, &item.DelFlg)
		tmp, err := strconv.ParseInt(item.Timestamp, 10, 0)
		if err != nil {
			tmp = time.Now().Unix()
		}
		item.Timestamp = time.Unix(tmp, 0).Format("2006-01-02")
		items = append(items, item)
	}
	if count == -1 {
		return nil, 0, errors.New("no data");
	}
	return items, GetArticleCount(), nil
}

func GetArticleCount() int64 {
	dbm := dao.GetDBM(conf.DATABASE_NAME)
	defer dbm.Close()
	sql := "SELECT count(*) as count FROM article WHERE delFlg=0 AND status=0"
	rows, err := dbm.Query(sql)
	if err != nil {
		return 0;
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

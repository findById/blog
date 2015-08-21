package service

import (
	"blog/dao"
	"blog/entity"
	"blog/utils"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

func SaveArticle(title, content, status string) {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sql := "INSERT INTO article (title,timestamp,content,status,delFlg) VALUES (?,?,?,?,?)"
	dbm.Exec(sql, title, timestamp, content, status, 0)
}

func UpdateArticle(id, title, content string) {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sql := "UPDATE article SET title=?,timestamp=?,content=? WHERE id=?"
	dbm.Exec(sql, title, timestamp, content, id)
}

func FindArticleById(id int) (entity.Article, error) {
	dbm := dao.GetDBM(db)
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
		rows.Scan(&item.Id, &item.Title, &item.Timestamp, &item.Content, &item.Status, &item.DelFlg)
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

func FindArticles(start int, size int) ([]entity.Article, error) {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	sql := "SELECT * FROM article WHERE delFlg=0 LIMIT ?,?"
	rows, err := dbm.Query(sql, start, size)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []entity.Article = make([]entity.Article, 0)
	count := -1
	for rows.Next() {
		count++
		var item entity.Article
		rows.Scan(&item.Id, &item.Title, &item.Timestamp, &item.Content, &item.Status, &item.DelFlg)
		tmp, err := strconv.ParseInt(item.Timestamp, 10, 0)
		if err != nil {
			tmp = time.Now().Unix()
		}
		item.Timestamp = time.Unix(tmp, 0).Format("2006-01-02")
		item.Content = utils.Substring(item.Content, 0, 100) + "..."
		item.Content = strings.Replace(item.Content, "```java", "", -1)
		item.Content = strings.Replace(item.Content, "```go", "", -1)
		item.Content = strings.Replace(item.Content, "```golang", "", -1)
		item.Content = strings.Replace(item.Content, "```php", "", -1)
		item.Content = strings.Replace(item.Content, "```c", "", -1)
		item.Content = strings.Replace(item.Content, "```javascript", "", -1)
		item.Content = strings.Replace(item.Content, "```js", "", -1)
		item.Content = strings.Replace(item.Content, "```css", "", -1)
		item.Content = strings.Replace(item.Content, "```html", "", -1)
		item.Content = strings.Replace(item.Content, "```sql", "", -1)
		item.Content = strings.Replace(item.Content, "```shell", "", -1)
		item.Content = strings.Replace(item.Content, "```", "", -1)
		item.Content = strings.Replace(item.Content, "<br>", "", -1)
		items = append(items, item)
	}
	if count == -1 {
		return nil, errors.New("no data")
	}
	return items, nil
}

func CountArticle() int64 {
	dbm := dao.GetDBM(db)
	defer dbm.Close()
	sql := "SELECT count(*) as count FROM article WHERE delFlg=0"
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

package controller

import (
	"blog/model/logger"
	"blog/service"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"blog/utils"
	"blog/entity"
	"encoding/json"
	"fmt"
)

func ArticleSaveHandler(w http.ResponseWriter, r *http.Request, action string) {
	if (action != "post") {
		http.NotFound(w, r);
		return;
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	keywords := r.FormValue("keywords")
	description := r.FormValue("description")
	lang := r.FormValue("lang")
	tag := r.FormValue("tag")
	status := r.FormValue("status")

	if (utils.IsEmpty(title)) {
		OnResponse(w, 201, "标题不能为空", nil);
		return;
	}
	if (utils.IsEmpty(content)) {
		OnResponse(w, 202, "内容不能为空", nil);
		return;
	}

	var article entity.Article;
	article.Title = title;
	article.Content = content;
	article.Keywords = keywords;
	article.Description = description;
	article.Lang = lang;
	article.Status = status;
	article.Tag = tag;

	service.SaveArticle(article)

	tmp := time.Unix(time.Now().Unix(), 0).Format("20060102")
	tmp2 := strings.Replace(title, " ", "-", -1)
	filename := "posted/" + tmp + "-" + tmp2 + ".txt"
	ioutil.WriteFile(filename, []byte(content), 0600)

	http.Redirect(w, r, "/articles", http.StatusFound)
}

func ArticlePostHandler(w http.ResponseWriter, r *http.Request, action string) {
	if r.Method == "GET" {
		model := make(map[string]interface{})
		model["title"] = "Post"
		model["action"] = "post"
		ExecuteTemplate(w, "post", model)
		return;
	}

	if r.Method != "POST" {
		return;
	}

	defer r.Body.Close();
	body, err := ioutil.ReadAll(r.Body);
	if  err != nil {
		http.NotFound(w, r);
		return;
	}

	var dat map[string]interface{};
	err = json.Unmarshal(body, &dat);
	if  err != nil {
		http.NotFound(w, r);
		return;
	}

	title := fmt.Sprint(dat["title"]);
	content := fmt.Sprint(dat["content"]);
	keywords := fmt.Sprint(dat["keywords"]);
	description := fmt.Sprint(dat["description"]);
	lang := fmt.Sprint(dat["lang"]);
	tag := fmt.Sprint(dat["tag"]);
	status := fmt.Sprint(dat["status"]);

	var article entity.Article;
	article.Title = title;
	article.Content = content;
	article.Keywords = keywords;
	article.Description = description;
	article.Lang = lang;
	article.Status = status;
	article.Tag = tag;

	service.SaveArticle(article)

	tmp := time.Unix(time.Now().Unix(), 0).Format("20060102")
	tmp2 := strings.Replace(title, " ", "-", -1)
	filename := "posted/" + tmp + "-" + tmp2 + ".txt"
	ioutil.WriteFile(filename, []byte(content), 0600)

	http.Redirect(w, r, "/articles", http.StatusFound)
}

func ArticleEditHandler(w http.ResponseWriter, r *http.Request, action string) {
	model := make(map[string]interface{})
	model["title"] = "Post"
	model["action"] = action
	ExecuteTemplate(w, "edit", model)
}

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("console", "request", "[" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05") + "][" +
	r.RemoteAddr + "][" + r.UserAgent() + "][" + r.Host + r.RequestURI + "]");

	p, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil {
		p = 1
	}
	size := 10
	offset := (p - 1) * size

	items, count, err := service.FindArticles(offset, size)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	model := make(map[string]interface{});
	model["timestamp"] = time.Unix(time.Now().Unix(), 0).Format("20060102150405");
	model["token"] = utils.GetRandomString(20);
	model["keywords"] = "";
	model["description"] = "";

	model["title"] = "Insert title here";

	model["items"] = items

	max, _ := strconv.ParseInt(strconv.Itoa(offset + size), 10, 64)

	var next int
	if max < count {
		next = p + 1
	} else {
		next = p
	}
	model["next"] = next

	var pre int
	if p > 1 {
		pre = p - 1
	} else {
		pre = 1
	}
	model["pre"] = pre

	ExecuteTemplate(w, "articles", model)
}

func ArticleByIdHandler(w http.ResponseWriter, r *http.Request, action string) {
	id, err := strconv.Atoi(action)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	item, err := service.FindArticleById(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	model := make(map[string]interface{})
	model["keywords"] = "" + item.Keywords;
	model["description"] = "" + item.Description;
	model["lang"] = "" + item.Lang;
	model["type"] = "" + item.Tag;

	model["timestamp"] = time.Unix(time.Now().Unix(), 0).Format("20060102150405");
	model["token"] = utils.GetRandomString(20);

	model["item"] = item
	model["title"] = item.Title
	ExecuteTemplate(w, "view", model)
}

func ArticleDeleteHandler(w http.ResponseWriter, r *http.Request, action string) {
	if (r.Method != "POST") {
		http.NotFound(w, r);
		return;
	}

	w.Header().Set("content-type", "application/json");

	defer r.Body.Close();
	body, err := ioutil.ReadAll(r.Body);
	if  err != nil {
		http.NotFound(w, r);
		return;
	}

	var dat map[string]interface{};
	err = json.Unmarshal(body, &dat);
	if  err != nil {
		http.NotFound(w, r);
		return;
	}

	id := fmt.Sprint(dat["id"]);
	service.DeleteArticle(id);

	http.Redirect(w, r, "/articles", http.StatusFound)
}


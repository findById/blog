package controller

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("views/index.html",
	"views/login.html", "views/articles/edit.html",
	"views/articles/view.html", "views/articles/articles.html",
	"views/articles/post.html"));

func ExecuteTemplate(response http.ResponseWriter, tmp string, data interface{}) {
	response.Header().Set("Content-Type", "text/html")
	err := templates.ExecuteTemplate(response, tmp + ".html", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

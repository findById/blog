package controller

import (
	"html/template"
	"net/http"
	"encoding/base64"
	"blog/model/logger"
	"time"
	"encoding/json"
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

func CheckCookie(response http.ResponseWriter, request *http.Request, action string) bool  {
	cookie, err := request.Cookie("account_email");
	if (err != nil || cookie == nil) {
		return false;
	}
	b, e := base64.StdEncoding.DecodeString(cookie.Value);
	if e != nil {
		return false;
	}
	logger.Info(action, "[" + string(b) + "][" + time.Unix(time.Now().Unix(), 0).Format("20060102150405") + "][" + request.RemoteAddr +
	"][" + request.UserAgent() + "][" + request.Host + request.RequestURI + "]");
	return true;
}

func OnResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	result := make(map[string]interface{})
	result["result"] = data
	result["message"] = message
	result["statusCode"] = statusCode

	//	enc := json.NewEncoder(w)
	//	enc.Encode(result)

	b, err := json.Marshal(result)
	if err != nil {
		return
	}
	w.Write(b)
}

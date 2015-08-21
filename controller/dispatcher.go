package controller

import (
	"blog/model/logger"
	"html/template"
	"net/http"
	"regexp"
)

func InitServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	// 用户
	mux.HandleFunc("/account", AccountHandler)
	// 博文
	mux.HandleFunc("/articles", ArticlesHandler)
	mux.HandleFunc("/view/", makeHandler(ArticleByIdHandler))
	mux.HandleFunc("/save/", makeHandler(ArticleSaveHandler))
	mux.HandleFunc("/edit/", makeHandler(ArticleEditHandler))
	mux.HandleFunc("/post/", ArticlePostHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", mainHandler)

	return mux
}

func mainHandler(response http.ResponseWriter, request *http.Request) {
	logger.Info("request:", "["+request.RemoteAddr+"]["+request.UserAgent()+"]["+request.Host+request.RequestURI+"]")
	if request.URL.Path != "/" {
		http.NotFound(response, request)
		return
	}
	model := make(map[string]interface{})
	model["title"] = "Insert title here"
	ExecuteTemplate(response, "index", model)
}

var validPath = regexp.MustCompile("^/(post|edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		logger.Info("request:", "["+request.RemoteAddr+"]["+request.UserAgent()+"]["+request.Host+request.RequestURI+"]")
		m := validPath.FindStringSubmatch(request.URL.Path)
		if m == nil {
			http.NotFound(response, request)
			return
		}
		fn(response, request, m[2])
	}
}

var templates = template.Must(template.ParseFiles("views/index.html", "views/login.html", "views/articles/edit.html", "views/articles/view.html", "views/articles/articles.html", "views/articles/post.html"))

func ExecuteTemplate(response http.ResponseWriter, tmp string, data interface{}) {
	response.Header().Set("Content-Type", "text/html")
	err := templates.ExecuteTemplate(response, tmp+".html", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

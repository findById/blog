package main

import (
	"blog/conf"
	"blog/controller"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"html/template"
	"blog/model/logger"
	"regexp"
	"strings"
	"time"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

func init() {
	path := flag.String("conf", "conf/conf.json", "path of conf.json")
	ip := flag.String("ip", "", "this will overwrite conf.IP if specified")
	port := flag.String("port", "", "this will overwrite conf.Port if specified")
	host := flag.String("host", "", "this will overwrite conf.Server if specified")
	context := flag.String("context", "", "this will overwrite conf.Server if specified")
	flag.Parse()

	conf.Init(*path, *ip, *port, *host, *context)

}

func main() {
	flag.Parse();

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0");
		if err != nil {
			log.Fatal(err);
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644);
		if err != nil {
			log.Fatal(err);
		}
		s := &http.Server{};
		s.Serve(l);
		return
	}

	mux := http.NewServeMux();

	// 系统
	mux.HandleFunc("/install", controller.InstallAction);
	mux.HandleFunc("/upload", controller.UploadHandler);
	mux.HandleFunc("/account", controller.AccountHandler);
	mux.HandleFunc("/account/save", controller.UserSaveHandler);
	mux.HandleFunc("/account/login", controller.LoginAction);

	// 博文
	mux.HandleFunc("/articles", controller.ArticlesHandler);
	mux.HandleFunc("/view/", makeHandler(controller.ArticleByIdHandler));

	mux.HandleFunc("/article/save/", RouterHandler(controller.ArticleSaveHandler));
	mux.HandleFunc("/article/edit/", RouterHandler(controller.ArticleEditHandler));
	mux.HandleFunc("/article/post/", RouterHandler(controller.ArticlePostHandler));
	mux.HandleFunc("/article/delete/", RouterHandler(controller.ArticleDeleteHandler));

	// 静态文件
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))));

	// 首页
	mux.HandleFunc("/", MainHandler);

	http.ListenAndServe(conf.Conf.Host, mux);
}

func MainHandler(response http.ResponseWriter, request *http.Request) {
	logger.Log("console", "request", "[" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05") + "][" +
	request.RemoteAddr + "][" + request.UserAgent() + "][" + request.Host + request.RequestURI + "]");

	if request.URL.Path != "/" {
		http.NotFound(response, request);
		return
	}
	model := make(map[string]interface{});
	model["title"] = "Insert title here";
	t, err := template.ParseFiles("views/index.html");
	if (err != nil) {
		http.NotFound(response, request);
		return;
	}
	t.Execute(response, model);
}


var valid = regexp.MustCompile("^/(view|articles)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		logger.Log("console", "request", "[" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05") + "][" +
		request.RemoteAddr + "][" + request.UserAgent() + "][" + request.Host + request.RequestURI + "]");
		m := valid.FindStringSubmatch(request.URL.Path);
		if m == nil {
			http.NotFound(response, request);
			return;
		}
		fn(response, request, m[2]);
	}
}

var validPath = regexp.MustCompile("^/([a-zA-Z]+)/(post|edit|save|view)/([a-zA-Z0-9]+)$")

func RouterHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		path := request.URL.Path;

		if (strings.Contains(path, "article") || strings.Contains(path, "save") || strings.Contains(path, "upload")) {
			if (!controller.CheckCookie(response, request, "request")) {
				http.Redirect(response, request, "/account?redirect=" + path, http.StatusFound)
				return;
			}
		}

		m := validPath.FindStringSubmatch(path);
		if m == nil {
			http.NotFound(response, request);
			return;
		}
		fn(response, request, m[3]);
	}
}

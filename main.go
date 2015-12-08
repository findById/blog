package main

import (
	"blog/conf"
	"blog/controller"
	"blog/service"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"html/template"
	"blog/model/logger"
	"regexp"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

func init() {
	service.Init()

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
	mux.HandleFunc("/upload", controller.UploadHandler);
	mux.HandleFunc("/account", controller.AccountHandler);

	// 博文
	mux.HandleFunc("/articles", controller.ArticlesHandler);
	mux.HandleFunc("/view/", makeHandler(controller.ArticleByIdHandler));
	mux.HandleFunc("/save/", makeHandler(controller.ArticleSaveHandler));
	mux.HandleFunc("/edit/", makeHandler(controller.ArticleEditHandler));
	mux.HandleFunc("/post/", controller.ArticlePostHandler);

	// 静态文件
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))));

	// 首页
	mux.HandleFunc("/", mainHandler);

	http.ListenAndServe(conf.Conf.Host, mux);
}

func mainHandler(response http.ResponseWriter, request *http.Request) {
	logger.Info("request:", "[" + request.RemoteAddr + "][" + request.UserAgent() + "][" + request.Host + request.RequestURI + "]")
	if request.URL.Path != "/" {
		http.NotFound(response, request);
		return
	}
	model := make(map[string]interface{});
	model["title"] = "Insert title here";
	t, err := template.ParseFiles("views/index.html");
	if (err != nil) {
		http.NotFound(response, request);
	}
	t.Execute(response, model);
}


var validPath = regexp.MustCompile("^/(post|edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		logger.Info("request:", "[" + request.RemoteAddr + "][" + request.UserAgent() + "][" + request.Host + request.RequestURI + "]")
		m := validPath.FindStringSubmatch(request.URL.Path)
		if m == nil {
			http.NotFound(response, request)
			return
		}
		fn(response, request, m[2])
	}
}

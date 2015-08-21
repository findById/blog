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
	flag.Parse()

	mux := controller.InitServeMux()

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(conf.Conf.Host, mux)
}

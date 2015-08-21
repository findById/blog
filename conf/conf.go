package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	// PathSeparator holds the OS-specific path separator.
	PathSeparator = string(os.PathSeparator)
	// PathListSeparator holds the OS-specific path list separator.
	PathListSeparator = string(os.PathListSeparator)

	VERSION = "0.0.1"
)

var (
	Conf *conf
)

// Configuration.
type conf struct {
	IP      string // server ip, ${ip}
	Port    string // server port
	Host    string // server host and port ({IP}:{Port})
	Context string // server context
}

func Load() {

}

func Init(path, ip, port, host, context string) {
	initConfig(path, ip, port, host, context)
}

func initConfig(path, ip, port, host, context string) {
	bytes, err := ioutil.ReadFile(path)
	if nil != err {
		log.Println(err)
		os.Exit(-1)
	}

	Conf = &conf{}

	err = json.Unmarshal(bytes, Conf)
	if err != nil {
		log.Println("Parses [conf.json] error: ", err)
		os.Exit(-1)
	}

	// IP
	Conf.IP = strings.Replace(Conf.IP, "${ip}", ip, 1)

	if "" != port {
		Conf.Port = port
	}

	// Host
	Conf.Host = strings.Replace(Conf.Host, "{IP}", Conf.IP, 1)
	Conf.Host = strings.Replace(Conf.Host, "{Port}", Conf.Port, 1)
	if "" != host {
		Conf.Host = host
	}

	// Context
	if "" != context {
		Conf.Context = context
	}
}

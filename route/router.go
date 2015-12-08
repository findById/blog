package route

import (
	"net/http"
)

type Dispatcher struct {
	handlers map[string][]*Handler
	NotFound http.HandlerFunc
}

type Handler struct {
	patt  string
	parts []string
	wild  bool
	http.HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{make(map[string][]*Handler), nil}
}

func AddAction(method, pattern string, handler interface{})  {

}

func Run()  {
	
}
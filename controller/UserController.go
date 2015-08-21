package controller

import (
	"net/http"
)

func UserSaveHandler(response http.ResponseWriter, request *http.Request) {
	model := make(map[string]interface{})
	model["title"] = "Insert title here"
	ExecuteTemplate(response, "index", model)
}

func AccountHandler(response http.ResponseWriter, request *http.Request) {
	model := make(map[string]interface{})
	model["title"] = "Sign In"
	ExecuteTemplate(response, "login", model)
}

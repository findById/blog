package controller

import (
	"blog/service"
	"encoding/json"
	"log"
	"net/http"
)

func LoginAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := r.ParseForm()
	if err != nil {
		response(w, 0, "参数错误", nil)
		return
	}
	email := r.FormValue("email")
	passwd := r.FormValue("passwd")
	log.Println(email + "   " + passwd)
	if email == "" || passwd == "" {
		response(w, 0, "参数错误", nil)
		return
	}
	user, err := service.FindUserByEmail(email)
	if err != nil {
		response(w, 0, "用户不存在", nil)
		return
	}
	if user.Password != passwd {
		response(w, 0, "密码错误", nil)
		return
	}
	// 存入cookie,使用cookie存储
	cookie := http.Cookie{Name: "email", Value: email, Path: "/"}
	http.SetCookie(w, &cookie)
}

func response(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	result := make(map[string]interface{})
	result["result"] = data
	result["message"] = message
	result["statusCode"] = statusCode
	b, err := json.Marshal(result)
	if err != nil {
		return
	}
	w.Write(b)
}

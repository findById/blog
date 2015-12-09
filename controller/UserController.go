package controller

import (
	"net/http"
	"blog/service"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func UserSaveHandler(response http.ResponseWriter, request *http.Request) {

	if !CheckCookie(response, request, "UserSave") {
		http.NotFound(response, request);
		return;
	}

	if (request.Method == "GET") {
		model := make(map[string]interface{})
		model["title"] = "Sign In"
		ExecuteTemplate(response, "login", model)
		return;
	}

	if (request.Method != "POST") {
		http.NotFound(response, request);
		return;
	}

	defer request.Body.Close();
	body, err := ioutil.ReadAll(request.Body);
	if  err != nil {
		OnResponse(response, 401, "请求参数不能为空", nil);
		return;
	}

	var dat map[string]interface{};
	err = json.Unmarshal(body, &dat);
	if  err != nil {
		OnResponse(response, 402, "JSON解析失败", nil);
		return;
	}

	firstName := fmt.Sprint(dat["firstName"]);
	lastName := fmt.Sprint(dat["lastName"]);
	email := fmt.Sprint(dat["email"]);
	password := fmt.Sprint(dat["password"]);

	_, e := service.FindUserByEmail(email);
	if e == nil {
		// http.NotFound(response, request);
		OnResponse(response, 201, "用户已存在", nil);
		return;
	}

	hash := md5.New();
	hash.Write([]byte(password));
	password = hex.EncodeToString(hash.Sum(nil));

	service.SaveUser(firstName, lastName, email, password);

	OnResponse(response, 200, "ok", nil);
}

func AccountHandler(response http.ResponseWriter, request *http.Request) {
	model := make(map[string]interface{})
	model["title"] = "Sign In"
	ExecuteTemplate(response, "login", model)
}

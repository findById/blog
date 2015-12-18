package controller

import (
	"blog/service"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/hex"
	"crypto/md5"
	"time"
	"blog/model/logger"
	"encoding/base64"
	"blog/conf"
)

func LoginAction(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "POST") {
		http.NotFound(w, r);
		return;
	}
	w.Header().Set("content-type", "application/json");

	defer r.Body.Close();
	body, err := ioutil.ReadAll(r.Body);
	if err != nil {
		OnResponse(w, 201, "用户名或密码错误", nil);
		return;
	}

	var dat map[string]interface{};
	err = json.Unmarshal(body, &dat);
	if err != nil {
		OnResponse(w, 201, "用户名或密码错误", nil);
		return;
	}

	email := fmt.Sprint(dat["email"]);
	password := fmt.Sprint(dat["passwd"]);

	hash := md5.New();
	hash.Write([]byte(password));
	password = hex.EncodeToString(hash.Sum(nil));
	// log.Println(email + "   " + password);

	if email == "" || password == "" {
		OnResponse(w, 201, "用户名或密码错误", nil);
		return;
	}
	user, err := service.FindUserByEmail(email);
	if err != nil {
		OnResponse(w, 201, "用户名或密码错误", nil);
		return;
	}
	if user.Password != password {
		OnResponse(w, 201, "用户名或密码错误", nil);
		return;
	}

	// 存入cookie,使用cookie存储
	t := time.Now();
	expires := time.Date(t.Year(), t.Month(), t.Day(), t.Hour() + 5, t.Minute(), t.Second(), 0, time.Local);
	cookie := http.Cookie{Name: conf.SESSION_KEY, Value: base64.StdEncoding.EncodeToString([]byte(email)), Path: "/", Expires: expires};
	http.SetCookie(w, &cookie);

	logger.Log("login", "SignIn", "[" + email + "][" + time.Unix(time.Now().Unix(), 0).Format("20060102150405") + "][" +
	r.RemoteAddr + "][" + r.UserAgent() + "][" + r.Host + r.RequestURI + "]");

	OnResponse(w, 200, "ok", nil);

}

func LogoutAction(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(conf.SESSION_KEY);
	if (err == nil && cookie != nil) {
		cookie := http.Cookie{Name: conf.SESSION_KEY, Value: "", Path: "/", Expires: time.Now().AddDate(0,0,0)};
		http.SetCookie(w, &cookie);
	}
	http.Redirect(w, r, "/", http.StatusFound);
}

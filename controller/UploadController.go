package controller

import (
	"os"
	"path/filepath"
	"net/http"
	"html/template"
	"strings"
	"fmt"
	"io"
	"log"
	"encoding/json"
	"math/rand"
	"time"
	"io/ioutil"
)

const (
	TemplateDir = "views/"
	uploadDir = "upload"
	RootDir = "/home/work/dev/golang/"
)

type FileBean struct {
	File     os.FileInfo;
	FileName string;
	Path     string;
	IsDir    bool;
}

type FileItem struct {
	Path       string;
	OriginName string;
	Url        string;
}

/*
 * upload
 */
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		http.NotFound(w, r);
		return;
	}
	if r.Method == "POST" {
		fmt.Println("Client:", r.RemoteAddr, "Method:", r.Method, "URL:", r.Host + r.RequestURI)

		token := r.URL.Query().Get("token")

		if ("90f8d5RyvmxPeHJ4RL6h" != token) {
			http.NotFound(w, r);
			return;
		}

		r.ParseForm();
		r.ParseMultipartForm(32 << 20);

		multipart := r.MultipartForm
		if multipart == nil {
			log.Println("Not MultipartForm.")
			w.Write(([]byte)("Not MultipartForm."))
			return
		}

		fileHeaders, findFile := multipart.File["file"]
		if !findFile || len(fileHeaders) == 0 {
			log.Println("No file befor uploading.")
			w.Write(([]byte)("No file befor uploading."))
			return
		}

		var items []FileItem = make([]FileItem, 0)
		for _, v := range fileHeaders {
			fileName := v.Filename

			fileExt := filepath.Ext(fileName)
			if checkContentType(fileExt) == false {
				log.Println("Content type refused.")
				w.Write(([]byte)("Content type refused."))
				return
			}

			src, err := v.Open()
			checkError(err, "Open file error." + fileName)
			defer src.Close()

			pathSeparator := string(os.PathSeparator)

			realPath := uploadDir + pathSeparator + fileName;

			output := RootDir + pathSeparator + realPath;
			dst, err := os.OpenFile(output, os.O_WRONLY | os.O_CREATE, 0666)
			checkError(err, "Open local file error")
			defer dst.Close()
			io.Copy(dst, src)

			var item FileItem;
			item.Path = strings.Replace(realPath, pathSeparator, "/", -1);
			item.OriginName = fileName;
			items = append(items, item);
		}

		model := make(map[string]interface{})
		model["result"] = items
		model["size"] = len(fileHeaders)
		model["stautsCode"] = "200"
		enc := json.NewEncoder(w)
		enc.Encode(model)
	}
}

func getString(length int) (string) {
	temp := "abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890";
	result := "";
	r := rand.New(rand.NewSource(time.Now().UnixNano()));
	for i := 0; i < length; i++ {
		result += string(temp[r.Intn(len(temp)) % len(temp)])
	}
	log.Println(result);
	return result;
}

func checkError(err error, message string) {
	if (err != nil) {
		println(err.Error() + message)
	}
}

func checkContentType(name string) bool {
	ext := []string{".apk", ".zip", ".png"}
	for _, v := range ext {
		if v == name {
			return true
		}
	}
	return false
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client:", r.RemoteAddr, "Method:", r.Method, "URL:", r.Host + r.RequestURI)

	separator := string(os.PathSeparator);

	dir := r.URL.Query().Get("type")

	var items []FileBean = make([]FileBean, 0)

	files, _ := ioutil.ReadDir(RootDir + separator + uploadDir + separator + dir);

	for _, f := range files {
		var bean FileBean;
		bean.File = f;
		bean.FileName = f.Name();
		bean.IsDir = f.IsDir();
		items = append(items, bean)
	}

	model := make(map[string]interface{})
	model["items"] = items
	model["title"] = "List"
	model["type"] = dir + "";
	t, _ := template.ParseFiles("views/list.html");
	t.Execute(w, model)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client:", r.RemoteAddr, "Method:", r.Method, "URL:", r.Host + r.RequestURI)

	separator := string(os.PathSeparator);

	name := r.URL.Query().Get("name")

	err := os.Remove(RootDir + separator + uploadDir + separator + name);
	if (err != nil) {
		err.Error();
	}

	http.Redirect(w, r, "/", 301)
}

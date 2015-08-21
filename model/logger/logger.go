package logger

import (
	"log"
	"os"
	"time"
)

func Debug(tag, message string) {
	Write("debug", tag, message)
}

func Info(tag, message string) {
	Write("info", tag, message)
}

func Warn(tag, message string) {
	Write("warn", tag, message)
}

func Error(tag, message string) {
	Write("error", tag, message)
}

func Write(level, tag, content string) {
	if level == "debug" {
		log.Println(content)
		return
	}
	tmp := time.Unix(time.Now().Unix(), 0).Format("20060102")

	WriteFile("logs/"+tmp+"_"+level+".out", ""+tag+" "+content)
}

func WriteFile(file, content string) bool {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Write([]byte(content + "\n"))
	if err != nil {
		return false
	}
	return true
}

package logger

import (
	"log"
	"os"
	"time"
	"fmt"
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

func Log(level, tag, message string) {
	Write(level, tag, message)
}

func Write(level, tag, content string) {
	if level == "debug" {
		log.Println(content)
		return
	}
	tmp := time.Unix(time.Now().Unix(), 0).Format("20060102")

	WriteFile("./build/logs/" + tmp + "_" + level + "%s_%d.out", "" + tag + " " + content)
}

func WriteFile(filePath, content string) bool {
	var out *os.File;
	for i := 0; i < 1024; i++ {
		f, err := os.OpenFile(fmt.Sprintf(filePath, "", i), os.O_CREATE | os.O_APPEND | os.O_RDWR, 0644);
		if err != nil {
			continue;
		}
		fi, err := f.Stat();
		if err == nil {
			if fi.Size() < (2 * (1 << 20)) {
				out = f;
				break;
			}
		}
		f.Close();
	}

	if out == nil {
		return false;
	}
	defer out.Close();
	out.Write([]byte(content + "\n"));
	return true;
}

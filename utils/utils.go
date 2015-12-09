package utils
import (
	"math/rand"
	"time"
)

func Substring(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func GetRandomString(length int) (string) {
	temp := "abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890";
	result := "";
	r := rand.New(rand.NewSource(time.Now().UnixNano()));
	for i := 0; i < length; i++ {
		result += string(temp[r.Intn(len(temp)) % len(temp)]);
	}
	return result;
}

func IsEmpty(obj string) bool {
	if (obj == "" || len(string(obj)) <= 0) {
		return true;
	}
	return false;
}


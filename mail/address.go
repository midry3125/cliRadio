package mail

import (
	"strings"
)

func GetServerAddr(addr string) (string, string) {
	idx := strings.Index(addr, "@")
	switch addr[idx+1:] {
	case "gmail.com":
		return "smtp.gmail.com", "465"
	case "yahoo.co.jp":
		return "smtp.mail.yahoo.com", "465"
	}
    return "", ""
}
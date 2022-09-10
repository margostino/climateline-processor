package security

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func IsAuthorized(r *http.Request) bool {
	jobSecret := os.Getenv("CLIMATELINE_JOB_SECRET")
	requestSecret := r.Header.Get("Authorization")
	return requestSecret == fmt.Sprintf("Bearer %s", jobSecret)
}

func IsAdmin(r *http.Request) bool {
	botSecret := os.Getenv("TELEGRAM_BOT_SECRET")
	requestSecret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	return requestSecret == botSecret
}

func IsAdminData(username string, chatId int64, r *http.Request) bool {
	return r.Host == os.Getenv("ALLOWED_HOST") &&
		r.Header.Get("User-Agent") == os.Getenv("ALLOWED_USER_AGENT") &&
		r.Method == os.Getenv("ALLOWED_METHOD") &&
		r.Proto == os.Getenv("ALLOWED_PROTOCOL") &&
		username == os.Getenv("ALLOWED_USERNAME") &&
		strconv.FormatInt(chatId, 10) == os.Getenv("ALLOWED_CHAT_ID")
}

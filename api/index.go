package api

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"log"
	"net/http"
	"os"
	"strconv"
)

var bot = NewBot()

func RunJob(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Hello World from Go"
	resp["language"] = "go"
	resp["cloud"] = "Hosted on Vercel!"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.Write(jsonResp)
	}
	Send()
	return
}

func NewBot() *tgbotapi.BotAPI {
	client, error := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	common.Check(error)
	//bot.Debug = true
	log.Printf("Authorized on account %s\n", client.Self.UserName)
	return client
}

func Send() {
	userId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ADMIN_USER"), 10, 64)
	msg := tgbotapi.NewMessage(userId, "Do you want to upload new article? [ Yes | No ]")
	msg.ReplyMarkup = nil
	bot.Send(msg)
}

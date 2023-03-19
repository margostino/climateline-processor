package api

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/bot"
	"github.com/margostino/climateline-processor/security"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Msg    string `json:"text"`
	ChatID int64  `json:"chat_id"`
	Method string `json:"method"`
}

var botApi *tgbotapi.BotAPI

func Bot(w http.ResponseWriter, r *http.Request) {
	var reply string

	log.Printf("Method: %s "+
		"Proto: %s "+
		"User-Agent: %s, "+
		"Host: %s, "+
		"RequestURI: %s, "+
		"RemoteAddr: %s",
		r.Method, r.Proto, r.Header.Get("User-Agent"), r.Host, r.RequestURI, r.RemoteAddr)

	body, _ := ioutil.ReadAll(r.Body)
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Fatal("Error updating →", err)
	}

	log.Printf("[%s@%d] %s - reply: %t", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text, update.Message.ReplyToMessage != nil)

	if security.IsAdmin(r) {
		w.Header().Add("Content-Type", "application/json")
		if bot.IsValidInput(update.Message) {
			reply = bot.Reply(update.Message)
		} else if "/start" == update.Message.Text {
			reply = "🌎 Welcome!"
		} else {
			reply = "Input is not valid"
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		reply = "Unauthorized to handle Bot communication"
	}

	data := Response{
		Msg:    reply,
		Method: "sendMessage",
		ChatID: update.Message.Chat.ID,
	}

	message, _ := json.Marshal(data)
	log.Printf("Response %s", string(message))
	fmt.Fprintf(w, string(message))
}

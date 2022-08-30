package api

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"os"
	"strconv"
)

var bot = NewBot()
var Items map[int]*gofeed.Item

func RunJob(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	log.Printf("ITEMS: %d\n", len(Items))
	Items = make(map[int]*gofeed.Item)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(os.Getenv("FEED_URL"))
	for id, item := range feed.Items {
		Notify(item, id)
		Items[id] = item
	}

	resp := make(map[string]string)
	resp["message"] = "Hello World from Go"
	resp["language"] = "go"
	resp["cloud"] = "Hosted on Vercel!"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		w.Write(jsonResp)
	}

	Ask()

	return
}

func NewBot() *tgbotapi.BotAPI {
	client, error := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	common.Check(error, "when creating a new BotAPI instance")
	//bot.Debug = true
	log.Printf("Authorized on account %s\n", client.Self.UserName)
	return client
}

func Notify(item *gofeed.Item, id int) {
	message := fmt.Sprintf("🔔 New article! \n🔑 ID: %d\n🗓 Date: %s\n📖 Content: %s\n", id, item.Updated, item.Content)
	Send(message)
}

func Ask() {
	message := "Do you want to upload new article? [ Yes | No ]"
	Send(message)
}

func Send(message string) {
	userId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ADMIN_USER"), 10, 64)
	msg := tgbotapi.NewMessage(userId, message)
	msg.ReplyMarkup = nil
	bot.Send(msg)
}

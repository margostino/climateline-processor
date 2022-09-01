package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"os"
	"strconv"
)

var bot = NewBot()

func Job(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	//log.Printf("Cached Items (RunJob): %d", len(cache.Items))

	var items = make([]*domain.Item, 0)

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(os.Getenv("FEED_URL"))
	for id, entry := range feed.Items {
		item := &domain.Item{
			Id:        id,
			Timestamp: entry.Updated,
			Title:     entry.Title,
			Link:      entry.Link,
			Content:   entry.Content,
		}
		Notify(item, id)
		items = append(items, item)
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

	UpdateCache(items)
	AskForUpdates()

	return
}

func NewBot() *tgbotapi.BotAPI {
	client, error := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	common.Check(error, "when creating a new BotAPI instance")
	//bot.Debug = true
	log.Printf("Authorized on account %s\n", client.Self.UserName)
	return client
}

func Notify(item *domain.Item, id int) {
	message := fmt.Sprintf("🔔 New article! \n🔑 ID: %d\n🗓 Date: %s\n📖 Content: %s\n", id, item.Timestamp, item.Content)
	Send(message)
}

func AskForUpdates() {
	message := "Do you want to upload new article? [ No | Yes {ID} ]\nExample: Yes 1"
	Send(message)
}

func Send(message string) {
	userId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ADMIN_USER"), 10, 64)
	msg := tgbotapi.NewMessage(userId, message)
	msg.ReplyMarkup = nil
	bot.Send(msg)
}

func UpdateCache(items []*domain.Item) {
	jsonData, err := json.Marshal(items)

	if !common.IsError(err, "when updating cache") {
		resp, err := http.Post(baseCacheUrl, "application/json", bytes.NewBuffer(jsonData))

		if resp.Status != "200" {
			log.Printf("Updating cache was not successful. Status: %d\n", resp.StatusCode)
		}
		common.SilentCheck(err, "in response of caching")
	}

}

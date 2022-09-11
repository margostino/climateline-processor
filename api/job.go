package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/security"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"os"
	"strconv"
)

var bot *tgbotapi.BotAPI

func Job(w http.ResponseWriter, r *http.Request) {
	bot, _ = NewBot()
	if security.IsAuthorized(r) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		var items = make([]*domain.Item, 0)

		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(os.Getenv("FEED_URL"))
		for id, entry := range feed.Items {
			item := &domain.Item{
				Id:        strconv.Itoa(id + 1),
				Timestamp: entry.Updated,
				Title:     entry.Title,
				Link:      entry.Link,
				Content:   entry.Content,
			}
			Notify(item)
			items = append(items, item)
		}

		response := domain.JobResponse{
			Items: len(items),
		}
		jsonResp, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
		} else {
			w.Write(jsonResp)
		}

		UpdateCache(items)
		AskForUpdates()

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}

func NewBot() (*tgbotapi.BotAPI, error) {
	client, error := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	//bot.Debug = true
	common.SilentCheck(error, "when creating a new BotAPI instance")
	//log.Printf("Authorized on account %s\n", client.Self.UserName)
	return client, error
}

func Notify(item *domain.Item) {
	message := fmt.Sprintf("🔔 New article! \n"+
		"🔑 ID: %s\n"+
		"🗓 Date: %s\n"+
		"💡 Title: %s\n"+
		"🔗 Link: <a href='%s'>Here</a>\n"+
		"📖 Content: %s\n",
		item.Id, item.Timestamp, item.Title, item.Link, item.Content)
	Send(message)
}

func AskForUpdates() {
	message := "❓ What do you want to do?\n" +
		"➡️ edit {id}\n" +
		"{new title}\n" +
		"{source name}\n" +
		"{location}\n" +
		"{category[agreements | assessment | awareness | warming | wildfires | floods | drought | health]}\n" +
		"⚡️️ Example:\n" +
		"edit 1\n" +
		"Massive heatwaves in Europe\n" +
		"Washington Post\n" +
		"Europe\n" +
		"warming\n" +
		"➡️ push {ids}\n" +
		"⚡️️ Example:\n" +
		"push 1 2"
	Send(message)
}

func Send(message string) {
	if bot != nil {
		userId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ADMIN_USER"), 10, 64)
		msg := tgbotapi.NewMessage(userId, message)
		msg.ReplyMarkup = nil
		msg.ParseMode = "HTML"
		bot.Send(msg)
	} else {
		log.Printf("Bot initialization failed")
	}
}

func UpdateCache(items []*domain.Item) {
	client := &http.Client{}
	json, err := json.Marshal(items)

	if !common.IsError(err, "when updating cache") {
		request, err := http.NewRequest(http.MethodPost, GetBaseCacheUrl(), bytes.NewBuffer(json))
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)

		common.SilentCheck(err, "in response of caching")

		if err == nil && response.StatusCode != 201 {
			log.Printf("Updating cache was not successful. Status: %d\n", response.StatusCode)
		}
	}

}

// GetBaseJobUrl Rather than a global and one-time assigment, this method is convenient for overriding on testing
func GetBaseJobUrl() string {
	return os.Getenv("JOB_BASE_URL")
}

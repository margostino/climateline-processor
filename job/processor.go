package job

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"github.com/mmcdole/gofeed"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

var botApi *tgbotapi.BotAPI

func Execute(writer *http.ResponseWriter) {
	var items = make([]*domain.Item, 0)
	botApi, _ = newBot()

	ctx := context.Background()
	sheets, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("GSHEET_API_KEY")))

	if !common.IsError(err, "when creating new Google API Service") {
		(*writer).WriteHeader(http.StatusOK)
		feeds := make([]string, 0)
		spreadsheetId := os.Getenv("SPREADSHEET_ID")
		readRange := os.Getenv("SPREADSHEET_RANGE")
		resp, err := sheets.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
		if !common.IsError(err, "unable to retrieve data from sheet") {
			if len(resp.Values) == 0 {
				(*writer).WriteHeader(http.StatusNotFound)
			} else {
				for _, row := range resp.Values {
					feeds = append(feeds, row[0].(string))
				}
			}

			for _, feedUrl := range feeds {
				fp := gofeed.NewParser()
				feed, _ := fp.ParseURL(feedUrl)

				if feed != nil {
					for id, entry := range feed.Items {
						item := &domain.Item{
							Id:        strconv.Itoa(id + 1),
							Timestamp: entry.Updated,
							Title:     entry.Title,
							Link:      entry.Link,
							Content:   entry.Content,
						}
						items = append(items, item)
					}
				} else {
					log.Printf("There are no feeds")
				}

			}

			for _, item := range items {
				notify(item)
			}

			response := domain.JobResponse{
				Items: len(items),
			}

			jsonResp, err := json.Marshal(response)
			if err != nil {
				fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
			} else {
				(*writer).Write(jsonResp)
			}

			updateCache(items)
			askForUpdates()

		} else {
			(*writer).WriteHeader(http.StatusBadRequest)
		}

	} else {
		(*writer).WriteHeader(http.StatusUnauthorized)
	}

}

func notify(item *domain.Item) {
	message := fmt.Sprintf("🔔 New article! \n"+
		"🔑 ID: %s\n"+
		"🗓 Date: %s\n"+
		"💡 Title: %s\n"+
		"🔗 Link: <a href='%s'>Here</a>\n"+
		"📖 Content: %s\n",
		item.Id, item.Timestamp, item.Title, item.Link, item.Content)
	send(message)
}

func askForUpdates() {
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
	send(message)
}

func send(message string) {
	if botApi != nil {
		userId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ADMIN_USER"), 10, 64)
		msg := tgbotapi.NewMessage(userId, message)
		msg.ReplyMarkup = nil
		msg.ParseMode = "HTML"
		botApi.Send(msg)
	} else {
		log.Printf("Bot initialization failed")
	}
}

func updateCache(items []*domain.Item) {
	client := &http.Client{}
	json, err := json.Marshal(items)

	if !common.IsError(err, "when updating cache") {
		request, err := http.NewRequest(http.MethodPost, cache.GetBaseCacheUrl(), bytes.NewBuffer(json))
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)

		common.SilentCheck(err, "in response of caching")

		if err == nil && response.StatusCode != 201 {
			log.Printf("Updating cache was not successful. Status: %d\n", response.StatusCode)
		}
	}

}

func newBot() (*tgbotapi.BotAPI, error) {
	client, error := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	//bot.Debug = true
	common.SilentCheck(error, "when creating a new BotAPI instance")
	//log.Printf("Authorized on account %s\n", client.Self.UserName)
	return client, error
}

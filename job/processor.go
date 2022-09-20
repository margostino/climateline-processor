package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var botApi *tgbotapi.BotAPI
var urls []string

func Execute(request *http.Request, writer *http.ResponseWriter) {
	var id = 0
	var items = make([]*domain.Item, 0)
	botApi, _ = newBot()

	category := strings.ToLower(request.URL.Query().Get("category"))
	if category == "" {
		category = "*"
	}
	urls = config.GetUrls(category)

	for _, feedUrl := range urls {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(feedUrl)

		if feed != nil {
			for _, entry := range feed.Items {
				var link, source string
				id += 1
				rawLink, err := url.Parse(entry.Link)

				if common.IsError(err, "when parsing feed link") {
					link = entry.Link
				} else {
					link = rawLink.Query().Get("url")
					sourceUrl, err := url.Parse(link)
					if !common.IsError(err, "when parsing source link") {
						source = common.NewString(sourceUrl.Hostname()).
							ReplaceAll("www.", "").
							ReplaceAll(".com", "").
							Value()
					}
				}

				item := &domain.Item{
					Id:         strconv.Itoa(id),
					Timestamp:  entry.Updated,
					Title:      entry.Title,
					Link:       link,
					Content:    entry.Content,
					SourceName: source,
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
		(*writer).WriteHeader(http.StatusNotFound)
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		(*writer).WriteHeader(http.StatusOK)
		(*writer).Write(jsonResp)
	}

	updateCache(items)
	askForUpdates()

}

func notify(item *domain.Item) {
	message := fmt.Sprintf("🔔 New article! \n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n", //"%s <a href='%s'>Here</a>\n",
		domain.ID_PREFIX, item.Id,
		domain.DATE_PREFIX, item.Timestamp,
		domain.TITLE_PREFIX, item.Title,
		domain.SOURCE_PREFIX, item.SourceName,
		domain.CONTENT_PREFIX, item.Content,
		domain.LINK_PREFIX, item.Link)
	send(message)
}

func askForUpdates() {
	message := "❓ What do you want to do?\n" +
		"➡️ edit {id}\n" +
		"{new title}\n" +
		"{source name}\n" +
		"{location}\n" +
		"{category[agreements | assessment | awareness | warming | wildfires | floods | drought | health | hurricane]}\n" +
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

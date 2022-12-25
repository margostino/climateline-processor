package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/urlshortener"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var botApi *tgbotapi.BotAPI
var twitterApi *twitter.Client
var urls []*config.UrlConfig
var bitlyDomain = "bit.ly"

func Execute(request *http.Request, writer *http.ResponseWriter) {
	var id = 0
	var items = make([]*domain.Item, 0)
	botApi, _ = newBot()
	twitterApi = newTwitterApi()

	category := strings.ToLower(request.URL.Query().Get("category"))
	if category == "" {
		category = "*"
	}
	urls = config.GetUrls(category)

	for _, feedUrl := range urls {
		if feedUrl.BotEnabled || feedUrl.TwitterEnabled {
			fp := gofeed.NewParser()
			feed, _ := fp.ParseURL(feedUrl.Url)

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
							source = strings.ReplaceAll(sourceUrl.Hostname(), "www.", "")
						}
					}

					item := &domain.Item{
						Id:                  strconv.Itoa(id),
						Timestamp:           entry.Updated,
						Title:               entry.Title,
						Link:                link,
						Content:             entry.Content,
						SourceName:          source,
						Tags:                feedUrl.Tags,
						ShouldNotifyBot:     feedUrl.BotEnabled,
						ShouldNotifyTwitter: feedUrl.TwitterEnabled,
					}
					items = append(items, item)
				}
			} else {
				log.Printf("There are no feeds")
			}
		}
	}

	for _, item := range items {
		if item.ShouldNotifyBot {
			notifyBot(item)
			updateCache(items)
		}
		if item.ShouldNotifyTwitter {
			notifyTwitter(item)
		}
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
}

func notifyTwitter(item *domain.Item) {
	var tweet string
	shorterLink := urlshortener.Shorten(item.Link)
	title := sanitizeTweet(item.Title)

	if shorterLink != "" {
		tweet = fmt.Sprintf("%s\nSource: %s (%s)\n%s", title, item.SourceName, shorterLink, item.Tags)
	} else {
		tweet = fmt.Sprintf("%s\nSource: %s\n%s", title, item.SourceName, item.Tags)
	}

	_, resp, err := twitterApi.Statuses.Update(tweet, nil)
	if err != nil && resp.StatusCode == 200 {
		log.Println("Tweet created")
	}
}

func sanitizeTweet(value string) string {
	return common.NewString(value).
		UnescapeString().
		ReplaceAll("<b>", "").
		ReplaceAll("</b>", "").
		Value()
}

func notifyBot(item *domain.Item) {
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

func newTwitterApi() *twitter.Client {
	var consumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	var consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	var token = os.Getenv("TWITTER_TOKEN")
	var tokenSecret = os.Getenv("TWITTER_TOKEN_SECRET")
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := config.Client(oauth1.NoContext, oauthToken)
	client := twitter.NewClient(httpClient)
	return client
}

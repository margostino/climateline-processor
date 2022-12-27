package job

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/internal"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var twitterApi *twitter.Client
var bitlyDomain = "bit.ly"

func Publish(request *http.Request, writer *http.ResponseWriter) {
	var items = make([]*domain.Item, 0)

	twitterApi = newTwitterApi()

	category := strings.ToLower(request.URL.Query().Get("category"))
	publishForced, err := strconv.ParseBool(request.URL.Query().Get("publish_forced"))
	if err != nil {
		publishForced = false
	}
	if category == "" {
		category = "*"
	}
	urls = config.GetUrls(category)

	items, err = internal.FetchNews(category)

	log.Printf("Category: %s\n", category)
	log.Printf("Publish Forced: %t\n", publishForced)

	for _, item := range items {
		if publishForced || item.ShouldNotifyTwitter {
			//notifyTwitter(item)
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
	//shorterLink := urlshortener.Shorten(item.Link)
	shorterLink := item.Link
	//params := &twitter.StatusUpdateParams{
	//	AttachmentURL: item.Link,
	//}
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

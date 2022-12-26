package news

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var urls []*config.UrlConfig

func GetBaseNewsUrl() string {
	return os.Getenv("NEWS_BASE_URL")
}

func Fetch(request *http.Request, writer *http.ResponseWriter) {
	var id = 0
	var items = make([]*domain.Item, 0)

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

	updateCache(items)

	jsonResp, err := json.Marshal(items)
	if err != nil {
		(*writer).WriteHeader(http.StatusNotFound)
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		(*writer).WriteHeader(http.StatusOK)
		(*writer).Write(jsonResp)
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

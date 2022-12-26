package bot

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/news"
	"log"
	"net/http"
	"os"
	"strings"
)

func Fetch(input string) string {
	var reply string
	var category string

	params := strings.Split(input, " ")
	if len(params) > 1 {
		category = params[1]
	} else {
		category = "*"
	}

	items, err := fetchItems(category)

	if !common.IsError(err, "when fetching news") {
		for _, item := range items {
			reply += fmt.Sprintf("🔔 New article! \n"+
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
		}
		//reply = "✅ Completed successfully"
	} else {
		reply = "🔴 Fetcher failed"
	}

	return reply
}

func fetchItems(category string) ([]*domain.Item, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s?category=%s", news.GetBaseNewsUrl(), category)
	log.Println("SARLANGA: " + url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	if err == nil && response.StatusCode == 200 {
		var items []*domain.Item
		err := json.NewDecoder(response.Body).Decode(&items)
		return items, err
	} else {
		return nil, err
	}
}

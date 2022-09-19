package bot

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"os"
)

func Show(input string) string {
	var id string
	var reply string

	if input == "/show" || input == "show" || input == "show all" || input == "show *" {
		id = "*"
	} else {
		id = extractIds(input, "show ")
	}

	items := getCachedItems(id)

	if len(items) > 0 {
		reply = buildShowReply(items[0])
	} else if id == "*" || id == "" || id == " " || id == "show" {
		reply = "Cache is empty"
	} else {
		reply = fmt.Sprintf("🤷‍ There is not item for ID %s", id)
	}

	return reply
}

func getCachedItems(ids string) []domain.Item {
	client := &http.Client{}
	var items []domain.Item
	url := fmt.Sprintf("%s?ids=%s", cache.GetBaseCacheUrl(), ids)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when getting cached item")
	err = json.NewDecoder(response.Body).Decode(&items)
	common.SilentCheck(err, "when decoding response from cache")
	return items
}

func buildShowReply(item domain.Item) string {
	return fmt.Sprintf("🔑 ID: %s\n"+
		"🗓 Date: %s\n"+
		"💡 Title: %s\n"+
		"🔗 Link: <a href='%s'>Here</a>\n"+
		"📖 Content: %s\n"+
		"🗳 Source: %s\n"+
		"📍 Location: %s\n"+
		"🏷 Category: %s\n",
		item.Id, item.Timestamp, item.Title, item.Link, item.Content, item.SourceName, item.Location, item.Category)
}

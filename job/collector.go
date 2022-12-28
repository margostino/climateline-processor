package job

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/internal"
	"net/http"
	"strings"
)

func Collect(request *http.Request, writer *http.ResponseWriter) {
	var items = make([]*domain.Item, 0)
	twitterApi = newTwitterApi()

	category := strings.ToLower(request.URL.Query().Get("category"))
	if category == "" {
		category = "*"
	}
	urls = config.GetUrls(category)

	items, err := internal.FetchNews(category)

	var botNotifications = 0
	for _, item := range items {
		if item.ShouldNotifyBot {
			botNotifications += 1
			internal.NotifyBot(item)
		}
	}

	response := domain.JobResponse{
		Items:            len(items),
		BotNotifications: botNotifications,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		(*writer).WriteHeader(http.StatusNotFound)
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		if botNotifications > 0 {
			(*writer).WriteHeader(http.StatusOK)
		} else {
			(*writer).WriteHeader(http.StatusNoContent)
		}
		
		(*writer).Write(jsonResp)
	}
}

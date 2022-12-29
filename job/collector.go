package job

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/internal"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Collect(request *http.Request, writer *http.ResponseWriter) {
	var items = make([]*domain.Item, 0)
	twitterApi = newTwitterApi()

	collectForced, err := strconv.ParseBool(request.URL.Query().Get("collect_forced"))
	if err != nil {
		collectForced = false
	}
	category := strings.ToLower(request.URL.Query().Get("category"))
	if category == "" {
		category = "*"
	}

	log.Printf("Category: %s\n", category)
	log.Printf("Collect Forced: %t\n", collectForced)

	urls = config.GetUrls(category)

	items, err = internal.FetchNews(category, collectForced)

	var botNotifications = 0
	for _, item := range items {
		if collectForced || item.ShouldNotifyBot {
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
			(*writer).Write(jsonResp)
		} else {
			(*writer).WriteHeader(http.StatusNoContent)
		}
	}
}

package bot

import (
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/internal"
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

	items, err := internal.FetchNews(category, true)

	if !common.IsError(err, "when fetching news") {
		//for _, item := range items {
			//internal.NotifyBot(item)
		//}
		reply = "✅ Completed successfully"
	} else {
		reply = "🔴 Fetcher failed"
	}

	return reply
}

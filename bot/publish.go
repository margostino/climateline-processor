package bot

import (
	"fmt"
	"github.com/margostino/climateline-processor/internal"
	"strings"
)

func Publish(input string) string {
	var reply string
	var category string

	params := strings.Split(input, " ")
	if len(params) > 1 {
		category = params[1]
	} else {
		category = "*"
	}

	response, err := internal.PublishNews(category, true)

	if err != nil {
		reply = fmt.Sprintf("🔴 Fetcher failed: %s", err.Error())
	} else {
		reply = fmt.Sprintf("✅ Completed successfully (%d items published)", response.Items)
	}

	return reply
}

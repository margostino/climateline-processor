package bot

import (
	"fmt"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/job"
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

	if fetchItems(category) {
		reply = "✅ Completed successfully"
	} else {
		reply = "🔴 Fetcher failed"
	}

	return reply
}

func fetchItems(category string) bool {
	client := &http.Client{}
	url := fmt.Sprintf("%s?category=%s", job.GetBaseJobUrl(), category)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when triggering job")
	return response.StatusCode == 200
}

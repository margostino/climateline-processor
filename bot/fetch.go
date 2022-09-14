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

	if fetchItems() {
		reply = "✅ Completed successfully"
	} else {
		reply = "🔴 Fetcher failed"
	}

	return reply
}

func ShouldFetch(input string) bool {
	sanitizedInput := SanitizeInput(input)
	return strings.Contains(sanitizedInput, "fetch")
}

func fetchItems() bool {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, job.GetBaseJobUrl(), nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when triggering job")
	return response.StatusCode == 200
}

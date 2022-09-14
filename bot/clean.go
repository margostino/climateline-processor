package bot

import (
	"fmt"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"net/http"
	"os"
)

func Clean() string {
	var reply string

	if cleanCache() {
		reply = "🧹 cache deleted"
	} else {
		reply = "🔴 cache deletion failed!"
	}

	return reply
}

func cleanCache() bool {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, cache.GetBaseCacheUrl(), nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when cleaning up the cache")
	return response.StatusCode == 200
}

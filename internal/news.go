package internal

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/news"
	"net/http"
	"os"
)

func FetchNews(category string) ([]*domain.Item, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s?category=%s", news.GetBaseNewsUrl(), category)
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

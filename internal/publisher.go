package internal

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"os"
)

func PublishNews(category string, publishForced bool) (*domain.JobResponse, error) {
	client := &http.Client{}
	baseUrl := os.Getenv("PUBLISHER_BASE_URL")
	url := fmt.Sprintf("%s?category=%s&publish_forced=%t", baseUrl, category, publishForced)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	if err == nil && response.StatusCode == 200 {
		var jobResponse domain.JobResponse
		err := json.NewDecoder(response.Body).Decode(&jobResponse)
		return &jobResponse, err
	} else {
		return nil, err
	}
}

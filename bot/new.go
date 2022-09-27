package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"log"
	"net/http"
	"os"
	"strings"
)

func New(input string) string {
	var reply string
	var items = make([]*domain.Item, 0)
	instructions := strings.Split(input, "\n")
	item := &domain.Item{
		Timestamp:  instructions[1],
		Title:      instructions[2],
		SourceName: instructions[3],
		Location:   instructions[4],
		Link:       instructions[5],
		Category:   strings.ToLower(instructions[6]),
		Content:    instructions[7],
	}

	items = append(items, item)

	ids, err := addCachedItem(items)
	if err == nil {
		for _, id := range ids {
			reply = fmt.Sprintf("✅ article with ID %s created", id)
		}

		if reply == "" {
			reply = "⚠️ no item created"
		}
	} else {
		reply = "🔴 article failed!"
	}

	return reply
}

func addCachedItem(items []*domain.Item) ([]string, error) {
	var ids []string
	client := &http.Client{}
	jsonRequest, err := json.Marshal(items)

	if !common.IsError(err, "when adding cache") {
		request, err := http.NewRequest(http.MethodPost, cache.GetBaseCacheUrl(), bytes.NewBuffer(jsonRequest))
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)

		common.SilentCheck(err, "in response of caching")

		if err == nil && response.StatusCode != 201 {
			log.Printf("Adding cache was not successful. Status: %d\n", response.StatusCode)
		}
		err = json.NewDecoder(response.Body).Decode(&ids)
	}

	return ids, err
}

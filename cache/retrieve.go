package cache

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
)

func Create(request *http.Request) (map[string]domain.Item, error) {
	var items []domain.Item
	var results = make(map[string]domain.Item)
	err := json.NewDecoder(request.Body).Decode(&items)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		results[item.Id] = item
	}

	return results, nil
}

package cache

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
)

func Create(request *http.Request, writer *http.ResponseWriter, db *map[string]domain.Item) {
	var items []domain.Item

	(*writer).WriteHeader(http.StatusCreated)

	err := json.NewDecoder(request.Body).Decode(&items)
	if err != nil {
		http.Error(*writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, item := range items {
		(*db)[item.Id] = item
	}

}

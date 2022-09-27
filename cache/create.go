package cache

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"strconv"
)

func Create(request *http.Request, writer *http.ResponseWriter, db *map[string]domain.Item) {
	var items []domain.Item
	var ids []string

	(*writer).WriteHeader(http.StatusCreated)

	err := json.NewDecoder(request.Body).Decode(&items)
	if err != nil {
		http.Error(*writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, item := range items {
		if item.Id == "" {
			item.Id = strconv.Itoa(len(*db) + 1)
		}
		(*db)[item.Id] = item
		ids = append(ids, item.Id)
	}

	jsonResp, err := json.Marshal(ids)
	(*writer).Write(jsonResp)

}

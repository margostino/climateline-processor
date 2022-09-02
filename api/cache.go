package api

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"log"
	"net/http"
	"os"
	"strings"
)

type Request []domain.Item

var cache = make(map[string]domain.Item)
var baseCacheUrl = os.Getenv("CACHE_BASE_URL")

func Cache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {

		defer r.Body.Close()
		w.WriteHeader(http.StatusCreated)
		var items []domain.Item
		err := json.NewDecoder(r.Body).Decode(&items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, item := range items {
			cache[item.Id] = item
		}

	} else if r.Method == "PUT" {

		defer r.Body.Close()
		w.WriteHeader(http.StatusNoContent)
		var edit domain.Edit
		id := r.URL.Query().Get("id")
		err := json.NewDecoder(r.Body).Decode(&edit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cache[id] = domain.Item{
			Id:         id,
			Timestamp:  cache[id].Timestamp,
			Link:       cache[id].Link,
			Content:    cache[id].Content,
			Title:      edit.Title,
			SourceName: edit.SourceName,
			Location:   edit.Location,
			Category:   edit.Category,
		}

	} else if r.Method == "GET" {
		var items = make([]*domain.Item, 0)
		idsQuery := r.URL.Query().Get("ids")
		ids := strings.Split(idsQuery, ",")

		for _, id := range ids {
			if item, ok := cache[id]; ok {
				items = append(items, &item)
			} else {
				log.Printf("Item %s not found\n", id)
			}
		}

		if len(items) > 0 {
			w.WriteHeader(http.StatusOK)
			response, marshalErr := json.Marshal(items)
			if common.IsError(marshalErr, "when marshaling item response") {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.Write(response)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	} else if r.Method == "DELETE" {
		w.WriteHeader(http.StatusOK)
		cache = make(map[string]domain.Item)
	} else {

		w.WriteHeader(http.StatusBadRequest)

	}

	return
}

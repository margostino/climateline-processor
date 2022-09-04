package api

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/security"
	"log"
	"net/http"
	"os"
	"strings"
)

type Request []domain.Item

var cache = make(map[string]domain.Item)
var baseCacheUrl = os.Getenv("CACHE_BASE_URL")

func Cache(w http.ResponseWriter, r *http.Request) {

	if security.IsAuthorized(r) {
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
			var edit domain.Edit
			id := r.URL.Query().Get("id")

			if _, ok := cache[id]; ok {
				w.WriteHeader(http.StatusNoContent)
				err := json.NewDecoder(r.Body).Decode(&edit)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				var title, sourceName, location, category string

				if edit.Title == "" {
					title = cache[id].Title
				} else {
					title = edit.Title
				}

				if edit.SourceName == "" {
					sourceName = cache[id].SourceName
				} else {
					sourceName = edit.SourceName
				}

				if edit.Location == "" {
					location = cache[id].Location
				} else {
					location = edit.Location
				}

				if edit.Category == "" {
					category = cache[id].Category
				} else {
					category = edit.Category
				}

				cache[id] = domain.Item{
					Id:         id,
					Timestamp:  cache[id].Timestamp,
					Link:       cache[id].Link,
					Content:    cache[id].Content,
					Title:      title,
					SourceName: sourceName,
					Location:   location,
					Category:   category,
				}

			} else {
				w.WriteHeader(http.StatusNotFound)
			}

		} else if r.Method == "GET" {
			var items = make([]domain.Item, 0)
			idsQuery := r.URL.Query().Get("ids")

			if idsQuery == "*" {
				for _, item := range cache {
					items = append(items, item)
				}
			} else {
				ids := strings.Split(idsQuery, ",")
				for _, id := range ids {
					if item, ok := cache[id]; ok {
						items = append(items, item)
					} else if id == "*" {
						log.Println("Cache is empty")
					} else {
						log.Printf("Item %s not found\n", id)
					}
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
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}

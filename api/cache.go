package api

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"os"
	"strconv"
)

type Request []domain.Item

var cache = make(map[int]domain.Item)
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

	} else if r.Method == "GET" {

		idStr := r.URL.Query().Get("id")
		id, parseErr := strconv.Atoi(idStr)

		if common.IsError(parseErr, "when parsing ID from path") {
			w.WriteHeader(http.StatusBadRequest)
		}

		if item, ok := cache[id]; ok {
			w.WriteHeader(http.StatusOK)
			response, marshalErr := json.Marshal(item)
			if common.IsError(marshalErr, "when marshaling item response") {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.Write(response)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	} else if r.Method == "DELETE" {
		w.WriteHeader(http.StatusCreated)
		cache = make(map[int]domain.Item)
	} else {

		w.WriteHeader(http.StatusBadRequest)

	}

	return
}

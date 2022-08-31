package api

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/common"
	"net/http"
	"strconv"
	"strings"
)

type Item struct {
	Id        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Status    string `json:"status"`
}

type Request []Item

var cache = make(map[int]Item)

func Cache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {

		defer r.Body.Close()
		w.WriteHeader(http.StatusCreated)
		var items []Item
		err := json.NewDecoder(r.Body).Decode(&items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, item := range items {
			cache[item.Id] = item
		}

	} else if r.Method == "GET" {

		idStr := strings.TrimPrefix(r.URL.Path, "/cache/items/")
		id, parseErr := strconv.Atoi(idStr)

		if common.Fail(parseErr, "when parsing ID from path") {
			w.WriteHeader(http.StatusBadRequest)
		}

		if item, ok := cache[id]; ok {
			w.WriteHeader(http.StatusOK)
			response, marshalErr := json.Marshal(item)
			if common.Fail(marshalErr, "when marshaling item response") {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.Write(response)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	} else if r.Method == "DELETE" {
		w.WriteHeader(http.StatusCreated)
		cache = make(map[int]Item)
	} else {

		w.WriteHeader(http.StatusBadRequest)

	}

	return
}

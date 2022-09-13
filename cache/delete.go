package cache

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"log"
	"net/http"
	"strings"
)

func Retrieve(request *http.Request, writer *http.ResponseWriter, db *map[string]domain.Item) {
	var items = make([]domain.Item, 0)
	idsQuery := request.URL.Query().Get("ids")

	if idsQuery == "*" {
		for _, item := range *db {
			items = append(items, item)
		}
	} else {
		ids := strings.Split(idsQuery, ",")
		for _, id := range ids {
			if item, ok := (*db)[id]; ok {
				items = append(items, item)
			} else if id == "*" {
				log.Println("Cache is empty")
			} else {
				log.Printf("Item %s not found\n", id)
			}
		}

	}

	if len(items) > 0 {
		(*writer).WriteHeader(http.StatusOK)
		response, marshalErr := json.Marshal(items)
		if common.IsError(marshalErr, "when marshaling item response") {
			(*writer).WriteHeader(http.StatusBadRequest)
		} else {
			(*writer).Write(response)
		}
	} else {
		(*writer).WriteHeader(http.StatusNotFound)
	}

}

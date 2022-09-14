package cache

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
)

func Update(request *http.Request, writer *http.ResponseWriter, db *map[string]domain.Item) {
	var edit domain.Update
	id := request.URL.Query().Get("id")

	if _, ok := (*db)[id]; ok {
		(*writer).WriteHeader(http.StatusNoContent)
		err := json.NewDecoder(request.Body).Decode(&edit)
		if err != nil {
			http.Error(*writer, err.Error(), http.StatusBadRequest)
			return
		}

		var title, sourceName, location, category string

		if edit.Title == "" {
			title = (*db)[id].Title
		} else {
			title = edit.Title
		}

		if edit.SourceName == "" {
			sourceName = (*db)[id].SourceName
		} else {
			sourceName = edit.SourceName
		}

		if edit.Location == "" {
			location = (*db)[id].Location
		} else {
			location = edit.Location
		}

		if edit.Category == "" {
			category = (*db)[id].Category
		} else {
			category = edit.Category
		}

		(*db)[id] = domain.Item{
			Id:         id,
			Timestamp:  (*db)[id].Timestamp,
			Link:       (*db)[id].Link,
			Content:    (*db)[id].Content,
			Title:      title,
			SourceName: sourceName,
			Location:   location,
			Category:   category,
		}

	} else {
		(*writer).WriteHeader(http.StatusNotFound)
	}
}

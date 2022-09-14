package api

import (
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/security"
	"net/http"
)

type Request []domain.Item

var db = make(map[string]domain.Item)

func Cache(w http.ResponseWriter, r *http.Request) {

	if security.IsAuthorized(r) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "POST":
			cache.Create(r, &w, &db)
		case "PUT":
			cache.Update(r, &w, &db)
		case "GET":
			cache.Retrieve(r, &w, &db)
		case "DELETE":
			cache.Delete(&w, &db)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}

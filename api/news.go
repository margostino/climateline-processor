package api

import (
	"github.com/margostino/climateline-processor/news"
	"github.com/margostino/climateline-processor/security"
	"net/http"
)

func News(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if security.IsAuthorized(r) {
		news.Fetch(r, &w)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}

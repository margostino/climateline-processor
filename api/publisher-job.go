package api

import (
	"github.com/margostino/climateline-processor/job"
	"github.com/margostino/climateline-processor/security"
	"net/http"
)

func ExecutePublisherJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if security.IsAuthorized(r) {
		job.Publish(r, &w)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}

package security

import (
	"fmt"
	"net/http"
	"os"
)

func IsAuthorized(r *http.Request) bool {
	jobSecret := os.Getenv("CLIMATELINE_JOB_SECRET")
	requestSecret := r.Header.Get("Authorization")
	return requestSecret == fmt.Sprintf("Bearer %s", jobSecret)
}

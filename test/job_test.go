package test

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/api"
	"github.com/margostino/climateline-processor/domain"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestJobUnauthorized(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/job", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Job)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestRunJobNewItem(t *testing.T) {
	request, err := mockJobRequest()
	if err != nil {
		t.Fatal(err)
	}

	feedContent := mockRssFeed()

	feedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(feedContent))
		require.NoError(t, err)
	}))

	cacheServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		var items []domain.Item
		err = json.NewDecoder(r.Body).Decode(&items)
		require.NoError(t, err)
		if len(items) != 1 {
			t.Errorf("unexpected caching items size: got %v want %v", len(items), 1)
		}
	}))

	os.Setenv("TELEGRAM_BOT_TOKEN", "dummy")
	os.Setenv("FEED_URL", feedServer.URL)
	os.Setenv("CACHE_BASE_URL", cacheServer.URL)

	defer feedServer.Close()
	defer cacheServer.Close()

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Job)

	handler.ServeHTTP(response, &request)

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var jobResponse domain.JobResponse
	err = json.NewDecoder(response.Body).Decode(&jobResponse)

	if jobResponse.Items != 1 {
		t.Errorf("handler returned unexpected response size: got %v want %v", jobResponse.Items, 1)
	}
}

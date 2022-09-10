package test

import (
	"github.com/margostino/climateline-processor/api"
	"net/http"
	"net/http/httptest"
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
}

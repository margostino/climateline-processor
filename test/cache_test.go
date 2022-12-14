package test

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/api"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCacheUnauthorized(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/cache", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Cache)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestNewItem(t *testing.T) {
	items := mockItems("mock.title")
	response, err := postCache(items)

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestGetItem(t *testing.T) {
	items := mockItems("mock.title")
	_, err := postCache(items)

	if err != nil {
		t.Fatal(err)
	}

	response, err := getCache("mock.id")

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var itemsResponse []domain.Item
	err = json.NewDecoder(response.Body).Decode(&itemsResponse)

	if len(itemsResponse) != 1 {
		t.Errorf("handler returned unexpected response size: got %v want %v", len(itemsResponse), 1)
	}

}

func TestDeleteItem(t *testing.T) {
	items := mockItems("mock.title")
	_, err := postCache(items)

	if err != nil {
		t.Fatal(err)
	}

	response, err := postDelete()

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	response, err = getCache("*")

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

}

func TestPutItem(t *testing.T) {
	items := mockItems("mock.title")
	_, err := postCache(items)

	if err != nil {
		t.Fatal(err)
	}

	response, err := putCache("mock.id", "another title update")

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	response, err = getCache("mock.id")

	if err != nil {
		t.Fatal(err)
	}

	var itemsResponse []domain.Item
	err = json.NewDecoder(response.Body).Decode(&itemsResponse)

	if len(itemsResponse) != 1 {
		t.Errorf("handler returned unexpected response size: got %v want %v", len(itemsResponse), 1)
	}

	oldTitle := "mock.title"
	newTitle := itemsResponse[0].Title

	if newTitle == oldTitle {
		t.Errorf("handler returned unexpected update title response: got %v want different than %v", newTitle, oldTitle)
	}

}

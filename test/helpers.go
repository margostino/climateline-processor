package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/api"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"net/http/httptest"
	"os"
)

func mockItems(title string) []domain.Item {
	item := mockItem(title)
	items := append(make([]domain.Item, 0), item)
	return items
}

func mockItem(title string) domain.Item {
	return domain.Item{
		Id:         "mock.id",
		Timestamp:  "2022-09-04T02:36:21Z",
		Title:      title,
		Link:       "mock.com",
		Content:    "mock some content",
		SourceName: "Test",
		Location:   "Somewhere",
		Category:   "warming",
	}
}

func mockCachePostRequest() (http.Request, error) {
	items := mockItems("mock.title")
	json, err := json.Marshal(items)
	body := bytes.NewBuffer(json)
	request, err := http.NewRequest(http.MethodPost, "/cache", body)
	setSecret(request)
	return *request, err
}

func mockCachePutRequest(id string, newTitle string) (http.Request, error) {
	item := mockItem(newTitle)
	json, err := json.Marshal(item)
	body := bytes.NewBuffer(json)
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/cache?id=%s", id), body)
	setSecret(request)
	return *request, err
}

func mockCacheDeleteRequest() (http.Request, error) {
	request, err := http.NewRequest(http.MethodDelete, "/cache", nil)
	setSecret(request)
	return *request, err
}

func mockCacheGetRequest(id string) (http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/cache?ids=%s", id), nil)
	setSecret(request)
	return *request, err
}

func setSecret(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	request.Header.Set("Content-Type", "application/json")
}

func postCache() (*httptest.ResponseRecorder, error) {
	request, err := mockCachePostRequest()
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Cache)

	handler.ServeHTTP(rr, &request)

	return rr, nil
}

func postDelete() (*httptest.ResponseRecorder, error) {
	request, err := mockCacheDeleteRequest()
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Cache)

	handler.ServeHTTP(rr, &request)

	return rr, nil
}

func getCache(id string) (*httptest.ResponseRecorder, error) {
	request, err := mockCacheGetRequest(id)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Cache)

	handler.ServeHTTP(rr, &request)

	return rr, nil
}

func putCache(id string, newTitle string) (*httptest.ResponseRecorder, error) {
	request, err := mockCachePutRequest(id, newTitle)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Cache)

	handler.ServeHTTP(rr, &request)
	return rr, nil
}

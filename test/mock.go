package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
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

func mockJobRequest() (http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, "/job", nil)
	setJobSecret(request)
	return *request, err
}

func mockBotRequest(message *BotRequest, secret string) (http.Request, error) {
	json, err := json.Marshal(message)
	body := bytes.NewBuffer(json)
	request, err := http.NewRequest(http.MethodPost, "/bot", body)
	setBotSecret(request, secret)
	return *request, err
}

func mockCachePostRequest(items []domain.Item) (http.Request, error) {
	json, err := json.Marshal(items)
	body := bytes.NewBuffer(json)
	request, err := http.NewRequest(http.MethodPost, "/cache", body)
	setJobSecret(request)
	return *request, err
}

func mockCachePutRequest(id string, newTitle string) (http.Request, error) {
	item := mockItem(newTitle)
	json, err := json.Marshal(item)
	body := bytes.NewBuffer(json)
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/cache?id=%s", id), body)
	setJobSecret(request)
	return *request, err
}

func mockCacheDeleteRequest() (http.Request, error) {
	request, err := http.NewRequest(http.MethodDelete, "/cache", nil)
	setJobSecret(request)
	return *request, err
}

func mockCacheGetRequest(id string) (http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/cache?ids=%s", id), nil)
	setJobSecret(request)
	return *request, err
}

func mockRssFeed() string {
	return `
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom" xmlns:idx="urn:atom-extension:indexing"> 
			<id>tag:google.com,2005:reader/user/12586952400243799274/state/com.google/alerts/6958349095823097994</id> 
			<title>Google Alert - climate change breaking news worldwide</title> 
			<link href="https://www.google.com/alerts/feeds/12586952400243799274/6958349095823097994" rel="self"></link> 
			<updated>2022-09-10T02:00:38Z</updated> 
			<entry> 
				<id>tag:google.com,2013:googlealerts/feed:13216910124206463966</id> 
				<title type="html">Something happened</title> 
				<link href="some.com"></link> 
				<published>2022-09-10T02:00:38Z</published> 
				<updated>2022-09-10T02:00:38Z</updated> 
				<content type="html">Some some some somewhere</content> 
				<author> <name></name> </author>         
			</entry>
		</feed> 
	`
}

func mockBotMessageRequest(message string) *BotRequest {
	return &BotRequest{
		UpdateId: 1,
		Message: &BotMessage{
			MessageId: 1,
			Text:      message,
			From: &BotFrom{
				Id:        1,
				FirstName: "mock.name",
				Username:  "mock.username",
			},
			Chat: &BotChat{
				Id: 1,
			},
		},
	}
}

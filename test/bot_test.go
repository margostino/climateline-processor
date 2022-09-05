package test

import (
	"bytes"
	"encoding/json"
	"github.com/margostino/climateline-processor/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBotUnauthorized(t *testing.T) {
	request := &BotRequest{
		UpdateId: 1,
		Message: &BotMessage{
			MessageId: 1,
			Text:      "testing mock",
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
	json, err := json.Marshal(request)
	body := bytes.NewBuffer(json)
	req, err := http.NewRequest(http.MethodPost, "/bot", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Bot)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

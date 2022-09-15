package test

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/api"
	"github.com/margostino/climateline-processor/domain"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestBotUnauthorized(t *testing.T) {
	message := mockBotMessageRequest("testing mock")

	request, err := mockBotRequest(message, "invalid-secret")
	require.NoError(t, err)

	response := call(&request)
	assertStatus(response, http.StatusUnauthorized, t)
}

func TestInvalidInput(t *testing.T) {
	message := mockBotMessageRequest("testing mock")

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	require.NoError(t, err)

	response := call(&request)
	assertStatus(response, http.StatusOK, t)

	botResponse, err := parseBotResponse(response)
	require.NoError(t, err)
	assertResponse(botResponse.ChatId != 1, botResponse.ChatId, 1, t)
	assertResponse(botResponse.Text != "Input is not valid", botResponse.Text, "Input is not valid", t)
}

func TestSingleShowInput(t *testing.T) {
	message := mockBotMessageRequest("show 1")

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	assertError(err, t)

	items := mockItems("mock.title")

	queries := make(map[string]string)
	queries["ids"] = "1"
	cacheServer := NewMockServer().
		withMethod(http.MethodGet).
		withStatus(http.StatusOK).
		withBody(items).
		withQueries(queries).
		start(t)

	os.Setenv("CACHE_BASE_URL", cacheServer.Url)

	response := call(&request)
	botResponse, err := parseBotResponse(response)
	require.NoError(t, err)
	assertStatus(response, http.StatusOK, t)
	assertResponse(botResponse.ChatId != 1, botResponse.ChatId, 1, t)
	assertResponse(!strings.Contains(botResponse.Text, "ID: mock.id"), botResponse.Text, "One item with ID: mock.id", t)
}

func TestShowAllInput(t *testing.T) {
	message := &BotRequest{
		UpdateId: 1,
		Message: &BotMessage{
			MessageId: 1,
			Text:      "show",
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

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	if err != nil {
		t.Fatal(err)
	}

	items := mockItems("mock.title")

	cacheServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("unexpected method for fetching caching: got %v want %v", r.Method, "GET")
		}
		if r.URL.Query().Get("ids") != "*" {
			t.Errorf("unexpected method for fetching caching: got %v want %v", r.URL.Query().Get("ids"), "*")
		}
		w.WriteHeader(http.StatusOK)
		res, err := json.Marshal(items)
		w.Write(res)
		require.NoError(t, err)
	}))

	os.Setenv("CACHE_BASE_URL", cacheServer.URL)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Bot)

	handler.ServeHTTP(response, &request)

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var botResponse domain.BotResponse
	err = json.NewDecoder(response.Body).Decode(&botResponse)

	if botResponse.ChatId != 1 {
		t.Errorf("handler returned unexpected chat ID: got %v want %v", botResponse.ChatId, 1)
	}

	if !strings.Contains(botResponse.Text, "ID: mock.id") {
		t.Errorf("handler returned unexpected chat ID: got %v want %v", botResponse.Text, "One item with ID: mock.id")
	}
}

func TestCleanInput(t *testing.T) {
	message := &BotRequest{
		UpdateId: 1,
		Message: &BotMessage{
			MessageId: 1,
			Text:      "clean",
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

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	if err != nil {
		t.Fatal(err)
	}

	cacheServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("unexpected method for fetching caching: got %v want %v", r.Method, "DELETE")
		}
		w.WriteHeader(http.StatusOK)
		require.NoError(t, err)
	}))

	os.Setenv("CACHE_BASE_URL", cacheServer.URL)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Bot)

	handler.ServeHTTP(response, &request)

	if err != nil {
		t.Fatal(err)
	}

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var botResponse domain.BotResponse
	err = json.NewDecoder(response.Body).Decode(&botResponse)

	if botResponse.ChatId != 1 {
		t.Errorf("handler returned unexpected chat ID: got %v want %v", botResponse.ChatId, 1)
	}

	if !strings.Contains(botResponse.Text, "cache deleted") {
		t.Errorf("handler returned unexpected chat ID: got %v want %v", botResponse.Text, "One item with ID: cache deleted")
	}
}

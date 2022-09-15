package test

import (
	"github.com/stretchr/testify/require"
	"net/http"
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
	message := mockBotMessageRequest("show")

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	assertError(err, t)

	items := mockItems("mock.title")

	queries := make(map[string]string)
	queries["ids"] = "*"
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
	assertResponse(botResponse.ChatId != 1, botResponse.ChatId, "*", t)
	assertResponse(!strings.Contains(botResponse.Text, "ID: mock.id"), botResponse.Text, "One item with ID: mock.id", t)
}

func TestCleanInput(t *testing.T) {
	message := mockBotMessageRequest("clean")

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	assertError(err, t)

	cacheServer := NewMockServer().
		withMethod(http.MethodDelete).
		withStatus(http.StatusOK).
		start(t)

	os.Setenv("CACHE_BASE_URL", cacheServer.Url)

	response := call(&request)
	botResponse, err := parseBotResponse(response)
	require.NoError(t, err)
	assertStatus(response, http.StatusOK, t)
	assertResponse(botResponse.ChatId != 1, botResponse.ChatId, 1, t)
	assertResponse(!strings.Contains(botResponse.Text, "cache deleted"), botResponse.Text, "One item with ID: cache deleted", t)
}

func TestFetchInput(t *testing.T) {
	message := mockBotMessageRequest("fetch")

	request, err := mockBotRequest(message, os.Getenv("TELEGRAM_BOT_SECRET"))
	assertError(err, t)

	items := mockItems("mock.title")
	cacheServer := NewMockServer().
		withMethod(http.MethodPost).
		withStatus(http.StatusOK).
		withBody(items).
		start(t)

	feeds := mockRssFeed()
	feedServer := NewMockServer().
		withMethod(http.MethodGet).
		withStatus(http.StatusOK).
		withBody(feeds).
		start(t)

	jobServer := NewMockServer().
		withMethod(http.MethodGet).
		withStatus(http.StatusOK).
		start(t)

	os.Setenv("CACHE_BASE_URL", cacheServer.Url)
	os.Setenv("FEED_URL", feedServer.Url)
	os.Setenv("JOB_BASE_URL", jobServer.Url)

	response := call(&request)
	botResponse, err := parseBotResponse(response)
	require.NoError(t, err)
	assertStatus(response, http.StatusOK, t)
	assertResponse(botResponse.ChatId != 1, botResponse.ChatId, 1, t)
	assertResponse(!strings.Contains(botResponse.Text, "Completed successfully"), botResponse.Text, "Completed successfully", t)
}

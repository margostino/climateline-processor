package test

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/api"
	"github.com/margostino/climateline-processor/bot"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setJobSecret(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	request.Header.Set("Content-Type", "application/json")
}

func setBotSecret(request *http.Request, secret string) {
	request.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)
	request.Header.Set("Content-Type", "application/json")
}

func postCache(items []domain.Item) (*httptest.ResponseRecorder, error) {
	request, err := mockCachePostRequest(items)
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

func parseBotResponse(response *httptest.ResponseRecorder) (*domain.BotResponse, error) {
	var botResponse domain.BotResponse
	err := json.NewDecoder(response.Body).Decode(&botResponse)
	return &botResponse, err
}

func createBotMessage(input string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Text: input,
	}
}

func validateInput(t *testing.T, input string) {
	assertValidInput(bot.IsValidInput(createBotMessage(input)), input, t)
}

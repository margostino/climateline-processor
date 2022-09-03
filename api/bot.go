package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/security"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Response struct {
	Msg    string `json:"text"`
	ChatID int64  `json:"chat_id"`
	Method string `json:"method"`
}

var githubClient = getGithubClient()

func Bot(w http.ResponseWriter, r *http.Request) {

	log.Printf("Method: %s "+
		"Proto: %s "+
		"User-Agent: %s, "+
		"Host: %s, "+
		"RequestURI: %s, "+
		"RemoteAddr: %s",
		r.Method, r.Proto, r.Header.Get("User-Agent"), r.Host, r.RequestURI, r.RemoteAddr)

	var reply string
	body, _ := ioutil.ReadAll(r.Body)
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Fatal("Error updating →", err)
	}

	log.Printf("[%s@%d] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

	if security.IsAdmin(update.Message.From.UserName, update.Message.Chat.ID, r) {
		defer r.Body.Close()
		w.Header().Add("Content-Type", "application/json")

		input := update.Message.Text
		if isValidInput(input) {
			if shouldPush(input) {
				ids := extractIds(input)
				items := getCachedItems(ids)
				reply = items[0].Link

				for _, item := range items {
					content := generateArticle(&item)
					message := "new article from workflow"
					options := &github.RepositoryContentFileOptions{
						Content: []byte(content),
						Message: &message,
					}
					path := fmt.Sprintf("articles/%s.md", strings.ReplaceAll(strings.ToLower(item.Title), " ", "-"))
					contentResponse, response, err := githubClient.Repositories.CreateFile(context.Background(), "margostino", "climateline", path, options)
					common.SilentCheck(err, "when creating new article on repository")
					println(contentResponse)
					println(response)
				}

			} else if shouldRunJob(input) {

				if runJob() {
					reply = "✅ Job completed successfully"
				} else {
					reply = "🔴 Job failed"
				}

			} else if shouldEdit(input) {
				instructions := strings.Split(input, "\n")
				edit := &domain.Edit{
					Title:      instructions[1],
					SourceName: instructions[2],
					Location:   instructions[3],
					Category:   instructions[4],
				}
				id := extractId(instructions[0])

				if updateCachedItems(id, edit) {
					reply = "✅ article updated"
				} else {
					reply = "🔴 article update failed!"
				}

			} else {
				reply = "👌"
			}
		} else {
			reply = "Input is not valid"
			log.Println(reply)
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		reply := "Unauthorized to handle Bot communication"
		log.Printf(reply)
	}

	data := Response{
		Msg:    reply,
		Method: "sendMessage",
		ChatID: update.Message.Chat.ID,
	}

	message, _ := json.Marshal(data)
	log.Printf("Response %s", string(message))
	fmt.Fprintf(w, string(message))
}

func sanitizeInput(input string) string {
	return strings.Trim(strings.ToLower(input), " ")
}

func isValidInput(input string) bool {
	sanitizedInput := sanitizeInput(input)
	match, err := regexp.MatchString(`^((push ([0-9]+\s*)+)|(edit [0-9]+\n.*?\n.*?\n.*?\n(agreements|assessment|awareness|warming|wildfires))|run)$`, sanitizedInput)
	common.SilentCheck(err, "when matching input with regex")
	return match
}

func shouldPush(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "push")
}

func shouldRunJob(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "run")
}

func shouldEdit(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "edit")
}

func extractIds(input string) string {
	ids := strings.Join(strings.Split(strings.TrimPrefix(input, "push "), " "), ",")
	return ids
}

func extractId(input string) string {
	return strings.TrimPrefix(input, "edit ")
}

func getCachedItems(ids string) []domain.Item {
	client := &http.Client{}
	var items []domain.Item
	url := fmt.Sprintf("%s?ids=%s", baseCacheUrl, ids)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when getting cached item")
	err = json.NewDecoder(response.Body).Decode(&items)
	common.SilentCheck(err, "when decoding response from cache")
	return items
}

func updateCachedItems(id string, edit *domain.Edit) bool {
	client := &http.Client{}
	url := fmt.Sprintf("%s?id=%s", baseCacheUrl, id)
	json, err := json.Marshal(edit)

	if !common.IsError(err, "when marshaling edit data") {
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
		response, err := client.Do(request)
		common.SilentCheck(err, "when updating cached item")
		return response.StatusCode == 204
	}
	return false
}

func runJob() bool {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, baseJobUrl, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when triggering job")
	return response.StatusCode == 200
}

func generateArticle(item *domain.Item) string {
	var icon string
	switch item.Category {
	case "agreements":
		icon = "handshake"
	case "assessment":
		icon = "file-text"
	case "awareness":
		icon = "seedling"
	case "warming":
		icon = "thermometer-three-quarters"
	case "wildfires":
		icon = "fires"
	}

	return fmt.Sprintf("---\n"+
		"title: '%s'\n"+
		"date: '%s'\n"+
		"source_url: '%s'\n"+
		"source_name: '%s'\n"+
		"location: '%s'\n"+
		"icon: %s\n"+
		"---\n\n"+
		"%s\n",
		item.Title, item.Timestamp, item.Link, item.SourceName, item.Location, icon, item.Content)
}

func getGithubClient() *github.Client {
	var githubAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

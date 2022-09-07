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

	log.Printf("[%s@%d] %s", update.Message.From.UserName, update.Message.Chat.ID, "update.Message.Text")

	if security.IsAdmin(update.Message.From.UserName, update.Message.Chat.ID, r) {
		defer r.Body.Close()
		w.Header().Add("Content-Type", "application/json")

		input := update.Message.Text
		if isValidInput(input) {
			if shouldPush(input) {
				ids := extractIds(input, "push ")
				items := getCachedItems(ids)

				if len(items) > 0 {
					for _, item := range items {
						content := generateArticle(&item)
						message := "new article from workflow"
						options := &github.RepositoryContentFileOptions{
							Content: []byte(content),
							Message: &message,
						}
						path := fmt.Sprintf("articles/%s.md", strings.ReplaceAll(strings.ToLower(item.Title), " ", "-"))
						_, response, err := githubClient.Repositories.CreateFile(context.Background(), "margostino", "climateline", path, options)
						common.SilentCheck(err, "when creating new article on repository")

						if response.StatusCode == 201 {
							reply = "✅ New article uploaded"
						} else {
							reply = fmt.Sprintf("🔴 Upload failed with status %s", response.Status)
						}
					}
				} else {
					reply = "⚠️ There are not items to upload"
				}

			} else if shouldFetch(input) {

				if fetchItems() {
					reply = "✅ Completed successfully"
				} else {
					reply = "🔴 Fetcher failed"
				}

			} else if shouldShow(input) {
				var id string
				if input == "/show" || input == "show" || input == "show all" || input == "show *" {
					id = "*"
				} else {
					id = extractIds(input, "show ")
				}

				items := getCachedItems(id)

				if len(items) > 0 {
					reply = buildShowReply(items[0])
				} else if id == "*" || id == "" || id == " " || id == "show" {
					reply = "Cache is empty"
				} else {
					reply = fmt.Sprintf("🤷‍ There is not item for ID %s", id)
				}

			} else if shouldEditProperty(input) {
				var edit *domain.Edit
				sanitizedInput := sanitizeInput(input)
				params := strings.Split(sanitizedInput, " ")
				property := params[0]
				id := params[1]
				value := common.NewString(input).
					TrimPrefix(fmt.Sprintf("%s %s", property, id)).
					TrimPrefix(" ").
					Value()

				if property == "category" {
					edit = &domain.Edit{
						Category: value,
					}
				} else if property == "location" {
					edit = &domain.Edit{
						Location: value,
					}
				} else if property == "title" {
					edit = &domain.Edit{
						Title: value,
					}
				} else {
					edit = &domain.Edit{
						SourceName: value,
					}
				}
				if updateCachedItems(id, edit) {
					reply = "✅ article updated"
				} else {
					reply = "🔴 article update failed!"
				}
			} else if shouldEdit(input) {
				instructions := strings.Split(input, "\n")
				edit := &domain.Edit{
					Title:      instructions[1],
					SourceName: instructions[2],
					Location:   instructions[3],
					Category:   instructions[4],
				}
				id := extractIds(instructions[0], "edit ")

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
		reply = "Unauthorized to handle Bot communication"
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
	return common.NewString(input).
		ToLower().
		Trim(" ").
		Value()
}

func isValidInput(input string) bool {
	sanitizedInput := sanitizeInput(input)
	match, err := regexp.MatchString(`^((push ([0-9]+\s*)+)|(edit [0-9]+\n.*?\n.*?\n.*?\n(agreements|assessment|awareness|warming|wildfires|floods|drought|health))|fetch|/fetch|show|/show|show [0-9]+|title [0-9]+ .*?|source [0-9]+ .*?|location [0-9]+ .*?|category [0-9]+ (agreements|assessment|awareness|warming|wildfires|floods|drought|health))$`, sanitizedInput)
	common.SilentCheck(err, "when matching input with regex")
	return match
}

func shouldPush(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "push")
}

func shouldEditProperty(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "category") ||
		strings.Contains(sanitizedInput, "title") ||
		strings.Contains(sanitizedInput, "source") ||
		strings.Contains(sanitizedInput, "location")
}

func shouldFetch(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "fetch")
}

func shouldEdit(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "edit")
}

func shouldShow(input string) bool {
	sanitizedInput := sanitizeInput(input)
	return strings.Contains(sanitizedInput, "show")
}

func extractIds(input string, prefix string) string {
	return common.NewString(input).
		ToLower().
		TrimPrefix(prefix).
		Split(" ").
		Join(",").
		Value()
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

func fetchItems() bool {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, baseJobUrl, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
	response, err := client.Do(request)
	common.SilentCheck(err, "when triggering job")
	return response.StatusCode == 200
}

func generateArticle(item *domain.Item) string {
	var icon string
	category := strings.ToLower(item.Category)
	switch category {
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
	case "floods":
		icon = "droplet"
	case "drought":
		icon = "droplet-slash"
	case "health":
		icon = "heart-pulse"
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

func buildShowReply(item domain.Item) string {
	return fmt.Sprintf("🔑 ID: %s\n"+
		"🗓 Date: %s\n"+
		"💡 Title: %s\n"+
		"🔗 Link: <a href='%s'>Here</a>\n"+
		"📖 Content: %s\n"+
		"🗳 Source: %s\n"+
		"📍 Location: %s\n"+
		"🏷 Category: %s\n",
		item.Id, item.Timestamp, item.Title, item.Link, item.Content, item.SourceName, item.Location, item.Category)
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

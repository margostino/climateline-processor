package api

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
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
	defer r.Body.Close()
	var reply string
	w.Header().Add("Content-Type", "application/json")

	//log.Printf("Cached Items (Reply): %d", len(cache.Items))

	body, _ := ioutil.ReadAll(r.Body)
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Fatal("Error updating →", err)
	}

	log.Printf("[%s@%d] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

	input := sanitizeInput(update.Message)
	if isValidInput(input) {
		if shouldUpload(input) {
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
				path := fmt.Sprintf("articles/%s.md", strings.ReplaceAll(strings.ToLower(item.Title), " ", "_"))
				contentResponse, response, err := githubClient.Repositories.CreateFile(context.Background(), "margostino", "climateline", path, options)
				common.SilentCheck(err, "when creating new article on repository")
				println(contentResponse)
				println(response)
			}

			// TODO: commit git new article
		} else {
			reply = "👌"
		}
	} else {
		reply = "Input is not valid"
		log.Println(reply)
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

func sanitizeInput(message *tgbotapi.Message) string {
	return strings.Trim(strings.ToLower(message.Text), " ")
}

func isValidInput(input string) bool {
	match, err := regexp.MatchString(`^((yes ([0-9]+\s*)+)|(no))$`, input)
	common.SilentCheck(err, "when matching input with regex")
	return match
}

func shouldUpload(input string) bool {
	return strings.Contains(input, "yes")
}

func extractIds(input string) string {
	ids := strings.Join(strings.Split(strings.TrimPrefix(input, "yes "), " "), ",")
	return ids
}

func getCachedItems(ids string) []domain.Item {
	var items []domain.Item
	url := fmt.Sprintf("%s?ids=%s", baseCacheUrl, ids)
	resp, err := http.Get(url)
	common.SilentCheck(err, "when getting cached item")
	err = json.NewDecoder(resp.Body).Decode(&items)
	common.SilentCheck(err, "when decoding response from cache")
	return items
}

func generateArticle(item *domain.Item) string {
	return fmt.Sprintf("---\n"+
		"title: '%s'\n"+
		"date: '%s'\n"+
		"source_url: '%s'\n"+
		"source_name: '%s'\n"+
		"location: '%s'\n"+
		"icon: %s\n"+
		"---\n\n"+
		"%s\n",
		item.Title, item.Timestamp, item.Link, "todo", "todo", "fire", item.Content)
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

package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/internal"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
)

func Collect(request *http.Request, writer *http.ResponseWriter) {
	var items = make([]*domain.Item, 0)
	twitterApi = newTwitterApi()

	category := strings.ToLower(request.URL.Query().Get("category"))
	if category == "" {
		category = "*"
	}
	urls = config.GetUrls(category)

	items, err := internal.FetchNews(category)

	githubClient := getGithubClient()
	inputs := map[string]interface{}{
		"environment": "production",
	}
	dispatcherEvent := github.CreateWorkflowDispatchEventRequest{
		Ref:    "master",
		Inputs: inputs,
	}
	r, err := githubClient.Actions.CreateWorkflowDispatchEventByFileName(context.TODO(), "margostino", "climateline-processor", "collector-dispatcher.yml", dispatcherEvent)
	common.SilentCheck(err, "when runner")
	println(r.StatusCode)

	for _, item := range items {
		if item.ShouldNotifyBot {
			internal.NotifyBot(item)
		}
	}

	response := domain.JobResponse{
		Items: len(items),
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		(*writer).WriteHeader(http.StatusNotFound)
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		(*writer).WriteHeader(http.StatusOK)
		(*writer).Write(jsonResp)
	}
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

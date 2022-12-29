package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/config"
	"github.com/margostino/climateline-processor/domain"
	"github.com/margostino/climateline-processor/internal"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var twitterApi *twitter.Client
var bitlyDomain = "bit.ly"

func Publish(request *http.Request, writer *http.ResponseWriter) {
	var items = make([]*domain.Item, 0)
	var twitterPosts = 0

	twitterApi = newTwitterApi()

	category := strings.ToLower(request.URL.Query().Get("category"))
	publishForced, err := strconv.ParseBool(request.URL.Query().Get("publish_forced"))
	if err != nil {
		publishForced = false
	}
	if category == "" {
		category = "*"
	}

	log.Printf("Category: %s\n", category)
	log.Printf("Publish Forced: %t\n", publishForced)

	urls = config.GetUrls(category)

	items, err = internal.FetchNews(category, publishForced)

	if len(items) > 0 {
		lastJobRun, err := getLastJobRun(category)
		if lastJobRun != nil && err == nil {
			log.Printf("Last Job Run %s\n", lastJobRun.String())
		}

		if !common.IsError(err, "when getting last job run") {
			for _, item := range items {
				shouldPublishByTime := lastJobRun != nil && (item.Published == lastJobRun || item.Published.After(*lastJobRun))
				shouldPublishByConfig := publishForced || item.ShouldNotifyTwitter
				if shouldPublishByConfig && shouldPublishByTime {
					_, err := notifyTwitter(item)
					if !common.IsError(err, "when posting tweet") {
						twitterPosts += 1
					}

				}
			}
		}
	}

	response := domain.JobResponse{
		Items:        len(items),
		TwitterPosts: twitterPosts,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		(*writer).WriteHeader(http.StatusNotFound)
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		if twitterPosts > 0 {
			(*writer).WriteHeader(http.StatusOK)
			(*writer).Write(jsonResp)
		} else {
			(*writer).WriteHeader(http.StatusNoContent)
			(*writer).Header().Set("Content-Length", "0")
		}
	}
}

func notifyTwitter(item *domain.Item) (*http.Response, error) {
	var tweet string
	//shorterLink := urlshortener.Shorten(item.Link)
	shorterLink := item.Link
	//params := &twitter.StatusUpdateParams{
	//	AttachmentURL: item.Link,
	//}
	title := sanitizeTweet(item.Title)

	if shorterLink != "" {
		tweet = fmt.Sprintf("%s\nSource: %s (%s)\n%s", title, item.SourceName, shorterLink, item.Tags)
	} else {
		tweet = fmt.Sprintf("%s\nSource: %s\n%s", title, item.SourceName, item.Tags)
	}

	_, resp, err := twitterApi.Statuses.Update(tweet, nil)
	if err != nil && resp.StatusCode == 200 {
		log.Println("Tweet created")
	}

	return resp, err
}

func sanitizeTweet(value string) string {
	return common.NewString(value).
		UnescapeString().
		ReplaceAll("<b>", "").
		ReplaceAll("</b>", "").
		Value()
}

func newTwitterApi() *twitter.Client {
	var consumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	var consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	var token = os.Getenv("TWITTER_TOKEN")
	var tokenSecret = os.Getenv("TWITTER_TOKEN_SECRET")
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := config.Client(oauth1.NoContext, oauthToken)
	client := twitter.NewClient(httpClient)
	return client
}

func getLastJobRun(category string) (*time.Time, error) {
	githubClient := getGithubClient()

	options := &github.ListWorkflowRunsOptions{
		Status: "success",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}

	workflowFilename := fmt.Sprintf("publisher-%s-job.yml", category)
	workflow, response, err := githubClient.Actions.ListWorkflowRunsByFileName(context.TODO(), "margostino", "climateline-processor", workflowFilename, options)

	if err == nil && response.StatusCode == 200 && len(workflow.WorkflowRuns) > 0 && &workflow.WorkflowRuns[0].UpdatedAt != nil {
		return &workflow.WorkflowRuns[0].UpdatedAt.Time, nil
	} else if response.StatusCode == 200 && len(workflow.WorkflowRuns) == 0 {
		veryOldDate := time.Now().AddDate(-1, 0, 0).UTC()
		return &veryOldDate, nil
	}

	return nil, err

}

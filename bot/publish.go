package bot

import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"strings"
)

func Publish(input string) string {
	var reply string
	var category string

	params := strings.Split(input, " ")
	if len(params) > 1 {
		category = params[1]
	} else {
		category = "*"
	}

	response, err := dispatchPublisherBy(category)

	if err != nil {
		reply = fmt.Sprintf("🔴 Publisher failed: %s", err.Error())
	} else {
		reply = fmt.Sprintf("✅ Completed successfully (status %d)", response.StatusCode)
	}

	return reply
}

func dispatchPublisherBy(category string) (*github.Response, error) {
	githubClient := getGithubClient()

	workflowFilename := fmt.Sprintf("publisher-%s-job.yml", category)
	inputs := map[string]interface{}{
		"category":      category,
		"publishForced": "true",
		"environment":   "production",
	}
	eventRequest := github.CreateWorkflowDispatchEventRequest{
		Ref:    "master",
		Inputs: inputs,
	}

	return githubClient.Actions.CreateWorkflowDispatchEventByFileName(context.TODO(), "margostino", "climateline-processor", workflowFilename, eventRequest)

}

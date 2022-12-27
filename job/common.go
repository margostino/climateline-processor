package job

import (
	"context"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/climateline-processor/config"
	"golang.org/x/oauth2"
	"os"
)

var urls []*config.UrlConfig

// GetBaseJobUrl Rather than a global and one-time assigment, this method is convenient for overriding on testing
func GetBaseJobUrl() string {
	return os.Getenv("JOB_BASE_URL")
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

package bot

import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

var githubClient *github.Client

func Push(input string, githubClient *github.Client) string {
	var reply string
	githubClient = getGithubClient()

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

	return reply
}

func ShouldPush(input string) bool {
	sanitizedInput := SanitizeInput(input)
	return strings.Contains(sanitizedInput, "push")
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

func getGithubClient() *github.Client {
	var githubAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

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

func PushReply(input string) string {
	var reply string

	lines := strings.Split(input, "\n")

	if len(lines) == 8 {
		var item domain.Item
		id := sanitizeReply(lines[0], domain.ID_PREFIX)
		timestamp := sanitizeReply(lines[1], domain.DATE_PREFIX)
		title := sanitizeReply(lines[2], domain.TITLE_PREFIX)
		link := sanitizeReply(lines[3], domain.LINK_PREFIX)
		content := sanitizeReply(lines[4], domain.CONTENT_PREFIX)
		source := sanitizeReply(lines[5], domain.SOURCE_PREFIX)
		location := sanitizeReply(lines[6], domain.LOCATION_PREFIX)
		category := sanitizeReply(lines[7], domain.CATEGORY_PREFIX)

		items := getCachedItems(id)

		if len(items) > 0 {
			item = items[0]
		} else {
			item = domain.Item{
				Id:         id,
				Timestamp:  timestamp,
				Title:      title,
				Link:       link,
				Content:    content,
				SourceName: source,
				Location:   location,
				Category:   category,
			}
		}
		reply = pushItem(&item)
	} else {
		reply = "No valid reply-push"
	}

	return reply
}

func Push(input string) string {
	var reply string

	ids := extractIds(input, "push ")
	items := getCachedItems(ids)

	if len(items) > 0 {
		for _, item := range items {
			reply = pushItem(&item)
		}
	} else {
		reply = "⚠️ There are not items to upload"
	}

	return reply
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
		icon = "fire"
	case "floods":
		icon = "droplet"
	case "drought":
		icon = "droplet-slash"
	case "health":
		icon = "heart-pulse"
	case "hurricane":
		icon = "hurricane"
	case "pollution":
		icon = "smog"
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
		sanitizeTitle(item.Title), item.Timestamp, item.Link, item.SourceName, item.Location, icon, item.Content)
}

func sanitizeReply(input string, prefix string) string {
	return common.NewString(input).
		TrimPrefix(prefix).
		Trim(" ").
		Value()
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

func sanitizeFilename(value string) string {
	return common.NewString(value).
		ToLower().
		ReplaceAll("<b>", "").
		ReplaceAll("</b>", "").
		ReplaceAll("&", "").
		ReplaceAll("#", "").
		ReplaceAll("|", "").
		ReplaceAll(";", "").
		ReplaceAll("\"", "").
		ReplaceAll("'", "").
		ReplaceAll(" ", "-").
		Value()
}

func sanitizeTitle(value string) string {
	return common.NewString(value).
		ReplaceAll("<b>", "").
		ReplaceAll("</b>", "").
		Value()
}

func pushItem(item *domain.Item) string {
	var reply string
	githubClient = getGithubClient()
	content := generateArticle(item)
	message := "new article from workflow"
	options := &github.RepositoryContentFileOptions{
		Content: []byte(content),
		Message: &message,
	}

	filename := sanitizeFilename(item.Title)
	path := fmt.Sprintf("articles/%s.md", filename)
	_, response, err := githubClient.Repositories.CreateFile(context.Background(), "margostino", "climateline", path, options)
	common.SilentCheck(err, "when creating new article on repository")

	if response.StatusCode == 201 {
		reply = "✅ New article uploaded"
	} else {
		reply = fmt.Sprintf("🔴 Upload failed with status %s", response.Status)
	}
	return reply
}

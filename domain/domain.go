package domain

import "time"

const (
	ID_PREFIX       = "🔑 ID:"
	DATE_PREFIX     = "🗓 Date:"
	TITLE_PREFIX    = "💡 Title:"
	LINK_PREFIX     = "🔗 Link:"
	SOURCE_PREFIX   = "📥 Source:"
	CONTENT_PREFIX  = "📖 Content:"
	LOCATION_PREFIX = "📍 Location:"
	CATEGORY_PREFIX = "🏷 Category:"
)

type Item struct {
	Id                  string     `json:"id"`
	Timestamp           string     `json:"timestamp"`
	Title               string     `json:"title"`
	Link                string     `json:"link"`
	Content             string     `json:"content"`
	SourceName          string     `json:"source_name"`
	Location            string     `json:"location"`
	Category            string     `json:"category"`
	Updated             *time.Time `json:"updated"`
	Published           *time.Time `json:"published"`
	Tags                string
	ShouldNotifyBot     bool
	ShouldNotifyTwitter bool
}

type Update struct {
	Title      string `json:"title"`
	SourceName string `json:"source_name"`
	Location   string `json:"location"`
	Category   string `json:"category"`
}

type JobResponse struct {
	Items            int `json:"items"`
	BotNotifications int `json:"bot_notifications"`
	TwitterPosts     int `json:"twitter_posts"`
}

type BotResponse struct {
	Text   string `json:"text"`
	ChatId int    `json:"chat_id"`
	Method string `json:"method"`
}

package domain

type Item struct {
	Id         string `json:"id"`
	Timestamp  string `json:"timestamp"`
	Title      string `json:"title"`
	Link       string `json:"link"`
	Content    string `json:"content"`
	SourceName string `json:"source_name"`
	Location   string `json:"location"`
	Category   string `json:"category"`
}

type Edit struct {
	Title      string `json:"title"`
	SourceName string `json:"source_name"`
	Location   string `json:"location"`
	Category   string `json:"category"`
}

type JobResponse struct {
	Items int `json:"items"`
}

type BotResponse struct {
	Text   string `json:"text"`
	ChatId int    `json:"chat_id"`
	Method string `json:"method"`
}

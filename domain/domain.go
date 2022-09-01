package domain

type Item struct {
	Id        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Content   string `json:"content"`
}

package domain

type Item struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Content   string `json:"content"`
}

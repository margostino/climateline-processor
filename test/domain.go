package test

type BotChat struct {
	Id int `json:"id"`
}

type BotFrom struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type BotMessage struct {
	MessageId int      `json:"message_id"`
	From      *BotFrom `json:"from"`
	Chat      *BotChat `json:"chat"`
	Text      string   `json:"text"`
}

type BotRequest struct {
	UpdateId int         `json:"update_id"`
	Message  *BotMessage `json:"message"`
}

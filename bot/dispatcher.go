package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"strings"
)

const (
	PUSH  = "push"
	FETCH = "fetch"
	CLEAN = "clean"
	EDIT  = "edit"
	SHOW  = "show"
)

func Reply(message *tgbotapi.Message) string {
	var reply string

	if message.ReplyToMessage != nil {
		input := message.ReplyToMessage.Text
		reply = PushReply(input)
	} else {
		input := message.Text
		command := common.NewString(input).
			ReplaceAll("/", "").
			ReplaceAll("_", " ").
			Value()
		sanitizedInput := SanitizeInput(command)
		commands := strings.Split(sanitizedInput, " ")

		if len(commands) > 0 {
			switch commands[0] {
			case PUSH:
				reply = Push(sanitizedInput)
			case EDIT:
				reply = Update(sanitizedInput)
			case FETCH:
				reply = Fetch(sanitizedInput)
			case SHOW:
				reply = Show(sanitizedInput)
			case CLEAN:
				reply = Clean()
			default:
				reply = "👌"
			}
		} else {
			reply = "🙈 Invalid command!"
		}
	}

	return reply
}

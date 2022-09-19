package bot

import (
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

func Reply(input string) string {
	var reply string
	command := common.NewString(input).
		ReplaceAll("/", "").
		ReplaceAll("_", " ").
		Value()
	sanitizedInput := SanitizeInput(command)
	commands := strings.Split(sanitizedInput, " ")

	if len(commands) > 0 {
		switch commands[0] {
		case PUSH:
			reply = Push(sanitizedInput, githubClient)
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

	return reply
}

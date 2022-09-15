package bot

import "strings"

const (
	PUSH  = "push"
	FETCH = "fetch"
	CLEAN = "clean"
	EDIT  = "edit"
	SHOW  = "show"
)

func Reply(input string) string {
	var reply string
	sanitizedInput := SanitizeInput(input)
	commands := strings.Split(sanitizedInput, " ")

	if len(commands) > 0 {
		switch commands[0] {
		case PUSH:
			reply = Push(input, githubClient)
		case EDIT:
			reply = Update(input)
		case FETCH:
			reply = Fetch()
		case SHOW:
			reply = Show(input)
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

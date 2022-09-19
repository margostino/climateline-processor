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
	var command string
	command = strings.ReplaceAll(input, "/", command)
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

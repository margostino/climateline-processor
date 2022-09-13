package bot

import "github.com/margostino/climateline-processor/common"

func SanitizeInput(input string) string {
	return common.NewString(input).
		ToLower().
		Trim(" ").
		Value()
}

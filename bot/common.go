package bot

import (
	"github.com/margostino/climateline-processor/common"
	"regexp"
)

func SanitizeInput(input string) string {
	return common.NewString(input).
		ToLower().
		Trim(" ").
		Value()
}

func IsValidInput(input string) bool {
	sanitizedInput := SanitizeInput(input)
	match, err := regexp.MatchString(`^((push ([0-9]+\s*)+)|(edit [0-9]+\n.*?\n.*?\n.*?\n(agreements|assessment|awareness|warming|wildfires|floods|drought|health))|fetch|/fetch|show|clean|/show|show [0-9]+|title [0-9]+ .*?|source [0-9]+ .*?|location [0-9]+ .*?|category [0-9]+ (agreements|assessment|awareness|warming|wildfires|floods|drought|health))$`, sanitizedInput)
	common.SilentCheck(err, "when matching input with regex")
	return match
}

func extractIds(input string, prefix string) string {
	return common.NewString(input).
		ToLower().
		TrimPrefix(prefix).
		Split(" ").
		Join(",").
		Value()
}

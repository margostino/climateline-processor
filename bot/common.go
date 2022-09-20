package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"regexp"
)

func SanitizeInput(input string) string {
	return common.NewString(input).
		ToLower().
		Trim(" ").
		Value()
}

func IsValidInput(message *tgbotapi.Message) bool {
	input := message.Text
	sanitizedInput := SanitizeInput(input)
	match, err := regexp.MatchString(`^((push ([0-9]+\s*)+)|(edit [0-9]+\n.*?\n.*?\n.*?\n(agreements|assessment|awareness|warming|wildfires|floods|drought|health|hurricane))|fetch|/fetch|/fetch_climate|/fetch_air|/fetch_effects|/fetch_drought|/fetch_floods|/fetch_greenhouse|/fetch_heatwave|/fetch_hurricanes|/fetch_wildfires|fetch climate|fetch air|fetch effects|fetch drought|fetch floods|fetch greenhouse|fetch heatwave|fetch hurricanes|fetch wildfires|/clean|clean|show|/show|show [0-9]+|title [0-9]+ .*?|source [0-9]+ .*?|location [0-9]+ .*?|category [0-9]+ (agreements|assessment|awareness|warming|wildfires|floods|drought|health))$`, sanitizedInput)
	common.SilentCheck(err, "when matching input with regex")
	return match || message.ReplyToMessage != nil
}

func extractIds(input string, prefix string) string {
	return common.NewString(input).
		ToLower().
		TrimPrefix(prefix).
		Split(" ").
		Join(",").
		Value()
}

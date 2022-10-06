package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"regexp"
	"strings"
)

var Categories = []string{
	"warming",
	"agreements",
	"assessment",
	"awareness",
	"wildfires",
	"floods",
	"drought",
	"health",
	"hurricane",
	"pollution",
}

var instructions = []string{
	`show [0-9]+`,
	`show`,
	`clean`,
	`help`,
	`push ([0-9]+\s*)+`,
	`fetch`,
	`title [0-9]+ .*?`,
	`source [0-9]+ .*?`,
	`location [0-9]+ .*?`,
	fmt.Sprintf(`category [0-9]+ (%s)`, strings.Join(Categories, "|")),
	regexBySourceConcat("fetch", " "),
	fmt.Sprintf(`edit [0-9]+\n.*?\n.*?\n.*?\n(%s)`, strings.Join(Categories, "|")),
	fmt.Sprintf(`new \n.*?\n.*?\n.*?\n.*?\n.*?\n(%s)\n.*?`, strings.Join(Categories, "|")),
}

var commands = []string{
	"/show",
	"/clean",
	"/fetch",
	"/help",
	regexBySourceConcat("/fetch", "_"),
}

var Sources = []string{
	"air",
	"climate",
	"effects",
	"drought",
	"floods",
	"greenhouse",
	"heatwave",
	"hurricanes",
	"wildfires",
}

func SanitizeInput(input string) string {
	return common.NewString(input).
		ToLower().
		Trim(" ").
		Value()
}

func IsValidInput(message *tgbotapi.Message) bool {
	input := message.Text
	sanitizedInput := SanitizeInput(input)
	regex := fmt.Sprintf(`^(%s)$`, buildInputRegex())
	match, err := regexp.MatchString(regex, sanitizedInput)
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

func buildInputRegex() string {
	var regex string
	var all = append(instructions, commands...)
	for _, pattern := range all {
		//partial := fmt.Sprintf(`(%s)`, instruction)
		if regex == "" {
			regex = pattern
		} else {
			regex += fmt.Sprintf(`|%s`, pattern)
		}
	}
	return regex
}

func regexBySourceConcat(command string, separator string) string {
	var regex string
	for _, source := range Sources {
		partial := fmt.Sprintf("%s%s%s", command, separator, source)
		if regex == "" {
			regex = partial
		} else {
			regex += fmt.Sprintf("|%s", partial)
		}

	}
	return regex
}

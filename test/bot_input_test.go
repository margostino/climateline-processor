package test

import (
	"fmt"
	"github.com/margostino/climateline-processor/bot"
	"testing"
)

func TestValidInputs(t *testing.T) {
	validateInput(t, "location 1 Sweden")
	validateInput(t, "title 1 some title")
	validateInput(t, "source 1 some source")
	validateInput(t, "push 1")
	validateInput(t, "show 1")
	validateInput(t, "clean")
	validateInput(t, "fetch")
	validateInput(t, "/fetch")
	validateInput(t, "/show")
	validateInput(t, "/clean")

	for _, category := range bot.Categories {
		validateInput(t, fmt.Sprintf("category 1 %s", category))
		validateInput(t, fmt.Sprintf("edit 1\nbreaking news\nwikipedia\nSweden\n%s", category))
		validateInput(t, fmt.Sprintf("new \n2022-09-10\nbreaking news\nwikipedia\nSweden\nsomelink.com\n%s\nThis is a new content", category))
	}

	for _, source := range bot.Sources {
		validateInput(t, fmt.Sprintf("fetch %s", source))
		validateInput(t, fmt.Sprintf("/fetch_%s", source))
	}

}

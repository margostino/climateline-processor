package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/cache"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"os"
	"strings"
)

func Update(input string) string {
	var reply string
	var id string
	var update *domain.Update

	sanitizedInput := SanitizeInput(input)

	if strings.Contains(sanitizedInput, "edit") {
		instructions := strings.Split(input, "\n")
		update = &domain.Update{
			Title:      instructions[1],
			SourceName: instructions[2],
			Location:   instructions[3],
			Category:   instructions[4],
		}
		id = extractIds(instructions[0], "edit ")
	} else {
		params := strings.Split(sanitizedInput, " ")
		property := params[0]
		id = params[1]
		value := common.NewString(input).
			TrimIndex(2).
			Value()

		if property == "category" {
			update = &domain.Update{
				Category: value,
			}
		} else if property == "location" {
			update = &domain.Update{
				Location: value,
			}
		} else if property == "title" {
			update = &domain.Update{
				Title: value,
			}
		} else {
			update = &domain.Update{
				SourceName: value,
			}
		}
	}

	if updateCachedItems(id, update) {
		reply = "✅ article updated"
	} else {
		reply = "🔴 article update failed!"
	}

	return reply
}

func ShouldUpdate(input string) bool {
	sanitizedInput := SanitizeInput(input)
	return strings.Contains(sanitizedInput, "category") ||
		strings.Contains(sanitizedInput, "title") ||
		strings.Contains(sanitizedInput, "source") ||
		strings.Contains(sanitizedInput, "location") ||
		strings.Contains(sanitizedInput, "edit")
}

func updateCachedItems(id string, edit *domain.Update) bool {
	client := &http.Client{}
	url := fmt.Sprintf("%s?id=%s", cache.GetBaseCacheUrl(), id)
	json, err := json.Marshal(edit)

	if !common.IsError(err, "when marshaling edit data") {
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLIMATELINE_JOB_SECRET")))
		response, err := client.Do(request)
		common.SilentCheck(err, "when updating cached item")
		return response.StatusCode == 204
	}
	return false
}

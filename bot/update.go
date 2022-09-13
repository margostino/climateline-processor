package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/api"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"net/http"
	"os"
	"strings"
)

func Edit(input string) string {
	var reply string
	var edit *domain.Edit
	sanitizedInput := SanitizeInput(input)
	params := strings.Split(sanitizedInput, " ")
	property := params[0]
	id := params[1]
	value := common.NewString(input).
		TrimIndex(2).
		Value()

	if property == "category" {
		edit = &domain.Edit{
			Category: value,
		}
	} else if property == "location" {
		edit = &domain.Edit{
			Location: value,
		}
	} else if property == "title" {
		edit = &domain.Edit{
			Title: value,
		}
	} else {
		edit = &domain.Edit{
			SourceName: value,
		}
	}
	if updateCachedItems(id, edit) {
		reply = "✅ article updated"
	} else {
		reply = "🔴 article update failed!"
	}

	return reply
}

func updateCachedItems(id string, edit *domain.Edit) bool {
	client := &http.Client{}
	url := fmt.Sprintf("%s?id=%s", api.GetBaseCacheUrl(), id)
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

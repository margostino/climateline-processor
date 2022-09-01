package api

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Response struct {
	Msg    string `json:"text"`
	ChatID int64  `json:"chat_id"`
	Method string `json:"method"`
}

func Bot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reply string
	w.Header().Add("Content-Type", "application/json")

	//log.Printf("Cached Items (Reply): %d", len(cache.Items))

	body, _ := ioutil.ReadAll(r.Body)
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Fatal("Error updating →", err)
	}

	log.Printf("[%s@%d] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

	input := sanitizeInput(update.Message)
	if isValidInput(input) {
		if shouldUpload(input) {
			id := extractId(input)
			item := getCachedItem(id)
			reply = item.Link
			// TODO: commit git new article
		} else {
			reply = "👌"
		}
	} else {
		reply = "Input is not valid"
		log.Println(reply)
	}

	data := Response{
		Msg:    reply,
		Method: "sendMessage",
		ChatID: update.Message.Chat.ID,
	}

	message, _ := json.Marshal(data)
	log.Printf("Response %s", string(message))
	fmt.Fprintf(w, string(message))

}

func sanitizeInput(message *tgbotapi.Message) string {
	return strings.Trim(strings.ToLower(message.Text), " ")
}

func isValidInput(input string) bool {
	match, err := regexp.MatchString(`^((yes [0-9]*)|(no))$`, input)
	common.SilentCheck(err, "when matching input with regex")
	return match
}

func shouldUpload(input string) bool {
	return strings.Contains(input, "yes")
}

func extractId(input string) int {
	id, err := strconv.Atoi(strings.Split(input, " ")[1])
	common.SilentCheck(err, "must not reach this state")
	return id
}

func getCachedItem(id int) domain.Item {
	var item domain.Item
	url := fmt.Sprintf("%s?id=%d", baseCacheUrl, id)
	resp, err := http.Get(url)
	common.SilentCheck(err, "when getting cached item")
	err = json.NewDecoder(resp.Body).Decode(&item)
	common.SilentCheck(err, "when decoding response from cache")
	return item
}

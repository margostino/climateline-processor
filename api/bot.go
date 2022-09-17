package api

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Msg    string `json:"text"`
	ChatID int64  `json:"chat_id"`
	Method string `json:"method"`
}

var botApi *tgbotapi.BotAPI

func Bot(w http.ResponseWriter, r *http.Request) {
	var reply string

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Printf("PATH: " + path)

	files, err := ioutil.ReadDir("/var/task")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Println(file.Name(), file.IsDir())
	}

	contents, err := os.ReadFile("./config/config.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		reply = "Error: " + err.Error()
		return
	} else {
		reply = "Content: " + string(contents)
	}
	body, _ := ioutil.ReadAll(r.Body)
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Fatal("Error updating →", err)
	}
	data := Response{
		Msg:    reply,
		Method: "sendMessage",
		ChatID: update.Message.Chat.ID,
	}
	jsonResp, err := json.Marshal(data)
	w.Write(jsonResp)
	return
	//log.Printf("Method: %s "+
	//	"Proto: %s "+
	//	"User-Agent: %s, "+
	//	"Host: %s, "+
	//	"RequestURI: %s, "+
	//	"RemoteAddr: %s",
	//	r.Method, r.Proto, r.Header.Get("User-Agent"), r.Host, r.RequestURI, r.RemoteAddr)
	//
	//body, _ := ioutil.ReadAll(r.Body)
	//var update tgbotapi.Update
	//if err := json.Unmarshal(body, &update); err != nil {
	//	log.Fatal("Error updating →", err)
	//}
	//
	//log.Printf("[%s@%d] %s", update.Message.From.UserName, update.Message.Chat.ID, "update.Message.Text")
	//
	//if security.IsAdmin(r) {
	//	w.Header().Add("Content-Type", "application/json")
	//	input := update.Message.Text
	//	if bot.IsValidInput(input) {
	//		reply = bot.Reply(input)
	//	} else {
	//		reply = "Input is not valid"
	//		log.Println(reply)
	//	}
	//} else {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	reply = "Unauthorized to handle Bot communication"
	//	log.Printf(reply)
	//}
	//
	//data := Response{
	//	Msg:    reply,
	//	Method: "sendMessage",
	//	ChatID: update.Message.Chat.ID,
	//}
	//
	//message, _ := json.Marshal(data)
	//log.Printf("Response %s", string(message))
	//fmt.Fprintf(w, string(message))
}

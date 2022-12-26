package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/climateline-processor/common"
	"github.com/margostino/climateline-processor/domain"
	"log"
	"os"
	"strconv"
)

var botApi, _ = newBot()

func NotifyBot(item *domain.Item) {
	message := fmt.Sprintf("🔔 New article! \n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n"+
		"%s %s\n", //"%s <a href='%s'>Here</a>\n",
		domain.ID_PREFIX, item.Id,
		domain.DATE_PREFIX, item.Timestamp,
		domain.TITLE_PREFIX, item.Title,
		domain.SOURCE_PREFIX, item.SourceName,
		domain.CONTENT_PREFIX, item.Content,
		domain.LINK_PREFIX, item.Link)
	send(message)
}

func send(message string) {
	if botApi != nil {
		userId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ADMIN_USER"), 10, 64)
		msg := tgbotapi.NewMessage(userId, message)
		msg.ReplyMarkup = nil
		msg.ParseMode = "HTML"
		botApi.Send(msg)
	} else {
		log.Printf("Bot initialization failed")
	}
}

func newBot() (*tgbotapi.BotAPI, error) {
	client, error := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	//bot.Debug = true
	common.SilentCheck(error, "when creating a new BotAPI instance")
	//log.Printf("Authorized on account %s\n", client.Self.UserName)
	return client, error
}

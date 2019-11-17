package bot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"os"
)

type botTelegram struct {
	api *tgbotapi.BotAPI
}

//SendMessage - send message in chat telegram
func SendMessage(chatID int64, msg string) {
	var bot = connect(os.Getenv("key_bot"))
	_, err := bot.api.Send(tgbotapi.NewMessage(chatID, msg))
	if err != nil {
		log.Fatal(err)
	}
}

func connect(token string) (b botTelegram) {
	var err error
	b.api, err = tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", b.api.Self.UserName)
	return b
}

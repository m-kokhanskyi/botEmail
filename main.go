package main

import (
	"botEmail/bot"
	"botEmail/email"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	loadEnv()

	messages := email.GetNewMessages()

	for _, m := range messages {
		idChat, _ := strconv.ParseInt(os.Getenv("id_telegram_chat"), 10, 64)
		bot.SendMessage(idChat, m)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

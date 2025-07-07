package main

import (
	"github.com/murad755/telegram-bot-lyrics/bot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not set in environment")
	}

	bot.Start(token)
}

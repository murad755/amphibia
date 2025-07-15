package main

import (
	"github.com/murad755/telegram-bot-lyrics/bot"
	"github.com/murad755/telegram-bot-lyrics/lyrics"
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
	baseURL := os.Getenv("LYRICS_API_URL")
	if baseURL == "" {
		log.Fatal("LYRICS_API_URL not set in environment")
	}

	lyricsClient := lyrics.NewURL(baseURL)

	bot.Start(token, lyricsClient)

}

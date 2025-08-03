package main

import (
	"log"
	"os"

	"github.com/murad755/amphibia/bot"
	"github.com/murad755/amphibia/lyrics"

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

	err := bot.Start(token, lyricsClient)
	if err != nil {
		log.Fatal(err)
	}
}

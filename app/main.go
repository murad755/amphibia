package main

import (
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not set")
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	menu := &tele.ReplyMarkup{}

	btnBohemian := menu.Data("Bohemian Rhapsody", "song_boh_rhap")
	btnImagine := menu.Data("Imagine", "song_imagine")
	btnBillie := menu.Data("Billie Jean", "song_billie")

	menu.Inline(
		menu.Row(btnBohemian),
		menu.Row(btnImagine),
		menu.Row(btnBillie),
	)

	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("👋 Welcome! Use /getlyrics to choose a song.")
	})

	bot.Handle("/getlyrics", func(c tele.Context) error {
		return c.Send("🎵 Which song's lyrics do you want to hear?", menu)
	})

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		print(c.Callback().Data)
		return nil
		//songId, err := strconv.Atoi(strings.TrimSpace(callbackData))
		//if err != nil {
		//	bh.Logger.Println(err.Error())
		//	return
		//}
	})

	log.Println("✅ Bot is running")
	bot.Start()
}

package bot

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Start(token string) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	RegisterHandlers(bot)

	log.Println("✅ Bot is running")
	bot.Start()
}

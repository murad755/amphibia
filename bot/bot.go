package bot

import (
	"github.com/murad755/telegram-bot-lyrics/lyrics"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Start(token string, lyricsClient *lyrics.Client) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	RegisterHandlers(bot, lyricsClient)

	bot.Start()
}

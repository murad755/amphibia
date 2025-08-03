package bot

import (
	"fmt"
	"github.com/murad755/amphibia/lyrics"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Start(token string, lyricsClient *lyrics.Client) error {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	h := NewHandler(bot, lyricsClient)
	h.bot.Start()

	return nil
}

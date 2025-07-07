package bot

import (
	"github.com/murad755/telegram-bot-lyrics/lyrics"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func RegisterHandlers(bot *tele.Bot) {
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("👋 Welcome! Type song name to get the song.")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		songName := strings.TrimSpace(c.Text())
		if songName == "" || strings.HasPrefix(songName, "/") {
			return nil
		}

		resp, err := lyrics.ListLyrics(songName)
		if err != nil {
			return c.Send("❌ Error fetching lyrics list")
		}

		if len(resp.Messages.Songlist) == 0 {
			return c.Send("😢 No songs found.")
		}

		menu := &tele.ReplyMarkup{}
		rows := make([]tele.Row, 0, len(resp.Messages.Songlist))
		for _, song := range resp.Messages.Songlist {
			rows = append(rows, menu.Row(menu.Data(song.Title, strconv.Itoa(song.ID))))
		}
		menu.Inline(rows...)

		return c.Send("🎵 Select a song from below:", menu)
	})

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		id := strings.TrimSpace(c.Callback().Data)
		resp, err := lyrics.GetLyrics(id)
		if err != nil {
			return c.Send("❌ Error fetching lyrics")
		}

		chunks := lyrics.ChunkString(resp.Messages.Lyrics, 4096)
		for _, part := range chunks {
			if err := c.Send(part); err != nil {
				return err
			}
		}
		return nil
	})
}

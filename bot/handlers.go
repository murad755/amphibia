package bot

import (
	"errors"
	"fmt"
	"github.com/murad755/amphibia/amphibia"
	"github.com/murad755/amphibia/lyrics"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (h *Handler) handleStart(c tele.Context) error {
	return c.Send("👋 Welcome! Type song name to get the song.")
}

func (h *Handler) handleText(c tele.Context) error {
	songName := strings.TrimSpace(c.Text())

	// Lucky search for getting first matched song
	if strings.HasSuffix(songName, "!") {
		cleanName := strings.ReplaceAll(songName, "!", "")

		firstLyrics, err := h.svc.FindFirstLyrics(strings.TrimSpace(cleanName))
		switch {
		case err == nil:
			return h.sendSplitText(c, firstLyrics)
		case errors.Is(err, amphibia.ErrLyricsUnavailable):
			return c.Send("Lyrics is not available")
		case errors.Is(err, amphibia.ErrNoSongsFound):
			return c.Send("Song not found")
		default:
			return c.Send("Unknown error occured")
		}
	}

	songList, err := h.svc.SearchLyrics(songName)
	switch {
	case errors.Is(err, amphibia.ErrLyricsUnavailable):
		return c.Send("Lyrics is not available")
	case errors.Is(err, amphibia.ErrNoSongsFound):
		return c.Send("Song not found")
	case err != nil:
		return c.Send("Unknown error occured")
	}

	menu := &tele.ReplyMarkup{}
	rows := make([]tele.Row, 0, len(songList))
	for _, song := range songList {
		rows = append(rows, menu.Row(menu.Data(song.Title, song.ID)))
	}
	menu.Inline(rows...)

	return c.Send("🎵 Select a song from below:", menu)
}

func (h *Handler) sendSplitText(c tele.Context, lyricsText string) error {
	chunks := lyrics.ChunkString(lyricsText, 4096)
	for _, part := range chunks {
		if err := c.Send(part); err != nil {
			return fmt.Errorf("sending lyrics: %w", err)
		}
	}

	return nil
}

func (h *Handler) handleCallback(c tele.Context) error {
	id := strings.TrimSpace(c.Callback().Data)

	l, err := h.svc.LyricsByID(id)
	switch {
	case err == nil:
		return h.sendSplitText(c, l)
	case errors.Is(err, amphibia.ErrLyricsUnavailable):
		return c.Send("Lyrics is not available")
	case errors.Is(err, amphibia.ErrNoSongsFound):
		return c.Send("Song not found")
	default:
		return c.Send("Unknown error occured")
	}
}

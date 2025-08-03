package bot

import (
	"fmt"
	"github.com/murad755/amphibia/lyrics"
	"log"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	bot          *tele.Bot
	lyricsClient *lyrics.Client
}

func NewHandler(bot *tele.Bot, client *lyrics.Client) *Handler {
	h := &Handler{bot: bot, lyricsClient: client}
	h.register()
	return h
}

func (h *Handler) register() {
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle(tele.OnText, h.handleText)
	h.bot.Handle(tele.OnCallback, h.handleCallback)
}

func (h *Handler) handleStart(c tele.Context) error {
	return c.Send("👋 Welcome! Type song name to get the song.")
}

func (h *Handler) handleText(c tele.Context) error {
	songName := strings.TrimSpace(c.Text())
	if songName == "" || strings.HasPrefix(songName, "/") {
		return ErrEmptySongName
	}

	resp, err := h.lyricsClient.ListLyrics(songName)
	if err != nil {
		log.Printf("Error getting lyrics: %v", err)
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
}

func (h *Handler) handleCallback(c tele.Context) error {
	id := strings.TrimSpace(c.Callback().Data)

	resp, err := h.lyricsClient.GetLyrics(id)
	if err != nil {
		log.Printf("Error getting lyrics: %v", err)
		return c.Send("❌ Error fetching lyrics")
	}

	lyricsText := strings.TrimSpace(resp.Messages.Lyrics)
	if lyricsText == "" {
		return c.Send("Sorry, no lyrics found for this song.")
	}

	chunks := lyrics.ChunkString(lyricsText, 4096)
	for _, part := range chunks {
		if err = c.Send(part); err != nil {
			return fmt.Errorf("sending lyrics: %w", err)
		}
	}

	return nil
}

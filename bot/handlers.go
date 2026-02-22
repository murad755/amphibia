package bot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/murad755/amphibia/amphibia"
	"github.com/murad755/amphibia/lyrics"

	tele "gopkg.in/telebot.v4"
)

const symfoniumTag = "#Symfonium"

var symfoniumPrefixes = []string{
	"Check out: ",              // e.g. "Check out: Bad Decisions • Girlhood #Symfonium"
	"I'm currently listening to: ", // e.g. "I'm currently listening to: Ariana Grande - 7 rings #Symfonium"
}

// parseListeningCaption extracts search query and lucky flag from a Symfonium-style caption.
// Supports "Check out: <query> #Symfonium" or "I'm currently listening to: <query> #Symfonium" with optional trailing "!" for lucky search.
// Returns query (trimmed), lucky (true if caption ends with !), and ok (true when prefix matches and query non-empty).
func parseListeningCaption(caption string) (query string, lucky bool, ok bool) {
	caption = strings.TrimSpace(caption)
	if caption == "" {
		return "", false, false
	}
	var rest string
	for _, prefix := range symfoniumPrefixes {
		if strings.HasPrefix(caption, prefix) {
			rest = strings.TrimSpace(caption[len(prefix):])
			break
		}
	}
	if rest == "" {
		return "", false, false
	}
	// Optional trailing #Symfonium
	rest = strings.TrimSpace(strings.TrimSuffix(rest, symfoniumTag))
	rest = strings.TrimSpace(rest)
	// Optional trailing ! for lucky search
	if strings.HasSuffix(rest, "!") {
		lucky = true
		rest = strings.TrimSpace(strings.TrimSuffix(rest, "!"))
	}
	if rest == "" {
		return "", false, false
	}
	return rest, lucky, true
}

// escapeBackticks replaces '`' with '\`' so caption can be wrapped in backticks for Markdown.
func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "\\`")
}

func (h *Handler) handleStart(c tele.Context) error {
	return c.Send("BEEP BOOP 🤖. Type song name to get the song. If you are feeling lucky, add ! in the end 👾.")
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

func (h *Handler) handlePhoto(c tele.Context) error {
	caption := c.Text()
	query, lucky, ok := parseListeningCaption(caption)
	if !ok {
		return nil
	}

	msg := c.Message()
	if msg.Photo == nil {
		return nil
	}
	// Shareable caption is the stripped query (artist + song only).
	shareableCaption := "`" + escapeBackticks(query) + "`"
	photo := &tele.Photo{File: tele.File{FileID: msg.Photo.FileID}, Caption: shareableCaption}
	if err := c.Send(photo, &tele.SendOptions{ParseMode: tele.ModeMarkdown}); err != nil {
		return err
	}

	// Search query: remove • and collapse spaces for lyrics API.
	searchQuery := strings.Join(strings.Fields(strings.ReplaceAll(query, "•", " ")), " ")

	// Run search using caption-derived query (same as handleText).
	if lucky {
		firstLyrics, err := h.svc.FindFirstLyrics(searchQuery)
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

	songList, err := h.svc.SearchLyrics(searchQuery)
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

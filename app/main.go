package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type ListLyricsResp struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Query    string   `json:"query"`
	Messages struct {
		Songlist []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"songlist"`
	} `json:"messages"`
}
type GetLyricsResp struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Query    string   `json:"query"`
	Messages struct {
		Lyrics string `json:"lyrics"`
	} `json:"messages"`
}

func listLyrics(query string) (*ListLyricsResp, error) {
	baseURL := "http://localhost:5001/api/v1/find-songs/?query="
	escapedQuery := url.QueryEscape(query)

	response, err := http.Get(baseURL + escapedQuery)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var resp ListLyricsResp

	err = decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func getLyrics(id string) (*GetLyricsResp, error) {
	baseURL := "http://localhost:5001/api/v1/song/"

	response, err := http.Get(baseURL + id)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var resp GetLyricsResp

	err = decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func chunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
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

	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("👋 Welcome! Use /getlyrics to choose a song.")
	})

	bot.Handle("/getlyrics", func(c tele.Context) error {
		songName := c.Message().Payload
		if songName == "" {
			return c.Send("❗ Please provide a song name. Example:\n/getlyrics Not Like Us")
		}
		lyrics, err := listLyrics(songName)
		if err != nil {
			return err
		}
		menu := &tele.ReplyMarkup{}
		rows := make([]tele.Row, 0, len(lyrics.Messages.Songlist))

		for _, song := range lyrics.Messages.Songlist {
			rows = append(rows, menu.Row(menu.Data(song.Title, strconv.Itoa(song.ID))))
		}
		menu.Inline(rows...)

		return c.Send("Test", menu)
	})

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		callbackData := c.Callback().Data
		songId := strings.TrimSpace(callbackData)

		resp, err := getLyrics(songId)
		if err != nil {
			return err
		}

		messageSplitted := chunkString(resp.Messages.Lyrics, 4096)
		for _, m := range messageSplitted {
			err = c.Send(m)
			if err != nil {
				return err
			}
		}

		return nil
	})

	log.Println("✅ Bot is running")
	bot.Start()
}

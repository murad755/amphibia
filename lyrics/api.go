package lyrics

import (
	"encoding/json"
	"github.com/murad755/telegram-bot-lyrics/models"
	"net/http"
	"net/url"
)

func ListLyrics(query string) (*models.ListLyricsResp, error) {
	baseURL := "http://localhost:5001/api/v1/find-songs/?query="
	escapedQuery := url.QueryEscape(query)

	response, err := http.Get(baseURL + escapedQuery)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var resp models.ListLyricsResp

	err = decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetLyrics(id string) (*models.GetLyricsResp, error) {
	baseURL := "http://localhost:5001/api/v1/song/"

	response, err := http.Get(baseURL + id)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var resp models.GetLyricsResp

	err = decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func ChunkString(s string, chunkSize int) []string {
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

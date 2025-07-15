package lyrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type Client struct {
	baseURL string
}

func NewURL(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) ListLyrics(query string) (*ListLyricsResp, error) {
	escapedQuery := url.QueryEscape(query)
	fullURL := c.baseURL + "/find-songs/?query=" + escapedQuery

	response, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var resp ListLyricsResp
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetLyrics(id string) (*GetLyricsResp, error) {
	fullURL := c.baseURL + "/song/" + id + "/"

	fmt.Printf("song ID %s", fullURL)
	response, err := http.Get(fullURL)
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

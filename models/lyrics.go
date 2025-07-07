package models

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

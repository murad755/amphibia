package amphibia

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Song struct {
	ID    string
	Title string
}

var (
	ErrEmptySongName     = errors.New("empty song name")
	ErrLyricsUnavailable = errors.New("lyrics unavailable")
	ErrNoSongsFound      = errors.New("no songs found")
)

func (s *Service) SearchLyrics(songName string) ([]Song, error) {
	if songName == "" || strings.HasPrefix(songName, "/") {
		return nil, ErrEmptySongName
	}

	resp, err := s.lyricsClient.ListLyrics(songName)
	if err != nil {
		log.Printf("Error getting lyrics: %v", err)
		return nil, fmt.Errorf("fetching lyrics list: %w", ErrLyricsUnavailable)
	}

	if len(resp.Messages.Songlist) == 0 {
		return nil, ErrNoSongsFound
	}

	songList := make([]Song, 0, len(resp.Messages.Songlist))
	for _, song := range resp.Messages.Songlist {
		songList = append(songList, Song{
			ID:    strconv.Itoa(song.ID),
			Title: song.Title,
		})
	}

	return songList, nil
}

func (s *Service) LyricsByID(id string) (string, error) {
	resp, err := s.lyricsClient.GetLyrics(id)
	if err != nil {
		log.Printf("Error getting lyrics: %v", err)
		return "", fmt.Errorf("fetching lyrics list: %w", ErrLyricsUnavailable)
	}

	lyricsText := strings.TrimSpace(resp.Messages.Lyrics)
	if lyricsText == "" {
		return "", ErrNoSongsFound
	}
	return lyricsText, nil
}

func (s *Service) FindFirstLyrics(songName string) (string, error) {
	resp, err := s.lyricsClient.ListLyrics(songName)
	if err != nil {
		log.Printf("Error listing lyrics (direct): %v", err)
		return "", fmt.Errorf("fetching lyrics list: %w", ErrLyricsUnavailable)
	}

	if len(resp.Messages.Songlist) == 0 {
		return "", ErrNoSongsFound
	}

	firstSong := resp.Messages.Songlist[0]
	return s.LyricsByID(strconv.Itoa(firstSong.ID))
}

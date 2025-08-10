package amphibia

import (
	"github.com/murad755/amphibia/lyrics"
)

type Service struct {
	lyricsClient *lyrics.Client
}

func NewService(lyricsClient *lyrics.Client) *Service {
	return &Service{lyricsClient: lyricsClient}
}

package bot

import (
	"fmt"
	"github.com/murad755/amphibia/amphibia"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	bot *tele.Bot
	svc *amphibia.Service
}

func NewHandler(bot *tele.Bot, client *amphibia.Service) *Handler {
	h := &Handler{bot: bot, svc: client}
	h.register()
	return h
}

func (h *Handler) register() {
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle(tele.OnText, h.handleText)
	h.bot.Handle(tele.OnCallback, h.handleCallback)
}

func Start(token string, svc *amphibia.Service) error {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	h := NewHandler(bot, svc)
	h.bot.Start()

	return nil
}

package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	domainbot "github.com/dvdxa/tg-notifier/internal/domain/bot"

	"github.com/dvdxa/tg-notifier/internal/usecase/bot"
)

type Handler struct {
	bot *tgbotapi.BotAPI
	u   *bot.Usecase
	log zerolog.Logger
}

func NewHandler(bot *tgbotapi.BotAPI, u *bot.Usecase, log zerolog.Logger) *Handler {
	return &Handler{
		bot: bot,
		u:   u,
		log: log.With().Str("pkg", "telegram").Logger(),
	}
}

func (h *Handler) Run(ctx context.Context) (err error) {
	h.log.Info().Msg("starting telegram bot...")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := h.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			var msg *tgbotapi.Message

			if update.Message != nil {
				msg = update.Message
			} else if update.ChannelPost != nil {
				msg = update.ChannelPost
			} else {
				continue
			}

			cmd := domainbot.Command(msg.Command())

			var res domainbot.Result
			res, err = h.u.HandleCommand(ctx, cmd)
			if err != nil {
				// TODO: do we need to send errs to chan?
				log.Err(err).Msg("failed to handle command")
				continue
			}

			text := fmt.Sprintf("%s (code %d)", res.Message, res.StatusCode)
			tgMsg := tgbotapi.NewMessage(msg.Chat.ID, text)

			_, err = h.bot.Send(tgMsg)
			if err != nil {
				log.Err(err).Msg("failed to send message")
			}
		}
	}
}

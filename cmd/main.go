package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/dvdxa/tg-notifier/internal/adapter/externalapi"
	"github.com/dvdxa/tg-notifier/internal/config"
	"github.com/dvdxa/tg-notifier/internal/transport/telegram"
	"github.com/dvdxa/tg-notifier/internal/usecase/bot"
)

func main() {
	log.Logger = config.Values().Logger()

	botAPI, err := tgbotapi.NewBotAPI(config.Values().Telegram.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create telegram bot")
	}

	sender := externalapi.NewSender(config.Values().ExternalApi.Addr)
	botUsecase := bot.NewUsecase(log.Logger, sender)
	tgHandler := telegram.NewHandler(botAPI, botUsecase, log.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		return tgHandler.Run(gctx)
	})

	g.Go(func() (err error) {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			log.Info().Str("source", "system.quit").Msgf("received signal: %v, shutting down...\n", sig)
			cancel()
			return nil
		}
	})

	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Err(err).Str("source", "g.Wait()").Msg("errgroup finished with error")
	}

	time.Sleep(2 * time.Second)
}

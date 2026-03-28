package telegram

import (
	"context"

	"github.com/dvdxa/tg-notifier/internal/domain/bot"
)

type botUsecase interface {
	HandleCommand(ctx context.Context, cmd bot.Command) (res bot.Result, err error)
}

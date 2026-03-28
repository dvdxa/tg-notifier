package bot

import (
	"context"

	"github.com/dvdxa/tg-notifier/internal/domain/bot"
)

type ExternalSender interface {
	Send(ctx context.Context, cmd bot.Command) (statusCode int, err error)
}

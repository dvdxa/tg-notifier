package bot

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/dvdxa/tg-notifier/internal/domain/bot"
)

type Usecase struct {
	sender ExternalSender
	log    zerolog.Logger
}

func NewUsecase(log zerolog.Logger, sender ExternalSender) *Usecase {
	return &Usecase{
		sender: sender,
		log:    log,
	}
}

func (s *Usecase) HandleCommand(ctx context.Context, cmd bot.Command) (res bot.Result, err error) {
	switch cmd {
	case bot.CommandStart, bot.CommandStop:
		var statusCode int
		statusCode, err = s.sender.Send(ctx, cmd)
		if err != nil {
			return
		}

		if statusCode == http.StatusOK {
			return bot.Result{
				StatusCode: statusCode,
				Message:    "Command successfully delivered to External API",
			}, nil
		}

		return bot.Result{
			StatusCode: statusCode,
			Message:    "External API returned an error",
		}, nil

	default:
		return res, errors.New("unknown command")
	}
}

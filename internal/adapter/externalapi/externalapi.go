package externalapi

import (
	"context"
	"net/http"
	"time"

	"github.com/dvdxa/tg-notifier/internal/domain/bot"
)

type ExternalAPI struct {
	client *http.Client
	addr   string
}

func NewSender(addr string) *ExternalAPI {
	return &ExternalAPI{
		addr: addr,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (e *ExternalAPI) Send(ctx context.Context, cmd bot.Command) (statusCode int, err error) {
	// TODO: implement me
	return http.StatusOK, nil
}

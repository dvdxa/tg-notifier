package config

import (
	"io"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const formatJSON = "json"

type config struct {
	Log         logger
	ExternalApi externalApi
	Telegram    telegram
}

type logger struct {
	LogLevel   string `envconfig:"LOGGER_LEVEL"   default:"debug"`
	LogFormat  string `envconfig:"LOGGER_FORMAT"  default:"console"`
	LogNoColor bool   `envconfig:"LOGGER_NOCOLOR" default:"true"`
}

type externalApi struct {
	Addr string `envconfig:"EXTERNAL_ADDR" required:"true"` // TODO: rename
}

type telegram struct {
	BotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
}

var configuration *config

func internalConfig() *config {
	if configuration == nil {

		configuration = &config{}

		if err := envconfig.Process("", configuration); err != nil {
			panic(err)
		}
	}
	return configuration
}

func Values() config {
	return *internalConfig()
}

func (cfg config) Logger() (logger zerolog.Logger) { // TODO: mv to sep pkg
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.Log.LogLevel); err == nil {
		level = newLevel
	}
	var out io.Writer = os.Stdout
	if cfg.Log.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.StampMicro,
			NoColor:    cfg.Log.LogNoColor,
		}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Timestamp().Logger()
}

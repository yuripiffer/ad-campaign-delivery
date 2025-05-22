package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func Init() zerolog.Logger {
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}

package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func Init() zerolog.Logger {
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}

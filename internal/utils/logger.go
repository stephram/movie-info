package utils

import (
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var _log zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlog.Print("Using default Zerolog instance")
	zlog.Logger = zlog.With().Caller().Logger()
	zlog.Printf("Logging with Caller")

	_log = zlog.Logger
	// _log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	_log.Printf("Created _log instance")
}

func GetLogger() zerolog.Logger {
	return _log
}

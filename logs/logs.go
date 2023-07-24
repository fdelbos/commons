package logs

import (
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	once sync.Once
)

func Start() {
	once.Do(func() {
		logger := log.Logger
		if viper.GetBool("log_caller") {
			logger = log.With().Caller().Logger()
		}

		if !viper.GetBool("log_json") {
			logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		}

		switch viper.GetString("log_level") {
		case "info":
			logger = logger.Level(zerolog.InfoLevel)
		case "warn":
			logger = logger.Level(zerolog.WarnLevel)
		case "error":
			logger = logger.Level(zerolog.ErrorLevel)
		}

		log.Logger = logger
	})
}

func Logger(moduleName string, withCaller bool) zerolog.Logger {
	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.WarnLevel)

	if withCaller {
		logger = logger.With().Caller().Logger()
	}

	if !viper.GetBool("log_json") {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	switch viper.GetString("log_level") {
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	}

	return logger.With().Str("module", moduleName).Logger()
}

func CallerMarshalFunc(pc uintptr, file string, line int) string {
	return file + ":" + strconv.Itoa(line)
}

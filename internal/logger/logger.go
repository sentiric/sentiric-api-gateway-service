// sentiric-api-gateway-service/internal/logger/logger.go
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DÜZELTME: Fonksiyon imzası diğer servislerle aynı olacak şekilde güncellendi.
func New(serviceName, env, logLevel string) zerolog.Logger {
	var logger zerolog.Logger

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
		log.Warn().Msgf("Geçersiz LOG_LEVEL '%s', varsayılan olarak 'info' kullanılıyor.", logLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339

	if env == "development" {
		// Geliştirme ortamı için renkli, okunabilir konsol logları
		output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		logger = log.Output(output).With().Timestamp().Str("service", serviceName).Logger()
	} else {
		// Üretim ortamı için yapılandırılmış JSON logları
		logger = zerolog.New(os.Stderr).With().Timestamp().Str("service", serviceName).Logger()
	}

	return logger.Level(level)
}
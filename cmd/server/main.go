// sentiric-api-gateway-service/cmd/server/main.go
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/sentiric/sentiric-api-gateway-service/internal/gateway"
	"github.com/sentiric/sentiric-api-gateway-service/internal/logger"
)

var (
	ServiceVersion string
	GitCommit      string
	BuildDate      string
)

const serviceName = "api-gateway-service"

func main() {
	godotenv.Load()

	cfg, err := gateway.LoadConfig()
	if err != nil {
		// DÜZELTME: Logger henüz başlatılmadığı için, standart bir zerolog
		// logger oluşturup hatayı basıyoruz.
		log := zerolog.New(os.Stderr).With().Timestamp().Logger()
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log := logger.New(serviceName, cfg.Env, cfg.LogLevel)

	log.Info().
		Str("version", ServiceVersion).
		Str("commit", GitCommit).
		Str("build_date", BuildDate).
		Str("profile", cfg.Env).
		Msg("🚀 Starting Sentiric API Gateway Service")

	if err := gateway.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("Gateway failed to run")
	}
}
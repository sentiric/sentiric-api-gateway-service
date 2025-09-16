package main

import (
	"github.com/joho/godotenv"
	"github.com/sentiric/sentiric-api-gateway-service/internal/gateway"
	"github.com/sentiric/sentiric-api-gateway-service/internal/logger"
)

// YENİ: ldflags ile doldurulacak değişkenler
var (
	ServiceVersion string
	GitCommit      string
	BuildDate      string
)

const serviceName = "api-gateway-service"

func main() {
	godotenv.Load()
	log := logger.New(serviceName)

	// YENİ: Başlangıçta versiyon bilgisini logla
	log.Info().
		Str("version", ServiceVersion).
		Str("commit", GitCommit).
		Str("build_date", BuildDate).
		Msg("🚀 Starting Sentiric API Gateway Service")

	cfg, err := gateway.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	if err := gateway.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("Gateway failed to run")
	}
}
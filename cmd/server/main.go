package main

import (
	"github.com/joho/godotenv"
	"github.com/sentiric/sentiric-api-gateway-service/internal/gateway"
	"github.com/sentiric/sentiric-api-gateway-service/internal/logger"
)

const serviceName = "api-gateway-service"

func main() {
	godotenv.Load()
	log := logger.New(serviceName)

	log.Info().Msg("Starting Sentiric API Gateway Service")

	cfg, err := gateway.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	if err := gateway.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("Gateway failed to run")
	}
}

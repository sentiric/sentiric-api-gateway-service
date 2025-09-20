// sentiric-api-gateway-service/cmd/server/main.go
package main

import (
	"github.com/joho/godotenv"
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

	// DÜZELTME: Config yüklemesi log'dan önce yapılır.
	cfg, err := gateway.LoadConfig()
	if err != nil {
		// Log nesnesi henüz yok, bu yüzden standart log kullan.
		// log.Fatal().Err(err).Msg("Failed to load configuration") -> Standart loglama ile değiştir.
		// Bu, zerolog'un varsayılan logger'ını kullanır.
		zerolog.New(os.Stderr).Fatal().Err(err).Msg("Failed to load configuration")
	}

	// DÜZELTME: Logger artık dinamik olarak ENV ve LOG_LEVEL'i alıyor.
	log := logger.New(serviceName, cfg.Env, cfg.LogLevel)

	log.Info().
		Str("version", ServiceVersion).
		Str("commit", GitCommit).
		Str("build_date", BuildDate).
		Str("profile", cfg.Env).
		Msg("🚀 Starting Sentiric API Gateway Service")
	
	// DÜZELTME: `Run` fonksiyonuna hem config hem de logger geçirilir.
	if err := gateway.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("Gateway failed to run")
	}
}
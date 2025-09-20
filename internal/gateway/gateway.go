// sentiric-api-gateway-service/internal/gateway/gateway.go
package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	userv1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// DÜZELTME: Config struct'ına Env ve LogLevel alanları eklendi.
type Config struct {
	HttpPort        string
	UserServiceAddr string
	CertPath        string
	KeyPath         string
	CaPath          string
	Env             string
	LogLevel        string
}

// DÜZELTME: LoadConfig artık Env ve LogLevel'i de okuyor.
func LoadConfig() (*Config, error) {
	return &Config{
		HttpPort:        os.Getenv("API_GATEWAY_HTTP_PORT"),
		UserServiceAddr: os.Getenv("USER_SERVICE_GRPC_URL"),
		CertPath:        os.Getenv("API_GATEWAY_CERT_PATH"),
		KeyPath:         os.Getenv("API_GATEWAY_KEY_PATH"),
		CaPath:          os.Getenv("GRPC_TLS_CA_PATH"),
		Env:             os.Getenv("ENV"),
		LogLevel:        os.Getenv("LOG_LEVEL"),
	}, nil
}

func loggingMiddleware(next http.Handler, log zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Info().
			Str("http_method", r.Method).
			Str("http_path", r.URL.Path).
			Dur("duration", duration).
			Msg("http.request.completed")
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func Run(cfg *Config, log zerolog.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/healthz", healthzHandler)

	go func() {
		log.Info().Str("dependency", cfg.UserServiceAddr).Msg("Arka uç gRPC servisine bağlanılmaya çalışılıyor...")

		var creds credentials.TransportCredentials
		var err error

		// DÜZELTME: Artık cfg.Env üzerinden kontrol yapılıyor.
		if cfg.CertPath == "" || cfg.Env == "development" {
			log.Warn().Msg("mTLS sertifikaları bulunamadı veya geliştirme ortamı. Güvensiz (insecure) gRPC bağlantısı kullanılıyor.")
			creds = insecure.NewCredentials()
		} else {
			creds, err = newClientTLS(cfg.CertPath, cfg.KeyPath, cfg.CaPath, cfg.UserServiceAddr)
			if err != nil {
				log.Error().Err(err).Msg("gRPC istemci TLS kimlik bilgileri oluşturulamadı. Gateway proxy olarak çalışmayacak.")
				return
			}
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}

		var conn *grpc.ClientConn
		for i := 0; i < 10; i++ {
			dialCtx, dialCancel := context.WithTimeout(ctx, 10*time.Second)
			conn, err = grpc.DialContext(dialCtx, cfg.UserServiceAddr, opts...)
			dialCancel()
			if err == nil {
				log.Info().Str("dependency", cfg.UserServiceAddr).Msg("Arka uç gRPC servisine başarıyla bağlanıldı.")
				break
			}
			log.Warn().Err(err).Int("attempt", i+1).Msg("Arka uç servisine bağlanılamadı, tekrar denenecek...")
			time.Sleep(5 * time.Second)
		}

		if err != nil {
			log.Error().Err(err).Msg("Maksimum deneme sonrası arka uç servisine bağlanılamadı. Gateway proxy olarak çalışmayacak.")
			return
		}
		defer conn.Close()

		grpcGatewayMux := runtime.NewServeMux()
		err = userv1.RegisterUserServiceHandler(ctx, grpcGatewayMux, conn)
		if err != nil {
			log.Error().Err(err).Msg("gRPC gateway handler'ı kaydedilemedi.")
			return
		}

		mainMux.Handle("/", loggingMiddleware(grpcGatewayMux, log))
		log.Info().Msg("gRPC Gateway başarıyla başlatıldı ve ana yönlendiriciye eklendi.")
	}()

	log.Info().Str("port", cfg.HttpPort).Msg("HTTP sunucusu başlatılıyor (/healthz endpoint'i aktif).")
	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), mainMux)
}

func newClientTLS(certPath, keyPath, caPath, serverAddr string) (credentials.TransportCredentials, error) {
	clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	serverName := strings.Split(serverAddr, ":")[0]

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		ServerName:   serverName,
	}), nil
}
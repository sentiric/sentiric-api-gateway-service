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
	"time" // YENİ

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	userv1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ... (Config struct ve LoadConfig fonksiyonu aynı) ...
type Config struct {
	HttpPort        string
	UserServiceAddr string
	CertPath        string
	KeyPath         string
	CaPath          string
}

func LoadConfig() (*Config, error) {
	return &Config{
		HttpPort:        os.Getenv("API_GATEWAY_HTTP_PORT"),
		UserServiceAddr: os.Getenv("USER_SERVICE_GRPC_URL"),
		CertPath:        os.Getenv("API_GATEWAY_CERT_PATH"),
		KeyPath:         os.Getenv("API_GATEWAY_KEY_PATH"),
		CaPath:          os.Getenv("GRPC_TLS_CA_PATH"),
	}, nil
}

// YENİ: Loglama Middleware'i
func loggingMiddleware(next http.Handler, log zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// İsteği bir sonraki handler'a (bizim durumumuzda gRPC gateway mux'ı) iletiyoruz.
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Info().
			Str("http_method", r.Method).
			Str("http_path", r.URL.Path).
			// Not: Durum kodunu almak için daha gelişmiş bir response writer sarmalayıcı gerekir.
			// Şimdilik bu temel loglama yeterlidir.
			Dur("duration", duration).
			Msg("http.request.completed")
	})
}

func Run(cfg *Config, log zerolog.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	creds, err := newClientTLS(cfg.CertPath, cfg.KeyPath, cfg.CaPath, cfg.UserServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to create client TLS credentials: %w", err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	err = userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.UserServiceAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register user service handler: %w", err)
	}

	// YENİ: Ana handler'ımızı loglama middleware'i ile sarmalıyoruz
	handlerWithLogging := loggingMiddleware(mux, log)

	log.Info().Str("port", cfg.HttpPort).Msg("Starting HTTP server for gRPC Gateway")
	// YENİ: Sarmalanmış handler'ı kullanıyoruz
	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), handlerWithLogging)
}

// ... (newClientTLS fonksiyonu aynı) ...
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

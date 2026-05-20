// [ARCH-COMPLIANCE] SOP-01: Generative Media API Routing
package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	audiov1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/audio_gen/v1"
	imagev1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/image/v1"
	userv1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/user/v1"
	videov1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/video/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	HttpPort         string
	UserServiceAddr  string
	VideoGatewayAddr string // YENİ
	ImageGatewayAddr string // YENİ
	AudioGatewayAddr string // YENİ
	CertPath         string
	KeyPath          string
	CaPath           string
	Env              string
	LogLevel         string
}

func LoadConfig() (*Config, error) {
	return &Config{
		HttpPort:         os.Getenv("API_GATEWAY_HTTP_PORT"),
		UserServiceAddr:  os.Getenv("USER_SERVICE_GRPC_URL"),
		VideoGatewayAddr: os.Getenv("VIDEO_GATEWAY_GRPC_URL"), // Örn: video-gateway-service:16101
		ImageGatewayAddr: os.Getenv("IMAGE_GATEWAY_GRPC_URL"), // Örn: image-gateway-service:16201
		AudioGatewayAddr: os.Getenv("AUDIO_GATEWAY_GRPC_URL"), // Örn: audio-gateway-service:16301
		CertPath:         os.Getenv("API_GATEWAY_CERT_PATH"),
		KeyPath:          os.Getenv("API_GATEWAY_KEY_PATH"),
		CaPath:           os.Getenv("GRPC_TLS_CA_PATH"),
		Env:              os.Getenv("ENV"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
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
		grpcGatewayMux := runtime.NewServeMux()
		var creds credentials.TransportCredentials
		var err error

		if cfg.CertPath == "" || cfg.Env == "development" {
			log.Warn().Msg("Insecure (Development) mode active for downstream gRPC connections.")
			creds = insecure.NewCredentials()
		} else {
			// Helper function to create creds (implemented below)
			creds, err = newClientTLS(cfg.CertPath, cfg.KeyPath, cfg.CaPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create TLS credentials")
			}
		}

		opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

		// 1. User Service Bağlantısı
		if cfg.UserServiceAddr != "" {
			conn, err := grpc.DialContext(ctx, cfg.UserServiceAddr, opts...)
			if err == nil {
				userv1.RegisterUserServiceHandler(ctx, grpcGatewayMux, conn)
				log.Info().Str("target", cfg.UserServiceAddr).Msg("User Service route registered.")
			}
		}

		// 2. Video Gateway Bağlantısı
		if cfg.VideoGatewayAddr != "" {
			conn, err := grpc.DialContext(ctx, cfg.VideoGatewayAddr, opts...)
			if err == nil {
				videov1.RegisterVideoGatewayServiceHandler(ctx, grpcGatewayMux, conn)
				log.Info().Str("target", cfg.VideoGatewayAddr).Msg("Video Gateway route registered.")
			}
		}

		// 3. Image Gateway Bağlantısı
		if cfg.ImageGatewayAddr != "" {
			conn, err := grpc.DialContext(ctx, cfg.ImageGatewayAddr, opts...)
			if err == nil {
				imagev1.RegisterImageGatewayServiceHandler(ctx, grpcGatewayMux, conn)
				log.Info().Str("target", cfg.ImageGatewayAddr).Msg("Image Gateway route registered.")
			}
		}

		// 4. Audio Gateway Bağlantısı
		if cfg.AudioGatewayAddr != "" {
			conn, err := grpc.DialContext(ctx, cfg.AudioGatewayAddr, opts...)
			if err == nil {
				audiov1.RegisterAudioGatewayServiceHandler(ctx, grpcGatewayMux, conn)
				log.Info().Str("target", cfg.AudioGatewayAddr).Msg("Audio Gateway route registered.")
			}
		}

		mainMux.Handle("/", loggingMiddleware(grpcGatewayMux, log))
		log.Info().Msg("gRPC Gateway successfully mounted all available endpoints.")
	}()

	log.Info().Str("port", cfg.HttpPort).Msg("API Gateway HTTP server running.")
	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), mainMux)
}

func newClientTLS(certPath, keyPath, caPath string) (credentials.TransportCredentials, error) {
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

	// "sentiric.cloud" gRPC sertifikalarımızdaki varsayılan SAN (Subject Alternative Name)
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		ServerName:   "sentiric.cloud",
	}), nil
}

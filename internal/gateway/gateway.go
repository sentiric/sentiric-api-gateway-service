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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	userv1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	HttpPort        string
	UserServiceAddr string
	// ... diğer servis adresleri ...
	CertPath string
	KeyPath  string
	CaPath   string
}

func LoadConfig() (*Config, error) {
	// Ortam değişkenlerinden konfigürasyonu yükle
	return &Config{
		HttpPort:        os.Getenv("API_GATEWAY_SERVICE_PORT"),
		UserServiceAddr: os.Getenv("USER_SERVICE_GRPC_URL"),
		CertPath:        os.Getenv("API_GATEWAY_SERVICE_CERT_PATH"),
		KeyPath:         os.Getenv("API_GATEWAY_SERVICE_KEY_PATH"),
		CaPath:          os.Getenv("GRPC_TLS_CA_PATH"),
	}, nil
}

// Run, HTTP sunucusunu başlatır ve gRPC isteklerini yönlendirir.
func Run(cfg *Config, log zerolog.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Gelen HTTP isteklerini yönlendirecek olan ana multiplexer'ı oluştur
	mux := runtime.NewServeMux()

	// user-service için mTLS istemci kimlik bilgilerini oluştur
	creds, err := newClientTLS(cfg.CertPath, cfg.KeyPath, cfg.CaPath, cfg.UserServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to create client TLS credentials: %w", err)
	}

	// user-service'e bağlanmak için gRPC dial options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	// HTTP isteklerini user-service'e yönlendirmek için register et
	err = userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.UserServiceAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register user service handler: %w", err)
	}

	log.Info().Str("port", cfg.HttpPort).Msg("Starting HTTP server for gRPC Gateway")
	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), mux)
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

	// 'serverAddr' genellikle 'hostname:port' formatındadır, sadece hostname kısmını alalım
	serverName := strings.Split(serverAddr, ":")[0]

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		ServerName:   serverName,
	}), nil
}

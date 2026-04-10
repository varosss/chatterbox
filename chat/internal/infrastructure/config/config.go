package config

import (
	"os"
	"strings"
)

type Database struct {
	DSN string
}

type HttpServer struct {
	Origins    string
	HostURL    string
	HostDomain string
	Port       string
}

type RabbitMQ struct {
	URL      string
	Exchange string
}

type JWT struct {
	Issuer string
}

type Security struct {
	PublicKeyPath string
}

type Config struct {
	Database   Database
	HttpServer HttpServer
	RabbitMQ   RabbitMQ
	Security   Security
	JWT        JWT
}

func Load() (*Config, error) {
	hostURL := os.Getenv("HTTP_SERVER_HOST_URL")
	domain := strings.TrimPrefix(hostURL, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	cfg := &Config{
		Database: Database{
			DSN: os.Getenv("POSTGRES_URL"),
		},
		HttpServer: HttpServer{
			Origins:    getEnv("HTTP_SERVER_ALLOW_ORIGIN", "*"),
			HostDomain: domain,
			HostURL:    hostURL,
			Port:       getEnv("HTTP_SERVER_PORT", "80"),
		},
		RabbitMQ: RabbitMQ{
			URL:      os.Getenv("RABBITMQ_URL"),
			Exchange: os.Getenv("RABBITMQ_EXCHANGE"),
		},
		Security: Security{
			PublicKeyPath: os.Getenv("SECURITY_PUBLIC_KEY_PATH"),
		},
		JWT: JWT{
			Issuer: os.Getenv("JWT_ISSUER"),
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package config

import (
	"os"
)

type Database struct {
	DSN string
}

type HttpServer struct {
	Port string
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
	cfg := &Config{
		Database: Database{
			DSN: os.Getenv("POSTGRES_URL"),
		},
		HttpServer: HttpServer{
			Port: getEnv("HTTP_SERVER_PORT", "80"),
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

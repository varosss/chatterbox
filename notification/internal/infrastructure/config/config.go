package config

import (
	"os"
	"strconv"
	"strings"
)

type Http struct {
	PublicURL string
	Host      string
	Port      string
}

type CORS struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

type RabbitMQ struct {
	URL      string
	Exchange string
	Queue    string
}

type JWT struct {
	Issuer string
}

type Security struct {
	PublicKeyPath string
}

type Config struct {
	Http     Http
	RabbitMQ RabbitMQ
	Security Security
	JWT      JWT
	CORS     CORS
}

func Load() (*Config, error) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	allowCreds, err := strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Http: Http{
			Host:      os.Getenv("HTTP_HOST"),
			Port:      getEnv("HTTP_SERVER_PORT", "80"),
			PublicURL: os.Getenv("PUBLIC_BASE_URL"),
		},
		CORS: CORS{
			AllowedOrigins:   allowedOrigins,
			AllowCredentials: allowCreds,
		},
		RabbitMQ: RabbitMQ{
			URL:      os.Getenv("RABBITMQ_URL"),
			Exchange: os.Getenv("RABBITMQ_EXCHANGE"),
			Queue:    os.Getenv("RABBITMQ_QUEUE"),
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

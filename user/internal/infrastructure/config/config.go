package config

import (
	"os"
	"strconv"
	"time"
)

type Database struct {
	DSN string
}

type JWT struct {
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type HttpServer struct {
	HostURL string
	Port    string
}

type Security struct {
	PublicKeyPath  string
	PrivateKeyPath string
	HashCost       int
}

type Config struct {
	HttpServer HttpServer
	Database   Database
	Security   Security
	JWT        JWT
}

func Load() (*Config, error) {
	accessTTL, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		return nil, err
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
	if err != nil {
		return nil, err
	}

	cost, err := strconv.Atoi(os.Getenv("SECURITY_HASH_COST"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Security: Security{
			PublicKeyPath:  os.Getenv("SECURITY_PUBLIC_KEY_PATH"),
			PrivateKeyPath: os.Getenv("SECURITY_PRIVATE_KEY_PATH"),
			HashCost:       cost,
		},
		Database: Database{
			DSN: os.Getenv("POSTGRES_URL"),
		},
		JWT: JWT{
			Issuer:     os.Getenv("JWT_ISSUER"),
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
		},
		HttpServer: HttpServer{
			Port: getEnv("HTTP_SERVER_PORT", "80"),
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

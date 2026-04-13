package config

import (
	"chatterbox/user/internal/infrastructure/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type Http struct {
	PublicURL string
	Host      string
	Port      string
}

type CORS struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

type Cookie struct {
	Secure   bool
	SameSite http.SameSite
	HttpOnly bool
	Domain   string
}

type Security struct {
	PublicKeyPath  string
	PrivateKeyPath string
	HashCost       int
}

type Config struct {
	Http     Http
	Database Database
	Security Security
	JWT      JWT
	CORS     CORS
	Cookie   Cookie
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

	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	allowCreds, err := strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))
	if err != nil {
		return nil, err
	}

	cookieSecure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		return nil, err
	}
	cookieHttpOnly, err := strconv.ParseBool(os.Getenv("COOKIE_HTTPONLY"))
	if err != nil {
		return nil, err
	}
	cookieSameSite := utils.ParseSameSite(getEnv("COOKIE_SAMESITE", "lax"))

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
		Http: Http{
			Host:      os.Getenv("HTTP_HOST"),
			Port:      getEnv("HTTP_SERVER_PORT", "80"),
			PublicURL: os.Getenv("PUBLIC_BASE_URL"),
		},
		CORS: CORS{
			AllowedOrigins:   allowedOrigins,
			AllowCredentials: allowCreds,
		},
		Cookie: Cookie{
			Secure:   cookieSecure,
			SameSite: cookieSameSite,
			Domain:   getEnv("COOKIE_DOMAIN", ""),
			HttpOnly: cookieHttpOnly,
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

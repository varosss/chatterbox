package app

import (
	"chatterbox/pkg/clock"
	"chatterbox/pkg/httpmiddleware"
	pkgsecurity "chatterbox/pkg/security"
	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/infrastructure/adapter/auth"
	"chatterbox/user/internal/infrastructure/adapter/postgres/repo"
	"chatterbox/user/internal/infrastructure/adapter/security"
	"chatterbox/user/internal/infrastructure/config"
	"chatterbox/user/internal/infrastructure/controller/httphandler"
	"chatterbox/user/internal/infrastructure/db"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server
	closers    []func() error
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	var closers []func() error

	database, err := db.NewPostgres(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}
	closers = append(closers, func() error {
		database.Close()
		return nil
	})

	userRepo := repo.NewUserPgxRepo(database)
	refreshTokenRepo := repo.NewRefreshTokenPgxRepo(database)

	publicKey, err := pkgsecurity.LoadPublicKey(cfg.Security.PublicKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := pkgsecurity.LoadPrivateKey(cfg.Security.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	passwordHasher := security.NewBcryptPasswordHasher(cfg.Security.HashCost)
	passwordVerifier := security.NewBcryptPasswordVerifier()

	jwtTokenSigner := auth.NewJWTSigner(
		privateKey,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	jwtTokenVerifier := auth.NewJWTVerifier(
		publicKey,
		cfg.JWT.Issuer,
	)

	systemClock := clock.NewSystemClock()

	registerUC := usecase.NewRegisterUseCase(
		userRepo,
		passwordHasher,
	)
	loginUC := usecase.NewLoginUseCase(
		userRepo,
		passwordVerifier,
		refreshTokenRepo,
		jwtTokenSigner,
		systemClock,
		cfg.JWT.RefreshTTL,
	)
	logoutUC := usecase.NewLogoutUseCase(
		jwtTokenVerifier,
		refreshTokenRepo,
	)
	refreshTokenUC := usecase.NewRefreshTokenUseCase(
		refreshTokenRepo,
		jwtTokenVerifier,
		jwtTokenSigner,
		systemClock,
		cfg.JWT.RefreshTTL,
	)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)

	authHandler := httphandler.NewAuthHandler(
		registerUC,
		loginUC,
		logoutUC,
		refreshTokenUC,
	)
	userHandler := httphandler.NewUserHandler(
		listUsersUC,
	)

	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(httpmiddleware.CORSMiddleware())

	ginEngine.POST("/register", authHandler.Register)
	ginEngine.POST("/login", authHandler.Login)
	ginEngine.POST("/logout", authHandler.Logout)
	ginEngine.POST("/refresh_token", authHandler.RefreshToken)

	ginEngine.GET(
		"/users",
		httpmiddleware.AuthMiddleware(auth.NewTokenVerifierWrapper(jwtTokenVerifier)),
		userHandler.List,
	)

	return &App{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.HttpServer.Port),
			Handler: ginEngine,
		},
		closers: closers,
	}, nil
}

func (a *App) Run() error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	errCh := make(chan error, 1)

	go func() {
		fmt.Printf("HTTP server running on %s\n", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case sig := <-stop:
		fmt.Printf("Received signal %s, shutting down...\n", sig)
		cancel()
		return a.Shutdown()

	case err := <-errCh:
		cancel()
		return err
	}
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	for _, closer := range a.closers {
		if err := closer(); err != nil {
			return err
		}
	}

	return nil
}

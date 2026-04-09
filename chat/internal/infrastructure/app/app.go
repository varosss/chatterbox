package app

import (
	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/infrastructure/adapter/event/producer"
	"chatterbox/chat/internal/infrastructure/adapter/postgres/repo"
	"chatterbox/chat/internal/infrastructure/config"
	"chatterbox/chat/internal/infrastructure/controller/httphandler"
	"chatterbox/chat/internal/infrastructure/db"
	"chatterbox/pkg/auth"
	"chatterbox/pkg/httpmiddleware"
	"chatterbox/pkg/security"
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

	chatRepo := repo.NewChatPgxRepo(database)
	messageRepo := repo.NewMessagePgxRepo(database)

	eventProducer, err := producer.NewRabbitMQProducer(
		cfg.RabbitMQ.URL,
		cfg.RabbitMQ.Exchange,
	)
	if err != nil {
		return nil, err
	}
	closers = append(closers, eventProducer.Close)

	createChatUC := usecase.NewCreateChatUseCase(eventProducer, chatRepo)
	listChatsUC := usecase.NewListChatsUseCase(chatRepo)
	createMessageUC := usecase.NewCreateMessageUseCase(eventProducer, messageRepo, chatRepo)
	listMessagesUC := usecase.NewListMessagesUseCase(messageRepo)

	publicKey, err := security.LoadPublicKey(cfg.Security.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(httpmiddleware.CORSMiddleware())
	ginEngine.Use(
		httpmiddleware.AuthMiddleware(
			auth.NewJWTVerifier(publicKey, cfg.JWT.Issuer),
		),
	)

	chatHandler := httphandler.NewChatHandler(createChatUC, listChatsUC)
	messageHandler := httphandler.NewMessageHandler(createMessageUC, listMessagesUC)

	ginEngine.POST("/chats", chatHandler.Create)
	ginEngine.GET("/chats", chatHandler.List)
	ginEngine.POST("/messages", messageHandler.Create)
	ginEngine.GET("/messages", messageHandler.List)

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

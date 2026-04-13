package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"chatterbox/notification/internal/application/usecase"
	"chatterbox/notification/internal/infrastructure/adapter/sender"
	"chatterbox/notification/internal/infrastructure/adapter/websocket/hub"
	"chatterbox/notification/internal/infrastructure/config"
	"chatterbox/notification/internal/infrastructure/controller/event/consumer"
	"chatterbox/notification/internal/infrastructure/controller/event/eventhandler"
	"chatterbox/notification/internal/infrastructure/controller/httphandler"
	"chatterbox/pkg/auth"
	"chatterbox/pkg/httpmiddleware"
	"chatterbox/pkg/security"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg           *config.Config
	httpServer    *http.Server
	eventConsumer *consumer.RabbitMQConsumer
	closers       []func() error
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	var closers []func() error

	eventConsumer, err := consumer.NewRabbitMQConsumer(
		cfg.RabbitMQ.URL,
		cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.Queue,
	)
	if err != nil {
		return nil, err
	}
	closers = append(closers, eventConsumer.Close)

	inMemoryHub := hub.NewInMemoryHub()
	closers = append(closers, inMemoryHub.Close)

	notificationSender := sender.NewWebSocketSender(inMemoryHub)

	notifyMessageUC := usecase.NewNotifyMessageUseCase(notificationSender)

	wsHandler := httphandler.NewWSHandler(inMemoryHub)
	messageCreatedHandler := eventhandler.NewMessageCreatedHandler(notifyMessageUC)

	publicKey, err := security.LoadPublicKey(cfg.Security.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(httpmiddleware.CORSMiddleware(cfg.CORS.AllowedOrigins))
	ginEngine.Use(
		httpmiddleware.AuthMiddleware(
			auth.NewJWTVerifier(publicKey, cfg.JWT.Issuer),
		),
	)

	ginEngine.GET("/ws", wsHandler.Handle)

	if err := eventConsumer.Register("message.created", messageCreatedHandler); err != nil {
		return nil, err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Http.Port),
		Handler: ginEngine,
	}

	return &App{
		cfg:           cfg,
		httpServer:    httpServer,
		eventConsumer: eventConsumer,
		closers:       closers,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	if err := a.eventConsumer.Start(ctx); err != nil {
		return err
	}

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

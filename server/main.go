package main

import (
	"context"
	"diploma/internal/auth"
	"diploma/internal/handlers"
	"diploma/pkg/config"
	"diploma/pkg/db"
	"diploma/pkg/logger"
	"diploma/pkg/logger/message"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func waitForShutdown(srv *http.Server, closeDB func()) {
	quit := make(chan os.Signal, 1) // Ограничил size канала, т.к. ловлю всего 1 сигнал
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info(&message.LogMessage{Message: "Gracefully shutting down server..."})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 секунд выглядит достаточно
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// Если таймаут прошел, то принудительное закрытие + логируем ошибку
		logger.Log.Error(&message.LogMessage{Message: fmt.Sprintf("Server forced to shutdown with error: %v", err)})
		if err = srv.Close(); err != nil {
			logger.Log.Error(&message.LogMessage{Message: "Error while closing server"})
		}
	}

	// Закрываем БД
	if closeDB != nil {
		closeDB()
	}
}

func addrFromBaseURL(raw, fallback string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return fallback // например ":8080"
	}
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return net.JoinHostPort(host, port)
}

func main() {
	flag.Parse()
	cfg := config.InitConfig()
	_ = logger.New("info")

	// Инициализация БД
	instance, err := db.InitSQL(cfg.DBConfig)
	if err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	db.PrepareDB(instance)

	api := handlers.New(handlers.Dependency{
		Cfg: cfg,
		DB:  instance,
	})

	// Запуск сервера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(auth.MiddlewareAuth(cfg))

	router.GET("health", api.GetHealthCheck)
	router.POST("/api/user/register", api.PostRegister)
	router.POST("/api/user/login", api.PostLogin)
	router.POST("/api/v1/texts", api.PostText)
	router.GET("/api/v1/texts/:id", api.GetTextHandler)
	router.GET("/api/v1/texts", api.GetTextsListHandler)
	router.POST("/api/v1/cards", api.PostCard)
	router.GET("/api/v1/cards/:id", api.GetCardHandler)
	router.GET("/api/v1/cards", api.GetCardsListHandler)
	router.POST("/api/v1/passwords", api.PostPassword)
	router.GET("/api/v1/passwords/:id", api.GetPasswordHandler)
	router.GET("/api/v1/passwords", api.GetPasswordsListHandler)
	router.POST("/api/v1/binaries", api.PostBinary)
	router.GET("/api/v1/binaries/:id", api.GetBinaryHandler)
	router.GET("/api/v1/binaries", api.GetBinariesListHandler)

	addr := addrFromBaseURL(cfg.BaseURL, cfg.ServerAddr)
	Srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Пробуем запуститься, только потом стартую горутину для graceful shutdown
	logger.Log.Info(&message.LogMessage{Message: fmt.Sprintf("Starting server on port %s", cfg.ServerAddr)})
	if err := Srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// graceful shutdown
	go waitForShutdown(Srv, func() {
		if instance == nil {
			logger.Log.Error(&message.LogMessage{Message: "DB instance is nil on shutdown"})
			return
		}
		instance.CloseSqlInstance()
	})
}

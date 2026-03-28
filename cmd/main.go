package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appdb "github.com/bakdaaswandi5818/pencatat-keuangan/pkg/database"
	applogger "github.com/bakdaaswandi5818/pencatat-keuangan/pkg/logger"

	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/handler"
	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/repository"
	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/service"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	log := applogger.New()

	// --- Database ---
	dsn := getEnv("DB_PATH", "finance.db")
	db, err := appdb.New(dsn)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	log.WithField("dsn", dsn).Info("database connected")

	// --- Layers ---
	txRepo := repository.NewGORMTransactionRepository(db)
	txSvc := service.NewTransactionService(txRepo)
	txHandler := handler.NewTransactionHandler(txSvc)

	// --- Echo ---
	e := echo.New()
	e.HideBanner = true

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			entry := log.WithFields(map[string]interface{}{
				"method": v.Method,
				"uri":    v.URI,
				"status": v.Status,
			})
			if v.Error != nil {
				entry.WithError(v.Error).Warn("request")
			} else {
				entry.Info("request")
			}
			return nil
		},
	}))

	txHandler.Register(e)

	// --- Graceful shutdown ---
	addr := getEnv("ADDR", ":8080")
	go func() {
		log.WithField("addr", addr).Info("starting server")
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.WithError(err).Error("shutdown error")
	}
	log.Info("server stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

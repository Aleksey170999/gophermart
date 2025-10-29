package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	config "github.com/Aleksey170999/go-loyaty/configs"
	"github.com/Aleksey170999/go-loyaty/internal/handler"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
	srv "github.com/Aleksey170999/go-loyaty/internal/server"
	"github.com/Aleksey170999/go-loyaty/internal/service"
)

func main() {
	cfg := config.NewConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if cfg.JWTSecret == "" {
		logger.Fatal("JWT_SECRET is not set")
	}

	db, err := repository.NewPostgresDB(cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal("failed to initialize db", zap.Error(err))
	}
	err = ApplyMigrations(db.DB)
	if err != nil {
		logger.Fatal(err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos, cfg.JWTSecret)
	handlers := handler.NewHandler(services, services.JWTService, logger)

	srv := new(srv.Server)
	go func() {
		logger.Info("Starting server", zap.String("address", cfg.RunAddr))
		if err := srv.Run(cfg.RunAddr, handlers.InitRoutes()); err != nil {
			logger.Fatal("failed to run server", zap.Error(err))
		}
	}()

	logger.Info("Application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("error occurred during server shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func ApplyMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "./migrations/"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}

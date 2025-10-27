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

	config "github.com/Aleksey170999/go-loyaty/configs"
	"github.com/Aleksey170999/go-loyaty/internal/handler"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
	srv "github.com/Aleksey170999/go-loyaty/internal/server"
	"github.com/Aleksey170999/go-loyaty/internal/service"
)

func main() {
	cfg := config.NewConfig()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	db, err := repository.NewPostgresDB(cfg.DatabaseDSN)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	err = ApplyMigrations(db.DB)
	if err != nil {
		fmt.Println(err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(srv.Server)
	go func() {
		if err := srv.Run(cfg.RunAddr, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

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

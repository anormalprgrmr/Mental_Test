package main

import (
	"context"
	"log/slog"
	"os"

	"mental_test/internal/config"
	"mental_test/internal/db"
	"mental_test/internal/handlers"
	"mental_test/internal/models"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	store, err := db.Open(context.Background(), cfg.DBPath)
	if err != nil {
		logger.Error("open storage", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			logger.Error("close storage", "error", err)
		}
	}()

	catalog, err := models.DefaultCatalog()
	if err != nil {
		logger.Error("load test catalog", "error", err)
		os.Exit(1)
	}

	bot, err := handlers.New(handlers.Config{
		Token:  cfg.BotToken,
		Poller: cfg.Poller,
		Admins: cfg.Admins,
	}, store, catalog, logger)
	if err != nil {
		logger.Error("create bot", "error", err)
		os.Exit(1)
	}

	bot.Start()
}

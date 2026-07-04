package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains runtime settings loaded from environment variables.
type Config struct {
	BotToken string
	Admins   []int64
	DBPath   string
	Poller   time.Duration
}

func Load() (Config, error) {
	botToken := firstNonEmpty(os.Getenv("BOT_TOKEN"), os.Getenv("BIT_TOKEN"))
	if botToken == "" {
		return Config{}, fmt.Errorf("BOT_TOKEN is required")
	}

	admins, err := parseAdmins(os.Getenv("ADMINS"))
	if err != nil {
		return Config{}, err
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "mental_test.db"
	}

	return Config{
		BotToken: botToken,
		Admins:   admins,
		DBPath:   dbPath,
		Poller:   10 * time.Second,
	}, nil
}

func parseAdmins(raw string) ([]int64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	admins := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid admin telegram id %q: %w", part, err)
		}
		admins = append(admins, id)
	}
	return admins, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"mental_test/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate(ctx context.Context) error {
	statements := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS users (
			telegram_id INTEGER PRIMARY KEY,
			username TEXT NOT NULL DEFAULT '',
			first_name TEXT NOT NULL DEFAULT '',
			last_name TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS test_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			test_id TEXT NOT NULL,
			result_json TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(telegram_id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_user_id ON test_results(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_test_id ON test_results(test_id);`,
	}

	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate sqlite database: %w", err)
		}
	}
	return nil
}

func (s *Store) UpsertUser(ctx context.Context, user models.TelegramUser) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (telegram_id, username, first_name, last_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(telegram_id) DO UPDATE SET
			username = excluded.username,
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			updated_at = CURRENT_TIMESTAMP;
	`, user.ID, user.Username, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("upsert user %d: %w", user.ID, err)
	}
	return nil
}

func (s *Store) SaveResult(ctx context.Context, userID int64, testID string, result models.Result) (models.SavedResult, error) {
	payload, err := json.Marshal(result)
	if err != nil {
		return models.SavedResult{}, fmt.Errorf("marshal result: %w", err)
	}

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO test_results (user_id, test_id, result_json, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP);
	`, userID, testID, string(payload))
	if err != nil {
		return models.SavedResult{}, fmt.Errorf("save result for user %d test %s: %w", userID, testID, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.SavedResult{}, fmt.Errorf("read saved result id: %w", err)
	}

	return models.SavedResult{
		ID:        id,
		UserID:    userID,
		TestID:    testID,
		Result:    result,
		CreatedAt: time.Now(),
	}, nil
}

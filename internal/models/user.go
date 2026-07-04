package models

import "time"

// TelegramUser represents a Telegram user stored by the application.
type TelegramUser struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string
}

// SavedResult is the persisted representation returned after saving a completed test.
type SavedResult struct {
	ID        int64
	UserID    int64
	TestID    string
	Result    Result
	CreatedAt time.Time
}

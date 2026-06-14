package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"basalam-feedback/internal/model"
)

type FeedbackRepository struct {
	db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func OpenWithRetry(databaseURL string, attempts int, delay time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	for attempt := 1; attempt <= attempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			return db, nil
		}
		time.Sleep(delay)
	}

	_ = db.Close()
	return nil, fmt.Errorf("database not ready after %d attempts: %w", attempts, err)
}

func EnsureSchema(db *sql.DB) error {
	_, err := db.Exec(`
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS feedbacks (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title VARCHAR(255) NOT NULL,
	message TEXT NOT NULL,
	status VARCHAR(50) NOT NULL DEFAULT 'submitted',
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`)
	return err
}

func (r *FeedbackRepository) Create(ctx context.Context, title, message string) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO feedbacks (title, message, status)
VALUES ($1, $2, $3)
`, title, message, model.StatusSubmitted)
	return err
}

func (r *FeedbackRepository) List(ctx context.Context) ([]model.Feedback, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, title, message, status, created_at
FROM feedbacks
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := make([]model.Feedback, 0)
	for rows.Next() {
		var feedback model.Feedback
		if err := rows.Scan(&feedback.ID, &feedback.Title, &feedback.Message, &feedback.Status, &feedback.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}
	return feedbacks, rows.Err()
}

func (r *FeedbackRepository) UpdateStatus(ctx context.Context, id, status string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
UPDATE feedbacks
SET status = $1
WHERE id = $2
`, status, id)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, nil
	}
	return true, nil
}

func ValidateStatus(status string) error {
	switch status {
	case model.StatusReviewing, model.StatusResolved:
		return nil
	default:
		return errors.New("invalid status")
	}
}

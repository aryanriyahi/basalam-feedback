package model

import "time"

const (
	StatusSubmitted = "submitted"
	StatusReviewing = "reviewing"
	StatusResolved  = "resolved"
)

type Feedback struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

package models

import (
    "time"

    "github.com/google/uuid"
)

type Poll struct {
    ID        uuid.UUID `json:"id"`
    Question  string    `json:"question"`
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt *time.Time `json:"expires_at,omitempty"`
    IsActive  bool      `json:"is_active"`
}

type Option struct {
    ID         uuid.UUID `json:"id"`
    PollID     uuid.UUID `json:"poll_id"`
    Value      string    `json:"value"`
    VotesCount int       `json:"votes_count"`
}

type Vote struct {
    PollID    uuid.UUID `json:"poll_id"`
    OptionID  uuid.UUID `json:"option_id"`
    IPAddress string    `json:"-"`
}

type CreatePollRequest struct {
    Question  string   `json:"question" binding:"required"`
    Options   []string `json:"options" binding:"required,min=2"`
    ExpiresAt *time.Time `json:"expires_at"`
}
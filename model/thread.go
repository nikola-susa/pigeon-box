package model

import (
	"time"
)

type Thread struct {
	ID                *int           `json:"id,omitempty" db:"id"`
	Name              string         `json:"name,omitempty" db:"name"`
	Description       *string        `json:"description,omitempty" db:"description"`
	UserID            int            `json:"user_id,omitempty" db:"user_id"`
	SlackID           string         `json:"slack_id,omitempty" db:"slack_id"`
	CreatedAt         string         `json:"created_at,omitempty" db:"created_at"`
	ExpiresAt         *time.Duration `json:"expires_at,omitempty" db:"expires_at"`
	MessagesExpiresAt *time.Duration `json:"messages_expires_at,omitempty" db:"messages_expires_at"`
	Key               string         `json:"key,omitempty" db:"key"`
	SlackTimestamp    *string        `json:"slack_timestamp,omitempty" db:"slack_timestamp"`
}

type RenderThread struct {
	ID                string
	Name              string
	Description       *string
	AuthorID          string
	IsAuthor          bool
	ExpiresAt         string
	MessagesExpiresAt string
	Version           string
}

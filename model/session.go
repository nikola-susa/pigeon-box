package model

import "time"

type Session struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	ThreadID  int       `db:"thread_id"`
	ExpiresAt time.Time `db:"expires_at"`
}

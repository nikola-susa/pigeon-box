package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/pigeon-box/model"
	"golang.org/x/net/context"
	"time"
)

func (s *Store) CreateSession(session model.Session) (*int, error) {
	row := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO session (user_id, thread_id, expires_at) VALUES (?, ?, ?) RETURNING id`,
		session.UserID,
		session.ThreadID,
		session.ExpiresAt,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert session: %w", err)
	}
	return &id, nil
}

func (s *Store) GetSession(id int) (*model.Session, error) {
	var session model.Session

	err := s.db.Get(&session, `SELECT * FROM session WHERE id = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get session by id: %w", err)
		}
		return nil, fmt.Errorf("get session by id: %w", err)
	}

	return &session, nil
}

func (s *Store) GetSessionByUserID(userID int) (*model.Session, error) {
	var session model.Session
	err := s.db.Get(&session, `SELECT * FROM session WHERE user_id = ?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get session by user_id: %w", err)
		}
		return nil, fmt.Errorf("get session by user_id: %w", err)
	}

	return &session, nil
}

func (s *Store) DeleteSession(id int) error {
	_, err := s.db.Exec(`DELETE FROM session WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s *Store) DeleteSessionByUserID(userID int) error {
	_, err := s.db.Exec(`DELETE FROM session WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("delete session by user_id: %w", err)
	}
	return nil
}

func (s *Store) UpdateSessionExpiresAt(id int, expiresAt time.Time) error {
	_, err := s.db.Exec(`UPDATE session SET expires_at = ? WHERE id = ?`, expiresAt, id)
	if err != nil {
		return fmt.Errorf("update session expires_at: %w", err)
	}
	return nil
}

func (s *Store) DeleteExpiredSessions() error {
	_, err := s.db.Exec(`DELETE FROM session WHERE expires_at < DATETIME('now')`)
	if err != nil {
		return fmt.Errorf("delete expired sessions: %w", err)
	}
	return nil
}

func (s *Store) GetExpiredSessions() ([]model.Session, error) {
	var sessions []model.Session
	err := s.db.Select(&sessions, `SELECT * FROM session WHERE expires_at < DATETIME('now')`)
	if err != nil {
		return nil, fmt.Errorf("get expired sessions: %w", err)
	}
	return sessions, nil
}

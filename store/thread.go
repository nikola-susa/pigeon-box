package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/pigeon-box/model"
)

func (s *Store) CreateThread(thread model.Thread) (*int, error) {
	row := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO thread (name, description, user_id, slack_id, key, expires_at, messages_expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		thread.Name,
		thread.Description,
		thread.UserID,
		thread.SlackID,
		thread.Key,
		thread.ExpiresAt,
		thread.MessagesExpiresAt,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert thread: %w", err)
	}
	return &id, nil
}

func (s *Store) SetThreadSlackTimestamp(threadID int, slackTimestamp string) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`UPDATE thread SET slack_timestamp = $1 WHERE id = $2`,
		slackTimestamp,
		threadID,
	)
	if err != nil {
		return fmt.Errorf("set thread slack timestamp: %w", err)
	}
	return nil
}

func (s *Store) GetThread(id int) (*model.Thread, error) {
	var thread model.Thread
	err := s.db.GetContext(
		context.Background(),
		&thread,
		`SELECT * FROM thread WHERE id = $1`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get thread by id: %w", err)
	}
	return &thread, nil
}

func (s *Store) DeleteThread(id int) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`DELETE FROM thread WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete thread: %w", err)
	}
	return nil
}

func (s *Store) DeleteStaleThread() error {
	_, err := s.db.ExecContext(
		context.Background(),
		`DELETE FROM thread WHERE id IN ( SELECT thread_id FROM message WHERE created_at < DATETIME('now', '-' || (thread.expires_at / 1000000000) || ' seconds') ) AND expires_at IS NOT NULL`,
	)
	if err != nil {
		return fmt.Errorf("delete thread with no new messages since: %w", err)
	}
	return nil
}

func (s *Store) GetExpiredThreads() ([]model.Thread, error) {
	var thread []model.Thread
	err := s.db.Select(
		&thread, `SELECT * FROM thread WHERE id IN ( SELECT thread_id FROM message WHERE created_at < DATETIME('now', '-' || (thread.expires_at / 1000000000) || ' seconds') ) AND expires_at IS NOT NULL`,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get expired threads: %w", err)
	}
	return thread, nil
}

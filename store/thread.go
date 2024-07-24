package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/secret-chat/model"
)

func (s *Store) CreateThread(thread model.Thread) (*int, error) {
	row := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO thread (name, description, user_id, slack_id, key) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		thread.Name,
		thread.Description,
		thread.UserID,
		thread.SlackID,
		thread.Key,
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

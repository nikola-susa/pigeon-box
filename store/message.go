package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/secret-chat/model"
)

func (s *Store) CreateMessage(message model.CreateMessageParams) (*int, error) {
	row := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO message (text, user_id, thread_id) VALUES ($1, $2, $3) RETURNING id`,
		message.Text,
		message.UserID,
		message.ThreadID,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert message: %w", err)
	}

	return &id, nil
}

func (s *Store) GetMessage(messageID int) (*model.Message, error) {
	var message model.Message
	err := s.db.Get(&message, `SELECT id, thread_id, user_id, file_id, text, created_at FROM message WHERE id = $1`, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select message: %w", err)
	}
	return &message, nil
}

func (s *Store) GetMessagesByThread(threadId int) ([]model.Message, error) {
	var messages []model.Message
	err := s.db.Select(&messages, `SELECT id, thread_id, user_id, file_id, text, created_at FROM message WHERE thread_id = $1`, threadId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select messages: %w", err)
	}
	return messages, nil
}

func (s *Store) SetMessageFileID(messageID int, fileID int) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`UPDATE message SET file_id = $1 WHERE id = $2`,
		fileID,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("update message: %w", err)
	}

	return nil
}

func (s *Store) DeleteMessage(messageID int) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`DELETE FROM message WHERE id = $1`,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}

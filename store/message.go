package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/pigeon-box/model"
	"time"
)

func (s *Store) CreateMessage(message model.CreateMessageParams) (*int, error) {
	row := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO message (text, user_id, thread_id) VALUES (?, ?, ?) RETURNING id`,
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

func (s *Store) UpdateMessage(message model.UpdateMessageParams) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`UPDATE message SET text = ?, updated_at = DATETIME('now') WHERE id = ?`,
		message.Text,
		message.ID,
	)
	if err != nil {
		return fmt.Errorf("update message: %w", err)
	}

	return nil
}

func (s *Store) GetMessage(messageID int) (*model.Message, error) {
	var message model.Message
	err := s.db.Get(&message, `SELECT id, thread_id, user_id, file_id, text, created_at, updated_at FROM message WHERE id = ?`, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select message: %w", err)
	}
	return &message, nil
}

func (s *Store) GetMessagesByThread(threadId int, lastId *int) ([]model.Message, error) {
	var messages []model.Message
	var err error
	if lastId != nil {
		err = s.db.Select(&messages, `SELECT id, thread_id, user_id, file_id, text, created_at, updated_at FROM message WHERE thread_id = ? AND id < ? ORDER BY created_at DESC LIMIT 25`, threadId, *lastId)
	} else {
		err = s.db.Select(&messages, `SELECT id, thread_id, user_id, file_id, text, created_at, updated_at FROM message WHERE thread_id = ? ORDER BY created_at DESC LIMIT 25`, threadId)
	}
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
		`UPDATE message SET file_id = ? WHERE id = ?`,
		fileID,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("update message: %w", err)
	}

	return nil
}

func (s *Store) SetMessageExpiresAt(messageID int, expireAt time.Time) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`UPDATE message SET expires_at = ? WHERE id = ?`,
		expireAt,
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
		`DELETE FROM message WHERE id = ?`,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}

func (s *Store) DeleteExpiredMessages() error {
	_, err := s.db.ExecContext(
		context.Background(),
		`DELETE FROM message WHERE expires_at > DATETIME('now')`,
	)
	if err != nil {
		return fmt.Errorf("delete expired messages: %w", err)
	}

	return nil
}

func (s *Store) GetExpiredMessages() ([]model.Message, error) {
	var messages []model.Message
	err := s.db.Select(&messages, `SELECT id, thread_id, user_id, file_id, text, created_at, updated_at FROM message WHERE expires_at > DATETIME('now')`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select expired messages: %w", err)
	}
	return messages, nil
}

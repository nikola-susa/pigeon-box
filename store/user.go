package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/pigeon-box/model"
)

func (s *Store) CreateUser(user model.CreateUser) (*int, error) {

	row := s.db.QueryRow(
		`INSERT INTO user (name, username, slack_id, avatar) VALUES (?, ?, ?, ?) RETURNING id`,
		user.Name,
		user.Username,
		user.SlackID,
		user.Avatar,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &id, nil
}

func (s *Store) GetUser(id int) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(
		context.Background(),
		&user,
		`SELECT * FROM user WHERE id = ?`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (s *Store) GetUserBySlackID(slackID string) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(
		context.Background(),
		&user,
		`SELECT * FROM user WHERE slack_id = ?`,
		slackID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by slack_id: %w", err)
	}
	return &user, nil
}

func (s *Store) UpdateUser(user model.CreateUser) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`UPDATE user SET slack_id = ?, name = ?, username = ?, avatar = ? WHERE id = ?`,
		user.SlackID,
		user.Name,
		user.Username,
		user.Avatar,
		user.SlackID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

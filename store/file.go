package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nikola-susa/pigeon-box/model"
)

func (s *Store) CreateFile(file model.File) (*int, error) {
	row := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO file (name, path, size, content_type, thread_id, user_id) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`,
		file.Name,
		file.Path,
		file.Size,
		file.ContentType,
		file.ThreadID,
		file.UserID,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert file: %w", err)
	}
	return &id, nil
}

func (s *Store) GetFile(id int) (*model.File, error) {
	var file model.File
	err := s.db.GetContext(
		context.Background(),
		&file,
		`SELECT * FROM file WHERE id = ?`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get file by id: %w", err)
	}
	return &file, nil
}

func (s *Store) GetFilesByThread(threadID int) ([]model.File, error) {
	var files []model.File
	err := s.db.SelectContext(
		context.Background(),
		&files,
		`SELECT * FROM file WHERE thread_id = ?`,
		threadID,
	)
	if err != nil {
		return nil, fmt.Errorf("get files by thread: %w", err)
	}
	return files, nil
}

func (s *Store) GetFileByPath(path string) (*model.File, error) {
	var file model.File
	err := s.db.GetContext(
		context.Background(),
		&file,
		`SELECT * FROM file WHERE path = ?`,
		path,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get file by path: %w", err)
	}
	return &file, nil
}

func (s *Store) DeleteFile(id int) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`DELETE FROM file WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

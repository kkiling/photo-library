package sqlitestorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DSN string
}

type Storage struct {
	config Config
	db     *sql.DB
}

func NewStorage(config Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, err
	}
	return &Storage{config: config, db: db}, nil
}

func (s *Storage) GetFileHash(ctx context.Context, filePath string, updateAt time.Time) (*string, error) {
	var hash string
	err := s.db.QueryRowContext(ctx, "SELECT hash FROM file_hash WHERE file_path = ? AND update_at = ?", filePath, updateAt).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &hash, nil
}

func (s *Storage) SaveFileHash(ctx context.Context, filePath, hash string, updateAt time.Time) error {
	_, err := s.db.ExecContext(ctx, "INSERT OR REPLACE INTO file_hash (file_path, update_at, hash) VALUES (?, ?, ?)", filePath, updateAt, hash)
	return err
}

func (s *Storage) SaveUploadFileResponse(ctx context.Context, hash string, uploadDate time.Time, success bool) error {
	_, err := s.db.ExecContext(ctx, "INSERT OR REPLACE INTO file_upload (hash, upload_at, success) VALUES (?, ?, ?)", hash, uploadDate, success)
	return err
}

func (s *Storage) FileAlreadyUpload(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT count(1) FROM file_upload WHERE hash = ? AND success = ?", hash, true).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

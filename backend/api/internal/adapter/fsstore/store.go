package fsstore

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseFilesDir string `yaml:"base_files_dir"`
}

type Store struct {
	cfg Config
}

func NewStore(cfg Config) *Store {
	return &Store{
		cfg: cfg,
	}
}

func (f *Store) SaveFileBody(ctx context.Context, ext string, body []byte) (filePath string, err error) {
	// Формируем новое имя файла
	filePath = fmt.Sprintf("%s/%s.%s", f.cfg.BaseFilesDir, uuid.New(), strings.ToLower(ext))

	// Создаем новый файл с новым именем
	newFile, err := os.Create(filePath)
	defer newFile.Close()

	if err != nil {
		return "", fmt.Errorf("failed to create new file: %w", err)
	}

	// Записываем данные в новый файл
	if _, err := newFile.Write(body); err != nil {
		return "", fmt.Errorf("failed to write to new file: %w", err)
	}

	return filePath, nil
}

func (f *Store) DeleteFile(ctx context.Context, filePath string) error {

	// Construct the absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	// Clean the path to get canonical path
	cleanPath := filepath.Clean(absPath)

	// Construct the absolute path for BaseFilesDir
	baseDir, err := filepath.Abs(f.cfg.BaseFilesDir)
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	// Clean the base directory path
	cleanBaseDir := filepath.Clean(baseDir)

	// Check if the filePath is inside BaseFilesDir
	if !strings.HasPrefix(cleanPath, cleanBaseDir) || cleanPath == cleanBaseDir {
		return fmt.Errorf("file %s is not in directory %s", filePath, f.cfg.BaseFilesDir)
	}

	// Delete the file
	if err := os.Remove(absPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %v", filePath, err)
	}

	return nil
}

func (f *Store) GetFileUrl(ctx context.Context, filepath string) error {
	//TODO implement me
	panic("implement me")
}

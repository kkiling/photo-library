package fsstore

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseFilesDir   string `yaml:"base_files_dir"`
	PhotoServerUrl string `yaml:"photo_server_url"`
}

type Store struct {
	cfg Config
}

func NewStore(cfg Config) *Store {
	return &Store{
		cfg: cfg,
	}
}

// SaveFile сохранение файла на диске
// fileName - имя файла например 123.jpg
// body - данные файла
// dirs - список каталогов в которых будет лежать файл
// fileKey - возвращает путь до файла относительно BaseFilesDir (dirs+fileName)
func (f *Store) SaveFile(_ context.Context, fileName string, body []byte, dirs ...string) (fileKey string, err error) {
	if _, err = os.Stat(f.cfg.BaseFilesDir); os.IsNotExist(err) {
		return "", fmt.Errorf("baseFileDir %s does not exist", f.cfg.BaseFilesDir)
	}

	fileKey = filepath.Join(dirs...)
	fullFilePath := filepath.Join(f.cfg.BaseFilesDir, fileKey)

	if err = os.MkdirAll(fullFilePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create directory %s: %w", fullFilePath, err)
	}

	// Формируем новое имя файла
	fileKey = filepath.Join(fileKey, fileName)
	fullFilePath = filepath.Join(fullFilePath, fileName)

	// Создаем новый файл с новым именем
	newFile, err := os.Create(fullFilePath)
	defer newFile.Close()

	if err != nil {
		return "", fmt.Errorf("failed to create new file: %w", err)
	}

	// Записываем данные в новый файл
	if _, err = newFile.Write(body); err != nil {
		return "", fmt.Errorf("failed to write to new file: %w", err)
	}

	return fileKey, nil
}

func (f *Store) DeleteFile(_ context.Context, fileKey string) error {
	filePath := filepath.Join(f.cfg.BaseFilesDir, fileKey)

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
	if err = os.Remove(absPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %v", filePath, err)
	}

	return nil
}

func (f *Store) DeleteFiles(ctx context.Context, fileKeys []string) error {
	for _, fileKey := range fileKeys {
		err := f.DeleteFile(ctx, fileKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Store) DeleteDirectory(_ context.Context, dirs ...string) error {
	dir := filepath.Join(dirs...)
	return os.RemoveAll(dir)
}

func (f *Store) GetFileBody(_ context.Context, fileKey string) ([]byte, error) {
	filePath := filepath.Join(f.cfg.BaseFilesDir, fileKey)
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create new file: %w", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create new file: %w", err)
	}

	return body, nil
}

func (f *Store) GetFileUrl(_ context.Context, fileKey string) string {
	fileName := filepath.Base(fileKey)
	return filepath.Join(f.cfg.PhotoServerUrl, fileName)
}

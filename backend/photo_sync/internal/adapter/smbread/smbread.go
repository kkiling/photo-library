package smbread

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter"
	"io"
	"net"
	"strings"
)

type Config struct {
	User       string
	Password   string
	Address    string
	ShareName  string
	DirPath    string
	Extensions []string
}

type SmbRead struct {
	config  Config
	session *smb2.Session
	share   *smb2.Share
}

func NewSmbRead(cfg Config) *SmbRead {
	return &SmbRead{
		config: cfg,
		share:  nil,
	}
}

func (s *SmbRead) readDir(fs *smb2.Share, dirPath string, files chan<- adapter.FileInfo) error {
	dir, err := fs.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		newPath := dirPath + "/" + name
		if entry.IsDir() {
			err = s.readDir(fs, newPath, files)
			if err != nil {
				return err
			}
		} else {
			for _, ext := range s.config.Extensions {
				if strings.HasSuffix(strings.ToLower(name), ext) {
					// Получение атрибутов файла
					info, err := fs.Stat(newPath)
					if err != nil {
						return fmt.Errorf("failed to stat file: %v", err)
					}
					files <- adapter.FileInfo{
						FilePath: newPath,
						UpdateAt: info.ModTime(),
					}
				}
			}
		}
	}

	return nil
}

func (s *SmbRead) Connect(ctx context.Context) error {
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.config.User,
			Password: s.config.Password,
		},
	}

	conn, err := net.Dial("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("fail net.Dial: %w", err)
	}

	session, err := dialer.DialContext(ctx, conn)
	if err != nil {
		return fmt.Errorf("dialer.DialContex: %w", err)
	}

	sessionWithContext := session.WithContext(ctx)
	share, err := sessionWithContext.Mount(s.config.ShareName)
	if err != nil {
		_ = session.Logoff()
		return fmt.Errorf("sessionWithContext.Mount: %w", err)
	}

	s.session = session
	s.share = share

	return nil
}

func (s *SmbRead) Disconnect() error {
	err := s.share.Umount()

	if err != nil {
		return fmt.Errorf("share.Unmount: %w", err)
	}

	err = s.session.Logoff()
	if err != nil {
		return fmt.Errorf("session.Logoff: %w", err)
	}

	return nil
}

func (s *SmbRead) ReadFiles(ctx context.Context, filesChan chan<- adapter.FileInfo) error {
	err := s.readDir(s.share, s.config.DirPath, filesChan)
	if err != nil {
		return fmt.Errorf("fail readDir: %w", err)
	}
	return nil
}

func (s *SmbRead) GetFileHash(ctx context.Context, filepath string) (string, error) {
	file, err := s.share.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("share.Open: %w", err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to copy content to hasher: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (s *SmbRead) GetFileData(ctx context.Context, filepath string) ([]byte, error) {
	file, err := s.share.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Чтение содержимого файла
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

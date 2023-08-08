package smbread

import (
	"context"
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter"
	"golang.org/x/crypto/blake2b"
	"io"
	"net"
	"strings"
	"time"
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

func (s *SmbRead) GetFileUpdateAt(ctx context.Context, filepath string) (time.Time, error) {
	info, err := s.share.Stat(filepath)
	if err != nil {
		return time.Time{}, fmt.Errorf("fail share.Sait: %w", err)
	}

	return info.ModTime(), nil
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
					files <- adapter.FileInfo{
						FilePath: newPath,
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
		err = fmt.Errorf("fail readDir: %w", err)
	}

	return nil
}

func (s *SmbRead) GetFileHash(ctx context.Context, filepath string) (string, error) {
	file, err := s.share.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("share.Open: %w", err)
	}
	defer file.Close()

	hash, err := blake2b.New256(nil) // Используем BLAKE2b с длиной хеша 256 бит.
	if err != nil {
		return "", fmt.Errorf("failed to create blake2b hasher: %w", err)
	}

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to copy content to hasher: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (s *SmbRead) GetFileBody(ctx context.Context, filepath string) ([]byte, error) {
	file, err := s.share.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Чтение содержимого файла
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return body, nil
}

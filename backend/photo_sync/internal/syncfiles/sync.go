package syncfiles

import (
	"context"
	"fmt"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Config struct {
	GrpcServerHost string
}

type FileRead interface {
	ReadFiles(ctx context.Context, filesChan chan<- adapter.FileInfo) error
	GetFileHash(ctx context.Context, filepath string) (string, error)
	GetFileData(ctx context.Context, filepath string) ([]byte, error)
}

type Storage interface {
	GetFileHash(ctx context.Context, filePath string, updateAt time.Time) (*string, error)
	SaveFileHash(ctx context.Context, filePath, hash string, updateAt time.Time) error
}

type SyncPhotos struct {
	fileRead FileRead
	storage  Storage
	cfg      Config
}

func createProgressBar(max int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("Find files"),
	)
}

func NewSyncPhotos(fileRead FileRead, storage Storage, cfg Config) *SyncPhotos {
	return &SyncPhotos{
		fileRead: fileRead,
		storage:  storage,
		cfg:      cfg,
	}
}

func (s *SyncPhotos) getFiles(ctx context.Context) ([]adapter.FileInfo, error) {
	bar := createProgressBar(-1, "Find files")

	var fileReadError error = nil

	files := make([]adapter.FileInfo, 0)
	filesChan := make(chan adapter.FileInfo)

	go func() {
		defer close(filesChan)
		err := s.fileRead.ReadFiles(ctx, filesChan)
		if err != nil {
			fileReadError = fmt.Errorf("fail RecursiveReadFiles: %w", err)
		}
	}()

	for filePath := range filesChan {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		files = append(files, filePath)

		if err := bar.Add(1); err != nil {
			return nil, fmt.Errorf("fail bar add: %w", err)
		}
	}

	if err := bar.Finish(); err != nil {
		return nil, fmt.Errorf("fail bar finish: %w", err)
	}

	return files, fileReadError
}

func (s *SyncPhotos) calculateFilesHash(ctx context.Context, files []adapter.FileInfo) (map[string][]string, error) {
	filesHash := make(map[string][]string)
	bar := createProgressBar(len(files), "Calculate files hash")
	for _, file := range files {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		fileHash := ""
		findHash, err := s.storage.GetFileHash(ctx, file.FilePath, file.UpdateAt)
		if err != nil {
			return nil, fmt.Errorf("storage.GetFileHash: %w", err)
		}

		if findHash == nil {

			hash, err := s.fileRead.GetFileHash(ctx, file.FilePath)
			if err != nil {
				return nil, fmt.Errorf("fileRead.GetFileHash: %w", err)
			}

			if err2 := s.storage.SaveFileHash(ctx, file.FilePath, hash, file.UpdateAt); err2 != nil {
				return nil, fmt.Errorf("storage.SaveFileHash: %w", err2)
			}

			fileHash = hash
		} else {
			fileHash = *findHash
		}

		filesHash[fileHash] = append(filesHash[fileHash], file.FilePath)

		if err := bar.Add(1); err != nil {
			return nil, fmt.Errorf("fail bar add: %w", err)
		}
	}

	if err := bar.Finish(); err != nil {
		return nil, fmt.Errorf("fail bar finish: %w", err)
	}

	return filesHash, nil
}

func (s *SyncPhotos) uploadFiles(ctx context.Context, files map[string][]string) error {
	bar := createProgressBar(len(files), "Upload files")

	// Создание подключения к gRPC серверу
	conn, err := grpc.Dial(s.cfg.GrpcServerHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	client := pbv1.NewPhotosServiceClient(conn)

	for hash, paths := range files {
		response, err := client.CheckHashPhoto(ctx, &pbv1.CheckHashPhotoRequest{Hash: hash})
		if err != nil {
			return fmt.Errorf("failed CheckHashPhotoRequest: %w", err)
		}

		if response.AlreadyUploaded {
			continue
		}

		data, err := s.fileRead.GetFileData(ctx, paths[0])
		if err != nil {
			return fmt.Errorf("fileRead.GetFileData: %w", err)
		}

		res, err := client.UploadPhoto(ctx, &pbv1.UploadPhotoRequest{
			Paths: paths,
			Body:  data,
		})

		if err != nil {
			return fmt.Errorf("failed UploadPhoto: %v", err)
		}

		if !res.Success {
			return fmt.Errorf("UploadPhoto: no success")
		}

		if err := bar.Add(1); err != nil {
			return fmt.Errorf("fail bar add: %w", err)
		}
	}

	if err := bar.Finish(); err != nil {
		return fmt.Errorf("fail bar finish: %w", err)
	}

	return nil
}

func (s *SyncPhotos) Sync(ctx context.Context) error {
	files, err := s.getFiles(ctx)
	if err != nil {
		return fmt.Errorf("fail getFiles: %w", err)
	}

	filesHash, err := s.calculateFilesHash(ctx, files)

	err = s.uploadFiles(ctx, filesHash)
	if err != nil {
		return fmt.Errorf("fail uploadFiles: %w", err)
	}

	return nil
}

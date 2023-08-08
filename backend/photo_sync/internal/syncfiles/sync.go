package syncfiles

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type Config struct {
	GrpcServerHost string
	ClientId       string
	NumWorkers     int
}

type uploadData struct {
	mainPath string
	paths    []string
	updateAt time.Time
	hash     string
}

type FileRead interface {
	ReadFiles(ctx context.Context, filesChan chan<- adapter.FileInfo) error
	GetFileUpdateAt(ctx context.Context, filepath string) (time.Time, error)
	GetFileHash(ctx context.Context, filepath string) (string, error)
	GetFileBody(ctx context.Context, filepath string) ([]byte, error)
}

type Storage interface {
	GetFileHash(ctx context.Context, filePath string, updateAt time.Time) (*string, error)
	SaveFileHash(ctx context.Context, filePath, hash string, updateAt time.Time) error
	SaveUploadFileResponse(ctx context.Context, hash string, uploadDate time.Time, success bool) error
	FileAlreadyUpload(ctx context.Context, hash string) (bool, error)
}

type SyncPhotos struct {
	fileRead FileRead
	storage  Storage
	cfg      Config
	mx       sync.Mutex
}

func createProgressBar(max int, description string) *pb.ProgressBar {
	fmt.Println(description)
	bar := pb.New(max)
	// bar.C(description)
	bar.Start()
	return bar
}

func NewSyncPhotos(fileRead FileRead, storage Storage, cfg Config) *SyncPhotos {
	return &SyncPhotos{
		fileRead: fileRead,
		storage:  storage,
		cfg:      cfg,
		mx:       sync.Mutex{},
	}
}

func (s *SyncPhotos) getFiles(ctx context.Context) ([]adapter.FileInfo, error) {
	bar := createProgressBar(-1, "Find files")
	defer bar.Finish()

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

	for file := range filesChan {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		files = append(files, file)

		bar.Increment()
	}

	return files, fileReadError
}

func (s *SyncPhotos) getUploadData(ctx context.Context, filePath string) (uploadData, error) {
	updateAt, err := s.fileRead.GetFileUpdateAt(ctx, filePath)
	if err != nil {
		return uploadData{}, fmt.Errorf("fileRead.GetFileUpdateAt: %w", err)
	}

	s.mx.Lock()
	findHash, err := s.storage.GetFileHash(ctx, filePath, updateAt)
	s.mx.Unlock()

	if err != nil {
		return uploadData{}, fmt.Errorf("storage.GetFileHash: %w", err)
	}

	if findHash != nil {
		return uploadData{
			paths:    []string{filePath},
			updateAt: updateAt,
			hash:     *findHash,
		}, nil
	}

	hash, err := s.fileRead.GetFileHash(ctx, filePath)
	if err != nil {
		return uploadData{}, fmt.Errorf("fileRead.GetFileHash: %w", err)
	}

	s.mx.Lock()
	if err := s.storage.SaveFileHash(ctx, filePath, hash, updateAt); err != nil {
		return uploadData{}, fmt.Errorf("storage.SaveFileHash: %w", err)
	}
	s.mx.Unlock()

	return uploadData{
		paths:    []string{filePath},
		updateAt: updateAt,
		hash:     hash,
	}, nil
}

func jointFilesWithData(dataByHash map[string]*uploadData, data uploadData) map[string]*uploadData {
	v, ok := dataByHash[data.hash]
	if !ok {
		dataByHash[data.hash] = &uploadData{
			mainPath: data.paths[0],
			paths:    data.paths,
			updateAt: data.updateAt,
			hash:     data.hash,
		}
		return dataByHash
	}

	v.paths = append(v.paths, data.paths[0])
	v.updateAt = data.updateAt
	v.mainPath = data.paths[0]

	return dataByHash
}

func (s *SyncPhotos) getUploadDataList(ctx context.Context, files []adapter.FileInfo) ([]uploadData, error) {
	bar := createProgressBar(len(files), "Calculate files hash")
	defer bar.Finish()

	dataByHash := make(map[string]*uploadData)
	filesChan := make(chan adapter.FileInfo)
	errorsChan := make(chan error, s.cfg.NumWorkers)

	go func() {
		defer close(filesChan)
		for _, file := range files {
			filesChan <- file
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < s.cfg.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesChan {
				if ctx.Err() != nil {
					errorsChan <- ctx.Err()
					return
				}

				data, err := s.getUploadData(ctx, file.FilePath)
				if err != nil {
					errorsChan <- fmt.Errorf("fail getFileHash: %w", err)
				}

				s.mx.Lock()
				dataByHash = jointFilesWithData(dataByHash, data)
				s.mx.Unlock()

				bar.Increment()
			}
		}()
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		// returning the first error, you might want to handle or log all errors
		return nil, <-errorsChan
	}

	result := make([]uploadData, 0, len(dataByHash))
	for _, data := range dataByHash {
		result = append(result, *data)
	}

	return result, nil
}

func (s *SyncPhotos) uploadFiles(ctx context.Context, uploadDataList []uploadData) error {
	bar := createProgressBar(len(uploadDataList), "Upload files")
	defer bar.Finish()

	// Создание подключения к gRPC серверу
	conn, err := grpc.Dial(s.cfg.GrpcServerHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	client := pbv1.NewSyncPhotosServiceClient(conn)

	uploadChan := make(chan uploadData)
	errorsChan := make(chan error, len(uploadDataList))

	go func() {
		defer close(uploadChan)
		for _, data := range uploadDataList {
			uploadChan <- data
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < s.cfg.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range uploadChan {
				if ctx.Err() != nil {
					errorsChan <- ctx.Err()
					return
				}
				if alreadyUpload, err := s.storage.FileAlreadyUpload(ctx, data.hash); err != nil {
					errorsChan <- fmt.Errorf("storage.FileAlreadyUpload: %w", err)
					bar.Increment()
					continue
				} else if alreadyUpload {
					bar.Increment()
					continue
				}

				body, err := s.fileRead.GetFileBody(ctx, data.mainPath)
				if err != nil {
					errorsChan <- fmt.Errorf("fileRead.GetFileBody: %w", err)
					bar.Increment()
					continue
				}

				res, err := client.UploadPhoto(ctx, &pbv1.UploadPhotoRequest{
					Paths: data.paths,
					Hash:  data.hash,
					Body:  body,
					UpdateAt: &timestamppb.Timestamp{
						Seconds: data.updateAt.Unix(),
					},
				})

				if err := s.storage.SaveUploadFileResponse(ctx, res.Hash, res.UploadedAt.AsTime(), err == nil); err != nil {
					errorsChan <- fmt.Errorf("storage.SaveUploadFileResponse: %w", err)
					bar.Increment()
					continue
				}

				if err != nil {
					errorsChan <- fmt.Errorf("failed UploadPhoto: %v", err)
					bar.Increment()
					continue
				}

				bar.Increment()
			}
		}()
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		for err := range errorsChan {
			fmt.Println(err)
		}
	}

	return nil
}

func (s *SyncPhotos) Sync(ctx context.Context) error {
	files, err := s.getFiles(ctx)
	if err != nil {
		return fmt.Errorf("fail getFiles: %w", err)
	}

	filesHash, err := s.getUploadDataList(ctx, files)

	err = s.uploadFiles(ctx, filesHash)
	if err != nil {
		return fmt.Errorf("fail uploadFiles: %w", err)
	}

	return nil
}

package syncfiles

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/model"
)

var ErrFileIsEmptyBody = fmt.Errorf("is empty body")

type Config struct {
	NumWorkers int
}

type FileRead interface {
	ReadFiles(ctx context.Context, filesChan chan<- model.FileInfo) error
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

type UploadClient interface {
	UploadPhoto(ctx context.Context, uploadData model.UploadData, body []byte) (model.UploadResult, error)
}

type SyncPhotos struct {
	fileRead FileRead
	storage  Storage
	cfg      Config
	mx       sync.Mutex
	client   UploadClient
}

func createProgressBar(max int, description string) *pb.ProgressBar {
	fmt.Println(description)
	bar := pb.New(max)
	// bar.C(description)
	bar.Start()
	return bar
}

func NewSyncPhotos(fileRead FileRead, storage Storage, client UploadClient, cfg Config) *SyncPhotos {
	return &SyncPhotos{
		fileRead: fileRead,
		storage:  storage,
		cfg:      cfg,
		client:   client,
		mx:       sync.Mutex{},
	}
}

func (s *SyncPhotos) getFiles(ctx context.Context) ([]model.FileInfo, error) {
	bar := createProgressBar(-1, "Find files")
	defer bar.Finish()

	var fileReadError error = nil

	files := make([]model.FileInfo, 0)
	filesChan := make(chan model.FileInfo)

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

func (s *SyncPhotos) getUploadData(ctx context.Context, filePath string) (model.UploadData, error) {
	updateAt, err := s.fileRead.GetFileUpdateAt(ctx, filePath)
	if err != nil {
		return model.UploadData{}, fmt.Errorf("fileRead.GetFileUpdateAt: %w", err)
	}

	s.mx.Lock()
	findHash, err := s.storage.GetFileHash(ctx, filePath, updateAt)
	s.mx.Unlock()

	if err != nil {
		return model.UploadData{}, fmt.Errorf("storage.GetFileHash: %w", err)
	}

	if findHash != nil {
		return model.UploadData{
			Paths:    []string{filePath},
			UpdateAt: updateAt,
			Hash:     *findHash,
		}, nil
	}

	hash, err := s.fileRead.GetFileHash(ctx, filePath)
	if err != nil {
		return model.UploadData{}, fmt.Errorf("fileRead.GetFileHash: %w", err)
	}

	s.mx.Lock()
	if err := s.storage.SaveFileHash(ctx, filePath, hash, updateAt); err != nil {
		return model.UploadData{}, fmt.Errorf("storage.SaveFileHash: %w", err)
	}
	s.mx.Unlock()

	return model.UploadData{
		Paths:    []string{filePath},
		UpdateAt: updateAt,
		Hash:     hash,
	}, nil
}

func jointByHash(dataByHash map[string]*model.UploadData, data model.UploadData) map[string]*model.UploadData {
	v, ok := dataByHash[data.Hash]
	if !ok {
		dataByHash[data.Hash] = &model.UploadData{
			MainPath: data.Paths[0],
			Paths:    data.Paths,
			UpdateAt: data.UpdateAt,
			Hash:     data.Hash,
		}
		return dataByHash
	}

	v.Paths = append(v.Paths, data.Paths[0])
	v.UpdateAt = data.UpdateAt
	v.MainPath = data.Paths[0]

	return dataByHash
}

func (s *SyncPhotos) getUploadDataList(ctx context.Context, files []model.FileInfo) ([]model.UploadData, error) {
	bar := createProgressBar(len(files), "Calculate files hash")
	defer bar.Finish()

	dataByHash := make(map[string]*model.UploadData)
	filesChan := make(chan model.FileInfo)
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
					bar.Increment()
					return
				}

				s.mx.Lock()
				dataByHash = jointByHash(dataByHash, data)
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

	result := make([]model.UploadData, 0, len(dataByHash))
	for _, data := range dataByHash {
		result = append(result, *data)
	}

	return result, nil
}

func (s *SyncPhotos) deletingAlreadyDownloadedFiles(ctx context.Context, uploadDataList []model.UploadData) ([]model.UploadData, error) {
	bar := createProgressBar(len(uploadDataList), "Delete already downloaded files")
	defer bar.Finish()

	result := make([]model.UploadData, 0)
	uploadChan := make(chan model.UploadData)
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

				if alreadyUpload, err := s.storage.FileAlreadyUpload(ctx, data.Hash); err != nil {
					errorsChan <- fmt.Errorf("storage.FileAlreadyUpload: %w", err)
					bar.Increment()
					continue
				} else if alreadyUpload {
					bar.Increment()
					continue
				}

				if data.MainPath == "" || len(data.Paths) == 0 {
					// Такого в принципе не должно быть, если это случилось, то это ошибка
					panic(fmt.Errorf("empty path"))
				}

				s.mx.Lock()
				result = append(result, data)
				s.mx.Unlock()
				bar.Increment()
			}
		}()
	}

	wg.Wait()
	close(errorsChan)
	bar.Finish()

	if len(errorsChan) > 0 {
		// returning the first error, you might want to handle or log all errors
		return nil, <-errorsChan
	}

	return result, nil
}

func (s *SyncPhotos) uploadFile(ctx context.Context, data model.UploadData) error {
	body, err := s.fileRead.GetFileBody(ctx, data.MainPath)
	if err != nil {
		return fmt.Errorf("fileRead.GetFileBody: %w", err)
	}

	if len(body) == 0 {
		return ErrFileIsEmptyBody
	}

	res, err := s.client.UploadPhoto(ctx, data, body)
	if err != nil {
		return fmt.Errorf("failed UploadPhoto: %v", err)
	}

	if err := s.storage.SaveUploadFileResponse(ctx, res.Hash, res.UploadedAt, err == nil); err != nil {
		return fmt.Errorf("storage.SaveUploadFileResponse: %w", err)
	}

	return nil
}

func (s *SyncPhotos) uploadFiles(ctx context.Context, uploadDataList []model.UploadData) error {
	bar := createProgressBar(len(uploadDataList), "Upload files")
	defer bar.Finish()

	uploadChan := make(chan model.UploadData)
	errorsChan := make(chan error, len(uploadDataList))
	warningsChan := make(chan error)
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

				if err := s.uploadFile(ctx, data); err != nil {
					if errors.Is(err, ErrFileIsEmptyBody) {
						warningsChan <- fmt.Errorf("empty body: %s", data.MainPath)
					} else {
						errorsChan <- fmt.Errorf("fileRead.GetFileBody: %w", err)
						bar.Increment()
						continue
					}
				}

				bar.Increment()
			}
		}()
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		// returning the first error, you might want to handle or log all errors
		return <-errorsChan
	}

	// Выводим все warning для информации
	for warning := range warningsChan {
		fmt.Println(warning.Error())
	}

	return nil
}

func (s *SyncPhotos) Sync(ctx context.Context) error {
	// Получение списка файлов
	files, err := s.getFiles(ctx)
	if err != nil {
		return fmt.Errorf("getFiles: %w", err)
	}

	// Формируем массив структур для загрузки файлов
	// В том числе для каждого файла рассчитываем хеш
	uploadDataList, err := s.getUploadDataList(ctx, files)
	if err != nil {
		return fmt.Errorf("getUploadDataList: %w", err)
	}

	// Удаление уже загруженных файлов
	uploadDataList, err = s.deletingAlreadyDownloadedFiles(ctx, uploadDataList)
	if err != nil {
		return fmt.Errorf("deletingAlreadyDownloadedFiles: %w", err)
	}

	// Загружаем файлы на сервер
	err = s.uploadFiles(ctx, uploadDataList)
	if err != nil {
		return fmt.Errorf("uploadFiles: %w", err)
	}

	return nil
}

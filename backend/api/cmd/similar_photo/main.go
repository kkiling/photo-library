package main

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/app"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	"sync"
)

// Чтение мета информации фотографий что бы определить какие свойства есть
// Например что бы создать миграцию таблицы со свойствами
func main() {
	var args config.Arguments
	_, err := flags.Parse(&args)
	if err != nil {
		panic(err)
	}

	cfgProvider, err := config.NewProvider(args)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.NewApp(cfgProvider)
	if err := application.Create(ctx); err != nil {
		panic(err)
	}

	similarPhotos := application.GetSimilarPhotos()
	database := application.GetDbAdapter()
	fileStorage := application.GetFileStorage()

	const limit = 1000
	var offset int64
	const maxGoroutines = 1

	countPhotos, err := database.GetPhotosCount(ctx)
	if err != nil {
		panic(err)
	}

	bar := pb.New(int(countPhotos)).Start()
	defer bar.Finish()

	photoChan := make(chan model.Photo)
	var wg sync.WaitGroup

	go func() {
		for offset = 0; offset < countPhotos; offset += limit {
			photos, err := database.GetPaginatedPhotos(ctx, offset, limit)
			if err != nil {
				panic(err)
			}
			for _, photo := range photos {
				photoChan <- photo
			}
		}
		close(photoChan)
	}()

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for photo := range photoChan {
				photoBody, err := fileStorage.GetFileBody(ctx, photo.FilePath)
				if err != nil {
					panic(fmt.Errorf("fileStorage.GetFileBody: %w", err))
				}
				if err = similarPhotos.SavePhotoVector(ctx, photo, photoBody); err != nil {
					application.Logger().Errorf("fail save photo vector: %s - %v", photo.ID, err)
				}
				bar.Increment()
			}
		}()
	}

	wg.Wait()
}

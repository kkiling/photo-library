package main

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/app"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
)

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

	processingPhotos := application.GetProcessingPhotos()
	statuses := []model.PhotoProcessingStatus{
		model.PhotoProcessingNew,
		model.PhotoProcessingExifData,
		model.PhotoProcessingTagsByMeta,
	}
	for _, status := range statuses {
		func(status model.PhotoProcessingStatus, limit int) {
			application.Logger().Infof("startProcessing photos with status %s", status)
			// Производим обработку
			for {
				count, getError := processingPhotos.ProcessingPhotos(ctx, status, limit)
				if getError != nil {
					panic(getError)
				}
				application.Logger().Infof("%d photos with %s status processed", count, status)
				if count == 0 {
					break
				}
			}
		}(status, 1000)
	}
}

package main

import (
	"context"

	"github.com/jessevdk/go-flags"

	"github.com/kkiling/photo-library/backend/api/internal/app"
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

	logger := application.GetLogger()

	logger.Infof("init processing photos")
	if err := processingPhotos.Init(ctx); err != nil {
		panic(err)
	}

	logger.Infof("start processing photos")
	for {
		res, processingErr := processingPhotos.ProcessingPhotos(ctx)
		if processingErr != nil {
			logger.Fatalf("fatal processing error")
		}

		logger.Infof("processing photos (%d/%d/%d)", res.SuccessProcessedPhotos, res.LockProcessedPhotos, res.ErrorProcessedPhotos)
		if res.EOF {
			break
		}
	}
	logger.Infof("stop processing photos")
}

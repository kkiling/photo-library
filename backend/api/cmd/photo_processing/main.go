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

	application.Logger().Infof("init processing photos")
	if err := processingPhotos.Init(ctx); err != nil {
		panic(err)
	}

	application.Logger().Infof("start processing photos")
	for {
		eof, processingErr := processingPhotos.ProcessingPhotos(ctx)
		if processingErr != nil {
			application.Logger().Fatalf("processingPhotos.ProcessingPhotos: %v", processingErr)
		}
		if eof == false {
			break
		}
		application.Logger().Infof("processing photos...")

	}

	application.Logger().Infof("stop processing photos")
}

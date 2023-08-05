package main

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var args config.Arguments
	_, err := flags.Parse(&args)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogger()

	cfgProvider, err := config.NewProvider(args)
	if err != nil {
		panic(err)
	}

	photoService := handler.NewPhotosServiceServer(logger, cfgProvider)

	go func() {
		err = photoService.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Infof("--- shutdown application ---")
		cancel()
	}()

	<-ctx.Done()
	logger.Infof("--- stopped application ---")
	photoService.Stop()
	logger.Infof("--- stop application ---")
}

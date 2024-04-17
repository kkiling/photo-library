package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"

	"github.com/kkiling/photo-library/backend/api/internal/app"
	"github.com/kkiling/photo-library/backend/api/internal/server"
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
	if err = application.Create(ctx); err != nil {
		panic(err)
	}

	log := application.GetLogger()

	srv := server.NewPhotoLibraryServer(
		application.GetLogger(),
		application.GetServerConfig(),
		application.GetPhotoGroupService(),
		application.GetPhotoTagsService(),
		application.GetPhotoMetadataService(),
	)

	go func() {
		err = srv.Start(ctx)
		if err != nil {
			log.Fatalf("fail start app: %v", err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Infof("--- shutdown application ---")
		cancel()
	}()

	<-ctx.Done()
	log.Infof("--- stopped application ---")
	srv.Stop()
	log.Infof("--- stop application ---")
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	go func() {
		err = application.StartPhotosServer(ctx)
		if err != nil {
			application.Logger().Fatalf("fail start app: %v", err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		application.Logger().Infof("--- shutdown application ---")
		cancel()
	}()

	<-ctx.Done()
	application.Logger().Infof("--- stopped application ---")
	application.StopPhotosServer()
	application.Logger().Infof("--- stop application ---")
}

package main

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/app"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
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
		err = application.Start(ctx)
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
	application.Stop()
	application.Logger().Infof("--- stop application ---")
}

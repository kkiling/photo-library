package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"time"
)

func main() {
	var args config.Arguments
	_, err := flags.Parse(&args)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	logger := log.NewLogger()

	cfgProvider, err := config.NewProvider(args)
	if err != nil {
		panic(err)
	}

	photoService := handler.NewPhotosServiceServer(logger, cfgProvider)

	err = photoService.Start(ctx)
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
	photoService.Stop()
	fmt.Println("EndStop")
}

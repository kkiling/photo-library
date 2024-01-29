package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/app"
	"github.com/kkiling/photo-library/backend/api/internal/service/exifphoto"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/systags"
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

	exifPhoto := application.GetExifPhoto()
	metaPhoto := application.GetMetaPhoto()
	sysTagPhoto := application.GetSysTagPhoto()
	database := application.GetDbAdapter()
	fileStorage := application.GetFileStorage()

	const limit = 1000
	var offset int64

	countPhotos, err := database.GetPhotosCount(ctx)
	if err != nil {
		panic(err)
	}

	bar := pb.New(int(countPhotos)).Start()
	defer bar.Finish()

	for offset = 0; offset < countPhotos; offset += limit {

		photos, err := database.GetPaginatedPhotos(ctx, offset, limit)
		if err != nil {
			panic(err)
		}

		for _, photo := range photos {

			func(photo model.Photo) {
				// Считали фото
				photoBody, err := fileStorage.GetFileBody(ctx, photo.FileName)
				if err != nil {
					panic(fmt.Errorf("fileStorage.GetFileBody: %w", err))
				}

				// Сохранили ее exif инфу в базу
				if err = exifPhoto.SavePhotoExifData(ctx, photo, photoBody); err != nil {
					if errors.Is(err, exifphoto.ExifCriticalErr) || errors.Is(err, exifphoto.ExifEOFErr) {
						// return
					} else {
						panic(err)
					}
				}

				// Сохранили мета данные в базу
				if err = metaPhoto.SavePhotoMetaData(ctx, photo, photoBody); err != nil {
					application.Logger().Errorf("fail save photo meta data: %s - %v", photo.ID, err)
				}

				// Сохранили теги фото
				if err = sysTagPhoto.CreateTagByMeta(ctx, photo); err != nil {
					if errors.Is(err, systags.ErrMetaNotFound) {
						return
					} else {
						panic(err)
					}
				}
			}(photo)

			bar.Increment()
		}
	}
}

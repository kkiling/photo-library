package main

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/photo-library/backend/api/internal/app"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
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

	// exifPhoto := application.GetExifPhoto()
	metaPhoto := application.GetMetaPhoto()
	// sysTagPhoto := application.GetSysTagPhoto()
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
				photoBody, err := fileStorage.GetFileBody(ctx, photo.FilePath)
				if err != nil {
					panic(fmt.Errorf("fileStorage.GetFileBody: %w", err))
				}

				/*if err = exifPhoto.SavePhotoExifData(ctx, photo, photoBody); err != nil {
					if errors.Is(err, exifphoto.ExifCriticalErr) || errors.Is(err, exifphoto.ExifEOFErr) {
						return
					} else {
						panic(err)
					}
				}*/

				if err = metaPhoto.SavePhotoMetaData(ctx, photo, photoBody); err != nil {
					application.Logger().Errorf("fail save photo meta data: %s - %v", photo.ID, err)
				}

				/*if err = sysTagPhoto.CreateTagByMeta(ctx, photo); err != nil {
					if errors.Is(err, systags.ErrMetaNotFound) {
						return
					} else {
						panic(err)
					}
				}*/
			}(photo)

			bar.Increment()
		}
	}
}

package server

import (
	"context"

	photogroupshandler "github.com/kkiling/photo-library/backend/api/internal/handler/photos_service/photo_groups_handler"
	photometadatahandler "github.com/kkiling/photo-library/backend/api/internal/handler/photos_service/photo_metadata_handler"
	phototagshandler "github.com/kkiling/photo-library/backend/api/internal/handler/photos_service/photo_tags_handler"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

const photoLibrarySwaggerName = "photo_library"

type PhotoLibraryServer struct {
	*CustomServer
}

func NewPhotoLibraryServer(
	logger log.Logger,
	cfg server.Config,
	photoGroupsService photogroupshandler.PhotoGroupsService,
	photoTagsService phototagshandler.PhotoTagsService,
	photoMetaDataService photometadatahandler.PhotoMetaDataService,
) *PhotoLibraryServer {
	return &PhotoLibraryServer{
		CustomServer: NewCustomServer(
			logger.Named("photos_library_server"),
			cfg,
			photogroupshandler.NewHandler(
				logger.Named("photos_group_handler"),
				photoGroupsService,
			),
			phototagshandler.NewHandler(
				logger.Named("photo_tags_handler"),
				photoTagsService,
			),
			photometadatahandler.NewHandler(
				logger.Named("photo_metadata_handler"),
				photoMetaDataService,
			),
		),
	}
}

func (s *PhotoLibraryServer) Start(ctx context.Context) error {
	return s.CustomServer.Start(ctx, photoLibrarySwaggerName)
}

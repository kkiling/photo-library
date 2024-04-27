package photo_metadata_handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/handler"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

func (p *PhotoMetadataHandler) GetPhotoMetaData(ctx context.Context, request *desc.GetPhotoMetaDataRequest) (*desc.Metadata, error) {
	photoID, err := uuid.Parse(request.PhotoId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("photoID is invalid: %w", err))
	}

	response, err := p.photosService.GetPhotoMetaData(ctx, photoID)

	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoMetaData")
	}

	return mapMetaData(&response), nil
}

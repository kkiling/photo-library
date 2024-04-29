package sync_photo_handler

import (
	"context"

	"github.com/kkiling/photo-library/backend/api/internal/handler"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func (p *SyncPhotoHandler) UploadPhoto(ctx context.Context, request *desc.UploadPhotoRequest) (*desc.UploadPhotoResponse, error) {
	response, err := p.syncPhoto.UploadPhoto(ctx, mapUploadPhotoRequest(request))

	if err != nil {
		return nil, handler.HandleError(err, "p.syncPhoto.UploadPhoto")
	}

	return mapUploadPhotoResponse(&response), nil
}

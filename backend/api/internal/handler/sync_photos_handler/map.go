package sync_photo_handler

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/sync_photos"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func mapUploadPhotoRequest(request *desc.UploadPhotoRequest) *sync_photos.SyncPhotoRequest {
	return &sync_photos.SyncPhotoRequest{
		Paths:    request.Paths,
		Hash:     request.Hash,
		Body:     request.Body,
		UpdateAt: request.UpdateAt.AsTime(),
		ClientId: request.ClientId,
	}
}

func mapUploadPhotoResponse(response *sync_photos.SyncPhotoResponse) *desc.UploadPhotoResponse {
	return &desc.UploadPhotoResponse{
		HasBeenUploadedBefore: response.HasBeenUploadedBefore,
		Hash:                  response.Hash,
	}
}

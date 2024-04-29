package sync_photo_handler

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func mapUploadPhotoRequest(request *desc.UploadPhotoRequest) *model.SyncPhotoRequest {
	return &model.SyncPhotoRequest{
		Paths:          request.Paths,
		Hash:           request.Hash,
		Body:           request.Body,
		PhotoUpdatedAt: request.PhotoUpdatedAt.AsTime(),
		ClientInfo:     request.ClientInfo,
	}
}

func mapUploadPhotoResponse(response *model.SyncPhotoResponse) *desc.UploadPhotoResponse {
	return &desc.UploadPhotoResponse{
		HasBeenUploadedBefore: response.HasBeenUploadedBefore,
		Hash:                  response.Hash,
	}
}

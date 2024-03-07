package mapper

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/syncphotos"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func UploadPhotoRequest(request *desc.UploadPhotoRequest) *syncphotos.SyncPhotoRequest {
	return &syncphotos.SyncPhotoRequest{
		Paths:    request.Paths,
		Hash:     request.Hash,
		Body:     request.Body,
		UpdateAt: request.UpdateAt.AsTime(),
		ClientId: request.ClientId,
	}
}

func UploadPhotoResponse(response *syncphotos.SyncPhotoResponse) *desc.UploadPhotoResponse {
	return &desc.UploadPhotoResponse{
		HasBeenUploadedBefore: response.HasBeenUploadedBefore,
		Hash:                  response.Hash,
	}
}

package mapper

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UploadPhotoRequest(request *desc.UploadPhotoRequest) *model.SyncPhotoRequest {
	return &model.SyncPhotoRequest{
		Paths:    request.Paths,
		Hash:     request.Hash,
		Body:     request.Body,
		UpdateAt: request.UpdateAt.AsTime(),
		ClientId: request.ClientId,
	}
}

func UploadPhotoResponse(response *model.SyncPhotoResponse) *desc.UploadPhotoResponse {
	return &desc.UploadPhotoResponse{
		HasBeenUploadedBefore: response.HasBeenUploadedBefore,
		Hash:                  response.Hash,
		UploadedAt: &timestamppb.Timestamp{
			Seconds: response.UploadAt.Unix(),
		},
	}
}
package photo_tags_handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kkiling/photo-library/backend/api/internal/handler"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

func (p *PhotoTagsHandler) GetPhotoTags(ctx context.Context, request *desc.GetPhotoTagsRequest) (*desc.GetPhotoTagsResponse, error) {
	photoID, err := uuid.Parse(request.PhotoId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("photoID is invalid: %w", err))
	}

	response, err := p.photosService.GetPhotoTags(ctx, photoID)

	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoMetaData")
	}

	return &desc.GetPhotoTagsResponse{
		Tags: mapTags(response),
	}, nil
}

func (p *PhotoTagsHandler) AddPhotoTag(ctx context.Context, request *desc.AddPhotoTagRequest) (*emptypb.Empty, error) {
	photoID, err := uuid.Parse(request.PhotoId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("photoID is invalid: %w", err))
	}
	categoryID, err := uuid.Parse(request.CategoryId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("categoryID is invalid: %w", err))
	}

	err = p.photosService.AddPhotoTag(ctx, photoID, categoryID, request.GetTagName())
	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.AddPhotoTag")
	}

	return &emptypb.Empty{}, nil
}

func (p *PhotoTagsHandler) DeletePhotoTag(ctx context.Context, request *desc.DeletePhotoTagRequest) (*emptypb.Empty, error) {
	photoID, err := uuid.Parse(request.PhotoId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("photoID is invalid: %w", err))
	}

	err = p.photosService.DeletePhotoTag(ctx, photoID)
	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.DeletePhotoTag")
	}

	return &emptypb.Empty{}, nil
}

func (p *PhotoTagsHandler) GetTagCategories(ctx context.Context, request *desc.GetTagCategoriesRequest) (*desc.GetTagCategoriesResponse, error) {
	response, err := p.photosService.GetTagCategories(ctx)

	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoGroups")
	}

	return mapGetTagCategories(response), nil
}

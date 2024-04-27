package photo_groups_handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kkiling/photo-library/backend/api/internal/handler"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

func (p *PhotoGroupsHandler) GetPhotoGroup(ctx context.Context, request *desc.GetPhotoGroupRequest) (*desc.GetPhotoGroupResponse, error) {
	groupID, err := uuid.ParseBytes([]byte(request.GroupId))
	if err != nil {
		return nil, server.ErrInvalidArgument(err)
	}

	response, err := p.photosService.GetPhotoGroup(ctx, groupID)
	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoGroup")
	}

	return &desc.GetPhotoGroupResponse{
		Item: mapPhotoGroup(&response),
	}, nil
}

func (p *PhotoGroupsHandler) GetPhotoGroups(ctx context.Context, request *desc.GetPhotoGroupsRequest) (*desc.GetPhotoGroupsResponse, error) {
	response, err := p.photosService.GetPhotoGroups(ctx, mapGetPhotoGroupsRequest(request))

	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoGroups")
	}

	return mapGetPhotoGroupsResponse(&response), nil
}

func (p *PhotoGroupsHandler) SetMainPhotoGroup(ctx context.Context, request *desc.SetMainPhotoGroupRequest) (*emptypb.Empty, error) {
	groupID, err := uuid.Parse(request.GroupId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("groupID is invalid: %w", err))
	}
	mainPhotoID, err := uuid.Parse(request.MainPhotoId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("mainPhotoID is invalid: %w", err))
	}

	err = p.photosService.SetMainPhotoGroup(ctx, groupID, mainPhotoID)
	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoMetaData")
	}

	return &emptypb.Empty{}, nil
}

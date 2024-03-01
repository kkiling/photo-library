package syncclient

import (
	"context"
	"fmt"

	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	client    pbv1.SyncPhotosServiceClient
	clientID  string
	accessKey string
}

func NewClient(client pbv1.SyncPhotosServiceClient, clientID string, accessKey string) *Client {
	return &Client{
		client:    client,
		clientID:  clientID,
		accessKey: accessKey,
	}
}

func (c *Client) UploadPhoto(ctx context.Context, data model.UploadData, body []byte) (model.UploadResult, error) {
	res, err := c.client.UploadPhoto(ctx, &pbv1.UploadPhotoRequest{
		Paths: data.Paths,
		Hash:  data.Hash,
		Body:  body,
		UpdateAt: &timestamppb.Timestamp{
			Seconds: data.UpdateAt.Unix(),
		},
		ClientId:  c.clientID,
		AccessKey: c.accessKey,
	})
	if err != nil {
		return model.UploadResult{}, fmt.Errorf("client.UploadPhoto: %w", err)
	}
	return model.UploadResult{
		HasBeenUploadedBefore: res.HasBeenUploadedBefore,
		Hash:                  res.Hash,
	}, nil
}

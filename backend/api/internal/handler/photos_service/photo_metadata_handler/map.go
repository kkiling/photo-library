package photo_metadata_handler

import (
	"github.com/kkiling/photo-library/backend/api/internal/handler/map_utils"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func mapMetaData(metadata *model.PhotoMetadata) *desc.Metadata {
	var geo *desc.Geo
	if metadata.Geo != nil {
		geo = &desc.Geo{
			Latitude:  metadata.Geo.Latitude,
			Longitude: metadata.Geo.Longitude,
		}
	}

	return &desc.Metadata{
		ModelInfo:   metadata.ModelInfo,
		SizeBytes:   int32(metadata.SizeBytes),
		WidthPixel:  int32(metadata.WidthPixel),
		HeightPixel: int32(metadata.HeightPixel),
		DataTime:    map_utils.ToTimestampPtr(metadata.DateTime),
		UpdatedAt:   map_utils.ToTimestamp(metadata.PhotoUpdatedAt),
		Geo:         geo,
	}
}

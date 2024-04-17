package model

import (
	"time"

	"github.com/google/uuid"
)

type ProcessingType string

const (
	ExifDataProcessing           ProcessingType = "EXIF_DATA"
	MetaDataProcessing           ProcessingType = "META_DATA"
	CatalogTagsProcessing        ProcessingType = "CATALOG_TAGS"
	MetaTagsProcessing           ProcessingType = "META_TAGS"
	PhotoVectorProcessing        ProcessingType = "PHOTO_VECTOR"
	SimilarCoefficientProcessing ProcessingType = "SIMILAR_COEFFICIENT"
	PhotoGroupProcessing         ProcessingType = "PHOTO_GROUP"
	PhotoPreviewProcessing       ProcessingType = "PHOTO_PREVIEW"
)

var ProcessingTypes = []ProcessingType{
	ExifDataProcessing,
	MetaDataProcessing,
	CatalogTagsProcessing,
	MetaTagsProcessing,
	PhotoVectorProcessing,
	SimilarCoefficientProcessing,
	PhotoPreviewProcessing,
	PhotoGroupProcessing,
}

type PhotoProcessingResult struct {
	EOF                    bool
	SuccessProcessedPhotos int
	ErrorProcessedPhotos   int
	LockProcessedPhotos    int
}

type PhotoProcessing struct {
	PhotoID        uuid.UUID
	ProcessedAt    time.Time
	ProcessingType ProcessingType
	Success        bool
}

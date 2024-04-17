// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package photo_library

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PhotoExtension string

const (
	PhotoExtensionJPEG PhotoExtension = "JPEG"
	PhotoExtensionPNG  PhotoExtension = "PNG"
)

func (e *PhotoExtension) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PhotoExtension(s)
	case string:
		*e = PhotoExtension(s)
	default:
		return fmt.Errorf("unsupported scan type for PhotoExtension: %T", src)
	}
	return nil
}

type NullPhotoExtension struct {
	PhotoExtension PhotoExtension
	Valid          bool // Valid is true if PhotoExtension is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPhotoExtension) Scan(value interface{}) error {
	if value == nil {
		ns.PhotoExtension, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PhotoExtension.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPhotoExtension) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PhotoExtension), nil
}

type PhotoStatus string

const (
	PhotoStatusACTIVE   PhotoStatus = "ACTIVE"
	PhotoStatusNOTVALID PhotoStatus = "NOT_VALID"
)

func (e *PhotoStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PhotoStatus(s)
	case string:
		*e = PhotoStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PhotoStatus: %T", src)
	}
	return nil
}

type NullPhotoStatus struct {
	PhotoStatus PhotoStatus
	Valid       bool // Valid is true if PhotoStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPhotoStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PhotoStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PhotoStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPhotoStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PhotoStatus), nil
}

type ProcessingType string

const (
	ProcessingTypeEXIFDATA           ProcessingType = "EXIF_DATA"
	ProcessingTypeMETADATA           ProcessingType = "META_DATA"
	ProcessingTypeCATALOGTAGS        ProcessingType = "CATALOG_TAGS"
	ProcessingTypeMETATAGS           ProcessingType = "META_TAGS"
	ProcessingTypePHOTOVECTOR        ProcessingType = "PHOTO_VECTOR"
	ProcessingTypeSIMILARCOEFFICIENT ProcessingType = "SIMILAR_COEFFICIENT"
	ProcessingTypePHOTOGROUP         ProcessingType = "PHOTO_GROUP"
	ProcessingTypePHOTOPREVIEW       ProcessingType = "PHOTO_PREVIEW"
)

func (e *ProcessingType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProcessingType(s)
	case string:
		*e = ProcessingType(s)
	default:
		return fmt.Errorf("unsupported scan type for ProcessingType: %T", src)
	}
	return nil
}

type NullProcessingType struct {
	ProcessingType ProcessingType
	Valid          bool // Valid is true if ProcessingType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProcessingType) Scan(value interface{}) error {
	if value == nil {
		ns.ProcessingType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProcessingType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProcessingType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ProcessingType), nil
}

type CoefficientsSimilarPhoto struct {
	PhotoId1    uuid.UUID
	PhotoId2    uuid.UUID
	Coefficient float64
}

type ExifPhotoDatum struct {
	PhotoID uuid.UUID
	Data    []byte
}

type GooseDbVersion struct {
	ID        int
	VersionID int64
	IsApplied bool
	Tstamp    pgtype.Timestamp
}

type MetaPhotoDatum struct {
	PhotoID      uuid.UUID
	ModelInfo    *string
	SizeBytes    int
	WidthPixel   int
	HeightPixel  int
	DateTime     pgtype.Timestamptz
	UpdatedAt    time.Time
	GeoLatitude  *float64
	GeoLongitude *float64
}

type Photo struct {
	ID        uuid.UUID
	FileKey   string
	Hash      string
	UpdatedAt time.Time
	Extension PhotoExtension
	Status    PhotoStatus
	Error     *string
}

type PhotoGroup struct {
	ID          uuid.UUID
	MainPhotoID uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PhotoGroupsPhoto struct {
	PhotoID uuid.UUID
	GroupID uuid.UUID
}

type PhotoLocation struct {
	PhotoID          uuid.UUID
	CreatedAt        time.Time
	GeoLatitude      float64
	GeoLongitude     float64
	FormattedAddress string
	Street           string
	HouseNumber      string
	Suburb           string
	Postcode         string
	State            string
	StateCode        string
	StateDistrict    string
	County           string
	Country          string
	CountryCode      string
	City             string
}

type PhotoPreview struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	FileKey     string
	SizePixel   int
	WidthPixel  int
	HeightPixel int
	Original    bool
}

type PhotoProcessing struct {
	PhotoID     uuid.UUID
	ProcessedAt time.Time
	Type        ProcessingType
	Success     bool
}

type PhotoTag struct {
	ID         uuid.UUID
	CategoryID uuid.UUID
	PhotoID    uuid.UUID
	Name       string
}

type PhotoUploadDatum struct {
	PhotoID  uuid.UUID
	UploadAt time.Time
	Paths    []string
	ClientID string
}

type PhotoVector struct {
	PhotoID uuid.UUID
	Vector  []float64
	Norm    float64
}

type RocketLock struct {
	Key         string
	LockedUntil time.Time
}

type TagCategory struct {
	ID    uuid.UUID
	Type  string
	Color string
}

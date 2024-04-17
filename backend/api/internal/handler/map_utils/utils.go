package map_utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToTimestampPtr(date *time.Time) *timestamppb.Timestamp {
	if date == nil {
		return nil
	}
	return timestamppb.New(*date)
}

func ToTimestamp(date time.Time) *timestamppb.Timestamp {
	return timestamppb.New(date)
}

func ToDatePtr(date *timestamppb.Timestamp) *time.Time {
	if date == nil {
		return nil
	}
	res := date.AsTime()
	return &res
}

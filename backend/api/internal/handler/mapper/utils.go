package mapper

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toTimestampPtr(date *time.Time) *timestamppb.Timestamp {
	if date == nil {
		return nil
	}
	return timestamppb.New(*date)
}

func toTimestamp(date time.Time) *timestamppb.Timestamp {
	return timestamppb.New(date)
}

func toDatePtr(date *timestamppb.Timestamp) *time.Time {
	if date == nil {
		return nil
	}
	res := date.AsTime()
	return &res
}

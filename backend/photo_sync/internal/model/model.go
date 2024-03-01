package model

import "time"

type UploadData struct {
	MainPath string
	Paths    []string
	UpdateAt time.Time
	Hash     string
}

type UploadResult struct {
	HasBeenUploadedBefore bool
	Hash                  string
}

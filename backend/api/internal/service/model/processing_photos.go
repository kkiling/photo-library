package model

type ProcessingPhotos struct {
	EOF                    bool
	SuccessProcessedPhotos int
	ErrorProcessedPhotos   int
	LockProcessedPhotos    int
}

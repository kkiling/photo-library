package photo_preview

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type imageCsv struct {
	mat gocv.Mat
}

func (icv *imageCsv) Size() (int, int) {
	dim := icv.mat.Size()
	return dim[1], dim[0]
}

func (icv *imageCsv) Load(inputData []byte) error {
	var err error
	icv.mat, err = gocv.IMDecode(inputData, gocv.IMReadUnchanged)
	if err != nil {
		return fmt.Errorf("gocv.IMDecode: %w", err)
	}
	return nil
}

func (icv *imageCsv) ToBytes(fileExt model.PhotoExtension) ([]byte, error) {
	ext := gocv.JPEGFileExt
	switch fileExt {
	case model.PhotoExtensionJpeg:
	case model.PhotoExtensionPng:
		ext = gocv.PNGFileExt
	default:
		return nil, fmt.Errorf("unknown fileExt: %s", fileExt)
	}

	byTeBuffer, err := gocv.IMEncode(ext, icv.mat)
	if err != nil {
		return nil, fmt.Errorf("gocv.IMDecode: %w", err)
	}
	defer byTeBuffer.Close()

	src := byTeBuffer.GetBytes()
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst, nil
}

func (icv *imageCsv) Close() {
	_ = icv.mat.Close()
}

//
//func (icv *imageCsv) Crop(left, top, right, bottom int) *imageCsv {
//	croppedMat := icv.mat.Region(image.Rect(left, top, right, bottom))
//	resultMat := croppedMat.Clone()
//	return &imageCsv{mat: resultMat}
//}

func (icv *imageCsv) Resize(width, height int) *imageCsv {
	resizeMat := gocv.NewMat()
	gocv.Resize(icv.mat, &resizeMat, image.Pt(width, height), 0, 0, gocv.InterpolationArea)
	// _ = icv.mat.Close()
	// icv.mat = resizeMat
	return &imageCsv{mat: resizeMat}
}

/*
func (icv *imageCsv) FlipTB() *imageCsv {
	dstMat := gocv.NewMatWithSize(icv.mat.Rows(), icv.mat.Cols(), icv.mat.Type())
	gocv.Flip(icv.mat, &dstMat, 0)
	return &imageCsv{mat: dstMat}
}

func (icv *imageCsv) FlipLR() *imageCsv {
	dstMat := gocv.NewMatWithSize(icv.mat.Rows(), icv.mat.Cols(), icv.mat.Type())
	gocv.Flip(icv.mat, &dstMat, 1)
	return &imageCsv{mat: dstMat}
}*/

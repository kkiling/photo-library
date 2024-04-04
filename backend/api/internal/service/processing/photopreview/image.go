package photopreview

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type ImageCV struct {
	mat gocv.Mat
}

func (icv *ImageCV) Size() (int, int) {
	dim := icv.mat.Size()
	return dim[1], dim[0]
}

func (icv *ImageCV) Load(inputData []byte) error {
	var err error
	icv.mat, err = gocv.IMDecode(inputData, gocv.IMReadUnchanged)
	if err != nil {
		return fmt.Errorf("gocv.IMDecode: %w", err)
	}
	return nil
}

func (icv *ImageCV) ToBytes(fileExt model.PhotoExtension) ([]byte, error) {
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

func (icv *ImageCV) Close() {
	_ = icv.mat.Close()
}

//
//func (icv *ImageCV) Crop(left, top, right, bottom int) *ImageCV {
//	croppedMat := icv.mat.Region(image.Rect(left, top, right, bottom))
//	resultMat := croppedMat.Clone()
//	return &ImageCV{mat: resultMat}
//}

func (icv *ImageCV) Resize(width, height int) *ImageCV {
	resizeMat := gocv.NewMat()
	gocv.Resize(icv.mat, &resizeMat, image.Pt(width, height), 0, 0, gocv.InterpolationArea)
	// _ = icv.mat.Close()
	// icv.mat = resizeMat
	return &ImageCV{mat: resizeMat}
}

/*
func (icv *ImageCV) FlipTB() *ImageCV {
	dstMat := gocv.NewMatWithSize(icv.mat.Rows(), icv.mat.Cols(), icv.mat.Type())
	gocv.Flip(icv.mat, &dstMat, 0)
	return &ImageCV{mat: dstMat}
}

func (icv *ImageCV) FlipLR() *ImageCV {
	dstMat := gocv.NewMatWithSize(icv.mat.Rows(), icv.mat.Cols(), icv.mat.Type())
	gocv.Flip(icv.mat, &dstMat, 1)
	return &ImageCV{mat: dstMat}
}*/

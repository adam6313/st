package source

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/jdeng/goheif"
	"github.com/tyr-tech-team/hawk/status"
)

var imageMap = map[string]Coder{
	"jpg":  new(JPG),
	"jpeg": new(JPEG),
	"png":  new(PNG),
	"heic": new(HEIC),
	"heif": new(HEIC),
	"gif":  new(GIF),
}

// Coder -
type Coder interface {
	// Decode -
	Decode(r io.Reader) (imageSource, error)

	// Encode -
	Encode(decoder imageSource, quality int) ([]byte, error)
}

func NewCorder(extension string) (Coder, error) {
	var (
		err error
		t   = []string{}
	)

	for key := range imageMap {
		t = append(t, key)
	}

	err = status.InvalidParameter.SetServiceCode(status.ServiceStorage).WithDetail(fmt.Sprintf("檔案格式不支持, 請使用%s格式檔案", strings.Join(t, " / "))).Err()

	coder, ok := imageMap[extension]
	if !ok {
		return nil, err
	}

	return coder, nil
}

// imageSource -
type imageSource interface{}

// JPG -
type JPG struct{}

// Decode -
func (j *JPG) Decode(r io.Reader) (imageSource, error) {
	i, _, err := image.Decode(r)
	return i, err
}

// Encode -
func (j *JPG) Encode(imageSource imageSource, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, imageSource.(image.Image), &jpeg.Options{Quality: quality}); err != nil {
		return nil, status.EncodedFailed.SetServiceCode(status.ServiceStorage).WithDetail("轉碼失敗").Err()
	}

	return buf.Bytes(), nil
}

// JPEG -
type JPEG struct{}

// Decode -
func (j *JPEG) Decode(r io.Reader) (imageSource, error) {
	i, _, err := image.Decode(r)
	return i, err
}

// Encode -
func (j *JPEG) Encode(imageSource imageSource, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, imageSource.(image.Image), &jpeg.Options{Quality: quality}); err != nil {
		return nil, status.EncodedFailed.SetServiceCode(status.ServiceStorage).WithDetail("轉碼失敗").Err()
	}

	return buf.Bytes(), nil
}

// PNG -
type PNG struct{}

// Decode -
func (p *PNG) Decode(r io.Reader) (imageSource, error) {
	i, _, err := image.Decode(r)
	return i, err
}

// Encode -
func (j *PNG) Encode(imageSource imageSource, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imageSource.(image.Image)); err != nil {
		return nil, status.EncodedFailed.SetServiceCode(status.ServiceStorage).WithDetail("轉碼失敗").Err()
	}

	return buf.Bytes(), nil
}

// HEIC -
type HEIC struct{}

// Decode -
func (h *HEIC) Decode(r io.Reader) (imageSource, error) {
	i, err := goheif.Decode(r)
	return i, err
}

// Encode -
func (h *HEIC) Encode(imageSource imageSource, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, imageSource.(image.Image), &jpeg.Options{Quality: quality}); err != nil {
		return nil, status.EncodedFailed.SetServiceCode(status.ServiceStorage).WithDetail("轉碼失敗").Err()
	}

	return buf.Bytes(), nil
}

// GIF -
type GIF struct{}

// Decode -
func (g *GIF) Decode(r io.Reader) (imageSource, error) {
	return gif.DecodeAll(r)
}

// Encode -
func (g *GIF) Encode(imageSource imageSource, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := gif.EncodeAll(buf, imageSource.(*gif.GIF)); err != nil {
		return nil, status.EncodedFailed.SetServiceCode(status.ServiceStorage).WithDetail("轉碼失敗").Err()
	}

	return buf.Bytes(), nil
}

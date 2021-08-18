package source

import (
	"image"
	"image/gif"
	"storage/app/domain/model"
	"sync"

	"github.com/andybons/gogif"
	"github.com/disintegration/imaging"
	"github.com/tyr-tech-team/hawk/status"
)

const (
	// ORIGIN -
	ORIGIN = "origin"

	// WEB -
	WEB = "web"

	// MOBILE -
	MOBILE = "mobile"

	// THUMBNAIL -
	THUMBNAIL = "thumbnail"
)

var (
	wg     sync.WaitGroup
	device = map[string]int{
		WEB:       1920,
		MOBILE:    1024,
		THUMBNAIL: 600,
	}
)

// CompressAmount -
func (f *file) CompressDevice() map[string]int {
	return map[string]int{
		"file": 0,
	}
}

// CompressDevice -
func (i *img) CompressDevice() map[string]int {
	switch v := i.imageSource.(type) {
	case image.Image:
		return customWidth(i.isResponsive, v.Bounds().Dx(), i.width)

	case *gif.GIF:
		return customWidth(i.isResponsive, v.Config.Width, i.width)
	}

	return nil
}

// Compression -
func (f *file) Compress(quality, width int, device string) (*model.Transfer, error) {
	return &model.Transfer{
		ID:        f.id,
		Name:      f.name,
		Prefix:    f.prefix,
		Extension: f.Extension(),
		Device:    ORIGIN,
		Data:      f.data,
	}, nil
}

// Compress -
func (i *img) Compress(quality, width int, device string) (*model.Transfer, error) {
	switch v := i.imageSource.(type) {
	case image.Image:
		i.hash, _ = blurhash(v)

		v = orientation(v, i.exif.Orientation)

		var data []byte
		if width == 0 {
			data = i.data
		} else {
			if v.Bounds().Dx() < width {
				width = v.Bounds().Dx()
			}

			i.height = int((float64(width) / float64(v.Bounds().Dx())) * float64(v.Bounds().Dy()))

			dst := imaging.Resize(v, width, i.height, imaging.Lanczos)
			data, _ = i.coder.Encode(dst, quality)
		}

		return &model.Transfer{
			ID:        i.id,
			Name:      i.name,
			Prefix:    i.prefix,
			Extension: i.Extension(),
			Device:    device,
			Hash:      i.hash,
			Data:      data,
		}, nil

	case *gif.GIF:
		i.hash, _ = blurhash(v.Image[0])

		if v.Config.Width < width {
			width = v.Config.Width
		}

		i.height = int((float64(width) / float64(v.Config.Width) * float64(v.Config.Height)))

		outGif := &gif.GIF{}
		for _, m := range v.Image {
			bounds := m.Bounds()

			palettedImage := image.NewPaletted(bounds, nil)
			quantizer := gogif.MedianCutQuantizer{NumColor: 64}
			quantizer.Quantize(palettedImage, bounds, m, image.ZP)

			// Add new frame to animated GIF
			outGif.Image = append(outGif.Image, palettedImage)
			outGif.Delay = append(outGif.Delay, 0)
		}

		data, _ := i.coder.Encode(outGif, quality)
		return &model.Transfer{
			ID:        i.id,
			Name:      i.name,
			Prefix:    i.prefix,
			Extension: i.Extension(),
			Device:    device,
			Hash:      i.hash,
			Data:      data,
		}, nil
	}

	return nil, status.BrandNotFound.SetServiceCode(status.ServiceStorage).WithDetail([]string{"compress failed"}...).Err()
}

func customWidth(isResponsive bool, originWidth, customWidth int) map[string]int {
	var m = map[string]int{}

	// set origin
	m[ORIGIN] = originWidth
	if customWidth != 0 && customWidth < originWidth {
		m[ORIGIN] = customWidth
	}

	if !isResponsive {
		return m
	}

	// set other
	for k, v := range device {
		if m[ORIGIN] < v {
			m[k] = m[ORIGIN]
			continue
		}

		m[k] = v
	}

	return m
}

func orientation(m image.Image, orient string) image.Image {
	switch orient {
	case "3", "4":
		m = imaging.Rotate180(m)
	case "5", "6":
		m = imaging.Rotate270(m)
	case "7", "8":
		m = imaging.Rotate90(m)
	}
	return m
}

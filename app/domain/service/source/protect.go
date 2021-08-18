package source

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

// Protect -
func (f *file) Protect(isProtect bool) error {
	return nil
}

// Protect -
func (i *img) Protect(isProtect bool) error {
	if !isProtect {
		return nil
	}

	var (
		widthGap  = 3
		heightGap = 0
		angle     = 20.0
	)

	switch v := i.imageSource.(type) {
	case image.Image:
		imageBounds := v.Bounds()

		watermark := makeWatermark(imageBounds.Dx()/widthGap, 0, angle)

		heightGap = imageBounds.Dy() / watermark.Bounds().Dy()

		m := image.NewRGBA(imageBounds)
		draw.Draw(m, imageBounds, v, imageBounds.Min, draw.Src)

		// write Watermark
		writeWatermark(widthGap, heightGap, watermark, m)

		i.imageSource = m

	case *gif.GIF:
		for i, p := range v.Image {
			imageBounds := p.Bounds()
			watermark := makeWatermark(imageBounds.Dx()/widthGap, 0, angle)
			heightGap = imageBounds.Dy() / watermark.Bounds().Dy()

			m := image.NewPaletted(imageBounds, p.Palette)

			draw.Draw(m, imageBounds, p, image.ZP, draw.Src)

			// write Watermark
			writeWatermark(widthGap, heightGap, watermark, m)

			v.Image[i] = m
		}

		i.imageSource = v
	}

	return nil
}

// makeWatermark - 建立浮水印
func makeWatermark(width, height int, angle float64) image.Image {
	wmb, _ := os.Open("watermark.png")
	defer wmb.Close()

	watermark, _ := png.Decode(wmb)
	watermark = imaging.Rotate(watermark, angle, color.RGBA{})
	watermark = imaging.Resize(watermark, width, height, imaging.Lanczos)

	return watermark
}

func writeWatermark(widthGap, heightGap int, watermark image.Image, d draw.Image) {
	var dx, dy = watermark.Bounds().Dx(), watermark.Bounds().Dy()

	for i := 0; i < widthGap; i++ {
		for j := 0; j < heightGap; j++ {
			x := ((dx * i) + (dx / widthGap * i)) - (dx / widthGap)
			y := (dy * j) + (dy / heightGap * j)

			offset := image.Pt(x, y)
			draw.Draw(d, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
		}
	}
}

func resizer(img image.Image, width, height, auto int) image.Image {
	if auto != 0 {
		r := img.Bounds()
		if r.Dx() > r.Dy() {
			return imaging.Resize(img, auto, 0, imaging.Lanczos)
		}
		return imaging.Resize(img, 0, auto, imaging.Lanczos)
	}

	return imaging.Resize(img, width, height, imaging.Lanczos)
}

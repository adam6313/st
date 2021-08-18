package source

import (
	"image"

	goblurhash "github.com/buckket/go-blurhash"
)

// Hash -
func (f *file) Hash() string {
	return ""
}

// Hash -
func (i *img) Hash() string {
	return i.hash
}

func blurhash(i image.Image) (string, error) {
	return goblurhash.Encode(4, 3, i)
}

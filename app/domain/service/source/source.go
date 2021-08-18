package source

import (
	"bytes"
	"storage/app/domain/model"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/tyr-tech-team/hawk/status"
)

const (
	// IMAGE -
	IMAGE = "image"
)

var (
	// Device -
	Device = map[string]int{
		"web":       1920,
		"mobile":    1020,
		"thumbnail": 600,
	}
)

// Source -
type Source interface {
	// Extension -
	Extension() string

	// Compress -
	Compress(quality, width int, device string) (*model.Transfer, error)

	// CompressDevice -
	CompressDevice() map[string]int

	// Hash -
	Hash() string

	// Protect -
	Protect(isPrivate bool) error

	// SetID -
	SetID(in string)

	// SetVersionPrefix -
	SetVersionPrefix(prefix string)
}

type file struct {
	// id -
	id string

	// name -
	name string

	// prefix -
	prefix string

	// extension -
	extension string

	// data - source data
	data []byte
}

type img struct {
	// id -
	id string

	// name -
	name string

	// prefix -
	prefix string

	// extension -
	extension string

	// data - source data
	data []byte

	// imageSource -
	imageSource imageSource

	// coder -
	coder Coder

	// width -
	width int

	// height -
	height int

	// hash -
	hash string

	// isResponsive -
	isResponsive bool

	// exif
	exif *exif
}

// NewSource -
func NewSource(in *model.UploadFileRequest) (Source, error) {
	// get file type
	kind, err := filetype.Match(in.Data)
	if err != nil {
		return nil, status.UploadFileFailed.SetServiceCode(status.ServiceStorage).WithDetail([]string{"檔案格式不支持", in.Name}...).Err()
	}

	// check extension is supported
	if err := extensionSupported(kind.Extension, in.Extensions); err != nil {
		return nil, err
	}

	// If unknown, it is not an image file
	switch kind.MIME.Type {
	case IMAGE:
		// newCoder
		coder, err := NewCorder(kind.Extension)
		if err != nil {
			return nil, err
		}

		// reader
		r := bytes.NewReader(in.Data)
		imageSource, _ := coder.Decode(r)

		return &img{
			name:         in.Name,
			prefix:       in.Prefix,
			extension:    kind.Extension,
			data:         in.Data,
			imageSource:  imageSource,
			coder:        coder,
			width:        in.Width,
			height:       in.Height,
			isResponsive: in.IsResponsive,
			exif:         getEXIF(in.Data),
		}, nil
	case types.Unknown.MIME.Type:
		return nil, status.UploadFileFailed.SetServiceCode(status.ServiceStorage).WithDetail([]string{"檔案格式不支持", in.Name}...).Err()
	}

	return &file{
		name:      in.Name,
		prefix:    in.Prefix,
		extension: kind.Extension,
		data:      in.Data,
	}, nil
}

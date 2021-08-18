package service

import (
	"context"
	"storage/app/domain/model"
	"storage/app/domain/service/source"
)

// Service -
type Service interface {
	// NewID -
	NewID() string

	// NewSource -
	NewSource(in *model.UploadFileRequest) (source.Source, error)

	// UploadFile -
	//UploadFile(ctx context.Context, ch <-chan *model.Transfer, devices map[string]int) ([]*model.File, error)
	UploadFile(ctx context.Context, transfer *model.Transfer) (*model.File, error)

	// UploadImageInbackground -
	UploadImageInbackground(ctx context.Context, in *model.File)

	// DownloadFile -
	DownloadFile(url string) ([]byte, error)

	// DeleteFile -
	DeleteFile(ctx context.Context, in []*model.File) error

	// SetFileInfo -
	SetFileInfo(f *model.File, t *model.Transfer)

	// FindFile -
	FindFile(ctx context.Context, fs []*model.File, target string) *model.File
}

// NewSource -
func (s *service) NewSource(in *model.UploadFileRequest) (source.Source, error) {
	return source.NewSource(in)
}

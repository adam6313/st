package usecase

import (
	"context"
	"storage/app/domain/model"
	"storage/app/domain/repository"
	"storage/app/domain/service"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

// StorageUsecase -
type StorageUsecase interface {
	// NewStorageID - 產生新ID
	NewStorageID(ctx context.Context) string

	// UploadFile - 上傳檔案
	UploadFile(ctx context.Context, in *model.UploadFileRequest) (*model.File, error)

	// DeleteFile - 刪除檔案
	DeleteFile(ctx context.Context, in *model.DeleteFileRequest) error
}

// NewStorageUsecase -
func NewStorageUsecase(service service.Service, googleRepo repository.StorageRepository, mongoRepo repository.MongoRepository) StorageUsecase {
	return &storageUsecase{
		service:    service,
		googleRepo: googleRepo,
		mongoRepo:  mongoRepo,
		log: func(ctx context.Context) *logrus.Entry {
			return ctxlogrus.Extract(ctx).WithField("entry", "storageUsecase")
		},
	}
}

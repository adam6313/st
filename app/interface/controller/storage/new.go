package storage

import (
	context "context"
	"storage/app/infra/config"
	"storage/app/usecase"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

type service struct {
	config         config.Config
	storageUsecase usecase.StorageUsecase
	log            func(context.Context) *logrus.Entry
}

// NewService -
func NewService(c config.Config, storageUsecase usecase.StorageUsecase) StorageServiceServer {
	return &service{
		config:         c,
		storageUsecase: storageUsecase,
		log: func(ctx context.Context) *logrus.Entry {
			return ctxlogrus.Extract(ctx).WithField("entry", "storageServer")
		},
	}
}

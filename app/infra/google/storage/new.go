package storage

import (
	"context"
	"storage/app/domain/repository"
	"storage/app/infra/config"

	"cloud.google.com/go/storage"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

type repo struct {
	Domain string
	Bucket string
	client *storage.Client
	log    func(context.Context) *logrus.Entry
}

// NewRepository -
func NewRepository(client *storage.Client, c config.Config) repository.StorageRepository {
	return &repo{
		Domain: c.Google.Domain,
		Bucket: c.Google.Bucket,
		client: client,
		log: func(ctx context.Context) *logrus.Entry {
			return ctxlogrus.Extract(ctx).WithField("entry", "storageRepository")
		},
	}
}

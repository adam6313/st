package storage

import (
	"context"
	"storage/app/domain/repository"
	"storage/app/infra/config"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	C = "storage"
)

type repo struct {
	db           *mongo.Client
	databaseName string
	log          func(context.Context) *logrus.Entry
}

// NewRepository -
func NewRepository(m *mongo.Client, c config.Config) repository.MongoRepository {
	return &repo{
		db:           m,
		databaseName: c.Mongo.Database,
		log: func(ctx context.Context) *logrus.Entry {
			return ctxlogrus.Extract(ctx).WithField("entry", "brandRepository")
		},
	}
}

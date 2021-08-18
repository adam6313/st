package service

import (
	"context"
	"storage/app/domain/repository"

	"github.com/bwmarrin/snowflake"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

type service struct {
	node       *snowflake.Node
	googleRepo repository.StorageRepository
	mongoRepo  repository.MongoRepository
	log        func(ctx context.Context) *logrus.Entry
}

// NewService -
func NewService(googleRepo repository.StorageRepository, mongoRepo repository.MongoRepository) Service {
	node, _ := snowflake.NewNode(1)

	return &service{
		node:       node,
		googleRepo: googleRepo,
		mongoRepo:  mongoRepo,
		log: func(ctx context.Context) *logrus.Entry {
			return ctxlogrus.Extract(ctx).WithField("entry", "service")
		},
	}
}

package repository

import (
	"context"
	"storage/app/domain/model"
)

// StorageRepository -
type StorageRepository interface {
	// ACL -
	ACL(ctx context.Context, object string) error

	// Upload -
	Upload(ctx context.Context, object string, data []byte) error

	// Attrs -
	Attrs(ctx context.Context, object string) (*model.File, error)

	// Delete -
	Delete(ctx context.Context, object string) error

	// GetDomain -
	GetDomain() string
}

// MongoRepository -
type MongoRepository interface {
	// FindByID -
	FindByID(ctx context.Context, id string) ([]*model.File, error)

	// FindByPrefix -
	FindByPrefix(ctx context.Context, prefix string) ([]*model.File, error)

	// GetByURL -
	GetByURL(ctx context.Context, url string) (*model.File, error)

	// Create -
	Create(ctx context.Context, in *model.File) error

	// Create -
	CreateMany(ctx context.Context, in []*model.File) error

	// DeleteByID -
	DeleteByID(ctx context.Context, id string) (*model.File, error)

	// DeleteByURL -
	DeleteByURL(ctx context.Context, url string) (*model.File, error)

	// DeleteManyByPrefix -
	DeleteManyByPrefix(ctx context.Context, prefix string) error
}

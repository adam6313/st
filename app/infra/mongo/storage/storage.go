package storage

import (
	"context"
	"storage/app/domain/model"

	"github.com/tyr-tech-team/hawk/status"
	"go.mongodb.org/mongo-driver/bson"
)

// FindByID -
func (r *repo) FindByID(ctx context.Context, id string) ([]*model.File, error) {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	filter := bson.M{
		"id": id,
	}

	result := []*model.File{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, status.NotFound.SetServiceCode(status.ServiceStorage).Err()

	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, status.NotFound.SetServiceCode(status.ServiceStorage).Err()
	}

	return result, nil
}

// FindByPrefix -
func (r *repo) FindByPrefix(ctx context.Context, prefix string) ([]*model.File, error) {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	filter := bson.M{
		"prefix": prefix,
	}

	result := []*model.File{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, status.NotFound.SetServiceCode(status.ServiceStorage).Err()

	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, status.NotFound.SetServiceCode(status.ServiceStorage).Err()
	}

	return result, nil
}

// GetByURL -
func (r *repo) GetByURL(ctx context.Context, url string) (*model.File, error) {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	filter := bson.M{
		"url": url,
	}

	result := &model.File{}
	if err := coll.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, status.NotFound.SetServiceCode(status.ServiceStorage).Err()
	}

	return result, nil
}

// Create -
func (r *repo) Create(ctx context.Context, in *model.File) error {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	if _, err := coll.InsertOne(ctx, in); err != nil {
		r.log(ctx).Error(err)
		return status.CreatedFailed.SetServiceCode(status.ServiceStorage).Err()
	}

	return nil
}

// CreateMany -
func (r *repo) CreateMany(ctx context.Context, in []*model.File) error {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	docs := make([]interface{}, len(in))
	for i, v := range in {
		docs[i] = v
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		r.log(ctx).Error(err)
		return status.CreatedFailed.SetServiceCode(status.ServiceStorage).Err()
	}

	return nil
}

// DeleteByID -
func (r *repo) DeleteByID(ctx context.Context, id string) (*model.File, error) {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	filter := bson.M{
		"id": id,
	}

	var result *model.File
	if err := coll.FindOneAndDelete(ctx, filter).Decode(&result); err != nil {
		return nil, status.DeletedFailed.SetServiceCode(status.ServiceStorage).Err()
	}

	return result, nil
}

// DeleteByURL -
func (r *repo) DeleteByURL(ctx context.Context, url string) (*model.File, error) {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	filter := bson.M{
		"url": url,
	}

	var result *model.File
	if err := coll.FindOneAndDelete(ctx, filter).Decode(&result); err != nil {
		return nil, status.DeletedFailed.SetServiceCode(status.ServiceStorage).Err()
	}

	return result, nil
}

// DeleteManyByPrefix -
func (r *repo) DeleteManyByPrefix(ctx context.Context, prefix string) error {
	// connection  database and collection
	coll := r.db.Database(r.databaseName).Collection(C)

	filter := bson.M{
		"prefix": prefix,
	}

	if _, err := coll.DeleteMany(ctx, filter); err != nil {
		r.log(ctx).Error(err)
		return status.DeletedFailed.SetServiceCode(status.ServiceStorage).Err()
	}

	return nil
}

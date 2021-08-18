package storage

import (
	"bytes"
	"context"
	"io"

	"storage/app/domain/model"

	"cloud.google.com/go/storage"
	"github.com/tyr-tech-team/hawk/status"
)

// ACL -
func (r *repo) ACL(ctx context.Context, object string) error {
	// set ACL
	acl := r.client.Bucket(r.Bucket).Object(object).ACL()

	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		r.log(ctx).WithError(err).Error("acl.Set failed")
		return status.ConnectFailed.Err()
	}

	return nil
}

// UploadFile -
func (r *repo) Upload(ctx context.Context, object string, data []byte) error {
	wc := r.client.Bucket(r.Bucket).Object(object).NewWriter(ctx)

	if _, err := io.Copy(wc, bytes.NewReader(data)); err != nil {
		r.log(ctx).WithError(err).Error("io.Copy failed")
		return status.UploadFileFailed.Err()
	}

	if err := wc.Close(); err != nil {
		r.log(ctx).WithError(err).Error("wc.Close() failed")
		return status.UploadFileFailed.Err()
	}

	return nil
}

// Attrs -
func (r *repo) Attrs(ctx context.Context, object string) (*model.File, error) {
	o := r.client.Bucket(r.Bucket).Object(object)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		r.log(ctx).WithError(err).Error("get file attrs failed")
		return nil, status.UploadFileNotFound.Err()
	}

	return &model.File{
		Bucket:      attrs.Bucket,
		Name:        attrs.Name,
		Size:        attrs.Size,
		URL:         attrs.MediaLink,
		ContentType: attrs.ContentType,
		CreatedAt:   attrs.Created.Unix(),
	}, nil
}

// DeleteFile -
func (r *repo) Delete(ctx context.Context, object string) error {
	o := r.client.Bucket(r.Bucket).Object(object)

	if err := o.Delete(ctx); err != nil {
		r.log(ctx).WithError(err).Error("delete file failed")
		return status.DeletedFailed.SetServiceCode(status.ServiceStorage).Err()
	}

	return nil
}

// GetDomain -
func (r *repo) GetDomain() string {
	return r.Domain
}

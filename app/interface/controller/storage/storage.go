package storage

import (
	"bytes"
	context "context"
	io "io"
	"storage/app/domain/model"
	"strings"

	"github.com/tyr-tech-team/hawk/status"
)

// NewStorageID -
func (s *service) NewStorageID(ctx context.Context, in *Empty) (*NewStorageIDResponse, error) {
	id := s.storageUsecase.NewStorageID(ctx)

	return &NewStorageIDResponse{
		Id: id,
	}, nil
}

// UploadFile -
func (s *service) UploadFile(stream StorageService_UploadFileServer) error {
	var (
		in   *model.UploadFileRequest
		data *bytes.Buffer
	)

	for {
		// recv data from stream
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}

		// first run
		if in == nil {
			data = bytes.NewBuffer(make([]byte, 0, int(reply.GetFileSize())))

			// set in
			in = &model.UploadFileRequest{
				Prefix:       reply.GetPrefix(),
				Name:         reply.GetName(),
				Extensions:   reply.GetExtensions(),
				Protect:      reply.GetProtect(),
				Width:        int(reply.GetWidth()),
				Height:       int(reply.GetHeight()),
				IsResponsive: reply.GetIsResponsive(),
			}
		}

		// write byte
		if _, err := data.Write(reply.GetChunk()); err != nil {
			s.log(stream.Context()).WithError(err).Error("write bytes.Buffer failed")
			return status.UploadFileFailed.Err()
		}
	}

	// sync data
	in.Data = data.Bytes()

	// upload file
	f, err := s.storageUsecase.UploadFile(stream.Context(), in)

	// send and close stream
	if err := stream.SendAndClose(&UploadFileResponse{
		Id:          f.ID,
		Name:        f.Name,
		Size_:       f.Size,
		Url:         f.URL,
		ContentType: f.ContentType,
		Prefix:      f.Prefix,
		Err:         strings.Join(status.ConvertStatus(err).Detail(), ""),
	}); err != nil {
		s.log(stream.Context()).WithError(err).Error("send and close failed")
		return status.UploadFileFailed.Err()
	}

	return nil
}

// DeleteFile -
func (s *service) DeleteFile(ctx context.Context, in *DeleteFileRequest) (*Empty, error) {
	err := s.storageUsecase.DeleteFile(ctx, &model.DeleteFileRequest{
		ID:     in.GetId(),
		Prefix: in.GetPrefix(),
		URL:    in.GetUrl(),
	})

	return &Empty{}, err
}

package usecase

import (
	"context"
	"io/ioutil"
	"storage/app/domain/model"
	"storage/app/domain/service"
	"storage/app/infra/config"
	"storage/app/infra/google"
	google_storage "storage/app/infra/google/storage"
	"storage/app/infra/mongo"
	mongo_storage "storage/app/infra/mongo/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

var (
	c = dig.New()
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	config.C.Info.RemoteHost = "127.0.0.1:8500"
	config.C.Info.Mode = "production"
	config.C.Google = &config.Google{
		Credentials: "../../second-hand-boutique-2cb7a6e6fef1.json",
	}

	c.Provide(ctx)
	if err := c.Provide(config.NewConsulClient); err != nil {
		panic(err)
	}

	if err := c.Provide(config.RemoteConfig); err != nil {
		panic(err)
	}

	// google
	c.Provide(google.NewDial)
	c.Provide(google_storage.NewRepository)

	c.Provide(mongo.NewDial)
	c.Provide(mongo_storage.NewRepository)

	c.Provide(service.NewService)

	c.Provide(NewStorageUsecase)
}

// TestUpdateFile -
func TestUpdateFile(t *testing.T) {
	c.Invoke(func(usecase StorageUsecase) {

		data, _ := ioutil.ReadFile("6.pdf")

		_, err := usecase.UploadFile(context.Background(), &model.UploadFileRequest{
			Data:         data,
			Protect:      true,
			Prefix:       "testpdf",
			Name:         "51_5184X3456.jpg",
			Extensions:   []string{"png", "jpg", "jpeg", "gif", "heif", "pdf"},
			IsResponsive: true,
			Width:        920,
		})

		assert.NoError(t, err)
	})
}

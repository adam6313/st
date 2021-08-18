package google

import (
	"context"
	"io/ioutil"
	"storage/app/infra/config"

	"storage/app/domain/model"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// NewDial -
func NewDial(c config.Config) (*storage.Client, error) {
	m := new(model.Mode)

	// model verify
	mode := m.ModelVerify(c.Info.Mode)

	// development mode ignore
	switch mode {
	case model.Mode_Development:
		return nil, nil
	}

	// new client(read []byte)
	client, err := storage.NewClient(
		context.Background(),
		option.WithCredentialsJSON(readCredentials(c)),
	)

	if err != nil {
		panic(err)
	}

	return client, nil
}

func readCredentials(c config.Config) []byte {
	if c.Google.Credentials == "" {
		panic("Credentials is required")
	}

	content, err := ioutil.ReadFile(c.Google.Credentials)
	if err != nil {
		panic(err)
	}

	return content
}

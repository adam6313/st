package mongo

import (
	"storage/app/infra/config"

	"github.com/tyr-tech-team/hawk/infra/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewDial -
func NewDial(c config.Config) (*mongo.Client, error) {
	return mongodb.NewDial(c.Mongo)
}

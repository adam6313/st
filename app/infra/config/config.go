package config

import (
	"context"

	"github.com/tyr-tech-team/hawk/config"
	"github.com/tyr-tech-team/hawk/config/source"
	"github.com/tyr-tech-team/hawk/pkg/consul"
	"github.com/tyr-tech-team/hawk/srv"
)

// C -
var C = Config{
	Info: config.Info{
		Name: "storage",
	},
	Google: &Google{},
}

// Config -
type Config struct {
	// Info -
	Info config.Info `yaml:"info"`

	// Mongo -
	Mongo config.Mongo `yaml:"mongo"`

	// Log -
	Log config.Log `yaml:"log"`

	// Google -
	Google *Google `yaml:"google"`

	// Trace -
	Trace Trace `yaml:"trace"`
}

// Trace -
type Trace struct {
	URL         string `yaml:"url"`
	Environment string `yaml:"environment"`
}

// Google -
type Google struct {
	// Credentials - 驗證書
	Credentials string `yaml:"credentials"`

	// Bucket - 區值
	Bucket string `yaml:"bucket"`

	// Domain - storage domain
	Domain string `yaml:"domain"`
}

// NewConsulClient -
func NewConsulClient(ctx context.Context) consul.Client {
	cc := consul.DefaultConsulConfig()
	cc.Address = C.Info.RemoteHost
	cli := consul.NewClient(ctx, cc)
	return cli
}

// RemoteConfig -
func RemoteConfig(cli consul.Client) (Config, error) {
	r := config.NewReader(source.NewConsul(cli, C.Info.Name), config.YAML)
	err := r.ReadWith(&C)

	return C, err
}

// RegisterClient -
func RegisterClient(cli consul.Client) srv.Register {
	return cli
}

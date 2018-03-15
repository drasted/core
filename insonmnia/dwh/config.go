package dwh

import (
	"github.com/jinzhu/configor"
	"github.com/sonm-io/core/accounts"
)

type Config struct {
	ListenAddr        string             `yaml:"address"`
	Eth               accounts.EthConfig `required:"true" yaml:"ethereum"`
	Storage           *storageConfig     `required:"true" yaml:"storage"`
	MetricsListenAddr string             `yaml:"metrics_listen_addr" default:"127.0.0.1:14004"`
}

type storageConfig struct {
	Backend  string `required:"true" yaml:"backend"`
	Endpoint string `required:"true" yaml:"endpoint"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := configor.Load(cfg, path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

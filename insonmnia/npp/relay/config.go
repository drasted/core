package relay

import (
	"fmt"
	"net"
	"os"

	"github.com/jinzhu/configor"
	"github.com/pborman/uuid"
	"github.com/sonm-io/core/insonmnia/logging"
	"go.uber.org/zap/zapcore"
)

// LoggingConfig represents a logging config.
type LoggingConfig struct {
	Level string `required:"true" default:"debug"`
	level zapcore.Level
}

// ClusterConfig represents a cluster membership config.
type ClusterConfig struct {
	Name      string
	Endpoint  string
	Announce  string
	SecretKey string `yaml:"secret_key" json:"-"`
	Members   []string
}

type config struct {
	Addr    string        `yaml:"endpoint" required:"true"`
	Cluster ClusterConfig `yaml:"cluster"`
	Logging LoggingConfig `yaml:"logging"`
}

// TODO: Docs.
type Config struct {
	Addr    net.Addr
	Cluster ClusterConfig
	Logging LoggingConfig
}

// NewConfig loads a new Relay server config from a file.
func NewConfig(path string) (*Config, error) {
	cfg := &config{}
	err := configor.Load(cfg, path)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveTCPAddr("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	lvl, err := logging.ParseLogLevel(cfg.Logging.Level)
	if err != nil {
		return nil, err
	}
	cfg.Logging.level = lvl

	if len(cfg.Cluster.Name) == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, err
		}

		cfg.Cluster.Name = fmt.Sprintf("%s-%s", hostname, uuid.New())
	}

	return &Config{
		Addr:    addr,
		Cluster: cfg.Cluster,
		Logging: cfg.Logging,
	}, nil
}

// LogLevel returns the minimum logging level configured.
func (c *Config) LogLevel() zapcore.Level {
	return c.Logging.level
}

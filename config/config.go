package config

import (
	"github.com/pkg/errors"
	"github.com/sencoder/go-loghouse/pkg/clickhouse"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var DefaultConfig = &Config{
	Http: &HttpConfig{
		Enabled: true,
		Listen:  ":8080",
	},
	Clickhouse: &clickhouse.Config{
		Address:  "127.0.0.1:9000",
		Username: "default",
	},
}

type Config struct {
	Http       *HttpConfig
	Clickhouse *clickhouse.Config
}

type HttpConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Listen  string `json:"listen" yaml:"listen"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "load config fail")
	}

	return LoadConfigBytes(bytes)
}

func LoadConfigBytes(configBytes []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, errors.Wrap(err, "parse config fail")
	}

	return &config, nil
}

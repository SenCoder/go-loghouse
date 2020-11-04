package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var DefaultConfig = &Config{
	Http: &HttpConfig{
		Enabled: true,
		Listen:  ":8080",
	},
}

type Config struct {
	Http *HttpConfig
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

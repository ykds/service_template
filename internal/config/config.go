package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"service_template/internal/server"
	"service_template/pkg/cache"
	"service_template/pkg/db"
	"service_template/pkg/logger"
)

type Config struct {
	Database   db.Option     `json:"database" yaml:"database"`
	Cache      cache.Option  `json:"cache" yaml:"cache"`
	HttpServer server.Option `json:"http_server" yaml:"http_server"`
	Log        logger.Option `json:"log" yaml:"log"`
}

func InitConfig(f string) (*Config, error) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return nil, err
	}
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(file, c)
	return c, err
}

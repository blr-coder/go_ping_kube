package config

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
)

type Config struct {
	AppPort string `yaml:"app_port"`
	Redis   *Redis `yaml:"redis"`
}

type Redis struct {
	RedisAddr string `yaml:"redis_addr"`
	RedisPass string `yaml:"redis_pass"`
	RedisDB   int    `yaml:"redis_db"`
}

func NewConfig(configPath string) (*Config, error) {
	var config Config

	configFile, err := os.OpenFile(configPath, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return nil, err
	}

	return &config, yaml.NewDecoder(configFile).Decode(&config)
}

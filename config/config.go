package config

import (
	errs "github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"postgres"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errs.WithMessage(err, "failed to read config file")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errs.WithMessage(err, "failed to unmarshal config")
	}

	return &config, nil
}

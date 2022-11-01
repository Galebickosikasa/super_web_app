package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"super_web_app/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type    string `yaml:"type" env-default:"port"`
		BindApi string `yaml:"bind_api" env-default:"0.0.0.0"`
		Port    string `yaml:"0.0.0.0" env-default:"3000"`
	} `yaml:"listen"`
	Database struct {
		Username string `yaml:"username" env-default:"developer"`
		Password string `env:"PASSWORD"`
		Host     string `yaml:"host" env-default:"physphile.ru"`
		Port     string `yaml:"port" env-default:"5432"`
		Database string `yaml:"database" env-default:"user"`
	} `yaml:"database"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}

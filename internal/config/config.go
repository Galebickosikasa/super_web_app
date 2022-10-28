package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"super_web_app/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type    string `yaml:"type"`
		BindApi string `yaml:"bind_api"`
		Port    string `yaml:"port"`
	} `yaml:"listen"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
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

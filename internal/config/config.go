package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Receiver struct {
	Name            string `yaml:"name"`
	NotificationURL string `yaml:"notification_url"`
}

type Config struct {
	Receivers []Receiver `yaml:"receivers"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

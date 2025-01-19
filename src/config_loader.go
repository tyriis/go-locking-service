package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	API APIConfig `yaml:"api"`
}

type APIConfig struct {
	Type string `yaml:"type"`
	Port string `yaml:"port"`
}

func NewConfig() (*Config, error) {

	// Read config file
	data, err := os.ReadFile("config/config.yaml")

	log.Debug().Msg(string(data))

	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}

package infrastructure

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tyriis/rest-go/internal/domain"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"path"`
		Prefix string `yaml:"keyPrefix"`
	} `yaml:"redis"`
	Api struct {
		Port string `yaml:"port"`
	} `yaml:"api"`
}

var (
	config *Config
	once   sync.Once
)

type YAMLConfigHandler struct {
	path      string
	validator domain.ConfigValidator
	logger    *Logger
}

func NewYAMLConfigHandler(path string, validator domain.ConfigValidator, logger *Logger) *YAMLConfigHandler {
	return &YAMLConfigHandler{
		path:      path,
		validator: validator,
		logger:    logger,
	}
}

func (h *YAMLConfigHandler) Load() (*Config, error) {
	// Retrieve the raw YAML data
	data, err := os.ReadFile(h.path)
	if err != nil {
		const msg = "YAMLConfigHandler.Load - os.ReadFile > %w"
		log.Error().Msg(fmt.Errorf(msg, err).Error())
		return nil, fmt.Errorf(msg, err)
	}

	// Unmarshal the YAML data to detect malformed YAML
	var rawData interface{}
	err = yaml.Unmarshal(data, &rawData)

	// Validate the raw data
	if err := h.validator.Validate(rawData); err != nil {
		const msg = "YAMLConfigHandler.Load - h.validator.Validate > %w"
		log.Error().Msg(fmt.Errorf(msg, err).Error())
		return nil, err
	}

	// Unmarshal the YAML data to the Config struct
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

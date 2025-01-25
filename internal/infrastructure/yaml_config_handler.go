package infrastructure

import (
	"fmt"
	"os"

	"github.com/tyriis/rest-go/internal/domain"
	"gopkg.in/yaml.v3"
)

type YAMLConfigHandler struct {
	path      string
	validator domain.ConfigValidator
	logger    domain.Logger
}

// NewYAMLConfigHandler creates a new YAMLConfigHandler with the given path, validator, and logger.
func NewYAMLConfigHandler(path string, validator domain.ConfigValidator, logger *Logger) *YAMLConfigHandler {
	return &YAMLConfigHandler{
		path:      path,
		validator: validator,
		logger:    logger,
	}
}

// Load reads the YAML file and returns the Config struct.
func (h *YAMLConfigHandler) Load() (*domain.Config, error) {
	// Retrieve the raw YAML data
	data, err := os.ReadFile(h.path)
	if err != nil {
		const msg = "YAMLConfigHandler.Load - os.ReadFile > %w"
		h.logger.Error(fmt.Errorf(msg, err).Error())
		return nil, fmt.Errorf(msg, err)
	}

	// Unmarshal the YAML data to detect malformed YAML
	var rawData interface{}
	err = yaml.Unmarshal(data, &rawData)
	if err != nil {
		const msg = "YAMLConfigHandler.Load - yaml.Unmarshal > %w"
		h.logger.Error(fmt.Errorf(msg, err).Error())
		return nil, err
	}

	// Validate the raw data
	if err := h.validator.Validate(rawData); err != nil {
		const msg = "YAMLConfigHandler.Load - h.validator.Validate > %w"
		h.logger.Error(fmt.Errorf(msg, err).Error())
		return nil, err
	}

	// Unmarshal the YAML data to the Config struct
	config := &domain.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		const msg = "YAMLConfigHandler.Load - yaml.Unmarshal > %w"
		h.logger.Error(fmt.Errorf(msg, err).Error())
		return nil, err
	}

	return config, nil
}

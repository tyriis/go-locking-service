package dao

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tyriis/rest-go/src/dto"
)

type ConfigDAO struct{}

func NewConfigDAO() *ConfigDAO {
	return &ConfigDAO{}
}

func (dao *ConfigDAO) RetrieveConfig() (*dto.Config, error) {
	// Read config file
	data, err := dao.RetrieveRawYAML()

	// Handle file not found
	if err != nil {
		return nil, err
	}

	// Parse YAML
	config := &dto.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, &dto.DAOError{
			Operation: "parse config",
			Err:       err,
		}
	}
	return config, nil
}

func (dao *ConfigDAO) RetrieveRawYAML() ([]byte, error) {
	// Read config file
	data, err := os.ReadFile("config/config.yaml")

	// Handle file not found
	if err != nil {
		return nil, &dto.DAOError{
			Operation: "load config",
			Err:       err,
		}
	}
	return data, nil
}

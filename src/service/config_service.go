package service

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/tyriis/rest-go/src/dao"
	"github.com/tyriis/rest-go/src/dto"
	"github.com/tyriis/rest-go/src/utils"
)

type ConfigService struct {
	dao *dao.ConfigDAO
}

func NewConfigService(dao *dao.ConfigDAO) *ConfigService {
	return &ConfigService{dao: dao}
}

func (service *ConfigService) GetConfig() (*dto.Config, error) {
	err := service.ValidateYAMLConfig()
	if err != nil {
		return nil, err
	}
	config, err := service.dao.RetrieveConfig()
	if err != nil {
		return nil, &dto.ServiceError{
			Service: utils.GetFunctionName(),
			Message: "failed to load configuration",
			Err:     err,
		}
	}
	return config, nil
}

func (service *ConfigService) ValidateYAMLConfig() error {
	// Retrieve the raw YAML data
	yamlData, err := service.dao.RetrieveRawYAML()
	if err != nil {
		return &dto.ServiceError{
			Service: utils.GetFunctionName(),
			Message: "Error reading YAML config file",
			Err:     err,
		}
	}

	// Unmarshal the YAML data to detect malformed YAML
	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return &dto.ServiceError{
			Service: utils.GetFunctionName(),
			Message: "Error parsing YAML config",
			Err:     err,
		}
	}

	// Compile the JSON schema
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile("config/schema.json")
	if err != nil {
		return &dto.ServiceError{
			Service: utils.GetFunctionName(),
			Message: "Error compiling jsonschema",
			Err:     err,
		}
	}

	// Validate the YAML data against the schema
	err = schema.Validate(data)
	if err != nil {
		if validationError, ok := err.(*jsonschema.ValidationError); ok {
			return &dto.ServiceError{
				Service: utils.GetFunctionName(),
				Message: "Validation error",
				Err:     validationError,
			}
		} else {
			return &dto.ServiceError{
				Service: utils.GetFunctionName(),
				Message: "Validation error",
				Err:     err,
			}
		}
	}
	return nil
}

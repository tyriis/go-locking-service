package infrastructure

import (
	"bytes"
	"embed"
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/tyriis/go-locking-service/internal/domain"
)

//go:embed assets/schemas/*
var assetsFS embed.FS

type JSONSchemaValidator struct {
	schemaPath string
	logger     domain.Logger
}

func NewJSONSchemaValidator(schemaPath string, logger domain.Logger) *JSONSchemaValidator {
	return &JSONSchemaValidator{
		schemaPath: schemaPath,
		logger:     logger,
	}
}

func (v *JSONSchemaValidator) Validate(data interface{}) error {
	// Convert data to JSON for validation
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Read schema
	schemaData, err := assetsFS.ReadFile(v.schemaPath)
	if err != nil {
		return err
	}

	// Unmarshal schema
	configSchema, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaData))
	if err != nil {
		return err
	}

	// Compile schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", configSchema); err != nil {
		return err
	}

	// Compile schema
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return err
	}

	// Validate
	var validationData interface{}
	if err := json.Unmarshal(jsonData, &validationData); err != nil {
		return err
	}

	return schema.Validate(validationData)
}

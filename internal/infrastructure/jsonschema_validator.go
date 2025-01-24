package infrastructure

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/tyriis/rest-go/internal/domain"
)

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

	// Load and compile schema
	schemaData, err := os.ReadFile(v.schemaPath)
	if err != nil {
		return err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(schemaData)); err != nil {
		return err
	}

	schema, err := compiler.Compile(v.schemaPath)
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

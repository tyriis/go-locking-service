package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

func validateConfig() {
	// Read and parse the YAML file
	yamlData, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading YAML config file")
		return
	}

	// Unmarshal the YAML data
	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing YAML config")
		return
	}

	// Compile the JSON schema
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile("config/schema.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Error compiling jsonschema")
		return
	}

	// Validate the YAML data against the schema
	err = schema.Validate(data)
	if err != nil {
		if validationError, ok := err.(*jsonschema.ValidationError); ok {
			log.Fatal().Msg(validationError.Error())
		} else {
			log.Fatal().Err(err).Msg("Validation error")
		}
		return
	}

	log.Info().Msg("Configuration is valid")
}

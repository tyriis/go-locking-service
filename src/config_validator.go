package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

func validateConfig() {
	// Read and parse the YAML file
	yamlData, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Error().Err(err).Msg("Error reading YAML file")
		return
	}

	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing YAML")
		return
	}

	// Compile the JSON schema
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile("config/schema.json")
	if err != nil {
		log.Error().Err(err).Msg("Error compiling schema")
		return
	}

	// Validate the YAML data against the schema
	err = schema.Validate(data)
	if err != nil {
		if validationError, ok := err.(*jsonschema.ValidationError); ok {
			log.Fatal().Msg(validationError.Error())
			// for _, outputUnit := range validationError.DetailedOutput().Errors {
			// 	handleErrors(outputUnit, validationError)
			// }
		} else {
			log.Fatal().Err(err).Msg("Validation error")
		}
		return
	}

	log.Info().Msg("Configuration is valid")
}

func handleErrors(outputUnit jsonschema.OutputUnit, validationErr *jsonschema.ValidationError) {
	if outputUnit.Errors != nil {
		for _, item := range outputUnit.Errors {
			if item.Errors != nil {
				handleErrors(item, validationErr)
			} else {
				handleError(item, item.Error, validationErr)
			}
		}
	} else {
		handleError(outputUnit, outputUnit.Error, validationErr)
	}
}

func handleError(detail jsonschema.OutputUnit, outputErr *jsonschema.OutputError, validationErr *jsonschema.ValidationError) {
	// if additionalProps, ok := outputErr.Kind.(*kind.AdditionalProperties); ok {
	// 	for _, property := range additionalProps.Properties {
	// 		path := formatYAMLKey(detail.InstanceLocation)
	// 		var propertyPath string
	// 		if path == "" {
	// 			propertyPath = property
	// 		} else {
	// 			propertyPath = fmt.Sprintf("%s.%s", path, property)
	// 		}
	// 		log.Fatal().
	// 			Msgf("invalid key: '%s'", propertyPath)
	// 	}
	// } else if typeErr, ok := outputErr.Kind.(*kind.Type); ok {
	// 	path := formatYAMLKey(detail.InstanceLocation)
	// 	log.Fatal().
	// 		Str("got", typeErr.Got).
	// 		Strs("want", typeErr.Want).
	// 		Msgf("invalid type for '%s'", path)
	// } else if requiredErr, ok := outputErr.Kind.(*kind.Required); ok {
	// 	path := formatYAMLKey(detail.InstanceLocation)
	// 	log.Fatal().
	// 		Msgf("missing key: '%s.%s'", path, requiredErr.Missing[0])
	// } else {
	log.Error().Msg("Config validation error")
	log.Fatal().Msg(validationErr.Error())
	// }
}

func formatYAMLKey(input string) string {
	return strings.ReplaceAll(strings.TrimPrefix(input, "/"), "/", ".")
}

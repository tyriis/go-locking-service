package main

import (
	"github.com/rs/zerolog/log"

	"github.com/tyriis/rest-go/src/api"
	"github.com/tyriis/rest-go/src/dao"
	"github.com/tyriis/rest-go/src/service"
)

func main() {
	initLogger()
	validateConfig()

	config, err := NewConfig()
	if err != nil {
		log.Info().Msg("failed to load configuration")
		log.Debug().Msg("hello world")
		log.Debug().Str("Value", "hello world").Send()
		log.Error().Err(err).Msg("failed to load configuration")
		log.Fatal().Err(err).Msg("failed to load configuration")
	}
	api.NewHelloApi(service.NewHelloService(dao.NewHelloDAO()))
	server(config)
}

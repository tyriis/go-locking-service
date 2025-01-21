package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/tyriis/rest-go/src/api"
	"github.com/tyriis/rest-go/src/dao"
	"github.com/tyriis/rest-go/src/service"
	"github.com/tyriis/rest-go/src/utils"
)

func main() {
	utils.InitLogger()
	// validateConfig()

	config, err := service.NewConfigService(dao.NewConfigDAO()).GetConfig()
	if err != nil {
		log.Error().
			// Stack().
			// Err(errors.WithStack(err)).
			Msg(err.Error())
		os.Exit(1)
		// log.Info().Msg("failed to load configuration")
		// log.Debug().Msg("hello world")
		// log.Debug().Str("Value", "hello world").Send()
		// log.Error().Err(err).Msg("failed to load configuration")
		// log.Fatal().Err(err).Msg("failed to load configuration")
	}
	api.NewHelloApi(service.NewHelloService(dao.NewHelloDAO()))
	server(config)
}

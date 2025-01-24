package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tyriis/rest-go/internal/delivery"
	"github.com/tyriis/rest-go/internal/infrastructure"
	"github.com/tyriis/rest-go/internal/repositories"
	"github.com/tyriis/rest-go/internal/usecases"
)

func main() {
	// initialize logger
	logger := infrastructure.NewLogger()

	// load config
	jsonValidator := infrastructure.NewJSONSchemaValidator("assets/schemas/config.json", logger)
	configHandler := infrastructure.NewYAMLConfigHandler("config/config.yaml", jsonValidator, logger)
	config, err := configHandler.Load()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	// initialize repository
	redisHandler := infrastructure.NewRedisHandler(*config, logger)
	lockRepo := repositories.NewLockRepository(redisHandler, logger)

	// initialize use case
	lockUseCase := usecases.NewLockUseCase(lockRepo, logger)

	// initialize http handler
	webserviceHandler := delivery.NewWebserviceHandler(lockUseCase, logger)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/locks", webserviceHandler.CreateLock).Methods("POST")
	r.HandleFunc("/api/v1/locks/{key}", webserviceHandler.DeleteLock).Methods("DELETE")
	r.HandleFunc("/api/v1/locks/{key}", webserviceHandler.ShowOneLock).Methods("GET")
	r.HandleFunc("/api/v1/locks", webserviceHandler.ShowAllLocks).Methods("GET")
	log.Printf("Server is running on port %s", config.Api.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Api.Port), r)
}

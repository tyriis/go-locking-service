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
	logger := infrastructure.NewLogger()

	validator := infrastructure.NewJSONSchemaValidator("assets/schemas/config.json")
	configHandler := infrastructure.NewYAMLConfigHandler("config/config.yaml", validator, logger)

	config, err := configHandler.Load()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	redisHandler := infrastructure.NewRedisHandler(config.Redis.Host + ":" + config.Redis.Port)
	lockRepo := repositories.NewLockRepository(redisHandler)

	// initialize use case
	lockUseCase := usecases.NewLockUseCase(lockRepo)

	// initialize http handler
	webserviceHandler := delivery.NewWebserviceHandler(lockUseCase, logger)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/locks/{name}", webserviceHandler.CreateLock).Methods("PUT")
	r.HandleFunc("/api/v1/locks/{name}", webserviceHandler.DeleteLock).Methods("DELETE")
	r.HandleFunc("/api/v1/locks/{name}", webserviceHandler.ShowOneLock).Methods("GET")
	r.HandleFunc("/api/v1/locks", webserviceHandler.ShowAllLocks).Methods("GET")
	log.Printf("Server is running on port %s", config.Api.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Api.Port), r)
}

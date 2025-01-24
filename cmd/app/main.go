package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Api.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

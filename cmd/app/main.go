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
	delivery "github.com/tyriis/go-locking-service/internal/delivery/http/service"
	"github.com/tyriis/go-locking-service/internal/infrastructure"
	"github.com/tyriis/go-locking-service/internal/metrics"
	"github.com/tyriis/go-locking-service/internal/repositories"
	"github.com/tyriis/go-locking-service/internal/usecases"
)

func main() {
	// initialize logger
	logger := infrastructure.NewLogger()

	// load config
	jsonValidator := infrastructure.NewJSONSchemaValidator("assets/schemas/config.json", logger)
	configHandler := infrastructure.NewYAMLConfigHandler("config/config.yaml", jsonValidator, logger)
	config, err := configHandler.Load()
	if err != nil {
		log.Fatal("App.main - Failed to load config")
	}

	// initialize repository
	redisHandler := infrastructure.NewRedisHandler(*config, logger)
	lockRepo := repositories.NewLockRepository(redisHandler, logger)

	// initialize use case
	lockUseCase := usecases.NewLockUseCase(lockRepo, logger)

	// initialize http handler
	webserviceHandler := delivery.NewWebserviceHandler(lockUseCase, logger)

	// initialize metrics service and middleware
	metricsService := metrics.NewPrometheusMetricsService()
	metricsMiddleware := metrics.NewMetricsMiddleware(metricsService)

	// Initialize and start metrics updater
	metricsUpdater := metrics.NewMetricsUpdater(lockRepo, metricsService, logger)
	metricsUpdater.Start()

	r := mux.NewRouter()

	// Apply metrics middleware to all routes
	r.Handle("/metrics", delivery.MetricsHandler())
	r.Handle("/api/v1/locks", metricsMiddleware.Middleware(http.HandlerFunc(webserviceHandler.CreateLock))).Methods("POST")
	r.Handle("/api/v1/locks/{key}", metricsMiddleware.Middleware(http.HandlerFunc(webserviceHandler.DeleteLock))).Methods("DELETE")
	r.Handle("/api/v1/locks/{key}", metricsMiddleware.Middleware(http.HandlerFunc(webserviceHandler.ShowOneLock))).Methods("GET")
	r.Handle("/api/v1/locks", metricsMiddleware.Middleware(http.HandlerFunc(webserviceHandler.ShowAllLocks))).Methods("GET")

	logger.Info("App.main - Server is running on http://" + config.Api.Host + ":" + config.Api.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Api.Host, config.Api.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("App.main - listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Stop metrics updater before shutting down
	metricsUpdater.Stop()

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("App.main - Server forced to shutdown:", err)
	}
}

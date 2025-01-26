package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	handler "main/internal/api/handlers"
	"main/internal/api/routes"
	cnf "main/internal/config"
	"main/internal/domain/usecases"
	pg "main/internal/repository/postgres"
	"main/internal/services"
	"main/pkg/auth"
	"main/pkg/database/migration"
	"main/pkg/database/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var postgresConfig postgres.Config
var config cnf.Config

func init() {
	config = *cnf.NewConfig()
}

func init() {
	postgresConfig = postgres.Config{
		URL:      config.Database.URL,
		User:     config.Database.User,
		Password: config.Database.Password,
		PoolSize: config.Database.PoolSize,
	}
}

func main() {

	dbPool, err := postgres.ConnectDB(postgresConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer postgres.CloseDB(dbPool)

	if err := migration.Migrate(dbPool, config.Migration.Directory); err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	port := config.Server.Port
	if port == 0 {
		log.Fatal("Server port is not set in the config file")
	}

	requestStorage := pg.NewRequestRepository(dbPool)
	requestService := services.NewRequestService(requestStorage)
	requestUC := usecases.NewRequestUC(requestService)
	requestHandler := handler.NewRequestHandler(requestUC)

	authS := auth.NewAuth(config.JwtAuth.Key)

	router := routes.NewRouter(
		authS,
		requestHandler,
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		log.Printf("Server is listening on port %d\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Fatalf("Server gracefully stopped")
}

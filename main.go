package main

import (
	"0xFarms-backend/config"
	"0xFarms-backend/internal/adapters"
	"0xFarms-backend/internal/core/services"
	"0xFarms-backend/internal/web"
	"0xFarms-backend/internal/web/handlers"
	"0xFarms-backend/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize the logger
	logger.InitLogger()

	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, _ := adapters.NewMongoAdapter(cfg.MONGO_URL)
	blogService := services.NewBlogService(db)
	farmService := services.NewFarmManagementSystemService(db)

	blogHandler := handlers.NewBlogHandler(blogService)
	farmHandler := handlers.NewFarmHandler(farmService)
	router := gin.Default()
	web.SetupAPIRoutes(router, blogHandler, farmHandler)

	// Define the server port
	PORT := fmt.Sprintf(":%s", cfg.PORT)
	gin.SetMode(gin.ReleaseMode)
	gracefulShutdown(router, PORT)

}

func gracefulShutdown(router *gin.Engine, port string) {
	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Create a server instance with a timeout
	srv := &http.Server{
		Addr:    port,
		Handler: router,
		// Optional: configure timeouts as needed
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Start the server in a goroutine
	go func() {

		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	log.Println("Shutting down server...")

	// Create a timeout context for shutting down the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
	log.Println("Server exiting")
}

package main

import (
	"log"

	"cidr-viewer/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cache-Control", "Pragma", "Expires"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		api.POST("/analyze", handlers.AnalyzeCIDRs)
		api.POST("/validate", handlers.ValidateCIDR)
		api.GET("/health", handlers.HealthCheck)
	}

	// Start server
	log.Println("Starting CIDR Viewer API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}